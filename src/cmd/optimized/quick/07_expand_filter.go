package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewExpandFilterCmd creates the Expand Filter removal command
func NewExpandFilterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-filter",
		Short: "Partial Clone 필터 제거 (모든 파일 다운로드)",
		Long: `Partial Clone 필터를 완전히 제거하여 모든 대용량 파일을 다운로드합니다.
Sparse Checkout은 유지하면서 blob:limit 필터만 해제하여 
현재 checkout된 경로의 모든 파일을 크기 제한 없이 다운로드합니다.
주의: 디스크 공간을 많이 사용할 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			executeExpandFilter()
		},
	}
}

func executeExpandFilter() {
	fmt.Println("🔄 Partial Clone 필터 제거 프로세스 시작")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ 오류: 현재 디렉토리는 Git 저장소가 아닙니다.")
		fmt.Println("   Git 저장소 루트에서 실행해주세요.")
		os.Exit(1)
	}
	
	// 1. 현재 필터 확인
	fmt.Println("\n1️⃣  현재 Partial Clone 필터 확인 중...")
	currentFilter := utils.GetPartialCloneFilter()
	
	if currentFilter == "" {
		fmt.Println("ℹ️  Partial Clone 필터가 설정되어 있지 않습니다.")
		fmt.Println("   이미 모든 파일이 다운로드된 상태입니다.")
		
		// Sparse Checkout 상태도 함께 표시
		sparseInfo := utils.GetSparseCheckoutInfo()
		if sparseInfo["enabled"].(bool) {
			fmt.Printf("\n📁 Sparse Checkout: 활성 (%d개 경로)\n", sparseInfo["count"])
			if paths, ok := sparseInfo["paths"].([]string); ok && len(paths) > 0 {
				fmt.Println("   설정된 경로:")
				for _, path := range paths {
					fmt.Printf("   • %s\n", path)
				}
			}
		}
		return
	}
	
	fmt.Printf("   현재 필터: %s\n", currentFilter)
	
	// 필터 크기 파싱하여 예상 다운로드 크기 안내
	var filterSize string
	if strings.Contains(currentFilter, "blob:limit=") {
		parts := strings.Split(currentFilter, "=")
		if len(parts) > 1 {
			filterSize = parts[1]
			fmt.Printf("   제외된 파일: %s 이상의 blob\n", filterSize)
		}
	}
	
	// 디스크 사용량 확인
	diskUsage := utils.GetDiskUsage()
	fmt.Printf("\n💾 현재 디스크 사용량:\n")
	if gitSize, ok := diskUsage[".git"]; ok {
		fmt.Printf("   .git 폴더: %s\n", gitSize)
	}
	if projectSize, ok := diskUsage["."]; ok {
		fmt.Printf("   프로젝트 전체: %s\n", projectSize)
	}
	
	// 사용자 확인
	fmt.Println("\n⚠️  경고: 필터를 제거하면 모든 대용량 파일이 다운로드됩니다.")
	fmt.Println("   이 작업은 상당한 디스크 공간과 네트워크 대역폭을 사용할 수 있습니다.")
	
	if !utils.ConfirmWithDefault("정말로 Partial Clone 필터를 제거하시겠습니까?", false) {
		fmt.Println("\n❌ 작업이 취소되었습니다.")
		return
	}
	
	// 2. 필터 제거
	fmt.Println("\n2️⃣  Partial Clone 필터 제거 중...")
	
	// remote.origin.partialclonefilter 제거
	cmd := exec.Command("git", "config", "--unset", "remote.origin.partialclonefilter")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "no such section") {
		// 설정이 없는 경우가 아닌 실제 오류만 처리
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 5 {
			fmt.Printf("❌ 필터 제거 실패: %v\n", err)
			if len(output) > 0 {
				fmt.Printf("   상세: %s\n", string(output))
			}
			os.Exit(1)
		}
	}
	fmt.Println("   ✓ 필터 설정 제거 완료")
	
	// 3. 모든 blob 다운로드
	fmt.Println("\n3️⃣  모든 파일 다운로드 중...")
	fmt.Println("   (네트워크 상황에 따라 시간이 걸릴 수 있습니다)")
	
	// git fetch --refetch 실행
	fetchCmd := exec.Command("git", "fetch", "--refetch")
	fetchCmd.Stdout = os.Stdout
	fetchCmd.Stderr = os.Stderr
	
	err = fetchCmd.Run()
	if err != nil {
		fmt.Printf("❌ 파일 다운로드 실패: %v\n", err)
		
		// 실패 시 필터 복원 제안
		fmt.Printf("\n💡 필터를 다시 설정하려면:\n")
		fmt.Printf("   git config remote.origin.partialclonefilter %s\n", currentFilter)
		os.Exit(1)
	}
	
	// 4. 결과 확인
	fmt.Println("\n4️⃣  작업 결과 확인")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 필터 제거 확인
	newFilter := utils.GetPartialCloneFilter()
	if newFilter == "" {
		fmt.Println("✅ Partial Clone 필터 제거 완료")
	} else {
		fmt.Printf("⚠️  필터가 여전히 설정됨: %s\n", newFilter)
	}
	
	// Sparse Checkout 상태
	sparseInfo := utils.GetSparseCheckoutInfo()
	if sparseInfo["enabled"].(bool) {
		fmt.Printf("📁 Sparse Checkout: 활성 유지 (%d개 경로)\n", sparseInfo["count"])
	}
	
	// 새로운 디스크 사용량
	newDiskUsage := utils.GetDiskUsage()
	fmt.Printf("\n💾 변경 후 디스크 사용량:\n")
	if gitSize, ok := newDiskUsage[".git"]; ok {
		fmt.Printf("   .git 폴더: %s", gitSize)
		if oldSize, ok := diskUsage[".git"]; ok && oldSize != gitSize {
			fmt.Printf(" (변경 전: %s)", oldSize)
		}
		fmt.Println()
	}
	if projectSize, ok := newDiskUsage["."]; ok {
		fmt.Printf("   프로젝트 전체: %s", projectSize)
		if oldSize, ok := diskUsage["."]; ok && oldSize != projectSize {
			fmt.Printf(" (변경 전: %s)", oldSize)
		}
		fmt.Println()
	}
	
	// 최종 안내
	fmt.Println("\n✅ Partial Clone 필터 제거가 완료되었습니다!")
	fmt.Println("   모든 파일이 크기 제한 없이 다운로드되었습니다.")
	
	if sparseInfo["enabled"].(bool) {
		fmt.Println("\n💡 Sparse Checkout은 여전히 활성 상태입니다.")
		fmt.Println("   전체 파일을 작업 트리에 체크아웃하려면:")
		fmt.Println("   ga opt quick to-full")
	}
}