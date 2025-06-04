package glfw

import "C"
import (
	"fmt"
	"syscall"
)

const (
	GLFW_DONT_CARE          = -1
	OpenGLProfile           = 0x00022008
	OpenGLCoreProfile       = 0x00032001
	OpenGLForwardCompatible = 0x00032002
	True                    = 1
	False                   = 0
	Resizable               = 0x00020003
	Focused                 = 0x00020001
	Iconified               = 0x00020002
	Resizeable              = 0x00020003
	Visible                 = 0x00020004
	Decorated               = 0x00020005
	AutoIconify             = 0x00020006
	Floating                = 0x00020007
	Maximized               = 0x00020008
	ContextVersionMajor     = 0x00022002
	ContextVersionMinor     = 0x00022003
	Samples                 = 0x0002100D
	ArrowCursor             = 0x00036001
	IbeamCursor             = 0x00036002
	CrosshairCursor         = 0x00036003
	HandCursor              = 0x00036004
	HResizeCursor           = 0x00036005
	VResizeCursor           = 0x00036006
	LR_CREATEDIBSECTION     = 0x00002000
	LR_DEFAULTCOLOR         = 0x00000000
	LR_DEFAULTSIZE          = 0x00000040
	LR_LOADFROMFILE         = 0x00000010
	LR_LOADMAP3DCOLORS      = 0x00001000
	LR_LOADTRANSPARENT      = 0x00000020
	LR_MONOCHROME           = 0x00000001
	LR_SHARED               = 0x00008000
	LR_VGACOLOR             = 0x00000080
	IMAGE_ICON              = 1
	CS_HREDRAW              = 0x0002
	CS_INSERTCHAR           = 0x2000
	CS_NOMOVECARET          = 0x4000
	CS_VREDRAW              = 0x0001
	CS_OWNDC                = 0x0020
	KF_EXTENDED             = 0x100
	GLFW_RELEASE            = 0
	GLFW_PRESS              = 1
	GLFW_REPEAT             = 2
	GLFW_CURSOR_NORMAL      = 0x00034001
	GLFW_CURSOR_HIDDEN      = 0x00034002
	GLFW_CURSOR_DISABLED    = 0x00034003
)

type Action int

type StandardCursor uint32

type Hint uint32

//
type GLFWvidmode struct {
	width       int
	height      int
	redBits     int
	greenBits   int
	blueBits    int
	refreshRate int
}

// Window represents a Window.
type Window struct {
	data                 *_GLFWwindow
	charCallback         CharCallback
	focusCallback        FocusCallback
	keyCallback          KeyCallback
	mouseButtonCallback  MouseButtonCallback
	cursorPosCallback    CursorPosCallback
	scrollCallback       ScrollCallback
	refreshCallback      RefreshCallback
	sizeCallback         SizeCallback
	dropCallback         DropCallback
	contentScaleCallback func(w *Window, x float32, y float32)

	fPosHolder             func(w *Window, xpos int, ypos int)
	fSizeHolder            func(w *Window, width int, height int)
	fFramebufferSizeHolder func(w *Window, width int, height int)
	fCloseHolder           func(w *Window)
	fMaximizeHolder        func(w *Window, maximized bool)
	fContentScaleHolder    func(w *Window, x float32, y float32)
	fRefreshHolder         func(w *Window)
	fFocusHolder           func(w *Window, focused bool)
	fIconifyHolder         func(w *Window, iconified bool)
	fMouseButtonHolder     func(w *Window, button MouseButton, action Action, mod ModifierKey)
	fCursorPosHolder       func(w *Window, xpos float64, ypos float64)
	fCursorEnterHolder     func(w *Window, entered bool)
	fScrollHolder          func(w *Window, xoff float64, yoff float64)
	fKeyHolder             func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
	fCharHolder            func(w *Window, char rune)
	fCharModsHolder        func(w *Window, char rune, mods ModifierKey)
	fDropHolder            func(w *Window, names []string)
}

type Cursor struct {
	data *_GLFWcursor
}

// iconID is the ID of the icon in the resource file.
const iconID = 1

var resources struct {
	handle syscall.Handle
	class  uint16
	cursor syscall.Handle
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
	// return C.glfwGetClipboardString(nil)
	return ""
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
// that can be set for a Window with SetCursor.
func CreateStandardCursor(shape StandardCursor) *Cursor {
	// c := C.glfwCreateStandardCursor(C.int(shape))
	panicError()
	return nil // &Cursor{c}
}

func SetProcessDPIAware() {
	_SetProcessDPIAware.Call()
}
func LoadCursor(curID uint16) (syscall.Handle, error) {
	h, _, err := _LoadCursor.Call(0, uintptr(curID))
	if h == 0 {
		return 0, fmt.Errorf("LoadCursorW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var s *_GLFWwindow
	if share != nil {
		s = share.data
	}
	w, err := glfwCreateWindow(width, height, title, monitor, s)
	if err != nil {
		return nil, fmt.Errorf("glfwCreateWindow failed: %v", err)
	}
	wnd := &Window{data: w}
	windowMap.put(wnd)
	return wnd, nil
}

func glfwCreateWindow(width, height int, title string, monitor *Monitor, share *_GLFWwindow) (*_GLFWwindow, error) {

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid width/heigth")
	}

	fbconfig := _glfw.hints.framebuffer
	ctxconfig := _glfw.hints.context
	wndconfig := _glfw.hints.window

	wndconfig.width = width
	wndconfig.height = height
	wndconfig.title = title
	ctxconfig.share = share
	// if _glfwIsValidContextConfig(&ctxconfig) != nil {
	//	return nil, fmt.Errorf("glfw context config is invalid: %v", ctxconfig)
	// }

	Window := &_GLFWwindow{}
	Window.next = _glfw.windowListHead
	_glfw.windowListHead = Window

	Window.videoMode.width = width
	Window.videoMode.height = height
	Window.videoMode.redBits = fbconfig.redBits
	Window.videoMode.greenBits = fbconfig.greenBits
	Window.videoMode.blueBits = fbconfig.blueBits
	Window.videoMode.refreshRate = _glfw.hints.refreshRate

	Window.monitor = monitor
	Window.resizable = wndconfig.resizable
	Window.decorated = wndconfig.decorated
	Window.autoIconify = wndconfig.autoIconify
	Window.floating = wndconfig.floating
	Window.focusOnShow = wndconfig.focusOnShow
	Window.cursorMode = GLFW_CURSOR_NORMAL

	Window.doublebuffer = fbconfig.doublebuffer

	Window.minwidth = GLFW_DONT_CARE
	Window.minheight = GLFW_DONT_CARE
	Window.maxwidth = GLFW_DONT_CARE
	Window.maxheight = GLFW_DONT_CARE
	Window.numer = GLFW_DONT_CARE
	Window.denom = GLFW_DONT_CARE

	SetProcessDPIAware()
	var err error
	Window.win32.handle, err = CreateWindowEx(
		WS_OVERLAPPED|WS_EX_APPWINDOW,
		0,
		"WindowName here",
		WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
		CW_USEDEFAULT, CW_USEDEFAULT, // Window position
		int32(width), int32(height), // Window width/heigth
		0, // No parent
		0, // No menu
		resources.handle,
		0)
	return Window, err
}

// SwapBuffers swaps the front and back buffers of the Window.
func (w *Window) SwapBuffers() {
	glfwSwapBuffers(w.data)
	panicError()
}

// SetCursor sets the cursor image to be used when the cursor is over the client area
func (w *Window) SetCursor(c *Cursor) {
	if c == nil {
		glfwSetCursor(w.data, nil)
	} else {
		// TODO glfwSetCursor(w.data, c)
	}
	panicError()
}

func glfwSetWindowPos(w *_GLFWwindow, xpos, ypos int) {

}

// SetPos sets the position, in screen coordinates, of the upper-left corner of the client area of the Window.
func (w *Window) SetPos(xpos, ypos int) {
	glfwSetWindowPos(w.data, xpos, ypos)
	panicError()
}

func glfwSetWindowSize(w *_GLFWwindow, xpos, ypos int) {

}

// SetSize sets the size, in screen coordinates, of the client area of the Window.
func (w *Window) SetSize(width, height int) {
	glfwSetWindowSize(w.data, width, height)
	panicError()
}

func glfwShowWindow(w *_GLFWwindow) {
	if w.monitor != nil {
		return
	}
	// TODO _glfwPlatformShowWindow(window);
	if w.focusOnShow {
		// TODO _glfwPlatformFocusWindow(window)
	}
}

// Show makes the Window visible, if it was previously hidden.
func (w *Window) Show() {
	glfwShowWindow(w.data)
	panicError()
}

func (w *Window) MakeContextCurrent() {
	// _GLFWWindow * Window = (_GLFWWindow *)handle;
	// _GLFWWindow * previous;
	// _GLFW_REQUIRE_INIT();
	// previous := glfwPlatformGetTls(&_glfw.contextSlot);
	// if w !=nil && w.client == GLFW_NO_API {
	//	panic("Cannot make current with a Window that has no OpenGL or OpenGL ES context");
	// }
	if w == nil {
		panic("Window is nil")
	}
	w.data.context.makeCurrent(&w.data.context)
	panicError()
}

// Focus brings the specified Window to front and sets input focus.
func (w *Window) Focus() {
	// TODO glfwFocusWindow(w.data)
}

// ShouldClose reports the value of the close flag of the specified Window.
func (w *Window) ShouldClose() bool {
	return w.data.shouldClose
}

// CursorPosCallback the cursor position callback.
type CursorPosCallback func(w *Window, xpos float64, ypos float64)

// SetCursorPosCallback sets the cursor position callback which is called
// when the cursor is moved. The callback is provided with the position relative
// to the upper-left corner of the client area of the Window.
func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	w.cursorPosCallback = cbfun
	return nil
}

// KeyCallback is the key callback.
type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

// SetKeyCallback sets the key callback which is called when a key is pressed, repeated or released.
func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	w.keyCallback = cbfun
	return nil
}

// CharCallback is the character callback.
type CharCallback func(w *Window, char rune)

// SetCharCallback sets the character callback which is called when a Unicode character is input.
func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	w.charCallback = cbfun
	return nil
}

// DropCallback is the drop callback which is called when an object is dropped over the Window.
type DropCallback func(w *Window, names []string)

// SetDropCallback sets the drop callback
func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	w.dropCallback = cbfun
	return nil
}

// ContentScaleCallback is the function signature for Window content scale
// callback functions.
type ContentScaleCallback func(w *Window, x float32, y float32)

// SetContentScaleCallback function sets the Window content scale callback of
// the specified Window, which is called when the content scale of the specified Window changes.
func (w *Window) SetContentScaleCallback(cbfun ContentScaleCallback) ContentScaleCallback {
	w.contentScaleCallback = cbfun
	return nil
}

// RefreshCallback is the Window refresh callback.
type RefreshCallback func(w *Window)

// SetRefreshCallback sets the refresh callback of the Window, which
// is called when the client area of the Window needs to be redrawn,
func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	w.refreshCallback = cbfun
	return nil
}

// FocusCallback is the Window focus callback.
type FocusCallback func(w *Window, focused bool)

// SetFocusCallback sets the focus callback of the Window, which is called when
// the Window gains or loses focus.
//
// After the focus callback is called for a Window that lost focus, synthetic key
// and mouse button release events will be generated for all such that had been
// pressed. For more information, see SetKeyCallback and SetMouseButtonCallback.
func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	w.focusCallback = cbfun
	return nil
}

// SizeCallback is the Window size callback.
type SizeCallback func(w *Window, width int, height int)

// SetSizeCallback sets the size callback of the Window, which is called when
// the Window is resized. The callback is provided with the size, in screen
// coordinates, of the client area of the Window.
func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.sizeCallback = cbfun
	return nil
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the Window and input
// callbacks associated with those events to be called.
func PollEvents() {
	glfwPollEvents()
	panicError()
}

// Terminate destroys all remaining Windows, frees any allocated resources and
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
	// TODO glfwTerminate()
}

func Init() error {
	return nil
}

// GetContentScale function retrieves the content scale for the specified
// Window. The content scale is the ratio between the current DPI and the
// platform's default DPI. If you scale all pixel dimensions by this scale then
// your content should appear at an appropriate size. This is especially
// important for text and any UI elements.
//
// This function may only be called from the main thread.
func (w *Window) GetContentScale() (float32, float32) {
	var x, y float32
	// C.glfwGetWindowContentScale(w.data, &x, &y)
	return float32(x), float32(y)
}

// GetFrameSize retrieves the size, in screen coordinates, of each edge of the frame
// of the specified Window. This size includes the title bar, if the Window has one.
// The size of the frame may vary depending on the Window-related hints used to create it.
//
// Because this function retrieves the size of each Window frame edge and not the offset
// along a particular coordinate axis, the retrieved values will always be zero or positive.
func (w *Window) GetFrameSize() (left, top, right, bottom int) {
	var l, t, r, b int
	// C.glfwGetWindowFrameSize(w.data, &l, &t, &r, &b)
	panicError()
	return int(l), int(t), int(r), int(b)
}

// SwapInterval sets the swap interval for the current context, i.e. the number
// of screen updates to wait before swapping the buffers of a Window and
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
	var xpos, ypos float32
	// C.glfwGetCursorPos(w.data, &xpos, &ypos)
	panicError()
	return float64(xpos), float64(ypos)
}

// GetSize returns the size, in screen coordinates, of the client area of the
// specified Window.
func (w *Window) GetSize() (width, height int) {
	var wi, h int
	// C.glfwGetWindowSize(w.data, &wi, &h)
	panicError()
	return int(wi), int(h)
}
