package avatar

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DiceBearClient handles communication with DiceBear API
type DiceBearClient struct {
	baseURL    string
	httpClient *http.Client
}

// DiceBearParams represents parameters for DiceBear API
type DiceBearParams struct {
	Seed                 string
	Size                 int
	BackgroundType       []string
	BackgroundRotation   []int
	Accessories          []string
	AccessoriesColor     []string
	HairColor            []string
	Flip                 bool
	Rotate               int
	Scale                int
	Radius               int
	TranslateX           int
	TranslateY           int
}

// NewDiceBearClient creates a new DiceBear API client
func NewDiceBearClient(timeout time.Duration) *DiceBearClient {
	return &DiceBearClient{
		baseURL: "https://api.dicebear.com/9.x",
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GenerateAvatar fetches an SVG avatar from DiceBear API
func (c *DiceBearClient) GenerateAvatar(style string, params DiceBearParams) ([]byte, error) {
	// Build the URL
	apiURL, err := c.buildURL(style, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build API URL: %w", err)
	}

	// Make the HTTP request
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to DiceBear API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DiceBear API returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	svgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Basic validation - check if it looks like SVG
	if !strings.Contains(string(svgData), "<svg") {
		return nil, fmt.Errorf("response does not appear to be valid SVG")
	}

	return svgData, nil
}

// GenerateAvatarURL returns the DiceBear API URL for the given style and parameters
func (c *DiceBearClient) GenerateAvatarURL(style string, params DiceBearParams) string {
	// Build the URL (ignore error for now, as we'll handle it in the caller)
	apiURL, err := c.buildURL(style, params)
	if err != nil {
		// Return a fallback URL or handle error appropriately
		return fmt.Sprintf("%s/%s/svg?seed=%s", c.baseURL, style, params.Seed)
	}
	return apiURL
}

// buildURL constructs the DiceBear API URL with parameters
func (c *DiceBearClient) buildURL(style string, params DiceBearParams) (string, error) {
	// Validate style
	if style == "" {
		return "", fmt.Errorf("style cannot be empty")
	}

	// Start with base URL
	baseURL := fmt.Sprintf("%s/%s/svg", c.baseURL, style)
	
	// Parse URL to add query parameters
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Build query parameters
	query := parsedURL.Query()

	// Add seed (required for deterministic generation)
	if params.Seed != "" {
		query.Set("seed", params.Seed)
	}

	// Add size if specified
	if params.Size > 0 {
		query.Set("size", strconv.Itoa(params.Size))
	}

	// Add optional parameters (simplified for now)
	// Note: Some parameters might not be compatible with all styles
	// We'll keep it simple initially and add more as needed
	
	if len(params.HairColor) > 0 && style == "pixel-art" {
		// Only add hair color for pixel-art style for now
		query.Set("hairColor", params.HairColor[0]) // Use first color only
	}

	if params.Flip {
		query.Set("flip", "true")
	}

	if params.Rotate != 0 {
		query.Set("rotate", strconv.Itoa(params.Rotate))
	}

	if params.Scale != 0 {
		query.Set("scale", strconv.Itoa(params.Scale))
	}

	if params.Radius != 0 {
		query.Set("radius", strconv.Itoa(params.Radius))
	}

	if params.TranslateX != 0 {
		query.Set("translateX", strconv.Itoa(params.TranslateX))
	}

	if params.TranslateY != 0 {
		query.Set("translateY", strconv.Itoa(params.TranslateY))
	}

	// Set the query and return the URL
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

// GetSupportedStyles returns a list of supported avatar styles
func (c *DiceBearClient) GetSupportedStyles() []string {
	return []string{
		"pixel-art",
		"bottts",
		"adventurer",
		"avataaars",
		"big-ears",
		"big-smile",
		"croodles",
		"fun-emoji",
		"icons",
		"identicon",
		"initials",
		"lorelei",
		"micah",
		"miniavs",
		"notionists",
		"open-peeps",
		"personas",
		"rings",
		"shapes",
		"thumbs",
	}
}

// ValidateStyle checks if a style is supported
func (c *DiceBearClient) ValidateStyle(style string) bool {
	supported := c.GetSupportedStyles()
	for _, s := range supported {
		if s == style {
			return true
		}
	}
	return false
}

// GetStyleInfo returns information about a specific style
func (c *DiceBearClient) GetStyleInfo(style string) StyleInfo {
	styleInfoMap := map[string]StyleInfo{
		"pixel-art": {
			Name:        "Pixel Art",
			Description: "Retro 8-bit style avatars perfect for terminal display",
			Recommended: true,
			Features:    []string{"accessories", "hair", "clothing"},
		},
		"bottts": {
			Name:        "Bottts",
			Description: "Robot-themed avatars with clean geometric shapes",
			Recommended: true,
			Features:    []string{"colors", "accessories", "antennas"},
		},
		"adventurer": {
			Name:        "Adventurer",
			Description: "Fantasy character avatars with medieval themes",
			Recommended: true,
			Features:    []string{"hair", "facial-hair", "accessories"},
		},
		"avataaars": {
			Name:        "Avataaars",
			Description: "Modern cartoon-style avatars with many customization options",
			Recommended: false,
			Features:    []string{"hair", "accessories", "clothing", "facial-hair"},
		},
	}

	if info, exists := styleInfoMap[style]; exists {
		return info
	}

	// Return default info for unknown styles
	return StyleInfo{
		Name:        style,
		Description: "Avatar style",
		Recommended: false,
		Features:    []string{},
	}
}

// StyleInfo contains information about an avatar style
type StyleInfo struct {
	Name        string
	Description string
	Recommended bool
	Features    []string
}

// TestConnection tests the connection to DiceBear API
func (c *DiceBearClient) TestConnection() error {
	// Try to generate a simple test avatar
	params := DiceBearParams{
		Seed: "test",
		Size: 32,
	}

	_, err := c.GenerateAvatar("pixel-art", params)
	if err != nil {
		return fmt.Errorf("DiceBear API connection test failed: %w", err)
	}

	return nil
} 