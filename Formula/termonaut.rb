class Termonaut < Formula
  desc "Gamified terminal productivity tracker with XP, achievements, and GitHub integration"
  homepage "https://github.com/oiahoon/termonaut"
  version "0.9.5"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.5/termonaut-0.9.5-darwin-amd64"
      sha256 "55ea6d3c934ffe4005e30dc7f49ea4c4c9f635f90c6e345750a6c171d936a4a6"
    end

    on_arm do
      url "https://github.com/oiahoon/termonaut/releases/download/v0.9.5/termonaut-0.9.5-darwin-arm64"
      sha256 "68298dd4439a7aaa07f87fcfac8b8dbe9908d2032eed1f5a72c5d72157f1be34"
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