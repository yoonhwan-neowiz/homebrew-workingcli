package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

func NewTagCmd() *cobra.Command {
	var annotated bool
	var message string
	var delete bool
	var cmd = &cobra.Command{
		Use:   "tag [<tagname>] [<commit>]",
		Short: "태그 관리",
		Long: `태그를 생성, 삭제하거나 목록을 조회합니다.
-a 옵션을 사용하면 주석이 있는 태그를 생성합니다.

사용법:
  ga tag                      # 태그 목록 조회
  ga tag <tagname>            # 태그 생성
  ga tag -a <tagname> -m <msg> # 주석이 있는 태그 생성
  ga tag -d <tagname>         # 태그 삭제`,
		Run: func(cmd *cobra.Command, args []string) {
			var gitArgs []string
			gitArgs = append(gitArgs, "tag")

			if delete {
				if len(args) == 0 {
					fmt.Println("삭제할 태그를 지정해주세요.")
					return
				}
				gitArgs = append(gitArgs, "-d")
				gitArgs = append(gitArgs, args...)
			} else if len(args) > 0 {
				if annotated {
					gitArgs = append(gitArgs, "-a")
					if message != "" {
						gitArgs = append(gitArgs, "-m", message)
					}
				}
				gitArgs = append(gitArgs, args...)
			}

			var out bytes.Buffer
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = &out
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				fmt.Println("Git tag 실행 중 오류:", err)
				return
			}

			// Git 출력의 한글 파일명 처리
			output := utils.ProcessGitOutput(out.String())
			fmt.Print(output)
		},
	}

	cmd.Flags().BoolVarP(&annotated, "annotate", "a", false, "주석이 있는 태그 생성")
	cmd.Flags().StringVarP(&message, "message", "m", "", "태그 메시지")
	cmd.Flags().BoolVarP(&delete, "delete", "d", false, "태그 삭제")
	return cmd
} 