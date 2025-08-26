package quick

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

// NewExpand100Cmd creates the Expand 100 depth command (deprecated - use expand command)
func NewExpand100Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "expand-100",
		Short: "히스토리 100개 확장 (deprecated - expand 100 사용)",
		Long: `이 명령어는 deprecated입니다. 
대신 'ga opt quick expand 100'을 사용하세요.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚠️  이 명령어는 deprecated입니다.")
			fmt.Println("   대신 사용: ga opt quick expand 100")
		},
	}
}