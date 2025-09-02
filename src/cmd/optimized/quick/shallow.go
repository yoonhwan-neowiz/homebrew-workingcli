package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewShallowCmd creates the Shallow Depth command
func NewShallowCmd() *cobra.Command {
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:   "shallow [depth]",
		Short: "히스토리를 지정된 depth로 줄이기",
		Long: `히스토리를 지정된 개수의 커밋만 유지하도록 줄입니다.
depth를 지정하지 않으면 기본값 1(최신 1개 커밋)로 설정됩니다.
과거 히스토리가 필요 없는 경우 디스크 공간을 절약할 수 있습니다.

예시:
  ga opt quick shallow        # depth=1로 설정 (기본값)
  ga opt quick shallow 5      # 최근 5개 커밋만 유지
  ga opt quick shallow 10     # 최근 10개 커밋만 유지
  ga opt quick shallow 5 -q   # quiet 모드로 자동 실행`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// quiet 모드 설정
			if quietMode {
				utils.SetQuietMode(true)
			}
			
			depth := 1
			if len(args) > 0 {
				if d, err := strconv.Atoi(args[0]); err == nil && d > 0 {
					depth = d
				} else {
					fmt.Println("❌ 잘못된 depth 값입니다. 양의 정수를 입력해주세요.")
					os.Exit(1)
				}
			}
			runShallow(depth)
		},
	}
	
	// -q 플래그 추가
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "자동 실행 모드 (확인 없음)")
	
	return cmd
}

                                              
func runShallow(targetDepth int) {
	// 색상 설정
	titleStyle := color.New(color.FgCyan, color.Bold)
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	
	titleStyle.Printf("\n✂️  히스토리 최소화 (Shallow depth=%d)\n", targetDepth)
	titleStyle.Println("=" + strings.Repeat("=", 39))
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		errorStyle.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. 현재 상태 백업 권고
	fmt.Println("\n📌 현재 상태:")
	
	// Shallow 상태 확인
	shallowInfo := utils.GetShallowInfo()
	isShallow := shallowInfo["isShallow"].(bool)
	
	if isShallow {
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("   ├─ Shallow 상태: %s\n", warningStyle.Sprint("이미 활성"))
			fmt.Printf("   └─ 현재 커밋 수: %s개\n", boldStyle.Sprint(depth))
			
			if depth == targetDepth {
				infoStyle.Printf("\n✅ 이미 목표 히스토리(depth=%d)를 유지하고 있습니다.\n", targetDepth)
				return
			} else if depth < targetDepth {
				warningStyle.Printf("\n⚠️  현재 depth(%d)가 목표 depth(%d)보다 작습니다.\n", depth, targetDepth)
				warningStyle.Println("   더 많은 히스토리가 필요하면 expand 명령어를 사용하세요.")
				return
			}
		}
	} else {
		// 전체 커밋 수 확인
		cmd := exec.Command("git", "rev-list", "--count", "HEAD")
		if output, err := cmd.Output(); err == nil {
			count := strings.TrimSpace(string(output))
			fmt.Printf("   ├─ Shallow 상태: %s\n", infoStyle.Sprint("비활성"))
			fmt.Printf("   └─ 전체 커밋 수: %s개\n", boldStyle.Sprint(count))
		}
	}
	
	// 3. 디스크 사용량 확인 (변환 전)
	diskUsageBefore := utils.GetDiskUsage()
	fmt.Println("\n💾 현재 디스크 사용량:")
	if gitSize, ok := diskUsageBefore["git"]; ok {
		fmt.Printf("   └─ .git 폴더: %s\n", boldStyle.Sprint(gitSize))
	}
	
	// 4. 사용자 경고 및 확인
	warningStyle.Println("\n⚠️  주의사항:")
	if targetDepth == 1 {
		warningStyle.Println("   • 과거 커밋 히스토리가 제거됩니다 (최신 1개만 유지)")
	} else {
		warningStyle.Printf("   • 최근 %d개 커밋만 유지됩니다\n", targetDepth)
	}
	warningStyle.Println("   • git blame이나 과거 조회가 제한됩니다")
	warningStyle.Println("   • 이 작업은 되돌릴 수 있습니다 (unshallow)")
	
	confirmMsg := fmt.Sprintf("\n히스토리를 depth=%d로 줄이시겠습니까?", targetDepth)
	// shallow는 안전한 작업이므로 quiet 모드에서 자동 수락
	if !utils.ConfirmForce(confirmMsg) {
		fmt.Println("취소되었습니다.")
		return
	}
	
	// 5. 현재 브랜치 확인
	currentBranch := utils.GetCurrentBranch()
	fmt.Printf("\n🌿 현재 브랜치: %s\n", boldStyle.Sprint(currentBranch))
	
	// 6. Shallow 변환 실행
	fmt.Printf("\n🔄 히스토리를 depth=%d로 조정 중... ", targetDepth)
	
	// git pull --depth=N 실행
	depthStr := strconv.Itoa(targetDepth)
	cmd := exec.Command("git", "pull", "--depth="+depthStr)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// 에러 처리
		if strings.Contains(string(output), "shallow") {
			// 이미 shallow인 경우 다른 방법 시도
			fmt.Print("(대체 방법 시도) ")
			
			// git fetch --depth=N으로 재시도
			cmd = exec.Command("git", "fetch", "--depth="+depthStr)
			output, err = cmd.CombinedOutput()
			
			if err != nil {
				errorStyle.Println("실패")
				errorStyle.Printf("❌ 오류: %s\n", strings.TrimSpace(string(output)))
				os.Exit(1)
			}
		} else {
			errorStyle.Println("실패")
			errorStyle.Printf("❌ 오류: %s\n", strings.TrimSpace(string(output)))
			os.Exit(1)
		}
	}
	
	infoStyle.Println("완료")
	
	// 7. Git GC 실행으로 불필요한 객체 정리
	fmt.Print("🧹 불필요한 객체 정리 중... ")
	
	cmd = exec.Command("git", "gc", "--prune=now", "--aggressive")
	if err := cmd.Run(); err != nil {
		warningStyle.Println("부분 성공")
		fmt.Println("   └─ 일부 객체가 정리되지 않았을 수 있습니다")
	} else {
		infoStyle.Println("완료")
	}
	
	// 8. 결과 확인
	fmt.Println("\n📊 최소화 결과:")
	
	// Shallow 상태 재확인
	shallowInfo = utils.GetShallowInfo()
	isShallow = shallowInfo["isShallow"].(bool)
	
	if isShallow {
		infoStyle.Println("   ├─ Shallow 상태: 활성화됨")
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("   └─ 유지된 커밋 수: %s개\n", boldStyle.Sprint(depth))
		}
	} else {
		warningStyle.Println("   └─ Shallow 변환이 완전히 적용되지 않았습니다")
	}
	
	// 9. 디스크 사용량 비교
	diskUsageAfter := utils.GetDiskUsage()
	fmt.Println("\n💾 디스크 사용량 변화:")
	
	if gitSizeBefore, ok1 := diskUsageBefore["git"]; ok1 {
		if gitSizeAfter, ok2 := diskUsageAfter["git"]; ok2 {
			fmt.Printf("   ├─ 변환 전: %s\n", gitSizeBefore)
			fmt.Printf("   ├─ 변환 후: %s\n", boldStyle.Sprint(gitSizeAfter))
			
			// 간단한 크기 비교 (문자열 비교)
			if gitSizeBefore != gitSizeAfter {
				infoStyle.Println("   └─ ✅ 디스크 공간이 절약되었습니다")
			}
		}
	}
	
	// 10. 추가 안내
	fmt.Println("\n💡 팁:")
	fmt.Println("   • 과거 히스토리가 필요한 경우: ga opt advanced unshallow")
	fmt.Println("   • 더 많은 커밋이 필요한 경우: ga opt quick expand [개수]")
	fmt.Println("   • 현재 상태 확인: ga opt advanced check-shallow")
	
	infoStyle.Println("\n✅ 히스토리 최소화가 완료되었습니다.")
}