package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type MemoState struct {
	Xpos     float32
	Ypos     float32
	Width    float32
	Max      float32
	dragging bool
	StartPos float32
	NotAtEnd bool
}

var memoStyle = &LabelStyle{
	Padding:  f32.Padding{5, 3, 1, 2},
	FontNo:   gpu.Mono,
	Color:    theme.OnSurface,
	FontSize: 0.9,
}

var MemoStateMap = make(map[any]*MemoState)

func Memo(text *[]string, style *LabelStyle) Wid {
	if style == nil {
		style = memoStyle
	}

	state := MemoStateMap[text]
	if state == nil {
		MemoStateMap[text] = &MemoState{}
		state = MemoStateMap[text]
	}

	f := font.Fonts[style.FontNo]
	lineHeight := f.Height(style.FontSize)
	fg := style.Color.Fg()

	return func(ctx Ctx) Dim {
		baseline := f.Baseline(style.FontSize) + style.Padding.T
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: ctx.H, Baseline: baseline}
		}

		MemoLineCount := int(ctx.Rect.H / lineHeight)
		TotalLineCount := len(*text)
		// Draw lines
		x := ctx.X + style.Padding.L
		y := ctx.Rect.Y + baseline + float32(MemoLineCount)*lineHeight
		StartLine := int(state.Ypos / lineHeight)
		if !state.NotAtEnd {
			StartLine = max(0, TotalLineCount-MemoLineCount)
		}

		if TotalLineCount < MemoLineCount {
			// Draw from top. Partially full area
			for i := TotalLineCount - 1; i >= 0; i-- {
				line := (*text)[i]
				f.DrawText(x, y, fg, style.FontSize, 0, gpu.LTR, line)
				y -= lineHeight
			}
		} else if StartLine+MemoLineCount < TotalLineCount {
			// Draw middle lines. Start at bottom
			for i := StartLine + MemoLineCount; i >= StartLine; i-- {
				line := (*text)[i]
				f.DrawText(x, y, fg, style.FontSize, 0, gpu.LTR, line)
				y -= lineHeight
			}
		} else {
			// Last lines
			for i := TotalLineCount - 1; i >= TotalLineCount-MemoLineCount; i-- {
				line := (*text)[i]
				f.DrawText(x, y, fg, style.FontSize, 0, gpu.LTR, line)
				y -= lineHeight
			}
		}

		if state.dragging {
			state.Ypos += mouse.Pos().Y - state.StartPos
			state.StartPos = mouse.Pos().Y
			state.dragging = mouse.LeftBtnDown()
		}
		if sys.ScrolledY != 0 {
			if sys.ScrolledY > 0 {
				state.NotAtEnd = true
			}
			state.Ypos -= sys.ScrolledY * lineHeight
			sys.ScrolledY = 0
			gpu.Invalidate(0)
		}
		state.Ypos = min(state.Ypos, lineHeight*float32(len(*text))-ctx.Rect.H)
		state.Ypos = max(state.Ypos, 0)

		if TotalLineCount > MemoLineCount {
			// Draw scrollbar track
			ctx2 := ctx
			ctx2.X += ctx2.W - 8
			ctx2.W = 8
			alpha := float32(0.3)
			if mouse.Hovered(ctx2.Rect) {
				alpha = 0.7
			}
			ctx2.H -= 2
			gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Fg().Alpha(alpha))

			// Draw thumb
			sumH := float32(len(*text)) * lineHeight
			ctx2.Rect.X += 1.0
			ctx2.Rect.W -= 2.0
			state.Ypos = max(0, min(state.Ypos, sumH))
			ctx2.Rect.Y += ctx.Rect.H * state.Ypos / sumH
			ctx2.Rect.H = max(15, ctx2.Rect.H*ctx2.Rect.H/sumH)
			if mouse.LeftBtnPressed(ctx2.Rect) && !state.dragging {
				state.dragging = true
				state.StartPos = mouse.StartDrag().Y
			}
			gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Fg().Alpha(alpha))
		}
		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
