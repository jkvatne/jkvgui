// sys is the only package that depends on glfw.
package sys

import (
	"flag"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jkvatne/jkvgui/f32"

	// Using my own purego-glfw implementation:
	glfw "github.com/jkvatne/purego-glfw"
	// Using standard go-gl from github:
	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/gpu"
)

var Monitors []*glfw.Monitor
var maxFps = flag.Int("maxfps", 60, "Set to maximum alowed frames pr second. Default to 60")

var (
	WindowList   []*Window
	WindowCount  atomic.Int32
	WinListMutex sync.Mutex
)

// Pr window global variables.
type Window struct {
	Window               *glfw.Window
	Name                 string
	Wno                  int
	UserScale            float32
	Mutex                sync.Mutex
	InvalidateCount      atomic.Int32
	Trigger              chan bool
	HintActive           bool
	Focused              bool
	BlinkState           atomic.Bool
	Blinking             atomic.Bool
	Cursor               int
	CurrentTag           interface{}
	LastTag              interface{}
	MoveToNext           bool
	MoveToPrevious       bool
	ToNext               bool
	SuppressEvents       bool
	MousePos             f32.Pos
	LeftBtnIsDown        bool
	LeftBtnReleased      bool
	Dragging             bool
	LeftBtnDownTime      time.Time
	LeftBtnUpTime        time.Time
	LeftBtnDoubleClicked bool
	ScrolledDistY        float32
	DialogVisible        bool
	redraws              int
	fps                  float64
	redrawStart          time.Time
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

func (w *Window) MakeContextCurrent() {
	w.Window.MakeContextCurrent()
	gpu.Gd = w.Gd
}

func (w *Window) SetCursor(c int) {
	w.Cursor = c
}

func (w *Window) Invalidate() {
	w.InvalidateCount.Add(1)
	glfw.PostEmptyMessage(w.Window)
	if len(w.Trigger) == 0 {
		w.Trigger <- true
	}
}

func (w *Window) Invalid() bool {
	if w.InvalidateCount.Load() != 0 {
		w.InvalidateCount.Store(0)
		return true
	}
	return false
}

func (w *Window) PollEvents() {
	w.ClearMouseBtns()
	// Tight loop, waiting for events, checking for events every minDelay
	// Break anyway if waiting more than MaxFrameDelay
	t := time.Now()
	for w.InvalidateCount.Load() == 0 && time.Since(t) < MaxFrameDelay {
		glfw.WaitEventsTimeout(float64(MaxFrameDelay) / 1e9)
	}
	if time.Since(t) < MinFrameDelay {
		time.Sleep(MinFrameDelay - time.Since(t))
	}
	w.InvalidateCount.Store(0)
	glfw.PollEvents()
}

func PollEvents() {
	glfw.WaitEventsTimeout(float64(MaxFrameDelay) / 1e9)
	glfw.PollEvents()
}

func Shutdown() {
	for _, win := range WindowList {
		win.Window.Destroy()
	}
	WindowList = WindowList[0:0]
	WindowList = WindowList[0:0]
	WindowCount.Store(0)
	glfw.Terminate()
	TerminateProfiling()
}

func GetMonitors() []*glfw.Monitor {
	return glfw.GetMonitors()
}

func focusCallback(w *glfw.Window, focused bool) {
	win := GetWindow(w)
	if win != nil {
		win.Focused = focused
		if !focused {
			slog.Info("Lost focus", "Wno ", win.Wno+1)
			win.ClearMouseBtns()
		} else {
			slog.Info("Got focus", "Wno", win.Wno+1)
		}
		win.Invalidate()
	} else {
		// slog.Info("Focus callback without any window", "Wno", win.Wno+1)
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
	Window.SetCloseCallback(closeCallback)
}

func closeCallback(w *glfw.Window) {
	// fmt.Printf("Close callback %v\n", w.ShouldClose())
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	slog.Debug("keyCallback", "key", key, "scancode", scancode, "action", action, "mods", mods)
	win := GetWindow(w)
	win.Invalidate()
	if key == glfw.KeyTab && action == glfw.Release {
		win.MoveByKey(mods != glfw.ModShift)
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
	win := GetWindow(w)
	win.Invalidate()
	LastRune = char
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	win := GetWindow(w)
	win.Invalidate()
	LastMods = mods
	x, y := w.GetCursorPos()
	win.MousePos.X = float32(x) / win.Gd.ScaleX
	win.MousePos.Y = float32(y) / win.Gd.ScaleY
	// wno := GetWindow(w)
	// slog.Info("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action, "FromWindow", wno)
	if button == glfw.MouseButtonLeft {
		if action == glfw.Release {
			win.LeftBtnIsDown = false
			win.LeftBtnReleased = true
			win.Dragging = false
			if time.Since(win.LeftBtnUpTime) < DoubleClickTime {
				win.LeftBtnDoubleClicked = true
			}
			win.LeftBtnUpTime = time.Now()
		} else if action == glfw.Press {
			win.LeftBtnIsDown = true
			win.LeftBtnDownTime = time.Now()
		}
	}
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(w *glfw.Window, xpos float64, ypos float64) {
	win := GetWindow(w)
	win.MousePos.X = float32(xpos) / win.Gd.ScaleX
	win.MousePos.Y = float32(ypos) / win.Gd.ScaleY
	win.Invalidate()
	// slog.Info("MouseMove callback", "wno", win.Wno, "InvalidateCount", win.InvalidateCount.Load())
}

func scrollCallback(w *glfw.Window, xoff float64, yOff float64) {
	slog.Debug("Scroll", "dx", xoff, "dy", yOff)
	win := GetWindow(w)
	if LastMods == glfw.ModControl {
		// ctrl+scrollwheel will zoom the whole window by changing gpu.UserScale.
		if yOff > 0 {
			win.UserScale *= ZoomFactor
		} else {
			win.UserScale /= ZoomFactor
		}
		win.UpdateSize()
	} else {
		win.ScrolledDistY = float32(yOff)
	}
	win.Invalidate()
}

func GetWindow(w *glfw.Window) *Window {
	for i, _ := range WindowList {
		if WindowList[i].Window == w {
			return WindowList[i]
		}
	}
	return nil
}

func sizeCallback(w *glfw.Window, width int, height int) {
	win := GetWindow(w)
	win.UpdateSize()
	gpu.UpdateResolution()
	win.Invalidate()
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

func glfwInit() error {
	return glfw.Init()
}

func DetachCurrentContext() {
	glfw.DetachCurrentContext()
}

func SwapInterval(n int) {
	glfw.SwapInterval(n)
}
