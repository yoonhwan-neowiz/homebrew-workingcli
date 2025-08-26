package advanced

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewUnshallowCmd creates the Unshallow command
func NewUnshallowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unshallow",
		Short: "히스토리 복원",
		Long: `전체 히스토리를 복원합니다.
과거 커밋 조회나 blame이 필요한 경우 사용합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Unshallow 로직 구현
			fmt.Println("21. Unshallow - 히스토리 복원")
		},
	}
}