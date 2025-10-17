package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/theme"
)

// Separator draws a separator.
// If dx or dy is 0, it will be the width or height of the parent.
func Separator(dx, dy float32) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: dx, H: dy, Baseline: 0}
		}
		d := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: dx, H: dy}
		if dx == 0 {
			d.W = ctx.Rect.W
		}
		if dy == 0 {
			d.H = ctx.Rect.H
		}
		// Separators do no drawing.
		return Dim{W: d.W, H: d.H}
	}
}

// Line draws a line of the given color.
// If dx or dy is 0, it will be the width or height of the parent.
func Line(dx, dy float32, color theme.UIRole) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: dx, H: dy, Baseline: 0}
		}
		d := ctx.Rect
		if dx != 0 {
			d.W = dx
		}
		if dy != 0 {
			d.H = dy
		}
		ctx.Win.Gd.SolidRect(d, color.Fg())
		return Dim{W: d.W, H: d.H}
	}
}
