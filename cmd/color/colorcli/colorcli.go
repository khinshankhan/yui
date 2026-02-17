package colorcli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/khinshankhan/yui/lib/colorconv"
)

var ErrHelpRequested = errors.New("help requested")

func isHelpArg(arg string) bool {
	switch arg {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func Help(app string) string {
	help := `%s - Color conversion tools

USAGE:
    %s <color>                    # Show all formats
    %s <target-format> <color>    # Convert to specific format

SUBCOMMANDS:
    help      Show this help message

TARGET FORMATS:
    hex       Hexadecimal (#rrggbb)
    rgb       RGB (rgb(r, g, b))
    hsl       HSL (hsl(h, s%, l%))
    hsv       HSV/HSB (hsv(h, s%, v%))
    cmyk      CMYK (cmyk(c%, m%, y%, k%))
    oklab     OKLab (oklab(l a b))
    oklch     OKLCH (oklch(l c h))
    lab       CIE Lab (lab(l a b))

INPUT FORMATS:
    Hex:      #rgb, #rrggbb, #rrggbbaa
    RGB:      rgb(255, 128, 0), rgba(255, 128, 0, 0.5)
    HSL:      hsl(30, 100%, 50%), hsla(30, 100%, 50%, 0.5)
    HSV:      hsv(30, 100%, 100%)
    CMYK:     cmyk(0%, 50%, 100%, 0%)
    OKLCH:    oklch(0.7 0.15 60)
    OKLab:    oklab(0.7 0.1 0.1)
    Named:    red, blue, green, etc.

EXAMPLES:
    %s help                        # Show help
    %s "#ff5500"                   # Show all formats
    %s rgb "#ff5500"               # rgb(255, 85, 0)
    %s hsl "#ff5500"               # hsl(20.0, 100.0%, 50.0%)
    %s hex "rgb(255, 85, 0)"       # #ff5500
    %s oklch "hsl(20, 100%, 50%)"  # oklch(0.655 0.203 41.3)
    %s red                         # Show all formats for red
    %s hex "oklch(0.7 0.15 60)"    # Convert OKLCH to hex`

	return strings.ReplaceAll(help, "%s", app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("color value required")
	}
	if isHelpArg(args[0]) {
		return ErrHelpRequested
	}

	subCmd := strings.ToLower(args[0])
	targetFormat := ""
	colorInput := ""

	switch subCmd {
	case "hex", "rgb", "hsl", "cmyk", "oklch", "oklab", "lab":
		targetFormat = subCmd
		colorInput = strings.Join(args[1:], " ")
	case "hsv", "hsb":
		targetFormat = "hsv"
		colorInput = strings.Join(args[1:], " ")
	default:
		colorInput = strings.Join(args, " ")
	}

	if colorInput == "" {
		return fmt.Errorf("color value required")
	}

	color, err := colorconv.Parse(colorInput)
	if err != nil {
		return err
	}

	if targetFormat == "" {
		fmt.Println(color.FormatAll())
		return nil
	}

	switch targetFormat {
	case "hex":
		fmt.Println(color.Hex())
	case "rgb":
		fmt.Println(color.FormatRGB())
	case "hsl":
		fmt.Println(color.FormatHSL())
	case "hsv":
		fmt.Println(color.FormatHSV())
	case "cmyk":
		fmt.Println(color.FormatCMYK())
	case "oklch":
		fmt.Println(color.FormatOKLCH())
	case "oklab":
		fmt.Println(color.FormatOKLab())
	case "lab":
		fmt.Println(color.FormatLab())
	}

	return nil
}
