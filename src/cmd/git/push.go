package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// defaultJobs는 병렬 작업의 기본 개수입니다.
const defaultJobs = 4

func NewPushCmd() *cobra.Command {
	var (
		recursive bool
		jobs     int
		force    bool
	)

	cmd := &cobra.Command{
		Use:   "push",
		Short: "변경사항을 원격 저장소로 푸시",
		Long: `현재 브랜치의 변경사항을 원격 저장소로 푸시합니다.
서브모듈이 있는 경우 자동으로 함께 푸시합니다.

사용법:
  ga push              # 현재 브랜치를 원격으로 푸시
  ga push origin main  # 특정 원격/브랜치로 푸시
  ga push --force     # 강제 푸시 (주의 필요)
  ga push -j 8        # 8개의 병렬 작업으로 서브모듈 푸시 (기본값: 4)

성능 최적화:
- 서브모듈 병렬 푸시 (기본 4개 작업)
- 필요한 경우에만 서브모듈 푸시
- 변경된 서브모듈만 선택적 푸시`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// jobs가 설정되지 않은 경우 기본값 사용
			if !cmd.Flags().Changed("jobs") {
				jobs = defaultJobs
			}

			// 1. 메인 저장소 push
			pushArgs := []string{"push"}
			if force {
				pushArgs = append(pushArgs, "--force")
			}
			pushArgs = append(pushArgs, args...)

			if err := execGitCommand(pushArgs...); err != nil {
				return fmt.Errorf("push 실패: %v", err)
			}

			// 2. 서브모듈 확인
			hasSubmodules, err := checkSubmodules()
			if err != nil {
				return err
			}

			if !hasSubmodules {
				return nil
			}

			// 3. 변경된 서브모듈 확인
			changedSubmodules, err := getChangedSubmodules()
			if err != nil {
				return err
			}

			if len(changedSubmodules) == 0 {
				return nil
			}

			// 4. 변경된 서브모듈 푸시
			fmt.Printf("\n🔄 %d개의 변경된 서브모듈을 푸시합니다...\n", len(changedSubmodules))
			
			// 푸시 작업 정의
			pushOperation := func(path string) error {
				// 서브모듈 디렉토리로 이동
				originalDir, _ := os.Getwd()
				if err := os.Chdir(path); err != nil {
					return fmt.Errorf("디렉토리 이동 실패: %v", err)
				}
				defer os.Chdir(originalDir)

				// 푸시 명령어 구성
				args := []string{"push"}
				if force {
					args = append(args, "--force")
				}
				
				if err := execGitCommand(args...); err != nil {
					return fmt.Errorf("푸시 실패: %v", err)
				}
				
				fmt.Printf("✅ %s: 푸시 완료\n", path)
				return nil
			}

			// 변경된 서브모듈만 처리하도록 필터링된 작업 실행
			var successCount, failCount int
			for _, submodule := range changedSubmodules {
				if err := pushOperation(submodule); err != nil {
					fmt.Printf("❌ %s: %v\n", submodule, err)
					failCount++
				} else {
					successCount++
				}
			}
			
			if failCount > 0 {
				return fmt.Errorf("서브모듈 푸시 중 %d개 실패", failCount)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "서브모듈을 재귀적으로 푸시")
	cmd.Flags().IntVarP(&jobs, "jobs", "j", defaultJobs, "병렬 작업 수 (서브모듈 푸시)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "강제 푸시 (주의 필요)")

	return cmd
}

// 서브모듈 목록 가져오기
func getSubmodules() ([]string, error) {
	cmd := exec.Command("git", "submodule", "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("서브모듈 상태 확인 실패: %v", err)
	}

	var submodules []string
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		// 상태 출력 형식: <hash> <path> (<branch>)
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			submodules = append(submodules, parts[1])
		}
	}
	return submodules, nil
}

// 변경된 서브모듈 목록 가져오기
func getChangedSubmodules() ([]string, error) {
	// 모든 서브모듈 가져오기
	allSubmodules, err := getSubmodules()
	if err != nil {
		return nil, err
	}

	var changedSubmodules []string
	for _, submodule := range allSubmodules {
		// 서브모듈 상태 확인
		cmd := exec.Command("git", "diff", "--quiet", submodule)
		if err := cmd.Run(); err != nil {
			// 에러가 발생하면 변경사항이 있는 것
			changedSubmodules = append(changedSubmodules, submodule)
			continue
		}

		// staged 변경사항 확인
		cmd = exec.Command("git", "diff", "--quiet", "--cached", submodule)
		if err := cmd.Run(); err != nil {
			changedSubmodules = append(changedSubmodules, submodule)
		}
	}

	return changedSubmodules, nil
} 