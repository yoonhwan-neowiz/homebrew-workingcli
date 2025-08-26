package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewDiffCmd() *cobra.Command {
	var staged bool
	var cmd = &cobra.Command{
		Use:   "diff [<path>...]",
		Short: "변경사항 확인",
		Long: `워킹 디렉토리의 변경사항을 확인합니다.
--staged 옵션을 사용하면 스테이징된 변경사항을 확인할 수 있습니다.

사용법:
  ga diff [<path>...]
  ga diff --staged [<path>...]`,
		Run: func(cmd *cobra.Command, args []string) {
			var gitArgs []string
			if staged {
				gitArgs = append(gitArgs, "diff", "--staged")
			} else {
				gitArgs = append(gitArgs, "diff")
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git diff 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVar(&staged, "staged", false, "스테이징된 변경사항 확인")
	return cmd
} 