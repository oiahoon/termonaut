class Termonaut < Formula
  desc "Gamified terminal productivity tracker with achievements and XP system"
  homepage "https://github.com/oiahoon/termonaut"
  url "https://github.com/oiahoon/termonaut/releases/download/v0.9.0/termonaut-v0.9.0-darwin-amd64.tar.gz"
  sha256 "PLACEHOLDER_SHA256_AMD64"
  license "MIT"
  version "0.9.0"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.0/termonaut-v0.9.0-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.0/termonaut-v0.9.0-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.0/termonaut-v0.9.0-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.0/termonaut-v0.9.0-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
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
      🚀 Termonaut has been installed successfully!

      Quick Start:
        1. Set up shell integration:
           #{bin}/termonaut advanced shell install

        2. Restart your terminal or source your shell config:
           source ~/.zshrc    # for zsh users
           source ~/.bashrc   # for bash users

        3. Start tracking your terminal activity:
           termonaut stats
           tn stats           # short alias

      For more information:
        termonaut --help
        https://github.com/oiahoon/termonaut

      Happy terminal tracking! 🎮
    EOS
  end

  test do
    # Test basic functionality
    assert_match "termonaut", shell_output("#{bin}/termonaut --version")
    assert_match "termonaut", shell_output("#{bin}/tn --version")

    # Test help command
    assert_match "Usage:", shell_output("#{bin}/termonaut --help")

    # Test that the binary is executable
    system "#{bin}/termonaut", "--version"
  end
end