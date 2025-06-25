class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "0.10.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.0/termonaut-0.10.0-darwin-amd64"
      sha256 "7144cf15ccc632c5f9d5a6bf4376024ae8e355cc01631f08979ac1bd6c82e21a"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.0/termonaut-0.10.0-darwin-arm64"
      sha256 "6f479a97351f8908e91d5cf4a42f697c83ec3fcc15566853f5fb379dbb166972"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-0.10.0-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-0.10.0-darwin-arm64" => "termonaut"
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
    assert_match "0.10.0", shell_output("#{bin}/termonaut version")
  end
end
