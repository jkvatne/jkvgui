package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	utf8 "golang.org/x/exp/utf8string"
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
	FontSize:           1.5,
	FontNo:             gpu.DefaultFont,
	InsideColor:        f32.Color{1.0, 1.0, 1.0, 1.0},
	BorderColor:        f32.Color{0, 0, 0, 1},
	FontColor:          f32.Color{0, 0, 0, 1},
	OutsidePadding:     f32.Padding{4, 4, 4, 4},
	InsidePadding:      f32.Padding{5, 2, 2, 2},
	BorderWidth:        1,
	BorderCornerRadius: 5,
	CursorWidth:        1.5,
}

type EditState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
}

var (
	StateMap = make(map[*string]*EditState)
	halfUnit int64
)

func Edit(text *string, action func(), style *EditStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultEdit
		}
		s := StateMap[text]
		if s == nil {
			StateMap[text] = &EditState{}
			s = StateMap[text]
			s.Buffer.Init(*text)
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
		if ctx.Rect.H == 0 {
			return Dim{w: width, h: height, baseline: baseline}
		}
		col := style.InsideColor
		outline := ctx.Rect
		outline.W = width
		outline.H = height
		focused := gpu.Focused(text)
		f := gpu.Fonts[style.FontNo]
		gpu.MoveFocus(text)
		if gpu.LeftMouseBtnPressed(outline) {
			gpu.Invalidate(0)
			col.A = 1
		}
		if gpu.LeftMouseBtnReleased(outline) {
			gpu.Invalidate(0)
			gpu.MouseBtnReleased = false
			halfUnit = time.Now().UnixMilli() % 333
			gpu.SetFocus(text)
			s.SelStart = f.RuneNo(gpu.MousePos.X-(r.X), style.FontSize, s.Buffer.String())
			s.SelEnd = s.SelStart
		}
		if focused {
			col.A *= 0.3
			gpu.Invalidate(111 * time.Millisecond)
			if gpu.MoveFocusToNext {
				gpu.FocusToNext = true
				gpu.MoveFocusToNext = false
			}
			if gpu.LastRune != 0 {
				s1 := s.Buffer.Slice(0, s.SelStart)
				s2 := s.Buffer.Slice(s.SelEnd, s.Buffer.RuneCount())
				s.Buffer.Init(s1 + string(gpu.LastRune) + s2)
				gpu.LastRune = 0
				s.SelStart++
				s.SelEnd++
			}
			if gpu.LastKey == glfw.KeyBackspace {
				str := *text
				*text = str[0 : len(str)-1]
			} else if gpu.LastKey == glfw.KeyLeft {
				s.SelStart--
				s.SelStart = max(0, s.SelStart)
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyRight {
				s.SelStart++
				s.SelStart = min(s.SelStart, s.Buffer.RuneCount())
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyEnd {
				s.SelStart = s.Buffer.RuneCount()
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyHome {
				s.SelStart = 0
				s.SelEnd = s.SelStart
			}
		} else if gpu.Hovered(outline) {
			col.A *= 0.1
		}

		gpu.RoundedRect(r, style.BorderCornerRadius, style.BorderWidth, col, style.BorderColor, 5, 0)
		f.SetColor(style.FontColor)
		x := ctx.Rect.X + style.OutsidePadding.L + style.InsidePadding.L + style.BorderWidth
		f.Printf(
			x,
			ctx.Rect.Y+baseline,
			style.FontSize,
			r.W-dwi-style.BorderWidth*2-fh,
			s.Buffer.String())
		f.SetColor(f32.Black)
		s.SelEnd = s.SelStart
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(x+dx, r.Y+style.InsidePadding.T, r.Y+baseline, 1, f32.Black)
		}
		return Dim{w: width, h: height, baseline: baseline}
	}
}
