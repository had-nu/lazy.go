package scaffold_test

import (
	"os"
	"testing"

	"github.com/hadnu/lazy.go/pkg/config"
	"github.com/hadnu/lazy.go/pkg/scaffold"
)

type tmplData struct {
	Config      *config.ProjectConfig
	Year        int
	LibName     string
	ServiceName string
}

func newTmplData(cfg *config.ProjectConfig) tmplData {
	return tmplData{Config: cfg, Year: 2026, LibName: "testapp", ServiceName: cfg.Name}
}

func apicfg() *config.ProjectConfig {
	return &config.ProjectConfig{
		Name:        "testapp",
		ModulePath:  "github.com/user/testapp",
		Description: "A test application",
		Author:      "Test Author",
		Type:        config.ProjectTypeAPI,
		Visibility:  config.VisibilityPublic,
		License:     config.LicenseMIT,
		Criticality: config.CriticalityProduction,
		Features: config.Features{
			Docker:        true,
			GitHubActions: true,
			SAST:          true,
		},
	}
}

func TestRenderTemplate_Readme(t *testing.T) {
	out, err := scaffold.RenderTemplate("readme.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("readme.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("readme.tmpl: empty output")
	}
}

func TestRenderTemplate_Gomod(t *testing.T) {
	out, err := scaffold.RenderTemplate("gomod.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("gomod.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("gomod.tmpl: empty output")
	}
}

func TestRenderTemplate_Workflow(t *testing.T) {
	out, err := scaffold.RenderTemplate("workflow.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("workflow.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("workflow.tmpl: empty output")
	}
}

func TestRenderTemplate_Golangci(t *testing.T) {
	out, err := scaffold.RenderTemplate("golangci.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("golangci.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("golangci.tmpl: empty output")
	}
}

func TestRenderTemplate_Contributing(t *testing.T) {
	out, err := scaffold.RenderTemplate("contributing.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("contributing.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("contributing.tmpl: empty output")
	}
}

func TestRenderTemplate_Dockerfile(t *testing.T) {
	out, err := scaffold.RenderTemplate("dockerfile.tmpl", newTmplData(apicfg()))
	if err != nil {
		t.Fatalf("dockerfile.tmpl failed: %v", err)
	}
	if len(out) == 0 {
		t.Error("dockerfile.tmpl: empty output")
	}
}

func TestGenerate_CreatesFiles(t *testing.T) {
	dir := t.TempDir()
	outDir := dir + "/testapp"

	gen := scaffold.New(apicfg(), outDir)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	for _, path := range []string{
		outDir + "/README.md",
		outDir + "/go.mod",
		outDir + "/.gitignore",
		outDir + "/Makefile",
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", path, err)
		}
	}
}
