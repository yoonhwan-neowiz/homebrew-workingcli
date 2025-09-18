class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.4"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.4/ga-darwin-arm64.tar.gz"
      sha256 "6ff645dfebb46137ff3207ff3089611d5522d3fe4560800dfbd38d56d2a67b83"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.4/ga-darwin-amd64.tar.gz"
      sha256 "4f194e7e88475934b081b26cce4f594c52555b76025ce1ec117e2d359bb7eefb"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.4/ga-linux-arm64.tar.gz"
      sha256 "750df85e65f4b198d32194903f0edf295668c130cd9c9cc6f6e69f4705b873c9"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.4/ga-linux-amd64.tar.gz"
      sha256 "ad8c409144d5b4760e9e8ec77d0ea3359cd98a334f4a60fd8620036c297f9958"
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
