package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewClearBranchScopeCmd creates the submodule Clear Branch Scope command
func NewClearBranchScopeCmd() *cobra.Command {
	var fetchFlag bool
	
	cmd := &cobra.Command{
		Use:     "clear-branch-scope",
		Aliases: []string{"cbs", "unscope", "show-all"},
		Short:   "서브모듈 브랜치 범위 제거 (모든 브랜치 표시)",
		Long: `서브모듈의 브랜치 범위를 제거하여 모든 로컬/원격 브랜치가 표시되도록 합니다.
set-branch-scope로 설정한 범위를 초기화합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runSubmoduleClearScope(fetchFlag)
		},
	}
	
	cmd.Flags().BoolVarP(&fetchFlag, "fetch", "f", false, "원격 브랜치를 다시 가져옴")
	
	return cmd
}

func runSubmoduleClearScope(fetchFlag bool) {
	fmt.Println("\n🔧 서브모듈 브랜치 범위 제거")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 서브모듈 존재 확인
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		fmt.Println("\nℹ️  서브모듈이 없습니다.")
		return
	}
	
	// 서브모듈 목록 가져오기
	cmd := exec.Command("git", "submodule", "foreach", "--quiet", "echo $path")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("\n❌ 서브모듈 목록을 가져올 수 없습니다: %v\n", err)
		return
	}
	
	submodulePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(submodulePaths) == 0 || (len(submodulePaths) == 1 && submodulePaths[0] == "") {
		fmt.Println("\nℹ️  초기화된 서브모듈이 없습니다.")
		return
	}
	
	// config에서 서브모듈 branch_scope 확인
	submoduleScope := config.GetSubmoduleBranchScope()
	if len(submoduleScope) == 0 {
		fmt.Println("\nℹ️  현재 설정된 브랜치 범위가 없습니다")
		return
	}
	
	fmt.Println("\n📋 현재 설정된 브랜치 범위:")
	for _, branch := range submoduleScope {
		fmt.Printf("   • %s\n", branch)
	}
	
	// 사용자 확인
	if !utils.ConfirmWithDefault("\n브랜치 범위를 제거하시겠습니까?", false) {
		fmt.Println("\n✨ 작업이 취소되었습니다")
		return
	}
	
	// 필터 제거
	clearSubmoduleBranchFilters(submodulePaths, fetchFlag)
}

func clearSubmoduleBranchFilters(submodulePaths []string, fetchFlag bool) {
	successCount := 0
	failCount := 0
	
	// config에서 서브모듈 branch_scope 제거
	if err := config.ClearSubmoduleBranchScope(); err != nil {
		fmt.Printf("⚠️ config.yaml 서브모듈 브랜치 스코프 제거 실패: %v\n", err)
	}
	
	// 각 서브모듈의 fetch refspec 복원
	for _, path := range submodulePaths {
		if path == "" {
			continue
		}
		
		// 서브모듈의 fetch refspec 복원
		if err := utils.RestoreFetchRefspecForSubmodule(path); err != nil {
			fmt.Printf("⚠️  %s fetch refspec 복원 실패: %v\n", path, err)
			failCount++
		} else {
			successCount++
		}
	}
	
	// fetch 플래그가 설정된 경우에만 각 서브모듈의 원격 브랜치 가져오기
	if fetchFlag {
		fmt.Println("\n🔄 서브모듈의 원격 브랜치를 가져오는 중...")
		for _, path := range submodulePaths {
			if path == "" {
				continue
			}
			
			cmd := exec.Command("git", "-C", path, "fetch", "origin", "--prune")
			if err := cmd.Run(); err != nil {
				fmt.Printf("⚠️  %s: 원격 브랜치 가져오기 실패: %v\n", path, err)
			} else {
				fmt.Printf("✅  %s: 원격 브랜치를 성공적으로 가져왔습니다\n", path)
			}
		}
	}
	
	fmt.Println("\n✅ 서브모듈 브랜치 범위가 제거되었습니다")
	fmt.Println("\n📋 결과:")
	fmt.Println("   • 모든 로컬 브랜치가 표시됩니다")
	fmt.Println("   • 모든 원격 브랜치가 표시됩니다")
	
	if successCount > 0 {
		fmt.Printf("\n🌳 필터 제거 상태:\n")
		fmt.Printf("   • 성공: %d개 서브모듈\n", successCount)
		if failCount > 0 {
			fmt.Printf("   • 실패: %d개 서브모듈\n", failCount)
		}
	}
}