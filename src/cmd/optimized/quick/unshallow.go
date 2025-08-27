package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewUnshallowCmd creates the Unshallow command
func NewUnshallowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unshallow",
		Short: "히스토리 복원",
		Long: `전체 히스토리를 복원합니다.
과거 커밋 조회나 blame이 필요한 경우 사용합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runUnshallow()
		},
	}
}

func runUnshallow() {
	// 색상 설정
	titleStyle := color.New(color.FgCyan, color.Bold)
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	
	titleStyle.Println("\n📚 히스토리 복원 (Unshallow)")
	titleStyle.Println("=" + strings.Repeat("=", 39))
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		errorStyle.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. 현재 Shallow 상태 확인
	shallowInfo := utils.GetShallowInfo()
	isShallow := shallowInfo["isShallow"].(bool)
	
	if !isShallow {
		infoStyle.Println("✅ 이미 전체 히스토리를 가지고 있습니다.")
		
		// 전체 커밋 개수 확인
		cmd := exec.Command("git", "rev-list", "--count", "HEAD")
		if output, err := cmd.Output(); err == nil {
			count := strings.TrimSpace(string(output))
			fmt.Printf("   └─ 전체 커밋 수: %s개\n", boldStyle.Sprint(count))
		}
		return
	}
	
	// 3. 현재 depth 표시
	fmt.Println("\n📊 현재 상태:")
	if depth, ok := shallowInfo["depth"].(int); ok {
		fmt.Printf("   ├─ Shallow 상태: %s\n", warningStyle.Sprint("활성"))
		fmt.Printf("   └─ 현재 커밋 수: %s개 (shallow)\n", boldStyle.Sprint(depth))
	}
	
	// 4. 사용자 확인
	warningStyle.Println("\n⚠️  전체 히스토리를 다운로드합니다.")
	warningStyle.Println("   이 작업은 시간이 오래 걸리고 디스크 사용량이 증가합니다.")
	
	if !utils.Confirm("계속하시겠습니까?") {
		fmt.Println("취소되었습니다.")
		return
	}
	
	// 5. 전체 히스토리 복원
	fmt.Print("\n🔄 히스토리 복원 중... ")
	
	cmd := exec.Command("git", "fetch", "--unshallow")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// 이미 unshallow인 경우 에러가 발생할 수 있음
		if strings.Contains(string(output), "already have") || 
		   strings.Contains(string(output), "unshallow") {
			infoStyle.Println("완료")
			fmt.Println("   └─ 이미 전체 히스토리를 가지고 있습니다.")
		} else {
			errorStyle.Println("실패")
			errorStyle.Printf("❌ 오류: %s\n", strings.TrimSpace(string(output)))
			os.Exit(1)
		}
	} else {
		infoStyle.Println("완료")
		
		// 출력 내용 표시
		if len(output) > 0 {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, line := range lines {
				if line != "" {
					fmt.Printf("   └─ %s\n", line)
				}
			}
		}
	}
	
	// 6. 결과 확인
	fmt.Println("\n📊 복원 결과:")
	
	// Shallow 상태 재확인
	shallowInfo = utils.GetShallowInfo()
	isShallow = shallowInfo["isShallow"].(bool)
	
	if isShallow {
		warningStyle.Println("   ├─ Shallow 상태: 여전히 활성")
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("   └─ 커밋 수: %s개\n", boldStyle.Sprint(depth))
		}
	} else {
		infoStyle.Println("   ├─ Shallow 상태: 비활성 (전체 히스토리)")
		
		// 전체 커밋 수 확인
		cmd := exec.Command("git", "rev-list", "--count", "HEAD")
		if output, err := cmd.Output(); err == nil {
			count := strings.TrimSpace(string(output))
			fmt.Printf("   └─ 전체 커밋 수: %s개\n", boldStyle.Sprint(count))
		}
	}
	
	// 7. 디스크 사용량 확인
	diskUsage := utils.GetDiskUsage()
	fmt.Println("\n💾 디스크 사용량:")
	if gitSize, ok := diskUsage["git"]; ok {
		fmt.Printf("   ├─ .git 폴더: %s\n", boldStyle.Sprint(gitSize))
	}
	if totalSize, ok := diskUsage["total"]; ok {
		fmt.Printf("   └─ 프로젝트 전체: %s\n", boldStyle.Sprint(totalSize))
	}
	
	fmt.Println("\n✅ 히스토리 복원이 완료되었습니다.")
	fmt.Println("   이제 모든 과거 커밋을 조회할 수 있습니다.")
}