class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.6"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-darwin-arm64.tar.gz"
      sha256 "8c95f3db047887f6d569b29eadef4ee456fa0110da143606d8d9a8550397f0da"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-darwin-amd64.tar.gz"
      sha256 "17eb81373108f2002fe89468bfd6f111db9d0606f3dc62768514739e0c8132ae"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-linux-arm64.tar.gz"
      sha256 "c549c02880f14e775149afe023c31236514a8c1686adfc2f78a48800035f1d4d"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.6/ga-linux-amd64.tar.gz"
      sha256 "11e50bab46dfe2c3ac1a3fba2388b15ba7e790e492488e612f89ce20ec112df4"
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
