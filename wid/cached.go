package wid

import (
	"flag"
	"log/slog"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

type CachedScrollState struct {
	ScrollState
	CacheStart   int
	Cache        []Wid
	CacheMaxSize int
	DbTotalCount int
	dbCount      func() int
	dbRead       func(n int) Wid
}

var doDbDebug = flag.Bool("debugDb", false, "Set to print db loggs")
var doScrollDebug = flag.Bool("debugScroll", false, "Set to print scrolling loggs")

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
func getBatchFromDb(s *CachedScrollState, start int, cnt int) (w []Wid) {
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
func DrawCached(ctx Ctx, state *CachedScrollState) []Dim {
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
	if n >= state.CacheMaxSize {
		state.CacheMaxSize = n + 4
		dbDebug(">> Increase Cache to", "size", state.CacheMaxSize, "n", n)
	}
	return dims
}

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
		ctx0 := ctx
		// If we are calculating sizes, just return the fixed Width/Height.
		if ctx.Mode != RenderChildren {
			return Dim{W: style.Width, H: style.Height, Baseline: 0}
		}
		// Estimated number of element is given by function n().
		// Typically it can be found from the source (database) as the total number of elements.
		state.Nmax = n()

		// Draw elements.
		DrawCached(ctx0, state)

		ctx0.Mode = CollectHeights
		yScroll := VertScollbarUserInput(ctx, &state.ScrollState, style)
		scrollUp(yScroll, &state.ScrollState, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		scrollDown(ctx, yScroll, &state.ScrollState, func(n int) float32 {
			return heightFromPos(ctx, n, f)
		})
		DrawVertScrollbar(ctx, &state.ScrollState, style)
		return Dim{ctx.W, ctx.H, 0}
	}
}
