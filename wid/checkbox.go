package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
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
		f := font.Fonts[style.FontNo]
		height := f.Height(style.FontSize) + style.OutsidePadding.T + style.OutsidePadding.B
		if ctx.Rect.H == 0 {
			return Dim{w: height, h: height, baseline: 0}
		}

		focus.Move(state)
		if focus.LeftMouseBtnPressed(ctx.Rect) {

		} else if focus.LeftMouseBtnReleased(ctx.Rect) {
			focus.MouseBtnReleased = false
			focus.Set(state)
			*state = !*state
		} else if focus.At(state) {
			if focus.MoveToNext {
				focus.ToNext = true
				focus.MoveToNext = false
			}
		} else if focus.Hovered(ctx.Rect) {
		}
		focus.AddFocusable(ctx.Rect, nil)

		if focus.Hovered(ctx.Rect) {
			Hint(hint, state)
		}
		r := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: height, H: height}
		if *state {
			gpu.DrawIcon(r.X, r.Y, r.W, gpu.BoxChecked, style.Color)
		} else {
			gpu.DrawIcon(r.X, r.Y, r.W, gpu.BoxUnchecked, style.Color)
		}
		return Dim{w: ctx.Rect.W, h: ctx.Rect.H, baseline: ctx.Baseline}
	}
}
