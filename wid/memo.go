package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/theme"
)

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

var MemoStateMap = make(map[any]*ScrollState)

func drawlines(ctx Ctx, text string, Wmax float32, f *font.Font, fg f32.Color) (sumH float32) {
	wrapedLines := font.Split(text, Wmax, f)
	lineHeight := f.Height
	for _, line := range wrapedLines {
		if fg != f32.Transparent {
			f.DrawText(ctx.Win.Gd, ctx.X, ctx.Y+f.Baseline, fg, ctx.Rect.W, gpu.LTR, line)
		}
		ctx.Rect.Y += lineHeight
		sumH += lineHeight
	}
	return sumH
}

func Memo(text *[]string, style *MemoStyle) Wid {
	if style == nil {
		style = DefMemo
	}
	StateMapMutex.RLock()
	state := MemoStateMap[text]
	StateMapMutex.RUnlock()
	if state == nil {
		StateMapMutex.Lock()
		MemoStateMap[text] = &ScrollState{}
		state = MemoStateMap[text]
		StateMapMutex.Unlock()
		// We want to show the last lines by default.
		state.AtEnd = true
	}

	f := font.Fonts[style.FontNo]
	fg := style.Color.Fg()

	return func(ctx Ctx) Dim {
		baseline := f.Baseline
		if ctx.Mode != RenderChildren {
			return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
		}
		ctx.Rect = ctx.Rect.Inset(style.OutsidePadding, style.BorderWidth)
		ctx.Win.Gd.RoundedRect(ctx.Rect, style.CornerRadius, style.BorderWidth, f32.Transparent, style.BorderRole.Fg())

		ctx.Rect = ctx.Rect.Inset(style.InsidePadding, 0)
		if *DebugWidgets {
			ctx.Win.Gd.RoundedRect(ctx.Rect, 0.0, 1.0, f32.Transparent, f32.Red)
		}
		heights := make([]float32, 64)
		Wmax := float32(0)
		if style.Wrap {
			Wmax = ctx.Rect.W
		}
		yScroll := VertScollbarUserInput(ctx, state)
		ctx.Win.Clip(ctx.Rect)
		ctx.Win.Mutex.Lock()
		defer ctx.Win.Mutex.Unlock()
		textLen := len(*text)
		ctx0 := ctx

		if state.Nmax < textLen {
			// If we do not have Ymax/Nmax, we need to calculate them.
			for i := state.Nmax; i < textLen; i++ {
				state.Ymax += drawlines(ctx0, (*text)[i], Wmax, f, f32.Transparent)
			}
			state.Nmax = textLen
		}
		if state.Ypos >= state.Ymax+ctx.H {
			state.AtEnd = true
		}
		if state.AtEnd && state.Ymax > ctx.H {
			// Start from the bottom
			ctx0.Y = ctx.Y + ctx.H
			sumH := float32(0.0)
			for i := textLen - 1; i >= 0 && ctx0.Y > ctx.Y; i-- {
				h := drawlines(ctx0, (*text)[i], Wmax, f, f32.Transparent)
				ctx0.Y -= h
				sumH += h
				_ = drawlines(ctx0, (*text)[i], Wmax, f, fg)
				state.Npos = i
				if i > state.Nmax {
					state.Ymax += h
					state.Nmax++
				}

			}
			state.Dy = ctx.Y - ctx0.Y
			state.Ypos = state.Ymax - ctx.H
		} else {
			// Start from Npos
			ctx0.Rect.Y -= state.Dy
			for i := state.Npos; i < textLen && ctx0.Y-ctx.Y < ctx.Rect.H; i++ {
				h := drawlines(ctx0, (*text)[i], Wmax, f, fg)
				heights[i-state.Npos] = h
				ctx0.Rect.Y += h
				if i >= state.Nmax {
					state.Ymax += h
					state.Nmax = i + 1
				}
			}
		}
		gpu.NoClip()

		if yScroll < 0 {
			scrollUp(yScroll, state,
				func(n int) float32 {
					return drawlines(ctx0, (*text)[n], Wmax, f, f32.Transparent)
				})
		} else if yScroll > 0 {
			scrollDown(ctx, yScroll, state,
				func(n int) float32 {
					return drawlines(ctx0, (*text)[n], Wmax, f, f32.Transparent)
				})
		}
		DrawVertScrollbar(ctx, state)
		return Dim{W: ctx.W, H: ctx.H, Baseline: baseline}
	}
}
