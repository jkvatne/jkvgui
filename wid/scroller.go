package wid

import (
	"log/slog"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
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
	// Nlast is the number of item actually drawn/calculated
	// Nlast will be Nmax when all items are drawn/calculated
	Nlast int
	// Ylast is the height of all items we have seen
	// Ylast will be equal to Ymax when all items are drawn/calculated
	Ylast        float32
	CacheStart   int
	Cache        []Wid
	CacheMaxSize int
	DbTotalCount int
	dbCount      func() int
	dbRead       func(n int) Wid
}

type ScrollStyle struct {
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
	ScrollbarWidth:    10.0,
	MinThumbHeight:    15.0,
	TrackAlpha:        0.15,
	NormalAlpha:       0.4,
	HoverAlpha:        0.8,
	ScrollerMargin:    1.0,
	ThumbCornerRadius: 5.0,
	ScrollFactor:      0.2,
}

var doDbDebug = false
var doScrollDebug = true

func dbDebug(msg string, args ...any) {
	if doDbDebug {
		slog.Info(msg, args...)
	}
}

func scrollDebug(msg string, args ...any) {
	if doScrollDebug {
		slog.Info(msg, args...)
	}
}

// getCachedWidget implements a cache of widget pointers
func getCachedWidget(s *ScrollState, idx int) Wid {
	s.DbTotalCount = s.dbCount()
	if s.CacheMaxSize == 0 {
		s.CacheMaxSize = 8
	}
	if idx >= s.DbTotalCount {
		return nil
	}
	if idx < s.CacheStart-s.CacheMaxSize {
		// We have jumped far before start. Re-fill cache starting a bit before idx, but not less than 0.
		s.CacheStart = max(0, idx-s.CacheMaxSize/5)
		s.Cache = getBatchFromDb(s, s.CacheStart, s.CacheMaxSize)
		dbDebug("Invalidate cache    ", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache), "added", s.CacheMaxSize)
	} else if idx > s.CacheStart+len(s.Cache)+s.CacheMaxSize {
		// We have jumped far after the start. Re-fill cache from a bit before idx
		s.CacheStart = min(idx-s.CacheMaxSize/5, s.DbTotalCount-s.CacheMaxSize)
		s.CacheStart = max(0, s.CacheStart)
		s.Cache = getBatchFromDb(s, s.CacheStart, s.CacheMaxSize)
		dbDebug("Invalidate cache    ", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache), "added", s.CacheMaxSize)
	} else if idx >= s.CacheStart+len(s.Cache) {
		// Moving past the end of the cache. Read inn more from database, ca 25% of the capacity
		// Repeat if needed
		for idx >= s.CacheStart+len(s.Cache) {
			w := getBatchFromDb(s, s.CacheStart+len(s.Cache), s.CacheMaxSize/4)
			s.Cache = append(s.Cache, w...)
			// If adding data made the cache too large, throw out the beginning
			overflowCount := len(s.Cache) - s.CacheMaxSize
			if overflowCount > 0 {
				s.Cache = s.Cache[overflowCount:]
				s.CacheStart = s.CacheStart + len(w)
				dbDebug("Adding to cache   ", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache), "added", len(w))
			} else {
				dbDebug("Reading beyond end", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache), "Read", len(w))
			}
		}
	} else if idx < s.CacheStart && (s.CacheStart-idx) < s.CacheMaxSize {
		// Read in either a full batch, or the number of items missing at the front.
		// This is only valid if the idx is within the CacheSize before start
		cnt := min(s.CacheMaxSize, s.CacheStart)
		// Starting at either 0 or the number
		w := getBatchFromDb(s, s.CacheStart-cnt, cnt)
		if len(w) != cnt {
			slog.Error("getBatchFromDb returned too few items")
		}
		s.CacheStart = s.CacheStart - cnt
		s.Cache = append(w, s.Cache...)
		s.Cache = s.Cache[:s.CacheMaxSize]
		dbDebug("Fill Cache front  ", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache), "cnt", cnt)

	}
	if idx-s.CacheStart < 0 {
		slog.Error("GetItem failed   ", "idx", idx, "CacheStart", s.CacheStart, "size", len(s.Cache))
		return nil
	}
	if idx-s.CacheStart >= len(s.Cache) {
		return nil
	}
	return s.Cache[idx-s.CacheStart]
}

// getBatchFromDb reads a number of items from the database
func getBatchFromDb(s *ScrollState, start int, cnt int) (w []Wid) {
	s.DbTotalCount = s.dbCount()
	if start >= s.DbTotalCount {
		return nil
	}
	for i := 0; i < cnt; i++ {
		if i+start >= s.DbTotalCount {
			break
		}
		w = append(w, s.dbRead(start+i))
	}
	return w
}

// VertScollbarUserInput will draw a bar at the right edge of the area r.
func VertScollbarUserInput(ctx Ctx, state *ScrollState, style *ScrollStyle) float32 {
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
		}
	}
	if ctx.Win.Hovered(ctx.Rect) {
		scr := ctx.Win.ScrolledY()
		if scr != 0 {
			ctx.Win.ScrolledDistY = 0
			// Handle mouse scroll-wheel. Scrolling down gives negative scr value
			// ScrollFactor is the fraction of the visible area that is scrolled.
			dy = -(scr * ctx.Rect.H) * style.ScrollFactor
			ctx.Win.Invalidate()
			scrollDebug("ScrollWheelInput:", "dy", int(dy))
		} else if sys.GetCurrentWindow().LastKey == sys.KeyHome {
			scrollDebug("Scroll KeyHome")
			dy = -999999
			ctx.Win.Invalidate()
		} else if sys.GetCurrentWindow().LastKey == sys.KeyEnd {
			scrollDebug("Scroll KeyEnd")
			dy = 999999
			ctx.Win.Invalidate()
		} else if sys.GetCurrentWindow().LastKey == sys.KeyDown {
			scrollDebug("Scroll KeyDown")
			dy = ctx.H / 5
		} else if sys.GetCurrentWindow().LastKey == sys.KeyUp {
			scrollDebug("Scroll KeyUp")
			dy = -ctx.H / 5
		} else if sys.GetCurrentWindow().LastKey == sys.KeyPageDown {
			scrollDebug("Scroll KeyDown")
			dy = ctx.H
		} else if sys.GetCurrentWindow().LastKey == sys.KeyPageUp {
			scrollDebug("Scroll KeyUp")
			dy = -ctx.H
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
func DrawVertScrollbar(ctx Ctx, state *ScrollState, style *ScrollStyle) {
	if ctx.Rect.H > state.Ymax {
		return
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
			scrollDebug("- At bottom of list   ", "yScroll", f32.F2(yScroll),
				"Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy), "Nmax", state.Nmax,
				"Npos", state.Npos, "Ymax", f32.F2(state.Ymax))
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
				scrollDebug("- Scroll down limited ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
					"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", f32.F2(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			} else {
				scrollDebug("- Scroll down partial ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
					"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", int(state.Ymax), "ItemHeight", f32.F2(currentItemHeight), "AboveEnd", int(-state.Ypos-ctx.H+state.Ymax))
			}
			yScroll = 0
		} else if state.Npos < state.Nmax {
			// Go down to the top of the next widget if there is space
			state.Npos++
			state.Ypos += currentItemHeight - state.Dy
			yScroll = max(0, yScroll-(currentItemHeight-state.Dy))
			state.Dy = 0
			updateYmax(state.Npos, state, currentItemHeight)
			scrollDebug("- Scroll down to next ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos), "Dy", f32.F2(state.Dy),
				"Npos", state.Npos, "Ymax", int(state.Ymax), "Nmax", state.Nmax, "ItemHeight", f32.F2(currentItemHeight))
		} else {
			// Should never come here.
			slog.Error("- Scroll down illegal state ", "yScroll", f32.F2(yScroll), "Ypos", f32.F2(state.Ypos),
				"Dy", f32.F2(state.Dy), "Npos", state.Npos, "Ymax", f32.F2(state.Ymax), "Nmax", state.Nmax, "ctx.H", f32.F2(ctx.H))
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

// DrawCached will draw all the visible elements.
// Drawing widget n is done by the drawWidget(n) function.
// It returns nil if no more elements are available
func DrawCached(ctx Ctx, state *ScrollState) []Dim {
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
	n := 0
	for i = state.Npos; sumH < ctx.Rect.H+100; i++ {
		w := getCachedWidget(state, i)
		// if getCachedWidget returns nil, it indicates the end of the element list.
		if w == nil {
			// We now know the exact total size Ymax
			break
		}
		n++
		// Do drawing and save the widget dimensions.
		dim := w(ctx0)
		dims = append(dims, dim)
		// Move down to next element
		ctx0.Rect.Y += dim.H
		// And reduce area
		ctx0.Rect.H -= dim.H
		// Update the total height drawn
		sumH += dim.H
		updateYmax(i, state, dim.H)
		// check that Nmax was ok, and update if not. Nmax is the number of elements so i should be less.
		if i >= state.Nmax {
			// Typically Nmax should always be correct, so this indicates an error
			state.Nmax = i + 1
			slog.Error("Nmax was too small and is increased", "Nmax", state.Nmax)
		}
	}
	if n >= state.CacheMaxSize {
		state.CacheMaxSize = n + 4
		dbDebug(">> Increase Cache to", "size", state.CacheMaxSize, "n", n)
	}
	return dims
}

func CashedScroller(state *ScrollState, style *ScrollStyle, f func(itemno int) Wid, n func() int) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	if style == nil {
		style = &DefaultScrollStyle
	}
	state.dbRead = f
	state.dbCount = n
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
		DrawCached(ctx0, state)

		ctx0.Mode = CollectHeights
		yScroll := VertScollbarUserInput(ctx, state, style)
		scrollUp(yScroll, state, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		scrollDown(ctx, yScroll, state, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		DrawVertScrollbar(ctx, state, style)
		return Dim{ctx.W, ctx.H, 0}
	}
}

func Scroller(state *ScrollState, style *ScrollStyle, widgets ...Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	if style == nil {
		style = &DefaultScrollStyle
	}
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

		yScroll := VertScollbarUserInput(ctx, state, style)
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
		DrawVertScrollbar(ctx, state, style)
		return Dim{ctx.W, ctx.H, 0}
	}
}
