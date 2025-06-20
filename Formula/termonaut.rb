class Termonaut < Formula
  desc "Gamified terminal productivity tracker with achievements, XP system, and personalized avatars"
  homepage "https://github.com/oiahoon/termonaut"
  url "https://github.com/oiahoon/termonaut/releases/download/v0.9.4/termonaut-0.9.4-darwin-amd64"
  sha256 "7377ed2c98cabf92b269b76bc35d958074e8bb99ae66eafe8fe68b903fe440c3"
  license "MIT"
  version "0.9.4"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.4/termonaut-0.9.4-darwin-amd64"
      sha256 "7377ed2c98cabf92b269b76bc35d958074e8bb99ae66eafe8fe68b903fe440c3"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.4/termonaut-0.9.4-darwin-arm64"
      sha256 "659c45152dc0b5f02ae6e7c912657e1d5dec73b57261b32fdceb40a1c6e02a02"
    end
  end

  # Linux builds temporarily disabled due to CGO cross-compilation issues
  # Will be re-enabled in future releases
  # on_linux do
  #   on_intel do
  #     url "https://github.com/oiahoon/termonaut/releases/download/v0.9.2/termonaut-0.9.2-linux-amd64"
  #     sha256 "TBD"
  #   end
  #
  #   on_arm do
  #     url "https://github.com/oiahoon/termonaut/releases/download/v0.9.2/termonaut-0.9.2-linux-arm64"
  #     sha256 "TBD"
  #   end
  # end

  def install
    # Install the binary directly (no tar.gz extraction needed)
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-0.9.4-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-0.9.4-darwin-arm64" => "termonaut"
    end

    # Create symlink for short command
    bin.install_symlink "termonaut" => "tn"
  end

  def post_install
    puts <<~EOS
      🚀 Termonaut v0.9.4 has been installed successfully!

      ✨ NEW: Enhanced Features & Final Polish for Production Readiness!

      Quick Start:
        1. Set up shell integration:
           #{bin}/termonaut advanced shell install

        2. Restart your terminal or source your shell config:
           source ~/.zshrc    # for zsh users
           source ~/.bashrc   # for bash users

        3. View your enhanced dashboard:
           termonaut stats    # Production-ready with avatar display
           tn stats           # short alias

        4. Test your system:
           termonaut avatar-test     # Test avatar system & network
           termonaut terminal-test   # Test terminal compatibility
           tn version               # Check detailed version info

        5. Explore GitHub integration:
           termonaut github sync now    # Sync your profile
           termonaut github badges      # Generate README badges
           tn github heatmap           # Create activity heatmap

      🚀 v0.9.4 Highlights:
        • 95% feature completeness for v1.0 readiness
        • Comprehensive project planning & documentation excellence
        • Enhanced development workflow & release process
        • Production-ready with comprehensive error handling
        • Safe shell configuration management system

      🎨 Avatar Styles: pixel-art, bottts, adventurer, avataaars
      📏 Smart Sizing: Auto-adjusts based on your terminal width

      For more information:
        termonaut --help
        https://github.com/oiahoon/termonaut

      Happy terminal tracking with enhanced features! 🚀✨
    EOS
  end

  test do
    # Test basic functionality
    assert_match "Termonaut", shell_output("#{bin}/termonaut version")
    assert_match "Termonaut", shell_output("#{bin}/tn version")

    # Test help command
    assert_match "Usage:", shell_output("#{bin}/termonaut --help")

    # Test that the binary is executable
    system "#{bin}/termonaut", "version"
  end
end