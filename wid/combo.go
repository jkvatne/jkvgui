package wid

import (
	"fmt"
	"log/slog"

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
	ScrollStyle
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
		BorderWidth:        1,
		BorderCornerRadius: 4,
		CursorWidth:        2,
		EditSize:           0.0,
		LabelSize:          0.0,
		LabelRightAdjust:   true,
		LabelSpacing:       2,
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

func (s *ComboStyle) Pad(p float32) *ComboStyle {
	ss := *s
	ss.OutsidePadding = f32.Pad(p)
	return &ss
}

func (s *ComboStyle) D(f *bool) *ComboStyle {
	ss := *s
	ss.Disabler = f
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
	// Initialize the state of the widget
	StateMapMutex.RLock()
	state := ComboStateMap[value]
	StateMapMutex.RUnlock()

	if state == nil {
		slog.Debug("Combo: Create new state")
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
			f32.Exit(1, "Combo with value that is not *int or  *string")
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
		if ctx.H < 0 {
			return Dim{}
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", &style.EditStyle, ctx.Rect)
		// Correct for icon at end
		valueRect.W -= fontHeight

		// Calculate the icon size and position for the drop-down arrow
		iconX := valueRect.X + valueRect.W
		iconY := frameRect.Y + style.InsidePadding.T
		focused := ctx.Win.At(value)

		if !style.Disabled() {
			if ctx.Win.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight * 1.2, H: fontHeight * 1.2}) {
				// Detect click on the "down arrow"
				slog.Debug("Combo: LeftBtnClick on down-arrow caused combo list to expand")
				state.expanded = true
				ctx.Win.Invalidate()
				ctx.Win.SetFocusedTag(value)
			}

			if ctx.Win.LeftBtnDoubleClick(ctx.Rect) {
				slog.Debug("Combo: LeftBtnClick on double-click caused combo list to expand")
				state.expanded = true
				ctx.Win.Invalidate()
				ctx.Win.SetFocusedTag(value)
			}
			EditMouseHandler(ctx, &state.EditState, valueRect, f, value)

			if state.expanded {
				if ctx.Win.LastKey == sys.KeyDown {
					state.index = min(state.index+1, len(list)-1)
				} else if ctx.Win.LastKey == sys.KeyUp {
					state.index = max(state.index-1, 0)
				} else if ctx.Win.LastKey == sys.KeyEnter || ctx.Win.LastKey == sys.KeyKPEnter {
					setValue(ctx, state.index, state, list, value)
					ctx.Win.LastKey = 0
				} else if ctx.Win.LastKey == sys.KeyEscape {
					slog.Debug("Combo: Esc key caused combo list to collapse")
					state.expanded = false
				}

				// This function is run after all other drawing commands
				dropDownBox := func() {
					state.ScrollState.Dragging = state.ScrollState.Dragging && ctx.Win.LeftBtnDown()
					lineHeight := fontHeight + style.InsidePadding.T + style.InsidePadding.B
					// Find the number of visible lines
					VisibleLines := min(len(list), int((ctx.Win.HeightDp-frameRect.Y-frameRect.H)/lineHeight))
					if VisibleLines >= len(list) {
						state.Npos = 0
						state.Dy = 0
						state.Ypos = 0
					}
					listHeight := float32(VisibleLines) * lineHeight
					// listRect is the rectangle where the list text is
					listRect := f32.Rect{X: frameRect.X, Y: frameRect.Y + frameRect.H, W: frameRect.W, H: listHeight}
					ctx.Win.Gd.Shade(listRect, 3, f32.Shade, 5)
					ctx.Win.Gd.SolidRect(listRect, theme.Surface.Bg())
					lineRect := f32.Rect{X: listRect.X, Y: listRect.Y, W: listRect.W, H: lineHeight}
					state.Ymax = float32(len(list)) * lineHeight
					state.Nmax = len(list)
					ctx0 := ctx
					ctx0.Rect = listRect
					VertScollbarUserInput(ctx0, &state.ScrollState, &style.ScrollStyle)
					doScrolling(ctx, &state.ScrollState, func(n int) float32 {
						return lineHeight
					})
					ctx.Win.Gd.Clip(listRect)
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
						if ctx.Win.LeftBtnClick(lineRect) {
							slog.Debug("Combo: LeftBtnPressed in expanded combo on", "line", i)
							state.expanded = false
							setValue(ctx, i, state, list, value)
						}
						f.DrawText(ctx.Win.Gd, lineRect.X+style.InsidePadding.L, lineRect.Y+baseline+style.InsidePadding.T, fg, lineRect.W, gpu.LTR, list[i])
						lineRect.Y += lineHeight
						if lineRect.Y > ctx.Win.HeightDp {
							break
						}
					}
					if len(list) > VisibleLines {
						DrawVertScrollbar(ctx0, &state.ScrollState, nil)
					}

					if ctx.Win.LeftBtnClick(f32.Rect{X: 0, Y: 0, W: 999999, H: 999999}) {
						slog.Debug("Combo: LeftBtnClick caused combo list to collapse")
						state.expanded = false
					}
					gpu.NoClip()
					ctx.Rect = listRect
				}
				ctx.Win.SuppressEvents = true
				ctx.Win.Defer(dropDownBox)
			}

			if focused {
				bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
				if !style.NotEditable {
					EditText(ctx, &state.EditState, nil)
				}
				if ctx.Win.LastKey == sys.KeyEnter {
					if state.expanded {
						setValue(ctx, state.index, state, list, value)
					} else {
						slog.Debug("Combo: Enter key caused combo list to expand")
						state.expanded = true
					}
				}

			} else if state.expanded {
				slog.Debug("Combo: Lost focus, not expanded")
				state.expanded = false
			}

			if ctx.Win.LeftBtnClick(frameRect) && !style.NotEditable && ctx.Win.At(value) {
				slog.Debug("Combo: LeftBtnClick, set focus.")
				ctx.Win.SetFocusedTag(value)
				state.SelStart = f.RuneNo(ctx.Win.MousePos().X-(frameRect.X), state.Buffer.String())
				state.SelEnd = state.SelStart
				ctx.Win.Invalidate()
			}
		}

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				dx := max(0.0, labelRect.W-f.Width(label)-style.LabelSpacing)
				f.DrawText(ctx.Win.Gd, labelRect.X+dx, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
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
		if style.Disabled() {
			fg = fg.Mute(0.3)
		}

		f.DrawText(ctx.Win.Gd, valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && !style.NotEditable {
			DrawCursor(ctx, &style.EditStyle, &state.EditState, valueRect, f)
		}

		// Draw dropdown arrow
		if !style.Disabled() {
			ctx.Win.Gd.DrawIcon(iconX, iconY, fontHeight*1.2, gpu.ArrowDropDown, fg)
		}
		// Draw frame around value
		bg := f32.Transparent
		if state.hovered {
			bg = fg.MultAlpha(0.05)
		}
		ctx.Win.Gd.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, style.BorderColor.Bg())

		// Draw debugging rectangles if wid.DebugWidgets is true
		DrawDebuggingInfo(ctx, labelRect, valueRect, ctx.Rect)

		return dim
	}
}
