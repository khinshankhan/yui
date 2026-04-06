package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Backend struct {
	Name             string
	Copy             []string
	Paste            []string
	DetachAfterStart bool
}

var (
	lookPath = exec.LookPath
	command  = exec.Command
	goos     = runtime.GOOS
)

func Copy(text string) error {
	backend, err := detectBackend()
	if err != nil {
		return err
	}

	cmd := command(backend.Copy[0], backend.Copy[1:]...)

	if backend.DetachAfterStart {
		applyDetachedProcessAttrs(cmd)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("%s copy failed: %w", backend.Name, err)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("%s copy failed: %w", backend.Name, err)
		}

		if _, err := io.WriteString(stdin, text); err != nil {
			_ = stdin.Close()
			return fmt.Errorf("%s copy failed: %w", backend.Name, err)
		}
		if err := stdin.Close(); err != nil {
			return fmt.Errorf("%s copy failed: %w", backend.Name, err)
		}

		waitCh := make(chan error, 1)
		go func() {
			waitCh <- cmd.Wait()
		}()

		select {
		case err := <-waitCh:
			if err != nil {
				return fmt.Errorf("%s copy failed: %w: %s", backend.Name, err, bytes.TrimSpace(stderr.Bytes()))
			}
			return nil
		case <-time.After(150 * time.Millisecond):
			return cmd.Process.Release()
		}
	}

	cmd.Stdin = bytes.NewBufferString(text)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s copy failed: %w: %s", backend.Name, err, bytes.TrimSpace(output))
	}

	return nil
}

func Paste() (string, error) {
	backend, err := detectBackend()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, backend.Paste[0], backend.Paste[1:]...)
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("%s paste timed out", backend.Name)
	}
	if err != nil {
		return "", fmt.Errorf("%s paste failed: %w: %s", backend.Name, err, bytes.TrimSpace(output))
	}

	return string(output), nil
}

func detectBackend() (Backend, error) {
	for _, candidate := range candidates() {
		if available(candidate.Copy[0]) && available(candidate.Paste[0]) {
			return candidate, nil
		}
	}

	return Backend{}, fmt.Errorf("no clipboard backend found; install pbcopy/pbpaste, wl-clipboard, xclip, xsel, or PowerShell")
}

func available(name string) bool {
	_, err := lookPath(name)
	return err == nil
}

func candidates() []Backend {
	switch goos {
	case "darwin":
		return []Backend{
			{
				Name:  "pbcopy/pbpaste",
				Copy:  []string{"pbcopy"},
				Paste: []string{"pbpaste"},
			},
		}
	case "windows":
		return []Backend{
			{
				Name:  "powershell",
				Copy:  []string{"powershell.exe", "-NoProfile", "-Command", "Set-Clipboard"},
				Paste: []string{"powershell.exe", "-NoProfile", "-Command", "Get-Clipboard"},
			},
			{
				Name:  "pwsh",
				Copy:  []string{"pwsh.exe", "-NoProfile", "-Command", "Set-Clipboard"},
				Paste: []string{"pwsh.exe", "-NoProfile", "-Command", "Get-Clipboard"},
			},
		}
	default:
		var backends []Backend

		if os.Getenv("WAYLAND_DISPLAY") != "" {
			backends = append(backends, Backend{
				Name:  "wl-clipboard",
				Copy:  []string{"wl-copy"},
				Paste: []string{"wl-paste", "--no-newline"},
			})
		}

		if os.Getenv("DISPLAY") != "" {
			backends = append(backends,
				Backend{
					Name:             "xclip",
					Copy:             []string{"xclip", "-silent", "-selection", "clipboard"},
					Paste:            []string{"xclip", "-selection", "clipboard", "-o"},
					DetachAfterStart: true,
				},
				Backend{
					Name:  "xsel",
					Copy:  []string{"xsel", "--clipboard", "--input"},
					Paste: []string{"xsel", "--clipboard", "--output"},
				},
			)
		}

		backends = append(backends,
			Backend{
				Name:  "wl-clipboard",
				Copy:  []string{"wl-copy"},
				Paste: []string{"wl-paste", "--no-newline"},
			},
			Backend{
				Name:             "xclip",
				Copy:             []string{"xclip", "-silent", "-selection", "clipboard"},
				Paste:            []string{"xclip", "-selection", "clipboard", "-o"},
				DetachAfterStart: true,
			},
			Backend{
				Name:  "xsel",
				Copy:  []string{"xsel", "--clipboard", "--input"},
				Paste: []string{"xsel", "--clipboard", "--output"},
			},
		)

		return backends
	}
}
