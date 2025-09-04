#!/bin/bash

# Git Optimization Benchmark - Scenario: Setup Master Workflow
# Master 클론 셋업 워크플로우 - 새로운 프로젝트 시작 시뮬레이션

# 시나리오 설명
SCENARIO_NAME="Setup Master Workflow"
SCENARIO_DESC="Master 클론 셋업 - Clone Master → Branch Scope → Shallow → Checkout → Submodule → Auto Merge-Base"

# 시나리오별 브랜치 설정
WORKFLOW_BRANCHES="master origin/live59.b/5904.7 origin/live59.b/5621.20"

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
    
    # 단계 계산
    local total_steps=7
    
    # Step 1: 저장소 초기화 (Clone Master 방식)
    show_progress 1 $total_steps "저장소 초기화 (Clone Master)"
    
    # 기존 테스트 디렉토리 삭제
    if [ -d "$TEST_DIR" ]; then
        log_output "기존 테스트 디렉토리 삭제 중..."
        rm -rf "$TEST_DIR"
    fi
    
    start_timer
    
    # Master 브랜치만 클론 (--single-branch 옵션)
    log_output "Master 브랜치 클론 중..."
    log_color "${CYAN}Clone master 진행 상황이 표시됩니다...${NC}"
    
    # 실제 원격 저장소가 있는 경우 사용
    if [ -n "$TEST_REPO" ] && [ "$USE_CLONE" = true ]; then
        GIT_TRACE=1 GIT_CURL_VERBOSE=1 git clone \
            --single-branch \
            --branch master \
            --progress \
            --verbose \
            "$TEST_REPO" "$TEST_DIR" 2>&1 | tee -a "$LOG_FILE"
    else
        # 로컬 템플릿에서 시뮬레이션
        local SOURCE_DIR="${SOURCE_TEMPLATE:-$HOME/Work/DesignB4}"
        
        if [ ! -d "$SOURCE_DIR" ]; then
            log_color "${RED}❌ 템플릿 저장소가 없습니다: $SOURCE_DIR${NC}"
            exit 1
        fi
        
        # 템플릿 복사
        if command -v rusync &> /dev/null; then
            rusync "$SOURCE_DIR" "$TEST_DIR" 2>&1 | tee -a "$LOG_FILE"
        else
            rsync -av --progress "$SOURCE_DIR/" "$TEST_DIR/" 2>&1 | tee -a "$LOG_FILE"
        fi
        
        cd "$TEST_DIR"
        
        # master 브랜치만 남기고 다른 브랜치 제거 시뮬레이션
        git checkout master 2>&1 | tee -a "$LOG_FILE"
        git branch | grep -v master | xargs -r git branch -D 2>&1 | tee -a "$LOG_FILE" || true
    fi
    
    cd "$TEST_DIR"
    
    # git 설정
    git config pull.rebase false
    
    local init_duration=$(end_timer)
    log_color "${GREEN}✓ Clone Master 완료 (${init_duration}s)${NC}"
    
    # Step 2: Baseline 측정
    show_progress 2 $total_steps "Baseline 측정"
    measure_step "BASELINE" 0 0
    local baseline_git=$LAST_GIT_STORE
    local baseline_wt=$LAST_WORKTREE
    local baseline_total=$LAST_TOTAL
    
    # Step 3: Branch Scope 설정
    show_progress 3 $total_steps "Branch Scope 설정"
    log_output "브랜치 범위 설정: $WORKFLOW_BRANCHES"
    
    start_timer
    
    # 원격 브랜치 fetch (master 외 추가 브랜치)
    log_output "추가 브랜치 fetch 중..."
    git fetch origin live59.b/5904.7:refs/remotes/origin/live59.b/5904.7 2>&1 | tee -a "$LOG_FILE" || true
    git fetch origin live59.b/5621.20:refs/remotes/origin/live59.b/5621.20 2>&1 | tee -a "$LOG_FILE" || true
    
    # 메인 저장소 branch-scope
    log_output "메인 저장소 branch-scope 적용 중..."
    $GA_CMD opt quick set-branch-scope $WORKFLOW_BRANCHES -q 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈 branch-scope
    log_output "서브모듈 branch-scope 적용 중..."
    $GA_CMD opt submodule set-branch-scope $WORKFLOW_BRANCHES -q 2>&1 | tee -a "$LOG_FILE"
    
    local scope_duration=$(end_timer)
    measure_step "BRANCH_SCOPE" 1 $scope_duration
    calculate_delta $baseline_git $LAST_GIT_STORE $baseline_wt $LAST_WORKTREE $baseline_total $LAST_TOTAL $scope_duration
    local after_scope_git=$LAST_GIT_STORE
    local after_scope_wt=$LAST_WORKTREE
    local after_scope_total=$LAST_TOTAL
    
    # Step 4: Shallow depth 1 적용
    show_progress 4 $total_steps "Shallow 최적화 (depth=1)"
    apply_shallow 1
    local shallow_duration=$OPT_DURATION
    measure_step "SHALLOW" 2 $shallow_duration
    calculate_delta $after_scope_git $LAST_GIT_STORE $after_scope_wt $LAST_WORKTREE $after_scope_total $LAST_TOTAL $shallow_duration
    local after_shallow_git=$LAST_GIT_STORE
    local after_shallow_wt=$LAST_WORKTREE
    local after_shallow_total=$LAST_TOTAL
    
    # Step 5: Local Checkout
    show_progress 5 $total_steps "Local Checkout"
    
    start_timer
    
    # 각 브랜치 checkout
    log_output "브랜치 체크아웃 중..."
    
    # master는 이미 체크아웃 되어 있음
    log_output "master 브랜치는 이미 체크아웃됨"
    
    # 추가 브랜치들 로컬로 체크아웃
    for branch in origin/live59.b/5904.7 origin/live59.b/5621.20; do
        local local_branch="${branch#origin/}"
        log_output "Checkout 브랜치: $local_branch"
        
        git checkout -b "$local_branch" "$branch" 2>&1 | tee -a "$LOG_FILE" || {
            log_color "${YELLOW}⚠ $local_branch 체크아웃 실패 - 이미 존재할 수 있음${NC}"
            git checkout "$local_branch" 2>&1 | tee -a "$LOG_FILE" || true
        }
    done
    
    # master로 돌아가기
    git checkout master 2>&1 | tee -a "$LOG_FILE"
    
    local checkout_duration=$(end_timer)
    log_color "${GREEN}✓ Checkout 완료 (${checkout_duration}s)${NC}"
    
    # Step 6: Submodule 동일 작업
    show_progress 6 $total_steps "Submodule 처리"
    
    start_timer
    
    # 서브모듈 초기화 및 최적화
    log_output "서브모듈 초기화 중..."
    
    # 서브모듈 초기화 (처음 클론한 경우)
    git submodule init 2>&1 | tee -a "$LOG_FILE" || true
    
    # ga opt submodule update 사용
    $GA_CMD opt submodule update -f 2>&1 | tee -a "$LOG_FILE" || {
        log_color "${YELLOW}⚠ 서브모듈 초기화 실패 - 계속 진행${NC}"
    }
    
    # 서브모듈에도 shallow 적용
    log_output "서브모듈 shallow 적용 중..."
    $GA_CMD opt submodule shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈에도 branch-scope 적용
    log_output "서브모듈 branch-scope 적용 중..."
    $GA_CMD opt submodule set-branch-scope $WORKFLOW_BRANCHES -q 2>&1 | tee -a "$LOG_FILE"
    
    local submodule_duration=$(end_timer)
    log_color "${GREEN}✓ Submodule 처리 완료 (${submodule_duration}s)${NC}"
    
    # Step 7: Auto Find Merge-Base
    show_progress 7 $total_steps "Auto Find Merge-Base"
    
    start_timer
    
    # auto_find_merge_base 실행
    log_output "Merge-base 자동 탐색 중..."
    log_output "대상 브랜치: master live59.b/5904.7 live59.b/5621.20"
    
    # ga opt quick auto 명령 실행
    $GA_CMD opt quick auto master origin/live59.b/5904.7 origin/live59.b/5621.20 -q 2>&1 | tee -a "$LOG_FILE" || {
        log_color "${YELLOW}⚠ Merge-base 탐색 실패 또는 시간 초과${NC}"
        
        # 실패 시 다시 로컬 브랜치로 시도
        log_output "로컬 브랜치로 재시도..."
        $GA_CMD opt quick auto master live59.b/5904.7 live59.b/5621.20 -q 2>&1 | tee -a "$LOG_FILE" || {
            log_color "${YELLOW}⚠ 두 번째 시도도 실패${NC}"
        }
    }
    
    local auto_duration=$(end_timer)
    log_color "${GREEN}✓ Auto Find Merge-Base 완료 (${auto_duration}s)${NC}"
    
    # 최종 측정
    measure_step "FINAL" 7 $auto_duration
    calculate_delta $after_shallow_git $LAST_GIT_STORE $after_shallow_wt $LAST_WORKTREE $after_shallow_total $LAST_TOTAL 0
    
    # 최종 결과 출력
    log_color "${CYAN}========================================${NC}"
    log_color "${GREEN}시나리오 완료: $SCENARIO_NAME${NC}"
    log_output "총 절감: $(bytes_to_human $TOTAL_SAVED)"
    log_output "전체 소요 시간: $(( init_duration + scope_duration + shallow_duration + checkout_duration + submodule_duration + auto_duration ))s"
    log_color "${CYAN}========================================${NC}"
}

# 시나리오별 설정 오버라이드
scenario_setup() {
    # 브랜치 목록을 워크플로우용으로 설정
    BRANCHES="$WORKFLOW_BRANCHES"
    
    # 기타 설정
    SKIP_SLIM=true  # to-slim은 이 워크플로우에서 제외
    SKIP_SPARSE=true  # sparse-checkout도 제외
}

# 시나리오 정리
scenario_cleanup() {
    # 필요시 정리 작업
    :
}