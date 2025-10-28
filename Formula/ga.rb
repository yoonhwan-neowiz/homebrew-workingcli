class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.5"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.5/ga-darwin-arm64.tar.gz"
      sha256 "275ee2bc6db6a0f01e9f11a55607172f1c308be3f9f0f8d90237fe3c03477881"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.5/ga-darwin-amd64.tar.gz"
      sha256 "cda6bec33716d6c2f717dd37b18c6f334bdbebab8002c1666ab39f104dcdd35a"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.5/ga-linux-arm64.tar.gz"
      sha256 "62793ff88ad9d32ce0137d6b2ae845ebaf8d065bf79f6b89d813bf08d103d377"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.5/ga-linux-amd64.tar.gz"
      sha256 "574b2678698ed0350b5eeae971d8ea269544a757aeed471d7096c7b6d79c5324"
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
