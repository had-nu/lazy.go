package wizard

// Step represents a single step in the interactive wizard.
type Step int

const (
	StepProjectName Step = iota
	StepModulePath
	StepDescription
	StepAuthor
	StepProjectType
	StepVisibility
	StepCriticality
	StepFeatures
	StepLicense
	StepGitHub
	StepDone
)

// String returns a human-readable label for the step.
func (s Step) String() string {
	switch s {
	case StepProjectName:
		return "Project Name"
	case StepModulePath:
		return "Module Path"
	case StepDescription:
		return "Description"
	case StepAuthor:
		return "Author"
	case StepProjectType:
		return "Project Type"
	case StepVisibility:
		return "Visibility"
	case StepCriticality:
		return "Criticality"
	case StepFeatures:
		return "Features"
	case StepLicense:
		return "License"
	case StepGitHub:
		return "GitHub Integration"
	case StepDone:
		return "Done"
	default:
		return "Unknown"
	}
}

// TotalSteps is the total number of wizard steps (excluding StepDone).
const TotalSteps = int(StepDone)

// WizardState holds all answers collected by the wizard so far.
type WizardState struct {
	CurrentStep  Step
	ProjectName  string
	ModulePath   string
	Description  string
	Author       string
	ProjectType  string
	Visibility   string
	Criticality  string
	Features     map[string]bool
	License      string
	GitHubEnable bool
	GitHubPush   bool
}

// NewWizardState initialises a fresh wizard state at the first step.
func NewWizardState() WizardState {
	return WizardState{
		CurrentStep: StepProjectName,
		Features:    make(map[string]bool),
	}
}
