class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.7"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.7/ga-darwin-arm64.tar.gz"
      sha256 "b6812588b836f22070f39554ade9dfe285db1a906069d9e10a1573ab871a83d5"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.7/ga-darwin-amd64.tar.gz"
      sha256 "3ad1846ac1f4d1a0dd65d9d2ac1fdd28f7e0562308c2daa4e06f27d8631118db"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.7/ga-linux-arm64.tar.gz"
      sha256 "d7e318f00ce646122e8e69f3b7e7315d6be4a5f2c8a7ab2b48014fa3bb1f1bd1"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.7/ga-linux-amd64.tar.gz"
      sha256 "2d5647e59bd3fe5251996930c9db9029991154e60dc995f5b3be65fd5b96ae0d"
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
