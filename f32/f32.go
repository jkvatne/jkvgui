package f32

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
	Lightgrey   = Color{0.9, 0.9, 0.9, 1.0}
	Blue        = Color{0, 0, 1, 1}
	Red         = Color{1, 0, 0, 1}
	Green       = Color{0, 1, 0, 1}
	White       = Color{1, 1, 1, 1}
	Shade       = Color{0, 0, 0, 0.1}
)

func (p Pos) Inside(r Rect) bool {
	return p.X > r.X && p.X < r.X+r.W && p.Y > r.Y && p.Y < r.Y+r.H
}

func WithAlpha(c Color, f float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: f * c.A}
}
