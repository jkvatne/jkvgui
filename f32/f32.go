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

type Pos struct {
	X float32
	Y float32
}

type RRect struct {
	X, Y, W, H, RR float32
}
type Rect struct {
	X, Y, W, H float32
}
type Padding struct {
	L, T, R, B float32
}

var (
	Transparent = Color{}
	Black       = Color{0, 0, 0, 1}
	Grey        = Color{0.4, 0.4, 0.4, 1}
	LightGrey   = Color{0.9, 0.9, 0.9, 1}
	Shadow      = Color{0.3, 0.3, 0.3, 0.2}
	Blue        = Color{0, 0, 1, 1}
	LightBlue   = Color{0.8, 0.8, 1.0, 1.0}
	Red         = Color{1, 0, 0, 1}
	Green       = Color{0, 1, 0, 1}
	White       = Color{1, 1, 1, 1}
	Yellow      = Color{1, 1, 0, 1}
	Shade       = Color{0.4, 0.4, 0.4, 0.5}
)

func (p Pos) Inside(r Rect) bool {
	return p.X > r.X && p.X < r.X+r.W && p.Y > r.Y && p.Y < r.Y+r.H
}

func (p Pos) Sub(d Pos) Pos {
	return Pos{p.X - d.X, p.Y - d.Y}
}

func WithAlpha(c Color, f float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: f * c.A}
}

func MultAlpha(c Color, f float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: f * c.A}
}

func (r Rect) Reduce(d float32) Rect {
	return Rect{r.X + d, r.Y + d, r.W - 2*d, r.H - 2*d}
}

func (r Rect) Inset(p Padding) Rect {
	return Rect{r.X + p.L, r.Y + p.R, r.W - p.L - p.R, r.H - p.T - p.B}
}

func (r Rect) Outset(p Padding) Rect {
	return Rect{r.X - p.L, r.Y - p.R, r.W + p.L + p.R, r.H + p.T + p.B}
}

func (r Rect) Move(x, y float32) Rect {
	return Rect{r.X + x, r.Y + y, r.W, r.H}
}

func FromRGB(c uint32) Color {
	col := Color{}
	col.R = float32(c>>16&0xFF) / 255.0
	col.G = float32(c>>8&0xFF) / 255.0
	col.B = float32(c&0xFF) / 255.0
	return col
}

func Emphasis(c Color) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// Tone is the Google material tone implementation
func (c Color) Tone(tone int) Color {
	h, s, _ := c.HSL()
	return Hsl2rgb(h, s, float64(tone)/100.0)
}

func (c Color) Alpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: a}
}

// Rgb2hsl is internal implementation converting RGB to HSL, HSV, or HSI.
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#General_approach
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
			panic(fmt.Errorf("hue input %v yielded sector %v", hueDegrees, hueSector))
		}
		m = light - (chroma / 2)
		r += m
		g += m
		b += m
	}
	return Color{R: float32(r), G: float32(g), B: float32(b), A: 1.0}
}
