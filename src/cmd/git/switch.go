package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewSwitchCmd() *cobra.Command {
	var create bool
	var force bool
	var cmd = &cobra.Command{
		Use:   "switch <branch>",
		Short: "브랜치 전환",
		Long: `다른 브랜치로 전환합니다.
-c 옵션을 사용하면 새 브랜치를 생성하고 전환합니다.

사용법:
  ga switch <branch>          # 브랜치 전환
  ga switch -c <new-branch>   # 새 브랜치 생성 및 전환
  ga switch -f <branch>       # 강제 브랜치 전환`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("전환할 브랜치를 지정해주세요.")
				return
			}

			var gitArgs []string
			gitArgs = append(gitArgs, "switch")
			if create {
				gitArgs = append(gitArgs, "-c")
			}
			if force {
				gitArgs = append(gitArgs, "--force")
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git switch 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVarP(&create, "create", "c", false, "새 브랜치 생성 및 전환")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "강제 브랜치 전환")
	return cmd
} 