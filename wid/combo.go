package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"time"
)

type ComboState struct {
	EditState
	index    int
	expanded bool
}

var DefaultCombo = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	Color:              theme.Surface,
	BorderColor:        theme.Outline,
	OutsidePadding:     f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:      f32.Padding{L: 4, T: 3, R: 2, B: 2},
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

func Combo(value *string, list []string, label string, style *EditStyle) Wid {
	// Make sure we have a style
	if style == nil {
		style = &DefaultCombo
	}
	f32.ExitIf(value == nil, "Combo with nil value")

	f := font.Get(style.FontNo)
	fontHeight := f.Height(style.FontSize)
	baseline := f.Baseline(style.FontSize)
	bg := style.Color.Bg()
	fg := style.Color.Fg()

	// Initialize the state of the widget
	state := ComboStateMap[value]
	if state == nil {
		ComboStateMap[value] = &ComboState{}
		state = ComboStateMap[value]
		state.Buffer.Init(*value)
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

		focused := focus.At(ctx.Rect, value)
		EditHandleMouse(&state.EditState, valueRect, f, style.FontSize, value)

		// Calculate the icon size and position for the drop-down arrow
		iconX := frameRect.X + frameRect.W - fontHeight
		iconY := frameRect.Y + style.InsidePadding.T

		// Detect click on the "down arrow"
		if mouse.LeftBtnClick(f32.Rect{X: iconX, Y: iconY, W: fontHeight, H: fontHeight}) {
			state.expanded = state.expanded
			gpu.Invalidate(0)
			focus.Set(value)
		} else if !focused {
			state.expanded = false
		}
		if state.expanded {
			if gpu.LastKey == glfw.KeyDown {
				state.index = min(state.index+1, len(list)-1)
			} else if gpu.LastKey == glfw.KeyUp {
				state.index = max(state.index-1, 0)
			} else if gpu.Return() {
				setValue(state.index, state, list)
				gpu.LastKey = 0
			}

			for i := range len(list) {
				itemRect := frameRect
				itemRect.Y = frameRect.Y + frameRect.H + float32(i)*itemRect.H
				if mouse.LeftBtnClick(itemRect) {
					setValue(i, state, list)
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
					if i == state.index {
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
			EditText(&state.EditState)
			if gpu.LastKey == glfw.KeyEnter {
				if state.expanded {
					setValue(state.index, state, list)
				} else {
					state.expanded = true
				}
				gpu.Invalidate(0)
			}
		} else if mouse.Hovered(frameRect) {
			bg = theme.Colors[theme.SurfaceContainer]
		}

		if mouse.LeftBtnPressed(frameRect) {
			gpu.Invalidate(0)
		}

		if mouse.LeftBtnClick(frameRect) {
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(value)
			state.SelStart = f.RuneNo(mouse.Pos().X-(frameRect.X), style.FontSize, state.Buffer.String())
			state.SelEnd = state.SelStart
			gpu.Invalidate(0)
		}

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, theme.Colors[style.BorderColor])
		// Draw label if it exists
		if label != "" {
			f.DrawText(labelRect.X+dx, valueRect.Y+baseline, fg, style.FontSize, labelRect.W-fontHeight, gpu.LTR, label)
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused {
			r := valueRect
			r.H--
			r.W = f.Width(style.FontSize, state.Buffer.Slice(state.SelStart, state.SelEnd))
			r.X += f.Width(style.FontSize, state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().Alpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, style.FontSize,
			valueRect.W-fontHeight, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx = f.Width(style.FontSize, state.Buffer.Slice(0, state.SelEnd))
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 1, style.Color.Fg())
		}

		// Draw dropdown arrow
		gpu.Draw(iconX, iconY, fontHeight, gpu.ArrowDropDown, fg)

		if gpu.DebugWidgets {
			gpu.Rect(labelRect, 1, f32.Transparent, f32.LightBlue)
			gpu.Rect(valueRect, 1, f32.Transparent, f32.LightRed)
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.Yellow)
		}

		return dim
	}
}
