package submodule

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewSetBranchScopeCmd creates the submodule Set Branch Scope command
func NewSetBranchScopeCmd() *cobra.Command {
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "set-branch-scope [브랜치1] [브랜치2] ...",
		Aliases: []string{"sbs", "scope", "branch-limit"},
		Short:   "서브모듈 브랜치 범위 설정 (특정 브랜치만 표시)",
		Long: `서브모듈의 브랜치 범위를 설정하여 선택한 브랜치만 표시되도록 합니다.
브랜치명을 입력하면 로컬과 origin 브랜치가 모두 필터링됩니다.

사용 예시:
  ga opt submodule set-branch-scope                # 대화형 모드
  ga opt submodule sbs main develop                # 짧은 별칭 사용
  ga opt submodule scope feature/test              # feature 브랜치만 표시
  ga opt submodule sbs main -q                     # quiet 모드로 자동 실행`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runSubmoduleSetBranchScope(args)
		},
	}
	
	// -q 플래그 추가
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runSubmoduleSetBranchScope(args []string) {
	// 서브모듈 존재 확인
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		fmt.Println("❌ 서브모듈이 없습니다")
		return
	}
	
	fmt.Println("\n🔧 서브모듈 브랜치 범위 설정")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 서브모듈 목록 가져오기
	cmd := exec.Command("git", "submodule", "foreach", "--quiet", "echo $path")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ 서브모듈 목록을 가져올 수 없습니다: %v\n", err)
		return
	}
	
	submodulePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(submodulePaths) == 0 || (len(submodulePaths) == 1 && submodulePaths[0] == "") {
		fmt.Println("ℹ️  초기화된 서브모듈이 없습니다.")
		return
	}
	
	// args가 있으면 바로 처리
	if len(args) > 0 {
		var branches []string
		for _, arg := range args {
			branch := strings.TrimSpace(arg)
			if branch != "" {
				// origin/ 접두사 제거
				branch = strings.TrimPrefix(branch, "origin/")
				if !utils.Contains(branches, branch) {
					branches = append(branches, branch)
				}
			}
		}
		
		if len(branches) > 0 {
			applySubmoduleBranchFilter(submodulePaths, branches)
			return
		}
	}
	
	// 현재 필터 설정 확인
	currentFilters := getSubmoduleFilters(submodulePaths)
	if len(currentFilters) > 0 {
		fmt.Println("\n📋 현재 필터링된 서브모듈:")
		for path, branches := range currentFilters {
			fmt.Printf("   • %s: %s\n", path, strings.Join(branches, ", "))
		}
		fmt.Println()
	}
	
	// 대화형 모드
	interactiveSubmoduleFilterMode(submodulePaths)
}

func getSubmoduleFilters(submodulePaths []string) map[string][]string {
	filters := make(map[string][]string)
	
	// 현재 config에서 branch_scope 읽기
	submoduleScope := config.GetSubmoduleBranchScope()
	if len(submoduleScope) > 0 {
		// 모든 서브모듈에 동일하게 적용
		for _, path := range submodulePaths {
			if path != "" {
				filters[path] = submoduleScope
			}
		}
	}
	
	return filters
}

func interactiveSubmoduleFilterMode(submodulePaths []string) {
	reader := bufio.NewReader(os.Stdin)
	
	// 모든 서브모듈의 브랜치를 수집 (중복 제거)
	branchMap := make(map[string]bool)
	
	fmt.Println("\n🔍 서브모듈 브랜치 확인 중...")
	for _, path := range submodulePaths {
		if path == "" {
			continue
		}
		
		// 로컬 브랜치
		localCmd := exec.Command("git", "-C", path, "branch", "--format=%(refname:short)")
		localOutput, _ := localCmd.Output()
		localBranches := strings.Split(strings.TrimSpace(string(localOutput)), "\n")
		
		for _, branch := range localBranches {
			if branch != "" && branch != "HEAD" {
				branchMap[branch] = true
			}
		}
		
		// 원격 브랜치
		remoteCmd := exec.Command("git", "-C", path, "branch", "-r", "--format=%(refname:short)")
		remoteOutput, _ := remoteCmd.Output()
		remoteBranches := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")
		
		for _, branch := range remoteBranches {
			if branch != "" && !strings.Contains(branch, "HEAD") {
				// origin/ 제거
				branch = strings.TrimPrefix(branch, "origin/")
				branchMap[branch] = true
			}
		}
	}
	
	// 브랜치 목록 생성
	var allBranches []string
	for branch := range branchMap {
		allBranches = append(allBranches, branch)
	}
	
	if len(allBranches) == 0 {
		fmt.Println("\n⚠️ 서브모듈에 브랜치가 없습니다")
		return
	}
	
	fmt.Println("\n📋 전체 브랜치 목록 (모든 서브모듈):")
	for i, branch := range allBranches {
		fmt.Printf("%2d. %s\n", i+1, branch)
	}
	
	fmt.Println("\n필터링할 브랜치를 선택하세요:")
	fmt.Println("• 단일 선택: 번호 또는 브랜치명 입력")
	fmt.Println("• 다중 선택: 공백으로 구분 (예: 1 3 5 또는 main develop)")
	fmt.Println("• 취소: q 또는 quit")
	fmt.Print("\n입력: ")
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "q" || input == "quit" || input == "" {
		fmt.Println("\n✨ 작업이 취소되었습니다")
		return
	}
	
	// 입력 파싱 (공백으로 구분)
	var selectedBranches []string
	parts := strings.Fields(input)  // 공백으로 분리
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		
		// 숫자인지 확인
		if idx := parseIndex(part); idx > 0 && idx <= len(allBranches) {
			branch := allBranches[idx-1]
			// origin/ 제거
			branch = strings.TrimPrefix(branch, "origin/")
			if !utils.Contains(selectedBranches, branch) {
				selectedBranches = append(selectedBranches, branch)
			}
		} else if part != "" {
			// 브랜치명 직접 입력
			branch := strings.TrimPrefix(part, "origin/")
			if !utils.Contains(selectedBranches, branch) {
				selectedBranches = append(selectedBranches, branch)
			}
		}
	}
	
	if len(selectedBranches) == 0 {
		fmt.Println("\n⚠️ 선택된 브랜치가 없습니다")
		return
	}
	
	// 브랜치 필터 적용
	applySubmoduleBranchFilter(submodulePaths, selectedBranches)
}

func applySubmoduleBranchFilter(submodulePaths []string, branches []string) {
	// 모든 서브모듈에 동일한 브랜치 필터 적용
	successCount := 0
	failCount := 0
	
	// .gaconfig/config.yaml에도 서브모듈 브랜치 스코프 저장
	if err := config.SetSubmoduleBranchScope(branches); err != nil {
		fmt.Printf("⚠️ config.yaml 서브모듈 브랜치 스코프 설정 실패: %v\n", err)
	}
	
	for _, path := range submodulePaths {
		if path == "" {
			continue
		}
		
		// 서브모듈의 fetch refspec 설정
		if err := utils.SetFetchRefspecForSubmodule(path, branches); err != nil {
			fmt.Printf("⚠️  %s fetch refspec 설정 실패: %v\n", path, err)
			failCount++
			continue
		}
		
		successCount++
	}
	
	fmt.Println("\n✅ 서브모듈 브랜치 범위가 설정되었습니다")
	fmt.Println("\n📋 필터링된 브랜치:")
	for _, branch := range branches {
		fmt.Printf("   • %s (로컬 및 origin/%s)\n", branch, branch)
	}
	
	fmt.Printf("\n📊 적용 결과:\n")
	fmt.Printf("   • 성공: %d개 서브모듈\n", successCount)
	if failCount > 0 {
		fmt.Printf("   • 실패: %d개 서브모듈\n", failCount)
	}
	
	// 실제 존재하는 브랜치 확인 (첫 번째 서브모듈 기준)
	if len(submodulePaths) > 0 && submodulePaths[0] != "" {
		checkBranchExistence(submodulePaths[0], branches)
	}
	
	fmt.Println("\n💡 팁:")
	fmt.Println("   • 필터를 제거하려면 'ga opt submodule clear-filter-branch' 명령을 사용하세요")
	fmt.Println("   • 이 설정은 모든 서브모듈에 동일하게 적용됩니다")
}

func checkBranchExistence(submodulePath string, branches []string) {
	// 로컬 브랜치 가져오기
	localCmd := exec.Command("git", "-C", submodulePath, "branch", "--format=%(refname:short)")
	localOutput, _ := localCmd.Output()
	localBranches := strings.Split(strings.TrimSpace(string(localOutput)), "\n")
	
	// 원격 브랜치 가져오기
	remoteCmd := exec.Command("git", "-C", submodulePath, "branch", "-r", "--format=%(refname:short)")
	remoteOutput, _ := remoteCmd.Output()
	remoteBranches := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")
	
	// 맵으로 변환
	localMap := make(map[string]bool)
	for _, branch := range localBranches {
		if branch != "" {
			localMap[branch] = true
		}
	}
	
	remoteMap := make(map[string]bool)
	for _, branch := range remoteBranches {
		if branch != "" {
			// origin/ 제거하고 저장
			branch = strings.TrimPrefix(branch, "origin/")
			remoteMap[branch] = true
		}
	}
	
	fmt.Println("\n🔍 실제 필터링 대상:")
	for _, branch := range branches {
		hasLocal := localMap[branch]
		hasRemote := remoteMap[branch]
		
		if hasLocal && hasRemote {
			fmt.Printf("   • %s (로컬 ✓, 원격 ✓)\n", branch)
		} else if hasLocal {
			fmt.Printf("   • %s (로컬 ✓)\n", branch)
		} else if hasRemote {
			fmt.Printf("   • %s (원격 ✓)\n", branch)
		} else {
			fmt.Printf("   • %s (⚠️ 아직 존재하지 않음)\n", branch)
		}
	}
}

func parseIndex(s string) int {
	var idx int
	fmt.Sscanf(s, "%d", &idx)
	return idx
}