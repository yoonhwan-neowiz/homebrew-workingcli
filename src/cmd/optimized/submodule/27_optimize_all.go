package submodule

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewOptimizeAllCmd creates the Optimize All Submodules command
func NewOptimizeAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "optimize-all",
		Short: "모든 서브모듈 SLIM 최적화",
		Long: `모든 서브모듈에 SLIM 최적화를 일괄 적용합니다.
Partial Clone, Sparse Checkout, Shallow를 모두 적용하여
서브모듈의 디스크 사용량을 최소화합니다.

적용되는 최적화:
- Shallow Clone (depth=1)로 변환
- Partial Clone 필터 적용 (blob:limit=1m)
- Sparse Checkout 설정 (필수 파일만)

실행 내용:
- git submodule foreach 'git config core.sparseCheckout true'
- git submodule foreach 'git fetch --depth=1'
- git submodule foreach 'git config remote.origin.partialclonefilter blob:limit=1m'

⚠️ 주의: 서브모듈의 전체 기능이 제한될 수 있습니다.
필요시 개별 서브모듈을 unshallow하여 복원하세요.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Optimize All Submodules 로직 구현
			fmt.Println("27. Optimize All - 서브모듈 SLIM 최적화")
		},
	}
}