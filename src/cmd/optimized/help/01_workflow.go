package help

import (
	"fmt"
	
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewWorkflowCmd creates the workflow guide command
func NewWorkflowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "workflow",
		Aliases: []string{"wf", "guide"},
		Short:   "Git 최적화 워크플로우 가이드 표시",
		Long: `Git 저장소 최적화 워크플로우를 안내합니다.
SLIM과 FULL 모드의 차이점과 각 워크플로우별 사용 시나리오를 설명합니다.`,
		Example: `  ga optimized workflow
  ga optimized wf`,
		Run: func(cmd *cobra.Command, args []string) {
			showWorkflowGuide()
		},
	}
}

func showWorkflowGuide() {
	// 색상 정의
	titleColor := color.New(color.FgCyan, color.Bold)
	headerColor := color.New(color.FgYellow, color.Bold)
	subHeaderColor := color.New(color.FgGreen)
	modeColor := color.New(color.FgMagenta, color.Bold)
	sizeColor := color.New(color.FgRed)
	commandColor := color.New(color.FgBlue)
	dimColor := color.New(color.FgHiBlack)

	// 제목
	titleColor.Println("\n🚀 Git 저장소 최적화 워크플로우 가이드")
	fmt.Println(string(make([]byte, 60)))

	// 모드 설명
	headerColor.Println("\n📊 최적화 모드 비교")
	fmt.Println()
	
	// FULL 모드
	modeColor.Print("  FULL 모드")
	fmt.Println(" (전체 저장소)")
	fmt.Print("  • 모든 파일 및 전체 히스토리 포함")
	sizeColor.Println(" (~103GB)")
	fmt.Println("  • 모든 브랜치, 태그, 커밋 접근 가능")
	fmt.Println("  • 네트워크 없이도 모든 작업 가능")
	dimColor.Println("  • 초기 클론 시간: 약 2-3시간")
	fmt.Println()
	
	// SLIM 모드
	modeColor.Print("  SLIM 모드")
	fmt.Println(" (최적화 저장소)")
	fmt.Print("  • 필수 파일과 최소 히스토리만 포함")
	sizeColor.Println(" (~30MB)")
	fmt.Println("  • 작업 필요시 선택적 확장 가능")
	fmt.Println("  • 빠른 클론 및 작업 시작")
	dimColor.Println("  • 초기 클론 시간: 약 30초")

	// 워크플로우
	headerColor.Println("\n\n🔄 주요 워크플로우")
	fmt.Println()

	// 1. INIT-SLIM
	subHeaderColor.Println("  1️⃣  INIT-SLIM: 신규 경량 클론 (∅ → SLIM)")
	fmt.Println("  ├─ 사용 시나리오:")
	fmt.Println("  │  • 새로운 개발 환경 설정")
	fmt.Println("  │  • 빠른 프로젝트 시작이 필요한 경우")
	fmt.Println("  │  • CI/CD 환경 구성")
	fmt.Println("  └─ 명령어:")
	commandColor.Println("     ga optimized clone-slim <repository-url> <folder>")
	fmt.Println()

	// 2. MIGRATE-SLIM
	subHeaderColor.Println("  2️⃣  MIGRATE-SLIM: 기존 저장소 경량화 (FULL → SLIM)")
	fmt.Println("  ├─ 사용 시나리오:")
	fmt.Println("  │  • 이미 클론한 대용량 저장소 최적화")
	fmt.Println("  │  • 디스크 공간 확보가 필요한 경우")
	fmt.Println("  │  • 저장소 성능 개선")
	fmt.Println("  └─ 명령어:")
	commandColor.Println("     ga optimized migrate      # 또는")
	commandColor.Println("     ga optimized to-slim")
	fmt.Println()

	// 3. RESTORE-FULL
	subHeaderColor.Println("  3️⃣  RESTORE-FULL: 전체 복원 (SLIM → FULL)")
	fmt.Println("  ├─ 사용 시나리오:")
	fmt.Println("  │  • 전체 히스토리 분석 필요")
	fmt.Println("  │  • 대규모 리팩토링 작업")
	fmt.Println("  │  • 오프라인 작업 준비")
	fmt.Println("  └─ 명령어:")
	commandColor.Println("     ga optimized to-full")
	fmt.Println()

	// 4. EXPAND-SLIM
	subHeaderColor.Println("  4️⃣  EXPAND-SLIM: 선택적 확장 (SLIM → SLIM+)")
	fmt.Println("  ├─ 사용 시나리오:")
	fmt.Println("  │  • 특정 폴더/파일만 추가 필요")
	fmt.Println("  │  • 히스토리 점진적 확장")
	fmt.Println("  │  • 병합을 위한 베이스 확장")
	fmt.Println("  └─ 명령어:")
	commandColor.Println("     ga optimized expand-slim             # 대화형 경로 선택")
	commandColor.Println("     ga optimized expand-path <path>      # 특정 경로 추가")
	commandColor.Println("     ga optimized expand-10/50/100        # 히스토리 확장")

	// 권장 사용 패턴
	headerColor.Println("\n\n💡 권장 사용 패턴")
	fmt.Println()
	
	fmt.Println("  🆕 신규 개발자:")
	commandColor.Println("     ga optimized clone-slim <url> → ga optimized expand-path src/")
	fmt.Println()
	
	fmt.Println("  🔧 기존 사용자:")
	commandColor.Println("     ga optimized status → ga optimized to-slim → ga optimized expand-slim")
	fmt.Println()
	
	fmt.Println("  🚀 CI/CD:")
	commandColor.Println("     ga optimized clone-slim --depth=1 → ga optimized performance")
	fmt.Println()
	
	fmt.Println("  📊 분석 작업:")
	commandColor.Println("     ga optimized to-full → (작업) → ga optimized to-slim")

	// 추가 도움말
	dimColor.Println("\n\n💬 추가 도움말")
	fmt.Println("  • 전체 명령어 목록: ga optimized commands")
	fmt.Println("  • 현재 상태 확인: ga optimized status")
	fmt.Println("  • 성능 최적화 설정: ga optimized performance")
	fmt.Println()
}