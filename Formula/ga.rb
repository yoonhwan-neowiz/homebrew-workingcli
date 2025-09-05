class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.0"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.0/ga-darwin-arm64.tar.gz"
      sha256 "055857e9fd878764b4e660a554b91b073acb6a8d2c4c4ff71ee545a4e471ea62"  # release.sh가 자동 업데이트
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.0/ga-darwin-amd64.tar.gz"
      sha256 "9b4ee989d0f1b1a441368f97d4c7af680d7f964295115dabcc98862b56340cc5"  # release.sh가 자동 업데이트
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.0/ga-linux-arm64.tar.gz"
      sha256 "18574392309448578c60a87789092fb73865db0084bfd73fbaf5fb5f9ca520e7"  # release.sh가 자동 업데이트
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.0/ga-linux-amd64.tar.gz"
      sha256 "075fc38c200857047ef9b1f1c03d3814098b102a4557247cd0c678770a585c53"  # release.sh가 자동 업데이트
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