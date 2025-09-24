package wid

import (
	"log/slog"
	"math"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type ScrollState struct {
	// Npos is the item number for the first widge
	Npos int
	// Ypos is the Y offset for the first widge (if it is partially visible)
	Ypos float32
	// Nmax is the total number of items
	Nmax int
	// Ymax is the total height of all items
	Ymax float32
	// Dy is the scroll distance
	Dy float32
	// Nest is the estimated number of entries
	Nest int
	// Yest is the estimated vertical size
	Yest float32
	// Dragging is a flag that is true while the mous button is down in the scrollbar
	Dragging bool
	// StartPos is the mouse position on start of dragging
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

func round(x float32) float32 {
	return float32(math.Round(float64(x)))
}

// VertScollbarUserInput will draw a bar at the right edge of the area r.
func VertScollbarUserInput(ctx Ctx, state *ScrollState) float32 {
	state.Dragging = state.Dragging && ctx.Win.LeftBtnDown()
	dy := float32(0.0)
	if state.Dragging {
		// Mouse dragging scroller thumb
		mouseY := ctx.Win.MousePos().Y
		mouseDelta := mouseY - state.StartPos
		thumbHeight := min(ctx.Rect.H, max(MinThumbHeight, ctx.Rect.H*ctx.Rect.H/state.Ymax))
		dy = mouseDelta * (state.Yest - ctx.Rect.H) / (ctx.Rect.H - thumbHeight)
		if dy != 0 && mouseY > ctx.Y && mouseY < ctx.Y+ctx.H {
			state.StartPos = mouseY
			ctx.Win.Invalidate()
			slog.Debug("Drag", "mouseDelta", mouseDelta, "dy", dy, "Ypos", int(state.Ypos), "Yest", int(state.Yest), "rect.H", int(ctx.Rect.H), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < state.Ymax-ctx.Rect.H-0.01)
		}
	}
	if scr := ctx.Win.ScrolledY(); scr != 0 {
		// Handle mouse scroll-wheel. Scrolling down gives negative scr value
		// ScrollFactor is the fraction of the visible area that is scrolled.
		dy = -(scr * ctx.Rect.H) * ScrollFactor
		ctx.Win.Invalidate()
	}
	if dy < 0 {
		// Scrolling up means no more at end
		state.AtEnd = false
	}
	return dy
}

// DrawVertScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is the position. (Ymax-Yvis) is max Ypos. Yvis is the visible part
func DrawVertScrollbar(ctx Ctx, state *ScrollState) {
	if ctx.Rect.H > state.Ymax {
		return
	}
	barRect := f32.Rect{
		X: ctx.Rect.X + ctx.Rect.W - ScrollbarWidth,
		Y: ctx.Rect.Y + ScrollerMargin,
		W: ScrollbarWidth,
		H: ctx.Rect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, ctx.Rect.H*barRect.H/state.Ymax))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (state.Ymax - ctx.Rect.H)
	if state.AtEnd {
		thumbPos = barRect.H - thumbHeight
	}
	thumbRect := f32.Rect{X: barRect.X + ScrollerMargin, Y: barRect.Y + thumbPos, W: ScrollbarWidth - ScrollerMargin*2, H: thumbHeight}
	// Draw scrollbar track
	ctx.Win.Gd.RoundedRect(barRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(TrackAlpha), f32.Transparent)
	// Draw thumb
	alpha := f32.Sel(ctx.Win.Hovered(thumbRect) || state.Dragging, NormalAlpha, HoverAlpha)
	ctx.Win.Gd.RoundedRect(thumbRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(alpha), f32.Transparent)
	// Start dragging if mouse pressed
	if ctx.Win.LeftBtnPressed(thumbRect) && !state.Dragging {
		state.Dragging = true
		state.StartPos = ctx.Win.StartDrag().Y
		slog.Info("Start dragging at", "StartPos", state.StartPos)
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
			slog.Info("Scroll up partial ", "yScroll", f32.F2S(yScroll, 1, 4), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 4), "Npos", state.Npos)
			state.Dy = state.Dy + yScroll
			state.Ypos = max(0, state.Ypos+yScroll)
			yScroll = 0
		} else if state.Npos > 0 && state.Ypos-yScroll > 0 {
			// Scroll up to previous line
			state.Npos--
			h := f(state.Npos)
			state.Ypos = max(0, state.Ypos-state.Dy)
			slog.Info("Scroll up one line", "yScroll", f32.F2S(yScroll, 1, 4), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 4), "Npos", state.Npos)
			yScroll = min(0, yScroll+state.Dy)
			state.Dy = h
		} else {
			slog.Info("At top", "yScroll", f32.F2S(yScroll, 1, 4), "Ypos", f32.F2S(state.Ypos, 1, 6), "Npos", state.Npos)
			state.Ypos = 0
			state.Dy = 0
			state.Npos = 0
			yScroll = 0
		}
	}
}

// scrollDown has yScroll>0
func scrollDown(ctx Ctx, yScroll float32, state *ScrollState, f func(n int) float32) {
	for yScroll > 0 {
		currentItemHeight := f(state.Npos)
		if state.Ypos+ctx.H >= state.Ymax {
			// At end
			state.AtEnd = true
			aboveEnd := state.Ymax - state.Ypos - ctx.H
			slog.Info("At bottom of list   ", "yScroll", f32.F2S(yScroll, 1, 5), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 5), "Npos", state.Npos, "Ymax", f32.F2S(state.Ymax, 0, 4), "AboveEnd", int(aboveEnd))
			yScroll = 0
			state.Yest = state.Ymax
			state.Ypos = state.Ymax - ctx.H
		} else if yScroll+state.Dy < currentItemHeight {
			// We are within the current widget. Calculate height of Npos
			state.Ypos += yScroll
			if state.Ypos > state.Ymax-ctx.H {
				state.Ypos = state.Ymax - ctx.H
			}
			state.Dy += yScroll
			slog.Info("Scroll down partial ", "yScroll", f32.F2S(yScroll, 1, 5), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 5), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", int(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			yScroll = 0
		} else if state.Npos < state.Nmax-1 {
			// Go down to the top of the next widget
			state.Npos++
			dy := currentItemHeight - state.Dy
			state.Ypos = min(state.Ypos+dy, state.Ypos+ctx.H)
			yScroll -= dy
			aboveEnd := state.Ymax - state.Ypos - ctx.H
			slog.Info("Scroll down to next ", "yScroll", f32.F2S(yScroll, 1, 5), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 5), "Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "AboveEnd", int(aboveEnd))
			state.Dy = 0.0
		} else {
			// Should never come here.
			slog.Error("Scroll down unknown ", "yScroll", f32.F2S(yScroll, 1, 5), "Ypos", f32.F2S(state.Ypos, 1, 6), "Dy", f32.F2S(state.Dy, 1, 5), "Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "ctx.H", ctx.H)
			yScroll = 0
		}
	}
}

func heightFromPos(ctx Ctx, pos int, f func(n int) Wid) float32 {
	ctx.Mode = CollectHeights
	if f == nil {
		return 0
	}
	w := f(pos)
	if w == nil {
		w = f(pos)
		return 0
	}
	return w(ctx).H
}

func DrawCachedFromPos(ctx Ctx, state *ScrollState, f func(n int) Wid) (dims []Dim) {
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
	for range 4 {
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
		yScroll := VertScollbarUserInput(ctx, state)
		state.Nest = n()
		DrawCachedFromPos(ctx0, state, f)
		if state.Nmax < state.Nest && state.Nmax > 1 {
			state.Yest = float32(state.Nest) * state.Ymax / float32(state.Nmax-1)
		}
		ctx0.Mode = CollectHeights
		if yScroll < 0 {
			scrollUp(yScroll, state, func(n int) float32 {
				return heightFromPos(ctx, n, f)
			})
		} else if yScroll > 0 {
			scrollDown(ctx, yScroll, state, func(n int) float32 {
				return heightFromPos(ctx, n, f)
			})
		}
		DrawVertScrollbar(ctx, state)
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
		yScroll := VertScollbarUserInput(ctx, state)
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
		scrollDown(ctx, yScroll, state, func(n int) float32 {
			return widgets[n](ctx0).H
		})
		DrawVertScrollbar(ctx, state)
		return Dim{ctx.W, ctx.H, 0}
	}
}
