package quick

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	
	"workingcli/src/config"
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// NewExpandSlimCmd creates the Expand SLIM command
func NewExpandSlimCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-slim",
		Short: "선택적 파일/폴더 확장 (SLIM 유지)",
		Long: `SLIM 상태를 유지하면서 특정 파일이나 폴더를 선택적으로 확장합니다.
Sparse Checkout 목록에 경로를 추가하고 Partial Clone 필터를 우회하여 
대용량 파일도 다운로드합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runExpandSlim()
		},
	}
}

func runExpandSlim() {
	fmt.Println("\n🔸 SLIM 선택적 확장")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 현재 상태 확인
	mode := utils.GetOptimizationMode()
	if mode != "SLIM" {
		fmt.Println("⚠️  경고: 현재 FULL 모드입니다. SLIM 모드에서만 선택적 확장이 유용합니다.")
		if !utils.Confirm("계속 진행하시겠습니까?") {
			fmt.Println("⏹  작업을 취소했습니다.")
			return
		}
	}

	// Sparse Checkout 상태 확인
	sparseInfo := utils.GetSparseCheckoutInfo()
	if sparseInfo["enabled"].(bool) {
		fmt.Printf("\n📋 현재 Sparse Checkout 목록 (%d개 경로):\n", sparseInfo["count"])
		if paths, ok := sparseInfo["paths"].([]string); ok {
			for _, path := range paths {
				fmt.Printf("   • %s\n", path)
			}
		}
	} else {
		fmt.Println("\n⚠️  Sparse Checkout이 비활성화되어 있습니다.")
		fmt.Println("   SLIM 모드로 전환하거나 Sparse Checkout을 수동으로 활성화하세요.")
	}

	// 확장할 경로 입력 받기
	fmt.Println("\n📂 확장할 경로 입력")
	fmt.Println("   • 폴더: 'Assets/Textures/' 형식")
	fmt.Println("   • 파일: 'Assets/Models/character.fbx' 형식")
	fmt.Println("   • 여러 경로: 공백으로 구분")
	fmt.Println("   • 취소: 빈 줄 입력")
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n경로 입력: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 오류: 입력 읽기 실패: %v\n", err)
		os.Exit(1)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("⏹  작업을 취소했습니다.")
		return
	}

	// 경로 파싱
	paths := strings.Fields(input)
	if len(paths) == 0 {
		fmt.Println("❌ 오류: 유효한 경로를 입력하세요.")
		os.Exit(1)
	}

	fmt.Printf("\n✅ 추가할 경로 (%d개):\n", len(paths))
	for _, path := range paths {
		fmt.Printf("   • %s\n", path)
	}

	if !utils.ConfirmWithDefault("이 경로들을 Sparse Checkout에 추가하시겠습니까?", true) {
		fmt.Println("⏹  작업을 취소했습니다.")
		return
	}

	// Sparse Checkout이 비활성화된 경우 활성화
	if !sparseInfo["enabled"].(bool) {
		fmt.Println("\n🔧 Sparse Checkout 활성화 중...")
		cmd := exec.Command("git", "sparse-checkout", "init", "--cone")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("❌ 오류: Sparse Checkout 활성화 실패: %v\n", err)
			if len(output) > 0 {
				fmt.Printf("   상세: %s\n", string(output))
			}
			os.Exit(1)
		}
		fmt.Println("✅ Sparse Checkout 활성화 완료")
	}

	// 경로 추가
	fmt.Println("\n🔧 경로 추가 중...")
	successCount := 0
	failCount := 0
	
	for _, path := range paths {
		fmt.Printf("   • %s 추가 중...", path)
		
		// git sparse-checkout add 명령 실행
		cmd := exec.Command("git", "sparse-checkout", "add", path)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf(" ❌ 실패\n")
			if len(output) > 0 {
				fmt.Printf("     오류: %s\n", strings.TrimSpace(string(output)))
			}
			failCount++
		} else {
			fmt.Printf(" ✅\n")
			successCount++
		}
	}

	// 결과 표시
	fmt.Printf("\n📊 처리 결과:\n")
	fmt.Printf("   • 성공: %d개\n", successCount)
	if failCount > 0 {
		fmt.Printf("   • 실패: %d개\n", failCount)
	}

	if successCount > 0 {
		// Config에 경로 저장
		if err := savePathsToConfig(paths); err != nil {
			fmt.Printf("⚠️  경고: Config 저장 실패: %v\n", err)
		} else {
			fmt.Println("\n💾 추가된 경로를 Config에 저장했습니다.")
		}

		// 필요한 객체 다운로드
		fmt.Println("\n🔄 필요한 파일 다운로드 중...")
		cmd := exec.Command("git", "read-tree", "-m", "-u", "HEAD")
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("⚠️  경고: 파일 업데이트 중 오류: %v\n", err)
			if len(output) > 0 {
				fmt.Printf("   상세: %s\n", string(output))
			}
		}

		// Partial Clone 필터가 있는 경우 대용량 파일도 다운로드
		if filter := utils.GetPartialCloneFilter(); filter != "" {
			fmt.Println("\n🔄 Partial Clone 필터 우회하여 대용량 파일 다운로드 중...")
			fmt.Printf("   (현재 필터: %s)\n", filter)
			
			// 추가된 경로의 모든 blob 다운로드
			for _, path := range paths {
				if successCount > 0 {
					cmd := exec.Command("git", "ls-files", "--sparse", path)
					filesOutput, err := cmd.Output()
					if err == nil && len(filesOutput) > 0 {
						files := strings.Split(strings.TrimSpace(string(filesOutput)), "\n")
						for _, file := range files {
							if file == "" {
								continue
							}
							// 개별 파일 다운로드 시도
							cmd = exec.Command("git", "checkout", "--", file)
							cmd.Run() // 오류 무시 (이미 존재하는 파일일 수 있음)
						}
					}
				}
			}
			fmt.Println("✅ 대용량 파일 다운로드 완료")
		}

		// 최종 상태 확인
		fmt.Println("\n📋 업데이트된 Sparse Checkout 목록:")
		cmd = exec.Command("git", "sparse-checkout", "list")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for i, line := range lines {
				if i < 10 { // 처음 10개만 표시
					fmt.Printf("   • %s\n", line)
				}
			}
			if len(lines) > 10 {
				fmt.Printf("   ... 외 %d개 경로\n", len(lines)-10)
			}
		}

		// 디스크 사용량 확인
		diskUsage := utils.GetDiskUsage()
		fmt.Println("\n💾 디스크 사용량:")
		if gitSize, ok := diskUsage["git"]; ok {
			fmt.Printf("   • .git 폴더: %s\n", gitSize)
		}
		if totalSize, ok := diskUsage["total"]; ok {
			fmt.Printf("   • 프로젝트 전체: %s\n", totalSize)
		}

		fmt.Println("\n✅ 선택적 확장이 완료되었습니다!")
		fmt.Println("   추가된 경로의 파일들이 작업 트리에 나타납니다.")
	} else {
		fmt.Println("\n⚠️  경로 추가에 모두 실패했습니다.")
		fmt.Println("   경로 형식을 확인하고 다시 시도하세요.")
	}
}

// savePathsToConfig saves the expanded paths to config file
func savePathsToConfig(newPaths []string) error {
	// 현재 config의 sparse paths 가져오기
	settings := config.GetAll()
	var existingPaths []string
	
	if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
		if sparse, ok := optimize["sparse"].(map[string]interface{}); ok {
			if paths, ok := sparse["paths"].([]interface{}); ok {
				for _, path := range paths {
					if p, ok := path.(string); ok {
						existingPaths = append(existingPaths, p)
					}
				}
			}
		}
	}
	
	// 중복 제거하며 새 경로 추가
	pathMap := make(map[string]bool)
	for _, p := range existingPaths {
		pathMap[p] = true
	}
	for _, p := range newPaths {
		pathMap[p] = true
	}
	
	// 맵을 슬라이스로 변환
	var allPaths []string
	for path := range pathMap {
		allPaths = append(allPaths, path)
	}
	
	// Config에 저장
	return config.Set("optimize.sparse.paths", allPaths)
}