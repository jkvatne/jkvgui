// Package sys is the only package that depends on glfw.
package sys

import (
	"flag"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"

	// Using my own purego-glfw implementation:
	glfw "github.com/jkvatne/purego-glfw"
	// Using standard go-gl from GitHub:
	// "github.com/go-gl/glfw/v3.3/glfw"
)

var (
	Monitors      []*glfw.Monitor
	maxFps        = flag.Int("maxfps", 60, "Set to maximum allowed frames pr second. Default to 60")
	NoScaling     bool
	WindowList    []*Window
	WindowCount   atomic.Int32
	WinListMutex  sync.RWMutex
	MinFrameDelay time.Duration
	MaxFrameDelay = time.Second * 5
	OpenGlStarted bool
)

type HintDef struct {
	WidgetRect f32.Rect // Original widgets size
	Text       string
	T          time.Time
	Tag        any
}

// Window variables.
type Window struct {
	Window               *glfw.Window
	Name                 string
	Wno                  int
	UserScale            float32
	Mutex                sync.Mutex
	Trigger              chan bool
	HintActive           bool
	Focused              bool
	Blinking             atomic.Bool
	Cursor               int
	CurrentTag           interface{}
	LastTag              interface{}
	MoveToNext           bool
	MoveToPrevious       bool
	ToNext               bool
	SuppressEvents       bool
	mousePos             f32.Pos
	LeftBtnIsDown        bool
	Dragging             bool
	DragStartPos         f32.Pos
	LeftBtnDownTime      time.Time
	LeftBtnUpTime        time.Time
	LeftBtnDoubleClicked bool
	LeftBtnClicked       bool
	ScrolledDistY        float32
	DialogVisible        bool
	redraws              int
	fps                  float64
	redrawStart          time.Time
	LastRune             rune
	LastKey              glfw.Key
	LastMods             glfw.ModifierKey
	NoScaling            bool
	CurrentHint          HintDef
	DeferredFunctions    []func()
	HeightPx             int
	HeightDp             float32
	WidthPx              int
	WidthDp              float32
	Gd                   gpu.GlData
}

var (
	pVResizeCursor   *glfw.Cursor
	pHResizeCursor   *glfw.Cursor
	pArrowCursor     *glfw.Cursor
	pHandCursor      *glfw.Cursor
	pCrosshairCursor *glfw.Cursor
	pIBeamCursor     *glfw.Cursor
)

//goland:noinspection ALL,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst
const (
	KeyRight     = glfw.KeyRight
	KeyLeft      = glfw.KeyLeft
	KeyUp        = glfw.KeyUp
	KeyDown      = glfw.KeyDown
	KeySpace     = glfw.KeySpace
	KeyEnter     = glfw.KeyEnter
	KeyKPEnter   = glfw.KeyKPEnter
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

type Cursor glfw.Cursor

func (w *Window) Destroy() {
	w.Window.Destroy()
}

func (w *Window) Invalidate() {
	glfw.PostEmptyEvent()
}

func (w *Window) PollEvents() {
	PollEvents()
}

func PollEvents() {
	t := time.Now()
	glfw.WaitEventsTimeout(float64(MaxFrameDelay) / 1e9)
	if time.Since(t) < MinFrameDelay {
		time.Sleep(MinFrameDelay - time.Since(t))
	}
}

func Shutdown() {
	WinListMutex.Lock()
	defer WinListMutex.Unlock()
	for _, win := range WindowList {
		win.Window.Destroy()
	}
	WindowList = nil
	WindowCount.Store(0)
	glfw.Terminate()
	TerminateProfiling()
	OpenGlStarted = false
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	win := GetWindow(w)
	if win == nil {
		slog.Error("Focus callback without any window")
		return
	}
	win.Focused = focused
	if !focused {
		slog.Debug("Lost focus", "Wno ", win.Wno+1)
		win.ClearMouseBtns()
	} else {
		slog.Debug("Got focus", "Wno", win.Wno+1)
	}
	win.Invalidate()
}

func setCallbacks(Window *glfw.Window) {
	Window.SetMouseButtonCallback(btnCallback)
	Window.SetCursorPosCallback(posCallback)
	Window.SetKeyCallback(keyCallback)
	Window.SetCharCallback(charCallback)
	Window.SetScrollCallback(scrollCallback)
	Window.SetContentScaleCallback(scaleCallback)
	Window.SetFocusCallback(focusCallback)
	Window.SetCloseCallback(closeCallback)
	Window.SetSizeCallback(sizeCallback)
}

func closeCallback(w *glfw.Window) {
	slog.Debug("Close callback", "ShouldClose", w.ShouldClose())
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	win := GetWindow(w)
	if win == nil {
		slog.Error("Key callback without any window")
		return
	}
	win.Invalidate()
	if key == glfw.KeyTab && action == glfw.Release {
		win.MoveByKey(mods != glfw.ModShift)
	}
	if action == glfw.Release || action == glfw.Repeat {
		win.LastKey = key
	}
	win.LastMods = mods
}

func charCallback(w *glfw.Window, char rune) {
	slog.Debug("charCallback()", "Rune", int(char))
	win := GetWindow(w)
	if win == nil {
		slog.Error("Char callback without any window")
		return
	}
	win.Invalidate()
	win.LastRune = char
}

// ClearMouseBtns is called when a window looses focus. It will reset the mouse button states.
func (w *Window) ClearMouseBtns() {
	w.LeftBtnIsDown = false
	w.Dragging = false
	w.LeftBtnDoubleClicked = false
	w.LeftBtnClicked = false
	w.ScrolledDistY = 0.0
	w.LeftBtnUpTime = time.Time{}
}

func (w *Window) leftBtnRelease() {
	w.LeftBtnIsDown = false
	w.Dragging = false
	if time.Since(w.LeftBtnUpTime) < DoubleClickTime {
		slog.Debug("MouseCb: - DoubleClick:")
		w.LeftBtnDoubleClicked = true
	} else {
		slog.Debug("MouseCb: - Click:")
		w.LeftBtnClicked = true
	}
	w.LeftBtnUpTime = time.Now()
}

func (w *Window) leftBtnPress() {
	w.LeftBtnIsDown = true
	w.LeftBtnClicked = false
	w.LeftBtnDownTime = time.Now()
}

func (w *Window) SimPos(x, y float32) {
	w.mousePos.X = x
	w.mousePos.Y = y
}

func (w *Window) SimLeftBtnPress(x, y float32) {
	w.mousePos.X = x
	w.mousePos.Y = y
	w.leftBtnPress()
}

func (w *Window) SimLeftBtnRelease(x, y float32) {
	w.mousePos.X = x
	w.mousePos.Y = y
	w.leftBtnRelease()
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	win := GetWindow(w)
	if win == nil {
		panic("Window callback without any window")
	}
	win.Invalidate()
	win.LastMods = mods
	x, y := w.GetCursorPos()
	win.mousePos.X = float32(x) / win.Gd.ScaleX
	win.mousePos.Y = float32(y) / win.Gd.ScaleY
	slog.Debug("MouseCb:", "Button", button, "X", x, "Y", y, "Action", action, "FromWindow", win.Wno, "Pos", win.mousePos)
	if button == glfw.MouseButtonLeft {
		if action == glfw.Release {
			win.leftBtnRelease()
		} else if action == glfw.Press {
			win.leftBtnPress()
		}
	}
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(w *glfw.Window, xPos float64, yPos float64) {
	win := GetWindow(w)
	if win == nil {
		slog.Error("Mouse position callback without any window")
		return
	}
	win.mousePos.X = float32(xPos) / win.Gd.ScaleX
	win.mousePos.Y = float32(yPos) / win.Gd.ScaleY
	win.Invalidate()
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("ScrollCb:", "dx", xoff, "dy", yOff)
	win := GetWindow(w)
	if win == nil {
		slog.Error("Scroll callback without any window")
		return
	}

	if win.LastMods == glfw.ModControl {
		// ctrl + scroll-wheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			win.UserScale *= ZoomFactor
		} else {
			win.UserScale /= ZoomFactor
		}
		win.UpdateSizeDp()
	} else {
		win.ScrolledDistY = float32(yOff)
	}
	win.Invalidate()
}

func GetWindow(w *glfw.Window) *Window {
	WinListMutex.RLock()
	defer WinListMutex.RUnlock()
	for i := range WindowList {
		if WindowList[i].Window == w {
			return WindowList[i]
		}
	}
	return nil
}

func sizeCallback(w *glfw.Window, width int, height int) {
	slog.Debug("sizeCallback", "width", width, "height", height)
	win := GetWindow(w)
	if win == nil {
		slog.Error("Size callback without any window")
		return
	}

	win.UpdateSize(width, height)
	win.Invalidate()
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	slog.Debug("scaleCallback", "x", x, "y", y)
	win := GetWindow(w)
	if win == nil {
		slog.Error("Scale callback without any window")
		return
	}

	win.UpdateSizeDp()
	win.UpdateResolution()
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

// PostEmptyEvent will post an empty event to the thread that initialized glfw.
// This is normally the background thread runnin main().
func PostEmptyEvent() {
	glfw.PostEmptyEvent()
}

func glfwInit() error {
	return glfw.Init()
}

func DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func SwapInterval(n int) {
	glfw.SwapInterval(n)
}

func GetCurrentContext() *glfw.Window {
	return glfw.GetCurrentContext()
}
