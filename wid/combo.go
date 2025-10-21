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
func setValue(ctx Ctx, i int, s *ComboState, list []string, value any) {
	s.index = i
	s.Buffer.Init(list[i])
	s.expanded = false
	ctx.Win.Invalidate()
	switch v := value.(type) {
	case *int:
		*v = s.index
	case *string:
		*v = list[s.index]
	}

}

var ComboStateMap = make(map[any]*ComboState)

func Combo(value any, list []string, label string, style *ComboStyle) Wid {
	// Make sure we have a style
	if style == nil {
		style = &DefaultCombo
	}
	f32.ExitIf(value == nil, "Combo with nil value")

	// Initialize the state of the widget
	StateMapMutex.RLock()
	state := ComboStateMap[value]
	StateMapMutex.RUnlock()
	if state == nil {
		StateMapMutex.Lock()
		s := ComboState{}
		s.Id = 999
		ComboStateMap[value] = &s
		state = ComboStateMap[value]
		StateMapMutex.Unlock()
		switch v := value.(type) {
		case *int:
			state.Buffer.Init(list[*v])
		case *string:
			state.Buffer.Init(fmt.Sprintf("%s", *v))
		default:
			f32.Exit("Combo with value that is not *int or  *string")
		}
	}
	// Precalculate some values
	f := font.Get(style.FontNo)
	fontHeight := f.Height
	baseline := f.Baseline
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := style.Dim(ctx.Rect.W, f)
		ctx.H = min(ctx.H, dim.H)
		if ctx.Mode != RenderChildren {
			return dim
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", &style.EditStyle, ctx.Rect)
		// Correct for icon at end
		valueRect.W -= fontHeight

		// Calculate the icon size and position for the drop-down arrow
		iconX := valueRect.X + valueRect.W
		iconY := frameRect.Y + style.InsidePadding.T

		if ctx.Win.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight * 1.2, H: fontHeight * 1.2}) {
			// Detect click on the "down arrow"
			state.expanded = true
			ctx.Win.Invalidate()
			ctx.Win.SetFocusedTag(value)
		}

		focused := ctx.Win.At(value)
		EditMouseHandler(ctx, &state.EditState, valueRect, f, value)

		if state.expanded {
			if ctx.Win.LastKey == sys.KeyDown {
				state.index = min(state.index+1, len(list)-1)
			} else if ctx.Win.LastKey == sys.KeyUp {
				state.index = max(state.index-1, 0)
			} else if ctx.Win.LastKey == sys.KeyEnter || ctx.Win.LastKey == sys.KeyKPEnter {
				setValue(ctx, state.index, state, list, value)
				ctx.Win.LastKey = 0
			}

			dropDownBox := func() {
				state.ScrollState.Dragging = state.ScrollState.Dragging && ctx.Win.LeftBtnDown()
				lineHeight := fontHeight + style.InsidePadding.T + style.InsidePadding.B
				// Find the number of visible lines
				Nvis := min(len(list), int((ctx.Win.HeightDp-frameRect.Y-frameRect.H)/lineHeight))
				if Nvis >= len(list) {
					state.Npos = 0
					state.Dy = 0
					state.Ypos = 0
				}
				listHeight := float32(Nvis) * lineHeight
				// listRect is the rectangle where the list text is
				listRect := f32.Rect{X: frameRect.X, Y: frameRect.Y + frameRect.H, W: frameRect.W, H: listHeight}
				ctx.Win.Gd.Shade(listRect, 3, f32.Shade, 5)
				ctx.Win.Gd.SolidRect(listRect, theme.Surface.Bg())
				lineRect := f32.Rect{X: listRect.X, Y: listRect.Y, W: listRect.W, H: lineHeight}
				state.Ymax = float32(len(list)) * lineHeight
				state.Yest = state.Ymax
				state.Nmax = len(list)
				ctx.Win.Clip(listRect)
				n := 0
				lineRect.Y -= state.Dy
				for i := state.Npos; i < len(list); i++ {
					n++
					if i == state.index {
						ctx.Win.Gd.SolidRect(lineRect, theme.SurfaceContainer.Bg())
					} else if ctx.Win.Hovered(lineRect) {
						ctx.Win.Gd.SolidRect(lineRect, theme.PrimaryContainer.Bg())
					} else {
						ctx.Win.Gd.SolidRect(lineRect, theme.Surface.Bg())
					}
					if ctx.Win.LeftBtnPressed(lineRect) {
						state.expanded = false
						setValue(ctx, i, state, list, value)
					}
					f.DrawText(ctx.Win.Gd, lineRect.X+style.InsidePadding.L, lineRect.Y+baseline+style.InsidePadding.T, fg, lineRect.W, gpu.LTR, list[i])
					lineRect.Y += lineHeight
					if lineRect.Y > ctx.Win.HeightDp {
						break
					}
				}
				if len(list) > Nvis {
					ctx0 := ctx
					ctx0.Rect = listRect
					DrawVertScrollbar(ctx0, &state.ScrollState)
				}

				if ctx.Win.LeftBtnClick(f32.Rect{X: 0, Y: 0, W: 999999, H: 999999}) {
					state.expanded = false
				}
				gpu.NoClip()
				ctx.Rect = listRect
				yScroll := VertScollbarUserInput(ctx, &state.ScrollState)
				scrollUp(yScroll, &state.ScrollState, func(n int) float32 {
					return lineHeight
				})
				scrollDown(ctx, yScroll, &state.ScrollState, func(n int) float32 {
					return lineHeight
				})

			}
			ctx.Win.SuppressEvents = true
			ctx.Win.Defer(dropDownBox)
		}

		if focused {
			bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
			if !style.NotEditable {
				EditText(ctx, &state.EditState)
			}
			if ctx.Win.LastKey == sys.KeyEnter {
				if state.expanded {
					setValue(ctx, state.index, state, list, value)
				} else {
					state.expanded = true
				}
				ctx.Win.Invalidate()
			}
		} else {
			state.expanded = false
		}
		if ctx.Win.LeftBtnClick(frameRect) && !style.NotEditable {
			ctx.Win.SetFocusedTag(value)
			state.SelStart = f.RuneNo(ctx.Win.MousePos().X-(frameRect.X), state.Buffer.String())
			state.SelEnd = state.SelStart
			ctx.Win.Invalidate()
		}

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				f.DrawText(ctx.Win.Gd, labelRect.X+labelRect.W-f.Width(label), valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			} else {
				f.DrawText(ctx.Win.Gd, labelRect.X, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			}
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused && !style.NotEditable {
			r := valueRect
			r.H--
			r.W = f.Width(state.Buffer.Slice(state.SelStart, max(state.SelStart, state.SelEnd)))
			r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().MultAlpha(0.8)
			ctx.Win.Gd.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.DrawText(ctx.Win.Gd, valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && !style.NotEditable {
			DrawCursor(ctx, &style.EditStyle, &state.EditState, valueRect, f)
		}

		// Draw dropdown arrow
		ctx.Win.Gd.DrawIcon(iconX, iconY, fontHeight*1.2, gpu.ArrowDropDown, fg)

		// Draw frame around value
		bg := f32.Transparent
		if state.hovered {
			bg = fg.MultAlpha(0.05)
		}
		ctx.Win.Gd.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, style.BorderColor.Fg())

		// Draw debugging rectngles if wid.DebugWidgets is true
		DrawDebuggingInfo(ctx, labelRect, valueRect, ctx.Rect)

		return dim
	}
}
