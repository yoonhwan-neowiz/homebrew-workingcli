package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

// defaultJobs는 병렬 작업의 기본 개수입니다.
const defaultJobs = 4

func NewPullCmd() *cobra.Command {
	var (
		recursive bool
		jobs     int
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
  ga pull --rebase    # rebase 방식으로 가져오기
  ga pull -j 8        # 8개의 병렬 작업으로 서브모듈 업데이트 (기본값: 4)

성능 최적화:
- 서브모듈 병렬 업데이트 (기본 4개 작업)
- 필요한 경우에만 서브모듈 업데이트
- 네트워크 최적화 (--depth 1)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// jobs가 설정되지 않은 경우 기본값 사용
			if !cmd.Flags().Changed("jobs") {
				jobs = defaultJobs
			}

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

			// 3. 서브모듈 업데이트가 필요한지 확인
			submodules, err := getSubmodules()
			if err != nil {
				return err
			}

			if len(submodules) == 0 {
				return nil
			}

			// 4. 병렬로 서브모듈 업데이트
			if err := updateSubmodulesParallel(submodules, jobs, recursive); err != nil {
				return fmt.Errorf("서브모듈 업데이트 실패: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "서브모듈을 재귀적으로 업데이트")
	cmd.Flags().IntVarP(&jobs, "jobs", "j", defaultJobs, "병렬 작업 수 (서브모듈 업데이트)")
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

// 서브모듈 병렬 업데이트
func updateSubmodulesParallel(submodules []string, jobs int, recursive bool) error {
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

			// 서브모듈 업데이트 명령어 구성
			args := []string{"submodule", "update", "--init"}
			if recursive {
				args = append(args, "--recursive")
			}
			args = append(args, "--depth", "1", path) // 네트워크 최적화

			if err := execGitCommand(args...); err != nil {
				errChan <- fmt.Errorf("서브모듈 '%s' 업데이트 실패: %v", path, err)
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
		return fmt.Errorf("서브모듈 업데이트 중 오류 발생:\n%s", strings.Join(errors, "\n"))
	}

	return nil
} 