package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/had-nu/lazy.go/pkg/wizard"
)

// ---- Styles ----------------------------------------------------------------

var (
	stylePrimary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	styleSecondary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	styleMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	styleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2)

	styleHeader = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 2).
			Bold(true)

	styleSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	styleUnselected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB"))

	styleProgress = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	styleHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4B5563")).
			Italic(true)
)

// ---- Model -----------------------------------------------------------------

// Model is the top-level BubbleTea model for the wizard.
type Model struct {
	state     wizard.WizardState
	textInput textinput.Model
	selection int    // cursor index for list/toggle steps
	toggles   []bool // for feature checkboxes
	validErr  string
	done      bool
	width     int
}

// New creates a fresh TUI Model.
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "type here..."
	ti.CharLimit = 128
	ti.Focus()

	state := wizard.NewWizardState()
	features := wizard.DefaultFeatures()
	state.Features = features

	// Build toggles slice from FeatureChoices order.
	fcs := wizard.FeatureChoices()
	toggles := make([]bool, len(fcs))
	for i, fc := range fcs {
		toggles[i] = features[fc.Key]
	}

	return Model{
		state:     state,
		textInput: ti,
		toggles:   toggles,
		width:     80,
	}
}

// ---- Messages --------------------------------------------------------------

type errMsg string
type doneMsg struct{}

// ---- Init ------------------------------------------------------------------

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// ---- Update ----------------------------------------------------------------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	isTextStep := isTextInputStep(m.state.CurrentStep)

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "enter":
		return m.advance()

	case "up", "k":
		if !isTextStep && m.selection > 0 {
			m.selection--
		}

	case "down", "j":
		if !isTextStep {
			max := maxSelection(m.state.CurrentStep)
			if m.selection < max {
				m.selection++
			}
		}

	case " ":
		if m.state.CurrentStep == wizard.StepFeatures {
			m.toggles[m.selection] = !m.toggles[m.selection]
		}
	}

	var cmd tea.Cmd
	if isTextStep {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

// advance validates and transitions to the next wizard step.
func (m Model) advance() (tea.Model, tea.Cmd) {
	m.validErr = ""

	if err := m.applyCurrentStep(); err != nil {
		m.validErr = err.Error()
		return m, nil
	}

	next := wizard.NextStep(m.state)
	if next == wizard.StepDone {
		m.done = true
		return m, tea.Quit
	}

	m.state.CurrentStep = next
	m.selection = 0
	m.prepareStepInput()
	return m, textinput.Blink
}

// applyCurrentStep reads controller input and stores answer in state.
func (m *Model) applyCurrentStep() error {
	switch m.state.CurrentStep {
	case wizard.StepProjectName:
		v := strings.TrimSpace(m.textInput.Value())
		if err := wizard.ValidateProjectName(v); err != nil {
			return err
		}
		m.state.ProjectName = wizard.SanitizeProjectName(v)

	case wizard.StepModulePath:
		v := strings.TrimSpace(m.textInput.Value())
		if err := wizard.ValidateModulePath(v); err != nil {
			return err
		}
		m.state.ModulePath = v

	case wizard.StepDescription:
		v := strings.TrimSpace(m.textInput.Value())
		if err := wizard.ValidateDescription(v); err != nil {
			return err
		}
		m.state.Description = v

	case wizard.StepAuthor:
		v := strings.TrimSpace(m.textInput.Value())
		if err := wizard.ValidateAuthor(v); err != nil {
			return err
		}
		m.state.Author = v

	case wizard.StepProjectType:
		choices := wizard.ProjectTypeChoices()
		if m.selection >= len(choices) {
			return fmt.Errorf("invalid selection")
		}
		m.state.ProjectType = choices[m.selection].Value

	case wizard.StepVisibility:
		choices := wizard.VisibilityChoices()
		if m.selection >= len(choices) {
			return fmt.Errorf("invalid selection")
		}
		m.state.Visibility = choices[m.selection].Value

	case wizard.StepCriticality:
		choices := wizard.CriticalityChoices()
		if m.selection >= len(choices) {
			return fmt.Errorf("invalid selection")
		}
		m.state.Criticality = choices[m.selection].Value

	case wizard.StepFeatures:
		fcs := wizard.FeatureChoices()
		for i, fc := range fcs {
			m.state.Features[fc.Key] = m.toggles[i]
		}

	case wizard.StepLicense:
		choices := wizard.LicenseChoices()
		if m.selection >= len(choices) {
			return fmt.Errorf("invalid selection")
		}
		m.state.License = choices[m.selection].Value

	case wizard.StepGitHub:
		choices := []string{"yes", "no"}
		m.state.GitHubEnable = choices[m.selection] == "yes"
		m.state.GitHubPush = m.state.GitHubEnable
	}

	return nil
}

// prepareStepInput initialises text input for text-based steps.
func (m *Model) prepareStepInput() {
	step := m.state.CurrentStep
	if isTextInputStep(step) {
		m.textInput.Reset()
		m.textInput.Focus()
		switch step {
		case wizard.StepProjectName:
			m.textInput.Placeholder = "e.g. my-service"
		case wizard.StepModulePath:
			m.textInput.Placeholder = "e.g. github.com/user/my-service"
		case wizard.StepDescription:
			m.textInput.Placeholder = "A short project description"
		case wizard.StepAuthor:
			m.textInput.Placeholder = "Your Name <email>"
		}
	}
}

// isTextInputStep returns true for steps that use a text input.
func isTextInputStep(s wizard.Step) bool {
	switch s {
	case wizard.StepProjectName, wizard.StepModulePath, wizard.StepDescription, wizard.StepAuthor:
		return true
	}
	return false
}

// maxSelection returns the max cursor index for list steps.
func maxSelection(s wizard.Step) int {
	switch s {
	case wizard.StepProjectType:
		return len(wizard.ProjectTypeChoices()) - 1
	case wizard.StepVisibility:
		return len(wizard.VisibilityChoices()) - 1
	case wizard.StepCriticality:
		return len(wizard.CriticalityChoices()) - 1
	case wizard.StepLicense:
		return len(wizard.LicenseChoices()) - 1
	case wizard.StepFeatures:
		return len(wizard.FeatureChoices()) - 1
	case wizard.StepGitHub:
		return 1
	}
	return 0
}

// State returns the completed wizard state (call after done == true).
func (m Model) State() wizard.WizardState {
	return m.state
}

// Done returns true when the wizard has completed.
func (m Model) Done() bool {
	return m.done
}
