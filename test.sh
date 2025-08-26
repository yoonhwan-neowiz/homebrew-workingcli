#!/bin/bash

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 테스트 결과 저장 변수
FAILED_TESTS=()
PASSED_TESTS=0
TOTAL_TESTS=0

# 테스트 실행 함수
run_test() {
    local test_path=$1
    local test_name=$(basename "$test_path")
    ((TOTAL_TESTS++))
    
    echo -e "${YELLOW}실행 중: $test_name${NC}"
    if go test -v "$test_path"; then
        echo -e "${GREEN}통과: $test_name${NC}"
        ((PASSED_TESTS++))
    else
        echo -e "${RED}실패: $test_name${NC}"
        FAILED_TESTS+=("$test_name")
    fi
    echo "----------------------------------------"
}

# 시작 메시지
echo "테스트 실행 시작..."
echo "----------------------------------------"

# test 디렉토리의 모든 테스트 패키지 실행
for test_pkg in $(find ./test -type f -name "*_test.go" | xargs -n1 dirname | sort -u); do
    run_test "$test_pkg"
done

# 결과 출력
echo "테스트 실행 완료"
echo "----------------------------------------"
echo -e "총 테스트: $TOTAL_TESTS"
echo -e "${GREEN}통과: $PASSED_TESTS${NC}"
echo -e "${RED}실패: ${#FAILED_TESTS[@]}${NC}"

# 실패한 테스트 목록 출력
if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo "실패한 테스트 목록:"
    for test in "${FAILED_TESTS[@]}"; do
        echo -e "${RED}- $test${NC}"
    done
    exit 1
fi

exit 0 