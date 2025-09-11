class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.8"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.8/ga-darwin-arm64.tar.gz"
      sha256 "123c8fcfacb54c1301cda6247c195bb5d0cce874fd1d39fd57171c4a15f6c6c2"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.8/ga-darwin-amd64.tar.gz"
      sha256 "0e9faebeda13ba280605658dc726d8d3e70615058a33641098e3891d64613b31"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.8/ga-linux-arm64.tar.gz"
      sha256 "54c401da463ee553264477919f96e548e327b335822dfd5b09e0cf02e47c1358"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.8/ga-linux-amd64.tar.gz"
      sha256 "9ecd3285851a42b86e5a39d7871204845c02b9634f7e3276594cdba7436dcc19"
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
