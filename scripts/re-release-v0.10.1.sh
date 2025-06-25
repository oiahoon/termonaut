#!/bin/bash

# Script to re-release v0.10.1 with all latest fixes
# This includes TUI layout fixes, Homebrew formula fixes, and GitHub Action improvements

set -e

VERSION="v0.10.1"
VERSION_NO_V="0.10.1"

echo "🚀 Re-releasing Termonaut ${VERSION} with latest fixes..."

# Ensure we're in the right directory
cd /Users/huangyuyao/OwnWork/termonaut

# Check if we have uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo "❌ Error: You have uncommitted changes. Please commit or stash them first."
    git status
    exit 1
fi

echo "✅ Working directory is clean"

# Delete existing tag if it exists
echo "🏷️ Managing release tag..."
if git tag -l | grep -q "^${VERSION}$"; then
    echo "📋 Deleting existing local tag ${VERSION}..."
    git tag -d ${VERSION}
fi

# Delete remote tag if it exists
if git ls-remote --tags origin | grep -q "refs/tags/${VERSION}$"; then
    echo "📋 Deleting existing remote tag ${VERSION}..."
    git push --delete origin ${VERSION}
fi

# Create new tag
echo "🏷️ Creating new tag ${VERSION}..."
git tag -a ${VERSION} -m "Termonaut ${VERSION} - Enhanced Experience with Fixes

🎨 **TUI Layout Fixes:**
- Fixed tab navigation visibility in narrow terminals
- Improved content height calculation and space allocation
- Added responsive design for different terminal sizes
- Fixed content overflow and truncation issues

🍺 **Homebrew Integration Fixes:**
- Fixed GitHub Action formula generation with correct SHA256 checksums
- Resolved YAML syntax errors in workflow files
- Updated formula URLs with proper version tags
- Improved error handling and validation

🔧 **Technical Improvements:**
- Removed duplicate function definitions
- Enhanced responsive layout system
- Better terminal size detection and adaptation
- Improved user experience across all terminal sizes

📱 **Responsive Design:**
- < 60 cols: Icon-only tabs, minimal layout
- 60-80 cols: Compact layout with short names
- 80-100 cols: Standard layout with full features
- > 100 cols: Wide layout with debug information

🧪 **Testing & Validation:**
- Added comprehensive test scripts
- Improved build and release processes
- Enhanced debugging and troubleshooting tools

This release ensures Termonaut works seamlessly across all terminal sizes and installation methods."

# Push the new tag
echo "📤 Pushing tag to remote..."
git push origin ${VERSION}

echo "✅ Tag ${VERSION} created and pushed successfully!"

echo ""
echo "📋 Release Summary:"
echo "  Version: ${VERSION}"
echo "  Commit: $(git rev-parse HEAD)"
echo "  Branch: $(git branch --show-current)"
echo "  Files changed since last release:"
git diff --name-only HEAD~5 HEAD | head -10

echo ""
echo "🎯 What happens next:"
echo "1. GitHub will automatically trigger the release workflow"
echo "2. Binaries will be built for all platforms"
echo "3. Release will be published with release notes"
echo "4. Homebrew formula will be automatically updated"
echo "5. Users can upgrade with: brew upgrade termonaut"

echo ""
echo "🔍 Monitor the release process:"
echo "  GitHub Actions: https://github.com/oiahoon/termonaut/actions"
echo "  Releases: https://github.com/oiahoon/termonaut/releases"
echo "  Homebrew Tap: https://github.com/oiahoon/homebrew-termonaut"

echo ""
echo "🧪 Test the release:"
echo "  # After release is published:"
echo "  brew upgrade termonaut"
echo "  termonaut version"
echo "  termonaut tui  # Test the layout fixes"

echo ""
echo "✅ Re-release ${VERSION} initiated successfully!"
echo "🎉 Check GitHub Actions for build progress."
