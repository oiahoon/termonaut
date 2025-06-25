#!/bin/bash

# Script to update Homebrew formula with correct SHA256 values
# Usage: ./scripts/update-homebrew-formula.sh [version]

set -e

VERSION=${1:-"v0.10.0"}
VERSION_NO_V=${VERSION#v}

echo "🍺 Updating Homebrew formula for Termonaut ${VERSION}"

# URLs for the binaries
AMD64_URL="https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-amd64"
ARM64_URL="https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-arm64"

echo "📥 Downloading binaries to calculate SHA256..."

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download binaries
echo "  Downloading AMD64 binary..."
curl -sSL "$AMD64_URL" -o termonaut-darwin-amd64

echo "  Downloading ARM64 binary..."
curl -sSL "$ARM64_URL" -o termonaut-darwin-arm64

# Calculate SHA256
echo "🔐 Calculating SHA256 checksums..."
DARWIN_AMD64_SHA=$(shasum -a 256 termonaut-darwin-amd64 | cut -d' ' -f1)
DARWIN_ARM64_SHA=$(shasum -a 256 termonaut-darwin-arm64 | cut -d' ' -f1)

echo "  AMD64 SHA256: ${DARWIN_AMD64_SHA}"
echo "  ARM64 SHA256: ${DARWIN_ARM64_SHA}"

# Go back to project directory
cd - > /dev/null

# Update the formula
echo "📝 Updating Formula/termonaut.rb..."

cat > Formula/termonaut.rb << EOF
class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "${VERSION_NO_V}"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-amd64"
      sha256 "${DARWIN_AMD64_SHA}"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-arm64"
      sha256 "${DARWIN_ARM64_SHA}"
    end
  end

  def install
    # Install the binary directly (no tar.gz extraction needed)
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-${VERSION_NO_V}-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-${VERSION_NO_V}-darwin-arm64" => "termonaut"
    end
  end

  def caveats
    <<~EOS
      🚀 Termonaut has been installed successfully!

      To get started:
      1. Initialize shell integration:
         termonaut advanced shell install

      2. Restart your terminal or run:
         source ~/.bashrc  # or ~/.zshrc

      3. Start tracking your productivity:
         termonaut stats
         termonaut tui

      4. Set up GitHub integration (optional):
         termonaut github auth
         termonaut github sync now

      📖 Documentation: https://github.com/oiahoon/termonaut
      🐛 Issues: https://github.com/oiahoon/termonaut/issues

      Happy terminal productivity tracking! 🎯
    EOS
  end

  test do
    assert_match "${VERSION_NO_V}", shell_output("#{bin}/termonaut version")
  end
end
EOF

# Clean up
rm -rf "$TEMP_DIR"

echo "✅ Formula updated successfully!"
echo ""
echo "📋 Summary:"
echo "  Version: ${VERSION_NO_V}"
echo "  AMD64 SHA256: ${DARWIN_AMD64_SHA}"
echo "  ARM64 SHA256: ${DARWIN_ARM64_SHA}"
echo ""
echo "🧪 To test the formula locally:"
echo "  brew install --build-from-source Formula/termonaut.rb"
echo ""
echo "🚀 To commit the changes:"
echo "  git add Formula/termonaut.rb"
echo "  git commit -m \"🍺 Update Homebrew formula to ${VERSION}\""
