package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewToFullCmd creates the To FULL restoration command
func NewToFullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to-full",
		Short: "FULL 모드로 복원",
		Long: `저장소를 FULL 모드로 복원합니다.
모든 최적화를 해제하고 전체 히스토리와 파일을 다운로드합니다.
주의: 대량의 디스크 공간이 필요할 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runToFull()
		},
	}
}

func runToFull() {
	fmt.Println("🔄 FULL 모드로 복원 시작 (SLIM → FULL)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ 오류: Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. 현재 모드 확인
	mode := utils.GetOptimizationMode()
	fmt.Printf("📊 현재 모드: %s\n", mode)
	
	if mode == "FULL" {
		fmt.Println("✅ 이미 FULL 모드입니다. 추가 작업이 필요하지 않습니다.")
		return
	}
	
	// 3. 디스크 공간 확인
	fmt.Println("\n💾 디스크 공간 확인:")
	availableSpace := utils.GetAvailableDiskSpaceFormatted()
	fmt.Printf("   사용 가능한 공간: %s\n", availableSpace)
	
	// 4. 경고 및 확인
	fmt.Println("\n⚠️  주의사항:")
	fmt.Println("• 전체 히스토리와 모든 파일을 다운로드합니다")
	fmt.Println("• 대량의 디스크 공간이 필요할 수 있습니다")
	fmt.Println("• 네트워크 속도에 따라 시간이 오래 걸릴 수 있습니다")
	
	if !utils.ConfirmWithDefault("\n계속하시겠습니까?", false) {
		fmt.Println("❌ 취소되었습니다.")
		return
	}
	
	// 5. 디스크 사용량 (복원 전)
	beforeDisk := utils.GetDiskUsage()
	fmt.Printf("\n📂 복원 전 크기:\n")
	fmt.Printf("   .git 폴더: %s\n", beforeDisk[".git"])
	fmt.Printf("   프로젝트 전체: %s\n", beforeDisk["total"])
	
	// 6. 복원 프로세스 시작
	fmt.Println("\n🔧 복원 작업 시작...")
	
	// 6-1. Sparse Checkout 해제
	if utils.IsSparseCheckoutEnabled() {
		fmt.Println("\n[1/4] Sparse Checkout 해제...")
		cmd := exec.Command("git", "sparse-checkout", "disable")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("⚠️  Sparse Checkout 해제 중 경고: %v\n", err)
			fmt.Printf("   출력: %s\n", string(output))
		} else {
			fmt.Println("✅ Sparse Checkout 해제 완료")
		}
	} else {
		fmt.Println("\n[1/4] Sparse Checkout이 이미 비활성화되어 있습니다")
	}
	
	// 6-2. Shallow 저장소인 경우 전체 히스토리 다운로드
	fmt.Println("\n[2/4] 전체 히스토리 다운로드...")
	if utils.IsShallowRepository() {
		fmt.Println("   Shallow 저장소 감지 - unshallow 수행 중...")
		cmd := exec.Command("git", "fetch", "--unshallow")
		output, err := cmd.CombinedOutput()
		if err != nil {
			// unshallow가 이미 된 경우 에러가 발생할 수 있음
			outputStr := string(output)
			if !strings.Contains(outputStr, "already have") {
				fmt.Printf("⚠️  히스토리 복원 중 경고: %v\n", err)
				fmt.Printf("   출력: %s\n", outputStr)
			}
		} else {
			fmt.Println("✅ 전체 히스토리 다운로드 완료")
		}
	}
	
	// 6-3. 모든 객체 다운로드 (refetch)
	fmt.Println("\n[3/4] 모든 객체 다운로드 (refetch)...")
	fmt.Println("   이 작업은 시간이 오래 걸릴 수 있습니다...")
	cmd := exec.Command("git", "fetch", "--refetch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// refetch가 지원되지 않는 Git 버전의 경우
		outputStr := string(output)
		if strings.Contains(outputStr, "unknown option") || strings.Contains(outputStr, "unrecognized") {
			fmt.Println("⚠️  refetch가 지원되지 않는 Git 버전입니다. 대체 방법 사용...")
			// 대체 방법: 모든 remote를 다시 fetch
			cmd = exec.Command("git", "fetch", "--all", "--prune")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("⚠️  객체 다운로드 중 경고: %v\n", err)
				fmt.Printf("   출력: %s\n", string(output))
			} else {
				fmt.Println("✅ 객체 다운로드 완료 (fetch --all)")
			}
		} else {
			fmt.Printf("⚠️  객체 다운로드 중 경고: %v\n", err)
			fmt.Printf("   출력: %s\n", outputStr)
		}
	} else {
		fmt.Println("✅ 모든 객체 다운로드 완료")
	}
	
	// 6-4. Partial Clone 필터 제거
	fmt.Println("\n[4/4] Partial Clone 필터 제거...")
	partialFilter := utils.GetPartialCloneFilter()
	if partialFilter != "" {
		// Partial Clone 관련 설정 제거
		configs := [][]string{
			{"--unset", "remote.origin.partialclonefilter"},
			{"--unset", "remote.origin.promisor"},
			{"--unset", "extensions.partialClone"},
		}
		
		for _, config := range configs {
			cmd = exec.Command("git", append([]string{"config"}, config...)...)
			err = cmd.Run()
			if err != nil {
				// 설정이 없는 경우 에러가 발생할 수 있음 (무시 가능)
				continue
			}
		}
		fmt.Println("✅ Partial Clone 필터 제거 완료")
	} else {
		fmt.Println("   Partial Clone 필터가 설정되어 있지 않습니다")
	}
	
	// 7. 결과 확인
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 복원 결과 확인")
	
	// 최종 모드 확인
	finalMode := utils.GetOptimizationMode()
	fmt.Printf("\n✅ 최종 모드: %s\n", finalMode)
	
	// 최적화 상태 확인
	fmt.Println("\n🔍 최적화 상태:")
	fmt.Printf("   Partial Clone: %s\n", func() string {
		if utils.GetPartialCloneFilter() == "" {
			return "비활성"
		}
		return "활성"
	}())
	fmt.Printf("   Sparse Checkout: %s\n", func() string {
		if utils.IsSparseCheckoutEnabled() {
			return "활성"
		}
		return "비활성"
	}())
	fmt.Printf("   Shallow: %s\n", func() string {
		if utils.IsShallowRepository() {
			return "활성"
		}
		return "비활성"
	}())
	
	// 디스크 사용량 (복원 후)
	afterDisk := utils.GetDiskUsage()
	fmt.Printf("\n📂 복원 후 크기:\n")
	fmt.Printf("   .git 폴더: %s\n", afterDisk[".git"])
	fmt.Printf("   프로젝트 전체: %s\n", afterDisk["total"])
	
	// 완료 메시지
	if finalMode == "FULL" {
		fmt.Println("\n✅ FULL 모드로 복원이 완료되었습니다!")
		fmt.Println("   모든 파일과 전체 히스토리를 사용할 수 있습니다.")
	} else {
		fmt.Println("\n⚠️  일부 최적화가 여전히 활성화되어 있습니다.")
		fmt.Println("   'ga optimized quick status' 명령으로 상태를 확인하세요.")
	}
}