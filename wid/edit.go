package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
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
	FontSize:           1.0,
	FontNo:             gpu.Normal,
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

func (s *EditStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

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
		f := font.Get(style.FontNo, style.FontColor)

		frameRect := ctx.Rect.Inset(style.OutsidePadding)
		textRect := frameRect.Inset(style.InsidePadding).Reduce(style.BorderWidth)

		fontHeight := f.Height(style.FontSize)
		baseline := f.Baseline(style.FontSize)

		if ctx.Rect.H == 0 {
			return Dim{w: textRect.W, h: fontHeight + style.TotalPaddingY(), baseline: baseline}
		}

		col := style.InsideColor
		focused := focus.At(ctx.Rect, text)

		if mouse.LeftBtnPressed(frameRect) {
			gpu.Invalidate(0)
			col.A = 1
		}
		if mouse.LeftBtnReleased(frameRect) {
			gpu.Invalidate(0)
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(text)
			s.SelStart = f.RuneNo(mouse.Pos().X-(frameRect.X), style.FontSize, s.Buffer.String())
			s.SelEnd = s.SelStart
		}
		bw := style.BorderWidth
		if focused {
			bw = min(style.BorderWidth*2, style.BorderWidth+2)
			gpu.Invalidate(111 * time.Millisecond)
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
		} else if mouse.Hovered(frameRect) {
			col.A *= 0.1
		}

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, col, style.BorderColor)
		f.SetColor(style.FontColor)
		// x := ctx.Rect.X + style.OutsidePadding.L + style.InsidePadding.L + style.BorderWidth
		f.Printf(
			textRect.X,
			textRect.Y+baseline,
			style.FontSize,
			textRect.W,
			s.Buffer.String())
		f.SetColor(f32.Black)
		s.SelEnd = s.SelStart
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(textRect.X+dx, textRect.Y, textRect.Y+textRect.H, 1, f32.Black)
		}
		return Dim{w: frameRect.W, h: frameRect.H, baseline: baseline}
	}
}
