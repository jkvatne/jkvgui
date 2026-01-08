package wid

import (
	"flag"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/sys"
)

// DebugWidgets is the flag set by command line argument -debug.
// Many widgets will then draw rectangles around their components (label/value).
var DebugWidgets = flag.Bool("debug", false, "Set to debug widgets and write font info")

type Mode int

const (
	RenderChildren Mode = iota
	CollectWidths
	CollectHeights
)

type Ctx struct {
	// Rect consists of the X,Y,W,H values. That is the size and position of the area to be drawn.
	f32.Rect
	Baseline float32
	Mode     Mode
	Win      *sys.Window
}

type Dim struct {
	W        float32
	H        float32
	Baseline float32
}

type Wid func(ctx Ctx) Dim

// SetCursor will update the cursor type in the current window
// This new cursor will be visible on next redraw
func (ctx Ctx) SetCursor(id int) {
	ctx.Win.SetCursor(sys.VResizeCursor)
	ctx.Win.Cursor = id
}

// NewCtx returns a new context with the current window size
func NewCtx(win *sys.Window) Ctx {
	return Ctx{Rect: f32.Rect{W: win.WidthDp, H: win.HeightDp}, Baseline: 0, Win: win}
}

// Show is used to display a form consisting of a widget.
// Typically, the widget is a column or a scroller.
func Show(w Wid) {
	win := sys.GetCurrentWindow()
	if win == nil || win.Window.ShouldClose() {
		return
	}
	ctx := NewCtx(win)
	if ctx.Rect.H > 0 && ctx.Rect.W > 0 {
		w(ctx)
	}
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

func Default[T any](ptr **T, def *T) {
	if *ptr == nil {
		*ptr = def
	}
}
