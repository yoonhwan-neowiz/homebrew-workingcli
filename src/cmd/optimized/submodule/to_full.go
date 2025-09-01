package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewToFullCmd creates the To FULL restoration command for submodules
func NewToFullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to-full",
		Short: "서브모듈을 FULL 모드로 복원 (recursive)",
		Long: `모든 서브모듈을 FULL 모드로 복원합니다 (recursive).
각 서브모듈의 Partial Clone 필터를 제거하고 Sparse Checkout을 비활성화하여 
전체 히스토리와 모든 파일을 복원합니다.

실행 내용:
1) Partial Clone 필터 제거
2) Sparse Checkout 비활성화
3) 모든 객체 다운로드 (fetch --unshallow)
4) 저장소 최적화 (repack + maintenance gc)

참고: 네트워크를 통해 모든 객체를 다운로드하므로 시간과 대역폭이 소요됩니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runToFull()
		},
	}
}

// runToFull restores all submodules to FULL mode in parallel
func runToFull() {
	// 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}

	fmt.Println("🚀 모든 서브모듈을 FULL 모드로 복원합니다...")
	fmt.Println("⚠️ 주의: 모든 객체를 다운로드하므로 시간과 네트워크 대역폭이 소요됩니다.")
	fmt.Println("⚠️ 디스크 사용량이 크게 증가할 수 있습니다.")

	if !utils.ConfirmWithDefault("계속하시겠습니까?", true) {
		fmt.Println("❌ 작업이 취소되었습니다.")
		return
	}

	fmt.Println()
	fmt.Printf("📦 총 %d개의 서브모듈을 병렬로 처리합니다.\n", count)
	fmt.Println("🔧 작업: Partial Clone 제거 + Sparse Checkout 비활성화 + 전체 히스토리 복원\n")

	// 결과 집계용 구조체 및 공유 변수
	type fullResult struct {
		path        string
		beforeHuman string
		afterHuman  string
		beforeBytes int64
		afterBytes  int64
		changed     bool
		wasShallow  bool
		wasSparse   bool
		wasPartial  bool
	}

	var (
		mu      sync.Mutex
		results []fullResult
		started = time.Now()
	)

	// 서브모듈 복원 작업 정의
	toFullOperation := func(path string) error {
		// 서브모듈 디렉토리로 이동
		originalDir, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("디렉토리 이동 실패: %v", err)
		}
		defer os.Chdir(originalDir)

		// Git 저장소인지 확인 (미초기화 서브모듈 대응)
		if !utils.IsGitRepository() {
			mu.Lock()
			results = append(results, fullResult{path: path, beforeHuman: "미초기화", afterHuman: "미초기화"})
			mu.Unlock()
			fmt.Printf("ℹ️ %s: 미초기화 서브모듈 (건너뜀)\n", path)
			return nil
		}

		// 현재 상태 확인
		wasShallow := false
		if output, _ := exec.Command("git", "rev-parse", "--is-shallow-repository").Output(); strings.TrimSpace(string(output)) == "true" {
			wasShallow = true
		}

		wasSparse := false
		if output, _ := exec.Command("git", "config", "core.sparseCheckout").Output(); strings.TrimSpace(string(output)) == "true" {
			wasSparse = true
		}

		wasPartial := false
		if output, _ := exec.Command("git", "config", "remote.origin.partialclonefilter").Output(); strings.TrimSpace(string(output)) != "" {
			wasPartial = true
		}

		// 이미 FULL 모드인 경우
		if !wasShallow && !wasSparse && !wasPartial {
			mu.Lock()
			beforeBytes, beforeHuman := utils.GetGitDirSize(".")
			results = append(results, fullResult{
				path:        path,
				beforeHuman: beforeHuman,
				afterHuman:  beforeHuman,
				beforeBytes: beforeBytes,
				afterBytes:  beforeBytes,
				changed:     false,
			})
			mu.Unlock()
			fmt.Printf("✅ %s: 이미 FULL 모드입니다 (%s)\n", path, beforeHuman)
			return nil
		}

		// 현재 .git 디렉토리 크기 측정
		beforeBytes, beforeHuman := utils.GetGitDirSize(".")

		// 1) Partial Clone 필터 제거
		if wasPartial {
			exec.Command("git", "config", "--unset", "remote.origin.partialclonefilter").Run()
			exec.Command("git", "config", "--unset", "remote.origin.promisor").Run()
			exec.Command("git", "config", "--unset", "extensions.partialClone").Run()
		}

		// 2) Sparse Checkout 비활성화
		if wasSparse {
			exec.Command("git", "sparse-checkout", "disable").Run()
			exec.Command("git", "config", "core.sparseCheckout", "false").Run()
		}

		// 3) 전체 히스토리 복원
		if wasShallow {
			// Shallow 저장소를 완전한 저장소로 변환
			cmd := exec.Command("git", "fetch", "--unshallow")
			if err := cmd.Run(); err != nil {
				// 이미 unshallow 상태일 수 있음
				exec.Command("git", "fetch", "--all").Run()
			}
		} else if wasPartial || wasSparse {
			// Partial Clone이나 Sparse였던 경우 모든 객체 다운로드
			exec.Command("git", "fetch", "--all", "--prune").Run()
		}

		// 4) 작업 트리 재설정 (Sparse Checkout 비활성화 후)
		if wasSparse {
			exec.Command("git", "read-tree", "-m", "-u", "HEAD").Run()
			exec.Command("git", "checkout", ".").Run()
		}

		// 5) 저장소 최적화
		exec.Command("git", "repack", "-a", "-d", "-f").Run()
		exec.Command("git", "maintenance", "run", "--task=gc").Run()
		exec.Command("git", "prune").Run()

		// 복원 후 크기 측정
		afterBytes, afterHuman := utils.GetGitDirSize(".")

		// 결과 기록
		mu.Lock()
		results = append(results, fullResult{
			path:        path,
			beforeHuman: beforeHuman,
			afterHuman:  afterHuman,
			beforeBytes: beforeBytes,
			afterBytes:  afterBytes,
			changed:     true,
			wasShallow:  wasShallow,
			wasSparse:   wasSparse,
			wasPartial:  wasPartial,
		})
		mu.Unlock()

		// 개별 결과 출력
		status := ""
		if wasShallow {
			status += "Shallow→Full "
		}
		if wasSparse {
			status += "Sparse→Full "
		}
		if wasPartial {
			status += "Partial→Full"
		}
		fmt.Printf("✅ %s: %s → %s [%s]\n", path, beforeHuman, afterHuman, strings.TrimSpace(status))
		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	successCount, failCount, err := utils.ExecuteOnSubmodulesParallel(toFullOperation, 4, true)

	// 요약 출력
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 작업 완료 요약")

	var totalBefore, totalAfter int64
	var restoredCount, skippedCount int
	for _, r := range results {
		totalBefore += r.beforeBytes
		totalAfter += r.afterBytes
		if r.changed {
			restoredCount++
		} else {
			skippedCount++
		}
	}

	increased := totalAfter - totalBefore
	if increased < 0 {
		increased = 0
	}

	fmt.Printf("✅ 성공: %d개 서브모듈\n", successCount)
	if restoredCount > 0 {
		fmt.Printf("🔄 복원됨: %d개 (SLIM → FULL)\n", restoredCount)
	}
	if skippedCount > 0 {
		fmt.Printf("⏭️ 건너뜀: %d개 (이미 FULL 모드)\n", skippedCount)
	}
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", failCount)
	}
	fmt.Printf("📦 총 크기: %s → %s (증가: %s)\n", utils.HumanizeBytes(totalBefore), utils.HumanizeBytes(totalAfter), utils.HumanizeBytes(increased))
	fmt.Printf("⏱  소요 시간: %v\n", time.Since(started).Round(time.Second))

	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}
}

