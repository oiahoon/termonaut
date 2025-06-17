package main

import (
	"fmt"
	"sort"

	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Show command usage by category",
	Long: `Display statistics about your command usage organized by categories
such as git, development, system administration, etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCategoriesCommand(cmd, args)
	},
}

func runCategoriesCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger := setupLogger(cfg.LogLevel)

	// Initialize database
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get all commands
	commands, err := db.GetAllCommands()
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	if len(commands) == 0 {
		fmt.Println("🤷 No commands found yet!")
		fmt.Println("Start using your terminal to see category statistics.")
		return nil
	}

	// Convert to string slice for analysis
	commandStrings := make([]string, len(commands))
	for i, cmd := range commands {
		commandStrings[i] = cmd.Command
	}

	// Analyze categories
	classifier := categories.NewCommandClassifier()
	categoryStats := classifier.AnalyzeCategories(commandStrings)

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		fmt.Printf("%+v\n", categoryStats)
		return nil
	}

	// Format output
	fmt.Println("📊 Command Usage by Category")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(categoryStats) == 0 {
		fmt.Println("\n🤷 No categorized commands found.")
		return nil
	}

	// Sort categories by count (descending)
	type sortableCategoryStats struct {
		category categories.Category
		stats    *categories.CategoryStats
		info     *categories.CategoryInfo
	}

	var sortableStats []sortableCategoryStats
	for category, stats := range categoryStats {
		info := classifier.GetCategoryInfo(category)
		sortableStats = append(sortableStats, sortableCategoryStats{
			category: category,
			stats:    stats,
			info:     info,
		})
	}

	sort.Slice(sortableStats, func(i, j int) bool {
		return sortableStats[i].stats.Count > sortableStats[j].stats.Count
	})

	fmt.Printf("\nTotal Commands Analyzed: %d\n\n", len(commandStrings))

	// Display category breakdown
	for i, item := range sortableStats {
		if i >= 10 { // Limit to top 10 categories
			break
		}

		// Create percentage bar
		barWidth := 30
		filledWidth := int(item.stats.Percentage * float64(barWidth) / 100)
		emptyWidth := barWidth - filledWidth

		bar := ""
		for j := 0; j < filledWidth; j++ {
			bar += "█"
		}
		for j := 0; j < emptyWidth; j++ {
			bar += "░"
		}

		fmt.Printf("%s %s\n", item.info.Icon, item.info.Name)
		fmt.Printf("   [%s] %.1f%%\n", bar, item.stats.Percentage)
		fmt.Printf("   Commands: %d | XP Bonus: %.1fx | Est. XP: %d\n",
			item.stats.Count, item.info.XPBonus, item.stats.TotalXP)
		fmt.Printf("   %s\n\n", item.info.Description)
	}

	// Show category mastery
	fmt.Println("🎯 Category Mastery:")
	masteryThresholds := []int{50, 25, 10, 5}
	masteryLevels := []string{"🏆 Master", "⭐ Expert", "💪 Proficient", "🌱 Beginner"}

	categoryCount := 0
	for i, item := range sortableStats {
		if i >= 10 {
			break
		}

		level := "🔰 Novice"
		for j, threshold := range masteryThresholds {
			if item.stats.Count >= threshold {
				level = masteryLevels[j]
				break
			}
		}

		fmt.Printf("  %s %s - %s (%d commands)\n",
			item.info.Icon, item.info.Name, level, item.stats.Count)
		categoryCount++
	}

	fmt.Printf("\n📈 You've mastered %d different command categories!\n", categoryCount)

	return nil
}
