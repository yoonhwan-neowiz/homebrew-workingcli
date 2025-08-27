package workspace

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewExpandPathCmd creates the command for expanding a specific path in a sparse checkout
func NewExpandPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expand-path [path]",
		Short: "특정 경로를 Sparse Checkout에 추가",
		Long: `SLIM 모드에서 특정 경로를 Sparse Checkout 목록에 추가하여 해당 경로의 파일들을 다운로드합니다.
이 명령어는 저장소가 SLIM 모드(sparse-checkout 활성화 상태)일 때만 동작합니다.

파일 경로가 추가되면 자동으로 non-cone 모드로 전환되며,
config.yaml의 sparse paths도 자동으로 동기화됩니다.`,
		Example: `  ga opt workspace expand-path Assets/Art/Textures
  ga op ws expand-path src/feature/new-module
  ga opt ws expand-path README.md`,
		Args: cobra.ExactArgs(1),
		Run:  runExpandPath,
	}
	return cmd
}

func runExpandPath(cmd *cobra.Command, args []string) {
	// Pre-flight checks
	if !utils.IsSparseCheckoutEnabled() {
		fmt.Println("❌ 오류: Sparse Checkout이 활성화되어 있지 않습니다.")
		fmt.Println("   ℹ️ 이 명령어는 SLIM 모드에서만 사용할 수 있습니다. 먼저 'ga opt quick to-slim'을 실행하세요.")
		os.Exit(1)
	}

	targetPath := args[0]
	fmt.Printf("▶️ 경로 확장 시작: %s\n\n", targetPath)

	// 1. 경로 유효성 확인
	fmt.Println("1. 저장소 내 경로 유효성 검증 중...")
	if !utils.PathExistsInRepo(targetPath) {
		fmt.Printf("❌ 오류: 경로 '%s'가 저장소의 현재 HEAD에 존재하지 않습니다.\n", targetPath)
		fmt.Println("   ℹ️ Git 저장소에 존재하는 유효한 파일 또는 디렉토리 경로를 입력해주세요.")
		os.Exit(1)
	}
	fmt.Println("   ✅ 경로가 유효합니다.")

	// 2. 현재 sparse-checkout list 가져오기
	fmt.Println("\n2. 현재 Sparse Checkout 목록 확인 중...")
	currentPaths := utils.GetCurrentSparsePaths()
	fmt.Printf("   현재 %d개의 경로가 설정되어 있습니다.\n", len(currentPaths))
	
	// 이미 추가된 경로인지 확인
	for _, path := range currentPaths {
		if path == targetPath {
			fmt.Printf("\n✅ 경로 '%s'는 이미 Sparse Checkout에 포함되어 있습니다.\n", targetPath)
			return
		}
	}
	
	// 3. 파일/폴더 여부 확인 및 기존 경로에서 파일 존재 여부 확인
	isNewPathFile := !strings.HasSuffix(targetPath, "/") && strings.Contains(targetPath, ".")
	hasExistingFiles := false
	
	for _, path := range currentPaths {
		if !strings.HasSuffix(path, "/") && strings.Contains(path, ".") {
			hasExistingFiles = true
			break
		}
	}
	
	// 4. cone/non-cone 모드 결정 및 전환
	fmt.Println("\n3. Sparse Checkout 모드 결정 중...")
	needsNonCone := hasExistingFiles || isNewPathFile
	currentConeMode := utils.CheckConeMode()
	
	if needsNonCone && currentConeMode {
		fmt.Println("   📋 파일 경로가 감지되어 non-cone 모드로 전환합니다...")
		if err := utils.RunGitCommand("sparse-checkout", "init", "--no-cone"); err != nil {
			fmt.Printf("   ⚠️ Non-cone 모드 전환 실패: %v\n", err)
		} else {
			fmt.Println("   ✅ Non-cone 모드로 전환 완료")
		}
	} else if needsNonCone {
		fmt.Println("   ✅ 이미 non-cone 모드입니다 (파일 경로 지원)")
	} else if !currentConeMode {
		fmt.Println("   ✅ Non-cone 모드 유지 (기존 설정)")
	} else {
		fmt.Println("   ✅ Cone 모드 유지 (디렉토리만)")
	}
	
	// 5. 경로 추가
	fmt.Println("\n4. Sparse Checkout에 경로 추가 및 파일 다운로드 중...")
	fmt.Println("   --------------------------------------------------")

	gitCmd := exec.Command("git", "sparse-checkout", "add", targetPath)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	err := gitCmd.Run()

	fmt.Println("   --------------------------------------------------")

	if err != nil {
		fmt.Printf("\n❌ 오류: 경로 '%s'를 Sparse Checkout에 추가하는 데 실패했습니다.\n", targetPath)
		os.Exit(1)
	}

	// 6. 갱신된 sparse-checkout list 가져오기
	fmt.Println("\n5. 갱신된 Sparse Checkout 목록 확인...")
	updatedPaths := utils.GetCurrentSparsePaths()
	fmt.Printf("   ✅ 총 %d개의 경로가 활성화되었습니다.\n", len(updatedPaths))
	
	// 7. config.yaml 업데이트
	fmt.Println("\n6. 설정 파일(.gaconfig/config.yaml) 동기화 중...")
	updateConfigWithSparsePaths(updatedPaths)
	
	// 8. 결과 표시
	fmt.Printf("\n✅ 성공: 경로 '%s'를 Sparse Checkout에 성공적으로 추가했습니다.\n", targetPath)
	
	// 간단한 목록 표시 (최대 10개)
	fmt.Println("\n📋 현재 활성화된 주요 경로:")
	fmt.Println("   --------------------")
	displayCount := len(updatedPaths)
	if displayCount > 10 {
		displayCount = 10
	}
	
	for i := 0; i < displayCount; i++ {
		path := updatedPaths[i]
		if !strings.HasSuffix(path, "/") && strings.Contains(path, ".") {
			fmt.Printf("   📄 %s (파일)\n", path)
		} else {
			fmt.Printf("   📁 %s\n", path)
		}
	}
	
	if len(updatedPaths) > 10 {
		fmt.Printf("   ... 외 %d개 경로\n", len(updatedPaths)-10)
	}
	fmt.Println("   --------------------")

	fmt.Println("\n🎉 경로 확장이 완료되었습니다.")
}

// updateConfigWithSparsePaths updates config.yaml with current sparse paths
func updateConfigWithSparsePaths(paths []string) {
	// config 패키지를 통해 sparse paths 업데이트
	if err := config.Set("optimize.sparse.paths", paths); err != nil {
		fmt.Printf("   ⚠️ 설정 파일 업데이트 실패: %v\n", err)
		return
	}
	
	fmt.Printf("   ✅ 설정 파일에 %d개 경로 동기화 완료\n", len(paths))
}