package colorcli

import (
	"fmt"
	"strings"

	"github.com/khinshankhan/yui/lib/colorconv"
)

func Usage(app string) string {
	return fmt.Sprintf(`Usage:
  %s <color>                    # Show all formats
  %s <target-format> <color>    # Convert to specific format

Target formats:
  hex   rgb   hsl   hsv/hsb   cmyk   lab   oklab   oklch

Examples:
  %s "#ff5500"
  %s rgb "#ff5500"
  %s hex "oklch(0.7 0.15 60)"`, app, app, app, app, app)
}

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("color value required")
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
