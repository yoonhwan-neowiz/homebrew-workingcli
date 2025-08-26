#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Git 명령어 테스트 시작${NC}"

# 테스트 실행
go test -v ./test/ -run "TestGit.*"

# 테스트 결과 확인
if [ $? -eq 0 ]; then
    echo -e "${GREEN}모든 Git 테스트가 성공적으로 완료되었습니다.${NC}"
else
    echo -e "${RED}일부 Git 테스트가 실패했습니다.${NC}"
    exit 1
fi 