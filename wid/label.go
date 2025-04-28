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
	Color     theme.UIRole
	Align     Alignment
	Multiline bool
	Width     float32
}

var DefaultLabel = LabelStyle{
	Padding: f32.Padding{2, 2, 1, 1},
	FontNo:  gpu.Normal14,
	Color:   theme.OnSurface,
}

var C = &LabelStyle{
	Padding: f32.Padding{2, 2, 1, 1},
	FontNo:  gpu.Normal14,
	Align:   AlignCenter,
	Color:   theme.OnSurface,
}

var H1C = &LabelStyle{
	Padding: f32.Padding{2, 3, 1, 2},
	FontNo:  gpu.Bold20,
	Color:   theme.OnSurface,
	Align:   AlignCenter,
}
var H1R = &LabelStyle{
	Padding: f32.Padding{2, 3, 1, 2},
	FontNo:  gpu.Bold20,
	Color:   theme.OnSurface,
	Align:   AlignRight,
}
var H1L = &LabelStyle{
	Padding: f32.Padding{2, 3, 1, 2},
	FontNo:  gpu.Bold20,
	Color:   theme.OnSurface,
	Align:   AlignLeft,
}
var H2C = &LabelStyle{
	Padding: f32.Padding{2, 3, 1, 2},
	FontNo:  gpu.Bold16,
	Color:   theme.OnSurface,
	Align:   AlignCenter,
}
var H2R = &LabelStyle{
	Padding: f32.Padding{2, 3, 1, 2},
	FontNo:  gpu.Bold16,
	Color:   theme.OnSurface,
	Align:   AlignRight,
}
var I = &LabelStyle{
	Padding: f32.Padding{5, 3, 1, 2},
	FontNo:  gpu.Italic14,
	Color:   theme.OnSurface,
}

func Label(text string, style *LabelStyle) Wid {
	if style == nil {
		style = &DefaultLabel
	}
	f := font.Fonts[style.FontNo]
	lineHeight := f.Height()
	return func(ctx Ctx) Dim {
		var lines []string
		if style.Multiline {
			lines = font.Split(text, ctx.Rect.W, f)
		} else {
			lines = append(lines, text)
		}
		height := lineHeight*float32(len(lines)) + style.Padding.T + style.Padding.B
		width := f.Width(text) + style.Padding.L + style.Padding.R
		if style.Multiline {
			width = ctx.Rect.W
		}
		baseline := f.Baseline() + style.Padding.T
		if ctx.Mode != RenderChildren {
			if style.Width > 0.0 {
				return Dim{W: style.Width, H: height, Baseline: baseline}
			} else {
				return Dim{W: width, H: height, Baseline: baseline}
			}
		}

		baseline = max(ctx.Baseline, baseline)
		for i, line := range lines {
			if style.Align == AlignCenter {
				f.DrawText(
					ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width)/2,
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), 0, gpu.LTR, line)
			} else if style.Align == AlignRight {
				f.DrawText(
					ctx.Rect.X+style.Padding.L+(ctx.Rect.W-width),
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), 0, gpu.LTR, line)
			} else if style.Align == AlignLeft {
				f.DrawText(
					ctx.Rect.X+style.Padding.L,
					ctx.Rect.Y+baseline+float32(i)*lineHeight,
					style.Color.Fg(), 0, gpu.LTR, line)
			} else {
				panic("Alignment out of range")
			}
		}
		if gpu.DebugWidgets {
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.Blue)
			gpu.HorLine(ctx.Rect.X, ctx.Rect.X+width, ctx.Rect.Y+baseline, 1, f32.Blue)
		}
		return Dim{W: width, H: height, Baseline: baseline}
	}
}
