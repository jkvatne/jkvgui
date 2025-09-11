// sys is the only package that depends on glfw.
package sys

import (
	"flag"
	"log/slog"
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

var Monitors []*glfw.Monitor
var maxFps = flag.Bool("maxfps", false, "Set to force redrawing as fast as possible")

type Window glfw.Window

var (
	WindowList       []*glfw.Window
	WinListMutex     sync.Mutex
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
	WindowList[wno].Cursor = c
}

func Invalidate() {
	WindowList[CurrentWno].InvalidateCount.Add(1)
	// WindowList[CurrentWno].PostEmptyEvent()
}

func gotInvalidate() bool {
	for _, info := range WindowList {
		if info.InvalidateCount.Load() != 0 {
			info.InvalidateCount.Store(0)
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
		glfw.WaitEventsTimeout(float64(MaxDelay) / 1e9)
	}
	glfw.PollEvents()
}

func Shutdown() {
	for _, win := range WindowList {
		win.Destroy()
	}
	WindowList = WindowList[0:0]
	WindowList = WindowList[0:0]
	WindowCount.Store(0)
	glfw.Terminate()
	TerminateProfiling()
}

func glfwInit() error {
	return glfw.Init()
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	wno := GetWno(w)
	if wno < len(WindowList) {
		WindowList[wno].Focused = focused
		if !focused {
			slog.Info("Lost focus", "Wno ", wno+1)
			ClearMouseBtns()
		} else {
			slog.Info("Got focus", "Wno", wno+1)
		}
		Invalidate()
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
	Window.SetCloseCallback(closeCallback)
}

func closeCallback(w *glfw.Window) {
	// fmt.Printf("Close callback %v\n", w.ShouldClose())
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	Invalidate()
	if key == glfw.KeyTab && action == glfw.Release {
		MoveByKey(mods != glfw.ModShift)
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
	Invalidate()
	LastRune = char
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	Invalidate()
	LastMods = mods
	x, y := w.GetCursorPos()
	wno := GetWno(w)
	info := WindowList[wno]
	info.MousePos.X = float32(x) / WindowList[wno].ScaleX
	info.MousePos.Y = float32(y) / WindowList[wno].ScaleY
	slog.Info("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action, "FromWindow", wno)
	if button == glfw.MouseButtonLeft {
		if action == glfw.Release {
			info.LeftBtnDown = false
			info.LeftBtnReleased = true
			info.Dragging = false
			if time.Since(info.LeftBtnUpTime) < DoubleClickTime {
				info.LeftBtnDoubleClick = true
			}
			info.LeftBtnUpTime = time.Now()
		} else if action == glfw.Press {
			info.LeftBtnDown = true
			info.LeftBtnDownTime = time.Now()
		}
	}
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(w *glfw.Window, xpos float64, ypos float64) {
	info := WindowList[GetWno(w)]
	info.MousePos.X = float32(xpos) / info.ScaleX
	info.MousePos.Y = float32(ypos) / info.ScaleY
	Invalidate()
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("Scroll", "dx", xoff, "dy", yOff)
	info := WindowList[GetWno(w)]
	if LastMods == glfw.ModControl {
		// ctrl+scrollwheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			info.UserScale *= ZoomFactor
		} else {
			info.UserScale /= ZoomFactor
		}
		UpdateSize(w)
	} else {
		info.ScrolledY = float32(yOff)
	}
	Invalidate()
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
	UpdateSize(w)
	gpu.UpdateResolution()
	Invalidate()
	slog.Info("sizeCallback", "wno", wno, "w", width, "h", height, "scaleX", f32.F2S(WindowList[wno].ScaleX, 3),
		"ScaleY", f32.F2S(WindowList[wno].ScaleY, 3), "UserScale", f32.F2S(WindowList[wno].UserScale, 3))
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

func GetClipboardString() (string, error) {
	return glfw.GetClipboardString(), nil
}

func MaximizeWindow(w *glfw.Window) {
	w.Maximize()
}

func MinimizeWindow(w *glfw.Window) {
	w.Iconify()
}

func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func UpdateSize(w *glfw.Window) {
	wno := GetWno(w)
	width, height := w.GetSize()
	if NoScaling {
		WindowList[wno].ScaleX = 1.0
		WindowList[wno].ScaleY = 1.0
	} else {
		WindowList[wno].ScaleX, WindowList[wno].ScaleY = WindowList[wno].GetContentScale()
		WindowList[wno].ScaleX *= WindowList[wno].UserScale
		WindowList[wno].ScaleY *= WindowList[wno].UserScale
	}
	gpu.ClientRectPx = gpu.IntRect{0, 0, width, height}
	gpu.ClientRectDp = f32.Rect{
		W: float32(width) / WindowList[wno].ScaleX,
		H: float32(height) / WindowList[wno].ScaleY}
	gpu.ScaleX, gpu.ScaleY = WindowList[wno].ScaleX, WindowList[wno].ScaleY
	gpu.UpdateResolution()
}
