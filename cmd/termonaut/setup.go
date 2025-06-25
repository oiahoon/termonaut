package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "🚀 Interactive setup wizard for new users",
	Long: `Welcome to Termonaut! This interactive setup wizard will help you:

• Install shell integration for automatic command tracking
• Configure your preferred UI mode and theme
• Set up your avatar and gamification preferences
• Test your terminal capabilities
• Get started with your first commands

Perfect for first-time users who want a guided experience!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetupWizard(cmd, args)
	},
}

func runSetupWizard(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Welcome to Termonaut Setup Wizard!")
	fmt.Println("=====================================")
	fmt.Println()
	
	// Check if already set up
	if isAlreadySetup() {
		fmt.Println("✅ Termonaut appears to be already set up!")
		fmt.Print("Do you want to run the setup wizard anyway? (y/N): ")
		
		if !askYesNo(false) {
			fmt.Println("Setup cancelled. Use 'termonaut tui' to start using Termonaut!")
			return nil
		}
		fmt.Println()
	}

	// Step 1: Welcome and explanation
	if err := showWelcomeStep(); err != nil {
		return err
	}

	// Step 2: Shell integration
	if err := setupShellIntegration(); err != nil {
		return err
	}

	// Step 3: UI preferences
	if err := setupUIPreferences(); err != nil {
		return err
	}

	// Step 4: Avatar preferences
	if err := setupAvatarPreferences(); err != nil {
		return err
	}

	// Step 5: Test and completion
	if err := completeSetup(); err != nil {
		return err
	}

	return nil
}

func isAlreadySetup() bool {
	// Check if shell integration exists
	binaryPath, err := shell.GetBinaryPath()
	if err != nil {
		return false
	}
	
	installer, err := shell.NewHookInstaller(binaryPath)
	if err != nil {
		return false
	}
	
	installed, err := installer.IsInstalled()
	return err == nil && installed
}

func showWelcomeStep() error {
	fmt.Println("📖 What is Termonaut?")
	fmt.Println("─────────────────────")
	fmt.Println("Termonaut is your terminal productivity companion that:")
	fmt.Println("• 📊 Tracks your command usage and productivity")
	fmt.Println("• 🎮 Gamifies your terminal experience with XP and levels")
	fmt.Println("• 🏆 Unlocks achievements as you explore new commands")
	fmt.Println("• 📈 Provides beautiful visualizations of your activity")
	fmt.Println("• 🎨 Features customizable avatars and themes")
	fmt.Println()
	fmt.Println("Let's get you set up! This will take about 2-3 minutes.")
	fmt.Println()
	
	fmt.Print("Ready to continue? (Y/n): ")
	if !askYesNo(true) {
		return fmt.Errorf("setup cancelled by user")
	}
	
	fmt.Println()
	return nil
}

func setupShellIntegration() error {
	fmt.Println("🔧 Step 1: Shell Integration")
	fmt.Println("────────────────────────────")
	fmt.Println("To track your commands, Termonaut needs to integrate with your shell.")
	fmt.Println("This will add a small hook to your ~/.bashrc or ~/.zshrc file.")
	fmt.Println()
	
	fmt.Print("Install shell integration? (Y/n): ")
	if !askYesNo(true) {
		fmt.Println("⚠️  Skipping shell integration. You can install it later with 'termonaut init'")
		fmt.Println()
		return nil
	}

	// Run the init command
	fmt.Println("Installing shell integration...")
	
	binaryPath, err := shell.GetBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get binary path: %w", err)
	}

	installer, err := shell.NewHookInstaller(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create hook installer: %w", err)
	}

	if err := installer.InstallWithForce(false); err != nil {
		return fmt.Errorf("failed to install shell hook: %w", err)
	}

	fmt.Printf("✅ Shell integration installed for %s!\n", installer.GetShellType())
	fmt.Println()
	return nil
}

func setupUIPreferences() error {
	fmt.Println("🎨 Step 2: UI Preferences")
	fmt.Println("─────────────────────────")
	fmt.Println("Termonaut offers different viewing modes:")
	fmt.Println()
	fmt.Println("1. 🧠 Smart Mode (Recommended) - Automatically adapts to your terminal size")
	fmt.Println("2. 📱 Compact Mode - Small avatars, efficient layout")
	fmt.Println("3. 🖥️  Full Mode - Large avatars, immersive experience")
	fmt.Println("4. 📝 Minimal Mode - Text-only stats output")
	fmt.Println()
	
	fmt.Print("Choose your preferred mode (1-4) [1]: ")
	choice := askChoice([]string{"1", "2", "3", "4"}, "1")
	
	var mode string
	var theme string
	
	switch choice {
	case "1":
		mode = "smart"
		theme = "space"
		fmt.Println("✅ Smart mode selected - Termonaut will adapt to your terminal!")
	case "2":
		mode = "compact"
		theme = "minimal"
		fmt.Println("✅ Compact mode selected - Efficient and clean!")
	case "3":
		mode = "full"
		theme = "space"
		fmt.Println("✅ Full mode selected - Maximum visual experience!")
	case "4":
		mode = "minimal"
		theme = "minimal"
		fmt.Println("✅ Minimal mode selected - Fast and lightweight!")
	}
	
	// Save preferences
	if err := saveUIConfig(mode, theme); err != nil {
		fmt.Printf("⚠️  Failed to save UI preferences: %v\n", err)
	} else {
		fmt.Println("💾 UI preferences saved!")
	}
	
	fmt.Println()
	return nil
}

func setupAvatarPreferences() error {
	fmt.Println("👤 Step 3: Avatar Preferences")
	fmt.Println("─────────────────────────────")
	fmt.Println("Termonaut can display personalized avatars that evolve with your level!")
	fmt.Println()
	
	fmt.Print("Enable avatar display? (Y/n): ")
	showAvatar := askYesNo(true)
	
	var avatarStyle string
	if showAvatar {
		fmt.Println()
		fmt.Println("Choose your avatar style:")
		fmt.Println("1. 🎮 Pixel Art (Recommended) - Retro gaming style")
		fmt.Println("2. 🤖 Bottts - Robot/bot style")
		fmt.Println("3. 🧑 Adventurer - Human adventurer style")
		fmt.Println("4. 😊 Avataaars - Cartoon style")
		fmt.Println()
		
		fmt.Print("Choose avatar style (1-4) [1]: ")
		styleChoice := askChoice([]string{"1", "2", "3", "4"}, "1")
		
		switch styleChoice {
		case "1":
			avatarStyle = "pixel-art"
			fmt.Println("✅ Pixel art style selected - Retro gaming vibes!")
		case "2":
			avatarStyle = "bottts"
			fmt.Println("✅ Bottts style selected - Beep boop robot mode!")
		case "3":
			avatarStyle = "adventurer"
			fmt.Println("✅ Adventurer style selected - Ready for quests!")
		case "4":
			avatarStyle = "avataaars"
			fmt.Println("✅ Avataaars style selected - Cartoon fun!")
		}
	} else {
		avatarStyle = "pixel-art" // Default even if disabled
		fmt.Println("✅ Avatar display disabled - You can enable it later in settings!")
	}
	
	// Save avatar preferences
	if err := saveAvatarConfig(showAvatar, avatarStyle); err != nil {
		fmt.Printf("⚠️  Failed to save avatar preferences: %v\n", err)
	} else {
		fmt.Println("💾 Avatar preferences saved!")
	}
	
	fmt.Println()
	return nil
}

func completeSetup() error {
	fmt.Println("🎉 Step 4: Setup Complete!")
	fmt.Println("──────────────────────────")
	fmt.Println("Congratulations! Termonaut is now set up and ready to use.")
	fmt.Println()
	
	fmt.Println("🚀 Next Steps:")
	fmt.Println("1. Restart your terminal or run: source ~/.bashrc (or ~/.zshrc)")
	fmt.Println("2. Use your terminal normally - commands will be tracked automatically")
	fmt.Println("3. Check your progress with: termonaut tui")
	fmt.Println("4. View quick stats with: termonaut stats")
	fmt.Println()
	
	fmt.Print("Would you like to test Termonaut now? (Y/n): ")
	if askYesNo(true) {
		fmt.Println()
		fmt.Println("🧪 Testing Termonaut...")
		
		// Test basic functionality
		if err := testTeremonaut(); err != nil {
			fmt.Printf("⚠️  Test failed: %v\n", err)
			fmt.Println("Don't worry, you can still use Termonaut. Try 'termonaut tui' manually.")
		} else {
			fmt.Println("✅ Test successful! Termonaut is working correctly.")
		}
	}
	
	fmt.Println()
	fmt.Println("🎊 Welcome to the Termonaut community!")
	fmt.Println("Happy terminal exploring! 🚀")
	
	return nil
}

// Helper functions
func askYesNo(defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	
	if input == "" {
		return defaultYes
	}
	
	return input == "y" || input == "yes"
}

func askChoice(options []string, defaultChoice string) string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultChoice
	}
	
	for _, option := range options {
		if input == option {
			return input
		}
	}
	
	return defaultChoice
}

func saveUIConfig(mode, theme string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	
	cfg.UI.DefaultMode = mode
	cfg.UI.Theme = theme
	
	return config.Save(cfg)
}

func saveAvatarConfig(showAvatar bool, avatarStyle string) error {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	
	cfg.UI.ShowAvatar = showAvatar
	cfg.UI.AvatarStyle = avatarStyle
	
	return config.Save(cfg)
}

func testTeremonaut() error {
	// Simple test - try to run stats command
	return runStatsCommand(nil, []string{})
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
