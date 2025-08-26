package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpand10Cmd creates the Expand 10 depth command
func NewExpand10Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-10",
		Short: "히스토리 10개 확장",
		Long: `현재 shallow 상태에서 10개의 커밋을 추가로 가져옵니다.
최근 브랜치 분기 확인이나 작은 규모의 병합 작업에 적합합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: 10개 히스토리 확장 로직 구현
			fmt.Println("8. Expand 10 - 히스토리 10개 확장")
		},
	}
}