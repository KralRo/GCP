package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type app struct {
	fs            *firestore.Client
	sessionSecret []byte
}

type adminConfig struct {
	Username     string `firestore:"username"`
	PasswordHash string `firestore:"password_hash"`
}

type pageContent struct {
	Body string `firestore:"body"`
}

const sessionCookieName = "admin_session"
const sessionTTL = 12 * time.Hour

func (a *app) loadAdminConfig(ctx context.Context) (*adminConfig, error) {
	doc, err := a.fs.Collection("admin").Doc("config").Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg adminConfig
	if err := doc.DataTo(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (a *app) loadContent(ctx context.Context) (string, error) {
	doc, err := a.fs.Collection("content").Doc("home").Get(ctx)
	if status.Code(err) == codes.NotFound {
		return "It works.", nil
	}
	if err != nil {
		return "", err
	}
	var c pageContent
	if err := doc.DataTo(&c); err != nil {
		return "", err
	}
	return c.Body, nil
}

func (a *app) signSession(username string, expiry int64) string {
	payload := fmt.Sprintf("%s.%d", username, expiry)
	mac := hmac.New(sha256.New, a.sessionSecret)
	mac.Write([]byte(payload))
	return payload + "." + hex.EncodeToString(mac.Sum(nil))
}

func (a *app) verifySession(value string) (string, bool) {
	parts := strings.SplitN(value, ".", 3)
	if len(parts) != 3 {
		return "", false
	}
	username, expiryStr, sig := parts[0], parts[1], parts[2]
	expected := hmac.New(sha256.New, a.sessionSecret)
	expected.Write([]byte(username + "." + expiryStr))
	if !hmac.Equal([]byte(sig), []byte(hex.EncodeToString(expected.Sum(nil)))) {
		return "", false
	}
	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil || time.Now().Unix() > expiry {
		return "", false
	}
	return username, true
}

func (a *app) setSessionCookie(w http.ResponseWriter, username string) {
	expiry := time.Now().Add(sessionTTL).Unix()
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    a.signSession(username, expiry),
		Path:     "/admin",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(expiry, 0),
	})
}

func (a *app) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/admin",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (a *app) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		if _, ok := a.verifySession(cookie.Value); !ok {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (a *app) handleIndex(w http.ResponseWriter, r *http.Request) {
	body, err := a.loadContent(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmplIndex.Execute(w, struct{ Body string }{body}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *app) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	cfg, err := a.loadAdminConfig(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cfg == nil {
		if err := tmplSetup.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := tmplLogin.Execute(w, struct{ Error bool }{r.URL.Query().Get("error") != ""}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username != cfg.Username || bcrypt.CompareHashAndPassword([]byte(cfg.PasswordHash), []byte(password)) != nil {
			http.Redirect(w, r, "/admin/login?error=1", http.StatusSeeOther)
			return
		}
		a.setSessionCookie(w, username)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *app) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cfg, err := a.loadAdminConfig(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cfg != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = a.fs.Collection("admin").Doc("config").Set(r.Context(), adminConfig{
		Username:     username,
		PasswordHash: string(hash),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.setSessionCookie(w, username)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (a *app) handleLogout(w http.ResponseWriter, r *http.Request) {
	a.clearSessionCookie(w)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func (a *app) handleAdmin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		body, err := a.loadContent(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmplEdit.Execute(w, struct {
			Body  string
			Saved bool
		}{body, r.URL.Query().Get("saved") != ""}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		_, err := a.fs.Collection("content").Doc("home").Set(r.Context(), pageContent{Body: r.FormValue("body")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin?saved=1", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
