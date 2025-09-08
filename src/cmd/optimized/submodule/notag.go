package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewNoTagCmd creates the command to remove submodule tags recursively
func NewNoTagCmd() *cobra.Command {
	var forceMode bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:     "notag",
		Aliases: []string{"no-tags", "remove-tags"},
		Short:   "서브모듈 태그 삭제 및 fetch 차단 (recursive)",
		Long: `모든 서브모듈의 로컬 태그를 재귀적으로 삭제하고 원격 태그 fetch를 차단합니다.
각 서브모듈의 저장소 크기를 줄이고 불필요한 태그 다운로드를 방지합니다.`,
		Example: `  ga opt submodule notag       # 서브모듈 태그 삭제
  ga opt submodule notag -f    # 확인 없이 강제 실행
  ga opt submodule notag -q    # 자동 실행 모드`,
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			runSubmoduleNoTag(forceMode)
		},
	}
	
	cmd.Flags().BoolVarP(&forceMode, "force", "f", false, "확인 없이 강제 실행")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

func runSubmoduleNoTag(forceMode bool) {
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
	submodules := getSubmodulePaths()
	
	fmt.Println("🏷️ 서브모듈 태그 최적화")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📦 서브모듈 개수: %d개\n", len(submodules))
	
	// 사용자 확인 (force 모드가 아닌 경우)
	if !forceMode {
		if !utils.ConfirmForce(fmt.Sprintf("\n%d개 서브모듈의 태그를 삭제하고 원격 태그 fetch를 차단하시겠습니까?", len(submodules))) {
			fmt.Println("\n✨ 작업이 취소되었습니다")
			return
		}
	}
	
	// 통계 변수
	totalDeletedTags := 0
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
		
		// 서브모듈 .git 크기 (삭제 전)
		gitPath := fmt.Sprintf("%s/.git", path)
		sizeBeforeCmd := exec.Command("du", "-sk", gitPath)
		sizeBeforeOutput, _ := sizeBeforeCmd.Output()
		var sizeBefore int64
		if len(sizeBeforeOutput) > 0 {
			fmt.Sscanf(string(sizeBeforeOutput), "%d", &sizeBefore)
			sizeBefore *= 1024 // KB to bytes
		}
		totalSizeBefore += sizeBefore
		
		// 태그 개수 확인
		tagCount := countSubmoduleTags(path)
		if tagCount > 0 {
			fmt.Printf("   📊 태그: %d개\n", tagCount)
			
			// 태그 삭제
			if err := deleteSubmoduleTags(path); err != nil {
				fmt.Printf("   ❌ 태그 삭제 실패: %v\n", err)
				failCount++
				continue
			}
			totalDeletedTags += tagCount
		} else {
			fmt.Println("   ℹ️  태그 없음")
		}
		
		// 태그 fetch 차단 설정
		if err := blockSubmoduleTagFetch(path); err != nil {
			fmt.Printf("   ❌ 태그 fetch 차단 실패: %v\n", err)
			failCount++
			continue
		}
		
		// 서브모듈 .git 크기 (삭제 후)
		sizeAfterCmd := exec.Command("du", "-sk", gitPath)
		sizeAfterOutput, _ := sizeAfterCmd.Output()
		var sizeAfter int64
		if len(sizeAfterOutput) > 0 {
			fmt.Sscanf(string(sizeAfterOutput), "%d", &sizeAfter)
			sizeAfter *= 1024 // KB to bytes
		}
		totalSizeAfter += sizeAfter
		
		// 크기 변화 표시
		if sizeBefore > sizeAfter {
			reduction := float64(sizeBefore-sizeAfter) / float64(sizeBefore) * 100
			fmt.Printf("   ✅ 완료 (%.1f%% 감소)\n", reduction)
		} else {
			fmt.Println("   ✅ 완료")
		}
		
		successCount++
	}
	
	// 전체 결과 표시
	fmt.Println("\n✅ 서브모듈 태그 최적화 완료")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("처리된 서브모듈: %d개 (성공: %d, 실패: %d)\n", 
		len(submodules), successCount, failCount)
	fmt.Printf("총 삭제된 태그: %d개\n", totalDeletedTags)
	
	if totalSizeBefore > totalSizeAfter {
		reduction := totalSizeBefore - totalSizeAfter
		fmt.Printf("전체 크기 감소: %s\n", utils.HumanizeBytes(reduction))
	}
	
	fmt.Println("\n💡 서브모듈 태그를 복원하려면: ga opt submodule alltag")
}

// countSubmoduleTags counts tags in a submodule
func countSubmoduleTags(path string) int {
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

// deleteSubmoduleTags deletes all tags in a submodule
func deleteSubmoduleTags(path string) error {
	// 태그 목록 가져오기
	cmd := exec.Command("git", "-C", path, "tag", "-l")
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
		
		cmd := exec.Command("git", "-C", path, "tag", "-d", tag)
		if err := cmd.Run(); err != nil {
			// 개별 태그 삭제 실패는 무시하고 계속
			continue
		}
	}
	
	return nil
}

// blockSubmoduleTagFetch blocks tag fetching for a submodule
func blockSubmoduleTagFetch(path string) error {
	cmd := exec.Command("git", "-C", path, "config", "remote.origin.tagOpt", "--no-tags")
	return cmd.Run()
}

// getSubmodulePaths returns list of submodule paths
func getSubmodulePaths() []string {
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