# WorkingCli (ga) - Git을 쉽고 빠르게

> 모든 Git 일을 더 쉽게. `ga` 한 가지 명령으로 큰 저장소도 가볍게, 커밋 메시지도 자동으로, 파일 선택도 대화형으로.

## 🤔 왜 필요한가?

### 문제들
- **큰 저장소가 너무 느림**: 매번 받는 데 오래 걸리고, 디스크도 빨리 찬다
- **커밋 메시지 쓰기 막막함**: 매번 제목·본문 고민
- **어떤 파일을 올릴지 헷갈림**: 파일 많으면 실수하기 쉬움
- **충돌이 무서움**: 어디부터 어떻게 풀어야 할지 모름

### 해결책
- ✅ 최적화된 클론과 히스토리 최소화로 **100GB → 500MB** 용량 절감
- ✅ AI가 변경사항을 읽고 적절한 커밋 메시지 자동 생성
- ✅ 화면에서 번호로 골라 담는 대화형 스테이징
- ✅ 충돌 해결 도우미가 단계별로 안내

## 🚀 5분 안에 시작하기

### 1단계: 설치 (macOS 기준)
```bash
brew tap yoonhwan-neowiz/workingcli
brew install ga
```

### 2단계: AI 기능 활성화 (선택사항)

#### 방법 1: 환경 변수로 설정 (임시)
```bash
# Claude 사용시
export CLAUDE_API_KEY=your-api-key

# 또는 ChatGPT 사용시
export OPENAI_API_KEY=your-api-key
```

#### 방법 2: 설정 파일로 저장 (영구)
```bash
# ~/.gaconfig/config.yaml 파일을 직접 편집
vi ~/.gaconfig/config.yaml

# AI 제공자와 API 키 설정
ai:
  provider: "claude"  # 또는 "openai"
  claude:
    api_key: "your-claude-api-key"
  openai:
    api_key: "your-openai-api-key"
```

💡 **팁**: 설정 파일은 처음 `ga` 명령 실행시 자동 생성됩니다.

### 3단계: 바로 써보기
```bash
# 큰 저장소를 가볍게 받기
ga opt setup clone https://github.com/your/repo.git my-project

# 파일 골라서 커밋 준비
ga stage

# AI가 커밋 메시지 작성
ga commit
```

## 📌 자주 쓰는 명령어 TOP 5

| 명령어 | 설명 | 언제 쓰나요? |
|--------|------|--------------|
| `ga stage` | 번호로 파일 선택 | 커밋할 파일을 고를 때 |
| `ga commit` | AI 커밋 메시지 자동 생성 | 커밋 메시지가 막막할 때 |
| `ga opt setup clone [repo] [folder]` | 최적화된 저장소 복사 | 큰 프로젝트 처음 받을 때 |
| `ga opt quick shallow 1` | 히스토리 최소화 | 저장소가 무거워졌을 때 |
| `ga resolve` | 충돌 해결 도우미 | 머지 충돌이 났을 때 |

## 💡 실전 시나리오

### 시나리오 1: "100GB 프로젝트를 5분 안에 받고 싶어요"
```bash
# 기존 방식 (30분 소요, 100GB 필요)
git clone https://github.com/company/huge-project.git

# ga 방식 (5분 소요, 500MB만!)
ga opt setup clone https://github.com/company/huge-project.git huge-project
cd huge-project
ga opt quick shallow 1  # 더 가볍게 만들기
```
**결과**: 30분 → 5분, 100GB → 500MB

### 시나리오 2: "오늘 작업한 파일만 골라서 커밋하고 싶어요"
```bash
$ ga stage

Git 상태:
  수정됨:   [1] src/main.go
  수정됨:   [2] src/utils.go  
  새 파일:  [3] test/new_test.go
  수정됨:   [4] docs/README.md

선택할 파일 번호를 입력하세요 (예: 1 2 3)
> 1 3

선택된 파일: [src/main.go, test/new_test.go]
파일을 스테이징하시겠습니까? (y/n) 
> y

$ ga commit

AI가 변경사항 분석 중...
생성된 커밋 메시지:
━━━━━━━━━━━━━━━━━━━━━
feat: 메인 로직 개선 및 테스트 추가

- 메인 함수 성능 최적화
- 새로운 단위 테스트 케이스 추가
━━━━━━━━━━━━━━━━━━━━━

이 메시지를 사용하시겠습니까? (Y/n)
> Y
```

### 시나리오 3: "머지했더니 충돌이 났어요!"
```bash
$ ga resolve

충돌 파일 목록:
1. src/config.go (3개 충돌)
2. src/api.go (1개 충돌)

해결 방법을 선택하세요:
1) 내 변경사항 사용 (현재 브랜치)
2) 상대방 변경사항 사용 (머지 브랜치)
3) 수동으로 하나씩 해결

선택 (1/2/3): 1

✓ 모든 충돌이 해결되었습니다!
커밋하려면 'ga commit'을 실행하세요.
```

## 💾 설치 방법

### macOS / Linux
```bash
# Homebrew로 설치 (권장)
brew tap yoonhwan-neowiz/workingcli
brew install ga
```

### Windows
```bash
# 준비 중입니다
# 현재는 WSL(Windows Subsystem for Linux) 환경에서 사용 가능
```

## 🎯 누가 쓰면 좋나요?

- **대용량 프로젝트 개발자**: Unity, Unreal 같은 큰 프로젝트
- **Git 초보자**: 복잡한 Git 명령어 대신 대화형 인터페이스
- **효율을 중시하는 팀**: 시간과 저장 공간 절약
- **AI 도구를 좋아하는 사람**: 커밋 메시지 자동화

## 🆘 도움이 필요하신가요?

```bash
# 전체 도움말 보기
ga help

# 최적화 명령어 도움말
ga opt help

# 특정 명령어 도움말
ga stage --help
```

## 📊 실제 효과

| 항목 | 기존 Git | WorkingCli (ga) | 개선율 |
|------|----------|-----------------|--------|
| 100GB 저장소 클론 | 30분 | 5분 | 6배 빠름 |
| 디스크 사용량 | 100GB | 500MB | 200배 절약 |
| 커밋 메시지 작성 | 5분 | 30초 | 10배 빠름 |
| 파일 스테이징 | 명령어 여러 번 | 대화형 한 번 | 3배 간편 |

---

**WorkingCli**로 Git 작업을 더 쉽고 빠르게! 🚀

문제가 있거나 제안사항이 있으시면 [이슈 제보](https://github.com/yoonhwan-neowiz/WorkingCli/issues)로 알려주세요.