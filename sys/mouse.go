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
	ZoomFactor = float32(math.Sqrt(math.Sqrt(2.0)))
)

// Pos is the mouse pointer location in device-independent screen coordinates
func Pos() f32.Pos {
	return CurrentInfo.MousePos
}

// StartDrag is called when a widget wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func StartDrag() f32.Pos {
	CurrentInfo.Dragging = true
	return CurrentInfo.MousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func Hovered(r f32.Rect) bool {
	if !CurrentInfo.Focused {
		return false
	}
	if CurrentInfo.SuppressEvents {
		return false
	}
	if CurrentInfo.Dragging {
		return false
	}
	if CurrentInfo.MousePos.Inside(r) {
		return true
	}
	return false
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	if !CurrentInfo.Focused {
		return false
	}
	if CurrentInfo.SuppressEvents {
		return false
	}
	return CurrentInfo.MousePos.Inside(r) && CurrentInfo.LeftBtnDown && !CurrentInfo.Dragging
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	if !CurrentInfo.Focused {
		slog.Info("LeftBtnDown but not focused", "Wno", CurrentInfo.Wno, "Name", CurrentInfo.Name)
		return false
	}
	if CurrentInfo.SuppressEvents {
		return false
	}
	return CurrentInfo.LeftBtnDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if CurrentInfo.SuppressEvents {
		return false
	}
	if !CurrentInfo.Focused {
		return false
	}
	if CurrentInfo.MousePos.Inside(r) && CurrentInfo.LeftBtnReleased && time.Since(CurrentInfo.LeftBtnDownTime) < LongPressTime {
		CurrentInfo.LeftBtnReleased = false
		return true
	}
	return false
}

// Reset is called when a window looses focus. It will reset the btn states.
func Reset() {
	CurrentInfo.LeftBtnDown = false
	CurrentInfo.LeftBtnReleased = false
	CurrentInfo.Dragging = false
	CurrentInfo.LeftBtnDoubleClick = false
}

func ClearMouseBtns() {
	CurrentInfo.LeftBtnReleased = false
	CurrentInfo.LeftBtnDoubleClick = false
}

// CurrentInfo.LeftBtnDoubleClick indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDoubleClick(r f32.Rect) bool {
	if !CurrentInfo.Focused {
		return false
	}
	if !CurrentInfo.SuppressEvents && CurrentInfo.MousePos.Inside(r) && CurrentInfo.LeftBtnDoubleClick {
		return CurrentInfo.LeftBtnDoubleClick
	}
	return false
}

func SimPos(x, y float32) {
	CurrentInfo.MousePos.X = x
	CurrentInfo.MousePos.Y = y
}

func SimLeftBtnPress() {
	CurrentInfo.LeftBtnDown = true
	CurrentInfo.LeftBtnDownTime = time.Now()
}

func SimLeftBtnRelease() {
	CurrentInfo.LeftBtnDown = false
	CurrentInfo.LeftBtnReleased = true
	CurrentInfo.Dragging = false
	if time.Since(CurrentInfo.LeftBtnUpTime) < DoubleClickTime {
		CurrentInfo.LeftBtnDoubleClick = true
	}
	CurrentInfo.LeftBtnUpTime = time.Now()
}

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func ScrolledY() float32 {
	if !CurrentInfo.Focused {
		return 0
	}
	if CurrentInfo.SuppressEvents {
		return 0.0
	}
	s := CurrentInfo.ScrolledY
	CurrentInfo.ScrolledY = 0.0
	return s
}
