package wizard

import (
	"fmt"
	"strings"

	"github.com/hadnu/lazy.go/pkg/config"
)

// NextStep returns the next wizard step based on current state.
// This enables conditional flow (e.g. skip some steps for library projects).
func NextStep(state WizardState) Step {
	switch state.CurrentStep {
	case StepProjectName:
		return StepModulePath
	case StepModulePath:
		return StepDescription
	case StepDescription:
		return StepAuthor
	case StepAuthor:
		return StepProjectType
	case StepProjectType:
		return StepVisibility
	case StepVisibility:
		return StepCriticality
	case StepCriticality:
		return StepFeatures
	case StepFeatures:
		return StepLicense
	case StepLicense:
		return StepGitHub
	case StepGitHub:
		return StepDone
	default:
		return StepDone
	}
}

// BuildConfig converts a completed WizardState into a ProjectConfig.
func BuildConfig(state WizardState) *config.ProjectConfig {
	cfg := &config.ProjectConfig{
		Name:        state.ProjectName,
		ModulePath:  state.ModulePath,
		Description: state.Description,
		Author:      state.Author,
		Type:        config.ProjectType(state.ProjectType),
		Visibility:  config.Visibility(state.Visibility),
		Criticality: config.CriticalityLevel(state.Criticality),
		License:     config.LicenseType(state.License),
		Features: config.Features{
			Docker:         state.Features["docker"],
			GitHubActions:  state.Features["github_actions"],
			Linting:        state.Features["linting"],
			StaticAnalysis: state.Features["static_analysis"],
			Dependabot:     state.Features["dependabot"],
			Tests:          state.Features["tests"],
			SAST:           state.Features["sast"],
		},
		GitHub: config.GitHubConfig{
			Enabled:    state.GitHubEnable,
			PushOnInit: state.GitHubPush,
		},
	}

	// Auto-enable security for production / security-critical projects.
	if cfg.IsSecure() {
		cfg.Features.StaticAnalysis = true
		cfg.Features.SAST = true
		cfg.Features.Tests = true
		if cfg.Features.GitHubActions {
			cfg.Features.Dependabot = true
		}
	}

	// Auto-suggest license when not set.
	if cfg.License == "" {
		cfg.License = SuggestLicense(cfg)
	}

	return cfg
}

// SuggestLicense returns the recommended license for a project configuration.
func SuggestLicense(cfg *config.ProjectConfig) config.LicenseType {
	switch {
	case cfg.Visibility == config.VisibilityPrivate:
		return config.LicenseProprietary
	case cfg.Type == config.ProjectTypeLibrary && cfg.Visibility == config.VisibilityPublic:
		return config.LicenseApache2
	case cfg.Visibility == config.VisibilityPublic:
		return config.LicenseMIT
	default:
		return config.LicenseProprietary
	}
}

// ProgressPercent returns wizard completion as 0-100.
func ProgressPercent(state WizardState) int {
	if TotalSteps == 0 {
		return 100
	}
	return int(state.CurrentStep) * 100 / TotalSteps
}

// ProjectTypeChoices returns display labels → values for the type selection.
func ProjectTypeChoices() []Choice {
	return []Choice{
		{Label: "CLI Tool", Value: string(config.ProjectTypeCLI)},
		{Label: "REST API", Value: string(config.ProjectTypeAPI)},
		{Label: "Microservice", Value: string(config.ProjectTypeMicroservice)},
		{Label: "Library", Value: string(config.ProjectTypeLibrary)},
		{Label: "Security Tool", Value: string(config.ProjectTypeSecurity)},
		{Label: "Concurrent Worker / Service", Value: string(config.ProjectTypeWorker)},
	}
}

// VisibilityChoices returns display labels → values for visibility.
func VisibilityChoices() []Choice {
	return []Choice{
		{Label: "Public (Open Source)", Value: string(config.VisibilityPublic)},
		{Label: "Private (Internal)", Value: string(config.VisibilityInternal)},
		{Label: "Private (Commercial)", Value: string(config.VisibilityPrivate)},
	}
}

// CriticalityChoices returns display labels → values for criticality.
func CriticalityChoices() []Choice {
	return []Choice{
		{Label: "Experimental", Value: string(config.CriticalityExperimental)},
		{Label: "Production", Value: string(config.CriticalityProduction)},
		{Label: "Security Critical", Value: string(config.CriticalitySecurity)},
	}
}

// LicenseChoices returns display labels → values for license selection.
func LicenseChoices() []Choice {
	return []Choice{
		{Label: fmt.Sprintf("Auto-suggest (recommended: %s)", "based on context"), Value: "auto"},
		{Label: "MIT", Value: string(config.LicenseMIT)},
		{Label: "Apache-2.0", Value: string(config.LicenseApache2)},
		{Label: "GPL-3.0", Value: string(config.LicenseGPL3)},
		{Label: "Proprietary (no license)", Value: string(config.LicenseProprietary)},
	}
}

// FeatureChoices returns all optional features with labels.
func FeatureChoices() []ToggleChoice {
	return []ToggleChoice{
		{Key: "tests", Label: "Unit Tests", Default: true},
		{Key: "linting", Label: "Linting (golangci-lint)"},
		{Key: "static_analysis", Label: "Static Analysis (staticcheck, gosec)"},
		{Key: "github_actions", Label: "GitHub Actions CI"},
		{Key: "docker", Label: "Docker"},
		{Key: "dependabot", Label: "Dependabot"},
		{Key: "sast", Label: "SAST / govulncheck"},
	}
}

// Choice is a labeled selection option.
type Choice struct {
	Label string
	Value string
}

// ToggleChoice is a feature toggle option.
type ToggleChoice struct {
	Key     string
	Label   string
	Default bool
}

// DefaultFeatures returns the default feature map (all features set to their default values).
func DefaultFeatures() map[string]bool {
	m := make(map[string]bool)
	for _, fc := range FeatureChoices() {
		m[fc.Key] = fc.Default
	}
	return m
}

// LicenseFromChoice resolves "auto" license to a concrete suggestion.
func LicenseFromChoice(value string, cfg *config.ProjectConfig) config.LicenseType {
	if strings.ToLower(value) == "auto" {
		return SuggestLicense(cfg)
	}
	return config.LicenseType(value)
}
