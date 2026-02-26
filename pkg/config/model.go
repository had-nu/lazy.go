package config

// ProjectType represents the type of Go project.
type ProjectType string

const (
	ProjectTypeCLI         ProjectType = "cli"
	ProjectTypeAPI         ProjectType = "api"
	ProjectTypeMicroservice ProjectType = "microservice"
	ProjectTypeLibrary     ProjectType = "library"
	ProjectTypeSecurity    ProjectType = "security"
	ProjectTypeWorker      ProjectType = "worker"
)

// Visibility controls repository access.
type Visibility string

const (
	VisibilityInternal Visibility = "internal"
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
)

// LicenseType identifies the desired open-source license.
type LicenseType string

const (
	LicenseMIT        LicenseType = "mit"
	LicenseGPL3       LicenseType = "gpl-3.0"
	LicenseApache2    LicenseType = "apache-2.0"
	LicenseProprietary LicenseType = "proprietary"
)

// CriticalityLevel describes the operational risk of the project.
type CriticalityLevel string

const (
	CriticalityExperimental CriticalityLevel = "experimental"
	CriticalityProduction   CriticalityLevel = "production"
	CriticalitySecurity     CriticalityLevel = "security-critical"
)

// Features represents optional capabilities to enable in the project.
type Features struct {
	Docker         bool `yaml:"docker"`
	GitHubActions  bool `yaml:"github_actions"`
	Linting        bool `yaml:"linting"`
	StaticAnalysis bool `yaml:"static_analysis"`
	Dependabot     bool `yaml:"dependabot"`
	Tests          bool `yaml:"tests"`
	SAST           bool `yaml:"sast"`
}

// ProjectConfig is the central configuration object for a lazy.go project.
type ProjectConfig struct {
	Name        string           `yaml:"name"`
	ModulePath  string           `yaml:"module_path"`
	Description string           `yaml:"description"`
	Author      string           `yaml:"author"`
	Type        ProjectType      `yaml:"type"`
	Visibility  Visibility       `yaml:"visibility"`
	License     LicenseType      `yaml:"license"`
	Criticality CriticalityLevel `yaml:"criticality"`
	Features    Features         `yaml:"features"`
	GitHub      GitHubConfig     `yaml:"github"`
}

// GitHubConfig holds repository creation settings.
type GitHubConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Topics     []string `yaml:"topics,omitempty"`
	PushOnInit bool     `yaml:"push_on_init"`
}

// IsPublic returns true if the project is intended for public consumption.
func (p *ProjectConfig) IsPublic() bool {
	return p.Visibility == VisibilityPublic
}

// IsSecure returns true if security tooling should be enforced.
func (p *ProjectConfig) IsSecure() bool {
	return p.Criticality == CriticalityProduction || p.Criticality == CriticalitySecurity
}

// AllProjectTypes returns all valid project type values.
func AllProjectTypes() []ProjectType {
	return []ProjectType{
		ProjectTypeCLI,
		ProjectTypeAPI,
		ProjectTypeMicroservice,
		ProjectTypeLibrary,
		ProjectTypeSecurity,
		ProjectTypeWorker,
	}
}

// AllLicenses returns all valid license values.
func AllLicenses() []LicenseType {
	return []LicenseType{
		LicenseMIT,
		LicenseGPL3,
		LicenseApache2,
		LicenseProprietary,
	}
}
