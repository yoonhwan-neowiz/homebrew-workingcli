class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.2"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-darwin-arm64.tar.gz"
      sha256 "8a03e97d2ec067d0922970cd1683683a2db7d3b1a671b44a178732c99280f970"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-darwin-amd64.tar.gz"
      sha256 "94742973ae245a73e4ac8a7ae0f4b4d050eff0deff7dd71bd2e20ad9df843d70"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-linux-arm64.tar.gz"
      sha256 "d427cc4c09e88bc2a9bac51668de2b72d6a1727de667962dba1944e301a40622"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-linux-amd64.tar.gz"
      sha256 "a488e5c1920c73231b7ccc8d6c472caa7718da3ba36837cd80223a0606bd0a5b"
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
