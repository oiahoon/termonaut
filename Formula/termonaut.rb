class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "0.10.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.0/termonaut-0.10.0-darwin-amd64"
      sha256 "PLACEHOLDER_INTEL_SHA256"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.0/termonaut-0.10.0-darwin-arm64"
      sha256 "PLACEHOLDER_ARM_SHA256"
    end
  end

  def install
    # Install the binary directly (no tar.gz extraction needed)
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-0.9.5-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-0.9.5-darwin-arm64" => "termonaut"
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
    assert_match "0.9.5", shell_output("#{bin}/termonaut version")
  end
end