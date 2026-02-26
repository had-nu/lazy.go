package scaffold_test

import (
	"testing"

	"github.com/hadnu/lazy.go/pkg/config"
	"github.com/hadnu/lazy.go/pkg/scaffold"
)

func TestGenerateLicense_MIT(t *testing.T) {
	text := scaffold.GenerateLicense(config.LicenseMIT, "Test Author", 2026)
	if text == "" {
		t.Fatal("expected non-empty license text")
	}
	if !contains(text, "MIT License") {
		t.Error("expected MIT header in license")
	}
	if !contains(text, "Test Author") {
		t.Error("expected author in license")
	}
	if !contains(text, "2026") {
		t.Error("expected year in license")
	}
}

func TestGenerateLicense_Apache(t *testing.T) {
	text := scaffold.GenerateLicense(config.LicenseApache2, "Corp", 2026)
	if !contains(text, "Apache License") {
		t.Error("expected Apache header")
	}
}

func TestGenerateLicense_GPL(t *testing.T) {
	text := scaffold.GenerateLicense(config.LicenseGPL3, "Corp", 2026)
	if !contains(text, "GNU GENERAL PUBLIC LICENSE") {
		t.Error("expected GPL header")
	}
}

func TestGenerateLicense_Proprietary(t *testing.T) {
	text := scaffold.GenerateLicense(config.LicenseProprietary, "Corp Inc.", 2026)
	if !contains(text, "PROPRIETARY") {
		t.Error("expected PROPRIETARY header")
	}
	if !contains(text, "Corp Inc.") {
		t.Error("expected author in proprietary notice")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := range s {
		if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
