package enhanced

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// AnimatedProgressBar represents an animated progress bar
type AnimatedProgressBar struct {
	current    float64
	target     float64
	width      int
	animSpeed  float64
	lastUpdate time.Time
	
	// Animation states
	isAnimating bool
	glowPhase   float64
	sparkles    []int
}

// NewAnimatedProgressBar creates a new animated progress bar
func NewAnimatedProgressBar(width int) *AnimatedProgressBar {
	return &AnimatedProgressBar{
		width:      width,
		animSpeed:  2.0, // Progress per second
		lastUpdate: time.Now(),
		sparkles:   make([]int, 0),
	}
}

// Update updates the progress bar animation
func (apb *AnimatedProgressBar) Update(targetPercent float64) {
	now := time.Now()
	elapsed := now.Sub(apb.lastUpdate).Seconds()
	apb.lastUpdate = now
	
	apb.target = targetPercent
	
	// Animate towards target
	if apb.current < apb.target {
		apb.current += apb.animSpeed * elapsed
		if apb.current > apb.target {
			apb.current = apb.target
		}
		apb.isAnimating = true
	} else {
		apb.isAnimating = false
	}
	
	// Update glow effect
	apb.glowPhase += elapsed * 3.0
	if apb.glowPhase > 2*3.14159 {
		apb.glowPhase -= 2*3.14159
	}
	
	// Update sparkles when animating
	if apb.isAnimating && len(apb.sparkles) < 3 {
		if elapsed > 0.1 { // Add sparkle every 100ms
			sparklePos := int(apb.current * float64(apb.width) / 100.0)
			if sparklePos > 0 && sparklePos < apb.width {
				apb.sparkles = append(apb.sparkles, sparklePos)
			}
		}
	}
	
	// Remove old sparkles
	newSparkles := make([]int, 0)
	for _, pos := range apb.sparkles {
		if time.Since(apb.lastUpdate) < time.Millisecond*500 {
			newSparkles = append(newSparkles, pos)
		}
	}
	apb.sparkles = newSparkles
}

// Render renders the animated progress bar
func (apb *AnimatedProgressBar) Render() string {
	filledWidth := int(apb.current * float64(apb.width) / 100.0)
	emptyWidth := apb.width - filledWidth
	
	// Create base progress bar
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)
	
	// Add glow effect when animating
	if apb.isAnimating {
		// Add pulsing glow
		glowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
		
		if filledWidth > 0 {
			// Make the last few characters glow
			glowWidth := min(3, filledWidth)
			normalWidth := filledWidth - glowWidth
			
			normalFilled := strings.Repeat("█", normalWidth)
			glowFilled := glowStyle.Render(strings.Repeat("█", glowWidth))
			filled = normalFilled + glowFilled
		}
	}
	
	// Add sparkles
	progressBar := filled + empty
	if len(apb.sparkles) > 0 {
		runes := []rune(progressBar)
		sparkleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
		
		for _, pos := range apb.sparkles {
			if pos < len(runes) {
				sparkleChar := sparkleStyle.Render("✨")
				// Replace character at position with sparkle
				before := string(runes[:pos])
				after := ""
				if pos+1 < len(runes) {
					after = string(runes[pos+1:])
				}
				progressBar = before + sparkleChar + after
				break // Only show one sparkle to avoid overlap
			}
		}
	}
	
	return progressBar
}

// IsAnimating returns true if the progress bar is currently animating
func (apb *AnimatedProgressBar) IsAnimating() bool {
	return apb.isAnimating
}

// XPProgressRenderer handles XP progress rendering with animations
type XPProgressRenderer struct {
	progressBar *AnimatedProgressBar
	lastXP      int
	lastLevel   int
}

// NewXPProgressRenderer creates a new XP progress renderer
func NewXPProgressRenderer() *XPProgressRenderer {
	return &XPProgressRenderer{
		progressBar: NewAnimatedProgressBar(20),
	}
}

// RenderXPProgress renders animated XP progress
func (xpr *XPProgressRenderer) RenderXPProgress(currentXP, currentLevel int) string {
	// Calculate progress
	currentLevelXP := (currentLevel - 1) * (currentLevel - 1) * 100
	nextLevelXP := currentLevel * currentLevel * 100
	progressXP := currentXP - currentLevelXP
	neededXP := nextLevelXP - currentLevelXP
	
	progressPercent := float64(progressXP) / float64(neededXP) * 100.0
	if progressPercent > 100 {
		progressPercent = 100
	}
	
	// Update animation
	xpr.progressBar.Update(progressPercent)
	
	// Check for level up
	levelUpMessage := ""
	if currentLevel > xpr.lastLevel && xpr.lastLevel > 0 {
		levelUpMessage = fmt.Sprintf("\n🎉 LEVEL UP! Welcome to Level %d! 🎉", currentLevel)
	}
	
	xpr.lastXP = currentXP
	xpr.lastLevel = currentLevel
	
	// Render with animation
	animatedBar := xpr.progressBar.Render()
	
	result := fmt.Sprintf(`🎮 Level Progress

🚀 Level %d Astronaut
📊 XP: %d / %d
🎯 Progress: [%s] %.1f%%

🌟 Next Level: %d XP needed
🏆 Total XP Earned: %d%s`,
		currentLevel, progressXP, neededXP, animatedBar, progressPercent,
		neededXP-progressXP, currentXP, levelUpMessage)
	
	return result
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
