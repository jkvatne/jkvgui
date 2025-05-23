package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
)

// VertResizer provides a draggable handle in between two widgets for resizing their area.

type ResizerState struct {
	// Pos can be -W/2..+W/2. Zero means divide in two equal parts.
	pos      float32
	dragging bool
	StartPos float32
}

var ResizerWidth = float32(2.0)

func VertResizer(state *ResizerState, widget1 Wid, widget2 Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: ctx.W, H: ctx.H, Baseline: ctx.Baseline}
		}
		state.dragging = state.dragging && mouse.LeftBtnDown()
		if state.dragging {
			// Mouse dragging divider
			if dx := mouse.Pos().X - state.StartPos; dx != 0 {
				state.pos = min(max(state.pos+dx, -ctx.W/2), ctx.W/2-ResizerWidth)
				gpu.Invalidate(0)
				slog.Info("Drag", "dy", dx, "pos", state.pos, "ctx.W", ctx.W, "ctx.H", ctx.H)
			}
			state.StartPos = mouse.StartDrag().X
		}

		ctx1 := ctx
		ctx2 := ctx
		ctx1.W = ctx.W/2 + state.pos - ResizerWidth/2
		ctx2.W = ctx.W - ctx1.W - ResizerWidth/2
		ctx2.X = ctx.X + ctx.W/2 + state.pos + ResizerWidth/2
		spacerRect := f32.Rect{X: ctx2.X - ResizerWidth/2, Y: ctx1.Y, W: ResizerWidth, H: ctx.H}
		widget1(ctx1)
		widget2(ctx2)
		gpu.Rect(spacerRect, 0.0, theme.SurfaceContainer.Fg(), theme.SurfaceContainer.Fg())
		// Start dragging if mouse pressed
		if mouse.LeftBtnPressed(spacerRect) && !state.dragging {
			state.dragging = true
			state.StartPos = mouse.StartDrag().X
			slog.Info("Start drag", "pos", state.pos, "state.StartPos", state.StartPos)
		}
		if mouse.Pos().Inside(spacerRect) {
			gpu.SetHresizeCursor()
		}
		return Dim{W: ctx.W, H: ctx.H, Baseline: ctx.Baseline}
	}
}

func HorResizer(state *ResizerState, widget1 Wid, widget2 Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{W: ctx.W, H: ctx.H, Baseline: ctx.Baseline}
		}
		state.dragging = state.dragging && mouse.LeftBtnDown()
		if state.dragging {
			// Mouse dragging divider
			if dy := mouse.Pos().Y - state.StartPos; dy != 0 {
				state.pos = min(max(state.pos+dy, -ctx.H/2), ctx.H/2-ResizerWidth)
				gpu.Invalidate(0)
				slog.Info("Drag", "dy", dy, "pos", state.pos, "ctx.W", ctx.W, "ctx.H", ctx.H)
			}
			state.StartPos = mouse.StartDrag().Y
		}

		ctx1 := ctx
		ctx2 := ctx
		ctx1.H = ctx.H/2 + state.pos - ResizerWidth/2
		ctx2.H = ctx.H - ctx1.H - ResizerWidth/2
		ctx2.Y = ctx.X + ctx.H/2 + state.pos + ResizerWidth/2
		spacerRect := f32.Rect{X: ctx.X, Y: ctx2.Y - ResizerWidth, W: ctx.W, H: ResizerWidth}
		widget1(ctx1)
		widget2(ctx2)
		gpu.Rect(spacerRect, 0.0, theme.SurfaceContainer.Fg(), theme.SurfaceContainer.Fg())
		// Start dragging if mouse pressed
		if mouse.LeftBtnPressed(spacerRect) && !state.dragging {
			state.dragging = true
			state.StartPos = mouse.StartDrag().Y
			slog.Info("Start drag", "pos", state.pos, "state.StartPos", state.StartPos)
		}
		if mouse.Pos().Inside(spacerRect) {
			gpu.SetVresizeCursor()
		}

		return Dim{W: ctx.W, H: ctx.H, Baseline: ctx.Baseline}
	}
}
