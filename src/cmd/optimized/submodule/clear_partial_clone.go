package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewExpandFilterCmd creates the submodule Expand Filter removal command
func NewExpandFilterCmd() *cobra.Command {
	var jobs int
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:   "expand-filter",
		Short: "서브모듈 Partial Clone 필터 제거",
		Long: `모든 서브모듈의 Partial Clone 필터를 제거하여 대용량 파일을 포함한 
모든 파일을 다운로드합니다. Sparse Checkout은 유지됩니다.

이 작업은 각 서브모듈의 디스크 사용량을 크게 증가시킬 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			executeSubmoduleExpandFilter(jobs)
		},
	}
	
	cmd.Flags().IntVar(&jobs, "jobs", 4, "병렬 처리할 작업 수")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func executeSubmoduleExpandFilter(jobs int) {
	fmt.Println("🔄 서브모듈 Partial Clone 필터 제거 프로세스")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 서브모듈 존재 확인
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		fmt.Println("ℹ️  서브모듈이 없습니다.")
		return
	}
	
	// 서브모듈 목록 가져오기
	cmd := exec.Command("git", "submodule", "foreach", "--quiet", "echo $path")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ 서브모듈 목록을 가져올 수 없습니다: %v\n", err)
		os.Exit(1)
	}
	
	submodulePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(submodulePaths) == 0 || (len(submodulePaths) == 1 && submodulePaths[0] == "") {
		fmt.Println("ℹ️  초기화된 서브모듈이 없습니다.")
		return
	}
	
	fmt.Printf("\n📊 대상 서브모듈: %d개\n", len(submodulePaths))
	
	// 각 서브모듈의 현재 필터 상태 확인
	fmt.Println("\n1️⃣  서브모듈 Partial Clone 필터 상태 확인")
	fmt.Println("─────────────────────────────────────────")
	
	type submoduleInfo struct {
		path         string
		filter       string
		diskUsageBefore string
		bytesBefore  int64
	}
	
	var submodules []submoduleInfo
	var totalFilteredCount int
	
	for _, path := range submodulePaths {
		if path == "" {
			continue
		}
		
		// 서브모듈 디렉토리로 이동하여 필터 확인
		filterCmd := exec.Command("git", "-C", path, "config", "remote.origin.partialclonefilter")
		filterOutput, _ := filterCmd.Output()
		filter := strings.TrimSpace(string(filterOutput))
		
		// 디스크 사용량 확인
		bytes, human := utils.GetSubmoduleGitSize(path)
		
		info := submoduleInfo{
			path:            path,
			filter:          filter,
			diskUsageBefore: human,
			bytesBefore:     bytes,
		}
		submodules = append(submodules, info)
		
		if filter != "" {
			totalFilteredCount++
			fmt.Printf("   📁 %s\n", path)
			fmt.Printf("      필터: %s\n", filter)
			fmt.Printf("      크기: %s\n", human)
		}
	}
	
	if totalFilteredCount == 0 {
		fmt.Println("   ✅ 모든 서브모듈이 이미 전체 다운로드 상태입니다.")
		return
	}
	
	// 사용자 확인
	fmt.Printf("\n⚠️  경고: %d개 서브모듈의 필터를 제거합니다.\n", totalFilteredCount)
	fmt.Println("   모든 대용량 파일이 다운로드되어 디스크 사용량이 크게 증가할 수 있습니다.")
	
	// Partial Clone 필터 제거는 안전한 작업이므로 quiet 모드에서 자동 수락
	if !utils.ConfirmForce("계속하시겠습니까?") {
		fmt.Println("\n❌ 작업이 취소되었습니다.")
		return
	}
	
	// 병렬로 필터 제거 및 fetch 실행
	fmt.Printf("\n2️⃣  서브모듈 필터 제거 중 (병렬 작업: %d)\n", jobs)
	fmt.Println("─────────────────────────────────────────")
	
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, jobs)
	resultChan := make(chan string, len(submodules))
	errorChan := make(chan error, len(submodules))
	
	for _, sm := range submodules {
		if sm.filter == "" {
			continue // 이미 필터가 없는 서브모듈은 건너뛰기
		}
		
		wg.Add(1)
		go func(info submoduleInfo) {
			defer wg.Done()
			semaphore <- struct{}{}        // 작업 슬롯 획득
			defer func() { <-semaphore }() // 작업 슬롯 반환
			
			// 필터 제거
			fmt.Printf("   🔧 %s 처리 중...\n", info.path)
			
			// 1. 필터 설정 제거
			unsetCmd := exec.Command("git", "-C", info.path, "config", "--unset", "remote.origin.partialclonefilter")
			if _, err := unsetCmd.CombinedOutput(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 5 {
					errorChan <- fmt.Errorf("%s: 필터 제거 실패 - %v", info.path, err)
					return
				}
			}
			
			// 2. 모든 blob 다운로드
			fetchCmd := exec.Command("git", "-C", info.path, "fetch", "--refetch")
			if err := fetchCmd.Run(); err != nil {
				// 실패 시 필터 복원 시도
				restoreCmd := exec.Command("git", "-C", info.path, "config", "remote.origin.partialclonefilter", info.filter)
				restoreCmd.Run()
				errorChan <- fmt.Errorf("%s: fetch 실패 - %v", info.path, err)
				return
			}
			
			resultChan <- info.path
		}(sm)
	}
	
	// 모든 작업 완료 대기
	wg.Wait()
	close(resultChan)
	close(errorChan)
	
	// 결과 집계
	successCount := len(resultChan)
	failCount := len(errorChan)
	
	// 에러 출력
	if failCount > 0 {
		fmt.Println("\n❌ 실패한 서브모듈:")
		for err := range errorChan {
			fmt.Printf("   %v\n", err)
		}
	}
	
	// 최종 결과 확인
	fmt.Println("\n3️⃣  작업 결과")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	fmt.Printf("✅ 성공: %d개\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개\n", failCount)
	}
	
	// 디스크 사용량 변화 표시
	if successCount > 0 {
		fmt.Println("\n💾 디스크 사용량 변화:")
		var totalBefore, totalAfter int64
		
		for _, sm := range submodules {
			if sm.filter == "" {
				continue
			}
			
			// 변경 후 크기 측정
			bytesAfter, humanAfter := utils.GetSubmoduleGitSize(sm.path)
			totalBefore += sm.bytesBefore
			totalAfter += bytesAfter
			
			if bytesAfter > sm.bytesBefore {
				increase := bytesAfter - sm.bytesBefore
				fmt.Printf("   📁 %s\n", sm.path)
				fmt.Printf("      %s → %s (증가: %s)\n", 
					sm.diskUsageBefore, humanAfter, utils.FormatSize(increase))
			}
		}
		
		if totalAfter > totalBefore {
			totalIncrease := totalAfter - totalBefore
			fmt.Printf("\n   📊 전체 증가량: %s\n", utils.FormatSize(totalIncrease))
		}
	}
	
	if successCount > 0 {
		fmt.Println("\n✅ 서브모듈 Partial Clone 필터 제거가 완료되었습니다!")
		fmt.Printf("   %d개 서브모듈의 모든 파일이 다운로드되었습니다.\n", successCount)
	}
}