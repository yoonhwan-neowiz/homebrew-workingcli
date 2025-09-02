package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"workingcli/src/cmd/git"
	"workingcli/src/cmd/optimized"
	"workingcli/src/config"
)

var rootCmd = &cobra.Command{
	Use:   "ga",
	Short: "Git Assistant - Git 작업 생산성을 높이기 위한 AI 기반 CLI 도구",
	Long: `Git Assistant (GA)는 Git 작업 생산성을 높이기 위한 AI 기반 CLI 도구입니다.
주요 기능:
- 대화형 Git stage/unstage 인터페이스
- AI 기반 커밋 메시지 자동 생성
- Git 히스토리 시각화
- 대화형 충돌 해결
- Git 기본 명령어 지원`,
	SilenceErrors: true,  // 에러 자동 출력 비활성화 (Git fallback 처리를 위해)
	SilenceUsage:  true,  // 에러 시 사용법 자동 출력 비활성화
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// clone-slim 같은 명령어는 Git 저장소 밖에서도 실행 가능해야 함
		// optimized setup clone-slim 명령인지 확인
		if len(os.Args) >= 4 && (os.Args[1] == "optimized" || os.Args[1] == "opt" || os.Args[1] == "op" || os.Args[1] == "optimize") && os.Args[2] == "setup" && os.Args[3] == "clone-slim" {
			// clone-slim은 설정 초기화 건너뛰기
			return nil
		}
		
		// 그 외 명령어는 설정 초기화
		if err := config.Initialize(); err != nil {
			return fmt.Errorf("설정 초기화 실패: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 기본적으로 stage 모드 실행
		stageCmd := NewStageCmd()
		stageCmd.Run(cmd, args)
	},
}

func Execute() {
	// 먼저 정의된 cobra 명령어 실행 시도
	err := rootCmd.Execute()
	if err != nil {
		// cobra에서 명령어를 찾지 못한 경우 Git fallback 시도
		if len(os.Args) > 1 {
			if gitErr := HandleGitFallback(os.Args[1:]); gitErr != nil {
				// Git fallback도 실패한 경우
				fmt.Printf("❌ 실행 실패한 명령어: %s\n", strings.Join(os.Args, " "))
				fmt.Printf("   Git 명령 실행 실패: %v\n", gitErr)
				fmt.Printf("   (원래 오류: %v)\n", err)
				os.Exit(1)
			}
			// Git fallback 성공 시 정상 종료
		} else {
			// 인자가 없는 경우 cobra 에러 출력
			fmt.Printf("❌ 실행 실패한 명령어: %s\n", strings.Join(os.Args, " "))
			fmt.Printf("   오류: %v\n", err)
			os.Exit(1)
		}
	}
}

func initGitConfig() {
	// 파일 시스템 모니터링 및 캐싱 설정
	
	// core.fsmonitor: 파일 시스템 변경 감지를 위한 모니터링 활성화
	// - 파일 변경 감지 성능 향상
	// - status 명령어 속도 개선
	// - 대규모 저장소에서 특히 효과적
	exec.Command("git", "config", "core.fsmonitor", "true").Run()

	// core.untrackedCache: untracked 파일 캐시 활성화
	// - untracked 파일 목록 캐싱으로 status 성능 향상
	// - 대규모 저장소에서 특히 유용
	// - 메모리 사용량 증가 가능성 있음
	exec.Command("git", "config", "core.untrackedCache", "true").Run()

	// core.fscache: 파일 시스템 캐시 활성화
	// - 파일 상태 정보 캐싱
	// - 반복적인 파일 시스템 접근 최소화
	// - status/add 명령어 성능 향상
	exec.Command("git", "config", "core.fscache", "true").Run()

	// core.fscachesize: 파일 시스템 캐시 크기 설정
	// - 기본값: 7000
	// - 대규모 저장소의 경우 더 큰 값 설정 가능
	// - 메모리 사용량과 성능의 트레이드오프 고려
	exec.Command("git", "config", "core.fscachesize", "100000").Run()

	// 인덱스 최적화

	// pack.compression: 압축 레벨 설정 (0-9)
	// - 0: 압축하지 않음 (가장 빠름, 용량 큼)
	// - 9: 최대 압축 (가장 느림, 용량 작음)
	// - 개발 중에는 0으로 설정하여 속도 향상
	exec.Command("git", "config", "pack.compression", "0").Run()

	// core.deltaBaseCacheLimit: delta base 캐시 제한
	// - delta base 캐시 크기 제한 (바이트 단위)
	// - 작은 값으로 설정하여 메모리 사용량 감소
	// - 큰 저장소에서 메모리 부족 방지
	exec.Command("git", "config", "core.deltaBaseCacheLimit", "1").Run()

	// index.threads: 인덱스 생성 시 사용할 스레드 수
	// - 멀티코어 활용을 위한 병렬 처리
	// - 기본값: 1
	// - CPU 코어 수에 따라 조정 (일반적으로 코어 수와 동일하게 설정)
	exec.Command("git", "config", "index.threads", "4").Run()

	// 네트워크 최적화

	// http.postBuffer: HTTP 전송 버퍼 크기
	// - 기본값: 1MB
	// - 대용량 파일 전송 시 성능 향상
	// - 네트워크 대역폭에 따라 조정
	// - 524288000 = 500MB
	exec.Command("git", "config", "http.postBuffer", "524288000").Run()

	// submodule.fetchJobs: 서브모듈 동시 fetch 작업 수
	// - 서브모듈 병렬 처리로 fetch 성능 향상
	// - 기본값: 1
	// - 네트워크 상태와 시스템 리소스에 따라 조정
	exec.Command("git", "config", "submodule.fetchJobs", "4").Run()

	// push.parallel: 병렬 push 작업 수
	// - 여러 refs를 동시에 push할 때의 병렬 처리 수
	// - 기본값: 1
	// - 네트워크 상태와 시스템 리소스에 따라 조정
	exec.Command("git", "config", "push.parallel", "4").Run()

	// 기타 최적화

	// remote.origin.prune: 원격 브랜치 자동 정리
	// - fetch 시 원격에서 삭제된 브랜치 자동 정리
	// - 저장소 정리 및 관리 자동화
	// - 불필요한 참조 제거로 성능 향상
	exec.Command("git", "config", "remote.origin.prune", "true").Run()

	// gc.auto: 자동 가비지 컬렉션 임계값
	// - 기본값: 6700
	// - loose object 수가 이 값을 초과하면 자동 gc 실행
	// - 높은 값으로 설정하여 불필요한 gc 방지
	exec.Command("git", "config", "gc.auto", "256").Run()

	// gc.aggressiveWindow: aggressive gc 윈도우 크기
	// - 기본값: 250
	// - 0으로 설정하여 aggressive gc 비활성화
	// - 성능 향상을 위해 비활성화 (일반 gc만 사용)
	exec.Command("git", "config", "gc.aggressiveWindow", "0").Run()
}

func init() {
	// Git 성능 최적화 설정 적용
	initGitConfig()

	// Git 한글 파일명 설정
	exec.Command("git", "config", "--global", "core.quotepath", "false").Run()

	// GA에서 특별히 처리하는 Git 명령어 (한글 처리, 병렬 처리, 안전 검사 등)
	rootCmd.AddCommand(
		git.NewStatusCmd(),   // 상태 확인 (한글 파일명 처리)
		git.NewPullCmd(),     // 원격 저장소에서 가져오기 (서브모듈 병렬 처리)
		git.NewPushCmd(),     // 원격 저장소로 푸시 (다양한 안전 검사)
		git.NewResetCmd(),    // 특정 커밋으로 되돌리기 (안전 모드, stash 백업)
		git.NewDiffCmd(),     // 변경사항 확인 (한글 파일명 처리)
		git.NewCheckoutCmd(), // 브랜치 전환/파일 복원 (한글 파일명 처리)
		git.NewMergeCmd(),    // 브랜치 병합 (병합 전략, 충돌 처리)
		git.NewRebaseCmd(),   // 브랜치 리베이스 (안전 검사, 중단/계속 지원)
		git.NewSwitchCmd(),   // 브랜치 전환 (한글 파일명 처리)
		git.NewTagCmd(),      // 태그 관리 (한글 파일명 처리)
		git.NewFetchCmd(),    // 원격 저장소 변경사항 가져오기 (한글 파일명 처리)
		git.NewSubmoduleCmd(), // 서브모듈 관리 (foreach --recursive 지원)
	)

	// GA 특화 기능
	rootCmd.AddCommand(
		NewStageCmd(),    // 대화형 stage/unstage
		NewCommitCmd(),   // AI 기반 커밋 메시지 생성
		NewResolveCmd(),  // 대화형 충돌 해결
		NewHistoryCmd(),  // Git 히스토리 시각화
		NewAnalyzeCmd(),  // 커밋 내역 분석
	)

	// 서브커맨드 추가
	rootCmd.AddCommand(NewConfigCmd())
	
	// Git 최적화 명령어
	rootCmd.AddCommand(optimized.NewOptimizedCmd())
} 