package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

type ButtonStyle struct {
	FontSize           float32
	FontNo             int
	FontWeight         float32
	FontColor          f32.Color
	InsideColor        f32.Color
	BorderColor        f32.Color
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	ShadowSize         float32
}

var OkBtn = ButtonStyle{
	FontSize:           2,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 1},
	FontColor:          f32.Color{0, 0, 0, 1},
	OutsidePadding:     f32.Padding{5, 5, 5, 5},
	InsidePadding:      f32.Padding{15, 5, 15, 5},
	BorderWidth:        1.143,
	BorderCornerRadius: 12,
	ShadowSize:         8,
}

var PrimaryBtn = ButtonStyle{
	FontSize:           2,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{0.5, 0.5, 1.0, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 0},
	FontColor:          f32.Color{1, 1, 1, 1},
	OutsidePadding:     f32.Padding{5, 5, 5, 5},
	InsidePadding:      f32.Padding{12, 4, 12, 4},
	BorderWidth:        0,
	BorderCornerRadius: 12,
}

func Button(text string, action func(), style ButtonStyle, hint string) Wid {
	return func(ctx Ctx) Dim {
		scale := style.FontSize / 2
		dho := style.OutsidePadding.T + style.OutsidePadding.B
		dhi := style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		dwi := style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth
		dwo := style.OutsidePadding.R + style.OutsidePadding.L
		height := (gpu.Fonts[style.FontNo].Ascent+gpu.Fonts[style.FontNo].Descent)*scale + dho + dhi
		width := gpu.Fonts[style.FontNo].Width(scale, text) + dwo + dwi
		baseline := gpu.Fonts[style.FontNo].Ascent*scale + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}

		ctx.Rect.W = width
		ctx.Rect.H = height

		gpu.MoveFocus(action)
		shadow := float32(0.0)
		col := style.InsideColor
		if gpu.LeftMouseBtnPressed(ctx.Rect) {
			col.A = 1
		} else if gpu.LeftMouseBtnReleased(ctx.Rect) {
			gpu.MouseBtnReleased = false
			gpu.SetFocus(action)
		} else if gpu.Focused(action) {
			col.A *= 0.3
			shadow = float32(1.0)
			if gpu.MoveFocusToNext {
				gpu.FocusToNext = true
				gpu.MoveFocusToNext = false
			}

		} else if gpu.Hovered(ctx.Rect) {
			col.A *= 0.1

		}
		gpu.AddFocusable(ctx.Rect, action)

		if gpu.Hovered(ctx.Rect) {
			Hint(hint, action)
		}

		r := ctx.Rect.Inset(style.OutsidePadding)
		gpu.RoundedRect(r, style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor, style.ShadowSize, shadow)
		gpu.Fonts[style.FontNo].SetColor(style.FontColor)
		gpu.Fonts[style.FontNo].Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+ctx.Baseline,
			style.FontSize, 0, text)
		gpu.Fonts[style.FontNo].SetColor(f32.Black)
		return Dim{}
	}
}
