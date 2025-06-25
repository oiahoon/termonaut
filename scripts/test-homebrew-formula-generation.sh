#!/bin/bash

# Test script to simulate the GitHub Action formula generation
# This helps verify the formula generation logic works correctly

set -e

VERSION="v0.10.0"
VERSION_NO_V="0.10.0"
DARWIN_AMD64_SHA="7144cf15ccc632c5f9d5a6bf4376024ae8e355cc01631f08979ac1bd6c82e21a"
DARWIN_ARM64_SHA="6f479a97351f8908e91d5cf4a42f697c83ec3fcc15566853f5fb379dbb166972"

echo "🧪 Testing Homebrew formula generation logic..."
echo "Version: ${VERSION}"
echo "AMD64 SHA: ${DARWIN_AMD64_SHA}"
echo "ARM64 SHA: ${DARWIN_ARM64_SHA}"

# Create test directory
TEST_DIR=$(mktemp -d)
cd "$TEST_DIR"

echo "📝 Generating formula using GitHub Action logic..."

# Simulate the GitHub Action formula generation
cat > termonaut.rb << 'EOF'
class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "VERSION_PLACEHOLDER"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/VERSION_TAG_PLACEHOLDER/termonaut-VERSION_PLACEHOLDER-darwin-amd64"
      sha256 "AMD64_SHA_PLACEHOLDER"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/VERSION_TAG_PLACEHOLDER/termonaut-VERSION_PLACEHOLDER-darwin-arm64"
      sha256 "ARM64_SHA_PLACEHOLDER"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-VERSION_PLACEHOLDER-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-VERSION_PLACEHOLDER-darwin-arm64" => "termonaut"
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
    assert_match "VERSION_PLACEHOLDER", shell_output("#{bin}/termonaut version")
  end
end
EOF

# Replace placeholders with actual values
sed -i.bak "s/VERSION_PLACEHOLDER/${VERSION_NO_V}/g" termonaut.rb
sed -i.bak "s/VERSION_TAG_PLACEHOLDER/${VERSION}/g" termonaut.rb
sed -i.bak "s/AMD64_SHA_PLACEHOLDER/${DARWIN_AMD64_SHA}/g" termonaut.rb
sed -i.bak "s/ARM64_SHA_PLACEHOLDER/${DARWIN_ARM64_SHA}/g" termonaut.rb
rm termonaut.rb.bak

echo "✅ Formula generated successfully!"
echo ""
echo "🔍 Validating Ruby syntax..."
if ruby -c termonaut.rb > /dev/null 2>&1; then
    echo "✅ Ruby syntax is valid!"
else
    echo "❌ Ruby syntax validation failed!"
    ruby -c termonaut.rb
    exit 1
fi

echo ""
echo "📋 Generated formula preview:"
echo "----------------------------------------"
head -20 termonaut.rb
echo "..."
echo "----------------------------------------"

echo ""
echo "💾 Copying validated formula back to project..."
cp termonaut.rb /Users/huangyuyao/OwnWork/termonaut/Formula/termonaut.rb

# Clean up
cd - > /dev/null
rm -rf "$TEST_DIR"

echo "✅ Test completed successfully!"
echo "🚀 The GitHub Action formula generation logic should now work correctly."
