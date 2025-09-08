package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewAllTagCmd creates the command to restore all submodule tags recursively
func NewAllTagCmd() *cobra.Command {
	var forceMode bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "alltag",
		Aliases: []string{"all-tags", "restore-tags"},
		Short:   "서브모듈 태그 복원 (recursive)",
		Long: `모든 서브모듈의 태그 fetch 차단을 해제하고 원격 태그를 재귀적으로 복원합니다.
notag로 제거한 서브모듈 태그들을 복원할 때 사용합니다.`,
		Example: `  ga opt submodule alltag       # 서브모듈 태그 복원
  ga opt submodule alltag -f    # 확인 없이 강제 실행
  ga opt submodule alltag -q    # 자동 실행 모드`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runSubmoduleAllTag(forceMode)
		},
	}
	
	cmd.Flags().BoolVarP(&forceMode, "force", "f", false, "확인 없이 강제 실행")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runSubmoduleAllTag(forceMode bool) {
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 서브모듈 목록 가져오기
	submoduleInfo := utils.GetSubmoduleInfo()
	submoduleCount, ok := submoduleInfo["count"].(int)
	if !ok || submoduleCount == 0 {
		fmt.Println("ℹ️  서브모듈이 없습니다.")
		return
	}
	
	// 서브모듈 경로 목록 가져오기
	submodules := getSubmodulePathsForAllTag()
	
	fmt.Println("🏷️ 서브모듈 태그 복원")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📦 서브모듈 개수: %d개\n", len(submodules))
	
	// 사용자 확인 (force 모드가 아닌 경우)
	if !forceMode {
		if !utils.ConfirmForce(fmt.Sprintf("\n%d개 서브모듈의 원격 태그를 모두 복원하시겠습니까?", len(submodules))) {
			fmt.Println("\n✨ 작업이 취소되었습니다")
			return
		}
	}
	
	// 통계 변수
	totalRestoredTags := 0
	totalSizeBefore := int64(0)
	totalSizeAfter := int64(0)
	successCount := 0
	failCount := 0
	
	// 각 서브모듈 처리
	for _, path := range submodules {
		name := path // 경로를 이름으로 사용
		
		fmt.Printf("\n📁 %s 처리 중...\n", name)
		
		// 서브모듈 디렉토리 확인
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("   ⚠️  경로를 찾을 수 없음: %s\n", path)
			failCount++
			continue
		}
		
		// 복원 전 태그 개수
		tagCountBefore := countSubmoduleTagsForAllTag(path)
		
		// 서브모듈 .git 크기 (복원 전)
		gitPath := fmt.Sprintf("%s/.git", path)
		sizeBeforeCmd := exec.Command("du", "-sk", gitPath)
		sizeBeforeOutput, _ := sizeBeforeCmd.Output()
		var sizeBefore int64
		if len(sizeBeforeOutput) > 0 {
			fmt.Sscanf(string(sizeBeforeOutput), "%d", &sizeBefore)
			sizeBefore *= 1024 // KB to bytes
		}
		totalSizeBefore += sizeBefore
		
		// 태그 fetch 차단 해제
		if err := unblockSubmoduleTagFetch(path); err != nil {
			// 설정이 없는 경우도 있으므로 경고만 표시
			fmt.Println("   ℹ️  태그 fetch 이미 활성화됨")
		}
		
		// 원격 태그 fetch
		fmt.Print("   📥 태그 가져오는 중...")
		if err := fetchSubmoduleTags(path); err != nil {
			fmt.Printf("\n   ❌ 태그 가져오기 실패: %v\n", err)
			failCount++
			continue
		}
		fmt.Println(" 완료")
		
		// 복원 후 태그 개수
		tagCountAfter := countSubmoduleTagsForAllTag(path)
		restoredCount := tagCountAfter - tagCountBefore
		if restoredCount > 0 {
			fmt.Printf("   📊 복원된 태그: %d개 (총 %d개)\n", restoredCount, tagCountAfter)
			totalRestoredTags += restoredCount
		} else {
			fmt.Printf("   ℹ️  새로운 태그 없음 (총 %d개)\n", tagCountAfter)
		}
		
		// 서브모듈 .git 크기 (복원 후)
		sizeAfterCmd := exec.Command("du", "-sk", gitPath)
		sizeAfterOutput, _ := sizeAfterCmd.Output()
		var sizeAfter int64
		if len(sizeAfterOutput) > 0 {
			fmt.Sscanf(string(sizeAfterOutput), "%d", &sizeAfter)
			sizeAfter *= 1024 // KB to bytes
		}
		totalSizeAfter += sizeAfter
		
		// 크기 변화 표시
		if sizeAfter > sizeBefore {
			increase := float64(sizeAfter-sizeBefore) / float64(sizeBefore) * 100
			fmt.Printf("   ✅ 완료 (%.1f%% 증가)\n", increase)
		} else {
			fmt.Println("   ✅ 완료")
		}
		
		successCount++
	}
	
	// 전체 결과 표시
	fmt.Println("\n✅ 서브모듈 태그 복원 완료")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("처리된 서브모듈: %d개 (성공: %d, 실패: %d)\n", 
		len(submodules), successCount, failCount)
	fmt.Printf("총 복원된 태그: %d개\n", totalRestoredTags)
	
	if totalSizeAfter > totalSizeBefore {
		increase := totalSizeAfter - totalSizeBefore
		fmt.Printf("전체 크기 증가: %s\n", utils.HumanizeBytes(increase))
	}
	
	fmt.Println("\n💡 서브모듈 태그를 제거하려면: ga opt submodule notag")
}

// countSubmoduleTagsForAllTag counts tags in a submodule
func countSubmoduleTagsForAllTag(path string) int {
	cmd := exec.Command("git", "-C", path, "tag")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	if len(output) == 0 {
		return 0
	}
	
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	count := 0
	for _, tag := range tags {
		if strings.TrimSpace(tag) != "" {
			count++
		}
	}
	return count
}

// unblockSubmoduleTagFetch removes tag fetch blocking for a submodule
func unblockSubmoduleTagFetch(path string) error {
	cmd := exec.Command("git", "-C", path, "config", "--unset", "remote.origin.tagOpt")
	return cmd.Run()
}

// fetchSubmoduleTags fetches all tags for a submodule
func fetchSubmoduleTags(path string) error {
	cmd := exec.Command("git", "-C", path, "fetch", "--tags")
	return cmd.Run()
}

// getSubmodulePathsForAllTag returns list of submodule paths
func getSubmodulePathsForAllTag() []string {
	var paths []string
	
	cmd := exec.Command("git", "submodule", "foreach", "--quiet", "echo $sm_path")
	output, err := cmd.Output()
	if err != nil {
		return paths
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	
	return paths
}