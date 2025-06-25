#!/bin/bash

# Termonaut v0.9.0-rc2 Release Script
# Builds and prepares release candidate 2 with user feedback fixes

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Version information
VERSION="0.9.0-rc2"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}🚀 Termonaut v${VERSION} Release Candidate 2 Build${NC}"
echo -e "══════════════════════════════════════════════════════"
echo -e "Version: ${GREEN}${VERSION}${NC}"
echo -e "Date: ${GREEN}${BUILD_DATE}${NC}"
echo -e "Commit: ${GREEN}${COMMIT_SHA}${NC}"
echo
echo -e "${YELLOW}🔧 User Feedback Fixes Included:${NC}"
echo -e "  • ⚡ Short 'tn' alias support with automatic symlink creation"
echo -e "  • 🔇 Enhanced job control suppression (complete elimination)"
echo -e "  • 🔧 User-friendly shell completion setup"
echo -e "  • 🔍 Fixed empty command stats detection logic"
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

# macOS (Intel)
echo -e "Building for macOS (Intel)..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-amd64 cmd/termonaut/*.go

# macOS (Apple Silicon)
echo -e "Building for macOS (Apple Silicon)..."
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-darwin-arm64 cmd/termonaut/*.go

# Linux (x64)
echo -e "Building for Linux (x64)..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-linux-amd64 cmd/termonaut/*.go

# Linux (ARM64)
echo -e "Building for Linux (ARM64)..."
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-linux-arm64 cmd/termonaut/*.go

# Windows (x64)
echo -e "Building for Windows (x64)..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o dist/termonaut-${VERSION}-windows-amd64.exe cmd/termonaut/*.go

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
cat > dist/RELEASE_NOTES_${VERSION}.md << EOF
# Termonaut ${VERSION} Release Candidate 2

## 🔧 User Feedback Fixes

This RC2 release addresses the three main user feedback issues from RC1:

### 1. 📝 Short Command Alias - FIXED ✅
**Problem**: \`termonaut xxx\` commands too long
**Solution**: Added \`tn\` as short alias for all commands

\`\`\`bash
# Before (long commands)
termonaut stats
termonaut config set theme emoji
termonaut advanced shell install

# After (short commands)
tn stats
tn config set theme emoji
tn advanced shell install
\`\`\`

### 2. 🔇 Job Control Messages - FIXED ✅
**Problem**: Seeing \`[1] + 91374 done\` messages
**Solution**: Enhanced shell hook with triple suppression

**New Hook Features:**
- Method 1: \`nohup\` for complete process detachment
- Method 2: Immediate \`disown\` to remove from job table
- Method 3: Temporary job control disable
- Zsh: \`setopt NO_NOTIFY\` and \`NO_HUP\`
- Bash: \`set +m\` to disable job control

**Apply Fix:**
\`\`\`bash
./fix_hook.sh
# OR
tn advanced shell install --force
\`\`\`

### 3. 🔍 Empty Command Stats - FIXED ✅
**Problem**: Pressing Enter on empty line doesn't show stats
**Solution**: Fixed empty command detection logic

**Improvements:**
- Handle case when no arguments provided
- Better string trimming and empty detection
- Enhanced error handling for edge cases

**Test Fix:**
\`\`\`bash
tn config set empty_command_stats true
# Press Enter on empty command line → Should show stats!
\`\`\`

## 🌟 All RC1 Features Still Included

### 💡 Empty Command Stats
- Quick stats when pressing Enter on empty command line
- Configurable minimal/rich display modes
- Privacy-aware integration

### 🎮 Complete Gamification System
- XP and leveling with space-themed progression
- 20+ achievements with rarity system
- Easter eggs and contextual surprises
- Command categorization (17 categories)

### 📊 Advanced Analytics
- Comprehensive productivity insights
- GitHub integration with heatmaps
- TUI dashboard interface
- Privacy-first command sanitization

## 🛠 Installation

### Quick Install (Recommended)
\`\`\`bash
# Download and install latest RC2
curl -sSL https://raw.githubusercontent.com/oiahoon/termonaut/main/install.sh | bash

# Or download specific binary for your platform from release assets
\`\`\`

### Manual Installation
1. Download the appropriate binary for your platform
2. Make it executable: \`chmod +x termonaut-${VERSION}-your-platform\`
3. Move to PATH: \`sudo mv termonaut-${VERSION}-your-platform /usr/local/bin/termonaut\`
4. Install shell hooks: \`tn advanced shell install\`

## 🧪 Testing the Fixes

\`\`\`bash
# 1. Test short alias
tn version
tn stats

# 2. Test job control fix (should be silent)
# Run normal commands - no [1] + done messages should appear

# 3. Test empty command stats
tn config set empty_command_stats true
# Press Enter on empty command line
\`\`\`

## 🔄 Breaking Changes
**None** - Fully backward compatible with RC1.

## 💬 What's Next?
- Final testing and feedback collection
- v1.0.0 stable release preparation
- Documentation polish

## 🙏 Feedback Welcome!
Test these fixes and let us know:
- Do the short commands work well in your workflow?
- Are the job control messages gone?
- Does empty command stats work as expected?

Report issues: https://github.com/oiahoon/termonaut/issues

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
echo -e "1. Test all three fixes thoroughly"
echo -e "2. Create GitHub release with assets in dist/"
echo -e "3. Update install.sh with new version"
echo -e "4. Announce RC2 fixes to community"
echo
echo -e "${BLUE}🔧 Quick Test Commands:${NC}"
echo -e "• Test short alias: ./dist/termonaut-${VERSION}-$(go env GOOS)-$(go env GOARCH) version"
echo -e "• Test empty stats: ./dist/termonaut-${VERSION}-$(go env GOOS)-$(go env GOARCH) log-command ''"
echo -e "• Test help: ./dist/termonaut-${VERSION}-$(go env GOOS)-$(go env GOARCH) --help"