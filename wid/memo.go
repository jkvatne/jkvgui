package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
)

type MemoState struct {
	Xpos      float32
	Ypos      float32
	Width     float32
	Max       float32
	dragging  bool
	StartLine int
	StartPos  float32
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

		// Draw lines
		x := ctx.X + style.Padding.L
		y := ctx.Rect.Y + baseline
		i := state.StartLine
		h := float32(0)
		for h < ctx.H && i < len(*text) {
			line := (*text)[i]
			f.DrawText(x, y, fg, style.FontSize, 0, gpu.LTR, line)
			y += lineHeight
			h += lineHeight
			i++
		}

		// Draw scrollbar
		ctx2 := ctx
		ctx2.X += ctx2.W - 8
		ctx2.W = 8
		alpha := float32(0.4)
		if mouse.Hovered(ctx2.Rect) {
			alpha = 1.0
		}
		gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Bg().Alpha(alpha))
		// Draw thumb
		sumH := float32(len(*text)) * lineHeight
		ctx2.Rect.X += 1.0
		ctx2.Rect.W -= 2.0
		ctx2.Rect.Y = state.Ypos * ctx.Rect.H / sumH
		ctx2.Rect.H *= ctx2.Rect.H / sumH
		if mouse.LeftBtnPressed(ctx2.Rect) && !state.dragging {
			state.dragging = true
			state.StartPos = mouse.StartDrag().Y
		}
		gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Fg().Alpha(alpha))
		if state.dragging {
			state.Ypos += mouse.Pos().Y - state.StartPos
			state.StartPos = mouse.Pos().Y
			state.dragging = mouse.LeftBtnDown()
		}

		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
