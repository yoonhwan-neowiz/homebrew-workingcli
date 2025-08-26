package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewFetchCmd() *cobra.Command {
	var all bool
	var prune bool
	var tags bool
	var cmd = &cobra.Command{
		Use:   "fetch [<remote>]",
		Short: "원격 저장소의 변경사항 가져오기",
		Long: `원격 저장소의 변경사항을 가져옵니다.
--all 옵션을 사용하면 모든 원격 저장소의 변경사항을 가져옵니다.

사용법:
  ga fetch                # 기본 원격 저장소의 변경사항 가져오기
  ga fetch <remote>      # 지정한 원격 저장소의 변경사항 가져오기
  ga fetch --all         # 모든 원격 저장소의 변경사항 가져오기
  ga fetch --prune       # 원격에서 삭제된 브랜치 정리
  ga fetch --tags        # 모든 태그 가져오기`,
		Run: func(cmd *cobra.Command, args []string) {
			var gitArgs []string
			gitArgs = append(gitArgs, "fetch")
			if all {
				gitArgs = append(gitArgs, "--all")
			}
			if prune {
				gitArgs = append(gitArgs, "--prune")
			}
			if tags {
				gitArgs = append(gitArgs, "--tags")
			}
			gitArgs = append(gitArgs, args...)

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git fetch 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "모든 원격 저장소의 변경사항 가져오기")
	cmd.Flags().BoolVar(&prune, "prune", false, "원격에서 삭제된 브랜치 정리")
	cmd.Flags().BoolVar(&tags, "tags", false, "모든 태그 가져오기")
	return cmd
} 