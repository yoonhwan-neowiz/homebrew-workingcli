#!/bin/bash

# Git Optimization Benchmark - Measurement Functions
# 저장소 크기 측정, 델타 계산 등

# ==================== 저장소 크기 측정 ====================

# Git 저장소 크기 측정 (Git objects)
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

# Git objects 디렉토리 크기만 측정
measure_git_objects() {
    local repo_path=${1:-.}
    local objects_size=0
    
    if [ -d "$repo_path/.git/objects" ]; then
        objects_size=$(du -sk "$repo_path/.git/objects" 2>/dev/null | cut -f1)
        echo $((objects_size * 1024))
    else
        echo 0
    fi
}

# ==================== 측정 기록 ====================

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

# 간단한 측정 (로그 없이)
measure_quick() {
    local repo_path=${1:-$TEST_DIR}
    
    local git_store=$(measure_git_store "$repo_path")
    local worktree=$(measure_worktree "$repo_path")
    local total=$((git_store + worktree))
    
    echo "$git_store $worktree $total"
}

# ==================== 델타 계산 ====================

# 델타 계산 및 출력
calculate_delta() {
    local prev_git=$1
    local curr_git=$2
    local prev_wt=$3
    local curr_wt=$4
    local prev_total=$5
    local curr_total=$6
    local duration=${7:-0}
    
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
    if [ $duration -gt 0 ]; then
        local efficiency=$(echo "scale=1; $delta_total / 1048576 / $duration" | bc)
        log_output "    효율:      ${efficiency} MB/s"
    fi
}

# 간단한 델타 계산 (퍼센트만)
calculate_simple_delta() {
    local prev=$1
    local curr=$2
    
    if [ $prev -gt 0 ]; then
        local delta=$((prev - curr))
        local percent=$(echo "scale=1; $delta * 100 / $prev" | bc)
        echo "${percent}%"
    else
        echo "N/A"
    fi
}

# ==================== 통계 함수 ====================

# 평균 계산
calculate_average() {
    local sum=0
    local count=0
    
    for val in "$@"; do
        sum=$((sum + val))
        count=$((count + 1))
    done
    
    if [ $count -gt 0 ]; then
        echo $((sum / count))
    else
        echo 0
    fi
}

# 최대값 찾기
find_max() {
    local max=0
    
    for val in "$@"; do
        if [ $val -gt $max ]; then
            max=$val
        fi
    done
    
    echo $max
}

# 최소값 찾기
find_min() {
    local min=999999999999
    
    for val in "$@"; do
        if [ $val -lt $min ]; then
            min=$val
        fi
    done
    
    echo $min
}

# ==================== Git 상태 측정 ====================

# shallow depth 확인
get_shallow_depth() {
    local repo_path=${1:-.}
    
    if [ -f "$repo_path/.git/shallow" ]; then
        local depth=$(git -C "$repo_path" rev-list --count HEAD 2>/dev/null || echo "unknown")
        echo $depth
    else
        echo "full"
    fi
}

# partial clone 필터 확인
get_partial_filter() {
    local repo_path=${1:-.}
    
    git -C "$repo_path" config remote.origin.partialclonefilter 2>/dev/null || echo "none"
}

# sparse-checkout 경로 수 확인
count_sparse_paths() {
    local repo_path=${1:-.}
    
    if [ -f "$repo_path/.git/info/sparse-checkout" ]; then
        wc -l < "$repo_path/.git/info/sparse-checkout"
    else
        echo 0
    fi
}

# Export functions for use in other modules
export -f measure_git_store measure_worktree measure_git_objects
export -f measure_step measure_quick
export -f calculate_delta calculate_simple_delta
export -f calculate_average find_max find_min
export -f get_shallow_depth get_partial_filter count_sparse_paths