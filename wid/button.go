package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type BtnStyle struct {
	FontNo         int
	BtnRole        theme.UIRole
	BorderColor    theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	Disabled       *bool
	IconPad        float32
	IconMagn       float32
	Width          float32
}

var Filled = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Primary,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 12, T: 5, R: 12, B: 7},
	BorderWidth:    0,
	CornerRadius:   6,
	Disabled:       nil,
	IconPad:        1,
	IconMagn:       1.3,
}

var Text = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Transparent,
	BorderColor:    theme.Transparent,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 5, T: 5, R: 5, B: 7},
	BorderWidth:    0,
	CornerRadius:   6,
	Disabled:       nil,
	IconPad:        1,
	IconMagn:       1.3,
}

var Outline = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Transparent,
	BorderColor:    theme.Outline,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 5, T: 5, R: 5, B: 7},
	BorderWidth:    1,
	CornerRadius:   6,
	Disabled:       nil,
	IconPad:        1,
	IconMagn:       1.3,
}

var Round = &BtnStyle{
	FontNo:         gpu.Normal14,
	BtnRole:        theme.Primary,
	BorderColor:    theme.Transparent,
	OutsidePadding: f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:  f32.Padding{L: 6, T: 6, R: 6, B: 6},
	BorderWidth:    0,
	CornerRadius:   -1,
	Disabled:       nil,
	IconMagn:       1.3,
}

var Header = &BtnStyle{
	FontNo:        gpu.Normal12,
	InsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
	BtnRole:       theme.PrimaryContainer,
	BorderColor:   theme.Outline,
	BorderWidth:   GridBorderWidth,
	Width:         0.3,
}

var CbHeader = &BtnStyle{
	FontNo:        gpu.Normal12,
	InsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
	BtnRole:       theme.PrimaryContainer,
	BorderColor:   theme.Outline,
	BorderWidth:   GridBorderWidth,
	Width:         18,
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
	baseline := f.Baseline + style.OutsidePadding.T + style.InsidePadding.T + style.BorderWidth
	height := f.Baseline + style.OutsidePadding.T + style.OutsidePadding.B +
		style.InsidePadding.T + style.InsidePadding.B
	width := font.Fonts[style.FontNo].Width(text) +
		style.InsidePadding.L + style.InsidePadding.R +
		style.OutsidePadding.R + style.OutsidePadding.L
	width = max(width, height)
	if ic != nil {
		if text == "" {
			width = height
		} else {
			width += f.Baseline*style.IconMagn + style.IconPad
		}
	}
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			if style.Width > 0 {
				width = style.Width
			}
			return Dim{W: width, H: height, Baseline: baseline}
		}
		if ctx.Rect.W > 1.0 {
			width = ctx.Rect.W
		}
		ctx.Baseline = max(ctx.Baseline, baseline)
		ctx.Rect.H = height
		bw := style.BorderWidth
		btnOutline := ctx.Rect.Inset(style.OutsidePadding, 0)
		btnOutline.Y += ctx.Baseline - baseline
		textRect := btnOutline.Inset(style.InsidePadding, 0)
		cr := style.CornerRadius
		if !ctx.Disabled {
			if sys.LeftBtnPressed(ctx.Rect) {
				gpu.Shade(btnOutline.Outset(f32.Padding{L: 4, T: 4, R: 4, B: 4}).Move(0, 0), cr, f32.Shade, 4)
				bw += 0.5
			} else if sys.Hovered(ctx.Rect) {
				gpu.Shade(btnOutline.Outset(f32.Pad(2)), cr, f32.Shade, 4)
				if hint != "" {
					Hint(hint, action)
				}
			}
			if action != nil && sys.LeftBtnClick(ctx.Rect) {
				sys.SetFocusedTag(action)
				if !ctx.Disabled {
					action()
					sys.Invalidate()
				}
			}
			if sys.At(ctx.Rect, action) {
				gpu.Shade(btnOutline.Outset(f32.Pad(2)).Move(0, 0),
					cr, f32.Shade, 4)
			}
		}
		fg := style.BtnRole.Fg().MultAlpha(ctx.Alpha())
		bg := style.BtnRole.Bg().MultAlpha(ctx.Alpha())

		btnOutline.X -= style.BorderWidth / 2
		btnOutline.Y -= style.BorderWidth / 2
		btnOutline.W += style.BorderWidth
		btnOutline.H += style.BorderWidth
		gpu.RoundedRect(btnOutline, cr, bw, bg, theme.Colors[style.BorderColor])

		iconRect := f32.Rect{X: textRect.X - textRect.H*0.15, Y: textRect.Y - textRect.H*0.15, W: textRect.H * style.IconMagn, H: textRect.H * style.IconMagn}
		if ic != nil {
			gpu.DrawIcon(iconRect.X, iconRect.Y, iconRect.W, ic, fg)
			textRect.X += iconRect.W + style.IconPad
			textRect.W -= iconRect.W + style.IconPad
		}
		f.DrawText(textRect.X, textRect.Y+f.Baseline, fg, 0, gpu.LTR, text)
		if *DebugWidgets {
			gpu.Rect(iconRect, 0.5, f32.Transparent, f32.Green)
			gpu.Rect(ctx.Rect, 0.5, f32.Transparent, f32.Red)
			gpu.Rect(textRect, 0.5, f32.Transparent, f32.Yellow)
			gpu.HorLine(textRect.X, textRect.X+textRect.W, textRect.Y+f.Baseline, 0.5, f32.Blue)
		}
		return Dim{W: ctx.Rect.W, H: ctx.Rect.H, Baseline: ctx.Baseline}
	}
}
