package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"log"
	"net/http"
)

type Dim struct {
	W        float32
	H        float32
	Baseline float32
}

type Mode int

const (
	RenderChildren Mode = iota
	CollectWidths
	CollectHeights
)

type Ctx struct {
	f32.Rect
	Baseline float32
	Disabled bool
	Mode     Mode
}

func (ctx Ctx) Alpha() float32 {
	if ctx.Disabled {
		return 0.3
	}
	return 1.0
}

func (ctx Ctx) Disable() Ctx {
	ctx.Disabled = true
	return ctx
}

func (ctx Ctx) Enable(enabled bool) Ctx {
	ctx.Disabled = !enabled
	return ctx
}

func DisableIf(disabler *bool, w Wid) Wid {
	return func(ctx Ctx) Dim {
		ctx.Disabled = *disabler
		return w(ctx)
	}
}

type Wid func(ctx Ctx) Dim

func NewCtx() Ctx {
	return Ctx{Rect: f32.Rect{X: 0, Y: 0, W: gpu.WindowWidthDp, H: gpu.WindowHeightDp}, Baseline: 0}
}

// Show is used to paint a given widget directly to the screen at
// given coordinates. Skipping all layout finctions.
func Show(x, y, w float32, widget Wid) {
	ctx := Ctx{Mode: CollectWidths}
	ctx.Rect.W = w
	// First calculate minimum dimensions by calling with empty ctx
	dim := widget(ctx)
	// Set minimum size and given x,y coordinates.
	ctx.Rect.W = dim.W
	ctx.Rect.H = dim.H
	ctx.Rect.X = x
	ctx.Rect.Y = y
	ctx.Baseline = dim.Baseline
	// Call again to paint the widget
	ctx.Mode = RenderChildren
	_ = widget(ctx)
	mouse.Reset()
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{H: 0.01}
	}
}

func StartProfiler() {
	// View at	http://localhost:6060/debug/pprof/heap
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Printf("pprof server failed: %v", err)
		}
	}()
}
