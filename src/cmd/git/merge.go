package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewMergeCmd() *cobra.Command {
	var noCommit bool
	var noFf bool
	var squash bool
	var cmd = &cobra.Command{
		Use:   "merge <branch>",
		Short: "브랜치 병합",
		Long: `현재 브랜치에 지정한 브랜치를 병합합니다.
충돌이 발생하면 resolve 명령어를 사용하여 해결할 수 있습니다.

사용법:
  ga merge <branch>                # 브랜치 병합
  ga merge --no-commit <branch>    # 병합 커밋 생성하지 않음
  ga merge --no-ff <branch>        # fast-forward 하지 않음
  ga merge --squash <branch>       # 스쿼시 병합`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("병합할 브랜치를 지정해주세요.")
				return
			}

			var gitArgs []string
			gitArgs = append(gitArgs, "merge")
			if noCommit {
				gitArgs = append(gitArgs, "--no-commit")
			}
			if noFf {
				gitArgs = append(gitArgs, "--no-ff")
			}
			if squash {
				gitArgs = append(gitArgs, "--squash")
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git merge 실행 중 오류:", err)
				fmt.Println("\n충돌이 발생한 경우 'ga resolve' 명령어를 사용하여 해결할 수 있습니다.")
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVar(&noCommit, "no-commit", false, "병합 커밋 생성하지 않음")
	cmd.Flags().BoolVar(&noFf, "no-ff", false, "fast-forward 하지 않음")
	cmd.Flags().BoolVar(&squash, "squash", false, "스쿼시 병합")
	return cmd
} 