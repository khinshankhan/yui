package clipboard

import (
	"fmt"
	"os"
	"time"

	"github.com/khinshankhan/yui/lib/sysexec"
)

func Copy(text string) error {
	backend, err := detectBackend()
	if err != nil {
		return err
	}
	return sysexec.RunInput(backend, "copy", text)
}

func Paste() (string, error) {
	backend, err := detectBackend()
	if err != nil {
		return "", err
	}
	return sysexec.RunOutput(backend, "paste", 2*time.Second)
}

func detectBackend() (sysexec.Backend, error) {
	b, err := sysexec.Detect(candidates())
	if err != nil {
		return sysexec.Backend{}, fmt.Errorf("no clipboard backend found; install pbcopy/pbpaste, wl-clipboard, xclip, xsel, or PowerShell")
	}
	return b, nil
}

func candidates() []sysexec.Backend {
	switch sysexec.GOOS {
	case "darwin":
		return []sysexec.Backend{
			{
				Name: "pbcopy/pbpaste",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"pbcopy"}},
					"paste": {Args: []string{"pbpaste"}},
				},
			},
		}
	case "windows":
		return []sysexec.Backend{
			{
				Name: "powershell",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"powershell.exe", "-NoProfile", "-Command", "Set-Clipboard"}},
					"paste": {Args: []string{"powershell.exe", "-NoProfile", "-Command", "Get-Clipboard"}},
				},
			},
			{
				Name: "pwsh",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"pwsh.exe", "-NoProfile", "-Command", "Set-Clipboard"}},
					"paste": {Args: []string{"pwsh.exe", "-NoProfile", "-Command", "Get-Clipboard"}},
				},
			},
		}
	default:
		var backends []sysexec.Backend

		if os.Getenv("WAYLAND_DISPLAY") != "" {
			backends = append(backends, sysexec.Backend{
				Name: "wl-clipboard",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"wl-copy"}},
					"paste": {Args: []string{"wl-paste", "--no-newline"}},
				},
			})
		}

		if os.Getenv("DISPLAY") != "" {
			backends = append(backends,
				sysexec.Backend{
					Name: "xclip",
					Cmds: map[string]sysexec.Cmd{
						"copy":  {Args: []string{"xclip", "-silent", "-selection", "clipboard"}, DetachAfterStart: true},
						"paste": {Args: []string{"xclip", "-selection", "clipboard", "-o"}},
					},
				},
				sysexec.Backend{
					Name: "xsel",
					Cmds: map[string]sysexec.Cmd{
						"copy":  {Args: []string{"xsel", "--clipboard", "--input"}},
						"paste": {Args: []string{"xsel", "--clipboard", "--output"}},
					},
				},
			)
		}

		backends = append(backends,
			sysexec.Backend{
				Name: "wl-clipboard",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"wl-copy"}},
					"paste": {Args: []string{"wl-paste", "--no-newline"}},
				},
			},
			sysexec.Backend{
				Name: "xclip",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"xclip", "-silent", "-selection", "clipboard"}, DetachAfterStart: true},
					"paste": {Args: []string{"xclip", "-selection", "clipboard", "-o"}},
				},
			},
			sysexec.Backend{
				Name: "xsel",
				Cmds: map[string]sysexec.Cmd{
					"copy":  {Args: []string{"xsel", "--clipboard", "--input"}},
					"paste": {Args: []string{"xsel", "--clipboard", "--output"}},
				},
			},
		)

		return backends
	}
}
