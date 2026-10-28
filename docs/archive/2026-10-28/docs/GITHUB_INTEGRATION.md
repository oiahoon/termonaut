# 🔗 GitHub Integration Guide

Termonaut v0.9.0 includes comprehensive GitHub integration features for showcasing your terminal productivity on your GitHub profile and social media.

## 🚀 Features Overview

### ✅ Available Features
- **Dynamic Badges**: Generate real-time badges showing your stats
- **Profile Generation**: Create comprehensive profile markdown
- **GitHub Actions Templates**: Automate updates with workflows
- **Export Functionality**: Save badges and profiles in multiple formats
- **Configuration Management**: Easy setup and management

### 🔄 Sync Features (Coming Soon)
- **Automatic Repository Sync**: Push stats to your GitHub repository
- **Scheduled Updates**: Automatic badge and profile updates
- **GitHub Actions Triggers**: Manual and automatic workflow triggers

## 📊 Badge Generation

### Basic Usage
```bash
# Generate badges in URL format (default)
tn github badges generate

# Generate in JSON format
tn github badges generate --format json

# Generate in Markdown format
tn github badges generate --format markdown

# Save to file
tn github badges generate --format json --output badges.json
```

### Available Badges
- **Commands**: Total commands executed
- **Level**: Current experience level
- **Streak**: Current activity streak
- **XP**: Experience points and level progress
- **Productivity**: Productivity score (placeholder)
- **Achievements**: Achievement progress
- **Last Active**: Time since last activity

### Example Output
```markdown
![Commands](https://img.shields.io/badge/Commands-124-green?style=flat-square&logo=terminal&logoColor=white)
![Level](https://img.shields.io/badge/Level-4-lightgrey?style=flat-square&logo=terminal&logoColor=white)
![Streak](https://img.shields.io/badge/Streak-2+days-red?style=flat-square&logo=terminal&logoColor=white)
```

## 📄 Profile Generation

### Basic Usage
```bash
# Generate profile in Markdown (default)
tn github profile generate

# Generate in JSON format
tn github profile generate --format json

# Save to file
tn github profile generate --output TERMONAUT_PROFILE.md
```

### Profile Sections
- **Stats Badges**: Visual representation of your stats
- **Overview**: Key metrics and achievements
- **Achievements**: Completed and in-progress achievements
- **Top Commands**: Most frequently used commands
- **Activity Summary**: Recent activity and trends

## 🤖 GitHub Actions Integration

### Available Workflows

#### 1. Stats Update Workflow
Automatically updates badges and stats every 6 hours.

```bash
tn github actions generate termonaut-stats-update
```

#### 2. Profile Sync Workflow
Syncs profile data to your repository (manual trigger).

```bash
tn github actions generate termonaut-profile-sync
```

#### 3. Weekly Report Workflow
Generates weekly productivity reports.

```bash
tn github actions generate termonaut-weekly-report
```

### Workflow Setup
1. Generate the workflow file
2. Commit to your repository's `.github/workflows/` directory
3. Configure any required secrets
4. The workflow will run based on its triggers

## ⚙️ Configuration

### GitHub Sync Configuration
```bash
# Enable GitHub sync
tn config set sync_enabled true

# Set your repository (username/repository)
tn config set sync_repo your-username/your-repository

# Set update frequency (hourly, daily, weekly)
tn config set badge_update_frequency daily
```

### Check Configuration
```bash
# View current sync settings
tn config get sync_enabled
tn config get sync_repo
tn config get badge_update_frequency
```

## 📋 Setup Guide

### 1. Repository Setup

#### Option A: Profile Repository
Create a repository with the same name as your username (e.g., `username/username`) for your GitHub profile README.

#### Option B: Dedicated Stats Repository
Create a separate repository for your Termonaut stats (e.g., `username/termonaut-stats`).

### 2. Basic Integration

1. **Configure Termonaut**:
   ```bash
   tn config set sync_enabled true
   tn config set sync_repo your-username/your-repository
   ```

2. **Generate badges**:
   ```bash
   tn github badges generate --format markdown
   ```

3. **Copy badges** to your README.md

4. **Generate profile**:
   ```bash
   tn github profile generate --output TERMONAUT_PROFILE.md
   ```

### 3. Advanced Integration with GitHub Actions

1. **Generate workflow**:
   ```bash
   tn github actions generate termonaut-stats-update
   ```

2. **Commit workflow** to your repository:
   ```bash
   git add .github/workflows/termonaut-stats-update.yml
   git commit -m "Add Termonaut stats automation"
   git push
   ```

3. **Configure repository** (if needed):
   - Go to your repository settings
   - Add any required secrets
   - Enable GitHub Actions if not already enabled

## 🎨 Customization

### Badge Styling
All badges use the Shields.io service with these default settings:
- Style: `flat-square`
- Logo: `terminal`
- Logo Color: `white`
- Dynamic colors based on values

### Profile Customization
The generated profile includes:
- Customizable sections
- Emoji indicators
- Progress bars for visual appeal
- Responsive design for different screen sizes

## 📱 Social Media Integration

### Sharing Your Stats
```bash
# Generate shareable profile
tn github profile generate --format markdown

# Generate badges for social media
tn github badges generate --format url
```

### Platform-Specific Tips

#### GitHub Profile README
- Use the full profile markdown
- Include badges at the top
- Add custom sections as needed

#### Twitter/X
- Share individual badges
- Use the stats summary for tweets
- Include screenshots of your terminal

#### LinkedIn
- Share productivity insights
- Use the overview section
- Highlight professional development

## 🔧 Troubleshooting

### Common Issues

#### Badges Not Updating
- Check if GitHub Actions workflow is enabled
- Verify repository permissions
- Ensure correct repository configuration

#### Profile Generation Errors
- Check database permissions
- Verify Termonaut is properly initialized
- Run `tn stats` to verify data availability

#### Sync Configuration Issues
- Verify repository format: `username/repository`
- Check if repository exists and is accessible
- Ensure sync is enabled: `tn config get sync_enabled`

### Debug Commands
```bash
# Check current configuration
tn config get sync_enabled
tn config get sync_repo

# Verify stats availability
tn stats

# Test badge generation
tn github badges generate --format json
```

## 🔮 Future Features

### Coming Soon
- **Real-time Sync**: Automatic repository updates
- **Heatmap Generation**: GitHub-style activity heatmaps
- **Advanced Analytics**: Detailed productivity insights
- **Custom Themes**: Personalized badge and profile styles
- **API Integration**: External service integration
- **Social Snippets**: Platform-specific sharing formats

### Planned Enhancements
- Interactive dashboards
- Team collaboration features
- Advanced achievement system
- Machine learning insights
- Custom badge creation
- Webhook integrations

## 📚 Examples

Check out these example files:
- `examples/exports/badges.json` - Badge data in JSON format
- `examples/exports/profile.md` - Generated profile markdown
- `.github/workflows/` - Example workflow files

## 🤝 Contributing

Want to improve GitHub integration? Check out:
- `internal/github/` - Core GitHub integration code
- `cmd/termonaut/main.go` - CLI command implementations
- `examples/` - Example scripts and demos

## 📞 Support

Need help with GitHub integration?
- Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
- Review the [API Documentation](API.md)
- Open an issue on GitHub
- Join our community discussions

---

*Happy coding and sharing your terminal productivity! 🚀*