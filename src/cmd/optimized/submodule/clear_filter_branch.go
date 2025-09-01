package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewClearFilterBranchCmd creates the submodule Clear Filter Branch command
func NewClearFilterBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear-filter-branch",
		Short: "서브모듈 브랜치 필터 제거 (모든 브랜치 표시)",
		Long: `서브모듈의 브랜치 필터를 제거하여 모든 로컬/원격 브랜치가 표시되도록 합니다.
filter-branch로 설정한 필터를 초기화합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runSubmoduleClearFilter()
		},
	}
}

func runSubmoduleClearFilter() {
	fmt.Println("\n🔧 서브모듈 브랜치 필터 제거")
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
	
	// 현재 필터가 설정된 서브모듈 찾기
	var hasFilter bool
	var filteredSubmodules []string
	filterInfo := make(map[string][]string)
	
	for _, path := range submodulePaths {
		if path == "" {
			continue
		}
		
		configKey := fmt.Sprintf("submodule.%s.branchFilter", path)
		getCmd := exec.Command("git", "config", "--get", configKey)
		output, err := getCmd.Output()
		
		if err == nil && len(output) > 0 {
			branchList := strings.TrimSpace(string(output))
			if branchList != "" {
				branches := strings.Split(branchList, ",")
				filterInfo[path] = branches
				filteredSubmodules = append(filteredSubmodules, path)
				hasFilter = true
			}
		}
	}
	
	if !hasFilter {
		fmt.Println("\nℹ️  현재 설정된 브랜치 필터가 없습니다")
		return
	}
	
	fmt.Println("\n📋 현재 필터링된 서브모듈:")
	for _, path := range filteredSubmodules {
		branches := filterInfo[path]
		fmt.Printf("   • %s (필터: %s)\n", path, strings.Join(branches, ", "))
	}
	
	// 사용자 확인
	if !utils.ConfirmWithDefault("\n브랜치 필터를 제거하시겠습니까?", false) {
		fmt.Println("\n✨ 작업이 취소되었습니다")
		return
	}
	
	// 필터 제거
	clearSubmoduleBranchFilters(filteredSubmodules)
}

func clearSubmoduleBranchFilters(submodules []string) {
	successCount := 0
	failCount := 0
	
	for _, path := range submodules {
		// 메인 저장소의 서브모듈 설정 제거
		configKey := fmt.Sprintf("submodule.%s.branchFilter", path)
		unsetCmd := exec.Command("git", "config", "--unset", configKey)
		if err := unsetCmd.Run(); err != nil {
			// Exit code 5는 키가 없는 경우 (이미 제거됨)
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 5 {
				// 이미 제거됨 - 성공으로 처리
				successCount++
			} else {
				fmt.Printf("\n⚠️  %s 필터 제거 중 경고: %v\n", path, err)
				failCount++
			}
		} else {
			successCount++
		}
		
		// 서브모듈 디렉토리의 설정도 제거
		submoduleUnsetCmd := exec.Command("git", "-C", path, "config", "--unset", "workingcli.branchFilter")
		submoduleUnsetCmd.Run() // 실패해도 무시 (서브모듈 내부 설정은 선택적)
	}
	
	fmt.Println("\n✅ 서브모듈 브랜치 필터가 제거되었습니다")
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