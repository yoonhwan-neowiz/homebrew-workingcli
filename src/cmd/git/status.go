package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Git 저장소 상태 확인",
		Long: `현재 Git 저장소의 상태를 확인합니다.
staged/unstaged 파일 목록과 브랜치 정보를 표시합니다.

사용법:
  ga status`,
		Run: func(cmd *cobra.Command, args []string) {
			var out bytes.Buffer
			gitCmd := exec.Command("git", "status")
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git status 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}
} 