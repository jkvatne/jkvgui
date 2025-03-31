package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type CheckboxStyle struct {
	FontSize float32
	FontNo   int
	Role     theme.UIRole
	Padding  f32.Padding
}

var DefaultCheckbox = CheckboxStyle{
	FontSize: 1.0,
	FontNo:   0,
	Role:     theme.Surface,
	Padding:  f32.Padding{L: 5, T: 3, R: 8, B: 3},
}

func Checkbox(text string, state *bool, style *CheckboxStyle, hint string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = &DefaultCheckbox
		}
		f := font.Fonts[style.FontNo]
		fontHeight := f.Height(style.FontSize)
		height := fontHeight + style.Padding.T + style.Padding.B
		width := f.Width(style.FontSize, text) + style.Padding.L + style.Padding.R + height
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		extRect := f32.Rect{X: ctx.Rect.X, Y: ctx.Rect.Y, W: width, H: height}
		iconRect := extRect.Inset(style.Padding, 0)
		iconRect.W = iconRect.H
		if ctx.Rect.H == 0 {
			return wid.Dim{W: extRect.W, H: extRect.H, Baseline: baseline}
		}
		if gpu.DebugWidgets {
			gpu.RoundedRect(extRect, 0, 0.5, f32.Transparent, f32.Blue)
		}

		focused := focus.At(ctx.Rect, state)

		if mouse.LeftBtnClick(extRect) {
			focus.Set(state)
			*state = !*state
		}
		if focused {
			gpu.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
		}
		if mouse.Hovered(extRect) {
			gpu.Shade(iconRect.Move(0, -1), 4, f32.Shade, 3)
			wid.Hint(hint, state)
		}
		// Icon checkbox is 3/4 of total size. Square is 45, box is 60 when H=17.2 and ScaleX=1.75. H=30. Ascenders=30
		if *state {
			icon.Draw(iconRect.X, iconRect.Y-1, iconRect.H, icon.BoxChecked, style.Role.Fg())
		} else {
			icon.Draw(iconRect.X, iconRect.Y-1, iconRect.H, icon.BoxUnchecked, style.Role.Fg())
		}
		f.DrawText(iconRect.X+fontHeight*6/5, extRect.Y+baseline, style.Role.Fg(), style.FontSize, 0, gpu.LeftToRight, text)

		return wid.Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
