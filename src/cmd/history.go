package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"github.com/spf13/cobra"
)

// ANSI 색상 코드 상수 추가
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[37m"
	colorBold    = "\033[1m"
)

// 패키지 레벨 변수로 변경
var lastChar rune
var lastKeyTime time.Time

// 패키지 레벨 변수로 merge commit 해시 맵 추가
var mergeCommits map[string]bool

// 패키지 레벨 변수로 충돌 해결 커밋 해시 맵 추가
var conflictCommits map[string]bool

type CommitNode struct {
	Hash            string
	Message         string
	Author          string
	Committer       string
	Date            time.Time
	Branches        []string
	Tags            []string
	Parents         []string
	IsHead          bool
	IsMerge         bool
	IsConflictResolved bool  // 충돌 해결 커밋 여부
	Graph           string
	IsGraphOnly     bool
}

type BranchInfo struct {
	Name       string
	Color      string
	Current    bool
	Remote     bool
	LastCommit string
}

// 페이지 단위 커밋 로딩을 위한 구조체
type CommitLoader struct {
	PageSize    int
	Skip        int
	TotalCount  int
	LastHash    string
	Filter      HistoryFilter
}

type HistoryView struct {
	Commits       []CommitNode
	Cursor        int
	WindowStart   int        // 현재 윈도우의 시작 인덱스
	WindowSize    int        // 윈도우 크기 (기본 30)
	TotalCommits  int        // 전체 커밋 수
	Filter        HistoryFilter
	Layout        ScreenLayout
	CurrentBranch string
	loader        *CommitLoader  // 커밋 로더 추가
	GlobalSearch  bool  // 전체 히스토리 검색 모드
}

type HistoryFilter struct {
	BranchPattern string
	Author        string
	Since         time.Duration
}

type ScreenLayout struct {
	Width      int
	Height     int
	GraphWidth int
	InfoWidth  int
}

// NewHistoryCmd returns the history command
func NewHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history [branch]",
		Short: "대화형 Git 히스토리 시각화",
		Long: `Git 히스토리를 대화형으로 시각화하여 보여줍니다.
브랜치 그래프와 커밋 정보를 터미널 크기에 맞춰 표시하고,
키보드로 네비게이션할 수 있습니다.

사용법:
  workingcli git history           # 현재 브랜치의 히스토리 표시
  workingcli git history [branch]  # 특정 브랜치의 히스토리 표시
  workingcli git history --all     # 모든 브랜치의 히스토리 표시
  workingcli git history --global  # 전체 히스토리에서 병합/충돌 커밋 검색`,
		Run: runHistory,
	}

	cmd.Flags().Bool("all", false, "모든 브랜치 표시")
	cmd.Flags().String("branch", "", "브랜치 패턴 필터 (예: feature/*)")
	cmd.Flags().String("author", "", "작성자 필터")
	cmd.Flags().String("since", "", "기간 필터 (예: 1.week)")
	cmd.Flags().Bool("global", false, "전체 히스토리에서 병합/충돌 커밋 검색")
	return cmd
}

func runHistory(cmd *cobra.Command, args []string) {
	// 필터 설정
	filter := HistoryFilter{}
	if pattern, _ := cmd.Flags().GetString("branch"); pattern != "" {
		filter.BranchPattern = pattern
	}
	if author, _ := cmd.Flags().GetString("author"); author != "" {
		filter.Author = author
	}
	if since, _ := cmd.Flags().GetString("since"); since != "" {
		// 기간 파싱 (예: 1.week, 2.days 등)
		parts := strings.Split(since, ".")
		if len(parts) == 2 {
			var duration time.Duration
			switch parts[1] {
			case "week", "weeks":
				duration = time.Hour * 24 * 7
			case "day", "days":
				duration = time.Hour * 24
			case "hour", "hours":
				duration = time.Hour
			}
			filter.Since = duration
		}
	}

	// 전체 검색 모드 설정
	globalSearch, _ := cmd.Flags().GetBool("global")

	// 히스토리 뷰 초기화
	view := &HistoryView{
		Filter:   filter,
		Layout: ScreenLayout{
			Width:      80,
			Height:    20,
			GraphWidth: 20,
			InfoWidth:  60,
		},
		GlobalSearch: globalSearch,
	}

	// 터미널 크기 가져오기
	if w, h, err := getTerminalSize(); err == nil {
		view.Layout.Width = w
		view.Layout.Height = h - 5 // 여유 공간
	}

	// 커밋 히스토리 로드
	loadCommits(view, args)

	// 대화형 모드 시작
	handleHistoryInteractive(view)
}

// 페이지 단위로 커밋을 로드하는 함수
func (l *CommitLoader) LoadPage() ([]CommitNode, error) {
	gitArgs := []string{
		"log",
		"--graph",
		"--all",
		fmt.Sprintf("--skip=%d", l.Skip),
		fmt.Sprintf("--max-count=%d", l.PageSize),
		"--pretty=format:%h -%d %s (%aI) <%cn>",
	}

	if l.Filter.Author != "" {
		gitArgs = append(gitArgs, "--author="+l.Filter.Author)
	}
	if l.Filter.Since != 0 {
		gitArgs = append(gitArgs, "--since="+l.Filter.Since.String())
	}

	cmd := exec.Command("git", gitArgs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Git 히스토리 로드 실패: %v", err)
	}

	var commits []CommitNode
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 { continue }

		// 그래프 문자와 나머지 부분 분리
		graphEnd := strings.Index(line, " ")
		if graphEnd == -1 { 
			// 그래프 라인만 있는 경우
			commits = append(commits, CommitNode{
				Graph:       line,
				IsGraphOnly: true,
			})
			continue
		}

		// 그래프 부분 추출 및 처리
		graphPart := line[:graphEnd]
		// 특수한 그래프 패턴 처리
		graphPart = strings.ReplaceAll(graphPart, "|\\ ", "|\\")
		graphPart = strings.ReplaceAll(graphPart, "|/ ", "|/")
		graphPart = strings.ReplaceAll(graphPart, "/\\ ", "/\\")
		graphPart = strings.ReplaceAll(graphPart, "\\| ", "\\|")
		graphPart = strings.ReplaceAll(graphPart, "/| ", "/|")
		graphPart = strings.ReplaceAll(graphPart, "|\\| ", "|\\")
		graphPart = strings.ReplaceAll(graphPart, "|/| ", "|/")
		graphPart = strings.ReplaceAll(graphPart, "\\| ", "\\|")
		graphPart = strings.ReplaceAll(graphPart, " | ", "|")
		graphPart = strings.ReplaceAll(graphPart, "\\/|", "\\|")
		graphPart = strings.ReplaceAll(graphPart, "/\\|", "/|")
		graphPart = strings.ReplaceAll(graphPart, "|/\\", "|/")
		graphPart = strings.ReplaceAll(graphPart, "\\|/", "\\|")
		graphPart = strings.ReplaceAll(graphPart, "/|\\", "/|")
		graphPart = strings.TrimRight(graphPart, " ")

		rest := strings.TrimSpace(line[graphEnd:])

		// 커밋 노드가 없는 그래프 라인인 경우
		if !strings.Contains(rest, " - ") {
			commits = append(commits, CommitNode{
				Graph:       graphPart,
				IsGraphOnly: true,
			})
			continue
		}

		// 해시와 나머지 부분 분리
		parts := strings.SplitN(rest, " -", 2)
		if len(parts) < 2 { continue }

		hash := strings.TrimSpace(parts[0])
		rest = strings.TrimSpace(parts[1])

		// merge commit 판단 로직 수정
		isMergeCommit := mergeCommits[hash]

		// refs, message, date, committer 추출
		var refs, message, dateStr, committer string

		// refs 추출 (있는 경우)
		if strings.HasPrefix(rest, " (") {
			refEnd := strings.Index(rest, ")")
			if refEnd != -1 {
				refs = strings.TrimSpace(rest[2:refEnd])
				rest = strings.TrimSpace(rest[refEnd+1:])
			}
		}

		// committer 추출
		committerStart := strings.LastIndex(rest, "<")
		committerEnd := strings.LastIndex(rest, ">")
		if committerStart != -1 && committerEnd != -1 && committerEnd > committerStart {
			committer = rest[committerStart+1 : committerEnd]
			rest = strings.TrimSpace(rest[:committerStart])
		}

		// date 추출
		dateStart := strings.LastIndex(rest, "(")
		dateEnd := strings.LastIndex(rest, ")")
		if dateStart != -1 && dateEnd != -1 && dateEnd > dateStart {
			dateStr = rest[dateStart+1 : dateEnd]
			rest = strings.TrimSpace(rest[:dateStart])
		}

		// 남은 부분이 메시지
		message = strings.TrimSpace(rest)

		// 커밋 노드 생성
		commit := CommitNode{
			Hash:      hash,
			Message:   message,
			Committer: committer,
			IsMerge:   isMergeCommit,
			Graph:     graphPart,
		}

		// refs 정보 설정
		if refs != "" {
			refList := strings.Split(refs, ", ")
			for _, ref := range refList {
				if strings.HasPrefix(ref, "HEAD") {
					commit.IsHead = true
				} else if strings.HasPrefix(ref, "tag: ") {
					commit.Tags = append(commit.Tags, strings.TrimPrefix(ref, "tag: "))
				} else {
					commit.Branches = append(commit.Branches, ref)
				}
			}
		}

		// 날짜 파싱
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			commit.Date = t
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// 전체 커밋 수를 가져오는 함수
func (l *CommitLoader) GetTotalCount() (int, error) {
	gitArgs := []string{"rev-list", "--count", "HEAD"}
	if l.Filter.Author != "" {
		gitArgs = append(gitArgs, "--author="+l.Filter.Author)
	}
	if l.Filter.Since != 0 {
		gitArgs = append(gitArgs, "--since="+l.Filter.Since.String())
	}

	cmd := exec.Command("git", gitArgs...)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("커밋 수 계산 실패: %v", err)
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("커밋 수 파싱 실패: %v", err)
	}

	return count, nil
}

// 윈도우 이동 시 필요한 커밋만 로드
func (view *HistoryView) LoadCommitsForWindow() error {
	if view.loader == nil {
		view.loader = &CommitLoader{
			PageSize: view.WindowSize * 2,
			Skip:     view.WindowStart,
			Filter:   view.Filter,
		}

		// 최초 로드 시 전체 커밋 수 계산
		total, err := view.loader.GetTotalCount()
		if err != nil {
			return err
		}
		view.TotalCommits = total
	}

	// 현재 윈도우에 필요한 커밋만 로드
	commits, err := view.loader.LoadPage()
	if err != nil {
		return err
	}

	// 현재 윈도우의 커밋들에 대해 병합/충돌 커밋 정보 로드
	mergeCommits = make(map[string]bool)
	conflictCommits = make(map[string]bool)

	if view.GlobalSearch {
		// 전체 히스토리에서 병합/충돌 커밋 검색
		cmd := exec.Command("git", "log", "--merges", "--format=%H")
		output, err := cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(output)))
			for scanner.Scan() {
				hash := strings.TrimSpace(scanner.Text())
				mergeCommits[hash] = true
			}
		}

		cmd = exec.Command("git", "log", "--grep=conflict", "--grep=resolve", "--grep=fix", "-i", "--format=%H")
		output, err = cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(output)))
			for scanner.Scan() {
				hash := strings.TrimSpace(scanner.Text())
				conflictCommits[hash] = true
			}
		}
	} else {
		// 현재 윈도우의 커밋들에 대해서만 검사
		for _, commit := range commits {
			// 병합 커밋 확인
			cmd := exec.Command("git", "rev-list", "--parents", "-n", "1", commit.Hash)
			output, err := cmd.Output()
			if err == nil {
				parents := strings.Fields(strings.TrimSpace(string(output)))
				if len(parents) > 2 { // 부모가 2개 이상이면 병합 커밋
					mergeCommits[commit.Hash] = true
				}
			}

			// 충돌 해결 커밋 확인
			cmd = exec.Command("git", "show", "-s", "--format=%B", commit.Hash)
			output, err = cmd.Output()
			if err == nil {
				message := strings.ToLower(string(output))
				if strings.Contains(message, "conflict") || strings.Contains(message, "resolve") || strings.Contains(message, "fix") {
					conflictCommits[commit.Hash] = true
				}
			}
		}
	}

	view.Commits = commits
	return nil
}

func loadCommits(view *HistoryView, args []string) {
	view.WindowSize = 30
	view.WindowStart = 0
	view.Cursor = 0

	// 초기 윈도우 로드
	if err := view.LoadCommitsForWindow(); err != nil {
		fmt.Printf("히스토리 로드 실패: %v\n", err)
		return
	}
}

// 윈도우 이동 시 호출되는 함수
func (view *HistoryView) MoveWindow(newStart int) error {
	if newStart < 0 {
		newStart = 0
	}
	if newStart >= view.TotalCommits {
		newStart = view.TotalCommits - view.WindowSize
		if newStart < 0 {
			newStart = 0
		}
	}

	// 윈도우가 실제로 이동했을 때만 새로 로드
	if newStart != view.WindowStart {
		view.WindowStart = newStart
		view.loader.Skip = newStart
		return view.LoadCommitsForWindow()
	}

	return nil
}

func handleHistoryInteractive(view *HistoryView) {
	for {
		displayHistory(view)
		action := handleHistoryInput(view)
		if action == "q" {
			return
		}
	}
}

func displayHistory(view *HistoryView) {
	clear := "\033[H\033[2J"
	fmt.Print(clear)

	// 현재 브랜치 정보 표시 (파란색)
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch, _ := cmd.Output()
	branchName := strings.TrimSpace(string(branch))
	view.CurrentBranch = branchName

	fmt.Printf("=== Git History (%s%s%s) ===\n", colorBlue, branchName, colorReset)
	fmt.Printf("커밋 %d-%d / 총 %d개\n", view.WindowStart+1, view.WindowStart+view.WindowSize, view.TotalCommits)
	fmt.Println("단축키:")
	fmt.Printf("%s↑%s/%sk%s: 위로   %s↓%s/%sj%s: 아래로   %s←%s/%sh%s: 이전 브랜치   %s→%s/%sl%s: 다음 브랜치\n", 
		colorBold, colorReset, colorBold, colorReset, 
		colorBold, colorReset, colorBold, colorReset,
		colorBold, colorReset, colorBold, colorReset,
		colorBold, colorReset, colorBold, colorReset)
	fmt.Printf("%sg%s: 윈도우 처음   %sG%s: 윈도우 끝   %sgg%s: 전체 처음   %sGG%s: 전체 끝\n",
		colorBold, colorReset, colorBold, colorReset,
		colorBold, colorReset, colorBold, colorReset)
	fmt.Printf("%sm%s: 다음 병합   %sM%s: 이전 병합   %sc%s: 다음 충돌   %sC%s: 이전 충돌   %sb%s: 브랜치   %sEnter%s: 상세\n\n",
		colorBold, colorReset, colorBold, colorReset,
		colorBold, colorReset, colorBold, colorReset,
		colorBold, colorReset, colorBold, colorReset)

	// 커밋 그래프 표시
	end := view.WindowStart + view.WindowSize
	if end > len(view.Commits) {
		end = len(view.Commits)
	}

	for i := view.WindowStart; i < end; i++ {
		commit := view.Commits[i]
		displayCommit(commit, view.Layout, i == view.Cursor)
	}

	// 브랜치 정보 표시
	branches := getBranchNames()
	currentIdx := -1
	for i, b := range branches {
		if b == branchName {
			currentIdx = i
			break
		}
	}

	fmt.Printf("\n%s=== 브랜치 목록 ===%s\n", colorBold, colorReset)
	for i, b := range branches {
		prefix := "  "
		if i == currentIdx {
			fmt.Printf("%s→ %s%s%s\n", colorBlue, b, colorReset, prefix)
		} else {
			fmt.Printf("%s%s%s\n", prefix, b, colorReset)
		}
	}

	// 현재 커서 위치와 전체 노드 수 표시
	fmt.Printf("\n%s=== 커서 위치 ===%s\n", colorBold, colorReset)
	fmt.Printf("(%d/%d)\n", view.Cursor+1, len(view.Commits))
}

func displayCommit(commit CommitNode, layout ScreenLayout, isCurrent bool) {
	// 그래프 라인만 있는 경우
	if commit.IsGraphOnly {
		// 그래프 색상 처리
		graph := commit.Graph
		if strings.Contains(graph, "*") {
			graph = strings.ReplaceAll(graph, "*", colorYellow+"*"+colorReset)
		}
		if strings.Contains(graph, "|") {
			graph = strings.ReplaceAll(graph, "|", colorBlue+"|"+colorReset)
		}
		if strings.Contains(graph, "/") {
			graph = strings.ReplaceAll(graph, "/", colorBlue+"/"+colorReset)
		}
		if strings.Contains(graph, "\\") {
			graph = strings.ReplaceAll(graph, "\\", colorBlue+"\\"+colorReset)
		}
		prefix := " "
		if isCurrent {
			prefix = colorCyan + "→" + colorReset
		}
		fmt.Printf("%s%s\n", prefix, graph)
		return
	}

	// 그래프 색상 처리
	graph := commit.Graph
	if strings.Contains(graph, "*") {
		graph = strings.ReplaceAll(graph, "*", colorYellow+"*"+colorReset)
	}
	if strings.Contains(graph, "|") {
		graph = strings.ReplaceAll(graph, "|", colorBlue+"|"+colorReset)
	}
	if strings.Contains(graph, "/") {
		graph = strings.ReplaceAll(graph, "/", colorBlue+"/"+colorReset)
	}
	if strings.Contains(graph, "\\") {
		graph = strings.ReplaceAll(graph, "\\", colorBlue+"\\"+colorReset)
	}

	// 현재 선택된 커밋 강조
	prefix := " "
	if isCurrent {
		prefix = colorCyan + "→" + colorReset
	}

	// refs 정보 구성
	var refs string
	if len(commit.Branches) > 0 || len(commit.Tags) > 0 || commit.IsHead {
		refs = " ("
		if commit.IsHead {
			refs += colorGreen + "HEAD -> " + colorReset
			if len(commit.Branches) > 0 {
				refs += colorBlue + commit.Branches[0] + colorReset
			}
		} else {
			var allRefs []string
			for _, branch := range commit.Branches {
				allRefs = append(allRefs, colorBlue+branch+colorReset)
			}
			for _, tag := range commit.Tags {
				allRefs = append(allRefs, colorYellow+"tag: "+tag+colorReset)
			}
			refs += strings.Join(allRefs, ", ")
		}
		refs += ")"
	}

	// Git log --graph --pretty=format:"%h -%d %s (%aI) <%cn>" 형식으로 출력
	fmt.Printf("%s%s %s -%s %s (%s) <%s>\n",
		prefix,
		graph,
		commit.Hash,
		refs,
		commit.Message,
		commit.Date.Format(time.RFC3339),
		commit.Committer)
}

func handleHistoryInput(view *HistoryView) string {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return ""
	}

	// 화살표 키 시퀀스 처리
	if char == '\x1b' {
		sequence := make([]rune, 2)
		for i := 0; i < 2; i++ {
			r, _, err := reader.ReadRune()
			if err != nil {
				return ""
			}
			sequence[i] = r
		}

		if sequence[0] == '[' {
			switch sequence[1] {
			case 'A': // 위 화살표
				if view.Cursor > 0 {
					view.Cursor--
					// 커서가 윈도우 시작점에 도달하면 이전 윈도우로 이동
					if view.Cursor < view.WindowStart {
						if err := view.MoveWindow(view.Cursor); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
						}
					}
				}
			case 'B': // 아래 화살표
				if view.Cursor < len(view.Commits)-1 {
					view.Cursor++
					// 커서가 윈도우 끝점에 도달하면 다음 윈도우로 이동
					if view.Cursor >= view.WindowStart + view.WindowSize {
						if err := view.MoveWindow(view.Cursor - view.WindowSize + 1); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
						}
					}
				}
			case 'C': // 오른쪽 화살표
				nextBranch(view)
			case 'D': // 왼쪽 화살표
				prevBranch(view)
			}
		}
		lastChar = 0  // 화살표 키 입력 후 lastChar 초기화
		return ""
	}

	// Enter 키 또는 빈 문자열 처리
	if char == '\r' || char == '\n' || char == 0 {
		showCommitDetails(view)
		return ""
	}

	// 기존 키 처리
	switch char {
	case 'q':
		return "q"
	case 'g': // g 키 처리
		if lastChar == 'g' {
			// gg: 전체 Git 히스토리의 최신 커밋으로 이동
			view.WindowStart = 0
			view.Cursor = 0
			if err := view.LoadCommitsForWindow(); err != nil {
				fmt.Printf("윈도우 이동 실패: %v\n", err)
			}
		} else {
			// g: 현재 윈도우의 최신 커밋으로 이동
			view.Cursor = view.WindowStart
		}
		lastChar = char
	case 'G': // G 키 처리
		if lastChar == 'G' {
			// GG: 전체 Git 히스토리의 가장 오래된 커밋으로 이동
			lastWindow := view.TotalCommits - view.WindowSize
			if lastWindow < 0 {
				lastWindow = 0
			}
			view.WindowStart = lastWindow
			if err := view.LoadCommitsForWindow(); err != nil {
				fmt.Printf("윈도우 이동 실패: %v\n", err)
			}
			view.Cursor = len(view.Commits) - 1
		} else {
			// G: 현재 윈도우의 가장 오래된 커밋으로 이동
			lastVisible := view.WindowStart + view.WindowSize
			if lastVisible > len(view.Commits) {
				lastVisible = len(view.Commits)
			}
			view.Cursor = lastVisible - 1
		}
		lastChar = char
	case 'j': // 아래로 이동
		if view.Cursor < len(view.Commits)-1 {
			view.Cursor++
			// 커서가 윈도우 끝점에 도달하면 다음 윈도우로 이동
			if view.Cursor >= view.WindowStart + view.WindowSize {
				if err := view.MoveWindow(view.Cursor - view.WindowSize + 1); err != nil {
					fmt.Printf("윈도우 이동 실패: %v\n", err)
				}
			}
		}
		lastChar = 0
	case 'k': // 위로 이동
		if view.Cursor > 0 {
			view.Cursor--
			// 커서가 윈도우 시작점에 도달하면 이전 윈도우로 이동
			if view.Cursor < view.WindowStart {
				if err := view.MoveWindow(view.Cursor); err != nil {
					fmt.Printf("윈도우 이동 실패: %v\n", err)
				}
			}
		}
		lastChar = 0
	case 'h': // 이전 브랜치
		prevBranch(view)
	case 'l': // 다음 브랜치
		nextBranch(view)
	case 'm': // 다음 merge commit
		findNextMergeCommit(view, true)
	case 'M': // 이전 merge commit
		findNextMergeCommit(view, false)
	case 'c': // 다음 충돌 해결 커밋
		findNextConflictCommit(view, true)
	case 'C': // 이전 충돌 해결 커밋
		findNextConflictCommit(view, false)
	case 'b': // 브랜치 목록
		showBranchList(view)
	case 'f': // 필터 설정
		showFilterSettings(view)
	default:
		lastChar = 0  // 다른 키가 입력되면 lastChar 초기화
	}

	return ""
}

func findNextMergeCommit(view *HistoryView, forward bool) {
	current := view.Cursor
	if forward {
		for i := current + 1; i < len(view.Commits); i++ {
			if mergeCommits[view.Commits[i].Hash] {
				view.Cursor = i
				if view.GlobalSearch {
					// 전체 검색 모드에서는 윈도우 이동
					if view.Cursor >= view.WindowStart + view.WindowSize {
						if err := view.MoveWindow(view.Cursor - view.WindowSize/2); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
							return
						}
						// 윈도우 이동 후 커밋 목록 새로 로드
						if err := view.LoadCommitsForWindow(); err != nil {
							fmt.Printf("커밋 로드 실패: %v\n", err)
							return
						}
					}
				}
				break
			}
		}
	} else {
		for i := current - 1; i >= 0; i-- {
			if mergeCommits[view.Commits[i].Hash] {
				view.Cursor = i
				if view.GlobalSearch {
					// 전체 검색 모드에서는 윈도우 이동
					if view.Cursor < view.WindowStart {
						if err := view.MoveWindow(view.Cursor - view.WindowSize/2); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
							return
						}
						// 윈도우 이동 후 커밋 목록 새로 로드
						if err := view.LoadCommitsForWindow(); err != nil {
							fmt.Printf("커밋 로드 실패: %v\n", err)
							return
						}
					}
				}
				break
			}
		}
	}
}

func showBranchList(view *HistoryView) {
	cmd := exec.Command("git", "branch", "--all")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	fmt.Print(string(output))
	fmt.Println("\nEnter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func showFilterSettings(view *HistoryView) {
	fmt.Println("\n=== Git History 필터 ===")
	fmt.Printf("1. 브랜치 패턴: %s (Enter: 유지)\n", view.Filter.BranchPattern)
	fmt.Printf("2. 작성자: %s (Enter: 유지)\n", view.Filter.Author)
	fmt.Printf("3. 기간: %v (Enter: 유지)\n", view.Filter.Since)
	fmt.Println("4. 필터 초기화")
	fmt.Println("q. 필터 설정 종료")
	fmt.Print("\n선택: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		fmt.Print("브랜치 패턴: ")
		pattern, _ := reader.ReadString('\n')
		view.Filter.BranchPattern = strings.TrimSpace(pattern)
	case "2":
		fmt.Print("작성자: ")
		author, _ := reader.ReadString('\n')
		view.Filter.Author = strings.TrimSpace(author)
	case "3":
		fmt.Print("기간 (예: 1.week): ")
		since, _ := reader.ReadString('\n')
		since = strings.TrimSpace(since)
		// 기간 파싱 (예: 1.week, 2.days 등)
		parts := strings.Split(since, ".")
		if len(parts) == 2 {
			var duration time.Duration
			switch parts[1] {
			case "week", "weeks":
				duration = time.Hour * 24 * 7
			case "day", "days":
				duration = time.Hour * 24
			case "hour", "hours":
				duration = time.Hour
			}
			view.Filter.Since = duration
		}
	case "4":
		view.Filter = HistoryFilter{}
	}

	// 필터 적용 후 커밋 다시 로드
	loadCommits(view, nil)
}

func showCommitDetails(view *HistoryView) bool {
	if view.Cursor >= len(view.Commits) {
		return false
	}

	commit := view.Commits[view.Cursor]
	
	// 화면 초기화
	clear := "\033[H\033[2J"
	fmt.Print(clear)

	// 헤더 출력
	fmt.Printf("\n%s=== 커밋 상세 정보 ===%s\n\n", colorBold, colorReset)
	
	// 커밋 기본 정보
	fmt.Printf("%s커밋 해시:%s %s\n", colorBold, colorReset, commit.Hash)
	if commit.IsHead {
		fmt.Printf("%sHEAD:%s %s\n", colorBold, colorReset, "→ "+view.CurrentBranch)
	}
	
	// 브랜치와 태그 정보
	if len(commit.Branches) > 0 {
		fmt.Printf("%s브랜치:%s %s\n", colorBold, colorReset, strings.Join(commit.Branches, ", "))
	}
	if len(commit.Tags) > 0 {
		fmt.Printf("%s태그:%s %s\n", colorBold, colorReset, strings.Join(commit.Tags, ", "))
	}

	// 커밋 메시지와 작성자 정보 가져오기
	cmd := exec.Command("git", "show", "-s", "--format=%B%n%n작성자: %an <%ae>%n작성일시: %ai%n%n커밋터: %cn <%ce>%n커밋일시: %ci", commit.Hash)
	output, err := cmd.Output()
	if err == nil {
		fmt.Printf("\n%s", string(output))
	}

	// 변경된 파일 목록
	fmt.Printf("\n%s변경된 파일 목록:%s\n", colorBold, colorReset)
	cmd = exec.Command("git", "show", "--stat", commit.Hash)
	if output, err := cmd.Output(); err == nil {
		// 첫 줄(커밋 해시 정보)을 제외하고 출력
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fmt.Println(strings.Join(lines[1:], "\n"))
		}
	}

	// 전체 diff 표시
	fmt.Printf("\n%s변경 내용:%s\n", colorBold, colorReset)
	cmd = exec.Command("git", "show", "--color", commit.Hash)
	if output, err := cmd.Output(); err == nil {
		// 첫 줄(커밋 해시 정보)을 제외하고 출력
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fmt.Println(strings.Join(lines[1:], "\n"))
		}
	}

	fmt.Printf("\n%s=== 네비게이션 ===%s\n", colorBold, colorReset)
	fmt.Printf("j/k: 다음/이전 커밋으로 이동\n")
	fmt.Printf("q/Enter: 이전 화면으로 돌아가기\n\n")

	// 사용자 입력 처리
	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return false
		}

		switch char {
		case 'j':
			if view.Cursor < len(view.Commits)-1 {
				view.Cursor++
				return showCommitDetails(view)
			}
		case 'k':
			if view.Cursor > 0 {
				view.Cursor--
				return showCommitDetails(view)
			}
		case 'q', '\r': // q 또는 Enter
			return false
		}
	}
}

func getBranchNames() []string {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var branches []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		branch := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimPrefix(branch, "* ")
		}
		branches = append(branches, branch)
	}

	return branches
}

func getTerminalSize() (width int, height int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return 80, 24, err
	}

	fmt.Sscanf(string(output), "%d %d", &height, &width)
	return width, height, nil
}

// 브랜치 전환 함수 추가
func switchBranch(view *HistoryView, branch string) {
	// 현재 브랜치 저장
	oldBranch := view.CurrentBranch

	// git checkout 실행
	cmd := exec.Command("git", "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n브랜치 전환 실패: %s -> %s\n", oldBranch, branch)
		fmt.Printf("에러: %s\n", string(output))
		fmt.Println("Enter를 누르면 계속합니다...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	// 브랜치 전환 성공
	view.CurrentBranch = branch
	fmt.Printf("\n브랜치 전환: %s -> %s\n", oldBranch, branch)
	fmt.Printf("결과: %s\n", string(output))
	
	// 브랜치 변경 후 커밋 목록 다시 로드
	loadCommits(view, []string{branch})
	view.Cursor = 0 // 커서 초기화

	// 잠시 메시지 표시
	fmt.Println("Enter를 누르면 계속합니다...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func nextBranch(view *HistoryView) {
	branches := getBranchNames()
	for i, branch := range branches {
		if branch == view.CurrentBranch && i < len(branches)-1 {
			switchBranch(view, branches[i+1])
			break
		}
	}
}

func prevBranch(view *HistoryView) {
	branches := getBranchNames()
	for i, branch := range branches {
		if branch == view.CurrentBranch && i > 0 {
			switchBranch(view, branches[i-1])
			break
		}
	}
}

// 충돌 해결 커밋 이동 함수 추가
func findNextConflictCommit(view *HistoryView, forward bool) {
	current := view.Cursor
	if forward {
		for i := current + 1; i < len(view.Commits); i++ {
			if conflictCommits[view.Commits[i].Hash] {
				view.Cursor = i
				if view.GlobalSearch {
					// 전체 검색 모드에서는 윈도우 이동
					if view.Cursor >= view.WindowStart + view.WindowSize {
						if err := view.MoveWindow(view.Cursor - view.WindowSize/2); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
							return
						}
						// 윈도우 이동 후 커밋 목록 새로 로드
						if err := view.LoadCommitsForWindow(); err != nil {
							fmt.Printf("커밋 로드 실패: %v\n", err)
							return
						}
					}
				}
				break
			}
		}
	} else {
		for i := current - 1; i >= 0; i-- {
			if conflictCommits[view.Commits[i].Hash] {
				view.Cursor = i
				if view.GlobalSearch {
					// 전체 검색 모드에서는 윈도우 이동
					if view.Cursor < view.WindowStart {
						if err := view.MoveWindow(view.Cursor - view.WindowSize/2); err != nil {
							fmt.Printf("윈도우 이동 실패: %v\n", err)
							return
						}
						// 윈도우 이동 후 커밋 목록 새로 로드
						if err := view.LoadCommitsForWindow(); err != nil {
							fmt.Printf("커밋 로드 실패: %v\n", err)
							return
						}
					}
				}
				break
			}
		}
	}
} 