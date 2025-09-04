#!/bin/bash

# Git Optimization Benchmark - Base Configuration
# 기본 설정, 색상 정의, 경로 설정

# ==================== 기본 설정 ====================
set -e  # 오류 시 즉시 중단

# 저장소 설정
TEST_REPO="${TEST_REPO:-git@git.nwz.kr:gamfs4/designb4}"
TEST_DIR="${TEST_DIR:-$HOME/Work/DesignB4_test}"

# 스크립트 경로
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LIB_DIR="$SCRIPT_DIR/lib"
SCENARIOS_DIR="$SCRIPT_DIR/scenarios"
GA_CMD="$(dirname "$SCRIPT_DIR")/ga"  # 상위 디렉토리의 ga 명령어

# 로그 및 결과 파일 경로
DATE_STR=$(date +"%Y%m%d_%H%M%S")
LOG_DIR="${LOG_DIR:-$SCRIPT_DIR/logs}"
RESULT_DIR="${RESULT_DIR:-$SCRIPT_DIR/results}"
LOG_FILE="${LOG_FILE:-$LOG_DIR/benchmark_${DATE_STR}.log}"
JSONL_FILE="${JSONL_FILE:-$LOG_DIR/benchmark_${DATE_STR}.jsonl}"
SUMMARY_FILE="${SUMMARY_FILE:-$RESULT_DIR/optimization_report_${DATE_STR}.txt}"

# ==================== 색상 정의 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# ==================== 기본 옵션 ====================
# 명령행에서 오버라이드 가능
MODE="${MODE:-fast}"  # fast 또는 accurate
SKIP_SLIM="${SKIP_SLIM:-false}"
SKIP_SPARSE="${SKIP_SPARSE:-false}"
BRANCHES="${BRANCHES:-master}"  # branch-scope에 사용할 브랜치
QUIET_MODE="${QUIET_MODE:-false}"
SOURCE_TEMPLATE="${SOURCE_TEMPLATE:-}"  # 템플릿 소스 디렉토리
USE_CLONE="${USE_CLONE:-false}"     # clone 방식 사용 여부

# ==================== 전역 변수 ====================
# 측정값 추적용
TOTAL_SAVED=0
LAST_GIT_STORE=0
LAST_WORKTREE=0
LAST_TOTAL=0

# ==================== 디렉토리 생성 ====================
ensure_directories() {
    mkdir -p "$LOG_DIR" "$RESULT_DIR"
}

# ==================== 환경 검증 ====================
validate_environment() {
    # ga 명령어 확인
    if [ ! -f "$GA_CMD" ]; then
        echo "❌ ga 명령어를 찾을 수 없습니다: $GA_CMD"
        exit 1
    fi
    
    # rusync/rsync 확인 (복사 모드일 때만)
    if [ "$USE_CLONE" = false ]; then
        if ! command -v rusync &> /dev/null && ! command -v rsync &> /dev/null; then
            echo "⚠️  rusync 또는 rsync가 필요합니다"
            echo "   rusync 설치: cargo install rusync"
            echo "   또는 --use-clone 옵션을 사용하세요"
            exit 1
        fi
    fi
}

# Export all variables for use in other modules
export TEST_REPO TEST_DIR SCRIPT_DIR LIB_DIR SCENARIOS_DIR GA_CMD
export DATE_STR LOG_DIR RESULT_DIR LOG_FILE JSONL_FILE SUMMARY_FILE
export RED GREEN YELLOW BLUE CYAN MAGENTA NC
export MODE SKIP_SLIM SKIP_SPARSE BRANCHES QUIET_MODE SOURCE_TEMPLATE USE_CLONE
export TOTAL_SAVED LAST_GIT_STORE LAST_WORKTREE LAST_TOTAL