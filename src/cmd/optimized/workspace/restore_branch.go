package workspace

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewRestoreBranchCmd creates the Restore Branch command
// DEPRECATED: This command is deprecated. Use filter-branch and clear-filter instead.
func NewRestoreBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore-branch",
		Short: "[DEPRECATED] 브랜치 전체 복원 (사용하지 않음)",
		Long: `[DEPRECATED] 이 명령어는 더 이상 사용되지 않습니다.

대신 다음 명령어를 사용하세요:
  • ga opt workspace filter-branch - 브랜치 필터 설정
  • ga opt workspace clear-filter - 브랜치 필터 제거`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚠️  DEPRECATED: 이 명령어는 더 이상 사용되지 않습니다.")
			fmt.Println("")
			fmt.Println("대신 다음 명령어를 사용하세요:")
			fmt.Println("  • ga opt workspace filter-branch - 브랜치 필터 설정")
			fmt.Println("  • ga opt workspace clear-filter - 브랜치 필터 제거")
		},
	}
}