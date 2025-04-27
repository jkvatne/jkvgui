package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
)

type BtnStyle struct {
	FontNo         int
	FontWeight     float32
	BtnRole        theme.UIRole
	BorderColor    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	Disabled       *bool
	IconPad        float32
}

var Filled = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 12, T: 5, R: 12, B: 5},
	BorderWidth:    0,
	CornerRadius:   6,
	Disabled:       nil,
	IconPad:        0.15,
}

var Text = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Transparent,
	BorderColor:    theme.Transparent,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 5, T: 5, R: 5, B: 5},
	BorderWidth:    0,
	CornerRadius:   0,
	Disabled:       nil,
	IconPad:        0.15,
}

var Outline = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Transparent,
	BorderColor:    theme.Outline,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 1},
	InsidePadding:  f32.Padding{L: 5, T: 5, R: 5, B: 5},
	BorderWidth:    1,
	CornerRadius:   6,
	Disabled:       nil,
	IconPad:        0.15,
}

var Round = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Transparent,
	OutsidePadding: f32.Padding{L: 5.5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 5, T: 5, R: 5, B: 5},
	BorderWidth:    0,
	CornerRadius:   -1,
	Disabled:       nil,
}

func (s *BtnStyle) Role(c theme.UIRole) *BtnStyle {
	ss := *s
	ss.BtnRole = c
	return &ss
}

func (s *BtnStyle) Font(n int) *BtnStyle {
	ss := *s
	ss.FontNo = n
	return &ss
}

func (s *BtnStyle) RR(r float32) *BtnStyle {
	ss := *s
	ss.CornerRadius = r
	return &ss
}

func Btn(text string, ic *gpu.Icon, action func(), style *BtnStyle, hint string) Wid {
	if style == nil {
		style = Filled
	}
	f := font.Fonts[style.FontNo]
	fontHeight := f.Ascent()
	baseline := f.Baseline() + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
	height := fontHeight + style.OutsidePadding.T + style.OutsidePadding.B +
		style.InsidePadding.T + style.InsidePadding.B + 2*style.BorderWidth
	width := font.Fonts[style.FontNo].Width(text) +
		style.InsidePadding.L + style.InsidePadding.R + 2*style.BorderWidth +
		style.OutsidePadding.R + style.OutsidePadding.L
	if ic != nil {
		if text == "" {
			width = height
		} else {
			width += fontHeight * 1.15
		}
	}

	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: width, H: height, Baseline: baseline}
		}
		ctx.Baseline = max(ctx.Baseline, baseline)
		ctx.Rect.W = width
		ctx.Rect.H = height
		b := style.BorderWidth
		btnOutline := ctx.Rect.Inset(style.OutsidePadding, 0)
		textRect := btnOutline.Inset(style.InsidePadding, style.BorderWidth)
		cr := style.CornerRadius
		if !ctx.Disabled {
			if mouse.LeftBtnPressed(ctx.Rect) {
				gpu.Shade(btnOutline.Outset(f32.Padding{L: 4, T: 4, R: 4, B: 4}).Move(0, 0), cr, f32.Shade, 4)
				b += 1
			} else if mouse.Hovered(ctx.Rect) {
				gpu.Shade(btnOutline.Outset(f32.Pad(2)), cr, f32.Shade, 4)
				Hint(hint, action)
			}
			if action != nil && mouse.LeftBtnClick(ctx.Rect) {
				focus.Set(action)
				if !ctx.Disabled {
					action()
					gpu.Invalidate(0)
				}
			}
			if focus.At(ctx.Rect, action) {
				b += 1
				gpu.Shade(btnOutline.Outset(f32.Pad(2)).Move(0, 0),
					cr, f32.Shade, 4)
			}
		}
		fg := style.BtnRole.Fg().Alpha(ctx.Alpha())
		bg := style.BtnRole.Bg().Alpha(ctx.Alpha())
		gpu.RoundedRect(btnOutline, cr, b, bg, theme.Colors[style.BorderColor])
		if ic != nil {
			gpu.DrawIcon(textRect.X, ctx.Rect.Y+baseline-0.85*fontHeight, fontHeight, ic, fg)
			textRect.X += fontHeight + style.IconPad*fontHeight
			textRect.W -= fontHeight + style.IconPad*fontHeight
		}
		f.DrawText(textRect.X, textRect.Y+f.Baseline(), fg, 0, gpu.LTR, text)
		if gpu.DebugWidgets {
			gpu.Rect(ctx.Rect, 1.0, f32.Transparent, f32.Red)
			gpu.Rect(textRect, 1.0, f32.Transparent, f32.Yellow)
			gpu.HorLine(textRect.X, textRect.X+textRect.W, textRect.Y+f.Baseline(), 1.0, f32.Blue)
		}
		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
