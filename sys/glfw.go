package sys

import (
	"flag"
	"github.com/jkvatne/jkvgui/buildinfo"
	"github.com/jkvatne/jkvgui/f32"
	// Using my own purego-glfw implementation:
	// glfw "github.com/jkvatne/purego-glfw"
	// Using standard go-gl from github:
	// "github.com/go-gl/glfw/v3.3/glfw"
	// Testing with glfw in local directory:
	"github.com/jkvatne/jkvgui/glfw"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/theme"
	"log/slog"
	"runtime"
	"time"
)

var Monitors []*glfw.Monitor
var maxFps = flag.Bool("maxfps", false, "Set to force redrawing as fast as possible")

type Window glfw.Window

var (
	WindowList       []*glfw.Window
	CurrentWindow    *glfw.Window
	CurrentWno       int
	pVResizeCursor   *glfw.Cursor
	pHResizeCursor   *glfw.Cursor
	pArrowCursor     *glfw.Cursor
	pHandCursor      *glfw.Cursor
	pCrosshairCursor *glfw.Cursor
	pIBeamCursor     *glfw.Cursor
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

const (
	ArrowCursor     = int(glfw.ArrowCursor)
	IBeamCursor     = int(glfw.IBeamCursor)
	CrosshairCursor = int(glfw.CrosshairCursor)
	HandCursor      = int(glfw.HandCursor)
	HResizeCursor   = int(glfw.HResizeCursor)
	VResizeCursor   = int(glfw.VResizeCursor)
)

var (
	LastRune  rune
	LastKey   glfw.Key
	LastMods  glfw.ModifierKey
	NoScaling bool
)

type Cursor glfw.Cursor

func SetCursor(wno int, c int) {
	gpu.Info[wno].Cursor = c
}

func Invalidate(w *glfw.Window) {
	wno := GetWno(w)
	n := gpu.Info[wno].InvalidateCount.Load()
	gpu.Info[wno].InvalidateCount.Store(n + 1)
}

func gotInvalidate() bool {
	for _, info := range gpu.Info {
		if info.InvalidateCount.Load() != 0 {
			n := info.InvalidateCount.Load()
			info.InvalidateCount.Store(n + 1)
			return true
		}
	}
	return false
}

func PollEvents() {
	t := time.Now()
	ClearMouseBtns()
	// Tight loop, waiting for events, checking for events every minDelay
	// Break anyway if waiting more than MaxDelay
	for !gotInvalidate() && time.Since(t) < MaxDelay {
		time.Sleep(minDelay)
	}
	glfw.PollEvents()
}

func Shutdown() {
	// glfw.DetachCurrentContext()
	for _, win := range WindowList {
		win.Destroy()
	}
	glfw.Terminate()
	TerminateProfiling()
}

// Init will initialize the system.
// The pallete is set to the default values
// The GLFW hints are set to the default values
// The connected monitors are put into the Monitors slice.
// Monitor info is printed to slog.
func Init() {
	runtime.LockOSThread()
	if *maxFps {
		MaxDelay = 0
	} else {
		MaxDelay = time.Second
	}
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	theme.SetDefaultPallete(true)
	SetDefaultHints()
	// Check all monitors and print size data
	Monitors = GetMonitors()
	// Select monitor as given, or use primary monitor.
	for i, m := range Monitors {
		SizeMmX, SizeMmY := m.GetPhysicalSize()
		mScaleX, mScaleY := m.GetContentScale()
		PosX, PosY, SizePxX, SizePxY := m.GetWorkarea()
		slog.Info("GetMonitors() for ", "Monitor", i+1,
			"WidthMm", SizeMmX, "HeightMm", SizeMmY,
			"WidthPx", SizePxX, "HeightPx", SizePxY, "PosX", PosX, "PosY", PosY,
			"ScaleX", f32.F2S(mScaleX, 3), "ScaleY", f32.F2S(mScaleY, 3))
	}
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	wno := GetWno(w)
	if wno < len(gpu.Info) {
		gpu.Info[wno].Focused = focused
		if !focused {
			resetFocus()
		}
		ClearMouseBtns()
		Invalidate(nil)
	}
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
	Invalidate(w)
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
	Invalidate(nil)
	LastRune = char
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	Invalidate(nil)
	LastMods = mods
	x, y := w.GetCursorPos()
	wno := GetWno(w)
	mousePos.X = float32(x) / gpu.Info[wno].ScaleX
	mousePos.Y = float32(y) / gpu.Info[wno].ScaleY
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
func posCallback(w *glfw.Window, xpos float64, ypos float64) {
	wno := GetWno(w)
	mousePos.X = float32(xpos) / gpu.Info[wno].ScaleX
	mousePos.Y = float32(ypos) / gpu.Info[wno].ScaleY
	Invalidate(w)
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("Scroll", "dx", xoff, "dy", yOff)
	if LastMods == glfw.ModControl {
		// ctrl+scrollwheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			gpu.CurrentInfo.UserScale *= ZoomFactor
		} else {
			gpu.CurrentInfo.UserScale /= ZoomFactor
		}
		UpdateSize(GetWno(w))
	} else {
		scrolledY = float32(yOff)
	}
	Invalidate(nil)
}

func UpdateSize(wno int) {
	width, height := WindowList[wno].GetSize()
	if NoScaling {
		gpu.Info[wno].ScaleX = 1.0
		gpu.Info[wno].ScaleY = 1.0
	} else {
		gpu.Info[wno].ScaleX, gpu.Info[wno].ScaleY = WindowList[wno].GetContentScale()
		gpu.Info[wno].ScaleX *= gpu.Info[wno].UserScale
		gpu.Info[wno].ScaleY *= gpu.Info[wno].UserScale
	}
	gpu.Info[wno].WindowOuterRectPx = gpu.IntRect{0, 0, width, height}
	gpu.Info[wno].WindowContentRectDp = f32.Rect{
		W: float32(width) / gpu.Info[wno].ScaleX,
		H: float32(height) / gpu.Info[wno].ScaleY}
	Invalidate(WindowList[wno])
	slog.Info("UpdateSize", "wno", wno, "w", width, "h", height, "scaleX", f32.F2S(gpu.Info[wno].ScaleX, 3),
		"ScaleY", f32.F2S(gpu.Info[wno].ScaleY, 3), "UserScale", f32.F2S(gpu.Info[wno].UserScale, 3))
}

func GetWno(w *glfw.Window) int {
	if w == nil {
		w = CurrentWindow
	}
	for i, _ := range WindowList {
		if WindowList[i] == w {
			return i
		}
	}
	return 0
}

func sizeCallback(w *glfw.Window, width int, height int) {
	wno := GetWno(w)
	UpdateSize(wno)
	gpu.UpdateResolution(wno)
	Invalidate(nil)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	width, height := w.GetSize()
	sizeCallback(w, width, height)
}

func SetDefaultHints() {
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

func createInvisibleWindow(w, h int, title string, monitor *glfw.Monitor) *glfw.Window {
	// Create invisible window so we can move it to correct monitor
	glfw.WindowHint(glfw.Visible, glfw.False)
	win, err := glfw.CreateWindow(w, h, title, monitor, nil)
	if err != nil || win == nil {
		panic(err)
	}
	return win
}

func SetupCursors() {
	pArrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	pVResizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)
	pHResizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pIBeamCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pCrosshairCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pHandCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
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

func MakeContextCurrent(wno int) {
	gpu.CurrentInfo = gpu.Info[wno]
	CurrentWindow = WindowList[wno]
	CurrentWno = wno
	WindowList[wno].MakeContextCurrent()
	gpu.UpdateResolution(wno)
	glfw.SwapInterval(1)
}

func init() {
	flag.Parse()
	slog.SetLogLoggerLevel(slog.Level(*logLevel))
	InitializeProfiling()
	buildinfo.Get()
}
