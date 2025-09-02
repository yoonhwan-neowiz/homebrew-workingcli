package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewToSlimCmd creates the To SLIM conversion command for submodules
func NewToSlimCmd() *cobra.Command {
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:   "to-slim",
		Short: "서브모듈을 SLIM 모드로 전환 (recursive)",
		Long: `모든 서브모듈을 SLIM 모드로 전환합니다 (recursive).
각 서브모듈에 Partial Clone 및 Sparse Checkout을 적용하고 GC를 수행하여 디스크 사용량을 줄입니다.

실행 내용:
1) Partial Clone 필터 설정: config 기반 (기본값: blob:limit=1m)
2) Sparse Checkout 활성화 (cone, 루트 최소 경로)
3) 불필요한 객체 정리 (repack + maintenance gc)

참고: 대용량 저장소의 경우 시간이 소요될 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runToSlim()
		},
	}
	
	// -q 플래그 추가
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

// runToSlim converts all submodules to SLIM mode in parallel
func runToSlim() {
	// 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}

	fmt.Println("🚀 모든 서브모듈을 SLIM 모드로 전환합니다...")
	fmt.Println("⚠️ 주의: 일부 서브모듈은 시간이 걸릴 수 있습니다.")

	// SLIM 전환은 안전한 작업이므로 quiet 모드에서 자동 수락
	if !utils.ConfirmForce("계속하시겠습니까?") {
		fmt.Println("❌ 작업이 취소되었습니다.")
		return
	}

	// config에서 서브모듈용 필터 설정 읽기
	submoduleFilter := config.GetString("optimize.submodule.filter.default")
	if submoduleFilter == "" {
		// 서브모듈용 기본값이 없으면 일반 기본값 사용
		submoduleFilter = config.GetString("optimize.filter.default")
		if submoduleFilter == "" {
			// 그래도 없으면 하드코딩된 기본값
			submoduleFilter = "blob:limit=1m"
		}
	}
	// blob:limit= 접두사 추가 (설정에 숫자만 있는 경우)
	if !strings.HasPrefix(submoduleFilter, "blob:") {
		submoduleFilter = "blob:limit=" + submoduleFilter
	}

	// config에서 sparse 경로 설정 읽기 (서브모듈용)
	sparsePaths := getSubmoduleSparsePaths()

	fmt.Println()
	fmt.Printf("📦 총 %d개의 서브모듈을 병렬로 처리합니다.\n", count)
	fmt.Printf("🔧 필터 설정: %s\n\n", submoduleFilter)

	// 결과 집계용 구조체 및 공유 변수
	type slimResult struct {
		path        string
		beforeHuman string
		afterHuman  string
		beforeBytes int64
		afterBytes  int64
		changed     bool
	}

	var (
		mu       sync.Mutex
		results  []slimResult
		started  = time.Now()
	)

	// 서브모듈 변환 작업 정의
	toSlimOperation := func(path string) error {
		// 서브모듈 디렉토리로 이동
		originalDir, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("디렉토리 이동 실패: %v", err)
		}
		defer os.Chdir(originalDir)

		// Git 저장소인지 확인 (미초기화 서브모듈 대응)
		if !utils.IsGitRepository() {
			mu.Lock()
			results = append(results, slimResult{path: path, beforeHuman: "미초기화", afterHuman: "미초기화"})
			mu.Unlock()
			fmt.Printf("ℹ️ %s: 미초기화 서브모듈 (건너뜀)\n", path)
			return nil
		}

		// 현재 .git 디렉토리 크기 측정
		beforeBytes, beforeHuman := utils.GetGitDirSize(".")

		// 1) Partial Clone 필터 설정 (config에서 읽은 값 사용)
		exec.Command("git", "config", "remote.origin.partialclonefilter", submoduleFilter).Run()
		exec.Command("git", "config", "remote.origin.promisor", "true").Run()
		exec.Command("git", "config", "extensions.partialClone", "origin").Run()

		// 2) Sparse Checkout 활성화
		exec.Command("git", "config", "core.sparseCheckout", "true").Run()
		
		// sparsePaths 설정에 따라 처리
		if len(sparsePaths) == 0 || (len(sparsePaths) == 1 && sparsePaths[0] == "*") {
			// 기본값: cone 모드로 루트만
			exec.Command("git", "sparse-checkout", "init", "--cone").Run()
			exec.Command("git", "sparse-checkout", "set", "/").Run()
		} else {
			// 사용자 지정 경로가 있는 경우
			hasFiles := false
			for _, path := range sparsePaths {
				if !strings.HasSuffix(path, "/") && strings.Contains(path, ".") {
					hasFiles = true
					break
				}
			}
			
			if hasFiles {
				exec.Command("git", "sparse-checkout", "init", "--no-cone").Run()
			} else {
				exec.Command("git", "sparse-checkout", "init", "--cone").Run()
			}
			
			args := append([]string{"sparse-checkout", "set"}, sparsePaths...)
			exec.Command("git", args...).Run()
		}

		// 네트워크 반영을 위한 안전한 fetch (필터 반영)
		exec.Command("git", "fetch", "--prune").Run()

		// 3) 불필요한 객체 정리 및 성능 설정 일부 적용
		exec.Command("git", "repack", "-a", "-d").Run()
		exec.Command("git", "maintenance", "run", "--task=gc").Run()

		// 전환 후 크기 측정
		afterBytes, afterHuman := utils.GetGitDirSize(".")

		// 결과 기록
		mu.Lock()
		results = append(results, slimResult{
			path:        path,
			beforeHuman: beforeHuman,
			afterHuman:  afterHuman,
			beforeBytes: beforeBytes,
			afterBytes:  afterBytes,
			changed:     afterBytes <= beforeBytes,
		})
		mu.Unlock()

		// 개별 결과 출력
		fmt.Printf("✅ %s: %s → %s\n", path, beforeHuman, afterHuman)
		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	successCount, failCount, err := utils.ExecuteOnSubmodulesParallel(toSlimOperation, 4, true)

	// 요약 출력
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 작업 완료 요약")

	var totalBefore, totalAfter int64
	for _, r := range results {
		totalBefore += r.beforeBytes
		totalAfter += r.afterBytes
	}

	saved := totalBefore - totalAfter
	if saved < 0 {
		saved = 0
	}

	fmt.Printf("✅ 성공: %d개 서브모듈\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", failCount)
	}
	fmt.Printf("📦 총 크기: %s → %s (절감: %s)\n", utils.HumanizeBytes(totalBefore), utils.HumanizeBytes(totalAfter), utils.HumanizeBytes(saved))
	fmt.Printf("⏱  소요 시간: %v\n", time.Since(started).Round(time.Second))

	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}
}

// getSubmoduleSparsePaths gets sparse checkout paths from config for submodules
func getSubmoduleSparsePaths() []string {
	settings := config.GetAll()
	var paths []string
	
	// 먼저 서브모듈 전용 설정을 찾기
	if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
		if submodule, ok := optimize["submodule"].(map[string]interface{}); ok {
			if sparse, ok := submodule["sparse"].(map[string]interface{}); ok {
				if configPaths, ok := sparse["paths"].([]interface{}); ok {
					for _, path := range configPaths {
						if p, ok := path.(string); ok {
							paths = append(paths, p)
						}
					}
					return paths
				}
			}
		}
		
		// 서브모듈 전용 설정이 없으면 일반 sparse 설정 사용
		if sparse, ok := optimize["sparse"].(map[string]interface{}); ok {
			if configPaths, ok := sparse["paths"].([]interface{}); ok {
				for _, path := range configPaths {
					if p, ok := path.(string); ok {
						paths = append(paths, p)
					}
				}
			}
		}
	}
	
	return paths
}
