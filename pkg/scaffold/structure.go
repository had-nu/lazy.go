package scaffold

import (
	"github.com/hadnu/lazy.go/pkg/config"
)

// BuildDirectoryTree returns the list of directories and files to create
// for a given project configuration. No unnecessary empty directories.
func BuildDirectoryTree(cfg *config.ProjectConfig) []DirEntry {
	data := commonData(cfg)
	var entries []DirEntry

	add := func(path, tmpl string, isDir bool) {
		entries = append(entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	dir := func(path string) { add(path, "", true) }
	file := func(path, tmpl string) { add(path, tmpl, false) }

	// Files common to all project types.
	file("README.md", "readme.tmpl")
	file("go.mod", "gomod.tmpl")
	file(".gitignore", "gitignore.tmpl")

	// LICENSE is written programmatically via scaffold.GenerateLicense().

	if cfg.IsPublic() {
		file("CONTRIBUTING.md", "contributing.tmpl")
		file("CODE_OF_CONDUCT.md", "coc.tmpl")
	}

	if cfg.IsSecure() {
		file("SECURITY.md", "security.tmpl")
	}

	if cfg.Features.Linting || cfg.Features.StaticAnalysis {
		file(".golangci.yml", "golangci.tmpl")
	}

	if cfg.Features.GitHubActions {
		dir(".github/workflows")
		file(".github/workflows/ci.yml", "workflow.tmpl")
	}

	if cfg.GitHub.Enabled {
		dir(".github")
		file(".github/PULL_REQUEST_TEMPLATE.md", "pr_template.tmpl")
	}

	if cfg.Features.Dependabot {
		dir(".github")
		file(".github/dependabot.yml", "dependabot.tmpl")
	}

	// Project-type specific structures.
	switch cfg.Type {
	case config.ProjectTypeLibrary:
		buildLibraryStructure(cfg, &entries, data)
	case config.ProjectTypeCLI:
		buildCLIStructure(cfg, &entries, data)
	case config.ProjectTypeAPI:
		buildAPIStructure(cfg, &entries, data)
	case config.ProjectTypeMicroservice:
		buildMicroserviceStructure(cfg, &entries, data)
	case config.ProjectTypeSecurity:
		buildSecurityToolStructure(cfg, &entries, data)
	case config.ProjectTypeWorker:
		buildWorkerStructure(cfg, &entries, data)
	}

	if cfg.Features.Docker {
		entries = append(entries, DirEntry{Path: "Dockerfile", IsDir: false, Template: "dockerfile.tmpl", Data: data})
		entries = append(entries, DirEntry{Path: ".dockerignore", IsDir: false, Template: "dockerignore.tmpl", Data: data})
	}

	return entries
}

func buildLibraryStructure(cfg *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	pkg := "pkg/" + data.LibName
	add(pkg+"/"+data.LibName+".go", "lib.tmpl", false)
	add(pkg+"/"+data.LibName+"_test.go", "lib_test.tmpl", false)
	if cfg.Features.Tests {
		add("Makefile", "makefile.tmpl", false)
	}
}

func buildCLIStructure(cfg *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	add("main.go", "main_cli.tmpl", false)
	add("cmd/root.go", "cmd_root.tmpl", false)
	add("internal/app/app.go", "internal_app.tmpl", false)
	add("internal/config/config.go", "internal_config.tmpl", false)
	add("Makefile", "makefile.tmpl", false)
}

func buildAPIStructure(cfg *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	add("cmd/server/main.go", "main_api.tmpl", false)
	add("internal/handler/handler.go", "handler.tmpl", false)
	add("internal/service/service.go", "service.tmpl", false)
	add("internal/repository/repository.go", "repository.tmpl", false)
	add("internal/middleware/middleware.go", "middleware.tmpl", false)
	add("internal/config/config.go", "internal_config.tmpl", false)
	add("api/openapi.yaml", "openapi.tmpl", false)
	add("Makefile", "makefile.tmpl", false)
}

func buildMicroserviceStructure(cfg *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	add("cmd/service/main.go", "main_api.tmpl", false)
	add("internal/handler/handler.go", "handler.tmpl", false)
	add("internal/service/service.go", "service.tmpl", false)
	add("internal/repository/repository.go", "repository.tmpl", false)
	add("internal/middleware/middleware.go", "middleware.tmpl", false)
	add("internal/config/config.go", "internal_config.tmpl", false)
	add("internal/worker/worker.go", "worker.tmpl", false)
	add("Makefile", "makefile.tmpl", false)
}

func buildSecurityToolStructure(cfg *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	add("main.go", "main_cli.tmpl", false)
	add("cmd/root.go", "cmd_root.tmpl", false)
	add("internal/scanner/scanner.go", "scanner.tmpl", false)
	add("internal/report/report.go", "report.tmpl", false)
	add("internal/config/config.go", "internal_config.tmpl", false)
	add("pkg/"+data.LibName+"/"+data.LibName+".go", "lib.tmpl", false)
	add("Makefile", "makefile.tmpl", false)
}

func buildWorkerStructure(_ *config.ProjectConfig, entries *[]DirEntry, data templateData) {
	add := func(path, tmpl string, isDir bool) {
		*entries = append(*entries, DirEntry{Path: path, IsDir: isDir, Template: tmpl, Data: data})
	}
	add("cmd/worker/main.go", "main_api.tmpl", false)
	add("internal/worker/worker.go", "worker.tmpl", false)
	add("internal/config/config.go", "internal_config.tmpl", false)
	add("Makefile", "makefile.tmpl", false)
}
