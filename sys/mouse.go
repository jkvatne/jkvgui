package sys

import (
	"log/slog"
	"math"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

var (
	LongPressTime   = time.Millisecond * 700
	DoubleClickTime = time.Millisecond * 330
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor = float32(math.Sqrt(math.Sqrt(2.0)))
)

// Pos is the mouse pointer location in device-independent screen coordinates
func Pos() f32.Pos {
	return gpu.CurrentInfo.MousePos
}

// StartDrag is called when a widget wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func StartDrag() f32.Pos {
	gpu.CurrentInfo.Dragging = true
	return gpu.CurrentInfo.MousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func Hovered(r f32.Rect) bool {
	if !gpu.CurrentInfo.Focused {
		return false
	}
	if gpu.CurrentInfo.SuppressEvents {
		return false
	}
	return gpu.CurrentInfo.MousePos.Inside(r) && !gpu.CurrentInfo.Dragging
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	if !gpu.CurrentInfo.Focused {
		return false
	}
	if gpu.CurrentInfo.SuppressEvents {
		return false
	}
	return gpu.CurrentInfo.MousePos.Inside(r) && gpu.CurrentInfo.LeftBtnDown && !gpu.CurrentInfo.Dragging
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	if !gpu.CurrentInfo.Focused {
		slog.Info("LeftBtnDown but not focused", "Wno", gpu.CurrentInfo.Wno, "Name", gpu.CurrentInfo.Name)
		return false
	}
	if gpu.CurrentInfo.SuppressEvents {
		return false
	}
	return gpu.CurrentInfo.LeftBtnDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if gpu.CurrentInfo.SuppressEvents {
		return false
	}
	if !gpu.CurrentInfo.Focused {
		return false
	}
	if gpu.CurrentInfo.MousePos.Inside(r) && gpu.CurrentInfo.LeftBtnReleased && time.Since(gpu.CurrentInfo.LeftBtnDownTime) < LongPressTime {
		gpu.CurrentInfo.LeftBtnReleased = false
		return true
	}
	return false
}

// Reset is called when a window looses focus. It will reset the btn states.
func Reset() {
	gpu.CurrentInfo.LeftBtnDown = false
	gpu.CurrentInfo.LeftBtnReleased = false
	gpu.CurrentInfo.Dragging = false
	gpu.CurrentInfo.LeftBtnDoubleClick = false
}

func ClearMouseBtns() {
	gpu.CurrentInfo.LeftBtnReleased = false
	gpu.CurrentInfo.LeftBtnDoubleClick = false
}

// gpu.CurrentInfo.LeftBtnDoubleClick indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDoubleClick(r f32.Rect) bool {
	if !gpu.CurrentInfo.Focused {
		return false
	}
	if !gpu.CurrentInfo.SuppressEvents && gpu.CurrentInfo.MousePos.Inside(r) && gpu.CurrentInfo.LeftBtnDoubleClick {
		return gpu.CurrentInfo.LeftBtnDoubleClick
	}
	return false
}

func SimPos(x, y float32) {
	gpu.CurrentInfo.MousePos.X = x
	gpu.CurrentInfo.MousePos.Y = y
}

func SimLeftBtnPress() {
	gpu.CurrentInfo.LeftBtnDown = true
	gpu.CurrentInfo.LeftBtnDownTime = time.Now()
}

func SimLeftBtnRelease() {
	gpu.CurrentInfo.LeftBtnDown = false
	gpu.CurrentInfo.LeftBtnReleased = true
	gpu.CurrentInfo.Dragging = false
	if time.Since(gpu.CurrentInfo.LeftBtnUpTime) < DoubleClickTime {
		gpu.CurrentInfo.LeftBtnDoubleClick = true
	}
	gpu.CurrentInfo.LeftBtnUpTime = time.Now()
}

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func ScrolledY() float32 {
	if !gpu.CurrentInfo.Focused {
		return 0
	}
	if gpu.CurrentInfo.SuppressEvents {
		return 0.0
	}
	s := gpu.CurrentInfo.ScrolledY
	gpu.CurrentInfo.ScrolledY = 0.0
	return s
}
