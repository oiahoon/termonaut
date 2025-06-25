package display

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// TerminalType represents different terminal types
type TerminalType int

const (
	TerminalBasic TerminalType = iota
	TerminalITerm2
	TerminalWarp
	TerminalVSCode
	TerminalModern
)

// NotificationPosition defines where to show the notification
type NotificationPosition int

const (
	PositionTop NotificationPosition = iota
	PositionBottom
	PositionCenter
)

// FloatingNotification represents a floating notification in the terminal
type FloatingNotification struct {
	Message  string
	Duration time.Duration
	Position NotificationPosition
	Style    lipgloss.Style
}

// TerminalDetector detects the current terminal type
type TerminalDetector struct{}

// DetectTerminal detects the current terminal type
func (td *TerminalDetector) DetectTerminal() TerminalType {
	// Check environment variables
	if term := os.Getenv("TERM_PROGRAM"); term != "" {
		switch term {
		case "iTerm.app":
			return TerminalITerm2
		case "vscode":
			return TerminalVSCode
		case "Warp":
			return TerminalWarp
		}
	}
	
	// Check for modern terminal features
	if td.supportsOSC() {
		return TerminalModern
	}
	
	return TerminalBasic
}

// supportsOSC checks if terminal supports OSC sequences
func (td *TerminalDetector) supportsOSC() bool {
	// Simple heuristic - check for common modern terminal indicators
	colorTerm := os.Getenv("COLORTERM")
	return colorTerm == "truecolor" || colorTerm == "24bit"
}

// FloatingNotifier handles floating notifications
type FloatingNotifier struct {
	detector     *TerminalDetector
	terminalType TerminalType
}

// NewFloatingNotifier creates a new floating notifier
func NewFloatingNotifier() *FloatingNotifier {
	detector := &TerminalDetector{}
	return &FloatingNotifier{
		detector:     detector,
		terminalType: detector.DetectTerminal(),
	}
}

// ShowEasterEgg shows a floating easter egg notification
func (fn *FloatingNotifier) ShowEasterEgg(message string, duration time.Duration) {
	notification := &FloatingNotification{
		Message:  message,
		Duration: duration,
		Position: PositionTop,
		Style:    fn.getEasterEggStyle(),
	}
	
	switch fn.terminalType {
	case TerminalITerm2:
		fn.showITerm2Notification(notification)
	case TerminalWarp:
		fn.showWarpNotification(notification)
	case TerminalVSCode:
		fn.showVSCodeNotification(notification)
	default:
		fn.showANSINotification(notification)
	}
}

// getEasterEggStyle returns the style for easter egg notifications
func (fn *FloatingNotifier) getEasterEggStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("226")).  // Bright yellow
		Foreground(lipgloss.Color("0")).    // Black text
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214"))
}

// showITerm2Notification shows notification using iTerm2 features
func (fn *FloatingNotifier) showITerm2Notification(notification *FloatingNotification) {
	// iTerm2 supports OSC 9 for notifications
	osc := fmt.Sprintf("\033]9;%s\007", notification.Message)
	fmt.Print(osc)
	
	// Also show ANSI version as fallback
	fn.showANSINotification(notification)
}

// showWarpNotification shows notification using Warp features
func (fn *FloatingNotifier) showWarpNotification(notification *FloatingNotification) {
	// Warp has special notification support
	// For now, fallback to ANSI
	fn.showANSINotification(notification)
}

// showVSCodeNotification shows notification using VS Code terminal features
func (fn *FloatingNotifier) showVSCodeNotification(notification *FloatingNotification) {
	// VS Code terminal supports special sequences
	// For now, fallback to ANSI
	fn.showANSINotification(notification)
}

// showANSINotification shows notification using ANSI escape sequences
func (fn *FloatingNotifier) showANSINotification(notification *FloatingNotification) {
	// Save current cursor position and hide cursor
	fmt.Print("\033[s\033[?25l")
	
	// Move to notification position
	switch notification.Position {
	case PositionTop:
		fmt.Print("\033[1;1H")
	case PositionBottom:
		fmt.Print("\033[999;1H") // Move to bottom
	case PositionCenter:
		fmt.Print("\033[12;1H") // Approximate center
	}
	
	// Clear the line and show notification
	fmt.Print("\033[K")
	styledMessage := notification.Style.Render(notification.Message)
	fmt.Print(styledMessage)
	
	// Set up auto-hide
	go func() {
		time.Sleep(notification.Duration)
		fn.hideANSINotification(notification)
	}()
}

// hideANSINotification hides the ANSI notification
func (fn *FloatingNotifier) hideANSINotification(notification *FloatingNotification) {
	// Move to notification position and clear
	switch notification.Position {
	case PositionTop:
		fmt.Print("\033[1;1H")
	case PositionBottom:
		fmt.Print("\033[999;1H")
	case PositionCenter:
		fmt.Print("\033[12;1H")
	}
	
	fmt.Print("\033[K")      // Clear line
	fmt.Print("\033[u")      // Restore cursor position
	fmt.Print("\033[?25h")   // Show cursor
}

// ShowFloatingEasterEgg is a convenience function for showing easter eggs
func ShowFloatingEasterEgg(message string) {
	notifier := NewFloatingNotifier()
	
	// Add some visual flair to the message
	decoratedMessage := fmt.Sprintf("🎉 %s 🎉", message)
	
	// Show for 3 seconds
	notifier.ShowEasterEgg(decoratedMessage, 3*time.Second)
}

// TestFloatingNotification tests the floating notification system
func TestFloatingNotification() {
	notifier := NewFloatingNotifier()
	
	fmt.Println("Testing floating notifications...")
	fmt.Printf("Detected terminal type: %v\n", notifier.terminalType)
	
	// Test different messages
	messages := []string{
		"🚀 Welcome to the space program!",
		"☕ Coffee break detected!",
		"🎮 Level up! You're now a Space Commander!",
		"🦆 Rubber duck debugging activated!",
	}
	
	for i, msg := range messages {
		fmt.Printf("Showing notification %d...\n", i+1)
		notifier.ShowEasterEgg(msg, 2*time.Second)
		time.Sleep(3 * time.Second) // Wait for notification to disappear
	}
	
	fmt.Println("Test complete!")
}

// SafeNotificationManager manages notifications safely
type SafeNotificationManager struct {
	notifier     *FloatingNotifier
	lastShown    time.Time
	cooldown     time.Duration
	isUserTyping bool
}

// NewSafeNotificationManager creates a safe notification manager
func NewSafeNotificationManager() *SafeNotificationManager {
	return &SafeNotificationManager{
		notifier: NewFloatingNotifier(),
		cooldown: 5 * time.Second, // 5 second cooldown between notifications
	}
}

// ShouldShowNotification checks if it's safe to show a notification
func (snm *SafeNotificationManager) ShouldShowNotification() bool {
	// Check cooldown
	if time.Since(snm.lastShown) < snm.cooldown {
		return false
	}
	
	// Check if user is typing (simple heuristic)
	if snm.isUserTyping {
		return false
	}
	
	return true
}

// ShowEasterEggSafely shows an easter egg notification safely
func (snm *SafeNotificationManager) ShowEasterEggSafely(message string) bool {
	if !snm.ShouldShowNotification() {
		return false
	}
	
	snm.notifier.ShowEasterEgg(message, 3*time.Second)
	snm.lastShown = time.Now()
	return true
}

// SetUserTyping sets the user typing state
func (snm *SafeNotificationManager) SetUserTyping(typing bool) {
	snm.isUserTyping = typing
}
