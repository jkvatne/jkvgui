package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type SwitchStyle struct {
	// Height          float32
	Padding         f32.Padding
	ShadowSize      float32
	BorderThickness float32
	Track           theme.UIRole
	Knob            theme.UIRole
	On              theme.UIRole
	FontNo          int
}

var DefaultSwitchStyle = &SwitchStyle{
	// Height:          15,
	Padding:         f32.Padding{L: 3, T: 1, R: 2, B: 1},
	ShadowSize:      4,
	BorderThickness: 1.0,
	Track:           theme.Surface,
	Knob:            theme.Outline,
	On:              theme.Primary,
	FontNo:          gpu.Normal12,
}

func Switch(label string, state *bool, action func(), style *SwitchStyle, hint string) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = DefaultSwitchStyle
		}
		f := font.Fonts[style.FontNo]
		h := f.Height
		labelWidth := f.Width(label) + style.Padding.L + style.Padding.R + 2
		baseline := f.Baseline + style.Padding.T
		width := h*13/8 + style.Padding.R + style.Padding.L
		height := h + style.Padding.T + style.Padding.B
		if h > height {
			height = h + style.Padding.T + style.Padding.B
		}
		if ctx.Mode != RenderChildren {
			return Dim{W: width + labelWidth, H: height, Baseline: baseline}
		}

		ctx.Rect.W = width
		ctx.Rect.H = height
		if *DebugWidgets {
			ctx.Win.Gd.RoundedRect(ctx.Rect, 0, 0.5, f32.Transparent, f32.Blue)
		}
		track := ctx.Rect.Inset(style.Padding, 0)
		knob := track.Reduce(height / 5).Square()
		knob.W = knob.H
		// Move knob to the right if it is on.
		if *state {
			knob.X += height / 2
		}
		if ctx.Win.Hovered(track) || ctx.Win.At(state) {
			ctx.Win.Gd.Shade(knob.Increase(style.ShadowSize), -1, f32.Shade, style.ShadowSize)
			Hint(ctx, hint, state)
		}
		if ctx.Win.LeftBtnPressed(ctx.Rect) {
			ctx.Win.SetFocusedTag(state)
		}
		if ctx.Win.LeftBtnClick(ctx.Rect) || ctx.Win.At(state) && IsKeyClick(ctx) {
			*state = !*state
			if action != nil {
				action()
			}
		}
		if *state == false {
			ctx.Win.Gd.RoundedRect(track, -1, style.BorderThickness, style.Track.Bg(), style.Knob.Bg())
			ctx.Win.Gd.RoundedRect(knob, -1, 0.0, style.Knob.Bg(), style.Knob.Bg())
		} else {
			ctx.Win.Gd.RoundedRect(track, -1, style.BorderThickness, style.On.Bg(), style.On.Bg())
			ctx.Win.Gd.RoundedRect(knob, -1, 0.0, style.On.Fg(), style.On.Fg())
		}
		f.DrawText(ctx.Win.Gd, track.X+width, knob.Y+knob.H+style.Padding.T, style.Track.Fg(), 0, gpu.LTR, label)

		return Dim{W: width, H: height, Baseline: 0}
	}
}
