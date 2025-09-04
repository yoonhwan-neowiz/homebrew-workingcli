#!/bin/bash

# Git Optimization Benchmark - Optimization Functions
# 저장소 초기화, 최적화 단계 실행

# ==================== 저장소 초기화 ====================

# 저장소 초기화 (rusync/rsync 복사 방식)
init_repository() {
    log_color "${YELLOW}>>> 저장소 초기화 (고속 복사 방식)${NC}"
    
    # 소스 템플릿 저장소
    local SOURCE_DIR="${SOURCE_TEMPLATE:-$HOME/Work/DesignB4}"
    
    # 템플릿 저장소 확인
    if [ ! -d "$SOURCE_DIR" ]; then
        log_color "${RED}❌ 템플릿 저장소가 없습니다: $SOURCE_DIR${NC}"
        log_output "DesignB4를 먼저 준비해주세요 (submodule 포함)"
        log_output "또는 --use-clone 옵션을 사용하세요"
        exit 1
    fi
    
    # 기존 테스트 디렉토리 삭제
    if [ -d "$TEST_DIR" ]; then
        log_output "기존 테스트 디렉토리 삭제 중..."
        rm -rf "$TEST_DIR"
    fi
    
    # rusync 또는 rsync로 고속 복사
    start_timer
    
    if command -v rusync &> /dev/null; then
        log_output "템플릿 저장소 고속 복사 중 (rusync): $SOURCE_DIR → $TEST_DIR"
        log_color "${CYAN}rusync를 사용한 병렬 복사 진행 중...${NC}"
        rusync "$SOURCE_DIR" "$TEST_DIR" 2>&1 | tee -a "$LOG_FILE"
        local copy_tool="rusync (고속)"
    else
        log_color "${YELLOW}rusync를 찾을 수 없어 rsync 사용${NC}"
        log_output "템플릿 저장소 복사 중 (rsync): $SOURCE_DIR → $TEST_DIR"
        rsync -av --progress "$SOURCE_DIR/" "$TEST_DIR/" 2>&1 | tee -a "$LOG_FILE"
        local copy_tool="rsync"
    fi
    
    cd "$TEST_DIR"
    
    # git 설정 초기화
    git config pull.rebase false
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ 저장소 초기화 완료 (${duration}s)${NC}"
    log_output "  복사 방식: $copy_tool"
    log_output "  소스: $SOURCE_DIR"
    log_output "  타겟: $TEST_DIR"
    
    # duration을 전역 변수로 저장하고, return은 성공(0) 반환
    INIT_DURATION=$duration
    return 0
}

# 저장소 초기화 (clone 방식)
init_repository_clone() {
    log_color "${YELLOW}>>> 저장소 초기화 (Clone 방식)${NC}"
    
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
    
    # duration을 전역 변수로 저장하고, return은 성공(0) 반환
    INIT_DURATION=$duration
    return 0
}

# ==================== 최적화 단계 함수 ====================

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
    
    # duration을 전역 변수로 저장
    OPT_DURATION=$duration
    return 0
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
    # 공백으로 구분된 브랜치를 배열로 변환하여 전달
    IFS=' ' read -ra branch_array <<< "$scope_branches"
    $GA_CMD opt quick set-branch-scope "${branch_array[@]}" 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 branch-scope 적용 중..."
    # 공백으로 구분된 브랜치를 배열로 변환하여 전달
    $GA_CMD opt submodule set-branch-scope "${branch_array[@]}" 2>&1 | tee -a "$LOG_FILE"
    
    # Accurate 모드에서만 prune 실행
    if [ "$MODE" = "accurate" ]; then
        log_output "Prune 실행 중 (accurate 모드)..."
        git fetch --prune --no-tags 2>&1 | tee -a "$LOG_FILE"
        
        # 서브모듈도 prune
        git submodule foreach 'git fetch --prune --no-tags' 2>&1 | tee -a "$LOG_FILE"
    fi
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Branch Scope 완료 (${duration}s)${NC}"
    
    # duration을 전역 변수로 저장
    OPT_DURATION=$duration
    return 0
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
    
    # duration을 전역 변수로 저장
    OPT_DURATION=$duration
    return 0
}

# Sparse-checkout 최적화
apply_sparse() {
    local patterns="${1:-src/ docs/ *.md *.json}"
    log_color "${YELLOW}>>> Sparse-checkout 최적화${NC}"
    
    start_timer
    
    log_output "Sparse 패턴 적용: $patterns"
    $GA_CMD opt sparse add $patterns 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Sparse-checkout 완료 (${duration}s)${NC}"
    
    # duration을 전역 변수로 저장
    OPT_DURATION=$duration
    return 0
}

# ==================== 복구 함수 ====================

# Shallow 복구 (전체 히스토리)
unshallow_repository() {
    log_color "${YELLOW}>>> Shallow 복구 (전체 히스토리)${NC}"
    
    start_timer
    
    # 메인 저장소
    log_output "메인 저장소 unshallow 중..."
    $GA_CMD opt quick unshallow 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 unshallow 중..."
    $GA_CMD opt submodule unshallow 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Unshallow 완료 (${duration}s)${NC}"
    
    return $duration
}

# Slim 모드에서 Full 모드로
to_full_repository() {
    log_color "${YELLOW}>>> Full 모드로 전환${NC}"
    
    start_timer
    
    # 메인 저장소
    log_output "메인 저장소 to-full 중..."
    $GA_CMD opt quick to-full 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 to-full 중..."
    $GA_CMD opt submodule to-full 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ To-full 완료 (${duration}s)${NC}"
    
    return $duration
}

# Branch scope 해제
clear_branch_scope() {
    local force="${1:-false}"
    log_color "${YELLOW}>>> Branch Scope 해제${NC}"
    
    start_timer
    
    local force_flag=""
    if [ "$force" = true ]; then
        force_flag="-f"
    fi
    
    # 메인 저장소
    log_output "메인 저장소 branch-scope 해제 중..."
    $GA_CMD opt quick clear-branch-scope $force_flag 2>&1 | tee -a "$LOG_FILE"
    
    # 서브모듈
    log_output "서브모듈 branch-scope 해제 중..."
    $GA_CMD opt submodule clear-branch-scope $force_flag 2>&1 | tee -a "$LOG_FILE"
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Branch Scope 해제 완료 (${duration}s)${NC}"
    
    return $duration
}

# Auto Find Merge Base - 브랜치 스코프로 병합점 자동 찾기
auto_find_merge_base() {
    local branches="${1:-$BRANCHES}"  # 브랜치 리스트 (공백 구분)
    local force="${2:-false}"          # 강제 실행 여부
    local quiet="${3:-$QUIET_MODE}"    # 조용한 모드
    
    log_color "${YELLOW}>>> Auto Find Merge Base (Branch Scope)${NC}"
    
    start_timer
    
    # 플래그 설정
    local flags=""
    if [ "$force" = true ]; then
        flags="$flags -f"
    fi
    if [ "$quiet" = true ]; then
        flags="$flags -q"
    fi
    
    # 브랜치가 지정되지 않으면 현재 브랜치와 main/master 사용
    if [ -z "$branches" ]; then
        local current_branch=$(git branch --show-current)
        local main_branch="main"
        
        # main 브랜치가 없으면 master 시도
        if ! git show-ref --quiet refs/heads/main; then
            if git show-ref --quiet refs/heads/master; then
                main_branch="master"
            else
                log_color "${RED}❌ main 또는 master 브랜치를 찾을 수 없습니다${NC}"
                return 1
            fi
        fi
        
        branches="$current_branch $main_branch"
        log_output "기본 브랜치 사용: $branches"
    fi
    
    log_output "병합점 찾기: $branches"
    
    # ga opt quick auto 명령 실행
    $GA_CMD opt quick auto $branches $flags 2>&1 | tee -a "$LOG_FILE"
    local exit_code=$?
    
    local duration=$(end_timer)
    
    if [ $exit_code -eq 0 ]; then
        log_color "${GREEN}✓ Auto Find Merge Base 완료 (${duration}s)${NC}"
    else
        log_color "${RED}✗ Auto Find Merge Base 실패 (${duration}s)${NC}"
    fi
    
    return $duration
}

# Auto Find Merge Base로 Branch Scope 설정
set_branch_scope_auto() {
    local quiet="${1:-$QUIET_MODE}"
    
    log_color "${YELLOW}>>> Branch Scope 자동 설정 (Auto Find)${NC}"
    
    start_timer
    
    # 현재 브랜치 확인
    local current_branch=$(git branch --show-current)
    local base_branch=""
    
    # 기본 브랜치 찾기
    for branch in main master develop; do
        if git show-ref --quiet refs/heads/$branch; then
            base_branch=$branch
            break
        fi
    done
    
    if [ -z "$base_branch" ]; then
        log_color "${RED}❌ 기본 브랜치를 찾을 수 없습니다${NC}"
        return 1
    fi
    
    log_output "브랜치 병합점 분석: $current_branch ← $base_branch"
    
    # 병합점 찾기 (quiet 모드)
    local merge_base=$($GA_CMD opt quick auto $current_branch $base_branch -q 2>/dev/null | head -1)
    
    if [ -n "$merge_base" ]; then
        log_output "병합점 발견: $merge_base"
        
        # 병합점 이후의 브랜치들 찾기
        local scope_branches="$current_branch $base_branch"
        
        # 관련 브랜치 추가 (병합점 이후에 분기된 브랜치들)
        local related_branches=$(git branch --contains "$merge_base" --no-color | sed 's/^[* ]*//' | tr '\n' ' ')
        
        if [ -n "$related_branches" ]; then
            scope_branches="$scope_branches $related_branches"
            # 중복 제거
            scope_branches=$(echo "$scope_branches" | tr ' ' '\n' | sort -u | tr '\n' ' ')
        fi
        
        log_output "Branch Scope 적용: $scope_branches"
        
        # Branch Scope 설정
        apply_branch_scope
    else
        log_color "${YELLOW}⚠ 병합점을 찾을 수 없어 전체 브랜치 사용${NC}"
    fi
    
    local duration=$(end_timer)
    log_color "${GREEN}✓ Branch Scope 자동 설정 완료 (${duration}s)${NC}"
    
    return $duration
}

# Export functions for use in other modules
export -f init_repository init_repository_clone
export -f apply_shallow apply_branch_scope apply_to_slim apply_sparse
export -f unshallow_repository to_full_repository clear_branch_scope
export -f auto_find_merge_base set_branch_scope_auto