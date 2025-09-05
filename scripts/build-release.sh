#!/bin/bash

# Quick build script for testing releases locally
# This script builds the binaries and creates tar.gz files for testing

set -e

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 설정
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TAP_ROOT="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$TAP_ROOT")"
BINARY_NAME="ga"

echo -e "${GREEN}🔨 Quick build for testing${NC}"

# 테스트 버전 (기본값)
TEST_VERSION="0.1.0-test"
if [ -n "$1" ]; then
    TEST_VERSION="$1"
fi

# 출력 디렉토리
OUTPUT_DIR="$TAP_ROOT/dist/test"
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# 메인 프로젝트에서 빌드
cd "$PROJECT_ROOT"

echo -e "${YELLOW}Building binaries...${NC}"

# 빌드 실행
if [ -f "$PROJECT_ROOT/build.command" ]; then
    bash "$PROJECT_ROOT/build.command"
else
    echo -e "${RED}❌ build.command not found${NC}"
    exit 1
fi

# 압축 파일 생성
echo -e "${YELLOW}Creating test archives...${NC}"

PLATFORMS=(
    "darwin_amd64"
    "darwin_arm64"
)

for platform in "${PLATFORMS[@]}"; do
    BINARY_PATH="$PROJECT_ROOT/build/${platform}/${BINARY_NAME}"
    
    if [ -f "$BINARY_PATH" ]; then
        ARCHIVE_NAME="${BINARY_NAME}-${platform//_/-}.tar.gz"
        ARCHIVE_PATH="$OUTPUT_DIR/$ARCHIVE_NAME"
        
        echo "  Creating $ARCHIVE_NAME..."
        cd "$PROJECT_ROOT/build/${platform}"
        tar -czf "$ARCHIVE_PATH" "${BINARY_NAME}"
        
        # SHA256 체크섬 계산 및 출력
        if command -v shasum > /dev/null; then
            SHA256=$(shasum -a 256 "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        else
            SHA256=$(sha256sum "$ARCHIVE_PATH" | cut -d ' ' -f 1)
        fi
        
        echo "    SHA256: $SHA256"
    fi
done

echo -e "${GREEN}✅ Test build complete!${NC}"
echo ""
echo "Test archives created in: $OUTPUT_DIR"
echo ""
echo "To test the Formula locally:"
echo "1. Update the Formula with these SHA256 values"
echo "2. Run: brew install --build-from-source ./Formula/ga.rb"