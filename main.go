package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/had-nu/lazy.go/pkg/config"
	ghpkg "github.com/had-nu/lazy.go/pkg/github"
	"github.com/had-nu/lazy.go/pkg/scaffold"
	"github.com/had-nu/lazy.go/pkg/security"
	"github.com/had-nu/lazy.go/pkg/tui"
	"github.com/had-nu/lazy.go/pkg/wizard"
)

const version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ---- Root Command ----------------------------------------------------------

var rootCmd = &cobra.Command{
	Use:   "lazy.go",
	Short: "Intelligent Go project and repository generator",
	Long: `lazy.go — Architectural coherence from intent.

An interactive TUI wizard that generates idiomatic Go project structures,
security policies, CI pipelines, and GitHub repositories — tailored to
your project's real purpose and risk profile.`,
}

func init() {
	rootCmd.AddCommand(initCmd, validateCmd, versionCmd)
}

// ---- init command ----------------------------------------------------------

var fromFile string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start the interactive project wizard",
	Long: `Start the lazy.go wizard to generate a new Go project.

Use --from to replay a saved lazygo.yml configuration without the wizard.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *config.ProjectConfig

		if fromFile != "" {
			// Headless mode: load config from YAML.
			loaded, err := config.LoadFromYAML(fromFile)
			if err != nil {
				return fmt.Errorf("loading config from %s: %w", fromFile, err)
			}
			cfg = loaded
			fmt.Println("✓ Loaded configuration from", fromFile)
		} else {
			// Interactive TUI wizard.
			m := tui.New()
			p := tea.NewProgram(m, tea.WithAltScreen())
			result, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI error: %w", err)
			}

			final, ok := result.(tui.Model)
			if !ok || !final.Done() {
				fmt.Println("Wizard cancelled.")
				return nil
			}

			// Print summary before generation.
			fmt.Println(tui.RenderSummary(final.State()))

			cfg = wizard.BuildConfig(final.State())
		}

		return runGeneration(cfg)
	},
}

func init() {
	initCmd.Flags().StringVar(&fromFile, "from", "", "Load configuration from a lazygo.yml file")
}

// ---- validate command ------------------------------------------------------

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a lazygo.yml configuration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		cfg, err := config.LoadFromYAML(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ Invalid configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Valid configuration: %s (%s/%s)\n",
			cfg.Name, cfg.Type, cfg.Criticality)
		return nil
	},
}

// ---- version command -------------------------------------------------------

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("lazy.go v%s\n", version)
	},
}

// ---- Generation Pipeline ---------------------------------------------------

func runGeneration(cfg *config.ProjectConfig) error {
	// Apply security defaults for production/critical projects.
	security.EnforceSecurity(cfg)

	// Determine output directory.
	outDir := filepath.Join(".", cfg.Name)

	fmt.Printf("\n⟳ Generating %s in ./%s ...\n\n", cfg.Name, cfg.Name)

	// Scaffold the project.
	gen := scaffold.New(cfg, outDir)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Write the license file explicitly (full text via license package).
	licenseText := scaffold.GenerateLicense(cfg.License, cfg.Author, 2026)
	if cfg.License != config.LicenseProprietary {
		if err := os.WriteFile(filepath.Join(outDir, "LICENSE"), []byte(licenseText), 0o644); err != nil {
			return fmt.Errorf("writing LICENSE: %w", err)
		}
	}

	// Export lazygo.yml.
	yamlPath := filepath.Join(outDir, "lazygo.yml")
	if err := config.ExportToYAML(cfg, yamlPath); err != nil {
		return fmt.Errorf("exporting config: %w", err)
	}

	printTree(outDir)

	fmt.Printf("\n✓ lazygo.yml exported to %s\n", yamlPath)
	fmt.Printf("✓ Project ready at %s\n\n", outDir)

	// GitHub integration.
	if cfg.GitHub.Enabled {
		fmt.Println("⟳ Creating GitHub repository...")
		opts := ghpkg.OptionsFromConfig(cfg, outDir)
		ctx := context.Background()
		if err := ghpkg.CreateRepository(ctx, opts); err != nil {
			fmt.Fprintf(os.Stderr, "⚠ GitHub integration failed: %v\n", err)
			fmt.Fprintln(os.Stderr, "  The project was generated locally. You can push manually.")
		} else {
			fmt.Println("✓ Repository created and pushed to GitHub.")
		}
	}

	fmt.Printf("🎉 Done! Start building:\n\n  cd %s && make build\n\n", cfg.Name)
	return nil
}

// printTree prints a simplified directory tree for the generated project.
func printTree(root string) {
	fmt.Printf("\nGenerated structure:\n\n")
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			fmt.Printf("  %s/\n", filepath.Base(root))
			return nil
		}
		depth := strings.Count(rel, string(os.PathSeparator))
		indent := indentStr(depth)
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		fmt.Printf("  %s%s\n", indent, name)
		return nil
	})
	if err != nil {
		_ = err // non-fatal, tree is cosmetic
	}
}

func indentStr(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "│   "
	}
	return s + "├── "
}
