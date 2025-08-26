package quick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"path/filepath"
	"bytes"
	
	"github.com/spf13/cobra"
)

// NewStatusCmd creates the status check command
func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
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
			if err := runStatus(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ 오류 발생: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// runStatus executes the status check logic
func runStatus() error {
	// 1. Git 저장소 확인
	if !isGitRepository() {
		return fmt.Errorf("현재 디렉토리는 Git 저장소가 아닙니다")
	}

	// 2. 각종 상태 정보 수집
	mode := getOptimizationMode()
	partialFilter := getPartialCloneFilter()
	sparseInfo := getSparseCheckoutInfo()
	shallowInfo := getShallowInfo()
	diskUsage := getDiskUsage()
	objectInfo := getObjectInfo()
	submoduleInfo := getSubmoduleInfo()
	excludedFiles := getExcludedLargeFiles(partialFilter)
	largestFiles := getLargestFilesInHistory()
	largestPack := getLargestPackInfo()
	dustAnalysis := getDustAnalysis()

	// 3. 결과 출력
	printStatusReport(
		mode,
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
	)

	return nil
}

// isGitRepository checks if current directory is a git repository
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// getOptimizationMode determines if repo is in SLIM or FULL mode
func getOptimizationMode() string {
	partialFilter := getPartialCloneFilter()
	sparseEnabled := isSparseCheckoutEnabled()
	isShallow := isShallowRepository()

	if partialFilter != "" || sparseEnabled || isShallow {
		return "SLIM"
	}
	return "FULL"
}

// getPartialCloneFilter returns the current partial clone filter
func getPartialCloneFilter() string {
	cmd := exec.Command("git", "config", "remote.origin.partialclonefilter")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// isSparseCheckoutEnabled checks if sparse checkout is enabled
func isSparseCheckoutEnabled() bool {
	cmd := exec.Command("git", "config", "core.sparseCheckout")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == "true"
}

// getSparseCheckoutInfo returns sparse checkout information
func getSparseCheckoutInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["enabled"] = isSparseCheckoutEnabled()
	
	if info["enabled"].(bool) {
		cmd := exec.Command("git", "sparse-checkout", "list")
		output, err := cmd.Output()
		if err == nil {
			paths := strings.Split(strings.TrimSpace(string(output)), "\n")
			info["paths"] = paths
			info["count"] = len(paths)
		} else {
			info["count"] = 0
		}
	} else {
		info["count"] = 0
	}
	
	return info
}

// isShallowRepository checks if repository is shallow
func isShallowRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == "true"
}

// getShallowInfo returns shallow repository information
func getShallowInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["isShallow"] = isShallowRepository()
	
	if info["isShallow"].(bool) {
		cmd := exec.Command("git", "rev-list", "--count", "HEAD")
		output, err := cmd.Output()
		if err == nil {
			count, _ := strconv.Atoi(strings.TrimSpace(string(output)))
			info["depth"] = count
		}
	}
	
	return info
}

// getDiskUsage returns disk usage information
func getDiskUsage() map[string]string {
	usage := make(map[string]string)
	
	// .git 폴더 크기
	cmd := exec.Command("du", "-sh", ".git")
	if output, err := cmd.Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			usage["git"] = fields[0]
		}
	}
	
	// 전체 프로젝트 크기
	cmd = exec.Command("du", "-sh", ".")
	if output, err := cmd.Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			usage["total"] = fields[0]
		}
	}
	
	return usage
}

// getObjectInfo returns git object statistics
func getObjectInfo() map[string]interface{} {
	info := make(map[string]interface{})
	
	cmd := exec.Command("git", "count-objects", "-v")
	output, err := cmd.Output()
	if err != nil {
		return info
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			key := parts[0]
			value := strings.TrimSpace(parts[1])
			
			switch key {
			case "count":
				info["looseObjects"], _ = strconv.Atoi(value)
			case "size":
				info["looseSize"] = value + " KB"
			case "in-pack":
				info["packObjects"], _ = strconv.Atoi(value)
			case "packs":
				info["packCount"], _ = strconv.Atoi(value)
			case "size-pack":
				sizeKB, _ := strconv.Atoi(value)
				sizeMB := sizeKB / 1024
				info["packSize"] = fmt.Sprintf("%d MB", sizeMB)
			case "prune-packable":
				info["prunePackable"], _ = strconv.Atoi(value)
			case "garbage":
				info["garbage"], _ = strconv.Atoi(value)
			}
		}
	}
	
	// Check for promisor objects
	cmd = exec.Command("git", "config", "remote.origin.promisor")
	if output, err := cmd.Output(); err == nil && strings.TrimSpace(string(output)) == "true" {
		// Estimate promisor objects (not directly countable)
		info["hasPromisor"] = true
	}
	
	return info
}

// getSubmoduleInfo returns submodule information
func getSubmoduleInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["count"] = 0
	info["optimized"] = 0
	
	// Check if .gitmodules exists
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return info
	}
	
	// Get submodule list
	cmd := exec.Command("git", "submodule", "status")
	output, err := cmd.Output()
	if err != nil {
		return info
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 && lines[0] != "" {
		info["count"] = len(lines)
	}
	
	// Count optimized submodules (simplified check)
	var optimizedCount int
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Extract submodule path
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			submodulePath := fields[1]
			// Check if submodule is shallow
			cmd := exec.Command("git", "-C", submodulePath, "rev-parse", "--is-shallow-repository")
			if output, err := cmd.Output(); err == nil {
				if strings.TrimSpace(string(output)) == "true" {
					optimizedCount++
				}
			}
		}
	}
	info["optimized"] = optimizedCount
	
	return info
}

// getExcludedLargeFiles returns sample of large files excluded by filter
func getExcludedLargeFiles(filter string) []map[string]string {
	var excludedFiles []map[string]string
	
	if filter == "" {
		return excludedFiles
	}
	
	// Parse filter to get size limit
	var sizeLimit int64 = 1048576 // Default 1MB
	if strings.HasPrefix(filter, "blob:limit=") {
		limitStr := strings.TrimPrefix(filter, "blob:limit=")
		// Parse size (e.g., "1m", "1M", "1048576")
		if strings.HasSuffix(strings.ToLower(limitStr), "m") {
			limitStr = strings.TrimSuffix(strings.ToLower(limitStr), "m")
			if val, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
				sizeLimit = val * 1048576
			}
		} else if strings.HasSuffix(strings.ToLower(limitStr), "k") {
			limitStr = strings.TrimSuffix(strings.ToLower(limitStr), "k")
			if val, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
				sizeLimit = val * 1024
			}
		} else {
			if val, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
				sizeLimit = val
			}
		}
	}
	
	// Find large files in git history
	cmd := exec.Command("git", "rev-list", "--objects", "--all")
	revListOutput, err := cmd.Output()
	if err != nil {
		return excludedFiles
	}
	
	// Get file sizes using git cat-file
	var largeFiles []map[string]string
	scanner := bytes.NewBuffer(revListOutput)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			objectID := fields[0]
			filePath := strings.Join(fields[1:], " ")
			
			// Get object size
			cmd = exec.Command("git", "cat-file", "-s", objectID)
			if sizeOutput, err := cmd.Output(); err == nil {
				if size, err := strconv.ParseInt(strings.TrimSpace(string(sizeOutput)), 10, 64); err == nil {
					if size >= sizeLimit {
						largeFiles = append(largeFiles, map[string]string{
							"path": filepath.Base(filePath),
							"size": formatSize(size),
						})
					}
				}
			}
		}
		
		// Limit to 5 sample files
		if len(largeFiles) >= 5 {
			break
		}
	}
	
	return largeFiles
}

// getLargestFilesInHistory returns the 5 largest files in git history with their paths
func getLargestFilesInHistory() []map[string]string {
	var largestFiles []map[string]string
	
	// Get all objects in history
	cmd := exec.Command("git", "rev-list", "--objects", "--all")
	revListOutput, err := cmd.Output()
	if err != nil {
		return largestFiles
	}
	
	// Collect all file sizes
	type fileInfo struct {
		path string
		size int64
		objectID string
	}
	var files []fileInfo
	
	scanner := bytes.NewBuffer(revListOutput)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			objectID := fields[0]
			filePath := strings.Join(fields[1:], " ")
			
			// Skip if path is empty
			if filePath == "" {
				continue
			}
			
			// Get object size
			cmd = exec.Command("git", "cat-file", "-s", objectID)
			if sizeOutput, err := cmd.Output(); err == nil {
				if size, err := strconv.ParseInt(strings.TrimSpace(string(sizeOutput)), 10, 64); err == nil {
					files = append(files, fileInfo{
						path:     filePath,
						size:     size,
						objectID: objectID,
					})
				}
			}
		}
	}
	
	// Sort by size descending
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if files[j].size > files[i].size {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
	
	// Take top 5
	for i := 0; i < 5 && i < len(files); i++ {
		// Check if file exists in working tree
		existsInWorking := "(deleted)"
		if _, err := os.Stat(files[i].path); err == nil {
			existsInWorking = "(exists)"
		}
		
		largestFiles = append(largestFiles, map[string]string{
			"path":     files[i].path,
			"size":     formatSize(files[i].size),
			"objectID": files[i].objectID[:7], // Short hash
			"status":   existsInWorking,
		})
	}
	
	return largestFiles
}

// getLargestPackInfo returns information about the largest pack file
func getLargestPackInfo() map[string]interface{} {
	packInfo := make(map[string]interface{})
	
	// Find all pack files
	packDir := ".git/objects/pack"
	files, err := os.ReadDir(packDir)
	if err != nil {
		return packInfo
	}
	
	var largestPack os.DirEntry
	var largestSize int64
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pack") {
			info, err := file.Info()
			if err == nil {
				if info.Size() > largestSize {
					largestSize = info.Size()
					largestPack = file
				}
			}
		}
	}
	
	if largestPack != nil {
		packInfo["name"] = largestPack.Name()
		packInfo["size"] = formatSize(largestSize)
		
		// Get pack file statistics
		packPath := filepath.Join(packDir, largestPack.Name())
		cmd := exec.Command("git", "verify-pack", "-v", packPath)
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			objectCount := 0
			for _, line := range lines {
				if strings.Contains(line, " blob ") || strings.Contains(line, " tree ") ||
				   strings.Contains(line, " commit ") || strings.Contains(line, " tag ") {
					objectCount++
				}
			}
			packInfo["objects"] = objectCount
		}
	}
	
	return packInfo
}

// getDustAnalysis runs dust command if available and returns disk usage analysis
func getDustAnalysis() map[string]interface{} {
	analysis := make(map[string]interface{})
	
	// Check if dust is installed
	cmd := exec.Command("which", "dust")
	if err := cmd.Run(); err != nil {
		analysis["available"] = false
		return analysis
	}
	
	analysis["available"] = true
	
	// Run dust with depth limit and reverse order
	cmd = exec.Command("dust", "-d", "2", "-r", "-n", "5")
	output, err := cmd.Output()
	if err != nil {
		return analysis
	}
	
	// Parse dust output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var topDirs []map[string]string
	
	for i, line := range lines {
		if i >= 5 { // Limit to top 5
			break
		}
		
		// Parse dust output format
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			size := fields[0]
			path := strings.Join(fields[1:], " ")
			topDirs = append(topDirs, map[string]string{
				"path": path,
				"size": size,
			})
		}
	}
	
	analysis["topDirs"] = topDirs
	return analysis
}

// formatSize formats bytes to human readable format
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1fGB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1fMB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1fKB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// printStatusReport prints the formatted status report
func printStatusReport(
	mode string,
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
	
	if gitSize, ok := diskUsage["git"]; ok {
		fmt.Printf("│ .git 폴더: %-28s │\n", gitSize)
	}
	if totalSize, ok := diskUsage["total"]; ok {
		fmt.Printf("│ 프로젝트 전체: %-24s │\n", totalSize)
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
				truncateString(file["path"], 20), 
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
				truncateString(file["path"], 25),
				file["size"])
			if fullPath := file["path"]; len(fullPath) > 25 {
				fmt.Printf("│    → %s %s\n│                                         │\n",
					truncateString(fullPath, 30),
					file["status"])
			}
		}
	}
	
	// Largest pack information
	if len(largestPack) > 0 {
		if name, ok := largestPack["name"].(string); ok {
			fmt.Println("│                                         │")
			fmt.Println("│ 📦 최대 Pack 파일:                      │")
			fmt.Printf("│  • 이름: %-30s │\n", truncateString(name, 30))
			if size, ok := largestPack["size"].(string); ok {
				fmt.Printf("│  • 크기: %-30s │\n", size)
			}
			if objects, ok := largestPack["objects"].(int); ok {
				fmt.Printf("│  • 오브젝트: %,d개                    │\n", objects)
			}
		}
	}
	
	// Dust analysis
	if available, ok := dustAnalysis["available"].(bool); ok && available {
		if topDirs, ok := dustAnalysis["topDirs"].([]map[string]string); ok && len(topDirs) > 0 {
			fmt.Println("│                                         │")
			fmt.Println("│ 💾 Dust 디스크 분석 (Top 5):           │")
			for _, dir := range topDirs {
				fmt.Printf("│  • %-25s %10s │\n",
					truncateString(dir["path"], 25),
					dir["size"])
			}
		}
	} else {
		fmt.Println("│                                         │")
		fmt.Println("│ ℹ️  dust 명령어가 설치되어 있지 않음    │")
	}
	
	fmt.Println("╰─────────────────────────────────────────╯")
}

// truncateString truncates string to specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}