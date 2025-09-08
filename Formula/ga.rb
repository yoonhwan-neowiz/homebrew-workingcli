class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.4"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-darwin-arm64.tar.gz"
      sha256 "59d075ca187091024b4ef5aaef55e621e9f524ddde45d9009e6132e803762e0a"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-darwin-amd64.tar.gz"
      sha256 "143ed1a6ce39f31501a663e77650cefa40b7484238b674d59f83fe11e881b0b7"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-linux-arm64.tar.gz"
      sha256 "808ae227948d82e12173bdd4fa8c0dd313fa1b7ce8e234738a6eb004493ad86c"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-linux-amd64.tar.gz"
      sha256 "0e495ba5458520dc33e291cc72f9e72330ac49cd380fd03034db933d6a251472"
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
