# 📦 Git Assistant (ga) 릴리스 가이드

## 🚀 빠른 시작

새 버전을 릴리스하려면 단 하나의 명령어만 실행하세요:

```bash
./scripts/release-all.sh 0.1.3
```

이 명령어 하나로 모든 작업이 자동으로 처리됩니다! ✨

## 📋 자동 처리 작업

`release-all.sh` 스크립트는 다음 작업들을 순서대로 자동 실행합니다:

1. **🔨 빌드 실행**
   - 모든 플랫폼 (macOS Intel/ARM, Linux Intel/ARM) 바이너리 빌드
   - 버전 정보 자동 삽입

2. **📦 아카이브 생성**
   - 각 플랫폼별 tar.gz 파일 생성
   - SHA256 체크섬 자동 계산

3. **📝 Formula 업데이트**
   - 버전 정보 업데이트
   - 다운로드 URL 업데이트
   - SHA256 체크섬 업데이트

4. **🔄 Git 작업**
   - Formula 변경사항 자동 커밋
   - origin/main으로 자동 Push

5. **🚀 GitHub Release 생성**
   - 릴리스 태그 생성
   - 릴리스 노트 자동 작성
   - 바이너리 파일 업로드

## 🛠️ 필수 사전 준비

### 1. GitHub CLI 설치
```bash
brew install gh
```

### 2. GitHub 인증
```bash
gh auth login
```

### 3. Go 언어 설치
```bash
brew install go
```

## 📂 프로젝트 구조

```
WorkingCli/
├── build.command           # 실제 빌드 스크립트
├── go.mod                  # Go 프로젝트 파일
└── homebrew-workingcli/    # Homebrew Tap 저장소
    ├── Formula/
    │   └── ga.rb           # Homebrew Formula
    ├── dist/               # 릴리스 파일 저장
    │   └── 0.1.3/
    │       ├── ga-darwin-amd64.tar.gz
    │       ├── ga-darwin-arm64.tar.gz
    │       ├── ga-linux-amd64.tar.gz
    │       └── ga-linux-arm64.tar.gz
    └── scripts/
        └── release-all.sh  # 통합 릴리스 스크립트
```

## 🔍 릴리스 후 테스트

### 기존 설치 업그레이드
```bash
brew update
brew upgrade ga
```

### 새로 설치
```bash
brew tap yoonhwan-neowiz/workingcli
brew install ga
```

### 설치 확인
```bash
ga --version
```

## ⚠️ 주의사항

1. **버전 형식**: 반드시 `X.Y.Z` 형식을 사용하세요 (예: 0.1.3)
2. **GitHub 인증**: 릴리스 생성을 위해 GitHub CLI 인증이 필요합니다
3. **중복 릴리스**: 같은 버전의 릴리스가 이미 존재하면 자동으로 삭제하고 재생성합니다

## 🐛 문제 해결

### SHA256 체크섬 불일치
- 스크립트가 자동으로 올바른 체크섬을 계산하여 적용합니다

### brew가 이전 버전을 설치하는 경우
```bash
# Homebrew 캐시 삭제
brew cleanup ga
rm -rf $(brew --cache)/ga*

# 다시 시도
brew update
brew upgrade ga
```

### Git Push 실패
- GitHub 저장소 권한을 확인하세요
- SSH 키 또는 Personal Access Token 설정을 확인하세요

## 📊 릴리스 체크리스트

- [ ] 버전 번호 결정
- [ ] `./scripts/release-all.sh X.Y.Z` 실행
- [ ] brew update & upgrade 테스트
- [ ] GitHub Release 페이지 확인
- [ ] 설치 테스트 완료

## 🔄 롤백 방법

문제가 발생한 경우:

1. Formula 파일 복구
```bash
git revert HEAD
git push origin main
```

2. GitHub Release 삭제
```bash
gh release delete vX.Y.Z --yes
```

3. 이전 버전으로 재설치
```bash
brew update
brew reinstall ga
```

## 📝 개별 스크립트 (레거시)

통합 스크립트를 사용하는 것을 권장하지만, 필요시 개별 스크립트도 사용 가능합니다:

- `scripts/backup/release.sh` - Formula 업데이트만
- `scripts/backup/build-release.sh` - 빌드 테스트용
- `scripts/backup/create-release.sh` - GitHub Release만 생성

---

**작성일**: 2025-09-08
**버전**: 1.0.0