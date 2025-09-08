class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.3"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.3/ga-darwin-arm64.tar.gz"
      sha256 "a3f0471bd1d389a147e1fdc425aa1efe892b413f0a0e9989a284876cd249a2e2"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.3/ga-darwin-amd64.tar.gz"
      sha256 "f7146b1eb195f65c524f24d6c01652e3a7f5b9801ca36d498866245f5d38205c"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.3/ga-linux-arm64.tar.gz"
      sha256 "8d2ad689f0ec1fd4d532db1084d9be16aaaaaadb1bbdfb4624ff968cf9a1e29e"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.3/ga-linux-amd64.tar.gz"
      sha256 "0ec9f87c67686f4bc7cfbe28cf433b6135a1c96a87c326ca58b09970094330bd"
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
