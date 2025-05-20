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
	Role      theme.UIRole
	Align     Alignment
	Multiline bool
	Width     float32
	Height    float32
}

var DefaultLabel = LabelStyle{
	Padding: f32.Padding{L: 2, T: 2, R: 1, B: 1},
	FontNo:  gpu.Normal12,
	Role:    theme.OnSurface,
}

var C = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 2, R: 1, B: 1},
	FontNo:  gpu.Normal12,
	Align:   AlignCenter,
	Role:    theme.OnSurface,
}

var H1C = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 3, R: 1, B: 2},
	FontNo:  gpu.Bold20,
	Role:    theme.OnSurface,
	Align:   AlignCenter,
}
var H1R = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 3, R: 1, B: 2},
	FontNo:  gpu.Bold20,
	Role:    theme.OnSurface,
	Align:   AlignRight,
}
var H1L = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 3, R: 1, B: 2},
	FontNo:  gpu.Bold20,
	Role:    theme.OnSurface,
	Align:   AlignLeft,
}
var H2C = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 3, R: 1, B: 2},
	FontNo:  gpu.Bold16,
	Role:    theme.OnSurface,
	Align:   AlignCenter,
}
var H2R = &LabelStyle{
	Padding: f32.Padding{L: 2, T: 3, R: 1, B: 2},
	FontNo:  gpu.Bold16,
	Role:    theme.OnSurface,
	Align:   AlignRight,
}
var I = &LabelStyle{
	Padding: f32.Padding{L: 5, T: 3, R: 1, B: 2},
	FontNo:  gpu.Italic12,
	Role:    theme.OnSurface,
}

func (l *LabelStyle) R(r theme.UIRole) *LabelStyle {
	ll := *l
	ll.Role = r
	return &ll
}

// BoxText will display a colored box with colored text inside it.
// The text will be centered inside the box.
func BoxText(text string, fg f32.Color, bg f32.Color, style *LabelStyle) Wid {
	f := font.Fonts[style.FontNo]
	baseline := f.Baseline() + (style.Height-f.Height())/2
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: style.Height, Baseline: baseline}
		}
		gpu.RoundedRect(ctx.Rect, 0, 0, bg, bg)
		f.DrawText(ctx.Rect.X+(ctx.Rect.W-f.Width(text))/2, ctx.Rect.Y+baseline, fg, 0, gpu.LTR, text)
		return Dim{W: style.Width, H: style.Height, Baseline: baseline}
	}
}

// Label will display a possibly multilied text
// Padding and alignment can be specified in the style.
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
			y := ctx.Rect.Y + baseline + float32(i)*lineHeight
			x := ctx.Rect.X
			if style.Align == AlignCenter {
				x += (ctx.Rect.W - width) / 2
			} else if style.Align == AlignRight {
				x += ctx.Rect.W - width + style.Padding.L
			} else if style.Align == AlignLeft {
				x += style.Padding.L
			} else {
				panic("Alignment out of range")
			}
			f.DrawText(x, y, style.Role.Fg(), 0, gpu.LTR, line)
			if *gpu.DebugWidgets {
				gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.Blue)
				gpu.HorLine(x, x+width, y, 1, f32.Blue)
			}
		}
		return Dim{W: width, H: height, Baseline: baseline}
	}
}
