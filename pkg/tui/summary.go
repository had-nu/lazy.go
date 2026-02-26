package tui

import (
	"strings"

	"github.com/hadnu/lazy.go/pkg/config"
	"github.com/hadnu/lazy.go/pkg/wizard"
)

// RenderSummary returns a human-readable summary of the collected configuration.
func RenderSummary(state wizard.WizardState) string {
	cfg := wizard.BuildConfig(state)
	return renderConfigTable(cfg)
}

func renderConfigTable(cfg *config.ProjectConfig) string {
	var sb strings.Builder

	sb.WriteString(styleHeader.Render(" 📋 Project Summary ") + "\n\n")

	rows := [][]string{
		{"Name", cfg.Name},
		{"Module", cfg.ModulePath},
		{"Description", cfg.Description},
		{"Author", cfg.Author},
		{"Type", string(cfg.Type)},
		{"Visibility", string(cfg.Visibility)},
		{"Criticality", string(cfg.Criticality)},
		{"License", string(cfg.License)},
	}

	for _, row := range rows {
		label := stylePrimary.Render(padRight(row[0]+":", 14))
		value := styleSecondary.Render(row[1])
		sb.WriteString("  " + label + " " + value + "\n")
	}

	sb.WriteString("\n  " + stylePrimary.Render("Features:") + "\n")
	appendFeature(&sb, "Tests", cfg.Features.Tests)
	appendFeature(&sb, "Linting", cfg.Features.Linting)
	appendFeature(&sb, "Static Analysis", cfg.Features.StaticAnalysis)
	appendFeature(&sb, "SAST", cfg.Features.SAST)
	appendFeature(&sb, "Docker", cfg.Features.Docker)
	appendFeature(&sb, "GitHub Actions", cfg.Features.GitHubActions)
	appendFeature(&sb, "Dependabot", cfg.Features.Dependabot)

	if cfg.GitHub.Enabled {
		sb.WriteString("\n  " + styleSuccess.Render("✓ GitHub repository will be created") + "\n")
	}

	return styleBox.Render(sb.String())
}

func appendFeature(sb *strings.Builder, name string, enabled bool) {
	icon := styleSuccess.Render("✓")
	if !enabled {
		icon = styleMuted.Render("✗")
	}
	label := styleMuted.Render(name)
	if enabled {
		label = styleUnselected.Render(name)
	}
	sb.WriteString("    " + icon + "  " + label + "\n")
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}
