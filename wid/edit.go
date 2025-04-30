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
	"log/slog"
	"strconv"
)

type EditStyle struct {
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
	Dp                 int
	ReadOnly           bool
}

var DefaultEdit = EditStyle{
	FontNo:             gpu.Normal12,
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
	Dp:                 2,
}

const GridBorderWidth = 1

var GridEdit = EditStyle{
	FontNo:        gpu.Normal12,
	EditSize:      0.5,
	Color:         theme.PrimaryContainer,
	BorderColor:   theme.Transparent,
	InsidePadding: f32.Padding{L: 2, T: 0, R: 2, B: 0},
	CursorWidth:   1,
	BorderWidth:   GridBorderWidth,
	Dp:            2,
}

type EditState struct {
	SelStart int
	SelEnd   int
	Buffer   utf8.String
	dragging bool
	modified bool
}

var (
	StateMap = make(map[any]*EditState)
)

func (s *EditStyle) Size(wl, we float32) *EditStyle {
	ss := *s
	ss.EditSize = we
	ss.LabelSize = wl
	return &ss
}

func (s *EditStyle) RO() *EditStyle {
	ss := *s
	ss.ReadOnly = true
	return &ss
}

func (s *EditStyle) TotalPaddingY() float32 {
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B
}

func (s *EditStyle) Top() float32 {
	return s.InsidePadding.T + s.InsidePadding.T
}

func (s *EditStyle) Dim(w float32, f *font.Font) Dim {
	if s.LabelSize > 1.0 || s.EditSize > 1.0 {
		w = s.LabelSize + s.EditSize
	} else if s.EditSize > 0.0 {
		w = s.EditSize
	}
	dim := Dim{W: w, H: f.Height() + s.TotalPaddingY(), Baseline: f.Baseline() + s.Top()}
	return dim
}

func DrawCursor(style *EditStyle, state *EditState, valueRect f32.Rect, f *font.Font) {
	if gpu.BlinkState.Load() {
		dx := f.Width(state.Buffer.Slice(0, state.SelEnd))
		if dx < valueRect.W {
			gpu.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 0.5+valueRect.H/10, style.Color.Fg())
		}
	}
}

func CalculateRects(hasLabel bool, style *EditStyle, r f32.Rect) (f32.Rect, f32.Rect, f32.Rect) {
	frameRect := r.Inset(style.OutsidePadding, 0)
	valueRect := frameRect.Inset(style.InsidePadding, 0)
	labelRect := valueRect
	if !hasLabel {
		labelRect.W = 0
		if style.EditSize > 1.0 {
			// Edit size given in device independent pixels. No label
			frameRect.W = style.EditSize
		} else if style.EditSize == 0.0 {
			// No size given. Use all
		} else {
			// Fractional edit size.
			// frameRect.W = style.EditSize * r.W
		}
	} else {
		// Have label
		ls, es := style.LabelSize, style.EditSize
		if ls == 0.0 && es == 0.0 {
			// No widht given, use 0.5/0.5
			ls, es = 0.5, 0.5
		} else if ls > 1.0 || es > 1.0 {
			// Use fixed sizes
			ls = ls / valueRect.W
			es = es / valueRect.W
		} else if ls == 0.0 && es < 1.0 {
			ls = 1 - es
		} else if es == 0.0 && ls < 1.0 {
			es = 1 - ls
		} else if ls < 1.0 && es < 1.0 {
			// Fractional sizes
		} else {
			panic("Edit can not have both fractional and absolute sizes for label/value")
		}
		labelRect.W = ls * valueRect.W
		valueRect.W = es * valueRect.W
		valueRect.X += labelRect.W
		frameRect.X += labelRect.W
		frameRect.W -= labelRect.W
	}
	frameRect.X -= style.BorderWidth / 2
	frameRect.Y -= style.BorderWidth / 2
	frameRect.W += style.BorderWidth
	frameRect.H += style.BorderWidth
	return frameRect, valueRect, labelRect
}

func ClearBuffers() {
	StateMap = make(map[any]*EditState)
}

func EditText(state *EditState) {
	if gpu.LastRune != 0 {
		p1 := min(state.SelStart, state.SelEnd, state.Buffer.RuneCount())
		p2 := min(max(state.SelStart, state.SelEnd), state.Buffer.RuneCount())
		s1 := state.Buffer.Slice(0, p1)
		s2 := state.Buffer.Slice(p2, state.Buffer.RuneCount())
		state.Buffer.Init(s1 + string(gpu.LastRune) + s2)
		gpu.LastRune = 0
		state.SelStart++
		state.SelEnd = state.SelStart
		state.modified = true
	} else if gpu.LastKey == glfw.KeyBackspace {
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
		state.modified = true
	} else if gpu.LastKey == glfw.KeyDelete {
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
		state.modified = true
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
	} else if gpu.LastKey == glfw.KeyC && sys.LastMods == glfw.ModControl {
		// Copy to clipboard
		glfw.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
	} else if gpu.LastKey == glfw.KeyX && sys.LastMods == glfw.ModControl {
		// Copy to clipboard
		glfw.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
	} else if gpu.LastKey == glfw.KeyV && sys.LastMods == glfw.ModControl {
		// Insert from clipboard
		s1 := state.Buffer.Slice(0, state.SelStart)
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + glfw.GetClipboardString() + s2)
		state.modified = true
	}
	if gpu.LastKey != 0 {
		gpu.Invalidate(0)
	}
}

func EditHandleMouse(state *EditState, valueRect f32.Rect, f *font.Font, value any) {
	if mouse.LeftBtnDoubleClick(valueRect) {
		state.SelStart = f.RuneNo(mouse.Pos().X-(valueRect.X), state.Buffer.String())
		state.SelStart = min(state.SelStart, state.Buffer.RuneCount())
		state.SelEnd = state.SelStart
		for state.SelStart > 0 && state.Buffer.At(state.SelStart-1) != rune(32) {
			state.SelStart--
		}
		for state.SelEnd < state.Buffer.RuneCount() && state.Buffer.At(state.SelEnd) != rune(32) {
			state.SelEnd++
		}
		slog.Info("Doubleclick")
		state.dragging = false

	} else if state.dragging {
		if mouse.LeftBtnDown() {
			state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), state.Buffer.String())
			// slog.Info("Dragging", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
		} else {
			state.SelEnd = f.RuneNo(mouse.Pos().X-(valueRect.X), state.Buffer.String())
			slog.Debug("Drag end", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
			state.dragging = false
			focus.SetFocusedTag(value)
		}
		gpu.Invalidate(0)

	} else if mouse.LeftBtnPressed(valueRect) {
		state.SelStart = f.RuneNo(mouse.Pos().X-(valueRect.X), state.Buffer.String())
		state.SelEnd = state.SelStart
		slog.Info("LeftBtnPressed", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
		state.dragging = true
		focus.SetFocusedTag(value)
		gpu.Invalidate(0)
	}
}

func DrawDebuggingInfo(labelRect f32.Rect, valueRect f32.Rect, WidgetRect f32.Rect) {
	if gpu.DebugWidgets {
		gpu.Rect(WidgetRect, 0.5, f32.Transparent, f32.MultAlpha(f32.Yellow, 0.25))
		gpu.Rect(labelRect, 0.5, f32.Transparent, f32.MultAlpha(f32.Green, 0.25))
		gpu.Rect(valueRect, 0.5, f32.Transparent, f32.MultAlpha(f32.Red, 0.25))
	}
}

func Edit(value any, label string, action func(), style *EditStyle) Wid {
	if style == nil {
		style = &DefaultEdit
	}

	state := StateMap[value]
	if state == nil {
		StateMap[value] = &EditState{}
		state = StateMap[value]
		gpu.GpuMutex.Lock()
		switch v := value.(type) {
		case *int:
			state.Buffer.Init(fmt.Sprintf("%d", *v))
		case *string:
			state.Buffer.Init(fmt.Sprintf("%s", *v))
		case *float32:
			state.Buffer.Init(strconv.FormatFloat(float64(*v), 'f', style.Dp, 32))
		case *float64:
			state.Buffer.Init(strconv.FormatFloat(*v, 'f', style.Dp, 64))
		default:
			f32.Exit("Edit with value that is not *int, *string *float32")
		}
		gpu.GpuMutex.Unlock()
	}

	// Precalculate some values
	f := font.Get(style.FontNo)
	baseline := f.Baseline()
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := style.Dim(ctx.W, f)
		if ctx.Mode != RenderChildren {
			return dim
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)

		labelWidth := f.Width(label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		focused := !style.ReadOnly && focus.At(ctx.Rect, value)
		EditHandleMouse(state, valueRect, f, value)

		if focused {
			bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
			EditText(state)
		} else if state.modified == true {
			// On loss of focus, update the actual values if they have changed
			state.modified = false
			switch v := value.(type) {
			case *int:
				n, err := strconv.Atoi(state.Buffer.String())
				if err == nil {
					gpu.GpuMutex.Lock()
					*v = n
					gpu.GpuMutex.Unlock()
				}
				state.Buffer.Init(fmt.Sprintf("%d", *v))
			case *string:
				gpu.GpuMutex.Lock()
				*v = state.Buffer.String()
				state.Buffer.Init(fmt.Sprintf("%s", *v))
				gpu.GpuMutex.Unlock()
			case *float32:
				f, err := strconv.ParseFloat(state.Buffer.String(), 64)
				if err == nil {
					gpu.GpuMutex.Lock()
					*v = float32(f)
					gpu.GpuMutex.Unlock()
				}
				state.Buffer.Init(strconv.FormatFloat(float64(*v), 'f', style.Dp, 32))
			case *float64:
				f, err := strconv.ParseFloat(state.Buffer.String(), 64)
				if err == nil {
					gpu.GpuMutex.Lock()
					*v = float64(f)
					gpu.GpuMutex.Unlock()
				}
				state.Buffer.Init(strconv.FormatFloat(*v, 'f', style.Dp, 64))
			}
		}

		state.SelEnd = min(state.SelEnd, state.Buffer.RuneCount())
		state.SelStart = min(state.SelStart, state.Buffer.RuneCount())

		// Draw frame around value
		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, f32.Transparent, style.BorderColor.Fg())

		// Draw label if it exists
		if label != "" {
			f.DrawText(labelRect.X+dx, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
		}

		// Draw selected rectangle
		if state.SelStart != state.SelEnd && focused {
			r := valueRect
			r.H--
			p1 := min(state.SelStart, state.SelEnd)
			p2 := max(state.SelStart, state.SelEnd)
			r.W = f.Width(state.Buffer.Slice(p1, p2))
			r.X += f.Width(state.Buffer.Slice(0, p1))
			c := theme.PrimaryContainer.Bg().Alpha(0.8)
			gpu.RoundedRect(r, 0, 0, c, c)
		}

		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if focused {
			DrawCursor(style, state, valueRect, f)
		}

		// Draw debugging rectangles if gpu.DebugWidgets is true
		DrawDebuggingInfo(labelRect, valueRect, ctx.Rect)

		return dim
	}
}
