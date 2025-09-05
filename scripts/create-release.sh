#!/bin/bash

# GitHub Release 생성 스크립트
# 사용법: ./create-release.sh

VERSION="0.1.0"
RELEASE_DIR="$(cd "$(dirname "$0")/../dist/${VERSION}" && pwd)"

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Creating GitHub Release for v${VERSION}${NC}"

# 파일 존재 확인
if [ ! -d "$RELEASE_DIR" ]; then
    echo -e "${RED}❌ Release directory not found: $RELEASE_DIR${NC}"
    exit 1
fi

# GitHub CLI 확인
if ! command -v gh &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI (gh) is not installed${NC}"
    echo "Install with: brew install gh"
    exit 1
fi

# 인증 확인
if ! gh auth status &> /dev/null; then
    echo -e "${YELLOW}⚠️  GitHub CLI not authenticated${NC}"
    echo "Please run: gh auth login"
    exit 1
fi

# Release 생성
echo -e "${YELLOW}📦 Creating release with binaries...${NC}"

gh release create "v${VERSION}" \
  --title "v${VERSION} - Git Assistant for Large Repositories" \
  --notes "## 🚀 First Release of Git Assistant (ga)

### Features
- 🎯 Smart Git workflow optimizer for large repositories (50GB+)
- ⚡ Performance-focused commands for massive codebases
- 🔄 Shallow clone and sparse-checkout management
- 🌳 Intelligent branch switching with cleanup
- 💻 Cross-platform support (macOS Intel/ARM, Linux)

### Installation

\`\`\`bash
# via Homebrew
brew tap yoonhwan-neowiz/workingcli
brew install ga
\`\`\`

### Quick Start

\`\`\`bash
ga                    # Interactive staging
ga commit             # AI-powered commit message
ga opt quick status   # Check repository optimization
ga opt quick to-slim  # Optimize large repository
\`\`\`

### SHA256 Checksums
- darwin_amd64: \`9b4ee989d0f1b1a441368f97d4c7af680d7f964295115dabcc98862b56340cc5\`
- darwin_arm64: \`055857e9fd878764b4e660a554b91b073acb6a8d2c4c4ff71ee545a4e471ea62\`
- linux_amd64: \`075fc38c200857047ef9b1f1c03d3814098b102a4557247cd0c678770a585c53\`
- linux_arm64: \`18574392309448578c60a87789092fb73865db0084bfd73fbaf5fb5f9ca520e7\`" \
  "$RELEASE_DIR/ga-darwin-amd64.tar.gz" \
  "$RELEASE_DIR/ga-darwin-arm64.tar.gz" \
  "$RELEASE_DIR/ga-linux-amd64.tar.gz" \
  "$RELEASE_DIR/ga-linux-arm64.tar.gz"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Release created successfully!${NC}"
    echo ""
    echo "View release at: https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/tag/v${VERSION}"
else
    echo -e "${RED}❌ Failed to create release${NC}"
    exit 1
fi