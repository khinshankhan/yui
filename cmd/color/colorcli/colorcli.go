package colorcli

import (
	"fmt"
	"strings"

	"github.com/khinshankhan/yui/lib/cli"
	"github.com/khinshankhan/yui/lib/colorconv"
)

func NewCommand(name string, aliases ...string) *cli.Command {
	return cli.New(name, "Color conversion tools").
		WithAliases(aliases...).
		WithArgs(
			cli.
				OptionalArg("target-format"),
			cli.
				RequiredArg("color"),
		).
		WithSections(
			cli.Section{
				Title: "TARGET FORMATS",
				Lines: []string{
					"hex       Hexadecimal (#rrggbb)",
					"rgb       RGB (rgb(r, g, b))",
					"hsl       HSL (hsl(h, s%, l%))",
					"hsv, hsb  HSV/HSB (hsv(h, s%, v%))",
					"cmyk      CMYK (cmyk(c%, m%, y%, k%))",
					"oklab     OKLab (oklab(l a b))",
					"oklch     OKLCH (oklch(l c h))",
					"lab       CIE Lab (lab(l a b))",
				},
			},
			cli.Section{
				Title: "INPUT FORMATS",
				Lines: []string{
					"Hex:      #rgb, #rrggbb, #rrggbbaa",
					"RGB:      rgb(255, 128, 0), rgba(255, 128, 0, 0.5)",
					"HSL:      hsl(30, 100%, 50%), hsla(30, 100%, 50%, 0.5)",
					"HSV:      hsv(30, 100%, 100%)",
					"CMYK:     cmyk(0%, 50%, 100%, 0%)",
					"OKLCH:    oklch(0.7 0.15 60)",
					"OKLab:    oklab(0.7 0.1 0.1)",
					"Named:    red, blue, green, etc.",
				},
			},
		).
		WithExamples(
			"%cmd% help                        # Show help",
			"%cmd% \"#ff5500\"                   # Show all formats",
			"%cmd% rgb \"#ff5500\"               # rgb(255, 85, 0)",
			"%cmd% hsl \"#ff5500\"               # hsl(20.0, 100.0%, 50.0%)",
			"%cmd% hex \"rgb(255, 85, 0)\"       # #ff5500",
			"%cmd% oklch \"hsl(20, 100%, 50%)\"  # oklch(0.655 0.203 41.3)",
			"%cmd% red                         # Show all formats for red",
			"%cmd% hex \"oklch(0.7 0.15 60)\"    # Convert OKLCH to hex",
		).
		WithRun(run)
}

func run(ctx *cli.Context, args []string) error {
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
		fmt.Fprintln(ctx.Stdout, color.FormatAll())
		return nil
	}

	switch targetFormat {
	case "hex":
		fmt.Fprintln(ctx.Stdout, color.Hex())
	case "rgb":
		fmt.Fprintln(ctx.Stdout, color.FormatRGB())
	case "hsl":
		fmt.Fprintln(ctx.Stdout, color.FormatHSL())
	case "hsv":
		fmt.Fprintln(ctx.Stdout, color.FormatHSV())
	case "cmyk":
		fmt.Fprintln(ctx.Stdout, color.FormatCMYK())
	case "oklch":
		fmt.Fprintln(ctx.Stdout, color.FormatOKLCH())
	case "oklab":
		fmt.Fprintln(ctx.Stdout, color.FormatOKLab())
	case "lab":
		fmt.Fprintln(ctx.Stdout, color.FormatLab())
	}

	return nil
}
