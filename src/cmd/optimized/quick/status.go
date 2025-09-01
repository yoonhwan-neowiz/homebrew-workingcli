package quick

import (
	"fmt"
	"os"
	"strings"
	"workingcli/src/utils"
	"workingcli/src/config"
	
	"github.com/spf13/cobra"
)

// NewStatusCmd creates the status check command
func NewStatusCmd() *cobra.Command {
	var verbose bool
	
	cmd := &cobra.Command{
		Use:   "status",
		Short: "현재 최적화 상태 확인",
		Long: `현재 저장소의 최적화 상태를 한눈에 확인합니다.

표시 정보:
- 모드: SLIM (최적화) / FULL (전체)
- 저장소 크기: 현재 디스크 사용량
- Git 오브젝트: .git 폴더 내 오브젝트 수와 팩 파일 상태
- 히스토리 깊이: shallow depth 상태
- Partial Clone: blob 필터 설정
- 제외된 파일: Partial Clone으로 제외된 대용량 파일 샘플
- Sparse Checkout: 활성 경로 수
- 브랜치 필터: 숨겨진 브랜치 수

출력 예시:
╭─────────────────────────────────────────╮
│ 📊 Git 저장소 최적화 상태              │
├─────────────────────────────────────────┤
│ 모드: SLIM                             │
│ 크기: 30MB (원본: 103GB)               │
│ .git: 28MB                             │
│                                         │
│ 📦 Git 오브젝트 상태:                  │
│  • 총 오브젝트: 1,234개                │
│  • 팩 파일: 3개 (25MB)                 │
│  • Loose 오브젝트: 45개 (3MB)          │
│  • Promisor 오브젝트: 98,765개 (원격)  │
│                                         │
│ 히스토리: depth=1 (shallow)            │
│ 필터: blob:limit=1m                    │
│                                         │
│ 📁 제외된 파일 (1MB 이상):             │
│  • Quest_Main_39.prefab (103MB)        │
│  • FMODStudioCache.asset (29MB)        │
│  • MainScene.unity (24MB)              │
│  • CharacterModel.fbx (15MB)           │
│  • BackgroundTexture.psd (8MB)         │
│  ... 외 2,347개 파일                   │
│                                         │
│ Sparse: 5개 경로 활성                  │
│ 브랜치: 3/50개 표시                    │
╰─────────────────────────────────────────╯

실행되는 명령어:
- git count-objects -v  (오브젝트 수 확인)
- git rev-list --count HEAD  (커밋 수 확인)
- git config --get remote.origin.partialclonefilter  (필터 확인)
- du -sh .git  (Git 폴더 크기 확인)`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runStatus(verbose); err != nil {
				fmt.Fprintf(os.Stderr, "❌ 오류 발생: %v\n", err)
				os.Exit(1)
			}
		},
	}
	
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "상세 정보 표시")
	return cmd
}

// runStatus executes the status check logic
func runStatus(verbose bool) error {
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		return fmt.Errorf("현재 디렉토리는 Git 저장소가 아닙니다")
	}

	// 2. 기본 상태 정보 수집 (항상)
	mode := utils.GetOptimizationMode()
	partialFilter := utils.GetPartialCloneFilter()
	sparseInfo := utils.GetSparseCheckoutInfo()
	shallowInfo := utils.GetShallowInfo()
	diskUsage := utils.GetDiskUsage()
	objectInfo := utils.GetObjectInfo()
	submoduleInfo := utils.GetSubmoduleInfo()
	
	// Config 정보 가져오기
	var configInfo *config.Config
	configInfo = config.Get()
	
	// 3. Verbose 모드에서만 무거운 분석
	var excludedFiles []map[string]string
	var largestFiles []map[string]string
	var largestPack map[string]interface{}
	var dustAnalysis map[string]interface{}
	var projectTotalSize string
	
	if verbose {
		excludedFiles = utils.GetExcludedLargeFiles(partialFilter)
		largestFiles = utils.GetLargestFilesInHistory()
		largestPack = utils.GetLargestPackInfo()
		dustAnalysis = utils.GetDustAnalysis()
		projectTotalSize = utils.GetProjectTotalSize()
	}

	// 4. 결과 출력
	printStatusReport(
		verbose,
		mode,
		configInfo,
		partialFilter,
		sparseInfo,
		shallowInfo,
		diskUsage,
		objectInfo,
		submoduleInfo,
		excludedFiles,
		largestFiles,
		largestPack,
		dustAnalysis,
		projectTotalSize,
	)

	return nil
}


// printStatusReport prints the formatted status report
func printStatusReport(
	verbose bool,
	mode string,
	configInfo *config.Config,
	partialFilter string,
	sparseInfo map[string]interface{},
	shallowInfo map[string]interface{},
	diskUsage map[string]string,
	objectInfo map[string]interface{},
	submoduleInfo map[string]interface{},
	excludedFiles []map[string]string,
	largestFiles []map[string]string,
	largestPack map[string]interface{},
	dustAnalysis map[string]interface{},
	projectTotalSize string,
) {
	// Header
	fmt.Println("╭─────────────────────────────────────────╮")
	fmt.Println("│ 📊 Git 저장소 최적화 상태              │")
	fmt.Println("├─────────────────────────────────────────┤")
	
	// Mode and size
	modeDisplay := mode
	if mode == "SLIM" {
		modeDisplay = "SLIM (최적화됨)"
	} else {
		modeDisplay = "FULL (전체)"
	}
	fmt.Printf("│ 모드: %-33s │\n", modeDisplay)
	
	// Config 정보 표시
	if configInfo != nil {
		if configInfo.Optimize.Mode != "" && configInfo.Optimize.Mode != strings.ToLower(mode) {
			fmt.Printf("│ ⚠️  Config 모드: %-25s │\n", fmt.Sprintf("%s (불일치)", configInfo.Optimize.Mode))
		} else if configInfo.Optimize.Mode != "" {
			fmt.Printf("│ Config 모드: %-28s │\n", configInfo.Optimize.Mode)
		}
	}
	
	if gitSize, ok := diskUsage["git"]; ok {
		fmt.Printf("│ .git 폴더: %-28s │\n", gitSize)
	}
	
	// Verbose 모드에서만 전체 프로젝트 크기 표시
	if verbose && projectTotalSize != "" && projectTotalSize != "N/A" {
		fmt.Printf("│ 프로젝트 전체: %-23s │\n", projectTotalSize)
	}
	
	// Config 상세 설정 표시
	if configInfo != nil {
		if configInfo.Optimize.Filter.Default != "" {
			fmt.Printf("│ 필터: %-32s │\n", fmt.Sprintf("blob:limit=%s", configInfo.Optimize.Filter.Default))
		}
		if len(configInfo.Optimize.BranchScope) > 0 {
			fmt.Printf("│ 브랜치 필터: %-25s │\n", fmt.Sprintf("%d개", len(configInfo.Optimize.BranchScope)))
			for _, branch := range configInfo.Optimize.BranchScope {
				fmt.Printf("│   - %-35s │\n", branch)
			}
		}
		if len(configInfo.Optimize.Sparse.Paths) > 0 {
			fmt.Printf("│ Sparse 경로: %-26s │\n", fmt.Sprintf("%d개", len(configInfo.Optimize.Sparse.Paths)))
			for _, path := range configInfo.Optimize.Sparse.Paths {
				fmt.Printf("│   - %-35s │\n", utils.TruncateString(path, 35))
			}
		}
	}
	
	fmt.Println("│                                         │")
	
	// Git objects
	if len(objectInfo) > 0 {
		fmt.Println("│ 📦 Git 오브젝트 상태:                  │")
		
		totalObjects := 0
		if loose, ok := objectInfo["looseObjects"].(int); ok {
			totalObjects += loose
		}
		if pack, ok := objectInfo["packObjects"].(int); ok {
			totalObjects += pack
		}
		
		if totalObjects > 0 {
			fmt.Printf("│  • 총 오브젝트: %-22s │\n", fmt.Sprintf("%,d개", totalObjects))
		}
		
		if packCount, ok := objectInfo["packCount"].(int); ok {
			if packSize, ok := objectInfo["packSize"].(string); ok {
				fmt.Printf("│  • 팩 파일: %d개 (%-18s) │\n", packCount, packSize)
			}
		}
		
		if looseCount, ok := objectInfo["looseObjects"].(int); ok {
			if looseSize, ok := objectInfo["looseSize"].(string); ok {
				fmt.Printf("│  • Loose 오브젝트: %d개 (%-11s) │\n", looseCount, looseSize)
			}
		}
		
		if hasPromisor, ok := objectInfo["hasPromisor"].(bool); ok && hasPromisor {
			fmt.Println("│  • Promisor 오브젝트: 활성 (원격)      │")
		}
		
		fmt.Println("│                                         │")
	}
	
	// History status
	if isShallow, ok := shallowInfo["isShallow"].(bool); ok && isShallow {
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("│ 히스토리: depth=%d (shallow)          │\n", depth)
		}
	} else {
		fmt.Println("│ 히스토리: 전체                         │")
	}
	
	// Partial Clone filter
	if partialFilter != "" {
		fmt.Printf("│ 필터: %-33s │\n", partialFilter)
	}
	
	// Excluded files
	if len(excludedFiles) > 0 {
		fmt.Println("│                                         │")
		fmt.Println("│ 📁 제외된 대용량 파일:                 │")
		for _, file := range excludedFiles {
			fmt.Printf("│  • %-20s %12s │\n", 
				utils.TruncateString(file["path"], 20), 
				file["size"])
		}
	}
	
	// Sparse checkout
	if sparseEnabled, ok := sparseInfo["enabled"].(bool); ok && sparseEnabled {
		if count, ok := sparseInfo["count"].(int); ok && count > 0 {
			fmt.Println("│                                         │")
			fmt.Printf("│ Sparse Checkout: %d개 경로 활성       │\n", count)
		}
	}
	
	// Submodules
	if count, ok := submoduleInfo["count"].(int); ok && count > 0 {
		if optimized, ok := submoduleInfo["optimized"].(int); ok {
			fmt.Printf("│ 서브모듈: %d개 (최적화: %d개)         │\n", count, optimized)
		}
	}
	
	// Largest files in history
	if len(largestFiles) > 0 {
		fmt.Println("│                                         │")
		fmt.Println("│ 🔍 히스토리 최대 파일 (Top 5):         │")
		for _, file := range largestFiles {
			fmt.Printf("│  • %-25s %10s │\n",
				utils.TruncateString(file["path"], 25),
				file["size"])
			if fullPath := file["path"]; len(fullPath) > 25 {
				fmt.Printf("│    → %s %s\n│                                         │\n",
					utils.TruncateString(fullPath, 30),
					file["status"])
			}
		}
	}
	
	// Largest pack information
	if len(largestPack) > 0 {
		if name, ok := largestPack["name"].(string); ok {
			fmt.Println("│                                         │")
			fmt.Println("│ 📦 최대 Pack 파일:                      │")
			fmt.Printf("│  • 이름: %-30s │\n", utils.TruncateString(name, 30))
			if size, ok := largestPack["size"].(string); ok {
				fmt.Printf("│  • 크기: %-30s │\n", size)
			}
			if objects, ok := largestPack["objects"].(int); ok {
				fmt.Printf("│  • 오브젝트: %,d개                    │\n", objects)
			}
		}
	}
	
	// Dust analysis (verbose에서만)
	if verbose {
		if available, ok := dustAnalysis["available"].(bool); ok && available {
			if topDirs, ok := dustAnalysis["topDirs"].([]map[string]string); ok && len(topDirs) > 0 {
				fmt.Println("│                                         │")
				fmt.Println("│ 💾 Dust 디스크 분석 (Top 5):           │")
				for _, dir := range topDirs {
					fmt.Printf("│  • %-25s %10s │\n",
						utils.TruncateString(dir["path"], 25),
						dir["size"])
				}
			}
		}
	}
	
	fmt.Println("╰─────────────────────────────────────────╯")
	
	// 기본 모드에서는 힌트 표시
	if !verbose {
		fmt.Println("\n💡 상세 정보: ga opt quick status -v")
	}
}

