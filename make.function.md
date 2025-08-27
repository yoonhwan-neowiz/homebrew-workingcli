# Git 저장소 최적화 명령어 구현 가이드

## 📊 구현 진행률: 26/33 (79%)

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
   src/cmd/optimized/quick/to_full.go 구현해줘. 
   utils/git.go와 utils/utils.go 유틸리티 활용. 
   src/cmd/optimized/quick/to_slim.go, src/config/config.go 구현 참고 및 이용."
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

# 서브모듈 명령어 (submodule 카테고리):
ga opt submodule status        # 서브모듈 상태 확인
ga op submodule to-slim        # 서브모듈 SLIM 전환
ga opt submodule to-full       # 서브모듈 FULL 복원
ga op submodule expand-slim    # 서브모듈 경로 확장
ga opt submodule filter-branch # 서브모듈 브랜치 필터
```

## 📁 파일명 규칙
- **모든 명령어**: `{명령어}.go` (번호 없이)
  - help: `workflow.go`, `commands.go`
  - quick: `status.go`, `to_slim.go`
  - submodule: `status.go`, `to_slim.go`
  
- **실제 명령어 사용**: 카테고리별로 동일한 명령어명 사용 가능
  - 메인: `ga opt quick status`
  - 서브모듈: `ga opt submodule status` (카테고리로 구분)

## 📋 개요
이 문서는 Git 저장소 최적화를 위한 33개 명령어의 구현 상세를 담고 있습니다.
각 명령어는 PRD 기반으로 구체적인 구현 방법이 정의되어 있습니다.

## 🎯 구현 진행 상황 (27/33)
- [x] help.workflow - Git 최적화 워크플로우 가이드
- [x] help.commands - 전체 명령어 목록
- [x] quick.status - 현재 최적화 상태 확인
- [x] quick.to-slim - SLIM 모드로 전환
- [x] quick.to-full - FULL 모드로 복원
- [x] quick.expand-slim - 선택적 경로 확장
- [x] quick.expand-filter - Partial Clone 필터 제거
- [x] advanced.expand - 히스토리 확장 (파라미터로 개수 지정)
- [x] advanced.expand-50 - (deprecated - expand 50 사용)
- [x] advanced.expand-100 - (deprecated - expand 100 사용)
- [x] quick.auto-find-merge-base - 브랜치 병합점 자동 찾기
- [x] advanced.check-merge - 병합 가능 여부 확인
- [x] setup.clone-slim - 최적화된 클론
- [x] setup.migrate - (deprecated - to-slim 사용)
- [x] setup.performance - 성능 최적화 설정
- [x] workspace.expand-path - 특정 경로 확장
- [x] quick.filter-branch - 브랜치 필터 설정 (특정 브랜치만 표시)
- [x] quick.clear-filter-branch - 브랜치 필터 제거 (모든 브랜치 표시)
- [x] workspace.restore-branch - (DEPRECATED - 사용하지 않음)
- [x] quick.shallow - 히스토리 줄이기
- [x] quick.unshallow - 히스토리 복원
- [x] advanced.check-shallow - 히스토리 상태 확인
- [x] advanced.check-filter - 브랜치 필터 확인
- [x] advanced.config - 설정 백업/복원/확인
- [x] submodule.status - 서브모듈별 최적화 상태 확인
- [ ] submodule.to-slim - 서브모듈을 SLIM 모드로 전환
- [ ] submodule.to-full - 서브모듈을 FULL 모드로 복원
- [ ] submodule.expand-slim - 서브모듈 선택적 경로 확장
- [ ] submodule.expand-filter - 서브모듈 Partial Clone 필터 제거
- [x] submodule.shallow - 서브모듈 shallow 변환 (recursive)
- [x] submodule.unshallow - 서브모듈 히스토리 복원 (recursive)
- [ ] submodule.filter-branch - 서브모듈 브랜치 필터 설정
- [ ] submodule.clear-filter-branch - 서브모듈 브랜치 필터 제거

---

## 📂 구현 파일 현황

### 카테고리별 명령어 구성

#### Quick 카테고리 (자주 사용하는 최적화 기능)
| 파일명 | 명령어 | 설명 | 상태 |
|--------|--------|------|------|
| `status.go` | `status` | 현재 최적화 상태 확인 | ✅ 구현 완료 |
| `to_slim.go` | `to-slim` | SLIM 모드로 전환 | ✅ 구현 완료 |
| `to_full.go` | `to-full` | FULL 모드로 복원 | ✅ 구현 완료 |
| `expand_slim.go` | `expand-slim` | 선택적 경로 확장 | ✅ 구현 완료 |
| `expand_filter.go` | `expand-filter` | Partial Clone 필터 제거 | ✅ 구현 완료 |
| `auto_find_merge_base.go` | `auto-find-merge-base` | 브랜치 병합점 자동 찾기 | ✅ 구현 완료 |
| `filter_branch.go` | `filter-branch` | 브랜치 필터 설정 | ✅ 구현 완료 |
| `clear_filter_branch.go` | `clear-filter` | 브랜치 필터 제거 | ✅ 구현 완료 |
| `shallow.go` | `shallow` | 히스토리 줄이기 | ✅ 구현 완료 |
| `unshallow.go` | `unshallow` | 히스토리 복원 | ✅ 구현 완료 |

#### Advanced 카테고리 (고급 최적화 기능)
| 파일명 | 명령어 | 설명 | 상태 |
|--------|--------|------|------|
| `expand.go` | `expand` | 히스토리 확장 (기본 10개) | ✅ 구현 완료 |
| `expand_50.go` | `expand-50` | 히스토리 50개 확장 (deprecated) | ✅ 구현 완료 |
| `expand_100.go` | `expand-100` | 히스토리 100개 확장 (deprecated) | ✅ 구현 완료 |
| `check_merge.go` | `check-merge` | 병합 가능 여부 확인 | ✅ 구현 완료 |
| `check_shallow.go` | `check-shallow` | 히스토리 상태 확인 | ✅ 구현 완료 |
| `check_filter.go` | `check-filter` | 브랜치 필터 확인 | ✅ 구현 완료 |
| `config.go` | `config` | 설정 백업/복원/확인 | ✅ 구현 완료 |

#### Workspace 카테고리 (작업 공간 관리)
| 파일명 | 명령어 | 설명 | 상태 |
|--------|--------|------|------|
| `expand_path.go` | `expand-path` | 특정 경로 확장 | ✅ 구현 완료 |
| `restore_branch.go` | `restore-branch` | (DEPRECATED) | ✅ DEPRECATED |

#### Submodule 카테고리 (서브모듈 최적화)
| 파일명 | 명령어 | 설명 | 상태 |
|--------|--------|------|------|
| `status.go` | `status` | 서브모듈 상태 확인 | ✅ 구현 완료 |
| `to_slim.go` | `to-slim` | SLIM 모드 전환 | ⏳ 대기 |
| `to_full.go` | `to-full` | FULL 모드 복원 | ⏳ 대기 |
| `expand_slim.go` | `expand-slim` | 경로 확장 | ⏳ 대기 |
| `expand_filter.go` | `expand-filter` | 필터 제거 | ⏳ 대기 |
| `shallow.go` | `shallow` | shallow 변환 (recursive) | ✅ 구현 완료 |
| `unshallow.go` | `unshallow` | 히스토리 복원 (recursive) | ✅ 구현 완료 |
| `filter_branch.go` | `filter-branch` | 브랜치 필터 | ⏳ 대기 |
| `clear_filter.go` | `clear-filter` | 필터 제거 | ⏳ 대기 |

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

### help.workflow (`src/cmd/optimized/help/workflow.go`)
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

### help.commands (`src/cmd/optimized/help/commands.go`)
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

### quick.status (`src/cmd/optimized/quick/status.go`)
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

### quick.to-slim (`src/cmd/optimized/quick/to_slim.go`)
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

### quick.to-full (`src/cmd/optimized/quick/to_full.go`)
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

### quick.expand-slim (`src/cmd/optimized/quick/expand_slim.go`)
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

### quick.expand-filter (`src/cmd/optimized/quick/expand_filter.go`)
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

### advanced.expand (`src/cmd/optimized/advanced/expand.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - Advanced로 이동
**목적**: 히스토리 확장 (파라미터로 개수 지정)
**구현 내용**:
```bash
# 사용법: ga opt advanced expand [depth]
# depth를 지정하지 않으면 기본값 10개

1. 현재 depth 확인
2. git fetch --deepen=N (N은 파라미터로 받은 값)
3. 확장된 히스토리 확인

# 사용 예시:
ga opt advanced expand        # 10개 확장 (기본값)
ga opt advanced expand 10     # 10개 확장
ga opt advanced expand 50     # 50개 확장  
ga opt advanced expand 100    # 100개 확장
ga opt advanced expand 66     # 66개 확장 (커스텀)
```

### advanced.expand-50 (`src/cmd/optimized/advanced/expand_50.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - deprecated, Advanced로 이동
**목적**: 히스토리 50개 커밋 확장 (deprecated - expand 50 사용)
**구현 내용**:
```bash
# deprecated - 대신 사용:
ga opt advanced expand 50
```

### advanced.expand-100 (`src/cmd/optimized/advanced/expand_100.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - deprecated, Advanced로 이동
**목적**: 히스토리 100개 커밋 확장 (deprecated - expand 100 사용)
**구현 내용**:
```bash
# deprecated - 대신 사용:
ga opt advanced expand 100
```

### quick.auto-find-merge-base (`src/cmd/optimized/quick/auto_find_merge_base.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
**목적**: 두 브랜치의 머지베이스 자동 찾기 (히스토리 자동 확장)
**구현 내용**:
```bash
# 현재 브랜치와 입력받은 타겟 브랜치 비교

1. 머지베이스 찾기 시도
   git merge-base <current-branch> <target-branch>

2. 실패시 점진적 확장
   - git fetch --deepen=10
   - 재시도
   - 필요시 계속 확장
   - 최종적으로 --unshallow

3. 결과 표시
   - 머지베이스 커밋 해시
   - 필요했던 depth
   - 각 브랜치까지의 거리
```

### advanced.check-merge (`src/cmd/optimized/advanced/check_merge.go`)
**상태**: ✅ 구현 완료 (2025-08-26) - Advanced로 이동
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

### setup.clone-slim (`src/cmd/optimized/setup/clone_slim.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
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

### setup.migrate (`src/cmd/optimized/setup/migrate.go`)
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

### setup.performance (`src/cmd/optimized/setup/performance.go`)
**상태**: ✅ 구현 완료 (2025-08-26)
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

### workspace.expand-path (`src/cmd/optimized/workspace/expand_path.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 특정 경로를 Sparse Checkout에 추가
**구현 내용**:
```bash
# 사용자 입력: 경로

1. 경로 유효성 확인
   - Git 저장소에 존재하는지 검증 (utils.PathExistsInRepo)
   - 이미 추가된 경로인지 중복 확인

2. Cone/Non-cone 모드 지능적 전환
   - 파일 경로 감지 시 자동으로 non-cone 모드로 전환
   - 기존 파일 경로가 있으면 non-cone 유지

3. Sparse Checkout에 추가
   git sparse-checkout add <경로>

4. Config 동기화
   - sparse-checkout list를 config.yaml에 자동 저장
   - config.Set() 활용하여 설정 파일 업데이트

5. 결과 표시
   - 활성화된 경로 목록 출력 (최대 10개)
   - 파일/폴더 구분 표시
```

### quick.filter-branch (`src/cmd/optimized/quick/filter_branch.go`)
**상태**: ✅ 구현 완료 (2025-08-27) - Quick으로 이동
**목적**: 브랜치 필터 설정 (특정 브랜치만 표시)
**구현 내용**:
```bash
# 브랜치 필터 설정으로 선택한 브랜치만 노출

1. 필터 모드 선택
   - single: 단일 브랜치만 표시
   - multi: 여러 브랜치 선택

2. 브랜치 타입 선택
   - local: 로컬 브랜치
   - remote: 원격 브랜치

3. 브랜치 선택
   - 리스트에서 번호로 선택
   - 직접 브랜치명 입력
   - multi 모드: 쉼표로 구분하여 여러 개 선택

4. 필터 적용
   git config custom.branchFilter <브랜치1,브랜치2,...>
   
5. 결과 확인
   - 필터링된 브랜치 목록 표시
   - 프로젝트별 설정 저장
```

### quick.clear-filter-branch (`src/cmd/optimized/quick/clear_filter_branch.go`)
**상태**: ✅ 구현 완료 (2025-08-27) - Quick으로 이동
**목적**: 브랜치 필터 제거 (모든 브랜치 표시)
**구현 내용**:
```bash
# 브랜치 필터를 제거하여 모든 브랜치 노출

1. 현재 필터 확인
   git config --get custom.branchFilter

2. 사용자 확인 프롬프트
   
3. 필터 제거
   git config --unset custom.branchFilter
   또는
   git config custom.branchFilter "*"

4. 결과 확인
   - 로컬 브랜치 개수 표시
   - 원격 브랜치 개수 표시
```

### workspace.restore-branch (`src/cmd/optimized/workspace/restore_branch.go`)
**상태**: ✅ DEPRECATED 처리 완료 (2025-08-27)
**목적**: ~~특정 브랜치만 전체 복원~~ (더 이상 사용하지 않음)
**구현 내용**:
```bash
# DEPRECATED - 이 기능은 더 이상 사용하지 않습니다
# 대신 17번 filter-branch와 18번 clear-filter를 사용하세요
```

### quick.shallow (`src/cmd/optimized/quick/shallow.go`)
**상태**: ✅ 구현 완료 (2025-08-27) - Quick으로 이동
**목적**: 히스토리를 지정된 depth로 줄이기 (기본값: 1)
**구현 내용**:
```bash
1. depth 파라미터 처리 (인자 없으면 기본값 1)
2. 현재 상태 백업
3. git pull --depth=[지정된 depth]
4. git gc --prune=now
5. 결과 확인 (새로운 커밋 수 표시)
```

### quick.unshallow (`src/cmd/optimized/quick/unshallow.go`)
**상태**: ✅ 구현 완료 (2025-08-27) - Quick으로 이동
**목적**: 전체 히스토리 복원
**구현 내용**:
```bash
1. git fetch --unshallow
2. 결과 확인
```

### advanced.check-shallow (`src/cmd/optimized/advanced/check_shallow.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
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

### advanced.check-filter (`src/cmd/optimized/advanced/check_filter.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 현재 필터 설정 확인
**구현 내용**:
```bash
1. Global 필터 확인
   git config remote.origin.partialclonefilter

2. 브랜치별 필터 확인
   git config --get-regexp branch.*.partialclonefilter

3. 결과 표시
```

### advanced.config (`src/cmd/optimized/advanced/config.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 최적화 설정 관리 (백업/복원/확인)
**구현 내용**:
```bash
# 사용자 선택: backup, restore, list, check

[Backup]
1. config.yaml 백업 (.gaconfig/backups/{timestamp}/config.yaml)
2. Sparse Checkout 목록 백업 (sparse-checkout.txt)
3. Git 최적화 설정 백업 (git-optimization.txt)

[Restore]
1. 백업 목록에서 선택
2. config.yaml 복원
3. Sparse Checkout 복원
4. Git 설정 복원

[List]
1. 모든 백업 타임스탬프 표시
2. 각 백업의 파일 목록과 크기 표시

[Check]
1. 현재 config.yaml 상태
2. Git 최적화 상태 (Partial Clone, Shallow, Sparse)
3. 백업 정보 요약
```

### submodule.status (`src/cmd/optimized/submodule/status.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 특정 서브모듈의 최적화 상태 확인 (quick.status의 서브모듈 버전)
**구현 내용**:
```bash
# 사용법: ga opt submodule shallow [depth]
# depth를 지정하지 않으면 기본값 1
# 모든 서브모듈에 recursive로 적용

1. depth 파라미터 처리 (기본값: 1)
2. 서브모듈 목록 확인
3. 각 서브모듈에 대해 (recursive):
   - 현재 shallow 상태 확인
   - git pull --depth=[depth] 실행
   - gc로 오래된 객체 정리
4. 결과 요약 표시 (성공/실패 카운트)

# 사용 예시:
ga opt submodule shallow        # depth=1 (기본값)
ga opt submodule shallow 5      # depth=5로 설정
ga opt submodule shallow 10     # depth=10으로 설정
```

### submodule.to-slim (`src/cmd/optimized/submodule/to_slim.go`)
**목적**: 특정 서브모듈을 SLIM 모드로 전환 (quick.to-slim의 서브모듈 버전)
**구현 내용**:
```bash
1. 사용자 확인 프롬프트 (대용량 다운로드 경고)
2. 서브모듈 목록 확인
3. 각 서브모듈에 대해:
   - 현재 shallow 상태 및 depth 확인
   - .git 폴더 크기 측정 (복원 전)
   - git fetch --unshallow 실행
   - 복원 후 크기 측정 및 비교
   - 총 커밋 수 표시
4. 결과 요약:
   - 성공/실패 카운트
   - 전체 크기 변화 표시
   - 각 서브모듈별 크기 증가량
```

### submodule.to-full (`src/cmd/optimized/submodule/to_full.go`)
**목적**: 특정 서브모듈을 FULL 모드로 복원 (quick.to-full의 서브모듈 버전)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름 (없으면 선택 메뉴)

1. 서브모듈 선택
2. 해당 서브모듈로 이동 후:
   - Partial Clone 필터 확인
   - Sparse Checkout 상태 확인
   - Shallow 상태 확인
   - 디스크 사용량 확인
3. 상태 표시 (quick.status와 동일 형식)
```

### submodule.expand-slim (`src/cmd/optimized/submodule/expand_slim.go`)
**목적**: 서브모듈의 선택적 경로 확장 (quick.expand-slim의 서브모듈 버전)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름 (없으면 선택 메뉴)

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - Partial Clone 필터 적용
   - Sparse Checkout 활성화
   - 불필요한 객체 정리
3. 결과 확인
```

### submodule.expand-filter (`src/cmd/optimized/submodule/expand_filter.go`)
**목적**: 서브모듈의 Partial Clone 필터 제거 (quick.expand-filter의 서브모듈 버전)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름 (없으면 선택 메뉴)

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - Sparse Checkout 해제
   - 모든 객체 다운로드
   - Partial Clone 필터 제거
3. 결과 확인
```

### submodule.shallow (`src/cmd/optimized/submodule/shallow.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 서브모듈을 Shallow Clone으로 변환 (depth 파라미터 지원, recursive)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름, 확장할 경로

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - 현재 Sparse Checkout 목록 확인
   - 경로 추가
   - 필요한 객체 다운로드
3. 확장 결과 표시
```

### submodule.unshallow (`src/cmd/optimized/submodule/unshallow.go`)
**상태**: ✅ 구현 완료 (2025-08-27)
**목적**: 서브모듈의 전체 히스토리 복원 (recursive)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름 (없으면 선택 메뉴)

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - 현재 필터 확인
   - 필터 제거
   - 모든 blob 다운로드
3. 결과 확인
```

### submodule.filter-branch (`src/cmd/optimized/submodule/filter_branch.go`)
**목적**: 서브모듈의 브랜치 필터 설정 (quick.filter-branch의 서브모듈 버전)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름, 필터 모드

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - 브랜치 목록 확인
   - 필터 모드 선택 (single/multi)
   - 브랜치 선택
   - 필터 적용
3. 필터링된 브랜치 목록 표시
```

### submodule.clear-filter (`src/cmd/optimized/submodule/clear_filter.go`)
**목적**: 서브모듈의 브랜치 필터 제거 (quick.clear-filter-branch의 서브모듈 버전)
**구현 내용**:
```bash
# 사용자 입력: 서브모듈 이름 (없으면 선택 메뉴)

1. 서브모듈 선택
2. 해당 서브모듈에서:
   - 현재 필터 확인
   - 사용자 확인 프롬프트
   - 필터 제거
3. 결과 확인 (모든 브랜치 표시)
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
feat(opt): implement status - Git repository optimization status check

# 테스트 단계
test(opt): add tests for status command functionality

# 문서화 단계
docs(opt): update documentation for status command usage
```

---

## 📝 함수별 커밋 메시지 템플릿

### Help 카테고리 (도움말)
```bash
# help.workflow
feat(opt): implement workflow - Git optimization workflow guide
test(opt): add tests for workflow help command
docs(opt): document workflow usage and examples

# help.commands  
feat(opt): implement commands - complete command list display
test(opt): add tests for commands help display
docs(opt): document commands help system
```

### Quick 카테고리 (빠른 실행)
```bash
# quick.status
feat(opt): implement status - repository optimization status check
test(opt): add tests for status metrics collection
docs(opt): document status output format and usage

# quick.to-slim
feat(opt): implement to-slim - convert repository to SLIM mode
test(opt): add tests for to-slim conversion process
docs(opt): document to-slim migration workflow

# quick.to-full
feat(opt): implement to-full - restore repository to FULL mode
test(opt): add tests for to-full restoration process
docs(opt): document to-full recovery workflow

# quick.expand-slim
feat(opt): implement expand-slim - selective path expansion
test(opt): add tests for expand-slim path addition
docs(opt): document expand-slim usage scenarios

# quick.expand-filter
feat(opt): implement expand-filter - remove Partial Clone filter
test(opt): add tests for expand-filter removal process
docs(opt): document expand-filter filter management

# advanced.expand
feat(opt): implement expand - extend history by 10 commits
test(opt): add tests for expand depth expansion
docs(opt): document expand-10 history extension

# advanced.expand-50
feat(opt): implement expand-50 - extend history by 50 commits
test(opt): add tests for expand-50 depth expansion
docs(opt): document expand-50 history extension

# advanced.expand-100
feat(opt): implement expand-100 - extend history by 100 commits
test(opt): add tests for expand-100 depth expansion
docs(opt): document expand-100 history extension

# quick.auto-find-merge-base
feat(opt): implement auto-find-merge-base - automatically locate merge base
test(opt): add tests for auto-find-merge-base detection
docs(opt): document auto-find-merge-base functionality

# advanced.check-merge
feat(opt): implement check-merge - verify merge base existence
test(opt): add tests for check-merge validation
docs(opt): document check-merge-base functionality
```

### Setup 카테고리 (초기 설정)
```bash
# setup.clone-slim
feat(opt): implement clone-slim - optimized repository cloning
test(opt): add tests for clone-slim initialization
docs(opt): document clone-slim setup process

# setup.migrate
feat(opt): implement migrate - convert existing repo to SLIM
test(opt): add tests for migrate conversion workflow
docs(opt): document migrate migration strategy

# setup.performance
feat(opt): implement performance - apply performance settings
test(opt): add tests for performance configuration
docs(opt): document performance optimization settings
```

### Workspace 카테고리 (작업공간)
```bash
# workspace.expand-path
feat(opt): implement expand-path - add specific paths to sparse
test(opt): add tests for expand-path path management
docs(opt): document expand-path selective expansion

# quick.filter-branch
feat(opt): implement filter-branch - branch-specific filters
test(opt): add tests for filter-branch filter application
docs(opt): document filter-branch branch filtering

# quick.clear-filter-branch
feat(opt): implement clear-filter-branch - remove all branch filters
test(opt): add tests for clear-filter-branch cleanup process
docs(opt): document clear-filter-branch filter removal

# workspace.restore-branch
feat(opt): implement restore-branch - restore specific branch
test(opt): add tests for restore-branch restoration
docs(opt): document restore-branch branch recovery
```

### Advanced 카테고리 (고급)
```bash
# quick.shallow
feat(opt): implement shallow - reduce history to depth 1
test(opt): add tests for shallow history reduction
docs(opt): document shallow shallow conversion

# quick.unshallow
feat(opt): implement unshallow - restore complete history
test(opt): add tests for unshallow history restoration
docs(opt): document unshallow full recovery

# advanced.check-shallow
feat(opt): implement check-shallow - verify shallow status
test(opt): add tests for check-shallow status detection
docs(opt): document check-shallow status checking

# advanced.check-filter
feat(opt): implement check-filter - inspect filter settings
test(opt): add tests for check-filter configuration check
docs(opt): document check-filter filter inspection

# advanced.config
feat(opt): implement config - backup/restore settings
test(opt): add tests for config save/load
docs(opt): document backup-config configuration management
```

### Submodule 카테고리 (서브모듈)
```bash
# submodule.shallow
feat(opt): implement shallow - shallow all submodules
test(opt): add tests for shallow batch processing
docs(opt): document shallow-all submodule optimization

# submodule.unshallow
feat(opt): implement unshallow - restore all submodules
test(opt): add tests for unshallow batch restoration
docs(opt): document unshallow-all submodule recovery

# submodule.status
feat(opt): implement submodule-status - check individual submodule status
test(opt): add tests for submodule-status status checking
docs(opt): document submodule-status usage

# submodule.to-slim
feat(opt): implement submodule-to-slim - convert submodule to SLIM
test(opt): add tests for submodule-to-slim conversion
docs(opt): document submodule-to-slim optimization

# submodule.to-full
feat(opt): implement submodule-to-full - restore submodule to FULL
test(opt): add tests for submodule-to-full restoration
docs(opt): document submodule-to-full recovery

# submodule.expand-slim
feat(opt): implement submodule-expand-slim - selective path expansion
test(opt): add tests for submodule-expand-slim path management
docs(opt): document submodule-expand-slim usage

# submodule.expand-filter
feat(opt): implement submodule-expand-filter - remove Partial Clone filter
test(opt): add tests for submodule-expand-filter filter removal
docs(opt): document submodule-expand-filter functionality

# submodule.filter-branch
feat(opt): implement submodule-filter-branch - branch filtering for submodule
test(opt): add tests for submodule-filter-branch filter application
docs(opt): document submodule-filter-branch usage

# submodule.clear-filter
feat(opt): implement submodule-clear-filter - clear branch filters
test(opt): add tests for submodule-clear-filter removal
docs(opt): document submodule-clear-filter functionality
```

---

## 🌿 브랜치 전략

### 브랜치 네이밍 규칙
```bash
# 기능 구현 브랜치
feature/opt-<함수명>

# 예시:
feature/opt-status
feature/opt-to-slim
feature/opt-clone-slim
```

### 브랜치 생성 및 작업 순서
```bash
# 1. 브랜치 생성
git checkout -b feature/opt-status

# 2. 구현 작업
# ... 코드 작성 ...

# 3. 커밋 (구현)
git add src/cmd/optimized/quick/status.go
git commit -m "feat(opt): implement status - repository optimization status check"

# 4. 테스트 추가
# ... 테스트 작성 ...

# 5. 커밋 (테스트)
git add src/cmd/optimized/quick/status_test.go
git commit -m "test(opt): add tests for status metrics collection"

# 6. 문서화
# ... 문서 작성 ...

# 7. 커밋 (문서)
git add docs/optimized/status.md
git commit -m "docs(opt): document status output format and usage"

# 8. Push
git push origin feature/opt-status
```

### 릴리스 태그

#### 태그 생성 시점
카테고리별 모든 명령어가 완료되면 즉시 태그를 생성합니다:

- **Help 카테고리**: 2개 모두 완료 시 → `v1.0.0-opt-help`
- **Quick 카테고리**: 10개 모두 완료 시 → `v1.1.0-opt-quick`
- **Setup 카테고리**: 3개 모두 완료 시 → `v1.2.0-opt-setup`
- **Workspace 카테고리**: 4개 모두 완료 시 → `v1.3.0-opt-workspace`
- **Advanced 카테고리**: 5개 모두 완료 시 → `v1.4.0-opt-advanced`
- **Submodule 카테고리**: 9개 모두 완료 시 → `v1.5.0-opt-submodule`
- **전체 완료**: 33개 모두 완료 시 → `v2.0.0-opt-complete`

#### 태그 생성 명령어
```bash
# 카테고리 완료 시 태그 생성
git tag -a v1.0.0-opt-help -m "Complete Help category implementation"
git push origin v1.0.0-opt-help

# 마일스톤별 태그 목록
v1.0.0-opt-help      # Help 카테고리 완료
v1.1.0-opt-quick     # Quick 카테고리 완료
v1.2.0-opt-setup     # Setup 카테고리 완료
v1.3.0-opt-workspace # Workspace 카테고리 완료
v1.4.0-opt-advanced  # Advanced 카테고리 완료
v1.5.0-opt-submodule # Submodule 카테고리 완료
v2.0.0-opt-complete  # 전체 최적화 기능 완료
```

---

## ⚠️ 중요: 체크리스트 업데이트

### 구현 완료 시 체크리스트 업데이트 필수
각 함수 구현이 완료되면 반드시 다음 항목들을 업데이트해야 합니다:

1. **진행 상황 업데이트** (상단 제목)
   ```markdown
   ## 🎯 구현 진행 상황 (1/33)  # 숫자 업데이트
   ```

2. **체크박스 업데이트** 
   ```markdown
   - [x] 01. workflow - Git 최적화 워크플로우 가이드  # 완료된 항목 체크
   ```

3. **구현 상태 표시** (함수별 구현 상세 섹션)
   ```markdown
   ### help.workflow (`src/cmd/optimized/help/workflow.go`)
   **상태**: ✅ 구현 완료 (2025-01-XX)  # 날짜 추가
   ```

### 업데이트 예시
```bash
# 구현 전
- [ ] workflow - Git 최적화 워크플로우 가이드

# 구현 후  
- [x] workflow - Git 최적화 워크플로우 가이드
```

### 체크리스트 업데이트 커밋
```bash
docs(opt): update checklist for workflow completion

- Mark workflow as completed
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
docs(opt): update checklist for <함수명> completion

- <함수명> 완료 표시
- 진행 상황 업데이트 (<현재/33>)
- 완료 날짜 추가

# 예시:
docs(opt): update checklist for status completion

- status 완료 표시
- 진행 상황 업데이트 (3/33)
- 완료 날짜 추가 (2025-08-26)

# 여러 함수 동시 완료 시
docs(opt): update checklist for multiple completions

- status, to-slim 완료 표시
- 진행 상황 업데이트 (4/33)
- 완료 날짜들 추가

# 카테고리 완료 시
docs(opt): complete Help category implementation

- Help 카테고리 전체 명령어 완료 표시
- 진행 상황 업데이트 (2/33)
- 마일스톤 달성 기록

# 부분 구현 또는 진행 중 상태 업데이트
docs(opt): update status implementation progress

- 부분 구현 내용 추가
- 남은 작업 TODO 업데이트
- 블로커나 의존성 문서화
```