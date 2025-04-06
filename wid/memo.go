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

type MemoStyle struct {
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	BorderRole     theme.UIRole
	BorderWidth    float32
	CornerRadius   float32

	FontNo   int
	FontSize float32
	Color    theme.UIRole
	Role     theme.UIRole
}

var DefMemo = &MemoStyle{
	InsidePadding:  f32.Padding{5, 3, 1, 4},
	OutsidePadding: f32.Padding{5, 3, 4, 3},
	FontNo:         gpu.Mono,
	FontSize:       0.9,
	Color:          theme.OnSurface,
	BorderRole:     theme.Outline,
	BorderWidth:    1.0,
	CornerRadius:   5.0,
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
	lineHeight := f.Height(style.FontSize)
	fg := style.Color.Fg()

	return func(ctx Ctx) Dim {
		baseline := f.Baseline(style.FontSize)
		if ctx.Mode != RenderChildren {
			return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
		}
		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		gpu.RoundedRect(ctx.Rect, style.CornerRadius, style.BorderWidth, f32.Transparent, style.BorderRole.Fg())

		ctx.Rect = ctx.Rect.Inset(style.InsidePadding, 0)
		if gpu.DebugWidgets {
			gpu.RoundedRect(ctx.Rect, 0.0, 1.0, f32.Transparent, f32.Red)
		}
		MemoLineCount := int(ctx.Rect.H / lineHeight)
		TotalLineCount := len(*text)

		if TotalLineCount < MemoLineCount {
			y := ctx.Rect.Y + baseline
			// Draw from top. Partially full area
			for i := 0; i < TotalLineCount; i++ {
				line := (*text)[i]
				f.DrawText(ctx.X, y, fg, style.FontSize, 0, gpu.LTR, line)
				y += lineHeight
			}
		} else if state.NotAtEnd {
			// Draw middle lines.
			// Start at bottom
			y := ctx.Rect.Y + ctx.Rect.H - lineHeight + baseline
			EndLine := int((state.Ypos + ctx.Rect.H) / lineHeight)
			EndLine = min(EndLine, TotalLineCount-1)
			if EndLine == TotalLineCount-1 {
				state.NotAtEnd = false
			}
			StartLine := max(0, EndLine-MemoLineCount+1)
			for i := EndLine; i >= StartLine; i-- {
				line := (*text)[i]
				f.DrawText(ctx.X, y, fg, style.FontSize, 0, gpu.LTR, line)
				y -= lineHeight
			}
		} else {
			// Last lines
			y := ctx.Rect.Y + ctx.Rect.H - lineHeight + baseline
			for i := TotalLineCount - 1; i >= TotalLineCount-MemoLineCount; i-- {
				line := (*text)[i]
				f.DrawText(ctx.X, y, fg, style.FontSize, 0, gpu.LTR, line)
				y -= lineHeight
			}
			state.Ypos = float32(TotalLineCount) * lineHeight
		}

		if state.dragging {
			// Mouse dragging scroller thumb
			dy := (mouse.Pos().Y - state.StartPos) / ctx.H * float32(TotalLineCount) * lineHeight
			if dy < 0 {
				state.NotAtEnd = true
			}
			state.Ypos += dy
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
			alpha := float32(0.2)
			if mouse.Hovered(ctx2.Rect) {
				alpha = 0.4
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
