package optimized

import (
	"github.com/spf13/cobra"
	"workingcli/src/cmd/optimized/advanced"
	"workingcli/src/cmd/optimized/help"
	"workingcli/src/cmd/optimized/quick"
	"workingcli/src/cmd/optimized/setup"
	"workingcli/src/cmd/optimized/submodule"
	"workingcli/src/cmd/optimized/workspace"
)

// NewOptimizedCmd creates the main optimized command
func NewOptimizedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimized",
		Aliases: []string{"opt", "op", "optimize"},
		Short: "Git 저장소 최적화 명령어 (Partial Clone + Sparse Checkout)",
		Long: `Git 저장소 최적화 명령어입니다.
대용량 저장소를 효율적으로 관리하기 위한 다양한 최적화 기능을 제공합니다.

주요 기능:
- SLIM/FULL 모드 전환 (103GB → 30MB)
- Partial Clone 필터 관리
- Sparse Checkout 경로 설정
- Smart Shallow 히스토리 관리`,
		Run: func(cmd *cobra.Command, args []string) {
			// 기본 실행 시 도움말 표시
			cmd.Help()
		},
	}

	// Help 카테고리
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "도움말 및 가이드",
	}
	helpCmd.AddCommand(
		help.NewWorkflowCmd(),    // 1. Workflow
		help.NewCommandsCmd(),     // 2. Commands
	)
	
	// Quick 카테고리
	quickCmd := &cobra.Command{
		Use:   "quick",
		Short: "자주 사용하는 최적화 기능",
	}
	quickCmd.AddCommand(
		quick.NewStatusCmd(),       // 3. Status
		quick.NewToSlimCmd(),        // 4. To SLIM
		quick.NewToFullCmd(),        // 5. To FULL
		quick.NewExpandSlimCmd(),    // 6. Expand SLIM
		quick.NewExpandFilterCmd(),  // 7. Expand Filter
		quick.NewExpand10Cmd(),      // 8. Expand 10
		quick.NewExpand50Cmd(),      // 9. Expand 50
		quick.NewExpand100Cmd(),     // 10. Expand 100
		quick.NewAutoFindMergeBaseCmd(),  // 11. Auto Find Merge Base
		quick.NewCheckMergeBaseCmd(),     // 12. Check Merge Base
	)
	
	// Setup 카테고리
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "초기 설정 및 마이그레이션",
	}
	setupCmd.AddCommand(
		setup.NewCloneSlimCmd(),    // 13. Clone SLIM
		setup.NewMigrateCmd(),       // 14. Migrate
		setup.NewPerformanceCmd(),   // 15. Performance
	)
	
	// Workspace 카테고리
	workspaceCmd := &cobra.Command{
		Use:   "workspace",
		Short: "작업 공간 관리",
	}
	workspaceCmd.AddCommand(
		workspace.NewExpandPathCmd(),      // 16. Expand Path
		workspace.NewFilterBranchCmd(),    // 17. Filter Branch
		workspace.NewClearFilterBranchCmd(),     // 18. Clear Filter Branch
		workspace.NewRestoreBranchCmd(),   // 19. Restore Branch
	)
	
	// Advanced 카테고리
	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "고급 최적화 기능",
	}
	advancedCmd.AddCommand(
		advanced.NewShallowCmd(),          // 20. Shallow
		advanced.NewUnshallowCmd(),        // 21. Unshallow
		advanced.NewCheckShallowCmd(),     // 22. Check Shallow
		advanced.NewCheckFilterCmd(),      // 23. Check Filter
		advanced.NewBackupConfigCmd(),     // 24. Backup Config
	)
	
	// Submodule 카테고리
	submoduleCmd := &cobra.Command{
		Use:   "submodule",
		Short: "서브모듈 최적화",
	}
	submoduleCmd.AddCommand(
		submodule.NewShallowAllCmd(),      // 25. Shallow All
		submodule.NewUnshallowAllCmd(),    // 26. Unshallow All
		submodule.NewOptimizeAllCmd(),     // 27. Optimize All
		submodule.NewStatusAllCmd(),       // 28. Status All
	)
	
	cmd.AddCommand(helpCmd, quickCmd, setupCmd, workspaceCmd, advancedCmd, submoduleCmd)
	return cmd
}