package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewAllTagCmd creates the command to restore all remote tags
func NewAllTagCmd() *cobra.Command {
	var forceMode bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "alltag",
		Aliases: []string{"all-tags", "restore-tags"},
		Short:   "모든 원격 태그 복원",
		Long: `태그 fetch 차단을 해제하고 모든 원격 태그를 다시 가져옵니다.
notag로 제거한 태그들을 복원할 때 사용합니다.`,
		Example: `  ga opt quick alltag       # 모든 태그 복원
  ga opt quick alltag -f    # 확인 없이 강제 실행
  ga opt quick alltag -q    # 자동 실행 모드`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runAllTag(forceMode)
		},
	}
	
	cmd.Flags().BoolVarP(&forceMode, "force", "f", false, "확인 없이 강제 실행")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runAllTag(forceMode bool) {
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	fmt.Println("🏷️ 태그 복원 (All-Tag 모드)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 1. 현재 태그 개수 확인
	tagCountBefore := countLocalTagsForAllTag()
	fmt.Printf("📊 현재 태그: %d개\n", tagCountBefore)
	
	// 사용자 확인 (force 모드가 아닌 경우)
	if !forceMode {
		if !utils.ConfirmForce("\n원격 태그를 모두 복원하시겠습니까?") {
			fmt.Println("\n✨ 작업이 취소되었습니다")
			return
		}
	}
	
	// 2. .git 폴더 크기 측정 (복원 전)
	diskUsageBefore := utils.GetDiskUsage()
	sizeBefore := diskUsageBefore["git"]
	if sizeBefore == "" {
		sizeBefore = "unknown"
	}
	
	// 3. 태그 필터 설정 제거
	fmt.Print("🔓 태그 fetch 차단 해제 중...")
	if err := unblockTagFetch(); err != nil {
		// 설정이 없는 경우도 있으므로 경고만 표시
		fmt.Println(" (이미 해제됨)")
	} else {
		fmt.Println(" 완료")
	}
	
	// 4. 모든 원격 태그 fetch
	fmt.Print("📥 원격 태그 가져오는 중...")
	if err := fetchAllTags(); err != nil {
		fmt.Printf("\n❌ 태그 가져오기 실패: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(" 완료")
	
	// 5. 복원 후 태그 개수 확인
	tagCountAfter := countLocalTagsForAllTag()
	restoredCount := tagCountAfter - tagCountBefore
	
	// 6. .git 폴더 크기 측정 (복원 후)
	diskUsageAfter := utils.GetDiskUsage()
	sizeAfter := diskUsageAfter["git"]
	if sizeAfter == "" {
		sizeAfter = "unknown"
	}
	
	// 7. 결과 표시
	fmt.Println("\n✅ 태그 복원 완료")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("복원된 태그: %d개\n", restoredCount)
	fmt.Printf("총 태그 수: %d개\n", tagCountAfter)
	fmt.Printf(".git 크기 변화: %s → %s", sizeBefore, sizeAfter)
	
	fmt.Println()
	
	// 태그 fetch 상태 확인
	if isTagFetchBlockedForAllTag() {
		fmt.Println("태그 fetch: 차단됨 ❌")
		fmt.Println("\n⚠️  태그 fetch가 여전히 차단되어 있습니다.")
		fmt.Println("   다음 fetch/pull에서 태그가 업데이트되지 않을 수 있습니다.")
	} else {
		fmt.Println("태그 fetch: 활성화됨 ✅")
	}
	
	fmt.Println("\n💡 태그를 다시 제거하려면: ga opt quick notag")
}

// unblockTagFetch removes the tag fetch blocking configuration
func unblockTagFetch() error {
	cmd := exec.Command("git", "config", "--unset", "remote.origin.tagOpt")
	return cmd.Run()
}

// fetchAllTags fetches all tags from remote
func fetchAllTags() error {
	cmd := exec.Command("git", "fetch", "--tags")
	return cmd.Run()
}

// countLocalTags counts the number of local tags (reused from notag.go)
func countLocalTagsForAllTag() int {
	cmd := exec.Command("git", "tag")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	if len(output) == 0 {
		return 0
	}
	
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	// 빈 문자열 제외
	count := 0
	for _, tag := range tags {
		if strings.TrimSpace(tag) != "" {
			count++
		}
	}
	return count
}

// isTagFetchBlocked checks if tag fetching is blocked (reused from notag.go)
func isTagFetchBlockedForAllTag() bool {
	cmd := exec.Command("git", "config", "--get", "remote.origin.tagOpt")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) == "--no-tags"
}