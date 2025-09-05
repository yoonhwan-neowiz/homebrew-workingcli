package quick

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	
	"workingcli/src/utils"
	"github.com/spf13/cobra"
)

// RefType represents the type of git reference
type RefType int

const (
	LocalBranch RefType = iota
	RemoteTrackingBranch
	Unknown
)

// RefInfo contains parsed information about a git reference
type RefInfo struct {
	Type       RefType
	Remote     string  // "origin" for remote branches
	BranchName string  // "master" or "live59.b/5904.7"
	FullRef    string  // "refs/heads/master" or "refs/remotes/origin/live59.b/5904.7"
	Original   string  // Original input string
}

// NewAutoFindMergeBaseCmd creates the Auto Find Merge Base command
func NewAutoFindMergeBaseCmd() *cobra.Command {
	var forceMode bool
	var quietMode bool
	
	cmd := &cobra.Command{
		Use:   "auto-find-merge-base [branch1] [branch2] ...",
		Aliases: []string{"auto", "auto-find"},
		Short: "브랜치 병합점 자동 찾기",
		Long: `두 개 이상의 브랜치가 만나는 공통 조상 커밋(merge-base)을 자동으로 찾습니다.
필요 시 히스토리를 자동 확장하며 병합 가능성을 판단하는 기준점을 제공합니다.

사용 예:
  ga opt quick auto                           # 대화형 모드
  ga opt quick auto master develop            # 두 브랜치 비교
  ga opt quick auto master origin/live59.a/5907.1 -f -q  # 여러 브랜치를 조용히 강제 실행`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAutoFindMergeBase(args, forceMode, quietMode)
		},
	}
	
	cmd.Flags().BoolVarP(&forceMode, "force", "f", false, "확인 없이 강제 실행")
	cmd.Flags().BoolVarP(&quietMode, "quiet", "q", false, "조용한 모드 (최소 출력)")
	
	return cmd
}

// Lock management for preventing concurrent git operations
var (
	lockMutex sync.Mutex
)

// getLockPath returns the lock file path for the current git repository
func getLockPath() string {
	// Get git directory
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return ".git/locks/deepen.lock" // fallback
	}
	gitDir := strings.TrimSpace(string(output))
	return filepath.Join(gitDir, "locks", "deepen.lock")
}

// acquireLock creates a lock file for exclusive git fetch operations
func acquireLock(quietMode bool) error {
	lockMutex.Lock()
	defer lockMutex.Unlock()
	
	lockPath := getLockPath()
	lockDir := filepath.Dir(lockPath)
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return fmt.Errorf("lock 디렉토리 생성 실패: %w", err)
	}
	
	// Try to create lock file immediately
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err == nil {
		file.Close()
		return nil
	}
	
	if !os.IsExist(err) {
		return fmt.Errorf("lock 파일 생성 실패: %w", err)
	}
	
	// Lock exists, use exponential backoff for retries
	maxRetries := 5
	backoff := 10 * time.Millisecond // Start with 10ms
	
	for i := 0; i < maxRetries; i++ {
		time.Sleep(backoff)
		backoff *= 2 // Double the backoff time: 10ms, 20ms, 40ms, 80ms, 160ms
		
		file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			file.Close()
			return nil
		}
		
		if !os.IsExist(err) {
			return fmt.Errorf("lock 파일 생성 실패: %w", err)
		}
	}
	
	// Ask user or force remove stale lock after brief retries
	if !quietMode {
		fmt.Printf("\n⚠️  락 파일이 사용 중입니다 (%s)\n", lockPath)
		fmt.Print("강제로 삭제하시겠습니까? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("락 획득 포기됨")
		}
	}
	
	// Force remove stale lock
	os.Remove(lockPath)
	file, err = os.Create(lockPath)
	if err != nil {
		return fmt.Errorf("lock 파일 강제 생성 실패: %w", err)
	}
	file.Close()
	return nil
}

// releaseLock removes the lock file
func releaseLock() {
	os.Remove(getLockPath())
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// resolveRef resolves a git reference to its full form
func resolveRef(ref string) (*RefInfo, error) {
	info := &RefInfo{
		Original: ref,
		Type:     Unknown,
	}
	
	// Use git rev-parse to get the full ref name
	cmd := exec.Command("git", "rev-parse", "--symbolic-full-name", ref)
	output, err := cmd.Output()
	
	if err == nil {
		fullRef := strings.TrimSpace(string(output))
		info.FullRef = fullRef
		
		if strings.HasPrefix(fullRef, "refs/heads/") {
			// Local branch
			info.Type = LocalBranch
			info.BranchName = strings.TrimPrefix(fullRef, "refs/heads/")
			info.Remote = "origin" // Default remote
			
			// Try to get upstream info
			upstreamCmd := exec.Command("git", "config", fmt.Sprintf("branch.%s.remote", info.BranchName))
			if remoteOut, err := upstreamCmd.Output(); err == nil {
				info.Remote = strings.TrimSpace(string(remoteOut))
			}
			
		} else if strings.HasPrefix(fullRef, "refs/remotes/") {
			// Remote tracking branch (e.g., refs/remotes/origin/live59.b/5904.7)
			info.Type = RemoteTrackingBranch
			trimmed := strings.TrimPrefix(fullRef, "refs/remotes/")
			
			// Find first '/' to separate remote from branch name
			if idx := strings.Index(trimmed, "/"); idx > 0 {
				info.Remote = trimmed[:idx]
				info.BranchName = trimmed[idx+1:]
			}
		}
	} else {
		// Fallback: parse manually
		if strings.HasPrefix(ref, "origin/") {
			info.Type = RemoteTrackingBranch
			info.Remote = "origin"
			info.BranchName = strings.TrimPrefix(ref, "origin/")
			info.FullRef = fmt.Sprintf("refs/remotes/origin/%s", info.BranchName)
		} else if strings.Contains(ref, "/") {
			parts := strings.SplitN(ref, "/", 2)
			if len(parts) == 2 {
				info.Type = RemoteTrackingBranch
				info.Remote = parts[0]
				info.BranchName = parts[1]
				info.FullRef = fmt.Sprintf("refs/remotes/%s/%s", info.Remote, info.BranchName)
			}
		} else {
			// Assume local branch
			info.Type = LocalBranch
			info.BranchName = ref
			info.Remote = "origin"
			info.FullRef = fmt.Sprintf("refs/heads/%s", ref)
		}
	}
	
	return info, nil
}

// deepenRepository progressively expands the shallow repository
func deepenRepository(refs []*RefInfo, quietMode bool) error {
	// Fixed increment of 100
	increment := 100
	maxDepth := 10000
	
	for depth := increment; depth <= maxDepth; depth += increment {
		if !quietMode {
			fmt.Printf("   📊 히스토리 확장 중... (depth: %d)\n", depth)
		}
		
		// Acquire lock before git operations
		if err := acquireLock(quietMode); err != nil {
			return fmt.Errorf("lock 획득 실패: %w", err)
		}
		
		// Try single fetch with --deepen for all remotes
		cmd := exec.Command("git", "fetch", "--all", 
			"--deepen", fmt.Sprintf("%d", depth),
			"--no-tags")
		
		output, err := cmd.CombinedOutput()
		releaseLock()
		
		if err != nil && !quietMode {
			// Check if it's a meaningful error
			outputStr := string(output)
			if !strings.Contains(outputStr, "already have") &&
			   !strings.Contains(outputStr, "shallow update not allowed") {
				fmt.Printf("   ⚠️  확장 중 경고: %s\n", strings.TrimSpace(outputStr))
			}
		}
		
		// Git operation complete - no artificial delay needed
		
		// Check if we can find merge-base now
		allBranches := make([]string, len(refs))
		for i, ref := range refs {
			allBranches[i] = ref.Original
		}
		
		if len(allBranches) >= 2 {
			mergeBase, err := utils.FindMergeBase(allBranches[0], allBranches[1])
			if err == nil && mergeBase != "" {
				return nil // Success!
			}
		}
	}
	
	// Last resort: unshallow the repository
	if !quietMode {
		fmt.Println("   🔄 전체 히스토리를 가져옵니다...")
	}
	
	if err := acquireLock(quietMode); err != nil {
		return fmt.Errorf("lock 획득 실패: %w", err)
	}
	defer releaseLock()
	
	cmd := exec.Command("git", "fetch", "--unshallow", "--all", "--no-tags")
	if output, err := cmd.CombinedOutput(); err != nil {
		// Try without --unshallow if already unshallow
		cmd = exec.Command("git", "fetch", "--all", "--no-tags")
		cmd.Run()
	} else if !quietMode && len(output) > 0 {
		fmt.Printf("   ℹ️  %s\n", strings.TrimSpace(string(output)))
	}
	
	return nil
}

// findMergeBase finds the merge base between branches
func findMergeBase(branch1, branch2 string, quietMode bool) (string, error) {
	if !quietMode {
		fmt.Printf("   🔍 %s와 %s의 병합점을 찾는 중...\n", branch1, branch2)
	}
	
	// Resolve references
	ref1, err := resolveRef(branch1)
	if err != nil {
		return "", fmt.Errorf("브랜치 %s 해석 실패: %w", branch1, err)
	}
	
	ref2, err := resolveRef(branch2)
	if err != nil {
		return "", fmt.Errorf("브랜치 %s 해석 실패: %w", branch2, err)
	}
	
	// Try to find merge base with current state
	mergeBase, err := utils.FindMergeBase(branch1, branch2)
	if err == nil && mergeBase != "" {
		return mergeBase, nil
	}
	
	// Check if shallow repository
	if !utils.IsShallowRepository() {
		return "", fmt.Errorf("병합점을 찾을 수 없습니다")
	}
	
	// Progressively deepen the repository
	if !quietMode {
		fmt.Println("   ℹ️  Shallow 저장소입니다. 히스토리를 확장합니다...")
	}
	
	if err := deepenRepository([]*RefInfo{ref1, ref2}, quietMode); err != nil {
		return "", fmt.Errorf("히스토리 확장 실패: %w", err)
	}
	
	// Try again after deepening
	mergeBase, err = utils.FindMergeBase(branch1, branch2)
	if err == nil && mergeBase != "" {
		return mergeBase, nil
	}
	
	return "", fmt.Errorf("병합점을 찾을 수 없습니다")
}

// showResult displays the merge base result
func showResult(branch1, branch2, mergeBase string, quietMode bool) {
	if quietMode {
		fmt.Println(mergeBase)
		return
	}
	
	fmt.Printf("   ✅ 병합점: %s\n", mergeBase)
	
	// Get distance from each branch
	dist1 := getDistance(branch1, mergeBase)
	dist2 := getDistance(branch2, mergeBase)
	
	fmt.Printf("   📐 거리: %s (%s), %s (%s)\n", branch1, dist1, branch2, dist2)
	
	// Show commit message
	cmd := exec.Command("git", "log", "-1", "--oneline", mergeBase)
	if output, err := cmd.Output(); err == nil {
		fmt.Printf("   📝 커밋: %s", output)
	}
}

// getDistance calculates the distance between a branch and base
func getDistance(branch, base string) string {
	count, err := utils.GetBranchDistance(branch, base)
	if err != nil {
		return "알 수 없음"
	}
	if count == 0 {
		return "동일"
	}
	return fmt.Sprintf("%d 커밋 ahead", count)
}

// runAutoFindMergeBase is the main execution function
func runAutoFindMergeBase(args []string, forceMode, quietMode bool) error {
	if !quietMode {
		fmt.Println("🔍 브랜치 병합점 자동 찾기")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")
	}
	
	// Check if we're in a git repository
	if !utils.IsGitRepository() {
		return fmt.Errorf("현재 디렉토리는 Git 저장소가 아닙니다")
	}
	
	// Get branches to compare
	var branches []string
	if len(args) > 0 {
		branches = args
	} else {
		// Interactive mode
		currentBranch := utils.GetCurrentBranch()
		if !quietMode {
			fmt.Printf("📍 현재 브랜치: %s\n", currentBranch)
		}
		
		// Get target branch from user
		fmt.Print("\n비교할 브랜치를 입력하세요: ")
		var target string
		fmt.Scanln(&target)
		
		if target == "" {
			return fmt.Errorf("브랜치명을 입력해주세요")
		}
		
		branches = []string{currentBranch, target}
	}
	
	// Validate we have at least 2 branches
	if len(branches) < 2 {
		return fmt.Errorf("최소 2개의 브랜치가 필요합니다")
	}
	
	// Find merge bases
	if len(branches) == 2 {
		// Simple two-branch comparison
		mergeBase, err := findMergeBase(branches[0], branches[1], quietMode)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return err
		}
		showResult(branches[0], branches[1], mergeBase, quietMode)
	} else {
		// Multiple branch comparison
		if !quietMode {
			fmt.Printf("\n📊 %d개 브랜치의 모든 병합점 분석 중...\n", len(branches))
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━")
		}
		
		hasError := false
		for i := 0; i < len(branches); i++ {
			for j := i + 1; j < len(branches); j++ {
				if !quietMode {
					fmt.Printf("\n▶ %s ↔ %s\n", branches[i], branches[j])
				}
				
				mergeBase, err := findMergeBase(branches[i], branches[j], quietMode)
				if err != nil {
					fmt.Printf("❌ %v\n", err)
					hasError = true
					continue
				}
				showResult(branches[i], branches[j], mergeBase, quietMode)
			}
		}
		
		if hasError {
			return fmt.Errorf("일부 브랜치에서 병합점을 찾을 수 없습니다")
		}
	}
	
	return nil
}