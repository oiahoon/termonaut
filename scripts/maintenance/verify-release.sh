#!/bin/bash

# Script to verify the v0.10.1 release status
# Checks GitHub Actions, release availability, and Homebrew formula

set -e

VERSION="v0.10.1"

echo "🔍 Verifying Termonaut ${VERSION} release status..."

echo ""
echo "1. 🏷️ Tag Status:"
echo "   Local tag: $(git tag -l | grep ${VERSION} || echo 'Not found')"
echo "   Remote tag: $(git ls-remote --tags origin | grep ${VERSION} | head -1 | cut -f2 | sed 's/refs\/tags\///' || echo 'Not found')"

echo ""
echo "2. 📦 GitHub Release:"
echo "   Checking release availability..."
RELEASE_STATUS=$(curl -s "https://api.github.com/repos/oiahoon/termonaut/releases/tags/${VERSION}" | jq -r '.tag_name // "not_found"')
if [ "$RELEASE_STATUS" = "${VERSION}" ]; then
    echo "   ✅ Release ${VERSION} is published"
    echo "   📋 Assets:"
    curl -s "https://api.github.com/repos/oiahoon/termonaut/releases/tags/${VERSION}" | jq -r '.assets[].name' | sed 's/^/     - /'
else
    echo "   ⏳ Release ${VERSION} is not yet published (GitHub Actions may still be running)"
fi

echo ""
echo "3. 🔨 GitHub Actions:"
echo "   Latest workflow runs:"
curl -s "https://api.github.com/repos/oiahoon/termonaut/actions/runs?per_page=3" | jq -r '.workflow_runs[] | "   - \(.name): \(.status) (\(.conclusion // "running"))"'

echo ""
echo "4. 🍺 Homebrew Formula:"
echo "   Checking homebrew-termonaut repository..."
FORMULA_VERSION=$(curl -s "https://raw.githubusercontent.com/oiahoon/homebrew-termonaut/main/termonaut.rb" | grep 'version' | head -1 | sed 's/.*"\(.*\)".*/\1/')
echo "   Formula version: ${FORMULA_VERSION}"
if [ "$FORMULA_VERSION" = "0.10.1" ]; then
    echo "   ✅ Homebrew formula is updated"
else
    echo "   ⏳ Homebrew formula may still be updating"
fi

echo ""
echo "5. 📥 Download URLs:"
echo "   Intel macOS: https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-0.10.1-darwin-amd64"
echo "   ARM64 macOS: https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-0.10.1-darwin-arm64"
echo "   Linux x64:   https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-0.10.1-linux-amd64"

echo ""
echo "6. 🧪 Test Commands:"
echo "   # After release is complete:"
echo "   brew upgrade termonaut"
echo "   termonaut version"
echo "   termonaut tui --mode smart"

echo ""
echo "🔗 Useful Links:"
echo "   GitHub Actions: https://github.com/oiahoon/termonaut/actions"
echo "   Releases: https://github.com/oiahoon/termonaut/releases"
echo "   Homebrew Tap: https://github.com/oiahoon/homebrew-termonaut"

echo ""
if [ "$RELEASE_STATUS" = "${VERSION}" ]; then
    echo "✅ Release ${VERSION} verification completed - Release is live!"
else
    echo "⏳ Release ${VERSION} is in progress - Check GitHub Actions for status"
fi
