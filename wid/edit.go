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
	Disabler           *bool
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
	LabelSpacing:       2,
	Dp:                 2,
}

const GridBorderWidth = 0.0

var GridEdit = EditStyle{
	FontNo:        gpu.Normal12,
	EditSize:      1.0,
	Color:         theme.PrimaryContainer,
	BorderColor:   theme.Transparent,
	InsidePadding: f32.Padding{L: 2, T: 1, R: 2, B: 1},
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
	value    any
	dp       int
}

var (
	StateMap      = make(map[any]*EditState)
	StateMapMutex sync.RWMutex
)

func (s *EditStyle) Disabled() bool {
	return s.Disabler != nil && *s.Disabler == true
}

func (s *EditStyle) Size(wl, we float32) *EditStyle {
	ss := *s
	ss.EditSize = we
	ss.LabelSize = wl
	return &ss
}

func (s *EditStyle) TopPad(p float32) *EditStyle {
	ss := *s
	ss.OutsidePadding.T = p
	return &ss
}

func (s *EditStyle) RO() *EditStyle {
	ss := *s
	ss.ReadOnly = true
	return &ss
}

func (s *EditStyle) Top() float32 {
	return s.OutsidePadding.T + s.InsidePadding.T + s.BorderWidth
}

func (s *EditStyle) D(flag *bool) *EditStyle {
	ss := *s
	ss.Disabler = flag
	return &ss
}

// Dim wil calculate the dimension of edit/combo/checkbox
// ctx.W is the maximum available space (unless it is 0)
func (s *EditStyle) Dim(ctx Ctx, f *font.Font) Dim {
	var w float32
	px, py := f32.TotalPadding(s.InsidePadding, s.OutsidePadding, s.BorderWidth)
	if s.LabelSize <= 1.0 && s.EditSize <= 1.0 {
		w = ctx.W
	} else if s.LabelSize > 1.0 && s.EditSize > 1.0 {
		// Sizes are given. Use them
		w = s.LabelSize + s.EditSize + s.LabelSpacing
	} else if s.EditSize > 0.0 && s.LabelSize > 0.0 {
		w = s.EditSize
	} else {
		w += px
	}
	h := f.Height + py
	return Dim{W: w, H: h, Baseline: f.Baseline + s.OutsidePadding.T + s.InsidePadding.T + s.BorderWidth}
}

func DrawCursor(ctx Ctx, style *EditStyle, state *EditState, valueRect f32.Rect, f *font.Font) {
	if sys.BlinkState.Load() {
		dx := f.Width(state.Buffer.Slice(0, min(state.Buffer.RuneCount(), state.SelEnd)))
		if dx < valueRect.W {
			ctx.Win.Gd.VertLine(valueRect.X+dx, valueRect.Y, valueRect.Y+valueRect.H, 0.5+valueRect.H/10, style.Color.Fg())
		}
	}
}

// CalculateRects returns frameRect, valueRect, labelRect based on available space in r
func CalculateRects(hasLabel bool, style *EditStyle, r f32.Rect) (dim, frameRect, valueRect, labelRect f32.Rect) {
	_, py := f32.TotalPadding(style.InsidePadding, style.OutsidePadding, style.BorderWidth)
	f := font.Get(style.FontNo)
	r.H = f.Height + py

	frameRect = r.Inset(style.OutsidePadding, 0)
	valueRect = frameRect.Inset(style.InsidePadding, style.BorderWidth)
	labelRect = valueRect
	if !hasLabel {
		labelRect.W = 0
		if style.EditSize > 1.0 {
			// Edit size given in device independent pixels. No label
			frameRect.W = style.EditSize
		} else if style.EditSize == 0.0 {
			// No size given. Use all
		} else {
			// Fractional edit size.
			frameRect.W = style.EditSize * r.W
			r.W = style.EditSize
		}
	} else {
		// Have label
		ls, es := style.LabelSize, style.EditSize
		if ls == 0.0 && es == 0.0 {
			// No width given, use 0.5/0.5
			ls, es = 0.5, 0.5
		} else if ls > 1.0 || es > 1.0 {
			// Use fixed sizes
			ls = ls / r.W
			es = es / r.W
		} else if ls == 0.0 && es < 1.0 {
			ls = 1 - es
		} else if es == 0.0 && ls < 1.0 {
			es = 1 - ls
		} else if ls < 1.0 && es < 1.0 {
			// Fractional sizes
		} else {
			f32.Exit(1, "Edit can not have both fractional and absolute sizes for label/value")
		}
		ls *= valueRect.W
		es *= valueRect.W
		frameRect.X += ls
		frameRect.W = es
		valueRect = frameRect.Inset(style.InsidePadding, style.BorderWidth)
		labelRect.W = ls - (style.InsidePadding.L + style.BorderWidth + style.InsidePadding.R)
		r.W = ls + es
	}
	dim = f32.Rect{X: 0, Y: 0, W: r.W, H: r.H}
	return dim, frameRect, valueRect, labelRect
}

func ClearBuffers() {
	StateMapMutex.Lock()
	defer StateMapMutex.Unlock()
	StateMap = make(map[any]*EditState)
}

func EditText(ctx Ctx, state *EditState, action func()) {
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
	} else if ctx.Win.LastKey == sys.KeyEnter || ctx.Win.LastKey == sys.KeyKPEnter {
		if action != nil {
			updateValue(&ctx, state)
			action()
		}
	}
	if ctx.Win.LastKey != 0 {
		ctx.Win.LastKey = 0
		ctx.Win.Invalidate()
	}
}

func EditMouseHandler(ctx Ctx, state *EditState, valueRect f32.Rect, f *font.Font, value any) {
	state.hovered = false
	if ctx.Win.LeftBtnDoubleClick(valueRect) {
		slog.Debug("EditMouseHandler:")
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
			ctx.Win.SetFocusedTag(value)
		}
		if newPos < state.SelStart {
			state.SelStart = newPos
		} else if newPos > state.SelEnd {
			state.SelEnd = newPos
		}
		ctx.Win.Invalidate()
		state.hovered = true

	} else if ctx.Win.LeftBtnPressed(valueRect) {
		state.SelStart = f.RuneNo(ctx.Win.MousePos().X-(valueRect.X), state.Buffer.String())
		state.SelEnd = state.SelStart
		if !ctx.Win.Dragging {
			ctx.Win.SetFocusedTag(value)
		}
		state.dragging = true
		slog.Debug("Start dragging because left btn pressed in edit.")
		ctx.Win.StartDrag()
		ctx.Win.Invalidate()
		state.hovered = true
	} else if ctx.Win.Hovered(ctx.Rect) {
		state.hovered = true
	}
}

func DrawDebuggingInfo(ctx Ctx, labelRect f32.Rect, valueRect f32.Rect, WidgetRect f32.Rect) {
	if *DebugWidgets {
		ctx.Win.Gd.OutlinedRect(WidgetRect, 1, f32.Magenta.MultAlpha(0.8))
		ctx.Win.Gd.OutlinedRect(labelRect, 1, f32.Green.MultAlpha(0.5))
		ctx.Win.Gd.OutlinedRect(valueRect, 1, f32.Red.MultAlpha(0.25))
	}
}

func updateValue(ctx *Ctx, state *EditState) {
	state.modified = false
	ctx.Win.Mutex.Lock()
	defer ctx.Win.Mutex.Unlock()
	switch v := state.value.(type) {
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
		state.Buffer.Init(strconv.FormatFloat(float64(*v), 'f', state.dp, 32))
	case *float64:
		f, err := strconv.ParseFloat(state.Buffer.String(), 64)
		if err == nil {
			*v = f
		}
		state.Buffer.Init(strconv.FormatFloat(*v, 'f', state.dp, 64))
	}
}

func Edit(value any, label string, action func(), style *EditStyle) Wid {
	if style == nil {
		style = &DefaultEdit
	}
	// Initialize the state of the widget
	StateMapMutex.RLock()
	state := StateMap[value]
	StateMapMutex.RUnlock()

	if state == nil {
		slog.Debug("Edit: Create new state")
		StateMapMutex.Lock()
		StateMap[value] = &EditState{value: value, dp: style.Dp}
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
			f32.Exit(1, "Edit with value that is not *int, *string *float32")
		}
	}

	// Pre-calculate some values
	f := font.Get(style.FontNo)
	fg := style.Color.Fg()
	bw := style.BorderWidth

	return func(ctx Ctx) Dim {
		// dim := style.Dim(ctx, f)
		r, frameRect, valueRect, labelRect := CalculateRects(label != "", style, ctx.Rect)
		ctx.H = min(ctx.H, r.H)
		dim := Dim{W: r.W, H: r.H}
		if ctx.Mode != RenderChildren {
			return dim
		}
		if ctx.H < 0 {
			return Dim{}
		}

		focused := !style.ReadOnly && ctx.Win.At(value)
		if !style.Disabled() {
			if ctx.Win.Focused {
				EditMouseHandler(ctx, state, valueRect, f, value)
			}
			if focused {
				bw = min(style.BorderWidth*1.5, style.BorderWidth+1)
				EditText(ctx, state, action)
			} else if state.modified == true {
				// On loss of focus, update the actual values if they have changed
				updateValue(&ctx, state)
			}
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
			bg = fg.WithAlpha(0.05)
		}
		ctx.Win.Gd.RoundedRect(frameRect, style.BorderCornerRadius, bw, bg, style.BorderColor.Bg())

		// Draw label if it exists
		if label != "" {
			if style.LabelRightAdjust {
				dx := max(0.0, labelRect.W-f.Width(label)-style.LabelSpacing)
				f.DrawText(ctx.Win.Gd, labelRect.X+dx, valueRect.Y+f.Baseline, fg, labelRect.W, gpu.LTR, label)
			} else {
				f.DrawText(ctx.Win.Gd, labelRect.X, valueRect.Y+f.Baseline, fg, labelRect.W, gpu.LTR, label)
			}
		}

		// Draw selected rectangle
		if focused && state.SelStart != state.SelEnd {
			if state.SelStart > state.SelEnd {
				slog.Error("SelStart>SelEnd!")
			} else {
				r := valueRect
				r.W = f.Width(state.Buffer.Slice(state.SelStart, state.SelEnd))
				r.X += f.Width(state.Buffer.Slice(0, state.SelStart))
				ctx.Win.Gd.SolidRect(r, theme.PrimaryContainer.Bg())
			}
		}

		// Draw value
		f.DrawText(ctx.Win.Gd, valueRect.X, valueRect.Y+f.Baseline, fg, valueRect.W, gpu.LTR, state.Buffer.String())

		// Draw cursor
		if !style.ReadOnly && ctx.Win.At(value) {
			DrawCursor(ctx, style, state, valueRect, f)
			if !ctx.Win.Blinking.Load() {
				ctx.Win.Blinking.Store(true)
			}
		}

		// Draw debugging rectangles if gpu.DebugWidgets is true
		DrawDebuggingInfo(ctx, labelRect, valueRect, ctx.Rect)

		return dim
	}
}
