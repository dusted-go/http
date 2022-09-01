package server

import (
	"html/template"
	"net/http"

	"github.com/dusted-go/fault/fault"
)

type ViewHandler struct {
	hotReload     bool
	layoutName    string
	templateFiles map[string][]string
	templates     map[string]*template.Template
}

func (h *ViewHandler) RenderView(
	statusCode int,
	key string,
	model interface{},
	w http.ResponseWriter,
	r *http.Request,
) error {
	var t *template.Template

	// In production settings use pre-created templates,
	// otherwise create a new template every time during
	// for a faster feedback loop during development:
	if h.hotReload {
		t = createTemplate(h.templateFiles[key]...)
	} else {
		t = h.templates[key]
	}

	w.WriteHeader(statusCode)
	err := t.ExecuteTemplate(w, h.layoutName, model)
	if err != nil {
		return fault.SystemWrapf(
			err, "server", "RenderView",
			"failed to execute template with key: %s", key)
	}
	return nil
}

func NewViewHandler(
	hotReload bool,
	layoutName string,
	templateFiles map[string][]string,
) *ViewHandler {

	templates := make(map[string]*template.Template)
	for key, files := range templateFiles {
		templates[key] = createTemplate(files...)
	}

	return &ViewHandler{
		hotReload:     hotReload,
		layoutName:    layoutName,
		templateFiles: templateFiles,
		templates:     templates,
	}
}

func createTemplate(files ...string) *template.Template {
	return template.Must(template.ParseFiles(files...))
}
