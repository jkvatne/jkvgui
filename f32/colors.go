package f32

import (
	"fmt"
	"math"
)

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

var (
	Transparent = Color{}
	Black       = Color{0.0, 0.0, 0.0, 1.0}
	Grey        = Color{0.4, 0.4, 0.4, 1.0}
	LightGrey   = Color{0.9, 0.9, 0.9, 1.0}
	Blue        = Color{0.0, 0.0, 1.0, 1.0}
	LightBlue   = Color{0.8, 0.8, 1.0, 1.0}
	Red         = Color{1.0, 0.0, 0.0, 1.0}
	LightRed    = Color{1.0, 0.8, 0.8, 1.0}
	Green       = Color{0.0, 1.0, 0.0, 1.0}
	LightGreen  = Color{0.8, 1.0, 0.8, 1.0}
	White       = Color{1.0, 1.0, 1.0, 1.0}
	Yellow      = Color{1.0, 1.0, 0.0, 1.0}
	Shade       = Color{0.6, 0.6, 0.6, 0.5}
	Cyan        = Color{0.0, 1.0, 1.0, 1.0}
	Magenta     = Color{1.0, 0.0, 1.0, 1.0}
)

// WithAlpha replaces the color c's alpha with the argument a
func (c Color) WithAlpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// MultAlpha multiplies the color c's alpha with the argument a
func (c Color) MultAlpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a * c.A}
}

// Mute will dim the color c by the constant k, which must be 0..1.0
// When k=1.0 there is no change. When k=0.0 the result is gray
func (c Color) Mute(k float32) Color {
	return Color{R: 0.5 + (c.R-0.5)*k, G: 0.5 + (c.G-0.5)*k, B: 0.5 + (c.B-0.5)*k, A: c.A}
}

// FromRGB will change a 24 bit color rgb code to float32 colors (type Color).
// Typically used to translate hex codes to a Color.
func FromRGB(c uint32) Color {
	col := Color{}
	col.R = float32(c>>16&0xFF) / 255.0
	col.G = float32(c>>8&0xFF) / 255.0
	col.B = float32(c&0xFF) / 255.0
	return col
}

// Tone is the Google material tone implementation
// It keeps the hue and saturation constant, but changes lightness
// It will panic with tone<0 or tone>100
func (c Color) Tone(tone int) Color {
	ExitIf(tone < 0 || tone > 100, "Tone() called with argument <0 or >100")
	h, s, _ := c.HSL()
	return Hsl2rgb(h, s, float64(tone)/100.0)
}

// HSL is internal implementation converting RGB to HSL, HSV, or HSI.
// Basically a direct implementation of https://en.wikipedia.org/wiki/HSL_and_HSV#General_approach
func (c Color) HSL() (float64, float64, float64) {
	var h, s, lvi float64
	var huePrime float64
	r := float64(c.R)
	g := float64(c.G)
	b := float64(c.B)
	maxCol := math.Max(math.Max(r, g), b)
	minCOl := math.Min(math.Min(r, g), b)
	chroma := maxCol - minCOl
	if chroma == 0 {
		h = 0
	} else {
		if r == maxCol {
			huePrime = math.Mod((g-b)/chroma, 6)
		} else if g == maxCol {
			huePrime = ((b - r) / chroma) + 2

		} else if b == maxCol {
			huePrime = ((r - g) / chroma) + 4

		}

		h = huePrime * 60
	}
	if r == g && g == b {
		lvi = r
	} else {
		lvi = (maxCol + minCOl) / 2
	}
	if lvi == 1 {
		s = 0
	} else {
		s = chroma / (1 - math.Abs(2*lvi-1))
	}

	if math.IsNaN(s) {
		s = 0
	}

	if h < 0 {
		h = 360 + h
	}

	return h, s, lvi
}

// Hsl2rgb is internal HSV->RGB function for doing conversions using float inputs (saturation, value) and
// outputs (for R, G, and B).
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#Converting_to_RGB
func Hsl2rgb(hueDegrees float64, saturation float64, light float64) Color {
	var r, g, b float64
	hueDegrees = math.Mod(hueDegrees, 360)
	if saturation == 0 {
		r = light
		g = light
		b = light
	} else {
		var chroma float64
		var m float64
		chroma = (1 - math.Abs((2*light)-1)) * saturation
		hueSector := hueDegrees / 60
		intermediate := chroma * (1 - math.Abs(math.Mod(hueSector, 2)-1))
		switch {
		case hueSector >= 0 && hueSector <= 1:
			r = chroma
			g = intermediate
			b = 0
		case hueSector > 1 && hueSector <= 2:
			r = intermediate
			g = chroma
			b = 0
		case hueSector > 2 && hueSector <= 3:
			r = 0
			g = chroma
			b = intermediate
		case hueSector > 3 && hueSector <= 4:
			r = 0
			g = intermediate
			b = chroma
		case hueSector > 4 && hueSector <= 5:
			r = intermediate
			g = 0
			b = chroma
		case hueSector > 5 && hueSector <= 6:
			r = chroma
			g = 0
			b = intermediate
		default:
			panic(fmt.Errorf("hue sys %v yielded sector %v", hueDegrees, hueSector))
		}
		m = light - (chroma / 2)
		r += m
		g += m
		b += m
	}
	return Color{R: float32(r), G: float32(g), B: float32(b), A: 1.0}
}
