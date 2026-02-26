package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// templateCache holds compiled templates.
var templateCache *template.Template

func init() {
	var err error
	templateCache, err = template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		panic(fmt.Sprintf("lazy.go: failed to parse embedded templates: %v", err))
	}
}

// RenderTemplate renders a named template with the given data.
func RenderTemplate(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := templateCache.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("executing template %q: %w", name, err)
	}
	return buf.String(), nil
}
