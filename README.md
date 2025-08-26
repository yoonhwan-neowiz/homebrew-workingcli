# WorkingCli (ga)

Git 작업을 혁신적으로 개선하는 AI 기반 CLI 도구입니다. 대용량 Git 저장소 최적화부터 AI 커밋 메시지 생성까지, 개발 생산성을 극대화하는 올인원 솔루션.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![AI Powered](https://img.shields.io/badge/AI-Claude%20%7C%20OpenAI-412991?logo=anthropic&logoColor=white)](https://www.anthropic.com)

## 🎯 왜 WorkingCli인가?

- **💾 대용량 저장소 최적화**: Partial Clone & Sparse Checkout으로 획기적인 용량 절감
- **🚀 생산성 향상**: AI 커밋 메시지 자동 생성으로 시간 절약
- **🎨 직관적 UX**: 복잡한 Git 작업을 대화형 인터페이스로 단순화
- **🤖 AI 통합**: Claude와 OpenAI를 자유롭게 전환하며 사용

## 🔑 핵심 기능

### 🗜️ Git 저장소 최적화 (NEW!)
- **Partial Clone & Sparse Checkout**: 필요한 파일만 다운로드
- **Smart Shallow**: 점진적 히스토리 확장으로 안전한 병합
- **28개 최적화 명령어**: SLIM/FULL 모드 전환, 선택적 확장
- **서브모듈 최적화**: 모든 서브모듈 일괄 처리

### 🤖 AI 커밋 메시지 자동 생성
- 변경사항을 분석하여 규칙에 맞는 메시지 생성
- 커밋 컨벤션 자동 적용 (feat, fix, refactor 등)
- 사용자 키워드 기반 커스터마이징

### 📂 대화형 Stage/Unstage
- 파일 번호로 간편한 선택/해제
- 실시간 상태 미리보기
- AI 커밋 메시지 즉시 생성

### 🔀 스마트 충돌 해결
- 충돌 파일 시각화
- 3-way 병합 옵션 제공
- 단계별 해결 가이드

## 🚀 빠른 시작

### 설치

```bash
# Go 설치 (권장)
go install github.com/workingcli@latest

# 또는 직접 빌드
git clone https://github.com/workingcli.git
cd workingcli
./build.sh
```

### 기본 사용법

```bash
# 대화형 스테이징
ga stage     # 또는 단순히 'ga'

# AI 커밋 메시지 생성
ga commit    # diff 기반 자동 생성
ga commit -k "로그인 기능"  # 키워드 지정

# 저장소 최적화 (NEW!)
ga optimized quick status     # 현재 최적화 상태 확인
ga optimized quick to-slim    # SLIM 모드로 전환
ga optimized help commands    # 28개 최적화 명령어 목록

# 충돌 해결
ga resolve

# 히스토리 시각화
ga history
```

## 💡 주요 사용 시나리오

### 1. 대화형 Stage/Unstage

```bash
$ ga stage

Git Status:
  Modified:   [1] src/main.go
  Modified:   [2] src/utils.go
  Untracked:  [3] test/new_test.go

선택된 파일: []
명령어: (숫자)선택, (y)적용, (m)AI커밋, (q)종료
> 1 2
선택된 파일: [src/main.go, src/utils.go]
> y
✓ 파일이 스테이징되었습니다
```

### 2. AI 기반 커밋

```bash
$ ga commit -v

AI가 변경사항을 분석 중...
생성된 커밋 메시지:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
feat: 사용자 인증 모듈 추가

- JWT 토큰 기반 인증 구현
- 로그인/로그아웃 API 엔드포인트 추가
- 세션 관리 미들웨어 구성
━━━━━━━━━━━━━━━━━━━━━━━━━━━━

이 메시지를 사용하시겠습니까? (Y/n)
```

### 3. 대용량 저장소 최적화 (NEW!)

```bash
$ ga optimized quick status

📊 저장소 최적화 상태
━━━━━━━━━━━━━━━━━━
모드: FULL (미최적화)
Partial Clone: 비활성
Sparse Checkout: 비활성
디스크 사용: 103GB
━━━━━━━━━━━━━━━━━━

$ ga optimized quick to-slim

🚀 SLIM 모드로 전환 중...
✓ Partial Clone 필터 적용 (blob:limit=1m)
✓ Sparse Checkout 활성화
✓ 불필요한 객체 정리 완료

📊 최적화 결과
━━━━━━━━━━━━━━━━━━
이전: 103GB → 현재: 30MB
절감: 99.97%
━━━━━━━━━━━━━━━━━━
```

### 4. 충돌 해결

```bash
$ ga resolve

충돌 파일 목록:
1. src/config.go (3개 충돌)
2. src/api.go (1개 충돌)

해결 방법 선택:
1) 현재 브랜치 변경사항 사용 (--ours)
2) 대상 브랜치 변경사항 사용 (--theirs)
3) 수동 병합 모드
> 
```

## 🏗️ 프로젝트 구조

```
WorkingCli/
├── src/
│   ├── cmd/
│   │   ├── optimized/    # Git 저장소 최적화 명령어 (28개)
│   │   │   ├── help/     # 도움말 및 가이드
│   │   │   ├── quick/    # 자주 사용하는 명령어
│   │   │   ├── setup/    # 초기 설정
│   │   │   ├── workspace/# 작업 공간 관리
│   │   │   ├── advanced/ # 고급 기능
│   │   │   └── submodule/# 서브모듈 관리
│   │   ├── git/          # Git 명령어 래퍼
│   │   ├── stage.go      # 대화형 스테이징
│   │   ├── commit.go     # AI 커밋 메시지
│   │   └── resolve.go    # 충돌 해결
│   ├── ai/               # AI 클라이언트
│   │   ├── client.go     # 인터페이스 정의
│   │   ├── claude.go     # Claude API 구현
│   │   └── openai.go     # OpenAI API 구현
│   ├── config/           # 설정 관리 (Viper)
│   └── utils/            # 공통 유틸리티
├── prompt/               # AI 프롬프트 템플릿
├── test/                 # 테스트 코드
└── .gaconfig/            # 프로젝트별 설정
```

## ⚙️ 설정

### 환경 변수

```bash
# AI API 키 설정 (택 1)
export CLAUDE_API_KEY=your-claude-api-key  # Claude (권장)
export OPENAI_API_KEY=your-openai-api-key  # OpenAI
```

### 프로젝트별 설정 (.gaconfig/config.yaml)

```yaml
# AI 제공자 선택
ai:
  provider: claude  # 또는 openai
  model: claude-3-sonnet  # 모델 지정

# 프롬프트 템플릿 경로
prompt:
  commit: .gaconfig/prompt/commit.md
  analyze: .gaconfig/prompt/analyze.md

# Git 설정
git:
  defaultBranch: main
  commitConvention: conventional  # conventional, angular, custom

# 저장소 최적화 설정
optimized:
  mode: slim               # slim 또는 full
  partialClone: blob:limit=1m  # Partial Clone 필터 크기
  sparseCheckout:          # Sparse Checkout 경로
    - src/
    - Assets/Scripts/
    - ProjectSettings/
  shallowDepth: 1          # Shallow Clone 깊이
```

## 🛠️ 기술 스택

### 핵심 기술
- **언어**: Go 1.19+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - 강력한 CLI 구축
- **설정 관리**: [Viper](https://github.com/spf13/viper) - 유연한 설정 관리
- **터미널 UI**: [fatih/color](https://github.com/fatih/color) - 컬러풀한 출력

### AI 통합
- **Claude API**: Anthropic의 최신 AI 모델
- **OpenAI API**: GPT 시리즈 지원

### 아키텍처 패턴
- **Command Pattern**: Cobra 기반 명령어 구조화
- **Factory Pattern**: AI 클라이언트 동적 생성
- **Interface-based Design**: 확장 가능한 AI 제공자 구조

## 🔧 개발 환경

### 필수 요구사항
- Go 1.19 이상
- Git 2.0 이상
- API 키 (Claude 또는 OpenAI)

### 개발 시작하기

```bash
# 저장소 클론
git clone https://github.com/workingcli.git
cd workingcli

# 의존성 설치
go mod download

# 테스트 실행
go test ./...
./test.sh        # 전체 테스트
./ai_test.sh     # AI 모듈 테스트
./git_test.sh    # Git 통합 테스트

# 빌드
./build.sh       # 전 플랫폼 빌드
go build -o ga   # 로컬 빌드

# 개발 모드 실행
go run main.go stage
```

### 코드 기여 가이드

1. **브랜치 전략**: `feature/기능명` 형식 사용
2. **커밋 컨벤션**: Conventional Commits 준수
3. **테스트**: 새 기능은 반드시 테스트 포함
4. **문서화**: 공개 API는 GoDoc 주석 필수

## 📚 상세 문서

- [사용자 가이드](docs/user-guide.md) - 상세한 사용법과 팁
- [개발자 가이드](docs/developer-guide.md) - 아키텍처와 확장 방법
- [API 문서](docs/api.md) - AI 클라이언트 인터페이스
- [시스템 아키텍처](docs/시스템-아키텍처.md) - 전체 구조 설명

## 🤝 기여하기

WorkingCli는 오픈소스 프로젝트입니다. 기여를 환영합니다!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/amazing-feature`)
3. Commit your Changes (`ga commit -k "새 기능"`)
4. Push to the Branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### 기여 영역
- 🐛 버그 수정
- ✨ 새 기능 제안
- 📝 문서 개선
- 🌐 다국어 지원
- 🧪 테스트 추가

## 🚀 로드맵

- [ ] 브랜치 관리 UI 개선
- [ ] 다중 저장소 동시 관리
- [ ] PR 템플릿 자동 생성
- [ ] 코드 리뷰 코멘트 AI 제안
- [ ] Git Hook 통합
- [ ] 팀 설정 공유 기능

## 📝 라이선스

MIT License - [LICENSE](LICENSE) 파일 참조

---

<div align="center">
  
**WorkingCli**로 Git 작업을 더 스마트하게! 🚀

[이슈 제보](https://github.com/workingcli/issues) · [기능 제안](https://github.com/workingcli/discussions)

</div>