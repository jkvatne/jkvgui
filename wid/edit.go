package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

const Ellipsis = string(rune(0x2026))

type EditStyle struct {
	FontSize           float32
	FontNo             int
	FontColor          theme.UIRole
	InsideColor        theme.UIRole
	BorderColor        theme.UIRole
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	CursorWidth        float32
}

var DefaultEdit = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	InsideColor:        theme.Surface,
	BorderColor:        theme.Outline,
	FontColor:          theme.OnSurface,
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
		f := font.Get(style.FontNo, theme.Colors[style.FontColor])

		frameRect := ctx.Rect.Inset(style.OutsidePadding)
		textRect := frameRect.Inset(style.InsidePadding).Reduce(style.BorderWidth)

		fontHeight := f.Height(style.FontSize)
		baseline := f.Baseline(style.FontSize)

		if ctx.Rect.H == 0 {
			return Dim{W: textRect.W, H: fontHeight + style.TotalPaddingY(), baseline: baseline}
		}

		bg := theme.Colors[style.InsideColor]
		focused := focus.At(ctx.Rect, text)

		if mouse.LeftBtnPressed(frameRect) {
			gpu.Invalidate(0)
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
				if s.SelStart > 0 {
					s.SelStart = max(0, s.SelStart-1)
					s.SelEnd = s.SelStart
					s1 := s.Buffer.Slice(0, max(s.SelStart, 0))
					s2 := s.Buffer.Slice(s.SelEnd+1, s.Buffer.RuneCount())
					s.Buffer.Init(s1 + s2)
				}
			} else if gpu.LastKey == glfw.KeyDelete {
				s1 := s.Buffer.Slice(0, max(s.SelStart, 0))
				s2 := s.Buffer.Slice(min(s.SelEnd+1, s.Buffer.RuneCount()), s.Buffer.RuneCount())
				s.Buffer.Init(s1 + s2)
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyLeft {
				s.SelStart = max(0, s.SelStart-1)
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyRight {
				s.SelStart = min(s.SelStart+1, s.Buffer.RuneCount())
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyEnd {
				s.SelStart = s.Buffer.RuneCount()
				s.SelEnd = s.SelStart
			} else if gpu.LastKey == glfw.KeyHome {
				s.SelStart = 0
				s.SelEnd = s.SelStart
			}
		} else if mouse.Hovered(frameRect) {
			bg = theme.Colors[theme.SurfaceContainer]
		}

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, theme.Colors[style.BorderColor])
		f.SetColor(theme.Colors[style.FontColor])
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
			gpu.VertLine(textRect.X+dx, textRect.Y, textRect.Y+textRect.H, 1, theme.Colors[theme.Primary])
		}
		return Dim{W: frameRect.W, H: frameRect.H, baseline: baseline}
	}
}
