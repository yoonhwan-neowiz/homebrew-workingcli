package submodule

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewStatusAllCmd creates the Status All Submodules command
func NewStatusAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status-all",
		Short: "모든 서브모듈 상태 확인",
		Long: `모든 서브모듈의 최적화 상태를 확인합니다.
각 서브모듈의 용량과 Shallow 상태를 표시합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Status All Submodules 로직 구현
			fmt.Println("28. Status All - 서브모듈 상태")
		},
	}
}