class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.6"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-darwin-arm64.tar.gz"
      sha256 "f8b35509580dceb7ed8c95838c76037a06d5c94c643a8e7cd394a903e681bb65"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-darwin-amd64.tar.gz"
      sha256 "97a434568cf0d4b9567123671f15a0f278d2f5603ced3d862a6caa024a9a209b"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-linux-arm64.tar.gz"
      sha256 "e40d6de7d46df337e87335c5cf2acd8ecfc302d594cf69e8444c5faa2cf1001d"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-linux-amd64.tar.gz"
      sha256 "73bb8fc6921b8d9e0a64535973e8ced7245d88d7b5aa02c1a91291d4a5df8f5d"
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
