#!/bin/bash

# Git 최적화 워크플로우 2-3-4 연속 실행
# 모든 시나리오에서 --use-clone 옵션 사용

echo "========================================="
echo "   Git 최적화 워크플로우 2-3-4 실행"
echo "========================================="
echo "시작 시간: $(date)"
echo ""

# 시나리오 02: Master 클론 셋업 (새 디렉토리 생성)
echo ">>> [1/3] 시나리오 02: Master 클론 셋업 시작"
./benchmark/ga-opt-benchmark.sh \
    --scenario 02-setup-master-workflow \
    --use-clone

if [ $? -ne 0 ]; then
    echo "❌ 시나리오 02 실패"
    exit 1
fi
echo ""

# 시나리오 03: 브랜치 전환 워크플로우 (기존 디렉토리 사용)
echo ">>> [2/3] 시나리오 03: 브랜치 전환 워크플로우 시작"
./benchmark/ga-opt-benchmark.sh \
    --scenario 03-branch-switch-workflow \
    --use-clone

if [ $? -ne 0 ]; then
    echo "❌ 시나리오 03 실패"
    exit 1
fi
echo ""

# 시나리오 04: Full Reset (기존 디렉토리 정리)
echo ">>> [3/3] 시나리오 04: Full Reset to Master Shallow 1 시작"
./benchmark/ga-opt-benchmark.sh \
    --scenario 04-full-reset-shallow \
    --use-clone \
    --clean-mode full

if [ $? -ne 0 ]; then
    echo "❌ 시나리오 04 실패"
    exit 1
fi

echo ""
echo "========================================="
echo "   워크플로우 완료!"
echo "========================================="
echo "종료 시간: $(date)"