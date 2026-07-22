package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

var (
	tmplIndex = template.Must(template.ParseFiles("frontend/templates/index.html"))
	tmplLogin = template.Must(template.ParseFiles("frontend/templates/admin_login.html"))
	tmplSetup = template.Must(template.ParseFiles("frontend/templates/admin_setup.html"))
	tmplEdit  = template.Must(template.ParseFiles("frontend/templates/admin_edit.html"))
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	databaseID := os.Getenv("FIRESTORE_DATABASE_ID")
	if databaseID == "" {
		databaseID = "(default)"
	}
	fsClient, err := firestore.NewClientWithDatabase(ctx, os.Getenv("FIRESTORE_PROJECT_ID"), databaseID)
	if err != nil {
		log.Fatalf("firestore client: %v", err)
	}
	defer fsClient.Close()

	app := &app{
		fs:            fsClient,
		sessionSecret: []byte(os.Getenv("SESSION_SECRET")),
		adminHost:     os.Getenv("ADMIN_HOST"),
	}

	mux := http.NewServeMux()
	// NOT /healthz - that path is intercepted by Google's frontend on *.run.app
	// domains and never reaches the container (returns a Google-branded 404).
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	mux.HandleFunc("/", app.handleIndex)

	mux.HandleFunc("/admin/login", app.handleLoginPage)
	mux.HandleFunc("/admin/setup", app.handleSetup)
	mux.HandleFunc("/admin/logout", app.handleLogout)
	mux.HandleFunc("/admin", app.requireAuth(app.handleAdmin))

	// Admin panel also reachable at its own subdomain root (e.g.
	// admin.example.com/) instead of the /admin path above. The /admin/*
	// routes stay registered as a fallback until DNS/domain mapping for
	// ADMIN_HOST is live.
	if app.adminHost != "" {
		mux.HandleFunc(app.adminHost+"/login", app.handleLoginPage)
		mux.HandleFunc(app.adminHost+"/setup", app.handleSetup)
		mux.HandleFunc(app.adminHost+"/logout", app.handleLogout)
		mux.HandleFunc(app.adminHost+"/", app.requireAuth(app.handleAdmin))
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
