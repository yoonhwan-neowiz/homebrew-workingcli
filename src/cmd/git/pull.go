package git

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	var (
		recursive bool
		rebase   bool
	)

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "원격 저장소에서 변경사항 가져오기",
		Long: `원격 저장소에서 변경사항을 가져와 현재 브랜치에 병합합니다.
서브모듈이 있는 경우 자동으로 함께 업데이트합니다.

사용법:
  ga pull              # 현재 브랜치의 변경사항 가져오기
  ga pull origin main  # 특정 원격/브랜치에서 가져오기
  ga pull --rebase    # rebase 방식으로 가져오기`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. 메인 저장소 pull
			pullArgs := []string{"pull", "--no-edit"}
			if rebase {
				pullArgs = append(pullArgs, "--rebase")
			}
			pullArgs = append(pullArgs, args...)

			if err := execGitCommand(pullArgs...); err != nil {
				return fmt.Errorf("pull 실패: %v", err)
			}

			// 2. 서브모듈 확인
			hasSubmodules, err := checkSubmodules()
			if err != nil {
				return err
			}

			if !hasSubmodules {
				return nil
			}

			// 3. 서브모듈 업데이트
			fmt.Println("\n🔄 서브모듈 업데이트 중...")
			
			// 서브모듈 업데이트 명령어 구성
			updateArgs := []string{"submodule", "update", "--init"}
			if recursive {
				updateArgs = append(updateArgs, "--recursive")
			}

			if err := execGitCommand(updateArgs...); err != nil {
				return fmt.Errorf("서브모듈 업데이트 실패: %v", err)
			}
			
			fmt.Println("✅ 서브모듈 업데이트 완료")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "서브모듈을 재귀적으로 업데이트")
	cmd.Flags().BoolVar(&rebase, "rebase", false, "rebase 방식으로 가져오기")
	
	return cmd
}

// Git 명령어 실행 헬퍼 함수
func execGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 서브모듈 존재 여부 확인
func checkSubmodules() (bool, error) {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf(".gitmodules 확인 실패: %v", err)
	}
	return true, nil
}

 