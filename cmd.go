package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// buildPrompt assembles the final prompt from pipe input and user arguments.
func buildPrompt(args []string, pipeContent string) string {
	prompt := strings.Join(args, " ")

	if pipeContent != "" && prompt != "" {
		return fmt.Sprintf("Here is the input data:\n\n```\n%s\n```\n\n%s", pipeContent, prompt)
	}
	if pipeContent != "" {
		return fmt.Sprintf("Here is some data. Please analyze it:\n\n```\n%s\n```", pipeContent)
	}
	return prompt
}

// onFirstWriteWriter wraps an io.Writer and calls a callback on the first Write.
type onFirstWriteWriter struct {
	w    io.Writer
	once sync.Once
	fn   func()
}

func (o *onFirstWriteWriter) Write(p []byte) (int, error) {
	o.once.Do(o.fn)
	return o.w.Write(p)
}

// runClaude executes the claude CLI in single-shot mode with the given prompt.
func runClaude(prompt, model string) error {
	claudePath, err := findClaude()
	if err != nil {
		return fmt.Errorf("claude CLI not found in PATH: %w\nInstall it from: https://docs.anthropic.com/en/docs/claude-code", err)
	}

	args := []string{"-p", prompt}

	if model != "" {
		args = append(args, "--model", model)
	}

	// Append passthrough flags for claude
	args = append(args, passthrough...)

	if dryRun {
		fmt.Printf("claude -p %q", prompt)
		if model != "" {
			fmt.Printf(" --model %s", model)
		}
		for _, a := range passthrough {
			fmt.Printf(" %s", a)
		}
		fmt.Println()
		return nil
	}

	cmd := exec.Command(claudePath, args...)
	cmd.Stdin = os.Stdin

	needRender := !rawOutput && isStdoutTerminal()

	var outBuf bytes.Buffer
	sp := startSpinner()

	if needRender {
		// Buffer output for styled rendering; spinner keeps running until done
		cmd.Stdout = &outBuf
	} else {
		// Raw / piped: stream directly, stop spinner on first byte
		cmd.Stdout = &onFirstWriteWriter{w: os.Stdout, fn: sp.Stop}
	}
	cmd.Stderr = &onFirstWriteWriter{w: os.Stderr, fn: sp.Stop}

	if err := cmd.Run(); err != nil {
		sp.Stop()
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("claude exited with code %d", exitErr.ExitCode())
		}
		return fmt.Errorf("failed to run claude: %w", err)
	}
	sp.Stop()

	if needRender {
		raw := outBuf.String()
		if raw == "" {
			return nil
		}
		rendered, err := renderMarkdown(strings.TrimSpace(raw))
		if err != nil {
			fmt.Print(raw)
			return nil
		}
		fmt.Print(rendered)
	}
	return nil
}
