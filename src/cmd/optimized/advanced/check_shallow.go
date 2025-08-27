package advanced

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewCheckShallowCmd creates the Check Shallow command
func NewCheckShallowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-shallow",
		Short: "히스토리 상태 확인",
		Long: `현재 커밋 수와 shallow 포인트를 표시합니다.
히스토리 상태를 진단합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runCheckShallow()
		},
	}
}

func runCheckShallow() {
	// 색상 설정
	titleStyle := color.New(color.FgCyan, color.Bold)
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	dimStyle := color.New(color.Faint)
	
	titleStyle.Println("\n🔍 히스토리 상태 확인")
	titleStyle.Println("=" + strings.Repeat("=", 39))
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		errorStyle.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. Shallow 상태 확인
	fmt.Println("\n📊 Shallow 상태:")
	shallowInfo := utils.GetShallowInfo()
	isShallow := shallowInfo["isShallow"].(bool)
	
	if isShallow {
		warningStyle.Println("   ├─ 상태: Shallow 활성")
		
		// 현재 depth 확인
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("   ├─ 현재 depth: %s개 커밋\n", boldStyle.Sprint(depth))
		}
		
		// Grafted 커밋 확인 (.git/shallow 파일)
		shallowFile := filepath.Join(".git", "shallow")
		if data, err := os.ReadFile(shallowFile); err == nil {
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			fmt.Printf("   ├─ Grafted 커밋 수: %s개\n", boldStyle.Sprint(len(lines)))
			
			// 첫 번째 grafted 커밋 표시
			if len(lines) > 0 && lines[0] != "" {
				shortHash := lines[0]
				if len(shortHash) > 7 {
					shortHash = shortHash[:7]
				}
				fmt.Printf("   └─ 가장 오래된 커밋: %s\n", dimStyle.Sprint(shortHash))
			}
		}
	} else {
		infoStyle.Println("   └─ 상태: 전체 히스토리 (Unshallow)")
	}
	
	// 3. 커밋 히스토리 분석
	fmt.Println("\n📈 히스토리 정보:")
	
	// 전체 커밋 수
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	if output, err := cmd.Output(); err == nil {
		count := strings.TrimSpace(string(output))
		fmt.Printf("   ├─ 현재 브랜치 커밋 수: %s개\n", boldStyle.Sprint(count))
	}
	
	// 현재 브랜치
	currentBranch := utils.GetCurrentBranch()
	fmt.Printf("   ├─ 현재 브랜치: %s\n", boldStyle.Sprint(currentBranch))
	
	// 모든 브랜치의 커밋 수 (--all)
	cmd = exec.Command("git", "rev-list", "--count", "--all")
	if output, err := cmd.Output(); err == nil {
		allCount := strings.TrimSpace(string(output))
		fmt.Printf("   └─ 모든 브랜치 총 커밋: %s개\n", boldStyle.Sprint(allCount))
	}
	
	// 4. 히스토리 제한 확인
	fmt.Println("\n🚧 히스토리 제한:")
	
	if isShallow {
		warningStyle.Println("   ├─ ⚠️  다음 기능이 제한됩니다:")
		fmt.Println("   ├─ • git blame (일부 제한)")
		fmt.Println("   ├─ • git log --all (일부 제한)")
		fmt.Println("   ├─ • 과거 커밋 체크아웃 불가")
		fmt.Println("   └─ • 머지베이스 찾기 제한")
	} else {
		infoStyle.Println("   └─ ✅ 제한 없음 (전체 히스토리 사용 가능)")
	}
	
	// 5. 최근 커밋 목록
	fmt.Println("\n📝 최근 커밋 (최대 5개):")
	
	cmd = exec.Command("git", "log", "--oneline", "-n", "5")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}
			
			// 마지막 줄 처리
			prefix := "├─"
			if i == len(lines)-1 {
				prefix = "└─"
			}
			
			// 커밋 해시와 메시지 분리
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				hash := parts[0]
				message := parts[1]
				fmt.Printf("   %s %s %s\n", prefix, warningStyle.Sprint(hash), message)
			} else {
				fmt.Printf("   %s %s\n", prefix, line)
			}
		}
	}
	
	// 6. 관련 설정
	fmt.Println("\n⚙️  관련 Git 설정:")
	
	// fetch.depth 설정 확인
	cmd = exec.Command("git", "config", "fetch.depth")
	if output, err := cmd.Output(); err == nil && len(output) > 0 {
		depth := strings.TrimSpace(string(output))
		fmt.Printf("   ├─ fetch.depth: %s\n", boldStyle.Sprint(depth))
	} else {
		fmt.Println("   ├─ fetch.depth: (설정 안 됨)")
	}
	
	// clone.depth 설정 확인
	cmd = exec.Command("git", "config", "clone.depth")
	if output, err := cmd.Output(); err == nil && len(output) > 0 {
		depth := strings.TrimSpace(string(output))
		fmt.Printf("   └─ clone.depth: %s\n", boldStyle.Sprint(depth))
	} else {
		fmt.Println("   └─ clone.depth: (설정 안 됨)")
	}
	
	// 7. 권장 사항
	fmt.Println("\n💡 권장 사항:")
	
	if isShallow {
		fmt.Println("   • 전체 히스토리 필요 시: ga opt advanced unshallow")
		fmt.Println("   • 일부 커밋 추가 시: ga opt quick expand [개수]")
		fmt.Println("   • 머지베이스 찾기: ga opt quick auto-find-merge-base")
	} else {
		fmt.Println("   • 히스토리 최소화: ga opt advanced shallow")
		fmt.Println("   • 디스크 절약이 필요한 경우 shallow 권장")
	}
	
	fmt.Println()
}