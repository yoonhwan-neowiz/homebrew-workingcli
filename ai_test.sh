#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# CLAUDE_API_KEY 환경변수 확인
if [ -z "$CLAUDE_API_KEY" ]; then
    echo -e "${RED}CLAUDE_API_KEY 환경변수가 설정되지 않았습니다.${NC}"
    echo -e "${YELLOW}테스트를 실행하기 전에 CLAUDE_API_KEY를 설정해주세요:${NC}"
    echo -e "export CLAUDE_API_KEY=your-api-key"
    exit 1
fi

echo -e "${GREEN}AI 명령어 테스트 시작${NC}"

# 테스트 실행
go test -v ./test/ -run "Test(Commit|Analyze|ParseCommits).*"

# 테스트 결과 확인
if [ $? -eq 0 ]; then
    echo -e "${GREEN}모든 AI 테스트가 성공적으로 완료되었습니다.${NC}"
else
    echo -e "${RED}일부 AI 테스트가 실패했습니다.${NC}"
    exit 1
fi 