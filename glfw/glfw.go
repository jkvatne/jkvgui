package glfw

import "C"

const OpenGLProfile = 0x00022008
const OpenGLCoreProfile = 0x00032001
const OpenGLForwardCompatible = 0x00032002
const True = 1
const False = 0
const Resizable = 0x00020003
const Focused = 0x00020001
const Iconified = 0x00020002
const Resizeable = 0x00020003
const Visible = 0x00020004
const Decorated = 0x00020005
const AutoIconify = 0x00020006
const Floating = 0x00020007
const Maximized = 0x00020008
const ContextVersionMajor = 0x00022002
const ContextVersionMinor = 0x00022003
const Samples = 0x0002100D
const ArrowCursor = 0x00036001
const IbeamCursor = 0x00036002
const CrosshairCursor = 0x00036003
const HandCursor = 0x00036004
const HResizeCursor = 0x00036005
const VResizeCursor = 0x00036006

// Cursor mode values.
const (
	CursorNormal   int = 0x00034001
	CursorHidden   int = 0x00034002
	CursorDisabled int = 0x00034003
)

type Action int

type StandardCursor uint32

type Hint uint32
type Monitor struct {
}

type Window struct {
}
type Cursor struct {
	// data *C.GLFWcursor
	data int
}

func panicError() {
	// err := acceptError()
	// if err != nil {
	//	panic(err)
	// }
}

func WindowHint(target Hint, hint int) {
	// C.glfwWindowHint(C.int(target), C.int(hint))
	panicError()
}

// GetClipboardString returns the contents of the system clipboard, if it
// contains or is convertible to a UTF-8 encoded string.
// This function may only be called from the main thread.
func GetClipboardString() string {
	// cs := C.glfwGetClipboardString(nil)
	// if cs == nil {
	//	acceptError(FormatUnavailable)
	//	return ""
	// }
	return "" // C.GoString(cs)
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
// This function may only be called from the main thread.
func SetClipboardString(str string) {
	// cp := C.CString(str)
	// defer C.free(unsafe.Pointer(cp))
	// C.glfwSetClipboardString(nil, cp)
	panicError()
}

// CreateStandardCursor returns a cursor with a standard shape,
// that can be set for a window with SetCursor.
func CreateStandardCursor(shape StandardCursor) *Cursor {
	// c := C.glfwCreateStandardCursor(C.int(shape))
	panicError()
	return nil // &Cursor{c}
}

// CreateWindow creates a window and its associated context. Most of the options
// controlling how the window and its context should be created are specified
// through Hint.
//
// Successful creation does not change which context is current. Before you can
// use the newly created context, you need to make it current using
// MakeContextCurrent.
//
// Note that the created window and context may differ from what you requested,
// as not all parameters and hints are hard constraints. This includes the size
// of the window, especially for full screen windows. To retrieve the actual
// attributes of the created window and context, use queries like
// Window.GetAttrib and Window.GetSize.
//
// To create the window at a specific position, make it initially invisible using
// the Visible window hint, set its position and then show it.
//
// If a fullscreen window is active, the screensaver is prohibited from starting.
//
// Windows: If the executable has an icon resource named ICON, it will be
// set as the icon for the window. If no such icon is present, the IDI_WINLOGO
// icon will be used instead.
//
// Mac OS X: The GLFW window has no icon, as it is not a document window, but the
// dock icon will be the same as the application bundle's icon. Also, the first
// time a window is opened the menu bar is populated with common commands like
// Hide, Quit and About. The (minimal) about dialog uses information from the
// application's bundle. For more information on bundles, see the Bundle
// Programming Guide provided by Apple.
//
// This function may only be called from the main thread.
func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	/*	var (
			m *C.GLFWmonitor
			s *C.GLFWwindow
		)
		t := C.CString(title)
		defer C.free(unsafe.Pointer(t))
		if monitor != nil {
			m = monitor.data
		}
		if share != nil {
			s = share.data
		}
		// w := C.glfwCreateWindow(C.int(width), C.int(height), t, m, s)
		if w == nil {
			return nil, acceptError(APIUnavailable, VersionUnavailable)
		}
		wnd := &Window{data: w}
		windows.put(wnd)
	*/
	return nil, nil
}

// SwapBuffers swaps the front and back buffers of the window. If the
// swap interval is greater than zero, the GPU driver waits the specified number
// of screen updates before swapping the buffers.
func (w *Window) SwapBuffers() {
	// C.glfwSwapBuffers(w.data)
	panicError()
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
// of the specified window. The set cursor will only be visible when the cursor mode of the
// window is CursorNormal.
//
// On some platforms, the set cursor may not be visible unless the window also has input focus.
func (w *Window) SetCursor(c *Cursor) {
	if c == nil {
		// C.glfwSetCursor(w.data, nil)
	} else {
		// C.glfwSetCursor(w.data, c.data)
	}
	panicError()
}

// SetPos sets the position, in screen coordinates, of the upper-left corner
// of the client area of the window.
// If it is a full screen window, this function does nothing.
// If you wish to set an initial window position you should create a hidden
// window (using Hint and Visible), set its position and then show it.
// This function may only be called from the main thread.
func (w *Window) SetPos(xpos, ypos int) {
	// C.glfwSetWindowPos(w.data, C.int(xpos), C.int(ypos))
	panicError()
}

// SetSize sets the size, in screen coordinates, of the client area of the window.
// For full screen windows, this function selects and switches to the resolution
// closest to the specified size, without affecting the window's context. As the
// context is unaffected, the bit depths of the framebuffer remain unchanged.
// This function may only be called from the main thread.
func (w *Window) SetSize(width, height int) {
	// C.glfwSetWindowSize(w.data, C.int(width), C.int(height))
	panicError()
}

// Show makes the window visible, if it was previously hidden. If the window is
// already visible or is in full screen mode, this function does nothing.
//
// This function may only be called from the main thread.
func (w *Window) Show() {
	// C.glfwShowWindow(w.data)
	panicError()
}

// MakeContextCurrent makes the context of the window current.
// Originally GLFW 3 passes a null pointer to detach the context.
// But since we're using receievers, DetachCurrentContext should
// be used instead.
func (w *Window) MakeContextCurrent() {
	// C.glfwMakeContextCurrent(w.data)
	panicError()
}

// Focus brings the specified window to front and sets input focus.
// The window should already be visible and not iconified.
//
// By default, both windowed and full screen mode windows are focused when initially created.
// Set the glfw.Focused to disable this behavior.
//
// Do not use this function to steal focus from other applications unless you are certain that
// is what the user wants. Focus stealing can be extremely disruptive.
func (w *Window) Focus() {
	// C.glfwFocusWindow(w.data)
}

// ShouldClose reports the value of the close flag of the specified window.
func (w *Window) ShouldClose() bool {
	// ret := glfwbool(C.glfwWindowShouldClose(w.data))
	panicError()
	return true
}

// CursorPosCallback the cursor position callback.
type CursorPosCallback func(w *Window, xpos float64, ypos float64)

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the window.
func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	/*previous = w.fCursorPosHolder
	w.fCursorPosHolder = cbfun
	if cbfun == nil {
		// C.glfwSetCursorPosCallback(w.data, nil)
	} else {
		// C.glfwSetCursorPosCallbackCB(w.data)
	}
	panicError()
	return previous
	*/
	return nil
}

// KeyCallback is the key callback.
type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

// SetKeyCallback sets the key callback which is called when a key is pressed,
// repeated or released.
//
// The key functions deal with physical keys, with layout independent key tokens
// named after their values in the standard US keyboard layout. If you want to
// input text, use the SetCharCallback instead.
//
// When a window loses focus, it will generate synthetic key release events for
// all pressed keys. You can tell these events from user-generated events by the
// fact that the synthetic ones are generated after the window has lost focus,
// i.e. Focused will be false and the focus callback will have already been
// called.
func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	/*previous = w.fKeyHolder
	w.fKeyHolder = cbfun
	if cbfun == nil {
		// C.glfwSetKeyCallback(w.data, nil)
	} else {
		// C.glfwSetKeyCallbackCB(w.data)
	}
	panicError()
	return previous
	*/
	return nil
}

// CharCallback is the character callback.
type CharCallback func(w *Window, char rune)

// SetCharCallback sets the character callback which is called when a
// Unicode character is input.
//
// The character callback is intended for Unicode text input. As it deals with
// characters, it is keyboard layout dependent, whereas the
// key callback is not. Characters do not map 1:1
// to physical keys, as a key may produce zero, one or more characters. If you
// want to know whether a specific physical key was pressed or released, see
// the key callback instead.
//
// The character callback behaves as system text input normally does and will
// not be called if modifier keys are held down that would prevent normal text
// input on that platform, for example a Super (Command) key on OS X or Alt key
// on Windows. There is a character with modifiers callback that receives these events.
func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	/*previous = w.fCharHolder
	w.fCharHolder = cbfun
	if cbfun == nil {
		// C.glfwSetCharCallback(w.data, nil)
	} else {
		// C.glfwSetCharCallbackCB(w.data)
	}
	panicError()
	return previous */
	return nil
}

// DropCallback is the drop callback.
type DropCallback func(w *Window, names []string)

// SetDropCallback sets the drop callback which is called when an object
// is dropped over the window.
func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	/*previous = w.fDropHolder
	w.fDropHolder = cbfun
	if cbfun == nil {
		// C.glfwSetDropCallback(w.data, nil)
	} else {
		// C.glfwSetDropCallbackCB(w.data)
	}
	panicError()
	return previous
	*/
	return nil
}

// ContentScaleCallback is the function signature for window content scale
// callback functions.
type ContentScaleCallback func(w *Window, x float32, y float32)

// SetContentScaleCallback function sets the window content scale callback of
// the specified window, which is called when the content scale of the specified
// window changes.
//
// This function must only be called from the main thread.
func (w *Window) SetContentScaleCallback(cbfun ContentScaleCallback) ContentScaleCallback {
	/*previous := w.fContentScaleHolder
	w.fContentScaleHolder = cbfun
	if cbfun == nil {
		// C.glfwSetWindowContentScaleCallback(w.data, nil)
	} else {
		// C.glfwSetWindowContentScaleCallbackCB(w.data)
	}
	return previous */
	return nil
}

// RefreshCallback is the window refresh callback.
type RefreshCallback func(w *Window)

// SetRefreshCallback sets the refresh callback of the window, which
// is called when the client area of the window needs to be redrawn, for example
// if the window has been exposed after having been covered by another window.
//
// On compositing window systems such as Aero, Compiz or Aqua, where the window
// contents are saved off-screen, this callback may be called only very
// infrequently or never at all.
func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	/*previous = w.fRefreshHolder
	w.fRefreshHolder = cbfun
	if cbfun == nil {
		// C.glfwSetWindowRefreshCallback(w.data, nil)
	} else {
		// C.glfwSetWindowRefreshCallbackCB(w.data)
	}
	panicError()
	return previous*/
	return nil
}

// FocusCallback is the window focus callback.
type FocusCallback func(w *Window, focused bool)

// SetFocusCallback sets the focus callback of the window, which is called when
// the window gains or loses focus.
//
// After the focus callback is called for a window that lost focus, synthetic key
// and mouse button release events will be generated for all such that had been
// pressed. For more information, see SetKeyCallback and SetMouseButtonCallback.
func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	/*previous = w.fFocusHolder
	w.fFocusHolder = cbfun
	if cbfun == nil {
		// C.glfwSetWindowFocusCallback(w.data, nil)
	} else {
		// C.glfwSetWindowFocusCallbackCB(w.data)
	}
	panicError()
	return previous*/
	return nil
}

// SizeCallback is the window size callback.
type SizeCallback func(w *Window, width int, height int)

// SetSizeCallback sets the size callback of the window, which is called when
// the window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the window.
func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	/*previous = w.fSizeHolder
	w.fSizeHolder = cbfun
	if cbfun == nil {
		// C.glfwSetWindowSizeCallback(w.data, nil)
	} else {
		// C.glfwSetWindowSizeCallbackCB(w.data)
	}
	panicError()
	return previous	 */
	return nil
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the window and input
// callbacks associated with those events to be called.
//
// This function is not required for joystick input to work.
//
// This function may not be called from a callback.
//
// This function may only be called from the main thread.
func PollEvents() {
	// C.glfwPollEvents()
	panicError()
}

// Terminate destroys all remaining windows, frees any allocated resources and
// sets the library to an uninitialized state. Once this is called, you must
// again call Init successfully before you will be able to use most GLFW
// functions.
//
// If GLFW has been successfully initialized, this function should be called
// before the program exits. If initialization fails, there is no need to call
// this function, as it is called by Init before it returns failure.
//
// This function may only be called from the main thread.
func Terminate() {
	// flushErrors()
	// C.glfwTerminate()
}

func Init() error {
	return nil
}

// GetContentScale function retrieves the content scale for the specified
// window. The content scale is the ratio between the current DPI and the
// platform's default DPI. If you scale all pixel dimensions by this scale then
// your content should appear at an appropriate size. This is especially
// important for text and any UI elements.
//
// This function may only be called from the main thread.
func (w *Window) GetContentScale() (float32, float32) {
	var x, y C.float
	// C.glfwGetWindowContentScale(w.data, &x, &y)
	return float32(x), float32(y)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// of the specified window. This size includes the title bar, if the window has one.
// The size of the frame may vary depending on the window-related hints used to create it.
//
// Because this function retrieves the size of each window frame edge and not the offset
// along a particular coordinate axis, the retrieved values will always be zero or positive.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b C.int
	// C.glfwGetWindowFrameSize(w.data, &l, &t, &r, &b)
	panicError()
	return int(l), int(t), int(r), int(b)
}

// SwapInterval sets the swap interval for the current context, i.e. the number
// of screen updates to wait before swapping the buffers of a window and
// returning from SwapBuffers. This is sometimes called
// 'vertical synchronization', 'vertical retrace synchronization' or 'vsync'.
//
// Contexts that support either of the WGL_EXT_swap_control_tear and
// GLX_EXT_swap_control_tear extensions also accept negative swap intervals,
// which allow the driver to swap even if a frame arrives a little bit late.
// You can check for the presence of these extensions using
// ExtensionSupported. For more information about swap tearing,
// see the extension specifications.
//
// Some GPU drivers do not honor the requested swap interval, either because of
// user settings that override the request or due to bugs in the driver.
func SwapInterval(interval int) {
	// C.glfwSwapInterval(C.int(interval))
	panicError()
}

// GetCursorPos returns the last reported position of the cursor.
//
// If the cursor is disabled (with CursorDisabled) then the cursor position is
// unbounded and limited only by the minimum and maximum values of a double.
//
// The coordinate can be converted to their integer equivalents with the floor
// function. Casting directly to an integer type works for positive coordinates,
// but fails for negative ones.
func (w *Window) GetCursorPos() (x, y float64) {
	var xpos, ypos C.double
	// C.glfwGetCursorPos(w.data, &xpos, &ypos)
	panicError()
	return float64(xpos), float64(ypos)
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified window.
func (w *Window) GetSize() (width, height int) {
	var wi, h C.int
	// C.glfwGetWindowSize(w.data, &wi, &h)
	panicError()
	return int(wi), int(h)
}
