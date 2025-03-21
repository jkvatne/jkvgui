package wid

import (
	"fmt"
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
	LabelFraction      float32
	LabelRightAdjust   bool
	LabelSpacing       float32
}

var DefaultEdit = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	InsideColor:        theme.Surface,
	BorderColor:        theme.Outline,
	FontColor:          theme.OnSurface,
	OutsidePadding:     f32.Padding{2, 3, 2, 3},
	InsidePadding:      f32.Padding{4, 1, 2, 1},
	BorderWidth:        0.66,
	BorderCornerRadius: 4,
	CursorWidth:        2,
	LabelFraction:      0.5,
	LabelRightAdjust:   true,
	LabelSpacing:       3,
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

func EditF32(label string, f32 *float32, action func(), style *EditStyle) Wid {
	ss := fmt.Sprintf("%0.2f", *f32)
	return Edit(label, &ss, action, style)
}

func EditInt(label string, i32 *int, action func(), style *EditStyle) Wid {
	ss := fmt.Sprintf("%d", *i32)
	return Edit(label, &ss, action, style)
}

func Edit(label string, text *string, action func(), style *EditStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultEdit
		}
		if text == nil {
			return Dim{}
		}
		s := StateMap[text]
		if s == nil {
			StateMap[text] = &EditState{}
			s = StateMap[text]
			s.Buffer.Init(*text)
		}
		f := font.Get(style.FontNo, theme.Colors[style.FontColor])

		widRect := ctx.Rect.Inset(style.OutsidePadding)
		frameRect := widRect
		if label != "" {
			frameRect.X += style.LabelFraction * widRect.W
			frameRect.W *= (1 - style.LabelFraction)
		}
		valueRect := frameRect.Inset(style.InsidePadding).Reduce(style.BorderWidth)
		labelRect := valueRect
		labelRect.X = widRect.X
		if label != "" {
			labelRect.W = style.LabelFraction * widRect.W
		} else {
			labelRect.W = 0
		}
		fontHeight := f.Height(style.FontSize)
		baseline := f.Baseline(style.FontSize)

		if ctx.Rect.H == 0 {
			return Dim{W: 32, H: fontHeight + style.TotalPaddingY(), Baseline: baseline}
		}

		bg := theme.Colors[style.InsideColor]
		focused := focus.At(ctx.Rect, text)

		if mouse.LeftBtnPressed(widRect) {
			gpu.Invalidate(0)
		}
		if mouse.LeftBtnReleased(widRect) {
			gpu.Invalidate(0)
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(text)
			s.SelStart = f.RuneNo(mouse.Pos().X-(valueRect.X), style.FontSize, s.Buffer.String())
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
		} else if mouse.Hovered(widRect) {
			bg = theme.Colors[theme.SurfaceContainer]
		}

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, theme.Colors[style.BorderColor])
		f.SetColor(theme.Colors[style.FontColor])
		labelWidth := f.Width(style.FontSize, label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		// Draw label
		f.Printf(
			labelRect.X+dx,
			valueRect.Y+baseline,
			style.FontSize,
			labelRect.W,
			label)
		// Draw value
		f.Printf(
			valueRect.X,
			valueRect.Y+baseline,
			style.FontSize,
			valueRect.W,
			s.Buffer.String())
		f.SetColor(f32.Black)
		s.SelEnd = s.SelStart
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 1, theme.Colors[theme.Primary])
		}
		return Dim{W: frameRect.W, H: frameRect.H, Baseline: baseline}
	}
}
