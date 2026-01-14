// Package f32 implements coordinates and colors using float32.
package f32

import (
	"math"
	"strconv"
)

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

// Abs is the absolute value of a float32
func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func CaEq(x, y, epsilon float32) bool {
	return Abs(x-y) < epsilon
}

// DiffXY is the sum of the absolute difference between points X and Y
func DiffXY(p1 Pos, p2 Pos) float32 {
	return Abs(p1.X-p2.X) + Abs(p1.Y-p2.Y)
}

// Diff is the distance between two points
//
//goland:noinspection GoUnusedExportedFunction
func Diff(p1 Pos, p2 Pos) float32 {
	d := (p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.X)*(p1.X-p2.X)
	return float32(math.Sqrt(float64(d)))
}

// Inside is true if the point p is inside the rectangle r
func (p Pos) Inside(r Rect) bool {
	return p.X >= r.X && p.X <= r.X+r.W && p.Y >= r.Y && p.Y <= r.Y+r.H
}

// Reduce will shrink the rectangle r by the amount d.
func (r Rect) Reduce(d float32) Rect {
	return Rect{r.X + d, r.Y + d, r.W - 2*d, r.H - 2*d}
}

// Square sets the rectangle's width equal to the height
func (r Rect) Square() Rect {
	return Rect{r.X, r.Y, r.H, r.H}
}

// Inset will shrink the rectangle r by the padding and the border width
func (r Rect) Inset(p Padding, bw float32) Rect {
	return Rect{r.X + p.L + bw, r.Y + p.T + bw,
		r.W - p.L - p.R - 2*bw, r.H - p.T - p.B - 2*bw}
}

// Outset will enlarge the rectangle by the padding
func (r Rect) Outset(p Padding) Rect {
	return Rect{r.X - p.L, r.Y - p.T, r.W + p.L + p.R, r.H + p.T + p.B}
}

// Increase will enlarge the rectangle r
func (r Rect) Increase(d float32) Rect {
	return Rect{r.X - d, r.Y - d, r.W + 2*d, r.H + 2*d}
}

// Move the rectangle
func (r Rect) Move(x, y float32) Rect {
	return Rect{r.X + x, r.Y + y, r.W, r.H}
}

// Pad returns a uniform padding
func Pad(pad float32) Padding {
	return Padding{pad, pad, pad, pad}
}

func (p Padding) H(bw float32) float32 {
	return p.T + p.B + 2*bw
}

func TotalPadding(p1 Padding, p2 Padding, bw float32) (float32, float32) {
	y := p1.T + p1.B + p2.T + p2.B + bw*2
	x := p1.L + p1.R + p2.L + p2.R + bw*2
	return x, y
}

// Sel will select either argument 2 or 3 depending on boolean argument 1
func Sel(condition bool, falseValue float32, trueValue float32) float32 {
	if condition {
		return trueValue
	}
	return falseValue
}

// F2S will format a float32 with db decimals
func F2S(x float32, dp int) string {
	s := strconv.FormatFloat(float64(x), 'f', dp, 32)
	return s
}

// F1 will format a float32 with 1b decimals
func F1(x float32) string {
	s := strconv.FormatFloat(float64(x), 'f', 1, 32)
	return s
}

// F2 will format a float32 with 2 decimals
func F2(x float32) string {
	s := strconv.FormatFloat(float64(x), 'f', 2, 32)
	return s
}

// F3 will format a float32 with 1b decimals and total width w
func F3(x float32) string {
	s := strconv.FormatFloat(float64(x), 'f', 3, 32)
	return s
}

// F6 will format a float32 with 1b decimals and total width w
func F6(x float32) string {
	s := strconv.FormatFloat(float64(x), 'f', 6, 32)
	return s
}

// Scale will multiply a number of *float32 by the factor
func Scale(factor float32, values ...*float32) {
	for _, x := range values {
		*x = *x * factor
	}
}

func Sum(x ...float32) float32 {
	sum := float32(0)
	for _, v := range x {
		sum += v
	}
	return sum
}
