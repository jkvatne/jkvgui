package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
	"math"
)

type ScrollState struct {
	Xpos     float32
	Ypos     float32
	Ymax     float32
	Dy       float32
	Npos     int
	dragging bool
	StartPos float32
	Width    float32
	Height   float32
	AtEnd    bool
}

var (
	ScrollbarWidth    = float32(10.0)
	MinThumbHeight    = float32(15.0)
	TrackAlpha        = float32(0.15)
	NormalAlpha       = float32(0.3)
	HoverAlpha        = float32(0.8)
	ScrollerMargin    = float32(1.0)
	ThumbCornerRadius = float32(5.0)
)

// DrawVertScrollbar will draw a bar at the right edge of the area r.
func VertScollbarUserInput(Yvis float32, state *ScrollState) float32 {
	state.dragging = state.dragging && mouse.LeftBtnDown()
	dy := float32(0.0)
	if state.dragging {
		// Mouse dragging scroller thumb
		dy = -(mouse.Pos().Y - state.StartPos) * state.Ymax / Yvis
		if dy != 0 {
			state.StartPos = mouse.Pos().Y
			gpu.Invalidate(0)
			slog.Debug("Drag", "dy", dy, "Ypos", int(state.Ypos), "state.Ymax", int(state.Ymax), "Yvis", int(Yvis), "state.StartPos", int(state.StartPos), "NotAtEnd", state.Ypos < state.Ymax-Yvis-0.01)
		}
	}
	if scr := sys.ScrolledY(); scr != 0 {
		// Handle mouse scroll-wheel. Scrolling down gives negative scr value
		dy = -(scr * Yvis) / 30
		gpu.Invalidate(0)
	}
	if dy < 0 {
		// Scrolling up means no more at end
		state.AtEnd = false
	}
	dy = float32(math.Round(float64(dy)))
	return dy
}

// DrawVertScrollbar will draw a bar at the right edge of the area r.
// state.Ypos is the position. (Ymax-Yvis) is max Ypos. Yvis is the visible part
func DrawVertScrollbar(barRect f32.Rect, Ymax float32, Yvis float32, state *ScrollState) {
	if Yvis > Ymax {
		return
	}
	barRect = f32.Rect{barRect.X + barRect.W - ScrollbarWidth, barRect.Y + ScrollerMargin, ScrollbarWidth, barRect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, Yvis*barRect.H/Ymax))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (Ymax - Yvis)
	if state.AtEnd {
		thumbPos = barRect.H - thumbHeight
	}
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

// DrawFromBottom will draw the last widgets from bottom up
// sumH will be the total heigth of the drawn widges, normally greater than the ctx.H
// dims is the dimension of the drawn widgets where dims[0] is at bottom.
func DrawFromBottom(ctx Ctx, widgets ...Wid) (sumH float32, dims []Dim) {
	ctx0 := ctx
	ctx0.Rect.Y += ctx0.Rect.H
	n := 0
	for i := len(widgets) - 1; i >= 0 && sumH < ctx.Rect.H; i-- {
		// Find height of current widget
		ctx0.Mode = CollectHeights
		ctx0.H = ctx.H
		dim := widgets[i](ctx0)
		sumH += dim.H
		// Draw it from y-H
		ctx0.Y -= dim.H
		ctx0.H = dim.H
		ctx0.Mode = RenderChildren
		dims = append(dims, widgets[i](ctx0))
		n++
	}

	// Verify sumH
	tempH := float32(0.0)
	for i := 0; i < len(dims); i++ {
		tempH += dims[i].H
	}
	if tempH != sumH {
		slog.Error("DrawFromBottom with diverging heights", "sumH", sumH, "tempH", tempH)
	}
	return sumH, dims
}

// DrawFromPos will draw widgets from state.Npos and downwards, with offset state.Dy
// It returns the total height and dimensions of all drawn widgets
func DrawFromPos(ctx Ctx, state *ScrollState, widgets ...Wid) (sumH float32, dims []Dim) {
	ctx0 := ctx
	ctx0.Rect.Y -= state.Dy
	gpu.Clip(ctx.Rect)
	for i := state.Npos; i < len(widgets) && sumH < ctx.Rect.H*2; i++ {
		ctx0.Rect.H = 0
		dim := widgets[i](ctx0)
		ctx0.Rect.Y += dim.H
		sumH += dim.H
		dims = append(dims, dim)
	}
	gpu.NoClip()
	return sumH, dims
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func Scroller(state *ScrollState, widgets ...Wid) Wid {
	f32.ExitIf(state == nil, "Scroller state must not be nil")

	return func(ctx Ctx) Dim {
		ctx0 := ctx
		if ctx.Mode != RenderChildren {
			return Dim{W: state.Width, H: state.Height, Baseline: 0}
		}
		dy := VertScollbarUserInput(ctx.Rect.H, state)
		if dy != 0.0 {
			slog.Info("Before Scroll", "dy", int(dy), "Dy", int(state.Dy), "state.Ypos", int(state.Ypos),
				"state.Ymax", int(state.Ymax), "Yvis", int(ctx.H), "Npos", state.Npos)
			state.Ypos += dy
			state.Dy += dy
			if state.Ypos < 0 {
				state.Ypos = 0
				state.Dy = 0
			}
			slog.Info("After Scroll", "dy", int(dy), "Dy", int(state.Dy), "state.Ypos", int(state.Ypos),
				"state.Ymax", int(state.Ymax), "Yvis", int(ctx.H), "Npos", state.Npos)
		}

		if state.AtEnd {
			// Draw from bottom up
			sumH, dims := DrawFromBottom(ctx, widgets...)
			state.Dy = sumH - ctx.Rect.H
			state.Npos = len(widgets) - len(dims)
			state.Ypos = state.Ymax - ctx.Rect.H
			DrawVertScrollbar(ctx.Rect, state.Ymax, ctx.Rect.H, state)
			return Dim{ctx.W, ctx.H, 0}
		}
		state.Npos = min(state.Npos, len(widgets)-1)

		sumH, dims := DrawFromPos(ctx0, state, widgets...)
		if state.Dy < 0 {
			// We need to reduce Npos
			state.Npos = max(0, state.Npos-1)
			state.Dy += dims[0].H
		}
		yvis := ctx.Rect.H
		newYmax := state.Ypos - state.Dy + sumH
		if abs(newYmax-state.Ymax) > 0.01 {
			slog.Info("Ymax updated", "was", int(state.Ymax), "new", int(newYmax), "state.Ypos",
				int(state.Ypos), "state.Dy", int(state.Dy), "sumH", int(sumH), "yvis", int(yvis), "state.Npos", state.Npos)
			state.Ymax = newYmax
		}

		if state.Ymax < yvis {
			slog.Info("Ymax set to yvis", "was", state.Ymax, "new", yvis)
			state.Ymax = yvis
		}
		if state.Npos+len(dims) < len(widgets) {
			// TODO state.Ymax = max(sumH, float32(len(widgets))*sumH/float32(len(dims)+state.Npos))
		}
		DrawVertScrollbar(ctx.Rect, state.Ymax, yvis, state)
		if state.Npos+len(dims) == len(widgets) && state.Ypos+yvis >= state.Ymax {
			// At end
			state.Dy = sumH - yvis
			state.Ypos = state.Ymax - yvis
			if state.Ypos < 0 {
				state.Ypos = 0
			}
			state.AtEnd = true
		} else if state.Dy > dims[0].H {
			// Ignore top widget, as it is no longer visible
			state.Dy -= dims[0].H
			state.Npos++
			gpu.Invalidate(0)
			slog.Info("Scrolled beyond top widget", "Npos", state.Npos, "Dy", state.Dy)
		} else if state.Dy < 0 {
			// Scrolling up.
			state.AtEnd = false
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
