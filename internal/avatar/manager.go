package avatar

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

// AvatarManager handles avatar generation, caching, and management
type AvatarManager struct {
	dicebear *DiceBearClient
	cache    *AvatarCache
	config   *Config
}

// Config holds avatar system configuration
type Config struct {
	CacheDir     string
	CacheTTL     time.Duration
	APITimeout   time.Duration
	DefaultStyle string
	DefaultSize  AvatarSize
}

// AvatarRequest represents a request for avatar generation
type AvatarRequest struct {
	Username string
	Level    int
	Style    string
	Size     AvatarSize
}

// AvatarSize defines the dimensions for avatar generation
type AvatarSize struct {
	SVGSize     int // 64, 128, 256
	ASCIIWidth  int // 20, 40, 60
	ASCIIHeight int // 10, 20, 30
}

// Avatar represents a generated avatar with metadata
type Avatar struct {
	Username    string
	Level       int
	Style       string
	Size        AvatarSize
	SVGData     []byte
	ASCIIArt    string
	Seed        string
	GeneratedAt time.Time
	CacheKey    string
}

// Predefined avatar sizes
var (
	SizeMini   = AvatarSize{SVGSize: 32, ASCIIWidth: 10, ASCIIHeight: 5}
	SizeSmall  = AvatarSize{SVGSize: 64, ASCIIWidth: 20, ASCIIHeight: 10}
	SizeMedium = AvatarSize{SVGSize: 128, ASCIIWidth: 40, ASCIIHeight: 20}
	SizeLarge  = AvatarSize{SVGSize: 256, ASCIIWidth: 60, ASCIIHeight: 30}
)

// NewAvatarManager creates a new avatar manager with the given configuration
func NewAvatarManager(config *Config) (*AvatarManager, error) {
	cache, err := NewAvatarCache(config.CacheDir, config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create avatar cache: %w", err)
	}

	dicebear := NewDiceBearClient(config.APITimeout)

	return &AvatarManager{
		dicebear: dicebear,
		cache:    cache,
		config:   config,
	}, nil
}

// Generate creates or retrieves an avatar for the given request
func (am *AvatarManager) Generate(request AvatarRequest) (*Avatar, error) {
	// Validate request
	if err := am.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid avatar request: %w", err)
	}

	// Generate cache key
	cacheKey := am.generateCacheKey(request)

	// Try to get from cache first
	if cached, err := am.cache.Get(cacheKey); err == nil && cached != nil {
		log.Printf("Avatar cache hit for key: %s", cacheKey)
		return cached, nil
	}

	log.Printf("Avatar cache miss for key: %s, generating new avatar", cacheKey)

	// Generate new avatar
	avatar, err := am.generateNew(request, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate avatar: %w", err)
	}

	// Cache the result
	if err := am.cache.Set(cacheKey, avatar); err != nil {
		log.Printf("Warning: failed to cache avatar: %v", err)
		// Don't fail the request if caching fails
	}

	return avatar, nil
}

// GetCached retrieves an avatar from cache if it exists
func (am *AvatarManager) GetCached(cacheKey string) (*Avatar, error) {
	return am.cache.Get(cacheKey)
}

// Refresh forces regeneration of an avatar, bypassing cache
func (am *AvatarManager) Refresh(username string) error {
	// Get user's current level from database
	level := am.getUserLevel()

	request := AvatarRequest{
		Username: username,
		Level:    level,
		Style:    am.config.DefaultStyle,
		Size:     am.config.DefaultSize,
	}

	cacheKey := am.generateCacheKey(request)

	// Remove from cache
	if err := am.cache.Delete(cacheKey); err != nil {
		log.Printf("Warning: failed to delete cached avatar: %v", err)
	}

	// Generate new avatar
	_, err := am.Generate(request)
	return err
}

// generateNew creates a new avatar from scratch
func (am *AvatarManager) generateNew(request AvatarRequest, cacheKey string) (*Avatar, error) {
	// Generate deterministic seed
	seed := am.generateSeed(request)

	// Create DiceBear parameters
	params := am.createDiceBearParams(request, seed)

	// Generate SVG URL for DiceBear
	svgURL := am.dicebear.GenerateAvatarURL(request.Style, params)

	// Fetch SVG data for caching with network error handling
	svgData, err := am.dicebear.GenerateAvatar(request.Style, params)
	if err != nil {
		// Check if it's a network-related error
		if am.isNetworkError(err) {
			return am.generateFallbackAvatar(request, cacheKey, seed,
				"🌐 Network issue: Unable to fetch avatar from DiceBear API. Using offline fallback.")
		}
		return nil, fmt.Errorf("failed to fetch SVG from DiceBear: %w", err)
	}

	// Convert SVG to ASCII using the new converter
	asciiArt, err := ConvertSVGToASCII(svgURL, request.Size)
	if err != nil {
		// If ASCII conversion fails, still return the avatar with SVG data
		log.Printf("Warning: ASCII conversion failed for %s: %v", request.Username, err)
		asciiArt = am.generateFallbackASCII(request)
	}

	avatar := &Avatar{
		Username:    request.Username,
		Level:       request.Level,
		Style:       request.Style,
		Size:        request.Size,
		SVGData:     svgData,
		ASCIIArt:    asciiArt,
		Seed:        seed,
		GeneratedAt: time.Now(),
		CacheKey:    cacheKey,
	}

	return avatar, nil
}

// generateSeed creates a deterministic seed for avatar generation
func (am *AvatarManager) generateSeed(request AvatarRequest) string {
	// Create seed that changes with level progression
	levelTier := request.Level / 5 // Changes every 5 levels
	return fmt.Sprintf("%s:%d:%d", request.Username, request.Level, levelTier)
}

// generateCacheKey creates a unique cache key for the request
func (am *AvatarManager) generateCacheKey(request AvatarRequest) string {
	data := fmt.Sprintf("%s:%d:%s:%dx%d",
		request.Username,
		request.Level,
		request.Style,
		request.Size.ASCIIWidth,
		request.Size.ASCIIHeight,
	)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// createDiceBearParams creates parameters for DiceBear API based on level
func (am *AvatarManager) createDiceBearParams(request AvatarRequest, seed string) DiceBearParams {
	params := DiceBearParams{
		Seed: seed,
		Size: request.Size.SVGSize,
	}

	// Enhanced level-based evolution system
	// Following the specification: Level 1-4: Basic, 5-9: Accessories, 10-19: Colors, 20-49: Backgrounds, 50+: Special
	switch {
	case request.Level >= 100:
		// Epic/Legendary variants (Level 100+)
		params.HairColor = []string{"ff6b6b", "4ecdc4", "45b7d1", "f39c12", "9b59b6", "e74c3c"}
		params.BackgroundType = []string{"gradientLinear", "gradientRadial"}
		params.BackgroundRotation = []int{45, 90, 135, 180}
		if request.Style == "pixel-art" {
			// Add special accessories for epic levels
			params.Accessories = []string{"glasses", "hat"}
			params.AccessoriesColor = []string{"ff6b6b", "4ecdc4", "45b7d1"}
		}
		params.Flip = request.Level%7 == 0 // Occasional flip for variety

	case request.Level >= 50:
		// Animated/Special elements (Level 50-99)
		params.HairColor = []string{"724133", "f59797", "65c9ff", "92d5ea", "fbbf24", "e67e22"}
		params.BackgroundType = []string{"gradientLinear"}
		params.BackgroundRotation = []int{0, 45, 90}
		if request.Style == "pixel-art" {
			params.Accessories = []string{"glasses"}
			params.AccessoriesColor = []string{"333333", "666666", "999999"}
		}
		params.Scale = 100 + (request.Level-50)*2 // Slightly larger scale

	case request.Level >= 20:
		// Special backgrounds (Level 20-49)
		params.HairColor = []string{"724133", "f59797", "65c9ff", "92d5ea", "fbbf24"}
		params.BackgroundType = []string{"solid"}
		if request.Style == "pixel-art" {
			params.Accessories = []string{"glasses"}
		}
		// Subtle rotation based on level
		params.Rotate = (request.Level - 20) * 2

	case request.Level >= 10:
		// Color theme changes (Level 10-19)
		params.HairColor = []string{"724133", "f59797", "65c9ff", "92d5ea"}
		// Start introducing subtle transformations
		params.TranslateX = (request.Level - 10) % 5
		params.TranslateY = (request.Level - 10) % 3

	case request.Level >= 5:
		// Add accessories (Level 5-9)
		params.HairColor = []string{"724133", "f59797", "65c9ff"}
		if request.Style == "pixel-art" {
			// Only add accessories for compatible styles
			if request.Level >= 7 {
				params.Accessories = []string{"glasses"}
			}
		}

	default:
		// Basic avatar (Level 1-4) - minimal customization
		// Use default parameters with just the seed
	}

	// Add some deterministic variety based on username
	usernameHash := 0
	for _, char := range request.Username {
		usernameHash += int(char)
	}

	// Use hash to select variants within level constraints
	if len(params.HairColor) > 0 {
		colorIndex := usernameHash % len(params.HairColor)
		params.HairColor = []string{params.HairColor[colorIndex]}
	}

	// Add subtle radius variation for higher levels
	if request.Level >= 10 {
		params.Radius = 10 + (usernameHash % 20) // 10-30 radius
	}

	return params
}

// validateRequest validates the avatar request parameters
func (am *AvatarManager) validateRequest(request AvatarRequest) error {
	if request.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if request.Level < 1 {
		return fmt.Errorf("level must be at least 1")
	}

	if request.Style == "" {
		request.Style = am.config.DefaultStyle
	}

	// Default to pixel-art style if not specified or invalid
	validStyles := []string{"pixel-art", "bottts", "adventurer", "avataaars"}
	valid := false
	for _, style := range validStyles {
		if request.Style == style {
			valid = true
			break
		}
	}
	if !valid {
		request.Style = "pixel-art" // Default to pixel-art
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	return &Config{
		CacheDir:     "~/.termonaut/avatars",
		CacheTTL:     7 * 24 * time.Hour, // 7 days
		APITimeout:   10 * time.Second,
		DefaultStyle: "pixel-art",
		DefaultSize:  SizeSmall,
	}
}

// isNetworkError checks if an error is network-related
func (am *AvatarManager) isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common network error patterns
	errStr := strings.ToLower(err.Error())

	// Network connectivity issues
	if strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection timeout") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "temporary failure in name resolution") {
		return true
	}

	// Check for specific error types
	if _, ok := err.(*net.DNSError); ok {
		return true
	}
	if _, ok := err.(*net.OpError); ok {
		return true
	}
	if _, ok := err.(*url.Error); ok {
		return true
	}

	return false
}

// generateFallbackAvatar creates a fallback avatar when network fails
func (am *AvatarManager) generateFallbackAvatar(request AvatarRequest, cacheKey, seed, message string) (*Avatar, error) {
	log.Printf("Network error for user %s: %s", request.Username, message)

	// Generate a simple fallback SVG
	fallbackSVG := am.generateFallbackSVG(request, seed)

	// Generate fallback ASCII art
	fallbackASCII := am.generateFallbackASCII(request)

	avatar := &Avatar{
		Username:    request.Username,
		Level:       request.Level,
		Style:       request.Style + "-fallback",
		Size:        request.Size,
		SVGData:     []byte(fallbackSVG),
		ASCIIArt:    fallbackASCII,
		Seed:        seed,
		GeneratedAt: time.Now(),
		CacheKey:    cacheKey,
	}

	return avatar, nil
}

// generateFallbackSVG creates a simple SVG when network is unavailable
func (am *AvatarManager) generateFallbackSVG(request AvatarRequest, seed string) string {
	// Create a simple geometric avatar based on username and level
	size := request.Size.SVGSize

	// Generate colors based on username hash
	hash := 0
	for _, char := range request.Username {
		hash += int(char)
	}

	// Level-based color progression
	hue := (hash + request.Level*30) % 360
	saturation := 60 + (request.Level % 40)
	lightness := 40 + (request.Level % 30)

	return fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
		<defs>
			<linearGradient id="grad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
				<stop offset="0%%" style="stop-color:hsl(%d,%d%%,%d%%);stop-opacity:1" />
				<stop offset="100%%" style="stop-color:hsl(%d,%d%%,%d%%);stop-opacity:1" />
			</linearGradient>
		</defs>
		<circle cx="%d" cy="%d" r="%d" fill="url(#grad)" />
		<text x="%d" y="%d" font-family="monospace" font-size="%d" text-anchor="middle" fill="white">%s</text>
		<text x="%d" y="%d" font-family="monospace" font-size="%d" text-anchor="middle" fill="white">Lv%d</text>
	</svg>`,
		size, size,
		hue, saturation, lightness,
		hue+60, saturation-10, lightness+10,
		size/2, size/2, size/2-5,
		size/2, size/2-5, size/8, string(request.Username[0]),
		size/2, size/2+15, size/12, request.Level)
}

// generateFallbackASCII creates fallback ASCII art when conversion fails
func (am *AvatarManager) generateFallbackASCII(request AvatarRequest) string {
	width := request.Size.ASCIIWidth
	height := request.Size.ASCIIHeight

	// Create a simple box avatar with username initial and level
	if width < 10 || height < 5 {
		// Very small size - minimal representation
		return fmt.Sprintf("[%c%d]", request.Username[0], request.Level)
	}

	// Generate a simple box avatar
	var lines []string

	// Top border
	lines = append(lines, strings.Repeat("=", width))

	// Empty lines
	for i := 1; i < height-3; i++ {
		lines = append(lines, "|"+strings.Repeat(" ", width-2)+"|")
	}

	// Content line with initial and level
	content := fmt.Sprintf("%c Lv%d", request.Username[0], request.Level)
	padding := (width - len(content) - 2) / 2
	if padding < 0 {
		padding = 0
	}
	contentLine := "|" + strings.Repeat(" ", padding) + content + strings.Repeat(" ", width-2-padding-len(content)) + "|"
	lines = append(lines, contentLine)

	// Bottom border
	lines = append(lines, strings.Repeat("=", width))

	return strings.Join(lines, "\n")
}

// TestConnection tests the connection to the avatar service
func (am *AvatarManager) TestConnection() error {
	return am.dicebear.TestConnection()
}

// GetNetworkStatus returns the current network status for avatar generation
func (am *AvatarManager) GetNetworkStatus() (bool, error) {
	err := am.TestConnection()
	if err != nil {
		if am.isNetworkError(err) {
			return false, fmt.Errorf("network connectivity issue: %w", err)
		}
		return false, fmt.Errorf("service error: %w", err)
	}
	return true, nil
}

// getUserLevel gets the user's current level from the database
func (am *AvatarManager) getUserLevel() int {
	// Try to get user progress from database
	// This requires access to the database, which we'll need to add to the manager
	// For now, return a reasonable default based on cache or config
	
	// Check if we have cached level information
	if am.cache != nil {
		if cachedAvatar := am.cache.Get("current_user_avatar"); cachedAvatar != nil {
			if avatar, ok := cachedAvatar.(*Avatar); ok && avatar.Level > 0 {
				return avatar.Level
			}
		}
	}
	
	// Default to level 1 if no information available
	return 1
}
