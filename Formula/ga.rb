class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.1.4"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-darwin-arm64.tar.gz"
      sha256 "6489d85f170abb0c24e99ac86c07aced48c235d4c8dcb932c7f3408e001809b9"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-darwin-amd64.tar.gz"
      sha256 "2d6648faddbb2d4d4a63d3ae5eedc33321d0c70cdbf32e0ca3d2f7b4427aabb1"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-linux-arm64.tar.gz"
      sha256 "ed94c8ab6e0e910799f6d5b4af805dc1fb62c518db4252771e8e137f0185c881"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.1.4/ga-linux-amd64.tar.gz"
      sha256 "3ff0f8452c327ca90f2cc79e662c4dcfae87b6315cc6c86835eefdf8d4b33975"
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
