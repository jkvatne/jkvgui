package sys

import (
	"math"
)

var (
	scrolledY float32
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor = float32(math.Sqrt(math.Sqrt(2.0)))
)

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
