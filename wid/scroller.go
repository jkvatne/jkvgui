package wid

import (
	"log/slog"
	"math"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type ScrollState struct {
	Xpos float32
	Ypos float32
	Ymax float32
	Dy   float32
	Npos int
	Nmax int
	// Nest is the estimated number of entries
	Nest int
	// Yest is the estimated vertical size
	Yest     float32
	dragging bool
	StartPos float32
	Width    float32
	Height   float32
	AtEnd    bool
}

var (
	ScrollbarWidth    = float32(10.0)
	MinThumbHeight    = float32(15.0)
	TrackAlpha        = float32(0.15)
	NormalAlpha       = float32(0.3)
	HoverAlpha        = float32(0.8)
	ScrollerMargin    = float32(1.0)
	ThumbCornerRadius = float32(5.0)
	// ScrollFactor is the fraction of the visible area that is scrolled.
	ScrollFactor = float32(0.25)
)

// VertScollbarUserInput will draw a bar at the right edge of the area r.
func VertScollbarUserInput(ctx Ctx, Yvis float32, state *ScrollState) float32 {
	state.dragging = state.dragging && ctx.Win.LeftBtnDown()
	dy := float32(0.0)
	if state.dragging {
		// Mouse dragging scroller thumb
		dy = (ctx.Win.MousePos().Y - state.StartPos) * state.Ymax / Yvis
		if dy != 0 {
			state.StartPos = ctx.Win.MousePos().Y
			ctx.Win.Invalidate()
			slog.Debug("Drag", "dy", dy, "Ypos", int(state.Ypos), "state.Ymax", int(state.Ymax), "Yvis", int(Yvis), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < state.Ymax-Yvis-0.01)
		}
	}
	if scr := ctx.Win.ScrolledY(); scr != 0 {
		// Handle mouse scroll-wheel. Scrolling down gives negative scr value
		// ScrollFactor is the fraction of the visible area that is scrolled.
		dy = -(scr * Yvis) * ScrollFactor
		ctx.Win.Invalidate()
	}
	if dy < 0 {
		// Scrolling up means no more at end
		state.AtEnd = false
	}
	dy = float32(math.Round(float64(dy)))
	return dy
}

// DrawVertScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is the position. (Ymax-Yvis) is max Ypos. Yvis is the visible part
func DrawVertScrollbar(ctx Ctx, barRect f32.Rect, Ymax float32, Yvis float32, state *ScrollState) {
	if Yvis > Ymax {
		return
	}
	barRect = f32.Rect{X: barRect.X + barRect.W - ScrollbarWidth, Y: barRect.Y + ScrollerMargin, W: ScrollbarWidth, H: barRect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, Yvis*barRect.H/Ymax))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (Ymax - Yvis)
	if state.AtEnd {
		thumbPos = barRect.H - thumbHeight
	}
	thumbRect := f32.Rect{X: barRect.X + ScrollerMargin, Y: barRect.Y + thumbPos, W: ScrollbarWidth - ScrollerMargin*2, H: thumbHeight}
	// Draw scrollbar track
	ctx.Win.Gd.RoundedRect(barRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(TrackAlpha), f32.Transparent)
	// Draw thumb
	alpha := f32.Sel(ctx.Win.Hovered(thumbRect) || state.dragging, NormalAlpha, HoverAlpha)
	ctx.Win.Gd.RoundedRect(thumbRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(alpha), f32.Transparent)
	// Start dragging if mouse pressed
	if ctx.Win.LeftBtnPressed(thumbRect) && !state.dragging {
		state.dragging = true
		state.StartPos = ctx.Win.StartDrag().Y
	}
}

// DrawFromPos will draw widgets from state.Npos and downwards, with offset state.Dy
// It returns the total height and dimensions of all drawn widgets
func DrawFromPos(ctx Ctx, state *ScrollState, widgets ...Wid) (dims []Dim) {
	ctx0 := ctx
	ctx0.Rect.Y -= state.Dy
	sumH := -state.Dy
	ctx0.Rect.H += state.Dy
	ctx.Win.Clip(ctx.Rect)
	for i := state.Npos; i < len(widgets) && sumH < ctx.Rect.H*2 && ctx0.H > 0; i++ {
		dim := widgets[i](ctx0)
		ctx0.Rect.Y += dim.H
		ctx0.Rect.H -= dim.H
		sumH += dim.H
		dims = append(dims, dim)
	}
	gpu.NoClip()
	return dims
}

// scrollUp with negative yScroll
func scrollUp(yScroll float32, state *ScrollState, f func(n int) float32) {
	for yScroll < 0 {
		state.AtEnd = false
		if -yScroll < state.Dy {
			// Scroll up less than the partial top line
			slog.Debug("Scroll up partial ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			state.Dy = state.Dy + yScroll
			state.Ypos = max(0, state.Ypos+yScroll)
			yScroll = 0
		} else if state.Npos > 0 && state.Ypos-yScroll > 0 {
			// Scroll up to previous line
			state.Npos--
			h := f(state.Npos)
			state.Ypos = max(0, state.Ypos-state.Dy)
			slog.Debug("Scroll up one line", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 2), "Npos", state.Npos)
			yScroll = min(0, yScroll+state.Dy)
			state.Dy = h
		} else {
			slog.Debug("At top", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Npos", state.Npos)
			state.Ypos = 0
			state.Dy = 0
			state.Npos = 0
			yScroll = 0
		}
	}
}

// scrollDown has yScroll>0
func scrollDown(yScroll float32, state *ScrollState, ctxH float32, f func(n int) float32) {
	for yScroll > 0 {
		if state.Ypos+ctxH >= state.Ymax {
			// At end
			state.AtEnd = true
			slog.Debug("At bottom of list   ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			yScroll = 0
		} else if yScroll+state.Dy < f(state.Npos) {
			// We are within the current widget.
			state.Ypos += yScroll
			state.Dy += yScroll
			slog.Debug("Scroll down partial ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			yScroll = 0
		} else if state.Npos < state.Nmax {
			// Go down to the top of the next widget
			height := f(state.Npos)
			state.Npos++
			state.Ypos += height - state.Dy + 1e-9
			state.Dy = 0.0
			slog.Debug("Scroll down to top of next line", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			yScroll = max(0, yScroll-(height-state.Dy))
		} else {
			// Should never come here.
			slog.Debug("Scroll down unknown", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			yScroll = 0
		}
	}
}

func heightFromPos(ctx Ctx, pos int, f func(n int) Wid) float32 {
	ctx.Mode = CollectHeights
	return f(pos)(ctx).H
}

func DrawCashedFromPos(ctx Ctx, state *ScrollState, f func(n int) Wid) (dims []Dim) {
	ctx0 := ctx
	ctx0.Rect.Y -= state.Dy
	sumH := -state.Dy
	ctx0.Rect.H += state.Dy
	ctx.Win.Clip(ctx.Rect)
	var i int
	for i = state.Npos; sumH < ctx.Rect.H*2 && ctx0.H > 0; i++ {
		w := f(i)
		if w == nil {
			break
		}
		dim := w(ctx0)
		ctx0.Rect.Y += dim.H
		ctx0.Rect.H -= dim.H
		sumH += dim.H
		dims = append(dims, dim)
		if i >= state.Nmax {
			state.Nmax = i + 1
			state.Ymax += dim.H
		}
	}
	// Go a bit longer than what is visible, without drawing
	// Just update Nmax/Ymax
	ctx0.Mode = CollectHeights
	for range 10 {
		i++
		w := f(i)
		if w == nil {
			break
		}
		dim := w(ctx0)
		sumH += dim.H
		dims = append(dims, dim)
		if i >= state.Nmax {
			state.Nmax = i + 1
			state.Ymax += dim.H
		}
		if state.Nmax > state.Nest {
			state.Nest = state.Nmax
		}
		if state.Nmax > 0 {
			state.Yest = state.Ymax * float32(state.Nest) / float32(state.Nmax)
		}
	}
	gpu.NoClip()
	return dims
}

func CashedScroller(state *ScrollState, f func(itemno int) Wid, n func() int) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	return func(ctx Ctx) Dim {
		ctx0 := ctx
		if ctx.Mode != RenderChildren {
			return Dim{W: state.Width, H: state.Height, Baseline: 0}
		}
		yScroll := VertScollbarUserInput(ctx, ctx.Rect.H, state)
		state.Nest = n()
		DrawCashedFromPos(ctx0, state, f)
		if state.Nmax < state.Nest && state.Nmax > 0 {
			state.Ymax = float32(state.Nest) * state.Ymax / float32(state.Nmax)
		}
		ctx0.Mode = CollectHeights
		if yScroll < 0 {
			scrollUp(yScroll, state, func(n int) float32 {
				return heightFromPos(ctx, n, f)
			})
		} else if yScroll > 0 {
			scrollDown(yScroll, state, ctx.H, func(n int) float32 {
				return heightFromPos(ctx, n, f)
			})
		}
		DrawVertScrollbar(ctx, ctx.Rect, state.Ymax, ctx.H, state)
		return Dim{ctx.W, ctx.H, 0}
	}
}

func Scroller(state *ScrollState, widgets ...Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")

	return func(ctx Ctx) Dim {
		ctx0 := ctx
		if ctx.Mode != RenderChildren {
			return Dim{W: state.Width, H: state.Height, Baseline: 0}
		}
		yScroll := VertScollbarUserInput(ctx, ctx.Rect.H, state)
		_ = DrawFromPos(ctx0, state, widgets...)

		if state.Nmax < len(widgets) {
			// If we do not have correct Ymax/Nmax, we need to calculate them.
			for i := state.Nmax - 1; i < len(widgets); i++ {
				ctx0.Mode = CollectHeights
				dim := widgets[i](ctx0)
				state.Ymax += dim.H
			}
			state.Nmax = len(widgets)
		}
		ctx0.Mode = CollectHeights
		scrollUp(yScroll, state, func(n int) float32 {
			return widgets[n](ctx0).H
		})
		scrollDown(yScroll, state, ctx.H, func(n int) float32 {
			return widgets[n](ctx0).H
		})
		DrawVertScrollbar(ctx, ctx.Rect, state.Ymax, ctx.H, state)
		return Dim{ctx.W, ctx.H, 0}
	}
}
