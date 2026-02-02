package colorconv

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Color represents a color that can be converted between different color spaces.
type Color struct {
	R, G, B float64 // RGB values normalized to 0-1
	A       float64 // Alpha channel 0-1
}

// NewRGB creates a color from RGB values (0-255).
func NewRGB(r, g, b int) Color {
	return Color{
		R: clamp01(float64(r) / 255),
		G: clamp01(float64(g) / 255),
		B: clamp01(float64(b) / 255),
		A: 1,
	}
}

// NewRGBA creates a color from RGBA values (0-255, alpha 0-1).
func NewRGBA(r, g, b int, a float64) Color {
	return Color{
		R: clamp01(float64(r) / 255),
		G: clamp01(float64(g) / 255),
		B: clamp01(float64(b) / 255),
		A: clamp01(a),
	}
}

// NewRGBFloat creates a color from RGB values (0-1).
func NewRGBFloat(r, g, b float64) Color {
	return Color{
		R: clamp01(r),
		G: clamp01(g),
		B: clamp01(b),
		A: 1,
	}
}

// RGB returns the color as RGB values (0-255).
func (c Color) RGB() (r, g, b int) {
	return int(math.Round(c.R * 255)),
		int(math.Round(c.G * 255)),
		int(math.Round(c.B * 255))
}

// RGBA returns the color as RGBA values (0-255, alpha 0-1).
func (c Color) RGBA() (r, g, b int, a float64) {
	return int(math.Round(c.R * 255)),
		int(math.Round(c.G * 255)),
		int(math.Round(c.B * 255)),
		c.A
}

// Hex returns the color as a hex string (e.g., "#ff5500").
func (c Color) Hex() string {
	r, g, b := c.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// HexAlpha returns the color as a hex string with alpha (e.g., "#ff550080").
func (c Color) HexAlpha() string {
	r, g, b := c.RGB()
	a := int(math.Round(c.A * 255))
	return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
}

// HSL returns the color as HSL values (h: 0-360, s: 0-100, l: 0-100).
func (c Color) HSL() (h, s, l float64) {
	r, g, b := c.R, c.G, c.B

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2

	if max == min {
		h, s = 0, 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return h, s * 100, l * 100
}

// HSV returns the color as HSV values (h: 0-360, s: 0-100, v: 0-100).
func (c Color) HSV() (h, s, v float64) {
	r, g, b := c.R, c.G, c.B

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	v = max

	if max == 0 {
		s = 0
	} else {
		s = (max - min) / max
	}

	if max == min {
		h = 0
	} else {
		d := max - min
		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return h, s * 100, v * 100
}

// CMYK returns the color as CMYK values (0-100).
func (c Color) CMYK() (cyan, magenta, yellow, key float64) {
	r, g, b := c.R, c.G, c.B

	key = 1 - math.Max(r, math.Max(g, b))
	if key == 1 {
		return 0, 0, 0, 100
	}

	cyan = (1 - r - key) / (1 - key) * 100
	magenta = (1 - g - key) / (1 - key) * 100
	yellow = (1 - b - key) / (1 - key) * 100
	key *= 100

	return cyan, magenta, yellow, key
}

// Linear converts sRGB to linear RGB.
func (c Color) Linear() Color {
	return Color{
		R: srgbToLinear(c.R),
		G: srgbToLinear(c.G),
		B: srgbToLinear(c.B),
		A: c.A,
	}
}

// XYZ returns the color in CIE XYZ color space (D65 illuminant).
func (c Color) XYZ() (x, y, z float64) {
	lin := c.Linear()
	r, g, b := lin.R, lin.G, lin.B

	// sRGB to XYZ matrix (D65)
	x = r*0.4124564 + g*0.3575761 + b*0.1804375
	y = r*0.2126729 + g*0.7151522 + b*0.0721750
	z = r*0.0193339 + g*0.1191920 + b*0.9503041

	return x, y, z
}

// Lab returns the color in CIE Lab color space.
func (c Color) Lab() (l, a, b float64) {
	x, y, z := c.XYZ()

	// D65 reference white
	const xn, yn, zn = 0.95047, 1.0, 1.08883

	x /= xn
	y /= yn
	z /= zn

	x = labF(x)
	y = labF(y)
	z = labF(z)

	l = 116*y - 16
	a = 500 * (x - y)
	b = 200 * (y - z)

	return l, a, b
}

// OKLab returns the color in OKLab color space.
func (c Color) OKLab() (l, a, b float64) {
	lin := c.Linear()
	r, g, bl := lin.R, lin.G, lin.B

	// Linear RGB to LMS
	L := 0.4122214708*r + 0.5363325363*g + 0.0514459929*bl
	M := 0.2119034982*r + 0.6806995451*g + 0.1073969566*bl
	S := 0.0883024619*r + 0.2817188376*g + 0.6299787005*bl

	// LMS to OKLab
	L_ := math.Cbrt(L)
	M_ := math.Cbrt(M)
	S_ := math.Cbrt(S)

	l = 0.2104542553*L_ + 0.7936177850*M_ - 0.0040720468*S_
	a = 1.9779984951*L_ - 2.4285922050*M_ + 0.4505937099*S_
	b = 0.0259040371*L_ + 0.7827717662*M_ - 0.8086757660*S_

	return l, a, b
}

// OKLCH returns the color in OKLCH color space (l: 0-1, c: 0-0.4+, h: 0-360).
func (c Color) OKLCH() (l, chroma, h float64) {
	l, a, b := c.OKLab()
	chroma = math.Sqrt(a*a + b*b)
	h = math.Atan2(b, a) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return l, chroma, h
}

// NewFromHex creates a color from a hex string.
func NewFromHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	hex = strings.ToLower(hex)

	var r, g, b, a int
	a = 255

	switch len(hex) {
	case 3: // RGB shorthand
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		r, g, b = r*17, g*17, b*17
	case 4: // RGBA shorthand
		_, err := fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		r, g, b, a = r*17, g*17, b*17, a*17
	case 6: // RRGGBB
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
	case 8: // RRGGBBAA
		_, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex color: %s", hex)
		}
	default:
		return Color{}, fmt.Errorf("invalid hex color length: %s", hex)
	}

	return NewRGBA(r, g, b, float64(a)/255), nil
}

// NewFromHSL creates a color from HSL values (h: 0-360, s: 0-100, l: 0-100).
func NewFromHSL(h, s, l float64) Color {
	s /= 100
	l /= 100
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return NewRGBFloat(r+m, g+m, b+m)
}

// NewFromHSV creates a color from HSV values (h: 0-360, s: 0-100, v: 0-100).
func NewFromHSV(h, s, v float64) Color {
	s /= 100
	v /= 100
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return NewRGBFloat(r+m, g+m, b+m)
}

// NewFromCMYK creates a color from CMYK values (0-100).
func NewFromCMYK(c, m, y, k float64) Color {
	c /= 100
	m /= 100
	y /= 100
	k /= 100

	r := (1 - c) * (1 - k)
	g := (1 - m) * (1 - k)
	b := (1 - y) * (1 - k)

	return NewRGBFloat(r, g, b)
}

// NewFromOKLCH creates a color from OKLCH values (l: 0-1, c: 0-0.4+, h: 0-360).
func NewFromOKLCH(l, c, h float64) Color {
	hRad := h * math.Pi / 180
	a := c * math.Cos(hRad)
	b := c * math.Sin(hRad)
	return NewFromOKLab(l, a, b)
}

// NewFromOKLab creates a color from OKLab values.
func NewFromOKLab(l, a, b float64) Color {
	// OKLab to LMS
	L_ := l + 0.3963377774*a + 0.2158037573*b
	M_ := l - 0.1055613458*a - 0.0638541728*b
	S_ := l - 0.0894841775*a - 1.2914855480*b

	L := L_ * L_ * L_
	M := M_ * M_ * M_
	S := S_ * S_ * S_

	// LMS to linear RGB
	r := +4.0767416621*L - 3.3077115913*M + 0.2309699292*S
	g := -1.2684380046*L + 2.6097574011*M - 0.3413193965*S
	bl := -0.0041960863*L - 0.7034186147*M + 1.7076147010*S

	// Linear to sRGB
	return Color{
		R: clamp01(linearToSrgb(r)),
		G: clamp01(linearToSrgb(g)),
		B: clamp01(linearToSrgb(bl)),
		A: 1,
	}
}

// Helper functions

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func srgbToLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func linearToSrgb(v float64) float64 {
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}

func labF(t float64) float64 {
	const delta = 6.0 / 29.0
	if t > delta*delta*delta {
		return math.Cbrt(t)
	}
	return t/(3*delta*delta) + 4.0/29.0
}

// Format functions for output

// FormatRGB formats as "rgb(r, g, b)".
func (c Color) FormatRGB() string {
	r, g, b := c.RGB()
	return fmt.Sprintf("rgb(%d, %d, %d)", r, g, b)
}

// FormatRGBA formats as "rgba(r, g, b, a)".
func (c Color) FormatRGBA() string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, a)
}

// FormatHSL formats as "hsl(h, s%, l%)".
func (c Color) FormatHSL() string {
	h, s, l := c.HSL()
	return fmt.Sprintf("hsl(%.1f, %.1f%%, %.1f%%)", h, s, l)
}

// FormatHSV formats as "hsv(h, s%, v%)".
func (c Color) FormatHSV() string {
	h, s, v := c.HSV()
	return fmt.Sprintf("hsv(%.1f, %.1f%%, %.1f%%)", h, s, v)
}

// FormatCMYK formats as "cmyk(c%, m%, y%, k%)".
func (c Color) FormatCMYK() string {
	cy, m, y, k := c.CMYK()
	return fmt.Sprintf("cmyk(%.1f%%, %.1f%%, %.1f%%, %.1f%%)", cy, m, y, k)
}

// FormatOKLCH formats as "oklch(l c h)".
func (c Color) FormatOKLCH() string {
	l, ch, h := c.OKLCH()
	return fmt.Sprintf("oklch(%.3f %.3f %.1f)", l, ch, h)
}

// FormatOKLab formats as "oklab(l a b)".
func (c Color) FormatOKLab() string {
	l, a, b := c.OKLab()
	return fmt.Sprintf("oklab(%.3f %.3f %.3f)", l, a, b)
}

// FormatLab formats as "lab(l a b)".
func (c Color) FormatLab() string {
	l, a, b := c.Lab()
	return fmt.Sprintf("lab(%.1f %.1f %.1f)", l, a, b)
}

// FormatAll returns all color formats.
func (c Color) FormatAll() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("hex:    %s\n", c.Hex()))
	sb.WriteString(fmt.Sprintf("rgb:    %s\n", c.FormatRGB()))
	sb.WriteString(fmt.Sprintf("hsl:    %s\n", c.FormatHSL()))
	sb.WriteString(fmt.Sprintf("hsv:    %s\n", c.FormatHSV()))
	sb.WriteString(fmt.Sprintf("cmyk:   %s\n", c.FormatCMYK()))
	sb.WriteString(fmt.Sprintf("lab:    %s\n", c.FormatLab()))
	sb.WriteString(fmt.Sprintf("oklab:  %s\n", c.FormatOKLab()))
	sb.WriteString(fmt.Sprintf("oklch:  %s", c.FormatOKLCH()))
	return sb.String()
}

// Parse attempts to parse a color from various formats.
func Parse(input string) (Color, error) {
	input = strings.TrimSpace(input)
	lower := strings.ToLower(input)

	// Try hex
	if strings.HasPrefix(lower, "#") || isHexColor(lower) {
		return NewFromHex(input)
	}

	// Try rgb/rgba
	if strings.HasPrefix(lower, "rgb") {
		return parseRGB(input)
	}

	// Try hsl/hsla
	if strings.HasPrefix(lower, "hsl") {
		return parseHSL(input)
	}

	// Try hsv/hsb
	if strings.HasPrefix(lower, "hsv") || strings.HasPrefix(lower, "hsb") {
		return parseHSV(input)
	}

	// Try cmyk
	if strings.HasPrefix(lower, "cmyk") {
		return parseCMYK(input)
	}

	// Try oklch
	if strings.HasPrefix(lower, "oklch") {
		return parseOKLCH(input)
	}

	// Try oklab
	if strings.HasPrefix(lower, "oklab") {
		return parseOKLab(input)
	}

	// Try named colors
	if c, ok := namedColors[lower]; ok {
		return NewFromHex(c)
	}

	// Try as hex without #
	if isHexColor(lower) {
		return NewFromHex(input)
	}

	return Color{}, fmt.Errorf("unable to parse color: %s", input)
}

func isHexColor(s string) bool {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 3 && len(s) != 4 && len(s) != 6 && len(s) != 8 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func parseRGB(input string) (Color, error) {
	re := regexp.MustCompile(`rgba?\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*([\d.]+))?\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid rgb format: %s", input)
	}

	r, _ := strconv.Atoi(matches[1])
	g, _ := strconv.Atoi(matches[2])
	b, _ := strconv.Atoi(matches[3])
	a := 1.0
	if matches[4] != "" {
		a, _ = strconv.ParseFloat(matches[4], 64)
	}

	return NewRGBA(r, g, b, a), nil
}

func parseHSL(input string) (Color, error) {
	re := regexp.MustCompile(`hsla?\s*\(\s*([\d.]+)\s*,\s*([\d.]+)%?\s*,\s*([\d.]+)%?\s*(?:,\s*([\d.]+))?\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid hsl format: %s", input)
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	s, _ := strconv.ParseFloat(matches[2], 64)
	l, _ := strconv.ParseFloat(matches[3], 64)

	c := NewFromHSL(h, s, l)
	if matches[4] != "" {
		c.A, _ = strconv.ParseFloat(matches[4], 64)
	}

	return c, nil
}

func parseHSV(input string) (Color, error) {
	re := regexp.MustCompile(`hs[vb]\s*\(\s*([\d.]+)\s*,\s*([\d.]+)%?\s*,\s*([\d.]+)%?\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid hsv format: %s", input)
	}

	h, _ := strconv.ParseFloat(matches[1], 64)
	s, _ := strconv.ParseFloat(matches[2], 64)
	v, _ := strconv.ParseFloat(matches[3], 64)

	return NewFromHSV(h, s, v), nil
}

func parseCMYK(input string) (Color, error) {
	re := regexp.MustCompile(`cmyk\s*\(\s*([\d.]+)%?\s*,\s*([\d.]+)%?\s*,\s*([\d.]+)%?\s*,\s*([\d.]+)%?\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid cmyk format: %s", input)
	}

	c, _ := strconv.ParseFloat(matches[1], 64)
	m, _ := strconv.ParseFloat(matches[2], 64)
	y, _ := strconv.ParseFloat(matches[3], 64)
	k, _ := strconv.ParseFloat(matches[4], 64)

	return NewFromCMYK(c, m, y, k), nil
}

func parseOKLCH(input string) (Color, error) {
	re := regexp.MustCompile(`oklch\s*\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid oklch format: %s", input)
	}

	l, _ := strconv.ParseFloat(matches[1], 64)
	c, _ := strconv.ParseFloat(matches[2], 64)
	h, _ := strconv.ParseFloat(matches[3], 64)

	return NewFromOKLCH(l, c, h), nil
}

func parseOKLab(input string) (Color, error) {
	re := regexp.MustCompile(`oklab\s*\(\s*([\d.-]+)\s+([\d.-]+)\s+([\d.-]+)\s*\)`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid oklab format: %s", input)
	}

	l, _ := strconv.ParseFloat(matches[1], 64)
	a, _ := strconv.ParseFloat(matches[2], 64)
	b, _ := strconv.ParseFloat(matches[3], 64)

	return NewFromOKLab(l, a, b), nil
}

// Common named colors
var namedColors = map[string]string{
	"black":   "#000000",
	"white":   "#ffffff",
	"red":     "#ff0000",
	"green":   "#00ff00",
	"blue":    "#0000ff",
	"yellow":  "#ffff00",
	"cyan":    "#00ffff",
	"magenta": "#ff00ff",
	"orange":  "#ffa500",
	"purple":  "#800080",
	"pink":    "#ffc0cb",
	"brown":   "#a52a2a",
	"gray":    "#808080",
	"grey":    "#808080",
	"silver":  "#c0c0c0",
	"gold":    "#ffd700",
	"navy":    "#000080",
	"teal":    "#008080",
	"olive":   "#808000",
	"maroon":  "#800000",
	"lime":    "#00ff00",
	"aqua":    "#00ffff",
	"fuchsia": "#ff00ff",
}

// NamedColorNames returns a sorted list of all named color names.
func NamedColorNames() []string {
	names := make([]string, 0, len(namedColors))
	for name := range namedColors {
		names = append(names, name)
	}
	// Sort for consistent output
	for i := 0; i < len(names)-1; i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return names
}
