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
	Max      float32
	dragging bool
	StartPos f32.Pos
}

func DrawScrollbar(r f32.Rect, sumH float32, state *ScrollState) {
	// Draw scrollbar track
	r.X += r.W - 8
	r.W = 8

	alpha := float32(0.4)
	if mouse.Hovered(r) {
		alpha = 1.0
	}
	gpu.RoundedRect(r, 2, 0.0, theme.SurfaceContainer.Bg().Alpha(alpha), f32.Transparent)
	// Draw thumb
	r.X += 1.0
	r.W -= 2.0
	r.Y += state.Ypos * r.H / sumH
	r.H *= r.H / sumH
	if mouse.LeftBtnPressed(r) && !state.dragging {
		state.dragging = true
		state.StartPos = mouse.StartDrag()
	}
	gpu.RoundedRect(r, 2, 0.0, theme.SurfaceContainer.Fg().Alpha(alpha), f32.Transparent)

}

func ScrollPane(state *ScrollState, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		sumH := float32(0.0)
		ctx0 := ctx
		ctx0.Rect.H = 0
		ne := 0
		maxW := float32(0)
		ctx0.Mode = CollectHeights
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxW = max(maxW, dims[i].W)
			sumH += dims[i].H
			if dims[i].W == 0 {
				ne++
			}
		}
		// Return height
		if ctx.Mode != RenderChildren {
			return Dim{W: maxW, H: sumH, Baseline: 0}
		}
		if ne > 0 {
			remaining := ctx.Rect.H - sumH
			for i, d := range dims {
				if d.H == 0 {
					dims[i].H = remaining / float32(ne)
				}
			}
		}
		// Draw children
		ctx1 := ctx
		oversize := max(0, sumH-ctx.Rect.H)
		state.Ypos = max(0, min(state.Ypos, oversize))
		ctx1.Rect.Y = -state.Ypos
		for i, w := range widgets {
			ctx1.Rect.H = dims[i].H
			w(ctx1)
			ctx1.Rect.Y += dims[i].H
		}
		if sumH >= ctx.Rect.H {
			DrawScrollbar(ctx.Rect, sumH, state)
		}
		if state.dragging {
			state.Ypos += mouse.Pos().Y - state.StartPos.Y
			state.StartPos = mouse.Pos()
			state.dragging = mouse.LeftBtnDown()
		}
		if sys.ScrolledY != 0 {
			state.Ypos -= sys.ScrolledY * 20
			sys.ScrolledY = 0
			gpu.Invalidate(0)
		}
		return Dim{0, 0, 0}
	}
}
