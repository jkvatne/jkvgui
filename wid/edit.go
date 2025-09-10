package wid

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	utf8 "golang.org/x/exp/utf8string"
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
	OutsidePadding:     f32.Padding{L: 2, T: 2, R: 2, B: 2},
	InsidePadding:      f32.Padding{L: 2, T: 2, R: 2, B: 2},
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

func (s *EditStyle) Dim(ctx *Ctx, f *font.Font) Dim {
	w := ctx.W
	if s.LabelSize > 1.0 || s.EditSize > 1.0 {
		w = s.LabelSize + s.EditSize
	} else if s.EditSize > 0.0 {
		w = s.EditSize
	}
	dim := Dim{W: w, H: f.Height + s.TotalPaddingY(), Baseline: f.Baseline + s.Top()}
	if ctx.H < dim.H {
		ctx.H = dim.H
	}
	return dim
}

func DrawCursor(style *EditStyle, state *EditState, valueRect f32.Rect, f *font.Font) {
	if sys.CurrentInfo.BlinkState.Load() {
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
			// No width given, use 0.5/0.5
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
			f32.Exit("Edit can not have both fractional and absolute sizes for label/value")
		}
		ls *= valueRect.W
		es *= valueRect.W
		frameRect.X += ls
		frameRect.W = es
		valueRect = frameRect.Inset(style.InsidePadding, style.BorderWidth)
		labelRect.W = ls
		labelRect.W = ls - (style.InsidePadding.L + style.BorderWidth + style.InsidePadding.R)
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
	if sys.LastRune != 0 {
		p1 := min(state.SelStart, state.SelEnd, state.Buffer.RuneCount())
		p2 := min(max(state.SelStart, state.SelEnd), state.Buffer.RuneCount())
		s1 := state.Buffer.Slice(0, p1)
		s2 := state.Buffer.Slice(p2, state.Buffer.RuneCount())
		state.Buffer.Init(s1 + string(sys.LastRune) + s2)
		sys.LastRune = 0
		state.SelStart++
		state.SelEnd = state.SelStart
		state.modified = true
	} else if sys.LastKey == sys.KeyBackspace {
		if state.SelStart == state.SelEnd && state.SelStart > 0 {
			// Delete single char backwards
			state.SelStart--
			state.SelEnd--
			s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
			s2 := state.Buffer.Slice(state.SelEnd+1, state.Buffer.RuneCount())
			state.Buffer.Init(s1 + s2)
		} else if state.SelStart > 0 && state.SelStart < state.SelEnd {
			// Delete multiple characters backwards
			s1 := state.Buffer.Slice(0, state.SelStart)
			s2 := state.Buffer.Slice(state.SelEnd, state.Buffer.RuneCount())
			state.Buffer.Init(s1 + s2)
			state.SelEnd = state.SelStart
		}
		state.modified = true
	} else if sys.LastKey == sys.KeyDelete {
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
		state.modified = true
	} else if sys.LastKey == sys.KeyRight && sys.LastMods == sys.ModShift {
		state.SelEnd = min(state.SelEnd+1, state.Buffer.RuneCount())
	} else if sys.LastKey == sys.KeyLeft && sys.LastMods == sys.ModShift {
		state.SelStart = max(0, state.SelStart-1)
	} else if sys.LastKey == sys.KeyLeft {
		state.SelStart = max(0, state.SelStart-1)
		state.SelEnd = state.SelStart
	} else if sys.LastKey == sys.KeyRight {
		state.SelStart = min(state.SelStart+1, state.Buffer.RuneCount())
		state.SelEnd = state.SelStart
	} else if sys.LastKey == sys.KeyEnd {
		state.SelEnd = state.Buffer.RuneCount()
		if sys.LastMods != sys.ModShift {
			state.SelStart = state.SelEnd
		}
	} else if sys.LastKey == sys.KeyHome {
		state.SelStart = 0
		if sys.LastMods != sys.ModShift {
			state.SelEnd = 0
		}
	} else if sys.LastKey == sys.KeyC && sys.LastMods == sys.ModControl {
		// Copy to clipboard
		sys.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
	} else if sys.LastKey == sys.KeyX && sys.LastMods == sys.ModControl {
		// Copy to clipboard
		sys.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
	} else if sys.LastKey == sys.KeyV && sys.LastMods == sys.ModControl {
		// Insert from clipboard
		s1 := state.Buffer.Slice(0, state.SelStart)
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		s3, _ := sys.GetClipboardString()
		state.Buffer.Init(s1 + s3 + s2)
		state.modified = true
	}
	if sys.LastKey != 0 {
		sys.Invalidate()
	}
}

func EditHandleMouse(state *EditState, valueRect f32.Rect, f *font.Font, value any, focused bool) {
	if sys.LeftBtnDoubleClick(valueRect) {
		state.SelStart = f.RuneNo(sys.Pos().X-(valueRect.X), state.Buffer.String())
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
		newPos := f.RuneNo(sys.Pos().X-(valueRect.X), state.Buffer.String())
		if sys.LeftBtnDown() {
			if newPos != state.SelEnd && newPos != state.SelStart {
				slog.Info("Dragging", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
			}
		} else {
			slog.Info("Drag end", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
			state.dragging = false
			sys.SetFocusedTag(value)
		}
		if newPos < state.SelStart {
			state.SelStart = newPos
		} else if newPos > state.SelEnd {
			state.SelEnd = newPos
		}
		sys.Invalidate()

	} else if sys.LeftBtnPressed(valueRect) && !state.dragging {
		state.SelStart = f.RuneNo(sys.Pos().X-(valueRect.X), state.Buffer.String())
		state.SelEnd = state.SelStart
		slog.Info("LeftBtnPressed", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
		state.dragging = true
		sys.SetFocusedTag(value)
		sys.Invalidate()
	}
}

func DrawDebuggingInfo(labelRect f32.Rect, valueRect f32.Rect, WidgetRect f32.Rect) {
	if *DebugWidgets {
		gpu.Rect(WidgetRect, 0.5, f32.Transparent, f32.Yellow.MultAlpha(0.25))
		gpu.Rect(labelRect, 0.5, f32.Transparent, f32.Green.MultAlpha(0.25))
		gpu.Rect(valueRect, 0.5, f32.Transparent, f32.Red.MultAlpha(0.25))
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
		sys.CurrentInfo.Mutex.Lock()
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
		sys.CurrentInfo.Mutex.Unlock()
	}

	// Pre-calculate some values
	f := font.Get(style.FontNo)
	baseline := f.Baseline
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := style.Dim(&ctx, f)
		if ctx.Mode != RenderChildren {
			return dim
		}

		frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)

		labelWidth := f.Width(label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		focused := !style.ReadOnly && sys.At(ctx.Rect, value)
		if sys.CurrentInfo.Focused {
			EditHandleMouse(state, valueRect, f, value, focused)
		}

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
					sys.CurrentInfo.Mutex.Lock()
					*v = n
					sys.CurrentInfo.Mutex.Unlock()
				}
				state.Buffer.Init(fmt.Sprintf("%d", *v))
			case *string:
				sys.CurrentInfo.Mutex.Lock()
				*v = state.Buffer.String()
				state.Buffer.Init(fmt.Sprintf("%s", *v))
				sys.CurrentInfo.Mutex.Unlock()
			case *float32:
				f, err := strconv.ParseFloat(state.Buffer.String(), 64)
				if err == nil {
					sys.CurrentInfo.Mutex.Lock()
					*v = float32(f)
					sys.CurrentInfo.Mutex.Unlock()
				}
				state.Buffer.Init(strconv.FormatFloat(float64(*v), 'f', style.Dp, 32))
			case *float64:
				f, err := strconv.ParseFloat(state.Buffer.String(), 64)
				if err == nil {
					sys.CurrentInfo.Mutex.Lock()
					*v = float64(f)
					sys.CurrentInfo.Mutex.Unlock()
				}
				state.Buffer.Init(strconv.FormatFloat(*v, 'f', style.Dp, 64))
			}
		}

		cnt := state.Buffer.RuneCount()
		if state.SelEnd > cnt {
			state.SelEnd = cnt
		}
		if state.SelStart > cnt {
			state.SelStart = cnt
		}

		// Draw frame around value
		gpu.RoundedRect(frameRect, style.BorderCornerRadius, bw, f32.Transparent, style.BorderColor.Fg())

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				f.DrawText(labelRect.X+dx, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			} else {
				f.DrawText(labelRect.X, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			}
		}

		// Draw selected rectangle
		if focused && state.SelStart != state.SelEnd {
			if state.SelStart > state.SelEnd {
				slog.Info("Selstart>Selend!")
			} else {
				r := valueRect
				r.H--
				r.W = f.Width(state.Buffer.Slice(state.SelStart, state.SelEnd))
				r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
				c := theme.PrimaryContainer.Bg().MultAlpha(0.8)
				gpu.Rect(r, 0, c, c)
			}
		}

		// Draw value
		f.DrawText(valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if !style.ReadOnly && sys.At(ctx.Rect, value) {
			DrawCursor(style, state, valueRect, f)
			if !sys.CurrentInfo.Blinking.Load() {
				sys.CurrentInfo.Blinking.Store(true)
			}
		}

		// Draw debugging rectangles if gpu.DebugWidgets is true
		DrawDebuggingInfo(labelRect, valueRect, ctx.Rect)

		return dim
	}
}
