package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type MemoState struct {
	ScrollState
}

type MemoStyle struct {
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	BorderRole     theme.UIRole
	BorderWidth    float32
	CornerRadius   float32
	Wrap           bool
	FontNo         int
	FontSize       float32
	Color          theme.UIRole
	Role           theme.UIRole
}

var DefMemo = &MemoStyle{
	InsidePadding:  f32.Padding{L: 2, T: 2, R: 1, B: 2},
	OutsidePadding: f32.Padding{L: 5, T: 3, R: 4, B: 3},
	FontNo:         gpu.Mono12,
	FontSize:       0.9,
	Color:          theme.OnSurface,
	BorderRole:     theme.Outline,
	BorderWidth:    1.0,
	CornerRadius:   0.0,
	Wrap:           true,
}

var MemoStateMap = make(map[any]*MemoState)

func Memo(text *[]string, style *MemoStyle) Wid {
	if style == nil {
		style = DefMemo
	}

	state := MemoStateMap[text]
	if state == nil {
		MemoStateMap[text] = &MemoState{}
		state = MemoStateMap[text]
		// We want to show the last lines by default.
		state.AtEnd = true
	}

	f := font.Fonts[style.FontNo]
	lineHeight := f.Height()
	fg := style.Color.Fg()

	return func(ctx Ctx) Dim {
		baseline := f.Baseline()
		if ctx.Mode != RenderChildren {
			return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
		}
		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		gpu.RoundedRect(ctx.Rect, style.CornerRadius, style.BorderWidth, f32.Transparent, style.BorderRole.Fg())

		ctx.Rect = ctx.Rect.Inset(style.InsidePadding, 0)
		if *gpu.DebugWidgets {
			gpu.RoundedRect(ctx.Rect, 0.0, 1.0, f32.Transparent, f32.Red)
		}
		gpu.Mutex.Lock()
		TotalLineCount := len(*text)
		gpu.Mutex.Unlock()
		Wmax := float32(0)
		if style.Wrap {
			Wmax = ctx.Rect.W
		}
		n := 0
		i := int(state.Ypos / lineHeight)
		y := ctx.Rect.Y
		dy := state.Ypos - float32(i)*lineHeight
		gpu.Clip(ctx.Rect)
		for y < ctx.Rect.Y+ctx.Rect.H+baseline && i < TotalLineCount {
			// Draw the wraped lines
			gpu.Mutex.Lock()
			wrapedLines := font.Split((*text)[i], Wmax, f)
			gpu.Mutex.Unlock()
			// Wrap long lines
			for _, line := range wrapedLines {
				f.DrawText(ctx.X, y+baseline-dy, fg, ctx.Rect.W, gpu.LTR, line)
				y += lineHeight
			}
			i++
			n++
		}
		gpu.NoClip()
		if i >= TotalLineCount && dy < lineHeight {
			state.AtEnd = true
		}
		state.Ymax = float32(TotalLineCount) * lineHeight
		dy = VertScollbarUserInput(ctx.Rect.H, &state.ScrollState)
		state.Ypos += dy
		if state.AtEnd {
			state.Ypos = state.Ymax - ctx.H
		}
		state.Ypos = max(0, min(state.Ypos, state.Ymax-ctx.H))
		DrawVertScrollbar(ctx.Rect, state.Ymax, ctx.H, &state.ScrollState)
		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
