package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/gpu"
	"log/slog"
	"math"
)

var (
	LastMods  glfw.ModifierKey
	scrolledY float32
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor = float32(math.Sqrt(math.Sqrt(2.0)))
)

func setCallbacks() {
	Window.SetMouseButtonCallback(BtnCallback)
	Window.SetCursorPosCallback(PosCallback)
	Window.SetKeyCallback(keyCallback)
	Window.SetCharCallback(charCallback)
	Window.SetScrollCallback(scrollCallback)
	Window.SetContentScaleCallback(scaleCallback)
	Window.SetFocusCallback(focusCallback)
	Window.SetSizeCallback(sizeCallback)
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	gpu.Invalidate(0)
	if key == glfw.KeyTab && action == glfw.Release {
		MoveByKey(mods != glfw.ModShift)
	}
	if action == glfw.Release {
		LastKey = key
	}
	LastMods = mods
}

// ScrolledY returns the amount of pixels scrolled vertically since the last call to this function.
// If gpu.SuppressEvents is true, the return value is always 0.0.
func ScrolledY() float32 {
	if gpu.SuppressEvents {
		return 0.0
	}
	s := scrolledY
	scrolledY = 0.0
	return s
}

func charCallback(w *glfw.Window, char rune) {
	slog.Debug("charCallback()", "Rune", int(char))
	gpu.Invalidate(0)
	gpu.LastRune = char
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("Scroll", "dx", xoff, "dy", yOff)
	if LastMods == glfw.ModControl {
		// ctrl+scrollwheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			gpu.UserScale *= ZoomFactor
		} else {
			gpu.UserScale /= ZoomFactor
		}
		UpdateSize(w)
	} else {
		scrolledY = float32(yOff)
	}
	gpu.Invalidate(0)
}

func focusCallback(w *glfw.Window, focused bool) {
	gpu.WindowHasFocus = focused
	if !focused {
		Reset()
	}
	gpu.Invalidate(0)
}

func sizeCallback(w *glfw.Window, width int, height int) {
	UpdateSize(w)
	gpu.UpdateResolution()
	gpu.Invalidate(0)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	sizeCallback(w, width, height)
}
