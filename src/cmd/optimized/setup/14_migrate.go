package setup

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewMigrateCmd creates the Migrate to SLIM command (DEPRECATED)
func NewMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "(DEPRECATED) 기존 저장소를 SLIM 모드로 변환 - 대신 'to-slim' 사용",
		Long: `⚠️ DEPRECATED: 이 명령어는 더 이상 사용되지 않습니다.
대신 'ga optimized quick to-slim' 명령어를 사용하세요.

migrate와 to-slim은 동일한 기능을 수행합니다:
- 기존 FULL 상태의 저장소를 SLIM 모드로 변환
- 작업 내용 보존하면서 최적화 적용

권장 사용법:
  ga optimized quick to-slim     # SLIM 모드로 전환
  ga opt quick to-slim           # 짧은 별칭`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚠️  DEPRECATED: 'migrate' 명령어는 더 이상 사용되지 않습니다.")
			fmt.Println()
			fmt.Println("대신 다음 명령어를 사용하세요:")
			fmt.Println("  ga optimized quick to-slim")
			fmt.Println()
			fmt.Println("migrate와 to-slim은 동일한 기능을 수행합니다.")
			fmt.Println("to-slim이 더 직관적이고 quick 카테고리에 있어 접근성이 좋습니다.")
		},
		Deprecated: "use 'ga optimized quick to-slim' instead",
	}
}