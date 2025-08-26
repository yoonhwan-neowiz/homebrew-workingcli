package config

import (
	"fmt"
	"os"
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
		Mode   string `mapstructure:"mode"` // slim 또는 full
		Filter struct {
			Default string            `mapstructure:"default"` // 기본 필터 (1m)
			Options map[string]string `mapstructure:"options"` // 필터 옵션들
		} `mapstructure:"filter"`
		Sparse struct {
			Paths []string `mapstructure:"paths"` // Sparse Checkout 경로들
		} `mapstructure:"sparse"`
	} `mapstructure:"optimize"`
}

var (
	// cfg는 현재 로드된 설정을 저장합니다.
	cfg *Config
	// v는 viper 인스턴스입니다.
	v *viper.Viper
)

// Initialize는 설정을 초기화합니다.
func Initialize() error {
	v = viper.New()

	// 1. Git 저장소 루트 찾기
	gitRoot, err := findGitRoot()
	if err != nil {
		return fmt.Errorf("git root를 찾을 수 없습니다: %w", err)
	}

	// 2. .gaconfig 디렉토리 생성
	configDir := filepath.Join(gitRoot, ".gaconfig")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf(".gaconfig 디렉토리 생성 실패: %w", err)
	}

	// 3. 프롬프트 디렉토리 생성
	promptDir := filepath.Join(configDir, "prompt")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return fmt.Errorf("prompt 디렉토리 생성 실패: %w", err)
	}

	// 4. 기본 설정 파일 생성
	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configFile); err != nil {
			return fmt.Errorf("기본 설정 파일 생성 실패: %w", err)
		}
	}

	// 5. Viper 설정
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// 6. 환경 변수 설정
	v.SetEnvPrefix("GA")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 7. 설정 파일 읽기
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("설정 파일 읽기 실패: %w", err)
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
  filter:
    default: "1m"  # 기본 Partial Clone 필터
    options:
      minimal: "1m"     # 소스코드만 (1MB 미만)
      basic: "25m"      # 코드 + 씬 파일
      extended: "50m"   # 대부분 리소스 포함
      full: "100m"      # 거의 전체
  sparse:
    paths: []  # Sparse Checkout 경로 목록
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
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
	gitRoot, err := findGitRoot()
	if err != nil {
		return "", err
	}

	switch name {
	case "analyze":
		return filepath.Join(gitRoot, ".gaconfig", cfg.Prompt.Analyze), nil
	case "commit":
		return filepath.Join(gitRoot, ".gaconfig", cfg.Prompt.Commit), nil
	default:
		return "", fmt.Errorf("알 수 없는 프롬프트: %s", name)
	}
}

// GetAll은 전체 설정을 반환합니다.
func GetAll() map[string]interface{} {
	return v.AllSettings()
} 