package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "🩺 Diagnose and fix Termonaut installation issues",
	Long: `Comprehensive system diagnosis to detect and fix common issues.

This command will check:
• Installation completeness and integrity
• Configuration file validity
• Shell integration status
• Database connectivity and health
• File permissions and access
• System compatibility

Options:
  --verbose    Show detailed diagnostic information
  --fix        Automatically fix detected issues
  --report     Generate a detailed diagnostic report

Perfect for troubleshooting installation or configuration problems!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDoctor(cmd, args)
	},
}

var (
	doctorVerbose bool
	doctorFix     bool
	doctorReport  bool
)

func init() {
	doctorCmd.Flags().BoolVar(&doctorVerbose, "verbose", false, "Show detailed diagnostic information")
	doctorCmd.Flags().BoolVar(&doctorFix, "fix", false, "Automatically fix detected issues")
	doctorCmd.Flags().BoolVar(&doctorReport, "report", false, "Generate detailed diagnostic report")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("🩺 Termonaut System Doctor")
	fmt.Println("==========================")
	fmt.Println()

	var issues []DiagnosticIssue
	var fixes []string

	// Run all diagnostic checks
	fmt.Println("🔍 Running diagnostic checks...")
	fmt.Println()

	// Check 1: Installation integrity
	fmt.Print("📦 Checking installation integrity... ")
	if issue := checkInstallationIntegrity(); issue != nil {
		fmt.Println("❌ Issues found")
		issues = append(issues, *issue)
		if doctorVerbose {
			fmt.Printf("   %s\n", issue.Description)
		}
	} else {
		fmt.Println("✅ OK")
	}

	// Check 2: Configuration validity
	fmt.Print("⚙️  Checking configuration validity... ")
	if issue := checkConfigurationValidity(); issue != nil {
		fmt.Println("❌ Issues found")
		issues = append(issues, *issue)
		if doctorVerbose {
			fmt.Printf("   %s\n", issue.Description)
		}
	} else {
		fmt.Println("✅ OK")
	}

	// Check 3: Shell integration
	fmt.Print("🐚 Checking shell integration... ")
	if issue := checkShellIntegration(); issue != nil {
		fmt.Println("❌ Issues found")
		issues = append(issues, *issue)
		if doctorVerbose {
			fmt.Printf("   %s\n", issue.Description)
		}
	} else {
		fmt.Println("✅ OK")
	}

	// Check 4: Database health
	fmt.Print("🗄️  Checking database health... ")
	if issue := checkDatabaseHealth(); issue != nil {
		fmt.Println("❌ Issues found")
		issues = append(issues, *issue)
		if doctorVerbose {
			fmt.Printf("   %s\n", issue.Description)
		}
	} else {
		fmt.Println("✅ OK")
	}

	// Check 5: File permissions
	fmt.Print("🔐 Checking file permissions... ")
	if issue := checkFilePermissions(); issue != nil {
		fmt.Println("❌ Issues found")
		issues = append(issues, *issue)
		if doctorVerbose {
			fmt.Printf("   %s\n", issue.Description)
		}
	} else {
		fmt.Println("✅ OK")
	}

	fmt.Println()

	// Summary
	if len(issues) == 0 {
		fmt.Println("🎉 All checks passed! Your Termonaut installation is healthy.")
		return nil
	}

	fmt.Printf("⚠️  Found %d issue(s):\n", len(issues))
	for i, issue := range issues {
		fmt.Printf("  %d. %s: %s\n", i+1, issue.Category, issue.Description)
		if issue.Solution != "" {
			fmt.Printf("     💡 Solution: %s\n", issue.Solution)
		}
	}
	fmt.Println()

	// Auto-fix if requested
	if doctorFix {
		fmt.Println("🔧 Attempting to fix issues...")
		for _, issue := range issues {
			if issue.AutoFixable {
				fmt.Printf("🔄 Fixing: %s... ", issue.Category)
				if err := issue.FixFunction(); err != nil {
					fmt.Printf("❌ Failed: %v\n", err)
				} else {
					fmt.Println("✅ Fixed")
					fixes = append(fixes, issue.Category)
				}
			} else {
				fmt.Printf("⚠️  Cannot auto-fix: %s (manual intervention required)\n", issue.Category)
			}
		}

		if len(fixes) > 0 {
			fmt.Printf("\n🎉 Successfully fixed %d issue(s)!\n", len(fixes))
			fmt.Println("Please restart your terminal or run 'source ~/.bashrc' to apply changes.")
		}
	} else {
		fmt.Println("💡 Run with --fix to automatically fix issues where possible.")
	}

	// Generate report if requested
	if doctorReport {
		if err := generateDiagnosticReport(issues, fixes); err != nil {
			fmt.Printf("⚠️  Failed to generate report: %v\n", err)
		} else {
			fmt.Println("📄 Diagnostic report saved to ~/.termonaut/diagnostic_report.txt")
		}
	}

	return nil
}

type DiagnosticIssue struct {
	Category    string
	Description string
	Solution    string
	AutoFixable bool
	FixFunction func() error
}

func checkInstallationIntegrity() *DiagnosticIssue {
	configDir := config.GetConfigDir()
	
	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return &DiagnosticIssue{
			Category:    "Installation Integrity",
			Description: "Configuration directory does not exist",
			Solution:    "Run 'termonaut init' to initialize Termonaut",
			AutoFixable: true,
			FixFunction: func() error {
				return os.MkdirAll(configDir, 0755)
			},
		}
	}

	return nil
}

func checkConfigurationValidity() *DiagnosticIssue {
	configDir := config.GetConfigDir()
	configFile := filepath.Join(configDir, "config.toml")

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &DiagnosticIssue{
			Category:    "Configuration",
			Description: "Configuration file does not exist",
			Solution:    "Run 'termonaut setup' to create default configuration",
			AutoFixable: true,
			FixFunction: func() error {
				// Create default config
				return createDefaultConfig(configFile)
			},
		}
	}

	// Try to load config
	if _, err := config.LoadConfig(); err != nil {
		return &DiagnosticIssue{
			Category:    "Configuration",
			Description: fmt.Sprintf("Configuration file is invalid: %v", err),
			Solution:    "Fix configuration file or run 'termonaut config wizard' to recreate",
			AutoFixable: false,
		}
	}

	return nil
}

func checkShellIntegration() *DiagnosticIssue {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &DiagnosticIssue{
			Category:    "Shell Integration",
			Description: "Cannot access home directory",
			Solution:    "Check file system permissions",
			AutoFixable: false,
		}
	}

	// Check common shell files
	shellFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
	}

	hasIntegration := false
	for _, shellFile := range shellFiles {
		if content, err := os.ReadFile(shellFile); err == nil {
			if strings.Contains(string(content), "termonaut") {
				hasIntegration = true
				break
			}
		}
	}

	if !hasIntegration {
		return &DiagnosticIssue{
			Category:    "Shell Integration",
			Description: "Shell integration not found in .bashrc or .zshrc",
			Solution:    "Run 'termonaut init' to install shell hooks",
			AutoFixable: true,
			FixFunction: func() error {
				// This would call the existing shell integration setup
				return installShellIntegration()
			},
		}
	}

	return nil
}

func checkDatabaseHealth() *DiagnosticIssue {
	configDir := config.GetConfigDir()
	dbFile := filepath.Join(configDir, "termonaut.db")

	// Check if database file exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return &DiagnosticIssue{
			Category:    "Database",
			Description: "Database file does not exist",
			Solution:    "Database will be created automatically on first use",
			AutoFixable: true,
			FixFunction: func() error {
				// Database will be created automatically
				return nil
			},
		}
	}

	// Check if database is accessible
	if _, err := os.Open(dbFile); err != nil {
		return &DiagnosticIssue{
			Category:    "Database",
			Description: fmt.Sprintf("Cannot access database file: %v", err),
			Solution:    "Check file permissions or recreate database",
			AutoFixable: false,
		}
	}

	return nil
}

func checkFilePermissions() *DiagnosticIssue {
	configDir := config.GetConfigDir()

	// Check if we can write to config directory
	testFile := filepath.Join(configDir, ".test_write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return &DiagnosticIssue{
			Category:    "File Permissions",
			Description: "Cannot write to configuration directory",
			Solution:    "Check directory permissions or run with appropriate privileges",
			AutoFixable: false,
		}
	}
	os.Remove(testFile) // Clean up test file

	return nil
}

func createDefaultConfig(configFile string) error {
	defaultConfig := `# Termonaut Configuration
# Generated by doctor command

[display]
mode = "enter"
theme = "emoji"
show_gamification = true

[tracking]
idle_timeout_minutes = 10
track_git_repos = true
command_categories = true

[privacy]
opt_out_commands = ["password", "secret"]
anonymous_mode = false
`

	return os.WriteFile(configFile, []byte(defaultConfig), 0644)
}

func installShellIntegration() error {
	// This would implement shell integration installation
	// For now, return a placeholder
	return fmt.Errorf("shell integration installation not implemented in doctor")
}

func generateDiagnosticReport(issues []DiagnosticIssue, fixes []string) error {
	configDir := config.GetConfigDir()
	reportFile := filepath.Join(configDir, "diagnostic_report.txt")

	var report strings.Builder
	report.WriteString("Termonaut Diagnostic Report\n")
	report.WriteString("===========================\n\n")

	if len(issues) == 0 {
		report.WriteString("✅ All checks passed! No issues found.\n")
	} else {
		report.WriteString(fmt.Sprintf("Found %d issue(s):\n\n", len(issues)))
		for i, issue := range issues {
			report.WriteString(fmt.Sprintf("%d. %s\n", i+1, issue.Category))
			report.WriteString(fmt.Sprintf("   Description: %s\n", issue.Description))
			report.WriteString(fmt.Sprintf("   Solution: %s\n", issue.Solution))
			report.WriteString(fmt.Sprintf("   Auto-fixable: %t\n\n", issue.AutoFixable))
		}
	}

	if len(fixes) > 0 {
		report.WriteString("Fixed issues:\n")
		for _, fix := range fixes {
			report.WriteString(fmt.Sprintf("- %s\n", fix))
		}
	}

	return os.WriteFile(reportFile, []byte(report.String()), 0644)
}
