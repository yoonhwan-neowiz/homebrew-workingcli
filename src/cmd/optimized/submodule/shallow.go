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

		// 현재 HEAD 커밋 SHA 가져오기 (detached HEAD 처리)
		headSHACmd := exec.Command("git", "rev-parse", "HEAD")
		headSHAOutput, err := headSHACmd.Output()
		if err != nil {
			return fmt.Errorf("HEAD 커밋 확인 실패: %v", err)
		}
		currentHeadSHA := strings.TrimSpace(string(headSHAOutput))
		
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
		}
		
		// Remote HEAD 기준으로 강제 shallow 변환
		// 먼저 remote에서 최신 정보를 가져옴
		fmt.Printf("🔄 %s: Remote HEAD 기준으로 shallow depth=%d 적용 중...\n", path, depth)
		
		// 1. fetch --depth로 remote의 최신 상태를 shallow로 가져옴
		fetchCmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("--depth=%d", depth), "--update-shallow")
		if err := fetchCmd.Run(); err != nil {
			// 실패 시 현재 커밋 SHA로 직접 시도
			fetchSHACmd := exec.Command("git", "fetch", "origin", currentHeadSHA, fmt.Sprintf("--depth=%d", depth))
			if err := fetchSHACmd.Run(); err != nil {
				// 그래도 실패하면 모든 참조를 shallow로 가져오기
				fetchAllCmd := exec.Command("git", "fetch", "--all", fmt.Sprintf("--depth=%d", depth))
				if err := fetchAllCmd.Run(); err != nil {
					return fmt.Errorf("Shallow fetch 실패: %v", err)
				}
			}
		}
		
		// 2. 브랜치 동기화를 shallow 변환 후에 실행 (quick/shallow.go처럼)
		// 현재 브랜치 저장
		currentBranch := "HEAD"
		currentBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		if currentBranchOutput, err := currentBranchCmd.Output(); err == nil {
			currentBranch = strings.TrimSpace(string(currentBranchOutput))
		}
		
		// 모든 로컬 브랜치를 remote와 동기화 (shallow 상태 적용을 위해)
		fmt.Printf("   ├─ 로컬 브랜치들을 shallow 상태로 동기화 중...\n")
		branchListCmd := exec.Command("git", "branch", "--format=%(refname:short)")
		if branchOutput, err := branchListCmd.Output(); err == nil {
			subBranches := strings.Split(strings.TrimSpace(string(branchOutput)), "\n")
			for _, subBranch := range subBranches {
				subBranch = strings.TrimSpace(subBranch)
				if subBranch == "" {
					continue
				}
				
				// 각 브랜치를 remote와 동기화
				checkoutCmd := exec.Command("git", "checkout", subBranch, "-q")
				if err := checkoutCmd.Run(); err == nil {
					// remote 브랜치가 있는지 확인
					remoteVerifyCmd := exec.Command("git", "rev-parse", "--verify", "origin/"+subBranch)
					if err := remoteVerifyCmd.Run(); err == nil {
						// reset --hard origin/branch로 shallow 상태 강제 적용
						resetCmd := exec.Command("git", "reset", "--hard", "origin/"+subBranch)
						resetCmd.Run()
					}
				}
			}
			
			// 원래 브랜치/커밋으로 돌아가기
			if currentBranch == "HEAD" {
				// Detached HEAD 상태였다면 원래 커밋으로
				checkoutSHACmd := exec.Command("git", "checkout", currentHeadSHA, "-q")
				checkoutSHACmd.Run()
			} else {
				// 브랜치였다면 브랜치로
				checkoutBackCmd := exec.Command("git", "checkout", currentBranch, "-q")
				checkoutBackCmd.Run()
			}
		}
		
		// 3. reflog 정리로 이전 히스토리 참조 제거
		reflogCmd := exec.Command("git", "reflog", "expire", "--expire=now", "--all")
		reflogCmd.Run()
		
		// 4. gc 실행으로 오래된 객체 완전 정리 (--aggressive 추가)
		gcCmd := exec.Command("git", "gc", "--prune=now", "--aggressive")
		gcCmd.Run()
		
		fmt.Printf("✅ %s: Shallow Clone으로 변환 완료 (depth=%d)\n", path, depth)
		
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