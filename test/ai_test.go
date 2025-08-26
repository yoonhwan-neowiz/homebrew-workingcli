package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"workingcli/src/cmd/ai"
)

// TestCommitCmd는 커밋 메시지 생성 명령어를 테스트합니다.
func TestCommitCmd(t *testing.T) {
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("CLAUDE_API_KEY가 설정되지 않아 테스트를 건너뜁니다")
	}

	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	// 파일 수정
	if err := os.WriteFile("test.txt", []byte("modified content"), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}

	// 변경사항 스테이징
	if err := exec.Command("git", "add", "test.txt").Run(); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}

	// 커밋 명령어 생성
	cmd := ai.NewCommitCmd()
	
	// 키워드 없는 커밋
	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("커밋 명령어 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("커밋 메시지가 비어있습니다")
	}

	// 키워드 있는 커밋
	if err := os.WriteFile("test.txt", []byte("another modification"), 0644); err != nil {
		t.Fatalf("파일 수정 실패: %v", err)
	}
	if err := exec.Command("git", "add", "test.txt").Run(); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}

	cmd = ai.NewCommitCmd()
	cmd.SetArgs([]string{"-k", "파일 내용 업데이트"})
	
	output, err = captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("키워드 있는 커밋 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("커밋 메시지가 비어있습니다")
	}
}

// TestAnalyzeCmd는 커밋 분석 명령어를 테스트합니다.
func TestAnalyzeCmd(t *testing.T) {
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("CLAUDE_API_KEY가 설정되지 않아 테스트를 건너뜁니다")
	}

	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	// 여러 커밋 생성
	for i := 1; i <= 3; i++ {
		content := []byte(fmt.Sprintf("content update %d", i))
		if err := os.WriteFile("test.txt", content, 0644); err != nil {
			t.Fatalf("파일 수정 실패: %v", err)
		}
		if err := exec.Command("git", "add", "test.txt").Run(); err != nil {
			t.Fatalf("git add 실패: %v", err)
		}
		if err := exec.Command("git", "commit", "-m", fmt.Sprintf("Update %d", i)).Run(); err != nil {
			t.Fatalf("git commit 실패: %v", err)
		}
	}

	tests := []struct {
		name string
		args []string
	}{
		{"최근 커밋", []string{"--last", "2"}},
		{"날짜 범위", []string{"--since", "2024-01-01"}},
		{"브랜치", []string{"--branch", "main"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ai.NewAnalyzeCmd()
			cmd.SetArgs(tt.args)

			output, err := captureOutput(func() error {
				return cmd.Execute()
			})

			if err != nil {
				t.Errorf("분석 명령어 실행 실패: %v", err)
			}
			if len(output) == 0 {
				t.Error("분석 결과가 비어있습니다")
			}
		})
	}
}

// TestParseCommits는 커밋 파싱 함수를 테스트합니다.
func TestParseCommits(t *testing.T) {
	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	// 테스트용 커밋 로그 생성
	output := `abc123
John Doe
2024-01-04T12:34:56Z
Test commit message
Detailed description
of the changes
---
def456
Jane Smith
2024-01-05T10:20:30Z
Another commit
With multiple
lines of
description
---`

	commits := ai.ParseCommits(output)

	if len(commits) != 2 {
		t.Errorf("파싱된 커밋 수가 잘못되었습니다. expected: 2, got: %d", len(commits))
	}

	// 첫 번째 커밋 검증
	if commits[0].Hash != "abc123" {
		t.Errorf("잘못된 커밋 해시. expected: abc123, got: %s", commits[0].Hash)
	}
	if commits[0].Author != "John Doe" {
		t.Errorf("잘못된 작성자. expected: John Doe, got: %s", commits[0].Author)
	}
	if commits[0].Message != "Test commit message" {
		t.Errorf("잘못된 커밋 메시지. expected: Test commit message, got: %s", commits[0].Message)
	}
	if !strings.Contains(commits[0].Description, "Detailed description") {
		t.Error("커밋 설명이 누락되었습니다")
	}
}

func TestAICommit(t *testing.T) {
	// CLAUDE_API_KEY 환경 변수 확인
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("CLAUDE_API_KEY 환경 변수가 설정되지 않았습니다")
	}

	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	// 테스트 파일 생성 및 커밋
	if err := os.WriteFile("test.txt", []byte("initial content"), 0644); err != nil {
		t.Fatalf("테스트 파일 생성 실패: %v", err)
	}
	if err := exec.Command("git", "add", "test.txt").Run(); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}

	cmd := ai.NewCommitCmd()
	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("commit 명령어 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("commit 메시지가 생성되지 않았습니다")
	}
}

func TestAIAnalyze(t *testing.T) {
	// CLAUDE_API_KEY 환경 변수 확인
	if os.Getenv("CLAUDE_API_KEY") == "" {
		t.Skip("CLAUDE_API_KEY 환경 변수가 설정되지 않았습니다")
	}

	originalDir := setupTestRepo(t)
	defer cleanupGitRepo(t, originalDir)

	// 테스트 커밋 생성
	if err := os.WriteFile("test.txt", []byte("initial content"), 0644); err != nil {
		t.Fatalf("테스트 파일 생성 실패: %v", err)
	}
	if err := exec.Command("git", "add", "test.txt").Run(); err != nil {
		t.Fatalf("git add 실패: %v", err)
	}
	if err := exec.Command("git", "commit", "-m", "Initial commit").Run(); err != nil {
		t.Fatalf("git commit 실패: %v", err)
	}

	cmd := ai.NewAnalyzeCmd()
	output, err := captureOutput(func() error {
		return cmd.Execute()
	})

	if err != nil {
		t.Errorf("analyze 명령어 실행 실패: %v", err)
	}
	if len(output) == 0 {
		t.Error("분석 결과가 생성되지 않았습니다")
	}
}

func TestAICommandRegistration(t *testing.T) {
	rootCmd := ai.NewAICmd()

	expectedCommands := []string{"commit", "analyze"}
	for _, cmdName := range expectedCommands {
		if cmd, _, _ := rootCmd.Find([]string{cmdName}); cmd == nil || cmd == rootCmd {
			t.Errorf("%s 명령어가 등록되지 않았습니다", cmdName)
		}
	}
} 