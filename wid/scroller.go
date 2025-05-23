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
	Nmax     int
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

// VertScollbarUserInput will draw a bar at the right edge of the area r.
func VertScollbarUserInput(Yvis float32, state *ScrollState) float32 {
	state.dragging = state.dragging && mouse.LeftBtnDown()
	dy := float32(0.0)
	if state.dragging {
		// Mouse dragging scroller thumb
		dy = (mouse.Pos().Y - state.StartPos) * state.Ymax / Yvis
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
	barRect = f32.Rect{X: barRect.X + barRect.W - ScrollbarWidth, Y: barRect.Y + ScrollerMargin, W: ScrollbarWidth, H: barRect.H - 2*ScrollerMargin}
	thumbHeight := min(barRect.H, max(MinThumbHeight, Yvis*barRect.H/Ymax))
	thumbPos := state.Ypos * (barRect.H - thumbHeight) / (Ymax - Yvis)
	if state.AtEnd {
		thumbPos = barRect.H - thumbHeight
	}
	thumbRect := f32.Rect{X: barRect.X + ScrollerMargin, Y: barRect.Y + thumbPos, W: ScrollbarWidth - ScrollerMargin*2, H: thumbHeight}
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

// DrawFromPos will draw widgets from state.Npos and downwards, with offset state.Dy
// It returns the total height and dimensions of all drawn widgets
func DrawFromPos(ctx Ctx, state *ScrollState, widgets ...Wid) (dims []Dim) {
	ctx0 := ctx
	ctx0.Rect.Y -= state.Dy
	sumH := -state.Dy
	gpu.Clip(ctx.Rect)
	for i := state.Npos; i < len(widgets) && sumH < ctx.Rect.H*1.5; i++ {
		ctx0.Rect.H = 0
		dim := widgets[i](ctx0)
		ctx0.Rect.Y += dim.H
		sumH += dim.H
		dims = append(dims, dim)
		if i >= state.Nmax {
			state.Ymax += dim.H
			state.Nmax = i + 1
		}
	}
	gpu.NoClip()
	return dims
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

		yScroll := VertScollbarUserInput(ctx.Rect.H, state)

		// Draw and return dimensions.
		// dims goes from Npos up to the last widget drawn.
		dims := DrawFromPos(ctx0, state, widgets...)
		dimsStart := state.Npos
		sumH := float32(0.0)
		for _, d := range dims {
			sumH += d.H
		}
		if yScroll < 0 {
			for yScroll < 0 {
				// Scroll up
				if -yScroll < state.Dy {
					state.Dy = max(0, state.Dy+yScroll)
					state.Ypos = max(0, state.Ypos+yScroll)
					slog.Info("Scroll up within widget", "yScroll", yScroll, "Ypos", state.Ypos, "Dy", state.Dy, "Npos", state.Npos)
					yScroll = 0
				} else if state.Npos > 0 {
					state.Npos = max(0, state.Npos-1)
					state.Ypos -= dims[state.Npos].H
					yScroll += dims[state.Npos].H
					if state.Ypos < 0 {
						state.Ypos = 0
					}
					state.Dy = 0
					slog.Info("Scroll up", "yScroll", yScroll, "Ypos", state.Ypos, "Dy", state.Dy, "Npos", state.Npos)
					yScroll = 0
				} else {
					slog.Info("At top", "Ypos was", state.Ypos, "Npos", state.Npos)
					state.Ypos = 0
					state.Dy = 0
					yScroll = 0
				}
			}
		} else if yScroll > 0 {
			i := 0
			j := dimsStart
			// Scroll down
			for yScroll > 0 {
				if yScroll+state.Dy < dims[i].H {
					// We are within the current widget.
					state.Ypos += yScroll
					state.Dy += yScroll
					slog.Info("Scroll down within widget", "yScroll", yScroll, "Ypos", state.Ypos, "Dy", state.Dy, "Npos", state.Npos)
					break
				} else {
					// Go down one widget
					state.Npos++
					state.Ypos += dims[i].H
					yScroll -= dims[i].H
					slog.Info("Scroll down one widget", "yScroll", yScroll, "Ypos", state.Ypos, "Dy", state.Dy, "Npos", state.Npos)
					i++
					j++
					if j >= len(widgets) {
						// No more widgets.
						slog.Info("No more widgets")
						break
					}
					if i >= len(dims) {
						ctx0.Mode = CollectHeights
						dims = append(dims, widgets[j](ctx0))
						sumH += dims[len(dims)-1].H
						state.Ymax += dims[len(dims)-1].H
						slog.Info("Next widget", "Ymax", state.Ymax)
					}
				}
			}
			// Handle end of widget list. We are now at Npos which starts at Ypos-Dy
			hTot := -state.Dy
			i = state.Npos
			j = 0
			// Walk down the rest of the visible widgets
			for i < len(widgets) && hTot < ctx.H {
				hTot += dims[j].H
				i++
				j++
				if j >= len(dims) && i < len(widgets) {
					ctx0.Mode = CollectHeights
					dims = append(dims, widgets[i](ctx0))
				}
			}
			// We terminated because we reached the end of the widget list.
			// That means we must re-align from the bottom,
			if i == len(widgets) && hTot <= ctx.H {

				h := float32(0.0)
				i = len(widgets) - 1
				j = len(dims) - 1
				// Loop up from bottom
				for j >= 0 && i >= 0 && h < ctx.H {
					h += dims[j].H
					i--
					j--
				}
				// Now recalculate Ypos and Npos
				state.Npos = max(0, i+1)
				state.Ypos = state.Ymax - ctx.H
				state.Dy = h - ctx.H
				slog.Info("At bottom", "Npos", state.Npos, "Ypos", state.Ypos, "Dy", state.Dy, "Ymax", state.Ymax)
			}

		}

		DrawVertScrollbar(ctx.Rect, state.Ymax, ctx.H, state)
		return Dim{ctx.W, ctx.H, 0}
	}
}
