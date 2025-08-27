package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewUnshallowAllCmd creates the Unshallow All Submodules command
func NewUnshallowAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unshallow-all",
		Short: "모든 서브모듈 히스토리 복원",
		Long: `모든 서브모듈의 전체 히스토리를 복원합니다.
Shallow 상태의 서브모듈을 완전한 저장소로 변환합니다.

실행 내용:
- git submodule foreach 'git fetch --unshallow'
- 모든 서브모듈의 전체 커밋 히스토리 다운로드

⚠️ 주의: 서브모듈 크기에 따라 상당한 디스크 공간이 필요할 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runUnshallowAll()
		},
	}
}

func runUnshallowAll() {
	// 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}

	fmt.Println("🔄 모든 서브모듈의 전체 히스토리를 복원합니다...")
	fmt.Println("⚠️ 주의: 대용량 저장소의 경우 시간이 오래 걸릴 수 있습니다.\n")

	// 사용자 확인
	if !utils.ConfirmWithDefault("계속하시겠습니까?", true) {
		fmt.Println("❌ 작업이 취소되었습니다.")
		return
	}

	fmt.Println()
	fmt.Printf("📦 총 %d개의 서브모듈을 병렬로 처리합니다.\n\n", count)

	// Unshallow 작업 정의
	unshallowOperation := func(path string) error {
		// 서브모듈 디렉토리로 이동
		originalDir, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("디렉토리 이동 실패: %v", err)
		}
		defer os.Chdir(originalDir)

		// 현재 shallow 상태 확인
		isShallowCmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
		output, _ := isShallowCmd.Output()
		isShallow := strings.TrimSpace(string(output)) == "true"

		if !isShallow {
			fmt.Printf("ℹ️ %s: 이미 전체 히스토리를 가지고 있습니다\n", path)
			return nil // 성공으로 처리
		}

		// 현재 depth 확인
		countCmd := exec.Command("git", "rev-list", "--count", "HEAD")
		countOutput, _ := countCmd.Output()
		currentDepth := strings.TrimSpace(string(countOutput))
		fmt.Printf("📊 %s: Shallow 상태 (depth: %s) → 전체 히스토리 다운로드 중...\n", path, currentDepth)
		
		// unshallow 실행
		fetchCmd := exec.Command("git", "fetch", "--unshallow")
		if err := fetchCmd.Run(); err != nil {
			// 실패 시 다른 방법 시도
			fetchAllCmd := exec.Command("git", "fetch", "--all")
			if err := fetchAllCmd.Run(); err != nil {
				return fmt.Errorf("히스토리 복원 실패: %v", err)
			}
		}
		
		// 복원 후 커밋 수 확인
		countCmd = exec.Command("git", "rev-list", "--count", "HEAD")
		countOutput, _ = countCmd.Output()
		totalCommits := strings.TrimSpace(string(countOutput))
		
		fmt.Printf("✅ %s: 전체 히스토리 복원 완료 (총 %s개 커밋)\n", path, totalCommits)
		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	successCount, failCount, err := utils.ExecuteOnSubmodulesParallel(unshallowOperation, 4, true)

	// 결과 요약
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 작업 완료 요약")
	fmt.Printf("✅ 성공: %d개 서브모듈\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", failCount)
	}
	
	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}
	
	fmt.Println("\n모든 서브모듈이 전체 히스토리를 가지게 되었습니다.")
}