package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

func Label(text string, size float32, p *f32.Padding, fontNo int) Wid {
	return func(ctx Ctx) Dim {
		if p == nil {
			p = &f32.Padding{5, 5, 5, 5}
		}
		height := gpu.Fonts[fontNo].Height(size) + p.T + p.B
		width := gpu.Fonts[fontNo].Width(size, text)/2 + p.L + p.R
		baseline := gpu.Fonts[fontNo].Ascent*size/2 + p.T
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		} else {
			r := ctx.Rect
			gpu.Rect(r, 1, f32.Lightgrey, f32.LightBlue)
			gpu.HorLine(ctx.Rect.X, ctx.Rect.X+width, ctx.Rect.Y+baseline, 1, f32.LightBlue)
			gpu.Fonts[fontNo].SetColor(f32.Black)
			gpu.Fonts[fontNo].Printf(ctx.Rect.X+p.L, ctx.Rect.Y+p.T+baseline, size, 0, text)
			return Dim{w: width, h: height, baseline: baseline}
		}
	}
}
