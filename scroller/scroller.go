package scroller

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type State struct {
	Xpos     float32
	Ypos     float32
	Width    float32
	Max      float32
	dragging bool
	StartPos f32.Pos
}

func W(state *State, widgets ...wid.Wid) wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		sumH := float32(0.0)
		ctx0 := ctx
		ctx0.Rect.H = 0
		ne := 0
		maxW := float32(0)
		dims := make([]wid.Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxW = max(maxW, dims[i].W)
			sumH += dims[i].H
			if dims[i].W == 0 {
				ne++
			}
		}
		// Return height
		if !ctx.Draw {
			return wid.Dim{W: maxW, H: sumH, Baseline: 0}
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
			// Draw scrollbar
			ctx2 := ctx
			ctx2.Rect.X += ctx2.Rect.W - 8
			ctx2.Rect.W = 8

			alpha := float32(0.4)
			if mouse.Hovered(ctx2.Rect) {
				alpha = 1.0
			}
			gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Bg().Alpha(alpha))
			// Draw thumb
			ctx2.Rect.X += 1.0
			ctx2.Rect.W -= 2.0
			ctx2.Rect.Y = state.Ypos * ctx.Rect.H / sumH
			ctx2.Rect.H *= ctx2.Rect.H / sumH
			if mouse.LeftBtnPressed(ctx2.Rect) && !state.dragging {
				state.dragging = true
				state.StartPos = mouse.StartDrag()
			}
			gpu.SolidRR(ctx2.Rect, 2, theme.SurfaceContainer.Fg().Alpha(alpha))
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
		return wid.Dim{0, 0, 0}
	}
}
