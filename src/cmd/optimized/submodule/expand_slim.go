package submodule

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewExpandSlimCmd creates the Expand SLIM command for submodules
func NewExpandSlimCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-slim",
		Short: "서브모듈 선택적 경로 확장 (SLIM 유지)",
		Long: `모든 서브모듈의 Sparse Checkout 경로를 선택적으로 확장합니다.
각 서브모듈에 동일한 경로를 추가하고 Partial Clone 필터를 우회하여
대용량 파일도 다운로드합니다.

실행 내용:
1) 각 서브모듈의 현재 Sparse Checkout 상태 확인
2) 사용자가 입력한 경로를 모든 서브모듈에 추가
3) Config에 경로 저장 (optimize.submodule.sparse.paths)
4) 필요한 파일 다운로드

참고: 병렬 처리로 빠르게 실행됩니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runExpandSlim()
		},
	}
}

// runExpandSlim expands sparse checkout paths for all submodules
func runExpandSlim() {
	// 서브모듈 확인
	submoduleInfo := utils.GetSubmoduleInfo()
	count, _ := submoduleInfo["count"].(int)
	if count == 0 {
		fmt.Println("ℹ️ 서브모듈이 없습니다.")
		return
	}

	fmt.Println("\n🔸 서브모듈 SLIM 선택적 확장")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📦 총 %d개의 서브모듈이 있습니다.\n", count)

	// 현재 config의 서브모듈 sparse 경로 표시
	currentPaths := getExpandSubmoduleSparsePaths()
	if len(currentPaths) > 0 {
		fmt.Println("\n📋 Config에 저장된 서브모듈 Sparse 경로:")
		for _, path := range currentPaths {
			fmt.Printf("   • %s\n", path)
		}
	}

	// 확장할 경로 입력 받기
	fmt.Println("\n📂 모든 서브모듈에 추가할 경로 입력")
	fmt.Println("   • 폴더: 'src/core/' 형식")
	fmt.Println("   • 파일: 'src/main.cpp' 형식")
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

	fmt.Printf("\n✅ 모든 서브모듈에 추가할 경로 (%d개):\n", len(paths))
	for _, path := range paths {
		fmt.Printf("   • %s\n", path)
	}

	if !utils.ConfirmWithDefault("이 경로들을 모든 서브모듈의 Sparse Checkout에 추가하시겠습니까?", true) {
		fmt.Println("⏹  작업을 취소했습니다.")
		return
	}

	// Config에 경로 저장
	if err := saveSubmodulePathsToConfig(paths); err != nil {
		fmt.Printf("⚠️  경고: Config 저장 실패: %v\n", err)
	} else {
		fmt.Println("\n💾 추가된 경로를 Config에 저장했습니다.")
	}

	fmt.Println("\n🚀 서브모듈 Sparse Checkout 경로 확장 시작...")
	started := time.Now()

	// 결과 집계용 구조체
	type expandResult struct {
		path    string
		success bool
		message string
	}

	var (
		mu      sync.Mutex
		results []expandResult
	)

	// 서브모듈 경로 확장 작업 정의
	expandOperation := func(submodulePath string) error {
		// 서브모듈 디렉토리로 이동
		originalDir, _ := os.Getwd()
		if err := os.Chdir(submodulePath); err != nil {
			return fmt.Errorf("디렉토리 이동 실패: %v", err)
		}
		defer os.Chdir(originalDir)

		// Git 저장소인지 확인
		if !utils.IsGitRepository() {
			mu.Lock()
			results = append(results, expandResult{
				path:    submodulePath,
				success: false,
				message: "미초기화 서브모듈",
			})
			mu.Unlock()
			return nil
		}

		// Sparse Checkout 상태 확인 및 활성화
		if !utils.IsSparseCheckoutEnabled() {
			// cone 모드 여부 결정 (파일이 포함되어 있으면 non-cone)
			hasFiles := false
			for _, path := range paths {
				if !strings.HasSuffix(path, "/") && strings.Contains(path, ".") {
					hasFiles = true
					break
				}
			}

			if hasFiles {
				exec.Command("git", "sparse-checkout", "init", "--no-cone").Run()
			} else {
				exec.Command("git", "sparse-checkout", "init", "--cone").Run()
			}
		}

		// 경로 추가
		successCount := 0
		for _, path := range paths {
			cmd := exec.Command("git", "sparse-checkout", "add", path)
			if err := cmd.Run(); err == nil {
				successCount++
			}
		}

		// 파일 업데이트
		exec.Command("git", "read-tree", "-m", "-u", "HEAD").Run()

		// Partial Clone 필터가 있는 경우 대용량 파일 다운로드
		if filter := utils.GetPartialCloneFilter(); filter != "" {
			for _, path := range paths {
				cmd := exec.Command("git", "ls-files", "--sparse", path)
				filesOutput, err := cmd.Output()
				if err == nil && len(filesOutput) > 0 {
					files := strings.Split(strings.TrimSpace(string(filesOutput)), "\n")
					for _, file := range files {
						if file != "" {
							exec.Command("git", "checkout", "--", file).Run()
						}
					}
				}
			}
		}

		// 결과 기록
		mu.Lock()
		results = append(results, expandResult{
			path:    submodulePath,
			success: successCount > 0,
			message: fmt.Sprintf("%d/%d 경로 추가", successCount, len(paths)),
		})
		mu.Unlock()

		if successCount > 0 {
			fmt.Printf("✅ %s: %d개 경로 추가 완료\n", submodulePath, successCount)
		} else {
			fmt.Printf("❌ %s: 경로 추가 실패\n", submodulePath)
		}

		return nil
	}

	// 병렬 실행 (최대 4개 작업, recursive 활성화)
	_, _, err = utils.ExecuteOnSubmodulesParallel(expandOperation, 4, true)

	// 요약 출력
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("📊 작업 완료 요약")
	
	// 성공/실패 집계
	actualSuccess := 0
	actualFail := 0
	for _, r := range results {
		if r.success {
			actualSuccess++
		} else {
			actualFail++
		}
	}

	fmt.Printf("✅ 성공: %d개 서브모듈\n", actualSuccess)
	if actualFail > 0 {
		fmt.Printf("❌ 실패: %d개 서브모듈\n", actualFail)
	}
	fmt.Printf("⏱  소요 시간: %v\n", time.Since(started).Round(time.Second))

	if err != nil {
		fmt.Printf("\n⚠️ 일부 작업 실패:\n%v\n", err)
	}

	// 개별 결과 표시
	if len(results) > 0 {
		fmt.Println("\n📋 세부 결과:")
		for _, r := range results {
			if r.success {
				fmt.Printf("   ✅ %s: %s\n", r.path, r.message)
			} else {
				fmt.Printf("   ❌ %s: %s\n", r.path, r.message)
			}
		}
	}

	fmt.Println("\n✅ 서브모듈 Sparse Checkout 경로 확장이 완료되었습니다!")
	fmt.Println("   추가된 경로의 파일들이 각 서브모듈에 나타납니다.")
}

// getExpandSubmoduleSparsePaths gets sparse checkout paths from config for submodules (for expand)
func getExpandSubmoduleSparsePaths() []string {
	settings := config.GetAll()
	var paths []string
	
	// 서브모듈 전용 sparse 설정 읽기
	if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
		if submodule, ok := optimize["submodule"].(map[string]interface{}); ok {
			if sparse, ok := submodule["sparse"].(map[string]interface{}); ok {
				if configPaths, ok := sparse["paths"].([]interface{}); ok {
					for _, path := range configPaths {
						if p, ok := path.(string); ok {
							paths = append(paths, p)
						}
					}
				}
			}
		}
	}
	
	return paths
}

// saveSubmodulePathsToConfig saves the expanded paths to config file for submodules
func saveSubmodulePathsToConfig(newPaths []string) error {
	// 현재 서브모듈 sparse paths 가져오기
	existingPaths := getExpandSubmoduleSparsePaths()
	
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
	return config.Set("optimize.submodule.sparse.paths", allPaths)
}