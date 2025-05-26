package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"log/slog"
	"time"
)

var (
	mousePos           f32.Pos
	leftBtnDown        bool
	leftBtnReleased    bool
	dragging           bool
	leftBtnDonwTime    time.Time
	LongPressTime      = time.Millisecond * 700
	DoubleClickTime    = time.Millisecond * 330
	leftBtnUpTime      = time.Now()
	leftBtnDoubleClick bool
)

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
	if gpu.SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && !dragging
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	if gpu.SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && leftBtnDown && !dragging
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	if gpu.SuppressEvents {
		return false
	}
	return leftBtnDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if gpu.SuppressEvents {
		return false
	}
	if mousePos.Inside(r) && leftBtnReleased && time.Since(leftBtnDonwTime) < LongPressTime {
		leftBtnReleased = false
		return true
	}
	return false
}

// Reset is called when a window looses focus. It will reset the btn states.
func Reset() {
	leftBtnDown = false
	leftBtnReleased = false
	dragging = false
	leftBtnDoubleClick = false
}

func FrameEnd() {
	leftBtnReleased = false
	leftBtnDoubleClick = false
}

// LeftBtnDoubleClick indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDoubleClick(r f32.Rect) bool {
	if !gpu.SuppressEvents && mousePos.Inside(r) && leftBtnDoubleClick {
		return leftBtnDoubleClick
	}
	return false
}

func SimPos(x, y float32) {
	mousePos.X = x
	mousePos.Y = y
}

func SimLeftBtnPress() {
	leftBtnDown = true
	leftBtnDonwTime = time.Now()
}

func SimLeftBtnRelease() {
	leftBtnDown = false
	leftBtnReleased = true
	dragging = false
	if time.Since(leftBtnUpTime) < DoubleClickTime {
		leftBtnDoubleClick = true
	}
	leftBtnUpTime = time.Now()
}

// BtnCallback is called from the glfw window handler when mouse buttons change states.
func BtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	gpu.Invalidate(0)
	x, y := w.GetCursorPos()
	mousePos.X = float32(x) / gpu.ScaleX
	mousePos.Y = float32(y) / gpu.ScaleY
	slog.Debug("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action)
	if button == glfw.MouseButtonLeft {
		if action == glfw.Release {
			leftBtnDown = false
			leftBtnReleased = true
			dragging = false
			if time.Since(leftBtnUpTime) < DoubleClickTime {
				leftBtnDoubleClick = true
			}
			leftBtnUpTime = time.Now()
		} else if action == glfw.Press {
			leftBtnDown = true
			leftBtnDonwTime = time.Now()
		}
	}
}

// PosCallback is called from the glfw window handler when the mouse moves.
func PosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	mousePos.X = float32(xpos) / gpu.ScaleX
	mousePos.Y = float32(ypos) / gpu.ScaleY
	gpu.Invalidate(0 * time.Millisecond)
}
