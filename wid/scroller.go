package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
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
	NormalAlpha       = float32(0.5)
	HoverAlpha        = float32(0.8)
	ScrollerMargin    = float32(1.0)
	ThumbCornerRadius = float32(2.0)
)

// DrawScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is posistion. Ymax is max Ypos. Yvis is the visible part
func DrawScrollbar(r f32.Rect, Ymax float32, Yvis float32, state *ScrollState) {
	state.Ypos = max(0, min(state.Ypos, Ymax))
	// Draw scrollbar track
	r = f32.Rect{r.X + r.W - ScrollbarWidth, r.Y + ScrollerMargin, ScrollbarWidth, r.H - 2*ScrollerMargin}
	alpha := f32.Sel(mouse.Hovered(r), NormalAlpha, HoverAlpha)
	gpu.RoundedRect(r, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Bg().Alpha(alpha), f32.Transparent)

	hThumb := min(r.H, max(MinThumbHeight, Yvis*r.H/Ymax))
	thumbPos := state.Ypos * (r.H - hThumb) / Ymax
	// Draw thumb
	rt := f32.Rect{r.X, r.Y + thumbPos, ScrollbarWidth - ScrollerMargin*2, hThumb}
	gpu.RoundedRect(rt, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().Alpha(alpha), f32.Transparent)

	if mouse.LeftBtnPressed(rt) && !state.dragging {
		state.dragging = true
		state.StartPos = mouse.StartDrag()
	}

	if state.dragging {
		state.Ypos = max(0, min(state.Ypos+mouse.Pos().Y-state.StartPos.Y))
		slog.Info("Scrolled", "Ypos", int(state.Ypos), "Ymax", int(Ymax), "r.H", int(r.H), "Startpos", int(state.StartPos.Y))
		state.StartPos = mouse.Pos()
		state.dragging = mouse.LeftBtnDown()
	}
	if scr := sys.ScrolledY(); scr != 0 {
		state.Ypos = max(0, state.Ypos-scr*20)
		slog.Info("Scrolled", "Ypos", int(state.Ypos), "Ymax", int(Ymax), "r.H", int(r.H))
		gpu.Invalidate(0)
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
			DrawScrollbar(ctx.Rect, sumH, ctx.Rect.H, state)
		}

		return Dim{0, 0, 0}
	}
}
