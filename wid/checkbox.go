package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type CheckboxStyle struct {
	EditStyle
}

var DefaultCheckbox = CheckboxStyle{
	EditStyle: EditStyle{
		FontNo:         gpu.Normal12,
		OutsidePadding: f32.Padding{L: 3, T: 1, R: 2, B: 1},
	},
}

var GridCb = CheckboxStyle{
	EditStyle: EditStyle{
		FontNo:        gpu.Normal12,
		EditSize:      18,
		Color:         theme.PrimaryContainer,
		BorderColor:   theme.Transparent,
		InsidePadding: f32.Padding{L: 2, T: 0, R: 2, B: 0},
		CursorWidth:   1,
		BorderWidth:   GridBorderWidth,
		Dp:            2,
	},
}

func (s *CheckboxStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func Checkbox(label string, state *bool, style *CheckboxStyle, hint string) Wid {
	if style == nil {
		style = &DefaultCheckbox
	}
	f := font.Fonts[style.FontNo]
	fontHeight := f.Height
	baseline := f.Baseline

	return func(ctx Ctx) Dim {
		dim := style.Dim(&ctx, f)
		if ctx.Mode != RenderChildren {
			return dim
		}
		if ctx.H < dim.H {
			ctx.H = dim.H
		}
		frameRect, _, labelRect := CalculateRects(label != "", &style.EditStyle, ctx.Rect)
		iconRect := labelRect
		iconRect.W = iconRect.H
		if ctx.Win.LeftBtnClick(ctx.Rect) {
			ctx.Win.SetFocusedTag(state)
			*state = !*state
		}
		if ctx.Win.At(ctx.Rect, state) {
			ctx.Win.Gd.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
		}
		if ctx.Win.Hovered(ctx.Rect) {
			ctx.Win.Gd.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
			Hint(ctx, hint, state)
		}
		if *state {
			ctx.Win.Gd.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxChecked, style.Color.Fg())
		} else {
			ctx.Win.Gd.DrawIcon(iconRect.X, iconRect.Y, iconRect.H, gpu.BoxUnchecked, style.Color.Fg())
		}
		labelRect.X += fontHeight * 6 / 5
		f.DrawText(ctx.Win.Gd, labelRect.X, labelRect.Y+baseline, style.Color.Fg(), 0, gpu.LTR, label)
		// Draw frame around value
		ctx.Win.Gd.RoundedRect(frameRect, style.BorderCornerRadius, style.BorderWidth, f32.Transparent, style.BorderColor.Fg())
		DrawDebuggingInfo(ctx, iconRect, iconRect, ctx.Rect)

		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
