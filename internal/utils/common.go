package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

// StringUtils provides string manipulation utilities
type StringUtils struct{}

// TruncateString truncates a string to the specified length with ellipsis
func (s StringUtils) TruncateString(str string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	
	if utf8.RuneCountInString(str) <= maxLen {
		return str
	}
	
	if maxLen <= 3 {
		return strings.Repeat(".", maxLen)
	}
	
	// Count runes to handle Unicode properly
	runes := []rune(str)
	return string(runes[:maxLen-3]) + "..."
}

// PadString pads a string to the specified width
func (s StringUtils) PadString(str string, width int, padChar rune) string {
	strLen := utf8.RuneCountInString(str)
	if strLen >= width {
		return str
	}
	
	padding := width - strLen
	return str + strings.Repeat(string(padChar), padding)
}

// CenterString centers a string within the specified width
func (s StringUtils) CenterString(str string, width int) string {
	strLen := utf8.RuneCountInString(str)
	if strLen >= width {
		return str
	}
	
	padding := width - strLen
	leftPad := padding / 2
	rightPad := padding - leftPad
	
	return strings.Repeat(" ", leftPad) + str + strings.Repeat(" ", rightPad)
}

// WrapText wraps text to fit within the specified width
func (s StringUtils) WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		// If adding this word would exceed width, start new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}
		
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	return lines
}

// TimeUtils provides time formatting utilities
type TimeUtils struct{}

// FormatDuration formats a duration in a human-readable way
func (t TimeUtils) FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
	} else {
		days := int(d.Hours()) / 24
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}
}

// FormatRelativeTime formats a time relative to now
func (t TimeUtils) FormatRelativeTime(timestamp time.Time) string {
	now := time.Now()
	diff := now.Sub(timestamp)
	
	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours()) / 24
		return fmt.Sprintf("%d day%s ago", days, pluralize(days))
	} else {
		return timestamp.Format("Jan 2, 2006")
	}
}

// FormatTimeOfDay formats time as HH:MM
func (t TimeUtils) FormatTimeOfDay(timestamp time.Time) string {
	return timestamp.Format("15:04")
}

// NumberUtils provides number formatting utilities
type NumberUtils struct{}

// FormatNumber formats a number with thousand separators
func (n NumberUtils) FormatNumber(num int) string {
	if num < 1000 {
		return fmt.Sprintf("%d", num)
	}
	
	str := fmt.Sprintf("%d", num)
	var result strings.Builder
	
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}
	
	return result.String()
}

// FormatBytes formats bytes in human-readable format
func (n NumberUtils) FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// FormatPercentage formats a percentage with specified decimal places
func (n NumberUtils) FormatPercentage(value, total float64, decimals int) string {
	if total == 0 {
		return "0%"
	}
	
	percentage := (value / total) * 100
	format := fmt.Sprintf("%%.%df%%%%", decimals)
	return fmt.Sprintf(format, percentage)
}

// ProgressUtils provides progress bar utilities
type ProgressUtils struct{}

// CreateProgressBar creates a text-based progress bar
func (p ProgressUtils) CreateProgressBar(current, total int, width int, filled, empty rune) string {
	if width <= 0 || total <= 0 {
		return ""
	}
	
	progress := float64(current) / float64(total)
	if progress > 1.0 {
		progress = 1.0
	}
	
	filledWidth := int(float64(width) * progress)
	emptyWidth := width - filledWidth
	
	return strings.Repeat(string(filled), filledWidth) + strings.Repeat(string(empty), emptyWidth)
}

// CreateProgressBarWithPercentage creates a progress bar with percentage
func (p ProgressUtils) CreateProgressBarWithPercentage(current, total int, width int) string {
	if width < 10 {
		return p.CreateProgressBar(current, total, width, '█', '░')
	}
	
	// Reserve space for percentage
	barWidth := width - 7 // " (100%)"
	if barWidth < 5 {
		barWidth = 5
	}
	
	bar := p.CreateProgressBar(current, total, barWidth, '█', '░')
	percentage := 0.0
	if total > 0 {
		percentage = float64(current) / float64(total) * 100
	}
	
	return fmt.Sprintf("%s (%.0f%%)", bar, percentage)
}

// ValidationUtils provides validation utilities
type ValidationUtils struct{}

// IsValidEmail checks if an email address is valid (basic check)
func (v ValidationUtils) IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// IsValidURL checks if a URL is valid (basic check)
func (v ValidationUtils) IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// SanitizeFilename removes invalid characters from a filename
func (v ValidationUtils) SanitizeFilename(filename string) string {
	// Replace invalid characters with underscore
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	return result
}

// MathUtils provides mathematical utilities
type MathUtils struct{}

// Clamp clamps a value between min and max
func (m MathUtils) Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampFloat clamps a float value between min and max
func (m MathUtils) ClampFloat(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// RoundToNearest rounds a number to the nearest multiple
func (m MathUtils) RoundToNearest(value, multiple int) int {
	return int(math.Round(float64(value)/float64(multiple))) * multiple
}

// Global utility instances
var (
	String     = StringUtils{}
	Time       = TimeUtils{}
	Number     = NumberUtils{}
	Progress   = ProgressUtils{}
	Validation = ValidationUtils{}
	Math       = MathUtils{}
)

// Helper functions

// pluralize returns "s" if count != 1
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// ChunkSlice splits a slice into chunks of specified size
func ChunkSlice(slice []string, chunkSize int) [][]string {
	if chunkSize <= 0 {
		return [][]string{slice}
	}
	
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	
	return chunks
}
