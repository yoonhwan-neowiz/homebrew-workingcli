package advanced

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewCheckMergeBaseCmd creates the Check Merge Base feasibility command
func NewCheckMergeBaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-merge-base",
		Short: "병합 베이스 확인",
		Long: `현재 히스토리에서 병합 베이스 존재 여부를 확인합니다.
불가능한 경우 auto-find-merge-base 명령 사용을 제안합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runCheckMergeBase()
		},
	}
}

func runCheckMergeBase() {
	fmt.Println("🔍 브랜치 병합 베이스 확인")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ 오류: 현재 디렉토리는 Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 현재 브랜치 확인
	currentBranch := utils.GetCurrentBranch()
	fmt.Printf("📍 현재 브랜치: %s\n", currentBranch)
	
	// 병합할 타겟 브랜치 입력받기
	targetBranch := getTargetBranchInput()
	
	fmt.Printf("\n📊 %s를 %s에 병합 가능한지 확인 중...\n\n", targetBranch, currentBranch)
	
	// 작업 중인 변경사항 확인
	stashed := false
	if utils.HasUncommittedChanges() {
		fmt.Println("⚠️  경고: 커밋되지 않은 변경사항이 있습니다.")
		fmt.Println("   병합 테스트를 진행하려면 변경사항을 처리해야 합니다.")
		fmt.Println()
		fmt.Println("   1. Stash - 변경사항을 임시 저장 (권장)")
		fmt.Println("   2. Reset - 변경사항을 버림 (주의!)")
		fmt.Println("   3. Cancel - 작업 취소")
		fmt.Println()
		
		choice := getUserChoice("선택하세요 (1/2/3)", []string{"1", "2", "3"})
		
		switch choice {
		case "1":
			// stash 저장
			cmd := exec.Command("git", "stash", "push", "-m", "check-merge temporary stash")
			err := cmd.Run()
			if err != nil {
				fmt.Println("❌ 오류: Stash 저장에 실패했습니다.")
				os.Exit(1)
			}
			stashed = true
			fmt.Println("   ✅ 변경사항이 임시 저장되었습니다.")
			defer func() {
				// stash 복원
				if stashed {
					cmd = exec.Command("git", "stash", "pop")
					err := cmd.Run()
					if err != nil {
						fmt.Println("   ⚠️  경고: Stash 복원에 실패했습니다.")
						fmt.Println("      'git stash pop' 명령을 수동으로 실행하세요.")
					} else {
						fmt.Println("   ✅ 변경사항이 복원되었습니다.")
					}
				}
			}()
		case "2":
			if !utils.ConfirmWithDefault("정말 모든 변경사항을 버리시겠습니까?", false) {
				fmt.Println("   작업을 취소합니다.")
				os.Exit(0)
			}
			// reset hard
			cmd := exec.Command("git", "reset", "--hard")
			err := cmd.Run()
			if err != nil {
				fmt.Println("❌ 오류: Reset 실행에 실패했습니다.")
				os.Exit(1)
			}
			fmt.Println("   ✅ 변경사항이 제거되었습니다.")
		case "3":
			fmt.Println("   작업을 취소합니다.")
			os.Exit(0)
		}
	}
	
	// 머지베이스 찾기 (확장 없이 현재 상태에서만)
	mergeBase, err := utils.FindMergeBase(currentBranch, targetBranch)
	
	if err != nil {
		fmt.Println("❌ 오류: 머지베이스를 찾을 수 없습니다.")
		fmt.Println("   두 브랜치가 관련이 없거나 히스토리가 부족합니다.")
		
		fmt.Println("\n💡 제안: 다음 방법을 시도해보세요:")
		fmt.Printf("   1. ga optimized quick auto-find-merge-base  # 자동으로 히스토리를 확장하며 머지베이스 찾기\n")
		if utils.IsShallowRepository() {
			fmt.Println("   2. ga optimized quick expand 100   # 수동으로 히스토리 확장")
			fmt.Println("   3. ga optimized quick to-full      # 전체 히스토리 복원 (시간 소요)")
		}
		os.Exit(1)
	}
	
	// 병합 시뮬레이션
	fmt.Println("🧪 병합 시뮬레이션 중...")
	conflictFiles, canMerge := simulateMerge(targetBranch)
	
	// 결과 표시
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 병합 가능성 분석 결과")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	fmt.Printf("🔗 머지베이스: %s\n", utils.GetShortCommit(mergeBase))
	
	// 브랜치 간 거리 표시
	showBranchDistance(currentBranch, targetBranch, mergeBase)
	
	fmt.Println("\n📊 병합 상태:")
	if canMerge {
		fmt.Println("✅ 병합 가능: 충돌 없이 병합할 수 있습니다!")
		fmt.Println("\n💡 병합 명령어:")
		fmt.Printf("   git merge %s\n", targetBranch)
	} else {
		fmt.Printf("⚠️  병합 시 충돌 발생: %d개 파일\n", len(conflictFiles))
		fmt.Println("\n🔥 충돌 파일 목록:")
		for _, file := range conflictFiles {
			fmt.Printf("   • %s\n", file)
		}
		
		fmt.Println("\n💡 권장 작업 순서:")
		fmt.Println("   1. 충돌 파일들을 미리 확인")
		fmt.Printf("   2. git merge %s 실행\n", targetBranch)
		fmt.Println("   3. 각 충돌을 수동으로 해결")
		fmt.Println("   4. git add <해결된 파일>")
		fmt.Println("   5. git commit")
	}
	
	// Shallow 저장소인 경우 추가 안내
	if utils.IsShallowRepository() {
		fmt.Println("\n📝 참고: Shallow 저장소 상태입니다.")
		fmt.Println("   복잡한 병합 작업 시 전체 히스토리가 필요할 수 있습니다.")
		fmt.Println("   필요시: ga optimized quick to-full")
	}
}


func getTargetBranchInput() string {
	reader := bufio.NewReader(os.Stdin)
	
	// 브랜치 목록 표시
	showAvailableBranches()
	
	fmt.Print("\n병합할 브랜치명을 입력하세요: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		fmt.Println("❌ 오류: 브랜치명을 입력해주세요.")
		os.Exit(1)
	}
	
	// 브랜치 존재 여부 확인
	if !checkBranchExists(input) {
		fmt.Printf("⚠️  경고: %s 브랜치를 찾을 수 없습니다.\n", input)
		if !utils.ConfirmWithDefault("계속 진행", false) {
			os.Exit(0)
		}
	}
	
	return input
}

func showAvailableBranches() {
	fmt.Println("\n📋 사용 가능한 브랜치:")
	
	localBranches, remoteCount := utils.GetBranches()
	currentBranch := utils.GetCurrentBranch()
	
	// 로컬 브랜치 표시
	if len(localBranches) > 0 {
		fmt.Println("  [로컬]")
		for _, branch := range localBranches {
			if branch == currentBranch {
				fmt.Printf("    %s (현재)\n", branch)
			} else {
				fmt.Printf("    %s\n", branch)
			}
		}
	}
	
	// 원격 브랜치 개수 표시
	if remoteCount > 0 {
		fmt.Printf("  [원격] %d개 브랜치\n", remoteCount)
	}
}

func checkBranchExists(branch string) bool {
	return utils.BranchExists(branch)
}



func simulateMerge(targetBranch string) ([]string, bool) {
	conflictFiles := []string{}
	
	// 병합 시뮬레이션 (실제로 병합하지 않음)
	cmd := exec.Command("git", "merge", "--no-commit", "--no-ff", targetBranch)
	output, err := cmd.CombinedOutput()
	
	// 병합 상태 확인
	if err != nil {
		// 충돌 발생 시
		if strings.Contains(string(output), "CONFLICT") {
			// 충돌 파일 목록 가져오기
			cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=U")
			conflictOutput, _ := cmd.Output()
			if len(conflictOutput) > 0 {
				files := strings.Split(strings.TrimSpace(string(conflictOutput)), "\n")
				for _, file := range files {
					if file != "" {
						conflictFiles = append(conflictFiles, file)
					}
				}
			}
		}
	}
	
	// 병합 취소 (시뮬레이션만 했으므로)
	cmd = exec.Command("git", "merge", "--abort")
	cmd.Run()
	
	// 충돌이 없으면 병합 가능
	return conflictFiles, len(conflictFiles) == 0
}


func getUserChoice(prompt string, validChoices []string) string {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Printf("%s: ", prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		for _, choice := range validChoices {
			if input == choice {
				return input
			}
		}
		
		fmt.Printf("❌ 잘못된 입력입니다. %s 중 하나를 선택하세요.\n", strings.Join(validChoices, ", "))
	}
}

func showBranchDistance(current, target, mergeBase string) {
	fmt.Println("\n📏 브랜치 간 거리:")
	
	// current 브랜치의 고유 커밋 수
	currentCount, _ := utils.GetBranchDistance(current, mergeBase)
	
	// target 브랜치의 고유 커밋 수
	targetCount, _ := utils.GetBranchDistance(target, mergeBase)
	
	fmt.Printf("   %s: 머지베이스로부터 %d개 커밋\n", current, currentCount)
	fmt.Printf("   %s: 머지베이스로부터 %d개 커밋\n", target, targetCount)
	
	// 전체 병합될 커밋 수
	if targetCount > 0 {
		fmt.Printf("   → 병합 시 %d개 커밋이 추가됩니다\n", targetCount)
	}
}