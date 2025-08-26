package test

import (
	"testing"

	"workingcli/src/cmd/git"
)

// TestGitStatus는 status 명령어를 테스트합니다.
func TestGitStatus(t *testing.T) {
	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	cmd := git.NewStatusCmd()
	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("status 명령어 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("status 출력이 비어있습니다")
	}
}

// TestGitBranch는 branch 명령어를 테스트합니다.
func TestGitBranch(t *testing.T) {
	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	cmd := git.NewBranchCmd()
	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("branch 명령어 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("branch 출력이 비어있습니다")
	}
}

// TestGitPush는 push 명령어를 테스트합니다.
func TestGitPush(t *testing.T) {
	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	cmd := git.NewPushCmd()
	_, err := captureOutput(func() error {
		return cmd.Execute()
	})

	// 원격 저장소가 없으므로 에러가 발생해야 함
	if err == nil {
		t.Error("push 명령어가 원격 저장소 없이 성공했습니다")
	}
}

// TestGitPull는 pull 명령어를 테스트합니다.
func TestGitPull(t *testing.T) {
	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	cmd := git.NewPullCmd()
	_, err := captureOutput(func() error {
		return cmd.Execute()
	})

	// 원격 저장소가 없으므로 에러가 발생해야 함
	if err == nil {
		t.Error("pull 명령어가 원격 저장소 없이 성공했습니다")
	}
}

// TestGitCommandRegistration은 모든 Git 명령어가 제대로 등록되었는지 테스트합니다.
func TestGitCommandRegistration(t *testing.T) {
	rootCmd := git.NewGitCmd()

	expectedCommands := []string{"status", "branch", "pull", "push"}
	for _, cmdName := range expectedCommands {
		if cmd, _, _ := rootCmd.Find([]string{cmdName}); cmd == nil || cmd == rootCmd {
			t.Errorf("%s 명령어가 등록되지 않았습니다", cmdName)
		}
	}
} 