package wid

import (
	"flag"
	"log/slog"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

// CachedScrollState is a ScrollState with additional data
// to implement a cache of widgets.
type CachedScrollState struct {
	ScrollState
	cacheStart   int
	cache        []Wid
	cacheMaxSize int
	dbTotalCount int
	dbCount      func() int
	dbRead       func(n int) Wid
}

var doDbDebug = flag.Bool("debug-db", false, "Set to print db logs")
var doScrollDebug = flag.Bool("debug-scroll", true, "Set to print scrolling logs")

func dbDebug(msg string, args ...any) {
	if *doDbDebug {
		slog.Info(msg, args...)
	}
}

func scrollDebug(msg string, args ...any) {
	if *doScrollDebug {
		slog.Info(msg, args...)
	}
}

// getCachedWidget implements a cache of widget pointers
func getCachedWidget(s *CachedScrollState, idx int) Wid {
	s.dbTotalCount = s.dbCount()
	if s.cacheMaxSize == 0 {
		s.cacheMaxSize = 8
	}
	if idx >= s.dbTotalCount {
		return nil
	}
	if idx < s.cacheStart-s.cacheMaxSize {
		// We have jumped far before start. Re-fill cache starting a bit before idx, but not less than 0.
		s.cacheStart = max(0, idx-s.cacheMaxSize/5)
		s.cache = getBatchFromDb(s, s.cacheStart, s.cacheMaxSize)
		dbDebug("Invalidate cache    ", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache), "added", s.cacheMaxSize)
	} else if idx > s.cacheStart+len(s.cache)+s.cacheMaxSize {
		// We have jumped far after the start. Re-fill cache from a bit before idx
		s.cacheStart = min(idx-s.cacheMaxSize/5, s.dbTotalCount-s.cacheMaxSize)
		s.cacheStart = max(0, s.cacheStart)
		s.cache = getBatchFromDb(s, s.cacheStart, s.cacheMaxSize)
		dbDebug("Invalidate cache    ", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache), "added", s.cacheMaxSize)
	} else if idx >= s.cacheStart+len(s.cache) {
		// Moving past the end of the cache. Read inn more from database, ca 25% of the capacity
		// Repeat if needed
		for idx >= s.cacheStart+len(s.cache) {
			w := getBatchFromDb(s, s.cacheStart+len(s.cache), s.cacheMaxSize/4)
			s.cache = append(s.cache, w...)
			// If adding data made the cache too large, throw out the beginning
			overflowCount := len(s.cache) - s.cacheMaxSize
			if overflowCount > 0 {
				s.cache = s.cache[overflowCount:]
				s.cacheStart = s.cacheStart + len(w)
				dbDebug("Adding to cache   ", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache), "added", len(w))
			} else {
				dbDebug("Reading beyond end", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache), "Read", len(w))
			}
		}
	} else if idx < s.cacheStart && (s.cacheStart-idx) < s.cacheMaxSize {
		// Read in either a full batch, or the number of items missing at the front.
		// This is only valid if the idx is within the CacheSize before start
		cnt := min(s.cacheMaxSize, s.cacheStart)
		// Starting at either 0 or the number
		w := getBatchFromDb(s, s.cacheStart-cnt, cnt)
		if len(w) != cnt {
			slog.Error("getBatchFromDb returned too few items")
		}
		s.cacheStart = s.cacheStart - cnt
		s.cache = append(w, s.cache...)
		s.cache = s.cache[:s.cacheMaxSize]
		dbDebug("Fill cache front  ", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache), "cnt", cnt)

	}
	if idx-s.cacheStart < 0 {
		slog.Error("GetItem failed   ", "idx", idx, "cacheStart", s.cacheStart, "size", len(s.cache))
		return nil
	}
	if idx-s.cacheStart >= len(s.cache) {
		return nil
	}
	return s.cache[idx-s.cacheStart]
}

// getBatchFromDb reads a number of items from the database
func getBatchFromDb(s *CachedScrollState, start int, cnt int) (w []Wid) {
	s.dbTotalCount = s.dbCount()
	if start >= s.dbTotalCount {
		return nil
	}
	for i := 0; i < cnt; i++ {
		if i+start >= s.dbTotalCount {
			break
		}
		w = append(w, s.dbRead(start+i))
	}
	return w
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

// drawCached will draw all the visible elements.
// Drawing widget n is done by the drawWidget(n) function.
// It returns nil if no more elements are available
func drawCached(ctx Ctx, state *CachedScrollState) []Dim {
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
		updateYmax(i, &state.ScrollState, dim.H)
		// check that Nmax was ok, and update if not. Nmax is the number of elements so i should be less.
		if i >= state.Nmax {
			// Typically Nmax should always be correct, so this indicates an error
			state.Nmax = i + 1
			slog.Error("Nmax was too small and is increased", "Nmax", state.Nmax)
		}
	}
	if n >= state.cacheMaxSize {
		state.cacheMaxSize = n + 4
		dbDebug(">> Increase cache to", "size", state.cacheMaxSize, "n", n)
	}
	return dims
}

// CaschedScroller is a scrollable container with vertical scrolling,
// it implements a cache for the elements in the container, suitable
// for large database tables etc.
func CashedScroller(state *CachedScrollState, style *ScrollStyle, f func(itemno int) Wid, n func() int) Wid {
	if style == nil {
		style = &DefaultScrollStyle
	}
	if state == nil {
		f32.Exit(1, "Scroller state must not be nil")
		return nil
	}
	state.dbRead = f
	state.dbCount = n
	return func(ctx Ctx) Dim {
		// If we are calculating sizes, just return the fixed Width/Height.
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: style.Height, Baseline: 0}
		}
		// Estimated number of element is given by function n().
		// Typically it can be found from the source (database) as the total number of elements.
		state.Nmax = n()

		// Draw elements.
		drawCached(ctx, state)
		VertScollbarUserInput(ctx, &state.ScrollState, style)
		doScrolling(ctx, &state.ScrollState, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		DrawVertScrollbar(ctx, &state.ScrollState, style)
		return Dim{ctx.W, ctx.H, 0}
	}
}
