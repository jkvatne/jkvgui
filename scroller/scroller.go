package scroller

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

type Style struct {
}

type State struct {
	Xpos  float32
	Ypos  float32
	Width float32
	Max   float32
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
		if ctx.Rect.H == 0 {
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
		for i, w := range widgets {
			ctx1.Rect.H = dims[i].H
			w(ctx1)
			ctx1.Rect.Y += dims[i].H
		}
		// Draw scrollbar
		ctx.Rect.X += ctx.Rect.W - 8
		ctx.Rect.W = 8
		gpu.SolidRR(ctx.Rect, 2, theme.SurfaceContainer.Bg().Alpha(0.2))
		// Draw thumb
		ctx.Rect.X += 0.5
		ctx.Rect.W -= 1.0
		ctx.Rect.Y = state.Ypos
		ctx.Rect.H *= 0.1
		gpu.SolidRR(ctx.Rect, 2, f32.Black.Alpha(0.9))
		return wid.Dim{100, sumH, 0}
	}
}
