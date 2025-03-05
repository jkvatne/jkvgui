package wid

import "github.com/jkvatne/jkvgui/gpu"

type EditStyle struct {
	FontSize           float32
	FontNo             int
	FontColor          gpu.Color
	InsideColor        gpu.Color
	BorderColor        gpu.Color
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      Padding
	OutsidePadding     Padding
	CursorWidth        float32
}

var DefaultEdit = EditStyle{
	FontSize:           16,
	FontNo:             0,
	InsideColor:        gpu.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:        gpu.Color{0, 0, 0, 1},
	FontColor:          gpu.Color{0, 0, 0, 1},
	OutsidePadding:     Padding{5, 5, 5, 5},
	InsidePadding:      Padding{8, 5, 5, 5},
	BorderWidth:        1,
	BorderCornerRadius: 5,
	CursorWidth:        1,
}

type EditState struct {
	SelStart int
	SelEnd   int
	Buffer   string
}

var StateMap = make(map[*string]*EditState)

func Edit(text *string, size int, action func(), style EditStyle) Wid {
	return func(ctx Ctx) Dim {
		s := StateMap[text]
		if s == nil {
			StateMap[text] = &EditState{}
			s = StateMap[text]
			s.Buffer = *text
		}
		scale := style.FontSize / gpu.InitialSize
		dho := style.OutsidePadding.T + style.OutsidePadding.B
		dhi := style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		dwi := style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth
		dwo := style.OutsidePadding.R + style.OutsidePadding.L
		height := (gpu.Fonts[style.FontNo].Ascent+gpu.Fonts[style.FontNo].Descent)*scale + dho + dhi
		width := float32(size)*gpu.Fonts[style.FontNo].Width(style.FontSize, "n")/gpu.InitialSize + dwo + dwi
		baseline := gpu.Fonts[style.FontNo].Ascent*scale + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
		s.Buffer = *text
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		col := style.InsideColor
		if gpu.LeftMouseBtnPressed(ctx.Rect) {
			col.A = 1
		} else if gpu.LeftMouseBtnReleased(ctx.Rect) {
			gpu.MouseBtnReleased = false
			gpu.SetFocus(action)
		} else if gpu.Focused(action) {
			col.A *= 0.3
			if gpu.MoveFocusToNext {
				gpu.FocusToNext = true
				gpu.MoveFocusToNext = false
			}
			if gpu.LastRune != 0 {
				*text = *text + string(gpu.LastRune)
				gpu.LastRune = 0
			}
		} else if gpu.Hovered(ctx.Rect) {
			col.A *= 0.1
		}

		gpu.RoundedRect(
			ctx.Rect.X+style.OutsidePadding.L,
			ctx.Rect.Y+style.OutsidePadding.T,
			width-style.OutsidePadding.L-style.OutsidePadding.R,
			height-style.OutsidePadding.T-style.OutsidePadding.B,
			style.BorderCornerRadius, style.BorderWidth, style.InsideColor, style.BorderColor)
		gpu.Fonts[style.FontNo].SetColor(style.FontColor)
		gpu.Fonts[style.FontNo].Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+baseline,
			style.FontSize, *text)
		return Dim{w: width, h: height, baseline: baseline}
	}
}
