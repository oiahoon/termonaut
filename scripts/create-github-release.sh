#!/bin/bash

# GitHub Release Creation Script
# Creates a GitHub release with built binaries

set -e

VERSION=${1:-"v0.9.0"}
BUILD_DIR="dist"
REPO="oiahoon/termonaut"

echo "🚀 Creating GitHub Release for ${VERSION}..."

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed"
    echo "Install it with: brew install gh"
    echo "Then authenticate with: gh auth login"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub"
    echo "Run: gh auth login"
    exit 1
fi

# Check if release artifacts exist
if [ ! -d "${BUILD_DIR}" ]; then
    echo "❌ Build directory ${BUILD_DIR} not found"
    echo "Run: ./scripts/build-release.sh ${VERSION}"
    exit 1
fi

# Create release notes
RELEASE_NOTES="## 🚀 Termonaut v0.9.0 Release Candidate

### 🎯 What's New

**🎮 Advanced Features & Analytics**
- **Custom Command Scoring** - Create custom XP multipliers for specific commands
- **Advanced Filtering** - Search commands by date, category, exit code, and regex patterns
- **Bulk Operations** - Recalculate XP, update categories, and export data in batch
- **API Server Framework** - REST endpoints for external tool integrations
- **Shell Integration Management** - Multi-shell support with installation/update commands

**🏆 Enhanced Gamification**
- **20+ Achievement Badges** including Shell Sprinter 🏃‍♂️, Git Commander 🧬, Night Coder 🌙
- **Failure Penalty System** - Smart XP penalties based on exit codes and command complexity
- **Easter Egg System** - 13+ contextual triggers with motivational quotes and ASCII art
- **Enhanced XP Calculation** - Complexity bonuses for pipes, redirections, and advanced commands

**🔒 Privacy & Security**
- **Command Sanitization** - Automatic detection and redaction of passwords, tokens, URLs
- **Configurable Privacy** - Granular control over what data is tracked and stored
- **Anonymous Mode** - Complete privacy mode for sensitive environments

**📊 Advanced Analytics**
- **Interactive TUI Dashboard** - Real-time stats with beautiful terminal interface
- **GitHub-style Heatmaps** - Activity visualization across days and months
- **Productivity Insights** - Deep analysis of command patterns and efficiency
- **Category Intelligence** - Automatic command classification with 17 categories

### 📦 Installation

\`\`\`bash
# macOS (Homebrew) - Coming Soon!
brew install termonaut

# Manual Installation (All Platforms)
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash

# Or download binaries directly from this release
\`\`\`

### 🔧 Quick Start

\`\`\`bash
# Install shell integration
termonaut advanced shell install

# View your stats
termonaut stats

# Check achievements
termonaut achievements

# Try the interactive dashboard
termonaut tui

# Test easter egg system
termonaut easter-egg --test
\`\`\`

### 🎮 Current Stats
- **114 commands tracked** across 3 terminal sessions
- **8 achievements earned** including Century, Shell Sprinter, Error Survivor
- **53 unique commands** with automatic categorization
- **Advanced features** all tested and operational

### 📈 Performance
- **< 1ms command logging** overhead
- **SQLite3 with WAL mode** for optimal performance
- **Memory efficient** with smart caching and background processing
- **Zero-dependency** TUI with beautiful ASCII visualizations

---

**Full Changelog**: https://github.com/oiahoon/termonaut/blob/main/CHANGELOG.md"

echo "📝 Creating GitHub release..."

# Create the release
gh release create "${VERSION}" \
    --repo "${REPO}" \
    --title "🚀 Termonaut ${VERSION} - Gamified Terminal Productivity Tracker" \
    --notes "${RELEASE_NOTES}" \
    --draft \
    ${BUILD_DIR}/termonaut-${VERSION}-*.tar.gz \
    ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt

echo "✅ GitHub release created successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Visit https://github.com/${REPO}/releases to review the draft release"
echo "2. Edit the release notes if needed"
echo "3. Publish the release"
echo "4. Submit the Homebrew formula to homebrew-core"
echo ""
echo "🍺 Homebrew submission command:"
echo "brew create https://github.com/${REPO}/releases/download/${VERSION}/termonaut-${VERSION}-darwin-arm64.tar.gz" 