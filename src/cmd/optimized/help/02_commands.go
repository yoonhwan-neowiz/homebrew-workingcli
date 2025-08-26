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
		Long: `Git 저장소 최적화를 위한 28개 명령어의 전체 목록을 표시합니다.
각 명령어의 용도와 사용 예시를 카테고리별로 정리하여 제공합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			printCommandsList()
		},
	}
}

func printCommandsList() {
	fmt.Println("📚 Git 저장소 최적화 명령어 목록 (28개)")
	fmt.Println("=====================================")
	
	// Help 카테고리
	fmt.Println("\n📖 Help (도움말)")
	fmt.Println("  1. workflow     - 최적화 워크플로우 가이드")
	fmt.Println("  2. commands     - 전체 명령어 목록 (현재 명령)")
	
	// Quick 카테고리
	fmt.Println("\n⚡ Quick (자주 사용)")
	fmt.Println("  3. status       - 현재 최적화 상태 확인")
	fmt.Println("  4. to-slim      - SLIM 모드로 전환 (103GB → 30MB)")
	fmt.Println("  5. to-full      - FULL 모드로 복원 (전체 저장소)")
	fmt.Println("  6. expand-slim  - 선택적 경로 확장")
	fmt.Println("  7. expand-filter- Partial Clone 필터 제거")
	fmt.Println("  8. expand       - 히스토리 확장 (기본 10개, 파라미터로 지정 가능)")
	fmt.Println("  9. expand-50    - (deprecated - expand 50 사용)")
	fmt.Println(" 10. expand-100   - (deprecated - expand 100 사용)")
	fmt.Println(" 11. find-merge   - 병합 베이스 찾기")
	fmt.Println(" 12. check-merge  - 병합 가능 여부 확인")
	
	// Setup 카테고리
	fmt.Println("\n🛠️ Setup (초기 설정)")
	fmt.Println(" 13. clone-slim   - 최적화된 클론 (처음부터 30MB)")
	fmt.Println(" 14. migrate      - (deprecated - to-slim 사용)")
	fmt.Println(" 15. performance  - 성능 최적화 설정")
	
	// Workspace 카테고리
	fmt.Println("\n💼 Workspace (작업 공간)")
	fmt.Println(" 16. expand-path  - 특정 경로 확장")
	fmt.Println(" 17. filter-branch- 브랜치별 필터 설정")
	fmt.Println(" 18. clear-filter - 필터 완전 제거")
	fmt.Println(" 19. restore-branch- 브랜치 전체 복원")
	
	// Advanced 카테고리
	fmt.Println("\n🔧 Advanced (고급 기능)")
	fmt.Println(" 20. shallow      - 히스토리 줄이기 (depth=1)")
	fmt.Println(" 21. unshallow    - 히스토리 복원")
	fmt.Println(" 22. check-shallow- 히스토리 상태 확인")
	fmt.Println(" 23. check-filter - 브랜치 필터 확인")
	fmt.Println(" 24. backup-config- 설정 백업/복원")
	
	// Submodule 카테고리
	fmt.Println("\n📦 Submodule (서브모듈)")
	fmt.Println(" 25. shallow-all  - 모든 서브모듈 shallow 변환")
	fmt.Println(" 26. unshallow-all- 모든 서브모듈 히스토리 복원")
	fmt.Println(" 27. optimize-all - 모든 서브모듈 SLIM 최적화")
	fmt.Println(" 28. status-all   - 모든 서브모듈 상태 확인")
	
	fmt.Println("\n💡 사용 예시:")
	fmt.Println("  ga optimized quick status        # 현재 상태 확인")
	fmt.Println("  ga optimized quick to-slim       # SLIM 모드로 전환")
	fmt.Println("  ga optimized quick expand 50     # 히스토리 50개 확장")
	fmt.Println("  ga optimized setup clone-slim <url> <folder>  # 최적화 클론")
	fmt.Println("\n자세한 사용법은 'ga optimized help workflow'를 참조하세요.")
}