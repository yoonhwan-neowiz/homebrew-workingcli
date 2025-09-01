package quick

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewSetBranchScopeCmd creates the Set Branch Scope command
func NewSetBranchScopeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set-branch-scope [브랜치1] [브랜치2] ...",
		Aliases: []string{"sbs", "scope", "branch-limit"},
		Short:   "브랜치 범위 설정 (특정 브랜치만 표시)",
		Long: `브랜치 범위를 설정하여 선택한 브랜치만 표시되도록 합니다.
브랜치명을 입력하면 로컬과 origin 브랜치가 모두 필터링됩니다.

사용 예시:
  ga opt quick set-branch-scope                # 대화형 모드
  ga opt quick sbs main develop                # 짧은 별칭 사용
  ga opt quick scope feature/test              # feature 브랜치만 표시`,
		Run: func(cmd *cobra.Command, args []string) {
			runSetBranchScope(args)
		},
	}
}

func runSetBranchScope(args []string) {
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다")
		return
	}

	fmt.Println("\n🔧 브랜치 범위 설정")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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
			applyBranchScope(branches)
			return
		}
	}

	// 현재 범위 설정 확인
	currentScope := utils.GetBranchScope()
	if len(currentScope) > 0 {
		fmt.Println("\n📋 현재 설정된 브랜치 범위:")
		for _, branch := range currentScope {
			fmt.Printf("   • %s\n", branch)
		}
		fmt.Println()
	}

	// 대화형 모드
	interactiveScopeMode()
}

func interactiveScopeMode() {
	reader := bufio.NewReader(os.Stdin)
	
	// 모든 브랜치 목록 가져오기 (중복 제거)
	allBranches := utils.GetAllUniqueBranches()
	
	if len(allBranches) == 0 {
		fmt.Println("\n⚠️ 브랜치가 없습니다")
		return
	}

	fmt.Println("\n📋 전체 브랜치 목록:")
	for i, branch := range allBranches {
		fmt.Printf("%2d. %s\n", i+1, branch)
	}

	fmt.Println("\n범위에 포함할 브랜치를 선택하세요:")
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

	// 브랜치 범위 적용
	applyBranchScope(selectedBranches)
}

func applyBranchScope(branches []string) {
	// Git config에 브랜치 범위 저장
	err := utils.SetBranchScope(branches)
	if err != nil {
		fmt.Printf("\n❌ 브랜치 범위 설정 실패: %v\n", err)
		return
	}

	fmt.Println("\n✅ 브랜치 범위가 설정되었습니다")
	fmt.Println("\n📋 설정된 브랜치 범위:")
	for _, branch := range branches {
		fmt.Printf("   • %s (로컬 및 origin/%s)\n", branch, branch)
	}

	// 실제 존재하는 브랜치 확인
	localBranches := utils.GetLocalBranches()
	remoteBranches := utils.GetRemoteBranches()
	
	fmt.Println("\n🔍 실제 범위 대상:")
	for _, branch := range branches {
		hasLocal := utils.Contains(localBranches, branch)
		hasRemote := utils.Contains(remoteBranches, "origin/"+branch)
		
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

	fmt.Println("\n💡 팁:")
	fmt.Println("   • 범위를 제거하려면 'ga opt quick clear-branch-scope' 명령을 사용하세요")
	fmt.Println("   • 이 설정은 프로젝트별로 저장됩니다")
}

func parseIndex(s string) int {
	var idx int
	fmt.Sscanf(s, "%d", &idx)
	return idx
}