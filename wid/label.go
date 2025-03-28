package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

type LabelStyle struct {
	Padding  f32.Padding
	FontNo   int
	FontSize float32
	Color    theme.UIRole
	Align    Alignment
}

var DefaultLabel = LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Normal,
	Color:    theme.OnSurface,
	FontSize: 1.0,
}

var H1C = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 2.0,
	Align:    AlignCenter,
}
var H1R = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 2.0,
	Align:    AlignRight,
}
var H1L = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 2.0,
	Align:    AlignLeft,
}
var H2C = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 1.5,
	Align:    AlignCenter,
}
var H2R = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 1.5,
	Align:    AlignRight,
}
var Center = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Bold,
	Color:    theme.OnSurface,
	FontSize: 1.0,
	Align:    AlignCenter,
}
var I = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Italic,
	Color:    theme.OnSurface,
	FontSize: 0.9,
}

func Label(text string, style *LabelStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultLabel
		}
		f := font.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text) + style.Padding.L + style.Padding.R
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		if ctx.Rect.H == 0 {
			return Dim{W: width, H: height, Baseline: baseline}
		}
		baseline = max(ctx.Baseline, baseline)
		f.SetColor(theme.Colors[style.Color])
		if style.Align == AlignCenter {
			f.Printf(ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width)/2, ctx.Rect.Y+baseline, style.FontSize, 0, text)
		} else if style.Align == AlignRight {
			f.Printf(ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width), ctx.Rect.Y+baseline, style.FontSize, 0, text)
		} else if style.Align == AlignLeft {
			f.Printf(ctx.Rect.X+style.Padding.L, ctx.Rect.Y+baseline, style.FontSize, 0, text)
		} else {
			panic("Alignment out of range")
		}
		if gpu.Debugging {
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.LightBlue)
			gpu.HorLine(ctx.Rect.X, ctx.Rect.X+width, ctx.Rect.Y+baseline, 1, f32.LightBlue)
		}
		return Dim{W: width, H: height, Baseline: baseline}
	}
}
