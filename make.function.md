# Git 저장소 최적화 명령어 구현 가이드

## 🎯 핵심 구현 전략 - AI 협업 워크플로우

### ⚠️ 필수 규칙: Zen MCP를 통한 구현 위임 ⚠️
**핵심 전략**: Claude의 컨텍스트를 절약하여 더 많은 작업을 수행할 수 있도록 합니다!
- ✅ **Claude 역할**: 파일 읽기, 분석, 검증, 테스트 (컨텍스트 소량 사용)
- ✅ **Gemini 역할**: 실제 코드 구현 작업 (Zen MCP를 통해 위임)
- 📌 **효과**: 한 세션에서 더 많은 명령어를 구현 가능

### 구현 프로세스 (3단계 사이클)

#### 1단계: Gemini 구현 (gemini-2.5-pro) - Zen MCP 활용
- **Claude 준비 작업**:
  - 필요한 파일들 확인 (Read 도구 사용 가능)
  - 구현 명세 파악
  - Zen MCP에 전달할 프롬프트 준비
- **Gemini 구현 (Zen MCP)**:
  - `mcp__zen__chat` 또는 `mcp__zen__thinkdeep` 사용
  - 입력: make.function.md 명세, 타겟 파일, 유틸리티 파일들
  - 출력: 완전히 구현된 코드
- **명령 예시**: 
  ```
  "mcp__zen__chat으로 make.function.md의 05번 to-full 명세에 따라 
   src/cmd/optimized/quick/05_to_full.go 구현해줘. 
   utils/git.go와 utils/utils.go 유틸리티 활용"
  ```

#### 2단계: Claude 검증 및 개선 (claude-opus-4.1)
- **검증 작업**:
  - 구현 코드 검토 및 목적 부합성 확인
  - 문법 및 로직 오류 검사
  - 컴파일 테스트 (`go build -o ga`)
  - 실행 테스트 및 동작 검증
  - 품질 승인 판단
- **개선 작업**:
  - 발견된 버그 수정
  - 누락된 기능 추가 구현
  - 코드 최적화 및 리팩토링
  - 에러 처리 강화
  - 사용자 경험 개선
- **산출물**: 검증 완료 및 개선된 최종 코드

#### 3단계: 사용자 최종 확인
- **검토 항목**:
  - Claude 검증 및 개선 결과 확인
  - 최종 구현 품질 승인
  - 커밋 지시
- **완료**: 커밋 및 문서 업데이트

### 역할 분담
- **Gemini (구현자)**: 
  - 초기 코드 작성
  - 핵심 로직 구현
  - 유틸리티 함수 활용
  
- **Claude (검증자 & 개선자)**: 
  - 코드 리뷰 및 품질 검증
  - 버그 수정 및 예외 처리
  - 누락 기능 추가 구현
  - 코드 개선 및 최적화
  - 테스트 실행 및 동작 확인
  - 최종 품질 보증
  
- **사용자 (승인자)**: 
  - 최종 판단 및 승인
  - 커밋 결정
  - 다음 작업 지시

### 협업 규칙
1. 각 명령어는 반드시 이 3단계 프로세스를 따름
2. Gemini는 항상 전체 컨텍스트(유틸리티 포함)를 받아 구현
3. Claude는 구현된 코드를 실제 환경에서 테스트하고 필요시 개선
4. Claude가 품질 승인한 코드만 사용자에게 제출
5. 사용자 승인 없이는 커밋하지 않음

### 품질 기준
- ✅ 컴파일 오류 없음
- ✅ 런타임 오류 없음
- ✅ 명세서 요구사항 100% 충족
- ✅ 적절한 에러 처리
- ✅ 사용자 친화적 메시지
- ✅ 코드 가독성 및 유지보수성

---

## 🚀 실행 방법
```bash
ga optimized {카테고리} {명령어}
ga opt {카테고리} {명령어}       # 짧은 별칭
ga op {카테고리} {명령어}        # 더 짧은 별칭

# 예시:
ga optimized help workflow      # 워크플로우 가이드 표시
ga opt help commands           # 전체 명령어 목록 표시  
ga op quick status             # 최적화 상태 확인
ga opt quick to-slim           # SLIM 모드로 전환
ga optimized quick to-full     # FULL 모드로 복원
```

## 📋 개요
이 문서는 Git 저장소 최적화를 위한 28개 명령어의 구현 상세를 담고 있습니다.
각 명령어는 PRD 기반으로 구체적인 구현 방법이 정의되어 있습니다.

## 🎯 구현 진행 상황 (12/28)
- [x] 01. workflow - Git 최적화 워크플로우 가이드
- [x] 02. commands - 전체 명령어 목록
- [x] 03. status - 현재 최적화 상태 확인
- [x] 04. to-slim - SLIM 모드로 전환
- [x] 05. to-full - FULL 모드로 복원
- [x] 06. expand-slim - 선택적 경로 확장
- [x] 07. expand-filter - Partial Clone 필터 제거
- [x] 08. expand (통합) - 히스토리 확장 (파라미터로 개수 지정)
- [x] 09. expand-50 - (deprecated - expand 50 사용)
- [x] 10. expand-100 - (deprecated - expand 100 사용)
- [x] 11. find-merge - 병합 베이스 찾기
- [ ] 12. check-merge - 병합 가능 여부 확인
- [ ] 13. clone-slim - 최적화된 클론
- [x] 14. migrate - (deprecated - to-slim 사용)
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

## 🔧 유틸리티 전략 및 지침

### 패키지 구조
```
src/
├── utils/           # 범용 유틸리티 패키지
│   ├── utils.go    # 기본 유틸리티 (UI, 파일 처리 등)
│   └── git.go      # Git 관련 유틸리티 함수
├── cmd/
│   ├── utils.go    # cmd 패키지 브릿지 (utils 패키지 재사용)
│   └── optimized/
│       └── quick/  # 최적화 명령어 구현
```

### 유틸리티 구성 방침
1. **Git 관련 함수** (`src/utils/git.go`)
   - Git 저장소 상태 확인 (IsGitRepository, GetOptimizationMode)
   - Git 설정 조회 (GetPartialCloneFilter, IsSparseCheckoutEnabled)
   - Git 정보 수집 (GetObjectInfo, GetSubmoduleInfo, GetDiskUsage)
   - 파일 분석 (GetExcludedLargeFiles, GetLargestFilesInHistory)
   - 포맷팅 헬퍼 (FormatSize, TruncateString)

2. **일반 유틸리티** (`src/utils/utils.go`)
   - 사용자 입력 처리 (Confirm, ConfirmWithDefault)
   - Git 경로 처리 (UnescapeGitPath, ProcessGitPaths, DecodeGitPath)
   - AI용 Diff 생성 (GetDiffForAI)
   - 파일 유형 판단 (IsSourceCodeFile)
   - 크기 변환 (HumanizeBytes)

3. **브릿지 파일** (`src/cmd/utils.go`)
   - utils 패키지의 필요한 함수들을 cmd 패키지로 재노출
   - 패키지 경계를 깔끔하게 유지

### 사용 가이드라인
- 새로운 명령어 구현 시 기존 유틸리티 재사용 우선
- Git 작업은 반드시 `utils/git.go`의 함수 활용
- 중복 코드 발견 시 즉시 유틸리티로 추출
- 유틸리티 함수는 단일 책임 원칙 준수
- 에러 처리는 호출하는 쪽에서 수행

### Import 경로 규칙
- **유틸리티 import**: `"workingcli/src/utils"`
- 절대 경로가 아닌 모듈 경로 사용
- 예시:
  ```go
  import (
      "fmt"
      "os"
      
      "workingcli/src/utils"  // ✅ 올바른 import
      "github.com/spf13/cobra"
  )
  ```

## 📚 함수별 구현 상세

### 01. workflow (`src/cmd/optimized/help/01_workflow.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26)
**목적**: 28개 전체 명령어 목록 표시
**구현 내용**: 
```
1. 카테고리별 명령어 그룹화
   - Help: 도움말 명령어 (workflow, commands)
   - Quick: 자주 사용하는 빠른 명령어 (status, to-slim, to-full 등)
   - Setup: 초기 설정 관련 (clone-slim, migrate, performance)
   - Workspace: 작업 공간 관리 (expand-path, filter-branch 등)
   - Advanced: 고급 기능 (shallow, unshallow, check-shallow 등)
   - Submodule: 서브모듈 관리 (shallow-all, optimize-all 등)

2. 각 명령어별 간단한 설명 포함
3. 사용 예시 제공
```

### 03. status (`src/cmd/optimized/quick/03_status.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26)
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

### 08. expand (통합 명령어) (`src/cmd/optimized/quick/08_expand_10.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
**목적**: 히스토리 확장 (파라미터로 개수 지정)
**구현 내용**:
```bash
# 사용법: ga opt quick expand [depth]
# depth를 지정하지 않으면 기본값 10개

1. 현재 depth 확인
2. git fetch --deepen=N (N은 파라미터로 받은 값)
3. 확장된 히스토리 확인

# 사용 예시:
ga opt quick expand        # 10개 확장 (기본값)
ga opt quick expand 10     # 10개 확장
ga opt quick expand 50     # 50개 확장  
ga opt quick expand 100    # 100개 확장
ga opt quick expand 66     # 66개 확장 (커스텀)
```

### 09. expand-50 (`src/cmd/optimized/quick/09_expand_50.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - deprecated
**목적**: 히스토리 50개 커밋 확장 (deprecated - expand 50 사용)
**구현 내용**:
```bash
# deprecated - 대신 사용:
ga opt quick expand 50
```

### 10. expand-100 (`src/cmd/optimized/quick/10_expand_100.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - deprecated
**목적**: 히스토리 100개 커밋 확장 (deprecated - expand 100 사용)
**구현 내용**:
```bash
# deprecated - 대신 사용:
ga opt quick expand 100
```

### 11. find-merge (`src/cmd/optimized/quick/11_find_merge.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
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
**상태**: ✅ 구현 완료 (2025-08-26) - deprecated
**목적**: 기존 저장소를 SLIM으로 변환 (deprecated - to-slim 사용)
**구현 내용**:
```bash
# deprecated - 대신 사용:
ga opt quick to-slim

# migrate와 to-slim은 동일한 기능:
- 기존 FULL 상태 저장소를 SLIM으로 변환
- 작업 내용 보존하면서 최적화 적용
- to-slim이 더 직관적이고 quick 카테고리에 있어 접근성 좋음
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

---

## 🔖 커밋 메시지 규약

### 기본 형식
```
<type>(<scope>): <subject>

[optional body]
[optional footer]
```

### 타입 정의
- `feat(opt)`: 새로운 최적화 기능 구현
- `test(opt)`: 최적화 기능 테스트 추가
- `docs(opt)`: 최적화 기능 문서화
- `fix(opt)`: 최적화 기능 버그 수정
- `refactor(opt)`: 최적화 기능 코드 리팩토링
- `perf(opt)`: 최적화 기능 성능 개선
- `chore(opt)`: 빌드, 설정 등 기타 변경

### 커밋 메시지 작성 규칙
1. **제목은 50자 이내**로 작성
2. **명령문 형태**로 작성 (implement, add, fix, update)
3. **함수 번호와 이름**을 명시
4. **구체적인 변경 내용** 포함
5. **본문은 한글로 작성** (제목과 기술적 용어 제외)

### 단계별 커밋 예시
```bash
# 구현 단계
feat(opt): implement 03-status - Git repository optimization status check

# 테스트 단계
test(opt): add tests for 03-status command functionality

# 문서화 단계
docs(opt): update documentation for 03-status command usage
```

---

## 📝 함수별 커밋 메시지 템플릿

### Help 카테고리 (도움말)
```bash
# 01. workflow
feat(opt): implement 01-workflow - Git optimization workflow guide
test(opt): add tests for 01-workflow help command
docs(opt): document 01-workflow usage and examples

# 02. commands  
feat(opt): implement 02-commands - complete command list display
test(opt): add tests for 02-commands help display
docs(opt): document 02-commands help system
```

### Quick 카테고리 (빠른 실행)
```bash
# 03. status
feat(opt): implement 03-status - repository optimization status check
test(opt): add tests for 03-status metrics collection
docs(opt): document 03-status output format and usage

# 04. to-slim
feat(opt): implement 04-to-slim - convert repository to SLIM mode
test(opt): add tests for 04-to-slim conversion process
docs(opt): document 04-to-slim migration workflow

# 05. to-full
feat(opt): implement 05-to-full - restore repository to FULL mode
test(opt): add tests for 05-to-full restoration process
docs(opt): document 05-to-full recovery workflow

# 06. expand-slim
feat(opt): implement 06-expand-slim - selective path expansion
test(opt): add tests for 06-expand-slim path addition
docs(opt): document 06-expand-slim usage scenarios

# 07. expand-filter
feat(opt): implement 07-expand-filter - remove Partial Clone filter
test(opt): add tests for 07-expand-filter removal process
docs(opt): document 07-expand-filter filter management

# 08. expand-10
feat(opt): implement 08-expand-10 - extend history by 10 commits
test(opt): add tests for 08-expand-10 depth expansion
docs(opt): document 08-expand-10 history extension

# 09. expand-50
feat(opt): implement 09-expand-50 - extend history by 50 commits
test(opt): add tests for 09-expand-50 depth expansion
docs(opt): document 09-expand-50 history extension

# 10. expand-100
feat(opt): implement 10-expand-100 - extend history by 100 commits
test(opt): add tests for 10-expand-100 depth expansion
docs(opt): document 10-expand-100 history extension

# 11. find-merge
feat(opt): implement 11-find-merge - locate merge base between branches
test(opt): add tests for 11-find-merge base detection
docs(opt): document 11-find-merge merge analysis

# 12. check-merge
feat(opt): implement 12-check-merge - verify merge compatibility
test(opt): add tests for 12-check-merge conflict detection
docs(opt): document 12-check-merge merge verification
```

### Setup 카테고리 (초기 설정)
```bash
# 13. clone-slim
feat(opt): implement 13-clone-slim - optimized repository cloning
test(opt): add tests for 13-clone-slim initialization
docs(opt): document 13-clone-slim setup process

# 14. migrate
feat(opt): implement 14-migrate - convert existing repo to SLIM
test(opt): add tests for 14-migrate conversion workflow
docs(opt): document 14-migrate migration strategy

# 15. performance
feat(opt): implement 15-performance - apply performance settings
test(opt): add tests for 15-performance configuration
docs(opt): document 15-performance optimization settings
```

### Workspace 카테고리 (작업공간)
```bash
# 16. expand-path
feat(opt): implement 16-expand-path - add specific paths to sparse
test(opt): add tests for 16-expand-path path management
docs(opt): document 16-expand-path selective expansion

# 17. filter-branch
feat(opt): implement 17-filter-branch - branch-specific filters
test(opt): add tests for 17-filter-branch filter application
docs(opt): document 17-filter-branch branch filtering

# 18. clear-filter
feat(opt): implement 18-clear-filter - remove all filters
test(opt): add tests for 18-clear-filter cleanup process
docs(opt): document 18-clear-filter filter removal

# 19. restore-branch
feat(opt): implement 19-restore-branch - restore specific branch
test(opt): add tests for 19-restore-branch restoration
docs(opt): document 19-restore-branch branch recovery
```

### Advanced 카테고리 (고급)
```bash
# 20. shallow
feat(opt): implement 20-shallow - reduce history to depth 1
test(opt): add tests for 20-shallow history reduction
docs(opt): document 20-shallow shallow conversion

# 21. unshallow
feat(opt): implement 21-unshallow - restore complete history
test(opt): add tests for 21-unshallow history restoration
docs(opt): document 21-unshallow full recovery

# 22. check-shallow
feat(opt): implement 22-check-shallow - verify shallow status
test(opt): add tests for 22-check-shallow status detection
docs(opt): document 22-check-shallow status checking

# 23. check-filter
feat(opt): implement 23-check-filter - inspect filter settings
test(opt): add tests for 23-check-filter configuration check
docs(opt): document 23-check-filter filter inspection

# 24. backup-config
feat(opt): implement 24-backup-config - backup/restore settings
test(opt): add tests for 24-backup-config save/load
docs(opt): document 24-backup-config configuration management
```

### Submodule 카테고리 (서브모듈)
```bash
# 25. shallow-all
feat(opt): implement 25-shallow-all - shallow all submodules
test(opt): add tests for 25-shallow-all batch processing
docs(opt): document 25-shallow-all submodule optimization

# 26. unshallow-all
feat(opt): implement 26-unshallow-all - restore all submodules
test(opt): add tests for 26-unshallow-all batch restoration
docs(opt): document 26-unshallow-all submodule recovery

# 27. optimize-all
feat(opt): implement 27-optimize-all - optimize all submodules
test(opt): add tests for 27-optimize-all batch optimization
docs(opt): document 27-optimize-all comprehensive optimization

# 28. status-all
feat(opt): implement 28-status-all - check all submodule status
test(opt): add tests for 28-status-all status collection
docs(opt): document 28-status-all status reporting
```

---

## 🌿 브랜치 전략

### 브랜치 네이밍 규칙
```bash
# 기능 구현 브랜치
feature/opt-<번호>-<함수명>

# 예시:
feature/opt-03-status
feature/opt-04-to-slim
feature/opt-13-clone-slim
```

### 브랜치 생성 및 작업 순서
```bash
# 1. 브랜치 생성
git checkout -b feature/opt-03-status

# 2. 구현 작업
# ... 코드 작성 ...

# 3. 커밋 (구현)
git add src/cmd/optimized/quick/03_status.go
git commit -m "feat(opt): implement 03-status - repository optimization status check"

# 4. 테스트 추가
# ... 테스트 작성 ...

# 5. 커밋 (테스트)
git add src/cmd/optimized/quick/03_status_test.go
git commit -m "test(opt): add tests for 03-status metrics collection"

# 6. 문서화
# ... 문서 작성 ...

# 7. 커밋 (문서)
git add docs/optimized/03-status.md
git commit -m "docs(opt): document 03-status output format and usage"

# 8. Push
git push origin feature/opt-03-status
```

### 릴리스 태그

#### 태그 생성 시점
카테고리별 모든 명령어가 완료되면 즉시 태그를 생성합니다:

- **Help 카테고리** (1-2번): 2개 모두 완료 시 → `v1.0.0-opt-help`
- **Quick 카테고리** (3-12번): 10개 모두 완료 시 → `v1.1.0-opt-quick`
- **Setup 카테고리** (13-15번): 3개 모두 완료 시 → `v1.2.0-opt-setup`
- **Workspace 카테고리** (16-19번): 4개 모두 완료 시 → `v1.3.0-opt-workspace`
- **Advanced 카테고리** (20-24번): 5개 모두 완료 시 → `v1.4.0-opt-advanced`
- **Submodule 카테고리** (25-28번): 4개 모두 완료 시 → `v1.5.0-opt-submodule`
- **전체 완료**: 28개 모두 완료 시 → `v2.0.0-opt-complete`

#### 태그 생성 명령어
```bash
# 카테고리 완료 시 태그 생성
git tag -a v1.0.0-opt-help -m "Complete Help category implementation (1-2)"
git push origin v1.0.0-opt-help

# 마일스톤별 태그 목록
v1.0.0-opt-help      # Help 카테고리 완료 (1-2번)
v1.1.0-opt-quick     # Quick 카테고리 완료 (3-12번)
v1.2.0-opt-setup     # Setup 카테고리 완료 (13-15번)
v1.3.0-opt-workspace # Workspace 카테고리 완료 (16-19번)
v1.4.0-opt-advanced  # Advanced 카테고리 완료 (20-24번)
v1.5.0-opt-submodule # Submodule 카테고리 완료 (25-28번)
v2.0.0-opt-complete  # 전체 최적화 기능 완료 (1-28번)
```

---

## ⚠️ 중요: 체크리스트 업데이트

### 구현 완료 시 체크리스트 업데이트 필수
각 함수 구현이 완료되면 반드시 다음 항목들을 업데이트해야 합니다:

1. **진행 상황 업데이트** (상단 제목)
   ```markdown
   ## 🎯 구현 진행 상황 (1/28)  # 숫자 업데이트
   ```

2. **체크박스 업데이트** 
   ```markdown
   - [x] 01. workflow - Git 최적화 워크플로우 가이드  # 완료된 항목 체크
   ```

3. **구현 상태 표시** (함수별 구현 상세 섹션)
   ```markdown
   ### 01. workflow (`src/cmd/optimized/help/01_workflow.go`)
   **상태**: ✅ 구현 완료 (2025-01-XX)  # 날짜 추가
   ```

### 업데이트 예시
```bash
# 구현 전
- [ ] 01. workflow - Git 최적화 워크플로우 가이드

# 구현 후  
- [x] 01. workflow - Git 최적화 워크플로우 가이드
```

### 체크리스트 업데이트 커밋
```bash
docs(opt): update checklist for 01-workflow completion

- Mark 01-workflow as completed
- Update progress count (1/28)
- Add completion date
```

### make.function.md 상태 업데이트 커밋 규약

**커밋 메시지 작성 시 주의사항:**
- 제목은 영문으로 작성 (GitHub 호환성)
- 본문은 한글로 작성하여 명확한 의미 전달
- 진행 상황 숫자는 정확히 업데이트

```bash
# 함수 구현 완료 시 체크리스트 업데이트
docs(opt): update checklist for <번호>-<함수명> completion

- <번호>-<함수명> 완료 표시
- 진행 상황 업데이트 (<현재/28>)
- 완료 날짜 추가

# 예시:
docs(opt): update checklist for 03-status completion

- 03-status 완료 표시
- 진행 상황 업데이트 (3/28)
- 완료 날짜 추가 (2025-08-26)

# 여러 함수 동시 완료 시
docs(opt): update checklist for multiple completions

- 03-status, 04-to-slim 완료 표시
- 진행 상황 업데이트 (4/28)
- 완료 날짜들 추가

# 카테고리 완료 시
docs(opt): complete Help category implementation

- Help 카테고리 전체 명령어 완료 표시 (1-2번)
- 진행 상황 업데이트 (2/28)
- 마일스톤 달성 기록

# 부분 구현 또는 진행 중 상태 업데이트
docs(opt): update 03-status implementation progress

- 부분 구현 내용 추가
- 남은 작업 TODO 업데이트
- 블로커나 의존성 문서화
```