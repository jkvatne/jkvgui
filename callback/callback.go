package callback

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/focus"
	"github.com/jkvatne/jkvgui/font"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/mouse"
	"log/slog"
)

var LastMods glfw.ModifierKey
var ScrolledY float32

// https://www.glfw.org/docs/latest/window_guide.html
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
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

func CharCallback(w *glfw.Window, char rune) {
	slog.Info("charCallback()", "Rune", int(char))
	gpu.Invalidate(0)
	gpu.LastRune = char
}

var N = 10000
var ScrollArea = f32.Rect{0, 0, 9999, 9999}

func ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	slog.Info("Scroll", "dx", xoff, "dy", yoff)
	if LastMods == glfw.ModControl {
		if yoff > 0 {
			gpu.UserScale *= 1.1
		} else {
			gpu.UserScale *= 0.9
		}
		gpu.UpdateSize(w, gpu.WindowWidthPx, gpu.WindowHeightPx)
	} else if mouse.Hovered(ScrollArea) {
		ScrolledY = float32(yoff)
	}
	gpu.Invalidate(0)
}

func Initialize(window *glfw.Window) {
	font.LoadFonts()
	window.SetMouseButtonCallback(mouse.BtnCallback)
	window.SetCursorPosCallback(mouse.PosCallback)
	window.SetKeyCallback(KeyCallback)
	window.SetCharCallback(CharCallback)
	window.SetScrollCallback(ScrollCallback)
	gpu.LoadIcons()
	gpu.UpdateResolution()
}
