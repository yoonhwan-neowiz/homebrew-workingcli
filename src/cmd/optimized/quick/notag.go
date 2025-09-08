package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewNoTagCmd creates the command to remove tags and block remote tag fetch
func NewNoTagCmd() *cobra.Command {
	var forceMode bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "notag",
		Aliases: []string{"no-tags", "remove-tags"},
		Short:   "로컬 태그 삭제 및 원격 태그 fetch 차단",
		Long: `모든 로컬 태그를 삭제하고 원격 태그 fetch를 차단합니다.
저장소 크기를 줄이고 불필요한 태그 다운로드를 방지합니다.`,
		Example: `  ga opt quick notag       # 태그 삭제 및 fetch 차단
  ga opt quick notag -f    # 확인 없이 강제 실행
  ga opt quick notag -q    # 자동 실행 모드`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runNoTag(forceMode)
		},
	}
	
	cmd.Flags().BoolVarP(&forceMode, "force", "f", false, "확인 없이 강제 실행")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runNoTag(forceMode bool) {
	// Git 저장소 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	fmt.Println("🏷️ 태그 최적화 (No-Tag 모드)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 1. 현재 태그 개수 확인
	tagCount := countLocalTags()
	if tagCount == 0 {
		fmt.Println("ℹ️  로컬 태그가 없습니다.")
	} else {
		fmt.Printf("📊 현재 태그: %d개\n", tagCount)
	}
	
	// 2. .git 폴더 크기 측정 (삭제 전)
	diskUsageBefore := utils.GetDiskUsage()
	sizeBefore := diskUsageBefore["git"]
	if sizeBefore == "" {
		sizeBefore = "unknown"
	}
	
	// 사용자 확인 (force 모드가 아닌 경우)
	if !forceMode && tagCount > 0 {
		if !utils.ConfirmForce(fmt.Sprintf("\n%d개의 태그를 삭제하고 원격 태그 fetch를 차단하시겠습니까?", tagCount)) {
			fmt.Println("\n✨ 작업이 취소되었습니다")
			return
		}
	}
	
	// 3. 모든 로컬 태그 삭제
	if tagCount > 0 {
		fmt.Print("🗑️  태그 삭제 중...")
		if err := deleteAllTags(); err != nil {
			fmt.Printf("\n❌ 태그 삭제 실패: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf(" 완료 (%d개 삭제됨)\n", tagCount)
	}
	
	// 4. 원격 태그 fetch 차단 설정
	fmt.Print("🚫 태그 fetch 차단 설정 중...")
	if err := blockTagFetch(); err != nil {
		fmt.Printf("\n❌ 태그 fetch 차단 설정 실패: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(" 완료")
	
	// 5. .git 폴더 크기 측정 (삭제 후)
	diskUsageAfter := utils.GetDiskUsage()
	sizeAfter := diskUsageAfter["git"]
	if sizeAfter == "" {
		sizeAfter = "unknown"
	}
	
	// 6. 결과 표시
	fmt.Println("\n✅ 태그 최적화 완료")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("삭제된 태그: %d개\n", tagCount)
	fmt.Printf(".git 크기 변화: %s → %s", sizeBefore, sizeAfter)
	
	fmt.Println()
	
	// 태그 fetch 차단 상태 확인
	if isTagFetchBlocked() {
		fmt.Println("태그 fetch: 차단됨 ❌")
	} else {
		fmt.Println("태그 fetch: 활성화됨 ✅")
	}
	
	fmt.Println("\n💡 원격 태그를 다시 받으려면: ga opt quick alltag")
}

// countLocalTags counts the number of local tags
func countLocalTags() int {
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

// deleteAllTags deletes all local tags
func deleteAllTags() error {
	// 태그 목록 가져오기
	cmd := exec.Command("git", "tag", "-l")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	
	if len(output) == 0 {
		return nil
	}
	
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// 각 태그 삭제
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		
		cmd := exec.Command("git", "tag", "-d", tag)
		if err := cmd.Run(); err != nil {
			// 개별 태그 삭제 실패는 무시하고 계속
			continue
		}
	}
	
	return nil
}

// blockTagFetch configures git to not fetch tags
func blockTagFetch() error {
	cmd := exec.Command("git", "config", "remote.origin.tagOpt", "--no-tags")
	return cmd.Run()
}

// isTagFetchBlocked checks if tag fetching is blocked
func isTagFetchBlocked() bool {
	cmd := exec.Command("git", "config", "--get", "remote.origin.tagOpt")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.TrimSpace(string(output)) == "--no-tags"
}