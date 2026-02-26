package security

import (
	"fmt"
	"strings"

	"github.com/hadnu/lazy.go/pkg/config"
)

// ShouldEnableSecurity returns true when the project warrants security tooling.
func ShouldEnableSecurity(cfg *config.ProjectConfig) bool {
	return cfg.IsSecure()
}

// EnforceSecurity mutates Features to activate all mandatory security options
// based on the project criticality level.
func EnforceSecurity(cfg *config.ProjectConfig) {
	if !ShouldEnableSecurity(cfg) {
		return
	}
	cfg.Features.StaticAnalysis = true
	cfg.Features.SAST = true
	cfg.Features.Tests = true
	if cfg.Features.GitHubActions {
		cfg.Features.Dependabot = true
	}
}

// GolangCIConfig generates a .golangci.yml configuration string.
func GolangCIConfig(cfg *config.ProjectConfig) string {
	var sb strings.Builder
	sb.WriteString("run:\n  timeout: 5m\n  go: \"1.22\"\n\n")
	sb.WriteString("linters:\n  enable:\n")

	base := []string{
		"errcheck", "gosimple", "govet", "ineffassign",
		"staticcheck", "unused", "gofmt", "misspell", "revive",
	}
	for _, l := range base {
		fmt.Fprintf(&sb, "    - %s\n", l)
	}
	if cfg.Features.SAST {
		sb.WriteString("    - gosec\n")
	}
	return sb.String()
}

// SecurityMD generates the content for SECURITY.md.
func SecurityMD(cfg *config.ProjectConfig) string {
	return fmt.Sprintf(`# Security Policy

## Reporting a Vulnerability

**Do not open a public GitHub issue.** Contact the maintainer at **%s** with:

- Description
- Steps to reproduce
- Impact assessment

You will receive acknowledgment within 48 hours.

## Practices

- Dependabot for dependency updates
- Static analysis: gosec, staticcheck
- govulncheck on every CI run
- Race condition detection: go test -race
`, cfg.Author)
}

// DependabotConfig returns a dependabot.yml string.
func DependabotConfig() string {
	return `version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
`
}
