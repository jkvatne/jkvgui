package wid

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
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
		OutsidePadding:     f32.Padding{L: 5, T: 5, R: 5, B: 5},
		InsidePadding:      f32.Padding{L: 4, T: 2, R: 2, B: 2},
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
	gpu.Invalidate(0)
	gpu.GpuMutex.Lock()
	defer gpu.GpuMutex.Unlock()
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
		style = &DefaultCombo
	}
	style.NotEditable = true
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
		gpu.GpuMutex.Lock()
		switch v := value.(type) {
		case *int:
			state.Buffer.Init(list[*v])
		case *string:
			state.Buffer.Init(fmt.Sprintf("%s", *v))
		default:
			f32.Exit("Combo with value that is not *int or  *string")
		}
		gpu.GpuMutex.Unlock()
	}
	// Precalculate some values
	f := font.Get(style.FontNo)
	fontHeight := f.Height()
	baseline := f.Baseline()
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

		labelWidth := f.Width(label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}

		// Calculate the icon size and position for the drop-down arrow
		iconX := frameRect.X + frameRect.W - fontHeight
		iconY := frameRect.Y + style.InsidePadding.T

		if mouse.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight, H: fontHeight}) {
			// Detect click on the "down arrow"
			state.expanded = true
			gpu.Invalidate(0)
			focus.SetFocusedTag(value)
		}

		focused := focus.At(ctx.Rect, value)
		EditHandleMouse(&state.EditState, valueRect, f, value)

		if state.expanded {
			if gpu.LastKey == glfw.KeyDown {
				state.index = min(state.index+1, len(list)-1)
			} else if gpu.LastKey == glfw.KeyUp {
				state.index = max(state.index-1, 0)
			} else if gpu.Return() {
				setValue(state.index, state, list, value)
				gpu.LastKey = 0
			}

			for i := range len(list) {
				itemRect := frameRect
				itemRect.Y = frameRect.Y + frameRect.H + float32(i)*itemRect.H
				if mouse.LeftBtnClick(itemRect) {
					setValue(i, state, list, value)
				}
			}

			dropDownBox := func() {
				baseline := f.Baseline()
				r1 := f32.Rect{frameRect.X, frameRect.Y + frameRect.H,
					frameRect.W, float32(min(len(list), style.MaxDropDown)) * frameRect.H}
				gpu.Shade(r1.Move(3, 3), 5, f32.Shade, 5)
				gpu.Rect(r1, 1, theme.Surface.Bg(), theme.Outline.Fg())
				r := frameRect
				r.Y += frameRect.H
				r.H = fontHeight + style.InsidePadding.T + style.InsidePadding.B
				sumH := r.H * float32(len(list))
				r.Y -= state.Ypos
				gpu.Clip(r1)
				for i := range len(list) {
					if i == state.index {
						gpu.Rect(r.Inset(f32.Pad(0), 1), 0, theme.SurfaceContainer.Bg(), theme.SurfaceContainer.Bg())
					} else if mouse.Hovered(r) {
						gpu.Rect(r.Inset(f32.Pad(0), 1), 0, theme.SurfaceContainer.Bg(), theme.SurfaceContainer.Bg())
					}
					if mouse.LeftBtnClick(r) {
						state.expanded = false
						setValue(i, state, list, value)
					}
					f.DrawText(valueRect.X, r.Y+baseline+style.InsidePadding.T, fg, r.W, gpu.LTR, list[i])
					r.Y += r.H
					if r.Y > r.Y+r.H {
						break
					}
				}
				if len(list) > 4 {
					DrawVertScrollbar(r1, sumH, r1.H, &state.ScrollState)
				}
				if mouse.LeftBtnClick(f32.Rect{X: 0, Y: 0, W: 999999, H: 999999}) {
					state.expanded = false
				}
				gpu.NoClip()
			}
			gpu.SupressEvents = true
			gpu.Defer(dropDownBox)
		}

		if focused {
			bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
			if !style.NotEditable {
				EditText(&state.EditState)
			}
			if gpu.LastKey == glfw.KeyEnter {
				if state.expanded {
					setValue(state.index, state, list, value)
				} else {
					state.expanded = true
				}
				gpu.Invalidate(0)
			}
		} else {
			state.expanded = false
		}
		if mouse.LeftBtnClick(frameRect) && !style.NotEditable {
			focus.SetFocusedTag(value)
			state.SelStart = f.RuneNo(mouse.Pos().X-(frameRect.X), state.Buffer.String())
			state.SelEnd = state.SelStart
			gpu.Invalidate(0)
		}

		// Draw label if it exists
		if label != "" {
			f.DrawText(labelRect.X+dx, valueRect.Y+baseline, fg, labelRect.W-fontHeight, gpu.LTR, label)
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused && !style.NotEditable {
			r := valueRect
			r.H--
			r.W = f.Width(state.Buffer.Slice(state.SelStart, state.SelEnd))
			r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().Alpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && !style.NotEditable {
			DrawCursor(&style.EditStyle, &state.EditState, valueRect, f)
		}

		// Draw dropdown arrow
		gpu.DrawIcon(iconX, iconY, fontHeight, gpu.ArrowDropDown, fg)

		// Draw frame around value
		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, f32.Transparent, style.BorderColor.Fg())

		// Draw debugging rectngles if gpu.DebugWidgets is true
		DrawDebuggingInfo(labelRect, valueRect, ctx.Rect)

		return dim
	}
}
