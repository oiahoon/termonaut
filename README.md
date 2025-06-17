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

### 🔍 **Core Tracking**
- **Command Logging**: Automatically track every command execution
- **Session Management**: Intelligent session detection and timing
- **Usage Analytics**: Daily, weekly, and monthly productivity insights
- **Smart Metrics**: Active time estimation, command categories, and streak tracking

### 🎮 **Gamification System**
- **XP & Levels**: Earn experience points for terminal usage with space-themed progression
- **Achievement Badges**: Unlock 🏅 badges for milestones (100 commands, 7-day streaks, etc.)
- **New Command Bonuses**: Extra XP for exploring new tools and commands
- **Streak Rewards**: Maintain daily/weekly usage streaks for motivation

### 📊 **Rich CLI Interface**
- **Interactive Stats**: Beautiful terminal-native data visualization
- **Multiple Views**: Session summaries, command breakdowns, and trend analysis
- **Customizable Display**: ASCII charts, emoji themes, and configurable output
- **Export Options**: JSON/CSV export for backup and integration

### 🔄 **GitHub Integration** (Optional)
- **Dynamic Badges**: Auto-updating Shields.io badges for your README
- **Stats Sync**: GitHub Actions integration for social sharing
- **Profile Enhancement**: Showcase your terminal prowess on your GitHub profile

### 🔒 **Privacy & Performance**
- **100% Local**: All data stays on your machine by default
- **Lightweight**: Minimal performance impact with async logging
- **SQLite Storage**: Fast, reliable, and portable data storage
- **No Dependencies**: Single binary with zero external requirements

## 🚀 Quick Start

### Installation

**GitHub Install (Recommended):**
```bash
# Install from GitHub releases
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
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

**Homebrew Installation (Recommended):**
```bash
# Install from Homebrew (coming soon to homebrew-core!)
brew install termonaut

# Or install from our custom tap (available now)
brew tap oiahoon/termonaut
brew install termonaut
```

### Setup

Initialize Termonaut in your shell:
```bash
termonaut advanced shell install
```

This automatically adds hooks to your `~/.bashrc` or `~/.zshrc`. Restart your terminal or run:
```bash
source ~/.bashrc  # or ~/.zshrc
```

**📖 Need detailed setup help?** Check our [Quick Start Guide](docs/QUICK_START.md) for step-by-step instructions!

## 📖 Usage

### Basic Commands

**View Your Stats:**
```bash
termonaut stats              # Today's overview
termonaut stats --weekly     # This week's stats
termonaut stats --alltime    # Lifetime statistics
```

**Check Your Progress:**
```bash
termonaut xp                 # Current XP and level
termonaut badges             # Earned achievements
termonaut sessions           # Recent terminal sessions
```

**Configuration:**
```bash
termonaut config set theme emoji       # Enable emoji theme
termonaut config set gamification true # Toggle XP system
termonaut config get                    # View all settings
```

**Data Management:**
```bash
termonaut export stats.json  # Export your data
termonaut import backup.json # Restore from backup
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

Unlock badges as you progress:

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

## 🔧 Advanced Features

### GitHub Badge Integration

Display your terminal stats on GitHub:

1. Enable sync in config: `termonaut config set sync_enabled true`
2. Set your repository: `termonaut config set sync_repo username/profile`
3. Add badges to your README:

```markdown
![Commands](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/profile/main/badges/commands.json)
![Level](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/profile/main/badges/level.json)
![Streak](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/username/profile/main/badges/streak.json)
```

### Shell Prompt Integration

Add your stats to your shell prompt:
```bash
# Add to ~/.bashrc or ~/.zshrc
export PS1="$(termonaut prompt) $PS1"
```

### API and Integrations

Export data for external tools:
```bash
termonaut export --json | jq '.stats.total_commands'  # Use with jq
termonaut export --csv > stats.csv                     # Spreadsheet analysis
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
- [ ] **v0.5**: Full gamification system and achievements
- [ ] **v0.9**: GitHub integration and badge generation
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