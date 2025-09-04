package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config는 전역 설정을 관리하는 구조체입니다.
type Config struct {
	AI struct {
		Provider string `mapstructure:"provider"` // openai 또는 claude
		OpenAI struct {
			APIKey string `mapstructure:"api_key"`
			Model  string `mapstructure:"model"` // 사용할 GPT 모델
		} `mapstructure:"openai"`
		Claude struct {
			APIKey string `mapstructure:"api_key"`
			Model  string `mapstructure:"model"` // 사용할 Claude 모델
		} `mapstructure:"claude"`
	} `mapstructure:"ai"`
	Prompt struct {
		Analyze string `mapstructure:"analyze"` // 분석 프롬프트 파일 경로
		Commit  string `mapstructure:"commit"`  // 커밋 프롬프트 파일 경로
	} `mapstructure:"prompt"`
	Optimize struct {
		Mode        string   `mapstructure:"mode"`         // slim 또는 full
		BranchScope []string `mapstructure:"branch_scope"` // 브랜치 필터 리스트
		Filter struct {
			Default string            `mapstructure:"default"` // 기본 필터 (1m)
			Options map[string]string `mapstructure:"options"` // 필터 옵션들
		} `mapstructure:"filter"`
		Sparse struct {
			Paths []string `mapstructure:"paths"` // Sparse Checkout 경로들
		} `mapstructure:"sparse"`
		Submodule struct {
			Mode        string   `mapstructure:"mode"`         // slim 또는 full
			BranchScope []string `mapstructure:"branch_scope"` // 서브모듈 브랜치 필터 리스트
			Filter struct {
				Default string            `mapstructure:"default"` // 서브모듈 기본 필터
				Options map[string]string `mapstructure:"options"` // 서브모듈 필터 옵션들
			} `mapstructure:"filter"`
			Sparse struct {
				Paths []string `mapstructure:"paths"` // 서브모듈 Sparse Checkout 경로들
			} `mapstructure:"sparse"`
		} `mapstructure:"submodule"`
	} `mapstructure:"optimize"`
}

var (
	// cfg는 현재 로드된 설정을 저장합니다.
	cfg *Config
	// v는 viper 인스턴스입니다.
	v *viper.Viper
)

// getGlobalConfigDir는 사용자 홈 디렉토리의 .gaconfig 경로를 반환합니다.
func getGlobalConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("사용자 정보를 가져올 수 없습니다: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".gaconfig"), nil
}

// getProjectConfigDir는 현재 프로젝트의 .gaconfig 경로를 반환합니다.
func getProjectConfigDir() (string, error) {
	gitRoot, err := findGitRoot()
	if err != nil {
		// Git 저장소가 아니면 현재 디렉토리 사용
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, ".gaconfig"), nil
	}
	return filepath.Join(gitRoot, ".gaconfig"), nil
}

// Initialize는 설정을 초기화합니다.
func Initialize() error {
	v = viper.New()

	// 1. 글로벌 설정 디렉토리 확인 및 생성
	globalConfigDir, err := getGlobalConfigDir()
	if err != nil {
		return fmt.Errorf("글로벌 설정 디렉토리 경로 가져오기 실패: %w", err)
	}
	
	// 글로벌 설정 디렉토리가 없으면 생성
	if err := os.MkdirAll(globalConfigDir, 0755); err != nil {
		return fmt.Errorf("글로벌 .gaconfig 디렉토리 생성 실패: %w", err)
	}

	// 글로벌 프롬프트 디렉토리 생성
	globalPromptDir := filepath.Join(globalConfigDir, "prompt")
	if err := os.MkdirAll(globalPromptDir, 0755); err != nil {
		return fmt.Errorf("글로벌 prompt 디렉토리 생성 실패: %w", err)
	}

	// 글로벌 설정 파일 생성 (없으면)
	globalConfigFile := filepath.Join(globalConfigDir, "config.yaml")
	if _, err := os.Stat(globalConfigFile); os.IsNotExist(err) {
		if err := createDefaultConfig(globalConfigFile); err != nil {
			return fmt.Errorf("글로벌 기본 설정 파일 생성 실패: %w", err)
		}
		// 기본 프롬프트 파일들도 생성
		if err := createDefaultPrompts(globalPromptDir); err != nil {
			return fmt.Errorf("기본 프롬프트 파일 생성 실패: %w", err)
		}
	}

	// 2. 프로젝트 설정 디렉토리 확인
	projectConfigDir, err := getProjectConfigDir()
	if err != nil {
		// 프로젝트 설정을 못 가져오면 글로벌 설정만 사용
		projectConfigDir = ""
	}

	// 3. Viper 설정 - 프로젝트 설정이 있으면 우선 사용
	var configFile string
	projectConfigFile := ""
	
	if projectConfigDir != "" {
		projectConfigFile = filepath.Join(projectConfigDir, "config.yaml")
		if _, err := os.Stat(projectConfigFile); err == nil {
			// 프로젝트 설정 파일이 있으면 사용
			configFile = projectConfigFile
		}
	}

	// 프로젝트 설정이 없으면 글로벌 설정 사용
	if configFile == "" {
		configFile = globalConfigFile
	}

	// 4. Viper 설정
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// 5. 환경 변수 설정
	v.SetEnvPrefix("GA")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 6. 설정 파일 읽기
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	// 7. 글로벌 설정을 기본값으로 병합 (프로젝트 설정이 우선)
	if configFile == projectConfigFile {
		// 프로젝트 설정을 사용 중이면 글로벌 설정을 기본값으로 병합
		globalViper := viper.New()
		globalViper.SetConfigFile(globalConfigFile)
		globalViper.SetConfigType("yaml")
		if err := globalViper.ReadInConfig(); err == nil {
			// 글로벌 설정의 각 키를 기본값으로 설정
			for key, value := range globalViper.AllSettings() {
				v.SetDefault(key, value)
			}
		}
	}

	// 8. 설정 구조체에 매핑
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("설정 매핑 실패: %w", err)
	}

	return nil
}

// findGitRoot는 현재 디렉토리에서 .git 디렉토리를 찾아 Git 저장소의 루트 경로를 반환합니다.
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("Git 저장소를 찾을 수 없습니다")
		}
		dir = parent
	}
}

// createDefaultConfig는 기본 설정 파일을 생성합니다.
func createDefaultConfig(path string) error {
	defaultConfig := `# GA CLI 설정 파일

# AI 설정
ai:
  provider: "claude"  # 또는 "openai"
  openai:
    api_key: ""  # GA_AI_OPENAI_API_KEY 환경 변수로 설정 가능
    model: "gpt-4-turbo-preview"  # 사용할 GPT 모델
  claude:
    api_key: ""  # GA_AI_CLAUDE_API_KEY 환경 변수로 설정 가능
    model: "claude-opus-4-20250514"  # 사용할 Claude 모델

# 프롬프트 설정
prompt:
  analyze: "prompt/analyze.md"  # 분석 프롬프트 파일 경로
  commit: "prompt/commit.md"   # 커밋 프롬프트 파일 경로

# Git 최적화 설정
optimize:
  mode: "full"  # slim 또는 full
  branch_scope: []  # 브랜치 필터 리스트 (빈 리스트 = 모든 브랜치)
  filter:
    default: "1m"  # 기본 Partial Clone 필터
    options:
      minimal: "1m"     # 소스코드만 (1MB 미만)
      basic: "25m"      # 코드 + 씬 파일
      extended: "50m"   # 대부분 리소스 포함
      full: "100m"      # 거의 전체
  sparse:
    paths: []  # Sparse Checkout 경로 목록
  submodule:
    mode: "full"  # slim 또는 full
    branch_scope: []  # 서브모듈 브랜치 필터 리스트
    filter:
      default: "1m"
      options:
        minimal: "1m"
        basic: "25m"
        extended: "50m"
        full: "100m"
    sparse:
      paths: []
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

// createDefaultPrompts는 기본 프롬프트 파일들을 생성합니다.
func createDefaultPrompts(promptDir string) error {
	// commit.md 기본 템플릿
	commitPrompt := `# Git 커밋 메시지 생성 프롬프트

아래 정보를 바탕으로 Git 커밋 메시지만 생성해주세요.
추가 설명이나 의견 없이 커밋 메시지 형식에 맞춰 작성해주세요.

## 입력 정보
{{ if .taskType }}
Task 이름(taskType):
{{.taskType}}
{{ end }}

{{ if .taskDesc }}
작업 설명(taskDesc):
{{.taskDesc}}
{{ end }}

파일 목록:
{{.files}}

변경 내용:
` + "```text" + `
{{.diff}}
` + "```" + `

키워드: {{.keyword}}

## 메시지 형식

첫 줄은 반드시 다음 형식을 지켜주세요:
` + "```" + `
{{ if .taskType }}{{.taskType}}: {{ end }}type[(scope)]: subject
` + "```" + `

Type (다음 중 하나 선택):
- feat: 새로운 기능 추가
- fix: 버그 수정
- docs: 문서 수정
- style: 코드 포맷팅
- refactor: 코드 리팩토링
- test: 테스트 코드
- chore: 빌드 업무 수정
- perf: 성능 개선
- ci: CI 설정 변경
- build: 빌드 시스템 변경
- revert: 이전 커밋 되돌리기

Scope: 변경 범위 (선택사항, 예: api, ui, db)

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

## 응답 형식
커밋 메시지만 작성해주세요. 설명이나 의견을 제외하고 아래 형식으로만 응답해주세요:

{{ if .taskType }}{{.taskType}}: {{ end }}type[(scope)]: subject

- 변경 내용 설명
- 변경 이유 설명
`

	// analyze.md 기본 템플릿
	analyzePrompt := `# Git 변경사항 분석 프롬프트

아래 정보를 바탕으로 Git 변경사항을 분석해주세요.

## 입력 정보

파일 목록:
{{.files}}

변경 내용:
` + "```text" + `
{{.diff}}
` + "```" + `

키워드: {{.keyword}}

## 분석 요청사항

다음 관점에서 변경사항을 분석해주세요:

1. **변경 요약**
   - 어떤 기능이 추가/수정/삭제되었는지
   - 주요 변경사항 목록

2. **코드 품질**
   - 잠재적 버그나 문제점
   - 개선할 수 있는 부분
   - 베스트 프랙티스 준수 여부

3. **영향 범위**
   - 다른 부분에 미칠 영향
   - 테스트가 필요한 부분
   - 성능에 미칠 영향

4. **추천사항**
   - 추가로 필요한 작업
   - 리팩토링 제안
   - 문서화 필요 사항
`

	// commit.md 파일 생성
	commitPath := filepath.Join(promptDir, "commit.md")
	if _, err := os.Stat(commitPath); os.IsNotExist(err) {
		if err := os.WriteFile(commitPath, []byte(commitPrompt), 0644); err != nil {
			return fmt.Errorf("commit.md 생성 실패: %w", err)
		}
	}

	// analyze.md 파일 생성
	analyzePath := filepath.Join(promptDir, "analyze.md")
	if _, err := os.Stat(analyzePath); os.IsNotExist(err) {
		if err := os.WriteFile(analyzePath, []byte(analyzePrompt), 0644); err != nil {
			return fmt.Errorf("analyze.md 생성 실패: %w", err)
		}
	}

	return nil
}

// Get은 설정 구조체를 반환합니다.
func Get() *Config {
	return cfg
}

// GetString은 문자열 설정 값을 반환합니다.
func GetString(key string) string {
	return v.GetString(key)
}

// GetBool은 불리언 설정 값을 반환합니다.
func GetBool(key string) bool {
	return v.GetBool(key)
}

// Set은 설정 값을 저장합니다.
func Set(key string, value interface{}) error {
	v.Set(key, value)
	return v.WriteConfig()
}

// GetPromptPath는 프롬프트 파일의 절대 경로를 반환합니다.
func GetPromptPath(name string) (string, error) {
	var promptFile string
	switch name {
	case "analyze":
		promptFile = cfg.Prompt.Analyze
	case "commit":
		promptFile = cfg.Prompt.Commit
	default:
		return "", fmt.Errorf("알 수 없는 프롬프트: %s", name)
	}

	// 1. 먼저 프로젝트 설정 확인
	projectConfigDir, err := getProjectConfigDir()
	if err == nil {
		projectPath := filepath.Join(projectConfigDir, promptFile)
		if _, err := os.Stat(projectPath); err == nil {
			return projectPath, nil
		}
	}

	// 2. 프로젝트에 없으면 글로벌 설정 사용
	globalConfigDir, err := getGlobalConfigDir()
	if err != nil {
		return "", fmt.Errorf("글로벌 설정 디렉토리를 가져올 수 없습니다: %w", err)
	}
	
	globalPath := filepath.Join(globalConfigDir, promptFile)
	if _, err := os.Stat(globalPath); err == nil {
		return globalPath, nil
	}

	// 3. 글로벌에도 없으면 기본 템플릿 생성
	globalPromptDir := filepath.Join(globalConfigDir, "prompt")
	if err := createDefaultPrompts(globalPromptDir); err != nil {
		return "", fmt.Errorf("기본 프롬프트 파일 생성 실패: %w", err)
	}
	
	return globalPath, nil
}

// GetAll은 전체 설정을 반환합니다.
func GetAll() map[string]interface{} {
	return v.AllSettings()
}

// GetBranchScope는 브랜치 스코프 설정을 반환합니다.
func GetBranchScope() []string {
	if cfg == nil || cfg.Optimize.BranchScope == nil {
		return []string{}
	}
	return cfg.Optimize.BranchScope
}

// SetBranchScope는 브랜치 스코프를 설정합니다.
func SetBranchScope(branches []string) error {
	v.Set("optimize.branch_scope", branches)
	cfg.Optimize.BranchScope = branches
	return v.WriteConfig()
}

// ClearBranchScope는 브랜치 스코프를 제거합니다.
func ClearBranchScope() error {
	return SetBranchScope([]string{})
}

// GetSubmoduleBranchScope는 서브모듈 브랜치 스코프 설정을 반환합니다.
func GetSubmoduleBranchScope() []string {
	if cfg == nil || cfg.Optimize.Submodule.BranchScope == nil {
		return []string{}
	}
	return cfg.Optimize.Submodule.BranchScope
}

// SetSubmoduleBranchScope는 서브모듈 브랜치 스코프를 설정합니다.
func SetSubmoduleBranchScope(branches []string) error {
	v.Set("optimize.submodule.branch_scope", branches)
	cfg.Optimize.Submodule.BranchScope = branches
	return v.WriteConfig()
}

// ClearSubmoduleBranchScope는 서브모듈 브랜치 스코프를 제거합니다.
func ClearSubmoduleBranchScope() error {
	return SetSubmoduleBranchScope([]string{})
} 