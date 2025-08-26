package help

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewWorkflowCmd creates the workflow guide command
func NewWorkflowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workflow",
		Short: "Git 최적화 워크플로우 가이드",
		Long: `Git 저장소 최적화 워크플로우를 안내합니다.
SLIM과 FULL 모드의 차이점과 각 워크플로우별 사용 시나리오를 설명합니다.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: 워크플로우 가이드 구현
			fmt.Println("1. Workflow - Git 최적화 워크플로우 가이드")
		},
	}
}