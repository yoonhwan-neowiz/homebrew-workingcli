#!/bin/bash

# Git Optimization Benchmark - Scenario: Branch Switch Workflow
# 작업 브랜치 전환 워크플로우 - 빈번한 태스크 브랜치 완료/추가/변경 시뮬레이션

# 시나리오 설명
SCENARIO_NAME="Branch Switch Workflow"
SCENARIO_DESC="작업 브랜치 전환 - 브랜치 완료 → Master Reset → 브랜치 제거/추가 → 재활성화"

# 초기 브랜치 (02 시나리오 완료 상태)
INITIAL_BRANCHES="master origin/live59.b/5904.7 origin/live59.b/5621.20"

# 제거할 브랜치
REMOVE_BRANCHES="origin/live59.b/5904.7"

# 추가할 브랜치
ADD_BRANCHES="origin/live59.b/51005.29 origin/live59.b/5914.13"

# 유지할 브랜치
KEEP_BRANCHES="origin/live59.b/5621.20"

# 최종 브랜치 목록
FINAL_BRANCHES="master origin/live59.b/5621.20 origin/live59.b/51005.29 origin/live59.b/5914.13"

# 시나리오 실행 함수
run_scenario() {
    log_color "${CYAN}========================================"
    log_color "${CYAN}   시나리오: $SCENARIO_NAME"
    log_color "${CYAN}========================================"
    log_output "$SCENARIO_DESC"
    log_output ""
    
    # 전역 변수 초기화
    TOTAL_SAVED=0
    LAST_GIT_STORE=0
    LAST_WORKTREE=0
    LAST_TOTAL=0
    
    # 단계 계산
    local total_steps=7
    
    # Step 1: 기존 작업 디렉토리 사용 (02 시나리오 완료 상태)
    show_progress 1 $total_steps "기존 작업 디렉토리 사용"
    
    start_timer
    
    # 테스트 디렉토리 확인
    if [ ! -d "$TEST_DIR" ]; then
        log_color "${RED}❌ 작업 디렉토리가 없습니다: $TEST_DIR"
        log_color "${YELLOW}먼저 02-setup-master-workflow 시나리오를 실행하세요"
        exit 1
    fi
    
    cd "$TEST_DIR"
    
    # 현재 상태 확인
    log_output "기존 작업 디렉토리 사용: $TEST_DIR"
    log_output "현재 브랜치 상태:"
    git branch -a | grep -E "master|live59" | tee -a "$LOG_FILE"
    
    # master 체크아웃
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    local init_duration=$(end_timer)
    log_color "${GREEN}✓ 작업 디렉토리 준비 완료 (${init_duration}s)"
    
    # Step 2: Baseline 측정
    show_progress 2 $total_steps "Baseline 측정 (작업 완료 상태)"
    measure_step "BASELINE" 0 0
    local baseline_git=$LAST_GIT_STORE
    local baseline_wt=$LAST_WORKTREE
    local baseline_total=$LAST_TOTAL
    
    # Step 3: 작업 브랜치 시뮬레이션 (커밋 추가)
    show_progress 3 $total_steps "작업 브랜치 개발 시뮬레이션"
    
    start_timer
    
    # live59.b/5904.7 브랜치에서 작업 완료 시뮬레이션
    if git branch -r | grep -q "origin/live59.b/5904.7"; then
        git checkout -b live59.b/5904.7 origin/live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
        echo "작업 완료: Feature 5904.7" > feature_5904.txt
        git add feature_5904.txt
        git commit -m "feat: Complete feature 5904.7" 2>&1 | tee -a "$LOG_FILE"
    fi
    
    # live59.b/5621.20 브랜치에서 작업 진행 시뮬레이션
    if git branch -r | grep -q "origin/live59.b/5621.20"; then
        git checkout -b live59.b/5621.20 origin/live59.b/5621.20 2>&1 | tee -a "$LOG_FILE"
        echo "작업 진행: Feature 5621.20" > feature_5621.txt
        git add feature_5621.txt
        git commit -m "wip: Ongoing work on feature 5621.20" 2>&1 | tee -a "$LOG_FILE"
    fi
    
    local work_duration=$(end_timer)
    log_color "${GREEN}✓ 작업 시뮬레이션 완료 (${work_duration}s)"
    
    # Step 4: Master Reset (Shallow 1)
    show_progress 4 $total_steps "Master Reset to Shallow 1"
    
    start_timer
    
    # master로 체크아웃
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    # 모든 로컬 브랜치 삭제 (master 제외)
    git branch | grep -v master | xargs -r git branch -D 2>&1 | tee -a "$LOG_FILE" || true
    
    # shallow 1로 리셋
    log_output "Master를 shallow depth=1로 리셋..."
    $GA_CMD opt quick to-shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈도 동일하게 처리
    if [ -f ".gitmodules" ]; then
        log_output "서브모듈 shallow 리셋..."
        $GA_CMD opt submodule to-shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
    fi
    
    local reset_duration=$(end_timer)
    measure_step "MASTER_RESET" 1 $reset_duration
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL $reset_duration
    local after_reset_git=$LAST_GIT_STORE
    local after_reset_wt=$LAST_WORKTREE
    local after_reset_total=$LAST_TOTAL
    
    # Step 5: 브랜치 제거 (5904.7)
    show_progress 5 $total_steps "완료된 브랜치 제거"
    
    start_timer
    
    log_output "완료된 브랜치 제거: $REMOVE_BRANCHES"
    
    # 새로운 브랜치 scope 설정 (5904.7 제외)
    local remaining_branches="master origin/live59.b/5621.20"
    $GA_CMD opt quick set-branch-scope $remaining_branches 2>&1 | tee -a "$LOG_FILE"
    
    # 불필요한 참조 정리
    git remote prune origin 2>&1 | tee -a "$LOG_FILE" || true
    
    local remove_duration=$(end_timer)
    measure_step "REMOVE_BRANCH" 2 $remove_duration
    calculate_delta $after_reset_git $LAST_GIT_STORE $after_reset_wt $LAST_WORKTREE $after_reset_total $LAST_TOTAL $remove_duration
    local after_remove_git=$LAST_GIT_STORE
    local after_remove_wt=$LAST_WORKTREE
    local after_remove_total=$LAST_TOTAL
    
    # Step 6: 새 브랜치 추가
    show_progress 6 $total_steps "새 작업 브랜치 추가"
    
    start_timer
    
    log_output "새 브랜치 추가: $ADD_BRANCHES"
    log_output "유지 브랜치: $KEEP_BRANCHES"
    
    # 최종 브랜치 scope 설정
    $GA_CMD opt quick set-branch-scope $FINAL_BRANCHES 2>&1 | tee -a "$LOG_FILE"
    
    # 새 브랜치들 fetch (shallow)
    for branch in $ADD_BRANCHES; do
        local branch_name=${branch#origin/}
        log_output "Fetching $branch_name..."
        git fetch origin $branch_name --depth=1 2>&1 | tee -a "$LOG_FILE" || {
            log_color "${YELLOW}⚠ Shallow fetch 실패, 전체 fetch 시도..."
            git fetch origin $branch_name 2>&1 | tee -a "$LOG_FILE"
        }
    done
    
    local add_duration=$(end_timer)
    measure_step "ADD_BRANCHES" 3 $add_duration
    calculate_delta $after_remove_git $LAST_GIT_STORE $after_remove_wt $LAST_WORKTREE $after_remove_total $LAST_TOTAL $add_duration
    local after_add_git=$LAST_GIT_STORE
    local after_add_wt=$LAST_WORKTREE
    local after_add_total=$LAST_TOTAL
    
    # Step 7: 기존 브랜치 재활성화
    show_progress 7 $total_steps "기존 브랜치 재활성화"
    
    start_timer
    
    log_output "브랜치 재활성화: live59.b/5621.20"
    
    # 5621.20 브랜치 체크아웃 (이미 있으므로 업데이트만)
    if git branch -r | grep -q "origin/live59.b/5621.20"; then
        git checkout -b live59.b/5621.20 origin/live59.b/5621.20 2>/dev/null || {
            git checkout live59.b/5621.20 2>&1 | tee -a "$LOG_FILE"
            git pull origin live59.b/5621.20 --depth=1 2>&1 | tee -a "$LOG_FILE" || {
                log_color "${YELLOW}⚠ Shallow pull 실패, 일반 pull 시도..."
                git pull origin live59.b/5621.20 2>&1 | tee -a "$LOG_FILE"
            }
        }
    fi
    
    # master로 돌아가기
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    local reactivate_duration=$(end_timer)
    log_color "${GREEN}✓ 브랜치 재활성화 완료 (${reactivate_duration}s)"
    
    # 최종 측정
    measure_step "FINAL" 7 $reactivate_duration
    calculate_delta $after_add_git $LAST_GIT_STORE $after_add_wt $LAST_WORKTREE $after_add_total $LAST_TOTAL 0
    
    # 최종 브랜치 상태 확인
    log_output ""
    log_output "최종 브랜치 상태:"
    git branch -a | grep -E "master|live59" | tee -a "$LOG_FILE"
    
    # 최종 결과 출력
    log_color "${CYAN}========================================"
    log_color "${GREEN}시나리오 완료: $SCENARIO_NAME"
    log_output "총 절감: $(bytes_to_human $TOTAL_SAVED)"
    log_output "전체 소요 시간: $(( init_duration + work_duration + reset_duration + remove_duration + add_duration + reactivate_duration ))s"
    log_color "${CYAN}========================================"
}

# 시나리오별 설정 오버라이드
scenario_setup() {
    # 브랜치 목록을 최종 상태로 설정
    BRANCHES="$FINAL_BRANCHES"
    
    # 기타 설정
    SKIP_SLIM=true  # to-slim은 이 워크플로우에서 제외
    SKIP_SPARSE=true  # sparse-checkout도 제외
}

# 시나리오 정리
scenario_cleanup() {
    # 필요시 정리 작업
    :
}