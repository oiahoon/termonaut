package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/display"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/spf13/cobra"
)

var easterEggCmd = &cobra.Command{
	Use:   "easter-egg",
	Short: "Test easter egg system or show motivational quote",
	Long: `Test the easter egg system with sample data or display a random 
motivational quote to brighten your terminal experience.

Flags:
  --floating    Test old floating notification system (deprecated)
  --smart       Test new smart notification system (recommended)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEasterEggCommand(cmd, args)
	},
}

var floatingTest bool
var smartTest bool

func init() {
	easterEggCmd.Flags().BoolVar(&floatingTest, "floating", false, "Test old floating notification system (deprecated)")
	easterEggCmd.Flags().BoolVar(&smartTest, "smart", false, "Test new smart notification system")
}

func runEasterEggCommand(cmd *cobra.Command, args []string) error {
	// Check if smart test is requested
	if smartTest {
		return runSmartNotificationTest()
	}
	
	// Check if floating test is requested (deprecated)
	if floatingTest {
		fmt.Println("⚠️  Warning: Floating notification system is deprecated due to terminal interference issues.")
		fmt.Println("   Please use --smart flag to test the new smart notification system.")
		fmt.Println()
		return runFloatingNotificationTest()
	}
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if easter eggs are enabled
	if !cfg.EasterEggsEnabled {
		fmt.Println("🚫 Easter eggs are disabled in configuration.")
		fmt.Println("Enable them with: termonaut config set easter_eggs_enabled true")
		return nil
	}

	// Get flags
	testMode, _ := cmd.Flags().GetBool("test")
	motivational, _ := cmd.Flags().GetBool("motivational")

	easterEggManager := gamification.NewEasterEggManager()

	if motivational {
		// Just show a motivational quote
		quote := easterEggManager.GetRandomMotivationalQuote()
		fmt.Printf("💫 %s\n", quote)
		return nil
	}

	if testMode {
		return runEasterEggTests(easterEggManager)
	}

	// Show easter egg status and recent triggers
	return showEasterEggStatus(cfg)
}

func runEasterEggTests(manager *gamification.EasterEggManager) error {
	fmt.Println("🧪 Testing Easter Egg System")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Test scenarios
	testScenarios := []struct {
		name        string
		context     *gamification.EasterEggContext
		description string
	}{
		{
			name: "Speed Run",
			context: &gamification.EasterEggContext{
				CommandsInTimeframe: 15,
				TimeframeDuration:   3 * time.Second,
				LastCommand:        "ls -la && cd test && git status",
			},
			description: "Rapid command execution",
		},
		{
			name: "Coffee Break Return",
			context: &gamification.EasterEggContext{
				IdleDuration: 10 * time.Minute,
				LastCommand:  "git pull",
			},
			description: "Long idle period",
		},
		{
			name: "Morning Greeting",
			context: &gamification.EasterEggContext{
				IsFirstCommandToday: true,
				LastCommand:        "ls",
			},
			description: "First command of the day",
		},
		{
			name: "Git Force Push",
			context: &gamification.EasterEggContext{
				LastCommand: "git push origin main --force",
			},
			description: "Dangerous git operation",
		},
		{
			name: "Exit Command",
			context: &gamification.EasterEggContext{
				LastCommand: "exit",
			},
			description: "Leaving the terminal",
		},
		{
			name: "Secret Command",
			context: &gamification.EasterEggContext{
				LastCommand: "sudo make me a sandwich",
			},
			description: "Hidden easter egg command",
		},
		{
			name: "4:20 Time",
			context: &gamification.EasterEggContext{
				LastCommand: "date",
			},
			description: "Special time trigger (simulated)",
		},
		{
			name: "Productivity Beast",
			context: &gamification.EasterEggContext{
				CommandsInTimeframe: 55,
				TimeframeDuration:   25 * time.Minute,
				LastCommand:        "git commit -m 'final push'",
			},
			description: "High productivity period",
		},
	}

	for i, scenario := range testScenarios {
		fmt.Printf("%d. %s - %s\n", i+1, scenario.name, scenario.description)
		
		// Try multiple times to account for probability
		triggered := false
		var message string
		
		for attempts := 0; attempts < 10 && !triggered; attempts++ {
			message = manager.CheckForEasterEgg(scenario.context)
			if message != "" {
				triggered = true
			}
		}
		
		if triggered {
			fmt.Printf("   ✅ %s\n", formatEasterEggOutput(message))
		} else {
			fmt.Printf("   ❌ No easter egg triggered (probability-based)\n")
		}
		fmt.Println()
	}

	fmt.Println("🎲 Random Motivational Quote:")
	quote := manager.GetRandomMotivationalQuote()
	fmt.Printf("   💫 %s\n", quote)
	
	return nil
}

func showEasterEggStatus(cfg *config.Config) error {
	fmt.Println("🥚 Easter Egg System Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Status: %s\n", enabledStatus(cfg.EasterEggsEnabled))
	fmt.Printf("Gamification: %s\n", enabledStatus(cfg.ShowGamification))
	fmt.Println()

	if cfg.EasterEggsEnabled {
		fmt.Println("📝 Available Easter Egg Triggers:")
		triggers := []string{
			"🏃‍♂️ Speed Run - Execute many commands quickly",
			"☕ Coffee Break - Return after long idle",
			"🌅 Morning Greeting - First command of the day",
			"⚠️ Git Force Push - Dangerous git operations",
			"🚪 Exit Command - Leaving the terminal",
			"🎮 Secret Commands - Hidden easter eggs",
			"🌿 Special Times - Time-based triggers",
			"🔥 Productivity Beast - High activity periods",
			"💥 Error Recovery - Consecutive command failures",
			"🌙 Night Coding - Late night activities",
		}

		for _, trigger := range triggers {
			fmt.Printf("  • %s\n", trigger)
		}
		
		fmt.Println()
		fmt.Println("🎯 Try these commands to trigger easter eggs:")
		fmt.Println("  • exit")
		fmt.Println("  • sudo make me a sandwich")
		fmt.Println("  • git push --force")
		fmt.Println("  • Execute commands quickly in succession")
		fmt.Println("  • Come back after a coffee break!")
	}

	return nil
}

func formatEasterEggOutput(message string) string {
	// Handle multi-line ASCII art
	if strings.Contains(message, "\n") {
		lines := strings.Split(message, "\n")
		formatted := lines[0] + "\n"
		for i, line := range lines[1:] {
			if i == 0 {
				formatted += "   " + line + "\n"
			} else {
				formatted += "   " + line
				if i < len(lines)-2 {
					formatted += "\n"
				}
			}
		}
		return formatted
	}
	return message
}

func enabledStatus(enabled bool) string {
	if enabled {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

// runSmartNotificationTest runs the smart notification test
func runSmartNotificationTest() error {
	fmt.Println("🔔 Smart Notification System Test")
	fmt.Println("=================================")
	fmt.Println()
	
	notifier := display.NewSmartNotifier()
	
	// Show available methods
	available := notifier.GetAvailableMethods()
	fmt.Printf("📋 Available notification methods (%d):\n", len(available))
	for i, method := range available {
		fmt.Printf("   %d. %s\n", i+1, notifier.MethodName(method))
	}
	fmt.Println()
	
	// Test messages
	testMessages := []string{
		"🚀 Welcome to Termonaut Space Program!",
		"☕ Coffee break detected! Caffeine levels optimal!",
		"🎮 Achievement Unlocked: Terminal Ninja!",
		"🦆 Rubber duck debugging mode activated!",
		"🌙 Late night coding session detected!",
		"🎉 Productivity celebration! You're on fire!",
	}
	
	fmt.Println("🧪 Testing smart notification selection...")
	fmt.Println("(Check for system notifications, terminal title changes, or bell sounds)")
	fmt.Println()
	
	for i, message := range testMessages {
		fmt.Printf("📱 Notification %d/%d: %s\n", i+1, len(testMessages), message)
		
		// Use smart notification
		if err := notifier.ShowEasterEgg(message); err != nil {
			fmt.Printf("   ❌ Failed: %v\n", err)
		} else {
			fmt.Printf("   ✅ Sent successfully\n")
		}
		
		// Wait between notifications
		if i < len(testMessages)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	
	fmt.Println()
	fmt.Println("✅ Smart notification test complete!")
	fmt.Println()
	fmt.Println("💡 How it works:")
	fmt.Println("   • Automatically detects your system capabilities")
	fmt.Println("   • Prefers system notifications (best user experience)")
	fmt.Println("   • Falls back to terminal title or bell notifications")
	fmt.Println("   • Never interferes with your terminal content")
	fmt.Println("   • Safe for all environments and use cases")
	fmt.Println()
	fmt.Println("🎯 This replaces the old floating notification system!")
	
	return nil
}

// runFloatingNotificationTest runs the floating notification test
func runFloatingNotificationTest() error {
	fmt.Println("🎭 Floating Easter Egg Notification Test")
	fmt.Println("========================================")
	fmt.Println()
	
	// Create notification manager
	notifier := display.NewFloatingNotifier()
	
	fmt.Println("📱 Testing different notification styles...")
	fmt.Println("(Watch the top of your terminal for floating notifications)")
	fmt.Println()
	
	// Test sequence
	testMessages := []struct {
		message string
		delay   time.Duration
	}{
		{"🚀 Welcome to Termonaut Space Program!", 1 * time.Second},
		{"☕ Coffee break detected! Caffeine levels optimal!", 4 * time.Second},
		{"🎮 Achievement Unlocked: Terminal Ninja!", 4 * time.Second},
		{"🦆 Rubber duck debugging mode activated!", 4 * time.Second},
		{"🌙 Late night coding session detected!", 4 * time.Second},
		{"🎉 Productivity celebration! You're on fire!", 4 * time.Second},
	}
	
	for i, test := range testMessages {
		fmt.Printf("⏰ Showing notification %d/%d: %s\n", i+1, len(testMessages), test.message)
		
		// Show the floating notification
		notifier.ShowEasterEgg(test.message, 3*time.Second)
		
		// Wait before next notification
		time.Sleep(test.delay)
	}
	
	fmt.Println()
	fmt.Println("✅ Test complete!")
	fmt.Println()
	fmt.Println("💡 How it works:")
	fmt.Println("   • Detects your terminal type automatically")
	fmt.Println("   • Uses modern terminal features when available")
	fmt.Println("   • Falls back to ANSI escape sequences for compatibility")
	fmt.Println("   • Notifications appear at the top and auto-disappear")
	fmt.Println("   • Safe cooldown prevents notification spam")
	fmt.Println()
	fmt.Println("🎯 In real usage, these would appear when you trigger easter eggs!")
	
	return nil
}

func init() {
	easterEggCmd.Flags().BoolP("test", "t", false, "Run easter egg system tests")
	easterEggCmd.Flags().BoolP("motivational", "m", false, "Show a random motivational quote")
	rootCmd.AddCommand(easterEggCmd)
} 