package advanced

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"workingcli/src/config"
	"workingcli/src/utils"
)

// NewConfigCmd creates the Config command
func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "최적화 설정 관리",
		Long: `최적화 설정을 백업, 복원, 확인합니다.
.gaconfig 디렉토리에 설정과 Git 최적화 정보를 저장합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			runConfig()
		},
	}
}

func runConfig() {
	// 색상 설정
	titleStyle := color.New(color.FgCyan, color.Bold)
	errorStyle := color.New(color.FgRed)
	
	titleStyle.Println("\n⚙️  최적화 설정 관리")
	titleStyle.Println("=" + strings.Repeat("=", 39))
	
	// 1. Git 저장소 확인
	if !utils.IsGitRepository() {
		errorStyle.Println("❌ Git 저장소가 아닙니다.")
		os.Exit(1)
	}
	
	// 2. 작업 선택
	fmt.Println("\n📋 작업 선택:")
	fmt.Println("   1. 현재 설정 백업")
	fmt.Println("   2. 백업에서 복원")
	fmt.Println("   3. 백업 목록 확인")
	fmt.Println("   4. 현재 설정 확인")
	fmt.Print("\n선택 (1-4): ")
	
	var choice string
	fmt.Scanln(&choice)
	
	switch choice {
	case "1":
		performBackup()
	case "2":
		performRestore()
	case "3":
		listBackups()
	case "4":
		showCurrentConfig()
	default:
		errorStyle.Println("❌ 잘못된 선택입니다.")
		os.Exit(1)
	}
}

func performBackup() {
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	boldStyle := color.New(color.Bold)
	
	fmt.Println("\n📦 현재 설정 백업")
	fmt.Println("─" + strings.Repeat("─", 39))
	
	// 백업 디렉토리 생성 (.gaconfig/backups)
	backupDir := filepath.Join(".gaconfig", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		errorStyle := color.New(color.FgRed)
		errorStyle.Printf("❌ 백업 디렉토리 생성 실패: %v\n", err)
		os.Exit(1)
	}
	
	// 타임스탬프로 백업 디렉토리 생성
	timestamp := time.Now().Format("20060102-150405")
	timestampDir := filepath.Join(backupDir, timestamp)
	if err := os.MkdirAll(timestampDir, 0755); err != nil {
		errorStyle := color.New(color.FgRed)
		errorStyle.Printf("❌ 타임스탬프 디렉토리 생성 실패: %v\n", err)
		os.Exit(1)
	}
	
	// 1. config.yaml 백업
	fmt.Print("📝 config.yaml 백업 중... ")
	configSource := filepath.Join(".gaconfig", "config.yaml")
	configBackup := filepath.Join(timestampDir, "config.yaml")
	
	if data, err := os.ReadFile(configSource); err != nil {
		warningStyle.Println("실패 (파일 없음)")
	} else {
		if err := os.WriteFile(configBackup, data, 0644); err != nil {
			warningStyle.Println("저장 실패")
		} else {
			infoStyle.Println("완료")
		}
	}
	
	// 2. Sparse Checkout 목록 백업
	fmt.Print("📂 Sparse Checkout 목록 백업 중... ")
	sparseBackup := filepath.Join(timestampDir, "sparse-checkout.txt")
	
	if utils.IsSparseCheckoutEnabled() {
		sparseInfo := utils.GetSparseCheckoutInfo()
		if paths, ok := sparseInfo["paths"].([]string); ok {
			sparseContent := strings.Join(paths, "\n")
			if err := os.WriteFile(sparseBackup, []byte(sparseContent), 0644); err != nil {
				warningStyle.Println("저장 실패")
			} else {
				infoStyle.Println("완료")
			}
		} else {
			warningStyle.Println("경로 없음")
		}
	} else {
		infoStyle.Println("건너뜀 (비활성)")
	}
	
	// 3. Git 최적화 설정 백업
	fmt.Print("🔧 Git 최적화 설정 백업 중... ")
	optimizationBackup := filepath.Join(timestampDir, "git-optimization.txt")
	
	var configLines []string
	
	// Partial Clone 필터
	if filter := utils.GetPartialCloneFilter(); filter != "" {
		configLines = append(configLines, fmt.Sprintf("partial-clone-filter=%s", filter))
	}
	
	// Shallow 상태
	shallowInfo := utils.GetShallowInfo()
	if isShallow := shallowInfo["isShallow"].(bool); isShallow {
		if depth, ok := shallowInfo["depth"].(int); ok {
			configLines = append(configLines, fmt.Sprintf("shallow-depth=%d", depth))
		}
	}
	
	// 브랜치 필터
	if branchFilter := utils.GetBranchFilter(); len(branchFilter) > 0 {
		configLines = append(configLines, fmt.Sprintf("branch-filter=%s", strings.Join(branchFilter, ",")))
	}
	
	// Sparse Checkout 상태
	if utils.IsSparseCheckoutEnabled() {
		configLines = append(configLines, "sparse-checkout=enabled")
		if utils.CheckConeMode() {
			configLines = append(configLines, "sparse-checkout-mode=cone")
		} else {
			configLines = append(configLines, "sparse-checkout-mode=non-cone")
		}
	}
	
	// 현재 모드
	mode := utils.GetOptimizationMode()
	configLines = append(configLines, fmt.Sprintf("mode=%s", mode))
	
	// 디스크 사용량
	diskUsage := utils.GetDiskUsage()
	if gitSize, ok := diskUsage["git"]; ok {
		configLines = append(configLines, fmt.Sprintf("git-size=%s", gitSize))
	}
	
	configContent := strings.Join(configLines, "\n")
	if err := os.WriteFile(optimizationBackup, []byte(configContent), 0644); err != nil {
		warningStyle.Println("저장 실패")
	} else {
		infoStyle.Println("완료")
	}
	
	// 3-1. 서브모듈 설정 백업
	fmt.Print("📦 서브모듈 설정 백업 중... ")
	submoduleBackup := filepath.Join(timestampDir, "submodule-settings.txt")
	
	var submoduleLines []string
	// config에서 서브모듈 설정 읽기
	settings := config.GetAll()
	if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
		if submodule, ok := optimize["submodule"].(map[string]interface{}); ok {
			if mode, ok := submodule["mode"].(string); ok {
				submoduleLines = append(submoduleLines, fmt.Sprintf("mode=%s", mode))
			}
			if filter, ok := submodule["filter"].(map[string]interface{}); ok {
				if defaultFilter, ok := filter["default"].(string); ok {
					submoduleLines = append(submoduleLines, fmt.Sprintf("filter.default=%s", defaultFilter))
				}
			}
			if sparse, ok := submodule["sparse"].(map[string]interface{}); ok {
				if paths, ok := sparse["paths"].([]interface{}); ok && len(paths) > 0 {
					var pathStrs []string
					for _, p := range paths {
						if ps, ok := p.(string); ok {
							pathStrs = append(pathStrs, ps)
						}
					}
					submoduleLines = append(submoduleLines, fmt.Sprintf("sparse.paths=%s", strings.Join(pathStrs, ",")))
				} else {
					submoduleLines = append(submoduleLines, "sparse.paths=")
				}
			}
		}
	}
	
	if len(submoduleLines) > 0 {
		content := strings.Join(submoduleLines, "\n")
		if err := os.WriteFile(submoduleBackup, []byte(content), 0644); err != nil {
			warningStyle.Println("저장 실패")
		} else {
			infoStyle.Println("완료")
		}
	} else {
		infoStyle.Println("건너뜀 (설정 없음)")
	}
	
	// 4. 백업 요약
	fmt.Println("\n✅ 백업 완료!")
	fmt.Printf("   ├─ 위치: %s\n", boldStyle.Sprint(backupDir))
	fmt.Printf("   ├─ 타임스탬프: %s\n", boldStyle.Sprint(timestamp))
	fmt.Println("   └─ 파일:")
	fmt.Printf("       • %s\n", filepath.Base(configBackup))
	if utils.IsSparseCheckoutEnabled() {
		fmt.Printf("       • %s\n", filepath.Base(sparseBackup))
	}
	fmt.Printf("       • %s\n", filepath.Base(optimizationBackup))
}

func performRestore() {
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	
	fmt.Println("\n♻️  백업에서 복원")
	fmt.Println("─" + strings.Repeat("─", 39))
	
	// 백업 목록 확인
	backupDir := ".gaconfig/backups"
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		errorStyle.Println("❌ 백업 디렉토리를 읽을 수 없습니다.")
		return
	}
	
	// 타임스탬프 추출 (디렉토리 기반)
	var timestamps []string
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if len(name) == 15 { // YYYYMMDD-HHMMSS 형식
				timestamps = append(timestamps, name)
			}
		}
	}
	
	if len(timestamps) == 0 {
		warningStyle.Println("⚠️  백업이 없습니다.")
		return
	}
	
	// 정렬 (최신 먼저)
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[j] > timestamps[i] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}
	
	// 백업 선택
	fmt.Println("\n📋 백업 목록:")
	for i, ts := range timestamps {
		fmt.Printf("   %d. %s\n", i+1, ts)
		if i >= 4 { // 최대 5개만 표시
			break
		}
	}
	
	fmt.Print("\n복원할 백업 번호: ")
	var choice int
	fmt.Scanln(&choice)
	
	if choice < 1 || choice > len(timestamps) {
		errorStyle.Println("❌ 잘못된 선택입니다.")
		return
	}
	
	selectedTimestamp := timestamps[choice-1]
	backupTimestampDir := filepath.Join(backupDir, selectedTimestamp)
	
	warningStyle.Println("\n⚠️  현재 설정이 백업 설정으로 교체됩니다.")
	if !utils.Confirm("계속하시겠습니까?") {
		fmt.Println("취소되었습니다.")
		return
	}
	
	// 1. config.yaml 복원
	configBackup := filepath.Join(backupTimestampDir, "config.yaml")
	if data, err := os.ReadFile(configBackup); err == nil {
		fmt.Print("📝 config.yaml 복원 중... ")
		
		configDest := ".gaconfig/config.yaml"
		if err := os.WriteFile(configDest, data, 0644); err != nil {
			warningStyle.Println("실패")
		} else {
			infoStyle.Println("완료")
			// 설정은 자동 로드됨
		}
	}
	
	// 2. Sparse Checkout 복원
	sparseBackup := filepath.Join(backupTimestampDir, "sparse-checkout.txt")
	if data, err := os.ReadFile(sparseBackup); err == nil {
		fmt.Print("📂 Sparse Checkout 목록 복원 중... ")
		
		// Sparse checkout 초기화
		cmd := exec.Command("git", "sparse-checkout", "init", "--cone")
		cmd.Run()
		
		// 경로 복원
		paths := strings.Split(strings.TrimSpace(string(data)), "\n")
		if len(paths) > 0 && paths[0] != "" {
			args := append([]string{"sparse-checkout", "set"}, paths...)
			cmd = exec.Command("git", args...)
			if err := cmd.Run(); err != nil {
				warningStyle.Println("일부 실패")
			} else {
				infoStyle.Println("완료")
			}
		}
	}
	
	// 3. Git 최적화 설정 복원
	optimizationBackup := filepath.Join(backupTimestampDir, "git-optimization.txt")
	if data, err := os.ReadFile(optimizationBackup); err == nil {
		fmt.Print("🔧 Git 최적화 설정 복원 중... ")
		
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				
				switch key {
				case "partial-clone-filter":
					cmd := exec.Command("git", "config", "remote.origin.partialclonefilter", value)
					cmd.Run()
					cmd = exec.Command("git", "config", "remote.origin.promisor", "true")
					cmd.Run()
				case "shallow-depth":
					// Shallow 복원은 별도 명령어로 처리
					fmt.Printf("\n   └─ Shallow depth %s 복원은 수동으로 실행: ga opt advanced shallow %s\n", value, value)
				case "branch-filter":
					// 브랜치 필터는 config.yaml에 저장됨
					fmt.Printf("\n   └─ 브랜치 필터 복원은 config.yaml을 통해 처리됨\n")
				}
			}
		}
		infoStyle.Println("완료")
	}
	
	// 3-1. 서브모듈 설정 복원 (있다면 config에 반영)
	submoduleBackup := filepath.Join(backupTimestampDir, "submodule-settings.txt")
	if data, err := os.ReadFile(submoduleBackup); err == nil {
		fmt.Print("📦 서브모듈 설정 복원 중... ")
		
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				
				switch key {
				case "mode":
					config.Set("optimize.submodule.mode", value)
				case "filter.default":
					config.Set("optimize.submodule.filter.default", value)
				case "sparse.paths":
					if value != "" {
						paths := strings.Split(value, ",")
						config.Set("optimize.submodule.sparse.paths", paths)
					}
				}
			}
		}
		infoStyle.Println("완료")
	}
	
	fmt.Printf("\n✅ 백업 %s에서 복원 완료!\n", boldStyle.Sprint(selectedTimestamp))
	fmt.Println("\n💡 상태 확인 명령어:")
	fmt.Println("   • 최적화 상태: ga opt quick status")
	fmt.Println("   • 필터 확인: ga opt advanced check-filter")
	fmt.Println("   • Shallow 상태: ga opt advanced check-shallow")
}

func listBackups() {
	warningStyle := color.New(color.FgYellow)
	errorStyle := color.New(color.FgRed)
	boldStyle := color.New(color.Bold)
	dimStyle := color.New(color.Faint)
	
	fmt.Println("\n📚 백업 목록")
	fmt.Println("─" + strings.Repeat("─", 39))
	
	backupDir := ".gaconfig/backups"
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		errorStyle.Println("❌ 백업 디렉토리를 읽을 수 없습니다.")
		return
	}
	
	// 타임스탬프별로 그룹화 (디렉토리 기반)
	backups := make(map[string][]string)
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 15 {
			timestamp := entry.Name()
			// 해당 디렉토리의 파일 목록 확인
			timestampDir := filepath.Join(backupDir, timestamp)
			if files, err := os.ReadDir(timestampDir); err == nil {
				for _, file := range files {
					if !file.IsDir() {
						backups[timestamp] = append(backups[timestamp], file.Name())
					}
				}
			}
		}
	}
	
	if len(backups) == 0 {
		warningStyle.Println("⚠️  백업이 없습니다.")
		return
	}
	
	// 타임스탬프 정렬
	var timestamps []string
	for ts := range backups {
		timestamps = append(timestamps, ts)
	}
	
	// 최신 순으로 정렬
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[j] > timestamps[i] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}
	
	// 백업 목록 표시
	for i, ts := range timestamps {
		fmt.Printf("\n%d. %s\n", i+1, boldStyle.Sprint(ts))
		
		for _, file := range backups[ts] {
			filePath := filepath.Join(backupDir, ts, file)
			info, err := os.Stat(filePath)
			if err == nil {
				size := info.Size()
				sizeStr := fmt.Sprintf("%d bytes", size)
				if size > 1024 {
					sizeStr = fmt.Sprintf("%.1f KB", float64(size)/1024)
				}
				
				fmt.Printf("   • %s %s\n", 
					file,
					dimStyle.Sprintf("(%s)", sizeStr))
			}
		}
		
		if i >= 9 { // 최대 10개만 표시
			if len(timestamps) > 10 {
				fmt.Printf("\n   ... 외 %d개 백업\n", len(timestamps)-10)
			}
			break
		}
	}
	
	// 디스크 사용량
	var totalSize int64
	for ts, files := range backups {
		for _, file := range files {
			filePath := filepath.Join(backupDir, ts, file)
			if info, err := os.Stat(filePath); err == nil {
				totalSize += info.Size()
			}
		}
	}
	
	fmt.Printf("\n💾 총 백업 크기: %s\n", boldStyle.Sprint(utils.FormatSize(totalSize)))
	fmt.Printf("📁 백업 위치: %s\n", dimStyle.Sprint(backupDir))
}

func showCurrentConfig() {
	infoStyle := color.New(color.FgGreen)
	warningStyle := color.New(color.FgYellow)
	boldStyle := color.New(color.Bold)
	dimStyle := color.New(color.Faint)
	
	fmt.Println("\n📋 현재 설정 상태")
	fmt.Println("─" + strings.Repeat("─", 39))
	
	// 1. config.yaml 상태
	fmt.Println("\n📄 Config 파일:")
	configFile := ".gaconfig/config.yaml"
	if _, err := os.Stat(configFile); err == nil {
		infoStyle.Printf("   └─ 상태: 존재함 (%s)\n", configFile)
		
		// 설정 내용 표시
		currentCfg := config.Get()
		
		if currentCfg == nil {
			warningStyle.Println("   └─ 설정 로드 실패")
			return
		}
		
		fmt.Println("\n⚙️  설정 내용:")
		if currentCfg.Optimize.Mode != "" {
			fmt.Printf("   ├─ 모드: %s\n", boldStyle.Sprint(currentCfg.Optimize.Mode))
		}
		
		// 브랜치 필터는 utils에서 가져옴
		if branchFilter := utils.GetBranchFilter(); len(branchFilter) > 0 {
			fmt.Printf("   ├─ 브랜치 필터: %s\n", boldStyle.Sprint(strings.Join(branchFilter, ", ")))
		}
		
		if len(currentCfg.Optimize.Sparse.Paths) > 0 {
			fmt.Printf("   └─ Sparse 경로: %d개\n", len(currentCfg.Optimize.Sparse.Paths))
			for i, path := range currentCfg.Optimize.Sparse.Paths {
				if i < 5 {
					fmt.Printf("       • %s\n", path)
				}
			}
			if len(currentCfg.Optimize.Sparse.Paths) > 5 {
				fmt.Printf("       • ... 외 %d개\n", len(currentCfg.Optimize.Sparse.Paths)-5)
			}
		}
	} else {
		warningStyle.Println("   └─ 상태: 없음")
	}
	
	// 2. Git 최적화 상태
	fmt.Println("\n🔧 Git 최적화:")
	
	// Partial Clone
	partialFilter := utils.GetPartialCloneFilter()
	if partialFilter != "" {
		fmt.Printf("   ├─ Partial Clone: %s\n", warningStyle.Sprint(partialFilter))
	} else {
		fmt.Printf("   ├─ Partial Clone: %s\n", dimStyle.Sprint("비활성"))
	}
	
	// Shallow 상태
	shallowInfo := utils.GetShallowInfo()
	if isShallow := shallowInfo["isShallow"].(bool); isShallow {
		if depth, ok := shallowInfo["depth"].(int); ok {
			fmt.Printf("   ├─ Shallow: %s (depth=%d)\n", warningStyle.Sprint("활성"), depth)
		}
	} else {
		fmt.Printf("   ├─ Shallow: %s\n", dimStyle.Sprint("비활성"))
	}
	
	// Sparse Checkout
	if utils.IsSparseCheckoutEnabled() {
		sparseInfo := utils.GetSparseCheckoutInfo()
		if count, ok := sparseInfo["count"].(int); ok {
			fmt.Printf("   └─ Sparse Checkout: %s (%d개 경로)\n", warningStyle.Sprint("활성"), count)
		}
	} else {
		fmt.Printf("   └─ Sparse Checkout: %s\n", dimStyle.Sprint("비활성"))
	}
	
	// 3. 백업 정보
	fmt.Println("\n💾 백업 정보:")
	backupDir := ".gaconfig/backups"
	if entries, err := os.ReadDir(backupDir); err == nil {
		backupCount := 0
		var latestBackup string
		for _, entry := range entries {
			if entry.IsDir() && len(entry.Name()) == 15 {
				backupCount++
				if entry.Name() > latestBackup {
					latestBackup = entry.Name()
				}
			}
		}
		
		if backupCount > 0 {
			fmt.Printf("   ├─ 백업 개수: %s개\n", boldStyle.Sprint(backupCount))
			fmt.Printf("   └─ 최근 백업: %s\n", boldStyle.Sprint(latestBackup))
		} else {
			dimStyle.Println("   └─ 백업 없음")
		}
	} else {
		dimStyle.Println("   └─ 백업 디렉토리 없음")
	}
	
	fmt.Println()
}