package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type app struct {
	fs            *firestore.Client
	sessionSecret []byte
	// Hostname the admin panel is also served on at its root path, e.g.
	// "admin.streckerova.kralroman.org". Empty disables host-based routing.
	adminHost string
}

type adminConfig struct {
	Username     string `firestore:"username"`
	PasswordHash string `firestore:"password_hash"`
}

type pageContent struct {
	Body string `firestore:"body"`
}

type blogPost struct {
	Title     string    `firestore:"title"`
	Date      string    `firestore:"date"`
	Body      string    `firestore:"body"`
	CreatedAt time.Time `firestore:"created_at"`
}

const sessionCookieName = "admin_session"
const sessionTTL = 12 * time.Hour

// Seed text shown on the homepage hero until the client edits it via /admin.
const defaultHomeBody = "Nastavíme společně bezpečnost tak, aby fungovala v každodenním provozu, splňovala zákonné požadavky a přirozeně zapadla do toho, jak vaše firma opravdu funguje. Zaměříme se na praktické kroky vycházející z NIS2 a AI Actu, které zapadnou do vašeho běžného provozu a budou srozumitelné a použitelné pro každého."

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
		return defaultHomeBody, nil
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

// path returns the right URL for an admin sub-page depending on which host
// the request came in on - the main site's /admin/<suffix> path, or the
// dedicated adminHost's /<suffix> root, when that's configured and matches
// the request. suffix must be empty or start with "/".
func (a *app) path(r *http.Request, suffix string) string {
	if a.adminHost != "" && r.Host == a.adminHost {
		if suffix == "" {
			return "/"
		}
		return suffix
	}
	return "/admin" + suffix
}

func (a *app) loginPath(r *http.Request) string  { return a.path(r, "/login") }
func (a *app) setupPath(r *http.Request) string  { return a.path(r, "/setup") }
func (a *app) logoutPath(r *http.Request) string { return a.path(r, "/logout") }
func (a *app) adminPath(r *http.Request) string  { return a.path(r, "") }
func (a *app) postsPath(r *http.Request) string  { return a.path(r, "/posts") }
func (a *app) postsNewPath(r *http.Request) string {
	return a.path(r, "/posts/new")
}
func (a *app) postEditPath(r *http.Request, id string) string {
	return a.path(r, "/posts/edit?id="+url.QueryEscape(id))
}
func (a *app) postDeletePath(r *http.Request, id string) string {
	return a.path(r, "/posts/delete?id="+url.QueryEscape(id))
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
		Path:     "/",
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
		Path:     "/",
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
			http.Redirect(w, r, a.loginPath(r), http.StatusSeeOther)
			return
		}
		if _, ok := a.verifySession(cookie.Value); !ok {
			http.Redirect(w, r, a.loginPath(r), http.StatusSeeOther)
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

// handleStatic renders a template that needs no request-specific data - the
// about/services pages, neither of which are Firestore-backed yet.
func (a *app) handleStatic(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// handleContact just acknowledges the submission - no SMTP wiring yet, see
// the "sent=unavailable" banner on the homepage. Never send credentials or
// mail here without going through Secret Manager first.
func (a *app) handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/#contact", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?sent=unavailable#contact", http.StatusSeeOther)
}

func (a *app) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	cfg, err := a.loadAdminConfig(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cfg == nil {
		if err := tmplSetup.Execute(w, struct{ SetupPath string }{a.setupPath(r)}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := tmplLogin.Execute(w, struct {
			Error     bool
			LoginPath string
		}{r.URL.Query().Get("error") != "", a.loginPath(r)}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username != cfg.Username || bcrypt.CompareHashAndPassword([]byte(cfg.PasswordHash), []byte(password)) != nil {
			http.Redirect(w, r, a.loginPath(r)+"?error=1", http.StatusSeeOther)
			return
		}
		a.setSessionCookie(w, username)
		http.Redirect(w, r, a.adminPath(r), http.StatusSeeOther)
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
		http.Redirect(w, r, a.loginPath(r), http.StatusSeeOther)
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
	http.Redirect(w, r, a.adminPath(r), http.StatusSeeOther)
}

func (a *app) handleLogout(w http.ResponseWriter, r *http.Request) {
	a.clearSessionCookie(w)
	http.Redirect(w, r, a.loginPath(r), http.StatusSeeOther)
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
			Body       string
			Saved      bool
			PostsPath  string
			LogoutPath string
		}{body, r.URL.Query().Get("saved") != "", a.postsPath(r), a.logoutPath(r)}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		_, err := a.fs.Collection("content").Doc("home").Set(r.Context(), pageContent{Body: r.FormValue("body")})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, a.adminPath(r)+"?saved=1", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// loadPosts returns blog posts newest-first. Firestore query results don't
// include the document ID, so it's copied onto ID after DataTo.
func (a *app) loadPosts(ctx context.Context) ([]blogPost, []string, error) {
	iter := a.fs.Collection("posts").OrderBy("created_at", firestore.Desc).Documents(ctx)
	defer iter.Stop()
	var posts []blogPost
	var ids []string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		var p blogPost
		if err := doc.DataTo(&p); err != nil {
			return nil, nil, err
		}
		posts = append(posts, p)
		ids = append(ids, doc.Ref.ID)
	}
	return posts, ids, nil
}

func (a *app) handleBlog(w http.ResponseWriter, r *http.Request) {
	posts, _, err := a.loadPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmplBlog.Execute(w, struct{ Posts []blogPost }{posts}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type postRow struct {
	blogPost
	EditPath   string
	DeletePath string
}

func (a *app) handlePosts(w http.ResponseWriter, r *http.Request) {
	posts, ids, err := a.loadPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows := make([]postRow, len(posts))
	for i, p := range posts {
		rows[i] = postRow{p, a.postEditPath(r, ids[i]), a.postDeletePath(r, ids[i])}
	}
	if err := tmplPosts.Execute(w, struct {
		Posts      []postRow
		NewPath    string
		AdminPath  string
		LogoutPath string
	}{rows, a.postsNewPath(r), a.adminPath(r), a.logoutPath(r)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostForm serves both /posts/new (no id) and /posts/edit?id=... - the
// same form self-submits back to whichever of the two it was loaded from.
func (a *app) handlePostForm(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	switch r.Method {
	case http.MethodGet:
		post := blogPost{}
		if id != "" {
			doc, err := a.fs.Collection("posts").Doc(id).Get(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := doc.DataTo(&post); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if err := tmplPostEdit.Execute(w, struct {
			blogPost
			IsNew     bool
			PostsPath string
		}{post, id == "", a.postsPath(r)}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		post := blogPost{
			Title: r.FormValue("title"),
			Date:  r.FormValue("date"),
			Body:  r.FormValue("body"),
		}
		if id != "" {
			// Preserve the original CreatedAt so edits don't reorder the list.
			existing, err := a.fs.Collection("posts").Doc(id).Get(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var current blogPost
			if err := existing.DataTo(&current); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post.CreatedAt = current.CreatedAt
			_, err = a.fs.Collection("posts").Doc(id).Set(r.Context(), post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			post.CreatedAt = time.Now().UTC()
			_, _, err := a.fs.Collection("posts").Add(r.Context(), post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, a.postsPath(r), http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *app) handlePostDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if _, err := a.fs.Collection("posts").Doc(id).Delete(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.postsPath(r), http.StatusSeeOther)
}
