package sysexec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"
)

var (
	LookPath = exec.LookPath
	NewCmd   = exec.Command
	GOOS     = runtime.GOOS
)

type Cmd struct {
	Args             []string
	DetachAfterStart bool
}

type Backend struct {
	Name string
	Cmds map[string]Cmd
}

func Detect(candidates []Backend) (Backend, error) {
	for _, b := range candidates {
		if allAvailable(b) {
			return b, nil
		}
	}
	return Backend{}, fmt.Errorf("no suitable backend found")
}

func allAvailable(b Backend) bool {
	seen := make(map[string]bool)
	for _, cmd := range b.Cmds {
		bin := cmd.Args[0]
		if seen[bin] {
			continue
		}
		if !Available(bin) {
			return false
		}
		seen[bin] = true
	}
	return true
}

func Available(name string) bool {
	_, err := LookPath(name)
	return err == nil
}

func RunInput(b Backend, op, text string) error {
	cmd, ok := b.Cmds[op]
	if !ok {
		return fmt.Errorf("%s: unknown operation %q", b.Name, op)
	}

	c := NewCmd(cmd.Args[0], cmd.Args[1:]...)

	if cmd.DetachAfterStart {
		applyDetachedProcessAttrs(c)

		stdin, err := c.StdinPipe()
		if err != nil {
			return fmt.Errorf("%s %s failed: %w", b.Name, op, err)
		}

		var stderr bytes.Buffer
		c.Stderr = &stderr

		if err := c.Start(); err != nil {
			return fmt.Errorf("%s %s failed: %w", b.Name, op, err)
		}

		if _, err := io.WriteString(stdin, text); err != nil {
			_ = stdin.Close()
			return fmt.Errorf("%s %s failed: %w", b.Name, op, err)
		}
		if err := stdin.Close(); err != nil {
			return fmt.Errorf("%s %s failed: %w", b.Name, op, err)
		}

		waitCh := make(chan error, 1)
		go func() {
			waitCh <- c.Wait()
		}()

		select {
		case err := <-waitCh:
			if err != nil {
				return fmt.Errorf("%s %s failed: %w: %s", b.Name, op, err, bytes.TrimSpace(stderr.Bytes()))
			}
			return nil
		case <-time.After(150 * time.Millisecond):
			return c.Process.Release()
		}
	}

	c.Stdin = bytes.NewBufferString(text)
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s failed: %w: %s", b.Name, op, err, bytes.TrimSpace(output))
	}
	return nil
}

func RunOutput(b Backend, op string, timeout time.Duration) (string, error) {
	cmd, ok := b.Cmds[op]
	if !ok {
		return "", fmt.Errorf("%s: unknown operation %q", b.Name, op)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)
	output, err := c.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("%s %s timed out", b.Name, op)
	}
	if err != nil {
		return "", fmt.Errorf("%s %s failed: %w: %s", b.Name, op, err, bytes.TrimSpace(output))
	}
	return string(output), nil
}
