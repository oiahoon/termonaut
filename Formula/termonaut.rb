class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "0.10.1"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.1/termonaut-0.10.1-darwin-amd64"
      sha256 "a152db963d0807ff8d58cb965f5c985df46c86b1d95aef0af9f5a4f42e4e0136"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.10.1/termonaut-0.10.1-darwin-arm64"
      sha256 "c8d8fcdbc5fd5b8c9f2828a9d4619250b9467003a611cc07e77d8accda6564db"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.intel?
      bin.install "termonaut-0.10.1-darwin-amd64" => "termonaut"
    elsif OS.mac? && Hardware::CPU.arm?
      bin.install "termonaut-0.10.1-darwin-arm64" => "termonaut"
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
    assert_match "0.10.1", shell_output("#{bin}/termonaut version")
  end
end
