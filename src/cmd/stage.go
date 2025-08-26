package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/utils"
	"path/filepath"
	"time"
	"sync"
)

type FileStatus struct {
	Index    int    // 표시 번호
	Path     string // 파일 경로
	Status   string // 상태 (A/M/D/R/C)
	Staged   bool   // stage 여부
	Selected bool   // 선택 여부
	Conflict bool   // 충돌 여부
}

// 전역 변수로 verbose 토글 상태 추가
var verboseEnabled bool = false

// 파일 상태 캐시
type StatusCache struct {
	Files     []FileStatus
	LastCheck time.Time
	Dirty     bool
}

var (
	statusCache StatusCache
	cacheMutex  sync.RWMutex
)

// 캐시 유효 시간 (500ms)
const cacheValidDuration = 500 * time.Millisecond

// NewStageCmd returns the stage command
func NewStageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stage",
		Short: "대화형 stage/unstage 인터페이스",
		Long: `Git 변경사항을 대화형으로 stage/unstage 할 수 있는 인터페이스를 제공합니다.
파일 상태를 시각적으로 표시하고 번호로 선택하여 상태를 변경할 수 있습니다.

사용법:
  ga stage           # 대화형 stage 모드 시작
  ga stage [files]   # 지정한 파일들을 바로 stage
  
필터링 옵션:
  --filter=modified   # 수정된 파일만
  --filter=added      # 새로 추가된 파일만
  --filter=deleted    # 삭제된 파일만
  --filter=renamed    # 이름이 변경된 파일만
  --filter=conflict   # 충돌이 발생한 파일만
  
파일 패턴:
  ga stage "*.go"              # Go 파일만
  ga stage "src/cmd/git/*"     # 특정 디렉토리`,
		Run: RunStage,
	}

	cmd.Flags().String("filter", "", "파일 타입 필터 (modified/added/deleted/renamed/conflict)")
	return cmd
}

// RunStage is exported for use in other packages
func RunStage(cmd *cobra.Command, args []string) {
	filter, _ := cmd.Flags().GetString("filter")

	if len(args) > 0 {
		// 파일 직접 지정 모드
		handleDirectStage(args)
		return
	}

	// 대화형 모드
	handleInteractiveStage(filter, cmd)
}

func handleDirectStage(files []string) {
	// 파일 패턴 처리
	var expandedFiles []string
	for _, pattern := range files {
		if strings.Contains(pattern, "*") {
			// 패턴이 있는 경우 glob 매칭
			matches, err := filepath.Glob(pattern)
			if err == nil && len(matches) > 0 {
				expandedFiles = append(expandedFiles, matches...)
			} else {
				fmt.Printf("경고: 패턴 '%s'와 일치하는 파일이 없습니다.\n", pattern)
			}
		} else {
			expandedFiles = append(expandedFiles, pattern)
		}
	}

	if len(expandedFiles) == 0 {
		fmt.Println("stage할 파일이 없습니다.")
		return
	}

	// 파일 존재 여부 확인
	var validFiles []string
	for _, file := range expandedFiles {
		if _, err := os.Stat(file); err == nil {
			validFiles = append(validFiles, file)
		} else {
			fmt.Printf("경고: 파일 '%s'를 찾을 수 없습니다.\n", file)
		}
	}

	if len(validFiles) == 0 {
		fmt.Println("stage할 유효한 파일이 없습니다.")
		return
	}

	args := []string{"add"}
	args = append(args, validFiles...)

	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		fmt.Printf("파일 stage 실패: %v\n", err)
		return
	}

	fmt.Printf("%d개 파일 stage 완료\n", len(validFiles))
}

func handleInteractiveStage(filter string, cmd *cobra.Command) {
	files := getGitStatus()
	if files == nil {
		return
	}

	// 필터 적용
	if filter != "" {
		files = filterFiles(files, filter)
		if len(files) == 0 {
			fmt.Printf("필터 '%s'에 해당하는 파일이 없습니다.\n", filter)
			return
		}
	}

	for {
		displayFiles(files)
		action := handleStageInput(files, cmd)
		if action == "q" {
			return
		} else if action == "d" {
			showDiff(files)
		} else if action == "s" {
			showSelectedFiles(files)
		} else if action == "r" {
			// 파일 상태 새로고침
			files = getGitStatus()
			continue
		} else if action == "y" {
			// 선택된 파일들 처리
			selectedFiles := getSelectedFiles(files)
			if len(selectedFiles) > 0 {
				// staged 파일은 unstage로, unstaged 파일은 stage로 처리
				var stagedFiles, unstagedFiles []string
				for _, file := range files {
					if file.Selected {
						if file.Staged {
							stagedFiles = append(stagedFiles, file.Path)
						} else {
							unstagedFiles = append(unstagedFiles, file.Path)
						}
					}
				}

				// unstage 처리
				if len(stagedFiles) > 0 {
					if confirmChanges(stagedFiles, true) {
						applyChanges(stagedFiles, true)
					}
				}

				// stage 처리
				if len(unstagedFiles) > 0 {
					if confirmChanges(unstagedFiles, false) {
						applyChanges(unstagedFiles, false)
					}
				}
			}
			// 파일 상태 새로고침
			files = getGitStatus()
			continue
		}
	}
}

func showDiff(files []FileStatus) {
	selected := getSelectedFiles(files)
	if len(selected) == 0 {
		fmt.Println("선택된 파일이 없습니다.")
		return
	}

	for _, file := range selected {
		fmt.Printf("\n=== %s의 변경사항 ===\n", file)
		// 파일의 staged 상태 확인
		var isStaged bool
		for _, f := range files {
			if f.Path == file {
				isStaged = f.Staged
				break
			}
		}

		// 파일 크기 및 타입 확인
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("파일 정보를 읽을 수 없습니다: %v\n", err)
			continue
		}

		// 1MB 크기 제한
		if fileInfo.Size() > 1024*1024 {
			fmt.Printf("파일이 너무 큽니다 (1MB 초과). diff를 표시하지 않습니다.\n")
			continue
		}

		// 파일 확장자 확인
		isSourceFile := utils.IsSourceCodeFile(file)

		if !isSourceFile {
			fmt.Printf("소스 코드 파일이 아닙니다. diff를 표시하지 않습니다.\n")
			continue
		}

		// staged 파일은 --cached 옵션 사용
		args := []string{"diff", "--color"}
		if isStaged {
			args = append(args, "--cached")
		}
		args = append(args, file)

		cmd := exec.Command("git", args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()

		// diff 출력의 한글 처리
		fmt.Print(utils.ProcessGitOutput(out.String()))
	}

	fmt.Println("\nEnter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func showSelectedFiles(files []FileStatus) {
	selected := getSelectedFiles(files)
	if len(selected) == 0 {
		fmt.Println("선택된 파일이 없습니다.")
		return
	}

	fmt.Println("\n=== 선택한 파일 목록 ===")
	for i, file := range selected {
		status := "[ ]"
		for _, f := range files {
			if f.Path == file {
				switch f.Status {
				case "A":
					status = "[A]"
				case "M":
					status = "[M]"
				case "D":
					status = "[D]"
				case "R":
					status = "[R]"
				case "C":
					status = "[C]"
				}
				break
			}
		}
		fmt.Printf("%d. %s %s\n", i+1, status, file)
	}

	fmt.Printf("\n총 %d개 파일이 선택됨\n", len(selected))
	fmt.Println("\nEnter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func confirmChanges(files []string, unstage bool) bool {
	action := "stage"
	if unstage {
		action = "unstage"
	}

	fmt.Println()
	for _, file := range files {
		fmt.Println(file)
	}
	
	return ConfirmWithDefault(fmt.Sprintf("\n다음 파일들을 %s 하시겠습니까?:", action), true)
}

func filterFiles(files []FileStatus, filter string) []FileStatus {
	var filtered []FileStatus
	for _, file := range files {
		switch filter {
		case "modified":
			if file.Status == "M" {
				filtered = append(filtered, file)
			}
		case "added":
			if file.Status == "A" {
				filtered = append(filtered, file)
			}
		case "deleted":
			if file.Status == "D" {
				filtered = append(filtered, file)
			}
		case "renamed":
			if file.Status == "R" {
				filtered = append(filtered, file)
			}
		case "conflict":
			if file.Conflict {
				filtered = append(filtered, file)
			}
		}
	}
	return filtered
}

func getSelectedFiles(files []FileStatus) []string {
	var selected []string
	for _, file := range files {
		if file.Selected {
			selected = append(selected, file.Path)
		}
	}
	return selected
}

func applyChanges(files []string, unstage bool) {
	if len(files) == 0 {
		return
	}

	// 한글 파일명 처리
	files = utils.ProcessGitPaths(files)
	
	args := []string{}
	if unstage {
		args = append(args, "restore", "--staged")
	} else {
		args = append(args, "add")
	}
	args = append(args, files...)

	// Git 명령어 실행
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("파일 %s 실패: %v\n%s\n", 
			map[bool]string{true: "unstage", false: "stage"}[unstage], 
			err, 
			string(output))
	}

	// 캐시를 dirty로 표시
	markCacheDirty()
}

func getGitStatus() []FileStatus {
	cacheMutex.RLock()
	if time.Since(statusCache.LastCheck) < cacheValidDuration && !statusCache.Dirty {
		defer cacheMutex.RUnlock()
		return statusCache.Files
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Git 저장소 확인
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("현재 디렉토리는 Git 저장소가 아닙니다.")
		return nil
	}

	// Git 설정 확인 및 적용
	quotepathCmd := exec.Command("git", "config", "core.quotepath")
	output, err := quotepathCmd.Output()
	if err != nil || strings.TrimSpace(string(output)) != "false" {
		// core.quotepath 설정이 false가 아니면 자동으로 설정
		exec.Command("git", "config", "core.quotepath", "false").Run()
	}

	var out bytes.Buffer
	cmd := exec.Command("git", "status", "--porcelain", "-z")  // NUL 문자로 구분
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("Git status 실행 중 오류:", err)
		return nil
	}

	// NUL 문자로 구분된 출력을 처리
	entries := bytes.Split(out.Bytes(), []byte{0})
	var files []FileStatus
	index := 1

	for _, entry := range entries {
		if len(entry) < 2 {
			continue
		}

		// 상태 코드와 파일 경로 분리
		statusCode := string(entry[:2])
		path := string(entry[3:])

		// 파일 경로 디코딩
		path = utils.DecodeGitPath(path)

		// 상태 확인
		staged := statusCode[0] != ' ' && statusCode[0] != '?' && statusCode[0] != 'U'
		status := "?"
		conflict := false

		switch {
		case statusCode[0] == 'A' || statusCode[1] == 'A':
			status = "A"
		case statusCode[0] == 'M' || statusCode[1] == 'M':
			status = "M"
		case statusCode[0] == 'D' || statusCode[1] == 'D':
			status = "D"
		case statusCode[0] == 'R' || statusCode[1] == 'R':
			status = "R"
		case statusCode[0] == 'U' || statusCode[1] == 'U':
			status = "C"
			conflict = true
		}

		// untracked 디렉토리인 경우 하위 파일들을 재귀적으로 추가
		if statusCode[0] == '?' && strings.HasSuffix(path, "/") {
			subFiles := getUntrackedFiles(path)
			for _, subFile := range subFiles {
				files = append(files, FileStatus{
					Index:    index,
					Path:     subFile,
					Status:   "?",
					Staged:   false,
					Selected: false,
					Conflict: false,
				})
				index++
			}
		} else {
			files = append(files, FileStatus{
				Index:    index,
				Path:     path,
				Status:   status,
				Staged:   staged,
				Selected: false,
				Conflict: conflict,
			})
			index++
		}
	}

	// 캐시 업데이트
	statusCache = StatusCache{
		Files:     files,
		LastCheck: time.Now(),
		Dirty:     false,
	}

	return files
}

// getUntrackedFiles는 지정된 디렉토리 내의 모든 untracked 파일을 재귀적으로 찾습니다.
func getUntrackedFiles(dir string) []string {
	var files []string
	entries, err := os.ReadDir(strings.TrimSuffix(dir, "/"))
	if err != nil {
		return files
	}

	for _, entry := range entries {
		path := dir + entry.Name()
		if entry.IsDir() {
			// 디렉토리인 경우 재귀적으로 탐색
			subFiles := getUntrackedFiles(path + "/")
			files = append(files, subFiles...)
		} else {
			// 일반 파일인 경우 목록에 추가
			files = append(files, path)
		}
	}

	return files
}

func displayFiles(files []FileStatus) {
	clear := "\033[H\033[2J"
	fmt.Print(clear)

	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)
	gray := color.New(color.FgHiBlack)
	orange := color.New(color.FgHiRed)

	// Git 상태 표시
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch, _ := cmd.Output()
	branchName := strings.TrimSpace(string(branch))

	// 특수 상태 확인
	var state string
	if _, err := os.Stat(".git/MERGE_HEAD"); err == nil {
		state = "MERGING"
	} else if _, err := os.Stat(".git/CHERRY_PICK_HEAD"); err == nil {
		state = "CHERRY-PICKING"
	} else if _, err := os.Stat(".git/REBASE_HEAD"); err == nil {
		state = "REBASING"
	} else if _, err := os.Stat(".git/REVERT_HEAD"); err == nil {
		state = "REVERTING"
	}

	fmt.Println("=== 변경된 파일 목록 ===")

	// Staged 파일
	fmt.Println("Changes to be committed:")
	stagedCount := 0
	for _, file := range files {
		if !file.Staged {
			continue
		}
		displayFileStatus(file, green, yellow, red, blue, gray, orange)
		stagedCount++
	}
	if stagedCount == 0 {
		fmt.Println("  (없음)")
	}

	// Modified 파일
	fmt.Println("\nChanges not staged for commit:")
	modifiedCount := 0
	for _, file := range files {
		if !file.Staged && file.Status != "A" {
			displayFileStatus(file, green, yellow, red, blue, gray, orange)
			modifiedCount++
		}
	}
	if modifiedCount == 0 {
		fmt.Println("  (없음)")
	}

	// Untracked 파일
	fmt.Println("\nUntracked files:")
	untrackedCount := 0
	for _, file := range files {
		if !file.Staged && file.Status == "A" {
			displayFileStatus(file, green, yellow, red, blue, gray, orange)
			untrackedCount++
		}
	}
	if untrackedCount == 0 {
		fmt.Println("  (없음)")
	}
	
	fmt.Println()
	fmt.Print("[")
	blue.Printf("%s", branchName)
	if state != "" {
		fmt.Print(" | ")
		red.Printf("%s", state)
	}
	fmt.Println("]")
	
	// verbose 상태 표시 추가
	fmt.Print("[Verbose: ")
	if verboseEnabled {
		green.Print("ON")
	} else {
		gray.Print("OFF")
	}
	fmt.Println("]")

	// 충돌 파일이 있는 경우 안내 메시지 표시
	hasConflict := false
	for _, file := range files {
		if file.Conflict {
			hasConflict = true
			break
		}
	}
	if hasConflict {
		fmt.Println("\n[!] 충돌이 발생한 파일이 있습니다.")
		red.Println("    'x'를 눌러 충돌 해결 모드로 전환할 수 있습니다.")
	}

	fmt.Println("\n선택할 파일 번호를 입력하세요.")
	fmt.Println("  - 번호 (예: 1,3-5): 파일 선택")
	fmt.Println("  - a: 전체 선택  c: 선택 취소  i: 선택 반전")
	fmt.Println("  - d: diff 보기  s: 선택 목록  h: 도움말")
	fmt.Println("  - r: 상태 새로고침  v: verbose 토글  m: AI 커밋 메시지 생성")
	fmt.Println("  - x: 충돌 해결 모드로 전환  y: 변경 적용  q: 종료")
	fmt.Print("선택: ")
}

func displayFileStatus(file FileStatus, green, yellow, red, blue, gray, orange *color.Color) {
	selected := " "
	if file.Selected {
		selected = "V"
	}

	statusColor := gray
	switch file.Status {
	case "A":
		statusColor = green
	case "M":
		statusColor = yellow
	case "D":
		statusColor = gray
	case "R":
		statusColor = blue
	case "C":
		statusColor = orange
	}

	fmt.Printf("[%2d] ", file.Index)
	if file.Selected {
		blue.Printf("[%s]", selected)
	} else {
		gray.Printf("[%s]", selected)
	}
	fmt.Printf(" ")
	statusColor.Printf("[%s]", file.Status)
	fmt.Printf(" %s\n", file.Path)
}

func handleStageInput(files []FileStatus, cmd *cobra.Command) string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "a":
		selectAll(files)
	case "c":
		clearSelection(files)
	case "i":
		invertSelection(files)
	case "r":
		return "r"
	case "h":
		showStageHelp()
	case "d":
		showDiff(files)
	case "v":
		verboseEnabled = !verboseEnabled  // verbose 토글
	case "s":
		showSelectedFiles(files)
	case "m":
		if handleAICommit(files, cmd) {
			return "r"  // 커밋 성공 시 새로고침
		}
	case "x":
		// 충돌 해결 모드로 전환
		hasConflict := false
		for _, file := range files {
			if file.Conflict {
				hasConflict = true
				break
			}
		}
		if !hasConflict {
			fmt.Println("\n현재 충돌 상태가 아닙니다.")
			fmt.Println("\nEnter를 누르면 계속합니다...")
			bufio.NewReader(os.Stdin).ReadString('\n')
			return ""
		}
		resolveCmd := &cobra.Command{}
		runResolve(resolveCmd, []string{})
		return "r"  // resolve 모드 종료 후 새로고침
	case "y", "q":
		return input
	default:
		handleNumberInput(input, files)
	}

	return ""
}

func selectAll(files []FileStatus) {
	for i := range files {
		files[i].Selected = true
	}
}

func clearSelection(files []FileStatus) {
	for i := range files {
		files[i].Selected = false
	}
}

func invertSelection(files []FileStatus) {
	for i := range files {
		files[i].Selected = !files[i].Selected
	}
}

func showStageHelp() {
	fmt.Println("\n=== Stage/Unstage 도움말 ===")
	fmt.Println("파일 선택:")
	fmt.Println("  - 단일 선택: \"1\" (1-N 사이의 번호)")
	fmt.Println("  - 다중 선택: \"1,3,4\" (쉼표로 구분)")
	fmt.Println("  - 범위 선택: \"1-3\" (시작-끝)")
	fmt.Println("  - 복합 선택: \"1-3,5,7-9\" (범위와 개별 선택 조합)")
	fmt.Println("\n명령어:")
	fmt.Println("  a: 모든 파일 선택")
	fmt.Println("  c: 모든 선택 취소")
	fmt.Println("  i: 선택된 항목 반전")
	fmt.Println("  r: 파일 상태 새로고침 (staged 파일 수정 시 필수)")
	fmt.Println("  d: 선택한 파일의 diff 보기")
	fmt.Println("  v: verbose 모드 토글 (AI 커밋 메시지에 diff 포함 여부)")
	fmt.Println("  s: 선택한 파일 목록 보기")
	fmt.Println("  h: 이 도움말 표시")
	fmt.Println("  x: 충돌 해결 모드로 전환")
	fmt.Println("  y: 변경 적용 (staged → unstage, unstaged → stage)")
	fmt.Println("  m: AI 커밋 메시지 생성 모드로 전환")
	fmt.Println("  q: 종료")
	fmt.Println("\nEnter를 누르면 이전 화면으로 돌아갑니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func handleNumberInput(input string, files []FileStatus) {
	parts := strings.Split(input, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			handleRangeSelection(part, files)
		} else {
			handleSingleSelection(part, files)
		}
	}
}

func handleRangeSelection(input string, files []FileStatus) {
	ranges := strings.Split(input, "-")
	if len(ranges) != 2 {
		fmt.Println("잘못된 범위 형식입니다. (예: 1-3)")
		return
	}

	start, err1 := parseFileIndex(ranges[0], len(files))
	end, err2 := parseFileIndex(ranges[1], len(files))
	if err1 != nil || err2 != nil {
		fmt.Printf("잘못된 파일 번호입니다. (1-%d 사이의 번호를 입력하세요)\n", len(files))
		return
	}

	if start > end {
		fmt.Println("잘못된 범위입니다. 시작 번호는 끝 번호보다 작아야 합니다.")
		return
	}

	for i := start - 1; i < end; i++ {
		if i >= 0 && i < len(files) {
			files[i].Selected = !files[i].Selected
		}
	}
}

func handleSingleSelection(input string, files []FileStatus) {
	index, err := parseFileIndex(input, len(files))
	if err != nil {
		fmt.Printf("잘못된 파일 번호입니다. (1-%d 사이의 번호를 입력하세요)\n", len(files))
		return
	}

	if index > 0 && index <= len(files) {
		files[index-1].Selected = !files[index-1].Selected
	}
}

func parseFileIndex(input string, maxFiles int) (int, error) {
	index := 0
	_, err := fmt.Sscanf(strings.TrimSpace(input), "%d", &index)
	if err != nil || index < 1 || index > maxFiles {
		return 0, fmt.Errorf("잘못된 파일 번호")
	}
	return index, nil
}

func handleAICommit(files []FileStatus, cmd *cobra.Command) bool {
	// staged 파일이 있는지 확인
	var hasStaged bool
	for _, file := range files {
		if file.Staged {
			hasStaged = true
			break
		}
	}

	if !hasStaged {
		fmt.Println("\nstage된 파일이 없습니다. 먼저 파일을 stage 해주세요.")
		fmt.Println("\nEnter를 누르면 계속합니다...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return false
	}

	// commit 명령어 실행
	commitCmd := NewCommitCmd()
	if verboseEnabled {
		commitCmd.SetArgs([]string{"-v"}) // verbose 모드가 활성화된 경우에만 -v 플래그 설정
	}
	if err := commitCmd.Execute(); err != nil {
		fmt.Printf("\n커밋 실행 실패: %v\n", err)
		fmt.Println("\nEnter를 누르면 계속합니다...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return false
	}

	return true
}

func markCacheDirty() {
	cacheMutex.Lock()
	statusCache.Dirty = true
	cacheMutex.Unlock()
} 
