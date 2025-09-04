package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewUpdateCmd creates the Update Submodules command
func NewUpdateCmd() *cobra.Command {
	var (
		forceUpdate  bool
		remoteUpdate bool
		quietMode    bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "서브모듈 업데이트 (참조 문제 자동 해결)",
		Long: `서브모듈을 업데이트합니다.
일반 업데이트 실패 시 자동으로 원격 최신 커밋으로 업데이트를 시도합니다.

실행 내용:
- git submodule update --init --recursive
- 실패 시: git submodule update --init --remote --force
- 또는 git submodule foreach를 사용한 개별 처리

사용 예시:
  ga opt submodule update          # 일반 업데이트
  ga opt submodule update -f       # 강제 업데이트 (원격 최신으로)
  ga opt submodule update --remote # 원격 추적 브랜치의 최신으로 업데이트

⚠️ 주의: 서브모듈 참조가 원격에 없는 경우 자동으로 최신 커밋으로 업데이트됩니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runUpdate(forceUpdate, remoteUpdate)
		},
	}

	cmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "강제 업데이트 (원격 최신 커밋으로)")
	cmd.Flags().BoolVar(&remoteUpdate, "remote", false, "원격 추적 브랜치의 최신으로 업데이트")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")

	return cmd
}

func runUpdate(forceUpdate bool, remoteUpdate bool) {
	fmt.Println("🔄 서브모듈 업데이트를 시작합니다...")
	
	// 먼저 서브모듈 초기화 (init) - 이게 있어야 카운트도 제대로 나옴
	fmt.Println("📥 서브모듈 초기화 중...")
	initCmd := exec.Command("git", "submodule", "update", "--init", "--recursive")
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	initErr := initCmd.Run()
	
	// 이제 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}
	
	fmt.Printf("\n📦 총 %d개의 서브모듈이 발견되었습니다.\n", count)
	
	// 강제 또는 원격 업데이트 모드 확인
	if forceUpdate || remoteUpdate {
		fmt.Println("\n⚠️ 강제 업데이트 모드: 원격 최신 커밋으로 업데이트합니다.")
		
		// 사용자 확인
		if !utils.ConfirmForce("모든 서브모듈을 원격 최신으로 업데이트하시겠습니까?") {
			fmt.Println("❌ 작업이 취소되었습니다.")
			return
		}
		
		performForceUpdate()
		return
	}
	
	// init 명령이 실패했다면 원격 최신으로 시도
	if initErr != nil {
		fmt.Println("\n⚠️ 일반 업데이트 실패 - 원격 최신 커밋으로 시도합니다...")
		
		// 사용자 확인 (quiet 모드가 아닌 경우)
		if !utils.ConfirmForce("원격 최신 커밋으로 업데이트를 시도하시겠습니까?") {
			fmt.Println("❌ 작업이 취소되었습니다.")
			return
		}
		
		performForceUpdate()
		return
	}
	
	fmt.Println("\n✅ 서브모듈 업데이트 완료!")
	showSubmoduleStatus()
}

func performForceUpdate() {
	// 서브모듈 동기화
	fmt.Println("\n🔄 서브모듈 설정 동기화 중...")
	syncCmd := exec.Command("git", "submodule", "sync", "--recursive")
	if err := syncCmd.Run(); err != nil {
		fmt.Printf("⚠️ 동기화 실패: %v\n", err)
	}

	// foreach를 사용한 개별 업데이트
	updateOperation := func(path string) error {
		originalDir, _ := os.Getwd()
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("디렉토리 이동 실패: %v", err)
		}
		defer os.Chdir(originalDir)

		fmt.Printf("📦 %s: 업데이트 중...\n", path)
		
		// 먼저 일반 업데이트 시도
		checkoutCmd := exec.Command("git", "checkout", "-f", "HEAD")
		checkoutCmd.Run() // 에러 무시
		
		// fetch all
		fetchCmd := exec.Command("git", "fetch", "--all")
		if err := fetchCmd.Run(); err != nil {
			fmt.Printf("  ⚠️ fetch 실패: %v\n", err)
		}
		
		// 원격 브랜치 확인 및 리셋
		// origin/HEAD 또는 origin/master, origin/main 시도
		var resetSuccess bool
		for _, ref := range []string{"origin/HEAD", "origin/master", "origin/main"} {
			resetCmd := exec.Command("git", "reset", "--hard", ref)
			if err := resetCmd.Run(); err == nil {
				fmt.Printf("✅ %s: %s로 업데이트 완료\n", path, ref)
				resetSuccess = true
				break
			}
		}
		
		if !resetSuccess {
			// 리셋 실패 시 최신 커밋으로 시도
			logCmd := exec.Command("git", "log", "--oneline", "-1", "--remotes")
			if output, err := logCmd.Output(); err == nil {
				parts := strings.Fields(string(output))
				if len(parts) > 0 {
					resetCmd := exec.Command("git", "reset", "--hard", parts[0])
					if err := resetCmd.Run(); err == nil {
						fmt.Printf("✅ %s: 최신 커밋으로 업데이트 완료\n", path)
						return nil
					}
				}
			}
			return fmt.Errorf("%s: 업데이트 실패", path)
		}
		
		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	successCount, failCount, err := utils.ExecuteOnSubmodulesParallel(updateOperation, 4, true)

	// 결과 요약
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 업데이트 완료 요약")
	fmt.Printf("✅ 성공: %d개 서브모듈\n", successCount)
	if failCount > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", failCount)
	}
	
	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}
	
	// 최종 상태 표시
	showSubmoduleStatus()
}

func showSubmoduleStatus() {
	fmt.Println("\n📋 서브모듈 상태:")
	statusCmd := exec.Command("git", "submodule", "status", "--recursive")
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	statusCmd.Run()
}