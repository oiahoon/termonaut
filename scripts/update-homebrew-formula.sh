#!/bin/bash

# Update Homebrew Formula with Real Checksums
# This script updates the Formula/termonaut.rb file with actual checksums from a release

set -e

VERSION=${1:-"v0.9.0"}
BUILD_VERSION=$(echo "$VERSION" | sed 's/^v//')
CHECKSUMS_FILE="dist/termonaut-${VERSION}-checksums.txt"

echo "🔄 Updating Homebrew formula for ${VERSION}..."

# Check if checksums file exists
if [[ ! -f "$CHECKSUMS_FILE" ]]; then
    echo "❌ Checksums file not found: $CHECKSUMS_FILE"
    echo "Please run ./scripts/build-release.sh ${VERSION} first"
    exit 1
fi

# Extract checksums
echo "📋 Extracting checksums..."

# macOS Intel
AMD64_SHA=$(grep "darwin-amd64.tar.gz" "$CHECKSUMS_FILE" | awk '{print $1}')
# macOS Apple Silicon
ARM64_SHA=$(grep "darwin-arm64.tar.gz" "$CHECKSUMS_FILE" | awk '{print $1}')
# Linux x86_64
LINUX_AMD64_SHA=$(grep "linux-amd64.tar.gz" "$CHECKSUMS_FILE" | awk '{print $1}')
# Linux ARM64
LINUX_ARM64_SHA=$(grep "linux-arm64.tar.gz" "$CHECKSUMS_FILE" | awk '{print $1}')

echo "Checksums found:"
echo "  macOS Intel:     $AMD64_SHA"
echo "  macOS ARM:       $ARM64_SHA"
echo "  Linux x86_64:    $LINUX_AMD64_SHA"
echo "  Linux ARM64:     $LINUX_ARM64_SHA"

# Update Formula/termonaut.rb
echo "📝 Updating Formula/termonaut.rb..."

# Create a temporary file
TEMP_FILE=$(mktemp)

# Update the formula with real checksums and version
sed -e "s/PLACEHOLDER_SHA256_AMD64/$AMD64_SHA/g" \
    -e "s/PLACEHOLDER_SHA256_ARM64/$ARM64_SHA/g" \
    -e "s/PLACEHOLDER_SHA256_LINUX_AMD64/$LINUX_AMD64_SHA/g" \
    -e "s/PLACEHOLDER_SHA256_LINUX_ARM64/$LINUX_ARM64_SHA/g" \
    -e "s/v0\.9\.0/$VERSION/g" \
    -e "s/\"0\.9\.0\"/\"$BUILD_VERSION\"/g" \
    Formula/termonaut.rb > "$TEMP_FILE"

# Replace the original file
mv "$TEMP_FILE" Formula/termonaut.rb

echo "✅ Formula updated successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Test the formula locally:"
echo "   brew install --formula Formula/termonaut.rb"
echo "   brew test termonaut"
echo "   brew audit --strict termonaut"
echo ""
echo "2. Commit the updated formula:"
echo "   git add Formula/termonaut.rb"
echo "   git commit -m \"Update Homebrew formula for ${VERSION}\""
echo ""
echo "3. Create or update your Homebrew tap:"
echo "   # If you have a tap repository:"
echo "   cp Formula/termonaut.rb /path/to/homebrew-tap/termonaut.rb"
echo "   cd /path/to/homebrew-tap"
echo "   git add termonaut.rb"
echo "   git commit -m \"Update termonaut to ${VERSION}\""
echo "   git push"
echo ""
echo "4. Or submit to homebrew-core:"
echo "   # Follow the instructions in docs/HOMEBREW_RELEASE.md"