// sys is the only package that depends on glfw.
// glfw is only imported in glfw_linux.go or glfw_windows.go
// Except for the imports, these files should be identical
// Use "github.com/go-gl/glfw/v3.3/glfw"

package sys

import (
	"log/slog"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
)

var (
	Monitors []*glfw.Monitor
)

type HintDef struct {
	WidgetRect f32.Rect // Original widgets size
	Text       string
	T          time.Time
	Tag        any
}

type (
	GlfwWindow  = glfw.Window
	Key         = glfw.Key
	Action      = glfw.Action
	ModifierKey = glfw.ModifierKey
	MouseButton = glfw.MouseButton
)

var (
	pVResizeCursor   *glfw.Cursor
	pHResizeCursor   *glfw.Cursor
	pArrowCursor     *glfw.Cursor
	pHandCursor      *glfw.Cursor
	pCrosshairCursor *glfw.Cursor
	pIBeamCursor     *glfw.Cursor
)

//goland:noinspection ALL,GoUnusedConst
const (
	KeyRight         = glfw.KeyRight
	KeyLeft          = glfw.KeyLeft
	KeyUp            = glfw.KeyUp
	KeyDown          = glfw.KeyDown
	KeyTab           = glfw.KeyTab
	KeySpace         = glfw.KeySpace
	KeyEnter         = glfw.KeyEnter
	KeyKPEnter       = glfw.KeyKPEnter
	KeyEscape        = glfw.KeyEscape
	KeyBackspace     = glfw.KeyBackspace
	KeyDelete        = glfw.KeyDelete
	KeyHome          = glfw.KeyHome
	KeyEnd           = glfw.KeyEnd
	KeyPageUp        = glfw.KeyPageUp
	KeyPageDown      = glfw.KeyPageDown
	KeyInsert        = glfw.KeyInsert
	KeyC             = glfw.KeyC
	KeyV             = glfw.KeyV
	KeyX             = glfw.KeyX
	ModShift         = glfw.ModShift
	ModControl       = glfw.ModControl
	ModAlt           = glfw.ModAlt
	Release          = glfw.Release
	Press            = glfw.Press
	Repeat           = glfw.Repeat
	MouseButtonLeft  = glfw.MouseButtonLeft
	MouseButtonRight = glfw.MouseButtonRight
)

const (
	ArrowCursor     = int(glfw.ArrowCursor)
	IBeamCursor     = int(glfw.IBeamCursor)
	CrosshairCursor = int(glfw.CrosshairCursor)
	HandCursor      = int(glfw.HandCursor)
	HResizeCursor   = int(glfw.HResizeCursor)
	VResizeCursor   = int(glfw.VResizeCursor)
)

func WaitEventsTimeout(secondsDelay float32) {
	glfw.WaitEventsTimeout(float64(secondsDelay))
}

func Terminate() {
	glfw.Terminate()
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
	win.HandleFocus(focused)
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

func closeCallback(w *glfw.Window) {
	slog.Debug("Close callback", "ShouldClose", w.ShouldClose())
}

// keyCallback see https://www.glfw.org/docs/latest/window_guide.html
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	GetWindow(w).HandleKey(key, scancode, action, mods)
}

func charCallback(w *glfw.Window, char rune) {
	GetWindow(w).HandleChar(char)
}

// btnCallback is called from the glfw window handler when mouse buttons change states.
func btnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	GetWindow(w).HandleMouseButton(button, action, mods)
}

// posCallback is called from the glfw window handler when the mouse moves.
func posCallback(w *glfw.Window, xPos float64, yPos float64) {
	GetWindow(w).HandleMousePos(xPos, yPos)
}

func scrollCallback(w *glfw.Window, xOff float64, yOff float64) {
	GetWindow(w).HandleMouseScroll(xOff, yOff)
}

func sizeCallback(w *glfw.Window, width int, height int) {
	slog.Debug("sizeCallback", "width", width, "height", height)
	GetWindow(w).UpdateSize(width, height)
}

func scaleCallback(w *glfw.Window, x float32, y float32) {
	slog.Debug("scaleCallback", "x", x, "y", y)
	GetWindow(w).UpdateSizeDp()
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
// This is normally the original thread running main().
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

func setupCursors() {
	pArrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	pVResizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)
	pHResizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pIBeamCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pCrosshairCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	pHandCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
}
