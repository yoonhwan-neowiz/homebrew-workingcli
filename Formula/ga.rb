class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.0"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.0/ga-darwin-arm64.tar.gz"
      sha256 "31c4267f6b53465c338b5236ece1a1f147eb2a92f29f9f2f6b86bb5f30b72f68"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.0/ga-darwin-amd64.tar.gz"
      sha256 "542716f3ffeb1bc40330537789763bc7e3475c2fa757e043d80071aa2ad77ccf"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.0/ga-linux-arm64.tar.gz"
      sha256 "b1c881fa155873fb896d5c377baa765e387fb6fa83a7991e0977383bd2140b62"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.0/ga-linux-amd64.tar.gz"
      sha256 "a69e63e4c3364540bf6ed942073534a0f3ad04905050e137f902fcff9fd74a0d"
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
