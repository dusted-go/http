package htmlview

import (
	"fmt"
	"html/template"
	"net/http"
)

type Writer struct {
	hotReload     bool
	layoutName    string
	templateFiles map[string][]string
	templates     map[string]*template.Template
}

func (hw *Writer) WriteView(
	w http.ResponseWriter,
	statusCode int,
	key string,
	model interface{},
) error {
	var t *template.Template

	// In production settings use pre-created templates,
	// otherwise create a new template every time during
	// for a faster feedback loop during development:
	if hw.hotReload {
		t = createTemplate(hw.templateFiles[key]...)
	} else {
		t = hw.templates[key]
	}

	w.WriteHeader(statusCode)
	err := t.ExecuteTemplate(w, hw.layoutName, model)

	if err != nil {
		return fmt.Errorf("error executing template with key '%s': %w", key, err)
	}
	return nil
}

func NewWriter(
	hotReload bool,
	layoutName string,
	templateFiles map[string][]string,
) *Writer {

	templates := make(map[string]*template.Template)
	for key, files := range templateFiles {
		templates[key] = createTemplate(files...)
	}

	return &Writer{
		hotReload:     hotReload,
		layoutName:    layoutName,
		templateFiles: templateFiles,
		templates:     templates,
	}
}

func createTemplate(files ...string) *template.Template {
	return template.Must(template.ParseFiles(files...))
}
