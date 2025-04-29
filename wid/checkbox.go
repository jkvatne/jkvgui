package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
)

type CheckboxStyle struct {
	EditStyle
}

var DefaultCheckbox = CheckboxStyle{
	EditStyle: EditStyle{
		FontNo:         gpu.Normal12,
		OutsidePadding: f32.Padding{L: 3, T: 0, R: 2, B: 0},
	},
}

var GridCb = CheckboxStyle{
	EditStyle: GridEdit,
}

func (s *CheckboxStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func Checkbox(label string, state *bool, style *CheckboxStyle, hint string) Wid {
	if style == nil {
		style = &DefaultCheckbox
	}
	f := font.Fonts[style.FontNo]
	fontHeight := f.Height()
	height := fontHeight + style.OutsidePadding.T + style.OutsidePadding.B
	width := f.Width(label) + style.OutsidePadding.L + style.OutsidePadding.R + style.InsidePadding.L + style.InsidePadding.R + height
	baseline := f.Baseline()

	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return style.Dim(width, f)
		}

		frameRect, _, labelRect := CalculateRects(label != "", &style.EditStyle, ctx.Rect)
		iconRect := labelRect
		iconRect.W = iconRect.H

		focused := focus.At(ctx.Rect, state)

		if mouse.LeftBtnClick(ctx.Rect) {
			focus.SetFocusedTag(state)
			*state = !*state
		}
		if focused {
			gpu.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
		}
		if mouse.Hovered(ctx.Rect) || (focused && !*state) {
			gpu.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
			Hint(hint, state)
		}
		gpu.RoundedRect(iconRect, 0, 0.5, f32.Transparent, f32.Blue)
		if *state {
			gpu.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxChecked, style.Color.Fg())
		} else {
			gpu.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxUnchecked, style.Color.Fg())
		}
		labelRect.X += fontHeight * 6 / 5
		f.DrawText(labelRect.X, labelRect.Y+baseline, style.Color.Fg(), 0, gpu.LTR, label)

		DrawDebuggingInfo(iconRect, iconRect, ctx.Rect)

		// Draw frame around value
		gpu.RoundedRect(frameRect, style.BorderCornerRadius, style.BorderWidth, f32.Transparent, style.BorderColor.Fg())

		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
