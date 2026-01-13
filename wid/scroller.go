package wid

import (
	"log/slog"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
)

type ScrollState struct {
	// Npos is the item number for the first visible widget
	// which can be only partially visible
	Npos int
	// Ypos is the Y position of the first pixel drawn. I.e. the amount scrolled down.
	// Minimum is 0 at the top, and maximum is Ymax-VisibleHeight.
	Ypos float32
	// Nmax is the total number of items
	Nmax int
	// Dy is the offset from the top of the first item down to the visible window.
	// i.e. the height not visible.
	Dy float32
	// Ymax is the total height of all items
	Ymax float32
	// Dragging is a flag that is true while the mous button is down in the scrollbar
	Dragging bool
	// StartPos is the mouse position on start of dragging
	StartPos float32
	// Width is typically a fraction like 0.5, used to divide available space
	// It could also be a fixed number of device independent pixel
	Width float32
	// Height is typically a fraction like 0.5, used to divide available space
	// It could also be a fixed number of device independent pixel
	Height float32
	AtEnd  bool
	Id     int
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
func VertScollbarUserInput(ctx Ctx, state *ScrollState) float32 {
	state.Dragging = state.Dragging && ctx.Win.LeftBtnDown()
	dy := float32(0.0)
	if state.Dragging {
		// Mouse dragging scroller thumb
		mouseY := ctx.Win.MousePos().Y
		mouseDelta := mouseY - state.StartPos
		thumbHeight := min(ctx.Rect.H, max(MinThumbHeight, ctx.Rect.H*ctx.Rect.H/state.Ymax))
		dy = mouseDelta * (state.Ymax - ctx.Rect.H) / (ctx.Rect.H - thumbHeight)
		if dy != 0 {
			if mouseY > ctx.Y && mouseY < ctx.Y+ctx.H {
				state.StartPos = mouseY
				ctx.Win.Invalidate()
				slog.Debug("Drag", "mouseDelta", mouseDelta, "dy", dy, "Ypos", int(state.Ypos), "Ymax", int(state.Ymax), "rect.H", int(ctx.Rect.H), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < state.Ymax-ctx.Rect.H-0.01)
			} else {
				slog.Debug("Dragging outside ctx.Rect", "MouseY", mouseY, "dy", dy, "ctx.Y", ctx.Y, "ctx.H", ctx.Rect.H)
				dy = 0
			}
		}
	}
	if ctx.Win.Hovered(ctx.Rect) {
		scr := ctx.Win.ScrolledY()
		if scr != 0 {
			ctx.Win.ScrolledDistY = 0
			// Handle mouse scroll-wheel. Scrolling down gives negative scr value
			// ScrollFactor is the fraction of the visible area that is scrolled.
			dy = -(scr * ctx.Rect.H) * ScrollFactor
			ctx.Win.Invalidate()
			// slog.Debug("ScrollWheelInput:", "dy", int(dy))
		}
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
	thumbHeight := ctx.Rect.H
	thumbPos := float32(0)
	if state.Ymax > ctx.Rect.H {
		thumbHeight = min(barRect.H, max(MinThumbHeight, ctx.Rect.H*barRect.H/state.Ymax))
		thumbPos = state.Ypos * (barRect.H - thumbHeight) / (state.Ymax - ctx.Rect.H)
	}
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
		slog.Debug("Scrollbar: Start dragging at", "StartPos", state.StartPos, "win.dragging", ctx.Win.Dragging)
	}
}

// scrollUp with negative yScroll
func scrollUp(yScroll float32, state *ScrollState, f func(n int) float32) {
	for yScroll < 0 {
		state.AtEnd = false
		if -yScroll < state.Dy {
			// Scroll up less than the partial top line. Reduce Dy/Ypos
			state.Dy = state.Dy + yScroll
			state.Ypos = state.Ypos + yScroll
			if state.Ypos < 0 {
				state.Ypos = 0
				state.Dy = 0
			}
			slog.Debug("- Scroll up partial   ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos)
			yScroll = 0
		} else if state.Npos > 0 && state.Ypos-yScroll > 0 {
			// Scroll up to previous line at its bottom edge
			state.Npos--
			h := f(state.Npos)
			ds := state.Dy
			state.Ypos = state.Ypos - ds
			state.Dy = h
			slog.Debug("- Scroll up one item  ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos, "h", h)
			yScroll = min(0, yScroll+ds)
		} else {
			slog.Debug("- At top              ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Npos", state.Npos)
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
		if state.Ypos+ctx.H > state.Ymax {
			// Below bottom of list
			state.AtEnd = true
			state.Ypos = state.Ymax - ctx.H
			if state.Ymax < ctx.H {
				state.Ymax = 0
				state.Dy = 0
				state.Npos = 0
			} else {
				h := float32(0)
				n := state.Nmax
				// Scan backwards to fill up available space
				for n >= 0 && h < ctx.H {
					n--
					h += f(n)
				}
				state.Dy = h - ctx.H
				state.Npos = n
			}
			slog.Debug("- At bottom of list   ", "yScroll", f32.F2S(yScroll, 1),
				"Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Nmax", state.Nmax,
				"Npos", state.Npos, "Ymax", f32.F2S(state.Ymax, 0))
			yScroll = 0
		} else if yScroll+state.Dy <= currentItemHeight {
			// Scrolling down within the top item.  No need to increment Npos, but we might reach the end.
			state.Ypos = state.Ypos + yScroll
			state.Dy = state.Dy + yScroll
			aboveEnd := state.Ymax - state.Ypos - ctx.H
			if aboveEnd < 0 {
				// Limit Ypos so we do not pass the end -Must also reduce Dy by the same amount
				state.Ypos = state.Ypos + aboveEnd
				state.Dy = state.Dy + aboveEnd
				slog.Debug("- Scroll down limited ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1),
					"Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", int(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			} else {
				slog.Debug("- Scroll down partial ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1),
					"Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", int(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			}
			yScroll = 0
		} else if state.Npos < state.Nmax {
			// Go down to the top of the next widget if there is space
			state.Npos++
			state.Ypos = state.Ypos + currentItemHeight - state.Dy
			state.Dy = 0
			slog.Debug("- Scroll down to next ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1),
				"Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax)
			yScroll = max(0, yScroll-currentItemHeight)
			/*} else if state.Npos < state.Nmax-1 {
			state.AtEnd = true
			state.Npos++
			state.Ypos = state.Ypos + aboveEnd
			state.Dy = state.Dy + aboveEnd
			slog.Debug("- Scroll down to btm  ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1), "Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "AboveEnd", int(aboveEnd))
			yScroll = 0*/
		} else {
			// Should never come here.
			slog.Error("- Scroll down illegal state ", "yScroll", f32.F2S(yScroll, 1), "Ypos", f32.F2S(state.Ypos, 1),
				"Dy", f32.F2S(state.Dy, 1), "Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "ctx.H", ctx.H)
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

// DrawCached will draw all the visible elements.
// Drawing widget n is done by the drawWidget(n) function.
// It returns nil if no more elements are available
func DrawCached(ctx Ctx, state *ScrollState, drawWidget func(n int) Wid) []Dim {
	var dims []Dim
	// Clip the drawing because elements can be partially above or below the allowed rectangle.
	ctx.Win.Gd.Clip(ctx.Rect)
	defer gpu.NoClip()

	// Start drawing above the top, the amount given in state.Dy.
	ctx0 := ctx
	ctx0.Rect.Y -= state.Dy
	ctx0.Rect.H += state.Dy
	sumH := -state.Dy

	// Now draw elements, starting at Npos. We draw well beyond the end of ctx.Rect
	// just to give a better estimate of the total list size Ymax
	var i int
	for i = state.Npos; sumH < ctx.Rect.H+200000; i++ {
		w := drawWidget(i)
		// if drawWidget returns nil, it indicates the end of the element list.
		if w == nil {
			// We now know the exact total size Ymax
			break
		}
		// Do drawing and save the widget dimensions.
		dim := w(ctx0)
		dims = append(dims, dim)
		// Move down to next element
		ctx0.Rect.Y += dim.H
		// And reduce area
		ctx0.Rect.H -= dim.H
		// Update the total height drawn
		sumH += dim.H
		// Update Ymax to be at least the size drawn.
		if state.Ypos-state.Dy+sumH > state.Ymax {
			slog.Debug("Increase", "yMax", f32.F2S(state.Ymax, 1),
				"by", f32.F2S(state.Ypos+sumH-state.Ymax, 1),
				"sumH", f32.F2S(sumH, 1), "i", i,
				"dim.H", f32.F2S(dim.H, 1),
				"Dy", f32.F2S(state.Dy, 1),
				"Ypos", f32.F2S(state.Ypos, 1),
				"Nmax", state.Nmax)
			state.Ymax = state.Ypos + sumH
		}
		// check that Nmax was ok, and update if not. Nmax is the number of elements so i should be less.
		if i >= state.Nmax {
			// Typically Nmax should always be correct, so this indicates an error
			state.Nmax = i + 1
			slog.Error("Nmax was too small and is increased", "Nmax", state.Nmax)
		}
	}
	return dims
}

func CashedScroller(state *ScrollState, f func(itemno int) Wid, n func() int) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	return func(ctx Ctx) Dim {
		ctx0 := ctx
		// If we are calculating sizes, just return the fixed Width/Height.
		if ctx.Mode != RenderChildren {
			return Dim{W: state.Width, H: state.Height, Baseline: 0}
		}
		// Estimated number of element is given by function n().
		// Typically it can be found from the source (database) as the total number of elements.
		state.Nmax = n()

		// Draw elements.
		DrawCached(ctx0, state, f)

		ctx0.Mode = CollectHeights
		yScroll := VertScollbarUserInput(ctx, state)
		scrollUp(yScroll, state, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		scrollDown(ctx, yScroll, state, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
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

		ctx0.Rect.Y -= state.Dy
		sumH := -state.Dy
		ctx0.Rect.H += state.Dy
		ctx.Win.Gd.Clip(ctx.Rect)
		for i := state.Npos; i < len(widgets) && sumH < ctx.Rect.H*2 && ctx0.H > 0; i++ {
			dim := widgets[i](ctx0)
			ctx0.Rect.Y += dim.H
			ctx0.Rect.H -= dim.H
			sumH += dim.H
		}
		gpu.NoClip()

		yScroll := VertScollbarUserInput(ctx, state)
		if state.Nmax < len(widgets) {
			// If we do not have correct Ymax/Nmax, we need to calculate them.
			for i := max(0, state.Nmax-1); i < len(widgets); i++ {
				ctx0.Mode = CollectHeights
				dim := widgets[i](ctx0)
				state.Ymax += dim.H
				state.Nmax = i + 1
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
