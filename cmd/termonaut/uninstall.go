package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "🗑️  Complete uninstallation of Termonaut",
	Long: `Completely remove Termonaut from your system.

This command will help you:
• Remove shell integration from your .bashrc/.zshrc
• Delete configuration files and directories
• Clean up data files and databases
• Remove aliases and symbolic links
• Clear cache and temporary files

Options:
  --config     Also remove configuration files
  --data       Also remove data files and databases
  --shell      Remove shell integration
  --all        Complete removal (equivalent to --config --data --shell)
  --dry-run    Show what would be removed without actually removing
  --force      Skip confirmation prompts

Use with caution! This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUninstall(cmd, args)
	},
}

var (
	uninstallConfig  bool
	uninstallData    bool
	uninstallShell   bool
	uninstallAll     bool
	uninstallDryRun  bool
	uninstallForce   bool
)

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallConfig, "config", false, "Remove configuration files")
	uninstallCmd.Flags().BoolVar(&uninstallData, "data", false, "Remove data files and databases")
	uninstallCmd.Flags().BoolVar(&uninstallShell, "shell", false, "Remove shell integration")
	uninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "Complete removal")
	uninstallCmd.Flags().BoolVar(&uninstallDryRun, "dry-run", false, "Show what would be removed")
	uninstallCmd.Flags().BoolVar(&uninstallForce, "force", false, "Skip confirmation prompts")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	fmt.Println("🗑️  Termonaut Uninstaller")
	fmt.Println("========================")
	fmt.Println()

	// If --all is specified, enable all removal options
	if uninstallAll {
		uninstallConfig = true
		uninstallData = true
		uninstallShell = true
	}

	// If no specific options, ask user what to remove
	if !uninstallConfig && !uninstallData && !uninstallShell {
		if err := askUninstallOptions(); err != nil {
			return err
		}
	}

	// Show what will be removed
	itemsToRemove := getUninstallItems()
	if len(itemsToRemove) == 0 {
		fmt.Println("ℹ️  Nothing selected for removal.")
		return nil
	}

	fmt.Println("📋 Items to be removed:")
	for _, item := range itemsToRemove {
		if uninstallDryRun {
			fmt.Printf("  🔍 [DRY RUN] %s\n", item)
		} else {
			fmt.Printf("  ❌ %s\n", item)
		}
	}
	fmt.Println()

	if uninstallDryRun {
		fmt.Println("✅ Dry run completed. No files were actually removed.")
		return nil
	}

	// Confirmation
	if !uninstallForce {
		fmt.Print("⚠️  Are you sure you want to proceed? This action cannot be undone! (y/N): ")
		if !askYesNo(false) {
			fmt.Println("❌ Uninstallation cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Perform uninstallation
	return performUninstall()
}

func askUninstallOptions() error {
	fmt.Println("What would you like to remove?")
	fmt.Println()

	fmt.Print("Remove shell integration (.bashrc/.zshrc hooks)? (y/N): ")
	uninstallShell = askYesNo(false)

	fmt.Print("Remove configuration files (~/.termonaut/config.toml)? (y/N): ")
	uninstallConfig = askYesNo(false)

	fmt.Print("Remove data files and databases (~/.termonaut/termonaut.db)? (y/N): ")
	uninstallData = askYesNo(false)

	fmt.Println()
	return nil
}

func getUninstallItems() []string {
	var items []string

	if uninstallShell {
		items = append(items, "Shell integration (bashrc/zshrc hooks)")
		items = append(items, "Termonaut alias ('tn' command)")
	}

	if uninstallConfig {
		configDir := config.GetConfigDir()
		items = append(items, fmt.Sprintf("Configuration directory: %s", configDir))
		items = append(items, "Configuration file: config.toml")
	}

	if uninstallData {
		configDir := config.GetConfigDir()
		items = append(items, fmt.Sprintf("Database file: %s/termonaut.db", configDir))
		items = append(items, fmt.Sprintf("Cache directory: %s/cache/", configDir))
		items = append(items, fmt.Sprintf("Log files: %s/*.log", configDir))
	}

	return items
}

func performUninstall() error {
	var errors []string

	fmt.Println("🔄 Starting uninstallation...")
	fmt.Println()

	// Remove shell integration
	if uninstallShell {
		fmt.Print("🔧 Removing shell integration... ")
		if err := removeShellIntegration(); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			errors = append(errors, fmt.Sprintf("Shell integration: %v", err))
		} else {
			fmt.Println("✅ Done")
		}

		fmt.Print("🔧 Removing termonaut alias... ")
		if err := removeAlias(); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			errors = append(errors, fmt.Sprintf("Alias removal: %v", err))
		} else {
			fmt.Println("✅ Done")
		}
	}

	// Remove configuration files
	if uninstallConfig {
		fmt.Print("📁 Removing configuration files... ")
		if err := removeConfigFiles(); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			errors = append(errors, fmt.Sprintf("Configuration files: %v", err))
		} else {
			fmt.Println("✅ Done")
		}
	}

	// Remove data files
	if uninstallData {
		fmt.Print("🗄️  Removing data files... ")
		if err := removeDataFiles(); err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			errors = append(errors, fmt.Sprintf("Data files: %v", err))
		} else {
			fmt.Println("✅ Done")
		}
	}

	fmt.Println()

	if len(errors) > 0 {
		fmt.Println("⚠️  Uninstallation completed with some errors:")
		for _, err := range errors {
			fmt.Printf("  • %s\n", err)
		}
		fmt.Println()
		fmt.Println("You may need to manually remove some items or run with sudo.")
		return fmt.Errorf("uninstallation completed with %d errors", len(errors))
	}

	fmt.Println("🎉 Uninstallation completed successfully!")
	fmt.Println()
	fmt.Println("Thank you for using Termonaut! 🚀")
	fmt.Println("If you decide to reinstall, you can download it from:")
	fmt.Println("https://github.com/oiahoon/termonaut")

	return nil
}

func removeShellIntegration() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Remove from common shell config files
	shellFiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".profile"),
	}

	for _, shellFile := range shellFiles {
		if err := removeTermonautFromShellFile(shellFile); err != nil {
			// Don't fail if file doesn't exist
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to clean %s: %w", shellFile, err)
			}
		}
	}

	return nil
}

func removeTermonautFromShellFile(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to remove
	}

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	skipNext := false

	for i, line := range lines {
		// Skip lines that contain termonaut hooks
		if strings.Contains(line, "termonaut") && 
		   (strings.Contains(line, "preexec") || 
		    strings.Contains(line, "precmd") ||
		    strings.Contains(line, "PROMPT_COMMAND")) {
			continue
		}

		// Skip termonaut comment blocks
		if strings.Contains(line, "# Termonaut") {
			skipNext = true
			continue
		}

		if skipNext && strings.TrimSpace(line) == "" {
			skipNext = false
			continue
		}

		if skipNext {
			continue
		}

		newLines = append(newLines, line)
	}

	// Write back the cleaned content
	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(filename, []byte(newContent), 0644)
}

func removeAlias() error {
	// Try to remove the 'tn' alias if it exists
	return removeAliasCommand()
}

func removeAliasCommand() error {
	// This would call the existing alias remove functionality
	// For now, we'll implement a simple version
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Check common locations for the 'tn' alias
	possiblePaths := []string{
		"/usr/local/bin/tn",
		"/usr/bin/tn",
		filepath.Join(homeDir, ".local/bin/tn"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove alias at %s: %w", path, err)
			}
		}
	}

	return nil
}

func removeConfigFiles() error {
	configDir := config.GetConfigDir()
	
	// Only remove config.toml, not the entire directory
	configFile := filepath.Join(configDir, "config.toml")
	if _, err := os.Stat(configFile); err == nil {
		if err := os.Remove(configFile); err != nil {
			return fmt.Errorf("failed to remove config file: %w", err)
		}
	}

	return nil
}

func removeDataFiles() error {
	configDir := config.GetConfigDir()

	// Remove database file
	dbFile := filepath.Join(configDir, "termonaut.db")
	if _, err := os.Stat(dbFile); err == nil {
		if err := os.Remove(dbFile); err != nil {
			return fmt.Errorf("failed to remove database: %w", err)
		}
	}

	// Remove cache directory
	cacheDir := filepath.Join(configDir, "cache")
	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.RemoveAll(cacheDir); err != nil {
			return fmt.Errorf("failed to remove cache directory: %w", err)
		}
	}

	// Remove log files
	logPattern := filepath.Join(configDir, "*.log")
	matches, err := filepath.Glob(logPattern)
	if err == nil {
		for _, match := range matches {
			if err := os.Remove(match); err != nil {
				return fmt.Errorf("failed to remove log file %s: %w", match, err)
			}
		}
	}

	// If directory is empty after removing files, remove it
	if isEmpty, _ := isDirEmpty(configDir); isEmpty {
		os.Remove(configDir)
	}

	return nil
}

func isDirEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		return false, nil // Directory is not empty
	}
	return true, nil // Directory is empty
}

func askYesNo(defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}
