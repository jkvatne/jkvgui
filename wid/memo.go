package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

type MemoState struct {
	ScrollState
	NotAtEnd bool
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
	InsidePadding:  f32.Padding{2, 2, 1, 2},
	OutsidePadding: f32.Padding{5, 3, 4, 3},
	FontNo:         gpu.Mono12,
	FontSize:       0.9,
	Color:          theme.OnSurface,
	BorderRole:     theme.Outline,
	BorderWidth:    1.0,
	CornerRadius:   0.0,
	Wrap:           false,
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
		if gpu.DebugWidgets {
			gpu.RoundedRect(ctx.Rect, 0.0, 1.0, f32.Transparent, f32.Red)
		}
		TotalLineCount := len(*text)

		// Find the number of lines from end of text that fit in window
		BottomLineCount := 0
		yBottom := float32(0)
		for yBottom < ctx.Rect.H && BottomLineCount < TotalLineCount {
			wrapedLines := font.Split((*text)[BottomLineCount], ctx.Rect.W, f)
			yBottom += lineHeight * float32(len(wrapedLines))
		}

		// Startline given by Ypos
		n := 0
		i := int(state.Ypos / lineHeight)
		y := ctx.Rect.Y
		dy := state.Ypos - float32(i)*lineHeight
		gpu.Clip(ctx.Rect)
		for y < ctx.Rect.Y+ctx.Rect.H+baseline && i < TotalLineCount {
			// Draw the wraped lines
			Wmax := float32(0)
			if style.Wrap {
				Wmax = ctx.Rect.W
			}
			wrapedLines := font.Split((*text)[i], Wmax, f)
			// Wrap long lines
			for _, line := range wrapedLines {
				f.DrawText(ctx.X, y+baseline-dy, fg, ctx.Rect.W, gpu.LTR, line)
				y += lineHeight
			}
			i++
			n++
		}
		gpu.NoClip()
		sumH := float32(len(*text)) * lineHeight
		DrawVertScrollbar(ctx.Rect, sumH, ctx.Rect.H, &state.ScrollState)
		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
