package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewShallowCmd creates the Shallow Submodules command
func NewShallowCmd() *cobra.Command {
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:   "shallow [depth]",
		Short: "서브모듈을 Shallow Clone으로 변환 (recursive)",
		Long: `모든 서브모듈을 Shallow Clone으로 변환합니다 (recursive).
depth를 지정하지 않으면 기본값 1(최신 1개 커밋)로 설정됩니다.
각 서브모듈의 히스토리를 제한하여 디스크 공간을 절약합니다.

예시:
  ga opt submodule shallow        # depth=1로 설정 (기본값)
  ga opt submodule shallow 5      # 최근 5개 커밋만 유지
  ga opt submodule shallow 10     # 최근 10개 커밋만 유지
  ga opt submodule shallow 5 -q   # quiet 모드로 자동 실행`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runShallow(args)
		},
	}
	
	// -q 플래그 추가
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runShallow(args []string) {
	// depth 파라미터 처리
	depth := 1
	if len(args) > 0 {
		if d, err := strconv.Atoi(args[0]); err == nil && d > 0 {
			depth = d
		} else {
			fmt.Printf("❌ 잘못된 depth 값: %s (양의 정수여야 합니다)\n", args[0])
			os.Exit(1)
		}
	}

	// 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}

	fmt.Printf("🔄 모든 서브모듈을 Shallow Clone으로 변환합니다 (depth=%d)...\n", depth)
	fmt.Printf("📦 총 %d개의 서브모듈을 병렬로 처리합니다.\n\n", count)

	// Shallow 변환 작업 정의
	shallowOperation := func(path string) error {
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

		if isShallow {
			// 이미 shallow인 경우 depth 확인
			countCmd := exec.Command("git", "rev-list", "--count", "HEAD")
			countOutput, _ := countCmd.Output()
			currentDepth := strings.TrimSpace(string(countOutput))
			
			if currentD, _ := strconv.Atoi(currentDepth); currentD == depth {
				fmt.Printf("ℹ️ %s: 이미 Shallow 상태 (depth: %s)\n", path, currentDepth)
				return nil // 성공으로 처리
			}
			
			// depth 업데이트 - fetch를 먼저 시도
			fetchCmd := exec.Command("git", "fetch", fmt.Sprintf("--depth=%d", depth))
			if err := fetchCmd.Run(); err != nil {
				// fetch 실패 시 pull with --allow-unrelated-histories
				pullCmd := exec.Command("git", "pull", fmt.Sprintf("--depth=%d", depth), "--allow-unrelated-histories")
				if err := pullCmd.Run(); err != nil {
					return fmt.Errorf("Shallow 업데이트 실패: %v", err)
				}
			}
			fmt.Printf("✅ %s: Depth를 %d로 변경\n", path, depth)
		} else {
			// shallow로 변환 - fetch를 먼저 시도 (더 안전)
			fetchCmd := exec.Command("git", "fetch", fmt.Sprintf("--depth=%d", depth))
			if err := fetchCmd.Run(); err != nil {
				// fetch 실패 시 pull with --allow-unrelated-histories
				pullCmd := exec.Command("git", "pull", fmt.Sprintf("--depth=%d", depth), "--allow-unrelated-histories")
				if err := pullCmd.Run(); err != nil {
					// 그래도 실패하면 origin과 현재 브랜치를 명시적으로 지정
					branch := "HEAD"
					branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
					if branchOutput, err := branchCmd.Output(); err == nil {
						branch = strings.TrimSpace(string(branchOutput))
					}
					
					fetchOriginCmd := exec.Command("git", "fetch", "origin", branch, fmt.Sprintf("--depth=%d", depth))
					if err := fetchOriginCmd.Run(); err != nil {
						return fmt.Errorf("Shallow 변환 실패: %v", err)
					}
				}
			}
			
			// gc 실행으로 오래된 객체 정리
			gcCmd := exec.Command("git", "gc", "--prune=now")
			gcCmd.Run()
			fmt.Printf("✅ %s: Shallow Clone으로 변환 완료\n", path)
		}
		
		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	successCount, failCount, err := utils.ExecuteOnSubmodulesParallel(shallowOperation, 4, true)

	// 요약
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 작업 완료 요약")
	fmt.Printf("✅ 성공: %d개 서브모듈\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", failCount)
	}
	
	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}
	
	fmt.Printf("\n모든 서브모듈의 depth가 %d로 설정되었습니다.\n", depth)
}