package github

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// SyncManager handles GitHub synchronization
type SyncManager struct {
	config           *config.Config
	statsCalculator  *stats.StatsCalculator
	badgeGenerator   *BadgeGenerator
	profileGenerator *ProfileGenerator
}

// NewSyncManager creates a new sync manager
func NewSyncManager(cfg *config.Config, statsCalculator *stats.StatsCalculator) *SyncManager {
	badgeGenerator := NewBadgeGenerator(statsCalculator, DefaultBadgeConfig())
	profileGenerator := NewProfileGenerator(statsCalculator, badgeGenerator)

	return &SyncManager{
		config:           cfg,
		statsCalculator:  statsCalculator,
		badgeGenerator:   badgeGenerator,
		profileGenerator: profileGenerator,
	}
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Success       bool      `json:"success"`
	Timestamp     time.Time `json:"timestamp"`
	FilesUpdated  []string  `json:"files_updated"`
	CommitHash    string    `json:"commit_hash,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	SyncDuration  string    `json:"sync_duration"`
	BadgesUpdated int       `json:"badges_updated"`
	ProfileSize   int       `json:"profile_size"`
}

// SyncToRepository syncs current stats to the configured GitHub repository
func (sm *SyncManager) SyncToRepository(userProgress *models.UserProgress) (*SyncResult, error) {
	startTime := time.Now()
	result := &SyncResult{
		Timestamp: startTime,
	}

	if !sm.config.SyncEnabled {
		result.ErrorMessage = "GitHub sync is disabled in configuration"
		return result, fmt.Errorf("sync disabled")
	}

	if sm.config.SyncRepo == "" {
		result.ErrorMessage = "No sync repository configured"
		return result, fmt.Errorf("no sync repository configured")
	}

	// Parse repository information
	repoParts := strings.Split(sm.config.SyncRepo, "/")
	if len(repoParts) != 2 {
		result.ErrorMessage = "Invalid repository format. Use: username/repository"
		return result, fmt.Errorf("invalid repository format")
	}

	repoOwner := repoParts[0]
	repoName := repoParts[1]

	// Create temporary directory for sync
	tempDir, err := os.MkdirTemp("", "termonaut-sync-*")
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to create temp directory: %v", err)
		return result, err
	}
	defer os.RemoveAll(tempDir)

	// Clone or update repository
	repoPath := filepath.Join(tempDir, repoName)
	if err := sm.cloneOrUpdateRepo(repoOwner, repoName, repoPath); err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to clone/update repository: %v", err)
		return result, err
	}

	// Generate and save badges
	badgesUpdated, err := sm.generateBadges(userProgress, repoPath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to generate badges: %v", err)
		return result, err
	}
	result.BadgesUpdated = badgesUpdated

	// Generate and save profile
	profileSize, err := sm.generateProfile(userProgress, repoPath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to generate profile: %v", err)
		return result, err
	}
	result.ProfileSize = profileSize

	// Generate heatmap
	if err := sm.generateHeatmap(repoPath); err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to generate heatmap: %v", err)
		return result, err
	}

	// Commit and push changes
	commitHash, filesUpdated, err := sm.commitAndPush(repoPath)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to commit and push: %v", err)
		return result, err
	}

	// Update result
	result.Success = true
	result.CommitHash = commitHash
	result.FilesUpdated = filesUpdated
	result.SyncDuration = time.Since(startTime).String()

	return result, nil
}

// TriggerGitHubAction triggers a GitHub Action workflow
func (sm *SyncManager) TriggerGitHubAction(workflowName string) error {
	if sm.config.SyncRepo == "" {
		return fmt.Errorf("no sync repository configured")
	}

	repoParts := strings.Split(sm.config.SyncRepo, "/")
	if len(repoParts) != 2 {
		return fmt.Errorf("invalid repository format")
	}

	// Use GitHub CLI to trigger workflow
	cmd := exec.Command("gh", "workflow", "run", workflowName,
		"--repo", sm.config.SyncRepo)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to trigger workflow: %v, output: %s", err, string(output))
	}

	return nil
}

// ScheduleSync schedules automatic sync based on configuration
func (sm *SyncManager) ScheduleSync(userProgress *models.UserProgress) error {
	if !sm.config.SyncEnabled {
		return nil
	}

	// Check if it's time to sync based on frequency
	lastSyncFile := filepath.Join(sm.config.DataDir, "last_sync.json")
	shouldSync, err := sm.shouldSync(lastSyncFile)
	if err != nil {
		return err
	}

	if !shouldSync {
		return nil
	}

	// Perform sync
	result, err := sm.SyncToRepository(userProgress)
	if err != nil {
		return err
	}

	// Save sync timestamp
	if err := sm.saveLastSync(lastSyncFile, result); err != nil {
		return err
	}

	return nil
}

// cloneOrUpdateRepo clones or updates the repository
func (sm *SyncManager) cloneOrUpdateRepo(owner, name, path string) error {
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, name)

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Clone repository
		cmd := exec.Command("git", "clone", repoURL, path)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to clone repository: %v, output: %s", err, string(output))
		}
	} else {
		// Update existing repository
		cmd := exec.Command("git", "-C", path, "pull", "origin", "main")
		if output, err := cmd.CombinedOutput(); err != nil {
			// Try master branch if main fails
			cmd = exec.Command("git", "-C", path, "pull", "origin", "master")
			if _, err2 := cmd.CombinedOutput(); err2 != nil {
				return fmt.Errorf("failed to update repository: %v, output: %s", err, string(output))
			}
		}
	}

	return nil
}

// generateBadges generates and saves badge files
func (sm *SyncManager) generateBadges(userProgress *models.UserProgress, repoPath string) (int, error) {
	badgesDir := filepath.Join(repoPath, "badges")
	if err := os.MkdirAll(badgesDir, 0755); err != nil {
		return 0, err
	}

	basicStats, err := sm.statsCalculator.GetBasicStats()
	if err != nil {
		return 0, err
	}

	badgeCount := 0

	// Commands badge
	commandsBadgeJSON := sm.generateBadgeJSON("Commands", fmt.Sprintf("%d", basicStats.TotalCommands), sm.getCommandsColor(basicStats.TotalCommands))
	if err := os.WriteFile(filepath.Join(badgesDir, "commands.json"), []byte(commandsBadgeJSON), 0644); err != nil {
		return badgeCount, err
	}
	badgeCount++

	// Level badge
	levelBadgeJSON := sm.generateBadgeJSON("Level", fmt.Sprintf("%d", userProgress.CurrentLevel), sm.getLevelColor(userProgress.CurrentLevel))
	if err := os.WriteFile(filepath.Join(badgesDir, "level.json"), []byte(levelBadgeJSON), 0644); err != nil {
		return badgeCount, err
	}
	badgeCount++

	// Streak badge
	streakBadgeJSON := sm.generateBadgeJSON("Streak", fmt.Sprintf("%d days", userProgress.CurrentStreak), sm.getStreakColor(userProgress.CurrentStreak))
	if err := os.WriteFile(filepath.Join(badgesDir, "streak.json"), []byte(streakBadgeJSON), 0644); err != nil {
		return badgeCount, err
	}
	badgeCount++

	// Productivity badge
	productivityScore := 80.0 // Placeholder calculation
	productivityBadgeJSON := sm.generateBadgeJSON("Productivity", fmt.Sprintf("%.1f%%", productivityScore), sm.getProductivityColor(productivityScore))
	if err := os.WriteFile(filepath.Join(badgesDir, "productivity.json"), []byte(productivityBadgeJSON), 0644); err != nil {
		return badgeCount, err
	}
	badgeCount++

	// Last active badge
	if userProgress.LastActivityDate != nil {
		lastActive := time.Since(*userProgress.LastActivityDate)
		lastActiveBadgeJSON := sm.generateBadgeJSON("Last Active", sm.formatDuration(lastActive), sm.getLastActiveColor(lastActive))
		if err := os.WriteFile(filepath.Join(badgesDir, "last-active.json"), []byte(lastActiveBadgeJSON), 0644); err != nil {
			return badgeCount, err
		}
		badgeCount++
	}

	return badgeCount, nil
}

// generateProfile generates and saves profile markdown
func (sm *SyncManager) generateProfile(userProgress *models.UserProgress, repoPath string) (int, error) {
	profile, err := sm.profileGenerator.GenerateProfile(userProgress)
	if err != nil {
		return 0, err
	}

	profilePath := filepath.Join(repoPath, "TERMONAUT_PROFILE.md")
	if err := os.WriteFile(profilePath, []byte(profile.ProfileMarkdown), 0644); err != nil {
		return 0, err
	}

	return len(profile.ProfileMarkdown), nil
}

// generateHeatmap generates and saves heatmap files
func (sm *SyncManager) generateHeatmap(repoPath string) error {
	heatmapDir := filepath.Join(repoPath, "heatmap")
	if err := os.MkdirAll(heatmapDir, 0755); err != nil {
		return err
	}

	// Generate current year heatmap
	currentYear := time.Now().Year()
	heatmapGenerator := NewHeatmapGenerator(sm.statsCalculator)

	heatmapData, err := heatmapGenerator.GenerateYearHeatmap(currentYear)
	if err != nil {
		return err
	}

	// Save as markdown
	markdownHeatmap := heatmapGenerator.GenerateMarkdownHeatmap(heatmapData)
	if err := os.WriteFile(filepath.Join(heatmapDir, fmt.Sprintf("%d.md", currentYear)), []byte(markdownHeatmap), 0644); err != nil {
		return err
	}

	// Save as HTML
	htmlHeatmap := heatmapGenerator.GenerateHTMLHeatmap(heatmapData)
	if err := os.WriteFile(filepath.Join(heatmapDir, fmt.Sprintf("%d.html", currentYear)), []byte(htmlHeatmap), 0644); err != nil {
		return err
	}

	// Save as SVG
	svgHeatmap := heatmapGenerator.GenerateSVGHeatmap(heatmapData)
	if err := os.WriteFile(filepath.Join(heatmapDir, fmt.Sprintf("%d.svg", currentYear)), []byte(svgHeatmap), 0644); err != nil {
		return err
	}

	return nil
}

// commitAndPush commits changes and pushes to repository
func (sm *SyncManager) commitAndPush(repoPath string) (string, []string, error) {
	// Configure git user
	exec.Command("git", "-C", repoPath, "config", "user.email", "termonaut@github.com").Run()
	exec.Command("git", "-C", repoPath, "config", "user.name", "Termonaut Bot").Run()

	// Add all changes
	cmd := exec.Command("git", "-C", repoPath, "add", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("failed to add files: %v, output: %s", err, string(output))
	}

	// Check if there are changes to commit
	cmd = exec.Command("git", "-C", repoPath, "diff", "--staged", "--quiet")
	if err := cmd.Run(); err == nil {
		// No changes to commit
		return "", []string{}, nil
	}

	// Get list of changed files
	cmd = exec.Command("git", "-C", repoPath, "diff", "--staged", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get changed files: %v", err)
	}
	filesUpdated := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Commit changes
	commitMessage := fmt.Sprintf("🚀 Update Termonaut stats - %s", time.Now().Format("2006-01-02 15:04:05"))
	cmd = exec.Command("git", "-C", repoPath, "commit", "-m", commitMessage)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", filesUpdated, fmt.Errorf("failed to commit: %v, output: %s", err, string(output))
	}

	// Get commit hash
	cmd = exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	output, err = cmd.Output()
	if err != nil {
		return "", filesUpdated, fmt.Errorf("failed to get commit hash: %v", err)
	}
	commitHash := strings.TrimSpace(string(output))

	// Push changes
	cmd = exec.Command("git", "-C", repoPath, "push", "origin", "HEAD")
	if output, err := cmd.CombinedOutput(); err != nil {
		return commitHash, filesUpdated, fmt.Errorf("failed to push: %v, output: %s", err, string(output))
	}

	return commitHash, filesUpdated, nil
}

// shouldSync determines if sync should happen based on frequency
func (sm *SyncManager) shouldSync(lastSyncFile string) (bool, error) {
	// Check if last sync file exists
	if _, err := os.Stat(lastSyncFile); os.IsNotExist(err) {
		return true, nil // First sync
	}

	// Read last sync data
	data, err := os.ReadFile(lastSyncFile)
	if err != nil {
		return true, nil // Error reading, assume we should sync
	}

	var lastSync SyncResult
	if err := json.Unmarshal(data, &lastSync); err != nil {
		return true, nil // Error parsing, assume we should sync
	}

	// Calculate time since last sync
	timeSinceSync := time.Since(lastSync.Timestamp)

	// Check frequency
	switch sm.config.BadgeUpdateFrequency {
	case "hourly":
		return timeSinceSync >= time.Hour, nil
	case "daily":
		return timeSinceSync >= 24*time.Hour, nil
	case "weekly":
		return timeSinceSync >= 7*24*time.Hour, nil
	default:
		return timeSinceSync >= 24*time.Hour, nil // Default to daily
	}
}

// saveLastSync saves the last sync result
func (sm *SyncManager) saveLastSync(lastSyncFile string, result *SyncResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lastSyncFile, data, 0644)
}

// Helper functions for badge generation
func (sm *SyncManager) generateBadgeJSON(label, message, color string) string {
	return fmt.Sprintf(`{
  "schemaVersion": 1,
  "label": "%s",
  "message": "%s",
  "color": "%s",
  "namedLogo": "terminal",
  "logoColor": "white"
}`, label, message, color)
}

func (sm *SyncManager) getCommandsColor(commands int) string {
	if commands >= 1000 {
		return "brightgreen"
	} else if commands >= 500 {
		return "green"
	} else if commands >= 100 {
		return "yellow"
	}
	return "lightgrey"
}

func (sm *SyncManager) getLevelColor(level int) string {
	if level >= 25 {
		return "purple"
	} else if level >= 10 {
		return "blue"
	} else if level >= 5 {
		return "green"
	}
	return "lightgrey"
}

func (sm *SyncManager) getStreakColor(streak int) string {
	if streak >= 30 {
		return "purple"
	} else if streak >= 7 {
		return "green"
	} else if streak >= 3 {
		return "yellow"
	}
	return "red"
}

func (sm *SyncManager) getProductivityColor(score float64) string {
	if score >= 80 {
		return "brightgreen"
	} else if score >= 60 {
		return "green"
	} else if score >= 40 {
		return "yellow"
	}
	return "red"
}

func (sm *SyncManager) getLastActiveColor(duration time.Duration) string {
	if duration < time.Hour {
		return "brightgreen"
	} else if duration < 24*time.Hour {
		return "green"
	} else if duration < 7*24*time.Hour {
		return "yellow"
	}
	return "red"
}

func (sm *SyncManager) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	} else if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	} else {
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
