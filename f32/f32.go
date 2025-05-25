// Package f32 implements coordinates and colors using float32.
package f32

import "strconv"

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

func (p Pos) Inside(r Rect) bool {
	return p.X > r.X && p.X < r.X+r.W && p.Y > r.Y && p.Y < r.Y+r.H
}

func (p Pos) Sub(d Pos) Pos {
	return Pos{p.X - d.X, p.Y - d.Y}
}

func (r Rect) Reduce(d float32) Rect {
	return Rect{r.X + d, r.Y + d, r.W - 2*d, r.H - 2*d}
}

func (r Rect) Square() Rect {
	return Rect{r.X, r.Y, r.H, r.H}
}

func (r Rect) Inset(p Padding, bw float32) Rect {
	return Rect{r.X + p.L + bw, r.Y + p.T + bw,
		r.W - p.L - p.R - 2*bw, r.H - p.T - p.B - 2*bw}
}

func (r Rect) Outset(p Padding) Rect {
	return Rect{r.X - p.L, r.Y - p.T, r.W + p.L + p.R, r.H + p.T + p.B}
}

func (r Rect) Out(d float32) Rect {
	return Rect{r.X - d, r.Y - d, r.W + 2*d, r.H + 2*d}
}

func (r Rect) Move(x, y float32) Rect {
	return Rect{r.X + x, r.Y + y, r.W, r.H}
}

func Pad(pad float32) Padding {
	return Padding{pad, pad, pad, pad}
}

func (p Padding) IsZero() bool {
	return p.L == 0 && p.T == 0 && p.R == 0 && p.B == 0
}

func (p Padding) H(bw float32) float32 {
	return p.T + p.B + 2*bw
}

func Sel(condition bool, falseValue float32, trueValue float32) float32 {
	if condition {
		return trueValue
	}
	return falseValue
}

func F2S(x float32, dp int) string {
	s := strconv.FormatFloat(float64(x), 'f', dp, 32)
	return s
}

func Scale(fact float32, values ...*float32) {
	for _, x := range values {
		*x = *x * fact
	}
}
