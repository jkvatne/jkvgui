package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
)

type SwitchStyle struct {
	// Height          float32
	Pad             f32.Padding
	ShadowSize      float32
	BorderThickness float32
	Track           theme.UIRole
	Knob            theme.UIRole
	On              theme.UIRole
	FontNo          int
}

var DefaultSwitchStyle = &SwitchStyle{
	// Height:          15,
	Pad:             f32.Padding{3, 1, 2, 1},
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
		h := f.Height()
		labelWidth := f.Width(label) + style.Pad.L + style.Pad.R + 2
		baseline := f.Baseline() + style.Pad.T
		width := h*13/8 + style.Pad.R + style.Pad.L
		height := h + style.Pad.T + style.Pad.B
		if h > height {
			height = h + style.Pad.T + style.Pad.B
		}
		if ctx.Mode != RenderChildren {
			return Dim{W: width + labelWidth, H: height, Baseline: baseline}
		}

		ctx.Rect.W = width
		ctx.Rect.H = height
		if gpu.DebugWidgets {
			gpu.RoundedRect(ctx.Rect, 0, 0.5, f32.Transparent, f32.Blue)
		}
		track := ctx.Rect.Inset(style.Pad, 0)
		knob := track.Reduce(height / 5).Square()
		knob.W = knob.H
		// Move konp to the right if it is on.
		if *state {
			knob.X += height / 2
		}
		if mouse.Hovered(track) || focus.At(track, state) {
			gpu.Shade(knob.Out(style.ShadowSize), -1, f32.Shade, style.ShadowSize)
		}
		if mouse.LeftBtnClick(ctx.Rect) {
			focus.SetFocusedTag(state)
			*state = !*state
		}
		if *state == false {
			gpu.RoundedRect(track, -1, style.BorderThickness, style.Track.Bg(), style.Knob.Fg())
			gpu.RoundedRect(knob, -1, 0.0, style.Knob.Fg(), style.Knob.Fg())
		} else {
			gpu.RoundedRect(track, -1, style.BorderThickness, style.On.Bg(), style.On.Bg())
			gpu.RoundedRect(knob, -1, 0.0, style.On.Fg(), style.On.Fg())
		}
		f.DrawText(track.X+width, knob.Y+knob.H+style.Pad.T, style.Track.Fg(), 0, gpu.LTR, label)

		return Dim{W: width, H: height, Baseline: 0}
	}
}
