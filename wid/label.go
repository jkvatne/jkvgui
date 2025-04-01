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
	Padding   f32.Padding
	FontNo    int
	FontSize  float32
	Color     theme.UIRole
	Align     Alignment
	Multiline bool
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
		lineHeight := f.Height(style.FontSize)
		var lines []string
		if style.Multiline {
			lines = font.Split(text, ctx.Rect.W, f, style.FontSize)
		} else {
			lines = append(lines, text)
		}
		height := lineHeight*float32(len(lines)) + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text) + style.Padding.L + style.Padding.R
		if style.Multiline {
			width = 20
		}
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		if ctx.Rect.H == 0 {
			return Dim{W: width, H: height, Baseline: baseline}
		}

		baseline = max(ctx.Baseline, baseline)
		for i, line := range lines {
			if style.Align == AlignCenter {
				f.DrawText(
					ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width)/2,
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), style.FontSize, 0, gpu.LeftToRight, line)
			} else if style.Align == AlignRight {
				f.DrawText(
					ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width),
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), style.FontSize, 0, gpu.LeftToRight, line)
			} else if style.Align == AlignLeft {
				f.DrawText(
					ctx.Rect.X+style.Padding.L,
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), style.FontSize, 0, gpu.LeftToRight, line)
			} else {
				panic("Alignment out of range")
			}
		}
		if gpu.Debugging {
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.LightBlue)
			gpu.HorLine(ctx.Rect.X, ctx.Rect.X+width, ctx.Rect.Y+baseline, 1, f32.LightBlue)
		}
		return Dim{W: width, H: height, Baseline: baseline}
	}
}
