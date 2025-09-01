package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// IsGitRepository checks if current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetOptimizationMode determines if repo is in SLIM or FULL mode
func GetOptimizationMode() string {
	partialFilter := GetPartialCloneFilter()
	sparseEnabled := IsSparseCheckoutEnabled()
	isShallow := IsShallowRepository()

	if partialFilter != "" || sparseEnabled || isShallow {
		return "SLIM"
	}
	return "FULL"
}

// GetPartialCloneFilter returns the current partial clone filter
func GetPartialCloneFilter() string {
	cmd := exec.Command("git", "config", "remote.origin.partialclonefilter")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// IsSparseCheckoutEnabled checks if sparse checkout is enabled
func IsSparseCheckoutEnabled() bool {
	cmd := exec.Command("git", "config", "core.sparseCheckout")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == "true"
}

// GetSparseCheckoutInfo returns sparse checkout information
func GetSparseCheckoutInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["enabled"] = IsSparseCheckoutEnabled()
	
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

// IsShallowRepository checks if repository is shallow
func IsShallowRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == "true"
}

// GetShallowInfo returns shallow repository information
func GetShallowInfo() map[string]interface{} {
	info := make(map[string]interface{})
	info["isShallow"] = IsShallowRepository()
	
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

// GetDiskUsage returns disk usage information (optimized for speed)
// Uses dust if available, otherwise falls back to du
func GetDiskUsage() map[string]string {
	usage := make(map[string]string)
	
	// Try dust first (faster and better output)
	if _, err := exec.LookPath("dust"); err == nil {
		// Use dust for .git folder
		cmd := exec.Command("dust", "-n", "1", "-d", "0", ".git")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			if len(lines) > 0 {
				fields := strings.Fields(lines[0])
				if len(fields) >= 1 {
					usage["git"] = fields[0]
				}
			}
		}
	} else {
		// Fallback to du if dust is not available
		cmd := exec.Command("du", "-sh", ".git")
		if output, err := cmd.Output(); err == nil {
			fields := strings.Fields(string(output))
			if len(fields) > 0 {
				usage["git"] = fields[0]
			}
		}
	}
	
	// 전체 프로젝트 크기는 verbose 모드에서만 (따로 호출)
	// 기본 모드에서는 생략하여 성능 향상
	
	return usage
}

// GetGitDirSize returns the size of a .git directory (bytes and human string)
// This function handles both regular .git directories and gitfile references
func GetGitDirSize(basePath string) (int64, string) {
	gitPath := filepath.Join(basePath, ".git")
	
	// .git이 gitfile인 경우 실제 경로 읽기
	if info, err := os.Stat(gitPath); err == nil && !info.IsDir() {
		if content, err := os.ReadFile(gitPath); err == nil {
			gitDir := strings.TrimPrefix(string(content), "gitdir: ")
			gitDir = strings.TrimSpace(gitDir)
			if !filepath.IsAbs(gitDir) {
				// 상대 경로는 basePath를 기준으로 해석
				gitDir = filepath.Join(basePath, gitDir)
			}
			gitPath = gitDir
		}
	}
	
	// du -sh로 사람이 읽기 쉬운 크기 확보
	human := "N/A"
	if output, err := exec.Command("du", "-sh", gitPath).Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			human = fields[0]
		}
	}
	
	// 사람이 읽는 단위를 바이트로 근사 변환
	bytes := ParseHumanSize(human)
	
	return bytes, human
}

// GetSubmoduleGitSize returns the size of a submodule's .git directory
// It changes to the submodule directory, gets the size, and returns to original directory
func GetSubmoduleGitSize(submodulePath string) (int64, string) {
	// 현재 디렉토리 저장
	originalDir, err := os.Getwd()
	if err != nil {
		return 0, "N/A"
	}
	
	// 서브모듈 디렉토리로 이동
	if err := os.Chdir(submodulePath); err != nil {
		return 0, "N/A"
	}
	defer os.Chdir(originalDir) // 원래 디렉토리로 복귀
	
	// 현재 디렉토리(".")를 기준으로 크기 측정
	return GetGitDirSize(".")
}

// ParseHumanSize converts human-readable size string to bytes
func ParseHumanSize(sizeStr string) int64 {
	if sizeStr == "N/A" || sizeStr == "" {
		return 0
	}
	
	var multiplier int64 = 1
	
	// 단위 확인 및 제거
	if strings.HasSuffix(sizeStr, "K") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "K")
	} else if strings.HasSuffix(sizeStr, "M") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "M")
	} else if strings.HasSuffix(sizeStr, "G") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "G")
	} else if strings.HasSuffix(sizeStr, "T") {
		multiplier = 1024 * 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "T")
	}
	
	// 숫자 파싱
	if val, err := strconv.ParseFloat(sizeStr, 64); err == nil {
		return int64(val * float64(multiplier))
	}
	
	return 0
}

// GetObjectInfo returns git object statistics
func GetObjectInfo() map[string]interface{} {
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

// GetSubmoduleInfo returns submodule information
func GetSubmoduleInfo() map[string]interface{} {
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

// GetExcludedLargeFiles returns sample of large files excluded by filter
func GetExcludedLargeFiles(filter string) []map[string]string {
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
	
	// Find large files in HEAD only (much faster than full history)
	cmd := exec.Command("git", "ls-tree", "-r", "-l", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return excludedFiles
	}
	
	// Parse ls-tree output to find large files
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// ls-tree format: <mode> <type> <object> <size>\t<path>
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			// Parse size (4th field)
			if size, err := strconv.ParseInt(fields[3], 10, 64); err == nil {
				if size >= sizeLimit {
					// Path is after the tab
					if tabIndex := strings.Index(line, "\t"); tabIndex > 0 {
						path := line[tabIndex+1:]
						excludedFiles = append(excludedFiles, map[string]string{
							"path": filepath.Base(path),
							"size": FormatSize(size),
						})
					}
				}
			}
		}
		
		// Limit to 5 sample files
		if len(excludedFiles) >= 5 {
			break
		}
	}
	
	return excludedFiles
}

// GetLargestFilesInHistory returns the 5 largest files in current HEAD (not full history)
func GetLargestFilesInHistory() []map[string]string {
	var largestFiles []map[string]string
	
	// Get files in HEAD only (much faster than full history)
	cmd := exec.Command("git", "ls-tree", "-r", "-l", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return largestFiles
	}
	
	// Parse ls-tree output
	type fileInfo struct {
		path string
		size int64
	}
	var files []fileInfo
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// ls-tree format: <mode> <type> <object> <size>\t<path>
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			// Parse size (4th field)
			if size, err := strconv.ParseInt(fields[3], 10, 64); err == nil {
				// Path is after the tab
				if tabIndex := strings.Index(line, "\t"); tabIndex > 0 {
					path := line[tabIndex+1:]
					files = append(files, fileInfo{
						path: path,
						size: size,
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
		largestFiles = append(largestFiles, map[string]string{
			"path": files[i].path,
			"size": FormatSize(files[i].size),
		})
	}
	
	return largestFiles
}

// GetLargestPackInfo returns information about the largest pack file
func GetLargestPackInfo() map[string]interface{} {
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
		packInfo["size"] = FormatSize(largestSize)
		
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

// GetDustAnalysis runs dust command if available and returns disk usage analysis
// This is for verbose mode detailed analysis
func GetDustAnalysis() map[string]interface{} {
	analysis := make(map[string]interface{})
	
	// Check if dust is installed
	if _, err := exec.LookPath("dust"); err != nil {
		analysis["available"] = false
		return analysis
	}
	
	analysis["available"] = true
	
	// Run dust with depth limit and reverse order
	cmd := exec.Command("dust", "-d", "2", "-r", "-n", "5")
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

// GetProjectTotalSize gets the total project size for verbose mode
// Uses dust if available for better performance, otherwise falls back to du
func GetProjectTotalSize() string {
	// Try dust first
	if _, err := exec.LookPath("dust"); err == nil {
		cmd := exec.Command("dust", "-n", "1", "-d", "0", ".")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			if len(lines) > 0 {
				fields := strings.Fields(lines[0])
				if len(fields) >= 1 {
					return fields[0]
				}
			}
		}
	}
	
	// Fallback to du
	cmd := exec.Command("du", "-sh", ".")
	if output, err := cmd.Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			return fields[0]
		}
	}
	
	return "N/A"
}

// FormatSize formats bytes to human readable format
func FormatSize(bytes int64) string {
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

// TruncateString truncates string to specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		// fallback: symbolic-ref 사용
		cmd = exec.Command("git", "symbolic-ref", "--short", "HEAD")
		output, err = cmd.Output()
		if err != nil {
			return "HEAD"
		}
	}
	return strings.TrimSpace(string(output))
}

// BranchExists checks if a branch exists (local or remote)
func BranchExists(branch string) bool {
	// 로컬 브랜치 확인
	cmd := exec.Command("git", "rev-parse", "--verify", branch)
	err := cmd.Run()
	if err == nil {
		return true
	}
	
	// 원격 브랜치 확인
	cmd = exec.Command("git", "rev-parse", "--verify", "origin/"+branch)
	err = cmd.Run()
	return err == nil
}

// HasUncommittedChanges checks if there are uncommitted changes
func HasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, _ := cmd.Output()
	return len(strings.TrimSpace(string(output))) > 0
}

// FindMergeBase attempts to find merge base between two branches
func FindMergeBase(branch1, branch2 string) (string, error) {
	cmd := exec.Command("git", "merge-base", branch1, branch2)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	mergeBase := strings.TrimSpace(string(output))
	if mergeBase == "" {
		return "", fmt.Errorf("머지베이스를 찾을 수 없습니다")
	}
	
	return mergeBase, nil
}

// GetShortCommit returns short version of commit hash
func GetShortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

// GetBranchDistance returns the number of commits between base and branch
func GetBranchDistance(branch, base string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", base+".."+branch)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// GetBranches returns list of branches (local and remote counts)
func GetBranches() ([]string, int) {
	localBranches := []string{}
	
	// 로컬 브랜치
	cmd := exec.Command("git", "branch")
	output, _ := cmd.Output()
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			localBranches = append(localBranches, branch)
		}
	}
	
	// 원격 브랜치 개수
	cmd = exec.Command("git", "branch", "-r")
	output, _ = cmd.Output()
	remoteCount := 0
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		remoteCount = len(lines)
	}
	
	return localBranches, remoteCount
}

// GetCurrentSparsePaths returns the current list of sparse checkout paths
func GetCurrentSparsePaths() []string {
	cmd := exec.Command("git", "sparse-checkout", "list")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var paths []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

// CheckConeMode checks if sparse-checkout is in cone mode
func CheckConeMode() bool {
	cmd := exec.Command("git", "config", "core.sparseCheckoutCone")
	output, err := cmd.Output()
	if err != nil {
		// If config doesn't exist, default is cone mode in newer Git versions
		return true
	}
	return strings.TrimSpace(string(output)) != "false"
}

// RunGitCommand runs a git command and returns the output
func RunGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

// PathExistsInRepo checks if a path exists in the git repository
func PathExistsInRepo(path string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", path)
	err := cmd.Run()
	return err == nil
}

// GetLocalBranches returns list of local branches
func GetLocalBranches() []string {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	var branches []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches
}

// GetRemoteBranches returns list of remote branches
func GetRemoteBranches() []string {
	cmd := exec.Command("git", "branch", "-r")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	var branches []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if branch != "" && !strings.Contains(branch, "HEAD") {
			branches = append(branches, branch)
		}
	}
	return branches
}



// CountLocalBranches returns the number of local branches
func CountLocalBranches() int {
	branches := GetLocalBranches()
	return len(branches)
}

// CountRemoteBranches returns the number of remote branches
func CountRemoteBranches() int {
	branches := GetRemoteBranches()
	return len(branches)
}

// GetAllUniqueBranches returns all unique branch names (local and remote without origin/ prefix)
func GetAllUniqueBranches() []string {
	branchMap := make(map[string]bool)
	
	// 로컬 브랜치 추가
	localBranches := GetLocalBranches()
	for _, branch := range localBranches {
		branchMap[branch] = true
	}
	
	// 원격 브랜치 추가 (origin/ 제거)
	remoteBranches := GetRemoteBranches()
	for _, branch := range remoteBranches {
		// origin/ 접두사 제거
		cleanBranch := strings.TrimPrefix(branch, "origin/")
		branchMap[cleanBranch] = true
	}
	
	// 맵을 슬라이스로 변환
	var uniqueBranches []string
	for branch := range branchMap {
		uniqueBranches = append(uniqueBranches, branch)
	}
	
	// 정렬
	for i := 0; i < len(uniqueBranches); i++ {
		for j := i + 1; j < len(uniqueBranches); j++ {
			if uniqueBranches[i] > uniqueBranches[j] {
				uniqueBranches[i], uniqueBranches[j] = uniqueBranches[j], uniqueBranches[i]
			}
		}
	}
	
	return uniqueBranches
}

// Contains checks if a slice contains a specific item
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IsFilePath checks if a path represents a file (not a directory)
func IsFilePath(path string) bool {
	return !strings.HasSuffix(path, "/") && strings.Contains(path, ".")
}

// HasFilePaths checks if any path in the list is a file path
func HasFilePaths(paths []string) bool {
	for _, path := range paths {
		if IsFilePath(path) {
			return true
		}
	}
	return false
}

// InitSparseCheckoutWithMode initializes sparse checkout with appropriate mode
func InitSparseCheckoutWithMode(paths []string) error {
	// Check if any path is a file
	hasFiles := HasFilePaths(paths)
	
	var cmd *exec.Cmd
	if hasFiles {
		// Use non-cone mode for file paths
		cmd = exec.Command("git", "sparse-checkout", "init", "--no-cone")
	} else {
		// Use cone mode for directories only
		cmd = exec.Command("git", "sparse-checkout", "init", "--cone")
	}
	
	return cmd.Run()
}

// EnsureNonConeMode ensures sparse checkout is in non-cone mode if needed
func EnsureNonConeMode(newPaths []string, existingPaths []string) error {
	// Check if we need non-cone mode
	needsNonCone := HasFilePaths(newPaths) || HasFilePaths(existingPaths)
	currentConeMode := CheckConeMode()
	
	if needsNonCone && currentConeMode {
		// Switch to non-cone mode
		return RunGitCommand("sparse-checkout", "init", "--no-cone")
	}
	
	return nil
}

// SubmoduleOperation은 각 서브모듈에서 실행할 작업을 정의
type SubmoduleOperation func(path string) error

// ExecuteOnSubmodulesParallel은 모든 서브모듈에 대해 병렬로 작업을 실행
func ExecuteOnSubmodulesParallel(operation SubmoduleOperation, jobs int, recursive bool) (int, int, error) {
	if jobs < 1 {
		jobs = 1
	}

	// 서브모듈 목록 가져오기
	args := []string{"submodule", "foreach"}
	if recursive {
		args = append(args, "--recursive")
	}
	args = append(args, "--quiet", "echo $path")
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("서브모듈 목록을 가져올 수 없습니다: %v", err)
	}

	submodulePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(submodulePaths) == 0 || (len(submodulePaths) == 1 && submodulePaths[0] == "") {
		return 0, 0, nil // 서브모듈이 없음
	}

	// 작업 풀 생성
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, jobs)
	errChan := make(chan error, len(submodulePaths))
	successChan := make(chan string, len(submodulePaths))

	for _, submodulePath := range submodulePaths {
		if submodulePath == "" {
			continue
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 작업 슬롯 획득
			defer func() { <-semaphore }() // 작업 슬롯 반환

			// 작업 실행
			if err := operation(path); err != nil {
				errChan <- fmt.Errorf("%s: %v", path, err)
			} else {
				successChan <- path
			}
		}(submodulePath)
	}

	// 모든 작업 완료 대기
	wg.Wait()
	close(errChan)
	close(successChan)

	// 결과 집계
	successCount := len(successChan)
	failCount := len(errChan)

	// 에러 수집
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return successCount, failCount, fmt.Errorf("일부 서브모듈 처리 실패:\n%s", strings.Join(errors, "\n"))
	}

	return successCount, failCount, nil
}