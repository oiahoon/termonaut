#!/bin/bash

# Script to manually fix the homebrew-termonaut repository
# This addresses the URL and SHA256 issues

set -e

VERSION="v0.10.1"
VERSION_NO_V="0.10.1"
DARWIN_AMD64_SHA="a152db963d0807ff8d58cb965f5c985df46c86b1d95aef0af9f5a4f42e4e0136"
DARWIN_ARM64_SHA="c8d8fcdbc5fd5b8c9f2828a9d4619250b9467003a611cc07e77d8accda6564db"

echo "🔧 Fixing homebrew-termonaut repository for ${VERSION}"

# Create temp directory and clone the homebrew-termonaut repo
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "📥 Cloning homebrew-termonaut repository..."
git clone https://github.com/oiahoon/homebrew-termonaut.git
cd homebrew-termonaut

echo "📝 Creating corrected formula..."

cat > termonaut.rb << EOF
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
    assert_match "${VERSION_NO_V}", shell_output("\#{bin}/termonaut version")
  end
end
EOF

echo "🔍 Validating formula syntax..."
ruby -c termonaut.rb && echo "✅ Formula syntax is valid"

echo "📋 Formula summary:"
echo "  Version: ${VERSION_NO_V}"
echo "  AMD64 URL: https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-amd64"
echo "  ARM64 URL: https://github.com/oiahoon/termonaut/releases/download/${VERSION}/termonaut-${VERSION_NO_V}-darwin-arm64"
echo "  AMD64 SHA256: ${DARWIN_AMD64_SHA}"
echo "  ARM64 SHA256: ${DARWIN_ARM64_SHA}"

echo ""
echo "🚀 Ready to commit and push. Please review the changes above."
echo "If everything looks correct, the script will commit and push the changes."
echo ""
read -p "Continue with commit and push? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📤 Committing and pushing changes..."
    
    git config --local user.email "action@github.com"
    git config --local user.name "GitHub Action Manual Fix"
    
    git add termonaut.rb
    git commit -m "🔧 Fix formula URLs and SHA256 checksums for v0.10.1

- Fixed URLs to include 'v' prefix in release tag
- Updated correct SHA256 checksums for both Intel and ARM64 binaries
- AMD64: ${DARWIN_AMD64_SHA}
- ARM64: ${DARWIN_ARM64_SHA}

This fixes the 404 error when installing via Homebrew."
    
    git push
    echo "✅ Successfully updated homebrew-termonaut repository!"
else
    echo "❌ Aborted. No changes were pushed."
fi

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"

echo ""
echo "🍺 After the fix, users can run:"
echo "  brew upgrade termonaut"
echo "  # or"
echo "  brew uninstall termonaut && brew install oiahoon/termonaut/termonaut"
