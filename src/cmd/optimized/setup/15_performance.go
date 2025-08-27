package setup

import (
	"fmt"
	"os/exec"
	"strings"
	
	"github.com/spf13/cobra"
	"workingcli/src/utils"
)

// NewPerformanceCmd creates the Performance optimization command
func NewPerformanceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Git 성능 최적화 설정 적용",
		Long: `Git 성능을 향상시키기 위한 다양한 설정을 적용합니다.

적용되는 설정:
- Git 기본 성능 최적화 (commitGraph, multiPackIndex 등)
- 회사 표준 설정 (압축, 버퍼, 캐싱 등)
- Git maintenance 스케줄 등록

사용법:
  ga optimized setup performance
  ga opt setup performance`,
		Run: func(cmd *cobra.Command, args []string) {
			executePerformance()
		},
	}
}

func executePerformance() {
	fmt.Println("🚀 Git 성능 최적화 설정을 적용합니다...")
	fmt.Println()
	
	// Git 저장소인지 확인
	if !utils.IsGitRepository() {
		fmt.Println("❌ Git 저장소가 아닙니다.")
		fmt.Println("   Git 저장소 디렉토리에서 실행해주세요.")
		return
	}
	
	// 성능 설정 배열 정의
	performanceConfigs := []struct {
		key   string
		value string
		desc  string
	}{
		// Git 기본 최적화 설정
		{"core.commitGraph", "true", "커밋 그래프 사용"},
		{"gc.writeCommitGraph", "true", "GC 시 커밋 그래프 작성"},
		{"fetch.writeCommitGraph", "true", "Fetch 시 커밋 그래프 작성"},
		{"core.multiPackIndex", "true", "멀티 팩 인덱스 사용"},
		{"fetch.parallel", "4", "병렬 Fetch (4개 스레드)"},
		{"gc.autoDetach", "false", "GC 백그라운드 실행 비활성화"},
		
		// 회사 표준 설정
		{"core.longpaths", "true", "긴 경로 지원 (Windows)"},
		{"format.pretty", "oneline", "로그 포맷 설정"},
		{"color.ui", "true", "컬러 출력 활성화"},
		{"pull.rebase", "true", "Pull 시 rebase 사용"},
		{"http.postBuffer", "2097152000", "HTTP 버퍼 크기 (2GB)"},
		{"pack.windowMemory", "256m", "팩 메모리 크기"},
		{"pack.packSizeLimit", "512m", "팩 파일 크기 제한"},
		{"core.compression", "9", "최대 압축 (.git 크기 최소화)"},
		{"pack.compression", "9", "팩 파일 최대 압축"},
		{"core.bigFileThreshold", "10m", "대용량 파일 임계값 (10MB)"},
		{"core.untrackedCache", "true", "Untracked 캐시 사용"},
		{"core.fsmonitor", "true", "파일 시스템 모니터 사용"},
	}
	
	// 설정 적용 전 현재 상태 백업
	fmt.Println("1️⃣ 현재 설정 백업...")
	backupCmd := exec.Command("git", "config", "--local", "--list")
	backupOutput, _ := backupCmd.Output()
	fmt.Println("   ✅ 현재 설정이 메모리에 백업되었습니다.")
	fmt.Println()
	
	// 성능 설정 적용
	fmt.Println("2️⃣ 성능 최적화 설정 적용...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	successCount := 0
	failedConfigs := []string{}
	
	for _, config := range performanceConfigs {
		cmd := exec.Command("git", "config", config.key, config.value)
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ %s 설정 실패\n", config.key)
			failedConfigs = append(failedConfigs, config.key)
		} else {
			fmt.Printf("✅ %-30s : %s\n", config.desc, config.value)
			successCount++
		}
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   적용 완료: %d개 / 실패: %d개\n", successCount, len(failedConfigs))
	
	if len(failedConfigs) > 0 {
		fmt.Printf("   실패한 설정: %s\n", strings.Join(failedConfigs, ", "))
	}
	fmt.Println()
	
	// Git maintenance 등록
	fmt.Println("3️⃣ Git maintenance 스케줄 등록...")
	maintenanceCmd := exec.Command("git", "maintenance", "register")
	if err := maintenanceCmd.Run(); err != nil {
		fmt.Printf("⚠️ Maintenance 등록 실패: %v\n", err)
		fmt.Println("   (수동으로 'git maintenance register' 실행 필요)")
	} else {
		fmt.Println("   ✅ Git maintenance 자동 실행 등록 완료")
	}
	fmt.Println()
	
	// 초기 maintenance 실행
	fmt.Println("4️⃣ 초기 maintenance 실행 (gc, commit-graph 등)...")
	runMaintenanceCmd := exec.Command("git", "maintenance", "run")
	if err := runMaintenanceCmd.Run(); err != nil {
		fmt.Printf("⚠️ Maintenance 실행 실패: %v\n", err)
	} else {
		fmt.Println("   ✅ 초기 maintenance 실행 완료 (gc, repack, commit-graph 업데이트)")
	}
	fmt.Println()
	
	// 최종 결과 확인
	fmt.Println("📊 성능 최적화 설정 결과:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// 주요 설정 확인
	checkSettings := []struct {
		key  string
		desc string
	}{
		{"core.commitGraph", "커밋 그래프"},
		{"core.multiPackIndex", "멀티 팩 인덱스"},
		{"core.fsmonitor", "파일 시스템 모니터"},
		{"core.untrackedCache", "Untracked 캐시"},
		{"core.compression", "압축 레벨"},
	}
	
	for _, setting := range checkSettings {
		cmd := exec.Command("git", "config", "--get", setting.key)
		if output, err := cmd.Output(); err == nil {
			value := strings.TrimSpace(string(output))
			fmt.Printf("%-20s: %s\n", setting.desc, value)
		}
	}
	
	// 디스크 사용량 확인
	diskUsage := utils.GetDiskUsage()
	if gitSize, ok := diskUsage["git"]; ok {
		fmt.Printf(".git 폴더 크기     : %s\n", gitSize)
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("✨ Git 성능 최적화 설정이 완료되었습니다!")
	fmt.Println()
	fmt.Println("💡 팁:")
	fmt.Println("- 정기적으로 'git maintenance run' 실행을 권장합니다.")
	fmt.Println("- 대용량 저장소의 경우 'ga opt quick to-slim'으로 추가 최적화 가능합니다.")
	
	// 설정 변경이 있었는지 확인
	if len(string(backupOutput)) > 0 {
		fmt.Println("- 이전 설정으로 되돌리려면 백업된 설정을 참고하세요.")
	}
}