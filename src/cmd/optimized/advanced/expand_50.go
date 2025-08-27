package advanced

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpand50Cmd creates the Expand 50 depth command (deprecated - use expand command)
func NewExpand50Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-50",
		Short: "히스토리 50개 확장 (deprecated - expand 50 사용)",
		Long: `이 명령어는 deprecated입니다. 
대신 'ga opt quick expand 50'을 사용하세요.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚠️  이 명령어는 deprecated입니다.")
			fmt.Println("   대신 사용: ga opt quick expand 50")
		},
	}
}