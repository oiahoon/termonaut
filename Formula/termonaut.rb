class Termonaut < Formula
  desc "Gamified terminal productivity tracker with achievements, XP system, and personalized avatars"
  homepage "https://github.com/oiahoon/termonaut"
  url "https://github.com/oiahoon/termonaut/releases/download/v0.9.2/termonaut-0.9.2-darwin-amd64"
  sha256 "069df6cbf275717490a5764f99b69f5be83660044a527079057fc448e98d39c8"
  license "MIT"
  version "0.9.2"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.2/termonaut-0.9.2-darwin-amd64"
      sha256 "069df6cbf275717490a5764f99b69f5be83660044a527079057fc448e98d39c8"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.2/termonaut-0.9.2-darwin-arm64"
      sha256 "dac85679330e9707420dc960deb176aa993a29bb3deb62c0414521bc9f24e7c0"
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
      bin.install "termonaut-0.9.2-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-0.9.2-darwin-arm64" => "termonaut"
    end

    # Create symlink for short command
    bin.install_symlink "termonaut" => "tn"
  end

  def post_install
    puts <<~EOS
      🎮 Termonaut v0.9.2 has been installed successfully!

      ✨ NEW: Optimized Easter Eggs & Network-Resilient Avatar System!

      Quick Start:
        1. Set up shell integration:
           #{bin}/termonaut advanced shell install

        2. Restart your terminal or source your shell config:
           source ~/.zshrc    # for zsh users
           source ~/.bashrc   # for bash users

        3. View your avatar dashboard:
           termonaut stats    # Enhanced with avatar display
           tn stats           # short alias

        4. Test your system:
           termonaut avatar-test     # Test avatar system & network
           termonaut terminal-test   # Test terminal compatibility
           tn avatar-test           # short alias

        5. Customize your avatar:
           termonaut avatar config --style pixel-art
           tn avatar preview -l 25

      🎮 v0.9.2 Highlights:
        • 60%+ reduced easter egg interruptions
        • New easter eggs: Python, JS, Database, Testing, AI Tools
        • Network-resilient avatar system with offline fallbacks
        • Enhanced modern terminal support (Warp, iTerm2, etc.)

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