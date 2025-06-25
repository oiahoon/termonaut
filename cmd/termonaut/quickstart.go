package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/spf13/cobra"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "⚡ Quick setup with sensible defaults",
	Long: `Quick setup for users who want to get started immediately.
This will:
• Install shell integration
• Use smart UI mode with space theme
• Enable avatars with pixel-art style
• Set up sensible defaults

For a guided experience, use 'termonaut setup' instead.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runQuickstart(cmd, args)
	},
}

func runQuickstart(cmd *cobra.Command, args []string) error {
	fmt.Println("⚡ Termonaut Quickstart")
	fmt.Println("======================")
	fmt.Println()
	
	// Step 1: Install shell integration
	fmt.Println("🔧 Installing shell integration...")
	if err := quickInstallShell(); err != nil {
		fmt.Printf("❌ Failed to install shell integration: %v\n", err)
		fmt.Println("You can install it manually later with 'termonaut init'")
	} else {
		fmt.Println("✅ Shell integration installed!")
	}
	
	// Step 2: Set up default configuration
	fmt.Println("⚙️  Setting up default configuration...")
	if err := quickSetupConfig(); err != nil {
		fmt.Printf("❌ Failed to setup configuration: %v\n", err)
	} else {
		fmt.Println("✅ Configuration saved!")
	}
	
	// Step 3: Show completion message
	fmt.Println()
	fmt.Println("🎉 Quickstart Complete!")
	fmt.Println("─────────────────────")
	fmt.Println("Termonaut is ready to use with these settings:")
	fmt.Println("• 🧠 Smart UI mode (adapts to terminal size)")
	fmt.Println("• 🎨 Space theme")
	fmt.Println("• 👤 Pixel-art avatars enabled")
	fmt.Println("• 🎮 Gamification enabled")
	fmt.Println()
	fmt.Println("🚀 Next Steps:")
	fmt.Println("1. Restart your terminal or run: source ~/.bashrc (or ~/.zshrc)")
	fmt.Println("2. Start using your terminal normally")
	fmt.Println("3. Check your progress: termonaut tui")
	fmt.Println()
	fmt.Println("💡 Tip: Use 'termonaut setup' for a guided configuration experience!")
	
	return nil
}

func quickInstallShell() error {
	binaryPath, err := shell.GetBinaryPath()
	if err != nil {
		return err
	}

	installer, err := shell.NewHookInstaller(binaryPath)
	if err != nil {
		return err
	}

	// Check if already installed
	installed, err := installer.IsInstalled()
	if err != nil {
		return err
	}
	
	if installed {
		return nil // Already installed, skip
	}

	return installer.InstallWithForce(false)
}

func quickSetupConfig() error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	
	// Set quickstart defaults
	cfg.UI.DefaultMode = "smart"
	cfg.UI.Theme = "space"
	cfg.UI.ShowAvatar = true
	cfg.UI.AvatarStyle = "pixel-art"
	cfg.UI.CompactLayout = false
	cfg.UI.AnimationsEnabled = true
	cfg.ShowGamification = true
	cfg.EasterEggsEnabled = true
	cfg.EmptyCommandStats = true
	
	return config.Save(cfg)
}

func init() {
	rootCmd.AddCommand(quickstartCmd)
}
