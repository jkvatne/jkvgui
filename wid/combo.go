package wid

import (
	"fmt"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type ComboState struct {
	EditState
	ScrollState
	index    int
	expanded bool
}

type ComboStyle struct {
	EditStyle
	MaxDropDown int
	NotEditable bool
}

var DefaultCombo = ComboStyle{
	EditStyle: EditStyle{
		FontNo:             gpu.Normal12,
		Color:              theme.Surface,
		BorderColor:        theme.Outline,
		OutsidePadding:     f32.Padding{L: 2, T: 2, R: 2, B: 2},
		InsidePadding:      f32.Padding{L: 2, T: 2, R: 2, B: 2},
		BorderWidth:        0.66,
		BorderCornerRadius: 4,
		CursorWidth:        2,
		EditSize:           0.0,
		LabelSize:          0.0,
		LabelRightAdjust:   true,
		LabelSpacing:       3,
	},
	MaxDropDown: 10,
	NotEditable: false,
}

var GridCombo = ComboStyle{
	EditStyle:   GridEdit,
	MaxDropDown: 10,
	NotEditable: true,
}

func (s *ComboStyle) Size(wl, we float32) *ComboStyle {
	ss := *s
	ss.EditSize = we
	ss.LabelSize = wl
	return &ss
}
func setValue(i int, s *ComboState, list []string, value any) {
	s.index = i
	s.Buffer.Init(list[i])
	s.expanded = false
	sys.Invalidate()
	sys.CurrentInfo.Mutex.Lock()
	defer sys.CurrentInfo.Mutex.Unlock()
	switch v := value.(type) {
	case *int:
		*v = s.index
	case *string:
		*v = list[s.index]
	}

}

var ComboStateMap = make(map[any]*ComboState)

func List(value any, list []string, label string, style *ComboStyle) Wid {
	if style == nil {
		s := ComboStyle{}
		s = DefaultCombo
		s.NotEditable = true
	}
	return Combo(value, list, label, style)
}

func Combo(value any, list []string, label string, style *ComboStyle) Wid {
	// Make sure we have a style
	if style == nil {
		style = &DefaultCombo
	}
	f32.ExitIf(value == nil, "Combo with nil value")

	// Initialize the state of the widget
	state := ComboStateMap[value]
	if state == nil {
		ComboStateMap[value] = &ComboState{}
		state = ComboStateMap[value]
		sys.CurrentInfo.Mutex.Lock()
		switch v := value.(type) {
		case *int:
			state.Buffer.Init(list[*v])
		case *string:
			state.Buffer.Init(fmt.Sprintf("%s", *v))
		default:
			f32.Exit("Combo with value that is not *int or  *string")
		}
		sys.CurrentInfo.Mutex.Unlock()
	}
	// Precalculate some values
	f := font.Get(style.FontNo)
	fontHeight := f.Height
	baseline := f.Baseline
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := style.Dim(&ctx, f)
		if ctx.Mode != RenderChildren {
			return dim
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", &style.EditStyle, ctx.Rect)
		// Correct for icon at end
		valueRect.W -= fontHeight

		// Calculate the icon size and position for the drop-down arrow
		iconX := valueRect.X + valueRect.W
		iconY := frameRect.Y + style.InsidePadding.T

		if sys.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight * 1.2, H: fontHeight * 1.2}) {
			// Detect click on the "down arrow"
			state.expanded = true
			sys.Invalidate()
			sys.SetFocusedTag(value)
		}

		focused := sys.At(ctx.Rect, value)
		EditHandleMouse(&state.EditState, valueRect, f, value, focused)

		if state.expanded {
			if sys.LastKey == sys.KeyDown {
				state.index = min(state.index+1, len(list)-1)
			} else if sys.LastKey == sys.KeyUp {
				state.index = max(state.index-1, 0)
			} else if sys.Return() {
				setValue(state.index, state, list, value)
				sys.LastKey = 0
			}

			dropDownBox := func() {
				state.ScrollState.dragging = state.ScrollState.dragging && sys.LeftBtnDown()
				baseline := f.Baseline
				lineHeight := fontHeight + style.InsidePadding.T + style.InsidePadding.B
				// Find the number of visible lines
				Nvis := min(len(list), int((gpu.ClientRectDp.H-frameRect.Y-frameRect.H)/lineHeight))
				if Nvis >= len(list) {
					state.Npos = 0
					state.Dy = 0
					state.Ypos = 0
				}
				listHeight := float32(Nvis) * lineHeight
				// listRect is the rectangle where the list text is
				listRect := f32.Rect{X: frameRect.X, Y: frameRect.Y + frameRect.H, W: frameRect.W, H: listHeight}
				gpu.Shade(listRect, 3, f32.Shade, 5)
				gpu.Rect(listRect, 0, theme.Surface.Bg(), theme.Surface.Bg())
				lineRect := f32.Rect{X: listRect.X, Y: listRect.Y, W: listRect.W, H: lineHeight}
				state.Ymax = float32(len(list)) * lineHeight
				state.Nmax = len(list)
				gpu.Clip(listRect)
				n := 0
				lineRect.Y -= state.Dy
				for i := state.Npos; i < len(list); i++ {
					n++
					if i == state.index {
						gpu.Rect(lineRect, 0, theme.SurfaceContainer.Bg(), theme.SurfaceContainer.Bg())
					} else if sys.Hovered(lineRect) {
						gpu.Rect(lineRect, 0, theme.PrimaryContainer.Bg(), theme.PrimaryContainer.Bg())
					} else {
						gpu.Rect(lineRect, 0, theme.Surface.Bg(), theme.Surface.Bg())
					}
					if sys.LeftBtnClick(lineRect) {
						state.expanded = false
						setValue(i, state, list, value)
					}
					f.DrawText(lineRect.X+style.InsidePadding.L, lineRect.Y+baseline+style.InsidePadding.T, fg, lineRect.W, gpu.LTR, list[i])
					lineRect.Y += lineHeight
					if lineRect.Y > gpu.ClientRectDp.H {
						break
					}
				}
				if len(list) > Nvis {
					DrawVertScrollbar(listRect, float32(len(list))*lineRect.H, float32(Nvis)*lineRect.H, &state.ScrollState)
				}

				if sys.LeftBtnClick(f32.Rect{X: 0, Y: 0, W: 999999, H: 999999}) {
					state.expanded = false
				}
				gpu.NoClip()

				yScroll := VertScollbarUserInput(listRect.H, &state.ScrollState)
				scrollUp(yScroll, &state.ScrollState, func(n int) float32 {
					return lineHeight
				})
				scrollDown(yScroll, &state.ScrollState, listRect.H, func(n int) float32 {
					return lineHeight
				})

			}
			sys.CurrentInfo.SuppressEvents = true
			gpu.Defer(dropDownBox)
		}

		if focused {
			bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
			if !style.NotEditable {
				EditText(&state.EditState)
			}
			if sys.LastKey == sys.KeyEnter {
				if state.expanded {
					setValue(state.index, state, list, value)
				} else {
					state.expanded = true
				}
				sys.Invalidate()
			}
		} else {
			state.expanded = false
		}
		if sys.LeftBtnClick(frameRect) && !style.NotEditable {
			sys.SetFocusedTag(value)
			state.SelStart = f.RuneNo(sys.Pos().X-(frameRect.X), state.Buffer.String())
			state.SelEnd = state.SelStart
			sys.Invalidate()
		}

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				f.DrawText(labelRect.X+labelRect.W-f.Width(label), valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			} else {
				f.DrawText(labelRect.X, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			}
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused && !style.NotEditable {
			r := valueRect
			r.H--
			r.W = f.Width(state.Buffer.Slice(state.SelStart, max(state.SelStart, state.SelEnd)))
			r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().MultAlpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && !style.NotEditable {
			DrawCursor(&style.EditStyle, &state.EditState, valueRect, f)
		}

		// Draw dropdown arrow
		gpu.DrawIcon(iconX, iconY, fontHeight*1.2, gpu.ArrowDropDown, fg)

		// Draw frame around value
		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, f32.Transparent, style.BorderColor.Fg())

		// Draw debugging rectngles if wid.DebugWidgets is true
		DrawDebuggingInfo(labelRect, valueRect, ctx.Rect)

		return dim
	}
}
