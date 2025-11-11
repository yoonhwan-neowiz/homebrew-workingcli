class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.8"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.8/ga-darwin-arm64.tar.gz"
      sha256 "b1e73f61de11569e979a9eb1297b4e4130a4ed8eb595d8c2d3c2a8dc99952c31"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.8/ga-darwin-amd64.tar.gz"
      sha256 "2cb83ade582a32f75931b847666e42d5f073da027a8a69ea9f7c900a3720a4a8"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.8/ga-linux-arm64.tar.gz"
      sha256 "ae7a3593e1f5090d5c40500f76b5da64769bae24c4591890bd235764f1158465"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.8/ga-linux-amd64.tar.gz"
      sha256 "a6f3a333bdebffa51aa59f3eadc69bd214f5068d6f05b7ad02b7e9f727fd8145"
    end
  end

  def install
    bin.install "ga"
  end

  test do
    # ga --help 명령이 정상 동작하는지 확인
    assert_match "Git Assistant", shell_output("#{bin}/ga --help 2>&1")
  end

  def caveats
    <<~EOS
      Git Assistant (ga) has been installed!
      
      Quick Start:
        ga                    # Interactive staging
        ga commit             # AI-powered commit message
        ga opt quick status   # Check repository optimization status
        ga opt quick to-slim  # Optimize large repository
      
      For more information:
        ga --help
        ga opt help workflow
    EOS
  end
end
