package colorcli

import (
	"fmt"
	"os"

	"github.com/khinshankhan/yui/cli/cobrawrapper"
	"github.com/khinshankhan/yui/lib/colorconv"
	"github.com/spf13/cobra"
	"strings"
)

func CreateColorCmd(prefixCmds []string) *cobra.Command {
	return cobrawrapper.CreateCmd(
		prefixCmds,
		&cobra.Command{
			Use:   "%s <subcommand>",
			Short: "Color conversion tools",
			Long: `USAGE:
    %s <color>                    # Show all formats
    %s <target-format> <color>    # Convert to specific format

TARGET FORMATS:
    hex       Hexadecimal (#rrggbb)
    rgb       RGB (rgb(r, g, b))
    hsl       HSL (hsl(h, s%, l%))
    hsv       HSV/HSB (hsv(h, s%, v%))
    cmyk      CMYK (cmyk(c%, m%, y%, k%))
    lab       CIE Lab (lab(l a b))
    oklab     OKLab (oklab(l a b))
    oklch     OKLCH (oklch(l c h))

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
    %s "#ff5500"                   # Show all formats
    %s rgb "#ff5500"               # rgb(255, 85, 0)
    %s hsl "#ff5500"               # hsl(20.0, 100.0%, 50.0%)
    %s hex "rgb(255, 85, 0)"       # #ff5500
    %s oklch "hsl(20, 100%, 50%)"  # oklch(0.655 0.203 41.3)
    %s red                         # Show all formats for red
    %s hex "oklch(0.7 0.15 60)"    # Convert OKLCH to hex`,
			Args: cobra.MinimumNArgs(1),

			Run: func(cmd *cobra.Command, args []string) {
				subCmd := args[0]

				// Check for format conversion: cmd <format> <color>
				targetFormat := ""
				colorInput := ""

				switch strings.ToLower(subCmd) {
				case "hex":
					targetFormat = "hex"
					colorInput = strings.Join(args[1:], " ")
				case "rgb":
					targetFormat = "rgb"
					colorInput = strings.Join(args[1:], " ")
				case "hsl":
					targetFormat = "hsl"
					colorInput = strings.Join(args[1:], " ")
				case "hsv", "hsb":
					targetFormat = "hsv"
					colorInput = strings.Join(args[1:], " ")
				case "cmyk":
					targetFormat = "cmyk"
					colorInput = strings.Join(args[1:], " ")
				case "oklch":
					targetFormat = "oklch"
					colorInput = strings.Join(args[1:], " ")
				case "oklab":
					targetFormat = "oklab"
					colorInput = strings.Join(args[1:], " ")
				case "lab":
					targetFormat = "lab"
					colorInput = strings.Join(args[1:], " ")
				default:
					// No format specified, show all formats
					colorInput = strings.Join(args, " ")
				}

				if colorInput == "" {
					fmt.Fprintln(os.Stderr, "Error: color value required")
					os.Exit(1)
				}

				color, err := colorconv.Parse(colorInput)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				if targetFormat == "" {
					// Show all formats
					fmt.Println(color.FormatAll())
				} else {
					// Show specific format
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
				}
			},
		},
	)
}
