package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "🧹 Clean up temporary files and caches",
	Long: `Clean up temporary files, caches, and old data to free up space.

This command can clean:
• Cache files and temporary data
• Old log files
• Temporary export files
• Old backup files

Options:
  --cache      Clean cache files
  --logs       Clean old log files
  --temp       Clean temporary files
  --old-data   Clean old data (older than 30 days)
  --all        Clean everything
  --dry-run    Show what would be cleaned without actually cleaning

This is safe to run and won't remove important configuration or data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCleanup(cmd, args)
	},
}

var (
	cleanupCache   bool
	cleanupLogs    bool
	cleanupTemp    bool
	cleanupOldData bool
	cleanupAll     bool
	cleanupDryRun  bool
)

func init() {
	cleanupCmd.Flags().BoolVar(&cleanupCache, "cache", false, "Clean cache files")
	cleanupCmd.Flags().BoolVar(&cleanupLogs, "logs", false, "Clean old log files")
	cleanupCmd.Flags().BoolVar(&cleanupTemp, "temp", false, "Clean temporary files")
	cleanupCmd.Flags().BoolVar(&cleanupOldData, "old-data", false, "Clean old data files")
	cleanupCmd.Flags().BoolVar(&cleanupAll, "all", false, "Clean everything")
	cleanupCmd.Flags().BoolVar(&cleanupDryRun, "dry-run", false, "Show what would be cleaned")
}

func runCleanup(cmd *cobra.Command, args []string) error {
	fmt.Println("🧹 Termonaut Cleanup")
	fmt.Println("===================")
	fmt.Println()

	// If --all is specified, enable all cleanup options
	if cleanupAll {
		cleanupCache = true
		cleanupLogs = true
		cleanupTemp = true
		cleanupOldData = true
	}

	// If no specific options, default to safe cleanup
	if !cleanupCache && !cleanupLogs && !cleanupTemp && !cleanupOldData {
		cleanupCache = true
		cleanupLogs = true
		cleanupTemp = true
	}

	configDir := config.GetConfigDir()
	var totalSize int64
	var totalFiles int

	if cleanupDryRun {
		fmt.Println("🔍 Dry run mode - showing what would be cleaned:")
		fmt.Println()
	}

	// Clean cache files
	if cleanupCache {
		size, files, err := cleanCacheFiles(configDir, cleanupDryRun)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to clean cache: %v\n", err)
		} else {
			totalSize += size
			totalFiles += files
			if cleanupDryRun {
				fmt.Printf("📁 Cache files: %d files, %s\n", files, formatSize(size))
			} else {
				fmt.Printf("✅ Cleaned cache files: %d files, %s freed\n", files, formatSize(size))
			}
		}
	}

	// Clean log files
	if cleanupLogs {
		size, files, err := cleanLogFiles(configDir, cleanupDryRun)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to clean logs: %v\n", err)
		} else {
			totalSize += size
			totalFiles += files
			if cleanupDryRun {
				fmt.Printf("📄 Log files: %d files, %s\n", files, formatSize(size))
			} else {
				fmt.Printf("✅ Cleaned log files: %d files, %s freed\n", files, formatSize(size))
			}
		}
	}

	// Clean temporary files
	if cleanupTemp {
		size, files, err := cleanTempFiles(configDir, cleanupDryRun)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to clean temp files: %v\n", err)
		} else {
			totalSize += size
			totalFiles += files
			if cleanupDryRun {
				fmt.Printf("🗂️  Temp files: %d files, %s\n", files, formatSize(size))
			} else {
				fmt.Printf("✅ Cleaned temp files: %d files, %s freed\n", files, formatSize(size))
			}
		}
	}

	// Clean old data
	if cleanupOldData {
		size, files, err := cleanOldData(configDir, cleanupDryRun)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to clean old data: %v\n", err)
		} else {
			totalSize += size
			totalFiles += files
			if cleanupDryRun {
				fmt.Printf("📊 Old data: %d files, %s\n", files, formatSize(size))
			} else {
				fmt.Printf("✅ Cleaned old data: %d files, %s freed\n", files, formatSize(size))
			}
		}
	}

	fmt.Println()
	if cleanupDryRun {
		fmt.Printf("📊 Total that would be cleaned: %d files, %s\n", totalFiles, formatSize(totalSize))
		fmt.Println("Run without --dry-run to actually perform the cleanup.")
	} else {
		fmt.Printf("🎉 Cleanup completed: %d files removed, %s freed\n", totalFiles, formatSize(totalSize))
		if totalSize == 0 {
			fmt.Println("💡 No files needed cleaning - your Termonaut installation is already tidy!")
		}
	}

	return nil
}

func cleanCacheFiles(configDir string, dryRun bool) (int64, int, error) {
	cacheDir := filepath.Join(configDir, "cache")
	return cleanDirectory(cacheDir, []string{"*.tmp", "*.cache", "export.json"}, dryRun)
}

func cleanLogFiles(configDir string, dryRun bool) (int64, int, error) {
	// Clean log files older than 7 days
	cutoffTime := time.Now().AddDate(0, 0, -7)
	return cleanOldFiles(configDir, []string{"*.log"}, cutoffTime, dryRun)
}

func cleanTempFiles(configDir string, dryRun bool) (int64, int, error) {
	patterns := []string{"*.tmp", "*.temp", ".DS_Store", "Thumbs.db"}
	return cleanDirectory(configDir, patterns, dryRun)
}

func cleanOldData(configDir string, dryRun bool) (int64, int, error) {
	// Clean backup files older than 30 days
	cutoffTime := time.Now().AddDate(0, 0, -30)
	backupDir := filepath.Join(configDir, "backups")
	return cleanOldFiles(backupDir, []string{"*.bak", "*.backup"}, cutoffTime, dryRun)
}

func cleanDirectory(dir string, patterns []string, dryRun bool) (int64, int, error) {
	var totalSize int64
	var totalFiles int

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0, 0, nil
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			if !info.IsDir() {
				totalSize += info.Size()
				totalFiles++

				if !dryRun {
					os.Remove(match)
				}
			}
		}
	}

	return totalSize, totalFiles, nil
}

func cleanOldFiles(dir string, patterns []string, cutoffTime time.Time, dryRun bool) (int64, int, error) {
	var totalSize int64
	var totalFiles int

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0, 0, nil
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			// Only remove files older than cutoff time
			if !info.IsDir() && info.ModTime().Before(cutoffTime) {
				totalSize += info.Size()
				totalFiles++

				if !dryRun {
					os.Remove(match)
				}
			}
		}
	}

	return totalSize, totalFiles, nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
