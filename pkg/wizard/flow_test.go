package wizard

import (
	"testing"

	"github.com/had-nu/lazy.go/pkg/config"
)

func TestBuildConfig_EnforcesSecurityForProduction(t *testing.T) {
	state := WizardState{
		ProjectName: "svc",
		ModulePath:  "github.com/x/svc",
		ProjectType: string(config.ProjectTypeAPI),
		Criticality: string(config.CriticalityProduction),
		Features:    map[string]bool{"github_actions": true},
	}

	cfg := BuildConfig(state)

	if !cfg.Features.StaticAnalysis {
		t.Error("production project must have StaticAnalysis")
	}
	if !cfg.Features.SAST {
		t.Error("production project must have SAST")
	}
	if !cfg.Features.Tests {
		t.Error("production project must have Tests")
	}
	if !cfg.Features.Dependabot {
		t.Error("production project with GH Actions must have Dependabot")
	}
}
