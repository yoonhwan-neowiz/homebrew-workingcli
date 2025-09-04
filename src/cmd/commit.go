package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"workingcli/src/ai"
	"workingcli/src/config"
	"workingcli/src/utils"
	"bufio"
)

func NewCommitCmd() *cobra.Command {
	var (
		keyword   string
		withDiff  bool
		modelType string
	)

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "AI 기반 커밋 메시지 생성",
		Long: `현재 Git 저장소의 변경사항을 분석하여 AI 기반으로 커밋 메시지를 생성합니다.
사용자가 키워드나 설명을 입력하면 이를 참고하여 더 정확한 메시지를 생성합니다.

사용법:
  ga commit                     # 파일 목록만으로 커밋 메시지 생성
  ga commit -k "기능 추가"       # 키워드와 함께 생성
  ga commit -v                  # verbose 모드로 생성 (diff 포함)
  ga commit -t gpt             # GPT 모델로 커밋 메시지 생성
  ga commit -t claude          # Claude 모델로 커밋 메시지 생성`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 스테이징된 파일 목록 확인
			var out bytes.Buffer
			gitCmd := exec.Command("git", "diff", "--staged", "--name-only")
			gitCmd.Stdout = &out
			if err := gitCmd.Run(); err != nil {
				return fmt.Errorf("스테이징된 파일 목록 확인 실패: %v", err)
			}

			// 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			if len(output) == 0 {
				return fmt.Errorf("스테이징된 변경사항이 없습니다. git add 명령어로 변경사항을 먼저 스테이징해주세요")
			}

			// 변경된 파일 목록 출력
			fmt.Println("\n스테이지된 파일 목록:")
			fmt.Print(output)

			// diff 가져오기
			fmt.Println("\n=== diff 가져오기 시작 ===")
			var diffContent string
			if diff, err := utils.GetDiffForAI(strings.Split(strings.TrimSpace(output), "\n"), "HEAD", withDiff); err == nil {
				diffContent = diff
				fmt.Println("diff 가져오기 성공")
			} else {
				fmt.Printf("diff 가져오기 실패: %v\n", err)
				return fmt.Errorf("diff 가져오기 실패: %v", err)
			}
			fmt.Println("=== diff 가져오기 끝 ===\n")
			// 사용자 입력 받기
			var taskType string
			if utils.ConfirmWithDefault("Task 이름(taskType)을 입력하시겠습니까?", false) {
				fmt.Print("\nTask 이름을 입력해주세요 (예: feature-auth, bugfix-login): ")
				reader := bufio.NewReader(os.Stdin)
				taskType, _ = reader.ReadString('\n')
				taskType = strings.TrimSpace(taskType)
				fmt.Println()
			}

			// taskDesc 입력 받기
			var taskDesc string
			if utils.ConfirmWithDefault("작업 의도나 추가 설명(taskDesc)을 입력하시겠습니까?", false) {
				fmt.Print("\n작업 의도나 추가 설명을 입력해주세요 (AI가 커밋 메시지 작성 시 참고합니다): ")
				reader := bufio.NewReader(os.Stdin)
				taskDesc, _ = reader.ReadString('\n')
				taskDesc = strings.TrimSpace(taskDesc)
				fmt.Println()
			}

			if !utils.ConfirmWithDefault("위 파일들의 변경사항으로 커밋 메시지를 생성하시겠습니까?", true) {
				return nil
			}

			// 프롬프트 템플릿 가져오기
			fmt.Println("\n=== 프롬프트 템플릿 가져오기 시작 ===")
			promptPath, err := config.GetPromptPath("commit")
			if err != nil {
				fmt.Printf("프롬프트 파일 경로 가져오기 실패: %v\n", err)
				return fmt.Errorf("프롬프트 파일 경로 가져오기 실패: %v", err)
			}
			fmt.Printf("프롬프트 파일 경로: %s\n", promptPath)

			// 프롬프트 파일 읽기
			promptContent, err := os.ReadFile(promptPath)
			if err != nil {
				fmt.Printf("프롬프트 파일 읽기 실패: %v\n", err)
				return fmt.Errorf("프롬프트 파일 읽기 실패: %v", err)
			}
			fmt.Printf("프롬프트 파일 크기: %d bytes\n", len(promptContent))
			fmt.Println("=== 프롬프트 템플릿 가져오기 끝 ===\n")

			// 템플릿 데이터 준비
			fmt.Println("\n=== 템플릿 데이터 준비 시작 ===")
			data := map[string]string{
				"taskType":    taskType,
				"taskDesc":    taskDesc,
				"files":       output,
				"diff":        diffContent,
				"keyword":     keyword,
			}
			fmt.Printf("files 길이: %d, diff 길이: %d, keyword: %s, taskType 길이: %d, taskDesc 길이: %d\n", len(output), len(diffContent), keyword, len(taskType), len(taskDesc))
			fmt.Println("=== 템플릿 데이터 준비 끝 ===\n")

			// 템플릿 적용
			fmt.Println("\n=== 템플릿 적용 시작 ===")
			tmpl, err := template.New("commit").Parse(string(promptContent))
			if err != nil {
				fmt.Printf("프롬프트 템플릿 파싱 실패: %v\n", err)
				return fmt.Errorf("프롬프트 템플릿 파싱 실패: %v", err)
			}
			fmt.Println("템플릿 파싱 성공")

			// 프롬프트 템플릿 적용
			var promptBuffer bytes.Buffer
			if err := tmpl.Execute(&promptBuffer, data); err != nil {
				return fmt.Errorf("프롬프트 템플릿 적용 실패: %v", err)
			}

			// 완성된 프롬프트 로깅
			fmt.Println("\n=== 생성된 프롬프트 ===")
			fmt.Println(promptBuffer.String())
			fmt.Println("=== 프롬프트 끝 ===\n")

			// AI 클라이언트 생성
			client, err := ai.NewClient(ai.ModelType(modelType))
			if err != nil {
				return fmt.Errorf("AI 클라이언트 초기화 실패: %v", err)
			}

			// 커밋 메시지 생성
			message, err := client.GenerateCommitMessage(promptBuffer.String())
			if err != nil {
				return fmt.Errorf("커밋 메시지 생성 실패: %v", err)
			}

			// AI 응답 로깅
			fmt.Println("\n=== AI 응답 ===")
			fmt.Println(message)
			fmt.Println("=== 응답 끝 ===\n")

			// 생성된 메시지 출력 및 확인
			fmt.Printf("\n생성된 커밋 메시지:\n%s\n\n", message)
			if !utils.ConfirmWithDefault("이 메시지로 커밋하시겠습니까?", false) {
				return nil
			}

			// 생성된 메시지로 커밋
			gitCmd = exec.Command("git", "commit", "-m", message)
			if err := gitCmd.Run(); err != nil {
				return fmt.Errorf("커밋 실행 실패: %v", err)
			}

			fmt.Println("\n커밋이 성공적으로 생성되었습니다.")
			return nil
		},
	}

	// 플래그 추가
	cmd.Flags().StringVarP(&keyword, "keyword", "k", "", "커밋 메시지 생성 시 참고할 키워드나 설명")
	cmd.Flags().BoolVarP(&withDiff, "verbose", "v", false, "전체 diff를 분석에 포함")
	cmd.Flags().StringVarP(&modelType, "type", "t", "claude", "사용할 AI 모델 (claude 또는 gpt)")

	return cmd
}

// 토큰 수 예측
func estimateTokenCount(diffs []string, files []string, keyword string) int {
	// 매우 간단한 추정: 평균적으로 영어 단어 1개가 1.3 토큰
	// 공백과 특수문자도 토큰으로 계산
	totalChars := 0
	for _, diff := range diffs {
		totalChars += len(diff)
	}
	for _, file := range files {
		totalChars += len(file)
	}
	totalChars += len(keyword)

	// 대략적인 토큰 수 추정
	return int(float64(totalChars) * 0.3)
}

// getDiffForAI 함수는 analyze.go로 이동되었습니다. 