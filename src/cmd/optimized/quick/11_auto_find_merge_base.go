package quick

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewAutoFindMergeBaseCmd creates the Auto Find Merge Base command
func NewAutoFindMergeBaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auto-find-merge-base",
		Short: "브랜치 병합점 자동 찾기",
		Long: `두 브랜치가 만나는 공통 조상 커밋(merge-base)을 자동으로 찾습니다.
필요 시 히스토리를 자동 확장하며 병합 가능성을 판단하는 기준점을 제공합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runAutoFindMergeBase()
		},
	}
}

func runAutoFindMergeBase() {
	fmt.Println("🔍 브랜치 병합점 자동 찾기")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ 오류: 현재 디렉토리는 Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 현재 브랜치 확인
	currentBranch := utils.GetCurrentBranch()
	fmt.Printf("📍 현재 브랜치: %s\n", currentBranch)
	
	// 비교할 브랜치 입력받기
	targetBranch := getBranchInput("비교할 브랜치명을 입력하세요")
	
	fmt.Printf("\n📊 %s와 %s의 병합점을 찾는 중...\n\n", currentBranch, targetBranch)
	
	// 머지베이스 찾기 시도
	mergeBase, depth, err := findMergeBase(currentBranch, targetBranch)
	
	if err != nil {
		fmt.Printf("❌ 오류: 병합점을 찾을 수 없습니다.\n")
		fmt.Printf("   상세: %v\n", err)
		os.Exit(1)
	}
	
	// 결과 표시
	fmt.Println("✅ 병합점 찾기 완료!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🔗 머지베이스: %s\n", mergeBase)
	
	if depth > 0 {
		fmt.Printf("📏 필요했던 depth: %d개 커밋\n", depth)
	} else {
		fmt.Println("📏 추가 히스토리 확장 없이 찾음")
	}
	
	// 커밋 정보 표시
	showCommitInfo(mergeBase)
	
	// 각 브랜치까지의 거리 표시
	showDistanceFromBase(currentBranch, targetBranch, mergeBase)
}

func getBranchInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	
	// 먼저 현재 브랜치 목록 표시
	showBranches()
	
	fmt.Printf("\n%s: ", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		fmt.Println("❌ 오류: 브랜치명을 입력해주세요.")
		os.Exit(1)
	}
	
	// 브랜치 존재 여부 확인
	if !branchExists(input) {
		// 원격 브랜치인지 확인
		if strings.Contains(input, "/") {
			fmt.Printf("ℹ️  원격 브랜치 %s를 사용합니다.\n", input)
		} else {
			fmt.Printf("⚠️  경고: %s 브랜치를 찾을 수 없습니다. 계속하시겠습니까?\n", input)
			if !utils.ConfirmWithDefault("계속 진행", false) {
				os.Exit(0)
			}
		}
	}
	
	return input
}

func showBranches() {
	fmt.Println("\n📋 사용 가능한 브랜치:")
	
	localBranches, remoteCount := utils.GetBranches()
	
	// 로컬 브랜치 표시
	if len(localBranches) > 0 {
		fmt.Println("  [로컬]")
		for _, branch := range localBranches {
			fmt.Printf("    %s\n", branch)
		}
	}
	
	// 원격 브랜치 개수 표시
	if remoteCount > 0 {
		fmt.Printf("  [원격] %d개 브랜치 (예: origin/main)\n", remoteCount)
	}
}

func branchExists(branch string) bool {
	return utils.BranchExists(branch)
}

func findMergeBase(branch1, branch2 string) (string, int, error) {
	totalDepth := 0
	maxAttempts := 10
	deepenStep := 10
	
	// 먼저 현재 상태에서 시도
	mergeBase, err := tryFindMergeBase(branch1, branch2)
	if err == nil && mergeBase != "" {
		return mergeBase, totalDepth, nil
	}
	
	// Shallow repository인 경우 점진적 확장
	if utils.IsShallowRepository() {
		fmt.Println("ℹ️  Shallow 저장소입니다. 히스토리를 점진적으로 확장합니다...")
		
		for i := 0; i < maxAttempts; i++ {
			fmt.Printf("   확장 중... (depth +%d)\n", deepenStep)
			
			// 히스토리 확장
			cmd := exec.Command("git", "fetch", "--deepen="+fmt.Sprintf("%d", deepenStep))
			err := cmd.Run()
			if err != nil {
				// unshallow 시도
				cmd = exec.Command("git", "fetch", "--unshallow")
				err = cmd.Run()
				if err != nil {
					break
				}
			}
			
			totalDepth += deepenStep
			
			// 다시 시도
			mergeBase, err = tryFindMergeBase(branch1, branch2)
			if err == nil && mergeBase != "" {
				return mergeBase, totalDepth, nil
			}
		}
	}
	
	// 마지막으로 전체 히스토리로 시도
	fmt.Println("ℹ️  전체 히스토리를 가져옵니다...")
	cmd := exec.Command("git", "fetch", "--unshallow")
	cmd.Run() // 이미 unshallow일 수 있으므로 에러 무시
	
	mergeBase, err = tryFindMergeBase(branch1, branch2)
	if err == nil && mergeBase != "" {
		return mergeBase, totalDepth, nil
	}
	
	return "", totalDepth, fmt.Errorf("공통 조상을 찾을 수 없습니다")
}

func tryFindMergeBase(branch1, branch2 string) (string, error) {
	return utils.FindMergeBase(branch1, branch2)
}

func showCommitInfo(commit string) {
	fmt.Printf("\n📝 커밋 정보:\n")
	
	// 커밋 정보 가져오기
	cmd := exec.Command("git", "log", "--oneline", "-n", "1", commit)
	output, err := cmd.Output()
	if err == nil {
		fmt.Printf("   %s\n", strings.TrimSpace(string(output)))
	}
	
	// 상세 정보
	cmd = exec.Command("git", "show", "--no-patch", "--format=%an <%ae>%n%ad", commit)
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) >= 2 {
			fmt.Printf("   작성자: %s\n", lines[0])
			fmt.Printf("   날짜: %s\n", lines[1])
		}
	}
}

func showDistanceFromBase(branch1, branch2, mergeBase string) {
	fmt.Printf("\n📏 머지베이스로부터의 거리:\n")
	
	// branch1까지의 거리
	distance1 := getDistanceFromBase(branch1, mergeBase)
	fmt.Printf("   %s: %s\n", branch1, distance1)
	
	// branch2까지의 거리
	distance2 := getDistanceFromBase(branch2, mergeBase)
	fmt.Printf("   %s: %s\n", branch2, distance2)
}

func getDistanceFromBase(branch, base string) string {
	count, err := utils.GetBranchDistance(branch, base)
	if err != nil {
		return "알 수 없음"
	}
	
	if count == 0 {
		return "동일함"
	}
	
	return fmt.Sprintf("%d개 커밋 ahead", count)
}