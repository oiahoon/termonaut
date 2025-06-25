#!/bin/bash

# Termonaut v0.9.0-rc Release Script
# Builds and prepares release candidate for distribution

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Version information
VERSION="0.9.0-rc"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Termonaut v${VERSION} Release Candidate Build${NC}"
echo -e "════════════════════════════════════════════════════"
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Date: ${GREEN}${BUILD_DATE}${NC}"
echo -e "Commit: ${GREEN}${COMMIT_SHA}${NC}"
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
if ! go test ./... -v; then
    echo -e "${RED}❌ Tests failed! Aborting release.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ All tests passed${NC}"

# Build for multiple platforms
echo -e "${BLUE}🔨 Building binaries...${NC}"

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.date=${BUILD_DATE}"

# macOS (Intel)
echo -e "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-amd64 cmd/termonaut/*.go

# macOS (Apple Silicon)
echo -e "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-arm64 cmd/termonaut/*.go

# Linux (x64)
echo -e "Building for Linux (x64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-linux-amd64 cmd/termonaut/*.go

# Linux (ARM64)
echo -e "Building for Linux (ARM64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-linux-arm64 cmd/termonaut/*.go

# Windows (x64)
echo -e "Building for Windows (x64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-windows-amd64.exe cmd/termonaut/*.go

echo -e "${GREEN}✅ All binaries built successfully${NC}"

# Test the local binary
echo -e "${BLUE}🧪 Testing local binary...${NC}"
LOCAL_BINARY="dist/termonaut-${VERSION}-$(go env GOOS)-$(go env GOARCH)"
if [[ "$(go env GOOS)" == "windows" ]]; then
    LOCAL_BINARY="${LOCAL_BINARY}.exe"
fi

if [[ -f "$LOCAL_BINARY" ]]; then
    echo -e "Testing version command..."
    $LOCAL_BINARY version
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
        sha256sum "$file" >> termonaut-${VERSION}-checksums.txt
    fi
done
cd ..
echo -e "${GREEN}✅ Checksums created${NC}"

# Create release notes
echo -e "${BLUE}📝 Creating release notes...${NC}"
cat > dist/RELEASE_NOTES_${VERSION}.md << EOF
# Termonaut ${VERSION} Release Candidate

## 🌟 What's New in ${VERSION}

### 💡 Empty Command Stats - The Game Changer
The standout feature of this release! Now when you press **Enter** on an empty command line, Termonaut shows you a quick stats summary instead of doing nothing.

**Key Features:**
- **Quick Access**: No more typing \`termonaut stats\` - just hit Enter!
- **Two Display Modes**: 
  - **Minimal**: \`📊 Lv.4 | 168 cmds | 1 streak | 1304 XP\`
  - **Rich**: Multi-line with progress bars, streaks, and top commands
- **Fully Configurable**: Enable/disable with \`termonaut config set empty_command_stats true/false\`
- **Privacy Aware**: Respects all your display and privacy settings

### 🔧 Configuration Improvements
- Fixed configuration saving for all privacy and gamification settings
- Better handling of \`easter_eggs_enabled\`, \`empty_command_stats\`, and privacy options
- Enhanced configuration validation and error handling

### 🎮 Easter Egg Enhancements
- Empty commands won't trigger Easter Eggs (by design - they show stats instead)
- Improved context-awareness for 22+ different Easter Egg conditions
- Better probability system for varied experiences

## 📋 Usage Examples

\`\`\`bash
# Enable the empty command stats feature
termonaut config set empty_command_stats true

# Choose your preferred theme
termonaut config set theme emoji     # Rich display with emojis
termonaut config set theme minimal   # Clean text-only

# Now just press Enter on an empty command line!
# Press Enter → See your stats instantly!

# Disable if you prefer the old behavior
termonaut config set empty_command_stats false
\`\`\`

## 🛠 Installation

### Quick Install (Recommended)
\`\`\`bash
# Download and install latest RC
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash

# Or download specific binary for your platform from the release assets
\`\`\`

### Manual Installation
1. Download the appropriate binary for your platform from the release assets
2. Make it executable: \`chmod +x termonaut-${VERSION}-your-platform\`
3. Move to your PATH: \`sudo mv termonaut-${VERSION}-your-platform /usr/local/bin/termonaut\`
4. Install shell hooks: \`termonaut advanced shell install\`

## 🔄 Breaking Changes
**None** - This is a fully backward-compatible release.

## 🐛 Known Issues
- None critical for RC testing
- Shell hook job control messages (cosmetic only)

## 🙏 What's Next?
- Your feedback on the empty command stats feature
- Final polish for v1.0.0 stable release
- Documentation improvements based on RC feedback

## 💬 Feedback
We'd love your feedback on this RC! Please test the empty command stats feature and let us know:
- How does the feature feel in daily use?
- Which display mode do you prefer?
- Any configuration improvements needed?

Report issues or share feedback at: https://github.com/oiahoon/termonaut/issues

---

**Full Changelog**: [CHANGELOG.md](https://github.com/oiahoon/termonaut/blob/main/CHANGELOG.md)
EOF

echo -e "${GREEN}✅ Release notes created${NC}"

# List all created files
echo -e "${BLUE}📦 Release artifacts:${NC}"
ls -la dist/

echo
echo -e "${GREEN}🎉 Release ${VERSION} build completed successfully!${NC}"
echo
echo -e "${YELLOW}📋 Next steps:${NC}"
echo -e "1. Test the release candidate thoroughly"
echo -e "2. Create GitHub release with assets in dist/"
echo -e "3. Update README.md with new version"
echo -e "4. Announce RC to community for testing"
echo
echo -e "${BLUE}🔗 Release files location: $(pwd)/dist/${NC}"
echo -e "${BLUE}🎯 Upload these files to GitHub release: termonaut-${VERSION}-*${NC}" 