package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

type ComboStyle struct {
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

type ComboState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
	index    int
	expanded bool
}

var DefaultCombo = ComboStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	InsideColor:        theme.Surface,
	BorderColor:        theme.Outline,
	FontColor:          theme.OnSurface,
	OutsidePadding:     f32.Padding{L: 2, T: 3, R: 2, B: 3},
	InsidePadding:      f32.Padding{L: 4, T: 2, R: 2, B: 2},
	BorderWidth:        1,
	BorderCornerRadius: 5,
	CursorWidth:        1.5,
}

func (s *ComboStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func setValue(i int, s *ComboState, list []string) {
	s.index = i
	s.Buffer.Init(list[i])
	s.expanded = false
	gpu.Invalidate(0)
}

var ComboStateMap = make(map[*string]*ComboState)

func Combo(text *string, list []string, style *ComboStyle) Wid {
	return func(ctx Ctx) Dim {
		// Make sure we have a style
		if style == nil {
			style = &DefaultCombo
		}
		// Initialize the state of the widget
		s := ComboStateMap[text]
		if s == nil {
			ComboStateMap[text] = &ComboState{}
			s = ComboStateMap[text]
			s.Buffer.Init(*text)
		}
		fg := theme.Colors[style.FontColor]
		bg := theme.Colors[style.InsideColor]
		f := font.Get(style.FontNo, fg)

		frameRect := ctx.Rect.Inset(style.OutsidePadding, 0)
		textRect := frameRect.Inset(style.InsidePadding, style.BorderWidth)
		fontHeight := f.Height(style.FontSize)
		baseline := f.Baseline(style.FontSize)

		if ctx.Rect.H == 0 {
			// Return minimum size
			return Dim{W: textRect.W, H: fontHeight + style.TotalPaddingY(), Baseline: baseline}
		}

		focused := focus.At(ctx.Rect, text)

		// Calculate the icon size and position for the drop-down arrow
		iconX := ctx.Rect.X + ctx.Rect.W - style.OutsidePadding.R - style.BorderWidth - style.InsidePadding.R - fontHeight
		iconY := ctx.Rect.Y + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		// Detect click on the "down arrow"
		if mouse.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight, H: fontHeight}) {
			s.expanded = !s.expanded
			gpu.Invalidate(0)
			focus.Set(text)
		}

		if !focused {
			s.expanded = false
		}
		if s.expanded {
			if gpu.LastKey == glfw.KeyDown {
				s.index = min(s.index+1, len(list)-1)
			} else if gpu.LastKey == glfw.KeyUp {
				s.index = max(s.index-1, 0)
			} else if gpu.Return() {
				setValue(s.index, s, list)
				gpu.LastKey = 0
			}

			for i := range len(list) {
				itemRect := frameRect
				itemRect.Y = frameRect.Y + frameRect.H + float32(i)*itemRect.H
				if mouse.LeftBtnClick(itemRect) {
					setValue(i, s, list)
				}
			}

			dropDownBox := func() {
				baseline := f.Baseline(style.FontSize) + style.InsidePadding.T
				for i := range len(list) {
					ibg := theme.Colors[theme.Surface]
					if i == s.index {
						ibg = theme.Colors[theme.SurfaceContainer]
					}
					itemRect := frameRect
					itemRect.Y = frameRect.Y + frameRect.H + float32(i)*itemRect.H
					gpu.RoundedRect(itemRect, 0, 0.5, ibg, theme.Colors[theme.Outline])
					x := textRect.X
					f.SetColor(theme.Colors[style.FontColor])
					f.Printf(x, itemRect.Y+baseline, style.FontSize, itemRect.W, list[i])
				}
			}
			gpu.Defer(dropDownBox)
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
				*text = str[0:max(len(str)-1, 0)]
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
			} else if gpu.LastKey == glfw.KeyEnter {
				if s.expanded {
					setValue(s.index, s, list)
				} else {
					s.expanded = true
				}
				gpu.Invalidate(0)
			} else if gpu.LastKey != 0 {
				gpu.Invalidate(0)
			}
		} else if mouse.Hovered(frameRect) {
			bg = theme.Colors[theme.SurfaceContainer]
		}

		if mouse.LeftBtnPressed(frameRect) {
			gpu.Invalidate(0)
			// col.A = 1
		}

		if mouse.LeftBtnClick(frameRect) {
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(text)
			s.SelStart = f.RuneNo(mouse.Pos().X-(frameRect.X), style.FontSize, s.Buffer.String())
			s.SelEnd = s.SelStart
			gpu.Invalidate(0)
		}

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, theme.Colors[style.BorderColor])
		f.SetColor(fg)
		f.Printf(
			textRect.X,
			textRect.Y+baseline,
			style.FontSize,
			textRect.W-fontHeight,
			s.Buffer.String())
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(textRect.X+dx, textRect.Y, textRect.Y+textRect.H, 1, fg)
		}

		icon.Draw(iconX, iconY, fontHeight, icon.ArrowDropDown, fg)

		return Dim{W: frameRect.W, H: frameRect.H, Baseline: baseline}
	}
}
