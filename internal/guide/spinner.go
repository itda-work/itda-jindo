package guide

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner handles animated loading indicator
type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		defer close(s.done)
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the spinner line
				fmt.Printf("\r\033[K")
				return
			default:
				s.mu.Lock()
				fmt.Printf("\r%s %s", spinnerFrames[i], s.message)
				s.mu.Unlock()
				i = (i + 1) % len(spinnerFrames)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	close(s.stop)
	<-s.done
}

// StopWithMessage stops the spinner and shows a final message
func (s *Spinner) StopWithMessage(message string) {
	s.Stop()
	fmt.Println(message)
}

// RunClaudeWithSpinner runs claude command with a spinner and returns the output
func RunClaudeWithSpinner(systemPrompt, userPrompt string) (string, error) {
	spinner := NewSpinner("Claude Code를 통해 가이드 작성 중...")
	spinner.Start()

	cmd := exec.Command("claude",
		"-p", userPrompt,
		"--system-prompt", systemPrompt,
		"--output-format", "text",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		spinner.Stop()
		return "", err
	}

	if err := cmd.Start(); err != nil {
		spinner.Stop()
		return "", err
	}

	// Read output
	var output strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		output.WriteString(scanner.Text())
		output.WriteString("\n")
	}

	err = cmd.Wait()
	spinner.StopWithMessage("✅ 가이드 작성 완료!")

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// PrintGuide prints the guide content with formatting
func PrintGuide(title string, content string, createdAt time.Time, cached bool) {
	// Header
	fmt.Println()
	fmt.Printf("📚 \033[1;35m%s\033[0m\n", title)

	if cached && !createdAt.IsZero() {
		fmt.Printf("   \033[90m📅 작성: %s  |  재생성: --refresh (-r)\033[0m\n", FormatAge(createdAt))
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	// Content
	fmt.Println(content)

	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
}

// OpenHTMLGuide generates HTML and opens in browser
func OpenHTMLGuide(guideType GuideType, id string, content string, createdAt time.Time) error {
	htmlPath, err := GenerateHTML(guideType, id, content, createdAt)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	fmt.Printf("📄 HTML 생성: %s\n", htmlPath)
	fmt.Println("🌐 브라우저에서 열기...")

	return OpenInBrowser(htmlPath)
}

// RunInteractiveGuide runs interactive guide session with claude
func RunInteractiveGuide(name, systemPrompt string) error {
	fmt.Println()
	fmt.Println("🤖 AI 주도형 가이드를 시작합니다...")
	fmt.Println("   - AI가 사용자 상황에 대해 질문합니다")
	fmt.Println("   - 답변에 따라 맞춤형 안내를 제공합니다")
	fmt.Println("   - 'exit' 또는 Ctrl+C로 종료")
	fmt.Println()

	initialPrompt := fmt.Sprintf("'%s'에 대한 맞춤형 가이드를 제공하겠습니다. 먼저 사용자의 상황과 요구사항을 파악하기 위해 몇 가지 질문을 드리겠습니다.", name)

	cmd := exec.Command("claude",
		"--system-prompt", systemPrompt,
		initialPrompt,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 130 { // Ctrl+C
				fmt.Println("\n⚠️  가이드가 취소되었습니다")
				return nil
			}
		}
		return fmt.Errorf("claude command failed: %w", err)
	}

	return nil
}
