package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/had-nu/lazy.go/pkg/config"
)

// DirEntry describes a single file or directory to create.
type DirEntry struct {
	Path     string // relative to project root
	IsDir    bool
	Template string // template name, empty for directories
	Data     any    // template data
}

// Generator orchestrates the full project scaffolding.
type Generator struct {
	cfg    *config.ProjectConfig
	outDir string
}

// New creates a new Generator.
func New(cfg *config.ProjectConfig, outDir string) *Generator {
	return &Generator{cfg: cfg, outDir: outDir}
}

// Generate runs the full scaffold pipeline.
func (g *Generator) Generate() error {
	if err := os.MkdirAll(g.outDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	entries := BuildDirectoryTree(g.cfg)

	for _, e := range entries {
		fullPath := filepath.Join(g.outDir, e.Path)

		if e.IsDir {
			if err := os.MkdirAll(fullPath, 0o755); err != nil {
				return fmt.Errorf("creating directory %s: %w", e.Path, err)
			}
			continue
		}

		// Ensure parent directory exists.
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("creating parent dir for %s: %w", e.Path, err)
		}

		if e.Template == "" {
			// Create empty file.
			if err := os.WriteFile(fullPath, []byte{}, 0o644); err != nil {
				return fmt.Errorf("creating file %s: %w", e.Path, err)
			}
			continue
		}

		content, err := RenderTemplate(e.Template, e.Data)
		if err != nil {
			return fmt.Errorf("rendering template %s: %w", e.Template, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("writing file %s: %w", e.Path, err)
		}
	}

	return nil
}

// templateData is the common data passed to all templates.
type templateData struct {
	Config      *config.ProjectConfig
	Year        int
	LibName     string
	ServiceName string
}

// commonData builds the shared template data struct.
func commonData(cfg *config.ProjectConfig) templateData {
	libName := strings.ToLower(strings.ReplaceAll(cfg.Name, "-", ""))
	return templateData{
		Config:      cfg,
		Year:        2026,
		LibName:     libName,
		ServiceName: cfg.Name,
	}
}

// RenderAll renders all templates for a config and returns path→content map.
func RenderAll(cfg *config.ProjectConfig) (map[string]string, error) {
	entries := BuildDirectoryTree(cfg)
	result := make(map[string]string)

	for _, e := range entries {
		if e.IsDir || e.Template == "" {
			continue
		}
		content, err := RenderTemplate(e.Template, e.Data)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", e.Template, err)
		}
		result[e.Path] = content
	}

	return result, nil
}

// funcMap provides helper functions available in all templates.
var funcMap = template.FuncMap{
	"upper":   strings.ToUpper,
	"lower":   strings.ToLower,
	"title":   strings.Title, //nolint:staticcheck
	"replace": strings.ReplaceAll,
	"join":    strings.Join,
}
