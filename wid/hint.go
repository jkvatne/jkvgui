package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"strings"
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
	FontColor       f32.Color
	CornerRadius    float32
	BorderColor     f32.Color
	BackgroundColor f32.Color
	BorderWidth     float32
	Padding         f32.Padding
	Delay           time.Duration
}

var DefaultHintStyle = HintStyle{
	FontNo:          gpu.DefaultFont,
	FontSize:        gpu.InitialSize * 0.75,
	FontColor:       f32.Color{0.0, 0.0, 0.0, 1.0},
	CornerRadius:    5,
	BorderColor:     f32.Color{R: 0.4, G: 0.4, B: 0.5, A: 1.0},
	BackgroundColor: f32.Color{R: 1.0, G: 1.0, B: 0.9, A: 1.0},
	BorderWidth:     1,
	Padding:         f32.Padding{3, 0, 1, 0},
	Delay:           time.Millisecond * 800,
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

func split(s string, maxWidth float32, font *gpu.Font, scale float32) []string {
	var width float32
	lines := make([]string, 0)
	words := strings.Split(s, " ")
	line := ""
	for _, word := range words {
		if word == "" {
			continue
		}
		width = font.Width(scale, line+" "+word)
		if width <= maxWidth {
			line = line + word + " "
		} else {
			if len(line) > 0 {
				// Use words up to the current word
				lines = append(lines, line)
				line = word + " "
			} else {
				// Hard break a very long word
				for j := len(word) - 1; j >= 1; j-- {
					if font.Width(scale, word[0:j]) > maxWidth {
						line = word[0:j]
						word = word[j:]
						break
					}
				}
				lines = append(lines, word)
			}
		}
	}
	lines = append(lines, line)
	return lines
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
		scale := style.FontSize / 2
		textHeight := (gpu.Fonts[style.FontNo].Ascent + gpu.Fonts[style.FontNo].Descent) * scale * 1.2

		w := textHeight * 8
		x := min(CurrentHint.Pos.X+w+style.Padding.L+style.Padding.R, gpu.WindowWidthDp)
		x = max(float32(0), x-w)

		lines := split(CurrentHint.Text, w-style.Padding.L-style.Padding.R, gpu.Fonts[style.FontNo], scale)
		gpu.Fonts[style.FontNo].SetColor(style.FontColor)

		h := textHeight*float32(len(lines)) + style.Padding.T + style.Padding.B
		y := min(CurrentHint.Pos.Y+h, gpu.WindowHeightDp)
		y = max(0, y-h)
		yb := y + style.Padding.T + textHeight
		r := f32.Rect{x, y, w, h}
		gpu.RoundedRect(r, style.CornerRadius, style.BorderWidth, style.BackgroundColor, style.BorderColor, 5, 0)
		for _, line := range lines {
			gpu.Fonts[style.FontNo].Printf(
				x+style.Padding.L+style.Padding.L+style.BorderWidth,
				yb, style.FontSize,
				0, line)
			yb = yb + style.FontSize
		}
		gpu.Fonts[style.FontNo].SetColor(f32.Black)
	}
	CurrentHint.Active = false
}
