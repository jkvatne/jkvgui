package sys

import (
	"log/slog"
	"math"
	"time"

	"github.com/jkvatne/jkvgui/f32"
)

var (
	LongPressTime   = time.Millisecond * 700
	DoubleClickTime = time.Millisecond * 330
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor   = float32(math.Sqrt(math.Sqrt(2.0)))
	DragMinDelta = float32(4)
)

// MousePos is the mouse pointer location in device-independent screen coordinates
func (win *Window) MousePos() f32.Pos {
	return win.mousePos
}

// StartDrag is called when a widget wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func (win *Window) StartDrag() f32.Pos {
	win.Dragging = true
	win.DragStartPos = win.mousePos
	return win.mousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func (win *Window) Hovered(r f32.Rect) bool {
	if win.SuppressEvents {
		return false
	}
	if win.Dragging {
		return false
	}
	if win.mousePos.Inside(r) {
		return true
	}
	return false
}

func HasMoved(p1, p2 f32.Pos) bool {
	return f32.DiffXY(p1, p2) > DragMinDelta
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func (win *Window) LeftBtnPressed(r f32.Rect) bool {
	if win.SuppressEvents || win.Dragging && HasMoved(win.DragStartPos, win.MousePos()) || !win.LeftBtnIsDown || !win.mousePos.Inside(r) {
		return false
	}
	slog.Debug("LeftBtnPressed", "MouseX", int(win.MousePos().X), "MouseY", int(win.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
	return true
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func (win *Window) LeftBtnDown() bool {
	if win.SuppressEvents {
		return false
	}
	return win.LeftBtnIsDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func (win *Window) LeftBtnClick(r f32.Rect) bool {
	if !win.SuppressEvents && win.mousePos.Inside(r) && time.Since(win.LeftBtnDownTime) < LongPressTime && win.LeftBtnClicked {
		slog.Debug("LeftBtnClick", "MouseX", int(win.MousePos().X), "MouseY", int(win.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
		win.LeftBtnClicked = false
		return true
	}
	return false
}

// LeftBtnDoubleClick indicates that the user is holding the left btn down
// independent of the mouse pointer location
func (win *Window) LeftBtnDoubleClick(r f32.Rect) bool {
	if !win.SuppressEvents && win.mousePos.Inside(r) && win.LeftBtnDoubleClicked {
		win.LeftBtnDoubleClicked = false
		slog.Debug("LeftBtnDoubleClick:", "X", int(win.MousePos().X), "Y", int(win.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
		return true
	}
	return false
}

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func (win *Window) ScrolledY() float32 {
	if !win.Focused {
		return 0
	}
	if win.SuppressEvents {
		return 0.0
	}
	s := win.ScrolledDistY
	win.ScrolledDistY = 0.0
	return s
}
