package colorconv

import (
	"math"
	"testing"
)

// tolerance for floating point comparisons
const epsilon = 0.01

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// ============================================================================
// Color Constructors
// ============================================================================

func TestNewRGB(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b int
		wantR   float64
		wantG   float64
		wantB   float64
	}{
		{"black", 0, 0, 0, 0, 0, 0},
		{"white", 255, 255, 255, 1, 1, 1},
		{"red", 255, 0, 0, 1, 0, 0},
		{"green", 0, 255, 0, 0, 1, 0},
		{"blue", 0, 0, 255, 0, 0, 1},
		{"mid gray", 128, 128, 128, 0.502, 0.502, 0.502},
		{"clamped high", 300, 300, 300, 1, 1, 1},
		{"clamped low", -10, -10, -10, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewRGB(tt.r, tt.g, tt.b)
			if !floatEqual(c.R, tt.wantR) {
				t.Errorf("R = %v, want %v", c.R, tt.wantR)
			}
			if !floatEqual(c.G, tt.wantG) {
				t.Errorf("G = %v, want %v", c.G, tt.wantG)
			}
			if !floatEqual(c.B, tt.wantB) {
				t.Errorf("B = %v, want %v", c.B, tt.wantB)
			}
			if c.A != 1 {
				t.Errorf("A = %v, want 1", c.A)
			}
		})
	}
}

func TestNewRGBA(t *testing.T) {
	c := NewRGBA(255, 128, 64, 0.5)
	if c.A != 0.5 {
		t.Errorf("A = %v, want 0.5", c.A)
	}
}

func TestNewRGBFloat(t *testing.T) {
	c := NewRGBFloat(0.5, 0.5, 0.5)
	if !floatEqual(c.R, 0.5) || !floatEqual(c.G, 0.5) || !floatEqual(c.B, 0.5) {
		t.Errorf("got (%v, %v, %v), want (0.5, 0.5, 0.5)", c.R, c.G, c.B)
	}
}

// ============================================================================
// Color Output Methods
// ============================================================================

func TestColorRGB(t *testing.T) {
	c := NewRGB(255, 128, 64)
	r, g, b := c.RGB()
	if r != 255 || g != 128 || b != 64 {
		t.Errorf("RGB() = (%d, %d, %d), want (255, 128, 64)", r, g, b)
	}
}

func TestColorRGBA(t *testing.T) {
	c := NewRGBA(255, 128, 64, 0.75)
	r, g, b, a := c.RGBA()
	if r != 255 || g != 128 || b != 64 || a != 0.75 {
		t.Errorf("RGBA() = (%d, %d, %d, %v), want (255, 128, 64, 0.75)", r, g, b, a)
	}
}

func TestColorHex(t *testing.T) {
	tests := []struct {
		name string
		c    Color
		want string
	}{
		{"black", NewRGB(0, 0, 0), "#000000"},
		{"white", NewRGB(255, 255, 255), "#ffffff"},
		{"red", NewRGB(255, 0, 0), "#ff0000"},
		{"orange", NewRGB(255, 165, 0), "#ffa500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Hex(); got != tt.want {
				t.Errorf("Hex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorHexAlpha(t *testing.T) {
	c := NewRGBA(255, 0, 0, 0.5)
	got := c.HexAlpha()
	want := "#ff000080"
	if got != want {
		t.Errorf("HexAlpha() = %v, want %v", got, want)
	}
}

func TestColorHSL(t *testing.T) {
	tests := []struct {
		name         string
		c            Color
		wantH, wantS float64
		wantL        float64
	}{
		{"red", NewRGB(255, 0, 0), 0, 100, 50},
		{"green", NewRGB(0, 255, 0), 120, 100, 50},
		{"blue", NewRGB(0, 0, 255), 240, 100, 50},
		{"white", NewRGB(255, 255, 255), 0, 0, 100},
		{"black", NewRGB(0, 0, 0), 0, 0, 0},
		{"gray", NewRGB(128, 128, 128), 0, 0, 50.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, s, l := tt.c.HSL()
			if !floatEqual(h, tt.wantH) {
				t.Errorf("H = %v, want %v", h, tt.wantH)
			}
			if !floatEqual(s, tt.wantS) {
				t.Errorf("S = %v, want %v", s, tt.wantS)
			}
			if !floatEqual(l, tt.wantL) {
				t.Errorf("L = %v, want %v", l, tt.wantL)
			}
		})
	}
}

func TestColorHSV(t *testing.T) {
	tests := []struct {
		name         string
		c            Color
		wantH, wantS float64
		wantV        float64
	}{
		{"red", NewRGB(255, 0, 0), 0, 100, 100},
		{"green", NewRGB(0, 255, 0), 120, 100, 100},
		{"blue", NewRGB(0, 0, 255), 240, 100, 100},
		{"white", NewRGB(255, 255, 255), 0, 0, 100},
		{"black", NewRGB(0, 0, 0), 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, s, v := tt.c.HSV()
			if !floatEqual(h, tt.wantH) {
				t.Errorf("H = %v, want %v", h, tt.wantH)
			}
			if !floatEqual(s, tt.wantS) {
				t.Errorf("S = %v, want %v", s, tt.wantS)
			}
			if !floatEqual(v, tt.wantV) {
				t.Errorf("V = %v, want %v", v, tt.wantV)
			}
		})
	}
}

func TestColorCMYK(t *testing.T) {
	tests := []struct {
		name                string
		c                   Color
		wantC, wantM, wantY float64
		wantK               float64
	}{
		{"red", NewRGB(255, 0, 0), 0, 100, 100, 0},
		{"green", NewRGB(0, 255, 0), 100, 0, 100, 0},
		{"blue", NewRGB(0, 0, 255), 100, 100, 0, 0},
		{"white", NewRGB(255, 255, 255), 0, 0, 0, 0},
		{"black", NewRGB(0, 0, 0), 0, 0, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, m, y, k := tt.c.CMYK()
			if !floatEqual(c, tt.wantC) {
				t.Errorf("C = %v, want %v", c, tt.wantC)
			}
			if !floatEqual(m, tt.wantM) {
				t.Errorf("M = %v, want %v", m, tt.wantM)
			}
			if !floatEqual(y, tt.wantY) {
				t.Errorf("Y = %v, want %v", y, tt.wantY)
			}
			if !floatEqual(k, tt.wantK) {
				t.Errorf("K = %v, want %v", k, tt.wantK)
			}
		})
	}
}

func TestColorLab(t *testing.T) {
	// Test that Lab values are in expected ranges
	c := NewRGB(255, 128, 64)
	l, a, b := c.Lab()

	// L should be 0-100
	if l < 0 || l > 100 {
		t.Errorf("L = %v, expected 0-100", l)
	}
	// a and b typically -128 to 128
	if a < -128 || a > 128 {
		t.Errorf("a = %v, expected -128 to 128", a)
	}
	if b < -128 || b > 128 {
		t.Errorf("b = %v, expected -128 to 128", b)
	}
}

func TestColorOKLab(t *testing.T) {
	c := NewRGB(255, 128, 64)
	l, a, b := c.OKLab()

	// L should be 0-1
	if l < 0 || l > 1 {
		t.Errorf("L = %v, expected 0-1", l)
	}
	// a and b typically -0.5 to 0.5
	if a < -0.5 || a > 0.5 {
		t.Errorf("a = %v, expected -0.5 to 0.5", a)
	}
	if b < -0.5 || b > 0.5 {
		t.Errorf("b = %v, expected -0.5 to 0.5", b)
	}
}

func TestColorOKLCH(t *testing.T) {
	c := NewRGB(255, 128, 64)
	l, ch, h := c.OKLCH()

	// L should be 0-1
	if l < 0 || l > 1 {
		t.Errorf("L = %v, expected 0-1", l)
	}
	// C should be >= 0
	if ch < 0 {
		t.Errorf("C = %v, expected >= 0", ch)
	}
	// H should be 0-360
	if h < 0 || h > 360 {
		t.Errorf("H = %v, expected 0-360", h)
	}
}

// ============================================================================
// Color Constructors from other color spaces
// ============================================================================

func TestNewFromHex(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantR   int
		wantG   int
		wantB   int
		wantErr bool
	}{
		{"6 digit", "#ff5500", 255, 85, 0, false},
		{"6 digit no hash", "ff5500", 255, 85, 0, false},
		{"3 digit", "#f50", 255, 85, 0, false},
		{"8 digit with alpha", "#ff550080", 255, 85, 0, false},
		{"4 digit with alpha", "#f508", 255, 85, 0, false},
		{"invalid length", "#ff", 0, 0, 0, true},
		{"invalid chars", "#gggggg", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewFromHex(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				r, g, b := c.RGB()
				if r != tt.wantR || g != tt.wantG || b != tt.wantB {
					t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)", r, g, b, tt.wantR, tt.wantG, tt.wantB)
				}
			}
		})
	}
}

func TestNewFromHSL(t *testing.T) {
	tests := []struct {
		name  string
		h, s  float64
		l     float64
		wantR int
		wantG int
		wantB int
	}{
		{"red", 0, 100, 50, 255, 0, 0},
		{"green", 120, 100, 50, 0, 255, 0},
		{"blue", 240, 100, 50, 0, 0, 255},
		{"white", 0, 0, 100, 255, 255, 255},
		{"black", 0, 0, 0, 0, 0, 0},
		{"negative hue wraps", -60, 100, 50, 255, 0, 255}, // magenta
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFromHSL(tt.h, tt.s, tt.l)
			r, g, b := c.RGB()
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)", r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestNewFromHSV(t *testing.T) {
	tests := []struct {
		name  string
		h, s  float64
		v     float64
		wantR int
		wantG int
		wantB int
	}{
		{"red", 0, 100, 100, 255, 0, 0},
		{"green", 120, 100, 100, 0, 255, 0},
		{"blue", 240, 100, 100, 0, 0, 255},
		{"white", 0, 0, 100, 255, 255, 255},
		{"black", 0, 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFromHSV(tt.h, tt.s, tt.v)
			r, g, b := c.RGB()
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)", r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestNewFromCMYK(t *testing.T) {
	tests := []struct {
		name       string
		c, m, y, k float64
		wantR      int
		wantG      int
		wantB      int
	}{
		{"red", 0, 100, 100, 0, 255, 0, 0},
		{"green", 100, 0, 100, 0, 0, 255, 0},
		{"blue", 100, 100, 0, 0, 0, 0, 255},
		{"white", 0, 0, 0, 0, 255, 255, 255},
		{"black", 0, 0, 0, 100, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewFromCMYK(tt.c, tt.m, tt.y, tt.k)
			r, g, b := col.RGB()
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)", r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestNewFromOKLCH(t *testing.T) {
	// Test roundtrip: RGB -> OKLCH -> RGB
	original := NewRGB(255, 128, 64)
	l, c, h := original.OKLCH()
	restored := NewFromOKLCH(l, c, h)

	r1, g1, b1 := original.RGB()
	r2, g2, b2 := restored.RGB()

	// Allow small rounding differences
	if abs(r1-r2) > 1 || abs(g1-g2) > 1 || abs(b1-b2) > 1 {
		t.Errorf("Roundtrip failed: (%d, %d, %d) -> (%d, %d, %d)", r1, g1, b1, r2, g2, b2)
	}
}

func TestNewFromOKLab(t *testing.T) {
	// Test roundtrip: RGB -> OKLab -> RGB
	original := NewRGB(128, 64, 192)
	l, a, b := original.OKLab()
	restored := NewFromOKLab(l, a, b)

	r1, g1, b1 := original.RGB()
	r2, g2, b2 := restored.RGB()

	// Allow small rounding differences
	if abs(r1-r2) > 1 || abs(g1-g2) > 1 || abs(b1-b2) > 1 {
		t.Errorf("Roundtrip failed: (%d, %d, %d) -> (%d, %d, %d)", r1, g1, b1, r2, g2, b2)
	}
}

// ============================================================================
// Parse function
// ============================================================================

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantR   int
		wantG   int
		wantB   int
		wantErr bool
	}{
		{"hex with hash", "#ff5500", 255, 85, 0, false},
		{"hex without hash", "ff5500", 255, 85, 0, false},
		{"hex shorthand", "#f50", 255, 85, 0, false},
		{"rgb", "rgb(255, 128, 64)", 255, 128, 64, false},
		{"rgba", "rgba(255, 128, 64, 0.5)", 255, 128, 64, false},
		{"hsl red", "hsl(0, 100%, 50%)", 255, 0, 0, false},
		{"hsv red", "hsv(0, 100%, 100%)", 255, 0, 0, false},
		{"cmyk red", "cmyk(0, 100, 100, 0)", 255, 0, 0, false},
		{"named red", "red", 255, 0, 0, false},
		{"named blue", "blue", 0, 0, 255, false},
		{"named orange", "orange", 255, 165, 0, false},
		// oklch is tested separately due to conversion differences
		{"invalid", "notacolor", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				r, g, b := c.RGB()
				// Allow some tolerance for conversions
				if abs(r-tt.wantR) > 2 || abs(g-tt.wantG) > 2 || abs(b-tt.wantB) > 2 {
					t.Errorf("RGB() = (%d, %d, %d), want (%d, %d, %d)", r, g, b, tt.wantR, tt.wantG, tt.wantB)
				}
			}
		})
	}
}

// ============================================================================
// Format functions
// ============================================================================

func TestFormatRGB(t *testing.T) {
	c := NewRGB(255, 128, 64)
	got := c.FormatRGB()
	want := "rgb(255, 128, 64)"
	if got != want {
		t.Errorf("FormatRGB() = %v, want %v", got, want)
	}
}

func TestFormatRGBA(t *testing.T) {
	c := NewRGBA(255, 128, 64, 0.5)
	got := c.FormatRGBA()
	want := "rgba(255, 128, 64, 0.50)"
	if got != want {
		t.Errorf("FormatRGBA() = %v, want %v", got, want)
	}
}

func TestFormatHSL(t *testing.T) {
	c := NewRGB(255, 0, 0)
	got := c.FormatHSL()
	if !contains(got, "hsl(") {
		t.Errorf("FormatHSL() = %v, expected hsl(...)", got)
	}
}

func TestFormatHSV(t *testing.T) {
	c := NewRGB(255, 0, 0)
	got := c.FormatHSV()
	if !contains(got, "hsv(") {
		t.Errorf("FormatHSV() = %v, expected hsv(...)", got)
	}
}

func TestFormatCMYK(t *testing.T) {
	c := NewRGB(255, 0, 0)
	got := c.FormatCMYK()
	if !contains(got, "cmyk(") {
		t.Errorf("FormatCMYK() = %v, expected cmyk(...)", got)
	}
}

func TestFormatOKLCH(t *testing.T) {
	c := NewRGB(255, 128, 64)
	got := c.FormatOKLCH()
	if !contains(got, "oklch(") {
		t.Errorf("FormatOKLCH() = %v, expected oklch(...)", got)
	}
}

func TestFormatOKLab(t *testing.T) {
	c := NewRGB(255, 128, 64)
	got := c.FormatOKLab()
	if !contains(got, "oklab(") {
		t.Errorf("FormatOKLab() = %v, expected oklab(...)", got)
	}
}

func TestFormatLab(t *testing.T) {
	c := NewRGB(255, 128, 64)
	got := c.FormatLab()
	if !contains(got, "lab(") {
		t.Errorf("FormatLab() = %v, expected lab(...)", got)
	}
}

func TestFormatAll(t *testing.T) {
	c := NewRGB(255, 128, 64)
	got := c.FormatAll()

	// Should contain all formats
	formats := []string{"hex:", "rgb:", "hsl:", "hsv:", "cmyk:", "lab:", "oklab:", "oklch:"}
	for _, f := range formats {
		if !contains(got, f) {
			t.Errorf("FormatAll() missing %v", f)
		}
	}
}

// ============================================================================
// NamedColorNames
// ============================================================================

func TestNamedColorNames(t *testing.T) {
	names := NamedColorNames()
	if len(names) == 0 {
		t.Error("NamedColorNames() returned empty slice")
	}

	// Check that it's sorted
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Error("NamedColorNames() is not sorted")
			break
		}
	}

	// Check some expected names
	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	expected := []string{"red", "green", "blue", "white", "black"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("NamedColorNames() missing %v", e)
		}
	}
}

// ============================================================================
// Helper functions
// ============================================================================

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkNewRGB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRGB(255, 128, 64)
	}
}

func BenchmarkColorHSL(b *testing.B) {
	c := NewRGB(255, 128, 64)
	for i := 0; i < b.N; i++ {
		c.HSL()
	}
}

func BenchmarkColorOKLCH(b *testing.B) {
	c := NewRGB(255, 128, 64)
	for i := 0; i < b.N; i++ {
		c.OKLCH()
	}
}

func BenchmarkParse(b *testing.B) {
	inputs := []string{
		"#ff5500",
		"rgb(255, 128, 64)",
		"hsl(30, 100%, 50%)",
		"red",
	}
	for i := 0; i < b.N; i++ {
		Parse(inputs[i%len(inputs)])
	}
}
