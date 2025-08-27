package quick

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewToSlimCmd creates the To SLIM conversion command
func NewToSlimCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to-slim",
		Short: "SLIM 모드로 전환 (103GB → 30MB)",
		Long: `저장소를 SLIM 모드로 전환합니다.
103GB → 30MB로 저장소 크기를 대폭 축소합니다.

적용되는 최적화:
- Partial Clone (blob:limit=1m) - 1MB 이상 파일 제외
- Sparse Checkout - 최소 경로만 체크아웃
- Shallow Clone (depth=1) - 최신 커밋만 유지
- GC 실행 - 불필요한 오브젝트 정리

실행 내용:
1. git config core.sparseCheckout true
2. git sparse-checkout init --cone
3. git config remote.origin.partialclonefilter blob:limit=1m
4. git fetch --depth=1
5. git gc --aggressive --prune=now

⚠️ 경고: 백업을 먼저 수행하세요!
예상 시간: 약 5-10분 (네트워크 속도에 따라 다름)`,
		Run: func(cmd *cobra.Command, args []string) {
			runToSlim()
		},
	}
}

// runToSlim executes the SLIM mode conversion
func runToSlim() {
	fmt.Println("🚀 SLIM 모드로 전환을 시작합니다...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 현재 상태 확인
	fmt.Println("\n📊 현재 상태 확인 중...")
	currentStatus := checkCurrentStatus()
	
	if currentStatus.isSlim {
		fmt.Println("✅ 이미 SLIM 모드입니다.")
		return
	}
	
	// 작업 중인 변경사항 확인
	if hasUncommittedChanges() {
		fmt.Println("\n⚠️  작업 중인 변경사항이 있습니다!")
		fmt.Println("다음 중 하나를 선택하세요:")
		fmt.Println("1. git stash로 임시 저장 후 진행")
		fmt.Println("2. 커밋 후 진행")
		fmt.Println("3. 취소")
		
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n선택 (1/2/3): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			fmt.Println("📦 변경사항을 stash에 저장합니다...")
			runGitCommand("stash", "push", "-m", "Auto-stash before SLIM conversion")
		case "2":
			fmt.Println("커밋을 먼저 수행하고 다시 실행해주세요.")
			return
		default:
			fmt.Println("취소되었습니다.")
			return
		}
	}
	
	// 백업 권장
	fmt.Println("\n⚠️  경고: SLIM 전환은 저장소를 크게 변경합니다.")
	fmt.Print("계속하시겠습니까? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("취소되었습니다.")
		return
	}
	
	// 현재 설정 백업
	fmt.Println("\n💾 현재 설정을 백업합니다...")
	backupConfig()
	
	// 설정에서 필터와 경로 읽기
	filterSize := config.GetString("optimize.filter.default")
	if filterSize == "" {
		filterSize = "blob:limit=1m"
	}
	
	sparsePaths := getSparsePaths()
	
	// SLIM 전환 프로세스
	fmt.Println("\n🔄 SLIM 모드 전환 중...")
	startTime := time.Now()
	
	// 1. Sparse Checkout 설정
	fmt.Println("\n[1/5] Sparse Checkout 설정...")
	
	// sparsePaths가 비어있거나 "*"만 있으면 README.md만 포함
	skipSparse := len(sparsePaths) == 0 || (len(sparsePaths) == 1 && sparsePaths[0] == "*")
	
	if skipSparse {
		// 기본값으로 README.md만 포함 (clone-slim과 동일)
		fmt.Println("  → 기본 설정 적용 (README.md만 포함)")
		runGitCommand("config", "core.sparseCheckout", "true")
		runGitCommand("sparse-checkout", "init", "--no-cone")
		runGitCommand("sparse-checkout", "set", "README.md")
		fmt.Println("  → README.md 파일만 유지")
	} else {
		// 사용자 지정 경로가 있는 경우
		runGitCommand("config", "core.sparseCheckout", "true")
		
		// 경로에 개별 파일이 포함되어 있는지 확인
		hasFiles := false
		for _, path := range sparsePaths {
			if !strings.HasSuffix(path, "/") && strings.Contains(path, ".") {
				hasFiles = true
				break
			}
		}
		
		if hasFiles {
			// 파일이 포함된 경우 non-cone 모드 사용
			runGitCommand("sparse-checkout", "init", "--no-cone")
		} else {
			// 디렉토리만 있는 경우 cone 모드 사용 (더 빠름)
			runGitCommand("sparse-checkout", "init", "--cone")
		}
		
		args := append([]string{"sparse-checkout", "set"}, sparsePaths...)
		runGitCommand(args...)
		fmt.Printf("  → %d개 경로 설정됨\n", len(sparsePaths))
	}
	
	// 2. Partial Clone 필터 적용
	fmt.Println("[2/5] Partial Clone 필터 적용...")
	runGitCommand("config", "remote.origin.partialclonefilter", filterSize)
	runGitCommand("config", "remote.origin.promisor", "true")
	runGitCommand("config", "extensions.partialClone", "origin")
	
	// 3. Shallow 설정
	fmt.Println("[3/5] 히스토리 최적화...")
	runGitCommand("pull", "--depth=1")
	
	// 4. 불필요한 객체 정리
	fmt.Println("[4/5] 불필요한 객체 정리...")
	runGitCommand("repack", "-a", "-d")
	runGitCommand("maintenance", "run", "--task=gc")
	
	// 5. 성능 설정 적용
	fmt.Println("[5/5] 성능 최적화 설정...")
	applyPerformanceSettings()
	
	// config에 모드 저장
	config.Set("optimize.mode", "slim")
	
	// 결과 확인
	elapsed := time.Since(startTime)
	fmt.Printf("\n✅ SLIM 모드 전환 완료! (소요 시간: %v)\n", elapsed.Round(time.Second))
	
	// 최종 상태 표시
	fmt.Println("\n📊 최종 상태:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	showFinalStatus()
	
	// stash 복원 여부 확인
	if hasStash() {
		fmt.Print("\n📦 이전에 저장한 변경사항을 복원하시겠습니까? (y/N): ")
		restore, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(restore)) == "y" {
			runGitCommand("stash", "pop")
			fmt.Println("✅ 변경사항이 복원되었습니다.")
		}
	}
}

// RepoStatus represents repository status
type RepoStatus struct {
	isSlim         bool
	hasPartialClone bool
	hasSparseCheckout bool
	isShallow      bool
	gitSize        string
	projectSize    string
}

// checkCurrentStatus checks current repository status
func checkCurrentStatus() RepoStatus {
	status := RepoStatus{}
	
	// Partial Clone 확인
	output, _ := runGitCommand("config", "remote.origin.partialclonefilter")
	status.hasPartialClone = output != ""
	
	// Sparse Checkout 확인
	output, _ = runGitCommand("config", "core.sparseCheckout")
	status.hasSparseCheckout = strings.TrimSpace(output) == "true"
	
	// Shallow 확인
	output, _ = runGitCommand("rev-parse", "--is-shallow-repository")
	status.isShallow = strings.TrimSpace(output) == "true"
	
	// SLIM 모드 판단
	status.isSlim = status.hasPartialClone || status.hasSparseCheckout || status.isShallow
	
	// 크기 확인
	gitUsage := utils.GetDiskUsage()
	if gitSize, ok := gitUsage["git"]; ok {
		status.gitSize = gitSize
	} else {
		status.gitSize = "N/A"
	}
	if totalSize, ok := gitUsage["total"]; ok {
		status.projectSize = totalSize
	} else {
		status.projectSize = "N/A"
	}
	
	return status
}

// hasUncommittedChanges checks for uncommitted changes
func hasUncommittedChanges() bool {
	output, _ := runGitCommand("status", "--porcelain")
	return strings.TrimSpace(output) != ""
}

// hasStash checks if there are stashed changes
func hasStash() bool {
	output, _ := runGitCommand("stash", "list")
	return strings.TrimSpace(output) != ""
}

// backupConfig backs up current Git configuration
func backupConfig() {
	output, _ := runGitCommand("config", "--local", "--list")
	
	backupFile := ".git-config-backup"
	file, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("⚠️  백업 파일 생성 실패: %v\n", err)
		return
	}
	defer file.Close()
	
	file.WriteString(output)
	fmt.Printf("✅ 설정이 %s에 백업되었습니다.\n", backupFile)
}

// getSparsePaths gets sparse checkout paths from config
func getSparsePaths() []string {
	settings := config.GetAll()
	var paths []string
	
	if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
		if sparse, ok := optimize["sparse"].(map[string]interface{}); ok {
			if configPaths, ok := sparse["paths"].([]interface{}); ok {
				for _, path := range configPaths {
					if p, ok := path.(string); ok {
						paths = append(paths, p)
					}
				}
			}
		}
	}
	
	return paths
}

// applyPerformanceSettings applies Git performance optimizations
func applyPerformanceSettings() {
	runGitCommand("config", "core.commitGraph", "true")
	runGitCommand("config", "gc.writeCommitGraph", "true")
	runGitCommand("config", "fetch.writeCommitGraph", "true")
	runGitCommand("config", "core.multiPackIndex", "true")
	runGitCommand("config", "fetch.parallel", "4")
	runGitCommand("config", "gc.autoDetach", "false")
}

// showFinalStatus displays final repository status
func showFinalStatus() {
	status := checkCurrentStatus()
	
	fmt.Printf("모드: SLIM\n")
	fmt.Printf("Partial Clone: %s\n", getStatusIcon(status.hasPartialClone))
	fmt.Printf("Sparse Checkout: %s\n", getStatusIcon(status.hasSparseCheckout))
	fmt.Printf("Shallow: %s\n", getStatusIcon(status.isShallow))
	fmt.Printf(".git 폴더: %s\n", status.gitSize)
	fmt.Printf("프로젝트 전체: %s\n", status.projectSize)
}

// getStatusIcon returns status icon
func getStatusIcon(enabled bool) string {
	if enabled {
		return "✅ 활성"
	}
	return "❌ 비활성"
}

// runGitCommand executes a Git command
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

