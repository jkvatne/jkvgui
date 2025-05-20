package sys

import (
	"flag"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/gpu/font"
	"github.com/jkvatne/jkvgui/mouse"
	"log/slog"
	"math"
	"time"
)

var (
	LastMods     glfw.ModifierKey
	scrolledY    float32
	RedrawsPrSec int
	// ZoomFactor is the factor by which the window is zoomed when ctrl+scrollwheel is used.
	ZoomFactor  = float32(math.Sqrt(math.Sqrt(2.0)))
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
		gpu.UpdateSize(w)
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
	gpu.UpdateSize(w)
	gpu.UpdateResolution()
	gpu.Invalidate(0)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	sizeCallback(w, width, height)
}

func InitializeWindow(window *glfw.Window) {
	font.LoadDefaultFonts()
	gpu.LoadIcons()
	gpu.UpdateResolution()

	window.SetMouseButtonCallback(mouse.BtnCallback)
	window.SetCursorPosCallback(mouse.PosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetCharCallback(charCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetContentScaleCallback(scaleCallback)
	window.SetFocusCallback(focusCallback)
	window.SetSizeCallback(sizeCallback)
}

var maxFps = flag.Bool("maxfps", false, "Set to force redrawing as fast as possible")
var logLevel = flag.Int("loglevel", 0, "Set log level (8=Error, 4=Warning, 0=Info(default), -4=Debug)")

// Initialize will initialize the gui system.
// It must be called before any other function in this package.
// It parses flags and sets up default logging
func Initialize() {
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	GetBuildInfo()
	if *maxFps {
		MaxDelay = 0
	}
}

func Shutdown() {
	glfw.Terminate()
	TerminateProfiling()
}
