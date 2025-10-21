package wid

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

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
	EditSize:      1.0,
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
	hovered  bool
}

var (
	StateMap      = make(map[any]*EditState)
	StateMapMutex sync.RWMutex
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
	return s.InsidePadding.T + s.InsidePadding.B + s.OutsidePadding.T + s.OutsidePadding.B + 2*s.BorderWidth
}

func (s *EditStyle) Top() float32 {
	return s.InsidePadding.T + s.InsidePadding.T + s.BorderWidth
}

// Dim wil calculate the dimension of edit/combo/checkbox
// Width is distributed between the label and the widget itself
func (s *EditStyle) Dim(w float32, f *font.Font) Dim {
	if s.LabelSize > 1.0 || s.EditSize > 1.0 {
		w = s.LabelSize + s.EditSize
	} else if s.EditSize > 0.0 {
		w = s.EditSize
	}
	h := f.Height + s.TotalPaddingY()
	return Dim{W: w, H: h, Baseline: f.Baseline + s.Top() + s.BorderWidth}
}

func DrawCursor(ctx Ctx, style *EditStyle, state *EditState, valueRect f32.Rect, f *font.Font) {
	if sys.BlinkState.Load() {
		dx := f.Width(state.Buffer.Slice(0, state.SelEnd))
		if dx < valueRect.W {
			ctx.Win.Gd.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 0.5+valueRect.H/10, style.Color.Fg())
		}
	}
}

// CalculateRects returns frameRect, valueRect, labelRect based on available spce in r
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
	StateMapMutex.Lock()
	defer StateMapMutex.Unlock()
	StateMap = make(map[any]*EditState)
}

func EditText(ctx Ctx, state *EditState) {
	if ctx.Win.LastRune != 0 {
		p1 := min(state.SelStart, state.SelEnd, state.Buffer.RuneCount())
		p2 := min(max(state.SelStart, state.SelEnd), state.Buffer.RuneCount())
		s1 := state.Buffer.Slice(0, p1)
		s2 := state.Buffer.Slice(p2, state.Buffer.RuneCount())
		state.Buffer.Init(s1 + string(ctx.Win.LastRune) + s2)
		ctx.Win.LastRune = 0
		state.SelStart++
		state.SelEnd = state.SelStart
		state.modified = true
	} else if ctx.Win.LastKey == sys.KeyBackspace {
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
	} else if ctx.Win.LastKey == sys.KeyDelete {
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
		state.modified = true
	} else if ctx.Win.LastKey == sys.KeyRight && ctx.Win.LastMods == sys.ModShift {
		state.SelEnd = min(state.SelEnd+1, state.Buffer.RuneCount())
	} else if ctx.Win.LastKey == sys.KeyLeft && ctx.Win.LastMods == sys.ModShift {
		state.SelStart = max(0, state.SelStart-1)
	} else if ctx.Win.LastKey == sys.KeyLeft {
		state.SelStart = max(0, state.SelStart-1)
		state.SelEnd = state.SelStart
	} else if ctx.Win.LastKey == sys.KeyRight {
		state.SelStart = min(state.SelStart+1, state.Buffer.RuneCount())
		state.SelEnd = state.SelStart
	} else if ctx.Win.LastKey == sys.KeyEnd {
		state.SelEnd = state.Buffer.RuneCount()
		if ctx.Win.LastMods != sys.ModShift {
			state.SelStart = state.SelEnd
		}
	} else if ctx.Win.LastKey == sys.KeyHome {
		state.SelStart = 0
		if ctx.Win.LastMods != sys.ModShift {
			state.SelEnd = 0
		}
	} else if ctx.Win.LastKey == sys.KeyC && ctx.Win.LastMods == sys.ModControl {
		// Copy to clipboard
		sys.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
	} else if ctx.Win.LastKey == sys.KeyX && ctx.Win.LastMods == sys.ModControl {
		// Copy to clipboard
		sys.SetClipboardString(state.Buffer.Slice(state.SelStart, state.SelEnd))
		s1 := state.Buffer.Slice(0, max(state.SelStart, 0))
		if state.SelEnd == state.SelStart {
			state.SelEnd++
		}
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		state.Buffer.Init(s1 + s2)
		state.SelEnd = state.SelStart
	} else if ctx.Win.LastKey == sys.KeyV && ctx.Win.LastMods == sys.ModControl {
		// Insert from clipboard
		s1 := state.Buffer.Slice(0, state.SelStart)
		s2 := state.Buffer.Slice(min(state.SelEnd, state.Buffer.RuneCount()), state.Buffer.RuneCount())
		s3, _ := sys.GetClipboardString()
		state.Buffer.Init(s1 + s3 + s2)
		state.modified = true
	}
	if ctx.Win.LastKey != 0 {
		ctx.Win.Invalidate()
	}
}

func EditMouseHandler(ctx Ctx, state *EditState, valueRect f32.Rect, f *font.Font, value any) {
	state.hovered = false
	if ctx.Win.LeftBtnDoubleClick(valueRect) {
		state.SelStart = f.RuneNo(ctx.Win.MousePos().X-(valueRect.X), state.Buffer.String())
		state.SelStart = min(state.SelStart, state.Buffer.RuneCount())
		state.SelEnd = state.SelStart
		for state.SelStart > 0 && state.Buffer.At(state.SelStart-1) != rune(32) {
			state.SelStart--
		}
		for state.SelEnd < state.Buffer.RuneCount() && state.Buffer.At(state.SelEnd) != rune(32) {
			state.SelEnd++
		}
		state.dragging = false

	} else if state.dragging {
		newPos := f.RuneNo(ctx.Win.MousePos().X-(valueRect.X), state.Buffer.String())
		if ctx.Win.LeftBtnDown() {
			if newPos != state.SelEnd && newPos != state.SelStart {
				slog.Debug("Dragging", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
			}
		} else {
			slog.Debug("Drag end", "SelStart", state.SelStart, "SelEnd", state.SelEnd)
			state.dragging = false
			// ctx.Win.SetFocusedTag(value)
		}
		if newPos < state.SelStart {
			state.SelStart = newPos
		} else if newPos > state.SelEnd {
			state.SelEnd = newPos
		}
		ctx.Win.Invalidate()
		state.hovered = true

	} else if ctx.Win.LeftBtnPressed(ctx.Rect) {
		state.SelStart = f.RuneNo(ctx.Win.MousePos().X-(valueRect.X), state.Buffer.String())
		state.SelEnd = state.SelStart
		if !ctx.Win.Dragging {
			ctx.Win.SetFocusedTag(value)
		}
		state.dragging = true
		ctx.Win.StartDrag()
		ctx.Win.Invalidate()
		state.hovered = true
	} else if ctx.Win.Hovered(ctx.Rect) {
		state.hovered = true
	}
}

func DrawDebuggingInfo(ctx Ctx, labelRect f32.Rect, valueRect f32.Rect, WidgetRect f32.Rect) {
	if *DebugWidgets {
		ctx.Win.Gd.OutlinedRect(WidgetRect, 0.5, f32.Yellow.MultAlpha(0.25))
		ctx.Win.Gd.OutlinedRect(labelRect, 0.5, f32.Green.MultAlpha(0.25))
		ctx.Win.Gd.OutlinedRect(valueRect, 0.5, f32.Red.MultAlpha(0.25))
	}
}

func updateValues(ctx *Ctx, state *EditState, style *EditStyle, value any) {
	state.modified = false
	ctx.Win.Mutex.Lock()
	defer ctx.Win.Mutex.Unlock()
	switch v := value.(type) {
	case *int:
		n, err := strconv.Atoi(state.Buffer.String())
		if err == nil {
			*v = n
		}
		state.Buffer.Init(fmt.Sprintf("%d", *v))
	case *string:
		*v = state.Buffer.String()
		state.Buffer.Init(fmt.Sprintf("%s", *v))
	case *float32:
		f, err := strconv.ParseFloat(state.Buffer.String(), 64)
		if err == nil {
			*v = float32(f)
		}
		state.Buffer.Init(strconv.FormatFloat(float64(*v), 'f', style.Dp, 32))
	case *float64:
		f, err := strconv.ParseFloat(state.Buffer.String(), 64)
		if err == nil {
			*v = f
		}
		state.Buffer.Init(strconv.FormatFloat(*v, 'f', style.Dp, 64))
	}
}

func Edit(value any, label string, action func(), style *EditStyle) Wid {
	if style == nil {
		style = &DefaultEdit
	}
	StateMapMutex.RLock()
	state := StateMap[value]
	StateMapMutex.RUnlock()
	if state == nil {
		StateMapMutex.Lock()
		StateMap[value] = &EditState{}
		state = StateMap[value]
		StateMapMutex.Unlock()
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
	}

	// Pre-calculate some values
	f := font.Get(style.FontNo)
	baseline := f.Baseline
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		dim := style.Dim(ctx.Rect.W, f)
		ctx.H = min(ctx.H, dim.H)
		if ctx.Mode != RenderChildren {
			return dim
		}
		if ctx.H < 0 {
			return Dim{}
		}
		frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)

		labelWidth := f.Width(label) + style.LabelSpacing + 1
		dx := float32(0)
		if style.LabelRightAdjust {
			dx = max(0.0, labelRect.W-labelWidth-style.LabelSpacing)
		}
		if ctx.Win.Focused {
			EditMouseHandler(ctx, state, valueRect, f, value)
		}
		focused := !style.ReadOnly && ctx.Win.At(value)
		if focused {
			bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
			EditText(ctx, state)
		} else if state.modified == true {
			// On loss of focus, update the actual values if they have changed
			updateValues(&ctx, state, style, value)
		}

		cnt := state.Buffer.RuneCount()
		if state.SelEnd > cnt {
			state.SelEnd = cnt
		}
		if state.SelStart > cnt {
			state.SelStart = cnt
		}

		// Draw frame around value with gray background when hovered
		bg := f32.Transparent
		if state.hovered {
			bg = fg.MultAlpha(0.05)
		}
		ctx.Win.Gd.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, style.BorderColor.Fg())

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				f.DrawText(ctx.Win.Gd, labelRect.X+dx, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			} else {
				f.DrawText(ctx.Win.Gd, labelRect.X, valueRect.Y+baseline, fg, labelRect.W, gpu.LTR, label)
			}
		}

		// Draw selected rectangle
		if focused && state.SelStart != state.SelEnd {
			if state.SelStart > state.SelEnd {
				slog.Error("Selstart>Selend!")
			} else {
				r := valueRect
				r.W = f.Width(state.Buffer.Slice(state.SelStart, state.SelEnd))
				r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
				ctx.Win.Gd.SolidRect(r, theme.PrimaryContainer.Bg())
			}
		}

		// Draw value
		f.DrawText(ctx.Win.Gd, valueRect.X, valueRect.Y+baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if !style.ReadOnly && ctx.Win.At(value) {
			DrawCursor(ctx, style, state, valueRect, f)
			if !ctx.Win.Blinking.Load() {
				ctx.Win.Blinking.Store(true)
			}
		}

		// Draw debugging rectangles if gpu.DebugWidgets is true
		DrawDebuggingInfo(ctx, labelRect, valueRect, ctx.Rect)
		if action != nil {
			action()
		}
		return dim
	}
}
