package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/had-nu/lazy.go/pkg/config"
)

func TestLoadFromYAML_RoundTrip(t *testing.T) {
	original := &config.ProjectConfig{
		Name:        "sentinel",
		ModulePath:  "github.com/user/sentinel",
		Description: "A security sentinel API",
		Author:      "Test Author",
		Type:        config.ProjectTypeAPI,
		Visibility:  config.VisibilityPublic,
		License:     config.LicenseApache2,
		Criticality: config.CriticalityProduction,
		Features: config.Features{
			Docker:         true,
			GitHubActions:  true,
			Tests:          true,
			StaticAnalysis: true,
		},
		GitHub: config.GitHubConfig{
			Enabled:    true,
			PushOnInit: false,
		},
	}

	dir := t.TempDir()
	yamlPath := filepath.Join(dir, "lazygo.yml")

	if err := config.ExportToYAML(original, yamlPath); err != nil {
		t.Fatalf("ExportToYAML: %v", err)
	}

	loaded, err := config.LoadFromYAML(yamlPath)
	if err != nil {
		t.Fatalf("LoadFromYAML: %v", err)
	}

	if !reflect.DeepEqual(original, loaded) {
		t.Errorf("round-trip mismatch:\ngot:  %+v\nwant: %+v", loaded, original)
	}
}

func TestExportToYAML_FileContents(t *testing.T) {
	cfg := &config.ProjectConfig{
		Name:       "testproj",
		ModulePath: "github.com/user/testproj",
		Type:       config.ProjectTypeCLI,
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "lazygo.yml")

	if err := config.ExportToYAML(cfg, path); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Fatal("exported file is empty")
	}
}

func TestValidate_MissingName(t *testing.T) {
	cfg := &config.ProjectConfig{
		ModulePath: "github.com/x/y",
		Type:       config.ProjectTypeCLI,
	}
	if err := config.Validate(cfg); err == nil {
		t.Error("expected error for missing name")
	}
}

func TestValidate_InvalidType(t *testing.T) {
	cfg := &config.ProjectConfig{
		Name:       "valid",
		ModulePath: "github.com/x/valid",
		Type:       config.ProjectType("bogus"),
	}
	if err := config.Validate(cfg); err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestIsPublic(t *testing.T) {
	pub := &config.ProjectConfig{Visibility: config.VisibilityPublic}
	priv := &config.ProjectConfig{Visibility: config.VisibilityPrivate}

	if !pub.IsPublic() {
		t.Error("expected public to be public")
	}
	if priv.IsPublic() {
		t.Error("expected private to not be public")
	}
}

func TestIsSecure(t *testing.T) {
	prod := &config.ProjectConfig{Criticality: config.CriticalityProduction}
	sec := &config.ProjectConfig{Criticality: config.CriticalitySecurity}
	exp := &config.ProjectConfig{Criticality: config.CriticalityExperimental}

	if !prod.IsSecure() {
		t.Error("production should be secure")
	}
	if !sec.IsSecure() {
		t.Error("security-critical should be secure")
	}
	if exp.IsSecure() {
		t.Error("experimental should not be secure")
	}
}
