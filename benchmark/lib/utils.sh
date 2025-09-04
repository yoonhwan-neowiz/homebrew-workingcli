#!/bin/bash

# Git Optimization Benchmark - Utility Functions
# 로그, 타이머, 변환 등 유틸리티 함수

# ==================== 로그 출력 함수 ====================

# 일반 로그 출력
log_output() {
    if [ "$QUIET_MODE" = false ]; then
        echo "$@" | tee -a "$LOG_FILE"
    else
        echo "$@" >> "$LOG_FILE"
    fi
}

# 색상 로그 출력
log_color() {
    if [ "$QUIET_MODE" = false ]; then
        echo -e "$@" | tee -a "$LOG_FILE"
    else
        echo -e "$@" >> "$LOG_FILE"
    fi
}

# 진행 상황 표시
show_progress() {
    local step=$1
    local total=$2
    local desc=$3
    log_color "${CYAN}[${step}/${total}] ${desc}${NC}"
}

# ==================== 시간 측정 함수 ====================

# 시간 측정 시작
start_timer() {
    START_TIME=$(date +%s)
}

# 시간 측정 종료 및 반환
end_timer() {
    local END_TIME=$(date +%s)
    local DURATION=$((END_TIME - START_TIME))
    echo $DURATION
}

# 시간 형식 변환 (초 → 분:초)
format_time() {
    local seconds=$1
    local mins=$((seconds / 60))
    local secs=$((seconds % 60))
    
    if [ $mins -gt 0 ]; then
        printf "%d분 %d초" $mins $secs
    else
        printf "%d초" $secs
    fi
}

# ==================== 크기 변환 함수 ====================

# 바이트를 사람이 읽기 쉬운 형식으로 변환
bytes_to_human() {
    local bytes=$1
    if [ $bytes -lt 1024 ]; then
        echo "${bytes}B"
    elif [ $bytes -lt 1048576 ]; then
        echo "$(echo "scale=1; $bytes/1024" | bc)KB"
    elif [ $bytes -lt 1073741824 ]; then
        echo "$(echo "scale=1; $bytes/1048576" | bc)MB"
    else
        echo "$(echo "scale=2; $bytes/1073741824" | bc)GB"
    fi
}

# KB 단위를 바이트로 변환
kb_to_bytes() {
    echo $(($1 * 1024))
}

# ==================== 경로 관련 함수 ====================

# 절대 경로 확인
get_absolute_path() {
    local path=$1
    if [ -d "$path" ]; then
        echo "$(cd "$path" && pwd)"
    elif [ -f "$path" ]; then
        echo "$(cd "$(dirname "$path")" && pwd)/$(basename "$path")"
    else
        echo "$path"
    fi
}

# ==================== 확인 및 프롬프트 ====================

# 사용자 확인
confirm() {
    local prompt="${1:-계속하시겠습니까?}"
    local default="${2:-n}"
    
    if [ "$QUIET_MODE" = true ]; then
        [ "$default" = "y" ] && return 0 || return 1
    fi
    
    local answer
    read -p "$prompt (y/n) [${default}]: " answer
    answer=${answer:-$default}
    
    [ "$answer" = "y" ] || [ "$answer" = "Y" ]
}

# ==================== 오류 처리 ====================

# 오류 메시지 출력 후 종료
die() {
    log_color "${RED}❌ $1${NC}" >&2
    exit ${2:-1}
}

# 경고 메시지 출력
warn() {
    log_color "${YELLOW}⚠️  $1${NC}" >&2
}

# 성공 메시지 출력
success() {
    log_color "${GREEN}✓ $1${NC}"
}

# ==================== 백업 관련 ====================

# 파일 백업
backup_file() {
    local file=$1
    local backup="${file}.backup.${DATE_STR}"
    
    if [ -f "$file" ]; then
        cp "$file" "$backup"
        log_output "백업 생성: $backup"
    fi
}

# ==================== 시스템 정보 ====================

# 디스크 사용량 확인
check_disk_space() {
    local path="${1:-.}"
    local available=$(df -k "$path" | awk 'NR==2 {print $4}')
    echo $((available * 1024))  # 바이트로 반환
}

# CPU 코어 수 확인
get_cpu_cores() {
    if [ "$(uname)" = "Darwin" ]; then
        sysctl -n hw.ncpu
    else
        nproc
    fi
}

# Export functions for use in other modules
export -f log_output log_color show_progress
export -f start_timer end_timer format_time
export -f bytes_to_human kb_to_bytes
export -f get_absolute_path confirm
export -f die warn success
export -f backup_file check_disk_space get_cpu_cores