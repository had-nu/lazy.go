package github

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// ValidateAuth checks that the user is authenticated with the GitHub CLI.
// It runs `gh auth status` and returns an error if not authenticated.
func ValidateAuth() error {
	cmd := exec.Command("gh", "auth", "status") //nolint:gosec
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("not authenticated with GitHub CLI: %w\nRun: gh auth login\nOutput: %s", err, string(out))
	}
	return nil
}

// TokenFromEnv attempts to read a GitHub token from environment variables.
// Returns empty string if not found — callers should fall back to gh CLI auth.
func TokenFromEnv() string {
	for _, key := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if v := lookupEnv(key); v != "" {
			return v
		}
	}
	return ""
}

// lookupEnv is os.LookupEnv extracted for testability.
var lookupEnv = func(key string) string {
	v, _ := os.LookupEnv(key)
	return v
}

// ErrNotAuthenticated is returned when GitHub CLI is not logged in.
var ErrNotAuthenticated = errors.New("not authenticated with GitHub CLI; run: gh auth login")
