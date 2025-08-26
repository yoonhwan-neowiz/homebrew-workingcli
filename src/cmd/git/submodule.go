package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewSubmoduleCmd() *cobra.Command {
	var (
		recursive bool
	)

	cmd := &cobra.Command{
		Use:     "submodule [command]",
		Aliases: []string{"sub"},  // ga sub로도 사용 가능
		Short:   "서브모듈 관리",
		Long: `Git 서브모듈을 관리합니다.
모든 서브모듈에서 명령어를 실행하거나 특정 서브모듈을 관리할 수 있습니다.

사용 예:
  ga submodule "git pull"     # 모든 서브모듈에서 git pull 실행
  ga submodule -r "git pull"  # 재귀적으로 모든 서브모듈에서 git pull 실행
  ga sub -r "git status"      # 별칭 사용 (재귀적으로 상태 확인)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 인자가 없으면 git submodule 기본 명령 실행
			if len(args) == 0 {
				gitCmd := exec.Command("git", "submodule")
				gitCmd.Stdout = os.Stdout
				gitCmd.Stderr = os.Stderr
				return gitCmd.Run()
			}

			// foreach를 위한 명령어가 있는 경우
			gitArgs := []string{"submodule", "foreach"}
			if recursive {
				gitArgs = append(gitArgs, "--recursive")
			}

			// 명령어를 그대로 전달
			command := strings.Join(args, " ")
			gitArgs = append(gitArgs, command)

			// git submodule foreach 실행
			gitCmd := exec.Command("git", gitArgs...)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			gitCmd.Stdin = os.Stdin
			
			if err := gitCmd.Run(); err != nil {
				return fmt.Errorf("git submodule foreach 실행 실패: %v", err)
			}

			return nil
		},
	}

	// 플래그 추가
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "중첩된 서브모듈에서도 명령어 실행")

	return cmd
}