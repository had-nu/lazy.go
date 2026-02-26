package wizard

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	validProjectName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]{0,63}$`)
	validModulePath  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-._/~]*$`)
)

// ValidateProjectName checks that name is safe for directory and module use.
func ValidateProjectName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if !validProjectName.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, digits, hyphens, or underscores (max 64 chars)")
	}
	return nil
}

// ValidateModulePath checks that the Go module path is valid.
func ValidateModulePath(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("module path must not contain '..'")
	}
	if !validModulePath.MatchString(path) {
		return fmt.Errorf("module path contains invalid characters")
	}
	// Must have at least one slash for a proper module path (e.g. github.com/user/project)
	// Exception: single-segment paths are allowed for local/simple modules
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return fmt.Errorf("module path should be in the form <host>/<user>/<project> (e.g. github.com/user/myapp)")
	}
	for _, p := range parts {
		if p == "" {
			return fmt.Errorf("module path contains empty segment")
		}
	}
	return nil
}

// ValidateDescription validates project description length.
func ValidateDescription(desc string) error {
	desc = strings.TrimSpace(desc)
	if len(desc) > 256 {
		return fmt.Errorf("description must be 256 characters or fewer")
	}
	return nil
}

// ValidateAuthor validates author/maintainer field.
func ValidateAuthor(author string) error {
	author = strings.TrimSpace(author)
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}
	for _, r := range author {
		if !unicode.IsPrint(r) {
			return fmt.Errorf("author contains non-printable characters")
		}
	}
	return nil
}

// SanitizeProjectName returns a safe version of the project name.
func SanitizeProjectName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name)
}
