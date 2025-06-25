#!/bin/bash

# Termonaut v0.9.0 Official Release Script
# Final stable release with comprehensive GitHub integration

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Version information
VERSION="0.9.0"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Termonaut v${VERSION} Official Release${NC}"
echo -e "════════════════════════════════════════════════"
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Date: ${GREEN}${BUILD_DATE}${NC}"
echo -e "Commit: ${GREEN}${COMMIT_SHA}${NC}"
echo
echo -e "${YELLOW}🌟 Official Release Features:${NC}"
echo -e "  • ⚡ Complete gamified terminal productivity tracker"
echo -e "  • 🎮 XP system with 20+ achievements and Easter eggs"
echo -e "  • 📊 Advanced analytics and TUI dashboard"
echo -e "  • 🔗 Comprehensive GitHub integration with badges"
echo -e "  • 🔇 Zero-noise shell integration (no job control messages)"
echo -e "  • 🛡️ Privacy-first command sanitization"
echo -e "  • 🎨 Beautiful themes and customizable display modes"
echo

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -d "cmd/termonaut" ]]; then
    echo -e "${RED}❌ Error: Must run from termonaut project root${NC}"
    exit 1
fi

# Clean previous builds
echo -e "${BLUE}🧹 Cleaning previous builds...${NC}"
rm -rf dist/
mkdir -p dist/

# Verify tests pass
echo -e "${BLUE}🧪 Running tests...${NC}"
if ! go test ./tests/unit/ -v; then
    echo -e "${RED}❌ Tests failed! Aborting release.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ All tests passed${NC}"

# Build for multiple platforms
echo -e "${BLUE}🔨 Building binaries...${NC}"

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.date=${BUILD_DATE}"

# Detect current platform and architecture
CURRENT_OS=$(go env GOOS)
CURRENT_ARCH=$(go env GOARCH)

echo -e "Building for current platform (${CURRENT_OS}-${CURRENT_ARCH})..."
CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-${CURRENT_OS}-${CURRENT_ARCH} cmd/termonaut/*.go

# macOS (Intel) - only if we're on macOS
if [[ "$CURRENT_OS" == "darwin" ]]; then
    echo -e "Building for macOS (Intel)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-amd64 cmd/termonaut/*.go

    echo -e "Building for macOS (Apple Silicon)..."
    CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-arm64 cmd/termonaut/*.go
fi

echo -e "${GREEN}✅ All binaries built successfully${NC}"

# Test the local binary
echo -e "${BLUE}🧪 Testing local binary...${NC}"
LOCAL_BINARY="dist/termonaut-${VERSION}-${CURRENT_OS}-${CURRENT_ARCH}"
if [[ "$CURRENT_OS" == "windows" ]]; then
    LOCAL_BINARY="${LOCAL_BINARY}.exe"
fi

if [[ -f "$LOCAL_BINARY" ]]; then
    echo -e "Testing version command..."
    $LOCAL_BINARY version
    echo -e "Testing GitHub integration..."
    $LOCAL_BINARY github --help
    echo -e "Testing short alias support..."
    $LOCAL_BINARY --help | grep -A1 "Aliases:"
    echo -e "${GREEN}✅ Binary test passed${NC}"
else
    echo -e "${RED}❌ Local binary not found: $LOCAL_BINARY${NC}"
    exit 1
fi

# Create checksums
echo -e "${BLUE}🔐 Creating checksums...${NC}"
cd dist/
for file in termonaut-${VERSION}-*; do
    if [[ -f "$file" ]]; then
        sha256sum "$file" >> termonaut-${VERSION}-checksums.txt 2>/dev/null || shasum -a 256 "$file" >> termonaut-${VERSION}-checksums.txt
    fi
done
cd ..
echo -e "${GREEN}✅ Checksums created${NC}"

# Create release notes
echo -e "${BLUE}📝 Creating release notes...${NC}"
cat > dist/RELEASE_NOTES_${VERSION}.md << 'EOF'
# Termonaut v0.9.0 - Official Release 🚀

## 🎉 Stable Release Highlights

Termonaut v0.9.0 is the first stable release of the gamified terminal productivity tracker! This release includes all features from the RC versions plus comprehensive GitHub integration.

### 🌟 Core Features

#### 🎮 Complete Gamification System
- **XP & Leveling**: Space-themed progression from Cadet to Cosmic Commander
- **20+ Achievements**: From Shell Sprinter to Docker Whale, unlock achievements based on your terminal habits
- **Easter Eggs**: 13+ contextual surprises for git commits, coffee breaks, and more
- **Command Categories**: 17 categories with smart auto-classification

#### 📊 Advanced Analytics & Insights
- **TUI Dashboard**: Beautiful terminal interface with real-time stats
- **Productivity Analytics**: Deep insights into your terminal usage patterns
- **Heatmaps**: Visual activity tracking like GitHub contribution graphs
- **Privacy-First**: Comprehensive command sanitization with configurable patterns

#### 🔗 GitHub Integration
- **Dynamic Badges**: Auto-generated shields.io badges for your README
- **Profile Generation**: Complete productivity profiles in Markdown
- **Repository Sync**: Automatic synchronization with your GitHub repos
- **GitHub Actions**: Workflow templates for automated stats updates

#### ⚡ User Experience Excellence
- **Short Alias**: Use `tn` instead of `termonaut` for all commands
- **Zero Noise**: Complete elimination of job control messages
- **Empty Command Stats**: Quick stats display when pressing Enter
- **Shell Completion**: User-friendly setup for bash/zsh completion

### 🛠 Installation

#### Quick Install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash
```

#### Manual Installation
1. Download the appropriate binary for your platform from the release assets
2. Make it executable: `chmod +x termonaut-0.9.0-your-platform`
3. Move to PATH: `sudo mv termonaut-0.9.0-your-platform /usr/local/bin/termonaut`
4. Install shell hooks: `tn advanced shell install`

### 🚀 Quick Start

```bash
# Install shell integration
tn advanced shell install

# View your stats
tn stats

# Open the TUI dashboard
tn tui

# Generate GitHub badges
tn github badges generate

# Create your productivity profile
tn github profile generate
```

### 🔧 Configuration

```bash
# Set your preferred theme
tn config set theme emoji

# Enable empty command stats
tn config set empty_command_stats true

# Configure GitHub integration
tn config set github.username your-username
tn config set github.sync_enabled true
```

### 📖 Documentation

- **README**: Complete setup and usage guide
- **GitHub Integration**: Comprehensive guide for social features
- **Troubleshooting**: Solutions for common issues
- **Contributing**: How to contribute to the project

### 🙏 Acknowledgments

Special thanks to all beta testers and contributors who provided feedback during the RC phases. Your input was invaluable in making this release stable and user-friendly.

### 🐛 Bug Reports & Feature Requests

Please report issues and suggest features on our [GitHub Issues](https://github.com/oiahoon/termonaut/issues) page.

---

**Happy Terminal Productivity! 🚀**
EOF

echo -e "${GREEN}✅ Release notes created${NC}"

# Create archive packages
echo -e "${BLUE}📦 Creating archive packages...${NC}"
cd dist/
for file in termonaut-${VERSION}-*; do
    if [[ -f "$file" && "$file" != *.txt && "$file" != *.md ]]; then
        platform=$(echo "$file" | sed "s/termonaut-${VERSION}-//" | sed 's/\.exe$//')
        if [[ "$file" == *.exe ]]; then
            zip "${file%.exe}.zip" "$file"
        else
            tar -czf "${file}.tar.gz" "$file"
        fi
    fi
done
cd ..
echo -e "${GREEN}✅ Archive packages created${NC}"

# Display build summary
echo
echo -e "${GREEN}🎉 Termonaut v${VERSION} Official Release Build Complete!${NC}"
echo -e "═══════════════════════════════════════════════════════════"
echo -e "📁 Build artifacts in: ${BLUE}dist/${NC}"
echo -e "📝 Release notes: ${BLUE}dist/RELEASE_NOTES_${VERSION}.md${NC}"
echo -e "🔐 Checksums: ${BLUE}dist/termonaut-${VERSION}-checksums.txt${NC}"
echo
echo -e "${YELLOW}📦 Built binaries:${NC}"
ls -la dist/termonaut-${VERSION}-* | grep -v checksums | grep -v RELEASE_NOTES
echo
echo -e "${YELLOW}🚀 Next steps:${NC}"
echo -e "  1. Test the binaries on different platforms"
echo -e "  2. Create GitHub release with: ${BLUE}scripts/create-github-release.sh${NC}"
echo -e "  3. Update documentation and changelog"
echo -e "  4. Announce the release!"
echo