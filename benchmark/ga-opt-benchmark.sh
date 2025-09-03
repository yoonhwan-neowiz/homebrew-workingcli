#!/bin/bash

# Git Optimization Benchmark Script
# 목적: ga 명령어를 활용한 단계별 Git 최적화 효과 측정
# 타겟: git@git.nwz.kr:gamfs4/designb4 → ~/Work/DesignB4_test

set -e  # 오류 시 즉시 중단

# ==================== 설정 ====================
TEST_REPO="git@git.nwz.kr:gamfs4/designb4"
TEST_DIR="$HOME/Work/DesignB4_test"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GA_CMD="$(dirname "$SCRIPT_DIR")/ga"  # 상위 디렉토리의 ga 명령어

# 로그 및 결과 파일 경로
DATE_STR=$(date +"%Y%m%d_%H%M%S")
LOG_DIR="$SCRIPT_DIR/logs"
RESULT_DIR="$SCRIPT_DIR/results"
LOG_FILE="$LOG_DIR/benchmark_${DATE_STR}.log"
JSONL_FILE="$LOG_DIR/benchmark_${DATE_STR}.jsonl"
SUMMARY_FILE="$RESULT_DIR/optimization_report_${DATE_STR}.txt"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# 기본 옵션
MODE="fast"  # fast 또는 accurate
SKIP_SLIM=false
SKIP_SPARSE=false
BRANCHES="master"  # branch-scope에 사용할 브랜치
QUIET_MODE=false

# ==================== 옵션 파싱 ====================
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode)
            MODE="$2"
            shift 2
            ;;
        --skip-slim)
            SKIP_SLIM=true
            shift
            ;;
        --skip-sparse)
            SKIP_SPARSE=true
            shift
            ;;
        --branches)
            BRANCHES="$2"
            shift 2
            ;;
        --quiet|-q)
            QUIET_MODE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --mode [fast|accurate]  측정 모드 (기본: fast)"
            echo "  --skip-slim            to-slim 단계 제외"
            echo "  --skip-sparse          sparse 단계 제외"
            echo "  --branches \"list\"      branch-scope 브랜치 목록"
            echo "  --quiet, -q            조용한 모드"
            echo "  --help, -h             도움말 표시"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# ==================== 유틸리티 함수 ====================

# 로그 출력 함수
log_output() {
    if [ "$QUIET_MODE" = false ]; then
        echo "$@" | tee -a "$LOG_FILE"
    else
        echo "$@" >> "$LOG_FILE"
    fi
}

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

# ==================== 측정 함수 ====================

# 저장소 크기 측정 (Git objects)
measure_git_store() {
    local repo_path=${1:-.}
    local total_size=0
    
    # 메인 저장소 .git 크기
    if [ -d "$repo_path/.git" ]; then
        local main_size=$(du -sk "$repo_path/.git" 2>/dev/null | cut -f1)
        total_size=$((total_size + main_size * 1024))
    fi
    
    # 서브모듈 .git 크기
    if [ -f "$repo_path/.gitmodules" ]; then
        while IFS= read -r submodule_path; do
            if [ -d "$repo_path/.git/modules/$submodule_path" ]; then
                local sub_size=$(du -sk "$repo_path/.git/modules/$submodule_path" 2>/dev/null | cut -f1)
                total_size=$((total_size + sub_size * 1024))
            fi
        done < <(git config --file "$repo_path/.gitmodules" --get-regexp path | awk '{print $2}')
    fi
    
    echo $total_size
}

# 워킹트리 크기 측정 (.git 제외)
measure_worktree() {
    local repo_path=${1:-.}
    local total_size=0
    
    # 전체 크기에서 .git 제외
    if [ -d "$repo_path" ]; then
        local full_size=$(du -sk "$repo_path" 2>/dev/null | cut -f1)
        local git_size=$(du -sk "$repo_path/.git" 2>/dev/null | cut -f1 || echo 0)
        total_size=$(((full_size - git_size) * 1024))
    fi
    
    echo $total_size
}

# 전체 측정 및 JSONL 기록
measure_step() {
    local step_name=$1
    local step_number=$2
    local duration=${3:-0}
    
    local git_store=$(measure_git_store "$TEST_DIR")
    local worktree=$(measure_worktree "$TEST_DIR")
    local total_size=$((git_store + worktree))
    
    # JSONL 형식으로 기록
    local json_entry=$(cat <<EOF
{"step": $step_number, "name": "$step_name", "timestamp": "$(date -Iseconds)", "git_store": $git_store, "worktree": $worktree, "total": $total_size, "duration": $duration}
EOF
)
    echo "$json_entry" >> "$JSONL_FILE"
    
    # 전역 변수에 저장 (델타 계산용)
    LAST_GIT_STORE=$git_store
    LAST_WORKTREE=$worktree
    LAST_TOTAL=$total_size
    
    # 사람이 읽기 쉬운 형식으로 로그
    log_output "  Git Store: $(bytes_to_human $git_store)"
    log_output "  Worktree:  $(bytes_to_human $worktree)"
    log_output "  Total:     $(bytes_to_human $total_size)"
    if [ $duration -gt 0 ]; then
        log_output "  Duration:  ${duration}s"
    fi
}

# 델타 계산 및 출력
calculate_delta() {
    local prev_git=$1
    local curr_git=$2
    local prev_wt=$3
    local curr_wt=$4
    local prev_total=$5
    local curr_total=$6
    
    local delta_git=$((prev_git - curr_git))
    local delta_wt=$((prev_wt - curr_wt))
    local delta_total=$((prev_total - curr_total))
    
    local percent_git=0
    local percent_wt=0
    local percent_total=0
    
    if [ $prev_git -gt 0 ]; then
        percent_git=$(echo "scale=1; $delta_git * 100 / $prev_git" | bc)
    fi
    if [ $prev_wt -gt 0 ]; then
        percent_wt=$(echo "scale=1; $delta_wt * 100 / $prev_wt" | bc)
    fi
    if [ $prev_total -gt 0 ]; then
        percent_total=$(echo "scale=1; $delta_total * 100 / $prev_total" | bc)
    fi
    
    log_color "${GREEN}  절감 효과:${NC}"
    log_output "    Git Store: $(bytes_to_human $delta_git) (${percent_git}%)"
    log_output "    Worktree:  $(bytes_to_human $delta_wt) (${percent_wt}%)"
    log_output "    Total:     $(bytes_to_human $delta_total) (${percent_total}%)"
    
    # 전역 누적 변수 업데이트
    TOTAL_SAVED=$((TOTAL_SAVED + delta_total))
    
    # 효율 계산 (MB/s)
    if [ ${7:-0} -gt 0 ]; then
        local efficiency=$(echo "scale=1; $delta_total / 1048576 / $7" | bc)
        log_output "    효율:      ${efficiency} MB/s"
    fi
}

# ==================== 최적화 단계 함수 ====================

# 저장소 초기화 (ga 명령어 활용)
init_repository() {
    log_color "${YELLOW}>>> 저장소 초기화${NC}"
    
    # 기존 디렉토리 삭제
    if [ -d "$TEST_DIR" ]; then
        log_output "기존 테스트 디렉토리 삭제 중..."
        rm -rf "$TEST_DIR"
    fi
    
    # 일반 클론 (ga로 최적화는 클론 후에)
    log_output "저장소 클론 중: $TEST_REPO"
    log_color "${CYAN}상세 진행 상황이 표시됩니다...${NC}"
    start_timer
    GIT_TRACE=1 GIT_CURL_VERBOSE=1 git clone \
        --progress \
        --verbose \
        "$TEST_REPO" "$TEST_DIR" 2>&1 | tee -a "$LOG_FILE"
    
    cd "$TEST_DIR"
    
    # 서브모듈 초기화 (ga opt submodule update -f 사용)
    log_output "서브모듈 초기화 중 (ga opt submodule update -f)..."
    log_color "${CYAN}서브모듈 진행 상황이 표시됩니다...${NC}"
    GIT_TRACE=1 $GA_CMD opt submodule update -f 2>&1 | tee -a "$LOG_FILE" || {
        log_color "${YELLOW}⚠ 서브모듈 초기화 실패 - 계속 진행${NC}"
    }
    
    # pull 설정 (rebase 대신 merge)
    git config pull.rebase false
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ 저장소 초기화 완료 (${duration}s)${NC}"
    
    return $duration
}

# Shallow 최적화
apply_shallow() {
    local depth=${1:-1}
    log_color "${YELLOW}>>> Shallow 최적화 (depth=$depth)${NC}"
    
    start_timer
    
    # 메인 저장소
    log_output "메인 저장소 shallow 적용 중..."
    $GA_CMD opt quick shallow $depth -q 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 shallow 적용 중..."
    $GA_CMD opt submodule shallow $depth -q 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Shallow 완료 (${duration}s)${NC}"
    
    return $duration
}

# Branch Scope 최적화
apply_branch_scope() {
    log_color "${YELLOW}>>> Branch Scope 최적화${NC}"
    
    start_timer
    
    # 현재 브랜치 확인
    local current_branch=$(git branch --show-current)
    local scope_branches="$BRANCHES"
    
    # 현재 브랜치가 목록에 없으면 추가
    if [[ ! " $scope_branches " =~ " $current_branch " ]]; then
        scope_branches="$scope_branches $current_branch"
    fi
    
    log_output "적용할 브랜치: $scope_branches"
    
    # 메인 저장소
    log_output "메인 저장소 branch-scope 적용 중..."
    $GA_CMD opt quick set-branch-scope $scope_branches 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 branch-scope 적용 중..."
    $GA_CMD opt submodule set-branch-scope $scope_branches 2>&1 | tee -a "$LOG_FILE"
    
    # Accurate 모드에서만 prune 실행
    if [ "$MODE" = "accurate" ]; then
        log_output "Prune 실행 중 (accurate 모드)..."
        git fetch --prune --no-tags 2>&1 | tee -a "$LOG_FILE"
        
        # 서브모듈도 prune
        git submodule foreach 'git fetch --prune --no-tags' 2>&1 | tee -a "$LOG_FILE"
    fi
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Branch Scope 완료 (${duration}s)${NC}"
    
    return $duration
}

# To-slim 최적화
apply_to_slim() {
    log_color "${YELLOW}>>> To-slim 최적화 (Partial Clone)${NC}"
    
    start_timer
    
    # 메인 저장소
    log_output "메인 저장소 to-slim 적용 중..."
    $GA_CMD opt quick to-slim -q 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 to-slim 적용 중..."
    $GA_CMD opt submodule to-slim -q 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ To-slim 완료 (${duration}s)${NC}"
    
    return $duration
}

# Sparse-checkout 최적화
apply_sparse() {
    log_color "${YELLOW}>>> Sparse-checkout 최적화${NC}"
    
    start_timer
    
    # 기본 sparse 패턴 (예시)
    local patterns="src/ docs/ *.md *.json"
    
    log_output "Sparse 패턴 적용: $patterns"
    $GA_CMD opt sparse add $patterns 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Sparse-checkout 완료 (${duration}s)${NC}"
    
    return $duration
}

# ==================== 결과 분석 함수 ====================

generate_summary() {
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}         최종 결과 요약${NC}"
    log_color "${CYAN}========================================${NC}"
    
    # 요약 파일 생성
    cat > "$SUMMARY_FILE" <<EOF
===========================================
Git 최적화 벤치마크 결과
===========================================
테스트 시각: $(date)
테스트 저장소: $TEST_REPO
테스트 경로: $TEST_DIR
측정 모드: $MODE
===========================================

EOF
    
    # JSONL 파일 분석
    if [ -f "$JSONL_FILE" ]; then
        # Python으로 분석 (jq가 없을 경우를 대비)
        python3 <<EOF >> "$SUMMARY_FILE"
import json
import sys

steps = []
with open('$JSONL_FILE', 'r') as f:
    for line in f:
        steps.append(json.loads(line))

if not steps:
    sys.exit(0)

# 초기 상태
baseline = steps[0]
print(f"[초기 상태]")
print(f"  Git Store: {baseline['git_store'] / 1073741824:.2f} GB")
print(f"  Worktree:  {baseline['worktree'] / 1073741824:.2f} GB")
print(f"  Total:     {baseline['total'] / 1073741824:.2f} GB")
print()

# 단계별 분석
print("[단계별 절감 효과]")
total_saved = 0
max_saving = 0
max_step = ""

for i in range(1, len(steps)):
    prev = steps[i-1]
    curr = steps[i]
    
    delta_git = prev['git_store'] - curr['git_store']
    delta_wt = prev['worktree'] - curr['worktree']
    delta_total = prev['total'] - curr['total']
    
    percent = (delta_total / prev['total'] * 100) if prev['total'] > 0 else 0
    
    print(f"\n{i}. {curr['name']}")
    print(f"   절감: {delta_total / 1073741824:.2f} GB ({percent:.1f}%)")
    
    if curr['duration'] > 0:
        efficiency = delta_total / 1048576 / curr['duration']
        print(f"   시간: {curr['duration']}초")
        print(f"   효율: {efficiency:.1f} MB/s")
    
    total_saved += delta_total
    
    if delta_total > max_saving:
        max_saving = delta_total
        max_step = curr['name']

# 최종 결과
final = steps[-1]
print("\n[최종 결과]")
print(f"  총 절감: {total_saved / 1073741824:.2f} GB")
print(f"  최종 크기: {final['total'] / 1073741824:.2f} GB")
print(f"  절감률: {(total_saved / baseline['total'] * 100):.1f}%")

print("\n[최대 절약 단계]")
print(f"  {max_step}: {max_saving / 1073741824:.2f} GB")

# 각 단계별 비중 계산
print("\n[단계별 비중]")
for i in range(1, len(steps)):
    prev = steps[i-1]
    curr = steps[i]
    delta = prev['total'] - curr['total']
    weight = (delta / total_saved * 100) if total_saved > 0 else 0
    print(f"  {curr['name']}: {weight:.1f}%")
EOF
    fi
    
    # 결과 출력
    log_color "${GREEN}=== 요약 ===${NC}"
    cat "$SUMMARY_FILE" | tee -a "$LOG_FILE"
    
    log_color "${GREEN}✓ 결과 파일 저장됨:${NC}"
    log_output "  - 로그: $LOG_FILE"
    log_output "  - 데이터: $JSONL_FILE"
    log_output "  - 요약: $SUMMARY_FILE"
}

# ==================== 메인 실행 ====================

main() {
    # 시작
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}   Git 최적화 벤치마크 시작${NC}"
    log_color "${CYAN}========================================${NC}"
    log_output "모드: $MODE"
    log_output "Skip Slim: $SKIP_SLIM"
    log_output "Skip Sparse: $SKIP_SPARSE"
    log_output ""
    
    # 전역 변수 초기화
    TOTAL_SAVED=0
    LAST_GIT_STORE=0
    LAST_WORKTREE=0
    LAST_TOTAL=0
    
    local step_count=0
    local total_steps=5
    if [ "$SKIP_SLIM" = true ]; then
        total_steps=$((total_steps - 1))
    fi
    if [ "$SKIP_SPARSE" = true ]; then
        total_steps=$((total_steps - 1))
    fi
    
    # Step 1: 저장소 초기화
    show_progress 1 $total_steps "저장소 초기화"
    init_repository
    init_duration=$?
    
    # Step 2: Baseline 측정
    show_progress 2 $total_steps "Baseline 측정"
    measure_step "BASELINE" 0 0
    local baseline_git=$LAST_GIT_STORE
    local baseline_wt=$LAST_WORKTREE
    local baseline_total=$LAST_TOTAL
    
    # Step 3: Branch Scope 최적화 (먼저 실행)
    show_progress 3 $total_steps "Branch Scope 최적화"
    apply_branch_scope
    local scope_duration=$?
    measure_step "BRANCH_SCOPE" 1 $scope_duration
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL $scope_duration
    local after_scope_git=$LAST_GIT_STORE
    local after_scope_wt=$LAST_WORKTREE
    local after_scope_total=$LAST_TOTAL
    
    # Step 4: Shallow 최적화 (Branch Scope 이후)
    show_progress 4 $total_steps "Shallow 최적화"
    apply_shallow 1
    local shallow_duration=$?
    measure_step "SHALLOW" 2 $shallow_duration
    calculate_delta $after_scope_git $LAST_GIT_STORE $after_scope_wt $LAST_WORKTREE $after_scope_total $LAST_TOTAL $shallow_duration
    local after_shallow_git=$LAST_GIT_STORE
    local after_shallow_wt=$LAST_WORKTREE
    local after_shallow_total=$LAST_TOTAL
    
    # Step 5: To-slim 최적화 (선택적)
    if [ "$SKIP_SLIM" = false ]; then
        show_progress 5 $total_steps "To-slim 최적화"
        apply_to_slim
        local slim_duration=$?
        measure_step "TO_SLIM" 3 $slim_duration
        calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL $slim_duration
    fi
    
    # Step 6: Sparse-checkout (선택적)
    if [ "$SKIP_SPARSE" = false ]; then
        local step_num=$((total_steps))
        show_progress $step_num $total_steps "Sparse-checkout 최적화"
        apply_sparse
        local sparse_duration=$?
        measure_step "SPARSE" 4 $sparse_duration
        calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL $sparse_duration
    fi
    
    # 최종 요약
    generate_summary
    
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}   벤치마크 완료!${NC}"
    log_color "${CYAN}========================================${NC}"
}

# 스크립트 실행
main "$@"