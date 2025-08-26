package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewRebaseCmd() *cobra.Command {
	var interactive bool
	var onto string
	var cmd = &cobra.Command{
		Use:   "rebase [<upstream>] [<branch>]",
		Short: "브랜치 리베이스",
		Long: `현재 브랜치를 다른 브랜치 위로 리베이스합니다.
--interactive 옵션을 사용하면 대화형 리베이스를 시작할 수 있습니다.

사용법:
  ga rebase <upstream>                    # 브랜치 리베이스
  ga rebase --interactive <upstream>      # 대화형 리베이스
  ga rebase --onto <newbase> <upstream>   # 특정 커밋으로 리베이스`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("리베이스할 브랜치를 지정해주세요.")
				return
			}

			var gitArgs []string
			gitArgs = append(gitArgs, "rebase")
			if interactive {
				gitArgs = append(gitArgs, "--interactive")
			}
			if onto != "" {
				gitArgs = append(gitArgs, "--onto", onto)
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git rebase 실행 중 오류:", err)
				fmt.Println("\n충돌이 발생한 경우:")
				fmt.Println("1. 'ga resolve' 명령어로 충돌을 해결하세요.")
				fmt.Println("2. 'git add' 명령어로 해결된 파일을 스테이징하세요.")
				fmt.Println("3. 'git rebase --continue'로 리베이스를 계속하세요.")
				fmt.Println("또는 'git rebase --abort'로 리베이스를 취소할 수 있습니다.")
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "대화형 리베이스")
	cmd.Flags().StringVar(&onto, "onto", "", "특정 커밋으로 리베이스")
	return cmd
} 