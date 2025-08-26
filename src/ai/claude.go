package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"workingcli/src/utils"
)

const (
	claudeAPIEndpoint = "https://api.anthropic.com/v1/messages"
	defaultClaudeModel = "claude-opus-4-20250514"
)

// ClaudeClient는 Claude API를 사용하는 AI 클라이언트입니다.
type ClaudeClient struct {
	config Config
}

// NewClaudeClient는 새로운 Claude API 클라이언트를 생성합니다.
func NewClaudeClient(config Config) (Client, error) {
	return &ClaudeClient{
		config: config,
	}, nil
}

// GenerateCommitMessage는 변경사항을 분석하여 커밋 메시지를 생성합니다.
func (c *ClaudeClient) GenerateCommitMessage(prompt string) (string, error) {
	return c.Complete(prompt)
}

// AnalyzeCommits는 커밋 내역을 분석하여 요약을 생성합니다.
func (c *ClaudeClient) AnalyzeCommits(prompt string) (string, error) {
	return c.Complete(prompt)
}

// Complete는 Claude API에 요청을 보내고 응답을 받습니다.
func (c *ClaudeClient) Complete(prompt string) (string, error) {
	// 토큰 수 예측 (간단한 휴리스틱: 영어 기준 1토큰 ≈ 4글자)
	estimatedTokens := len(prompt) / 4
	if estimatedTokens > 4000 {
		fmt.Println("\n%s", prompt)
		fmt.Printf("\n토큰 수가 제한을 초과합니다.\n")
		fmt.Printf("현재 설정: 4000 토큰\n")
		fmt.Printf("예상 사용: %d 토큰\n", estimatedTokens)
		
		if !utils.ConfirmWithDefault("계속 진행하시겠습니까?", false) {
			return "", fmt.Errorf("사용자가 작업을 취소했습니다")
		}
	}

	// API 요청 데이터 구성
	model := c.config.Model
	if model == "" {
		model = defaultClaudeModel
	}
	requestData := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 4000,
	}

	// JSON 인코딩
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("JSON 인코딩 실패: %v", err)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest("POST", claudeAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}

	// 헤더 설정
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// HTTP 요청 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("응답 읽기 실패: %v", err)
	}

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	// JSON 응답 파싱
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("JSON 디코딩 실패: %v", err)
	}

	// 응답 메시지 추출 (Claude API v3 구조)
	content, ok := response["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("응답 메시지가 없습니다")
	}

	firstContent, ok := content[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("잘못된 메시지 형식")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("메시지 텍스트가 없습니다")
	}

	return text, nil
} 