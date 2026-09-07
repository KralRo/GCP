package main

import (
	"context"
	"html/template"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fieldSpec describes one editable text field on a page: where it lives in
// Firestore, what the admin form should show for it, and how the public
// template should treat it (a single line, a paragraph, or a newline-per-item
// bullet list).
type fieldSpec struct {
	Key     string
	Label   string
	Kind    string // "line" | "text" | "list"
	Section string // groups fields under a subheading in the admin form
}

type fieldView struct {
	Key   string
	Label string
	Kind  string
	Value string
}

type sectionView struct {
	Heading string
	Fields  []fieldView
}

// loadPageContent merges whatever's stored in Firestore over the page's
// defaults, so a field nobody has edited yet still renders the original
// copy instead of coming up blank.
func (a *app) loadPageContent(ctx context.Context, page string, defaults map[string]string) (map[string]string, error) {
	content := make(map[string]string, len(defaults))
	for k, v := range defaults {
		content[k] = v
	}
	doc, err := a.fs.Collection("content").Doc(page).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return content, nil
	}
	if err != nil {
		return nil, err
	}
	var stored map[string]string
	if err := doc.DataTo(&stored); err != nil {
		return nil, err
	}
	for k, v := range stored {
		if v != "" {
			content[k] = v
		}
	}
	return content, nil
}

func buildSections(spec []fieldSpec, content map[string]string) []sectionView {
	var sections []sectionView
	for _, f := range spec {
		if len(sections) == 0 || sections[len(sections)-1].Heading != f.Section {
			sections = append(sections, sectionView{Heading: f.Section})
		}
		i := len(sections) - 1
		sections[i].Fields = append(sections[i].Fields, fieldView{f.Key, f.Label, f.Kind, content[f.Key]})
	}
	return sections
}

// lines splits a textarea's raw value into non-empty, trimmed lines - used
// by public templates to render "list" fields as bullet points.
func lines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// handlePublicPage renders a public page from its Firestore-backed content
// map, falling back to defaults for anything not yet edited.
func (a *app) handlePublicPage(t *template.Template, page string, defaults map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := a.loadPageContent(r.Context(), page, defaults)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, struct{ Content map[string]string }{content}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// handlePageEditor serves the admin editor for one page's content map -
// GET shows the current values grouped by section, POST replaces the whole
// Firestore document with the submitted form.
func (a *app) handlePageEditor(page, title string, spec []fieldSpec, defaults map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			content, err := a.loadPageContent(r.Context(), page, defaults)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := tmplPageEdit.Execute(w, struct {
				Title        string
				Sections     []sectionView
				Saved        bool
				HomePath     string
				AboutPath    string
				ServicesPath string
				PostsPath    string
				LogoutPath   string
			}{
				title, buildSections(spec, content), r.URL.Query().Get("saved") != "",
				a.adminPath(r), a.aboutEditPath(r), a.servicesEditPath(r), a.postsPath(r), a.logoutPath(r),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case http.MethodPost:
			data := make(map[string]string, len(spec))
			for _, f := range spec {
				data[f.Key] = r.FormValue(f.Key)
			}
			if _, err := a.fs.Collection("content").Doc(page).Set(r.Context(), data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			editPath := a.adminPath(r)
			if page != "home" {
				editPath = a.path(r, "/"+page)
			}
			http.Redirect(w, r, editPath+"?saved=1", http.StatusSeeOther)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
