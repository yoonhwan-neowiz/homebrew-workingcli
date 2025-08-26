package ai

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"workingcli/src/config"
)

// Config는 AI 클라이언트 설정입니다.
type Config struct {
	APIKey string
	Model  string // 사용할 AI 모델
}

// Client는 AI 서비스와의 통신을 담당하는 인터페이스입니다.
type Client interface {
	// GenerateCommitMessage는 변경사항을 분석하여 커밋 메시지를 생성합니다.
	GenerateCommitMessage(prompt string) (string, error)

	// AnalyzeCommits는 커밋 내역을 분석하여 요약을 생성합니다.
	AnalyzeCommits(prompt string) (string, error)
}

// CommitInfo는 커밋 정보를 담는 구조체입니다.
type CommitInfo struct {
	Hash        string
	Author      string
	Date        string
	Message     string
	Description string
	Files       []string
	Diff        string
}

// ModelType은 사용할 AI 모델 타입입니다.
type ModelType string

const (
	ModelClaude ModelType = "claude"
	ModelGPT    ModelType = "gpt"
)

// NewClient는 새로운 AI 클라이언트를 생성합니다.
func NewClient(modelType ModelType) (Client, error) {
	switch modelType {
	case ModelGPT:
		// OpenAI API 키 확인
		if apiKey := config.GetString("ai.openai.api_key"); apiKey != "" {
			return NewOpenAIClient(Config{
				APIKey: apiKey,
				Model:  config.GetString("ai.openai.model"),
			})
		}
		return promptForAPIKey("ai.openai.api_key", "OpenAI")

	case ModelClaude, "": // 기본값은 Claude
		// Claude API 키 확인
		if apiKey := config.GetString("ai.claude.api_key"); apiKey != "" {
			return NewClaudeClient(Config{
				APIKey: apiKey,
				Model:  config.GetString("ai.claude.model"),
			})
		}
		return promptForAPIKey("ai.claude.api_key", "Claude")

	default:
		return nil, fmt.Errorf("지원하지 않는 모델 타입입니다: %s", modelType)
	}
}

// promptForAPIKey는 사용자에게 API 키를 입력받습니다.
func promptForAPIKey(configKey, clientType string) (Client, error) {
	fmt.Printf("%s API 키가 설정되지 않았습니다.\nAPI 키를 입력해주세요: ", clientType)
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("API 키 입력 읽기 실패: %v", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errors.New("API 키가 입력되지 않았습니다")
	}

	// 설정 저장
	if err := config.Set(configKey, apiKey); err != nil {
		return nil, fmt.Errorf("API 키 설정 저장 실패: %v", err)
	}

	fmt.Printf("API 키가 설정되었습니다. 다음에는 환경 변수로 설정하는 것을 권장합니다:\nexport GA_AI_%s_API_KEY=your-api-key\n", strings.ToUpper(strings.Split(configKey, ".")[1]))

	switch configKey {
	case "ai.openai.api_key":
		return NewOpenAIClient(Config{
			APIKey: apiKey,
			Model:  config.GetString("ai.openai.model"),
		})
	case "ai.claude.api_key":
		return NewClaudeClient(Config{
			APIKey: apiKey,
			Model:  config.GetString("ai.claude.model"),
		})
	default:
		return nil, fmt.Errorf("지원하지 않는 설정 키입니다: %s", configKey)
	}
} 