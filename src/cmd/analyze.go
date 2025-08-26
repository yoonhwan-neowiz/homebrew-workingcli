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
)

const maxTokens = 4000  // 기본 토큰 제한

func NewAnalyzeCmd() *cobra.Command {
	var (
		last      int
		since     string
		until     string
		branch    string
		modelType string
		withDiff  bool
	)

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Git 커밋 내역 분석",
		Long: `Git 저장소의 커밋 내역을 분석하여 종합적인 요약을 생성합니다.
기간을 지정하여 특정 범위의 커밋만 분석할 수 있습니다.

사용법:
  ga analyze --last 10           # 최근 10개 커밋 분석
  ga analyze --since 1.week      # 1주일 이내 커밋 분석
  ga analyze --branch feature/*  # 특정 브랜치 패턴의 커밋 분석
  ga analyze -t gpt             # GPT 모델로 분석
  ga analyze -t claude          # Claude 모델로 분석
  ga analyze -v                 # 전체 diff 포함하여 분석`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Git 로그 명령어 구성
			gitArgs := []string{"log", "--pretty=format:%H%n%an%n%ad%n%s%n%b%n---"}
			if last > 0 {
				gitArgs = append(gitArgs, fmt.Sprintf("-%d", last))
			}
			if since != "" {
				gitArgs = append(gitArgs, "--since="+since)
			}
			if until != "" {
				gitArgs = append(gitArgs, "--until="+until)
			}
			if branch != "" {
				gitArgs = append(gitArgs, branch)
			}

			// Git 로그 실행
			output, err := exec.Command("git", gitArgs...).Output()
			if err != nil {
				return fmt.Errorf("Git 로그 가져오기 실패: %v", err)
			}

			// 커밋 정보 파싱
			commits := ParseCommits(string(output), withDiff)
			if len(commits) == 0 {
				return fmt.Errorf("분석할 커밋이 없습니다")
			}

			// diff 포함 여부에 따라 토큰 수 예상
			estimatedTokens := estimateTokens(commits)
			if estimatedTokens > maxTokens {
				fmt.Printf("\n토큰 수가 제한을 초과합니다.\n")
				fmt.Printf("현재 설정: %d 토큰\n", maxTokens)
				fmt.Printf("예상 사용: %d 토큰\n\n", estimatedTokens)
				
				if !utils.ConfirmWithDefault("계속 진행하시겠습니까?", false) {
					return fmt.Errorf("사용자가 작업을 취소했습니다")
				}
			}

			// 프롬프트 템플릿 가져오기
			promptPath, err := config.GetPromptPath("analyze")
			if err != nil {
				return fmt.Errorf("프롬프트 파일 경로 가져오기 실패: %v", err)
			}

			// 프롬프트 파일 읽기
			promptContent, err := os.ReadFile(promptPath)
			if err != nil {
				return fmt.Errorf("프롬프트 파일 읽기 실패: %v", err)
			}

			// 템플릿 데이터 준비
			data := map[string]interface{}{
				"Commits": commits,
			}

			// 템플릿 적용
			tmpl, err := template.New("analyze").Parse(string(promptContent))
			if err != nil {
				return fmt.Errorf("프롬프트 템플릿 파싱 실패: %v", err)
			}

			// 프롬프트 템플릿 적용
			var promptBuffer bytes.Buffer
			if err := tmpl.Execute(&promptBuffer, data); err != nil {
				return fmt.Errorf("프롬프트 템플릿 적용 실패: %v", err)
			}

			// 완성된 프롬프트 로깅
			fmt.Println("\n=== 생성된 분석 프롬프트 ===")
			fmt.Println(promptBuffer.String())
			fmt.Println("=== 프롬프트 끝 ===\n")

			// AI 클라이언트 생성
			client, err := ai.NewClient(ai.ModelType(modelType))
			if err != nil {
				return fmt.Errorf("AI 클라이언트 초기화 실패: %v", err)
			}

			// 커밋 분석 수행
			analysis, err := client.AnalyzeCommits(promptBuffer.String())
			if err != nil {
				return fmt.Errorf("커밋 분석 실패: %v", err)
			}

			// AI 응답 로깅
			fmt.Println("\n=== AI 분석 결과 ===")
			fmt.Println(analysis)
			fmt.Println("=== 분석 끝 ===\n")

			// 분석 결과 출력
			fmt.Printf("\n분석 결과:\n%s\n", analysis)

			return nil
		},
	}

	// 플래그 추가
	cmd.Flags().IntVarP(&last, "last", "n", 0, "분석할 최근 커밋 수")
	cmd.Flags().StringVarP(&since, "since", "s", "", "시작 시점 (예: 1.week, 2.days)")
	cmd.Flags().StringVarP(&until, "until", "u", "", "종료 시점")
	cmd.Flags().StringVarP(&branch, "branch", "b", "", "분석할 브랜치 패턴")
	cmd.Flags().StringVarP(&modelType, "type", "t", "claude", "사용할 AI 모델 (claude 또는 gpt)")
	cmd.Flags().BoolVarP(&withDiff, "verbose", "v", false, "전체 diff를 분석에 포함")

	return cmd
}

// ParseCommits는 Git 로그 출력을 파싱하여 커밋 정보 구조체로 변환합니다.
func ParseCommits(output string, withDiff bool) []ai.CommitInfo {
	var commits []ai.CommitInfo
	sections := strings.Split(output, "\n---\n")

	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		lines := strings.Split(strings.TrimSpace(section), "\n")
		if len(lines) < 4 {
			continue
		}

		commit := ai.CommitInfo{
			Hash:    lines[0],
			Author:  lines[1],
			Date:    lines[2],
			Message: lines[3],
		}

		// 커밋 설명이 있는 경우
		if len(lines) > 4 {
			commit.Description = strings.Join(lines[4:], "\n")
		}

		// 변경된 파일 목록 가져오기
		if files, err := getCommitFiles(commit.Hash); err == nil {
			commit.Files = files
		} else {
			fmt.Printf("경고: 커밋 %s의 변경 파일 목록 가져오기 실패: %v\n", commit.Hash, err)
			continue
		}

		// diff 가져오기
		fmt.Println("\n=== diff 가져오기 시작 ===")
		var diffContent string
		if diff, err := utils.GetDiffForAI(commit.Files, commit.Hash, withDiff); err == nil {
			diffContent = diff
			commit.Diff = diffContent
			fmt.Println("diff 가져오기 성공")
		} else {
			fmt.Printf("diff 가져오기 실패: %v\n", err)
			return nil
		}
		fmt.Println("=== diff 가져오기 끝 ===\n")

		commits = append(commits, commit)
	}

	return commits
}

// getCommitFiles는 특정 커밋의 변경된 파일 목록을 가져옵니다.
func getCommitFiles(hash string) ([]string, error) {
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", hash)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return utils.ProcessGitPaths(files), nil
}

// estimateTokens는 커밋 정보를 기반으로 예상 토큰 수를 계산합니다.
func estimateTokens(commits []ai.CommitInfo) int {
	// 기본 토큰 (메타데이터, 구조 등)
	tokens := 500

	for _, commit := range commits {
		// 커밋 메시지 및 기본 정보
		tokens += len(commit.Message) / 4
		tokens += len(commit.Description) / 4
		tokens += len(commit.Files) * 10
		if commit.Diff != "" {
			tokens += len(commit.Diff) / 4
		}
	}

	return tokens
}
