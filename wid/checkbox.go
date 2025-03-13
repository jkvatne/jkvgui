package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

type CheckboxStyle struct {
	FontSize       float32
	FontNo         int
	Color          f32.Color
	OutsidePadding f32.Padding
}

var DefaultCheckbox = CheckboxStyle{
	FontSize:       1,
	FontNo:         0,
	Color:          f32.Color{R: 0, G: 0, B: 0, A: 1},
	OutsidePadding: f32.Padding{L: 3, T: 3, R: 3, B: 3},
}

func Checkbox(text string, state *bool, style *CheckboxStyle, hint string) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultCheckbox
		}
		f := gpu.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.OutsidePadding.T + style.OutsidePadding.B
		if ctx.Rect.H == 0 {
			return Dim{w: height, h: height, baseline: 0}
		}

		gpu.MoveFocus(state)
		if gpu.LeftMouseBtnPressed(ctx.Rect) {

		} else if gpu.LeftMouseBtnReleased(ctx.Rect) {
			gpu.MouseBtnReleased = false
			gpu.SetFocus(state)
			*state = !*state
		} else if gpu.Focused(state) {
			if gpu.MoveFocusToNext {
				gpu.FocusToNext = true
				gpu.MoveFocusToNext = false
			}
		} else if gpu.Hovered(ctx.Rect) {
		}
		gpu.AddFocusable(ctx.Rect, nil)

		if gpu.Hovered(ctx.Rect) {
			Hint(hint, state)
		}
		r := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: height, H: height}
		if *state {
			DrawIcon(r.X, r.Y, r.W, BoxChecked, style.Color)
		} else {
			DrawIcon(r.X, r.Y, r.W, BoxUnchecked, style.Color)
		}
		return Dim{w: ctx.Rect.W, h: ctx.Rect.H, baseline: ctx.Baseline}
	}
}
