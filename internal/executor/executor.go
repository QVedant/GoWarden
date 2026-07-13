package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/QVedant/GoWarden/internal/registry"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	TimedOut bool
}

func Run(ctx context.Context, lang registry.Language, source string) (*Result, error) {
	workDir, err := os.MkdirTemp("", "gowarden-*")
	if err != nil {
		return nil, fmt.Errorf("creating workspace: %w", err)
	}
	defer os.RemoveAll(workDir)

	sourceFile := filepath.Join(workDir, "main"+lang.Extension)
	if err := os.WriteFile(sourceFile, []byte(source), 0644); err != nil {
		return nil, fmt.Errorf("writing source file: %w", err)
	}

	binaryPath := filepath.Join(workDir, "binary")

	//compile
	if len(lang.Compile) > 0 {
		compileArgv := substitute(lang.Compile, sourceFile, binaryPath)
		if err := runStep(ctx, workDir, compileArgv, lang.TimeoutSeconds); err != nil {
			return nil, fmt.Errorf("compile failed: %w", err)
		}
	}

	runArgv := substitute(lang.Run, sourceFile, binaryPath)
	return runCaptured(ctx, workDir, runArgv, lang.TimeoutSeconds)
}

func substitute(argv []string, file, binary string) []string {
	out := make([]string, len(argv))
	for i, arg := range argv {
		arg = strings.ReplaceAll(arg, "{file}", file)
		arg = strings.ReplaceAll(arg, "{binary}", binary)
		out[i] = arg
	}
	return out
}

func runStep(ctx context.Context, dir string, argv []string, timeoutSeconds int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}

func runCaptured(ctx context.Context, dir string, argv []string, timeoutSeconds int) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.ExitCode = -1
		return result, nil
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil // non-zero exit is a valid result, not a Go error
		}
		return nil, fmt.Errorf("running command: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}
