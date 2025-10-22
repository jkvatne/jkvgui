package sys

import (
	"go-pat/lib"
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
func (w *Window) MousePos() f32.Pos {
	return w.mousePos
}

// StartDrag is called when a widget wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func (w *Window) StartDrag() f32.Pos {
	w.Dragging = true
	w.DragStartPos = w.mousePos
	return w.mousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func (w *Window) Hovered(r f32.Rect) bool {
	if w.SuppressEvents {
		return false
	}
	if w.Dragging {
		return false
	}
	if w.mousePos.Inside(r) {
		return true
	}
	return false
}

func HasMoved(p1, p2 f32.Pos) bool {
	diff := lib.Abs(p1.X-p2.X) + lib.Abs(p1.Y-p2.Y)
	return diff > DragMinDelta
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func (w *Window) LeftBtnPressed(r f32.Rect) bool {
	if w.SuppressEvents || w.Dragging && HasMoved(w.DragStartPos, w.MousePos()) || !w.LeftBtnIsDown || !w.mousePos.Inside(r) {
		return false
	}
	slog.Debug("LeftBtnPressed", "MouseX", int(w.MousePos().X), "MouseY", int(w.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
	return true
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func (w *Window) LeftBtnDown() bool {
	if w.SuppressEvents {
		return false
	}
	return w.LeftBtnIsDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func (w *Window) LeftBtnClick(r f32.Rect) bool {
	if !w.SuppressEvents && w.mousePos.Inside(r) && time.Since(w.LeftBtnDownTime) < LongPressTime && w.LeftBtnClicked {
		slog.Debug("LeftBtnClick", "MouseX", int(w.MousePos().X), "MouseY", int(w.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
		w.LeftBtnClicked = false
		return true
	}
	return false
}

// ClearMouseBtns is called when a window looses focus. It will reset the mouse button states.
func (w *Window) ClearMouseBtns() {
	w.LeftBtnIsDown = false
	w.Dragging = false
	w.LeftBtnDoubleClicked = false
	w.LeftBtnClicked = false
	w.ScrolledDistY = 0.0
}

// LeftBtnDoubleClick indicates that the user is holding the left btn down
// independent of the mouse pointer location
func (w *Window) LeftBtnDoubleClick(r f32.Rect) bool {
	if !w.SuppressEvents && w.mousePos.Inside(r) && w.LeftBtnDoubleClicked {
		w.LeftBtnDoubleClicked = false
		slog.Debug("LeftBtnDoubleClick:", "X", int(w.MousePos().X), "Y", int(w.MousePos().Y), "r.x", int(r.X), "r.y", int(r.Y), "r.W", int(r.W), "r.H", int(r.H))
		return true
	}
	return false
}

func (w *Window) SimPos(x, y float32) {
	w.mousePos.X = x
	w.mousePos.Y = y
}

func (w *Window) SimLeftBtnPress() {
	w.LeftBtnIsDown = true
	w.LeftBtnDownTime = time.Now()
}

func (w *Window) SimLeftBtnRelease() {
	w.LeftBtnIsDown = false
	w.Dragging = false
	if time.Since(w.LeftBtnUpTime) < DoubleClickTime {
		w.LeftBtnDoubleClicked = true
	} else {
		w.LeftBtnClicked = false
	}
	w.LeftBtnUpTime = time.Now()
}

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func (w *Window) ScrolledY() float32 {
	if !w.Focused {
		return 0
	}
	if w.SuppressEvents {
		return 0.0
	}
	s := w.ScrolledDistY
	w.ScrolledDistY = 0.0
	return s
}
