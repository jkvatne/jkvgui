package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

// Separator draws a separator.
// If dx or dy is 0, it will be the width or height of the parent.
func Separator(dx, dy float32) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: dx, H: dy, Baseline: 0}
		}
		d := f32.Rect{ctx.Rect.X, ctx.Rect.Y, dx, dy}
		if dx == 0 {
			d.W = ctx.Rect.W
		}
		if dy == 0 {
			d.H = ctx.Rect.H
		}
		// Separators do no drawing.
		return Dim{d.W, d.H, 0}
	}
}

// Line draws a line of the given color.
// If dx or dy is 0, it will be the width or height of the parent.
func Line(dx, dy float32, color theme.UIRole) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: dx, H: dy, Baseline: 0}
		}
		d := f32.Rect{ctx.Rect.X, ctx.Rect.Y, dx, dy}
		if dx == 0 {
			d.W = ctx.Rect.W
		}
		if dy == 0 {
			d.H = ctx.Rect.H
		}
		gpu.Rect(ctx.Rect, dy, color.Fg(), color.Fg())
		return Dim{d.W, d.H, 0}
	}
}
