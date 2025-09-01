package quick

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewClearBranchScopeCmd creates the Clear Branch Scope command
func NewClearBranchScopeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "clear-branch-scope",
		Aliases: []string{"cbs", "unscope", "show-all"},
		Short:   "브랜치 범위 제거 (모든 브랜치 표시)",
		Long: `브랜치 범위를 제거하여 모든 로컬/원격 브랜치가 표시되도록 합니다.
set-branch-scope로 설정한 범위를 초기화합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runClearScope()
		},
	}
}

func runClearScope() {
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다")
		return
	}

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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n브랜치 범위를 제거하시겠습니까? (y/N): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("\n✨ 작업이 취소되었습니다")
		return
	}

	// 범위 제거
	clearBranchScope()
}

func clearBranchScope() {
	// Config 파일에서 브랜치 범위 제거
	err := config.ClearBranchScope()
	if err != nil {
		fmt.Printf("\n⚠️ 브랜치 범위 제거 중 경고: %v\n", err)
		// 경고만 표시하고 계속 진행
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