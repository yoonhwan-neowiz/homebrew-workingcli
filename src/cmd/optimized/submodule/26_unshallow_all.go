package submodule

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewUnshallowAllCmd creates the Unshallow All Submodules command
func NewUnshallowAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unshallow-all",
		Short: "모든 서브모듈 히스토리 복원",
		Long: `모든 서브모듈의 전체 히스토리를 복원합니다.
Shallow 상태의 서브모듈을 완전한 저장소로 변환합니다.

실행 내용:
- git submodule foreach 'git fetch --unshallow'
- 모든 서브모듈의 전체 커밋 히스토리 다운로드

⚠️ 주의: 서브모듈 크기에 따라 상당한 디스크 공간이 필요할 수 있습니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Unshallow All Submodules 로직 구현
			fmt.Println("26. Unshallow All - 서브모듈 히스토리 복원")
		},
	}
}