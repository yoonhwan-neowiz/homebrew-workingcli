#!/bin/bash

# Git Optimization Benchmark - Scenario: Full Optimization
# 전체 최적화 시나리오 - 모든 최적화 단계를 순차적으로 실행

# 시나리오 설명
SCENARIO_NAME="Full Optimization"
SCENARIO_DESC="전체 최적화 프로세스 - Branch Scope → Shallow → To-slim → Sparse"

# 시나리오 실행 함수
run_scenario() {
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}   시나리오: $SCENARIO_NAME${NC}"
    log_color "${CYAN}========================================${NC}"
    log_output "$SCENARIO_DESC"
    log_output ""
    
    # 전역 변수 초기화
    TOTAL_SAVED=0
    LAST_GIT_STORE=0
    LAST_WORKTREE=0
    LAST_TOTAL=0
    
    # 단계 수 계산
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
    if [ "$USE_CLONE" = true ]; then
        init_repository_clone
    else
        init_repository
    fi
    # INIT_DURATION 전역 변수 사용 (optimize.sh에서 설정됨)
    init_duration=$INIT_DURATION
    
    # Step 2: Baseline 측정
    show_progress 2 $total_steps "Baseline 측정"
    measure_step "BASELINE" 0 0
    local baseline_git=$LAST_GIT_STORE
    local baseline_wt=$LAST_WORKTREE
    local baseline_total=$LAST_TOTAL
    
    # Step 3: Branch Scope 최적화 (먼저 실행)
    show_progress 3 $total_steps "Branch Scope 최적화"
    apply_branch_scope
    local scope_duration=$OPT_DURATION
    measure_step "BRANCH_SCOPE" 1 $scope_duration
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL $scope_duration
    local after_scope_git=$LAST_GIT_STORE
    local after_scope_wt=$LAST_WORKTREE
    local after_scope_total=$LAST_TOTAL
    
    # Step 4: Shallow 최적화 (Branch Scope 이후)
    show_progress 4 $total_steps "Shallow 최적화"
    apply_shallow 1
    local shallow_duration=$OPT_DURATION
    measure_step "SHALLOW" 2 $shallow_duration
    calculate_delta $after_scope_git $LAST_GIT_STORE $after_scope_wt $LAST_WORKTREE $after_scope_total $LAST_TOTAL $shallow_duration
    local after_shallow_git=$LAST_GIT_STORE
    local after_shallow_wt=$LAST_WORKTREE
    local after_shallow_total=$LAST_TOTAL
    
    # Step 5: To-slim 최적화 (선택적)
    if [ "$SKIP_SLIM" = false ]; then
        show_progress 5 $total_steps "To-slim 최적화"
        apply_to_slim
        local slim_duration=$OPT_DURATION
        measure_step "TO_SLIM" 3 $slim_duration
        calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL $slim_duration
        local after_slim_git=$LAST_GIT_STORE
        local after_slim_wt=$LAST_WORKTREE
        local after_slim_total=$LAST_TOTAL
    fi
    
    # Step 6: Sparse-checkout (선택적)
    if [ "$SKIP_SPARSE" = false ]; then
        local step_num=$((total_steps))
        show_progress $step_num $total_steps "Sparse-checkout 최적화"
        apply_sparse
        local sparse_duration=$OPT_DURATION
        measure_step "SPARSE" 4 $sparse_duration
        
        # 델타 계산시 이전 단계 값 사용
        if [ "$SKIP_SLIM" = false ]; then
            calculate_delta $after_slim_git $LAST_GIT_STORE $after_slim_wt $LAST_WORKTREE $after_slim_total $LAST_TOTAL $sparse_duration
        else
            calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL $sparse_duration
        fi
    fi
    
    # 최종 결과 출력
    log_color "${CYAN}========================================${NC}"
    log_color "${GREEN}시나리오 완료: $SCENARIO_NAME${NC}"
    log_output "총 절감: $(bytes_to_human $TOTAL_SAVED)"
    log_color "${CYAN}========================================${NC}"
}

# 시나리오별 설정 오버라이드 (선택적)
scenario_setup() {
    # 이 시나리오에 특화된 설정이 있다면 여기서 설정
    # 예: BRANCHES="main develop"
    :
}

# 시나리오 정리
scenario_cleanup() {
    # 시나리오 종료 후 정리 작업이 필요하면 여기서 수행
    :
}