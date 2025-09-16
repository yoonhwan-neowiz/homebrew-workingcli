class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.2"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-darwin-arm64.tar.gz"
      sha256 "f24cc0ac70ca303ecb619a035291c372b38dfa584919db3e4c15cab94bf12738"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-darwin-amd64.tar.gz"
      sha256 "6ba0bca070e58a17b1d32b957c70cd13cee14c5fd2e72c2f65c8a4cb0f2446d4"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-linux-arm64.tar.gz"
      sha256 "b92ad54b61aec050ce999ebaf3c0062a938cd720ec350ffbef1aa75862d6f752"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.2/ga-linux-amd64.tar.gz"
      sha256 "dc8ee39f50ec4a2b20484225c6485b1c58ad411ea38f36b33ba6fd8491ff9e83"
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
