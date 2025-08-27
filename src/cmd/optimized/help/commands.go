package help

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewCommandsCmd creates the commands list command
func NewCommandsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
        Aliases: []string{"help", "-h"},  // ga sub로도 사용 가능
		Short: "전체 최적화 명령어 목록",
		Long: `Git 저장소 최적화를 위한 명령어 체계를 표시합니다.
각 명령어를 카테고리와 용도별로 구성하여 제공합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			printCommandsList()
		},
	}
}

func printCommandsList() {
	fmt.Println("📚 Git 저장소 최적화 명령어 체계")
	fmt.Println("=====================================")
	
	// Help 카테고리
	fmt.Println("\n📖 Help (도움말 및 가이드)")
	fmt.Println("  • workflow     - 최적화 워크플로우 가이드")
	fmt.Println("  • commands     - 전체 명령어 목록 (현재 명령)")
	
	// Quick 카테고리
	fmt.Println("\n⚡ Quick (자주 사용하는 기능)")
	fmt.Println("  [모드 전환]")
	fmt.Println("  • status       - 현재 최적화 상태 확인")
	fmt.Println("  • to-slim      - SLIM 모드로 전환 (103GB → 30MB)")
	fmt.Println("  • to-full      - FULL 모드로 복원 (전체 저장소)")
	fmt.Println("  \n  [확장 및 필터]")
	fmt.Println("  • expand-slim  - 선택적 경로 확장")
	fmt.Println("  • expand-filter- Partial Clone 필터 제거")
	fmt.Println("  • filter-branch- 브랜치별 필터 설정")
	fmt.Println("  • clear-filter - 필터 완전 제거")
	fmt.Println("  \n  [히스토리 관리]")
	fmt.Println("  • shallow      - 히스토리 줄이기 (depth=1)")
	fmt.Println("  • unshallow    - 히스토리 복원")
	fmt.Println("  • auto-find-merge-base - 병합 베이스 자동 찾기")
	
	// Setup 카테고리
	fmt.Println("\n🛠️ Setup (초기 설정 및 마이그레이션)")
	fmt.Println("  • clone-slim   - 최적화된 클론 (처음부터 30MB)")
	fmt.Println("  • migrate      - 기존 저장소 마이그레이션")
	fmt.Println("  • performance  - 성능 최적화 설정")
	
	// Workspace 카테고리
	fmt.Println("\n💼 Workspace (작업 공간 관리)")
	fmt.Println("  • expand-path  - 특정 경로 확장")
	fmt.Println("  • restore-branch- 브랜치 전체 복원")
	
	// Advanced 카테고리
	fmt.Println("\n🔧 Advanced (고급 최적화 기능)")
	fmt.Println("  [히스토리 확장]")
	fmt.Println("  • expand-10    - 히스토리 10개 확장")
	fmt.Println("  • expand-50    - 히스토리 50개 확장")
	fmt.Println("  • expand-100   - 히스토리 100개 확장")
	fmt.Println("  \n  [상태 확인]")
	fmt.Println("  • check-merge-base - 병합 베이스 확인")
	fmt.Println("  • check-shallow    - 히스토리 상태 확인")
	fmt.Println("  • check-filter     - 브랜치 필터 확인")
	fmt.Println("  \n  [설정 관리]")
	fmt.Println("  • config       - 설정 백업/복원")
	
	// Submodule 카테고리
	fmt.Println("\n📦 Submodule (서브모듈 최적화)")
	fmt.Println("  • shallow      - 모든 서브모듈 shallow 변환")
	fmt.Println("  • unshallow    - 모든 서브모듈 히스토리 복원")
	
	fmt.Println("\n💡 사용 예시:")
	fmt.Println("  ga opt quick status              # 현재 상태 확인")
	fmt.Println("  ga opt quick to-slim             # SLIM 모드로 전환")
	fmt.Println("  ga opt advanced expand-50        # 히스토리 50개 확장")
	fmt.Println("  ga opt setup clone-slim <url>    # 최적화 클론")
	fmt.Println("\n📌 자세한 사용법: 'ga optimized help workflow'")
}