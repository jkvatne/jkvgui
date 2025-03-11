package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

type EditStyle struct {
	FontSize           float32
	FontNo             int
	FontColor          f32.Color
	InsideColor        f32.Color
	BorderColor        f32.Color
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	CursorWidth        float32
}

var DefaultEdit = EditStyle{
	FontSize:           12,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 1},
	FontColor:          f32.Color{0, 0, 0, 1},
	OutsidePadding:     f32.Padding{5, 5, 5, 5},
	InsidePadding:      f32.Padding{8, 5, 5, 5},
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

func Edit(text *string, action func(), style *EditStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultEdit
		}
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
		innerWidth := gpu.Fonts[style.FontNo].Width(style.FontSize, "n") * scale
		width := innerWidth + dwo + dwi
		baseline := gpu.Fonts[style.FontNo].Ascent*scale + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
		s.Buffer = *text
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		col := style.InsideColor
		outline := ctx.Rect
		outline.W = width
		outline.H = height
		gpu.MoveFocus(text)
		if gpu.LeftMouseBtnPressed(outline) {
			col.A = 1
		} else if gpu.LeftMouseBtnReleased(outline) {
			gpu.MouseBtnReleased = false
			gpu.SetFocus(text)
		} else if gpu.Focused(text) {
			col.A *= 0.3
			if gpu.MoveFocusToNext {
				gpu.FocusToNext = true
				gpu.MoveFocusToNext = false
			}
			if gpu.LastRune != 0 {
				*text = *text + string(gpu.LastRune)
				gpu.LastRune = 0
			}
			if gpu.Backspace {
				str := *text
				*text = str[0 : len(str)-1]
				gpu.Backspace = false
			}
		} else if gpu.Hovered(outline) {
			col.A *= 0.1
		}

		gpu.RoundedRect(
			ctx.Rect.X+style.OutsidePadding.L,
			ctx.Rect.Y+style.OutsidePadding.T,
			width-style.OutsidePadding.L-style.OutsidePadding.R,
			height-style.OutsidePadding.T-style.OutsidePadding.B,
			style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor, 5, 0)
		gpu.Fonts[style.FontNo].SetColor(style.FontColor)
		gpu.Fonts[style.FontNo].Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+baseline,
			style.FontSize, innerWidth, *text)
		return Dim{w: width, h: height, baseline: baseline}
	}
}

const ellipsis = string(rune(0x2026))
