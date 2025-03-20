package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type ButtonStyle struct {
	FontSize       float32
	FontNo         int
	FontWeight     float32
	BtnRole        theme.UIRole
	BorderColor    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
}

var TextBtn = ButtonStyle{
	FontSize:       1.0,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Secondary,
	BorderColor:    theme.Secondary,
	OutsidePadding: f32.Padding{5, 5, 5, 5},
	InsidePadding:  f32.Padding{12, 4, 12, 4},
	BorderWidth:    0,
	CornerRadius:   12,
}

var RoundBtn = ButtonStyle{
	FontSize:       1.5,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{5, 5, 5, 5},
	InsidePadding:  f32.Padding{12, 4, 12, 4},
	BorderWidth:    0,
	CornerRadius:   9999,
}

var Btn = ButtonStyle{
	FontSize:       1.5,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{5, 5, 5, 5},
	InsidePadding:  f32.Padding{12, 4, 12, 4},
	BorderWidth:    0,
	CornerRadius:   6,
}

func (s *ButtonStyle) Role(c theme.UIRole) *ButtonStyle {
	ss := *s
	ss.BtnRole = c
	return &ss
}

var fg f32.Color
var bg f32.Color

func Filled(text string, icon *gpu.Icon, action func(), style *ButtonStyle, hint string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = &Btn
		}
		f := font.Fonts[style.FontNo]
		fontHeight := f.Height(style.FontSize)

		height := fontHeight + style.OutsidePadding.T + style.OutsidePadding.B + style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		width := font.Fonts[style.FontNo].Width(style.FontSize, text) +
			style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth +
			style.OutsidePadding.R + style.OutsidePadding.L
		baseline := f.Baseline(style.FontSize) + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		if ctx.Rect.H == 0 {
			return wid.Dim{W: width, H: height, Baseline: baseline}
		}
		if icon != nil {
			width += fontHeight * 1.15
		}
		ctx.Rect.W = width
		ctx.Rect.H = height
		r := ctx.Rect
		r.Y += ctx.Baseline - baseline
		b := style.BorderWidth
		r = r.Inset(style.OutsidePadding)
		cr := min(style.CornerRadius, r.H/2)
		if mouse.LeftBtnPressed(ctx.Rect) {
			gpu.Shade(r.Outset(f32.Padding{4, 4, 4, 4}).Move(0, 0), cr, f32.Shade, 4)
			b += 1
		} else if mouse.Hovered(ctx.Rect) {
			gpu.Shade(r.Outset(f32.Padding{4, 4, 4, 4}).Move(2, 2), cr, f32.Shade, 4)
		}
		if mouse.LeftBtnReleased(ctx.Rect) {
			focus.Set(action)
		}
		if focus.At(ctx.Rect, action) {
			b += 1
		}

		if mouse.Hovered(ctx.Rect) {
			wid.Hint(hint, action)
		}
		fg = style.BtnRole.Fg()
		bg = style.BtnRole.Bg()
		gpu.RoundedRect(r, cr, b, style.BtnRole.Bg(), theme.Colors[style.BorderColor])
		// x0 := ctx.Rect.X + style.OutsidePadding.L + style.InsidePadding.L + style.BorderWidth
		// y0 := ctx.Rect.Y + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
		r = r.Inset(style.InsidePadding).Reduce(style.BorderWidth)
		if icon != nil {
			gpu.DrawIcon(r.X, ctx.Rect.Y+baseline-0.85*fontHeight, fontHeight, icon, fg)
			r.X += fontHeight * 1.15
		}
		f.SetColor(style.BtnRole.Fg())
		f.Printf(
			r.X,
			ctx.Rect.Y+baseline,
			style.FontSize, 0, text)

		return wid.Dim{}
	}
}
