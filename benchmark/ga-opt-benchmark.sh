#!/bin/bash

# Git Optimization Benchmark Script - Modular Version
# 목적: ga 명령어를 활용한 단계별 Git 최적화 효과 측정
# 구조: 모듈화된 라이브러리 + 시나리오 기반 실행

set -e  # 오류 시 즉시 중단

# ==================== 스크립트 경로 설정 ====================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LIB_DIR="$SCRIPT_DIR/lib"
SCENARIOS_DIR="$SCRIPT_DIR/scenarios"

# ==================== 라이브러리 로드 ====================
# 순서 중요: base → utils → measure → optimize → report
source "$LIB_DIR/base.sh"
source "$LIB_DIR/utils.sh"
source "$LIB_DIR/measure.sh"
source "$LIB_DIR/optimize.sh"
source "$LIB_DIR/report.sh"

# ==================== 옵션 파싱 ====================
SCENARIO="99-full-optimization"  # 기본 시나리오
SCENARIO_FILE=""  # 직접 지정된 시나리오 파일

parse_options() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --scenario|-s)
                SCENARIO="$2"
                shift 2
                ;;
            --scenario-file|-f)
                SCENARIO_FILE="$2"
                shift 2
                ;;
            --list-scenarios|-l)
                list_scenarios
                exit 0
                ;;
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
            --source)
                SOURCE_TEMPLATE="$2"
                shift 2
                ;;
            --use-clone)
                USE_CLONE=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
}

# ==================== 도움말 ====================
show_help() {
    cat <<EOF
사용법: $0 [OPTIONS]

기본 옵션:
  --scenario, -s NAME    실행할 시나리오 (기본: 01-full-optimization)
  --scenario-file FILE   커스텀 시나리오 파일 경로
  --list-scenarios, -l   사용 가능한 시나리오 목록 표시
  
벤치마크 옵션:
  --mode [fast|accurate] 측정 모드 (기본: fast)
  --skip-slim           to-slim 단계 제외
  --skip-sparse         sparse 단계 제외
  --branches "list"     branch-scope 브랜치 목록
  --source PATH         템플릿 소스 디렉토리 (기본: ~/Work/DesignB4)
  --use-clone           기존 clone 방식 사용 (호환성)
  --quiet, -q           조용한 모드
  --help, -h            도움말 표시

시나리오 예제:
  $0                                    # 기본 전체 최적화 실행
  $0 --scenario 01-full-optimization    # 전체 최적화
  $0 --scenario 02-ci-build            # CI/CD 최적화
  $0 --scenario-file custom.sh         # 커스텀 시나리오

사용 가능한 시나리오:
EOF
    list_scenarios
}

# ==================== 시나리오 관리 ====================
list_scenarios() {
    echo ""
    echo "내장 시나리오:"
    for scenario in "$SCENARIOS_DIR"/*.sh; do
        if [ -f "$scenario" ]; then
            local name=$(basename "$scenario" .sh)
            local desc=""
            
            # 시나리오 파일에서 설명 추출
            if grep -q "^SCENARIO_DESC=" "$scenario"; then
                desc=$(grep "^SCENARIO_DESC=" "$scenario" | head -1 | cut -d'"' -f2)
            fi
            
            printf "  %-25s %s\n" "$name" "$desc"
        fi
    done
    echo ""
}

# 시나리오 로드
load_scenario() {
    local scenario_path=""
    
    # 시나리오 파일 경로 결정
    if [ -n "$SCENARIO_FILE" ]; then
        # 직접 지정된 파일
        if [ -f "$SCENARIO_FILE" ]; then
            scenario_path="$SCENARIO_FILE"
        else
            die "시나리오 파일을 찾을 수 없습니다: $SCENARIO_FILE"
        fi
    else
        # 이름으로 지정된 시나리오
        scenario_path="$SCENARIOS_DIR/${SCENARIO}.sh"
        if [ ! -f "$scenario_path" ]; then
            # .sh 확장자 없이 시도
            scenario_path="$SCENARIOS_DIR/${SCENARIO}"
            if [ ! -f "$scenario_path" ]; then
                die "시나리오를 찾을 수 없습니다: $SCENARIO"
            fi
        fi
    fi
    
    log_output "시나리오 로드: $scenario_path"
    source "$scenario_path"
    
    # 필수 함수 확인
    if ! declare -f run_scenario >/dev/null; then
        die "시나리오에 run_scenario 함수가 정의되어 있지 않습니다"
    fi
}

# ==================== 메인 실행 ====================
main() {
    # 옵션 파싱
    parse_options "$@"
    
    # 환경 검증
    validate_environment
    
    # 디렉토리 생성
    ensure_directories
    
    # 시작
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}   Git 최적화 벤치마크 시작${NC}"
    log_color "${CYAN}========================================${NC}"
    log_output "시간: $(date)"
    log_output "모드: $MODE"
    log_output "Skip Slim: $SKIP_SLIM"
    log_output "Skip Sparse: $SKIP_SPARSE"
    log_output "복사 방식: $([ "$USE_CLONE" = true ] && echo "Clone" || echo "rusync/rsync")"
    log_output ""
    
    # 시나리오 로드 및 실행
    load_scenario
    
    # 시나리오 설정 (있으면)
    if declare -f scenario_setup >/dev/null; then
        scenario_setup
    fi
    
    # 시나리오 실행
    run_scenario
    
    # 시나리오 정리 (있으면)
    if declare -f scenario_cleanup >/dev/null; then
        scenario_cleanup
    fi
    
    # 최종 요약 생성
    generate_summary
    
    # CSV 리포트 생성 (선택적)
    if [ "$MODE" = "accurate" ]; then
        generate_csv_report
    fi
    
    # ASCII 그래프 생성 (선택적)
    if [ "$QUIET_MODE" = false ]; then
        generate_ascii_graph
    fi
    
    log_color "${CYAN}========================================${NC}"
    log_color "${CYAN}   벤치마크 완료!${NC}"
    log_color "${CYAN}========================================${NC}"
}

# 스크립트 실행
main "$@"