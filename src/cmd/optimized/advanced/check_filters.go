package advanced

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewCheckFilterCmd creates the Check Filter command
func NewCheckFilterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-filter",
		Short: "현재 필터 설정 확인",
		Long: `Partial Clone 필터와 브랜치 필터 설정을 확인합니다.
현재 적용된 모든 필터 상태를 진단합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runCheckFilter()
		},
	}
}

func runCheckFilter() {
	// 색상 설정
	titleStyle := color.New(color.FgCyan, color.Bold)
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	dimStyle := color.New(color.Faint)
	
	titleStyle.Println("\n🔍 필터 설정 확인")
	titleStyle.Println("=" + strings.Repeat("=", 39))
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		errorStyle.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. Partial Clone 필터 확인
	fmt.Println("\n📦 Partial Clone 필터:")
	
	// Global 필터 (remote.origin.partialclonefilter)
	partialFilter := utils.GetPartialCloneFilter()
	if partialFilter != "" {
		warningStyle.Println("   ├─ 상태: 활성")
		fmt.Printf("   ├─ 전역 필터: %s\n", boldStyle.Sprint(partialFilter))
		
		// 필터 해석
		if strings.HasPrefix(partialFilter, "blob:limit=") {
			sizeLimit := strings.TrimPrefix(partialFilter, "blob:limit=")
			fmt.Printf("   ├─ 제한 크기: %s 이상의 파일 제외\n", boldStyle.Sprint(sizeLimit))
		} else if strings.HasPrefix(partialFilter, "blob:none") {
			fmt.Println("   ├─ 모든 blob 제외 (트리만 포함)")
		} else if strings.HasPrefix(partialFilter, "tree:") {
			depth := strings.TrimPrefix(partialFilter, "tree:")
			fmt.Printf("   ├─ 트리 깊이: %s\n", boldStyle.Sprint(depth))
		}
		
		// Promisor 설정 확인
		cmd := exec.Command("git", "config", "remote.origin.promisor")
		if output, err := cmd.Output(); err == nil {
			promisor := strings.TrimSpace(string(output))
			if promisor == "true" {
				fmt.Printf("   └─ Promisor: %s\n", infoStyle.Sprint("활성"))
			}
		}
	} else {
		infoStyle.Println("   └─ 상태: 비활성 (모든 객체 포함)")
	}
	
	// 3. 브랜치별 필터 확인
	fmt.Println("\n🌿 브랜치별 필터:")
	
	// git config --get-regexp으로 브랜치별 필터 찾기
	cmd := exec.Command("git", "config", "--get-regexp", "branch\\..*\\.partialclonefilter")
	output, err := cmd.Output()
	
	if err == nil && len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}
			
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				config := parts[0]
				filter := strings.Join(parts[1:], " ")
				
				// 브랜치 이름 추출
				branchName := strings.TrimPrefix(config, "branch.")
				branchName = strings.TrimSuffix(branchName, ".partialclonefilter")
				
				prefix := "├─"
				if i == len(lines)-1 {
					prefix = "└─"
				}
				
				fmt.Printf("   %s %s: %s\n", prefix, boldStyle.Sprint(branchName), filter)
			}
		}
	} else {
		dimStyle.Println("   └─ (브랜치별 필터 없음)")
	}
	
	// 4. 커스텀 브랜치 필터 확인 (filter-branch 명령어용)
	fmt.Println("\n🎯 커스텀 브랜치 필터:")
	
	branchScope := config.GetBranchScope()
	if len(branchScope) > 0 {
		warningStyle.Println("   ├─ 상태: 활성")
		fmt.Printf("   ├─ 필터링된 브랜치: %s\n", boldStyle.Sprint(strings.Join(branchScope, ", ")))
		
		// 숨겨진 브랜치 계산
		allBranches := utils.GetAllUniqueBranches()
		hiddenCount := 0
		var hiddenBranches []string
		
		for _, branch := range allBranches {
			if !utils.Contains(branchScope, branch) {
				hiddenCount++
				hiddenBranches = append(hiddenBranches, branch)
			}
		}
		
		if hiddenCount > 0 {
			fmt.Printf("   ├─ 숨겨진 브랜치 수: %s개\n", warningStyle.Sprint(hiddenCount))
			
			// 처음 5개만 표시
			if len(hiddenBranches) <= 5 {
				fmt.Printf("   └─ 숨겨진 브랜치: %s\n", dimStyle.Sprint(strings.Join(hiddenBranches, ", ")))
			} else {
				first5 := hiddenBranches[:5]
				fmt.Printf("   └─ 숨겨진 브랜치 (일부): %s ...\n", dimStyle.Sprint(strings.Join(first5, ", ")))
			}
		}
	} else {
		infoStyle.Println("   └─ 상태: 비활성 (모든 브랜치 표시)")
	}
	
	// 5. Sparse Checkout 필터
	fmt.Println("\n📂 Sparse Checkout 필터:")
	
	sparseInfo := utils.GetSparseCheckoutInfo()
	if enabled, ok := sparseInfo["enabled"].(bool); ok && enabled {
		warningStyle.Println("   ├─ 상태: 활성")
		
		if count, ok := sparseInfo["count"].(int); ok {
			fmt.Printf("   ├─ 활성 경로 수: %s개\n", boldStyle.Sprint(count))
		}
		
		// Cone mode 확인
		coneMode := utils.CheckConeMode()
		if coneMode {
			fmt.Printf("   ├─ 모드: %s\n", infoStyle.Sprint("Cone"))
		} else {
			fmt.Printf("   ├─ 모드: %s\n", warningStyle.Sprint("Non-cone"))
		}
		
		// 경로 목록 (최대 5개)
		if paths, ok := sparseInfo["paths"].([]string); ok && len(paths) > 0 {
			fmt.Println("   └─ 활성 경로:")
			max := 5
			if len(paths) < max {
				max = len(paths)
			}
			
			for i := 0; i < max; i++ {
				fmt.Printf("       • %s\n", paths[i])
			}
			
			if len(paths) > 5 {
				fmt.Printf("       • ... 외 %d개\n", len(paths)-5)
			}
		}
	} else {
		infoStyle.Println("   └─ 상태: 비활성 (모든 파일 포함)")
	}
	
	// 6. 필터 영향 분석
	fmt.Println("\n💡 필터 영향:")
	
	hasFilters := partialFilter != "" || len(branchScope) > 0 || 
	              (sparseInfo["enabled"] != nil && sparseInfo["enabled"].(bool))
	
	if hasFilters {
		fmt.Println("   활성화된 필터로 인해:")
		
		if partialFilter != "" {
			fmt.Printf("   • 큰 파일이 제외됨 (%s)\n", partialFilter)
		}
		if len(branchScope) > 0 {
			fmt.Printf("   • 일부 브랜치가 숨겨짐 (%d개)\n", len(branchScope))
		}
		if sparseInfo["enabled"] != nil && sparseInfo["enabled"].(bool) {
			fmt.Println("   • 작업 디렉토리가 부분적으로만 체크아웃됨")
		}
	} else {
		infoStyle.Println("   • 모든 필터가 비활성 상태입니다")
		fmt.Println("   • 전체 저장소 내용을 사용할 수 있습니다")
	}
	
	// 7. 권장 사항
	fmt.Println("\n🔧 관련 명령어:")
	fmt.Println("   • Partial Clone 필터 제거: ga opt quick expand-filter")
	fmt.Println("   • 브랜치 필터 설정: ga opt workspace filter-branch")
	fmt.Println("   • 브랜치 필터 제거: ga opt workspace clear-filter-branch")
	fmt.Println("   • Sparse Checkout 경로 추가: ga opt workspace expand-path")
	
	fmt.Println()
}