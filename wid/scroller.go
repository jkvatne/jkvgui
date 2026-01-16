package wid

import (
	"log/slog"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type ScrollStyle struct {
	Height            float32
	Width             float32
	ScrollbarWidth    float32
	MinThumbHeight    float32
	TrackAlpha        float32
	NormalAlpha       float32
	HoverAlpha        float32
	ScrollerMargin    float32
	ThumbCornerRadius float32
	// ScrollFactor is the fraction of the visible area that is scrolled.
	ScrollFactor float32
}

var DefaultScrollStyle = ScrollStyle{
	Width:             0.5,
	ScrollbarWidth:    10.0,
	MinThumbHeight:    15.0,
	TrackAlpha:        0.15,
	NormalAlpha:       0.4,
	HoverAlpha:        0.8,
	ScrollerMargin:    1.0,
	ThumbCornerRadius: 3.0,
	ScrollFactor:      0.2,
}

type ScrollState struct {
	// Npos is the item number for the first visible widget
	// which can be only partially visible
	Npos int
	// Nmax is the total number of items
	Nmax int
	// Nlast is the number of item actually drawn/calculated, will be Nmax when all items are drawn/calculated
	Nlast int
	// Ypos is the Y position of the first pixel drawn. I.e. the amount scrolled down,
	// minimum is 0 at the top, and maximum is Ymax-VisibleHeight.
	Ypos float32
	// Ylast is the height of all items we have seen, will be equal to Ymax when all items are drawn/calculated
	Ylast float32
	// Ymax is the total height of all items
	Ymax float32
	// Dy is the offset from the top of the first item down to the visible window.
	// i.e. the height not visible.
	Dy float32
	// Dragging is a flag that is true while the mous button is down in the scrollbar
	Dragging bool
	// StartPos is the mouse position on start of dragging
	StartPos      float32
	AtEnd         bool
	Id            int
	PendingScroll float32
}

func doScrolling(ctx Ctx, state *ScrollState, f func(n int) float32) {
	ds := f32.Abs(state.PendingScroll)
	if f32.Abs(state.PendingScroll) < 2*ctx.H {
		ds = min(ds, ctx.H/12)
	}
	if state.PendingScroll > 0 {
		state.PendingScroll -= ds
		ctx.Mode = CollectHeights
		scrollDown(ctx, ds, state, f)
		sys.Invalidate()
	} else if state.PendingScroll < 0 {
		state.PendingScroll += ds
		ctx.Mode = CollectHeights
		scrollUp(-ds, state, f)
		sys.Invalidate()
	}
}

// VertScollbarUserInput will draw a bar at the right edge of the area r.
func VertScollbarUserInput(ctx Ctx, state *ScrollState, style *ScrollStyle) {
	state.Dragging = state.Dragging && ctx.Win.LeftBtnDown()
	dy := float32(0.0)
	if state.Dragging {
		// Mouse dragging scroller thumb
		mouseY := ctx.Win.MousePos().Y
		mouseDelta := mouseY - state.StartPos
		thumbHeight := min(ctx.Rect.H, max(style.MinThumbHeight, ctx.Rect.H*ctx.Rect.H/state.Ymax))
		dy = mouseDelta * (state.Ymax - ctx.Rect.H) / (ctx.Rect.H - thumbHeight)
		if dy != 0 {
			if mouseY > ctx.Y && mouseY < ctx.Y+ctx.H {
				state.StartPos = mouseY
				ctx.Win.Invalidate()
				scrollDebug("Drag", "mouseDelta", mouseDelta, "dy", dy, "Ypos", int(state.Ypos), "Ymax", int(state.Ymax), "rect.H", int(ctx.Rect.H), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < state.Ymax-ctx.Rect.H-0.01)
			} else {
				scrollDebug("Dragging outside ctx.Rect", "MouseY", mouseY, "dy", dy, "ctx.Y", ctx.Y, "ctx.H", ctx.Rect.H)
				dy = 0
			}
			state.PendingScroll += dy
		}
	} else if ctx.Win.Hovered(ctx.Rect) {
		scr := ctx.Win.ScrolledY()
		w := sys.GetCurrentWindow()
		if w == nil {
			slog.Error("Current window is nil")
		} else if scr != 0 {
			ctx.Win.ScrolledDistY = 0
			// Handle mouse scroll-wheel. Scrolling down gives negative scr value
			// ScrollFactor is the fraction of the visible area that is scrolled.
			dy = -(scr * ctx.Rect.H) * style.ScrollFactor
			ctx.Win.Invalidate()
			if dy < 0 {
				state.AtEnd = false
			}
			state.PendingScroll += dy
			scrollDebug("ScrollWheelInput:", "dy", int(dy))
		} else if w.LastKey == sys.KeyHome {
			scrollDebug("Scroll KeyHome")
			state.AtEnd = false
			ctx.Win.Invalidate()
			state.PendingScroll = -999999
		} else if w.LastKey == sys.KeyEnd {
			scrollDebug("Scroll KeyEnd")
			ctx.Win.Invalidate()
			state.PendingScroll = 999999
		} else if w.LastKey == sys.KeyDown {
			scrollDebug("Scroll KeyDown")
			state.PendingScroll = ctx.H / 5
		} else if w.LastKey == sys.KeyUp {
			scrollDebug("Scroll KeyUp")
			state.AtEnd = false
			state.PendingScroll -= ctx.H / 5
		} else if w.LastKey == sys.KeyPageDown {
			scrollDebug("Scroll KeyDown")
			state.PendingScroll += ctx.H
		} else if w.LastKey == sys.KeyPageUp {
			scrollDebug("Scroll KeyUp")
			state.AtEnd = false
			state.PendingScroll -= ctx.H
		}
	}
}

// DrawVertScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is the position. (Ymax-Yvis) is max Ypos. Yvis is the visible part
func DrawVertScrollbar(ctx Ctx, state *ScrollState, style *ScrollStyle) {
	if ctx.Rect.H > state.Ymax {
		return
	}
	if style == nil {
		style = &DefaultScrollStyle
	}
	barRect := f32.Rect{
		X: ctx.Rect.X + ctx.Rect.W - style.ScrollbarWidth,
		Y: ctx.Rect.Y + style.ScrollerMargin,
		W: style.ScrollbarWidth,
		H: ctx.Rect.H - 2*style.ScrollerMargin}
	thumbHeight := ctx.Rect.H
	thumbPos := float32(0)
	if state.Ymax > ctx.Rect.H {
		thumbHeight = min(barRect.H, max(style.MinThumbHeight, ctx.Rect.H*barRect.H/state.Ymax))
		thumbPos = state.Ypos * (barRect.H - thumbHeight) / (state.Ymax - ctx.Rect.H)
	}
	if state.AtEnd {
		thumbPos = barRect.H - thumbHeight
	}
	thumbRect := f32.Rect{X: barRect.X + style.ScrollerMargin, Y: barRect.Y + thumbPos, W: style.ScrollbarWidth - style.ScrollerMargin*2, H: thumbHeight}
	// Draw scrollbar track
	ctx.Win.Gd.RoundedRect(barRect, style.ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(style.TrackAlpha), f32.Transparent)
	// Draw thumb
	alpha := f32.Sel(ctx.Win.Hovered(thumbRect) || state.Dragging, style.NormalAlpha, style.HoverAlpha)
	ctx.Win.Gd.RoundedRect(thumbRect, style.ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(alpha), f32.Transparent)
	// Start dragging if mouse pressed
	if ctx.Win.LeftBtnPressed(thumbRect) && !state.Dragging {
		state.Dragging = true
		state.StartPos = ctx.Win.StartDrag().Y
		scrollDebug("Scrollbar: Start dragging at", "StartPos", state.StartPos, "win.dragging", ctx.Win.Dragging)
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
			scrollDebug("- Scroll up partial   ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy), "Npos", state.Npos)
			yScroll = 0
		} else if state.Npos > 0 && state.Ypos-yScroll > 0 {
			// Scroll up to previous line at its bottom edge
			state.Npos--
			h := f(state.Npos)
			ds := state.Dy
			state.Ypos = state.Ypos - state.Dy
			state.Dy = h
			scrollDebug("- Scroll up one item  ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy), "Npos", state.Npos, "h", h)
			yScroll = min(0, yScroll+ds)
		} else {
			scrollDebug("- At top              ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos), "Npos", state.Npos)
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
			state.Ypos = max(0, state.Ymax-ctx.H)
			if state.Ymax < ctx.H {
				state.Dy = 0
				state.Npos = 0
				state.Ypos = 0
				scrollDebug("- Too few elements    ", "yScroll", f32.F2(yScroll),
					"Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy), "Nmax", state.Nmax,
					"Npos", state.Npos, "Ymax", f32.F2(state.Ymax))
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
				scrollDebug("- At bottom of list   ", "yScroll", f32.F2(yScroll),
					"Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy), "Nmax", state.Nmax,
					"Npos", state.Npos, "Ymax", f32.F2(state.Ymax))
			}
			yScroll = 0
		} else if yScroll+state.Dy <= currentItemHeight {
			// Scrolling down within the top item.  No need to increment Npos, but we might reach the end.
			state.Ypos = state.Ypos + yScroll
			state.Dy = state.Dy + yScroll
			aboveEnd := state.Ymax - state.Ypos - ctx.H
			if state.Ymax-state.Ypos < ctx.H {
				// Limit Ypos so we do not pass the end -Must also reduce Dy by the same amount
				state.Ypos = state.Ymax - ctx.H
				state.Dy = state.Dy + aboveEnd
				state.PendingScroll = 0
				scrollDebug("- Scroll down limited ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
					"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", f32.F2(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			} else {
				scrollDebug("- Scroll down partial ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
					"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", f32.F2(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			}
			yScroll = 0
		} else if state.Npos < state.Nmax-1 {
			// Go down to the top of the next widget if there is space
			state.Npos++
			state.Ypos += currentItemHeight - state.Dy
			yScroll = max(0, yScroll-(currentItemHeight-state.Dy))
			state.Dy = 0
			updateYmax(state.Npos, state, currentItemHeight)
			scrollDebug("- Scroll down to next ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy),
				"Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "ItemHeight", f32.F2(currentItemHeight), "ctx.H", ctx.H)
		} else {
			// Should never come here.
			slog.Error("- Scroll down illegal state ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
				"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", f32.F2(state.Ymax), "Nmax", state.Nmax, "ctx.H", f32.F2(ctx.H))
			yScroll = 0
		}
	}
}

func updateYmax(i int, state *ScrollState, h float32) {
	// Update Ymax to be at least the size drawn.
	if i+1 > state.Nlast {
		state.Nlast = i + 1
		if i == 0 {
			state.Ylast = h
		} else {
			state.Ylast = state.Ylast + h
		}
		// Ymax is estimated. Will only be correct when Nlast=Nmax-1
		yest := state.Ylast / float32(state.Nlast) * float32(state.Nmax)
		scrollDebug("Estimate size",
			"yest", f32.F2(yest),
			"Nlast", state.Nlast,
			"Nmax", state.Nmax,
			"Ylast", f32.F2(state.Ylast),
			"dim.H", h,
			"Dy", f32.F2(state.Dy))
		state.Ymax = yest
	}
}

// Scroller is a scrollable container (vertical scrolling only)
func Scroller(state *ScrollState, style *ScrollStyle, widgets ...Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	if style == nil {
		style = &DefaultScrollStyle
	}
	return func(ctx Ctx) Dim {
		ctx0 := ctx
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: style.Height, Baseline: 0}
		}

		// focused := ctx.Win.At(value)

		ctx0.Rect.Y -= state.Dy
		sumH := -state.Dy
		ctx0.Rect.H += state.Dy
		ctx.Win.Gd.Clip(ctx.Rect)
		for i := state.Npos; i < len(widgets) && sumH < ctx.Rect.H*2 && ctx0.H > 0; i++ {
			dim := widgets[i](ctx0)
			ctx0.Rect.Y += dim.H
			ctx0.Rect.H -= dim.H
			sumH += dim.H
			updateYmax(i, state, dim.H)
		}
		gpu.NoClip()

		VertScollbarUserInput(ctx, state, style)
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
		doScrolling(ctx, state, func(n int) float32 {
			return widgets[n](ctx0).H
		})
		DrawVertScrollbar(ctx, state, style)
		return Dim{ctx.W, ctx.H, 0}
	}
}
