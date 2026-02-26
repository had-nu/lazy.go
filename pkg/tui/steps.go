package tui

import (
	"fmt"
	"strings"

	"github.com/had-nu/lazy.go/pkg/wizard"
)

// View renders the current state to a string for BubbleTea.
func (m Model) View() string {
	if m.done {
		return renderDone()
	}

	var sb strings.Builder
	sb.WriteString(renderHeader(m.state))
	sb.WriteString("\n\n")
	sb.WriteString(renderStep(m))
	sb.WriteString("\n\n")
	if m.validErr != "" {
		sb.WriteString(styleError.Render("✗ "+m.validErr) + "\n\n")
	}
	sb.WriteString(renderHints(m.state.CurrentStep))
	return sb.String()
}

func renderHeader(state wizard.WizardState) string {
	title := styleHeader.Render(" lazy.go — Go Project Generator ")
	progress := renderProgressBar(wizard.ProgressPercent(state), 40)
	stepLabel := styleMuted.Render(fmt.Sprintf(" Step %d/%d — %s",
		int(state.CurrentStep)+1, wizard.TotalSteps, state.CurrentStep.String()))
	return title + "\n" + progress + stepLabel
}

func renderProgressBar(percent, width int) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return styleProgress.Render(bar) + " " + styleMuted.Render(fmt.Sprintf("%d%%", percent)) + "\n"
}

func renderStep(m Model) string {
	step := m.state.CurrentStep
	switch {
	case isTextInputStep(step):
		return renderTextInput(m)
	case step == wizard.StepFeatures:
		return renderFeatureToggles(m)
	case step == wizard.StepGitHub:
		return renderGitHubStep(m)
	default:
		return renderListSelection(m)
	}
}

func renderTextInput(m Model) string {
	prompt := stepPrompt(m.state.CurrentStep)
	return styleBox.Render(
		stylePrimary.Render(prompt) + "\n\n" +
			m.textInput.View() + "\n",
	)
}

func renderListSelection(m Model) string {
	choices := stepChoices(m.state.CurrentStep)
	prompt := stepPrompt(m.state.CurrentStep)

	var sb strings.Builder
	sb.WriteString(stylePrimary.Render(prompt) + "\n\n")
	for i, c := range choices {
		if i == m.selection {
			sb.WriteString(styleSelected.Render("▶  "+c.Label) + "\n")
		} else {
			sb.WriteString(styleUnselected.Render("   "+c.Label) + "\n")
		}
	}

	return styleBox.Render(sb.String())
}

func renderFeatureToggles(m Model) string {
	fcs := wizard.FeatureChoices()
	var sb strings.Builder
	sb.WriteString(stylePrimary.Render("Select features to enable:") + "\n\n")
	for i, fc := range fcs {
		cursor := "  "
		if i == m.selection {
			cursor = styleSelected.Render("▶ ")
		}
		toggle := "☐"
		if m.toggles[i] {
			toggle = styleSuccess.Render("☑")
		}
		sb.WriteString(fmt.Sprintf("%s%s  %s\n", cursor, toggle,
			styleUnselected.Render(fc.Label)))
	}
	return styleBox.Render(sb.String())
}

func renderGitHubStep(m Model) string {
	opts := []string{"Yes — create GitHub repository and push", "No — local project only"}
	var sb strings.Builder
	sb.WriteString(stylePrimary.Render("Create a GitHub repository?") + "\n\n")
	for i, opt := range opts {
		if i == m.selection {
			sb.WriteString(styleSelected.Render("▶  "+opt) + "\n")
		} else {
			sb.WriteString(styleUnselected.Render("   "+opt) + "\n")
		}
	}
	return styleBox.Render(sb.String())
}

func renderHints(step wizard.Step) string {
	var hints []string
	if isTextInputStep(step) {
		hints = append(hints, "enter → next")
	} else if step == wizard.StepFeatures {
		hints = append(hints, "↑/↓ move   space toggle   enter → next")
	} else {
		hints = append(hints, "↑/↓ move   enter → next")
	}
	hints = append(hints, "ctrl+c → quit")
	return styleHint.Render("  " + strings.Join(hints, "   "))
}

func renderDone() string {
	return "\n" + styleSuccess.Render("  ✓ Project configuration complete!") +
		"\n" + styleMuted.Render("  Generating your project...") + "\n\n"
}

// ---- Helpers ---------------------------------------------------------------

type labeledChoice struct {
	Label string
	Value string
}

func stepChoices(step wizard.Step) []labeledChoice {
	var out []labeledChoice
	switch step {
	case wizard.StepProjectType:
		for _, c := range wizard.ProjectTypeChoices() {
			out = append(out, labeledChoice{c.Label, c.Value})
		}
	case wizard.StepVisibility:
		for _, c := range wizard.VisibilityChoices() {
			out = append(out, labeledChoice{c.Label, c.Value})
		}
	case wizard.StepCriticality:
		for _, c := range wizard.CriticalityChoices() {
			out = append(out, labeledChoice{c.Label, c.Value})
		}
	case wizard.StepLicense:
		for _, c := range wizard.LicenseChoices() {
			out = append(out, labeledChoice{c.Label, c.Value})
		}
	}
	return out
}

func stepPrompt(step wizard.Step) string {
	switch step {
	case wizard.StepProjectName:
		return "What is the name of your project?"
	case wizard.StepModulePath:
		return "What is the Go module path?"
	case wizard.StepDescription:
		return "Briefly describe your project:"
	case wizard.StepAuthor:
		return "Your name / maintainer:"
	case wizard.StepProjectType:
		return "What type of project is this?"
	case wizard.StepVisibility:
		return "Who is this project for?"
	case wizard.StepCriticality:
		return "What is the criticality level?"
	case wizard.StepLicense:
		return "Choose a license:"
	default:
		return step.String()
	}
}
