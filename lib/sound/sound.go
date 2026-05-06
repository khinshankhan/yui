package sound

import (
	"fmt"
	"os/exec"
	"runtime"
)

var (
	lookPath = exec.LookPath
	command  = exec.Command
	goos     = runtime.GOOS
)

type player struct {
	bin  string
	args func(file string) []string
}

func Play(file string) error {
	for _, p := range players() {
		if _, err := lookPath(p.bin); err == nil {
			cmd := command(p.bin, p.args(file)...)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("%s failed: %w: %s", p.bin, err, output)
			}
			return nil
		}
	}
	return fmt.Errorf("no sound player found; install afplay, paplay, aplay, ffplay, sox, or mpv")
}

func players() []player {
	switch goos {
	case "darwin":
		return []player{
			{bin: "afplay", args: func(f string) []string { return []string{f} }},
		}
	case "windows":
		return []player{
			{bin: "powershell.exe", args: func(f string) []string {
				return []string{"-NoProfile", "-Command",
					fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", f)}
			}},
		}
	default:
		return []player{
			{bin: "paplay", args: func(f string) []string { return []string{f} }},
			{bin: "pw-play", args: func(f string) []string { return []string{f} }},
			{bin: "aplay", args: func(f string) []string { return []string{f} }},
			{bin: "ffplay", args: func(f string) []string { return []string{"-nodisp", "-autoexit", f} }},
			{bin: "play", args: func(f string) []string { return []string{f} }},
			{bin: "mpv", args: func(f string) []string { return []string{"--no-video", f} }},
		}
	}
}
