class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.5"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.5/ga-darwin-arm64.tar.gz"
      sha256 "b38ec6c5a5facb5669a1ec91e7f31e8c80e2718538f775e34a6af569574032ed"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.5/ga-darwin-amd64.tar.gz"
      sha256 "31c404f8c876ff070d2f25dba0c9bdda4383c458338c3fb6344e857645675325"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.5/ga-linux-arm64.tar.gz"
      sha256 "33ac1dc0a3aff3545c70587740b7f5f2a9b4bfd366013aa65204204f8e8037ed"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.5/ga-linux-amd64.tar.gz"
      sha256 "cb4962b05128cd9f817a2070980792135c6c497e6f4c127cf51d2ad286cecea7"
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
