package avatar

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

// ConvertSVGToASCII converts SVG from URL to ASCII art using ascii-image-converter
func ConvertSVGToASCII(svgURL string, size AvatarSize) (string, error) {
	if svgURL == "" {
		return "", fmt.Errorf("empty SVG URL")
	}

	// Get dimensions from the AvatarSize struct
	width := size.ASCIIWidth
	height := size.ASCIIHeight

	// Since ascii-image-converter doesn't support SVG directly,
	// we need to convert SVG to PNG first
	pngData, err := convertSVGToPNG(svgURL, size.SVGSize)
	if err != nil {
		return "", fmt.Errorf("failed to convert SVG to PNG: %w", err)
	}

	// Create a temporary PNG file
	tempFile, err := createTempPNG(pngData)
	if err != nil {
		return "", fmt.Errorf("failed to create temp PNG file: %w", err)
	}
	defer os.Remove(tempFile) // Clean up

	// Configure ascii-image-converter options with enhanced quality settings
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{width, height}
	
	// Color and visual quality settings
	flags.Colored = true         // Enable 24-bit colors for vivid avatars
	flags.Complex = true         // Use extended 69-character set for maximum detail
	flags.CustomMap = ""         // Use default rich character set for best quality
	flags.Negative = false       // Don't invert colors
	flags.Grayscale = false      // Keep original colors, not grayscale
	flags.CharBackgroundColor = false // Use foreground colors only
	
	// Advanced quality settings
	flags.Braille = false        // Use regular ASCII for better compatibility
	flags.Threshold = 128        // Optimal threshold for detail preservation
	flags.Dither = false         // No dithering for cleaner look
	flags.Full = false           // Use specified dimensions exactly
	flags.FlipX = false          // No horizontal flip
	flags.FlipY = false          // No vertical flip
	
	// Size-specific optimizations for better quality
	if width >= 40 { // Large and medium avatars
		// Use custom character map with more detail for larger sizes
		flags.CustomMap = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	} else if width >= 20 { // Small avatars
		// Balanced character set for small avatars
		flags.CustomMap = " .:-=+*#%@"
	} else { // Mini avatars
		// Simple character set for very small avatars
		flags.CustomMap = " .:+#@"
	}
	
	// Font color enhancement (RGB values for better contrast)
	// Keep default font color for best automatic color mapping

	// Convert PNG to ASCII
	asciiArt, err := aic_package.Convert(tempFile, flags)
	if err != nil {
		return "", fmt.Errorf("failed to convert PNG to ASCII: %w", err)
	}

	// Clean up the ASCII art
	asciiArt = strings.TrimSpace(asciiArt)
	
	// Ensure we have some content
	if asciiArt == "" {
		return "", fmt.Errorf("conversion resulted in empty ASCII art")
	}
	
	return asciiArt, nil
}

// convertSVGToPNG converts SVG from URL to PNG data
func convertSVGToPNG(svgURL string, size int) ([]byte, error) {
	// For now, we'll try to get PNG directly from DiceBear by changing the format
	// DiceBear supports PNG format by changing the endpoint
	pngURL := strings.Replace(svgURL, "/svg", "/png", 1)
	
	// Add size parameter if not present
	if !strings.Contains(pngURL, "size=") {
		separator := "?"
		if strings.Contains(pngURL, "?") {
			separator = "&"
		}
		pngURL = fmt.Sprintf("%s%ssize=%d", pngURL, separator, size)
	}

	// Fetch PNG data
	resp, err := http.Get(pngURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PNG from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PNG request failed with status: %d", resp.StatusCode)
	}

	pngData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PNG data: %w", err)
	}

	return pngData, nil
}

// createTempPNG creates a temporary PNG file from PNG data
func createTempPNG(pngData []byte) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "avatar_*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Write PNG data to the file
	_, err = io.Copy(tempFile, bytes.NewReader(pngData))
	if err != nil {
		os.Remove(tempFile.Name()) // Clean up on error
		return "", fmt.Errorf("failed to write PNG data: %w", err)
	}

	return tempFile.Name(), nil
} 