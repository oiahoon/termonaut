class Termonaut < Formula
  desc "Gamified terminal productivity tracker with achievements, XP system, and personalized avatars"
  homepage "https://github.com/oiahoon/termonaut"
  url "https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-darwin-amd64.tar.gz"
  sha256 "ba539d7ed329c8729bf71a7f2ea3cb975dddc721a49ba7b2745574571dbe7c6b"
  license "MIT"
  version "0.9.1"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-darwin-amd64.tar.gz"
      sha256 "ba539d7ed329c8729bf71a7f2ea3cb975dddc721a49ba7b2745574571dbe7c6b"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-darwin-arm64.tar.gz"
      sha256 "c5732407ef4383d13a297313a714e276b3848f81da857a75bb0511fcd0db2a37"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-linux-amd64.tar.gz"
      sha256 "f49f9ae5504b0ae056d9b806a4bc8bdb6867619cda29b00193b2c524b84defff"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.1/termonaut-v0.9.1-linux-arm64.tar.gz"
      sha256 "4d15e72e85da870f5988cc2b98145fe66b5bf63988dfa9034428e1e17c3b292b"
    end
  end

  def install
    bin.install "termonaut-darwin-amd64" => "termonaut" if OS.mac? && Hardware::CPU.intel?
    bin.install "termonaut-darwin-arm64" => "termonaut" if OS.mac? && Hardware::CPU.arm?
    bin.install "termonaut-linux-amd64" => "termonaut" if OS.linux? && Hardware::CPU.intel?
    bin.install "termonaut-linux-arm64" => "termonaut" if OS.linux? && Hardware::CPU.arm?

    # Create symlink for short command
    bin.install_symlink "termonaut" => "tn"
  end

  def post_install
    puts <<~EOS
      🎨 Termonaut v0.9.1 has been installed successfully!

      ✨ NEW: Avatar System with personalized ASCII art!

      Quick Start:
        1. Set up shell integration:
           #{bin}/termonaut advanced shell install

        2. Restart your terminal or source your shell config:
           source ~/.zshrc    # for zsh users
           source ~/.bashrc   # for bash users

        3. View your new avatar dashboard:
           termonaut stats    # Enhanced with avatar display
           tn stats           # short alias

        4. Customize your avatar:
           termonaut avatar config --style pixel-art
           termonaut avatar show
           tn avatar preview -l 25

      🎮 Avatar Styles: pixel-art, bottts, adventurer, avataaars
      📏 Smart Sizing: Auto-adjusts based on your terminal width

      For more information:
        termonaut --help
        https://github.com/oiahoon/termonaut

      Happy terminal tracking with personalized avatars! 🎨✨
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