# 🚀 Termonaut Quick Start Guide

Welcome to Termonaut! This guide will help you get up and running with the gamified terminal productivity tracker in just a few minutes.

## 📦 Installation

### Option 1: Homebrew (Recommended for macOS/Linux)

```bash
# Install from our custom tap
brew tap oiahoon/termonaut
brew install termonaut
```

### Option 2: Quick Install Script

```bash
# One-line installation for all platforms
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

### Option 3: Manual Download

Visit our [GitHub Releases](https://github.com/oiahoon/termonaut/releases) and download the binary for your platform.

## ⚡ First-Time Setup

### 1. Initialize Shell Integration

After installation, set up shell hooks to track your commands:

```bash
# This will add hooks to your ~/.bashrc or ~/.zshrc
termonaut advanced shell install

# Restart your terminal or source your shell config
source ~/.zshrc  # or ~/.bashrc for bash users
```

### 2. Verify Installation

Check that everything is working:

```bash
# Check version
termonaut --version

# View help
termonaut --help

# Check your first stats (should be minimal)
termonaut stats
```

## 🎮 Understanding the Basics

### XP and Leveling System

Termonaut tracks your terminal usage and awards XP (Experience Points) for:
- **Running commands** (base XP)
- **Using different categories** (bonus multipliers)
- **Complex commands** (pipes, redirections, etc.)
- **Successful execution** (no penalty for exit code 0)

### Achievement System

Unlock badges by reaching milestones:
- 🚀 **First Launch** - Your first command
- 🌟 **Explorer** - Use 10 different commands  
- 🏆 **Century** - Execute 100 commands
- 🔥 **Pro Streaker** - 7+ day usage streak
- And many more!

## 📊 Essential Commands

### Check Your Progress

```bash
# View today's statistics
termonaut stats

# Check your current level and XP
termonaut analytics

# See your achievements
termonaut achievements

# Interactive dashboard
termonaut tui
```

### Configuration Commands

```bash
# View current configuration
termonaut config show

# Enable privacy mode
termonaut config set privacy_sanitizer true

# Change theme
termonaut config set theme emoji

# View all available settings
termonaut config --help
```

## ⚙️ Key Configuration Options

Termonaut stores its configuration in `~/.termonaut/config.toml`. Here are the most important settings:

### Display & Theme
```toml
# Display mode for command output
display_mode = "enter"          # Options: off, enter, ps1, floating

# Visual theme
theme = "emoji"                 # Options: minimal, emoji, ascii

# Show gamification features
show_gamification = true
```

### Privacy Settings
```toml
# Enable command sanitization
privacy_sanitizer = true

# Sanitize specific data types
sanitize_passwords = true
sanitize_urls = true
sanitize_file_paths = true

# Commands to completely ignore
opt_out_commands = ["password", "secret", "token"]

# Anonymous mode (strips personal paths)
anonymous_mode = false
```

### Tracking Behavior
```toml
# Session timeout in minutes
idle_timeout_minutes = 10

# Track git repository context
track_git_repos = true

# Enable automatic command categorization
command_categories = true
```

### Easter Eggs
```toml
# Enable fun easter eggs and motivational quotes
easter_eggs_enabled = true
```

## 🎯 Getting Started Workflow

### Day 1: Installation & Setup
1. Install Termonaut using your preferred method
2. Run `termonaut advanced shell install`
3. Restart your terminal
4. Execute a few commands and run `termonaut stats`

### Day 2-7: Explore Features
1. Check your progress daily with `termonaut stats`
2. Try the interactive dashboard: `termonaut tui`
3. Test easter eggs: `termonaut easter-egg --test`
4. Explore command categories: `termonaut categories`

### Week 2+: Advanced Features
1. Set up custom scoring: `termonaut advanced scoring list`
2. Try advanced filtering: `termonaut advanced filter search --help`
3. Generate activity heatmaps: `termonaut heatmap generate`
4. Explore analytics: `termonaut analytics`

## 🔧 Customization Examples

### Example 1: Minimal Setup for Focused Users
```bash
# Disable easter eggs and reduce visual noise
termonaut config set easter_eggs_enabled false
termonaut config set theme minimal
termonaut config set display_mode off
```

### Example 2: Full Gamification Experience
```bash
# Enable all features for maximum engagement
termonaut config set theme emoji
termonaut config set show_gamification true
termonaut config set easter_eggs_enabled true
termonaut config set display_mode enter
```

### Example 3: Privacy-First Configuration
```bash
# Maximum privacy settings
termonaut config set anonymous_mode true
termonaut config set privacy_sanitizer true
termonaut config set sanitize_passwords true
termonaut config set sanitize_urls true
termonaut config set sanitize_file_paths true
```

## 📈 Understanding Your Data

### File Locations
- **Configuration**: `~/.termonaut/config.toml`
- **Database**: `~/.termonaut/termonaut.db`
- **Logs**: `~/.termonaut/termonaut.log`

### Data Export
```bash
# Export all data as JSON
termonaut export --json > my-stats.json

# Export specific time range
termonaut analytics --from 2024-01-01 --to 2024-01-31
```

## 🆘 Troubleshooting

### Shell Integration Not Working
```bash
# Check if hooks are installed
termonaut advanced shell status

# Reinstall hooks
termonaut advanced shell install

# Manual verification
echo $PROMPT_COMMAND  # bash
echo $precmd_functions  # zsh
```

### Performance Issues
```bash
# Check database size
ls -lh ~/.termonaut/termonaut.db

# View recent logs
tail -f ~/.termonaut/termonaut.log
```

### Command Not Tracking
```bash
# Check if command is in opt-out list
termonaut config show | grep opt_out

# Test command tracking
termonaut test-command "echo hello"
```

## 🎉 What's Next?

### Explore Advanced Features
- **TUI Dashboard**: `termonaut tui` for interactive exploration
- **GitHub Integration**: `termonaut github --help` for profile badges
- **Advanced Analytics**: `termonaut analytics` for deep insights
- **Bulk Operations**: `termonaut advanced bulk --help` for data management

### Community & Support
- **Documentation**: Check `/docs` folder for detailed guides
- **GitHub Issues**: Report bugs or request features
- **Development**: See `DEVELOPMENT.md` for contributing

### Pro Tips
1. **Use aliases**: Create aliases for frequently used Termonaut commands
2. **Check achievements regularly**: They provide motivation and track progress
3. **Experiment with themes**: Find the display mode that works for you
4. **Enable privacy mode**: If you work with sensitive data
5. **Explore easter eggs**: They make the experience more fun!

---

## 📋 Quick Reference

| Command | Description |
|---------|-------------|
| `termonaut stats` | View current statistics |
| `termonaut achievements` | Check earned badges |
| `termonaut tui` | Interactive dashboard |
| `termonaut config show` | View configuration |
| `termonaut analytics` | Deep productivity insights |
| `termonaut easter-egg --test` | Test easter egg system |
| `termonaut advanced --help` | Power user features |
| `termonaut --help` | Full command reference |

Happy tracking! 🚀✨

---

*Need help? Check out our [Troubleshooting Guide](TROUBLESHOOTING.md) or open an issue on [GitHub](https://github.com/oiahoon/termonaut/issues).* 