package scaffold_test

import (
	"testing"

	"github.com/hadnu/lazy.go/pkg/config"
	"github.com/hadnu/lazy.go/pkg/scaffold"
)

func cfg(t config.ProjectType) *config.ProjectConfig {
	return &config.ProjectConfig{
		Name:       "myapp",
		ModulePath: "github.com/user/myapp",
		Type:       t,
		Author:     "Test Author",
	}
}

func TestBuildDirectoryTree_Library(t *testing.T) {
	entries := scaffold.BuildDirectoryTree(cfg(config.ProjectTypeLibrary))
	assertContainsPath(t, entries, "go.mod")
	assertContainsPath(t, entries, "README.md")
	assertContainsPath(t, entries, ".gitignore")
	assertContainsPath(t, entries, "pkg/myapp/myapp.go")
	assertContainsPath(t, entries, "pkg/myapp/myapp_test.go")
	// Libraries should NOT have cmd/
	assertNotContainsPrefix(t, entries, "cmd/")
}

func TestBuildDirectoryTree_CLI(t *testing.T) {
	entries := scaffold.BuildDirectoryTree(cfg(config.ProjectTypeCLI))
	assertContainsPath(t, entries, "main.go")
	assertContainsPath(t, entries, "cmd/root.go")
	assertContainsPath(t, entries, "internal/app/app.go")
	assertContainsPath(t, entries, "internal/config/config.go")
}

func TestBuildDirectoryTree_API(t *testing.T) {
	entries := scaffold.BuildDirectoryTree(cfg(config.ProjectTypeAPI))
	assertContainsPath(t, entries, "cmd/server/main.go")
	assertContainsPath(t, entries, "internal/handler/handler.go")
	assertContainsPath(t, entries, "internal/service/service.go")
	assertContainsPath(t, entries, "internal/repository/repository.go")
	assertContainsPath(t, entries, "internal/middleware/middleware.go")
	assertContainsPath(t, entries, "api/openapi.yaml")
}

func TestBuildDirectoryTree_SecurityTool(t *testing.T) {
	entries := scaffold.BuildDirectoryTree(cfg(config.ProjectTypeSecurity))
	assertContainsPath(t, entries, "internal/scanner/scanner.go")
	assertContainsPath(t, entries, "internal/report/report.go")
}

func TestBuildDirectoryTree_Docker(t *testing.T) {
	c := cfg(config.ProjectTypeCLI)
	c.Features.Docker = true
	entries := scaffold.BuildDirectoryTree(c)
	assertContainsPath(t, entries, "Dockerfile")
	assertContainsPath(t, entries, ".dockerignore")
}

func TestBuildDirectoryTree_GolangCI(t *testing.T) {
	c := cfg(config.ProjectTypeCLI)
	c.Features.Linting = true
	entries := scaffold.BuildDirectoryTree(c)
	assertContainsPath(t, entries, ".golangci.yml")
}

func TestBuildDirectoryTree_GitHubActions(t *testing.T) {
	c := cfg(config.ProjectTypeAPI)
	c.Features.GitHubActions = true
	entries := scaffold.BuildDirectoryTree(c)
	assertContainsPath(t, entries, ".github/workflows/ci.yml")
}

func TestBuildDirectoryTree_PublicProject(t *testing.T) {
	c := cfg(config.ProjectTypeLibrary)
	c.Visibility = config.VisibilityPublic
	entries := scaffold.BuildDirectoryTree(c)
	assertContainsPath(t, entries, "CONTRIBUTING.md")
	assertContainsPath(t, entries, "CODE_OF_CONDUCT.md")
}

func TestBuildDirectoryTree_SecureProject(t *testing.T) {
	c := cfg(config.ProjectTypeAPI)
	c.Criticality = config.CriticalityProduction
	entries := scaffold.BuildDirectoryTree(c)
	assertContainsPath(t, entries, "SECURITY.md")
}

// assertContainsPath fails if no entry has the given path.
func assertContainsPath(t *testing.T, entries []scaffold.DirEntry, path string) {
	t.Helper()
	for _, e := range entries {
		if e.Path == path {
			return
		}
	}
	t.Errorf("expected entry with path %q, not found in tree", path)
}

// assertNotContainsPrefix fails if any file entry starts with prefix.
func assertNotContainsPrefix(t *testing.T, entries []scaffold.DirEntry, prefix string) {
	t.Helper()
	for _, e := range entries {
		if !e.IsDir && len(e.Path) >= len(prefix) && e.Path[:len(prefix)] == prefix {
			t.Errorf("unexpected entry with prefix %q: %s", prefix, e.Path)
		}
	}
}
