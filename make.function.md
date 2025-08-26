# Git 저장소 최적화 명령어 구현 가이드

## 📋 개요
이 문서는 Git 저장소 최적화를 위한 28개 명령어의 구현 상세를 담고 있습니다.
각 명령어는 PRD 기반으로 구체적인 구현 방법이 정의되어 있습니다.

## 🎯 구현 진행 상황 (0/28)
- [ ] 01. workflow - Git 최적화 워크플로우 가이드
- [ ] 02. commands - 전체 명령어 목록
- [ ] 03. status - 현재 최적화 상태 확인
- [ ] 04. to-slim - SLIM 모드로 전환
- [ ] 05. to-full - FULL 모드로 복원
- [ ] 06. expand-slim - 선택적 경로 확장
- [ ] 07. expand-filter - Partial Clone 필터 제거
- [ ] 08. expand-10 - 히스토리 10개 확장
- [ ] 09. expand-50 - 히스토리 50개 확장
- [ ] 10. expand-100 - 히스토리 100개 확장
- [ ] 11. find-merge - 병합 베이스 찾기
- [ ] 12. check-merge - 병합 가능 여부 확인
- [ ] 13. clone-slim - 최적화된 클론
- [ ] 14. migrate - 기존 저장소 SLIM 변환
- [ ] 15. performance - 성능 최적화 설정
- [ ] 16. expand-path - 특정 경로 확장
- [ ] 17. filter-branch - 브랜치별 필터 설정
- [ ] 18. clear-filter - 필터 완전 제거
- [ ] 19. restore-branch - 브랜치 전체 복원
- [ ] 20. shallow - 히스토리 줄이기
- [ ] 21. unshallow - 히스토리 복원
- [ ] 22. check-shallow - 히스토리 상태 확인
- [ ] 23. check-filter - 브랜치 필터 확인
- [ ] 24. backup-config - 설정 백업/복원
- [ ] 25. shallow-all - 모든 서브모듈 shallow 변환
- [ ] 26. unshallow-all - 모든 서브모듈 히스토리 복원
- [ ] 27. optimize-all - 모든 서브모듈 SLIM 최적화
- [ ] 28. status-all - 모든 서브모듈 상태 확인

---

## 📚 함수별 구현 상세

### 01. workflow (`src/cmd/optimized/help/01_workflow.go`)
**목적**: Git 최적화 워크플로우 가이드 표시
**구현 내용**:
```
1. SLIM과 FULL 모드의 차이점 설명
   - FULL: 전체 파일 히스토리와 모든 파일 포함 (약 103GB)
   - SLIM: 필수 파일과 최소 히스토리만 유지 (약 30MB)

2. 주요 워크플로우 4가지 안내:
   - INIT-SLIM: 신규 경량 클론 (∅ → SLIM)
   - MIGRATE-SLIM: 기존 저장소 경량화 (FULL → SLIM)
   - RESTORE-FULL: 전체 복원 (SLIM → FULL)
   - EXPAND-SLIM: 선택적 확장 (SLIM → SLIM+)

3. 각 워크플로우별 사용 시나리오 설명
4. 권장 사용 패턴과 예시 명령어 제공
```

### 02. commands (`src/cmd/optimized/help/02_commands.go`)
**목적**: 28개 전체 명령어 목록 표시
**구현 내용**: ✅ 이미 구현됨

### 03. status (`src/cmd/optimized/quick/03_status.go`)
**목적**: 현재 저장소의 최적화 상태 확인
**구현 내용**:
```bash
# 확인 항목:
1. Partial Clone 필터 상태
   git config remote.origin.partialclonefilter

2. Sparse Checkout 상태
   git sparse-checkout list

3. Shallow 상태
   git rev-parse --is-shallow-repository

4. 디스크 사용량
   du -sh .git        # .git 폴더 크기
   du -sh .           # 전체 프로젝트 크기

5. 서브모듈 상태
   git submodule foreach 'echo $name: $(du -sh .git)'

# 출력 형식:
📊 저장소 최적화 상태
━━━━━━━━━━━━━━━━━━
모드: [SLIM/FULL]
Partial Clone: [활성/비활성] (필터: blob:limit=1m)
Sparse Checkout: [활성/비활성] (N개 경로)
Shallow: [활성/비활성] (depth: N)
.git 폴더: XX MB
프로젝트 전체: XX MB
서브모듈: N개 (최적화: N개)
```

### 04. to-slim (`src/cmd/optimized/quick/04_to_slim.go`)
**목적**: FULL → SLIM 모드 전환
**구현 내용**:
```bash
# 실행 순서:
1. 현재 상태 백업
   git config --local --list > .git-config-backup

2. Partial Clone 필터 적용
   git config remote.origin.partialclonefilter blob:limit=1m
   git config remote.origin.promisor true
   git config extensions.partialClone origin

3. Sparse Checkout 활성화
   git sparse-checkout init --cone
   git sparse-checkout set src/ Assets/Scripts/ Assets/Shaders/ ProjectSettings/

4. 불필요한 객체 정리
   git repack -a -d
   git maintenance run --task=gc

5. 결과 확인
   du -sh .git
```

### 05. to-full (`src/cmd/optimized/quick/05_to_full.go`)
**목적**: SLIM → FULL 모드 복원
**구현 내용**:
```bash
# 실행 순서:
1. Sparse Checkout 해제
   git sparse-checkout disable

2. 모든 객체 다운로드
   git fetch --unshallow 2>/dev/null || true
   git fetch --refetch

3. Partial Clone 필터 제거
   git config --unset remote.origin.partialclonefilter
   git config --unset remote.origin.promisor
   git config --unset extensions.partialClone

4. 결과 확인
   du -sh .git
```

### 06. expand-slim (`src/cmd/optimized/quick/06_expand_slim.go`)
**목적**: SLIM 상태에서 선택적 경로 확장
**구현 내용**:
```bash
# 사용자 입력 받기: 확장할 경로

1. 현재 Sparse Checkout 목록 확인
   git sparse-checkout list

2. 경로 추가
   git sparse-checkout add <경로>

3. 필요한 객체 다운로드
   git read-tree -m -u HEAD

4. 확장 결과 표시
```

### 07. expand-filter (`src/cmd/optimized/quick/07_expand_filter.go`)
**목적**: Partial Clone 필터 제거 (Sparse는 유지)
**구현 내용**:
```bash
1. 현재 필터 확인
   git config remote.origin.partialclonefilter

2. 필터 제거
   git config --unset remote.origin.partialclonefilter

3. 모든 blob 다운로드
   git fetch --refetch

4. 결과 확인
```

### 08. expand-10 (`src/cmd/optimized/quick/08_expand_10.go`)
**목적**: 히스토리 10개 커밋 확장
**구현 내용**:
```bash
1. 현재 depth 확인
2. git fetch --deepen=10
3. 확장된 히스토리 확인
```

### 09. expand-50 (`src/cmd/optimized/quick/09_expand_50.go`)
**목적**: 히스토리 50개 커밋 확장
**구현 내용**:
```bash
1. 현재 depth 확인
2. git fetch --deepen=50
3. 확장된 히스토리 확인
```

### 10. expand-100 (`src/cmd/optimized/quick/10_expand_100.go`)
**목적**: 히스토리 100개 커밋 확장
**구현 내용**:
```bash
1. 현재 depth 확인
2. git fetch --deepen=100
3. 확장된 히스토리 확인
```

### 11. find-merge (`src/cmd/optimized/quick/11_find_merge.go`)
**목적**: 두 브랜치의 머지베이스 찾기
**구현 내용**:
```bash
# 사용자 입력: branch1, branch2

1. 머지베이스 찾기 시도
   git merge-base <branch1> <branch2>

2. 실패시 점진적 확장
   - git fetch --deepen=10
   - 재시도
   - 필요시 계속 확장

3. 결과 표시
   - 머지베이스 커밋 해시
   - 필요했던 depth
```

### 12. check-merge (`src/cmd/optimized/quick/12_check_merge.go`)
**목적**: 브랜치 병합 가능 여부 확인
**구현 내용**:
```bash
# 사용자 입력: target-branch

1. 머지베이스 확인
2. 병합 시뮬레이션
   git merge --no-commit --no-ff <branch>
3. 충돌 여부 확인
4. 결과 표시
   git merge --abort
```

### 13. clone-slim (`src/cmd/optimized/setup/13_clone_slim.go`)
**목적**: 처음부터 최적화된 클론
**구현 내용**:
```bash
# 사용자 입력: URL, 폴더명

1. Partial Clone으로 클론
   git clone --filter=blob:limit=1m --sparse <url> <folder>

2. Sparse Checkout 설정
   cd <folder>
   git sparse-checkout init --cone
   git sparse-checkout set src/ Assets/Scripts/ Assets/Shaders/ ProjectSettings/

3. Shallow 설정
   git pull --depth=1

4. 서브모듈 초기화
   git submodule update --init --filter=blob:limit=50k --depth=1

5. 성능 설정 적용
   git config core.commitGraph true
   git config gc.writeCommitGraph true
```

### 14. migrate (`src/cmd/optimized/setup/14_migrate.go`)
**목적**: 기존 저장소를 SLIM으로 변환
**구현 내용**:
```bash
1. 작업 중인 변경사항 확인
   git status

2. stash가 필요한 경우 저장
   git stash

3. to-slim 프로세스 실행 (04번 참조)

4. stash 복원
   git stash pop
```

### 15. performance (`src/cmd/optimized/setup/15_performance.go`)
**목적**: 성능 최적화 설정 적용
**구현 내용**:
```bash
1. Git 성능 설정
   git config core.commitGraph true
   git config gc.writeCommitGraph true
   git config fetch.writeCommitGraph true
   git config core.multiPackIndex true
   git config fetch.parallel 4
   git config gc.autoDetach false

2. Maintenance 스케줄 등록
   git maintenance register

3. 초기 maintenance 실행
   git maintenance run

4. 설정 확인 표시
```

### 16. expand-path (`src/cmd/optimized/workspace/16_expand_path.go`)
**목적**: 특정 경로를 Sparse Checkout에 추가
**구현 내용**:
```bash
# 사용자 입력: 경로

1. 경로 유효성 확인
2. Sparse Checkout에 추가
   git sparse-checkout add <경로>
3. 파일 다운로드
4. 결과 표시
```

### 17. filter-branch (`src/cmd/optimized/workspace/17_filter_branch.go`)
**목적**: 브랜치별로 다른 Partial Clone 필터 적용
**구현 내용**:
```bash
# 사용자 입력: 브랜치명, 필터 크기

1. 브랜치 전환
   git checkout <브랜치>

2. 브랜치별 필터 설정
   git config branch.<브랜치>.partialclonefilter blob:limit=<크기>

3. 필터 적용
   git fetch --refetch

4. 설정 확인
```

### 18. clear-filter (`src/cmd/optimized/workspace/18_clear_filter.go`)
**목적**: 모든 필터 제거
**구현 내용**:
```bash
1. Partial Clone 필터 제거
2. Sparse Checkout 해제
3. 모든 객체 다운로드
4. 결과 확인
```

### 19. restore-branch (`src/cmd/optimized/workspace/19_restore_branch.go`)
**목적**: 특정 브랜치만 전체 복원
**구현 내용**:
```bash
# 사용자 입력: 브랜치명

1. 브랜치 전환
2. 해당 브랜치의 모든 파일 다운로드
3. 히스토리 복원
4. 결과 확인
```

### 20. shallow (`src/cmd/optimized/advanced/20_shallow.go`)
**목적**: 히스토리를 depth=1로 줄이기
**구현 내용**:
```bash
1. 현재 상태 백업
2. git pull --depth=1
3. git gc --prune=now
4. 결과 확인
```

### 21. unshallow (`src/cmd/optimized/advanced/21_unshallow.go`)
**목적**: 전체 히스토리 복원
**구현 내용**:
```bash
1. git fetch --unshallow
2. 결과 확인
```

### 22. check-shallow (`src/cmd/optimized/advanced/22_check_shallow.go`)
**목적**: 현재 shallow 상태 확인
**구현 내용**:
```bash
1. Shallow 여부 확인
   git rev-parse --is-shallow-repository

2. Shallow인 경우 depth 확인
   git rev-list --count HEAD

3. Grafted 커밋 확인
   cat .git/shallow

4. 결과 표시
```

### 23. check-filter (`src/cmd/optimized/advanced/23_check_filter.go`)
**목적**: 현재 필터 설정 확인
**구현 내용**:
```bash
1. Global 필터 확인
   git config remote.origin.partialclonefilter

2. 브랜치별 필터 확인
   git config --get-regexp branch.*.partialclonefilter

3. 결과 표시
```

### 24. backup-config (`src/cmd/optimized/advanced/24_backup_config.go`)
**목적**: 최적화 설정 백업/복원
**구현 내용**:
```bash
# 사용자 선택: backup 또는 restore

[Backup]
1. Git 설정 백업
   git config --local --list > .git-optimization-backup

2. Sparse Checkout 목록 백업
   git sparse-checkout list > .git-sparse-backup

[Restore]
1. 백업 파일 읽기
2. 설정 복원
3. Sparse Checkout 복원
```

### 25. shallow-all (`src/cmd/optimized/submodule/25_shallow_all.go`)
**목적**: 모든 서브모듈 shallow 변환
**구현 내용**:
```bash
1. git submodule foreach 'git pull --depth=1'
2. 각 서브모듈 결과 표시
```

### 26. unshallow-all (`src/cmd/optimized/submodule/26_unshallow_all.go`)
**목적**: 모든 서브모듈 히스토리 복원
**구현 내용**:
```bash
1. git submodule foreach 'git fetch --unshallow'
2. 각 서브모듈 결과 표시
```

### 27. optimize-all (`src/cmd/optimized/submodule/27_optimize_all.go`)
**목적**: 모든 서브모듈 SLIM 최적화
**구현 내용**:
```bash
1. 각 서브모듈에 대해:
   - Partial Clone 필터 적용
   - Sparse Checkout 설정
   - Shallow 설정
   
2. git submodule foreach 실행
3. 결과 통계 표시
```

### 28. status-all (`src/cmd/optimized/submodule/28_status_all.go`)
**목적**: 모든 서브모듈 상태 확인
**구현 내용**:
```bash
1. git submodule foreach 실행:
   - Partial Clone 상태
   - Sparse Checkout 상태
   - Shallow 상태
   - 디스크 사용량

2. 테이블 형식으로 출력:
   서브모듈명 | 모드 | 필터 | Sparse | Shallow | 크기
```

---

## 🔧 공통 유틸리티 함수

### 에러 처리
```go
func handleError(err error, msg string) {
    if err != nil {
        fmt.Printf("❌ 오류: %s\n", msg)
        fmt.Printf("   상세: %v\n", err)
        os.Exit(1)
    }
}
```

### Git 명령 실행
```go
func runGitCommand(args ...string) (string, error) {
    cmd := exec.Command("git", args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### 진행 상황 표시
```go
func showProgress(current, total int, message string) {
    percentage := (current * 100) / total
    fmt.Printf("\r[%d%%] %s", percentage, message)
}
```

### 디스크 사용량 확인
```go
func getDiskUsage(path string) string {
    cmd := exec.Command("du", "-sh", path)
    output, _ := cmd.Output()
    return strings.TrimSpace(string(output))
}
```

---

## 📝 구현 시 주의사항

1. **에러 처리**: 모든 Git 명령어 실행 시 에러 체크 필수
2. **사용자 확인**: 위험한 작업(전체 복원, 필터 제거) 시 확인 프롬프트
3. **진행 표시**: 시간이 오래 걸리는 작업은 진행 상황 표시
4. **백업**: 설정 변경 전 현재 상태 백업
5. **서브모듈**: 서브모듈 작업 시 재귀적 처리
6. **성능**: 대용량 저장소 처리 시 메모리 효율성 고려

## 🧪 테스트 방법

각 함수 구현 후:
1. 테스트 저장소에서 실행
2. 예상 결과와 실제 결과 비교
3. 에러 케이스 테스트
4. 실제 프로젝트에 적용 전 백업

---

## 📊 완료 기준

- [ ] 모든 함수가 PRD 명세대로 구현됨
- [ ] 에러 처리가 완벽함
- [ ] 사용자 친화적인 출력
- [ ] 테스트 완료
- [ ] 문서화 완료