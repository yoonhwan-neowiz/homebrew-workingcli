package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpand100Cmd creates the Expand 100 depth command
func NewExpand100Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-100",
		Short: "히스토리 100개 확장",
		Long: `현재 shallow 상태에서 100개의 커밋을 추가로 가져옵니다.
대규모 작업이나 릴리즈 브랜치 병합에 적합합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: 100개 히스토리 확장 로직 구현
			fmt.Println("10. Expand 100 - 히스토리 100개 확장")
		},
	}
}