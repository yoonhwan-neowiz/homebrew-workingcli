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
		help.NewWorkflowCmd(),    // Workflow
		help.NewCommandsCmd(),     // Commands
	)
	
	// Quick 카테고리
	quickCmd := &cobra.Command{
		Use:   "quick",
		Short: "자주 사용하는 최적화 기능",
	}
	quickCmd.AddCommand(
		quick.NewStatusCmd(),       // Status
		quick.NewToSlimCmd(),        // To SLIM
		quick.NewToFullCmd(),        // To FULL
		quick.NewExpandSlimCmd(),    // Expand SLIM
		quick.NewClearPartialCloneCmd(),  // Clear Partial Clone
		quick.NewAutoFindMergeBaseCmd(),  // Auto Find Merge Base
		quick.NewSetBranchScopeCmd(),    // Set Branch Scope
		quick.NewClearBranchScopeCmd(),  // Clear Branch Scope
		quick.NewShallowCmd(),          // Shallow
		quick.NewUnshallowCmd(),        // Unshallow
	)
	
	// Setup 카테고리
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "초기 설정 및 마이그레이션",
	}
	setupCmd.AddCommand(
		setup.NewCloneSlimCmd(),    // Clone SLIM
		setup.NewCloneMasterCmd(),  // Clone Master
		setup.NewMigrateCmd(),       // Migrate
		setup.NewPerformanceCmd(),   // Performance
	)
	
	// Workspace 카테고리
	workspaceCmd := &cobra.Command{
		Use:   "workspace",
		Short: "작업 공간 관리",
	}
	workspaceCmd.AddCommand(
		workspace.NewExpandPathCmd(),      // Expand Path
		workspace.NewRestoreBranchCmd(),   // Restore Branch
	)
	
	// Advanced 카테고리
	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "고급 최적화 기능",
	}
	advancedCmd.AddCommand(
		advanced.NewExpand10Cmd(),      // Expand 10
		advanced.NewExpand50Cmd(),      // Expand 50
		advanced.NewExpand100Cmd(),     // Expand 100
		advanced.NewCheckMergeBaseCmd(),     // Check Merge Base
		advanced.NewCheckShallowCmd(),     // Check Shallow
		advanced.NewCheckFilterCmd(),      // Check Filter
		advanced.NewConfigCmd(),            // Config
	)
	
	// Submodule 카테고리
	submoduleCmd := &cobra.Command{
		Use:   "submodule",
        Aliases: []string{"sub"},
		Short: "서브모듈 최적화",
	}
	submoduleCmd.AddCommand(
		// 개별 서브모듈 제어
		submodule.NewStatusCmd(),        // Status (submodule)
		submodule.NewToSlimCmd(),         // To-Slim (submodule)
		submodule.NewToFullCmd(),         // To-Full (submodule)
		submodule.NewExpandSlimCmd(),    // Expand-Slim (submodule)
		submodule.NewExpandFilterCmd(),   // Expand-Filter (submodule)
		
		// 전체 서브모듈 제어
		submodule.NewShallowCmd(),         // Shallow (recursive)
		submodule.NewUnshallowCmd(),       // Unshallow (recursive)
		submodule.NewUpdateCmd(),           // Update (서브모듈 업데이트)
		
		// 브랜치 필터
		submodule.NewSetBranchScopeCmd(),    // Set Branch Scope (submodule)
		submodule.NewClearBranchScopeCmd(), // Clear Branch Scope (submodule)
	)
	
	cmd.AddCommand(helpCmd, quickCmd, setupCmd, workspaceCmd, advancedCmd, submoduleCmd)
	return cmd
}