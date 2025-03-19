package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

func Separator(dx, dy float32, color theme.UIRole) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Rect.H == 0 {
			return Dim{W: dx, H: dy, baseline: 0}
		}
		d := f32.Rect{ctx.Rect.X, ctx.Rect.Y, dx, dy}
		if dx == 0 {
			d.W = ctx.Rect.W
		}
		if dy == 0 {
			d.H = ctx.Rect.H
		}
		col := theme.Colors[color]
		gpu.Rect(d, 0, col, col)
		return Dim{d.W, d.H, 0}
	}
}
