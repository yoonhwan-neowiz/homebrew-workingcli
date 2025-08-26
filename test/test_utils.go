package test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"
)

// setupTestRepo는 테스트용 Git 저장소를 생성합니다.
func setupTestRepo(t *testing.T) string {
	// 임시 디렉토리 생성
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("임시 디렉토리 생성 실패: %v", err)
	}

	// 현재 디렉토리 저장
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("현재 디렉토리 확인 실패: %v", err)
	}

	// 테스트 디렉토리로 이동
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("디렉토리 이동 실패: %v", err)
	}

	// Git 저장소 초기화
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Fatalf("Git 저장소 초기화 실패: %v", err)
	}

	// Git 설정
	exec.Command("git", "config", "user.name", "test").Run()
	exec.Command("git", "config", "user.email", "test@example.com").Run()

	return currentDir
}

// cleanupGitRepo는 테스트용 Git 저장소를 정리합니다.
func cleanupGitRepo(t *testing.T, originalDir string) {
	// 원래 디렉토리로 복귀
	if err := os.Chdir(originalDir); err != nil {
		t.Errorf("원래 디렉토리로 복귀 실패: %v", err)
	}
}

// captureOutput은 명령어 실행 결과를 캡처합니다.
func captureOutput(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}