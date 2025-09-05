#!/bin/bash

# Git Optimization Benchmark - Scenario: Full Reset to Master Shallow 1
# 전체 리셋 워크플로우 - 모든 브랜치 제거 후 Master Shallow 1만 유지

# 시나리오 설명
SCENARIO_NAME="Full Reset to Master Shallow 1"
SCENARIO_DESC="전체 리셋 - 모든 브랜치 삭제 → Master Shallow 1 → 서브모듈 동기화 → Clean 상태"

# 초기 브랜치 (복잡한 상태 시뮬레이션)
INITIAL_BRANCHES="master origin/live59.b/5904.7 origin/live59.b/5621.20 origin/live59.b/51005.29 origin/live59.b/5914.13 origin/live59.b/5822.15"

# 최종 브랜치 (master만)
FINAL_BRANCHES="master"

# Clean 모드 (safe 또는 full)
CLEAN_MODE="${CLEAN_MODE:-safe}"

# 시나리오 실행 함수
run_scenario() {
    log_color "${CYAN}========================================"
    log_color "${CYAN}   시나리오: $SCENARIO_NAME"
    log_color "${CYAN}========================================"
    log_output "$SCENARIO_DESC"
    log_output "Clean Mode: $CLEAN_MODE"
    log_output ""
    
    # 전역 변수 초기화
    TOTAL_SAVED=0
    LAST_GIT_STORE=0
    LAST_WORKTREE=0
    LAST_TOTAL=0
    
    # 단계 계산
    local total_steps=7
    
    # Step 1: 기존 작업 디렉토리 사용
    show_progress 1 $total_steps "작업 디렉토리 확인"
    
    start_timer
    
    # 테스트 디렉토리 확인
    if [ ! -d "$TEST_DIR" ]; then
        log_color "${RED}❌ 작업 디렉토리가 없습니다: $TEST_DIR"
        log_color "${YELLOW}먼저 02 또는 03 시나리오를 실행하세요"
        exit 1
    fi
    
    cd "$TEST_DIR"
    
    # 현재 상태 확인
    log_output "기존 작업 디렉토리 사용: $TEST_DIR"
    log_output "현재 브랜치 상태:"
    git branch -a | tee -a "$LOG_FILE"
    
    # master 체크아웃
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    local init_duration=$(end_timer)
    log_color "${GREEN}✓ 작업 디렉토리 준비 완료 (${init_duration}s)"
    
    # Step 2: Baseline 측정 (복잡한 상태)
    show_progress 2 $total_steps "Baseline 측정 (복잡한 상태)"
    measure_step "BASELINE_COMPLEX" 0 0
    local baseline_git=$LAST_GIT_STORE
    local baseline_wt=$LAST_WORKTREE
    local baseline_total=$LAST_TOTAL
    
    log_output "복잡한 상태 저장소 크기:"
    log_output "  Git Store: $(bytes_to_human $baseline_git)"
    log_output "  Worktree: $(bytes_to_human $baseline_wt)"
    log_output "  Total: $(bytes_to_human $baseline_total)"
    
    # Step 3: 모든 로컬 브랜치 삭제
    show_progress 3 $total_steps "모든 로컬 브랜치 삭제"
    
    start_timer
    
    # master 체크아웃
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    # 모든 로컬 브랜치 삭제 (master 제외)
    log_output "모든 로컬 브랜치 삭제 중..."
    local deleted_count=0
    for branch in $(git branch | grep -v master | sed 's/*//g'); do
        git branch -D $branch 2>&1 | tee -a "$LOG_FILE"
        ((deleted_count++))
    done
    log_output "$deleted_count개 브랜치 삭제됨"
    
    local delete_duration=$(end_timer)
    measure_step "DELETE_BRANCHES" 1 $delete_duration
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL $delete_duration
    local after_delete_git=$LAST_GIT_STORE
    local after_delete_wt=$LAST_WORKTREE
    local after_delete_total=$LAST_TOTAL
    
    # Step 4: Master를 Shallow Depth=1로 리셋
    show_progress 4 $total_steps "Master를 Shallow Depth=1로 리셋"
    
    start_timer
    
    log_output "Master를 shallow depth=1로 변환..."
    
    # Branch scope을 master만으로 설정
    $GA_CMD opt quick set-branch-scope master 2>&1 | tee -a "$LOG_FILE"
    
    # Shallow 변환
    $GA_CMD opt quick to-shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
    
    local shallow_duration=$(end_timer)
    measure_step "SHALLOW_RESET" 2 $shallow_duration
    calculate_delta $after_delete_git $LAST_GIT_STORE $after_delete_wt $LAST_WORKTREE $after_delete_total $LAST_TOTAL $shallow_duration
    local after_shallow_git=$LAST_GIT_STORE
    local after_shallow_wt=$LAST_WORKTREE
    local after_shallow_total=$LAST_TOTAL
    
    # Step 5: 서브모듈 동기화 (Master Shallow 1)
    show_progress 5 $total_steps "서브모듈 Master Shallow 1 동기화"
    
    start_timer
    
    if [ -f ".gitmodules" ]; then
        log_output "서브모듈 초기화 및 shallow 설정..."
        
        # 서브모듈 초기화
        git submodule update --init 2>&1 | tee -a "$LOG_FILE"
        
        # 각 서브모듈을 master shallow 1로 설정
        git submodule foreach --recursive '
            git checkout master 2>/dev/null || git checkout -b master 2>/dev/null || true
            git branch | grep -v master | xargs -r git branch -D 2>/dev/null || true
        ' 2>&1 | tee -a "$LOG_FILE"
        
        # 서브모듈 shallow 변환
        $GA_CMD opt submodule to-shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
        
        # 서브모듈 branch-scope 설정
        $GA_CMD opt submodule set-branch-scope master -q 2>&1 | tee -a "$LOG_FILE"
    else
        log_output "서브모듈이 없습니다."
    fi
    
    local submodule_duration=$(end_timer)
    measure_step "SUBMODULE_SYNC" 3 $submodule_duration
    calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL $submodule_duration
    local after_submodule_git=$LAST_GIT_STORE
    local after_submodule_wt=$LAST_WORKTREE
    local after_submodule_total=$LAST_TOTAL
    
    # Step 6: Clean 작업 (옵션에 따라)
    show_progress 6 $total_steps "Clean 작업 ($CLEAN_MODE mode)"
    
    start_timer
    
    if [ "$CLEAN_MODE" = "full" ]; then
        log_output "Full clean mode - 모든 불필요한 객체 제거..."
        
        # 깊은 정리
        git reflog expire --expire=now --all 2>&1 | tee -a "$LOG_FILE"
        git gc --prune=now --aggressive 2>&1 | tee -a "$LOG_FILE"
        git repack -Ad 2>&1 | tee -a "$LOG_FILE"
        
        # 원격 참조 정리
        git remote prune origin 2>&1 | tee -a "$LOG_FILE"
        
        # 불필요한 파일 제거
        git clean -fdx 2>&1 | tee -a "$LOG_FILE"
        
    elif [ "$CLEAN_MODE" = "safe" ]; then
        log_output "Safe clean mode - 기본 정리만 수행..."
        
        # 안전한 정리
        git gc --auto 2>&1 | tee -a "$LOG_FILE"
        git remote prune origin 2>&1 | tee -a "$LOG_FILE"
        git clean -fd 2>&1 | tee -a "$LOG_FILE"
    fi
    
    local clean_duration=$(end_timer)
    measure_step "CLEAN_WORK" 4 $clean_duration
    calculate_delta $after_submodule_git $LAST_GIT_STORE $after_submodule_wt $LAST_WORKTREE $after_submodule_total $LAST_TOTAL $clean_duration
    local after_clean_git=$LAST_GIT_STORE
    local after_clean_wt=$LAST_WORKTREE
    local after_clean_total=$LAST_TOTAL
    
    # Step 7: 최종 검증
    show_progress 7 $total_steps "최종 상태 검증"
    
    start_timer
    
    log_output "최종 브랜치 상태:"
    git branch -a | tee -a "$LOG_FILE"
    
    log_output ""
    log_output "최종 커밋 히스토리 (depth=1 확인):"
    git log --oneline -5 2>&1 | tee -a "$LOG_FILE"
    
    log_output ""
    log_output "Remote 상태:"
    git remote -v | tee -a "$LOG_FILE"
    
    if [ -f ".gitmodules" ]; then
        log_output ""
        log_output "서브모듈 상태:"
        git submodule status | tee -a "$LOG_FILE"
    fi
    
    local verify_duration=$(end_timer)
    log_color "${GREEN}✓ 검증 완료 (${verify_duration}s)"
    
    # 최종 측정
    measure_step "FINAL" 7 0
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL 0
    
    # 최종 결과 출력
    log_color "${CYAN}========================================"
    log_color "${GREEN}시나리오 완료: $SCENARIO_NAME"
    log_output ""
    log_output "저장소 크기 변화:"
    log_output "  초기 (복잡한 상태): $(bytes_to_human $baseline_total)"
    log_output "  최종 (Clean 상태): $(bytes_to_human $LAST_TOTAL)"
    log_output "  총 절감: $(bytes_to_human $TOTAL_SAVED)"
    log_output "  절감률: $(( TOTAL_SAVED * 100 / baseline_total ))%"
    log_output ""
    log_output "전체 소요 시간: $(( init_duration + delete_duration + shallow_duration + submodule_duration + clean_duration + verify_duration ))s"
    log_output "Clean Mode: $CLEAN_MODE"
    log_color "${CYAN}========================================"
}

# 시나리오별 설정 오버라이드
scenario_setup() {
    # 최종 브랜치는 master만
    BRANCHES="master"
    
    # 기타 설정
    SKIP_SLIM=true  # to-slim은 이 워크플로우에서 제외
    SKIP_SPARSE=true  # sparse-checkout도 제외
    
    # Clean 모드 옵션 처리
    if [ -n "$SCENARIO_CLEAN_MODE" ]; then
        CLEAN_MODE="$SCENARIO_CLEAN_MODE"
    fi
}

# 시나리오 정리
scenario_cleanup() {
    # 필요시 정리 작업
    :
}