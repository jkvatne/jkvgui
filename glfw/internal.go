package glfw

import (
	"errors"
	"golang.org/x/sys/windows"
	"log"
	"sync"
	"syscall"
	"unicode"
	"unsafe"
)

// Monitor structure
//
type _GLFWmonitor = struct {
	name        [128]byte
	userPointer unsafe.Pointer
	widthMM     int
	heightMM    int
	// The window whose video mode is current on this monitor
	// window      *glfwWindow
	modes       *GLFWvidmode
	modeCount   int
	currentMode GLFWvidmode
	// originalRamp GLFWgammaramp
	// currentRamp  GLFWgammaramp
	// This is defined in the window API's platform.h
	// _GLFW_PLATFORM_MONITOR_STATE;
}

// Cursor structure
//
type _GLFWcursor struct {
	next   *_GLFWcursor
	handle syscall.Handle
}

type _GLFWmakecontextcurrentfun = func(w *_GLFWwindow) error
type _GLFWswapbuffersfun = func(w *_GLFWwindow)
type _GLFWswapintervalfun = func(interval int)
type _GLFWextensionsupportedfun = func(x byte) bool
type _GLFWgetprocaddressfun = func()
type _GLFWdestroycontextfun = func(w *_GLFWwindow)

// Context structure
//
type _GLFWcontext struct {
	client                  int
	source                  int
	major, minor, revision  int
	forward, debug, noerror bool
	profile                 int
	robustness              int
	release                 int
	// PFNGLGETSTRINGIPROC  GetStringi;
	// PFNGLGETINTEGERVPROC GetIntegerv;
	// PFNGLGETSTRINGPROC   GetString;
	makeCurrent        _GLFWmakecontextcurrentfun
	swapBuffers        _GLFWswapbuffersfun
	swapInterval       _GLFWswapintervalfun
	extensionSupported _GLFWextensionsupportedfun
	getProcAddress     _GLFWgetprocaddressfun
	destroy            _GLFWdestroycontextfun
	wgl                struct {
		dc       HDC
		handle   HANDLE
		interval int
	}
}

type _GLFWwindow struct {
	next *_GLFWwindow
	// Window settings and state
	resizable          bool
	decorated          bool
	autoIconify        bool
	floating           bool
	focusOnShow        bool
	shouldClose        bool
	userPointer        unsafe.Pointer
	doublebuffer       bool
	videoMode          GLFWvidmode
	monitor            *Monitor
	cursor             *_GLFWcursor
	minwidth           int
	minheight          int
	maxwidth           int
	maxheight          int
	numer              int
	denom              int
	stickyKeys         bool
	stickyMouseButtons bool
	lockKeyMods        bool
	cursorMode         int
	mouseButtons       [MouseButtonLast + 1]byte
	keys               [KeyLast + 1]byte
	// Virtual cursor position when cursor is disabled
	virtualCursorPosX float64
	virtualCursorPosY float64
	rawMouseMotion    bool
	context           _GLFWcontext
	lastCursorPosX    float64 // The last received cursor position, regardless of source
	lastCursorPosY    float64 // The last received cursor position, regardless of source

	charCallback           CharCallback
	focusCallback          FocusCallback
	keyCallback            KeyCallback
	mouseButtonCallback    MouseButtonCallback
	cursorPosCallback      CursorPosCallback
	scrollCallback         ScrollCallback
	refreshCallback        RefreshCallback
	sizeCallback           SizeCallback
	dropCallback           DropCallback
	contentScaleCallback   ContentScaleCallback
	fFramebufferSizeHolder func(w *_GLFWwindow, width int, height int)
	fCloseHolder           func(w *_GLFWwindow)
	fMaximizeHolder        func(w *_GLFWwindow, maximized bool)
	fIconifyHolder         func(w *_GLFWwindow, iconified bool)
	fCursorEnterHolder     func(w *_GLFWwindow, entered bool)
	fCharModsHolder        func(w *_GLFWwindow, char rune, mods ModifierKey)

	Win32 _GLFWwindowWin32
}

type _GLFWwindowWin32 = struct {
	handle         syscall.Handle
	bigIcon        syscall.Handle
	smallIcon      syscall.Handle
	cursorTracked  bool
	frameAction    bool
	iconified      bool
	maximized      bool
	transparent    bool // Whether to enable framebuffer transparency on DWM
	scaleToMonitor bool
	width          int    // Cached size used to filter out duplicate events
	height         int    // Cached size used to filter out duplicate events
	highSurrogate  uint16 // The last recevied high surrogate when decoding pairs of UTF-16 messages
}

type _GLFWinitconfig = struct {
	hatButtons bool
	ns         struct {
		menubar bool
		chdir   bool
	}
	wl struct {
		libdecorMode int
	}
}
type _GLFWwndconfig = struct {
	xpos           int
	ypos           int
	width          int
	height         int
	title          string
	resizable      bool
	visible        bool
	decorated      bool
	focused        bool
	autoIconify    bool
	floating       bool
	maximized      bool
	centerCursor   bool
	focusOnShow    bool
	scaleToMonitor bool
	ns             struct {
		retina    bool
		frameName string
	}
}

type _GLFWctxconfig = struct {
	client     int
	source     int
	major      int
	minor      int
	forward    bool
	debug      bool
	noerror    bool
	profile    int
	robustness int
	release    int
	share      *_GLFWwindow
	nsgl       struct {
		offline bool
	}
}

type hints = struct {
	init        _GLFWinitconfig
	framebuffer _GLFWfbconfig
	window      _GLFWwndconfig
	context     _GLFWctxconfig
	refreshRate int
}

type _GLFWfbconfig = struct {
	redBits        int
	greenBits      int
	blueBits       int
	alphaBits      int
	depthBits      int
	stencilBits    int
	accumRedBits   int
	accumGreenBits int
	accumBlueBits  int
	accumAlphaBits int
	auxBuffers     int
	stereo         bool
	samples        int
	sRGB           bool
	doublebuffer   bool
	transparent    bool
	handle         uintptr
}

type _GLFWerror struct {
	next        *_GLFWerror
	code        int
	description string
}

type _GLFWtls = struct {
	allocated bool
	index     int
}

// Library global Data
var _glfw struct {
	hints
	class          uint16
	available      bool
	instance       syscall.Handle
	initialized    bool
	errorListHead  *_GLFWerror
	cursorListHead *_GLFWcursor
	windowListHead *_GLFWwindow
	monitors       []_GLFWmonitor
	errorSlot      _GLFWtls
	contextSlot    _GLFWtls
	errorLock      sync.Mutex
	win32          struct {
		helperWindowHandle syscall.Handle
		helperWindowClass  uint16
		mainWindowClass    uint16
	}
	wgl struct {
		dc                         HDC
		handle                     syscall.Handle
		interval                   int
		instance                   *windows.LazyDLL
		wglCreateContextAttribsARB *windows.LazyProc
		wglDeleteContext           *windows.LazyProc
		wglGetProcAddress          *windows.LazyProc
		wglGetCurrentDC            *windows.LazyProc
		wglGetCurrentContext       *windows.LazyProc
		wglMakeCurrent             *windows.LazyProc
		wglShareLists              *windows.LazyProc
		wglSwapBuffers             *windows.LazyProc
		wglCreateContext           *windows.LazyProc
		wglSetPixelFormat          *windows.LazyProc
		wglChoosePixelFormat       *windows.LazyProc
		wglDescribePixelFormat     *windows.LazyProc
		getProcAddress             *windows.LazyProc
		GetExtensionsStringEXT     *windows.LazyProc
		GetExtensionsStringARB     *windows.LazyProc
		GetPixelFormatAttribivARB  *windows.LazyProc
		ARB_pixel_format           int
		ARB_multisample            bool
		ARB_framebuffer_sRGB       bool
		EXT_framebuffer_sRGB       bool
		EXT_colorspace             bool
	}
}

/*
func getModifiers() key.Modifiers {
	var kmods key.Modifiers
	if GetKeyState(VK_LWIN)&0x1000 != 0 || GetKeyState(VK_RWIN)&0x1000 != 0 {
		kmods |= key.ModSuper
	}
	if GetKeyState(VK_MENU)&0x1000 != 0 {
		kmods |= key.ModAlt
	}
	if GetKeyState(VK_CONTROL)&0x1000 != 0 {
		kmods |= key.ModCtrl
	}
	if GetKeyState(VK_SHIFT)&0x1000 != 0 {
		kmods |= key.ModShift
	}
	return kmods
}
*/
func glfwInputKey(window *_GLFWwindow, key Key, scancode int, action int, mods ModifierKey) {
	var repeated bool
	if key >= 0 && key <= KeyLast {
		repeated = false

		if action == GLFW_RELEASE && window.keys[key] == GLFW_RELEASE {
			return
		}

		if action == GLFW_PRESS && window.keys[key] == GLFW_PRESS {
			repeated = true
		}

		/*		if (action == GLFW_RELEASE && window.stickyKeys) {
					window.keys[key] = _GLFW_STICK;
				} else {
					window.keys[key] = (char)
					action
				}
		*/
		if repeated {
			action = GLFW_REPEAT
		}
	}
	/*
		if (!window.lockKeyMods) {
			mods &= ~(GLFW_MOD_CAPS_LOCK | GLFW_MOD_NUM_LOCK)
		} */

	if window.keyCallback != nil {
		w := windowMap.get(window)
		window.keyCallback(w, key, scancode, Action(action), mods)
	}

}
func glfwInputMouseClick(window *_GLFWwindow, button MouseButton, action Action, mods ModifierKey) {
	// TODO if (!window.lockKeyMods)	mods &= ~(GLFW_MOD_CAPS_LOCK | GLFW_MOD_NUM_LOCK);
	// TODO if (action == GLFW_RELEASE && window.stickyMouseButtons) window.mouseButtons[button] = _GLFW_STICK; else window.mouseButtons[button] = (char) action;
	w := windowMap.get(window)
	if window.mouseButtonCallback != nil {
		window.mouseButtonCallback(w, button, action, mods)
	}
}

// Notifies shared code that a window has lost or received input focus
func glfwInputWindowFocus(window *_GLFWwindow, focused bool) {
	if window == nil {
		return
	}
	if window.focusCallback != nil {
		w := windowMap.get(window)
		window.focusCallback(w, focused)
	}
	if !focused {
		// Force release of buttons
		/* TODO
		for k := Key(0);  k <= KeyLast;  k++ {
			if (window.keys[k] == GLFW_PRESS) {
				scancode := glfwPlatformGetKeyScancode(k);
				glfwInputKey(window, k, scancode, GLFW_RELEASE, 0);
			}
		}*/
		for button := MouseButton(0); button <= MouseButtonLast; button++ {
			if window.mouseButtons[button] == GLFW_PRESS {
				glfwInputMouseClick(window, button, GLFW_RELEASE, 0)
			}
		}
	}
}

func glfwInputCursorPos(window *_GLFWwindow, xpos, ypos float64) {
	w := windowMap.get(window)
	if window.virtualCursorPosX == xpos && window.virtualCursorPosY == ypos {
		return
	}
	window.virtualCursorPosX = xpos
	window.virtualCursorPosY = ypos
	if window.cursorPosCallback != nil {
		window.cursorPosCallback(w, xpos, ypos)
	}
}

func glfwInputScroll(window *_GLFWwindow, xoffset, yoffset float64) {
	w := windowMap.get(window)
	if window.scrollCallback != nil {
		window.scrollCallback(w, xoffset, yoffset)
	}
}

func glfwInputWindowDamage(window *_GLFWwindow) {
	w := windowMap.get(window)
	if window.refreshCallback != nil {
		window.refreshCallback(w)
	}
}

func getKeyMods() ModifierKey {
	var mods ModifierKey
	if GetKeyState(VK_SHIFT)&0x8000 != 0 {
		mods |= ModShift
	}
	if GetKeyState(VK_CONTROL)&0x8000 != 0 {
		mods |= ModControl
	}
	if GetKeyState(VK_MENU)&0x8000 != 0 {
		mods |= ModAlt
	}
	if (GetKeyState(VK_LWIN)|GetKeyState(VK_RWIN))&0x8000 != 0 {
		mods |= ModSuper
	}
	if (GetKeyState(VK_CAPITAL) & 1) != 0 {
		mods |= ModCapsLock
	}
	if (GetKeyState(VK_NUMLOCK) & 1) != 0 {
		mods |= ModNumLock
	}
	return mods
}

var winMap map[syscall.Handle]*_GLFWwindow

func GetProp(hwnd syscall.Handle) *_GLFWwindow {
	return winMap[hwnd]
}

func SetProp(hwnd syscall.Handle, prop *_GLFWwindow) {
	if winMap == nil {
		winMap = make(map[syscall.Handle]*_GLFWwindow)
	}
	winMap[hwnd] = prop
}

func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	window := GetProp(hwnd)
	log.Printf("msg=%d\n", msg)
	switch msg {
	case WM_CLOSE:
		window.shouldClose = true
	case WM_UNICHAR:
		if wParam == UNICODE_NOCHAR {
			// Tell the system that we accept WM_UNICHAR messages.
			return TRUE
		}
		fallthrough
	case WM_CHAR, WM_SYSCHAR:
		if r := rune(wParam); unicode.IsPrint(r) {
			window.charCallback(nil, r)
		}
		return TRUE
	case WM_DPICHANGED:
		// Let Windows know we're prepared for runtime DPI changes.
		return TRUE
	case WM_ERASEBKGND:
		// Avoid flickering between GPU content and background color.
		return TRUE
	case WM_KEYDOWN, WM_KEYUP, WM_SYSKEYDOWN, WM_SYSKEYUP:
		var key Key
		action := GLFW_PRESS
		if (lParam>>16)&0x100 != 0 {
			action = GLFW_RELEASE
		}
		mods := getKeyMods()
		scancode := int((lParam >> 16) & 0x1ff)
		switch scancode {
		case 0: // scancode = MapVirtualKeyW((UINT) wParam, MAPVK_VK_TO_VSC);
		case 0x54:
			scancode = 0x137 // Alt+PrtSc
		case 0x146:
			scancode = 0x45 // Ctrl+Pause
		case 0x136:
			scancode = 0x36 // CJK IME sets the extended bit for right Shift
		}

		key = Key(scancode) // TODO Needs keycodes[scancode]
		if wParam == VK_CONTROL {
			if lParam>>16&KF_EXTENDED != 0 {
				// Right side keys have the extended key bit set
				key = KeyRightControl
			} else {
				/*
					// NOTE: Alt Gr sends Left Ctrl followed by Right Alt
					// HACK: We only want one event for Alt Gr, so if we detect
					//       this sequence we discard this Left Ctrl message now
					//       and later report Right Alt normally
					MSG next;
					const DWORD time = GetMessageTime();

					if (PeekMessageW(&next, NULL, 0, 0, PM_NOREMOVE)) {
						if (next.message == WM_KEYDOWN ||
							next.message == WM_SYSKEYDOWN ||
							next.message == WM_KEYUP ||
							next.message == WM_SYSKEYUP)
						{
							if (next.wParam == VK_MENU &&
								(HIWORD(next.lParam) & KF_EXTENDED) &&
								next.time == time)
							{
								// Next message is Right Alt down so discard this
								break;
							}
						}
					}
				*/
				// This is a regular Left Ctrl message
				key = KeyLeftControl
			}
		}

		if action == GLFW_RELEASE && wParam == VK_SHIFT {
			// HACK: Release both Shift keys on Shift up event, as when both
			//       are pressed the first release does not emit any event
			// NOTE: The other half of this is in _glfwPlatformPollEvents
			glfwInputKey(window, KeyLeftShift, scancode, action, mods)
			glfwInputKey(window, KeyRightShift, scancode, action, mods)
		} else if wParam == VK_SNAPSHOT {
			// HACK: Key down is not reported for the Print Screen key
			glfwInputKey(window, key, scancode, GLFW_PRESS, mods)
			glfwInputKey(window, key, scancode, GLFW_RELEASE, mods)
		} else {
			glfwInputKey(window, key, scancode, action, mods)
		}
		break

	case WM_LBUTTONDOWN, WM_LBUTTONUP, WM_RBUTTONDOWN, WM_RBUTTONUP, WM_MBUTTONDOWN, WM_MBUTTONUP:
		var button MouseButton
		if msg == WM_LBUTTONDOWN || msg == WM_LBUTTONUP {
			button = MouseButtonLeft
		} else if msg == WM_RBUTTONDOWN || msg == WM_RBUTTONUP {
			button = MouseButtonRight
		} else if msg == WM_MBUTTONDOWN || msg == WM_MBUTTONUP {
			button = MouseButtonMiddle
		}
		var action Action
		if msg == WM_LBUTTONDOWN || msg == WM_RBUTTONDOWN || msg == WM_MBUTTONDOWN {
			action = GLFW_PRESS
		} else {
			action = GLFW_RELEASE
		}
		var i MouseButton
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == GLFW_PRESS {
				break
			}
		}
		// if i > MouseButtonLast {
		// TODO SetCapture(hWnd);
		// }

		glfwInputMouseClick(window, button, action, getKeyMods())
		for i = MouseButtonFirst; i <= MouseButtonLast; i++ {
			if window.mouseButtons[i] == GLFW_PRESS {
				break
			}
		}
		// if (i > MouseButtonLast)
		// TODO ReleaseCapture();
		// }

		return 0

	// TODO case WM_CANCELMODE:

	case WM_SETFOCUS:
		glfwInputWindowFocus(window, true)
		// HACK: Do not disable cursor while the user is interacting with a caption button
		// TODO if (window.Win32.frameAction) break;
		// TODO if (window.cursorMode == GLFW_CURSOR_DISABLED)	disableCursor(window);
		return 0
	case WM_KILLFOCUS:
		// TODO if (window.cursorMode == GLFW_CURSOR_DISABLED) enableCursor(window);
		// TODO if (window.monitor && window.autoIconify) _glfwPlatformIconifyWindow(window);
		glfwInputWindowFocus(window, false)
		return 0

	case WM_MOUSEMOVE:
		x := float64(int(lParam & 0xffff))
		y := float64(int((lParam >> 16) & 0xffff))
		if !window.Win32.cursorTracked {
			// tme.dwFlags = TME_LEAVE;
			// tme.hwndTrack = window.handle;
			// TrackMouseEvent(&tme);
			// window.cursorTracked = true;
			// glfwInputCursorEnter(window, GLFW_TRUE);
		}

		if window.cursorMode == CursorDisabled {
			dx := float64(x) - window.lastCursorPosX
			dy := float64(y) - window.lastCursorPosY
			// TODO if _glfw.Win32.disabledCursorWindow != window {			break			}
			glfwInputCursorPos(window, window.virtualCursorPosX+dx, window.virtualCursorPosY+dy)
		} else {
			glfwInputCursorPos(window, x, y)
		}
		window.lastCursorPosX = x
		window.lastCursorPosY = y
		return 0

	case WM_MOUSEWHEEL:
		glfwInputScroll(window, 0.0, float64(int16(wParam>>16))/120.0)
		return 0

	case WM_MOUSEHWHEEL:
		glfwInputScroll(window, -float64(int16(wParam>>16))/120.0, 0.0)
		return 0

	case WM_PAINT:
		glfwInputWindowDamage(window)

	case WM_SIZE:
		// TODO
		// return TRUE

	case WM_GETMINMAXINFO:
		// TODO

	case WM_SETCURSOR:
		// TODO
	}

	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func glfwPlatformPollEvents() {
	var msg Msg
	for PeekMessage(&msg, 0, 0, 0, PM_REMOVE) {
		if msg.Message == WM_QUIT {
			// NOTE: While GLFW does not itself post WM_QUIT, other processes may post it to this one, for example Task Manager
			// HACK: Treat WM_QUIT as a close on all windows
			// window = _glfw.windowListHead;
			// while (window {
			//	_glfwInputWindowCloseRequest(window);
			//	window = window- > next;
			// }
		} else {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	// HACK: Release modifier keys that the system did not emit KEYUP for
	// NOTE: Shift keys on Windows tend to "stick" when both are pressed as no key up message is generated by the first key release
	// NOTE: Windows key is not reported as released by the Win+V hotkey. Other Win hotkeys are handled implicitly by _glfwInputWindowFocus
	//       because they change the input focus
	// NOTE: The other half of this is in the WM_*KEY* handler in windowProc
	/* TODO
	handle = GetActiveWindow();
	if (handle!=nil) {
		window := 74W(handle, "GLFW");
		if window != nil {
			//const keys[4][2] = {{ VK_LSHIFT, GLFW_KEY_LEFT_SHIFT },    { VK_RSHIFT, GLFW_KEY_RIGHT_SHIFT },{ VK_LWIN, GLFW_KEY_LEFT_SUPER },{ VK_RWIN, GLFW_KEY_RIGHT_SUPER }}
			for i := 0; i < 4; i++ {
				vk := keys[i][0];
				key := keys[i][1];
				// scancode := Win32.scancodes[key];
				if GetKeyState(vk) & 0x8000 || window.keys[key] != GLFW_PRESS {
					continue;
				}
				_glfwInputKey(window, key, scancode, GLFW_RELEASE, getKeyMods());
			}
		}
	}
	window := _glfw.Win32.disabledCursorWindow;
	if window!=nil {
		var width, height int
		glfwPlatformGetWindowSize(window, &width, &height);
		// NOTE: Re-center the cursor only if it has moved since the last call, to avoid breaking glfwWaitEvents with WM_MOUSEMOVE
		if window.lastCursorPosX != width / 2 || window.lastCursorPosY != height / 2 {
			glfwPlatformSetCursorPos(window, width / 2, height / 2);
		}
	}
	*/
}

func glfwSwapBuffers(window *_GLFWwindow) {
	if window.context.client == 0 {
		panic("Cannot swap buffers of a window that has no OpenGL or OpenGL ES context")
	}
	window.context.swapBuffers(window)
}

func cursorInContentArea(window *_GLFWwindow) bool {
	/*var area RECT
	var pos	POINT
	if (!GetCursorPos(&pos)) {
		return false
	}
	if WindowFromPoint(pos) != window.Win32.handle {
		return false;
	}
	GetClientRect(window.Win32.handle, &area);
	ClientToScreen(window.Win32.handle, (POINT*) &area.left);
	ClientToScreen(window.Win32.handle, (POINT*) &area.right);
	return PtInRect(&area, pos);
	*/
	return true
}

func glfwSetCursor(window *_GLFWwindow, cursor *_GLFWcursor) {
	window.cursor = cursor
	if cursorInContentArea(window) {
		// TODO updateCursorImage(window)
	}
}

func SetFocus(window *_GLFWwindow) {
	r1, _, err := _SetFocus.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetFocus failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetFocus failed")
	}
}

func BringWindowToTop(window *_GLFWwindow) {
	r1, _, err := _BringWindowToTop.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("BringWindowToTop failed, " + err.Error())
	}
	if r1 == 0 {
		panic("BringWindowToTop failed")
	}
}

func SetForegroundWindow(window *_GLFWwindow) {
	r1, _, err := _SetForegroundWindow.Call(uintptr(unsafe.Pointer(window.Win32.handle)))
	if r1 == 0 || err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetForegroundWindow failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetForegroundWindow failed")
	}
}

func glfwFocusWindow(window *_GLFWwindow) {
	BringWindowToTop(window)
	SetForegroundWindow(window)
	SetFocus(window)
}
