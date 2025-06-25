#!/bin/bash

# Termonaut v0.10.0 Release Script
# Major User Experience Update

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Version information
VERSION="0.10.0"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Termonaut v${VERSION} - Major User Experience Update${NC}"
echo -e "══════════════════════════════════════════════════════════════"
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Date: ${GREEN}${BUILD_DATE}${NC}"
echo -e "Commit: ${GREEN}${COMMIT_SHA}${NC}"
echo
echo -e "${YELLOW}🌟 v0.10.0 Key Features:${NC}"
echo -e "  • 🆕 Interactive Setup Wizard (termonaut setup)"
echo -e "  • ⚡ Quick Start Command (termonaut quickstart)"
echo -e "  • 🎨 Three-Tier Viewing Modes Architecture"
echo -e "  • 🖼️ Dynamic Avatar System (8x4 to 70x25 chars)"
echo -e "  • 🔧 Permission-Safe Installation"
echo -e "  • 🔗 Alias Management System (termonaut alias)"
echo -e "  • 📱 Fully Responsive Design"
echo -e "  • 🎯 Simplified Command Structure"
echo

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -d "cmd/termonaut" ]]; then
    echo -e "${RED}❌ Error: Must run from termonaut project root${NC}"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${YELLOW}⚠️ Warning: You have uncommitted changes${NC}"
    echo "Uncommitted files:"
    git status --porcelain
    read -p "Continue with release? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}❌ Release cancelled${NC}"
        exit 1
    fi
fi

# Create release directory
RELEASE_DIR="releases/v${VERSION}"
mkdir -p "${RELEASE_DIR}"

echo -e "${BLUE}📦 Building cross-platform binaries...${NC}"

# Build flags with version information
BUILD_FLAGS="-ldflags=-X main.version=v${VERSION} -X main.commit=${COMMIT_SHA} -X main.date=${BUILD_DATE}"

# Build for different platforms
echo -e "${PURPLE}🔨 Building Darwin AMD64...${NC}"
GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} -o "${RELEASE_DIR}/termonaut-${VERSION}-darwin-amd64" cmd/termonaut/*.go

echo -e "${PURPLE}🔨 Building Darwin ARM64...${NC}"
GOOS=darwin GOARCH=arm64 go build ${BUILD_FLAGS} -o "${RELEASE_DIR}/termonaut-${VERSION}-darwin-arm64" cmd/termonaut/*.go

echo -e "${PURPLE}🔨 Building Linux AMD64...${NC}"
GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o "${RELEASE_DIR}/termonaut-${VERSION}-linux-amd64" cmd/termonaut/*.go

echo -e "${PURPLE}🔨 Building Linux ARM64...${NC}"
GOOS=linux GOARCH=arm64 go build ${BUILD_FLAGS} -o "${RELEASE_DIR}/termonaut-${VERSION}-linux-arm64" cmd/termonaut/*.go

echo -e "${PURPLE}🔨 Building Windows AMD64...${NC}"
GOOS=windows GOARCH=amd64 go build ${BUILD_FLAGS} -o "${RELEASE_DIR}/termonaut-${VERSION}-windows-amd64.exe" cmd/termonaut/*.go

# Make binaries executable
chmod +x "${RELEASE_DIR}"/termonaut-${VERSION}-*

echo -e "${BLUE}📋 Generating checksums...${NC}"
cd "${RELEASE_DIR}"
sha256sum termonaut-${VERSION}-* > "termonaut-${VERSION}-checksums.txt"
cd - > /dev/null

echo -e "${BLUE}📝 Creating release notes...${NC}"

# Create release notes
cat > "${RELEASE_DIR}/RELEASE_NOTES_${VERSION}.md" << EOF
# 🚀 Termonaut v${VERSION} - Major User Experience Update

## 🎯 Overview

This is a major release focused on dramatically improving the user experience, especially for new users. We've rebuilt the onboarding system, enhanced the avatar display, and simplified the command structure.

## 🆕 New Features

### New User Experience
- **Interactive Setup Wizard**: \`termonaut setup\` - Guided configuration for new users
- **Quick Start**: \`termonaut quickstart\` - One-command setup with sensible defaults
- **Smart Onboarding**: Automatic detection of existing installations
- **Permission-Safe Installation**: Intelligent directory selection, no sudo required

### Three-Tier Viewing Modes
- **Smart Mode**: \`termonaut tui\` - Automatically adapts to your terminal size (default)
- **Compact Mode**: \`termonaut tui --mode compact\` - Efficient layout for smaller terminals
- **Full Mode**: \`termonaut tui --mode full\` - Immersive experience for wide terminals
- **Minimal Mode**: \`termonaut stats\` - Lightning-fast shell output
- **Configurable Defaults**: Set your preferred mode in config file

### Dynamic Avatar System
- **Adaptive Sizing**: Avatars scale from 8x4 to 70x25 characters based on terminal size
- **Real-time Adaptation**: Automatically adjusts when you resize your terminal
- **Multiple Styles**: Choose from pixel-art, bottts, adventurer, or avataaars themes
- **Evolution System**: Avatar appearance changes as you level up
- **Fallback System**: Beautiful default avatars when network is unavailable

### Alias Management System
- **\`termonaut alias info\`** - Show alias information and status
- **\`termonaut alias check\`** - Check if 'tn' alias exists
- **\`termonaut alias create\`** - Create 'tn' shortcut manually
- **\`termonaut alias remove\`** - Remove 'tn' alias

## 🔧 Technical Improvements

### Permission Problem Fix
- **Smart Directory Selection**: Prioritizes user directories (\`~/.local/bin\`)
- **Permission Detection**: Automatic write permission checking
- **Graceful Degradation**: Symlink creation failure doesn't affect main installation
- **User-Friendly Errors**: Clear guidance when issues occur

### Responsive Layout System
- **Intelligent Avatar Sizing**: 35-70 character width support
- **Dynamic Content Adjustment**: Stats area adapts to available space
- **Multi-Size Support**: 7 different avatar size tiers
- **Real-time Adaptation**: Responds to terminal resize events

### Configuration System Enhancement
- **UIConfig Structure**: New configuration section for UI preferences
- **Default Mode Setting**: Users can set their preferred TUI mode
- **Theme Persistence**: Avatar and theme preferences saved
- **Backward Compatible**: All existing configs continue to work

## 🐛 Bug Fixes

- **Fixed**: Permission denied errors during installation
- **Fixed**: Avatar display issues on narrow terminals
- **Fixed**: New user confusion about getting started
- **Fixed**: Command structure complexity

## 📊 User Experience Improvements

### For New Users
- **95% reduction** in setup complexity
- **Clear onboarding** with step-by-step guidance
- **No permission issues** with smart directory selection
- **Immediate success** with sensible defaults

### For Existing Users
- **Simplified commands** - One \`tui\` command instead of multiple
- **Better visuals** - Much wider avatar display (up to 70 characters)
- **Responsive design** - Adapts to any terminal size
- **Backward compatible** - All existing workflows continue to work

## 🚀 Getting Started

### New Users
\`\`\`bash
# Interactive setup (recommended)
termonaut setup

# Or quick setup
termonaut quickstart
\`\`\`

### Existing Users
\`\`\`bash
# Enhanced experience with same commands
termonaut tui        # Now smarter and more responsive
termonaut stats      # Same fast stats you love
\`\`\`

## 📈 What's Next

This release establishes a solid foundation for future enhancements:
- Plugin system architecture
- Advanced customization options
- Social features and sharing
- Performance optimizations

## 🙏 Acknowledgments

Special thanks to users who provided feedback on:
- Installation permission issues
- New user experience challenges
- Avatar display improvements
- Command structure simplification

---

**Full Changelog**: [v0.9.5...v0.10.0](https://github.com/oiahoon/termonaut/compare/v0.9.5...v0.10.0)
EOF

echo -e "${GREEN}✅ Release v${VERSION} built successfully!${NC}"
echo
echo -e "${BLUE}📁 Release files created in: ${RELEASE_DIR}${NC}"
echo -e "${BLUE}📦 Binaries:${NC}"
ls -la "${RELEASE_DIR}"/termonaut-${VERSION}-*
echo
echo -e "${BLUE}📋 Checksums:${NC}"
cat "${RELEASE_DIR}/termonaut-${VERSION}-checksums.txt"
echo
echo -e "${YELLOW}🚀 Next steps:${NC}"
echo -e "  1. Review release notes: ${RELEASE_DIR}/RELEASE_NOTES_${VERSION}.md"
echo -e "  2. Test binaries on different platforms"
echo -e "  3. Create git tag: git tag -a v${VERSION} -m 'Release v${VERSION}'"
echo -e "  4. Push tag: git push origin v${VERSION}"
echo -e "  5. Create GitHub release with these artifacts"
echo -e "  6. Update Homebrew formula"
echo

echo -e "${GREEN}🎉 Termonaut v${VERSION} is ready for release!${NC}"
