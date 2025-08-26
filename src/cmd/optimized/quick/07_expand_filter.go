package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpandFilterCmd creates the Expand Filter removal command
func NewExpandFilterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-filter",
		Short: "Partial Clone 필터 제거 (모든 파일 다운로드)",
		Long: `Partial Clone 필터를 완전히 제거하여 모든 대용량 파일을 다운로드합니다.
Sparse Checkout은 유지하면서 blob:limit 필터만 해제하여 
현재 checkout된 경로의 모든 파일을 크기 제한 없이 다운로드합니다.
주의: 디스크 공간을 많이 사용할 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: 필터 제거 로직 구현
			fmt.Println("7. Expand Filter - Partial Clone 필터 제거")
		},
	}
}