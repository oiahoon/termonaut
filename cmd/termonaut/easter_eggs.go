package main

import (
	"fmt"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/spf13/cobra"
)

var easterEggCmd = &cobra.Command{
	Use:   "easter-egg",
	Short: "Test easter egg system or show motivational quote",
	Long: `Test the easter egg system with sample data or display a random 
motivational quote to brighten your terminal experience.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEasterEggCommand(cmd, args)
	},
}

func runEasterEggCommand(cmd *cobra.Command, args []string) error {
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

	// Check for flags
	test, _ := cmd.Flags().GetBool("test")
	motivational, _ := cmd.Flags().GetBool("motivational")

	if motivational {
		return showMotivationalQuote()
	}

	if test {
		return runEasterEggTests()
	}

	// Default behavior - show a random easter egg
	return showRandomEasterEgg()
}

func showMotivationalQuote() error {
	quotes := []string{
		"🚀 \"The best way to predict the future is to create it.\" - Peter Drucker",
		"💡 \"Code is like humor. When you have to explain it, it's bad.\" - Cory House",
		"🎯 \"First, solve the problem. Then, write the code.\" - John Johnson",
		"⚡ \"Experience is the name everyone gives to their mistakes.\" - Oscar Wilde",
		"🌟 \"The only way to learn a new programming language is by writing programs in it.\" - Dennis Ritchie",
		"🔥 \"Talk is cheap. Show me the code.\" - Linus Torvalds",
		"🎨 \"Clean code always looks like it was written by someone who cares.\" - Robert C. Martin",
		"🚀 \"Any fool can write code that a computer can understand. Good programmers write code that humans can understand.\" - Martin Fowler",
	}

	// Select random quote
	quote := quotes[time.Now().UnixNano()%int64(len(quotes))]
	
	fmt.Println("💭 Motivational Quote")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println(quote)
	fmt.Println()
	fmt.Println("🎯 Keep coding and stay motivated!")

	return nil
}

func runEasterEggTests() error {
	fmt.Println("🧪 Easter Egg System Test")
	fmt.Println("=========================")
	fmt.Println()

	// Create easter egg manager
	eggManager := gamification.NewEasterEggManager()

	// Test different contexts
	testContexts := []struct {
		name    string
		context *gamification.EasterEggContext
	}{
		{
			name: "High Activity",
			context: &gamification.EasterEggContext{
				CommandsInTimeframe: 50,
				TimeframeDuration:   30 * time.Minute,
				LastCommand:         "git commit -m 'feat: awesome feature'",
				CommandHistory:      []string{"git add .", "git commit", "git push"},
			},
		},
		{
			name: "Coffee Break",
			context: &gamification.EasterEggContext{
				IdleDuration: 45 * time.Minute,
				LastCommand:  "ls",
			},
		},
		{
			name: "First Command Today",
			context: &gamification.EasterEggContext{
				IsFirstCommandToday: true,
				LastCommand:         "termonaut stats",
			},
		},
		{
			name: "Docker Usage",
			context: &gamification.EasterEggContext{
				LastCommand:    "docker run -it ubuntu bash",
				CommandHistory: []string{"docker build", "docker run"},
			},
		},
	}

	fmt.Println("🎭 Testing easter egg triggers...")
	fmt.Println()

	for i, test := range testContexts {
		fmt.Printf("📋 Test %d: %s\n", i+1, test.name)
		
		easterEgg := eggManager.CheckForEasterEgg(test.context)
		if easterEgg != "" {
			fmt.Printf("   🎉 Triggered: %s\n", easterEgg)
		} else {
			fmt.Printf("   ⚪ No easter egg triggered\n")
		}
		fmt.Println()
	}

	fmt.Println("✅ Easter egg system test complete!")
	fmt.Println()
	fmt.Printf("📊 System Status: %s\n", getEasterEggStatus())
	
	return nil
}

func showRandomEasterEgg() error {
	// Create a sample context
	context := &gamification.EasterEggContext{
		CommandsInTimeframe: 10,
		TimeframeDuration:   15 * time.Minute,
		LastCommand:         "termonaut easter-egg",
		IsFirstCommandToday: false,
	}

	eggManager := gamification.NewEasterEggManager()
	easterEgg := eggManager.CheckForEasterEgg(context)

	if easterEgg != "" {
		fmt.Println("🎉 Easter Egg!")
		fmt.Println("==============")
		fmt.Println()
		fmt.Println(easterEgg)
		fmt.Println()
	} else {
		fmt.Println("🎯 No easter egg this time!")
		fmt.Println("Try running more commands or come back later.")
		fmt.Println()
		fmt.Println("💡 Tip: Easter eggs are triggered by various activities:")
		fmt.Println("   • High command frequency")
		fmt.Println("   • Specific command patterns")
		fmt.Println("   • Time-based conditions")
		fmt.Println("   • Special occasions")
	}

	return nil
}

func getEasterEggStatus() string {
	cfg, err := config.Load()
	if err != nil {
		return "❓ Unknown"
	}
	
	if cfg.EasterEggsEnabled {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

func init() {
	easterEggCmd.Flags().BoolP("test", "t", false, "Run easter egg system tests")
	easterEggCmd.Flags().BoolP("motivational", "m", false, "Show a random motivational quote")
	rootCmd.AddCommand(easterEggCmd)
}
