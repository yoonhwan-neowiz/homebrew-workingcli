class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.7"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.7/ga-darwin-arm64.tar.gz"
      sha256 "6e59e200db9972cc0f43adf451b01af149c85b61fe53a7daeb58c773c688d843"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.7/ga-darwin-amd64.tar.gz"
      sha256 "e5ca74b663d72b4c118dc4ed42e7a103b7ae075055d2afd1a384ad6935776773"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.7/ga-linux-arm64.tar.gz"
      sha256 "dd8003aef58991621343c906dfa4e505e57fbe3396121128f0a291702a911feb"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.7/ga-linux-amd64.tar.gz"
      sha256 "db67013c6e03c4fe76f1286cf4f4923249a29cdf0bdf49153ea039b2f111c5c8"
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
