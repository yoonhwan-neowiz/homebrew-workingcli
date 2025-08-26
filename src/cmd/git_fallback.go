package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// HandleGitFallback는 GA에서 처리하지 않는 Git 명령어를 직접 git으로 전달
func HandleGitFallback(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("명령어가 지정되지 않았습니다")
	}
	
	// Git 명령어 목록 (대부분의 일반적인 Git 명령어)
	gitCommands := map[string]bool{
		// 기본 명령어
		"add": true, "status": true, "commit": true, "log": true,
		"diff": true, "branch": true, "checkout": true, "merge": true,
		"rebase": true, "reset": true, "push": true, "pull": true,
		"fetch": true, "clone": true, "init": true,
		
		// 고급 명령어
		"stash": true, "tag": true, "show": true, "blame": true,
		"cherry-pick": true, "revert": true, "clean": true,
		"describe": true, "grep": true, "bisect": true,
		"submodule": true, "remote": true, "worktree": true,
		
		// 파일/디렉토리 작업
		"mv": true, "rm": true,
		
		// 정보 확인
		"ls-files": true, "ls-tree": true, "cat-file": true,
		"rev-parse": true, "rev-list": true, "reflog": true,
		"shortlog": true, "whatchanged": true,
		
		// 설정 및 관리
		"config": true, "gc": true, "fsck": true, "prune": true,
		"notes": true, "archive": true, "bundle": true,
		
		// 협업
		"format-patch": true, "am": true, "apply": true,
		"request-pull": true, "send-email": true,
		
		// Git 2.23+ 신규 명령어
		"switch": true, "restore": true,
		
		// 기타
		"help": true, "version": true,
	}
	
	command := args[0]
	
	// Git 명령어인지 확인
	if gitCommands[command] {
		// git 명령어로 실행
		gitCmd := exec.Command("git", args...)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		gitCmd.Stdin = os.Stdin
		
		if err := gitCmd.Run(); err != nil {
			// Git 실행 오류를 그대로 전달 (Git의 오류 메시지가 더 유용함)
			return nil
		}
		return nil
	}
	
	// Git 명령어도 아닌 경우
	return fmt.Errorf("알 수 없는 명령어: %s\n사용 가능한 명령어를 보려면 'ga --help'를 실행하세요", command)
}