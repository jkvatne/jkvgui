package mouse

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"log/slog"
	"time"
)

var (
	mousePos        f32.Pos
	leftBtnDown     bool
	leftBtnReleased bool
	dragging        bool
	leftBtnDonwTime time.Time
)

var LongPressTime = time.Millisecond * 700

// Pos is the mouse pointer location in device-independent screen coordinates
func Pos() f32.Pos {
	return mousePos
}

// StartDrag is called when a widges wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func StartDrag() f32.Pos {
	dragging = true
	return mousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func Hovered(r f32.Rect) bool {
	return mousePos.Inside(r) && !dragging
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the button is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	return mousePos.Inside(r) && leftBtnDown && !dragging
}

// LeftBtnDown indicates that the user is holding the left button down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	return leftBtnDown
}

// LeftBtnClick returns true if the left button has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if mousePos.Inside(r) && leftBtnReleased && time.Since(leftBtnDonwTime) < LongPressTime {
		leftBtnReleased = false
		return true
	}
	return false
}

// Reset is called when a window looses focus. It will reset the button states.
func Reset() {
	leftBtnDown = false
	leftBtnReleased = false
	dragging = false
}

// BtnCallback is called from the glfw window handler when mouse buttons change states.
func BtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	gpu.Invalidate(0)
	x, y := w.GetCursorPos()
	mousePos.X = float32(x) / gpu.ScaleX
	mousePos.Y = float32(y) / gpu.ScaleY
	slog.Debug("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action)
	if action == glfw.Release {
		leftBtnDown = false
		leftBtnReleased = true
		dragging = false
	} else if action == glfw.Press {
		leftBtnDown = true
		leftBtnDonwTime = time.Now()
	}
}

// PosCallback is called from the glfw window handler when the mouse moves.
func PosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	mousePos.X = float32(xpos) / gpu.ScaleX
	mousePos.Y = float32(ypos) / gpu.ScaleY
	gpu.Invalidate(20 * time.Millisecond)
}
