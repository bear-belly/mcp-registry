package templates

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/bear-belly/mcp-registry/internal/models"
)

var (
	// Templates holds all parsed templates
	Templates *template.Template
)

// InitTemplates initializes and parses all templates
func InitTemplates(config models.Config) error {
	var err error

	// First, parse the layout template
	Templates, err = template.New("layout.html").ParseFiles(filepath.Join(config.TemplatePath, "layout.html"))
	if err != nil {
		return fmt.Errorf("parsing layout template: %w", err)
	}

	// Then parse all other templates
	pattern := filepath.Join(config.TemplatePath, "*.html")
	Templates, err = Templates.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("parsing other templates: %w", err)
	}

	return nil
}

// ExecuteTemplate executes a template with tracing
func ExecuteTemplate(ctx context.Context, w http.ResponseWriter, name string, data interface{}) error {
	err := Templates.ExecuteTemplate(w, name, data)
	if err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}

	return nil
}

// PageData holds common data for all pages
type PageData struct {
	Title        string
	PageTemplate string      // Specifies which template content to render
	Data         interface{} // For passing page-specific data
	ConfigJSON   string      // JSON string representation of config data
}
