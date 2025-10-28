class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "0.2.6"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.6/ga-darwin-arm64.tar.gz"
      sha256 "389bd077fc0cf5fec78d874615663fb65c0b0323eb1170e9e8af7e67650048c4"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.6/ga-darwin-amd64.tar.gz"
      sha256 "01159e8fdbcaa0c9fd188bf12a00fa921859dc4224f978e190170d3d697bd6a1"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.6/ga-linux-arm64.tar.gz"
      sha256 "64a98dbb61efc44e36115a742cbff80d1aea743457c783f284039053a2c75dfc"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/v0.2.6/ga-linux-amd64.tar.gz"
      sha256 "48c81d93ed8c1f2964894509c9cca7a1ecd36e103e08407d55f9091bbb12bd59"
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
