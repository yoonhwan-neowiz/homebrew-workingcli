#!/bin/bash

# ============================================================================
# Homebrew 전체 릴리스 자동화 스크립트
# 사용법: ./release-all.sh <version>
# 예제: ./release-all.sh 0.1.3
# 
# 이 스크립트는 다음 작업을 자동으로 수행합니다:
# 1. 빌드 실행 (모든 플랫폼)
# 2. tar.gz 아카이브 생성
# 3. SHA256 체크섬 계산 및 Formula 업데이트
# 4. Git 커밋 및 Push
# 5. GitHub Release 생성
# ============================================================================

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 설정
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TAP_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$TAP_ROOT")"
FORMULA_FILE="$TAP_ROOT/Formula/ga.rb"
BINARY_NAME="ga"

# 버전 검증
if [ -z "$1" ]; then
    echo -e "${RED}❌ 버전이 필요합니다${NC}"
    echo -e "${YELLOW}사용법: ./release-all.sh <version>${NC}"
    echo "예제: ./release-all.sh 0.1.3"
    exit 1
fi

VERSION="$1"
TAG_NAME="v${VERSION}"

# 버전 형식 검증
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}❌ 잘못된 버전 형식입니다. 예: 0.1.3${NC}"
    exit 1
fi

echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║${NC}  ${GREEN}🚀 Homebrew 릴리스 자동화 v${VERSION}${NC}  ${CYAN}║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# GitHub CLI 확인
if ! command -v gh &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI (gh)가 설치되지 않았습니다${NC}"
    echo "설치: brew install gh"
    exit 1
fi

# GitHub 인증 확인
if ! gh auth status &> /dev/null; then
    echo -e "${RED}❌ GitHub CLI가 인증되지 않았습니다${NC}"
    echo "실행: gh auth login"
    exit 1
fi

# ============================================================================
# 1단계: 빌드 준비
# ============================================================================
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📁 [1/5] 릴리스 디렉토리 준비${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

RELEASE_DIR="$TAP_ROOT/dist/${VERSION}"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"
echo -e "${GREEN}✓${NC} 릴리스 디렉토리 생성: $RELEASE_DIR"

# ============================================================================
# 2단계: 빌드 실행
# ============================================================================
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}🔨 [2/5] 바이너리 빌드${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

cd "$PROJECT_ROOT"

# 빌드 스크립트 실행
if [ ! -f "$PROJECT_ROOT/build.command" ]; then
    echo -e "${RED}❌ build.command를 찾을 수 없습니다${NC}"
    exit 1
fi

# VERSION 환경변수로 빌드 실행
export VERSION="$VERSION"
echo -e "빌드 버전: ${GREEN}${VERSION}${NC}"
bash "$PROJECT_ROOT/build.command"
unset VERSION

echo -e "${GREEN}✓${NC} 빌드 완료"

# ============================================================================
# 3단계: 아카이브 생성 및 SHA256 계산
# ============================================================================
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📦 [3/5] 릴리스 아카이브 생성${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

PLATFORMS=(
    "darwin_amd64"
    "darwin_arm64"
    "linux_amd64"
    "linux_arm64"
)

# SHA256 체크섬 저장 변수
declare -A SHA256_CHECKSUMS

for platform in "${PLATFORMS[@]}"; do
    BINARY_PATH="$PROJECT_ROOT/build/${platform}/${BINARY_NAME}"
    
    if [ -f "$BINARY_PATH" ]; then
        ARCHIVE_NAME="${BINARY_NAME}-${platform//_/-}.tar.gz"
        ARCHIVE_PATH="$RELEASE_DIR/$ARCHIVE_NAME"
        
        echo -n "  📄 $ARCHIVE_NAME ... "
        cd "$PROJECT_ROOT/build/${platform}"
        tar -czf "$ARCHIVE_PATH" "${BINARY_NAME}"
        
        # SHA256 체크섬 계산
        if command -v shasum > /dev/null; then
            SHA256=$(shasum -a 256 "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        else
            SHA256=$(sha256sum "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        fi
        
        SHA256_CHECKSUMS[$platform]="$SHA256"
        echo -e "${GREEN}✓${NC}"
        echo "     SHA256: ${CYAN}${SHA256}${NC}"
    else
        echo -e "  ${YELLOW}⚠${NC} ${platform} 바이너리를 찾을 수 없습니다"
    fi
done

cd "$TAP_ROOT"

# ============================================================================
# 4단계: Formula 파일 업데이트
# ============================================================================
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📝 [4/5] Formula 파일 업데이트 및 Git 작업${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Formula 백업
cp "$FORMULA_FILE" "$FORMULA_FILE.bak"

# Formula 파일 생성 (템플릿 사용)
cat > "$FORMULA_FILE" << EOF
class Ga < Formula
  desc "Git Assistant - Smart Git workflow optimizer for large repositories"
  homepage "https://github.com/yoonhwan-neowiz/WorkingCli"
  version "${VERSION}"
  license "MIT"
  
  # macOS 플랫폼별 설정
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/${TAG_NAME}/ga-darwin-arm64.tar.gz"
      sha256 "${SHA256_CHECKSUMS[darwin_arm64]}"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/${TAG_NAME}/ga-darwin-amd64.tar.gz"
      sha256 "${SHA256_CHECKSUMS[darwin_amd64]}"
    end
  end

  # Linux 플랫폼별 설정
  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/${TAG_NAME}/ga-linux-arm64.tar.gz"
      sha256 "${SHA256_CHECKSUMS[linux_arm64]}"
    else
      url "https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/download/${TAG_NAME}/ga-linux-amd64.tar.gz"
      sha256 "${SHA256_CHECKSUMS[linux_amd64]}"
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
EOF

echo -e "${GREEN}✓${NC} Formula 파일 업데이트 완료"

# Git 커밋
echo -n "  Git 커밋 중 ... "
git add "$FORMULA_FILE"
git commit -m "feat: Release version ${VERSION}

- Updated Formula to version ${VERSION}
- Updated download URLs
- Updated SHA256 checksums for all platforms" > /dev/null 2>&1
echo -e "${GREEN}✓${NC}"

# Git Push
echo -n "  Git Push 중 ... "
git push origin main > /dev/null 2>&1
echo -e "${GREEN}✓${NC}"

# ============================================================================
# 5단계: GitHub Release 생성
# ============================================================================
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}🚀 [5/5] GitHub Release 생성${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# 기존 릴리스 확인 및 삭제
if gh release view "${TAG_NAME}" &> /dev/null; then
    echo -e "${YELLOW}⚠${NC} 기존 릴리스 ${TAG_NAME}를 삭제합니다..."
    gh release delete "${TAG_NAME}" --yes > /dev/null 2>&1
fi

# 릴리스 노트 생성
RELEASE_NOTES="## 🚀 Git Assistant v${VERSION}

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
"

# SHA256 체크섬 추가
for platform in "${PLATFORMS[@]}"; do
    if [ -n "${SHA256_CHECKSUMS[$platform]}" ]; then
        RELEASE_NOTES="${RELEASE_NOTES}- ${platform//_/-}: \`${SHA256_CHECKSUMS[$platform]}\`
"
    fi
done

# 릴리스 생성
echo -n "  GitHub Release 생성 중 ... "
cd "$TAP_ROOT"

# 릴리스 파일 목록 생성
RELEASE_FILES=""
for platform in "${PLATFORMS[@]}"; do
    ARCHIVE_NAME="${BINARY_NAME}-${platform//_/-}.tar.gz"
    if [ -f "$RELEASE_DIR/$ARCHIVE_NAME" ]; then
        RELEASE_FILES="$RELEASE_FILES $RELEASE_DIR/$ARCHIVE_NAME"
    fi
done

# GitHub Release 생성
if gh release create "${TAG_NAME}" \
    --title "${TAG_NAME} - Git Assistant for Large Repositories" \
    --notes "${RELEASE_NOTES}" \
    $RELEASE_FILES > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${RED}✗${NC}"
    echo -e "${RED}❌ GitHub Release 생성 실패${NC}"
    exit 1
fi

# ============================================================================
# 완료 메시지
# ============================================================================
echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║${NC}           ${GREEN}✅ 릴리스 프로세스 완료!${NC}            ${CYAN}║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}📦 릴리스 버전:${NC} ${VERSION}"
echo -e "${GREEN}🔗 GitHub Release:${NC} https://github.com/yoonhwan-neowiz/homebrew-workingcli/releases/tag/${TAG_NAME}"
echo ""
echo -e "${YELLOW}테스트 방법:${NC}"
echo -e "  ${CYAN}brew update${NC}"
echo -e "  ${CYAN}brew upgrade ga${NC}"
echo ""
echo -e "${YELLOW}새로 설치:${NC}"
echo -e "  ${CYAN}brew tap yoonhwan-neowiz/workingcli${NC}"
echo -e "  ${CYAN}brew install ga${NC}"
echo ""

# 백업 파일 삭제
rm -f "$FORMULA_FILE.bak"

echo -e "${GREEN}✨ 모든 작업이 성공적으로 완료되었습니다!${NC}"