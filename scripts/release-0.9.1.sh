#!/bin/bash

# Termonaut v0.9.1 Release Script
# Avatar System & Intelligent UI Release

set -e

VERSION="v0.9.1"
RELEASE_NAME="Avatar System & Intelligent UI"
BUILD_DIR="dist"

echo "🎨 Termonaut v0.9.1 Release Script"
echo "==================================="
echo "Version: ${VERSION}"
echo "Release: ${RELEASE_NAME}"
echo ""

# Check if we're on main branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "❌ Error: Must be on main branch for release. Currently on: $CURRENT_BRANCH"
    exit 1
fi

# Check if tag exists
if ! git tag -l | grep -q "^${VERSION}$"; then
    echo "❌ Error: Tag ${VERSION} does not exist. Please create it first."
    exit 1
fi

# Check if build directory exists
if [ ! -d "$BUILD_DIR" ]; then
    echo "❌ Error: Build directory not found. Please run build-release.sh first."
    exit 1
fi

echo "✅ Pre-flight checks passed"
echo ""

# Display release artifacts
echo "📦 Release Artifacts:"
ls -la ${BUILD_DIR}/*.tar.gz
echo ""

# Display checksums
echo "🔐 Checksums:"
cat ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt
echo ""

# Create release notes
echo "📝 Creating release notes..."

RELEASE_NOTES=$(cat << 'EOF'
# 🎨 Termonaut v0.9.1: Avatar System & Intelligent UI

## 🌟 Major Features

### 🎨 Complete Avatar System
- **Personalized ASCII Art Avatars** that evolve with your level
- **4 Distinct Styles**: pixel-art, bottts, adventurer, avataaars
- **Smart Sizing**: Auto-adjusts based on terminal width (mini/small/medium/large)
- **Evolution System**: Visual progression every 5 levels
- **High-Quality Rendering**: 24-bit color with optimized character sets

### 🖥️ Intelligent Dashboard Layout
- **Side-by-Side Display**: Avatar on left, stats on right
- **Terminal Width Detection**: Responsive design for any terminal size
- **Rich Statistics**: Progress bars, achievements, productivity insights
- **Professional UI**: Clean borders, separators, and visual hierarchy

### ⚙️ Complete Management System
- **Full CLI Interface**: `show`, `config`, `preview`, `refresh`, `stats` commands
- **Smart Caching**: 7-day TTL with automatic management
- **Easy Configuration**: Style and size preferences
- **Preview System**: See avatars at different levels

## 🚀 Quick Start

```bash
# View your new avatar dashboard
termonaut stats

# Configure your avatar style
termonaut avatar config --style pixel-art

# Preview your avatar at level 25
termonaut avatar preview -l 25

# Show your current avatar
termonaut avatar show
```

## 📦 Installation

### Fresh Installation
```bash
# macOS (Homebrew)
brew tap oiahoon/termonaut
brew install termonaut

# Linux/macOS (curl)
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

### Upgrade from v0.9.0
```bash
# Homebrew users
brew upgrade termonaut

# Manual download
wget https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-<platform>.tar.gz
```

## 🔧 Technical Improvements

- **Enhanced Dependencies**: ascii-image-converter v1.13.1, golang.org/x/term
- **Quality Optimizations**: Size-specific character sets, advanced parameters
- **Performance**: Parallel processing, efficient caching, minimal network requests
- **Error Handling**: Graceful fallbacks to regular stats display

## 🎮 Avatar Styles Preview

Try different styles to find your perfect terminal identity:
- `pixel-art`: Retro 8-bit gaming style (default)
- `bottts`: Modern robot/android theme
- `adventurer`: Fantasy RPG character style
- `avataaars`: Contemporary cartoon style

## 📊 What's New

✨ **New Features**:
- Complete avatar system with 4 styles
- Intelligent side-by-side layout
- Terminal width detection and responsive sizing
- Avatar evolution system
- Rich statistics display with progress bars

🔧 **Technical**:
- 24-bit color support with ANSI codes
- Size-specific character optimization
- Smart caching with 7-day TTL
- Professional UI with visual separators

📖 **Documentation**:
- Updated README with avatar examples
- Complete avatar system specification
- Configuration guide and troubleshooting

---

**🎨 Experience the future of terminal productivity with personalized avatars!**

*Termonaut v0.9.1 - Where productivity meets personality* ✨
EOF
)

echo "$RELEASE_NOTES" > ${BUILD_DIR}/release-notes-${VERSION}.md

echo "✅ Release notes created: ${BUILD_DIR}/release-notes-${VERSION}.md"
echo ""

# Create GitHub release command
echo "🚀 GitHub Release Command:"
echo ""
echo "gh release create ${VERSION} \\"
echo "  --title \"🎨 Termonaut ${VERSION}: ${RELEASE_NAME}\" \\"
echo "  --notes-file ${BUILD_DIR}/release-notes-${VERSION}.md \\"
echo "  --draft \\"
for file in ${BUILD_DIR}/*.tar.gz; do
    echo "  \"$file\" \\"
done
echo "  ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt"
echo ""

echo "📋 Manual Release Steps:"
echo "1. Review the release notes in ${BUILD_DIR}/release-notes-${VERSION}.md"
echo "2. Run the GitHub release command above (or use GitHub web interface)"
echo "3. Update Homebrew formula with new checksums"
echo "4. Test installation on different platforms"
echo "5. Announce the release"
echo ""

echo "🎯 Homebrew Formula Update:"
echo "Update Formula/termonaut.rb with these checksums:"
echo ""
echo "  # macOS Intel"
echo "  sha256 \"$(grep darwin-amd64 ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt | cut -d' ' -f1)\""
echo ""
echo "  # macOS Apple Silicon"  
echo "  sha256 \"$(grep darwin-arm64 ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt | cut -d' ' -f1)\""
echo ""
echo "  # Linux x86_64"
echo "  sha256 \"$(grep linux-amd64 ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt | cut -d' ' -f1)\""
echo ""
echo "  # Linux ARM64"
echo "  sha256 \"$(grep linux-arm64 ${BUILD_DIR}/termonaut-${VERSION}-checksums.txt | cut -d' ' -f1)\""
echo ""

echo "✅ Release preparation complete!"
echo "🎨 Ready to ship Termonaut v0.9.1 with Avatar System!" 