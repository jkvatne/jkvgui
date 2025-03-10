package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"time"
)

var CurrentHint struct {
	Pos    f32.Pos // Mouse position at time of pop-up
	Text   string
	T      time.Time
	Tag    any
	Active bool
}

type HintStyle struct {
	FontNo          int
	FontSize        float32
	CornerRadius    float32
	BorderColor     f32.Color
	BackgroundColor f32.Color
	BorderWidth     float32
	Padding         f32.Padding
	Delay           time.Duration
}

var DefaultHintStyle = HintStyle{
	FontNo:          0,
	FontSize:        9,
	CornerRadius:    5,
	BorderColor:     f32.Color{R: 0.2, G: 0.2, B: 0.2, A: 0.7},
	BackgroundColor: f32.Color{R: 0.0, G: 0.2, B: 0.2, A: 0.7},
	BorderWidth:     1,
	Padding:         f32.Padding{5, 5, 5, 5},
	Delay:           2,
}

// Hint is called if the mouse is inside a clickable widget
// i.e. when it is hovered.
func Hint(text string, tag any) {
	if !lib.TagsEqual(CurrentHint.Tag, tag) {
		CurrentHint.Pos = gpu.MousePos
		CurrentHint.Text = text
		CurrentHint.Tag = tag
		CurrentHint.T = time.Now()
	}
	CurrentHint.Active = true
}

// ShowHint is called at the end of the display loop.
// It will show the hint on top of everything else.
func ShowHint(style *HintStyle) {
	if CurrentHint.Tag == nil {
		return
	}
	if style == nil {
		style = &DefaultHintStyle
	}
	if time.Since(CurrentHint.T) > style.Delay && CurrentHint.Active {
		scale := style.FontSize / gpu.InitialSize
		textHeight := (gpu.Fonts[style.FontNo].Ascent + gpu.Fonts[style.FontNo].Descent) * scale * 1.2
		h := textHeight*4 + style.Padding.T + style.Padding.B
		w := textHeight * 30
		x := min(CurrentHint.Pos.X+w, gpu.WindowWidthDp)
		x = max(float32(0), x-w)
		y := min(CurrentHint.Pos.Y+h, gpu.WindowHeightDp)
		y = max(0, y-h)
		gpu.RoundedRect(x, y, w, h, style.CornerRadius, style.BorderWidth, style.BackgroundColor, style.BorderColor)
	}
	CurrentHint.Active = false
}
