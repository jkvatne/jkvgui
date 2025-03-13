package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"time"
)

const Ellipsis = string(rune(0x2026))

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
	FontSize:           1,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{0.9, 0.9, 0.9, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 1},
	FontColor:          f32.Color{0, 0, 0, 1},
	OutsidePadding:     f32.Padding{3, 3, 3, 3},
	InsidePadding:      f32.Padding{5, 3, 1, 3},
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
		r := ctx.Rect.Inset(style.OutsidePadding)
		dho := style.OutsidePadding.T + style.OutsidePadding.B
		dhi := style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		dwi := style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth
		dwo := style.OutsidePadding.R + style.OutsidePadding.L
		fh := gpu.Fonts[style.FontNo].Height(style.FontSize)
		height := fh + dho + dhi
		width := r.W + dwo
		baseline := gpu.Fonts[style.FontNo].Baseline(style.FontSize) + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
		s.Buffer = *text
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		col := style.InsideColor
		outline := ctx.Rect
		outline.W = width
		outline.H = height
		focused := gpu.Focused(text)
		gpu.MoveFocus(text)
		if gpu.LeftMouseBtnPressed(outline) {
			col.A = 1
		} else if gpu.LeftMouseBtnReleased(outline) {
			gpu.MouseBtnReleased = false
			gpu.SetFocus(text)
		} else if focused {
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
		f := gpu.Fonts[style.FontNo]
		gpu.RoundedRect(r, style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor, 5, 0)
		f.SetColor(style.FontColor)
		f.Printf(
			ctx.Rect.X+style.OutsidePadding.L+style.InsidePadding.L+style.BorderWidth,
			ctx.Rect.Y+baseline,
			style.FontSize,
			r.W-dwi-style.BorderWidth*2-fh,
			*text)
		f.SetColor(f32.Black)
		s.SelStart = 6
		s.SelEnd = s.SelStart
		if s.SelStart > 0 && focused {
			k := time.Now().UnixMilli() / 500
			if k&1 == 0 {
				dx := f.Width(style.FontSize, (*text)[0:s.SelStart])
				gpu.VertLine(r.X+dx, r.Y+style.InsidePadding.T, r.Y+baseline, 1, f32.Black)
			}
		}
		return Dim{w: width, h: height, baseline: baseline}
	}
}
