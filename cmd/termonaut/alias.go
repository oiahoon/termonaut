package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "🔗 Manage 'tn' command alias",
	Long: `Manage the 'tn' shortcut alias for termonaut.

This command helps you create, check, or remove the 'tn' symbolic link
that allows you to use 'tn' instead of 'termonaut'.

Subcommands:
  create   Create the 'tn' alias
  check    Check if 'tn' alias exists
  remove   Remove the 'tn' alias
  info     Show alias information`,
}

var aliasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create 'tn' alias",
	Long: `Create a symbolic link 'tn' that points to termonaut.
This allows you to use 'tn' as a shortcut for 'termonaut'.

The command will try different locations in this order:
1. ~/.local/bin (recommended, no sudo needed)
2. /usr/local/bin (requires sudo)
3. Other directories in your PATH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return createAlias(cmd, args)
	},
}

var aliasCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check 'tn' alias status",
	Long:  `Check if the 'tn' alias exists and where it points to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkAlias(cmd, args)
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove 'tn' alias",
	Long:  `Remove the 'tn' symbolic link if it exists.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAlias(cmd, args)
	},
}

var aliasInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show alias information",
	Long:  `Show detailed information about the 'tn' alias system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showAliasInfo(cmd, args)
	},
}

func createAlias(cmd *cobra.Command, args []string) error {
	fmt.Println("🔗 Creating 'tn' alias...")
	
	// Get binary path
	binaryPath, err := shell.GetBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get binary path: %w", err)
	}

	// Try different locations
	locations := []struct {
		path        string
		description string
		needsSudo   bool
	}{
		{filepath.Join(os.Getenv("HOME"), ".local/bin"), "User local bin (recommended)", false},
		{"/usr/local/bin", "System local bin", true},
	}

	for _, loc := range locations {
		fmt.Printf("📁 Trying %s (%s)...\n", loc.path, loc.description)
		
		// Create directory if it doesn't exist (for user directories)
		if !loc.needsSudo {
			if err := os.MkdirAll(loc.path, 0755); err != nil {
				fmt.Printf("❌ Cannot create directory: %v\n", err)
				continue
			}
		}

		// Check if directory exists
		if _, err := os.Stat(loc.path); os.IsNotExist(err) {
			fmt.Printf("❌ Directory doesn't exist: %s\n", loc.path)
			continue
		}

		linkPath := filepath.Join(loc.path, "tn")

		// Check if already exists
		if _, err := os.Lstat(linkPath); err == nil {
			fmt.Printf("✅ 'tn' alias already exists at %s\n", linkPath)
			return nil
		}

		// Try to create symlink
		var createErr error
		if loc.needsSudo {
			fmt.Printf("🔐 Creating symlink requires sudo privileges...\n")
			createErr = createSymlinkWithSudo(binaryPath, linkPath)
		} else {
			createErr = os.Symlink(binaryPath, linkPath)
		}

		if createErr == nil {
			fmt.Printf("✅ Created 'tn' alias at %s\n", linkPath)
			
			// Check if location is in PATH
			if !isInPath(loc.path) {
				fmt.Printf("⚠️  Warning: %s is not in your PATH\n", loc.path)
				fmt.Printf("💡 Add this to your shell config (~/.bashrc or ~/.zshrc):\n")
				fmt.Printf("   export PATH=\"%s:$PATH\"\n", loc.path)
			}
			
			return nil
		}

		fmt.Printf("❌ Failed to create symlink: %v\n", createErr)
	}

	return fmt.Errorf("failed to create 'tn' alias in any location")
}

func checkAlias(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking 'tn' alias status...")
	
	// Check if 'tn' command is available
	tnPath, err := exec.LookPath("tn")
	if err != nil {
		fmt.Println("❌ 'tn' alias not found in PATH")
		fmt.Println("💡 Run 'termonaut alias create' to create it")
		return nil
	}

	fmt.Printf("✅ 'tn' alias found at: %s\n", tnPath)

	// Check if it's a symlink
	if linkInfo, err := os.Lstat(tnPath); err == nil {
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(tnPath); err == nil {
				fmt.Printf("🔗 Points to: %s\n", target)
			}
		} else {
			fmt.Printf("📄 Type: Regular file (not a symlink)\n")
		}
	}

	// Test if it works
	fmt.Println("🧪 Testing 'tn' command...")
	if output, err := exec.Command("tn", "--version").Output(); err == nil {
		fmt.Printf("✅ 'tn' command works: %s", string(output))
	} else {
		fmt.Printf("❌ 'tn' command test failed: %v\n", err)
	}

	return nil
}

func removeAlias(cmd *cobra.Command, args []string) error {
	fmt.Println("🗑️  Removing 'tn' alias...")
	
	// Find the alias
	tnPath, err := exec.LookPath("tn")
	if err != nil {
		fmt.Println("✅ 'tn' alias not found (already removed)")
		return nil
	}

	fmt.Printf("📍 Found 'tn' at: %s\n", tnPath)

	// Check if we can remove it
	dir := filepath.Dir(tnPath)
	if !isWritable(dir) {
		fmt.Printf("🔐 Removing alias requires sudo privileges...\n")
		cmd := exec.Command("sudo", "rm", tnPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to remove with sudo: %w", err)
		}
	} else {
		if err := os.Remove(tnPath); err != nil {
			return fmt.Errorf("failed to remove: %w", err)
		}
	}

	fmt.Printf("✅ Removed 'tn' alias from %s\n", tnPath)
	return nil
}

func showAliasInfo(cmd *cobra.Command, args []string) error {
	fmt.Println("ℹ️  'tn' Alias Information")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("The 'tn' alias is a shortcut that allows you to use 'tn' instead of 'termonaut'.")
	fmt.Println()
	fmt.Println("📍 Preferred locations (in order):")
	fmt.Println("1. ~/.local/bin (user-specific, no sudo needed)")
	fmt.Println("2. /usr/local/bin (system-wide, requires sudo)")
	fmt.Println()
	fmt.Println("🔧 Commands:")
	fmt.Println("• termonaut alias create  - Create the alias")
	fmt.Println("• termonaut alias check   - Check alias status")
	fmt.Println("• termonaut alias remove  - Remove the alias")
	fmt.Println()
	fmt.Println("💡 If ~/.local/bin is used, make sure it's in your PATH:")
	fmt.Println("   export PATH=\"$HOME/.local/bin:$PATH\"")
	fmt.Println()
	
	// Show current status
	fmt.Println("📊 Current Status:")
	checkAlias(cmd, args)
	
	return nil
}

// Helper functions
func createSymlinkWithSudo(target, link string) error {
	cmd := exec.Command("sudo", "ln", "-sf", target, link)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func isInPath(dir string) bool {
	pathDirs := strings.Split(os.Getenv("PATH"), ":")
	for _, pathDir := range pathDirs {
		if pathDir == dir {
			return true
		}
	}
	return false
}

func isWritable(dir string) bool {
	testFile := filepath.Join(dir, ".termonaut_write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

func init() {
	// Add subcommands
	aliasCmd.AddCommand(aliasCreateCmd)
	aliasCmd.AddCommand(aliasCheckCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)
	aliasCmd.AddCommand(aliasInfoCmd)
	
	// Add to root command
	rootCmd.AddCommand(aliasCmd)
}
