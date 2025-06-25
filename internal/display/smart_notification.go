package display

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// SmartNotificationMethod represents different notification methods
type SmartNotificationMethod int

const (
	MethodSystemNotification SmartNotificationMethod = iota
	MethodTerminalTitle
	MethodTerminalBell
	MethodInlineMessage
	MethodStatusBar
)

// SmartNotifier provides intelligent notification selection
type SmartNotifier struct {
	preferredMethods []SmartNotificationMethod
	fallbackMethod   SmartNotificationMethod
}

// NewSmartNotifier creates a new smart notifier
func NewSmartNotifier() *SmartNotifier {
	return &SmartNotifier{
		preferredMethods: []SmartNotificationMethod{
			MethodSystemNotification,
			MethodTerminalTitle,
			MethodTerminalBell,
		},
		fallbackMethod: MethodInlineMessage,
	}
}

// ShowEasterEgg shows an easter egg using the best available method
func (sn *SmartNotifier) ShowEasterEgg(message string) error {
	// Try preferred methods in order
	for _, method := range sn.preferredMethods {
		if err := sn.tryMethod(method, message); err == nil {
			return nil
		}
	}
	
	// Fall back to inline message
	return sn.tryMethod(sn.fallbackMethod, message)
}

// tryMethod attempts to show notification using specific method (internal)
func (sn *SmartNotifier) tryMethod(method SmartNotificationMethod, message string) error {
	switch method {
	case MethodSystemNotification:
		return sn.showSystemNotification(message)
	case MethodTerminalTitle:
		return sn.showTerminalTitle(message)
	case MethodTerminalBell:
		return sn.showTerminalBell(message)
	case MethodInlineMessage:
		return sn.showInlineMessage(message)
	case MethodStatusBar:
		return sn.showStatusBar(message)
	default:
		return fmt.Errorf("unknown notification method")
	}
}

// showSystemNotification shows OS-level notification
func (sn *SmartNotifier) showSystemNotification(message string) error {
	switch runtime.GOOS {
	case "darwin": // macOS
		return sn.showMacOSNotification(message)
	case "linux":
		return sn.showLinuxNotification(message)
	case "windows":
		return sn.showWindowsNotification(message)
	default:
		return fmt.Errorf("system notifications not supported on %s", runtime.GOOS)
	}
}

// showMacOSNotification shows macOS system notification
func (sn *SmartNotifier) showMacOSNotification(message string) error {
	// Clean message for AppleScript
	cleanMessage := strings.ReplaceAll(message, `"`, `\"`)
	
	script := fmt.Sprintf(`display notification "%s" with title "Termonaut" subtitle "Easter Egg" sound name "Glass"`, cleanMessage)
	
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// showLinuxNotification shows Linux desktop notification
func (sn *SmartNotifier) showLinuxNotification(message string) error {
	// Try notify-send first
	if _, err := exec.LookPath("notify-send"); err == nil {
		cmd := exec.Command("notify-send", "Termonaut", message, "-t", "3000")
		return cmd.Run()
	}
	
	// Try zenity as fallback
	if _, err := exec.LookPath("zenity"); err == nil {
		cmd := exec.Command("zenity", "--info", "--text="+message, "--title=Termonaut")
		return cmd.Run()
	}
	
	return fmt.Errorf("no notification system found")
}

// showWindowsNotification shows Windows notification
func (sn *SmartNotifier) showWindowsNotification(message string) error {
	// Use PowerShell to show notification
	script := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('%s', 'Termonaut')`, message)
	cmd := exec.Command("powershell", "-Command", script)
	return cmd.Run()
}

// showTerminalTitle shows notification in terminal title
func (sn *SmartNotifier) showTerminalTitle(message string) error {
	// Save current title (if possible)
	originalTitle := os.Getenv("TERM_TITLE")
	
	// Set new title with easter egg
	fmt.Printf("\033]0;🎉 Termonaut: %s\007", message)
	
	// Restore title after delay
	go func() {
		time.Sleep(3 * time.Second)
		if originalTitle != "" {
			fmt.Printf("\033]0;%s\007", originalTitle)
		} else {
			// Reset to default
			fmt.Printf("\033]0;Terminal\007")
		}
	}()
	
	return nil
}

// showTerminalBell shows notification with bell sound and brief message
func (sn *SmartNotifier) showTerminalBell(message string) error {
	// Ring the bell
	fmt.Print("\a")
	
	// Show brief inline message that doesn't interfere
	fmt.Printf("\r🎉 %s\n", message)
	
	return nil
}

// showInlineMessage shows a safe inline message
func (sn *SmartNotifier) showInlineMessage(message string) error {
	// Show message on new line without interfering with cursor
	fmt.Printf("\n🎉 %s\n", message)
	return nil
}

// showStatusBar shows notification in status bar (tmux/screen)
func (sn *SmartNotifier) showStatusBar(message string) error {
	// Try tmux first
	if os.Getenv("TMUX") != "" {
		cmd := exec.Command("tmux", "display-message", "-d", "3000", fmt.Sprintf("🎉 %s", message))
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	
	// Try screen
	if os.Getenv("STY") != "" {
		cmd := exec.Command("screen", "-X", "echo", fmt.Sprintf("🎉 %s", message))
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	
	return fmt.Errorf("no status bar available")
}

// GetAvailableMethods returns list of available notification methods
func (sn *SmartNotifier) GetAvailableMethods() []SmartNotificationMethod {
	var available []SmartNotificationMethod
	
	// Test each method
	testMethods := []SmartNotificationMethod{
		MethodSystemNotification,
		MethodTerminalTitle,
		MethodTerminalBell,
		MethodStatusBar,
		MethodInlineMessage,
	}
	
	for _, method := range testMethods {
		if sn.isMethodAvailable(method) {
			available = append(available, method)
		}
	}
	
	return available
}

// isMethodAvailable checks if a notification method is available
func (sn *SmartNotifier) isMethodAvailable(method SmartNotificationMethod) bool {
	switch method {
	case MethodSystemNotification:
		switch runtime.GOOS {
		case "darwin":
			_, err := exec.LookPath("osascript")
			return err == nil
		case "linux":
			_, err1 := exec.LookPath("notify-send")
			_, err2 := exec.LookPath("zenity")
			return err1 == nil || err2 == nil
		case "windows":
			_, err := exec.LookPath("powershell")
			return err == nil
		}
		return false
	case MethodTerminalTitle:
		return true // Always available
	case MethodTerminalBell:
		return true // Always available
	case MethodStatusBar:
		return os.Getenv("TMUX") != "" || os.Getenv("STY") != ""
	case MethodInlineMessage:
		return true // Always available
	default:
		return false
	}
}

// MethodName returns human-readable name for notification method
func (sn *SmartNotifier) MethodName(method SmartNotificationMethod) string {
	switch method {
	case MethodSystemNotification:
		return "System Notification"
	case MethodTerminalTitle:
		return "Terminal Title"
	case MethodTerminalBell:
		return "Terminal Bell + Message"
	case MethodStatusBar:
		return "Status Bar (tmux/screen)"
	case MethodInlineMessage:
		return "Inline Message"
	default:
		return "Unknown"
	}
}

// TestAllMethods tests all available notification methods
func (sn *SmartNotifier) TestAllMethods() {
	fmt.Println("🧪 Testing all available notification methods...")
	fmt.Println()
	
	available := sn.GetAvailableMethods()
	
	for i, method := range available {
		methodName := sn.MethodName(method)
		fmt.Printf("📱 Testing %d/%d: %s\n", i+1, len(available), methodName)
		
		testMessage := fmt.Sprintf("Test notification %d - %s", i+1, methodName)
		
		if err := sn.tryMethod(method, testMessage); err != nil {
			fmt.Printf("   ❌ Failed: %v\n", err)
		} else {
			fmt.Printf("   ✅ Success\n")
		}
		
		// Wait between tests
		if i < len(available)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	
	fmt.Println()
	fmt.Println("✅ Testing complete!")
}

// ShowEasterEggSmart shows easter egg with smart method selection and user preference
func ShowEasterEggSmart(message string) {
	notifier := NewSmartNotifier()
	
	if err := notifier.ShowEasterEgg(message); err != nil {
		// If all methods fail, show simple message
		fmt.Printf("🎉 %s\n", message)
	}
}

// TryMethod attempts to show notification using specific method (exported version)
func (sn *SmartNotifier) TryMethod(method SmartNotificationMethod, message string) error {
	return sn.tryMethod(method, message)
}
