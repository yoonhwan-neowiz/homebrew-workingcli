package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewCloneSlimCmd creates the Clone SLIM command
func NewCloneSlimCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone-slim [URL] [folder]",
		Short: "새로운 저장소를 최적화 모드로 클론",
		Long: `처음부터 SLIM 모드로 최적화된 상태로 저장소를 클론합니다.
Partial Clone (1MB), Sparse Checkout, Shallow depth=1을 모두 적용합니다.

사용법:
  ga optimized setup clone-slim <URL> [folder]
  
예시:
  ga opt setup clone-slim https://github.com/user/repo.git myproject
  ga opt setup clone-slim https://github.com/user/repo.git`,
		// Args: cobra.RangeArgs(1, 2),  // 일단 Args 검증 제거
		Run: func(cmd *cobra.Command, args []string) {
			// args 수동 검증
			if len(args) < 1 {
				fmt.Println("❌ URL을 입력해주세요.")
				fmt.Println("\n사용법:")
				fmt.Println("  ga optimized setup clone-slim <URL> [folder]")
				return
			}
			executeCloneSlim(args)
		},
	}
}

func executeCloneSlim(args []string) {
	url := args[0]
	
	// 폴더명 결정 (지정안하면 URL에서 추출)
	var folder string
	if len(args) > 1 {
		folder = args[1]
	} else {
		// URL에서 폴더명 추출
		parts := strings.Split(url, "/")
		repoName := parts[len(parts)-1]
		folder = strings.TrimSuffix(repoName, ".git")
	}
	
	// 폴더가 이미 존재하는지 확인
	if _, err := os.Stat(folder); err == nil {
		fmt.Printf("❌ '%s' 폴더가 이미 존재합니다.\n", folder)
		
		// 덮어쓸지 확인
		if !utils.Confirm("기존 폴더를 삭제하고 다시 클론하시겠습니까?") {
			fmt.Println("작업이 취소되었습니다.")
			return
		}
		
		// 기존 폴더 삭제
		if err := os.RemoveAll(folder); err != nil {
			fmt.Printf("❌ 기존 폴더 삭제 실패: %v\n", err)
			return
		}
	}
	
	fmt.Printf("🚀 SLIM 모드로 '%s' 저장소를 클론합니다...\n", url)
	fmt.Printf("   타겟 폴더: %s\n", folder)
	fmt.Println()
	
	// 1. Partial Clone으로 클론 (filter + sparse + no-checkout)
	fmt.Println("1️⃣ Partial Clone으로 저장소 클론 (no-checkout)...")
	cloneCmd := exec.Command("git", "clone", 
		"--filter=blob:limit=1m",
		"--sparse",
		"--no-checkout",  // 파일을 아직 체크아웃하지 않음
		url, 
		folder)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("❌ 클론 실패: %v\n", err)
		return
	}
	fmt.Println("   ✅ 클론 완료 (파일 체크아웃 대기중)")
	
	// 클론된 디렉토리로 이동
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ 현재 디렉토리 저장 실패: %v\n", err)
		return
	}
	
	if err := os.Chdir(folder); err != nil {
		fmt.Printf("❌ '%s' 디렉토리로 이동 실패: %v\n", folder, err)
		return
	}
	defer os.Chdir(originalDir)
	
	// 2. Sparse Checkout 설정
	fmt.Println("\n2️⃣ Sparse Checkout 설정...")
	
	// cone 모드로 초기화 (디렉토리 기반 제어)
	initCmd := exec.Command("git", "sparse-checkout", "init", "--cone")
	if output, err := initCmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ Sparse Checkout 초기화 실패: %v\n", err)
		fmt.Printf("   출력: %s\n", string(output))
		return
	}
	
	// 루트 디렉토리만 포함 (1 depth - .gitmodules, README.md 등 모든 루트 파일)
	setCmd := exec.Command("git", "sparse-checkout", "set", "/")
	if output, err := setCmd.CombinedOutput(); err != nil {
		fmt.Printf("⚠️ Sparse Checkout 기본 경로 설정 실패: %v\n", err)
		fmt.Printf("   출력: %s\n", string(output))
	}
	
	fmt.Println("   ✅ Sparse Checkout 초기화 완료 (루트 파일들 포함)")
	fmt.Println("   ℹ️ 추가 경로는 'ga opt quick expand-slim'으로 설정 가능")
	
	// 3. Shallow 설정 (depth=1) 및 checkout
	fmt.Println("\n3️⃣ Shallow 히스토리 설정 및 파일 체크아웃...")
	
	// Shallow fetch로 최신 커밋만 가져오기
	fetchCmd := exec.Command("git", "fetch", "--depth=1")
	if output, err := fetchCmd.CombinedOutput(); err != nil {
		fmt.Printf("⚠️ Shallow fetch 실패: %v\n", err)
		fmt.Printf("   출력: %s\n", string(output))
	} else {
		fmt.Println("   ✅ Shallow fetch 완료")
	}
	
	// Sparse Checkout 재적용 (checkout 전에 한번 더 확인)
	reapplyCmd := exec.Command("git", "sparse-checkout", "reapply")
	if output, err := reapplyCmd.CombinedOutput(); err != nil {
		fmt.Printf("⚠️ Sparse Checkout 재적용 실패: %v\n", err)
		fmt.Printf("   출력: %s\n", string(output))
	}
	
	// 이제 checkout 수행
	checkoutCmd := exec.Command("git", "checkout", "HEAD")
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ 체크아웃 실패: %v\n", err)
		fmt.Printf("   출력: %s\n", string(output))
		return
	} else {
		fmt.Println("   ✅ 파일 체크아웃 완료")
	}
	
	// 4. 서브모듈 초기화 (무조건 실행)
	fmt.Println("\n4️⃣ 서브모듈 초기화...")
	
	// .gitmodules 체크 없이 무조건 실행
	submoduleCmd := exec.Command("git", "submodule", "update", 
		"--init",
		"--depth=1",
		"--recursive")
	submoduleCmd.Stdout = os.Stdout
	submoduleCmd.Stderr = os.Stderr
	
	if err := submoduleCmd.Run(); err != nil {
		// 서브모듈이 없는 경우도 에러가 아니므로 경고만 표시
		fmt.Printf("ℹ️ 서브모듈 처리 완료 (서브모듈이 없을 수 있음)\n")
	} else {
		fmt.Println("   ✅ 서브모듈 초기화 완료")
	}
	
	// 5. 태그 제거 및 fetch 차단
	fmt.Println("\n5️⃣ 태그 최적화 (No-Tag 모드)...")
	
	// 로컬 태그 개수 확인
	tagCountCmd := exec.Command("git", "tag")
	tagOutput, _ := tagCountCmd.Output()
	var tagCount int
	if len(tagOutput) > 0 {
		tags := strings.Split(strings.TrimSpace(string(tagOutput)), "\n")
		for _, tag := range tags {
			if strings.TrimSpace(tag) != "" {
				tagCount++
			}
		}
	}
	
	if tagCount > 0 {
		fmt.Printf("   🏷️ %d개의 태그 삭제 중...", tagCount)
		// 모든 태그 삭제
		if tags := strings.Split(strings.TrimSpace(string(tagOutput)), "\n"); len(tags) > 0 {
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					delCmd := exec.Command("git", "tag", "-d", tag)
					delCmd.Run() // 에러 무시
				}
			}
		}
		fmt.Println(" 완료")
	}
	
	// 태그 fetch 차단 설정
	fmt.Print("   🚫 원격 태그 fetch 차단 설정...")
	blockTagCmd := exec.Command("git", "config", "remote.origin.tagOpt", "--no-tags")
	if err := blockTagCmd.Run(); err != nil {
		fmt.Printf(" 실패: %v\n", err)
	} else {
		fmt.Println(" 완료")
	}
	
	// 서브모듈도 태그 제거
	if _, err := os.Stat(".gitmodules"); err == nil {
		fmt.Println("   🔄 서브모듈 태그 제거 중...")
		submoduleNoTagCmd := exec.Command("git", "submodule", "foreach", "--recursive",
			"git tag -l | xargs -r git tag -d && git config remote.origin.tagOpt --no-tags")
		if err := submoduleNoTagCmd.Run(); err != nil {
			fmt.Printf("   ⚠️ 서브모듈 태그 제거 실패: %v\n", err)
		} else {
			fmt.Println("   ✅ 서브모듈 태그 제거 완료")
		}
	}

	// 6. 성능 설정 적용
	fmt.Println("\n6️⃣ Git 성능 최적화 설정...")
	performanceConfigs := [][]string{
		// 기존 최적화 설정
		{"core.commitGraph", "true"},
		{"gc.writeCommitGraph", "true"},
		{"fetch.writeCommitGraph", "true"},
		{"core.multiPackIndex", "true"},
		{"fetch.parallel", "4"},
		{"gc.autoDetach", "false"},
		
		// 회사 표준 설정 추가
		{"core.longpaths", "true"},          // 긴 경로 지원 (Windows)
		{"format.pretty", "oneline"},         // 로그 포맷
		{"color.ui", "true"},                 // 컬러 출력
		{"pull.rebase", "true"},              // pull 시 rebase 사용
		{"http.postBuffer", "2097152000"},    // HTTP 버퍼 크기 (2GB)
		{"pack.windowMemory", "256m"},        // 팩 메모리 크기 (줄여서 압축률 향상)
		{"pack.packSizeLimit", "512m"},       // 팩 파일 크기 제한 (작게 나눠서 관리)
		{"core.compression", "9"},            // 최대 압축 (.git 크기 최소화)
		{"pack.compression", "9"},            // 팩 파일도 최대 압축
		{"core.bigFileThreshold", "10m"},     // 10MB 이상은 delta 압축 제외
		{"core.untrackedCache", "true"},      // untracked 캐시 사용
		{"core.fsmonitor", "true"},           // 파일 시스템 모니터 사용
	}
	
	for _, config := range performanceConfigs {
		configCmd := exec.Command("git", "config", config[0], config[1])
		if err := configCmd.Run(); err != nil {
			fmt.Printf("⚠️ 설정 실패 (%s): %v\n", config[0], err)
		}
	}
	fmt.Println("   ✅ 성능 설정 완료")
	
	// 6. .gaconfig 디렉토리 및 설정 파일 생성
	fmt.Println("\n6️⃣ GA 설정 파일 생성...")
	if err := createGAConfig(); err != nil {
		fmt.Printf("⚠️ GA 설정 파일 생성 실패 (수동으로 생성 필요): %v\n", err)
	} else {
		fmt.Println("   ✅ .gaconfig/config.yaml 생성 완료")
	}
	
	// 7. 결과 확인
	fmt.Println("\n📊 최적화 결과:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Partial Clone 필터 확인
	filterCmd := exec.Command("git", "config", "remote.origin.partialclonefilter")
	if output, err := filterCmd.Output(); err == nil {
		fmt.Printf("Partial Clone: ✅ (필터: %s)\n", strings.TrimSpace(string(output)))
	}
	
	// Sparse Checkout 확인
	sparseCmd := exec.Command("git", "config", "core.sparseCheckout")
	if output, err := sparseCmd.Output(); err == nil && strings.TrimSpace(string(output)) == "true" {
		listCmd := exec.Command("git", "sparse-checkout", "list")
		if listOutput, err := listCmd.Output(); err == nil {
			paths := strings.Split(strings.TrimSpace(string(listOutput)), "\n")
			fmt.Printf("Sparse Checkout: ✅ (%d개 경로)\n", len(paths))
		}
	}
	
	// Shallow 상태 확인
	shallowCmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	if output, err := shallowCmd.Output(); err == nil && strings.TrimSpace(string(output)) == "true" {
		depthCmd := exec.Command("git", "rev-list", "--count", "HEAD")
		if depthOutput, err := depthCmd.Output(); err == nil {
			fmt.Printf("Shallow: ✅ (depth: %s)\n", strings.TrimSpace(string(depthOutput)))
		}
	}
	
	// 디스크 사용량 확인
	diskUsage := utils.GetDiskUsage()
	if gitSize, ok := diskUsage["git"]; ok {
		fmt.Printf(".git 폴더: %s\n", gitSize)
	}
	if totalSize, ok := diskUsage["total"]; ok {
		fmt.Printf("프로젝트 전체: %s\n", totalSize)
	}
	
	fmt.Println("\n✨ SLIM 모드 클론이 완료되었습니다!")
	fmt.Printf("📁 클론된 경로: %s\n", folder)
	fmt.Println("\n💡 팁: 필요한 경로를 추가하려면:")
	fmt.Println("   ga opt quick expand-slim")
}

// createGAConfig creates .gaconfig directory and config.yaml for SLIM mode
func createGAConfig() error {
	// .gaconfig 디렉토리 생성
	configDir := ".gaconfig"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf(".gaconfig 디렉토리 생성 실패: %w", err)
	}
	
	// prompt 디렉토리 생성
	promptDir := filepath.Join(configDir, "prompt")
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		return fmt.Errorf("prompt 디렉토리 생성 실패: %w", err)
	}
	
	// SLIM 모드에 최적화된 config.yaml 생성
	configContent := `# GA CLI 설정 파일 (SLIM 모드)

# AI 설정
ai:
  provider: "claude"  # 또는 "openai"
  openai:
    api_key: ""  # GA_AI_OPENAI_API_KEY 환경 변수로 설정 가능
    model: "gpt-4-turbo-preview"
  claude:
    api_key: ""  # GA_AI_CLAUDE_API_KEY 환경 변수로 설정 가능
    model: "claude-opus-4-20250514"

# 프롬프트 설정
prompt:
  analyze: "prompt/analyze.md"
  commit: "prompt/commit.md"

# Git 최적화 설정 (SLIM 모드)
optimize:
  mode: "slim"  # SLIM 모드로 설정
  filter:
    default: "1m"  # 1MB 필터 적용
    options:
      minimal: "1m"     # 소스코드만 (1MB 미만)
      basic: "25m"      # 코드 + 씬 파일
      extended: "50m"   # 대부분 리소스 포함
      full: "100m"      # 거의 전체
  sparse:
    paths: []  # Sparse Checkout 경로 (프로젝트별로 다르므로 비워둠)
`
	
	configFile := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("config.yaml 파일 생성 실패: %w", err)
	}
	
	// 기본 프롬프트 파일 생성
	analyzePrompt := `# Git Diff 분석 프롬프트

변경사항을 분석하고 다음 내용을 포함하여 설명해주세요:
1. 주요 변경사항 요약
2. 파일별 변경 내용
3. 잠재적 영향도
4. 개선 제안사항
`
	
	commitPrompt := `# Git Commit 메시지 생성 프롬프트

다음 규칙에 따라 커밋 메시지를 생성해주세요:
- 첫 줄: 50자 이내의 요약
- 빈 줄
- 상세 설명 (필요한 경우)
- 변경 이유와 영향도 포함
`
	
	if err := os.WriteFile(filepath.Join(promptDir, "analyze.md"), []byte(analyzePrompt), 0644); err != nil {
		return fmt.Errorf("analyze.md 파일 생성 실패: %w", err)
	}
	
	if err := os.WriteFile(filepath.Join(promptDir, "commit.md"), []byte(commitPrompt), 0644); err != nil {
		return fmt.Errorf("commit.md 파일 생성 실패: %w", err)
	}
	
	return nil
}