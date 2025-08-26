package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewToFullCmd creates the To FULL restoration command
func NewToFullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "to-full",
		Short: "FULL 모드로 복원 (30MB → 103GB)",
		Long: `저장소를 FULL 모드로 복원합니다.
모든 최적화를 해제하고 전체 히스토리와 파일을 다운로드합니다.
주의: 103GB의 디스크 공간이 필요합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: FULL 복원 로직 구현
			fmt.Println("5. To FULL - 전체 복원 (30MB → 103GB)")
		},
	}
}