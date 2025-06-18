package avatar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AvatarCache handles local caching of avatars
type AvatarCache struct {
	cacheDir string
	ttl      time.Duration
}

// CacheEntry represents a cached avatar with metadata
type CacheEntry struct {
	Avatar      *Avatar   `json:"avatar"`
	CachedAt    time.Time `json:"cached_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	AccessCount int       `json:"access_count"`
	LastAccess  time.Time `json:"last_access"`
}

// NewAvatarCache creates a new avatar cache
func NewAvatarCache(cacheDir string, ttl time.Duration) (*AvatarCache, error) {
	// Expand home directory if needed
	if strings.HasPrefix(cacheDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, cacheDir[2:])
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create subdirectories
	svgDir := filepath.Join(cacheDir, "svg")
	asciiDir := filepath.Join(cacheDir, "ascii")
	metaDir := filepath.Join(cacheDir, "meta")

	for _, dir := range []string{svgDir, asciiDir, metaDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache subdirectory %s: %w", dir, err)
		}
	}

	cache := &AvatarCache{
		cacheDir: cacheDir,
		ttl:      ttl,
	}

	// Clean up expired entries on startup
	go cache.cleanupExpired()

	return cache, nil
}

// Get retrieves an avatar from cache
func (c *AvatarCache) Get(cacheKey string) (*Avatar, error) {
	metaPath := c.getMetaPath(cacheKey)
	
	// Check if metadata file exists
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cache entry not found")
	}

	// Read metadata
	entry, err := c.readCacheEntry(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache entry: %w", err)
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		// Clean up expired entry
		c.Delete(cacheKey)
		return nil, fmt.Errorf("cache entry expired")
	}

	// Update access statistics
	entry.AccessCount++
	entry.LastAccess = time.Now()
	c.writeCacheEntry(metaPath, entry)

	// Load SVG and ASCII data
	avatar := entry.Avatar
	
	// Load SVG data
	svgPath := c.getSVGPath(cacheKey)
	if svgData, err := ioutil.ReadFile(svgPath); err == nil {
		avatar.SVGData = svgData
	}

	// Load ASCII data
	asciiPath := c.getASCIIPath(cacheKey)
	if asciiData, err := ioutil.ReadFile(asciiPath); err == nil {
		avatar.ASCIIArt = string(asciiData)
	}

	return avatar, nil
}

// Set stores an avatar in cache
func (c *AvatarCache) Set(cacheKey string, avatar *Avatar) error {
	now := time.Now()
	
	entry := &CacheEntry{
		Avatar:      avatar,
		CachedAt:    now,
		ExpiresAt:   now.Add(c.ttl),
		AccessCount: 1,
		LastAccess:  now,
	}

	// Save metadata
	metaPath := c.getMetaPath(cacheKey)
	if err := c.writeCacheEntry(metaPath, entry); err != nil {
		return fmt.Errorf("failed to write cache metadata: %w", err)
	}

	// Save SVG data
	if len(avatar.SVGData) > 0 {
		svgPath := c.getSVGPath(cacheKey)
		if err := ioutil.WriteFile(svgPath, avatar.SVGData, 0644); err != nil {
			return fmt.Errorf("failed to write SVG data: %w", err)
		}
	}

	// Save ASCII data
	if avatar.ASCIIArt != "" {
		asciiPath := c.getASCIIPath(cacheKey)
		if err := ioutil.WriteFile(asciiPath, []byte(avatar.ASCIIArt), 0644); err != nil {
			return fmt.Errorf("failed to write ASCII data: %w", err)
		}
	}

	return nil
}

// Delete removes an avatar from cache
func (c *AvatarCache) Delete(cacheKey string) error {
	// Remove all associated files
	paths := []string{
		c.getMetaPath(cacheKey),
		c.getSVGPath(cacheKey),
		c.getASCIIPath(cacheKey),
	}

	var errors []string
	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to delete cache files: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Clear removes all cached avatars
func (c *AvatarCache) Clear() error {
	// Remove all files in cache directory
	return os.RemoveAll(c.cacheDir)
}

// GetStats returns cache statistics
func (c *AvatarCache) GetStats() (CacheStats, error) {
	stats := CacheStats{}

	// Walk through metadata files to collect statistics
	metaDir := filepath.Join(c.cacheDir, "meta")
	err := filepath.Walk(metaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".json") {
			stats.TotalEntries++
			
			entry, err := c.readCacheEntry(path)
			if err != nil {
				return nil // Skip invalid entries
			}

			if time.Now().After(entry.ExpiresAt) {
				stats.ExpiredEntries++
			} else {
				stats.ValidEntries++
			}

			stats.TotalAccessCount += entry.AccessCount
			stats.TotalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return stats, fmt.Errorf("failed to collect cache stats: %w", err)
	}

	// Calculate hit rate (this would need to be tracked separately in a real implementation)
	if stats.TotalAccessCount > 0 {
		stats.HitRate = float64(stats.ValidEntries) / float64(stats.TotalAccessCount)
	}

	return stats, nil
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalEntries     int     `json:"total_entries"`
	ValidEntries     int     `json:"valid_entries"`
	ExpiredEntries   int     `json:"expired_entries"`
	TotalSize        int64   `json:"total_size"`
	TotalAccessCount int     `json:"total_access_count"`
	HitRate          float64 `json:"hit_rate"`
}

// cleanupExpired removes expired cache entries
func (c *AvatarCache) cleanupExpired() {
	metaDir := filepath.Join(c.cacheDir, "meta")
	
	filepath.Walk(metaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			entry, err := c.readCacheEntry(path)
			if err != nil {
				return nil
			}

			if time.Now().After(entry.ExpiresAt) {
				// Extract cache key from filename
				filename := filepath.Base(path)
				cacheKey := strings.TrimSuffix(filename, ".json")
				c.Delete(cacheKey)
			}
		}

		return nil
	})
}

// Helper methods for file paths
func (c *AvatarCache) getMetaPath(cacheKey string) string {
	return filepath.Join(c.cacheDir, "meta", cacheKey+".json")
}

func (c *AvatarCache) getSVGPath(cacheKey string) string {
	return filepath.Join(c.cacheDir, "svg", cacheKey+".svg")
}

func (c *AvatarCache) getASCIIPath(cacheKey string) string {
	return filepath.Join(c.cacheDir, "ascii", cacheKey+".txt")
}

// readCacheEntry reads cache metadata from file
func (c *AvatarCache) readCacheEntry(path string) (*CacheEntry, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// writeCacheEntry writes cache metadata to file
func (c *AvatarCache) writeCacheEntry(path string, entry *CacheEntry) error {
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

// GetCacheSize returns the total size of the cache in bytes
func (c *AvatarCache) GetCacheSize() (int64, error) {
	var size int64
	
	err := filepath.Walk(c.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// CleanupBySize removes least recently used entries until size is under limit
func (c *AvatarCache) CleanupBySize(maxSize int64) error {
	currentSize, err := c.GetCacheSize()
	if err != nil {
		return err
	}

	if currentSize <= maxSize {
		return nil // No cleanup needed
	}

	// Get all cache entries sorted by last access time
	type entryInfo struct {
		cacheKey   string
		lastAccess time.Time
		size       int64
	}

	var entries []entryInfo
	metaDir := filepath.Join(c.cacheDir, "meta")
	
	err = filepath.Walk(metaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if strings.HasSuffix(path, ".json") {
			entry, err := c.readCacheEntry(path)
			if err != nil {
				return nil
			}

			filename := filepath.Base(path)
			cacheKey := strings.TrimSuffix(filename, ".json")
			
			entries = append(entries, entryInfo{
				cacheKey:   cacheKey,
				lastAccess: entry.LastAccess,
				size:       info.Size(),
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Sort by last access time (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].lastAccess.After(entries[j].lastAccess) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove oldest entries until we're under the size limit
	for _, entry := range entries {
		if currentSize <= maxSize {
			break
		}

		if err := c.Delete(entry.cacheKey); err == nil {
			currentSize -= entry.size
		}
	}

	return nil
} 