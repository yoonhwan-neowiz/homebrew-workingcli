package quick

import (
	"fmt"
	"os/exec"
	
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewClearBranchScopeCmd creates the Clear Branch Scope command
func NewClearBranchScopeCmd() *cobra.Command {
	var fetchFlag bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "clear-branch-scope",
		Aliases: []string{"cbs", "unscope", "show-all"},
		Short:   "브랜치 범위 제거 (모든 브랜치 표시)",
		Long: `브랜치 범위를 제거하여 모든 로컬/원격 브랜치가 표시되도록 합니다.
set-branch-scope로 설정한 범위를 초기화합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runClearScope(fetchFlag)
		},
	}
	
	cmd.Flags().BoolVarP(&fetchFlag, "fetch", "f", false, "원격 브랜치를 다시 가져옴")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runClearScope(fetchFlag bool) {

	fmt.Println("\n🔧 브랜치 범위 제거")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 현재 범위 설정 확인
	currentScope := config.GetBranchScope()
	if len(currentScope) == 0 {
		fmt.Println("\nℹ️ 현재 설정된 브랜치 범위가 없습니다")
		return
	}

	fmt.Println("\n📋 현재 설정된 브랜치 범위:")
	for _, branch := range currentScope {
		fmt.Printf("   • %s\n", branch)
	}

	// 사용자 확인
	// 브랜치 범위 제거는 안전한 작업이므로 quiet 모드에서 자동 수락
	if !utils.ConfirmForce("\n브랜치 범위를 제거하시겠습니까?") {
		fmt.Println("\n✨ 작업이 취소되었습니다")
		return
	}

	// 범위 제거
	clearBranchScope(fetchFlag)
}

func clearBranchScope(fetchFlag bool) {
	// Config 파일에서 브랜치 범위 제거
	err := config.ClearBranchScope()
	if err != nil {
		fmt.Printf("\n⚠️ 브랜치 범위 제거 중 경고: %v\n", err)
		// 경고만 표시하고 계속 진행
	}
	
	// Git fetch refspec 복원 (백업에서 복원하거나 기본값으로 설정)
	if err := utils.RestoreFetchRefspec(); err != nil {
		fmt.Printf("\n⚠️ fetch refspec 복원 실패: %v\n", err)
	}

	// fetch 플래그가 설정된 경우에만 원격 브랜치 가져오기
	if fetchFlag {
		fmt.Println("\n🔄 원격 브랜치를 가져오는 중...")
		cmd := exec.Command("git", "fetch", "origin", "--prune")
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️ 원격 브랜치 가져오기 실패: %v\n", err)
		} else {
			fmt.Println("✅ 원격 브랜치를 성공적으로 가져왔습니다")
		}
	}

	fmt.Println("\n✅ 브랜치 범위가 제거되었습니다")
	fmt.Println("\n📋 결과:")
	fmt.Println("   • 모든 로컬 브랜치가 표시됩니다")
	fmt.Println("   • 모든 원격 브랜치가 표시됩니다")

	// 현재 브랜치 수 표시
	localCount := utils.CountLocalBranches()
	remoteCount := utils.CountRemoteBranches()

	fmt.Printf("\n🌳 브랜치 상태:\n")
	fmt.Printf("   • 로컬 브랜치: %d개\n", localCount)
	fmt.Printf("   • 원격 브랜치: %d개\n", remoteCount)
}