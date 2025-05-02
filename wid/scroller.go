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
	Dy       float32
	Npos     int
	dragging bool
	StartPos float32
	Width    float32
	Height   float32
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
	state.dragging = state.dragging && mouse.LeftBtnDown()
	if state.dragging {
		// Mouse dragging scroller thumb
		if dy := (mouse.Pos().Y - state.StartPos) * Ymax / Yvis; dy != 0 {
			state.Dy += dy
			state.Ypos += dy
			state.StartPos = mouse.Pos().Y
			gpu.Invalidate(0)
			slog.Debug("Drag", "dy", dy, "Ypos", int(state.Ypos), "Ymax", int(Ymax), "Yvis", int(Yvis), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < Ymax-Yvis-0.01)
		}
	}
	if scr := sys.ScrolledY(); scr != 0 {
		// Handle mouse scroll-wheel. Scrolling down gives negative scr value
		state.Dy -= scr * Yvis / 10
		state.Ypos -= scr * Yvis / 10
		gpu.Invalidate(0)
		slog.Info("Scroll", "Dy", int(state.Dy), "Ymax", int(Ymax), "Yvis", int(Yvis), "Npos", state.Npos)
	}
	state.Ypos = max(min(Ymax-Yvis, state.Ypos), 0)
	barRect = f32.Rect{barRect.X + barRect.W - ScrollbarWidth, barRect.Y + ScrollerMargin, ScrollbarWidth, barRect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, Yvis*barRect.H/Ymax))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (Ymax - Yvis)
	thumbRect := f32.Rect{barRect.X, barRect.Y + thumbPos, ScrollbarWidth - ScrollerMargin*2, thumbHeight}
	// Draw scrollbar track
	gpu.RoundedRect(barRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(TrackAlpha), f32.Transparent)
	// Draw thumb
	alpha := f32.Sel(mouse.Hovered(thumbRect) || state.dragging, NormalAlpha, HoverAlpha)
	gpu.RoundedRect(thumbRect, ThumbCornerRadius, 0.0, theme.SurfaceContainer.Fg().MultAlpha(alpha), f32.Transparent)
	// Start dragging if mouse pressed
	if mouse.LeftBtnPressed(thumbRect) && !state.dragging {
		state.dragging = true
		state.StartPos = mouse.StartDrag().Y
	}
}

func Scroller(state *ScrollState, widgets ...Wid) Wid {
	dims := make([]Dim, len(widgets))
	f32.ExitIf(state == nil, "Scroller state must not be nil")

	return func(ctx Ctx) Dim {
		ctx0 := ctx
		// Collect Height for all children
		// Fractional heights are not supported for scroller.
		ctx0.Mode = CollectHeights
		if ctx.Mode != RenderChildren {
			return Dim{W: state.Width, H: state.Height, Baseline: 0}
		}
		ctx0 = ctx
		ctx0.Rect.Y -= state.Dy
		sumH := float32(0.0)
		gpu.Clip(ctx.Rect)
		i := min(state.Npos, len(widgets))
		oldSum := float32(0)
		for i < len(widgets) && oldSum < ctx.Rect.H {
			ctx0.Rect.H = 0
			dims[i-state.Npos] = widgets[i](ctx0)
			ctx0.Rect.Y += dims[i-state.Npos].H
			oldSum = sumH
			sumH += dims[i-state.Npos].H
			i++
		}
		gpu.NoClip()
		ymax := float32(len(widgets)) * dims[0].H
		DrawVertScrollbar(ctx.Rect, ymax, ctx.Rect.H, state)
		if i == len(widgets) && state.Dy > dims[0].H {
			// At end
			state.Dy = sumH - ctx.Rect.H
			state.Ypos = ymax - ctx.Rect.H
		} else if state.Dy > dims[0].H {
			state.Dy -= dims[0].H
			state.Npos++
			gpu.Invalidate(0)
		} else if state.Dy < 0 {
			state.Dy += dims[0].H
			state.Npos--
			if state.Npos <= 0 {
				state.Npos = 0
				state.Dy = 0
				state.Ypos = 0
			}
			gpu.Invalidate(0)
		}
		return Dim{ctx.W, ctx.H, 0}
	}
}
