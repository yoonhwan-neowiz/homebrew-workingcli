package setup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewCloneMasterCmd creates the Clone Master command
func NewCloneMasterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone-master [URL] [folder]",
		Aliases: []string{"clone"},
		Short: "Master 브랜치만 최적화 모드로 클론",
		Long: `Master 브랜치만 빠르게 클론합니다.
브랜치 스코프를 master로 제한하고 shallow depth=1을 적용합니다.

사용법:
  ga optimized setup clone-master <URL> [folder]
  
예시:
  ga opt setup clone-master https://github.com/user/repo.git myproject
  ga opt setup clone-master https://github.com/user/repo.git`,
		Run: func(cmd *cobra.Command, args []string) {
			// args 수동 검증
			if len(args) < 1 {
				fmt.Println("❌ URL을 입력해주세요.")
				fmt.Println("\n사용법:")
				fmt.Println("  ga optimized setup clone-master <URL> [folder]")
				return
			}
			executeCloneMaster(args)
		},
	}
}

func executeCloneMaster(args []string) {
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
	
	fmt.Printf("🚀 Master 브랜치만 '%s' 저장소를 클론합니다...\n", url)
	fmt.Printf("   타겟 폴더: %s\n", folder)
	fmt.Println()
	
	// 1. Shallow Clone으로 master 브랜치만 클론 (no-checkout)
	fmt.Println("1️⃣ Master 브랜치만 Shallow Clone (no-checkout)...")
	cloneCmd := exec.Command("git", "clone", 
		"--depth=1",
		"--single-branch",  // master 브랜치만 가져옴
		"--branch", "master",
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
	
	// 2. 브랜치 스코프를 master로 제한
	fmt.Println("\n2️⃣ 브랜치 스코프를 master로 제한...")
	
	// fetch 설정을 master 브랜치만으로 제한
	fetchConfigCmd := exec.Command("git", "config", 
		"remote.origin.fetch", 
		"+refs/heads/master:refs/remotes/origin/master")
	if err := fetchConfigCmd.Run(); err != nil {
		fmt.Printf("⚠️ 브랜치 스코프 설정 실패: %v\n", err)
	} else {
		fmt.Println("   ✅ 브랜치 스코프 설정 완료 (master만 fetch)")
	}
	
	// 3. 체크아웃 수행
	fmt.Println("\n3️⃣ 파일 체크아웃...")
	checkoutCmd := exec.Command("git", "checkout", "master")
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		// master가 없으면 HEAD로 체크아웃
		checkoutCmd = exec.Command("git", "checkout", "HEAD")
		if err := checkoutCmd.Run(); err != nil {
			fmt.Printf("❌ 체크아웃 실패: %v\n", err)
			fmt.Printf("   출력: %s\n", string(output))
			return
		}
	}
	fmt.Println("   ✅ 파일 체크아웃 완료")
	
	// 4. 서브모듈 초기화 및 최적화
	fmt.Println("\n4️⃣ 서브모듈 초기화 및 최적화...")
	
	// 한 줄로 서브모듈 초기화 + shallow clone + single-branch 설정
	updateCmd := exec.Command("git", "submodule", "update", 
		"--init",           // 초기화
		"--recursive",      // 재귀적으로 모든 서브모듈
		"--depth=1",        // shallow clone
		"--single-branch")  // master 브랜치만
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	
	if err := updateCmd.Run(); err != nil {
		// 서브모듈이 없거나 실패한 경우
		if strings.Contains(err.Error(), "No submodule mapping") {
			fmt.Println("   ℹ️ 서브모듈이 없습니다.")
		} else {
			fmt.Printf("⚠️ 서브모듈 업데이트 실패: %v\n", err)
		}
	} else {
		fmt.Println("   ✅ 서브모듈 최적화 완료 (master only, shallow, single-branch)")
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
		// 기본 최적화 설정
		{"core.commitGraph", "true"},
		{"gc.writeCommitGraph", "true"},
		{"fetch.writeCommitGraph", "true"},
		{"core.multiPackIndex", "true"},
		{"fetch.parallel", "4"},
		{"gc.autoDetach", "false"},
		
		// 추가 성능 설정
		{"core.longpaths", "true"},          // 긴 경로 지원 (Windows)
		{"color.ui", "true"},                 // 컬러 출력
		{"pull.rebase", "true"},              // pull 시 rebase 사용
		{"http.postBuffer", "524288000"},     // HTTP 버퍼 크기 (500MB)
	}
	
	for _, config := range performanceConfigs {
		configCmd := exec.Command("git", "config", config[0], config[1])
		if err := configCmd.Run(); err != nil {
			fmt.Printf("⚠️ 설정 실패 (%s): %v\n", config[0], err)
		}
	}
	fmt.Println("   ✅ 성능 설정 완료")
	
	// 6. 결과 확인
	fmt.Println("\n📊 최적화 결과:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 브랜치 스코프 확인
	fetchCmd := exec.Command("git", "config", "remote.origin.fetch")
	if output, err := fetchCmd.Output(); err == nil {
		fmt.Printf("브랜치 스코프: ✅ (%s)\n", strings.TrimSpace(string(output)))
	}
	
	// Shallow 상태 확인
	shallowCmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	if output, err := shallowCmd.Output(); err == nil && strings.TrimSpace(string(output)) == "true" {
		depthCmd := exec.Command("git", "rev-list", "--count", "HEAD")
		if depthOutput, err := depthCmd.Output(); err == nil {
			fmt.Printf("Shallow: ✅ (depth: %s)\n", strings.TrimSpace(string(depthOutput)))
		}
	}
	
	// 현재 브랜치 확인
	branchCmd := exec.Command("git", "branch", "--show-current")
	if output, err := branchCmd.Output(); err == nil {
		fmt.Printf("현재 브랜치: %s\n", strings.TrimSpace(string(output)))
	}
	
	// 리모트 브랜치 확인
	remoteBranchCmd := exec.Command("git", "branch", "-r")
	if output, err := remoteBranchCmd.Output(); err == nil {
		remoteBranches := strings.Split(strings.TrimSpace(string(output)), "\n")
		fmt.Printf("리모트 브랜치: %d개\n", len(remoteBranches))
		for _, branch := range remoteBranches {
			fmt.Printf("  %s\n", strings.TrimSpace(branch))
		}
	}
	
	// 서브모듈 전체 reset --hard 수행 (최종 정리)
	if _, err := os.Stat(".gitmodules"); err == nil {
		fmt.Println("\n🔄 서브모듈 전체 reset --hard 수행...")
		resetCmd := exec.Command("git", "submodule", "foreach", "--recursive", 
			"git", "reset", "--hard")
		resetCmd.Stdout = os.Stdout
		resetCmd.Stderr = os.Stderr
		
		if err := resetCmd.Run(); err != nil {
			fmt.Printf("⚠️ 서브모듈 reset 실패: %v\n", err)
		} else {
			fmt.Println("✅ 서브모듈 reset 완료")
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
	
	fmt.Println("\n✨ Master 브랜치 클론이 완료되었습니다!")
	fmt.Printf("📁 클론된 경로: %s\n", folder)
	fmt.Println("\n💡 팁: 다른 브랜치가 필요한 경우:")
	fmt.Println("   git fetch origin <branch-name>:<branch-name>")
	fmt.Println("   git checkout <branch-name>")
}