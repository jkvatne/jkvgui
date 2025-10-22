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

var GridCheckBox = CheckboxStyle{
	EditStyle: EditStyle{
		FontNo:        gpu.Normal12,
		EditSize:      0,
		Color:         theme.PrimaryContainer,
		InsidePadding: f32.Padding{L: 1, T: 1, R: 1, B: 1},
		BorderWidth:   GridBorderWidth,
	},
}

func Checkbox(label string, state *bool, style *CheckboxStyle, hint string) Wid {
	if style == nil {
		style = &DefaultCheckbox
	}
	f := font.Fonts[style.FontNo]
	fontHeight := f.Height
	baseline := f.Baseline

	return func(ctx Ctx) Dim {
		ctx0 := ctx
		if ctx.H <= 0 {
			return Dim{}
		}
		dim := style.Dim(ctx.Rect.W, f)
		ctx.H = min(ctx.H, dim.H)
		if ctx.Mode != RenderChildren {
			return dim
		}
		frameRect, _, labelRect := CalculateRects(label != "", &style.EditStyle, ctx0.Rect)
		iconRect := labelRect
		iconRect.W = fontHeight
		iconRect.H = fontHeight
		if ctx.Win.LeftBtnClick(ctx.Rect) {
			ctx.Win.SetFocusedTag(state)
			*state = !*state
		}
		if ctx.Win.At(state) {
			ctx.Win.Gd.Shade(iconRect, 4, f32.Shade, 3)
		}
		r := labelRect
		r.W = iconRect.W + f.Width(label)
		if ctx.Win.Hovered(r) {
			ctx.Win.Gd.Shade(iconRect, 4, f32.Shade, 3)
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
