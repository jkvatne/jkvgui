package wid

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/callback"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

type EditStyle struct {
	FontSize           float32
	FontNo             int
	FontColor          theme.UIRole
	InsideColor        theme.UIRole
	BorderColor        theme.UIRole
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	CursorWidth        float32
	LabelFraction      float32
	LabelRightAdjust   bool
	LabelSpacing       float32
}

var DefaultEdit = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	InsideColor:        theme.Surface,
	BorderColor:        theme.Outline,
	FontColor:          theme.OnSurface,
	OutsidePadding:     f32.Padding{L: 2, T: 3, R: 2, B: 3},
	InsidePadding:      f32.Padding{L: 4, T: 1, R: 2, B: 1},
	BorderWidth:        0.66,
	BorderCornerRadius: 4,
	CursorWidth:        2,
	LabelFraction:      0.5,
	LabelRightAdjust:   true,
	LabelSpacing:       3,
}

type EditState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
	dragging bool
}

var (
	StateMap = make(map[*string]*EditState)
	halfUnit int64
)

func (s *EditStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func EditF32(label string, f32 *float32, action func(), style *EditStyle) Wid {
	ss := fmt.Sprintf("%0.2f", *f32)
	return Edit(label, &ss, action, style)
}

func EditInt(label string, i32 *int, action func(), style *EditStyle) Wid {
	ss := fmt.Sprintf("%d", *i32)
	return Edit(label, &ss, action, style)
}

func Edit(label string, text *string, action func(), style *EditStyle) Wid {
	return func(ctx Ctx) Dim {
		if style == nil {
			style = &DefaultEdit
		}
		if text == nil {
			return Dim{}
		}
		state := StateMap[text]
		if state == nil {
			StateMap[text] = &EditState{}
			state = StateMap[text]
			state.Buffer.Init(*text)
		}
		f := font.Get(style.FontNo, theme.Colors[style.FontColor])

		widRect := ctx.Rect.Inset(style.OutsidePadding)
		frameRect := widRect
		if label != "" {
			frameRect.X += style.LabelFraction * widRect.W
			frameRect.W *= 1 - style.LabelFraction
		}
		valueRect := frameRect.Inset(style.InsidePadding).Reduce(style.BorderWidth)
		labelRect := valueRect
		labelRect.X = widRect.X
		if label != "" {
			labelRect.W = style.LabelFraction * widRect.W
		} else {
			labelRect.W = 0
		}
		fontHeight := f.Height(style.FontSize)
		baseline := f.Baseline(style.FontSize)

		if ctx.Rect.H == 0 {
			return Dim{W: 32, H: fontHeight + style.TotalPaddingY(), Baseline: baseline}
		}

		bg := theme.Colors[style.InsideColor]
		focused := focus.At(ctx.Rect, text)

		if mouse.LeftBtnPressed(widRect) {
			state.SelStart = f.RuneNo(mouse.Pos().X-(valueRect.X), style.FontSize, state.Buffer.String())
			state.dragging = true
			mouse.Lock()
		}
		if mouse.LeftBtnDown() && state.dragging {
			state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), style.FontSize, state.Buffer.String())
			state.SelEnd = max(state.SelStart, state.SelEnd)
		}

		if mouse.LeftBtnReleased(widRect) {
			gpu.Invalidate(0)
			state.dragging = false
			halfUnit = time.Now().UnixMilli() % 333
			focus.Set(text)
			state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), style.FontSize, state.Buffer.String())
		}
		bw := style.BorderWidth
		if focused {
			bw = min(style.BorderWidth*2, style.BorderWidth+2)
			gpu.Invalidate(111 * time.Millisecond)
			if gpu.LastRune != 0 {
				s1 := state.Buffer.Slice(0, state.SelStart)
				s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
				state.Buffer.Init(s1 + string(gpu.LastRune) + s2)
				gpu.LastRune = 0
				state.SelStart++
				state.SelEnd++
			}
			if gpu.LastKey == glfw.KeyBackspace {
				if state.SelStart > 0 {
					state.SelStart = max(0, state.SelStart-1)
					state.SelEnd = state.SelStart
					s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
					s2 := state.Buffer.Slice(state.SelEnd+1, state.Buffer.RuneCount())
					state.Buffer.Init(s1 + s2)
				}
			} else if gpu.LastKey == glfw.KeyDelete {
				s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
				s2 := state.Buffer.Slice(min(state.SelEnd+1, state.Buffer.RuneCount()), state.Buffer.RuneCount())
				state.Buffer.Init(s1 + s2)
				state.SelEnd = state.SelStart
			} else if gpu.LastKey == glfw.KeyRight && callback.LastMods == glfw.ModShift {
				state.SelEnd = min(state.SelEnd+1, state.Buffer.RuneCount())
			} else if gpu.LastKey == glfw.KeyLeft && callback.LastMods == glfw.ModShift {
				if state.SelStart <= state.SelEnd {
					state.SelStart = max(0, state.SelStart-1)
				} else {
					state.SelEnd--
				}
			} else if gpu.LastKey == glfw.KeyLeft {
				state.SelStart = max(0, state.SelStart-1)
				state.SelEnd = state.SelStart
			} else if gpu.LastKey == glfw.KeyRight {
				state.SelStart = min(state.SelStart+1, state.Buffer.RuneCount())
				state.SelEnd = state.SelStart
			} else if gpu.LastKey == glfw.KeyEnd {
				state.SelStart = state.Buffer.RuneCount()
				state.SelEnd = state.SelStart
			} else if gpu.LastKey == glfw.KeyHome {
				state.SelStart = 0
				state.SelEnd = 0
			}
		} else if mouse.Hovered(widRect) {
			bg = theme.Colors[theme.SurfaceContainer]
		}
		state.SelEnd = min(state.SelEnd, state.Buffer.RuneCount())
		state.SelStart = min(state.SelStart, state.Buffer.RuneCount(), state.SelEnd)

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, theme.Colors[style.BorderColor])
		f.SetColor(theme.Colors[style.FontColor])
		labelWidth := f.Width(style.FontSize, label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		// Draw label
		if label != "" {
			gpu.Rect(labelRect, 1, f32.Transparent, f32.LightBlue)
			f.Printf(
				labelRect.X+dx,
				valueRect.Y+baseline,
				style.FontSize,
				labelRect.W,
				label)
		}
		// Draw selected rectangle
		if state.SelStart != state.SelEnd {
			r := valueRect
			r.H--
			r.W = f.Width(style.FontSize, state.Buffer.Slice(state.SelStart, state.SelEnd))
			r.X += f.Width(style.FontSize, state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().Alpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.Printf(
			valueRect.X,
			valueRect.Y+baseline,
			style.FontSize,
			valueRect.W,
			state.Buffer.String())
		f.SetColor(f32.Black)

		dx = f.Width(style.FontSize, state.Buffer.Slice(0, state.SelStart))

		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 1, theme.Colors[theme.Primary])
		}
		return Dim{W: frameRect.W, H: frameRect.H, Baseline: baseline}
	}
}
