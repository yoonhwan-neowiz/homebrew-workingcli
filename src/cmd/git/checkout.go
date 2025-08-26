package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewCheckoutCmd() *cobra.Command {
	var force bool
	var cmd = &cobra.Command{
		Use:   "checkout <branch> | -- <path>...",
		Short: "브랜치 전환 또는 파일 복원",
		Long: `브랜치를 전환하거나 워킹 디렉토리의 파일을 복원합니다.

사용법:
  ga checkout <branch>          # 브랜치 전환
  ga checkout -b <new-branch>   # 새 브랜치 생성 및 전환
  ga checkout -- <path>...      # 파일 복원
  ga checkout -f <branch>       # 강제 브랜치 전환`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("브랜치 이름이나 파일 경로를 지정해주세요.")
				return
			}

			var gitArgs []string
			if force {
				gitArgs = append(gitArgs, "checkout", "-f")
			} else {
				gitArgs = append(gitArgs, "checkout")
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git checkout 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "강제 브랜치 전환")
	return cmd
} 