package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpand50Cmd creates the Expand 50 depth command
func NewExpand50Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-50",
		Short: "히스토리 50개 확장",
		Long: `현재 shallow 상태에서 50개의 커밋을 추가로 가져옵니다.
중간 규모 작업이나 스프린트 단위 병합에 적합합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: 50개 히스토리 확장 로직 구현
			fmt.Println("9. Expand 50 - 히스토리 50개 확장")
		},
	}
}