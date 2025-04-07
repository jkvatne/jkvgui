package wid

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
	"time"
)

type EditStyle struct {
	FontSize           float32
	FontNo             int
	Color              theme.UIRole
	BorderColor        theme.UIRole
	BorderWidth        float32
	BorderCornerRadius float32
	InsidePadding      f32.Padding
	OutsidePadding     f32.Padding
	CursorWidth        float32
	EditSize           float32
	LabelSize          float32
	LabelRightAdjust   bool
	LabelSpacing       float32
}

var DefaultEdit = EditStyle{
	FontSize:           1.0,
	FontNo:             gpu.Normal,
	Color:              theme.Surface,
	BorderColor:        theme.Outline,
	OutsidePadding:     f32.Padding{L: 5, T: 5, R: 5, B: 5},
	InsidePadding:      f32.Padding{L: 4, T: 3, R: 2, B: 2},
	BorderWidth:        1,
	BorderCornerRadius: 4,
	CursorWidth:        2,
	EditSize:           0.0,
	LabelSize:          0.0,
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
	StateMap = make(map[any]*EditState)
	halfUnit int64
)

func (s *EditStyle) Size(wl, we float32) *EditStyle {
	ss := *s
	ss.EditSize = we
	ss.LabelSize = wl
	return &ss
}

func (s *EditStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func (s *EditStyle) Top() float32 {
	return s.InsidePadding.T + s.InsidePadding.T + s.BorderWidth
}

func CalculateRects(hasLabel bool, style *EditStyle, r f32.Rect) (f32.Rect, f32.Rect, f32.Rect) {
	widRect := f32.Rect{r.X, r.Y, r.W, r.H}
	widRect = widRect.Inset(style.OutsidePadding, style.BorderWidth)
	frameRect := widRect
	labelRect := widRect
	if !hasLabel {
		labelRect.W = 0
		if style.EditSize > 1.0 {
			// Edit size given in device independent pixels. No label
			frameRect.W = style.EditSize
		} else if style.EditSize == 0.0 {
			// No size given. Use all
		} else {
			// Fractional edit size.
			frameRect.W *= style.EditSize
		}
	} else {
		// Have label
		ls, es := style.LabelSize, style.EditSize
		if ls == 0.0 && es == 0.0 {
			// No widht given, use 0.5/0.5
			ls, es = 0.5, 0.5
		} else if ls > 1.0 || es > 1.0 {
			// Use fixed sizes
			ls = ls / widRect.W
			es = es / widRect.W
		} else if ls == 0.0 && es < 1.0 {
			ls = 1 - es
		} else if es == 0.0 && ls < 1.0 {
			es = 1 - ls
		} else if ls < 1.0 && es < 1.0 {
			// Fractional sizes
		} else {
			panic("Edit can not have both fractional and absolute sizes for label/value")
		}
		labelRect.W = ls * widRect.W
		frameRect.W = es * widRect.W
		frameRect.X += labelRect.W
	}
	valueRect := frameRect.Inset(style.InsidePadding, 0)
	return frameRect, valueRect, labelRect
}

func EditText(state *EditState) {
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
		if state.SelStart > 0 && state.SelStart == state.SelEnd {
			state.SelStart--
			state.SelEnd = state.SelStart
			s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
			s2 := state.Buffer.Slice(state.SelEnd+1, state.Buffer.RuneCount())
			state.Buffer.Init(s1 + s2)
		} else if state.SelStart > 0 && state.SelStart < state.SelEnd {
			s1 := state.Buffer.Slice(0, state.SelStart)
			s2 := state.Buffer.Slice(state.SelEnd, state.Buffer.RuneCount())
			state.Buffer.Init(s1 + s2)
		}
	} else if gpu.LastKey == glfw.KeyDelete {
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
	} else if gpu.LastKey == glfw.KeyRight && sys.LastMods == glfw.ModShift {
		state.SelEnd = min(state.SelEnd+1, state.Buffer.RuneCount())
	} else if gpu.LastKey == glfw.KeyLeft && sys.LastMods == glfw.ModShift {
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
	state.SelEnd = min(state.SelEnd, state.Buffer.RuneCount())
	state.SelStart = min(state.SelEnd, state.SelStart)
	if gpu.LastKey != 0 {
		gpu.Invalidate(0)
	}
}

func EditHandleMouse(state *EditState, valueRect f32.Rect, f *font.Font, fontSize float32, value any) {
	if mouse.LeftBtnPressed(valueRect) {
		state.SelStart = f.RuneNo(mouse.Pos().X-(valueRect.X), fontSize, state.Buffer.String())
		state.dragging = true
		mouse.StartDrag()
	}
	if mouse.LeftBtnDown() && state.dragging {
		state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), fontSize, state.Buffer.String())
		state.SelEnd = max(state.SelStart, state.SelEnd)
	}

	if mouse.LeftBtnClick(valueRect) {
		gpu.Invalidate(0)
		state.dragging = false
		halfUnit = time.Now().UnixMilli() % 333
		focus.Set(value)
		state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), fontSize, state.Buffer.String())
	}
}

func Edit(value any, label string, action func(), style *EditStyle) Wid {
	if style == nil {
		style = &DefaultEdit
	}
	f32.ExitIf(value == nil, "Edit with nil value")
	state := StateMap[value]
	if state == nil {
		StateMap[value] = &EditState{}
		state = StateMap[value]
		switch v := value.(type) {
		case *int:
			state.Buffer.Init(fmt.Sprintf("%d", *v))
		case *string:
			state.Buffer.Init(fmt.Sprintf("%s", *v))
		case *float32:
			state.Buffer.Init(fmt.Sprintf("%f", *v))
		}
	}
	f := font.Get(style.FontNo)
	fontHeight := f.Height(style.FontSize)
	baseline := f.Baseline(style.FontSize)
	bg := style.Color.Bg()
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := Dim{W: ctx.W, H: fontHeight + style.TotalPaddingY(), Baseline: baseline + style.Top()}
		if ctx.Mode != RenderChildren {
			return dim
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)

		focused := focus.At(ctx.Rect, value)
		EditHandleMouse(state, valueRect, f, style.FontSize, value)
		if focused {
			bw = min(style.BorderWidth*2, style.BorderWidth+2)
			EditText(state)
		} else if mouse.Hovered(ctx.Rect) {
			bg = style.Color.Bg().Mute(0.8)
		}
		state.SelEnd = min(state.SelEnd, state.Buffer.RuneCount())
		state.SelStart = min(state.SelStart, state.Buffer.RuneCount(), state.SelEnd)

		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, style.BorderColor.Fg())
		labelWidth := f.Width(style.FontSize, label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}

		// Draw label
		if label != "" {
			f.DrawText(labelRect.X+dx, valueRect.Y+baseline, fg, style.FontSize, labelRect.W, gpu.LTR, label)
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused {
			r := valueRect
			r.H--
			r.W = f.Width(style.FontSize, state.Buffer.Slice(state.SelStart, state.SelEnd))
			r.X += f.Width(style.FontSize, state.Buffer.Slice(0, state.SelStart))
			c := theme.PrimaryContainer.Bg().Alpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}
		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, style.FontSize,
			valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused && (time.Now().UnixMilli()-halfUnit)/333&1 == 1 {
			dx = f.Width(style.FontSize, state.Buffer.Slice(0, state.SelEnd))
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 1, style.Color.Fg())
		}

		if gpu.DebugWidgets {
			gpu.Rect(labelRect, 1, f32.Transparent, f32.LightBlue)
			gpu.Rect(valueRect, 1, f32.Transparent, f32.LightRed)
			gpu.Rect(ctx.Rect, 1, f32.Transparent, f32.Yellow)
		}

		return dim
	}
}
