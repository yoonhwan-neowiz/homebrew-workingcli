package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

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

			// 4. 병렬로 서브모듈 푸시
			if err := pushSubmodulesParallel(changedSubmodules, jobs, recursive, force); err != nil {
				return fmt.Errorf("서브모듈 푸시 실패: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "서브모듈을 재귀적으로 푸시")
	cmd.Flags().IntVarP(&jobs, "jobs", "j", defaultJobs, "병렬 작업 수 (서브모듈 푸시)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "강제 푸시 (주의 필요)")

	return cmd
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

// 서브모듈 병렬 푸시
func pushSubmodulesParallel(submodules []string, jobs int, recursive bool, force bool) error {
	if jobs < 1 {
		jobs = 1
	}

	// 작업 풀 생성
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, jobs)
	errChan := make(chan error, len(submodules))

	for _, submodule := range submodules {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			semaphore <- struct{}{} // 작업 슬롯 획득
			defer func() { <-semaphore }() // 작업 슬롯 반환

			// 서브모듈 디렉토리로 이동
			if err := os.Chdir(path); err != nil {
				errChan <- fmt.Errorf("서브모듈 '%s' 디렉토리 이동 실패: %v", path, err)
				return
			}
			defer os.Chdir("..") // 원래 디렉토리로 복귀

			// 푸시 명령어 구성
			args := []string{"push"}
			if force {
				args = append(args, "--force")
			}
			if recursive {
				// 재귀적 푸시를 위한 추가 처리
				if err := pushRecursive(args, force); err != nil {
					errChan <- fmt.Errorf("서브모듈 '%s' 재귀적 푸시 실패: %v", path, err)
					return
				}
			} else {
				// 일반 푸시
				if err := execGitCommand(args...); err != nil {
					errChan <- fmt.Errorf("서브모듈 '%s' 푸시 실패: %v", path, err)
					return
				}
			}
		}(submodule)
	}

	// 모든 작업 완료 대기
	wg.Wait()
	close(errChan)

	// 에러 수집
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("서브모듈 푸시 중 오류 발생:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// 재귀적 푸시 처리
func pushRecursive(args []string, force bool) error {
	// 현재 저장소 푸시
	if err := execGitCommand(args...); err != nil {
		return err
	}

	// 하위 서브모듈 확인
	submodules, err := getSubmodules()
	if err != nil {
		return err
	}

	// 각 서브모듈에 대해 재귀적으로 푸시
	for _, submodule := range submodules {
		if err := os.Chdir(submodule); err != nil {
			return fmt.Errorf("서브모듈 '%s' 디렉토리 이동 실패: %v", submodule, err)
		}
		if err := pushRecursive(args, force); err != nil {
			return fmt.Errorf("서브모듈 '%s' 재귀적 푸시 실패: %v", submodule, err)
		}
		if err := os.Chdir(".."); err != nil {
			return fmt.Errorf("상위 디렉토리 이동 실패: %v", err)
		}
	}

	return nil
} 