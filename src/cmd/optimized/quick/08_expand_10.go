package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewExpand10Cmd creates the unified Expand depth command
// Usage: ga opt quick expand [depth]
func NewExpand10Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand [depth]",
		Short: "히스토리 확장 (기본값: 10개)",
		Long: `현재 shallow 상태에서 지정한 개수만큼 커밋을 추가로 가져옵니다.
depth를 지정하지 않으면 기본값 10개를 확장합니다.

사용 예시:
  ga opt quick expand        # 10개 확장 (기본값)
  ga opt quick expand 10     # 10개 확장
  ga opt quick expand 50     # 50개 확장
  ga opt quick expand 100    # 100개 확장
  ga opt quick expand 66     # 66개 확장 (커스텀)`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// depth 파라미터 처리
			depth := 10 // 기본값
			
			if len(args) > 0 {
				var err error
				depth, err = strconv.Atoi(args[0])
				if err != nil || depth <= 0 {
					fmt.Println("❌ 오류: depth는 양수여야 합니다.")
					fmt.Printf("   입력값: %s\n", args[0])
					os.Exit(1)
				}
			}
			
			// 용도 설명 결정
			var purpose string
			switch {
			case depth <= 10:
				purpose = "최근 브랜치 분기 확인이나 작은 규모의 병합 작업"
			case depth <= 50:
				purpose = "중간 규모 작업이나 스프린트 단위 병합"
			case depth <= 100:
				purpose = "대규모 작업이나 릴리즈 브랜치 병합"
			default:
				purpose = fmt.Sprintf("대규모 히스토리 탐색 (커밋 %d개)", depth)
			}
			
			executeExpandHistory(depth, purpose)
		},
	}
}

// executeExpandHistory는 지정된 개수만큼 히스토리를 확장하는 함수
func executeExpandHistory(count int, purpose string) {
	fmt.Printf("📚 히스토리 %d개 확장 프로세스 시작\n", count)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   용도: %s\n", purpose)
	
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ 오류: 현재 디렉토리는 Git 저장소가 아닙니다.")
		fmt.Println("   Git 저장소 루트에서 실행해주세요.")
		os.Exit(1)
	}
	
	// 1. 현재 shallow 상태 확인
	fmt.Println("\n1️⃣  현재 히스토리 상태 확인 중...")
	shallowInfo := utils.GetShallowInfo()
	
	if !shallowInfo["isShallow"].(bool) {
		fmt.Println("ℹ️  현재 저장소는 shallow 상태가 아닙니다.")
		fmt.Println("   전체 히스토리를 이미 가지고 있습니다.")
		
		// 전체 커밋 개수 표시
		cmd := exec.Command("git", "rev-list", "--count", "HEAD")
		output, err := cmd.Output()
		if err == nil {
			totalCount := strings.TrimSpace(string(output))
			fmt.Printf("   전체 커밋 수: %s개\n", totalCount)
		}
		return
	}
	
	// 현재 depth 확인
	currentDepth := 0
	if depth, ok := shallowInfo["depth"].(int); ok {
		currentDepth = depth
		fmt.Printf("   현재 depth: %d개 커밋\n", currentDepth)
	}
	
	// 2. 히스토리 확장
	fmt.Printf("\n2️⃣  %d개 커밋 추가 다운로드 중...\n", count)
	fmt.Println("   (네트워크 상황에 따라 시간이 걸릴 수 있습니다)")
	
	// git fetch --deepen=N 실행
	fetchCmd := exec.Command("git", "fetch", fmt.Sprintf("--deepen=%d", count))
	output, err := fetchCmd.CombinedOutput()
	
	if err != nil {
		// 에러 처리
		errorMsg := string(output)
		
		// 이미 unshallow인 경우
		if strings.Contains(errorMsg, "unshallow") || strings.Contains(errorMsg, "no longer shallow") {
			fmt.Println("ℹ️  저장소가 더 이상 shallow 상태가 아닙니다.")
			fmt.Println("   전체 히스토리를 가지고 있습니다.")
			
			// 전체 커밋 수 표시
			countCmd := exec.Command("git", "rev-list", "--count", "HEAD")
			if countOutput, err := countCmd.Output(); err == nil {
				totalCount := strings.TrimSpace(string(countOutput))
				fmt.Printf("   전체 커밋 수: %s개\n", totalCount)
			}
			return
		}
		
		// 다른 에러
		fmt.Printf("❌ 히스토리 확장 실패: %v\n", err)
		if len(errorMsg) > 0 {
			fmt.Printf("   상세: %s\n", errorMsg)
		}
		
		// 대안 제시
		fmt.Println("\n💡 대안:")
		fmt.Println("   • 전체 히스토리 복원: ga opt advanced unshallow")
		fmt.Printf("   • 더 많은 커밋 확장: ga opt quick expand %d\n", count*2)
		os.Exit(1)
	}
	
	// 성공 메시지 파싱
	if len(output) > 0 {
		outputStr := string(output)
		if strings.Contains(outputStr, "deepening") {
			// deepening 메시지가 있으면 표시
			lines := strings.Split(strings.TrimSpace(outputStr), "\n")
			for _, line := range lines {
				if strings.Contains(line, "deepening") || strings.Contains(line, "commit") {
					fmt.Printf("   %s\n", line)
				}
			}
		}
	}
	
	// 3. 결과 확인
	fmt.Println("\n3️⃣  확장 결과 확인")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 새로운 depth 확인
	newDepthCmd := exec.Command("git", "rev-list", "--count", "HEAD")
	newDepthOutput, err := newDepthCmd.Output()
	if err == nil {
		newDepth, _ := strconv.Atoi(strings.TrimSpace(string(newDepthOutput)))
		fmt.Printf("✅ 히스토리 확장 완료\n")
		fmt.Printf("   이전 depth: %d개 커밋\n", currentDepth)
		fmt.Printf("   현재 depth: %d개 커밋\n", newDepth)
		fmt.Printf("   추가된 커밋: %d개\n", newDepth-currentDepth)
	}
	
	// shallow 상태 재확인
	isShallowCmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	shallowOutput, _ := isShallowCmd.Output()
	if strings.TrimSpace(string(shallowOutput)) == "false" {
		fmt.Println("\n🎉 저장소가 더 이상 shallow 상태가 아닙니다!")
		fmt.Println("   전체 히스토리를 보유하게 되었습니다.")
	} else {
		fmt.Printf("\n💡 추가 확장이 필요한 경우:\n")
		fmt.Printf("   • %d개 더 확장: ga opt quick expand %d\n", count, count)
		
		// 추천 확장 옵션
		fmt.Println("\n📊 권장 확장 옵션:")
		fmt.Println("   • 10개: ga opt quick expand 10    (브랜치 분기 확인)")
		fmt.Println("   • 50개: ga opt quick expand 50    (스프린트 병합)")
		fmt.Println("   • 100개: ga opt quick expand 100  (릴리즈 병합)")
		fmt.Println("   • 전체: ga opt advanced unshallow  (모든 히스토리)")
	}
}