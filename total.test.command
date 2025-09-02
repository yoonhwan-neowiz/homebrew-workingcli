#!/bin/bash

# Git Optimization Test Script
# Target: ~/Work/DesignB5_test (20GB: main 4GB + submodules 16GB)

TEST_DIR="$HOME/Work/DesignB5_test"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LOG_DIR="$SCRIPT_DIR/logs"

# Use built ga command from WorkingCli
GA_CMD="$SCRIPT_DIR/ga"
DATE_STR=$(date +"%Y%m%d_%H%M%S")
LOG_FILE="$LOG_DIR/${DATE_STR}.log"
RESULT_FILE="$TEST_DIR/optimization_results_${DATE_STR}.txt"

# Create logs directory if it doesn't exist
mkdir -p "$LOG_DIR"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to log both to console and file
log_output() {
    echo "$@" | tee -a "$LOG_FILE"
}

log_color() {
    echo -e "$@" | tee -a "$LOG_FILE"
}

# Start test
log_color "${CYAN}========================================${NC}"
log_color "${CYAN}   Git Optimization Test Start${NC}"
log_color "${CYAN}========================================${NC}"
log_output ""
log_output "Log file: $LOG_FILE"
log_output "Result file: $RESULT_FILE"
log_output ""

# Check if TEST_DIR exists, if not clone it
if [ ! -d "$TEST_DIR" ]; then
    log_color "${YELLOW}TEST_DIR이 존재하지 않습니다. Git clone을 시작합니다...${NC}"
    log_output "Repository: git@git.nwz.kr:gamfs4/designb4"
    log_output "Target: $TEST_DIR"
    
    # Clone repository
    git clone git@git.nwz.kr:gamfs4/designb4 "$TEST_DIR" 2>&1 | tee -a "$LOG_FILE"
    
    if [ $? -ne 0 ]; then
        log_color "${RED}Git clone 실패!${NC}"
        exit 1
    fi
    
    # Initialize submodules
    log_color "${BLUE}서브모듈 초기화 중...${NC}"
    cd "$TEST_DIR"
    
    # Set git pull configuration to avoid divergent branches error
    git config pull.rebase false 2>&1 | tee -a "$LOG_FILE"
    
    git submodule update --init 2>&1 | tee -a "$LOG_FILE"
    
    if [ $? -ne 0 ]; then
        log_color "${RED}서브모듈 초기화 실패!${NC}"
        exit 1
    fi
    
    log_color "${GREEN}저장소 clone 및 서브모듈 초기화 완료!${NC}"
else
    log_color "${GREEN}TEST_DIR이 이미 존재합니다: $TEST_DIR${NC}"
fi

cd "$TEST_DIR" || exit 1

# Ensure git pull configuration is set
git config pull.rebase false

# Initialize result file
echo "Git Optimization Test Results" > "$RESULT_FILE"
echo "Test Date: $(date)" >> "$RESULT_FILE"
echo "Repository: $TEST_DIR" >> "$RESULT_FILE"
echo "Log File: $LOG_FILE" >> "$RESULT_FILE"
echo "=======================================" >> "$RESULT_FILE"
echo "" >> "$RESULT_FILE"

# Size measurement function
measure_size() {
    local label=$1
    log_color "${YELLOW}>>> $label${NC}"
    echo "$label" >> "$RESULT_FILE"
    echo "$label" >> "$LOG_FILE"
    
    # dust command output
    echo "Dust output:" >> "$RESULT_FILE"
    echo "Dust output:" >> "$LOG_FILE"
    dust -d 1 | tee -a "$RESULT_FILE" | tee -a "$LOG_FILE"
    
    # .git folder size
    echo "" >> "$RESULT_FILE"
    echo "Git folder size:" >> "$RESULT_FILE"
    echo "" >> "$LOG_FILE"
    echo "Git folder size:" >> "$LOG_FILE"
    du -sh .git 2>/dev/null | tee -a "$RESULT_FILE" | tee -a "$LOG_FILE"
    
    # Submodules size
    echo "Submodules size:" >> "$RESULT_FILE"
    echo "Submodules size:" >> "$LOG_FILE"
    du -sh */ 2>/dev/null | grep -v "^du:" | tee -a "$RESULT_FILE" | tee -a "$LOG_FILE"
    
    echo "---------------------------------------" >> "$RESULT_FILE"
    echo "---------------------------------------" >> "$LOG_FILE"
    echo "" >> "$RESULT_FILE"
    echo "" >> "$LOG_FILE"
}

# Clean local branches function
clean_local_branches() {
    log_color "${BLUE}Cleaning local branches...${NC}"
    
    # Get current branch
    CURRENT_BRANCH=$(git branch --show-current)
    
    # Switch to master or main if not already
    if [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "main" ]; then
        # Try master first
        if git show-ref --verify --quiet refs/heads/master; then
            git checkout master 2>&1 | tee -a "$LOG_FILE"
        elif git show-ref --verify --quiet refs/heads/main; then
            git checkout main 2>&1 | tee -a "$LOG_FILE"
        else
            # Create and checkout to HEAD detached state
            git checkout --detach HEAD 2>&1 | tee -a "$LOG_FILE"
        fi
    fi
    
    # Delete all local branches except master/main
    git branch | grep -v "master\|main\|*" | xargs -r git branch -D 2>&1 | tee -a "$LOG_FILE" || true
    
    # Clean up submodules branches too
    git submodule foreach 'git checkout --detach HEAD 2>/dev/null || git checkout master 2>/dev/null || git checkout main 2>/dev/null' 2>&1 | tee -a "$LOG_FILE"
    
    sleep 1
}

# Restore function (including full space recovery)
restore_full() {
    log_color "${BLUE}Performing full restore...${NC}"
    
    # Clean branches first
    clean_local_branches
    
    # Run gc and prune
    git gc --prune=now --aggressive 2>&1 | tee -a "$LOG_FILE"
    git fetch origin --prune 2>&1 | tee -a "$LOG_FILE"
    
    # Same for submodules
    git submodule foreach 'git gc --prune=now --aggressive && git fetch origin --prune' 2>&1 | tee -a "$LOG_FILE"
    
    sleep 2
}

# Initial cleanup
clean_local_branches

# ====================================
# 1. Baseline measurement (current to_full state)
# ====================================
log_color "${GREEN}[1/6] Baseline measurement (to_full state)${NC}"
measure_size "1. BASELINE (to_full state)"

# ====================================
# 2. Shallow test (Main + Submodule)
# ====================================
log_color "${GREEN}[2/6] Shallow optimization test${NC}"
log_color "${BLUE}Applying: shallow depth=1${NC}"

log_color "${CYAN}[2-1] Main shallow depth=1${NC}"
$GA_CMD opt quick shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[2-2] Submodule shallow depth=1${NC}"
$GA_CMD opt submodule shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
sleep 2

measure_size "2. SHALLOW (depth=1)"

# Restore
log_color "${BLUE}Restoring: unshallow${NC}"
log_color "${CYAN}[2-3] Main unshallow${NC}"
$GA_CMD opt quick unshallow -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[2-4] Submodule unshallow${NC}"
$GA_CMD opt submodule unshallow -q 2>&1 | tee -a "$LOG_FILE"
restore_full
measure_size "2. RESTORED from Shallow"

# ====================================
# 3. Branch Scope test (Main + Submodule)
# ====================================
log_color "${GREEN}[3/6] Branch Scope test${NC}"
log_color "${BLUE}Applying: branch scope (master, live59.b/5904.7)${NC}"

log_color "${CYAN}[3-1] Main set-branch-scope${NC}"
$GA_CMD opt quick set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[3-2] Submodule set-branch-scope${NC}"
$GA_CMD opt submodule set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
sleep 2

measure_size "3. BRANCH SCOPE (2 branches)"

# Restore
log_color "${BLUE}Restoring: clear branch scope${NC}"
log_color "${CYAN}[3-3] Main clear-branch-scope${NC}"
$GA_CMD opt quick clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[3-4] Submodule clear-branch-scope${NC}"
$GA_CMD opt submodule clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
restore_full
measure_size "3. RESTORED from Branch Scope"

# ====================================
# 4. Shallow + Branch Scope combo (without to-slim)
# ====================================
log_color "${GREEN}[4/6] Shallow + Branch Scope combo test${NC}"
log_color "${BLUE}Applying: shallow + branch scope${NC}"

log_color "${CYAN}[4-1] Main shallow 1${NC}"
$GA_CMD opt quick shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-2] Submodule shallow 1${NC}"
$GA_CMD opt submodule shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-3] Main set-branch-scope${NC}"
$GA_CMD opt quick set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-4] Submodule set-branch-scope${NC}"
$GA_CMD opt submodule set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
sleep 2

measure_size "4. SHALLOW + BRANCH SCOPE"

# Restore
log_color "${BLUE}Restoring: unshallow + clear scope${NC}"
log_color "${CYAN}[4-5] Main unshallow${NC}"
$GA_CMD opt quick unshallow -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-6] Submodule unshallow${NC}"
$GA_CMD opt submodule unshallow -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-7] Main clear-branch-scope${NC}"
$GA_CMD opt quick clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[4-8] Submodule clear-branch-scope${NC}"
$GA_CMD opt submodule clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
restore_full
measure_size "4. RESTORED from Combo"

# ====================================
# 5. Partial Clone (to-slim) test
# ====================================
log_color "${GREEN}[5/6] Partial Clone (to-slim) test${NC}"
log_color "${BLUE}Applying: to-slim${NC}"

log_color "${CYAN}[5-1] Main to-slim${NC}"
$GA_CMD opt quick to-slim -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[5-2] Submodule to-slim${NC}"
$GA_CMD opt submodule to-slim -q 2>&1 | tee -a "$LOG_FILE"
sleep 2

measure_size "5. TO-SLIM (partial clone)"

# Restore
log_color "${BLUE}Restoring: to-full${NC}"
log_color "${CYAN}[5-3] Main to-full${NC}"
$GA_CMD opt quick to-full -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[5-4] Submodule to-full${NC}"
$GA_CMD opt submodule to-full -q 2>&1 | tee -a "$LOG_FILE"
restore_full
measure_size "5. RESTORED from Slim"

# ====================================
# 6. Final 3-stage combo (Slim + Shallow + Scope)
# ====================================
log_color "${GREEN}[6/6] Final 3-stage combo test${NC}"
log_color "${BLUE}Applying: to-slim -> shallow -> branch scope${NC}"

log_color "${CYAN}[6-1] Main to-slim${NC}"
$GA_CMD opt quick to-slim -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-2] Submodule to-slim${NC}"
$GA_CMD opt submodule to-slim -q 2>&1 | tee -a "$LOG_FILE"
sleep 1

log_color "${CYAN}[6-3] Main shallow 1${NC}"
$GA_CMD opt quick shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-4] Submodule shallow 1${NC}"
$GA_CMD opt submodule shallow 1 -q 2>&1 | tee -a "$LOG_FILE"
sleep 1

log_color "${CYAN}[6-5] Main set-branch-scope${NC}"
$GA_CMD opt quick set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-6] Submodule set-branch-scope${NC}"
$GA_CMD opt submodule set-branch-scope master live59.b/5904.7 2>&1 | tee -a "$LOG_FILE"
sleep 2

measure_size "6. FULL COMBO (slim+shallow+scope)"

# Final restore
log_color "${BLUE}Final restore in progress...${NC}"
log_color "${CYAN}[6-7] Main to-full${NC}"
$GA_CMD opt quick to-full -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-8] Submodule to-full${NC}"
$GA_CMD opt submodule to-full -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-9] Main unshallow${NC}"
$GA_CMD opt quick unshallow -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-10] Submodule unshallow${NC}"
$GA_CMD opt submodule unshallow -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-11] Main clear-branch-scope${NC}"
$GA_CMD opt quick clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
log_color "${CYAN}[6-12] Submodule clear-branch-scope${NC}"
$GA_CMD opt submodule clear-branch-scope -f -q 2>&1 | tee -a "$LOG_FILE"
restore_full
measure_size "6. FINAL RESTORE"

# ====================================
# Result summary
# ====================================
log_color "${CYAN}========================================${NC}"
log_color "${CYAN}   Test Completed${NC}"
log_color "${CYAN}========================================${NC}"

echo "" >> "$RESULT_FILE"
echo "=======================================" >> "$RESULT_FILE"
echo "TEST COMPLETED: $(date)" >> "$RESULT_FILE"

log_color "${GREEN}All tests completed successfully!${NC}"
log_color "${GREEN}Result file: $RESULT_FILE${NC}"
log_color "${GREEN}Log file: $LOG_FILE${NC}"

# Simple summary output
log_output ""
log_color "${YELLOW}Key Results Summary:${NC}"
grep -E "^[0-9]\." "$RESULT_FILE" | while read line; do
    echo "  $line"
    grep -A1 "^$line" "$RESULT_FILE" | grep "Git folder size:" | head -1
done
