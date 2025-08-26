package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"workingcli/src/config"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "설정 관리",
		Long: `GA CLI의 설정을 관리합니다.
API 키 설정, 프롬프트 경로 설정 등을 수행할 수 있습니다.`,
	}

	// get 서브커맨드
	getCmd := &cobra.Command{
		Use:   "get [key]",
		Short: "설정 값 조회",
		Long: `설정 값을 조회합니다.
키를 지정하지 않으면 전체 설정을 조회합니다.

사용법:
  ga config get           # 전체 설정 조회
  ga config get [key]    # 특정 설정 값 조회

예시:
  ga config get                    # 전체 설정 조회
  ga config get ai.provider        # AI 제공자 설정 조회
  ga config get prompt.analyze     # 분석 프롬프트 경로 조회`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// 전체 설정 출력
				settings := config.GetAll()
				fmt.Println("\n=== GA CLI 설정 상태 ===")
				
				// AI 설정
				fmt.Println("\n[AI 설정]")
				fmt.Printf("ai.provider = %v\n", settings["ai"].(map[string]interface{})["provider"])
				
				// API 키는 보안을 위해 마스킹 처리
				aiSettings := settings["ai"].(map[string]interface{})
				if openai, ok := aiSettings["openai"].(map[string]interface{}); ok {
					if key, exists := openai["api_key"].(string); exists && key != "" {
						fmt.Printf("ai.openai.api_key = %s***%s\n", key[:3], key[len(key)-4:])
					} else {
						fmt.Println("ai.openai.api_key = <not set>")
					}
					if model, exists := openai["model"].(string); exists && model != "" {
						fmt.Printf("ai.openai.model = %s\n", model)
					}
				}
				
				if claude, ok := aiSettings["claude"].(map[string]interface{}); ok {
					if key, exists := claude["api_key"].(string); exists && key != "" {
						fmt.Printf("ai.claude.api_key = %s***%s\n", key[:3], key[len(key)-4:])
					} else {
						fmt.Println("ai.claude.api_key = <not set>")
					}
					if model, exists := claude["model"].(string); exists && model != "" {
						fmt.Printf("ai.claude.model = %s\n", model)
					}
				}

				// 프롬프트 설정
				fmt.Println("\n[프롬프트 설정]")
				if prompt, ok := settings["prompt"].(map[string]interface{}); ok {
					if analyze, exists := prompt["analyze"].(string); exists {
						fmt.Printf("prompt.analyze = %s\n", analyze)
					}
					if commit, exists := prompt["commit"].(string); exists {
						fmt.Printf("prompt.commit = %s\n", commit)
					}
				}

				// 최적화 설정
				fmt.Println("\n[Git 최적화 설정]")
				if optimize, ok := settings["optimize"].(map[string]interface{}); ok {
					if mode, exists := optimize["mode"].(string); exists {
						fmt.Printf("optimize.mode = %s\n", mode)
					}
					
					if filter, ok := optimize["filter"].(map[string]interface{}); ok {
						if defaultFilter, exists := filter["default"].(string); exists {
							fmt.Printf("optimize.filter.default = %s\n", defaultFilter)
						}
						if options, ok := filter["options"].(map[string]interface{}); ok {
							fmt.Println("optimize.filter.options:")
							for key, value := range options {
								fmt.Printf("  - %s: %v\n", key, value)
							}
						}
					}
					
					if sparse, ok := optimize["sparse"].(map[string]interface{}); ok {
						if paths, ok := sparse["paths"].([]interface{}); ok {
							if len(paths) > 0 {
								fmt.Println("optimize.sparse.paths:")
								for _, path := range paths {
									fmt.Printf("  - %v\n", path)
								}
							} else {
								fmt.Println("optimize.sparse.paths = []")
							}
						}
					}
				}

				return nil
			}

			// 특정 키 값 조회
			key := args[0]
			value := config.GetString(key)
			if strings.Contains(key, "api_key") && value != "" {
				fmt.Printf("%s = %s***%s\n", key, value[:3], value[len(value)-4:])
			} else {
				fmt.Printf("%s = %s\n", key, value)
			}
			return nil
		},
	}

	// set 서브커맨드
	setCmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "설정 값 저장",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]
			if err := config.Set(key, value); err != nil {
				return fmt.Errorf("설정 저장 실패: %w", err)
			}
			fmt.Printf("%s = %s\n", key, value)
			return nil
		},
	}

	// init 서브커맨드
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "설정 초기화",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Initialize(); err != nil {
				return fmt.Errorf("설정 초기화 실패: %w", err)
			}
			fmt.Println("설정이 초기화되었습니다.")
			return nil
		},
	}

	cmd.AddCommand(getCmd, setCmd, initCmd)
	return cmd
} 