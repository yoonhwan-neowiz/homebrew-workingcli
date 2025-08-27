package submodule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"workingcli/src/utils"
	
	"github.com/spf13/cobra"
)

// NewStatusCmd creates the submodule status command
func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "서브모듈별 최적화 상태 확인",
		Long: `각 서브모듈의 최적화 상태를 개별적으로 확인합니다.

표시 정보:
- 모드: SLIM (최적화) / FULL (전체)
- 크기: 각 서브모듈의 .git 폴더 크기
- Shallow 상태: depth 정보
- Partial Clone: blob 필터 설정
- Sparse Checkout: 활성 경로 수
- 커밋 수: 각 서브모듈의 히스토리 깊이

출력 예시:
╭─────────────────────────────────────────╮
│ 🔍 서브모듈 최적화 상태                │
├─────────────────────────────────────────┤
│ 📦 Packages/DesignSystem                │
│ • 모드: SLIM                           │
│ • 크기: 15MB                           │
│ • Shallow: depth=1                     │
│ • 필터: blob:limit=500k                │
│ • Sparse: 3개 경로                     │
├─────────────────────────────────────────┤
│ 📦 Packages/NetworkingLayer             │
│ • 모드: FULL                           │
│ • 크기: 230MB                          │
│ • 전체 히스토리 (커밋 1,234개)         │
├─────────────────────────────────────────┤
│ 📊 요약                                 │
│ • 전체 서브모듈: 2개                   │
│ • 최적화됨: 1개                        │
│ • 전체 크기: 245MB                     │
│ • 절약 가능: ~200MB                    │
╰─────────────────────────────────────────╯`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runSubmoduleStatus(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ 오류 발생: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// SubmoduleStatus holds status information for a single submodule
type SubmoduleStatus struct {
	Name          string
	Path          string
	Mode          string
	Size          string
	IsShallow     bool
	Depth         int
	PartialFilter string
	SparseEnabled bool
	SparsePaths   int
	CommitCount   int
	SizeBytes     int64
}

// runSubmoduleStatus executes the submodule status check logic
func runSubmoduleStatus() error {
	// Check if .gitmodules exists
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		fmt.Println("ℹ️  서브모듈이 없습니다")
		return nil
	}
	
	// Get list of submodules
	submodules, err := getSubmoduleList()
	if err != nil {
		return fmt.Errorf("서브모듈 목록을 가져올 수 없습니다: %v", err)
	}
	
	if len(submodules) == 0 {
		fmt.Println("ℹ️  서브모듈이 없습니다")
		return nil
	}
	
	// Collect status for each submodule
	statuses := make([]SubmoduleStatus, 0)
	for _, submodule := range submodules {
		status := getSubmoduleStatusInfo(submodule)
		statuses = append(statuses, status)
	}
	
	// Print report
	printSubmoduleStatusReport(statuses)
	
	return nil
}

// getSubmoduleList returns list of submodule paths
func getSubmoduleList() ([]string, error) {
	cmd := exec.Command("git", "submodule", "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	paths := make([]string, 0)
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Format: " 1234567890abcdef path (version)"
		// or: "+1234567890abcdef path (version)" (modified)
		// or: "-1234567890abcdef path (version)" (uninitialized)
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			paths = append(paths, fields[1])
		}
	}
	
	return paths, nil
}

// getSubmoduleStatusInfo collects status information for a single submodule
func getSubmoduleStatusInfo(submodulePath string) SubmoduleStatus {
	status := SubmoduleStatus{
		Name: filepath.Base(submodulePath),
		Path: submodulePath,
		Mode: "FULL",
	}
	
	// Enter submodule directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(submodulePath); err != nil {
		status.Size = "N/A"
		return status
	}
	
	// Check if initialized
	if !utils.IsGitRepository() {
		status.Size = "미초기화"
		status.Mode = "미초기화"
		return status
	}
	
	// Get .git folder size
	gitPath := ".git"
	if info, err := os.Stat(gitPath); err == nil && !info.IsDir() {
		// It's a gitfile, read actual path
		if content, err := os.ReadFile(gitPath); err == nil {
			gitDir := strings.TrimPrefix(string(content), "gitdir: ")
			gitDir = strings.TrimSpace(gitDir)
			if !filepath.IsAbs(gitDir) {
				gitDir = filepath.Join(originalDir, gitDir)
			}
			gitPath = gitDir
		}
	}
	
	cmd := exec.Command("du", "-sh", gitPath)
	if output, err := cmd.Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			status.Size = fields[0]
			// Parse size for bytes (rough estimation)
			sizeStr := fields[0]
			multiplier := int64(1)
			if strings.HasSuffix(sizeStr, "K") {
				multiplier = 1024
				sizeStr = strings.TrimSuffix(sizeStr, "K")
			} else if strings.HasSuffix(sizeStr, "M") {
				multiplier = 1024 * 1024
				sizeStr = strings.TrimSuffix(sizeStr, "M")
			} else if strings.HasSuffix(sizeStr, "G") {
				multiplier = 1024 * 1024 * 1024
				sizeStr = strings.TrimSuffix(sizeStr, "G")
			}
			if val, err := strconv.ParseFloat(sizeStr, 64); err == nil {
				status.SizeBytes = int64(val * float64(multiplier))
			}
		}
	}
	
	// Check Partial Clone filter
	status.PartialFilter = utils.GetPartialCloneFilter()
	
	// Check Sparse Checkout
	sparseInfo := utils.GetSparseCheckoutInfo()
	if enabled, ok := sparseInfo["enabled"].(bool); ok {
		status.SparseEnabled = enabled
		if count, ok := sparseInfo["count"].(int); ok {
			status.SparsePaths = count
		}
	}
	
	// Check Shallow status
	shallowInfo := utils.GetShallowInfo()
	if isShallow, ok := shallowInfo["isShallow"].(bool); ok {
		status.IsShallow = isShallow
		if depth, ok := shallowInfo["depth"].(int); ok {
			status.Depth = depth
		}
	}
	
	// Get commit count
	cmd = exec.Command("git", "rev-list", "--count", "HEAD")
	if output, err := cmd.Output(); err == nil {
		if count, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
			status.CommitCount = count
		}
	}
	
	// Determine mode
	if status.PartialFilter != "" || status.SparseEnabled || status.IsShallow {
		status.Mode = "SLIM"
	}
	
	return status
}

// printSubmoduleStatusReport prints the formatted status report
func printSubmoduleStatusReport(statuses []SubmoduleStatus) {
	fmt.Println("╭─────────────────────────────────────────╮")
	fmt.Println("│ 🔍 서브모듈 최적화 상태                │")
	fmt.Println("├─────────────────────────────────────────┤")
	
	totalSize := int64(0)
	optimizedCount := 0
	
	for i, status := range statuses {
		if i > 0 {
			fmt.Println("├─────────────────────────────────────────┤")
		}
		
		// Submodule name
		fmt.Printf("│ 📦 %-36s │\n", utils.TruncateString(status.Path, 36))
		
		// Mode
		modeDisplay := status.Mode
		if status.Mode == "SLIM" {
			modeDisplay = "SLIM (최적화됨)"
			optimizedCount++
		} else if status.Mode == "FULL" {
			modeDisplay = "FULL (전체)"
		}
		fmt.Printf("│ • 모드: %-31s │\n", modeDisplay)
		
		// Size
		if status.Size != "" {
			fmt.Printf("│ • 크기: %-31s │\n", status.Size)
			totalSize += status.SizeBytes
		}
		
		// Shallow status
		if status.IsShallow {
			fmt.Printf("│ • Shallow: depth=%-22d │\n", status.Depth)
		} else if status.CommitCount > 0 {
			fmt.Printf("│ • 전체 히스토리 (커밋 %s개) %s│\n", 
				formatNumber(status.CommitCount),
				strings.Repeat(" ", 13-len(fmt.Sprintf("%d", status.CommitCount))))
		}
		
		// Partial Clone filter
		if status.PartialFilter != "" {
			fmt.Printf("│ • 필터: %-31s │\n", status.PartialFilter)
		}
		
		// Sparse Checkout
		if status.SparseEnabled {
			fmt.Printf("│ • Sparse: %d개 경로 %s│\n", 
				status.SparsePaths,
				strings.Repeat(" ", 20-len(fmt.Sprintf("%d", status.SparsePaths))))
		}
	}
	
	// Summary
	fmt.Println("├─────────────────────────────────────────┤")
	fmt.Println("│ 📊 요약                                 │")
	fmt.Printf("│ • 전체 서브모듈: %d개 %s│\n", 
		len(statuses),
		strings.Repeat(" ", 18-len(fmt.Sprintf("%d", len(statuses)))))
	fmt.Printf("│ • 최적화됨: %d개 %s│\n", 
		optimizedCount,
		strings.Repeat(" ", 23-len(fmt.Sprintf("%d", optimizedCount))))
	
	// Total size
	totalSizeStr := utils.FormatSize(totalSize)
	fmt.Printf("│ • 전체 크기: %-26s │\n", totalSizeStr)
	
	// Potential savings
	if optimizedCount < len(statuses) {
		// Estimate: SLIM mode saves ~90% of space
		potentialSavings := totalSize * int64(len(statuses)-optimizedCount) * 9 / 10
		savingsStr := utils.FormatSize(potentialSavings)
		fmt.Printf("│ • 절약 가능: ~%-25s │\n", savingsStr)
	}
	
	fmt.Println("╰─────────────────────────────────────────╯")
	
	// Recommendations
	if optimizedCount < len(statuses) {
		fmt.Println("\n💡 권장사항:")
		fmt.Println("• 모든 서브모듈 최적화: ga opt submodule to-slim")
		fmt.Println("• 특정 서브모듈 최적화: cd <서브모듈경로> && ga opt quick to-slim")
	}
}

// formatNumber formats a number with comma separators
func formatNumber(n int) string {
	str := fmt.Sprintf("%d", n)
	if n < 1000 {
		return str
	}
	
	result := ""
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}