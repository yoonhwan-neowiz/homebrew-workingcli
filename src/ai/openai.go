package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"workingcli/src/utils"
)

const (
	openaiAPIEndpoint = "https://api.openai.com/v1/chat/completions"
	defaultGPTModel   = "gpt-4-turbo-preview"
)

// OpenAIClient는 OpenAI API를 사용하는 AI 클라이언트입니다.
type OpenAIClient struct {
	config Config
}

// NewOpenAIClient는 새로운 OpenAI API 클라이언트를 생성합니다.
func NewOpenAIClient(config Config) (Client, error) {
	return &OpenAIClient{
		config: config,
	}, nil
}

// GenerateCommitMessage는 변경사항을 분석하여 커밋 메시지를 생성합니다.
func (c *OpenAIClient) GenerateCommitMessage(prompt string) (string, error) {
	return c.Complete(prompt)
}

// AnalyzeCommits는 커밋 내역을 분석하여 요약을 생성합니다.
func (c *OpenAIClient) AnalyzeCommits(prompt string) (string, error) {
	return c.Complete(prompt)
}

// Complete는 OpenAI API에 요청을 보내고 응답을 받습니다.
func (c *OpenAIClient) Complete(prompt string) (string, error) {
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
		model = defaultGPTModel
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
	req, err := http.NewRequest("POST", openaiAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}

	// 헤더 설정
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

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

	// 디버깅: API 응답 출력
	fmt.Printf("OpenAI API 응답:\n%s\n", string(body))

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 오류 (상태 코드: %d): %s", resp.StatusCode, string(body))
	}

	// JSON 응답 파싱
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("JSON 디코딩 실패: %v", err)
	}

	// 응답 메시지 추출
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("응답 메시지가 없습니다")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("잘못된 메시지 형식")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("메시지 객체가 없습니다")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("메시지 내용이 없습니다")
	}

	// 백틱이 있는 경우에만 코드 블록 처리
	if strings.Contains(content, "```") {
		lines := strings.Split(content, "\n")
		var result []string
		inCodeBlock := false

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			
			// 백틱 블록의 시작이나 끝인 경우 스킵
			if strings.HasPrefix(trimmed, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}

			// 코드 블록 안에 있는 내용만 포함
			if inCodeBlock {
				result = append(result, line)
			} else {
				result = append(result, trimmed)
			}
		}

		content = strings.Join(result, "\n")
	}

	content = strings.TrimSpace(content)
	return content, nil
}