#!/bin/bash

# Git Optimization Benchmark - Reporting Functions
# 결과 분석 및 리포트 생성

# ==================== 결과 요약 ====================

# 최종 요약 생성
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
복사 방식: $([ "$USE_CLONE" = true ] && echo "Clone" || echo "rusync/rsync")
===========================================

EOF
    
    # JSONL 파일 분석
    if [ -f "$JSONL_FILE" ]; then
        analyze_jsonl >> "$SUMMARY_FILE"
    fi
    
    # 결과 출력
    log_color "${GREEN}=== 요약 ===${NC}"
    cat "$SUMMARY_FILE" | tee -a "$LOG_FILE"
    
    log_color "${GREEN}✓ 결과 파일 저장됨:${NC}"
    log_output "  - 로그: $LOG_FILE"
    log_output "  - 데이터: $JSONL_FILE"
    log_output "  - 요약: $SUMMARY_FILE"
}

# JSONL 파일 분석
analyze_jsonl() {
    python3 <<EOF
import json
import sys

steps = []
with open('$JSONL_FILE', 'r') as f:
    for line in f:
        if line.strip():
            steps.append(json.loads(line))

if not steps:
    print("데이터가 없습니다.")
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

# 성능 지표
total_duration = sum(step.get('duration', 0) for step in steps)
print(f"\n[성능 지표]")
print(f"  총 실행 시간: {total_duration}초")
print(f"  평균 효율: {(total_saved / 1048576 / total_duration):.1f} MB/s" if total_duration > 0 else "  평균 효율: N/A")
EOF
}

# ==================== 개별 리포트 생성 ====================

# 간단한 리포트 생성
generate_quick_report() {
    local title="${1:-Quick Report}"
    
    echo "=== $title ===" >> "$LOG_FILE"
    echo "시간: $(date)" >> "$LOG_FILE"
    echo "저장소: $TEST_DIR" >> "$LOG_FILE"
    
    if [ -d "$TEST_DIR" ]; then
        local sizes=($(measure_quick))
        echo "Git Store: $(bytes_to_human ${sizes[0]})" >> "$LOG_FILE"
        echo "Worktree: $(bytes_to_human ${sizes[1]})" >> "$LOG_FILE"
        echo "Total: $(bytes_to_human ${sizes[2]})" >> "$LOG_FILE"
    fi
    
    echo "" >> "$LOG_FILE"
}

# CSV 형식 리포트 생성
generate_csv_report() {
    local csv_file="${RESULT_DIR}/benchmark_${DATE_STR}.csv"
    
    # CSV 헤더
    echo "Step,Name,Timestamp,Git_Store_GB,Worktree_GB,Total_GB,Duration_s,Delta_GB,Delta_Percent" > "$csv_file"
    
    # JSONL을 CSV로 변환
    if [ -f "$JSONL_FILE" ]; then
        python3 <<EOF
import json
import csv

steps = []
with open('$JSONL_FILE', 'r') as f:
    for line in f:
        if line.strip():
            steps.append(json.loads(line))

with open('$csv_file', 'a', newline='') as csvfile:
    writer = csv.writer(csvfile)
    
    for i, step in enumerate(steps):
        git_store_gb = step['git_store'] / 1073741824
        worktree_gb = step['worktree'] / 1073741824
        total_gb = step['total'] / 1073741824
        
        if i > 0:
            prev = steps[i-1]
            delta_gb = (prev['total'] - step['total']) / 1073741824
            delta_percent = (delta_gb / (prev['total'] / 1073741824)) * 100 if prev['total'] > 0 else 0
        else:
            delta_gb = 0
            delta_percent = 0
        
        writer.writerow([
            step['step'],
            step['name'],
            step['timestamp'],
            f"{git_store_gb:.2f}",
            f"{worktree_gb:.2f}",
            f"{total_gb:.2f}",
            step.get('duration', 0),
            f"{delta_gb:.2f}",
            f"{delta_percent:.1f}"
        ])
EOF
        log_output "CSV 리포트 생성됨: $csv_file"
    fi
}

# ==================== 비교 리포트 ====================

# 두 개의 벤치마크 결과 비교
compare_benchmarks() {
    local file1="$1"
    local file2="$2"
    
    if [ ! -f "$file1" ] || [ ! -f "$file2" ]; then
        die "비교할 파일을 찾을 수 없습니다"
    fi
    
    python3 <<EOF
import json

def load_jsonl(filename):
    steps = []
    with open(filename, 'r') as f:
        for line in f:
            if line.strip():
                steps.append(json.loads(line))
    return steps

steps1 = load_jsonl('$file1')
steps2 = load_jsonl('$file2')

print("=== 벤치마크 비교 ===")
print(f"파일 1: $file1")
print(f"파일 2: $file2")
print()

# 최종 크기 비교
if steps1 and steps2:
    final1 = steps1[-1]['total'] / 1073741824
    final2 = steps2[-1]['total'] / 1073741824
    diff = final2 - final1
    diff_percent = (diff / final1 * 100) if final1 > 0 else 0
    
    print(f"최종 크기:")
    print(f"  파일 1: {final1:.2f} GB")
    print(f"  파일 2: {final2:.2f} GB")
    print(f"  차이: {diff:+.2f} GB ({diff_percent:+.1f}%)")
    
    # 실행 시간 비교
    time1 = sum(s.get('duration', 0) for s in steps1)
    time2 = sum(s.get('duration', 0) for s in steps2)
    time_diff = time2 - time1
    
    print(f"\n실행 시간:")
    print(f"  파일 1: {time1}초")
    print(f"  파일 2: {time2}초")
    print(f"  차이: {time_diff:+d}초")
EOF
}

# ==================== 실시간 모니터링 ====================

# 실시간 진행 상황 표시
show_live_progress() {
    local current_step="$1"
    local total_steps="$2"
    
    if [ -d "$TEST_DIR" ]; then
        local sizes=($(measure_quick))
        local percent=$((current_step * 100 / total_steps))
        
        printf "\r[%3d%%] Step %d/%d | Git: %s | Worktree: %s | Total: %s" \
            $percent $current_step $total_steps \
            "$(bytes_to_human ${sizes[0]})" \
            "$(bytes_to_human ${sizes[1]})" \
            "$(bytes_to_human ${sizes[2]})"
    fi
}

# ==================== 그래프 생성 (선택적) ====================

# ASCII 그래프 생성
generate_ascii_graph() {
    if [ ! -f "$JSONL_FILE" ]; then
        return
    fi
    
    python3 <<EOF
import json

steps = []
with open('$JSONL_FILE', 'r') as f:
    for line in f:
        if line.strip():
            steps.append(json.loads(line))

if not steps:
    exit()

# 크기 그래프
print("\n[크기 변화 그래프]")
max_size = max(s['total'] for s in steps)

for step in steps:
    size_percent = int(step['total'] * 50 / max_size) if max_size > 0 else 0
    bar = '#' * size_percent
    gb = step['total'] / 1073741824
    print(f"{step['name']:15} [{bar:50}] {gb:.1f}GB")
EOF
}

# Export functions for use in other modules
export -f generate_summary analyze_jsonl
export -f generate_quick_report generate_csv_report
export -f compare_benchmarks show_live_progress
export -f generate_ascii_graph