package sys

import (
	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"math"
	"time"
)

var (
	mousePos           f32.Pos
	leftBtnDown        bool
	leftBtnReleased    bool
	dragging           bool
	leftBtnDownTime    time.Time
	LongPressTime      = time.Millisecond * 700
	DoubleClickTime    = time.Millisecond * 330
	leftBtnUpTime      = time.Now()
	leftBtnDoubleClick bool
	scrolledY          float32
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor = float32(math.Sqrt(math.Sqrt(2.0)))
)

// Pos is the mouse pointer location in device-independent screen coordinates
func Pos() f32.Pos {
	return mousePos
}

// StartDrag is called when a widget wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func StartDrag() f32.Pos {
	dragging = true
	return mousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func Hovered(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && !dragging
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && leftBtnDown && !dragging
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	if SuppressEvents {
		return false
	}
	return leftBtnDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	if mousePos.Inside(r) && leftBtnReleased && time.Since(leftBtnDownTime) < LongPressTime {
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
	if !SuppressEvents && mousePos.Inside(r) && leftBtnDoubleClick {
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
	leftBtnDownTime = time.Now()
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

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func ScrolledY() float32 {
	if SuppressEvents {
		return 0.0
	}
	s := scrolledY
	scrolledY = 0.0
	return s
}
