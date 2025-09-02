package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 전역 quiet 모드 플래그
var quietMode bool = false

// SetQuietMode quiet 모드를 설정하는 함수
func SetQuietMode(quiet bool) {
	quietMode = quiet
}

// IsQuietMode 현재 quiet 모드 상태를 확인하는 함수
func IsQuietMode() bool {
	return quietMode
}

// ConfirmWithDefault Y/N 입력을 처리하는 함수
// defaultValue가 true면 Y가 기본값, false면 N이 기본값
// quiet 모드에서는 defaultValue를 자동으로 반환
func ConfirmWithDefault(prompt string, defaultValue bool) bool {
	// quiet 모드에서는 기본값을 자동으로 반환
	if quietMode {
		if defaultValue {
			fmt.Printf("%s (자동 선택: Y)\n", prompt)
		} else {
			fmt.Printf("%s (자동 선택: N)\n", prompt)
		}
		return defaultValue
	}
	
	reader := bufio.NewReader(os.Stdin)
	defaultStr := "Y/n"
	if !defaultValue {
		defaultStr = "y/N"
	}
	
	fmt.Printf("%s (%s): ", prompt, defaultStr)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

// Confirm Y/N 입력을 처리하는 함수 (기본값 없음)
// quiet 모드에서는 파괴적인 작업에 대해 false를 반환 (안전 우선)
func Confirm(prompt string) bool {
	// quiet 모드에서는 안전을 위해 false 반환
	if quietMode {
		fmt.Printf("%s (자동 선택: N - quiet 모드에서는 안전을 위해 거부)\n", prompt)
		return false
	}
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// ConfirmForce quiet 모드에서도 강제로 true를 반환하는 확인 함수
// 안전한 작업에 대해서만 사용
func ConfirmForce(prompt string) bool {
	if quietMode {
		fmt.Printf("%s (자동 선택: Y)\n", prompt)
		return true
	}
	return Confirm(prompt)
}

// UnescapeGitPath Git 출력에서 이스케이프된 한글 파일명을 원래 문자열로 변환
func UnescapeGitPath(s string) string {
	var result strings.Builder
	for len(s) > 0 {
		// Git의 8진수 이스케이프 시퀀스 처리 (\xxx 형식)
		if strings.HasPrefix(s, "\\") && len(s) >= 4 {
			// 8진수 3자리 패턴 확인
			if match := regexp.MustCompile(`^\\([0-7]{3})`).FindStringSubmatch(s); len(match) > 1 {
				if num, err := strconv.ParseInt(match[1], 8, 32); err == nil {
					result.WriteRune(rune(num))
					s = s[4:]
					continue
				}
			}
		}
		// Git의 유니코드 이스케이프 시퀀스 처리 (\uXXXX 형식)
		if strings.HasPrefix(s, "\\u") && len(s) >= 6 {
			if num, err := strconv.ParseInt(s[2:6], 16, 32); err == nil {
				result.WriteRune(rune(num))
				s = s[6:]
				continue
			}
		}
		r, size := utf8.DecodeRuneInString(s)
		result.WriteRune(r)
		s = s[size:]
	}
	return result.String()
}

// ProcessGitOutput Git 명령어 출력을 처리하여 한글 파일명을 올바르게 변환
func ProcessGitOutput(output string) string {
	if output == "" {
		return output
	}

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		// 따옴표로 둘러싸인 부분 처리
		if strings.Contains(line, "\"") {
			parts := strings.Split(line, "\"")
			for j := range parts {
				if j%2 == 1 { // 따옴표 안의 내용만 처리
					parts[j] = UnescapeGitPath(parts[j])
				}
			}
			lines[i] = strings.Join(parts, "\"")
			continue
		}

		// 일반적인 경로 패턴 처리
		pathPattern := regexp.MustCompile(`(?:^|\s)([^\s]+/[^\s]*)`)
		lines[i] = pathPattern.ReplaceAllStringFunc(line, func(match string) string {
			trimmed := strings.TrimSpace(match)
			return strings.Replace(match, trimmed, UnescapeGitPath(trimmed), 1)
		})
	}
	return strings.Join(lines, "\n")
}

// ProcessGitPaths Git 파일 경로 목록을 처리하여 한글 파일명을 올바르게 변환
func ProcessGitPaths(paths []string) []string {
	if len(paths) == 0 {
		return paths
	}

	result := make([]string, len(paths))
	for i, path := range paths {
		// 경로에서 이스케이프된 문자 처리
		unescaped := UnescapeGitPath(path)
		
		// 경로 구분자 정규화
		normalized := strings.ReplaceAll(unescaped, "\\", "/")
		
		result[i] = normalized
	}
	return result
}

// DecodeGitPath Git 파일 경로를 처리하여 한글 파일명을 올바르게 변환
func DecodeGitPath(path string) string {
	// 경로에서 이스케이프된 문자 처리
	unescaped := UnescapeGitPath(path)
	
	// 경로 구분자 정규화
	normalized := strings.ReplaceAll(unescaped, "\\", "/")
	
	return normalized
}

// HumanizeBytes 바이트 크기를 사람이 읽기 쉬운 형태로 변환
func HumanizeBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// IsSourceCodeFile 파일 확장자를 기준으로 소스 코드 파일 여부를 판단
func IsSourceCodeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	sourceExts := map[string]bool{
		".go":    true,
		".js":    true,
		".ts":    true,
		".py":    true,
		".java":  true,
		".c":     true,
		".cpp":   true,
		".h":     true,
		".hpp":   true,
		".cs":    true,
		".rb":    true,
		".php":   true,
		".swift": true,
		".kt":    true,
		".rs":    true,
		".scala": true,
		".md":    true,
	}
	return sourceExts[ext]
}

// GetDiffForAI는 Git diff를 가져와서 AI 분석에 적합한 형태로 반환합니다.
func GetDiffForAI(files []string, ref string, withDiff bool) (string, error) {
	var result strings.Builder
	var diffs []string
	var fileInfos []string
	totalSize := int64(0)

	for _, rawFile := range files {
		// Git 파일 경로 디코딩
		file := DecodeGitPath(strings.TrimSpace(rawFile))
		if file == "" {
			continue
		}

		// 파일 크기 및 타입 확인
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("파일 정보 가져오기 실패 '%s': %v\n", file, err)
			continue
		}

		totalSize += fileInfo.Size()
		isSourceFile := IsSourceCodeFile(file)

		// Git status 가져오기
		cmd := exec.Command("git", "status", "--porcelain", file)
		output, err := cmd.Output()
		if err != nil || len(output) < 2 {
			fmt.Printf("Git status 가져오기 실패 '%s': %v\n", file, err)
			continue
		}
		status := string(output)[:2]

		statusText := ""
		switch status[0] {
		case 'A':
			statusText = "추가"
		case 'M':
			statusText = "수정"
		case 'D':
			statusText = "삭제"
		case 'R':
			statusText = "이름변경"
		}

		info := fmt.Sprintf("- %s (%s, %s", file, statusText, HumanizeBytes(fileInfo.Size()))
		if isSourceFile {
			info += ", 소스코드"
			}
		info += ")"
		fileInfos = append(fileInfos, info)

			if withDiff {
				// 1MB 크기 제한
				if fileInfo.Size() > 1024*1024 {
					fmt.Printf("경고: '%s' 파일이 너무 큽니다 (1MB 초과). diff를 포함하지 않습니다.\n", file)
					continue
				}

			// 소스 코드 파일만 diff 포함
			if !IsSourceCodeFile(file) {
				fmt.Printf("정보: '%s' 파일은 소스 코드가 아닙니다. diff를 포함하지 않습니다.\n", file)
				continue
		}

			var diffCmd *exec.Cmd
			if ref == "HEAD" {
				// staged 변경사항의 diff
				diffCmd = exec.Command("git", "diff", "--cached", "--no-prefix", "--", file)
			} else {
				// 특정 커밋의 diff
				diffCmd = exec.Command("git", "show", "--pretty=format:", "--no-prefix", ref, "--", file)
			}

			output, err := diffCmd.Output()
			if err != nil {
				fmt.Printf("Diff 가져오기 실패 '%s': %v\n", file, err)
				continue
			}
			
			if len(output) > 0 {
				diffs = append(diffs, string(output))
			}
		}
	}

	// 파일 정보 추가
	result.WriteString("파일 정보:\n")
	result.WriteString(strings.Join(fileInfos, "\n"))
	result.WriteString("\n\n")

	// diff 추가
	if withDiff && len(diffs) > 0 {
		result.WriteString("변경 내용:\n")
		result.WriteString(strings.Join(diffs, "\n"))
	}

	return result.String(), nil
}

// GetAvailableDiskSpace returns available disk space for the current directory
func GetAvailableDiskSpace() (int64, error) {
	cmd := exec.Command("df", "-k", ".")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("디스크 공간 확인 실패: %v", err)
	}
	
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("df 출력 파싱 실패")
	}
	
	// df 출력의 두 번째 줄 파싱
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0, fmt.Errorf("df 출력 형식 오류")
	}
	
	// 네 번째 필드가 사용 가능한 공간 (KB 단위)
	availableKB, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("사용 가능한 공간 파싱 실패: %v", err)
	}
	
	// KB를 바이트로 변환
	return availableKB * 1024, nil
}

// GetAvailableDiskSpaceFormatted returns formatted available disk space
func GetAvailableDiskSpaceFormatted() string {
	available, err := GetAvailableDiskSpace()
	if err != nil {
		return "확인 불가"
	}
	return HumanizeBytes(available)
} 