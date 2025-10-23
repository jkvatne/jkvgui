package wid

import (
	"flag"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/sys"
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
	Win      *sys.Window
}

var DebugWidgets = flag.Bool("debug", true, "Set to debug widgets and write font info")

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

func (ctx Ctx) SetCursor(id int) {
	ctx.Win.SetCursor(sys.VResizeCursor)
	ctx.Win.Cursor = id
}

type Wid func(ctx Ctx) Dim

func Show(w Wid) {
	win := sys.GetCurrentWindow()
	if win == nil || win.Window.ShouldClose() {
		return
	}
	w(NewCtx(win))
}

func NewCtx(win *sys.Window) Ctx {
	return Ctx{Rect: f32.Rect{W: win.WidthDp, H: win.HeightDp}, Baseline: 0, Win: win}
}

// Display is used to paint a given widget directly to the screen at
// given coordinates. Skipping all layout functions.
func Display(win *sys.Window, x, y, w float32, widget Wid) {
	ctx := NewCtx(win)
	ctx.Mode = CollectWidths
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
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		if ctx.Mode != RenderChildren {
			return Dim{H: 0.01, W: 0}
		}
		return Dim{H: ctx.H, W: ctx.W}
	}
}
