package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/github"
	"github.com/oiahoon/termonaut/internal/stats"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub integration features",
	Long: `GitHub integration features including badges, profiles, and Actions workflows.

This command provides tools to:
- Generate dynamic badges for your GitHub profile
- Create shareable profile summaries
- Set up GitHub Actions workflows
- Export data for social sharing`,
}

// badgesCmd represents the badges command
var badgesCmd = &cobra.Command{
	Use:   "badges",
	Short: "Generate GitHub badges",
	Long: `Generate dynamic GitHub badges for your Termonaut stats.

These badges can be used in your GitHub profile README, project documentation,
or anywhere you want to showcase your terminal productivity.`,
}

// badgesGenerateCmd generates badge URLs
var badgesGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate badge URLs",
	Long:  `Generate Shields.io badge URLs for your current stats.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		logger := logrus.New()
		logger.SetLevel(logrus.WarnLevel) // Reduce noise

		db, err := database.New(config.GetDataDir(cfg), logger)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		statsCalc := stats.New(db)
		badgeConfig := github.DefaultBadgeConfig()
		badgeGen := github.NewBadgeGenerator(statsCalc, badgeConfig)

		// Get user progress
		userProgress, err := db.GetUserProgress()
		if err != nil {
			return fmt.Errorf("failed to get user progress: %w", err)
		}

		// Get basic stats
		basicStats, err := statsCalc.GetBasicStats()
		if err != nil {
			return fmt.Errorf("failed to get basic stats: %w", err)
		}

		// Generate badge URLs with simple color logic
		badges := map[string]string{
			"XP":           badgeGen.GenerateXPBadge(userProgress),
			"Commands":     badgeGen.GenerateCommandsBadge(basicStats.TotalCommands),
			"Streak":       badgeGen.GenerateStreakBadge(userProgress.CurrentStreak),
			"Productivity": badgeGen.GenerateProductivityBadge(0.8),   // TODO: Calculate actual productivity
			"Achievements": badgeGen.GenerateAchievementsBadge(5, 10), // TODO: Get actual achievements
		}

		if userProgress.LastActivityDate != nil {
			badges["Last Active"] = badgeGen.GenerateLastActiveBadge(*userProgress.LastActivityDate)
		}

		// Output format
		output, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		if format == "json" {
			data, err := json.MarshalIndent(badges, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}

			if output != "" {
				return os.WriteFile(output, data, 0644)
			}
			fmt.Println(string(data))
			return nil
		}

		// Default markdown format
		fmt.Println("# Termonaut Badges")
		for label, url := range badges {
			fmt.Printf("![%s](%s)\n", label, url)
		}

		fmt.Println()
		fmt.Println("## Markdown for README:")
		fmt.Println()
		for label, url := range badges {
			fmt.Printf("![%s](%s) ", label, url)
		}
		fmt.Println()

		return nil
	},
}

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Generate shareable profiles",
	Long: `Generate shareable profile summaries and social media snippets.

Create comprehensive profiles that can be shared on GitHub, social media,
or anywhere you want to showcase your terminal productivity.`,
}

// profileGenerateCmd generates a profile
var profileGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a profile",
	Long:  `Generate a comprehensive profile with stats, badges, and achievements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		logger := logrus.New()
		logger.SetLevel(logrus.WarnLevel)

		db, err := database.New(config.GetDataDir(cfg), logger)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		statsCalc := stats.New(db)
		badgeConfig := github.DefaultBadgeConfig()
		badgeGen := github.NewBadgeGenerator(statsCalc, badgeConfig)
		profileGen := github.NewProfileGenerator(statsCalc, badgeGen)

		// Get user progress
		userProgress, err := db.GetUserProgress()
		if err != nil {
			return fmt.Errorf("failed to get user progress: %w", err)
		}

		// Generate profile
		profileData, err := profileGen.GenerateProfile(userProgress)
		if err != nil {
			return fmt.Errorf("failed to generate profile: %w", err)
		}

		output, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		var content []byte
		switch format {
		case "json":
			content, err = profileGen.ExportProfile(profileData, "json")
		default:
			content, err = profileGen.ExportProfile(profileData, "markdown")
		}

		if err != nil {
			return fmt.Errorf("failed to export profile: %w", err)
		}

		if output != "" {
			return os.WriteFile(output, content, 0644)
		}

		fmt.Print(string(content))
		return nil
	},
}

// actionsCmd represents the actions command
var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "GitHub Actions integration",
	Long: `GitHub Actions integration tools for automated stats updates.

Set up GitHub Actions workflows to automatically update your badges,
generate reports, and sync your Termonaut data.`,
}

// actionsListCmd lists available workflow templates
var actionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available workflow templates",
	Long:  `List all available GitHub Actions workflow templates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoOwner := viper.GetString("github.repo_owner")
		repoName := viper.GetString("github.repo_name")

		if repoOwner == "" || repoName == "" {
			fmt.Println("Configure GitHub repository first:")
			fmt.Println("  termonaut config set github.repo_owner your-username")
			fmt.Println("  termonaut config set github.repo_name your-repo")
			return nil
		}

		actionsManager := github.NewActionsManager(repoOwner, repoName)
		templates := actionsManager.GetWorkflowTemplates()

		fmt.Println("Available GitHub Actions workflow templates:")
		fmt.Println()
		for _, template := range templates {
			fmt.Printf("• %s\n  %s\n\n", template.Name, template.Description)
		}

		return nil
	},
}

// actionsGenerateCmd generates a workflow file
var actionsGenerateCmd = &cobra.Command{
	Use:   "generate [template-name]",
	Short: "Generate a workflow file",
	Long:  `Generate a GitHub Actions workflow file from a template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]

		repoOwner := viper.GetString("github.repo_owner")
		repoName := viper.GetString("github.repo_name")

		if repoOwner == "" || repoName == "" {
			return fmt.Errorf("configure GitHub repository first:\n" +
				"  termonaut config set github.repo_owner your-username\n" +
				"  termonaut config set github.repo_name your-repo")
		}

		actionsManager := github.NewActionsManager(repoOwner, repoName)

		workflowContent, err := actionsManager.GenerateWorkflowFile(templateName, nil)
		if err != nil {
			return fmt.Errorf("failed to generate workflow: %w", err)
		}

		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = actionsManager.GetWorkflowFilePath(templateName)
		}

		// Create directory if needed
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(output, []byte(workflowContent), 0644); err != nil {
			return fmt.Errorf("failed to write workflow file: %w", err)
		}

		fmt.Printf("Generated workflow: %s\n", output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(githubCmd)

	// Badges commands
	githubCmd.AddCommand(badgesCmd)
	badgesCmd.AddCommand(badgesGenerateCmd)

	badgesGenerateCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	badgesGenerateCmd.Flags().StringP("format", "f", "markdown", "Output format (markdown, json)")

	// Profile commands
	githubCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileGenerateCmd)

	profileGenerateCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	profileGenerateCmd.Flags().StringP("format", "f", "markdown", "Output format (markdown, json)")

	// Actions commands
	githubCmd.AddCommand(actionsCmd)
	actionsCmd.AddCommand(actionsListCmd)
	actionsCmd.AddCommand(actionsGenerateCmd)

	actionsGenerateCmd.Flags().StringP("output", "o", "", "Output file (default: .github/workflows/[template-name].yml)")
}
