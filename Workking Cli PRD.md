# WorkingCli

Git 작업을 더 쉽고 효율적으로 만들어주는 CLI 도구입니다.

## 1. 핵심 기능 정의
### 1.1 필수 구현 기능
1. ✓ Cobra 기반 CLI 구조
2. ✓ AI API 연동:
   - ✓ Claude API (환경변수: GA_AI_CLAUDE_API_KEY)
   - ✓ OpenAI API (환경변수: GA_AI_OPENAI_API_KEY)
   - ✓ API 키가 설정되지 않은 경우 사용자에게 입력 프롬프트 표시
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
   - ✓ stage: 대화형 stage/unstage 인터페이스
     - ✓ 파일 상태 시각화 (Added/Modified/Deleted/Renamed/Conflict)
     - ✓ 다중 선택 및 범위 선택 지원
     - ✓ Git 상태 및 브랜치 컨텍스트 표시
     - ✓ 한글 파일명 처리:
       - ✓ Git 설정 자동 적용: core.quotepath false
       - ✓ 한글 파일명 인코딩/디코딩 처리
       - ✓ 한글 파일명이 포함된 diff 출력 처리
     - ✓ Untracked 파일 처리:
       - ✓ 디렉토리의 경우 하위 파일 재귀적 탐색
       - ✓ 개별 파일 단위로 stage 가능
       - ✓ 디렉토리 전체 stage 지원
     - ✓ 충돌 파일 감지 및 해결 모드 전환
       - ✓ 충돌 파일이 있을 경우 알림 표시: "[!] 충돌이 발생한 파일이 있습니다."
       - ✓ 'x' 키를 통한 충돌 해결 모드 진입
       - ✓ 충돌 해결 후 stage 모드로 자동 복귀
     - ✓ AI 커밋 메시지 생성 ('m' 키)
       - ✓ stage된 파일 목록 표시
       - ✓ 키워드 입력 프롬프트
       - ✓ 생성된 메시지 확인 및 커밋 여부 선택
       - ✓ 커밋 완료 시 프로그램 종료
       - ✓ 실패 시 stage 모드로 복귀
   - ✓ resolve: 대화형 충돌 해결
     - ✓ 브랜치 컨텍스트 기반 충돌 해결
     - ✓ 단계별 충돌 해결 가이드
     - ✓ stage 모드와의 양방향 전환
       - ✓ stage 모드에서 'x' 키로 진입
       - ✓ resolve 모드에서 's' 키로 복귀
     - ✓ 충돌 해결 옵션:
       - ✓ 현재 브랜치 변경사항 사용
       - ✓ 대상 브랜치 변경사항 사용
       - ✓ 수동 병합 모드
         - ✓ 백업 파일 자동 생성
         - ✓ 외부 편집기 지원
         - ✓ 충돌 상태 복원 기능
     - ✓ 진행 상황 시각화:
       - ✓ 해결된 파일: [✓]
       - ✓ 현재 파일: [>]
       - ✓ 미해결 파일: [ ]
     - ✓ 충돌 해결 완료 후 처리:
       ```bash
       모든 충돌이 해결되었습니다!
       다음 중 선택해주세요:
       1. merge commit 생성
       2. stage로 돌아가기
       3. 종료

       선택: _
       ```
       - ✓ 잘못된 선택시 다시 선택하도록 루프 처리
       - ✓ 각 선택에 따른 동작:
         1. merge commit 생성 후 종료
         2. stage 모드로 전환
         3. merge commit 생성하지 않고 종료

4. AI 기능:
   - ✓ commit: 변경사항 분석 후 커밋 메시지 자동 생성
        - ✓ 커밋 메시지 생성 프롬프트 템플릿 기반 동작
        - ✓ 성능 최적화:
          - ✓ 첫 단계에서는 스테이징된 파일 목록만 분석
          - ✓ `--verbose`/`-v` 옵션으로 전체 diff 분석 활성화
          - ✓ 사용자 확인 단계 추가로 불필요한 API 호출 방지
          - ✓ 파일 제한:
            - ✓ 크기 제한: 1MB 초과 파일은 diff 제외
            - ✓ 타입 제한: 소스 코드 파일만 diff 포함
            - ✓ 제외된 파일은 파일명과 크기만 전달
        - ✓ AI 모델 선택:
          - ✓ `-t`/`--type` 옵션으로 AI 모델 선택 (claude/gpt)
          - ✓ 기본값은 Claude 사용
        - ✓ 응답 언어:
          - ✓ 기본적으로 한글로 응답
          - ✓ 파일명이나 코드 관련 전문 용어는 영어 원문 유지
        - ✓ 사용 예:
       ```bash
          ga commit                     # 파일 목록만으로 커밋 메시지 생성 (Claude)
          ga commit -k "기능 추가"       # 키워드와 파일 목록으로 생성
          ga commit -v                  # 전체 diff 분석하여 생성
          ga commit -k "버그 수정" -v    # 키워드와 전체 diff로 생성
          ga commit -t gpt             # GPT 모델로 커밋 메시지 생성
          ga commit -t claude          # Claude 모델로 커밋 메시지 생성
       ```
   - ✓ analyze: 지정된 기간의 커밋 내역 분석
        - ✓ 기간 지정 방식:
          - ✓ 최근 N개 커밋 (예: --last 10)
          - ✓ 날짜 범위 (예: --since 2024-01-01 --until 2024-01-31)
          - ✓ 브랜치 범위 (예: --branch feature/*)
        - ✓ 파일 제한:
          - ✓ 크기 제한: 1MB 초과 파일은 diff 제외
          - ✓ 타입 제한: 소스 코드 파일만 diff 포함
          - ✓ 제외된 파일은 파일명과 크기만 전달
        - ✓ AI 모델 선택:
          - ✓ `-t`/`--type` 옵션으로 AI 모델 선택 (claude/gpt)
          - ✓ 기본값은 Claude 사용
        - ✓ 응답 언어:
          - ✓ 기본적으로 한글로 응답
          - ✓ 파일명이나 코드 관련 전문 용어는 영어 원문 유지
     - ✓ 분석 내용:
          - ✓ 개발 활동 요약 (커밋 수, 기여자, 주요 수정 파일)
          - ✓ 주요 변경사항 (기능 추가/수정, 버그 수정, 성능 개선)
          - ✓ 기술적 변경사항 (아키텍처, API, DB 스키마 변경)
          - ✓ 품질 및 유지보수 (테스트, 코드 품질, 기술 부채)
          - ✓ 영향도 분석 (영향 범위, 호환성, 마이그레이션)
          - ✓ 다음 단계 제안 (필요 작업, 잠재적 문제, 개선사항)
          - ✓ 특이사항 (주의 필요 사항, 실험적 기능)

### 1.2 프로젝트 구조
```
workingcli/
├── main.go                    # 엔트리포인트
├── src/
│   ├── ai/                   # AI 관련
│   │   ├── claude.go        # ✓ Claude 구현
│   │   ├── openai.go        # ✓ OpenAI 구현
│   │   └── client.go        # ✓ AI 클라이언트 인터페이스
│   │
│   ├── cmd/                 # CLI 명령어
│   │   ├── root.go         # ✓ 루트 명령어
│   │   ├── commit.go       # ✓ AI 기반 커밋
│   │   ├── analyze.go      # ✓ 커밋 분석
│   │   ├── stage.go        # ✓ 대화형 stage
│   │   └── resolve.go      # ✓ 충돌 해결
│   │
│   ├── config/             # ✓ 설정 관리
│   │   └── config.go       # ✓ Viper 기반 설정
│   │
│   └── utils/              # ✓ 유틸리티
│       └── utils.go        # ✓ 공통 함수
│
└── .gaconfig/              # ✓ 설정 디렉토리
    ├── config.yaml         # ✓ 설정 파일
    └── prompt/            # ✓ 프롬프트 템플릿
        ├── commit.md      # ✓ 커밋 메시지 생성용
        └── analyze.md     # ✓ 분석용
```

### 1.3 명령어 구조
```bash
ga
├── config           # ✓ 설정 관리
│   ├── init        # ✓ 설정 초기화
│   ├── get         # ✓ 설정 값 조회
│   └── set         # ✓ 설정 값 저장
├── stage           # ✓ 대화형 stage/unstage
├── commit          # ✓ AI 기반 커밋 메시지 생성
├── resolve         # ✓ 대화형 충돌 해결
└── analyze         # ✓ 커밋 내역 분석
```

## 2. 설정 관리
### 2.1 설정 파일 구조
```yaml
# .gaconfig/config.yaml

# AI 설정
ai:
  provider: claude  # 기본 AI 제공자 (claude 또는 gpt)
  openai:
    api_key: ""    # OpenAI API 키
  claude:
    api_key: ""    # Claude API 키

# 프롬프트 설정
prompt:
  commit: prompt/commit.md    # 커밋 메시지 생성 프롬프트
  analyze: prompt/analyze.md  # 커밋 분석 프롬프트
```

### 2.2 환경 변수
```bash
# AI API 키 설정
export GA_AI_OPENAI_API_KEY="your-openai-api-key"
export GA_AI_CLAUDE_API_KEY="your-claude-api-key"
```

### 2.3 설정 명령어
```bash
# 설정 초기화
ga config init

# 설정 값 조회
ga config get ai.provider
ga config get ai.openai.api_key
ga config get prompt.analyze

# 설정 값 저장
ga config set ai.provider claude
ga config set prompt.commit "prompt/commit.md"
```

## 3. Git 성능 최적화 설정
### 3.1 파일 시스템 모니터링
```bash
# 파일 시스템 이벤트 감지 활성화
git config core.fsmonitor true

# 파일 시스템 캐시 만료 시간 설정 (초)
git config core.fsmonitorhookversion 2
```

### 3.2 인덱스 최적화
```bash
# 인덱스 버전 설정
git config core.indexversion 4

# 변경 감지 최적화
git config core.untrackedcache true
```

### 3.3 네트워크 최적화
```bash
# 병렬 인덱스 프리로드
git config core.preloadindex true

# 압축 수준 설정
git config core.compression 1

# 델타 압축 설정
git config pack.deltacachesize 128m
git config pack.deltacacheLimit 100
```

### 3.4 기타 성능 설정
```bash
# 커밋 그래프 생성 최적화
git config core.commitgraph true

# 멀티스레드 팩 생성
git config pack.threads 0
```

## 4. 프롬프트 관리
### 4.1 프롬프트 파일 구조
프롬프트 파일은 `.gaconfig/prompt` 디렉토리에 위치하며, 각각의 용도에 맞는 프롬프트 템플릿을 포함합니다.

#### 4.1.1 커밋 메시지 생성 프롬프트
```markdown
# .gaconfig/prompt/commit.md

# Git 커밋 메시지 생성 프롬프트

아래 정보를 바탕으로 Git 커밋 메시지만 생성해주세요.
추가 설명이나 의견 없이 커밋 메시지 형식에 맞춰 작성해주세요.

## 입력 정보
파일 목록:
{{.files}}

{{.diff}}

키워드: {{.keyword}}

## 메시지 형식
type[(scope)]: subject

Type (다음 중 하나 선택):
feat|fix|docs|style|refactor|test|chore|perf|ci|build|revert

Subject:
- 50자 이내
- 한글로 작성 (기술 용어는 영어)
- 마침표 없이
- 현재 시제

본문:
- 한 줄 띄우고 작성
- 각 줄은 72자 이내
- "-" 목록 형식으로 작성
- 변경한 이유와 변경 내용을 상세히 설명
- 여러 줄로 작성 가능
```

#### 4.1.2 커밋 분석 프롬프트
```markdown
# .gaconfig/prompt/analyze.md

# Git 커밋 분석 프롬프트

아래 커밋 정보를 분석하여 개발 활동을 요약해주세요.

## 입력 정보
{{range .Commits}}
커밋: {{.Hash}}
작성자: {{.Author}}
날짜: {{.Date}}
메시지: {{.Message}}

파일 목록:
{{range .Files}}
- {{.}}{{end}}

{{.Diff}}
{{end}}

## 분석 요구사항
다음 관점에서 분석해주세요:

1. 개발 활동 요약
   - 커밋 빈도와 규모
   - 주요 기여자와 작업 영역
   - 핵심 변경 파일과 모듈

2. 주요 변경사항
   - 기능 추가/수정
   - 버그 수정
   - 성능 개선
   - 기술적 변경

3. 품질과 유지보수
   - 코드 품질 개선
   - 테스트 추가/수정
   - 문서화
   - 리팩토링

4. 영향도 분석
   - 변경 범위와 영향
   - 잠재적 위험
   - 주의 필요 사항

5. 다음 단계
   - 필요한 후속 작업
   - 개선 제안
   - 검토 필요 사항

## 응답 형식
- 한글로 작성 (기술 용어는 영어)
- 중요 발견 사항은 굵게 표시
- 각 섹션은 번호로 구분
```

### 4.2 프롬프트 관리 명령어
```bash
# 프롬프트 파일 경로 설정
ga config set prompt.commit "prompt/commit.md"
ga config set prompt.analyze "prompt/analyze.md"

# 프롬프트 파일 경로 확인
ga config get prompt.commit
ga config get prompt.analyze
```