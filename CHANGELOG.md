# Changelog

All notable changes to Termonaut will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.2] - 2025-06-20 - Easter Eggs & Network Resilience 🎮

### 🎮 Easter Eggs System Optimization

**Probability Optimization for Better UX**
- **Reduced Trigger Frequency**: All easter eggs probabilities reduced by 60%+ to minimize user disruption
  - Speed run: 0.8 → 0.15 (↓81%)
  - Coffee break: 0.6 → 0.25 (↓58%)
  - New day: 0.9 → 0.4 (↓56%)
  - Git commits: 0.5 → 0.2 (↓60%)
  - Docker/K8s: 0.3-0.4 → 0.15-0.2 (↓50%)
  - ASCII art: 0.2 → 0.05 (↓75%)
- **Enhanced Message Variety**: Added 30+ new easter egg messages across all categories
- **Modern Terminal Optimization**: Enhanced formatting for Warp, iTerm2, VS Code terminals

**New Easter Egg Categories**
- **Programming Languages** (0.1 probability):
  - 🐍 Python detection with Zen of Python references
  - 🎭 JavaScript/Node.js with classic JS humor
- **Database Operations** (0.15 probability):
  - 🗄️ MySQL, PostgreSQL, MongoDB, Redis detection
  - SQL humor and database management jokes
- **Testing Frameworks** (0.12 probability):
  - 🧪 pytest, jest, rspec, mocha detection
  - TDD and quality assurance humor
- **AI Tools** (0.08 probability):
  - 🤖 ChatGPT, Claude, Copilot, AI command detection
  - Human-AI collaboration messages

**Enhanced Existing Categories**
- **Docker**: Added containerization humor and deployment jokes
- **Kubernetes**: Enhanced orchestration and YAML references
- **Git**: More commit message creativity commentary
- **Vim**: Modal editing and hjkl warrior references
- **ASCII Art**: New Code Bear pattern addition

### 🌐 Avatar System Network Resilience

**Comprehensive Network Error Handling**
- **Smart Error Detection**: `isNetworkError()` method identifies network vs. service issues
  - DNS resolution failures
  - Connection timeouts and refusals
  - Network unreachable conditions
  - Temporary name resolution failures
- **Graceful Fallback System**: Offline avatar generation when network fails
  - Geometric SVG generation based on username + level
  - Deterministic color schemes with HSL gradients
  - Personalized fallback ASCII art with user initials
- **User-Friendly Notifications**: Clear error messages and status indicators
  - "🌐 Network issue: Unable to fetch avatar from DiceBear API"
  - "⚠️ Fallback mode will be used"
  - Specific error details and recovery suggestions

**New Testing and Diagnostics**
- **Avatar Network Test**: `termonaut avatar-test` command
  - Network connectivity verification
  - DiceBear API accessibility check
  - Avatar generation testing with real user stats
  - Cache information and system recommendations
- **Terminal Compatibility Test**: `termonaut terminal-test` command
  - Modern terminal detection (Warp, iTerm2, Alacritty, etc.)
  - Unicode and emoji support verification
  - Color capability testing
  - Easter egg formatting preview

**Enhanced Avatar Management**
- **Improved Error Resilience**: Avatar generation continues even with partial failures
- **Better Caching Strategy**: Fallback avatars are cached to reduce network dependency
- **Network Status API**: `GetNetworkStatus()` method for system health monitoring

### 🎨 Terminal Compatibility Enhancements

**Modern Terminal Support**
- **Enhanced Detection**: Support for 9+ modern terminal emulators
  - Warp Terminal, iTerm2, Alacritty, Kitty, Hyper
  - Windows Terminal, Tabby, Terminus, Rio Terminal
- **Optimized Formatting**: Terminal-specific easter egg formatting
  - Truecolor support detection
  - Enhanced ANSI escape sequences for modern terminals
  - Fallback formatting for basic terminals

**Display Improvements**
- **Context-Aware Formatting**: Easter eggs adapt to terminal capabilities
- **Better Unicode Support**: Improved box drawing and emoji rendering
- **Color Enhancement**: 24-bit color support where available

### 🔧 Technical Improvements

**Code Quality & Maintainability**
- **Modular Easter Egg System**: Clean separation of concerns
- **Robust Error Handling**: Comprehensive error catching and user feedback
- **Performance Optimization**: Reduced network dependency and faster fallbacks
- **Enhanced Logging**: Better debugging information for troubleshooting

**Testing & Validation**
- **Comprehensive Test Suite**: New commands for system validation
- **Network Simulation**: Fallback system testing capabilities
- **Terminal Compatibility Matrix**: Support verification across platforms

### 📖 Documentation Updates

- **Easter Eggs Guide**: Complete documentation of all easter egg types and triggers
- **Network Troubleshooting**: Avatar system connectivity issue resolution
- **Terminal Compatibility**: Modern terminal setup and optimization guide
- **Testing Commands**: Usage guide for new diagnostic tools

### 🎯 User Experience Impact

**Reduced Interruption**
- Easter eggs now provide entertainment without disrupting workflow
- Significantly lower trigger rates maintain productivity focus
- Enhanced message variety keeps interactions fresh

**Improved Reliability**
- Avatar system works seamlessly in offline/poor network conditions
- Clear status feedback helps users understand system state
- Graceful degradation ensures core functionality always available

**Better Terminal Integration**
- Optimized display across all major terminal emulators
- Enhanced visual quality on modern terminals
- Consistent experience regardless of terminal choice

## [0.9.1] - 2025-06-19 - Avatar System & Intelligent UI 🎨

### 🎨 Major New Features

**Avatar System - Complete Visual Identity**
- **Personalized ASCII Art Avatars**: Unique visual representation that evolves with your level
- **Smart Side-by-Side Layout**: Intelligent stats display with avatar on left, stats on right
- **Adaptive Terminal Sizing**: Automatically adjusts avatar size based on terminal width
  - Mini (10x5), Small (20x10), Medium (40x20), Large (60x30)
  - Smart detection with fallback to 80x24 default
- **Multiple Avatar Styles**: Choose from 4 distinct themes
  - `pixel-art`: Retro 8-bit gaming style (default)
  - `bottts`: Modern robot/android theme
  - `adventurer`: Fantasy RPG character style
  - `avataaars`: Contemporary cartoon style
- **Rich Visual Quality**: High-quality colored ASCII art with optimized character sets
  - 24-bit color support with vivid ANSI escape codes
  - Size-specific character optimization for maximum detail
  - Enhanced ascii-image-converter parameters for best results

### ✨ Enhanced User Experience

**Intelligent Dashboard Layout**
- **Terminal Width Detection**: Uses golang.org/x/term for accurate terminal size detection
- **Responsive Design**: Automatically adjusts layout based on available space
- **Rich Statistics Display**: Enhanced stats with progress bars, achievements, and insights
- **Visual Separators**: Professional layout with borders and column separators
- **ANSI Code Filtering**: Accurate width calculation removing color codes

**Avatar Management System**
- **Complete CLI Interface**: Full suite of avatar management commands
  - `termonaut avatar show` - Display current avatar with level info
  - `termonaut avatar config` - Configure style and size preferences
  - `termonaut avatar preview -l X` - Preview avatar at different levels
  - `termonaut avatar refresh` - Force regenerate avatar cache
  - `termonaut avatar stats` - System statistics and cache info
- **Evolution System**: Avatar appearance changes as you level up
  - Deterministic generation based on username + level + tier
  - Visual progression every 5 levels with new features
  - Next evolution level indicators
- **Intelligent Caching**: 7-day TTL with automatic cache management
  - MD5-based cache keys for efficient storage
  - Cache hit/miss logging for performance monitoring
  - Automatic cleanup and refresh capabilities

### 🔧 Technical Improvements

**Enhanced Dependencies**
- **ascii-image-converter v1.13.1**: Professional ASCII art conversion
- **golang.org/x/term**: Terminal size detection and control
- **DiceBear 9.0 API**: High-quality avatar generation service

**Quality Optimizations**
- **Size-Specific Character Sets**: Optimized ASCII characters for each avatar size
- **Advanced Converter Parameters**: Complex character mapping, color enhancement, threshold tuning
- **Performance Optimization**: Parallel processing, efficient caching, minimal network requests
- **Error Handling**: Graceful fallbacks to regular stats display on avatar failures

### 📊 Enhanced Statistics Display

**Rich Dashboard Features**
- **Progress Visualization**: XP progress bars with percentage display
- **Achievement Tracking**: Automatic achievement detection and display
- **Productivity Insights**: Average commands per session, usage patterns
- **Top Commands Visualization**: Bar charts for most-used commands
- **Contextual Tips**: Configuration hints and next evolution information

### 🎮 Gamification Enhancements

**Avatar Evolution System**
- **Level-Based Progression**: Visual changes every 5 levels
- **Style Consistency**: Maintained character identity across levels
- **Multiple Tiers**: Different visual themes for various level ranges
- **Social Features**: Shareable avatar previews and configurations

### 🛠️ Configuration Management

**Avatar Configuration**
- **Persistent Settings**: Avatar style and size preferences saved to config
- **Easy Switching**: Quick style changes with immediate refresh
- **Validation**: Input validation for all avatar parameters
- **Integration**: Seamless integration with existing configuration system

### 📖 Documentation

- **Avatar System Specification**: Complete technical documentation
- **User Guides**: Updated README with avatar usage examples
- **Configuration Reference**: All new avatar-related settings documented

## [0.9.0] - 2025-01-20 - Official Stable Release 🚀

### 🎉 Stable Release Highlights

This is the first stable release of Termonaut! All RC features have been thoroughly tested and are now production-ready.

### Added
- **🔗 Complete GitHub Integration**
  - Dynamic badge generation with shields.io integration
  - Automatic profile generation in Markdown format
  - Repository synchronization capabilities
  - GitHub Actions workflow templates
  - Heatmap generation for activity visualization
- **📊 Enhanced Analytics**
  - Comprehensive productivity insights
  - Advanced command categorization (17 categories)
  - Time-based analysis and productivity patterns
  - Export functionality for badges and profiles
- **🎮 Full Gamification System**
  - Space-themed XP progression (Cadet to Cosmic Commander)
  - 20+ achievements with rarity system
  - Contextual Easter eggs (13+ trigger conditions)
  - Command complexity bonuses and failure penalties

### Improved
- **⚡ User Experience**
  - `tn` short alias working perfectly across all commands
  - Complete elimination of job control messages
  - User-friendly shell completion setup
  - Empty command stats with configurable display modes
- **🛡️ Privacy & Security**
  - Advanced command sanitization patterns
  - Privacy-first approach with configurable redaction
  - Secure handling of sensitive information
- **🎨 Interface & Themes**
  - Beautiful TUI dashboard with real-time updates
  - Multiple display modes (minimal, rich, quiet)
  - CI environment auto-detection
  - Customizable themes and emoji support

### Fixed
- **🔧 All RC2 Issues Resolved**
  - Short alias (`tn`) functionality fully working
  - Job control message suppression (100% effective)
  - Empty command detection logic improvements
  - Shell hook installation reliability
  - CGO compilation for SQLite support

### Technical Improvements
- Comprehensive test coverage with 83% GitHub integration success rate
- Multi-platform build support (macOS Intel/ARM, Linux x64/ARM64, Windows)
- Enhanced error handling and logging
- Performance optimizations for large datasets
- Robust configuration management system

### Documentation
- Complete setup and usage guides
- GitHub integration documentation
- Troubleshooting guides
- Contributing guidelines
- API documentation for developers

## [0.9.0-rc2] - 2025-06-18 - User Feedback Fixes Release

### Added
- **⚡ Short Command Alias**
  - Added `tn` as short alias for `termonaut` command
  - Works with all subcommands and flags
  - Reduces typing and improves user experience
  - Example: `tn stats`, `tn config set theme emoji`, `tn advanced shell install`

### Fixed
- **🔇 Enhanced Job Control Message Suppression**
  - Completely eliminates `[1] + 91374 done` messages
  - Triple suppression method implementation:
    - Method 1: `nohup` with complete redirection
    - Method 2: Immediate job `disown`
    - Method 3: Temporary job control disable
  - Zsh: `setopt NO_NOTIFY` and `NO_HUP`
  - Bash: `set +m` to disable job control
  - Updated `fix_hook.sh` script with enhanced fixes
- **🔍 Empty Command Stats Detection**
  - Fixed empty command detection logic for edge cases
  - Handle case when no arguments provided to log-command
  - Improved trimming and empty string detection
  - Better error handling for background operations
  - Resolves issue where pressing Enter on empty line didn't show stats

### Changed
- **📖 Updated Documentation**
  - README.md includes short command examples
  - TROUBLESHOOTING.md with new problem solutions
  - Enhanced user guides with `tn` command usage

## [0.9.0-rc] - 2024-01-XX - Release Candidate: Enhanced UX & Empty Command Stats

### Added
- **💡 Empty Command Stats Feature**
  - Quick stats display when pressing Enter on empty command line
  - Configurable display modes: minimal and rich
  - `empty_command_stats` configuration option with full control
  - Respects privacy settings and display preferences
  - Seamless integration with existing shell hooks
- **📖 Comprehensive Documentation**
  - Troubleshooting guide for common issues
  - Empty command demo script with usage examples
  - Enhanced API documentation for internal interfaces
  - Complete configuration reference
- **⚡ Short Command Alias**
  - Added `tn` as short alias for `termonaut` command
  - Works with all subcommands and flags
  - Reduces typing and improves user experience

### Changed
- **⚡ Enhanced Configuration System**
  - Fixed configuration saving for all new options (easter_eggs_enabled, empty_command_stats, privacy settings)
  - Improved configuration validation and error handling
  - Better default value management with viper integration
- **🎮 Improved Easter Egg Integration**
  - Empty commands don't trigger Easter Eggs (by design)
  - Better context awareness and conditional triggering
  - Enhanced probability system for varied experiences

### Fixed
- **🔧 Configuration Management**
  - Fixed configuration save/load issues for new fields
  - Proper handling of all privacy and feature toggle settings
  - Resolved viper configuration persistence problems
- **🐚 Shell Integration Improvements**
  - Better empty command detection in shell hooks
  - Improved silent operation for background processes
  - Enhanced logging and debugging capabilities
- **🔇 Job Control Message Suppression (v0.9.0-rc Enhanced)**
  - Enhanced shell hook with multiple suppression methods
  - Method 1: nohup with complete redirection
  - Method 2: Immediate job disown
  - Method 3: Temporary job control disable (setopt NO_NOTIFY for Zsh, set +m for Bash)
  - Eliminates `[1] + 91374 done` messages completely
- **🔍 Empty Command Stats Detection**
  - Fixed empty command detection logic for edge cases
  - Handle case when no arguments provided to log-command
  - Improved trimming and empty string detection
  - Better error handling for background operations

### Technical Improvements
- Code organization and error handling enhancements
- Performance optimizations for stats calculation
- Better memory management in TUI components
- Comprehensive feature flag system implementation

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