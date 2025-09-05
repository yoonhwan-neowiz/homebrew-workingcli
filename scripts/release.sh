#!/bin/bash

# Release automation script for homebrew-workingcli
# This script builds, packages, and updates the Formula for new releases

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 설정
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TAP_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$TAP_ROOT")"
FORMULA_FILE="$TAP_ROOT/Formula/ga.rb"
BINARY_NAME="ga"

# 버전 정보 가져오기 (인자로 받거나 현재 버전 사용)
if [ -z "$1" ]; then
    echo -e "${YELLOW}Usage: ./release.sh <version>${NC}"
    echo "Example: ./release.sh 0.1.1"
    exit 1
fi

VERSION="$1"
TAG_NAME="v${VERSION}"

echo -e "${GREEN}🚀 Starting release process for version ${VERSION}${NC}"

# 1. 작업 디렉토리 준비
echo -e "${YELLOW}📁 Preparing release directory...${NC}"
RELEASE_DIR="$TAP_ROOT/dist/${VERSION}"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# 2. 메인 프로젝트에서 빌드 실행
echo -e "${YELLOW}🔨 Building binaries...${NC}"
cd "$PROJECT_ROOT"

# 빌드 스크립트 실행
if [ -f "$PROJECT_ROOT/build.command" ]; then
    # build.command의 버전 정보 임시 업데이트
    sed -i.bak "s/VERSION=\".*\"/VERSION=\"${VERSION}\"/" "$PROJECT_ROOT/build.command"
    bash "$PROJECT_ROOT/build.command"
    # 원본 복원
    mv "$PROJECT_ROOT/build.command.bak" "$PROJECT_ROOT/build.command"
else
    echo -e "${RED}❌ build.command not found${NC}"
    exit 1
fi

# 3. 각 플랫폼별 압축 파일 생성
echo -e "${YELLOW}📦 Creating release archives...${NC}"

PLATFORMS=(
    "darwin_amd64"
    "darwin_arm64"
    "linux_amd64"
    "linux_arm64"
)

# SHA256 체크섬을 저장할 변수들 개별 선언
SHA256_darwin_amd64=""
SHA256_darwin_arm64=""
SHA256_linux_amd64=""
SHA256_linux_arm64=""

for platform in "${PLATFORMS[@]}"; do
    BINARY_PATH="$PROJECT_ROOT/build/${platform}/${BINARY_NAME}"
    
    if [ -f "$BINARY_PATH" ]; then
        ARCHIVE_NAME="${BINARY_NAME}-${platform//_/-}.tar.gz"
        ARCHIVE_PATH="$RELEASE_DIR/$ARCHIVE_NAME"
        
        # tar.gz 생성
        echo "  Creating $ARCHIVE_NAME..."
        cd "$PROJECT_ROOT/build/${platform}"
        tar -czf "$ARCHIVE_PATH" "${BINARY_NAME}"
        
        # SHA256 체크섬 계산
        if command -v shasum > /dev/null; then
            SHA256=$(shasum -a 256 "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        else
            SHA256=$(sha256sum "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        fi
        
        # platform별로 변수에 저장
        eval "SHA256_${platform}=\"$SHA256\""
        echo "    SHA256: $SHA256"
    else
        echo -e "${YELLOW}  Skipping ${platform} (binary not found)${NC}"
    fi
done

cd "$TAP_ROOT"

# 4. Formula 파일 업데이트
echo -e "${YELLOW}📝 Updating Formula file...${NC}"

# Formula 백업
cp "$FORMULA_FILE" "$FORMULA_FILE.bak"

# 버전 업데이트
sed -i '' "s/version \".*\"/version \"${VERSION}\"/" "$FORMULA_FILE"

# URL 업데이트
sed -i '' "s|download/v[0-9.]*|download/${TAG_NAME}|g" "$FORMULA_FILE"

# SHA256 체크섬 업데이트
if [ -n "${SHA256_darwin_arm64}" ]; then
    sed -i '' "s/PENDING_ARM64_SHA256/${SHA256_darwin_arm64}/" "$FORMULA_FILE"
fi

if [ -n "${SHA256_darwin_amd64}" ]; then
    sed -i '' "s/PENDING_AMD64_SHA256/${SHA256_darwin_amd64}/" "$FORMULA_FILE"
fi

if [ -n "${SHA256_linux_arm64}" ]; then
    sed -i '' "s/PENDING_LINUX_ARM64_SHA256/${SHA256_linux_arm64}/" "$FORMULA_FILE"
fi

if [ -n "${SHA256_linux_amd64}" ]; then
    sed -i '' "s/PENDING_LINUX_AMD64_SHA256/${SHA256_linux_amd64}/" "$FORMULA_FILE"
fi

# 5. Git 커밋
echo -e "${YELLOW}📝 Committing Formula changes...${NC}"
git add "$FORMULA_FILE"
git commit -m "feat: Release version ${VERSION}

- Updated Formula to version ${VERSION}
- Updated download URLs
- Updated SHA256 checksums"

# 6. GitHub Release 생성 안내
echo -e "${GREEN}✅ Release preparation complete!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Push changes to homebrew tap:"
echo "   git push origin main"
echo ""
echo "2. Create GitHub Release in the main project:"
echo "   cd $PROJECT_ROOT"
echo "   git tag -a ${TAG_NAME} -m \"Release version ${VERSION}\""
echo "   git push origin ${TAG_NAME}"
echo ""
echo "3. Upload the following files to GitHub Release:"
for platform in "${PLATFORMS[@]}"; do
    ARCHIVE_NAME="${BINARY_NAME}-${platform//_/-}.tar.gz"
    if [ -f "$RELEASE_DIR/$ARCHIVE_NAME" ]; then
        echo "   - $RELEASE_DIR/$ARCHIVE_NAME"
    fi
done
echo ""
echo "4. Test the installation:"
echo "   brew tap yoonhwan-neowiz/workingcli"
echo "   brew install ga"

# SHA256 체크섬 정보 저장
echo -e "${YELLOW}📋 SHA256 Checksums:${NC}"
echo "================================="
for platform in "${PLATFORMS[@]}"; do
    eval "SHA256_VAL=\$SHA256_${platform}"
    if [ -n "${SHA256_VAL}" ]; then
        echo "${platform}: ${SHA256_VAL}"
    fi
done
echo "================================="