class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.0"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/WorkingCli/releases/download/v0.1.0/ga-darwin-arm64.tar.gz"
      sha256 "PENDING_ARM64_SHA256"  # release.sh가 자동 업데이트
    else
      url "https://github.com/yoonhwan-neowiz/WorkingCli/releases/download/v0.1.0/ga-darwin-amd64.tar.gz"
      sha256 "PENDING_AMD64_SHA256"  # release.sh가 자동 업데이트
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/WorkingCli/releases/download/v0.1.0/ga-linux-arm64.tar.gz"
      sha256 "PENDING_LINUX_ARM64_SHA256"  # release.sh가 자동 업데이트
    else
      url "https://github.com/yoonhwan-neowiz/WorkingCli/releases/download/v0.1.0/ga-linux-amd64.tar.gz"
      sha256 "PENDING_LINUX_AMD64_SHA256"  # release.sh가 자동 업데이트
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