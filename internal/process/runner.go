// Package process owns bounded external-process execution.
package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// Spec describes an external process.
type Spec struct {
	Name        string
	Args        []string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	GracePeriod time.Duration
}

// Runner executes external processes.
type Runner interface {
	Run(context.Context, Spec) error
}

// OSRunner executes operating-system processes and forwards cancellation.
type OSRunner struct{}

// Run executes a process until completion or context cancellation.
func (OSRunner) Run(ctx context.Context, spec Spec) error {
	path, err := exec.LookPath(spec.Name)
	if err != nil {
		return fmt.Errorf("find %s: %w", spec.Name, err)
	}
	command := exec.Command(path, spec.Args...)
	command.Stdin = spec.Stdin
	command.Stdout = spec.Stdout
	command.Stderr = spec.Stderr
	if err := command.Start(); err != nil {
		return fmt.Errorf("start %s: %w", spec.Name, err)
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- command.Wait()
	}()

	select {
	case err := <-waitCh:
		if err != nil {
			return fmt.Errorf("%s failed: %w", spec.Name, err)
		}
		return nil
	case <-ctx.Done():
		if command.Process == nil {
			return fmt.Errorf("%s cancelled: %w", spec.Name, ctx.Err())
		}
		if err := command.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("interrupt %s: %w", spec.Name, err)
		}
		grace := spec.GracePeriod
		if grace <= 0 {
			grace = 10 * time.Second
		}
		timer := time.NewTimer(grace)
		defer timer.Stop()
		select {
		case <-waitCh:
			return fmt.Errorf("%s cancelled: %w", spec.Name, ctx.Err())
		case <-timer.C:
			if err := command.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
				return fmt.Errorf("kill %s after cancellation: %w", spec.Name, err)
			}
			<-waitCh
			return fmt.Errorf("%s killed after cancellation: %w", spec.Name, ctx.Err())
		}
	}
}
