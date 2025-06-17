# Changelog

All notable changes to Termonaut will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - v0.9.0 Release Candidate

### Added
- Troubleshooting documentation for common issues
- API documentation for internal interfaces
- Enhanced error handling and validation
- Release preparation and build optimization

### Changed
- Performance optimizations for TUI and CLI commands
- Improved shell hook integration (job control message fix in progress)

### Fixed
- Various UI/UX improvements in TUI interface
- Configuration validation enhancements

## [0.8.0] - 2024-01-XX - Advanced Features & User Enhancements

### Added
- **🔒 Privacy & Command Sanitization System**
  - Comprehensive command sanitization with configurable patterns
  - Smart detection of passwords, tokens, URLs, emails, file paths
  - Privacy-first approach with pattern-based redaction
- **⚡ Enhanced XP System with Failure Penalties**
  - Exit code-based failure penalty calculation
  - Complexity bonuses for pipes, redirections, arguments
  - Category-specific XP adjustments and smart scaling
- **🏆 Extended Achievement System (20+ achievements)**
  - Shell Sprinter 🏃‍♂️, Config Whisperer 🧙‍♂️, Night Coder 🌙
  - Git Commander 🧬, Pro Streaker 🔥, Sudo Smasher 🛡️
  - Docker Whale 🐳, Vim Escape Artist 🎭, Error Survivor 💪
  - Time-based and behavior-based achievement triggers
- **🎭 Comprehensive Easter Egg System**
  - 13+ contextual trigger conditions
  - Speed run, coffee break, morning greeting triggers
  - Git, Docker, Kubernetes, Vim command triggers
  - ASCII art celebrations and motivational quotes
- **🎯 Advanced CLI Commands**
  - `termonaut tui` - Interactive terminal dashboard
  - `termonaut analytics` - Deep productivity insights
  - `termonaut heatmap` - Activity visualization
  - `termonaut dashboard` - Comprehensive overview
  - `termonaut easter-egg` - Test easter egg system
  - `termonaut github` - GitHub integration commands
  - `termonaut categories` - Command categorization view

## [0.7.0] - 2024-01-XX - Performance & Reliability + Enhanced Features

### Added
- **🎲 Randomized Easter Eggs**
  - Context-sensitive easter egg system with 13 trigger conditions
  - Probabilistic trigger system with varied rarity
  - Support for git, docker, kubernetes, vim commands
- **🎨 Display Modes (三种显示模式)**
  - Minimal, Rich, and Quiet modes for different use cases
  - CI environment auto-detection and adaptation
  - Visual progress bars and dynamic emoji selection
- **🤖 CI Environment Auto-Detection**
  - Support for 15+ CI platforms (GitHub Actions, GitLab CI, Jenkins, etc.)
  - Automatic quiet mode activation for CI environments
  - Configurable override options
- **🎮 Enhanced Gaming System**
  - XP Multiplier system with time-based bonuses
  - Power-up system (Double XP, Command Frenzy, etc.)
  - Daily Quest and Weekly Challenge systems
  - Command Rarity system (Common to Legendary)
- **🔥 GitHub Activity Heatmaps**
  - HTML, SVG, and Markdown format generation
  - GitHub-style activity visualization
  - Monthly and yearly statistics breakdown
- **📦 Updated Installation Methods**
  - GitHub-based installation script with multi-platform support
  - Automatic platform detection and version management

## [0.6.0] - 2024-01-XX - GitHub Integration

### Added
- **GitHub Actions Support**
  - Workflow templates for automated stats updates
  - Badge generation system
  - Repository integration capabilities
- **Social Features**
  - Shareable stat summaries and profile generation
  - Dynamic badge creation for GitHub README
  - Social media snippet generation

## [0.5.0] - 2024-01-XX - Beta Release

### Added
- **Achievement System**
  - 17+ core achievements implemented
  - Dynamic achievement unlocking and progress tracking
  - Achievement categories and rarity indicators
- **Data Management**
  - Comprehensive command and session storage
  - Real-time statistics calculation
  - Efficient database operations and data integrity
- **User Customization**
  - Rich theme system (emoji/unicode)
  - Configurable XP rates and category multipliers
  - Flexible display preferences and JSON output

## [0.4.0] - 2024-01-XX - Rich CLI Experience

### Added
- **Enhanced UI/UX**
  - Rich terminal formatting with colors and emojis
  - Interactive progress bars and charts
  - Responsive layout design and beautiful dashboard interfaces
- **Command Categories** ⭐
  - Automatic command classification (17 categories)
  - Category-based statistics and visualization
  - Custom category definitions with XP multipliers
- **Advanced Stats** ⭐
  - Comprehensive productivity analysis engine
  - Time pattern analysis (daily/weekly/hourly)
  - Efficiency metrics and automation suggestions

## [0.3.0] - 2024-01-XX - Gamification Core

### Added
- **XP System**
  - Experience point calculation with bonuses
  - Level progression system with mathematical scaling
  - Level-up notifications and themed titles
  - Progress visualization with Unicode bars
- **Achievement Framework**
  - 17+ predefined achievements across categories
  - Progress tracking and unlock detection
  - Achievement categories and rarity system
  - Real-time achievement notifications

## [0.2.0] - 2024-01-XX - Stats & Display

### Added
- **Statistics Engine**
  - Basic statistics calculation and session management
  - Command counting and frequency analysis
  - Time-based statistics and analysis
- **Display System**
  - Rich terminal formatting and ASCII art
  - Configurable output formats and JSON export

## [0.1.0] - 2024-01-XX - MVP Foundation

### Added
- **Core CLI Framework**
  - Cobra CLI setup and basic command structure
  - Configuration management (TOML-based)
  - Logging infrastructure and version information
- **Database Foundation**
  - SQLite3 with WAL mode for performance
  - Database schema design and migrations
  - Command logging and session tracking
- **Shell Integration**
  - Shell hook system (Zsh/Bash support)
  - Command interception and logging
  - Silent background operation with performance optimization

## Project Status

**Current Phase**: Phase 4 - v0.9.0 Release Candidate
**Next Milestone**: v1.0.0 Stable Release
**Target Release**: v0.9.0 RC

---

## Version History Template

### [X.Y.Z] - YYYY-MM-DD

#### Added
- New features and functionality

#### Changed
- Changes in existing functionality

#### Deprecated
- Soon-to-be removed features

#### Removed
- Features removed in this release

#### Fixed
- Bug fixes

#### Security
- Security vulnerability fixes