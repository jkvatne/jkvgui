package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/glfw"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

type ComboState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
	index    int
	expanded bool
}

var DefaultCombo = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	Color:              theme.Surface,
	BorderColor:        theme.Outline,
	OutsidePadding:     f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:      f32.Padding{L: 4, T: 2, R: 2, B: 2},
	BorderWidth:        0.66,
	BorderCornerRadius: 4,
	CursorWidth:        2,
	EditSize:           0.0,
	LabelSize:          0.0,
	LabelRightAdjust:   true,
	LabelSpacing:       3,
}

func setValue(i int, s *ComboState, list []string) {
	s.index = i
	s.Buffer.Init(list[i])
	s.expanded = false
	gpu.Invalidate(0)
}

var ComboStateMap = make(map[*string]*ComboState)

func Combo(text *string, list []string, label string, style *EditStyle) Wid {
	// Make sure we have a style
	if style == nil {
		style = &DefaultCombo
	}
	f32.ExitIf(text == nil, "Combo with nil value")

	f := font.Get(style.FontNo)
	fontHeight := f.Height(style.FontSize)
	baseline := f.Baseline(style.FontSize)
	bg := style.Color.Bg()
	fg := style.Color.Fg()

	// Initialize the state of the widget
	s := ComboStateMap[text]
	if s == nil {
		ComboStateMap[text] = &ComboState{}
		s = ComboStateMap[text]
		s.Buffer.Init(*text)
	}

	return func(ctx Ctx) Dim {
		dim := Dim{W: ctx.W, H: fontHeight + style.TotalPaddingY(), Baseline: baseline + style.Top()}
		if ctx.Mode != RenderChildren {
			return dim
		}
		frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)

		labelWidth := f.Width(style.FontSize, label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		// Draw label if it exists
		if label != "" {
			f.DrawText(
				labelRect.X+dx,
				labelRect.Y+baseline,
				fg,
				style.FontSize,
				labelRect.W-fontHeight, gpu.LTR,
				label)
		}

		if style.LabelSize > 1.0 {
			frameRect.X += style.LabelSize
			frameRect.W = style.EditSize
			labelRect.W -= style.LabelSize
			labelRect.X += style.LabelSize
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
		} else if !focused {
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
				baseline := f.Baseline(style.FontSize)
				r := f32.Rect{frameRect.X, frameRect.Y + frameRect.H,
					frameRect.W, float32(len(list)) * frameRect.H}
				gpu.Shade(r.Move(3, 3), 5, f32.Shade, 5)
				gpu.Rect(r, 1, theme.Surface.Bg(), theme.Outline.Fg())
				r = frameRect
				r.Y += frameRect.H
				r.H = fontHeight + style.InsidePadding.T + style.InsidePadding.B
				for i := range len(list) {
					if i == s.index {
						gpu.Rect(r.Inset(f32.Pad(0), 1), 0, theme.SurfaceContainer.Bg(), theme.SurfaceContainer.Bg())
					}
					f.DrawText(valueRect.X, r.Y+baseline+style.InsidePadding.T, fg, style.FontSize, r.W, gpu.LTR, list[i])
					r.Y += r.H
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

		f.DrawText(
			valueRect.X,
			valueRect.Y+baseline,
			fg,
			style.FontSize,
			valueRect.W-fontHeight, gpu.LTR,
			s.Buffer.String())
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			n := s.Buffer.RuneCount()
			if s.SelStart > n {
				s.SelStart = n
			}
			dx := f.Width(style.FontSize, s.Buffer.Slice(0, s.SelStart))
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 1, fg)
		}

		gpu.Draw(iconX, iconY, fontHeight, gpu.ArrowDropDown, fg)

		if gpu.DebugWidgets {
			gpu.Rect(labelRect, 1, f32.Transparent, f32.LightBlue)
			gpu.Rect(valueRect, 1, f32.Transparent, f32.LightRed)
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.Yellow)
		}

		return dim
	}
}
