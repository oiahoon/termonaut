package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/oiahoon/termonaut/internal/display"
)

var notificationTestCmd = &cobra.Command{
	Use:   "notification-test",
	Short: "🔔 Test smart notification system",
	Long: `Test the improved smart notification system that uses multiple methods:

• System notifications (macOS/Linux/Windows)
• Terminal title updates
• Terminal bell + message
• Status bar notifications (tmux/screen)
• Safe inline messages

The system automatically selects the best available method for your environment.`,
	Run: runNotificationTest,
}

var (
	testAllMethods bool
	testMethod     string
	testMessage    string
)

func init() {
	notificationTestCmd.Flags().BoolVar(&testAllMethods, "all", false, "Test all available notification methods")
	notificationTestCmd.Flags().StringVar(&testMethod, "method", "", "Test specific method (system, title, bell, status, inline)")
	notificationTestCmd.Flags().StringVar(&testMessage, "message", "🎉 Test notification from Termonaut!", "Custom test message")
	
	rootCmd.AddCommand(notificationTestCmd)
}

func runNotificationTest(cmd *cobra.Command, args []string) {
	fmt.Println("🔔 Smart Notification System Test")
	fmt.Println("=================================")
	fmt.Println()
	
	notifier := display.NewSmartNotifier()
	
	if testAllMethods {
		// Test all available methods
		notifier.TestAllMethods()
		return
	}
	
	if testMethod != "" {
		// Test specific method
		runSpecificMethodTest(notifier, testMethod)
		return
	}
	
	// Default: show available methods and test smart selection
	runSmartSelectionTest(notifier)
}

func runSpecificMethodTest(notifier *display.SmartNotifier, methodName string) {
	fmt.Printf("🎯 Testing specific method: %s\n", methodName)
	fmt.Println()
	
	var method display.SmartNotificationMethod
	switch methodName {
	case "system":
		method = display.MethodSystemNotification
	case "title":
		method = display.MethodTerminalTitle
	case "bell":
		method = display.MethodTerminalBell
	case "status":
		method = display.MethodStatusBar
	case "inline":
		method = display.MethodInlineMessage
	default:
		fmt.Printf("❌ Unknown method: %s\n", methodName)
		fmt.Println("Available methods: system, title, bell, status, inline")
		return
	}
	
	fmt.Printf("📱 Testing: %s\n", notifier.MethodName(method))
	fmt.Printf("📝 Message: %s\n", testMessage)
	fmt.Println()
	
	if err := notifier.TryMethod(method, testMessage); err != nil {
		fmt.Printf("❌ Test failed: %v\n", err)
	} else {
		fmt.Printf("✅ Test successful!\n")
	}
}

func runSmartSelectionTest(notifier *display.SmartNotifier) {
	fmt.Println("🧠 Smart Notification Selection Test")
	fmt.Println()
	
	// Show available methods
	available := notifier.GetAvailableMethods()
	fmt.Printf("📋 Available notification methods (%d):\n", len(available))
	for i, method := range available {
		fmt.Printf("   %d. %s\n", i+1, notifier.MethodName(method))
	}
	fmt.Println()
	
	// Test smart selection
	fmt.Println("🎯 Testing smart method selection...")
	fmt.Printf("📝 Message: %s\n", testMessage)
	fmt.Println()
	
	fmt.Println("⏰ Showing notification in 3 seconds...")
	time.Sleep(3 * time.Second)
	
	if err := notifier.ShowEasterEgg(testMessage); err != nil {
		fmt.Printf("❌ Smart selection failed: %v\n", err)
	} else {
		fmt.Printf("✅ Smart selection successful!\n")
		fmt.Println()
		fmt.Println("💡 The system automatically chose the best available method.")
		fmt.Println("   Check your system notifications, terminal title, or listen for bell sound.")
	}
	
	fmt.Println()
	fmt.Println("🔧 Advanced testing options:")
	fmt.Println("   --all              Test all available methods")
	fmt.Println("   --method system    Test system notifications only")
	fmt.Println("   --method title     Test terminal title only")
	fmt.Println("   --method bell      Test terminal bell only")
	fmt.Println("   --method status    Test status bar only")
	fmt.Println("   --method inline    Test inline message only")
	fmt.Println("   --message \"text\"   Use custom test message")
}
