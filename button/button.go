package button

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type Style struct {
	FontSize       float32
	FontNo         int
	FontWeight     float32
	BtnRole        theme.UIRole
	BorderColor    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	Disabled       *bool
}

var TextBtn = Style{
	FontSize:       1.0,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Secondary,
	BorderColor:    theme.Secondary,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 12, T: 4, R: 12, B: 4},
	BorderWidth:    0,
	CornerRadius:   12,
}

var RoundBtn = Style{
	FontSize:       1.5,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 12, T: 4, R: 12, B: 4},
	BorderWidth:    0,
	CornerRadius:   9999,
}

var Btn = Style{
	FontSize:       1.5,
	FontNo:         gpu.Normal,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 12, T: 4, R: 12, B: 4},
	BorderWidth:    0,
	CornerRadius:   6,
	Disabled:       nil,
}

func (s *Style) Role(c theme.UIRole) *Style {
	ss := *s
	ss.BtnRole = c
	return &ss
}

func Filled(text string, ic *icon.Icon, action func(), style *Style, hint string) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		if style == nil {
			style = &Btn
		}
		f := font.Fonts[style.FontNo]
		fontHeight := f.Height(style.FontSize)

		height := fontHeight + style.OutsidePadding.T + style.OutsidePadding.B +
			style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
		width := font.Fonts[style.FontNo].Width(style.FontSize, text) +
			style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth +
			style.OutsidePadding.R + style.OutsidePadding.L
		baseline := f.Baseline(style.FontSize) + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth

		if ctx.Rect.H == 0 {
			return wid.Dim{W: width, H: height, Baseline: baseline}
		}
		if ic != nil {
			width += fontHeight * 1.15
		}
		ctx.Rect.W = width
		ctx.Rect.H = height
		b := style.BorderWidth
		r := ctx.Rect
		r.Y += ctx.Baseline - baseline
		r = r.Inset(style.OutsidePadding)
		cr := style.CornerRadius
		if !ctx.Disabled {
			if mouse.LeftBtnPressed(ctx.Rect) {
				gpu.Shade(r.Outset(f32.Padding{L: 4, T: 4, R: 4, B: 4}).Move(0, 0), cr, f32.Shade, 4)
				b += 1
			} else if mouse.Hovered(ctx.Rect) {
				gpu.Shade(r.Outset(f32.Padding{L: 4, T: 4, R: 4, B: 4}).Move(2, 2), cr, f32.Shade, 4)
			}
			if action != nil && mouse.LeftBtnReleased(ctx.Rect) {
				focus.Set(action)
				if !ctx.Disabled {
					action()
					gpu.Invalidate(0)
				}
			}
			if focus.At(ctx.Rect, action) {
				b += 1
				gpu.Shade(r.Outset(f32.UniformPad(2)).Move(0, 0),
					cr, f32.Shade, 4)
			}

			if mouse.Hovered(ctx.Rect) {
				wid.Hint(hint, action)
			}
		}
		fg := style.BtnRole.Fg().Alpha(ctx.Alpha())
		bg := style.BtnRole.Bg().Alpha(ctx.Alpha())
		gpu.RoundedRect(r, cr, b, bg, theme.Colors[style.BorderColor])
		r = r.Inset(style.InsidePadding).Reduce(style.BorderWidth)
		if ic != nil {
			icon.Draw(r.X, ctx.Rect.Y+baseline-0.85*fontHeight, fontHeight, ic, fg)
			r.X += fontHeight * 1.15
		}
		f.SetColor(fg)
		f.Printf(
			r.X,
			ctx.Rect.Y+baseline,
			style.FontSize, 0,
			text)

		return wid.Dim{}
	}
}
