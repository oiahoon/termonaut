# 🚀 Termonaut
*Your Terminal Journey Companion - Track, Gamify, and Level Up Your CLI Productivity*

```
    🚀
   /|\
  / | \
 |  T  |
 |     |
 ||   ||
 /\   /\
```

Termonaut is a lightweight, privacy-focused terminal productivity tracker that gamifies your command-line experience. Transform your daily shell usage into an engaging RPG-like journey with XP, levels, badges, and shareable stats—all without leaving your terminal.

## ✨ Features

### 🆕 **New User Experience** ⭐ *New in Latest Version!*
- **Interactive Setup Wizard**: `termonaut setup` - Guided configuration for new users
- **Quick Start**: `termonaut quickstart` - One-command setup with sensible defaults
- **Smart Onboarding**: Automatic detection of existing installations
- **Permission-Safe Installation**: Intelligent directory selection, no sudo required

### 🎨 **Three-Tier Viewing Modes** ⭐ *Enhanced Architecture!*
- **Smart Mode**: `termonaut tui` - Automatically adapts to your terminal size (default)
- **Compact Mode**: `termonaut tui --mode compact` - Efficient layout for smaller terminals
- **Full Mode**: `termonaut tui --mode full` - Immersive experience for wide terminals
- **Minimal Mode**: `termonaut stats` - Lightning-fast shell output
- **Configurable Defaults**: Set your preferred mode in config file

### 🖼️ **Dynamic Avatar System** ⭐ *Responsive Design!*
- **Adaptive Sizing**: Avatars scale from 8x4 to 70x25 characters based on terminal size
- **Multiple Styles**: Choose from pixel-art, bottts, adventurer, or avataaars themes
- **Real-time Adaptation**: Automatically adjusts when you resize your terminal
- **Evolution System**: Avatar appearance changes and gains new features as you level up
- **Fallback System**: Beautiful default avatars when network is unavailable
- **Alias Management**: `termonaut alias` commands for easy 'tn' shortcut management

### 💡 **Empty Command Stats** ⭐ *Stable Feature!*
- **Instant Stats**: Press Enter on empty command line to see quick stats
- **Dual Modes**: Minimal one-liner or rich multi-line display
- **Smart Integration**: Respects your theme and privacy settings
- **Fully Configurable**: Enable/disable with simple config setting

### 🔍 **Core Tracking**
- **Command Logging**: Automatically track every command execution
- **Session Management**: Intelligent session detection and timing
- **Usage Analytics**: Daily, weekly, and monthly productivity insights
- **Smart Metrics**: Active time estimation, command categories, and streak tracking

### 🎮 **Gamification System** ⭐ *Enhanced!*
- **XP & Levels**: Earn experience points for terminal usage with space-themed progression
- **Achievement Badges**: Unlock 🏅 badges for milestones (20+ achievements available)
- **Easter Eggs**: 22+ contextual surprises for git, docker, kubernetes, AI tools, and special moments
- **New Command Bonuses**: Extra XP for exploring new tools and commands
- **Streak Rewards**: Maintain daily/weekly usage streaks for motivation
- **Category Mastery**: 17 command categories with specialized XP multipliers

### 📊 **Rich CLI Interface**
- **Interactive Stats**: Beautiful terminal-native data visualization
- **Multiple Views**: Session summaries, command breakdowns, and trend analysis
- **Interactive Dashboard**: Modern TUI with Bubble Tea framework (`termonaut tui`)
- **Customizable Display**: ASCII charts, emoji themes, and configurable output
- **Export Options**: JSON/CSV export for backup and integration
- **Short Aliases**: Use `tn` instead of `termonaut` for all commands

### 🔄 **GitHub Integration** ⭐ *Stable in v0.9.2!*
- **Dynamic Badges**: Auto-updating Shields.io badges for your README (6 badge types)
- **Profile Generation**: Complete productivity profiles in Markdown with avatar integration
- **Repository Sync**: Automatic synchronization with your GitHub repos
- **GitHub Actions**: Workflow templates for automated stats updates
- **Heatmap Generation**: GitHub-style activity heatmaps in HTML/SVG/Markdown formats
- **Stats Export**: JSON and Markdown export for social sharing

### 🎭 **Easter Eggs & Fun** ⭐ *Enhanced in v0.9.2!*
- **22+ Trigger Conditions**: Speed runs, coffee breaks, git operations, docker/k8s, AI tools
- **Context-Aware**: Smart detection of programming languages, databases, testing frameworks
- **Motivational Messages**: 30+ unique messages across all categories
- **Modern Terminal Support**: Optimized for Warp, iTerm2, VS Code terminals
- **Configurable Frequency**: Reduced interruption with balanced probability settings

### 🔒 **Privacy & Performance**
- **100% Local**: All data stays on your machine by default
- **Command Sanitization**: Smart detection and redaction of passwords, tokens, URLs
- **Lightweight**: Minimal performance impact with async logging
- **SQLite Storage**: Fast, reliable, and portable data storage
- **No Dependencies**: Single binary with zero external requirements
- **CI Environment Detection**: Automatic quiet mode for 15+ CI platforms

## 🚀 Quick Start

### Installation

**GitHub Install (Recommended):**
```bash
# Install from GitHub releases
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

**Homebrew Installation:**
```bash
# Install from our custom tap
brew tap oiahoon/termonaut
brew install termonaut

# Or install from homebrew-core (coming soon!)
brew install termonaut
```

**Manual Installation:**
```bash
# Download latest release for your platform
# Linux (x64)
wget https://github.com/oiahoon/termonaut/releases/latest/download/termonaut-linux-amd64
chmod +x termonaut-linux-amd64
sudo mv termonaut-linux-amd64 /usr/local/bin/termonaut

# macOS (Intel)
wget https://github.com/oiahoon/termonaut/releases/latest/download/termonaut-darwin-amd64
chmod +x termonaut-darwin-amd64
sudo mv termonaut-darwin-amd64 /usr/local/bin/termonaut

# macOS (Apple Silicon)
wget https://github.com/oiahoon/termonaut/releases/latest/download/termonaut-darwin-arm64
chmod +x termonaut-darwin-arm64
sudo mv termonaut-darwin-arm64 /usr/local/bin/termonaut
```

**Build from Source:**
```bash
git clone https://github.com/oiahoon/termonaut.git
cd termonaut
go build -o termonaut cmd/termonaut/*.go
sudo mv termonaut /usr/local/bin/
```

### Setup

🆕 **New User? Start Here:**

**Interactive Setup (Recommended):**
```bash
termonaut setup
# Guided configuration with explanations
```

**Quick Setup:**
```bash
termonaut quickstart
# One-command setup with sensible defaults
```

**Manual Setup:**
```bash
termonaut init
# Install shell integration manually
```

This automatically adds hooks to your `~/.bashrc` or `~/.zshrc`. Restart your terminal or run:
```bash
source ~/.bashrc  # or ~/.zshrc
```

**💡 Pro Tip**: Use `tn` as a short alias for `termonaut` in all commands!

**📖 Need detailed setup help?** Check our [Quick Start Guide](docs/QUICK_START.md) for step-by-step instructions!

## 📖 Usage

### Basic Commands

**View Your Stats:**
```bash
termonaut stats              # Today's overview (minimal mode)
termonaut stats --weekly     # This week's stats
termonaut stats --monthly    # This month's statistics

# Or use the short alias:
tn stats                     # Today's overview
tn stats --weekly            # This week's stats
tn stats --monthly           # This month's statistics
```

**Interactive Dashboard:**
```bash
termonaut tui                # Smart mode (adapts to terminal size)
termonaut tui --mode compact # Compact mode for smaller terminals
termonaut tui --mode full    # Full mode for wide terminals
termonaut tui --mode minimal # Text-only stats output

# Short commands:
tn tui                       # Smart mode dashboard
tn tui -m compact           # Compact mode
tn tui -m full              # Full mode
```

**Avatar System:**
```bash
# Avatar system is integrated into TUI and stats display
# Avatar management via configuration:
termonaut config set avatar_style pixel-art    # Change avatar style
termonaut config set avatar_size large         # Set avatar size preference

# Short commands:
tn config set avatar_style bottts              # Change to robot style
tn config get avatar                           # View avatar settings
```

**Alias Management:**
```bash
termonaut alias info         # Show alias information and status
termonaut alias check        # Check if 'tn' alias exists
termonaut alias create       # Create 'tn' shortcut manually
termonaut alias remove       # Remove 'tn' alias

# Short commands:
tn alias info                # Show alias information
tn alias check               # Check alias status
```

**Configuration:**
```bash
termonaut config set theme emoji       # Enable emoji theme
termonaut config set gamification true # Toggle XP system
termonaut config get                    # View all settings

# Short commands:
tn config set theme emoji              # Enable emoji theme
tn config set gamification true        # Toggle XP system
tn config get                          # View all settings
```

**Data Management:**
```bash
# Data export/import features are planned for future release
# Current data location: ~/.termonaut/termonaut.db
# Manual backup: cp ~/.termonaut/termonaut.db ~/backup/

# Advanced data operations:
termonaut advanced bulk --help    # Bulk operations on command data
termonaut advanced filter --help  # Advanced filtering and search
```

### Example Output

```bash
$ termonaut stats --today

🚀 Today's Terminal Stats (2024-01-15)
─────────────────────────────────────
Commands Executed: 127 🎯
Active Time: 3h 42m ⏱️
Session Count: 4 📱
New Commands: 3 ⭐
Current Streak: 12 days 🔥

Top Commands:
git (23) ████████████████████████████████
ls (18)  ███████████████████████████
cd (15)  ██████████████████████
vim (12) ████████████████

🎮 Level 8 Astronaut (2,150 XP)
Progress to Level 9: ████████████░░░░ 75%
```

## ⚙️ Configuration

Termonaut uses a TOML configuration file at `~/.termonaut/config.toml`:

```toml
# Display and Theme
display_mode = "enter"          # Options: off, enter, ps1, floating
theme = "emoji"                 # Options: minimal, emoji, ascii
show_gamification = true        # Enable XP and leveling system

# Tracking Behavior
idle_timeout_minutes = 10       # Session timeout
track_git_repos = true          # Include git repository context
command_categories = true       # Categorize commands automatically

# GitHub Integration (Optional)
sync_enabled = false            # Enable GitHub sync
sync_repo = "username/termonaut-profile"
badge_update_frequency = "daily"

# Privacy
opt_out_commands = ["password", "secret"]  # Commands to ignore
anonymous_mode = false          # Strip personal paths from logs
```

## 🎖️ Achievement System

Unlock 20+ badges as you progress:

| Badge | Description | Criteria |
|-------|-------------|----------|
| 🚀 **First Launch** | Welcome aboard! | Execute your first command |
| 🌟 **Explorer** | Command discoverer | Use 50 unique commands |
| 🏆 **Century** | Daily powerhouse | 100 commands in one day |
| 🔥 **Streak Keeper** | Consistency master | 7-day usage streak |
| 👨‍🚀 **Space Commander** | Terminal veteran | Reach level 10 |
| 🪐 **Cosmic Explorer** | Universe navigator | 30-day usage streak |
| ⚡ **Lightning Fast** | Speed demon | 500 commands in one day |
| 🛸 **Master Navigator** | Elite astronaut | Reach level 25 |
| 🏃‍♂️ **Shell Sprinter** | Speed demon | Execute commands rapidly |
| 🧙‍♂️ **Config Whisperer** | Configuration master | Edit config files frequently |
| 🌙 **Night Coder** | Night owl | Code during late hours |
| 🧬 **Git Commander** | Version control expert | Master git operations |
| 🔥 **Pro Streaker** | Consistency champion | Maintain long streaks |
| 🛡️ **Sudo Smasher** | Admin privileges user | Use sudo commands |
| 🐳 **Docker Whale** | Container expert | Work with Docker |
| 🎭 **Vim Escape Artist** | Editor ninja | Master vim commands |
| 💪 **Error Survivor** | Resilience master | Handle command failures |

*View all achievements with `tn tui` and navigate to the Achievements tab!*

## 🔧 Advanced Features

### 🔄 GitHub Integration & Sync

Automatically sync your terminal stats to GitHub and display dynamic badges:

#### Quick Setup
```bash
# Interactive setup (recommended)
termonaut github sync setup

# Or manual configuration:
termonaut config set sync_enabled true
termonaut config set sync_repo username/termonaut-profile
termonaut config set badge_update_frequency daily
```

#### Manual & Automatic Sync
```bash
# Manual sync (immediate)
termonaut github sync now

# Check sync status
termonaut github sync status

# View available commands
termonaut github --help
```

#### Generate Dynamic Badges
```bash
# Generate all badges
termonaut github badges generate

# Generate profile
termonaut github profile generate

# Setup GitHub Actions for automation
termonaut github actions generate termonaut-stats-update
```

#### Add Badges to Your README
```markdown
![Commands](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/commands.json)
![Level](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/level.json)
![Streak](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/streak.json)
![Productivity](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/productivity.json)
![Last Active](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/termonaut-profile/main/badges/last-active.json)
```

#### Sync Features
- **🔄 Automatic Sync**: Background sync based on frequency (hourly/daily/weekly)
- **📊 Profile Generation**: Complete productivity profiles in Markdown
- **🏷️ Dynamic Badges**: 5+ badge types with real-time data
- **📈 Heatmaps**: Visual activity heatmaps (HTML/SVG/Markdown)
- **⚡ GitHub Actions**: Automated workflows for updates
- **🔒 Privacy**: Only syncs aggregated stats, no sensitive data

### Shell Prompt Integration

Add your stats to your shell prompt:
```bash
# Add to ~/.bashrc or ~/.zshrc
export PS1="$(termonaut prompt) $PS1"
```

### API and Integrations

Advanced data operations:
```bash
termonaut advanced analytics --help    # Advanced analytics features
termonaut advanced api --help          # API server for integrations
termonaut github sync --help           # GitHub integration features
```

## 🏗️ Project Structure

```
~/.termonaut/
├── config.toml           # User configuration
├── termonaut.db         # SQLite database
├── termonaut.log        # Application logs
├── cache/
│   ├── export.json      # Latest export data
│   └── badges/          # Generated badge files
└── backups/             # Automatic daily backups
```

## 🚧 Feature Status

### ✅ Fully Implemented
- **Core Tracking**: Command logging, session management, analytics
- **Gamification**: XP system, achievements, easter eggs
- **TUI Interface**: Interactive dashboard with all tabs functional
- **GitHub Integration**: Badges, profiles, sync, GitHub Actions
- **New User Experience**: Setup wizard, quickstart, alias management
- **Advanced Features**: Shell integration, API server, bulk operations

### 🔄 In Development
- **Avatar CLI Commands**: Direct avatar management via command line
- **Data Export/Import**: JSON/CSV export and backup restoration
- **Enhanced Configuration**: More granular config management

### 📋 Planned Features
- **Plugin System**: Extensible architecture for custom features
- **Social Features**: Leaderboards and community sharing
- **Advanced Analytics**: Machine learning insights
- **Mobile Companion**: Stats viewing on mobile devices

## 🤝 Contributing

We welcome contributions! Please see our [Development Guide](DEVELOPMENT.md) for:
- Development setup and workflow
- Architecture overview
- Testing guidelines
- Code style and standards

### Quick Development Setup

```bash
git clone https://github.com/oiahoon/termonaut.git
cd termonaut
make dev-setup    # Install dependencies and dev tools
make test         # Run test suite
make build        # Build binary
```

## 📊 Roadmap

- [x] **v0.1**: Basic command logging and SQLite storage
- [x] **v0.2**: Enhanced stats and session detection
- [x] **v0.5**: Full gamification system and achievements
- [x] **v0.9**: GitHub integration and badge generation
- [x] **v0.9.2**: Avatar system, Easter eggs, Network resilience
- [ ] **v1.0**: Stable release with comprehensive documentation
- [ ] **v1.x**: Social features and leaderboards

See our [Project Planning](PROJECT_PLANNING.md) for detailed milestones.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by terminal productivity tools and gaming mechanics
- Built with ❤️ for the command-line community
- Special thanks to contributors and beta testers

---

**"Hack your habits from the shell."** 🚀

*Transform your terminal from a tool into an adventure. Every command is a step toward mastery.*

For detailed documentation, visit our [Wiki](https://github.com/oiahoon/termonaut/wiki) or run `termonaut help`.