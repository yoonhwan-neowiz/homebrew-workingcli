#!/bin/bash

# Change directory to the script's location
cd "$(dirname "$0")"

# Display the current working directory
echo "Current directory: $PWD"

# 인자 파싱
INSTALL_GLOBAL=false
for arg in "$@"; do
    case $arg in
        --install)
            INSTALL_GLOBAL=true
            shift
            ;;
        *)
            ;;
    esac
done

# 빌드 환경 설정
BINARY_NAME="ga"
BUILD_TIME=$(date +%Y%m%d_%H%M%S)
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")


# 테스트 버전 (기본값)
if [ -n "$VERSION" ]; then
    VERSION=$VERSION
else
    VERSION="0.1.0"
fi

# 색상 정의
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "=== Git Assistant 빌드 시작 ==="
echo "버전: ${VERSION}"
echo "빌드 시간: ${BUILD_TIME}"
echo "커밋 해시: ${COMMIT_HASH}"

# 빌드 디렉토리 생성 또는 초기화
BUILD_DIR="build"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 의존성 다운로드
echo "의존성 다운로드 중..."
go mod download

# 테스트 실행
# echo "테스트 실행 중..."
# go test ./... || { echo "테스트 실패"; exit 1; }

# 빌드 함수
build() {
    local OS=$1
    local ARCH=$2
    local OUTPUT_DIR="${BUILD_DIR}/${OS}_${ARCH}"
    local BINARY="${OUTPUT_DIR}/${BINARY_NAME}"
    
    # Windows의 경우 .exe 확장자 추가
    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY}.exe"
    fi

    echo "빌드 중: ${OS}/${ARCH}"
    
    # 디렉토리 생성
    mkdir -p ${OUTPUT_DIR}

    # 환경 변수 설정 및 빌드
    GOOS=${OS} GOARCH=${ARCH} go build \
        -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitHash=${COMMIT_HASH}" \
        -o ${BINARY}
    
    if [ $? -eq 0 ]; then
        # macOS의 경우 권한 설정
        if [ "$OS" = "darwin" ]; then
            echo "macOS 보안 설정 적용: ${BINARY}"
            chmod +x ${BINARY}
        fi
        echo -e "${GREEN}완료: ${BINARY}${NC}"
    else
        echo "Error: ${OS}/${ARCH} 빌드 실패"
        exit 1
    fi
}

# 다양한 플랫폼용 빌드
build "darwin" "amd64"  # macOS Intel
build "darwin" "arm64"  # macOS Apple Silicon
build "linux" "amd64"   # Linux Intel/AMD
build "linux" "arm64"   # Linux ARM64
build "windows" "amd64" # Windows Intel/AMD

# 현재 플랫폼용 심볼릭 링크 생성
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
CURRENT_ARCH=$(uname -m)
if [ "${CURRENT_ARCH}" = "x86_64" ]; then
    CURRENT_ARCH="amd64"
elif [ "${CURRENT_ARCH}" = "arm64" ] || [ "${CURRENT_ARCH}" = "aarch64" ]; then
    CURRENT_ARCH="arm64"
fi

# 현재 플랫폼용 바이너리 심볼릭 링크
CURRENT_BINARY="${BUILD_DIR}/${CURRENT_OS}_${CURRENT_ARCH}/${BINARY_NAME}"
if [ -f "${CURRENT_BINARY}" ]; then
    ln -sf "${CURRENT_BINARY}" "${BINARY_NAME}"
    echo -e "${GREEN}심볼릭 링크 생성: ${BINARY_NAME} -> ${CURRENT_BINARY}${NC}"
fi

echo "=== 빌드 완료 ==="
echo "빌드된 파일은 ${BUILD_DIR} 디렉토리에 있습니다."
echo ""
echo "빌드된 파일 목록:"
ls -lh ${BUILD_DIR}/*

# 실행 권한 부여
chmod +x "${BINARY_NAME}"

# --install 인자가 있을 때만 전역 설치
if [ "$INSTALL_GLOBAL" = true ]; then
    echo ""
    echo "=== 전역 설치 ==="
    if sudo ln -sf "$(pwd)/${BINARY_NAME}" /usr/local/bin/${BINARY_NAME}; then
        echo -e "${GREEN}시스템 전역 명령어 설치 완료: ${BINARY_NAME}${NC}"
        echo "이제 터미널 어디서나 '${BINARY_NAME}' 명령어를 사용할 수 있습니다."
    else
        echo "Error: 시스템 전역 명령어 설치 실패. sudo 권한이 필요합니다."
    fi
else
    echo ""
    echo "=== 로컬 빌드 완료 ==="
    echo -e "${GREEN}현재 디렉토리에서 './${BINARY_NAME}'로 실행할 수 있습니다.${NC}"
    echo "전역 설치를 원하시면 './build.command --install'을 실행하세요."
fi 