package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
)

var templateFuncs = template.FuncMap{"lines": lines}

func mustParse(path string) *template.Template {
	return template.Must(template.New(filepath.Base(path)).Funcs(templateFuncs).ParseFiles(path))
}

var (
	tmplIndex    = mustParse("frontend/templates/index.html")
	tmplAbout    = mustParse("frontend/templates/about.html")
	tmplServices = mustParse("frontend/templates/services.html")
	tmplBlog     = mustParse("frontend/templates/blog.html")
	tmplLogin    = mustParse("frontend/templates/admin_login.html")
	tmplSetup    = mustParse("frontend/templates/admin_setup.html")
	tmplPageEdit = mustParse("frontend/templates/admin_page_edit.html")
	tmplPosts    = mustParse("frontend/templates/admin_posts.html")
	tmplPostEdit = mustParse("frontend/templates/admin_post_edit.html")
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
	mux.HandleFunc("/", app.handlePublicPage(tmplIndex, "home", homeDefaults))
	mux.HandleFunc("/about", app.handlePublicPage(tmplAbout, "about", aboutDefaults))
	mux.HandleFunc("/services", app.handlePublicPage(tmplServices, "services", servicesDefaults))
	mux.HandleFunc("/blog", app.handleBlog)
	mux.HandleFunc("/contact", app.handleContact)

	mux.HandleFunc("/admin/login", app.handleLoginPage)
	mux.HandleFunc("/admin/setup", app.handleSetup)
	mux.HandleFunc("/admin/logout", app.handleLogout)
	mux.HandleFunc("/admin", app.requireAuth(app.handlePageEditor("home", "Texty úvodní stránky", homeFields, homeDefaults)))
	mux.HandleFunc("/admin/about", app.requireAuth(app.handlePageEditor("about", "Texty stránky O mně", aboutFields, aboutDefaults)))
	mux.HandleFunc("/admin/services", app.requireAuth(app.handlePageEditor("services", "Texty stránky Služby", servicesFields, servicesDefaults)))
	mux.HandleFunc("/admin/posts", app.requireAuth(app.handlePosts))
	mux.HandleFunc("/admin/posts/new", app.requireAuth(app.handlePostForm))
	mux.HandleFunc("/admin/posts/edit", app.requireAuth(app.handlePostForm))
	mux.HandleFunc("/admin/posts/delete", app.requireAuth(app.handlePostDelete))

	// Admin panel also reachable at its own subdomain root (e.g.
	// admin.streckerova.kralroman.org/) instead of the /admin path above.
	// The /admin/* routes stay registered as a fallback until DNS/domain
	// mapping for ADMIN_HOST is live.
	if app.adminHost != "" {
		mux.HandleFunc(app.adminHost+"/login", app.handleLoginPage)
		mux.HandleFunc(app.adminHost+"/setup", app.handleSetup)
		mux.HandleFunc(app.adminHost+"/logout", app.handleLogout)
		mux.HandleFunc(app.adminHost+"/", app.requireAuth(app.handlePageEditor("home", "Texty úvodní stránky", homeFields, homeDefaults)))
		mux.HandleFunc(app.adminHost+"/about", app.requireAuth(app.handlePageEditor("about", "Texty stránky O mně", aboutFields, aboutDefaults)))
		mux.HandleFunc(app.adminHost+"/services", app.requireAuth(app.handlePageEditor("services", "Texty stránky Služby", servicesFields, servicesDefaults)))
		mux.HandleFunc(app.adminHost+"/posts", app.requireAuth(app.handlePosts))
		mux.HandleFunc(app.adminHost+"/posts/new", app.requireAuth(app.handlePostForm))
		mux.HandleFunc(app.adminHost+"/posts/edit", app.requireAuth(app.handlePostForm))
		mux.HandleFunc(app.adminHost+"/posts/delete", app.requireAuth(app.handlePostDelete))
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
