package soundcli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/khinshankhan/yui/lib/cli"
	"github.com/khinshankhan/yui/lib/sound"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Play sound files").
		WithAliases(aliases...).
		WithSubcommandName("action").
		WithExample("Play the default notification sound", "ping").
		WithExample("Play a specific file", "play", "alert.wav").
		WithDefaultSubcommand("ping").
		Register(
			NewPingCommand("ping"),
			NewPlayCommand("play"),
		)
}

func NewPingCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Play the default notification sound").
		WithAliases(aliases...).
		WithRun(runPing)
}

func NewPlayCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Play a specific sound file").
		WithAliases(aliases...).
		WithArgs(cli.RequiredArg("file")).
		WithRun(runPlay)
}

func runPing(ctx *cli.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("ping does not accept arguments")
	}

	file := defaultSound()
	if file == "" {
		return fmt.Errorf("no default notification sound found for this platform")
	}

	return sound.Play(file)
}

func runPlay(ctx *cli.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("file argument required")
	}
	return sound.Play(args[0])
}

func defaultSound() string {
	candidates := defaultSoundCandidates()
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func defaultSoundCandidates() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/System/Library/Sounds/Ping.aiff",
			"/System/Library/Sounds/Glass.aiff",
			"/System/Library/Sounds/Tink.aiff",
		}
	case "windows":
		return []string{
			`C:\Windows\Media\notify.wav`,
			`C:\Windows\Media\chimes.wav`,
		}
	default:
		return []string{
			"/usr/share/sounds/freedesktop/stereo/bell.oga",
			"/usr/share/sounds/freedesktop/stereo/message.oga",
			"/usr/share/sounds/freedesktop/stereo/complete.oga",
		}
	}
}
