# WorkingCli PRD

> 참고: (todo)로 표시된 항목들은 향후 구현 예정이며 현재는 무시합니다.

## 1. 핵심 기능 정의
### 1.1 필수 구현 기능
1. ✓ Cobra 기반 CLI 구조
2. AI API 연동 (todo):
   - Claude API (환경변수: CLAUDE_API_KEY)
   - OpenAI API (환경변수: OPENAI_API_KEY)
   - API 키가 설정되지 않은 경우 사용자에게 입력 프롬프트 표시
3. Git 기본 명령어 래핑:
   - ✓ status: 현재 상태 확인
     - ✓ 기본 상태 표시
     - ✓ 한글 파일명 처리
   - branch list (todo): 브랜치 목록 표시
   - ✓ pull: 원격 저장소에서 가져오기
     - ✓ 서브모듈 자동 업데이트
     - ✓ 병렬 처리 최적화
     - ✓ 재귀적 업데이트 지원
   - ✓ push: 원격 저장소로 푸시
     - ✓ 서브모듈 자동 푸시
     - ✓ 병렬 처리 최적화
     - ✓ 재귀적 푸시 지원
   - ✓ reset: 현재 브랜치를 특정 커밋으로 되돌리기
     - ✓ --hard: 워킹 디렉토리까지 모두 되돌리기
     - ✓ --soft: 커밋만 되돌리고 변경사항은 스테이징 영역에 유지
     - ✓ 변경사항 미리보기
     - ✓ 안전한 실행을 위한 확인 절차
   - ✓ diff: 변경사항 확인
     - ✓ 워킹 디렉토리 변경사항 확인
     - ✓ staged 변경사항 확인 (--staged)
     - ✓ 특정 파일의 변경사항 확인
     - ✓ 한글 파일명 처리
   - ✓ checkout: 브랜치 전환 또는 파일 복원
     - ✓ 브랜치 전환
     - ✓ 새 브랜치 생성 및 전환 (-b)
     - ✓ 파일 복원
     - ✓ 강제 브랜치 전환 (-f)
   - ✓ merge: 브랜치 병합
     - ✓ 일반 병합
     - ✓ 병합 커밋 생성 안 함 (--no-commit)
     - ✓ fast-forward 하지 않음 (--no-ff)
     - ✓ 스쿼시 병합 (--squash)
     - ✓ 충돌 발생 시 resolve 명령어로 연계
   - ✓ rebase: 브랜치 리베이스
     - ✓ 일반 리베이스
     - ✓ 대화형 리베이스 (--interactive)
     - ✓ 특정 커밋으로 리베이스 (--onto)
     - ✓ 충돌 발생 시 상세한 해결 방법 안내
   - ✓ switch: 브랜치 전환 (Git 2.23+)
     - ✓ 브랜치 전환
     - ✓ 새 브랜치 생성 및 전환 (-c)
     - ✓ 강제 브랜치 전환 (-f)
   - ✓ tag: 태그 관리
     - ✓ 태그 목록 조회
     - ✓ 태그 생성
     - ✓ 주석이 있는 태그 생성 (-a)
     - ✓ 태그 삭제 (-d)
   - ✓ fetch: 원격 저장소의 변경사항 가져오기
     - ✓ 기본 원격 저장소의 변경사항 가져오기
     - ✓ 특정 원격 저장소의 변경사항 가져오기
     - ✓ 모든 원격 저장소의 변경사항 가져오기 (--all)
     - ✓ 원격에서 삭제된 브랜치 정리 (--prune)
     - ✓ 모든 태그 가져오기 (--tags)
   - submodule (todo): 모든 서브모듈에서 명령어 실행
     - -r/--recursive: 중첩된 서브모듈에서도 명령어 실행
   - stage: 대화형 stage/unstage 인터페이스
     - 파일 상태 시각화 (Added/Modified/Deleted/Renamed/Conflict)
     - 다중 선택 및 범위 선택 지원
     - Git 상태 및 브랜치 컨텍스트 표시
     - 한글 파일명 처리:
       - Git 설정 자동 적용: core.quotepath false
       - 한글 파일명 인코딩/디코딩 처리
       - 한글 파일명이 포함된 diff 출력 처리
     - Untracked 파일 처리:
       - 디렉토리의 경우 하위 파일 재귀적 탐색
       - 개별 파일 단위로 stage 가능
       - 디렉토리 전체 stage 지원
     - 충돌 파일 감지 및 해결 모드 전환
       - 충돌 파일이 있을 경우 알림 표시: "[!] 충돌이 발생한 파일이 있습니다."
       - 'x' 키를 통한 충돌 해결 모드 진입
       - 충돌 해결 후 stage 모드로 자동 복귀
     - AI 커밋 메시지 생성 ('m' 키)
       - stage된 파일 목록 표시
       - 키워드 입력 프롬프트
       - 생성된 메시지 확인 및 커밋 여부 선택
       - 커밋 완료 시 프로그램 종료
       - 실패 시 stage 모드로 복귀
   - resolve: 대화형 충돌 해결
     - 브랜치 컨텍스트 기반 충돌 해결
     - 단계별 충돌 해결 가이드
     - stage 모드와의 양방향 전환
       - stage 모드에서 'x' 키로 진입
       - resolve 모드에서 's' 키로 복귀
     - 충돌 해결 옵션:
       - 현재 브랜치 변경사항 사용
       - 대상 브랜치 변경사항 사용
       - 수동 병합 모드
         - 백업 파일 자동 생성
         - 외부 편집기 지원
         - 충돌 상태 복원 기능
     - 진행 상황 시각화:
       - 해결된 파일: [✓]
       - 현재 파일: [>]
       - 미해결 파일: [ ]
     - 충돌 해결 완료 후 처리:
       ```bash
       모든 충돌이 해결되었습니다!
       다음 중 선택해주세요:
       1. merge commit 생성
       2. stage로 돌아가기
       3. 종료

       선택: _
       ```
       - 잘못된 선택시 다시 선택하도록 루프 처리
       - 각 선택에 따른 동작:
         1. merge commit 생성 후 종료
         2. stage 모드로 전환
         3. merge commit 생성하지 않고 종료

4. AI 기능:
   - auto commit: 변경사항 분석 후 커밋 메시지 자동 생성
        - 커밋은 사용자가 대략의 키워드를 입력할 수 있어 어떤작업을 했는지 빈값으로 들어가면 기본 프롬프트를 사용자 키워드나 문장이 들어오면 내용을 참고해서 판단하는 기능 고도화.
        - 성능 최적화:
          - 첫 단계에서는 스테이징된 파일 목록만 분석
          - `--with-diff`/`-v` 옵션으로 전체 diff 분석 활성화 (verbose 모드)
          - 사용자 확인 단계 추가로 불필요한 API 호출 방지
          - 토일 제한:
            - 크기 제한: 1MB 초과 파일은 diff 제외
            - 타입 제한: 소스 코드 파일만 diff 포함
            - 제외된 파일은 파일명과 크기만 전달
          - 토큰 제한 관리:
            - 기본 max_tokens: 4000
            - 토큰 수 초과 시 대화형 확인:
       ```bash
              토큰 수가 제한을 초과합니다.
              현재 설정: 4000 토큰
              예상 사용: 5200 토큰
              계속 진행하시겠습니까? (y/n):
       ```
            - 사용자가 'y'를 선택하면 요청 진행
            - 'n'을 선택하면 작업 취소 또는 대안 제시
        - AI 모델 선택:
          - `-t`/`--type` 옵션으로 AI 모델 선택 (claude/gpt)
          - 기본값은 Claude 사용
        - 응답 언어:
          - 기본적으로 한글로 응답
          - 파일명이나 코드 관련 전문 용어는 영어 원문 유지
        - 사용 예:
       ```bash
          workingcli ai commit                     # 파일 목록만으로 커밋 메시지 생성 (Claude)
          workingcli ai commit -k "기능 추가"       # 키워드와 파일 목록으로 생성
          workingcli ai commit --with-diff         # 전체 diff 분석하여 생성
          workingcli ai commit -k "버그 수정" -v    # 키워드와 전체 diff로 생성
          workingcli ai commit -t gpt             # GPT 모델로 커밋 메시지 생성
          workingcli ai commit -t claude          # Claude 모델로 커밋 메시지 생성
       ```
   - analyze: 지정된 기간의 커밋 내역 분석
        - 기간 지정 방식:
          - 최근 N개 커밋 (예: --last 10)
          - 날짜 범위 (예: --since 2024-01-01 --until 2024-01-31)
          - 브랜치 범위 (예: --branch feature/*)
        - 파일 제한:
          - 크기 제한: 1MB 초과 파일은 diff 제외
          - 타입 제한: 소스 코드 파일만 diff 포함
          - 제외된 파일은 파일명과 크기만 전달
        - AI 모델 선택:
          - `-t`/`--type` 옵션으로 AI 모델 선택 (claude/gpt)
          - 기본값은 Claude 사용
        - 응답 언어:
          - 기본적으로 한글로 응답
          - 파일명이나 코드 관련 전문 용어는 영어 원문 유지
     - 분석 내용:
       - 개발 활동 요약
       - 주요 변경사항
       - 기술적 변경사항
       - 품질 및 유지보수
       - 영향도 분석
       - 다음 단계 제안
       - 특이사항

### 1.2 프로젝트 구조
```
workingcli/
├── main.go                    # 엔트리포인트
├── src/
│   ├── ai/                   # AI 관련
│   │   ├── claude.go        # Claude 구현
│   │   ├── openai.go        # OpenAI 구현
│   │   └── client.go        # AI 클라이언트
│   │
│   ├── cmd/                 # CLI 명령어
│   │   ├── root.go         # ✓ 루트 명령어
│   │   ├── git/            # ✓ Git 관련 명령어
│   │   │   ├── status.go   # ✓ 구현 완료
│   │   │   ├── branch.go   # ✓ 구현 완료
│   │   │   ├── pull.go     # ✓ 구현 완료
│   │   │   ├── push.go     # ✓ 구현 완료
│   │   │   ├── reset.go    # ✓ 구현 완료
│   │   │   ├── submodule.go # ✓ 구현 완료
│   │   │   ├── stage.go    # ✓ 구현 완료
│   │   │   └── resolve.go  # ✓ 구현 완료
│   │   └── ai/             # AI 관련 명령어 (todo)
│   │       ├── commit.go
│   │       └── analyze.go
│   │
│   └── git/                # Git 작업 처리
│       └── client.go
│
└── prompt/                 # 프롬프트 템플릿
    ├── commit.md          # 커밋 메시지 생성용
    └── analyze.md         # 분석용
```

### 1.3 명령어 구조
```bash
ga
├── config           # 설정 관리
│   ├── init        # 설정 초기화
│   ├── get         # 설정 값 조회
│   └── set         # 설정 값 저장
├── stage            # 대화형 stage/unstage
├── commit          # AI 기반 커밋 메시지 생성
├── resolve         # 대화형 충돌 해결
├── history         # Git 히스토리 시각화
├── analyze         # 커밋 내역 분석
├── status          # 상태 확인
├── branch          # 브랜치 목록
├── pull            # 원격 저장소에서 가져오기
├── push            # 원격 저장소로 푸시
├── reset           # 특정 커밋으로 되돌리기
└── submodule       # 서브모듈 명령어 실행
```

#### 3.1.2 명령어 구조
```bash
ga stage           # 대화형 stage/unstage 모드 시작
ga stage [files]   # 지정한 파일들을 바로 stage

# 파일 타입별 필터
ga stage --filter=modified   # 수정된 파일만
ga stage --filter=added      # 새로 추가된 파일만
ga stage --filter=deleted    # 삭제된 파일만

# 파일 패턴 지원
ga stage "*.go"              # Go 파일만
ga stage "src/cmd/git/*"     # 특정 디렉토리
```

## 2. 기술 스택
- 언어: Go
- CLI 프레임워크: Cobra
- AI:
  - Claude API (claude-3-opus-20240229 모델)
  - OpenAI API (gpt-4-turbo-preview 모델)
- 의존성 관리: Go Modules

## 3. 구현 우선순위
1. CLI 기본 구조 (Cobra)
2. AI API 연동:
   - Claude API
   - OpenAI API
3. Git 기본 명령어
4. AI 커밋 메시지 생성
5. AI 분석 기능

### 3.1 Git stage/unstage 기능
#### 3.1.1 기능 개요
- `git status` 결과를 기반으로 대화형 stage/unstage 인터페이스 제공
- 파일 상태를 시각적으로 표시하고 번호로 선택하여 상태 변경
- 다중 선택 및 범위 선택 지원

#### 3.1.2 명령어 구조
```bash
workingcli git stage           # 대화형 stage/unstage 모드 시작
workingcli git stage [files]   # 지정한 파일들을 바로 stage

# 파일 타입별 필터
workingcli git stage --filter=modified   # 수정된 파일만
workingcli git stage --filter=added      # 새로 추가된 파일만
workingcli git stage --filter=deleted    # 삭제된 파일만

# 파일 패턴 지원
workingcli git stage "*.go"              # Go 파일만
workingcli git stage "src/cmd/git/*"     # 특정 디렉토리
```

#### 3.1.3 파일 상태 구분
- Changes to be committed: staged 상태의 파일들
  - 선택 시 unstage 동작 수행
  - 상태 표시:
    - [A] Added: 새로 추가된 파일
    - [M] Modified: 수정된 파일
    - [D] Deleted: 삭제된 파일
    - [R] Renamed: 이름이 변경된 파일
- Changes not staged for commit: 수정되었지만 unstaged 상태의 파일들
  - 선택 시 stage 동작 수행
  - staged 파일이 추가로 수정/삭제된 경우 두 번 표시됨:
    - MM: staged 수정 파일이 다시 수정됨
    - AM: staged 새 파일이 수정됨
    - MD: staged 수정 파일이 삭제됨
    - AD: staged 새 파일이 삭제됨
    - RM: staged 이름 변경 파일이 수정됨
    - RD: staged 이름 변경 파일이 삭제됨
- Untracked files: 새로 추가되었지만 unstaged 상태의 파일들
  - 선택 시 stage 동작 수행
  - 디렉토리 처리:
    - 디렉토리가 untracked인 경우 하위 파일 목록 표시
    - 디렉토리 단위 또는 개별 파일 단위로 stage 가능
    - 재귀적 탐색으로 모든 하위 파일 처리
  - 파일 필터링:
    - 특정 확장자만 선택 가능
    - 패턴 매칭으로 파일 선택 지원

#### 3.1.4 선택 방식
- 파일 선택 입력 형식:
  - 단일 선택: "1" (존재하는 파일 번호만 가능)
  - 다중 선택: "1,3,4" (존재하지 않는 번호는 무시)
  - 범위 선택: "1-3" (범위가 유효하지 않으면 에러 메시지 표시)
  - 복합 선택: "1-3,5,7-9" (범위와 개별 선택 조합)

#### 3.1.5 상호작용 기능
1. 상태 표시
   - Git 상태 표시:
     - 현재 브랜치 이름 (파란색)
     - 특수 상태: MERGING, CHERRY-PICKING, REBASING, REVERTING (빨간색)
     - verbose 토글 상태: ON (초록색) / OFF (회색)
   - 각 파일 앞에 번호 부여 (1부터 시작)
   - 선택 상태 표시: [V]선택됨 (파란색), [ ]선택안됨 (회색)
   - 파일 상태 표시: [A]추가, [M]수정, [D]삭제, [R]이름변경, [C]충돌
   - 색상으로 구분:
     - 선택 상태: [V] 파란색, [ ] 회색
     - 파일 상태: Added: 초록, Modified: 노랑, Deleted: 회색, Renamed: 파란색, Conflict: 주황
     - Git 상태: MERGING/CHERRY-PICKING/REBASING/REVERTING: 빨강
     - verbose 토글: ON: 초록, OFF: 회색

2. 명령어 입력
   - 번호 입력: 파일 선택/선택 해제 토글
     - 잘못된 번호 입력 시 에러 메시지 표시: "잘못된 파일 번호입니다. (1-N 사이의 번호를 입력하세요)"
     - 범위가 유효하지 않을 경우 에러 메시지 표시: "잘못된 범위입니다. 시작 번호는 끝 번호보다 작아야 합니다."
   - a: 전체 선택
   - c: 선택 취소
   - i: 선택 반전
   - r: 파일 상태 새로고침 (staged 파일 수정 시 필수)
     - staged 파일이 수정된 경우 unstaged 상태로도 표시됨
     - 이런 경우 'r'을 통해 최신 상태로 업데이트 필요
     - Git 상태가 변경될 때마다 사용 권장
   - d: 선택한 파일의 diff 보기
   - v: verbose 모드 토글 (AI 커밋 메시지에 diff 포함 여부)
     - AI 커밋 메시지 생성 시 diff 포함 여부를 토글
     - 상태가 화면에 표시됨 (ON/OFF)
   - s: 선택한 파일 목록 보기
   - h: 이 도움말 표시
   - x: 충돌 해결 모드로 전환
   - y: 변경 적용 (staged → unstage, unstaged → stage)
   - m: AI 커밋 메시지 생성 모드로 전환
     - 현재 staged 파일들을 기준으로 커밋 메시지 생성
     - 키워드 입력 프롬프트 제공:
       ```bash
       커밋 메시지 생성을 위한 키워드나 설명을 입력하세요 (선택사항):
       ```
     - verbose 토글 상태에 따라 변경 내용 포함 여부 결정
     - 생성된 메시지 확인 후 커밋 여부 선택:
       ```bash
       다음 커밋 메시지로 커밋하시겠습니까? (y/n):
       feat: stage 모드에서 AI 커밋 메시지 생성 기능 추가

       - stage 모드에서 'm' 키를 통해 AI 커밋 메시지 생성 가능
       - 키워드 입력을 통한 컨텍스트 제공 지원
       - staged 파일 목록과 diff 기반 메시지 생성
       ```
   - q: 종료

3. 도움말 표시
```bash
=== Stage/Unstage 도움말 ===
파일 선택:
  - 단일 선택: "1" (1-N 사이의 번호)
  - 다중 선택: "1,3,4" (쉼표로 구분)
  - 범위 선택: "1-3" (시작-끝)
  - 복합 선택: "1-3,5,7-9" (범위와 개별 선택 조합)

명령어:
  a: 모든 파일 선택
  c: 모든 선택 취소
  i: 선택된 항목 반전
  r: 파일 상태 새로고침 (staged 파일 수정 시 필수)
  d: 선택한 파일의 diff 보기
  v: verbose 모드 토글 (AI 커밋 메시지에 diff 포함 여부)
  s: 선택한 파일 목록 보기
  h: 이 도움말 표시
  x: 충돌 해결 모드로 전환
  y: 변경 적용 (staged → unstage, unstaged → stage)
  m: AI 커밋 메시지 생성 모드로 전환
  q: 종료

Enter를 누르면 이전 화면으로 돌아갑니다...
```

4. 변경 적용 확인
```bash
# staged 파일이 있는 경우
다음 파일들을 unstage 하시겠습니까? (y/n):
[M] src/main.go
[A] src/config.go

# unstaged 파일이 있는 경우
다음 파일들을 stage 하시겠습니까? (y/n):
[M] src/utils.go
[D] src/old.go

# 작업 완료 후 이전 화면으로 돌아가 파일 목록 새로고침
```

#### 3.1.6 구현 세부사항
1. Git 상태 관리
   ```go
   type GitState struct {
       Branch    string   // 현재 브랜치
       State     string   // NORMAL, MERGING, CHERRY-PICKING, REBASING, REVERTING
       StateDesc string   // 상태 설명 (예: "develop에서 merge 중")
   }
   ```

2. 파일 상태 관리
   ```go
   type FileStatus struct {
       Index    int      // 표시 번호
       Path     string   // 파일 경로
       Status   string   // 상태 (A/M/D/R/C)
       Staged   bool     // stage 여부
       Selected bool     // 선택 여부
       Conflict bool     // 충돌 여부
   }
   ```

3. 상태 확인 함수
   ```go
   type StatusChecker interface {
       GetGitState() (*GitState, error)      // Git 상태 확인
       GetConflicts() ([]string, error)      // 충돌 파일 목록
       IsInProgress() bool                   // 작업 진행 중 여부
   }
   ```

#### 3.1.7 에러 처리
- 잘못된 파일 번호 입력
- 존재하지 않는 파일 지정
- Git 명령어 실패
- 권한 문제

#### 3.1.8 성능 최적화
- 파일 상태 캐싱
- 대용량 저장소 처리
- 긴 파일 목록 페이징

### 3.2 Git 충돌 해결 기능
#### 3.2.1 기능 개요
- 선택된 파일 중 충돌 상태인 파일들을 대화형으로 해결
- 충돌이 발생한 브랜치 컨텍스트를 명확하게 표시
- 단계별 해결 과정 안내

#### 3.2.2 명령어 구조
```bash
# stage 명령어에서 충돌 파일이 선택된 경우 자동으로 충돌 해결 모드로 전환
workingcli git stage           # 대화형 모드에서 충돌 파일 선택 시

# 직접 충돌 해결 모드 시작
workingcli git resolve        # 모든 충돌 파일에 대해 대화형 해결
workingcli git resolve [files] # 지정한 파일들만 해결
```

#### 3.2.3 충돌 해결 인터페이스
```bash
=== 충돌 해결 모드 ===
[Git 작업 컨텍스트]
작업 유형: MERGING
현재 브랜치: feature/ui (abc1234)
대상 브랜치: develop (def5678)

[충돌 파일 목록] (2/5)
  [✓] src/cmd/git/status.go   # 해결됨
  [>] src/cmd/git/merge.go    # 현재 해결 중
  [ ] src/cmd/git/stage.go    # 미해결
  [ ] src/cmd/git/utils.go    # 미해결
  [ ] src/git/client.go       # 미해결

[현재 파일 상태]
파일: src/cmd/git/merge.go
[충돌 상황]
작업 중이던 브랜치: feature/ui (abc1234)
• 최근 커밋1: "merge" (2024-01-04 15:30) (abc1234)
    • merge file list
• 최근 커밋2: "feat: 대화형 UI 구현" (2024-01-04 15:00) (852ksv2)
    • 주요 변경: UI 로직 추가, 상태 관리 개선

가져오려는 브랜치: develop (def5678)
• 최근 커밋1: "merge" (2024-01-05 09:15) (def5678)
    • merge file list
• 최근 커밋2: "cc: bug fix" (2024-01-05 09:09) (002911)
    • 주요 변경: 에러 핸들링, 로깅 추가

현재 파일: src/cmd/git/merge.go
해결 방법을 선택하세요:
1. feature/ui의 변경사항 사용 (내 브랜치)
2. develop의 변경사항 사용 (대상 브랜치)
3. 수동 병합 모드
4. 다음 파일로 건너뛰기
5. 도움말
q. 종료

선택: _

=== 도움말 ===
파일 해결 방법:
  1: 현재 브랜치 변경사항 사용
  2: 대상 브랜치 변경사항 사용
  3: 수동 병합 모드 진입
  4: 다음 파일로 건너뛰기

네비게이션:
  n: 다음 충돌 파일
  p: 이전 충돌 파일
  r: 상태 새로고침
  l: 충돌 파일 목록

기타:
  h: 도움말 표시
  q: 종료 (모든 충돌이 해결된 경우 merge commit 생성)

종료 시 동작:
1. 모든 충돌이 해결된 경우:
   - 자동으로 merge commit 생성
   - Git의 기본 merge commit 메시지 사용
   - 메시지가 없는 경우 자동 생성:
     ```
     Merge branch 'feature/ui' into develop
     
     * 충돌 해결:
       - src/cmd/git/merge.go
       - src/cmd/git/stage.go
     ```

2. 미해결된 충돌이 있는 경우:
   ```bash
   아직 해결되지 않은 충돌이 있습니다:
   [ ] src/cmd/git/stage.go
   [ ] src/cmd/git/utils.go
   [ ] src/git/client.go
   
   모든 충돌을 해결한 후 다시 시도해주세요.
   종료하시겠습니까? (y/n): _
   ```
```

#### 3.2.4 수동 병합 모드
1. 백업 파일 생성 및 초기화
   ```bash
   === 수동 병합 모드 ===
   현재 파일: src/cmd/git/merge.go

   어느 버전을 기준으로 병합하시겠습니까?
   1. feature/ui의 변경사항 (내 브랜치)
   2. develop의 변경사항 (대상 브랜치)
   q. 종료

   선택: 1

   선택하신 버전으로 초기화했습니다.
   백업 파일이 생성되었습니다:
   - 내 브랜치 버전: /absolute/path/to/merge.go.ours
   - 대상 브랜치 버전: /absolute/path/to/merge.go.theirs
   ```

2. 병합 완료/취소
   ```bash
   파일을 수동 병합한 후 아래 옵션을 선택하세요:
   1. 병합 완료
   2. 취소 (충돌 상태로 복원)
   q. 종료

   선택: _
   ```
   완료시 항상 백업 파일 정리
   취소시 백벙 파일 정리하고 파일을 충돌 상태로 복원.
   충돌 해결 모드로 복귀
  


#### 3.2.5 구현 세부사항
1. 충돌 정보 관리
   ```go
   type ConflictContext struct {
       CurrentBranch  string   // 현재 브랜치 (mine)
       TargetBranch   string   // 대상 브랜치 (theirs)
       Operation      string   // MERGING/CHERRY-PICK/REBASE
       Description    string   // 사용자 친화적 설명
   }

   type ConflictFile struct {
       Path          string   // 파일 경로
       Status        string   // 해결 상태
       LocalChanges  string   // 현재 브랜치 변경사항
       RemoteChanges string   // 대상 브랜치 변경사항
       Resolved      bool     // 해결 여부
       BackupOurs    string   // 내 브랜치 백업 파일 경로
       BackupTheirs  string   // 대상 브랜치 백업 파일 경로
   }

   type ConflictResolver interface {
       GetContext() (*ConflictContext, error)
       GetConflictFiles() ([]ConflictFile, error)
       ResolveUsingTheirs(path string) error
       ResolveUsingMine(path string) error
       CreateBackupFiles(path string) error
       RestoreFromBackup(path string) error
       CleanupBackupFiles(path string) error
       IsResolved(path string) (bool, error)
   }
   ```

#### 3.2.6 에러 처리
- 충돌 파일 접근 권한 문제
- 잘못된 해결 방법 선택
- 파일 상태 확인 실패
- 백업 파일 생성/삭제 실패

#### 3.2.7 성능 최적화
- 충돌 정보 캐싱
- 변경사항 미리 파싱
- 상태 확인 최적화
- 백업 파일 관리 최적화

### 3.3 Git History 시각화
#### 3.3.1 기능 개요
- ✓ 터미널 기반 대화형 Git 히스토리 시각화
- ✓ 브랜치 그래프와 커밋 정보를 터미널 크기에 맞춰 표시
- ✓ 키보드 기반 네비게이션과 필터링 지원
- ✓ 브랜치 노드와 병합 지점 강조 표시
- ✓ Git의 기본 로그 형식과 동일한 그래프 구조 유지
- ✓ 윈도우 기반 스크롤링으로 효율적인 커밋 히스토리 탐색

#### 3.3.2 명령어 구조
```bash
ga history           # 현재 브랜치의 히스토리 표시
ga history [branch]  # 특정 브랜치의 히스토리 표시
ga history --all     # 모든 브랜치의 히스토리 표시

# 필터링 옵션
ga history --branch=feature/*  # 특정 브랜치 패턴만 표시
ga history --since=1.week     # 1주일 이내 커밋만 표시
ga history --author=yoonhwan  # 특정 작성자의 커밋만 표시
```

#### 3.3.3 화면 구조
```bash
=== Git History (feature/stage-ui) ===
최근 30개 커밋 표시 중
단축키:
↑/k: 위로 이동   ↓/j: 아래로 이동
←/h: 이전 브랜치 →/l: 다음 브랜치
g: 처음으로     G: 마지막으로
m: merge commit   b: 브랜치 목록
q: 종료          Enter: 상세 정보

→* ff52664 (HEAD -> test-branch-1) feat: 브랜치 전환 시 상세 정보 출력 및 UX 개선
 * 9ce73f7 build: GA 바이너리 업데이트
 * 9b26e7b feat: 브랜치 전환 및 vim 스타일 단축키 추가
 * 960c96a feat: GA 바이너리 파일 업데이트
 * e568caf feat: Git 명령어 추가 (status, pull, push, reset)
 * 6c5f42f feat: test.txt 파일 내용 업데이트
 |\
 | * b4dabb4 (test-branch-2) New changes on branch 2
 * | 6674b8f New changes on branch 1
 * | 79984f6 fix: Handle empty target commits in resolve command
 * | eecc5de Merge branch 'test-branch-2' into test-branch-1
 |\|
 | * 387d042 Changes on test-branch-2
 * | a2c6d59 More changes on test-branch-1
 |/
 * b017c6e Initial commit on test-branch-1

=== 브랜치 목록 ===
  master
→ test-branch-1  
  test-branch-2
  test-branch1
  test-branch2
```

#### 3.3.4 상호작용 기능
1. 네비게이션
   - ✓ ↑/k: 위로 이동 (윈도우 중간 지점 초과 시 윈도우 이동)
   - ✓ ↓/j: 아래로 이동 (윈도우 중간 지점 초과 시 윈도우 이동)
   - ✓ ←/h: 이전 브랜치로 전환
   - ✓ →/l: 다음 브랜치로 전환
   - ✓ g: 처음으로 이동 (윈도우 시작점 재설정)
   - ✓ G: 마지막으로 이동 (윈도우 끝점 재설정)
   - ✓ m: 다음 merge commit으로 이동 (필요시 윈도우 이동)
   - ✓ M: 이전 merge commit으로 이동 (필요시 윈도우 이동)
   - ✓ b: 브랜치 목록 표시
   - ✓ Enter: 커밋 상세 정보 표시
   - ✓ q: 종료

2. 윈도우 기반 스크롤링
   - 윈도우 크기: 30개 커밋
   - 중간 지점: 15번째 커밋 (윈도우 크기 * 0.5)
   - 스크롤 동작:
     - 커서가 중간 지점을 넘어가면 윈도우 이동
     - 위로 이동 시: 윈도우를 위로 15개 커밋만큼 이동
     - 아래로 이동 시: 윈도우를 아래로 15개 커밋만큼 이동
   - 경계 처리:
     - 처음/마지막 커밋에 도달 시 윈도우 고정
     - 남은 커밋이 30개 미만일 경우 마지막 윈도우 유지
   - 성능 최적화:
     - 윈도우 이동 시 새로운 커밋만 추가 로드
     - 이전 윈도우 커밋 캐싱으로 빠른 이동 지원

3. 시각화 요소
   - ✓ 브랜치 그래프: Git의 기본 로그와 동일한 구조
   - ✓ 현재 선택: 파란색 화살표(→)로 표시
   - ✓ HEAD 위치: (HEAD -> branch-name) 형식으로 표시
   - ✓ 브랜치 정보: 브랜치명과 태그 정보 표시
   - ✓ 커밋 정보: 해시, 메시지 형식으로 간단히 표시

4. 브랜치 전환
   - ✓ 화살표 키나 h/l로 브랜치 간 이동
   - ✓ 브랜치 전환 시 자동으로 checkout
   - ✓ 전환 결과 메시지 표시
   - ✓ 실패 시 에러 메시지 출력

5. 필터 설정 (f 키)
   - ✓ 브랜치 패턴 필터
   - ✓ 작성자 필터
   - ✓ 기간 필터 (1.week, 2.days 등)
   - ✓ 필터 초기화 옵션

#### 3.3.5 구현 세부사항
1. 데이터 구조
   ```go
   type CommitNode struct {
       Hash        string
       Message     string
       Branches    []string
       Tags        []string
       IsHead      bool
       IsMerge     bool
       Graph       string
       IsGraphOnly bool
   }

   type HistoryView struct {
       Commits       []CommitNode
       Cursor        int        // 현재 선택된 커밋 인덱스
       WindowStart   int        // 현재 윈도우의 시작 인덱스
       WindowSize    int        // 윈도우 크기 (기본 30)
       TotalCommits  int        // 전체 커밋 수
       Filter        HistoryFilter
       Layout        ScreenLayout
       CurrentBranch string
   }

   type WindowManager interface {
       MoveWindow(direction int) error    // 윈도우 이동
       UpdateCursor(newPos int) error     // 커서 위치 업데이트
       IsWindowMove(cursorPos int) bool   // 윈도우 이동 필요 여부 확인
       GetVisibleCommits() []CommitNode   // 현재 윈도우에 표시할 커밋 목록
   }
   ```

2. 성능 최적화
   - ✓ 윈도우 기반 페이징 처리 (기본 30개 커밋)
   - ✓ 터미널 크기에 맞춘 출력
   - ✓ 그래프 라인 최적화
   - ✓ 필터링 결과 캐싱
   - ✓ 윈도우 이동 시 증분 로딩
   - ✓ 이전/이후 윈도우 프리페칭

## 4. 환경 설정
### 4.1 필수 환경변수
```bash
# AI 설정
GA_AI_PROVIDER=claude     # 기본 AI 제공자 (claude 또는 openai)
GA_AI_CLAUDE_API_KEY=""  # Claude API 키
GA_AI_OPENAI_API_KEY=""  # OpenAI API 키

# Git 설정
git config --global core.quotepath false  # 한글 파일명 처리
```

### 4.1.1 설정 파일 구조
```yaml
# .gaconfig/config.yaml
ai:
  provider: "claude"  # 또는 "openai"
  openai:
    api_key: ""  # GA_AI_OPENAI_API_KEY 환경 변수로 설정 가능
  claude:
    api_key: ""  # GA_AI_CLAUDE_API_KEY 환경 변수로 설정 가능
   prompt:
  analyze: "prompt/analyze.md"  # 분석 프롬프트 파일 경로
  commit: "prompt/commit.md"   # 커밋 프롬프트 파일 경로
```

### 4.2 Git 성능 최적화 설정
```bash
# 파일 시스템 모니터링 및 캐싱
git config core.fsmonitor true      # 파일 시스템 변경 감지
git config core.untrackedCache true # untracked 파일 캐시
git config core.fscache true        # 파일 시스템 캐시
git config core.fscachesize 100000  # 캐시 크기 설정

# 인덱스 최적화
git config pack.compression 0       # 개발 중 압축 비활성화
git config core.deltaBaseCacheLimit 1  # delta base 캐시 제한
git config index.threads 4          # 인덱스 생성 병렬 처리

# 네트워크 최적화
git config http.postBuffer 524288000  # HTTP 버퍼 크기 (500MB)
git config submodule.fetchJobs 4    # 서브모듈 병렬 fetch
git config push.parallel 4          # 병렬 push 작업 수

# 기타 최적화
git config remote.origin.prune true # 원격 브랜치 자동 정리
git config gc.auto 256             # 자동 gc 임계값 증가
git config gc.aggressiveWindow 0   # aggressive gc 비활성화
```

### 4.3 프롬프트 관리
- 위치: `.gaconfig/prompt/` 디렉토리
- 기본 제공 템플릿:
  - `commit.md`: 커밋 메시지 생성 프롬프트
    - 커밋 메시지 형식 정의
    - 타입 및 스코프 가이드라인
    - 예시 및 모범 사례
  - `analyze.md`: 커밋 분석 프롬프트
    - 분석 범위 및 관점
    - 출력 형식 정의
    - 중점 분석 항목

## 5. 에러 처리
- API 키 미설정
- Git 저장소 미초기화
- API 연결 실패
- Git 명령어 실패
- 사용자 입력 검증

## 6. 성능 고려사항
### 6.1 API 호출 최적화
- 필요한 경우에만 API 호출
- 단계적 정보 수집 (파일 목록 → diff)
- 사용자 확인 단계를 통한 불필요한 호출 방지
- 파일 크기 및 타입 기반 최적화:
  - 1MB 초과 파일 제외
  - 소스 코드 파일만 포함
  - 제외된 파일은 메타데이터만 전달

### 6.2 Git 명령어 최적화
- 필요한 정보만 가져오기 (예: --name-only 옵션 활용)
- 대용량 diff 처리 전략 수립

### 6.3 빌드 및 배포
- 플랫폼별 바이너리 생성 (darwin, linux, windows)
- 실행 파일은 항상 build.sh, test.sh 등 스크립트를 통해 처리

## 7. 구현 현황 및 다음 작업

### 7.1 구현 완료된 기능
1. CLI 기본 구조 (Cobra)
   - ✓ 기본 명령어 구조 설정
   - ✓ 서브커맨드 지원
   - ✓ 플래그 처리

2. Git 기본 명령어:
   - ✓ status: 현재 상태 확인
     - ✓ 기본 상태 표시
     - ✓ 한글 파일명 처리
   
   - ✓ pull: 원격 저장소에서 가져오기
     - ✓ 병렬 처리 최적화 (기본값: 4개 작업)
     - ✓ 서브모듈 지원
     - ✓ 진행 상황 표시
     - ✓ 에러 처리 및 복구 방법 안내
   
   - ✓ push: 원격 저장소로 푸시
     - ✓ 병렬 처리 최적화 (기본값: 4개 작업)
     - ✓ 서브모듈 지원
     - ✓ 진행 상황 표시
     - ✓ 에러 처리 및 복구 방법 안내
   
   - ✓ reset: 현재 브랜치를 특정 커밋으로 되돌리기
     - ✓ 모드별 구현:
       - ✓ --hard: 워킹 디렉토리까지 모두 되돌리기
       - ✓ --soft: 커밋만 되돌리고 변경사항은 스테이징 영역에 유지
       - ✓ --mixed (기본): 커밋과 스테이징 영역을 되돌리고 워킹 디렉토리는 유지
     - ✓ 기능:
       - ✓ 현재 변경사항 초기화
       - ✓ 특정 커밋으로 되돌리기
       - ✓ 변경사항 미리보기
       - ✓ 사용자 확인 절차
       - ✓ 복구 방법 안내 (이전 커밋 해시 제공)
     - ✓ 안전장치:
       - ✓ 변경사항 미리보기
       - ✓ 경고 메시지
       - ✓ 확인 절차
       - ✓ 복구 방법 안내

   - ✓ diff: 변경사항 확인
     - ✓ 워킹 디렉토리 변경사항 확인
     - ✓ staged 변경사항 확인 (--staged)
     - ✓ 특정 파일의 변경사항 확인
     - ✓ 한글 파일명 처리

   - ✓ checkout: 브랜치 전환 또는 파일 복원
     - ✓ 브랜치 전환
     - ✓ 새 브랜치 생성 및 전환 (-b)
     - ✓ 파일 복원
     - ✓ 강제 브랜치 전환 (-f)

   - ✓ merge: 브랜치 병합
     - ✓ 일반 병합
     - ✓ 병합 커밋 생성 안 함 (--no-commit)
     - ✓ fast-forward 하지 않음 (--no-ff)
     - ✓ 스쿼시 병합 (--squash)
     - ✓ 충돌 발생 시 resolve 명령어로 연계

   - ✓ rebase: 브랜치 리베이스
     - ✓ 일반 리베이스
     - ✓ 대화형 리베이스 (--interactive)
     - ✓ 특정 커밋으로 리베이스 (--onto)
     - ✓ 충돌 발생 시 상세한 해결 방법 안내

   - ✓ switch: 브랜치 전환 (Git 2.23+)
     - ✓ 브랜치 전환
     - ✓ 새 브랜치 생성 및 전환 (-c)
     - ✓ 강제 브랜치 전환 (-f)

   - ✓ tag: 태그 관리
     - ✓ 태그 목록 조회
     - ✓ 태그 생성
     - ✓ 주석이 있는 태그 생성 (-a)
     - ✓ 태그 삭제 (-d)

   - ✓ fetch: 원격 저장소의 변경사항 가져오기
     - ✓ 기본 원격 저장소의 변경사항 가져오기
     - ✓ 특정 원격 저장소의 변경사항 가져오기
     - ✓ 모든 원격 저장소의 변경사항 가져오기 (--all)
     - ✓ 원격에서 삭제된 브랜치 정리 (--prune)
     - ✓ 모든 태그 가져오기 (--tags)

### 7.2 다음 구현 작업
1. Git 명령어 최적화:
   - status: 현재 상태 확인
     - 성능 최적화:
       - 파일 상태 캐싱
       - 대용량 저장소 처리
       - 긴 파일 목록 페이징
     - 한글 파일명 처리:
       - core.quotepath 설정 자동화
       - 인코딩/디코딩 처리
       - 상태 표시 개선:
         - 색상 구분
         - 아이콘 추가
         - 상태별 그룹화
   
   - submodule: 서브모듈 명령어 실행
     - 병렬 처리 최적화:
       - 기본 4개 작업
       - 작업 수 조절 옵션
     - 중첩된 서브모듈:
       - -r/--recursive 옵션
       - 깊이 제한 설정
     - 진행 상황 시각화:
       - 진행률 표시
       - 현재 작업 중인 모듈 표시
       - 완료/실패 상태 구분
     - 에러 처리 개선:
       - 상세한 에러 메시지
       - 복구 방법 안내
       - 실패한 모듈 재시도 옵션

2. AI 기능 구현:
   - auto commit: 변경사항 분석 후 커밋 메시지 자동 생성
   - analyze: 지정된 기간의 커밋 내역 분석

3. 성능 최적화:
   - 파일 상태 캐싱
   - 대용량 저장소 처리
   - API 호출 최적화
   - 병렬 처리 개선

4. 사용자 경험 개선:
   - 컬러 테마 지원
   - 진행 상황 표시 개선
   - 에러 메시지 상세화
   - 도움말 시스템 강화
```



