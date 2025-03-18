package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
)

type LabelStyle struct {
	Padding  f32.Padding
	FontNo   int
	FontSize float32
	Color    f32.Color
}

var DefaultLabel = LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Normal,
	Color:    f32.Black,
	FontSize: 1.0,
}

var H1 = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    f32.Black,
	FontSize: 2.0,
}

var I = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Italic,
	Color:    f32.Black,
	FontSize: 0.9,
}

func Label(text string, style *LabelStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultLabel
		}
		f := font.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text)/2 + style.Padding.L + style.Padding.R
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		f.SetColor(style.Color)
		f.Printf(ctx.Rect.X+style.Padding.L, ctx.Rect.Y+baseline, style.FontSize, 0, text)
		if gpu.Debugging {
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.LightBlue)
			gpu.HorLine(ctx.Rect.X, ctx.Rect.X+width, ctx.Rect.Y+baseline, 1, f32.LightBlue)
		}
		return Dim{w: width, h: height, baseline: baseline}
	}
}
