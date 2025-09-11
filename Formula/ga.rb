class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.1"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.1/ga-darwin-arm64.tar.gz"
      sha256 "b6410588b461e53ac15bcf76c9d25ffea1ce2b2f52aec84eb6e9a932de45fbee"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.1/ga-darwin-amd64.tar.gz"
      sha256 "2987af74994cbe88ddd3a138e5ed4631998b5a1e5fa2146bb1768a021674b5c6"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.1/ga-linux-arm64.tar.gz"
      sha256 "c6aad07c0e85e4d581b998f4cea1cac2e5af516d699e358942c4942e09b889ed"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.1/ga-linux-amd64.tar.gz"
      sha256 "addd1d6413e0db8d32f6911e2b16085fa5c259fa205a273bd2f039f6aefcdad6"
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
