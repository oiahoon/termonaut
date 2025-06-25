package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/oiahoon/termonaut/internal/display"
)

var floatingTestCmd = &cobra.Command{
	Use:   "floating-test",
	Short: "🎭 Test floating easter egg notifications",
	Long: `Test the floating notification system for easter eggs.

This command demonstrates how easter eggs could appear as floating
notifications at the top of your terminal, similar to desktop notifications.

The system automatically detects your terminal type and uses the best
available method to display notifications.`,
	Run: runFloatingTest,
}

func runFloatingTest(cmd *cobra.Command, args []string) {
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
}

func init() {
	rootCmd.AddCommand(floatingTestCmd)
}
