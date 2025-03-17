package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
)

type ButtonStyle struct {
	FontSize       float32
	FontNo         int
	FontWeight     float32
	FontColor      f32.Color
	InsideColor    f32.Color
	BorderColor    f32.Color
	BorderWidth    float32
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	ShadowSize     float32
}

var OkBtn = ButtonStyle{
	FontSize:       1.5,
	FontNo:         gpu.Normal,
	InsideColor:    f32.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:    f32.Color{0, 0, 0, 1},
	FontColor:      f32.Color{0, 0, 0, 1},
	OutsidePadding: f32.Padding{5, 5, 5, 5},
	InsidePadding:  f32.Padding{15, 5, 15, 5},
	BorderWidth:    1.143,
	CornerRadius:   12,
	ShadowSize:     8,
}

var PrimaryBtn = ButtonStyle{
	FontSize:       2,
	FontNo:         gpu.Normal,
	InsideColor:    f32.Color{0.5, 0.5, 1.0, 1.0},
	BorderColor:    f32.Color{0, 0, 0, 1.0},
	FontColor:      f32.Color{1, 1, 1, 1},
	OutsidePadding: f32.Padding{5, 5, 5, 5},
	InsidePadding:  f32.Padding{12, 4, 12, 4},
	BorderWidth:    0,
	CornerRadius:   12,
}

func Button(text string, action func(), style ButtonStyle, hint string) Wid {
	return func(ctx Ctx) Dim {
		f := font.Fonts[style.FontNo]
		dho := style.OutsidePadding.T + style.OutsidePadding.B
		dhi := style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		dwi := style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth
		dwo := style.OutsidePadding.R + style.OutsidePadding.L
		height := f.Height(style.FontSize) + dho + dhi
		width := font.Fonts[style.FontNo].Width(style.FontSize, text) + dwo + dwi
		baseline := f.Baseline(style.FontSize) + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}

		ctx.Rect.W = width
		ctx.Rect.H = height
		b := style.BorderWidth
		focus.Move(action)
		// shadow := float32(0.0)
		col := style.InsideColor
		if focus.LeftMouseBtnPressed(ctx.Rect) {
			gpu.Shade(ctx.Rect.Move(0, 0), style.CornerRadius, f32.Black, 3)
			b += 1
		} else if focus.Hovered(ctx.Rect) {
			gpu.Shade(ctx.Rect.Move(2, 2), style.CornerRadius, f32.Black, 3)
		}
		if focus.LeftMouseBtnReleased(ctx.Rect) {
			focus.MouseBtnReleased = false
			focus.Set(action)
		}
		if focus.At(action) {
			b += 1
		}
		focus.AddFocusable(ctx.Rect, action)

		if focus.Hovered(ctx.Rect) {
			Hint(hint, action)
		}

		r := ctx.Rect.Inset(style.OutsidePadding)
		gpu.RoundedRect(r, style.CornerRadius, b, col, style.BorderColor)
		f.SetColor(style.FontColor)
		f.Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+ctx.Baseline,
			style.FontSize, 0, text)
		f.SetColor(f32.Black)
		return Dim{}
	}
}
