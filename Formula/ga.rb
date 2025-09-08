class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.2"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-darwin-arm64.tar.gz"
      sha256 "18e85ad34c01ae709fb74966f54d84d0aa6362b6cf7ce66a724ed57583d1f969"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-darwin-amd64.tar.gz"
      sha256 "11d4fea6e5da2e757a934a4aa264e916f062ba2c43f8ba549ae4c9549e6c54ef"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-linux-arm64.tar.gz"
      sha256 "00df0670e45a7fffbfc5f624dc13949538f83bebbba55ca1c7fb40879e144041"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-linux-amd64.tar.gz"
      sha256 "636f7e1f2d19213129abbb73b67f77a35f9842e649892a13aae03e7a80946590"
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