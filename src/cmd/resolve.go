package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

type CommitInfo struct {
	Hash        string
	Message     string
	Time        string	
	Description string
}

type ConflictContext struct {
	CurrentBranch     string      // 현재 브랜치 (mine)
	TargetBranch      string      // 대상 브랜치 (theirs)
	Operation         string      // MERGE/CHERRY-PICK/REBASE
	CurrentCommits    []CommitInfo // 현재 브랜치 최근 커밋들
	TargetCommits     []CommitInfo // 대상 브랜치 최근 커밋들
}

type ConflictFile struct {
	Path           string   // 파일 경로
	Status         string   // 해결 상태 (Resolved/Unresolved)
	LocalChanges   string   // 현재 브랜치 변경사항
	RemoteChanges  string   // 대상 브랜치 변경사항
	Resolved       bool     // 해결 여부
	BackupOurs     string   // 내 브랜치 백업 파일 경로
	BackupTheirs   string   // 대상 브랜치 백업 파일 경로
	ManualMergeRef string   // 수동 병합 시 기준 버전 (ours/theirs)
	FileType       string   // 파일 타입 (소스코드/텍스트/바이너리)
	FileSize       int64    // 파일 크기
	ConflictCount  int      // 충돌 발생 횟수
	LastResolved   string   // 마지막 해결 시간
}

type ConflictHistory struct {
	FilePath    string    // 파일 경로
	ResolveTime time.Time // 해결 시간
	Resolution  string    // 해결 방법 (ours/theirs/manual)
	Branch1     string    // 첫 번째 브랜치
	Branch2     string    // 두 번째 브랜치
	Operation   string    // 작업 종류 (merge/rebase/cherry-pick)
}

// 충돌 해결 히스토리 저장 경로
const conflictHistoryFile = ".git/ga_conflict_history.json"

// 외부 도구 설정
type ExternalTool struct {
	Name    string   // 도구 이름
	Command string   // 실행 명령어
	Args    []string // 기본 인자
}

var externalTools = []ExternalTool{
	{
		Name:    "VS Code",
		Command: "code",
		Args:    []string{"--wait", "--diff"},
	},
	{
		Name:    "Vim",
		Command: "vimdiff",
		Args:    []string{},
	},
	{
		Name:    "Meld",
		Command: "meld",
		Args:    []string{},
	},
}

func NewResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "대화형 충돌 해결",
		Long: `Git 충돌을 대화형으로 해결할 수 있는 인터페이스를 제공합니다.
충돌이 발생한 파일들을 단계별로 해결하고, 외부 도구를 통한 해결도 지원합니다.

사용법:
  ga resolve        # 모든 충돌 파일에 대해 대화형 해결
  ga resolve [files] # 지정한 파일들만 해결`,
		Run: runResolve,
	}
	return cmd
}

func runResolve(cmd *cobra.Command, args []string) {
	context, err := getConflictContext()
	if err != nil {
		fmt.Println("충돌 상태를 확인할 수 없습니다:", err)
		return
	}

	if context == nil {
		fmt.Println("현재 충돌 상태가 아닙니다.")
		return
	}

	files := getConflictFiles()
	if len(files) == 0 {
		fmt.Println("해결할 충돌이 없습니다.")
		return
	}

	// 특정 파일만 해결
	if len(args) > 0 {
		files = filterConflictFiles(files, args)
		if len(files) == 0 {
			fmt.Println("지정한 파일에 충돌이 없습니다.")
			return
		}
	}

	handleConflictResolution(context, files)
}

func getConflictContext() (*ConflictContext, error) {
	// HEAD 파일 확인으로 현재 작업 상태 파악
	_, err := os.ReadFile(".git/HEAD")
	if err != nil {
		return nil, fmt.Errorf("Git 저장소를 확인할 수 없습니다: %v", err)
	}

	// MERGE_HEAD, CHERRY_PICK_HEAD, REBASE_HEAD 확인
	var operation, targetBranch string
	if _, err := os.Stat(".git/MERGE_HEAD"); err == nil {
		operation = "MERGING"
		// MERGE_MSG에서 브랜치 이름 추출
		if mergeMsg, err := os.ReadFile(".git/MERGE_MSG"); err == nil {
			// "Merge branch 'branch-name'" 형식에서 브랜치 이름 추출
			msgLines := strings.Split(string(mergeMsg), "\n")
			if len(msgLines) > 0 {
				if matches := regexp.MustCompile(`Merge branch '([^']+)'`).FindStringSubmatch(msgLines[0]); len(matches) > 1 {
					targetBranch = matches[1]
				}
			}
		}
		if targetBranch == "" {
			// MERGE_HEAD에서 커밋 해시 읽기
			targetBranch = readGitFile(".git/MERGE_HEAD")
		}
	} else if _, err := os.Stat(".git/CHERRY_PICK_HEAD"); err == nil {
		operation = "CHERRY-PICKING"
		targetBranch = readGitFile(".git/CHERRY_PICK_HEAD")
	} else if _, err := os.Stat(".git/REBASE_HEAD"); err == nil {
		operation = "REBASING"
		targetBranch = readGitFile(".git/REBASE_HEAD")
	} else {
		return nil, nil // 충돌 상태 아님
	}

	// 현재 브랜치 이름 가져오기
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	currentBranch, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("현재 브랜치를 확인할 수 없습니다: %v", err)
	}

	// 현재 브랜치의 최근 커밋들
	currentCommits := getRecentCommits(strings.TrimSpace(string(currentBranch)))

	// 대상 브랜치의 최근 커밋들
	targetCommits := getRecentCommits(strings.TrimSpace(targetBranch))

	return &ConflictContext{
		CurrentBranch:  strings.TrimSpace(string(currentBranch)),
		TargetBranch:   strings.TrimSpace(targetBranch),
		Operation:      operation,
		CurrentCommits: currentCommits,
		TargetCommits:  targetCommits,
	}, nil
}

func getRecentCommits(ref string) []CommitInfo {
	cmd := exec.Command("git", "log", "-n", "2", "--format=%H%n%s%n%ai%n%b%n---", ref)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var commits []CommitInfo
	sections := strings.Split(string(output), "---\n")
	for _, section := range sections {
		if section == "" {
			continue
		}

		lines := strings.Split(strings.TrimSpace(section), "\n")
		if len(lines) < 3 {
			continue
		}

		description := ""
		if len(lines) > 3 {
			description = strings.Join(lines[3:], "\n")
		}

		commits = append(commits, CommitInfo{
			Hash:        lines[0],
			Message:     lines[1],
			Time:        lines[2],
			Description: description,
		})
	}

	return commits
}

func readGitFile(path string) string {
	var out bytes.Buffer
	cmd := exec.Command("git", "show", path)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return ProcessGitOutput(out.String())
}

func getConflictFiles() []ConflictFile {
	var out bytes.Buffer
	cmd := exec.Command("git", "status", "--porcelain", "-z")  // NUL 문자로 구분
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("Git status 실행 중 오류:", err)
		return nil
	}

	// NUL 문자로 구분된 출력을 처리
	output := out.String()
	entries := strings.Split(strings.TrimRight(output, "\x00"), "\x00")
	
	var files []ConflictFile
	for _, entry := range entries {
		if len(entry) < 4 {
			continue
		}
		
		// 충돌 파일 확인 (UU = unmerged, both modified)
		if entry[0] == 'U' || entry[1] == 'U' || (entry[0] == 'D' && entry[1] == 'D') {
			// 파일 경로 디코딩
			path := utils.DecodeGitPath(strings.TrimSpace(entry[3:]))
			
			files = append(files, ConflictFile{
				Path:     path,
				Status:   "Unresolved",
				Resolved: false,
			})
		}
	}
	return files
}

func filterConflictFiles(files []ConflictFile, targets []string) []ConflictFile {
	var filtered []ConflictFile
	for _, file := range files {
		for _, target := range targets {
			if file.Path == target {
				filtered = append(filtered, file)
				break
			}
		}
	}
	return filtered
}

func handleConflictResolution(context *ConflictContext, files []ConflictFile) {
	currentIndex := 0
	for currentIndex < len(files) {
		displayConflictStatus(context, files, currentIndex)
		action := handleUserInput()

		switch action {
		case "1": // 현재 브랜치 변경사항 사용
			if err := resolveUsingMine(files[currentIndex].Path); err != nil {
				fmt.Printf("오류: %v\n", err)
				continue
			}
			files[currentIndex].Resolved = true
			files[currentIndex].LastResolved = time.Now().Format("2006-01-02 15:04:05")
			saveConflictHistory(&files[currentIndex], context, "ours")
			currentIndex++
		case "2": // 대상 브랜치 변경사항 사용
			if err := resolveUsingTheirs(files[currentIndex].Path); err != nil {
				fmt.Printf("오류: %v\n", err)
				continue
			}
			files[currentIndex].Resolved = true
			files[currentIndex].LastResolved = time.Now().Format("2006-01-02 15:04:05")
			saveConflictHistory(&files[currentIndex], context, "theirs")
			currentIndex++
		case "3": // 수동 병합 모드
			if err := handleManualMerge(&files[currentIndex], context); err != nil {
				fmt.Printf("오류: %v\n", err)
				continue
			}
			if files[currentIndex].Resolved {
				files[currentIndex].LastResolved = time.Now().Format("2006-01-02 15:04:05")
				saveConflictHistory(&files[currentIndex], context, "manual")
				currentIndex++
			}
		case "4": // 다음 파일
			currentIndex++
		case "p": // 이전 파일
			if currentIndex > 0 {
				currentIndex--
			}
		case "r": // 상태 새로고침
			checkResolved(&files[currentIndex])
		case "l": // 목록 보기
			displayFileList(files)
		case "h": // 도움말
			showHelp()
		case "t": // 히스토리 보기
			displayConflictHistory(&files[currentIndex])
			fmt.Println("\nEnter를 누르면 계속합니다...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		case "s": // stage 모드로 돌아가기
			if !allResolved(files) {
				fmt.Println("\n아직 해결되지 않은 충돌이 있습니다:")
				displayUnresolvedFiles(files)
				if !ConfirmWithDefault("stage 모드로 돌아가시겠습니까?", false) {
					continue
				}
			}
			cmd := exec.Command("ga", "stage")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Run()
			return
		case "q": // 종료
			if !allResolved(files) {
				fmt.Println("\n아직 해결되지 않은 충돌이 있습니다:")
				displayUnresolvedFiles(files)
				if !ConfirmWithDefault("종료하시겠습니까?", false) {
					continue
				}
			}
			return
		}

		if currentIndex >= len(files) {
			if allResolved(files) {
				fmt.Println("\n모든 충돌이 해결되었습니다!")
				for {
					fmt.Println("다음 중 선택해주세요:")
					fmt.Println("1. merge commit 생성")
					fmt.Println("2. stage로 돌아가기")
					fmt.Println("3. 종료")
					fmt.Print("\n선택: ")

					reader := bufio.NewReader(os.Stdin)
					choice, _ := reader.ReadString('\n')
					choice = strings.TrimSpace(choice)

					switch choice {
					case "1":
						if err := createMergeCommit(); err != nil {
							fmt.Printf("merge commit 생성 실패: %v\n", err)
						} else {
							fmt.Println("merge commit이 생성되었습니다.")
						}
						return
					case "2":
						cmd := exec.Command("ga", "stage")
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Stdin = os.Stdin
						cmd.Run()
						return
					case "3":
						fmt.Println("종료합니다. 나중에 merge commit을 생성하세요.")
						return
					default:
						fmt.Println("잘못된 선택입니다. 다시 선택해주세요.")
						continue
					}
				}
			} else {
				currentIndex = 0 // 미해결된 파일이 있으면 처음으로 돌아감
			}
		}
	}
}

func allResolved(files []ConflictFile) bool {
	for _, file := range files {
		if !file.Resolved {
			return false
		}
	}
	return true
}

func displayUnresolvedFiles(files []ConflictFile) {
	for _, file := range files {
		if !file.Resolved {
			fmt.Printf("[ ] %s\n", file.Path)
		}
	}
}

func confirmExit() bool {
	return ConfirmWithDefault("종료하시겠습니까?", true)
}

func createMergeCommit() error {
	// 기본 merge commit 메시지를 사용하여 커밋
	cmd := exec.Command("git", "commit", "--no-edit")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("merge commit 생성 실패: %v", err)
	}
	return nil
}

func displayConflictStatus(context *ConflictContext, files []ConflictFile, currentIndex int) {
	clear := "\033[H\033[2J"
	fmt.Print(clear)

	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	// 상단 상태바
	fmt.Println("=== Git Assistant - 충돌 해결 모드 ===")
	fmt.Printf("진행률: [")
	resolvedCount := 0
	for _, f := range files {
		if f.Resolved {
			resolvedCount++
		}
	}
	progress := float64(resolvedCount) / float64(len(files)) * 100
	green.Printf("%d%%", int(progress))
	fmt.Printf("] (%d/%d 해결됨)\n\n", resolvedCount, len(files))
	
	// 충돌 파일 목록
	fmt.Printf("충돌 파일 목록 (%d/%d):\n", currentIndex+1, len(files))
	for i, file := range files {
		status := " "
		if file.Resolved {
			status = green.Sprintf("✓")
		} else if i == currentIndex {
			status = yellow.Sprintf("▶")
		}
		
		// 파일 정보
		fileInfo := fmt.Sprintf("%s", file.Path)
		if file.FileType != "" {
			fileInfo += fmt.Sprintf(" (%s, %s)", file.FileType, utils.HumanizeBytes(file.FileSize))
		}
		if file.ConflictCount > 0 {
			fileInfo += red.Sprintf(" [충돌 %d개]", file.ConflictCount)
		}
		
		fmt.Printf("  [%s] %s\n", status, fileInfo)
	}

	fmt.Println("\n=== 충돌 상황 ===")
	fmt.Printf("작업 중이던 브랜치: ")
	green.Printf("%s", context.CurrentBranch)
	fmt.Printf(" (%s)\n", context.CurrentCommits[0].Hash[:7])

	// 현재 브랜치 커밋 정보
	for i, commit := range context.CurrentCommits {
		fmt.Printf("• 최근 커밋%d: ", i+1)
		yellow.Printf("\"%s\"", commit.Message)
		fmt.Printf(" (%s) (%s)\n", commit.Time[:16], commit.Hash[:7])
		if commit.Description != "" {
			cyan.Printf("    • %s\n", commit.Description)
		}
	}

	fmt.Printf("\n가져오려는 브랜치: ")
	blue.Printf("%s", context.TargetBranch)
	if len(context.TargetCommits) > 0 {
		fmt.Printf(" (%s)\n", context.TargetCommits[0].Hash[:7])

		// 대상 브랜치 커밋 정보
		for i, commit := range context.TargetCommits {
			fmt.Printf("• 최근 커밋%d: ", i+1)
			yellow.Printf("\"%s\"", commit.Message)
			fmt.Printf(" (%s) (%s)\n", commit.Time[:16], commit.Hash[:7])
			if commit.Description != "" {
				cyan.Printf("    • %s\n", commit.Description)
			}
		}
	} else {
		fmt.Println()
	}

	red.Printf("\n현재 작업: %s\n", context.Operation)
	
	// 현재 파일 상세 정보
	currentFile := files[currentIndex]
	fmt.Printf("\n=== 현재 파일 정보 ===\n")
	fmt.Printf("파일: %s\n", currentFile.Path)
	fmt.Printf("타입: %s (%s)\n", currentFile.FileType, utils.HumanizeBytes(currentFile.FileSize))
	if currentFile.ConflictCount > 0 {
		fmt.Printf("충돌: %d개\n", currentFile.ConflictCount)
	}
	if currentFile.LastResolved != "" {
		fmt.Printf("마지막 해결: %s\n", currentFile.LastResolved)
	}

	fmt.Println("\n=== 해결 방법 ===")
	fmt.Printf("1. %s의 변경사항 사용 (내 브랜치)\n", context.CurrentBranch)
	if context.TargetBranch != "" {
		fmt.Printf("2. %s의 변경사항 사용 (대상 브랜치)\n", context.TargetBranch)
	} else {
		fmt.Println("2. 대상 브랜치의 변경사항 사용")
	}
	fmt.Println("3. 수동 병합 모드")
	fmt.Println("4. 다음 파일로 건너뛰기")
	fmt.Println("p. 이전 파일")
	fmt.Println("r. 상태 새로고침")
	fmt.Println("l. 파일 목록")
	fmt.Println("h. 도움말")
	fmt.Println("t. 히스토리 보기")
	fmt.Println("s. stage 모드로 돌아가기")
	fmt.Println("q. 종료")
	fmt.Print("\n선택: ")
}

func handleUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func resolveUsingMine(path string) error {
	path = utils.UnescapeGitPath(path)
	cmd := exec.Command("git", "checkout", "--ours", path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return exec.Command("git", "add", path).Run()
}

func resolveUsingTheirs(path string) error {
	path = utils.UnescapeGitPath(path)
	cmd := exec.Command("git", "checkout", "--theirs", path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return exec.Command("git", "add", path).Run()
}

func launchExternalTool(path string) error {
	// 파일명 디코딩
	decodedPath := utils.DecodeGitPath(path)
	
	// 임시 파일 생성
	localFile := decodedPath + ".LOCAL"
	remoteFile := decodedPath + ".REMOTE"
	baseFile := decodedPath + ".BASE"
	mergedFile := decodedPath + ".MERGED"

	// 각 버전 추출
	if err := extractVersions(decodedPath, localFile, remoteFile, baseFile, mergedFile); err != nil {
		return err
	}

	// 사용 가능한 외부 도구 확인
	availableTools := getAvailableTools()
	if len(availableTools) == 0 {
		return fmt.Errorf("사용 가능한 외부 도구가 없습니다")
	}

	// 도구 선택 UI 표시
	fmt.Println("\n=== 외부 도구 선택 ===")
	for i, tool := range availableTools {
		fmt.Printf("%d. %s\n", i+1, tool.Name)
	}
	fmt.Print("\n도구 선택 (기본: 1): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selectedTool ExternalTool
	if input == "" {
		selectedTool = availableTools[0]
	} else {
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(availableTools) {
			return fmt.Errorf("잘못된 선택입니다")
		}
		selectedTool = availableTools[index-1]
	}

	// 선택된 도구로 파일 열기
	args := append(selectedTool.Args, mergedFile)
	if selectedTool.Name == "VS Code" {
		args = append(selectedTool.Args, localFile, remoteFile)
	} else if selectedTool.Name == "Vim" {
		args = append([]string{localFile, remoteFile}, selectedTool.Args...)
	} else if selectedTool.Name == "Meld" {
		args = append(args, localFile, remoteFile)
	}

	cmd := exec.Command(selectedTool.Command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("외부 도구 실행 실패: %v", err)
	}

	// 병합된 파일 복사
	if err := copyMergedFile(mergedFile, decodedPath); err != nil {
		return err
	}

	fmt.Printf("\n%s로 파일을 수정했습니다.\n", selectedTool.Name)
	fmt.Println("'R'을 눌러 상태를 새로고침하세요.")

	// 임시 파일 정리
	cleanupTempFiles(localFile, remoteFile, baseFile, mergedFile)
	return nil
}

func extractVersions(path, localFile, remoteFile, baseFile, mergedFile string) error {
	// LOCAL 버전 (현재 브랜치)
	cmd := exec.Command("git", "show", ":2:"+path)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("LOCAL 버전 추출 실패: %v", err)
	}
	if err := os.WriteFile(localFile, output, 0644); err != nil {
		return err
	}

	// REMOTE 버전 (대상 브랜치)
	cmd = exec.Command("git", "show", ":3:"+path)
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("REMOTE 버전 추출 실패: %v", err)
	}
	if err := os.WriteFile(remoteFile, output, 0644); err != nil {
		return err
	}

	// BASE 버전 (공통 조상)
	cmd = exec.Command("git", "show", ":1:"+path)
	output, err = cmd.Output()
	if err == nil {
		if err := os.WriteFile(baseFile, output, 0644); err != nil {
			return err
		}
	}

	// 현재 작업 중인 파일을 MERGED로 복사
	input, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("현재 파일 읽기 실패: %v", err)
	}
	return os.WriteFile(mergedFile, input, 0644)
}

func getAvailableTools() []ExternalTool {
	var available []ExternalTool
	for _, tool := range externalTools {
		if isToolAvailable(tool.Command) {
			available = append(available, tool)
		}
	}
	return available
}

func isToolAvailable(command string) bool {
	cmd := exec.Command("which", command)
	return cmd.Run() == nil
}

func copyMergedFile(mergedFile, originalFile string) error {
	input, err := os.ReadFile(mergedFile)
	if err != nil {
		return fmt.Errorf("병합된 파일 읽기 실패: %v", err)
	}
	return os.WriteFile(originalFile, input, 0644)
}

func cleanupTempFiles(files ...string) {
	for _, file := range files {
		os.Remove(file)
	}
}

func checkResolved(file *ConflictFile) {
	cmd := exec.Command("git", "status", "--porcelain", file.Path)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// UU가 없으면 해결된 것
	file.Resolved = !strings.HasPrefix(string(output), "UU")
}

func displayFileList(files []ConflictFile) {
	fmt.Println("\n=== 충돌 파일 목록 ===")
	for i, file := range files {
		status := "[ ]"
		if file.Resolved {
			status = "[✓]"
		}
		fmt.Printf("%s %d. %s\n", status, i+1, file.Path)
	}
	fmt.Println("\nEnter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func showHelp() {
	fmt.Println("\n=== 충돌 해결 도움말 ===")
	fmt.Println("파일 해결 방법:")
	fmt.Println("  1: 현재 브랜치 변경사항 사용")
	fmt.Println("  2: 대상 브랜치 변경사항 사용")
	fmt.Println("  3: 수동 병합 모드 진입")
	fmt.Println("  4: 다음 파일로 건너뛰기")
	fmt.Println("\n네비게이션:")
	fmt.Println("  n: 다음 충돌 파일")
	fmt.Println("  p: 이전 충돌 파일")
	fmt.Println("  r: 상태 새로고침")
	fmt.Println("  l: 충돌 파일 목록")
	fmt.Println("\n기타:")
	fmt.Println("  h: 도움말 표시")
	fmt.Println("  q: 종료")
	fmt.Println("\n수동 병합 모드:")
	fmt.Println("  - VS Code나 vimdiff를 통한 시각적 비교")
	fmt.Println("  - 백업 파일 자동 생성 및 관리")
	fmt.Println("  - 충돌 상태 복원 가능")
	fmt.Println("\n종료 시 동작:")
	fmt.Println("  - 모든 충돌 해결 시 자동으로 merge commit 생성")
	fmt.Println("  - 미해결 충돌 있을 시 확인 후 종료")
	fmt.Println("\nEnter를 누르면 이전 화면으로 돌아갑니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func handleResolveInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func showResolveHelp() {
	fmt.Println("\n=== 충돌 해결 도움말 ===")
	fmt.Println("1. 현재 브랜치(mine)의 변경사항을 사용")
	fmt.Println("2. 대상 브랜치(theirs)의 변경사항을 사용")
	fmt.Println("3. VS Code로 수동 해결")
	fmt.Println("4. 다음 파일로 이동")
	fmt.Println("p: 이전 파일로 이동")
	fmt.Println("r: 현재 파일 상태 새로고침")
	fmt.Println("l: 전체 파일 목록 보기")
	fmt.Println("q: 종료")
	fmt.Println("\nEnter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func handleManualMerge(file *ConflictFile, context *ConflictContext) error {
	// 백업 파일 생성
	if err := createBackupFiles(file); err != nil {
		return fmt.Errorf("백업 파일 생성 실패: %v", err)
	}

	// 기준 버전 선택
	fmt.Printf("\n=== 수동 병합 모드 ===\n")
	fmt.Printf("현재 파일: %s\n\n", file.Path)
	fmt.Println("어느 버전을 기준으로 병합하시겠습니까?")
	fmt.Println("1. 현재 브랜치의 변경사항 (내 브랜치)")
	fmt.Println("2. 대상 브랜치의 변경사항 (대상 브랜치)")
	fmt.Println("q. 종료")
	fmt.Print("\n선택: ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var useOurs bool
	switch choice {
	case "1":
		useOurs = true
	case "2":
		useOurs = false
	case "q":
		cleanupBackupFiles(file)
		return nil
	default:
		cleanupBackupFiles(file)
		return fmt.Errorf("잘못된 선택입니다")
	}

	// 선택한 버전으로 초기화
	if err := initializeManualMerge(file, useOurs); err != nil {
		cleanupBackupFiles(file)
		return err
	}

	fmt.Printf("\n선택하신 버전으로 초기화했습니다.\n")
	fmt.Println("백업 파일이 생성되었습니다:")
	fmt.Printf("- 내 브랜치 버전: %s\n", file.BackupOurs)
	fmt.Printf("- 대상 브랜치 버전: %s\n", file.BackupTheirs)

	// 외부 편집기로 파일 열기
	if err := launchExternalTool(file.Path); err != nil {
		cleanupBackupFiles(file)
		return err
	}

	// 병합 완료 여부 확인
	fmt.Println("\n파일을 수동 병합한 후 아래 옵션을 선택하세요:")
	fmt.Println("1. 병합 완료")
	fmt.Println("2. 취소 (충돌 상태로 복원)")
	fmt.Println("q. 종료")
	fmt.Print("\n선택: ")

	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		// 병합 완료 처리
		file.Resolved = true
		cleanupBackupFiles(file)
		return nil
	case "2":
		// 충돌 상태로 복원
		if err := restoreConflictState(file); err != nil {
			return err
		}
		cleanupBackupFiles(file)
		return nil
	default:
		cleanupBackupFiles(file)
		return fmt.Errorf("잘못된 선택입니다")
	}
}

func initializeManualMerge(file *ConflictFile, useOurs bool) error {
	// 백업 파일 경로 설정
	absPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("현재 디렉토리 경로를 가져올 수 없습니다: %v", err)
	}
	
	file.BackupOurs = fmt.Sprintf("%s/%s.ours", absPath, file.Path)
	file.BackupTheirs = fmt.Sprintf("%s/%s.theirs", absPath, file.Path)

	// 백업 파일 생성
	if err := createBackupFiles(file); err != nil {
		return fmt.Errorf("백업 파일 생성 실패: %v", err)
	}

	// 선택한 버전으로 초기화
	var cmd *exec.Cmd
	if useOurs {
		cmd = exec.Command("git", "checkout", "--ours", file.Path)
	} else {
		cmd = exec.Command("git", "checkout", "--theirs", file.Path)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("파일 초기화 실패: %v", err)
	}

	return nil
}

func createBackupFiles(file *ConflictFile) error {
	// 파일명 디코딩
	decodedPath := utils.DecodeGitPath(file.Path)
	
	// 백업 파일 경로 설정
	file.BackupOurs = decodedPath + ".ours"
	file.BackupTheirs = decodedPath + ".theirs"

	// 현재 브랜치 버전 백업
	cmd := exec.Command("git", "show", ":2", decodedPath)
	oursOut, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("현재 브랜치 버전 백업 실패: %v", err)
	}
	if err := os.WriteFile(file.BackupOurs, oursOut, 0644); err != nil {
		return err
	}

	// 대상 브랜치 버전 백업
	cmd = exec.Command("git", "show", ":3", decodedPath)
	theirsOut, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("대상 브랜치 버전 백업 실패: %v", err)
	}
	return os.WriteFile(file.BackupTheirs, theirsOut, 0644)
}

func cleanupBackupFiles(file *ConflictFile) error {
	if err := os.Remove(file.BackupOurs); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(file.BackupTheirs); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func restoreConflictState(file *ConflictFile) error {
	return exec.Command("git", "checkout", "-m", file.Path).Run()
}

// 충돌 해결 히스토리 저장
func saveConflictHistory(file *ConflictFile, context *ConflictContext, resolution string) error {
	history := ConflictHistory{
		FilePath:    file.Path,
		ResolveTime: time.Now(),
		Resolution:  resolution,
		Branch1:     context.CurrentBranch,
		Branch2:     context.TargetBranch,
		Operation:   context.Operation,
	}

	// 기존 히스토리 로드
	histories := loadConflictHistories()
	histories = append(histories, history)

	// 최근 100개만 유지
	if len(histories) > 100 {
		histories = histories[len(histories)-100:]
	}

	// JSON 파일로 저장
	data, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return fmt.Errorf("히스토리 JSON 변환 실패: %v", err)
	}

	if err := os.WriteFile(conflictHistoryFile, data, 0644); err != nil {
		return fmt.Errorf("히스토리 파일 저장 실패: %v", err)
	}

	return nil
}

// 충돌 해결 히스토리 로드
func loadConflictHistories() []ConflictHistory {
	data, err := os.ReadFile(conflictHistoryFile)
	if err != nil {
		return []ConflictHistory{}
	}

	var histories []ConflictHistory
	if err := json.Unmarshal(data, &histories); err != nil {
		fmt.Printf("히스토리 파일 읽기 실패: %v\n", err)
		return []ConflictHistory{}
	}

	return histories
}

// 파일의 이전 충돌 해결 히스토리 조회
func getFileConflictHistory(path string) []ConflictHistory {
	histories := loadConflictHistories()
	var fileHistories []ConflictHistory
	
	for _, h := range histories {
		if h.FilePath == path {
			fileHistories = append(fileHistories, h)
		}
	}
	
	return fileHistories
}

// 충돌 해결 히스토리 표시
func displayConflictHistory(file *ConflictFile) {
	histories := getFileConflictHistory(file.Path)
	if len(histories) == 0 {
		fmt.Println("\n이 파일의 충돌 해결 히스토리가 없습니다.")
		return
	}

	fmt.Printf("\n=== 충돌 해결 히스토리 (%s) ===\n", file.Path)
	for i := len(histories) - 1; i >= 0; i-- {
		h := histories[i]
		fmt.Printf("• %s\n", h.ResolveTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  - 작업: %s (%s ↔ %s)\n", h.Operation, h.Branch1, h.Branch2)
		fmt.Printf("  - 해결: %s\n", h.Resolution)
		if i > 0 {
			fmt.Println()
		}
	}
} 