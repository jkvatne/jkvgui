package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
)

type ScrollState struct {
	Xpos     float32
	Ypos     float32
	Width    float32
	dragging bool
	StartPos f32.Pos
}

var (
	ScrollbarWidth    = float32(10.0)
	MinThumbHeight    = float32(15.0)
	TrackAlpha        = float32(0.15)
	NormalAlpha       = float32(0.3)
	HoverAlpha        = float32(0.8)
	ScrollerMargin    = float32(1.0)
	ThumbCornerRadius = float32(2.0)
)

// DrawVertScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is the position. (Ymax-Yvis) is max Ypos. Yvis is the visible part
func DrawVertScrollbar(barRect f32.Rect, Ymax float32, Yvis float32, state *ScrollState) {
	if Yvis > Ymax {
		return
	}
	if state.dragging {
		// Mouse dragging scroller thumb
		dy := (mouse.Pos().Y - state.StartPos.Y) * Ymax / Yvis
		// if dy < 0 {		state.NotAtEnd = true	}
		state.Ypos += dy
		state.StartPos = mouse.Pos()
		state.dragging = mouse.LeftBtnDown()
	}
	if scr := sys.ScrolledY(); scr != 0 {
		// Handle mouse scroll-wheel
		// state.NotAtEnd = scr > 0
		state.Ypos -= scr * Yvis / 10
		gpu.Invalidate(0)
	}
	state.Ypos = max(0, min(state.Ypos, Ymax-Yvis))
	barRect = f32.Rect{barRect.X + barRect.W - ScrollbarWidth, barRect.Y + ScrollerMargin, ScrollbarWidth, barRect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, Yvis*barRect.H/(Ymax-Yvis)))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (Ymax - Yvis)
	thumbRect := f32.Rect{barRect.X, barRect.Y + thumbPos, ScrollbarWidth - ScrollerMargin*2, thumbHeight}
	// Draw scrollbar track
	gpu.RoundedRect(barRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().Alpha(TrackAlpha), f32.Transparent)
	// Draw thumb
	alpha := f32.Sel(mouse.Hovered(thumbRect) || state.dragging, NormalAlpha, HoverAlpha)
	gpu.RoundedRect(thumbRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().Alpha(alpha), f32.Transparent)
	// Start dragging if mouse pressed
	if mouse.LeftBtnPressed(thumbRect) && !state.dragging {
		state.dragging = true
		state.StartPos = mouse.StartDrag()
	}

}

func Scroller(state *ScrollState, widgets ...Wid) Wid {
	dims := make([]Dim, len(widgets))
	return func(ctx Ctx) Dim {
		ctx0 := ctx
		ctx0.Rect.H = 0
		maxW := float32(0)
		ctx0.Mode = CollectHeights
		sumH := float32(0.0)
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxW = max(maxW, dims[i].W)
			sumH += dims[i].H
		}
		// Return height
		if ctx.Mode != RenderChildren {
			return Dim{W: maxW, H: sumH, Baseline: 0}
		}
		// Draw children
		ctx0 = ctx
		ctx0.Rect.Y = -state.Ypos
		sumH = float32(0.0)
		for i, w := range widgets {
			ctx0.Rect.H = dims[i].H
			dims[i] = w(ctx0)
			ctx0.Rect.Y += dims[i].H
			sumH += dims[i].H
		}
		if sumH >= ctx.Rect.H {
			DrawVertScrollbar(ctx.Rect, sumH, ctx.Rect.H, state)
		}

		return Dim{0, 0, 0}
	}
}
