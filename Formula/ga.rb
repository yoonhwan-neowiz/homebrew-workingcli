class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.2"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-darwin-arm64.tar.gz"
      sha256 "718b23f29cb17ac448f2cd022b801bd45ef6f6b00548bfb164546aee40006567"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-darwin-amd64.tar.gz"
      sha256 "6cc261bd4254692bb09d8510d3b53a56b7640f5383d7bbef630d90ca30943aa5"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-linux-arm64.tar.gz"
      sha256 "20e260072bc1cbd7131ae1eaa1e405ebfe0dc08c5e6d085a884159c70d5e6fd6"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.2/ga-linux-amd64.tar.gz"
      sha256 "b11274eb9460690c82028fe0053722d50dc145ac47e04132c8510e6f4fdff309"
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