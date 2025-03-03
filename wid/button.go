package wid

import "github.com/jkvatne/jkvgui/gpu"

type ButtonStyle struct {
	FontSize           float32
	FontNo             int
	FontWeight         float32
	FontColor          gpu.Color
	InsideColor        gpu.Color
	BorderColor        gpu.Color
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      Padding
	OutsidePadding     Padding
}

var OkBtn = ButtonStyle{
	FontSize:           32,
	FontNo:             0,
	InsideColor:        gpu.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:        gpu.Color{0, 0, 0, 1},
	FontColor:          gpu.Color{0, 0, 0, 1},
	OutsidePadding:     Padding{5, 5, 5, 5},
	InsidePadding:      Padding{15, 5, 15, 5},
	BorderWidth:        2,
	BorderCornerRadius: 7,
}

func Button(text string, action func(), style ButtonStyle) Wid {
	return func(ctx Ctx) Dim {
		scale := style.FontSize / gpu.InitialSize
		height := (gpu.Fonts[style.FontNo].Ascent+gpu.Fonts[style.FontNo].Descent)*scale +
			style.OutsidePadding.T + style.OutsidePadding.B + style.InsidePadding.T + style.InsidePadding.B + style.BorderWidth*2
		width := gpu.Fonts[style.FontNo].Width(style.FontSize, text)/gpu.InitialSize +
			style.OutsidePadding.R + style.OutsidePadding.L + style.InsidePadding.L + style.InsidePadding.R + style.BorderWidth*2
		baseline := gpu.Fonts[style.FontNo].Ascent*scale + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}

		ctx.Rect.W = width
		ctx.Rect.H = height
		if FocusToPrevious {
			InFocus = LastFocusable
		}
		LastFocusable = nil
		if FocusToNext {
			FocusToNext = false
			InFocus = action
		}
		col := style.InsideColor
		if Pressed(ctx.Rect) {
			col.A = 1
		} else if Released(ctx.Rect) {
			MouseBtnReleased = false
			InFocus = action
		} else if Hovered(ctx.Rect) {
			col.A *= 0.1
		} else if Focused(action) {
			col.A *= 0.3
			if gpu.MoveFocusToNext {
				FocusToNext = true
				gpu.MoveFocusToNext = false
			}
		}
		Clickables = append(Clickables, Clickable{Rect: ctx.Rect, Action: action})

		gpu.RoundedRect(
			ctx.Rect.X+style.OutsidePadding.L,
			ctx.Rect.Y+style.OutsidePadding.T,
			ctx.Rect.W-style.OutsidePadding.L-style.OutsidePadding.R,
			ctx.Rect.H-style.OutsidePadding.T-style.OutsidePadding.B,
			style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor)
		gpu.Fonts[style.FontNo].SetColor(style.FontColor)
		gpu.Fonts[style.FontNo].Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+ctx.Baseline,
			style.FontSize, text)

		return Dim{}
	}
}
