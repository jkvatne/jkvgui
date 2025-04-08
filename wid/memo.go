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
	ScrollState
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
			dy := (mouse.Pos().Y - state.StartPos.Y) / ctx.H * float32(TotalLineCount) * lineHeight
			if dy < 0 {
				state.NotAtEnd = true
			}
			state.Ypos += dy
			state.StartPos = mouse.Pos()
			state.dragging = mouse.LeftBtnDown()
		}
		scr := sys.ScrolledY()
		if scr != 0 {
			if scr > 0 {
				state.NotAtEnd = true
			}
			state.Ypos -= scr * lineHeight
			gpu.Invalidate(0)
		}
		state.Ypos = min(state.Ypos, lineHeight*float32(len(*text))-ctx.Rect.H)
		state.Ypos = max(state.Ypos, 0)

		if TotalLineCount > MemoLineCount {
			sumH := float32(len(*text)) * lineHeight
			DrawScrollbar(ctx.Rect, sumH, &state.ScrollState)
		}
		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
