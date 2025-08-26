package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewResetCmd() *cobra.Command {
	var (
		hard bool
		soft bool
	)

	cmd := &cobra.Command{
		Use:   "reset [commit]",
		Short: "특정 커밋으로 되돌리기 또는 변경사항 초기화",
		Long: `두 가지 주요 용도로 사용됩니다:

1. 특정 커밋으로 되돌리기 (커밋 지정 시)
   ga reset [commit]        # mixed 모드 (기본값)
   ga reset --soft [commit] # 커밋만 되돌리기
   ga reset --hard [commit] # 워킹 디렉토리까지 되돌리기

2. 현재 변경사항 초기화 (커밋 미지정 시)
   ga reset --mixed        # 스테이징 영역 초기화
   ga reset --soft         # 아무 동작 안함
   ga reset --hard         # 모든 변경사항 초기화

모드 설명:
  --mixed (기본값): 커밋과 스테이징 영역을 되돌림. 워킹 디렉토리는 유지
  --soft: 커밋만 되돌림. 변경사항은 스테이징 영역에 유지
  --hard: 커밋, 스테이징 영역, 워킹 디렉토리 모두 되돌림 (주의 필요)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 현재 HEAD 커밋 해시 저장 (실행 취소용)
			currentHash, err := exec.Command("git", "rev-parse", "HEAD").Output()
			if err != nil {
				return fmt.Errorf("현재 커밋 확인 실패: %v", err)
			}
			currentHashStr := strings.TrimSpace(string(currentHash))

			// 모드 및 경고 메시지 설정
			mode := "mixed"
			if hard {
				mode = "hard"
			} else if soft {
				mode = "soft"
			}

			// Case 1: 특정 커밋으로 되돌리기
			if len(args) > 0 {
				targetCommit := args[0]
				
				// 영향받는 파일 목록 확인 (커밋 간 차이)
				affectedFiles, err := getAffectedFiles(targetCommit)
				if err != nil {
					return err
				}

				// 변경사항 미리보기 표시
				fmt.Printf("\n=== Reset 미리보기 (커밋 되돌리기) ===\n")
				fmt.Printf("모드: %s\n", mode)
				switch mode {
				case "hard":
					fmt.Println("영향: 커밋, 스테이징 영역, 워킹 디렉토리 모두 되돌립니다.")
				case "soft":
					fmt.Println("영향: 커밋만 되돌리고 변경사항은 스테이징 영역에 유지됩니다.")
				case "mixed":
					fmt.Println("영향: 커밋과 스테이징 영역을 되돌리고 워킹 디렉토리는 유지됩니다.")
				}
				fmt.Printf("현재 커밋: %s\n", currentHashStr[:8])
				fmt.Printf("대상 커밋: %s\n", targetCommit)
				fmt.Printf("\n커밋 간 변경된 파일 목록:\n")
				for _, file := range affectedFiles {
					fmt.Printf("  %s\n", file)
				}

				// 현재 워킹 디렉토리 변경사항 확인
				workingChanges, err := getWorkingChanges()
				if err != nil {
					return err
				}
				if len(workingChanges) > 0 {
					fmt.Printf("\n현재 워킹 디렉토리 변경사항:\n")
					for _, change := range workingChanges {
						fmt.Printf("  %s\n", change)
					}
				}

			// Case 2: 현재 변경사항 초기화
			} else {
				// 현재 워킹 디렉토리 변경사항만 확인
				workingChanges, err := getWorkingChanges()
				if err != nil {
					return err
				}

				if len(workingChanges) == 0 {
					return fmt.Errorf("초기화할 변경사항이 없습니다.")
				}

				fmt.Printf("\n=== Reset 미리보기 (변경사항 초기화) ===\n")
				fmt.Printf("모드: %s\n", mode)
				switch mode {
				case "hard":
					fmt.Println("영향: 모든 변경사항이 초기화됩니다.")
				case "soft":
					fmt.Println("영향: 변경사항이 유지됩니다.")
				case "mixed":
					fmt.Println("영향: 스테이징된 변경사항이 초기화됩니다.")
				}
				fmt.Printf("\n초기화될 변경사항 목록:\n")
				for _, change := range workingChanges {
					fmt.Printf("  %s\n", change)
				}
			}

			// hard 모드일 경우 특별 경고
			if hard {
				fmt.Printf("\n⚠️  경고: hard reset은 위 변경사항을 모두 삭제합니다!\n")
				fmt.Printf("이 작업은 되돌릴 수 없습니다. 진행하시겠습니까? (y/N): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					return fmt.Errorf("작업이 취소되었습니다.")
				}
			} else {
				fmt.Printf("\n계속하시려면 Enter를 누르세요 (Ctrl+C로 취소)...")
				fmt.Scanln()
			}

			// git reset 실행
			gitArgs := []string{"reset"}
			switch {
			case hard:
				gitArgs = append(gitArgs, "--hard")
			case soft:
				gitArgs = append(gitArgs, "--soft")
			}
			if len(args) > 0 {
				gitArgs = append(gitArgs, args[0])
			}

			output, err := exec.Command("git", gitArgs...).CombinedOutput()
			if err != nil {
				return fmt.Errorf("git reset 실행 실패: %v\n%s", err, string(output))
			}

			// 성공 메시지 및 복구 방법 안내
			if len(args) > 0 {
				fmt.Printf("\n✓ Reset 완료: %s 모드로 %s로 되돌렸습니다.\n", mode, args[0])
				if hard {
					fmt.Printf("\n복구가 필요한 경우 다음 명령어로 되돌릴 수 있습니다:\n")
					fmt.Printf("  git reset --hard %s\n", currentHashStr[:8])
				}
			} else {
				fmt.Printf("\n✓ Reset 완료: %s 모드로 변경사항을 초기화했습니다.\n", mode)
			}

			return nil
		},
	}

	// 플래그 추가
	cmd.Flags().BoolVar(&hard, "hard", false, "워킹 디렉토리까지 모두 되돌리기 (주의: 변경사항 삭제)")
	cmd.Flags().BoolVar(&soft, "soft", false, "커밋만 되돌리고 변경사항은 스테이징 영역에 유지")

	return cmd
}

// 커밋 간의 변경된 파일 목록 가져오기
func getAffectedFiles(targetCommit string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s...HEAD", targetCommit))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("변경된 파일 목록 확인 실패: %v", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

// 현재 워킹 디렉토리의 변경사항 목록 가져오기
func getWorkingChanges() ([]string, error) {
	// staged 변경사항
	cmd := exec.Command("git", "diff", "--name-status", "--cached")
	stagedOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("staged 변경사항 확인 실패: %v", err)
	}

	// unstaged 변경사항
	cmd = exec.Command("git", "diff", "--name-status")
	unstagedOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("unstaged 변경사항 확인 실패: %v", err)
	}

	// untracked 파일
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	untrackedOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("untracked 파일 확인 실패: %v", err)
	}

	var changes []string
	
	// staged 변경사항 처리
	for _, line := range strings.Split(string(stagedOutput), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			status := line[0:1]
			file := line[2:]
			prefix := ""
			switch status {
			case "M":
				prefix = "[Modified/Staged]"
			case "A":
				prefix = "[Added/Staged]"
			case "D":
				prefix = "[Deleted/Staged]"
			case "R":
				prefix = "[Renamed/Staged]"
			}
			changes = append(changes, fmt.Sprintf("%s %s", prefix, file))
		}
	}

	// unstaged 변경사항 처리
	for _, line := range strings.Split(string(unstagedOutput), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			status := line[0:1]
			file := line[2:]
			prefix := ""
			switch status {
			case "M":
				prefix = "[Modified]"
			case "D":
				prefix = "[Deleted]"
			}
			changes = append(changes, fmt.Sprintf("%s %s", prefix, file))
		}
	}

	// untracked 파일 처리
	for _, line := range strings.Split(string(untrackedOutput), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			changes = append(changes, fmt.Sprintf("[Untracked] %s", line))
		}
	}

	return changes, nil
} 