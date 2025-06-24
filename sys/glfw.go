package sys

import (
	"flag"
	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/glfw"
	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
	"runtime"
	"time"
)

type Window glfw.Window

var (
	WindowList    []*glfw.Window
	VResizeCursor *glfw.Cursor
	HResizeCursor *glfw.Cursor
)

const (
	KeyRight     = glfw.KeyRight
	KeyLeft      = glfw.KeyLeft
	KeyUp        = glfw.KeyUp
	KeyDown      = glfw.KeyDown
	KeySpace     = glfw.KeySpace
	KeyEnter     = glfw.KeyEnter
	KeyEscape    = glfw.KeyEscape
	KeyBackspace = glfw.KeyBackspace
	KeyDelete    = glfw.KeyDelete
	KeyHome      = glfw.KeyHome
	KeyEnd       = glfw.KeyEnd
	KeyPageUp    = glfw.KeyPageUp
	KeyPageDown  = glfw.KeyPageDown
	KeyInsert    = glfw.KeyInsert
	KeyC         = glfw.KeyC
	KeyV         = glfw.KeyV
	KeyX         = glfw.KeyX
	ModShift     = glfw.ModShift
	ModControl   = glfw.ModControl
	ModAlt       = glfw.ModAlt
)

var (
	LastRune rune
	LastKey  glfw.Key
	LastMods glfw.ModifierKey
)

func PollEvents() {
	glfw.PollEvents()
}

func Shutdown() {
	glfw.Terminate()
	TerminateProfiling()
}

var Monitors []*glfw.Monitor

func init() {
	runtime.LockOSThread()
	flag.Parse()
	if *maxFps {
		MaxDelay = 0
	}
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	theme.SetDefaultPallete(true)
	setHints()
	// Check all monitors and print size data
	Monitors = GetMonitors()
	// Select monitor as given, or use primary monitor.
	for i, m := range Monitors {
		SizeMmX, SizeMmY := m.GetPhysicalSize()
		ScaleX, ScaleY := m.GetContentScale()
		PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
		slog.Info("CreateWindow()", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(ScaleX, 3), "ScaleY", f32.F2S(ScaleY, 3))
	}
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	windowHasFocus = focused
	if !focused {
		resetFocus()
	}
	gpu.Invalidate(0)
}

func setCallbacks(Window *glfw.Window) {
	Window.SetMouseButtonCallback(btnCallback)
	Window.SetCursorPosCallback(posCallback)
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
		moveByKey(mods != glfw.ModShift)
	}
	if action == glfw.Release || action == glfw.Repeat {
		LastKey = key
	}
	LastMods = mods
}

func Return() bool {
	return LastKey == glfw.KeyEnter || LastKey == glfw.KeyKPEnter
}

func charCallback(w *glfw.Window, char rune) {
	slog.Debug("charCallback()", "Rune", int(char))
	gpu.Invalidate(0)
	LastRune = char
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	gpu.Invalidate(0)
	LastMods = mods
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
			leftBtnDownTime = time.Now()
		}
	}
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(xw *glfw.Window, xpos float64, ypos float64) {
	mousePos.X = float32(xpos) / gpu.ScaleX
	mousePos.Y = float32(ypos) / gpu.ScaleY
	gpu.Invalidate(0 * time.Millisecond)
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

func UpdateSize(w *glfw.Window) {
	width, height := WindowList[0].GetSize()
	gpu.WindowHeightPx = height
	gpu.WindowWidthPx = width
	gpu.ScaleX, gpu.ScaleY = w.GetContentScale()
	gpu.ScaleX *= gpu.UserScale
	gpu.ScaleY *= gpu.UserScale
	WindowWidthDp = float32(width) / gpu.ScaleX
	WindowHeightDp = float32(height) / gpu.ScaleY
	gpu.WindowRect = f32.Rect{W: WindowWidthDp, H: WindowHeightDp}
	slog.Info("UpdateSize", "w", width, "h", height, "scaleX", f32.F2S(gpu.ScaleX, 3),
		"ScaleY", f32.F2S(gpu.ScaleY, 3), "UserScale", f32.F2S(gpu.UserScale, 3))
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

func setHints() {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.Floating, glfw.False) // True will keep the window on top
}

func SetMaximizedHint(maximized bool) {
	if maximized {
		glfw.WindowHint(glfw.Maximized, glfw.True)
	} else {
		glfw.WindowHint(glfw.Maximized, glfw.False)
	}

}

func createWindow(w, h int, title string, monitor *glfw.Monitor) *glfw.Window {
	// Create invisible window so we can move it to correct monitor
	glfw.WindowHint(glfw.Visible, glfw.False)
	win, err := glfw.CreateWindow(w, h, title, monitor, nil)
	if err != nil || win == nil {
		panic(err)
	}
	return win
}

func SetupCursors() {
	VResizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)
	HResizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
}

func SetClipboardString(s string) {
	glfw.SetClipboardString(s)
}

func GetClipboardString() string {
	return glfw.GetClipboardString()
}

func MaximizeWindow(w *glfw.Window) {
	w.Maximize()
}
func MinimizeWindow(w *glfw.Window) {
	w.Iconify()
}

func MakeContextCurrent(w *glfw.Window) {
	w.MakeContextCurrent()
}
