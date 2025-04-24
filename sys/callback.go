package sys

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"log/slog"
	"time"
)

var (
	LastMods     glfw.ModifierKey
	scrolledY    float32
	RedrawsPrSec int
)

var (
	startTime   time.Time
	redraws     int
	redrawStart time.Time
)

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	gpu.Invalidate(0)
	if key == glfw.KeyTab && action == glfw.Release {
		focus.MoveByKey(mods != glfw.ModShift)
	}
	if action == glfw.Release {
		gpu.LastKey = key
	}
	LastMods = mods
}

func ScrolledY() float32 {
	if gpu.SupressEvents {
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
			gpu.UserScale *= 1.05
		} else {
			gpu.UserScale *= 0.95
		}
		gpu.UpdateSize(w, gpu.WindowWidthPx, gpu.WindowHeightPx)
	} else {
		scrolledY = float32(yOff)
	}
	gpu.Invalidate(0)
}

func focusCallback(w *glfw.Window, focused bool) {
	gpu.WindowHasFocus = focused
	if !focused {
		mouse.Reset()
	}
	gpu.Invalidate(0)
}

func sizeCallback(w *glfw.Window, width int, height int) {
	gpu.UpdateSize(w, width, height)
	gpu.UpdateResolution()
	gpu.Invalidate(0)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	sizeCallback(w, width, height)
}

func Initialize(window *glfw.Window) {
	window.SetMouseButtonCallback(mouse.BtnCallback)
	window.SetCursorPosCallback(mouse.PosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetCharCallback(charCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetContentScaleCallback(scaleCallback)
	window.SetFocusCallback(focusCallback)
	window.SetSizeCallback(sizeCallback)
	font.LoadDefaultFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution()
}
