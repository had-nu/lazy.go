package wizard_test

import (
	"testing"

	"github.com/had-nu/lazy.go/pkg/wizard"
)

func TestValidateProjectName_Valid(t *testing.T) {
	cases := []string{"myapp", "my-app", "my_app", "App123", "a"}
	for _, c := range cases {
		if err := wizard.ValidateProjectName(c); err != nil {
			t.Errorf("expected %q to be valid: %v", c, err)
		}
	}
}

func TestValidateProjectName_Invalid(t *testing.T) {
	cases := []string{"", "1app", "my app", "my.app", "../evil"}
	for _, c := range cases {
		if err := wizard.ValidateProjectName(c); err == nil {
			t.Errorf("expected %q to be invalid", c)
		}
	}
}

func TestValidateModulePath_Valid(t *testing.T) {
	cases := []string{
		"github.com/user/project",
		"github.com/user/my-project",
		"example.com/org/svc",
	}
	for _, c := range cases {
		if err := wizard.ValidateModulePath(c); err != nil {
			t.Errorf("expected %q to be valid: %v", c, err)
		}
	}
}

func TestValidateModulePath_Invalid(t *testing.T) {
	cases := []string{
		"",
		"../traversal",
		"noSlash",
		"host//double",
	}
	for _, c := range cases {
		if err := wizard.ValidateModulePath(c); err == nil {
			t.Errorf("expected %q to be invalid", c)
		}
	}
}

func TestValidateDescription_TooLong(t *testing.T) {
	long := make([]byte, 300)
	for i := range long {
		long[i] = 'a'
	}
	if err := wizard.ValidateDescription(string(long)); err == nil {
		t.Error("expected error for description > 256 chars")
	}
}

func TestValidateAuthor_Empty(t *testing.T) {
	if err := wizard.ValidateAuthor(""); err == nil {
		t.Error("expected error for empty author")
	}
}

func TestValidateAuthor_Valid(t *testing.T) {
	if err := wizard.ValidateAuthor("John Doe <john@example.com>"); err != nil {
		t.Errorf("expected valid author: %v", err)
	}
}

func TestSanitizeProjectName(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"My App", "my-app"},
		{"UPPER", "upper"},
		{"already-fine", "already-fine"},
	}
	for _, c := range cases {
		got := wizard.SanitizeProjectName(c.in)
		if got != c.out {
			t.Errorf("SanitizeProjectName(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

func TestNextStep_Sequence(t *testing.T) {
	state := wizard.NewWizardState()
	// Walk through all steps.
	for state.CurrentStep != wizard.StepDone {
		next := wizard.NextStep(state)
		if next <= state.CurrentStep && next != wizard.StepDone {
			t.Fatalf("NextStep did not advance: %v → %v", state.CurrentStep, next)
		}
		state.CurrentStep = next
	}
}
