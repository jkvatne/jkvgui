package glfw

import (
	"errors"
	"fmt"
	"golang.design/x/clipboard"
	"golang.org/x/sys/windows"
	"log/slog"
	"syscall"
	"unsafe"
)

const (
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
	SWP_NOREDRAW       = 0x0008
	SWP_NOACTIVATE     = 0x0010
	SWP_FRAMECHANGED   = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOOWNERZORDER  = 0x0200
	SWP_NOSENDCHANGING = 0x0400
)

var resources struct {
	handle syscall.Handle
	class  uint16
	cursor syscall.Handle
}

type Point struct {
	X, Y int32
}

type Msg struct {
	Hwnd     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

var (
	kernel32                = windows.NewLazySystemDLL("kernel32.dll")
	_procGetModuleHandleExW = kernel32.NewProc("GetModuleHandleExW")
	_GetModuleHandleW       = kernel32.NewProc("GetModuleHandleW")
	_GlobalAlloc            = kernel32.NewProc("GlobalAlloc")
	_GlobalFree             = kernel32.NewProc("GlobalFree")
	_GlobalLock             = kernel32.NewProc("GlobalLock")
	_GlobalUnlock           = kernel32.NewProc("GlobalUnlock")

	user32                         = windows.NewLazySystemDLL("user32.dll")
	_SetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	_EnumDisplayMonitors           = user32.NewProc("EnumDisplayMonitors")
	_EnumDisplayAdapters           = user32.NewProc("EnumDisplayAdaptersA")
	_EnumDisplayDevices            = user32.NewProc("EnumDisplayDevicesW")
	_EnumDisplaySettings           = user32.NewProc("EnumDisplaySettingsW")
	_GetMonitorInfo                = user32.NewProc("GetMonitorInfoW")
	_ToUnicode                     = user32.NewProc("ToUnicode")
	_MapVirtualKeyW                = user32.NewProc("MapVirtualKeyW")
	_AdjustWindowRectEx            = user32.NewProc("AdjustWindowRectEx")
	_CallMsgFilter                 = user32.NewProc("CallMsgFilterW")
	_CloseClipboard                = user32.NewProc("CloseClipboard")
	_CreateWindowEx                = user32.NewProc("CreateWindowExW")
	_DefWindowProc                 = user32.NewProc("DefWindowProcW")
	_DestroyWindow                 = user32.NewProc("DestroyWindow")
	_DispatchMessage               = user32.NewProc("DispatchMessageW")
	_EmptyClipboard                = user32.NewProc("EmptyClipboard")
	_GetWindowRect                 = user32.NewProc("GetWindowRect")
	_GetClientRect                 = user32.NewProc("GetClientRect")
	_GetClipboardData              = user32.NewProc("GetClipboardData")
	_GetDC                         = user32.NewProc("GetDC")
	_GetDpiForWindow               = user32.NewProc("GetDpiForWindow")
	_GetKeyState                   = user32.NewProc("GetKeyState")
	_GetMessage                    = user32.NewProc("GetMessageW")
	_GetMessageTime                = user32.NewProc("GetMessageTime")
	_GetSystemMetrics              = user32.NewProc("GetSystemMetrics")
	_GetWindowLong                 = user32.NewProc("GetWindowLongPtrW")
	_GetWindowLong32               = user32.NewProc("GetWindowLongW")
	_GetWindowPlacement            = user32.NewProc("GetWindowPlacement")
	_KillTimer                     = user32.NewProc("KillTimer")
	_LoadCursor                    = user32.NewProc("LoadCursorW")
	_LoadImage                     = user32.NewProc("LoadImageW")
	_MonitorFromPoint              = user32.NewProc("MonitorFromPoint")
	_MonitorFromWindow             = user32.NewProc("MonitorFromWindow")
	_MoveWindow                    = user32.NewProc("MoveWindow")
	_MsgWaitForMultipleObjectsEx   = user32.NewProc("MsgWaitForMultipleObjectsEx")
	_OpenClipboard                 = user32.NewProc("OpenClipboard")
	_PeekMessage                   = user32.NewProc("PeekMessageW")
	_PostMessage                   = user32.NewProc("PostMessageW")
	_PostQuitMessage               = user32.NewProc("PostQuitMessage")
	_ReleaseCapture                = user32.NewProc("ReleaseCapture")
	_RegisterClassExW              = user32.NewProc("RegisterClassExW")
	_ReleaseDC                     = user32.NewProc("releaseDC")
	_ScreenToClient                = user32.NewProc("ScreenToClient")
	_ShowWindow                    = user32.NewProc("ShowWindow")
	_SetCapture                    = user32.NewProc("SetCapture")
	_SetCursor                     = user32.NewProc("SetCursor")
	_SetClipboardData              = user32.NewProc("SetClipboardData")
	_SetForegroundWindow           = user32.NewProc("SetForegroundWindow")
	_SetFocus                      = user32.NewProc("SetFocus")
	_SetProcessDPIAware            = user32.NewProc("SetProcessDPIAware")
	_SetTimer                      = user32.NewProc("SetTimer")
	_SetWindowLong                 = user32.NewProc("SetWindowLongPtrW")
	_SetWindowLong32               = user32.NewProc("SetWindowLongW")
	_SetWindowPlacement            = user32.NewProc("SetWindowPlacement")
	_SetWindowPos                  = user32.NewProc("SetWindowPos")
	_SetWindowText                 = user32.NewProc("SetWindowTextW")
	_TranslateMessage              = user32.NewProc("TranslateMessage")
	_UnregisterClass               = user32.NewProc("UnregisterClassW")
	_UpdateWindow                  = user32.NewProc("UpdateWindow")
	_BringWindowToTop              = user32.NewProc("BringWindowToTop")
	_GetCursorPos                  = user32.NewProc("GetCursorPos")

	shcore                    = windows.NewLazySystemDLL("shcore")
	_GetDpiForMonitor         = shcore.NewProc("GetDpiForMonitor")
	_GetScaleFactorForMonitor = shcore.NewProc("GetScaleFactorForMonitor")
)

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

func GetKeyState(nVirtKey int) uint16 {
	c, _, _ := _GetKeyState.Call(uintptr(nVirtKey))
	return uint16(c)
}
func GetModuleHandle() (syscall.Handle, error) {
	h, _, err := _GetModuleHandleW.Call(uintptr(0))
	if h == 0 {
		return 0, fmt.Errorf("GetModuleHandleW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func RegisterClassEx(cls *WndClassEx) (uint16, error) {
	a, _, err := _RegisterClassExW.Call(uintptr(unsafe.Pointer(cls)))
	if a == 0 {
		return 0, fmt.Errorf("RegisterClassExW failed: %v", err)
	}
	return uint16(a), nil
}

/*
func GetModuleHandle(modulename *uint16) (hMonitor syscall.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(_GetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(modulename)), 0, 0)
	hMonitor = syscall.Handle(r0)
	if hMonitor == 0 {
		err = fmt.Errorf("GetModuleHandle error %v", e1)
	}
	return
}
*/
func LoadImage(hInst syscall.Handle, res uint32, typ uint32, cx, cy int, fuload uint32) (syscall.Handle, error) {
	h, _, err := _LoadImage.Call(uintptr(hInst), uintptr(res), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuload))
	if h == 0 {
		return 0, fmt.Errorf("LoadImageW failed: %v", err)
	}
	return syscall.Handle(h), nil
}

func glfwSetSize(window *Window, width, height int) {
	if glfwIsWindows10Version1607OrGreater() {
		// AdjustWindowRectExForDpi(&rect, getWindowStyle(window),	FALSE, getWindowExStyle(window), GetDpiForWindow(window.win32.hMonitor));
	} else {
		// AdjustWindowRectEx(&rect, getWindowStyle(window), FALSE, getWindowExStyle(window));
	}
	// glfwSetWi	r1, _, err := _SetWindowPos.Call(uintptr(hWnd), uintptr(after), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(flags))
	//	if err != nil && !errors.Is(err, syscall.Errno(0)) {
	//		panic("SetWindowPos failed, " + err.Error())
	//	}ndowPos(window)
	_, _, err := _SetWindowPos.Call(uintptr(window.Win32.handle), 0, 0, 0, uintptr(width), uintptr(height), uintptr(SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOMOVE|SWP_NOZORDER))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetWindowPos failed, " + err.Error())
	}
	// SetWindowPos(window.Win32.handle, 0, 0, 0, width, height, SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_NOMOVE|SWP_NOZORDER)
	// SetWindowPos(window->win32.hMonitor, HWND_TOP, 0, 0, rect.right - rect.left, rect.bottom - rect.top,SWP_NOACTIVATE | SWP_NOOWNERZORDER | SWP_NOMOVE | SWP_NOZORDER);
}

func CreateWindowEx(dwExStyle uint32, lpClassName uint16, lpWindowName string, dwStyle uint32, x, y, w, h int32, hWndParent, hMenu, hInstance syscall.Handle, lpParam uintptr) (syscall.Handle, error) {
	wname, _ := syscall.UTF16PtrFromString(lpWindowName)
	hwnd, _, err := _CreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(lpClassName),
		uintptr(unsafe.Pointer(wname)),
		uintptr(dwStyle),
		uintptr(x), uintptr(y),
		uintptr(w), uintptr(h),
		uintptr(hWndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(lpParam))
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowEx failed: %v", err)
	}
	return syscall.Handle(hwnd), nil
}

func PeekMessage(m *Msg, hwnd syscall.Handle, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	r, _, _ := _PeekMessage.Call(uintptr(unsafe.Pointer(m)), uintptr(hwnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax), uintptr(wRemoveMsg))
	return r != 0
}

func TranslateMessage(m *Msg) {
	_TranslateMessage.Call(uintptr(unsafe.Pointer(m)))
}

func DispatchMessage(m *Msg) {
	_DispatchMessage.Call(uintptr(unsafe.Pointer(m)))
}

// PollEvents processes only those events that have already been received and
// then returns immediately. Processing events will cause the Window and input
// callbacks associated with those events to be called.
// this was called glfwPollEvents()
func PollEvents() {
	var msg Msg
	for PeekMessage(&msg, 0, 0, 0, PM_REMOVE) {
		if msg.Message == WM_QUIT {
			// NOTE: While GLFW does not itself post WM_QUIT, other processes
			//       may post it to this one, for example Task Manager
			// HACK: Treat WM_QUIT as a close on all windows
			window := _glfw.windowListHead
			for window != nil {
				// TODO _glfwInputWindowCloseRequest(window)
				window = window.next
			}
		} else {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}

	// HACK: Release modifier keys that the system did not emit KEYUP for
	// NOTE: Shift keys on Windows tend to "stick" when both are pressed as
	//       no key up message is generated by the first key release
	// NOTE: Windows key is not reported as released by the Win+V hotkey
	//       Other Win hotkeys are handled implicitly by _glfwInputWindowFocus
	//       because they change the input focus
	// NOTE: The other half of this is in the WM_*KEY* handler in windowProc
	/* TODO
	hMonitor = GetActiveWindow()
	if hMonitor != 0 {
		window = GetPropW(hMonitor, "GLFW")
		if window != nil {
			keys := [4][2]int{{VK_LSHIFT, KeyLeftShift}, {VK_RSHIFT, KeyRightShift}, {VK_LWIN, KeyLeftSuper}, {VK_RWIN, KeyRightSuper}}
			for i := 0; i < 4; i++ {
				vk := keys[i][0]
				key := keys[i][1]
				scancode := _glfw.scancodes[key]
				if (GetKeyState(vk)&0x8000 != 0) || (window.keys[key] != GLFW_PRESS) {
					continue
				}
				_glfwInputKey(window, key, scancode, GLFW_RELEASE, getKeyMods())
			}
		}
	}
	window = _glfw.disabledCursorWindow
	if window != nil {
		var width, height int
		// TODO _glfwPlatformGetWindowSize(window, &width, &height);
		// NOTE: Re-center the cursor only if it has moved since the last call,
		//       to avoid breaking glfwWaitEvents with WM_MOUSEMOVE
		if window.Win32.lastCursorPosX != width/2 || window.Win32.lastCursorPosY != height/2 {
			// TODO _glfwPlatformSetCursorPos(window, width / 2, height / 2);
		}
	}*/
}

func glfwIsValidContextConfig(ctxconfig *_GLFWctxconfig) error {
	if ctxconfig.source != GLFW_NATIVE_CONTEXT_API && ctxconfig.source != GLFW_EGL_CONTEXT_API && ctxconfig.source != GLFW_OSMESA_CONTEXT_API {
		return fmt.Errorf("Invalid context creation API")
	}
	if ctxconfig.client != GLFW_NO_API && ctxconfig.client != GLFW_OPENGL_API && ctxconfig.client != GLFW_OPENGL_ES_API {
		return fmt.Errorf("Invalid client API")
	}
	if ctxconfig.share != nil {
		if ctxconfig.client == GLFW_NO_API || ctxconfig.share.context.client == GLFW_NO_API {
			return fmt.Errorf("Invalid share")
		}
		if ctxconfig.source != ctxconfig.share.context.source {
			return fmt.Errorf("Invalid share")
		}
	}

	if ctxconfig.client == GLFW_OPENGL_API {
		if (ctxconfig.major < 1 || ctxconfig.minor < 0) ||
			(ctxconfig.major == 1 && ctxconfig.minor > 5) ||
			(ctxconfig.major == 2 && ctxconfig.minor > 1) ||
			(ctxconfig.major == 3 && ctxconfig.minor > 3) {
			// OpenGL 1.0 is the smallest valid version
			// OpenGL 1.x series ended with version 1.5
			// OpenGL 2.x series ended with version 2.1
			// OpenGL 3.x series ended with version 3.3
			// For now, let everything else through
			return fmt.Errorf("Invalid OpenGL version %i.%i", ctxconfig.major, ctxconfig.minor)
		}

		if ctxconfig.profile != 0 {
			if ctxconfig.profile != GLFW_OPENGL_CORE_PROFILE && ctxconfig.profile != GLFW_OPENGL_COMPAT_PROFILE {
				return fmt.Errorf("Invalid OpenGL profile 0x%08X", ctxconfig.profile)
			}
			if ctxconfig.major <= 2 || (ctxconfig.major == 3 && ctxconfig.minor < 2) {
				// Desktop OpenGL context profiles are only defined for version 3.2 and above
				return fmt.Errorf("Context profiles are only defined for OpenGL version 3.2 and above")
			}
		}
		if ctxconfig.forward && ctxconfig.major <= 2 {
			// Forward-compatible contexts are only defined for OpenGL version 3.0 and above
			return fmt.Errorf("Forward-compatibility is only defined for OpenGL version 3.0 and above")
		}
	} else if ctxconfig.client == GLFW_OPENGL_ES_API {
		if ctxconfig.major < 1 || ctxconfig.minor < 0 || (ctxconfig.major == 1 && ctxconfig.minor > 1) || (ctxconfig.major == 2 && ctxconfig.minor > 0) {
			// OpenGL ES 1.0 is the smallest valid version
			// OpenGL ES 1.x series ended with version 1.1
			// OpenGL ES 2.x series ended with version 2.0
			// For now, let everything else through
			return fmt.Errorf("Invalid OpenGL ES version %i.%i", ctxconfig.major, ctxconfig.minor)
		}
	}
	// if ctxconfig.robustness > 0 && ctxconfig.robustness != GLFW_NO_RESET_NOTIFICATION && ctxconfig.robustness != GLFW_LOSE_CONTEXT_ON_RESET {
	//	return fmt.Errorf("Invalid context robustness mode 0x%08X", ctxconfig.robustness)
	// }

	// if ctxconfig.release > 0 && ctxconfig.release != GLFW_RELEASE_BEHAVIOR_NONE && ctxconfig.release != GLFW_RELEASE_BEHAVIOR_FLUSH {
	//	return fmt.Errorf("Invalid context release behavior 0x%08X", ctxconfig.release)
	// }
	return nil
}

func getWindowStyle(window *_GLFWwindow) uint32 {
	var style uint32 = WS_CLIPSIBLINGS | WS_CLIPCHILDREN
	if window.monitor != nil {
		style |= WS_POPUP
	} else {
		style |= WS_SYSMENU | WS_MINIMIZEBOX
	}
	if window.decorated {
		style |= WS_CAPTION
	}
	if window.resizable {
		style |= WS_MAXIMIZEBOX | WS_THICKFRAME
	} else {
		style |= WS_POPUP
	}
	return style
}

func getWindowExStyle(w *_GLFWwindow) uint32 {
	var style uint32 = WS_EX_APPWINDOW
	if w.monitor != nil || w.floating {
		style |= WS_EX_TOPMOST
	}
	return style
}

func _glfwRegisterWindowClassWin32() error {
	/*var wc WNDCLASSEXW
	wc.cbSize        = sizeof(wc);
	wc.style         = CS_HREDRAW | CS_VREDRAW | CS_OWNDC;
	wc.lpfnWndProc   = windowProc;
	wc.hInstance     = _glfw.Win32.instance;
	wc.hCursor       = LoadCursorW(NULL, IDC_ARROW);
	wc.lpszClassName = _GLFW_WNDCLASSNAME;
	// Load user-provided icon if available
	//wc.hIcon = LoadImageW(GetModuleHandleW(NULL),"GLFW_ICON", IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);
	//if (!wc.hIcon) {
		// No user-provided icon found, load default icon
		//wc.hIcon = LoadImageW(NULL,	IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);
	//}*/
	icon := syscall.Handle(0)
	wcls := WndClassEx{
		CbSize:        uint32(unsafe.Sizeof(WndClassEx{})),
		Style:         CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
		LpfnWndProc:   syscall.NewCallback(windowProc),
		HInstance:     _glfw.instance,
		HIcon:         icon,
		LpszClassName: syscall.StringToUTF16Ptr("GLFW"),
	}
	var err error
	_glfw.class, err = RegisterClassEx(&wcls)
	return err
}

func createNativeWindow(window *_GLFWwindow, wndconfig *_GLFWwndconfig, fbconfig *_GLFWfbconfig) error {
	var err error
	var frameX, frameY, frameWidth, frameHeight int32
	style := getWindowStyle(window)
	exStyle := getWindowExStyle(window)

	if _glfw.win32.mainWindowClass == 0 {
		err = _glfwRegisterWindowClassWin32()
		if err != nil {
			panic(err)
		}
		_glfw.win32.mainWindowClass = _glfw.class
	}
	if window.monitor != nil {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))
		_, _, err := _GetMonitorInfo.Call(uintptr(window.monitor.hMonitor), uintptr(unsafe.Pointer(&mi)))
		if errors.Is(err, syscall.Errno(0)) {
			return err
		}
		// NOTE: This window placement is temporary and approximate, as the
		//       correct position and size cannot be known until the monitor
		//       video mode has been picked in _glfwSetVideoModeWin32
		frameX = mi.RcMonitor.Left
		frameY = mi.RcMonitor.Top
		frameWidth = mi.RcMonitor.Right - mi.RcMonitor.Left
		frameHeight = mi.RcMonitor.Bottom - mi.RcMonitor.Top
	} else {
		rect := RECT{0, 0, int32(wndconfig.width), int32(wndconfig.height)}
		window.Win32.maximized = wndconfig.maximized
		if wndconfig.maximized {
			style |= WS_MAXIMIZE
		}
		// TODO AdjustWindowRectEx(&rect, style, FALSE, exStyle);
		frameX = CW_USEDEFAULT
		frameY = CW_USEDEFAULT
		frameWidth = rect.Right - rect.Left
		frameHeight = rect.Bottom - rect.Top
	}

	window.Win32.handle, err = CreateWindowEx(
		exStyle,
		_glfw.class,
		wndconfig.title,
		style,
		frameX, frameY,
		frameWidth, frameHeight,
		0, // No parent
		0, // No menu
		resources.handle,
		uintptr(unsafe.Pointer(wndconfig)))
	setProp(window.Win32.handle, window)
	return err
}

type PIXELFORMATDESCRIPTOR = struct {
	nSize           uint16
	nVersion        uint16
	dwFlags         uint32
	iPixelType      uint8
	cColorBits      uint8
	cRedBits        uint8
	cRedShift       uint8
	cGreenBits      uint8
	cGreenShift     uint8
	cBlueBits       uint8
	cBlueShift      uint8
	cAlphaBits      uint8
	cAlphaShift     uint8
	cAccumBits      uint8
	cAccumRedBits   uint8
	cAccumGreenBits uint8
	cAccumBlueBits  uint8
	cAccumAlphaBits uint8
	cDepthBits      uint8
	cStencilBits    uint8
	cAuxBuffers     uint8
	iLayerType      uint8
	bReserved       uint8
	dwLayerMask     uint32
	dwVisibleMask   uint32
	dwDamageMask    uint32
}

var (
	gdi32          = windows.NewLazySystemDLL("gdi32.dll")
	_GetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
	_CreateDC      = gdi32.NewProc("CreateDCW")
	_DeleteDC      = gdi32.NewProc("DeleteDC")

	ntdll                 = windows.NewLazySystemDLL("ntdll.dll")
	_RtlVerifyVersionInfo = ntdll.NewProc("RtlVerifyVersionInfo")
)

const (
	PFD_DRAW_TO_WINDOW = 0x04
	PFD_SUPPORT_OPENGL = 0x20
	PFD_DOUBLEBUFFER   = 0x01
	PFD_TYPE_RGBA      = 0x00
)

func glfwTerminate() {
	/* TODO
	   if (_glfw.Win32.deviceNotificationHandle) {
	   	UnregisterDeviceNotification(_glfw.Win32.deviceNotificationHandle);
	   }
	*/
	if _glfw.win32.helperWindowHandle != 0 {
		_, _, err := _DestroyWindow.Call(uintptr(_glfw.win32.helperWindowHandle))
		if !errors.Is(err, syscall.Errno(0)) {
			slog.Error("UnregisterClass failed, " + err.Error())
		}
	}
	if _glfw.win32.helperWindowClass != 0 {
		_, _, err := _UnregisterClass.Call(uintptr(_glfw.win32.helperWindowClass), uintptr(_glfw.win32.instance))
		if !errors.Is(err, syscall.Errno(0)) {
			slog.Error("UnregisterClass failed, " + err.Error())
		}
	}
	if _glfw.win32.mainWindowClass != 0 {
		_, _, err := _UnregisterClass.Call(uintptr(_glfw.win32.mainWindowClass), uintptr(_glfw.win32.instance))
		if !errors.Is(err, syscall.Errno(0)) {
			slog.Error("UnregisterClass failed, " + err.Error())
		}
	}
}

func glfwPlatformInit() error {
	var err error
	createKeyTables()
	if isWindows10Version1703OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if isWindows8Point1OrGreater() {
		_, _, err := _SetProcessDpiAwarenessContext.Call(uintptr(PROCESS_PER_MONITOR_DPI_AWARE))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("SetProcessDpiAwarenessContext failed, " + err.Error())
		}
	} else if IsWindowsVistaOrGreater() {
		_, _, _ = _SetProcessDPIAware.Call()
	}

	/* This is not in C version
	if err := _glfwRegisterWindowClassWin32(); err != nil {
		return fmt.Errorf("glfw platform init failed, _glfwRegisterWindowClassWin32 failed, %v ", err.Error())
	}*/
	// _, _, err := _procGetModuleHandleExW.Call(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS|GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT, uintptr(unsafe.Pointer(&_glfw)), uintptr(unsafe.Pointer(&_glfw.instance)))

	_glfw.instance, err = GetModuleHandle()
	if err != nil {
		return fmt.Errorf("glfw platform init failed %v ", err.Error())
	}

	err = createHelperWindow()
	if err != nil {
		return err
	}
	glfwPollMonitorsWin32()
	// TODO? _glfwPlatformSetTls(&_glfw.errorSlot, &_glfwMainThreadError)
	glfwDefaultWindowHints()
	_glfw.initialized = true
	return nil
}

func glfwPlatformCreateWindow(window *_GLFWwindow, wndconfig *_GLFWwndconfig, ctxconfig *_GLFWctxconfig, fbconfig *_GLFWfbconfig) error {
	err := createNativeWindow(window, wndconfig, fbconfig)
	if err != nil {
		return err
	}
	if ctxconfig.client != GLFW_NO_API {
		if ctxconfig.source == GLFW_NATIVE_CONTEXT_API {
			if err := _glfwInitWGL(); err != nil {
				return fmt.Errorf("_glfwInitWGL error " + err.Error())
			}
			if err := glfwCreateContextWGL(window, ctxconfig, fbconfig); err != nil {
				return fmt.Errorf("glfwCreateContextWGL error " + err.Error())
			}
		} else if ctxconfig.source == GLFW_EGL_CONTEXT_API {
			if err := glfwInitEGL(); err != nil {
				return err
			}
			if err := _glfwCreateContextEGL(window, ctxconfig, fbconfig); err != nil {
				return err
			}
		} else if ctxconfig.source == GLFW_OSMESA_CONTEXT_API {
			if err := glfwInitOSMesa(); err != nil {
				return err
			}
			if err := _glfwCreateContextOSMesa(window, ctxconfig, fbconfig); err != nil {
				return err
			}
		}
		if err := _glfwRefreshContextAttribs(window, ctxconfig); err != nil {
			return err
		}
	}
	if window.monitor != nil {
		glfwShowWindow(window)
		glfwFocusWindow(window)
		// acquireMonitor(window)
		// fitToMonitor(window)
		if wndconfig.centerCursor {
			// _glfwCenterCursorInContentArea(window)
		}
	} else if wndconfig.visible {
		glfwShowWindow(window)
		if wndconfig.focused {
			glfwFocusWindow(window)
		}
	}
	return nil
}

func glfwCreateWindow(width, height int, title string, monitor *Monitor, share *_GLFWwindow) (*_GLFWwindow, error) {

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid width/heigth")
	}
	// End of _glfwPlatformCreateWindow
	fbconfig := _glfw.hints.framebuffer
	ctxconfig := _glfw.hints.context
	wndconfig := _glfw.hints.window
	wndconfig.width = width
	wndconfig.height = height

	wndconfig.title = title
	ctxconfig.share = share
	if glfwIsValidContextConfig(&ctxconfig) != nil {
		return nil, fmt.Errorf("glfw context config is invalid: %v", ctxconfig)
	}

	window := &_GLFWwindow{}
	window.next = _glfw.windowListHead
	_glfw.windowListHead = window

	window.videoMode.width = width
	window.videoMode.height = height
	window.videoMode.redBits = fbconfig.redBits
	window.videoMode.greenBits = fbconfig.greenBits
	window.videoMode.blueBits = fbconfig.blueBits
	window.videoMode.refreshRate = _glfw.hints.refreshRate

	window.monitor = monitor
	window.resizable = wndconfig.resizable
	window.decorated = wndconfig.decorated
	window.autoIconify = wndconfig.autoIconify
	window.floating = wndconfig.floating
	window.focusOnShow = wndconfig.focusOnShow
	window.cursorMode = GLFW_CURSOR_NORMAL
	window.doublebuffer = fbconfig.doublebuffer
	window.minwidth = GLFW_DONT_CARE
	window.minheight = GLFW_DONT_CARE
	window.maxwidth = GLFW_DONT_CARE
	window.maxheight = GLFW_DONT_CARE
	window.numer = GLFW_DONT_CARE
	window.denom = GLFW_DONT_CARE

	if err := glfwPlatformCreateWindow(window, &wndconfig, &ctxconfig, &fbconfig); err != nil {
		// glfwDestroyWindow(window)
		return nil, fmt.Errorf("Error creating window, " + err.Error())
	}
	return window, nil
}

func glfwDefaultWindowHints() {
	_glfw.hints.context.client = GLFW_OPENGL_API
	_glfw.hints.context.source = GLFW_NATIVE_CONTEXT_API
	_glfw.hints.context.major = 1
	_glfw.hints.context.minor = 0
	// The default is a focused, visible, resizable window with decorations
	_glfw.hints.window.resizable = true
	_glfw.hints.window.visible = true
	_glfw.hints.window.decorated = true
	_glfw.hints.window.focused = true
	_glfw.hints.window.autoIconify = true
	_glfw.hints.window.centerCursor = true
	_glfw.hints.window.focusOnShow = true
	// The default is 24 bits of color, 24 bits of depth and 8 bits of stencil, double buffered
	_glfw.hints.framebuffer.redBits = 8
	_glfw.hints.framebuffer.greenBits = 8
	_glfw.hints.framebuffer.blueBits = 8
	_glfw.hints.framebuffer.alphaBits = 8
	_glfw.hints.framebuffer.depthBits = 24
	_glfw.hints.framebuffer.stencilBits = 8
	_glfw.hints.framebuffer.doublebuffer = true
	// The default is to select the highest available refresh rate
	_glfw.hints.refreshRate = GLFW_DONT_CARE
	// The default is to use full Retina resolution framebuffers
	_glfw.hints.window.ns.retina = true
}

func helperWindowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	/*	switch msg	{
		//case WM_DISPLAYCHANGE:
		//    _glfwPollMonitorsWin32();

		case WM_DEVICECHANGE:
		if (!_glfw.joysticksInitialized)
				break;

		if (wParam == DBT_DEVICEARRIVAL)
		{
		DEV_BROADCAST_HDR* dbh = (DEV_BROADCAST_HDR*) lParam;
		if (dbh && dbh->dbch_devicetype == DBT_DEVTYP_DEVICEINTERFACE)
		_glfwDetectJoystickConnectionWin32();
		}
		else if (wParam == DBT_DEVICEREMOVECOMPLETE)
		{
		DEV_BROADCAST_HDR* dbh = (DEV_BROADCAST_HDR*) lParam;
		if (dbh && dbh->dbch_devicetype == DBT_DEVTYP_DEVICEINTERFACE)
		_glfwDetectJoystickDisconnectionWin32();
		}

		break;
		}
		}
	*/
	r1, _, _ := _DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r1
}

func glfwShowWindow(w *_GLFWwindow) {
	mode := windows.SW_NORMAL
	if w.Win32.iconified {
		mode = windows.SW_MINIMIZE
	} else if w.Win32.maximized {
		mode = windows.SW_MAXIMIZE
	}
	_, _, err := _ShowWindow.Call(uintptr(w.Win32.handle), uintptr(mode))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("ShowWindow failed, " + err.Error())
	}
}

func createHelperWindow() error {
	var err error
	var wc WndClassEx
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = CS_OWNDC
	wc.LpfnWndProc = syscall.NewCallback(helperWindowProc)
	wc.HInstance = _glfw.instance
	wc.LpszClassName = syscall.StringToUTF16Ptr("GLFW3 Helper")

	_glfw.win32.helperWindowClass, err = RegisterClassEx(&wc)
	if _glfw.win32.helperWindowClass == 0 || err != nil {
		panic("Win32: Failed to register helper window class")
	}
	_glfw.win32.helperWindowHandle, err =
		CreateWindowEx(WS_OVERLAPPED,
			_glfw.win32.helperWindowClass,
			"Helper window",
			WS_OVERLAPPED|WS_CLIPSIBLINGS|WS_CLIPCHILDREN,
			0, 0, 500, 500,
			0, 0,
			resources.handle,
			0)

	if _glfw.win32.helperWindowHandle == 0 || err != nil {
		panic("Win32: Failed to create helper window")
	}
	_, _, err = _ShowWindow.Call(uintptr(_glfw.win32.helperWindowHandle), windows.SW_HIDE)
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		return err
	}

	// TODO Register for HID device notifications
	/*
		{
			dbi DEV_BROADCAST_DEVICEINTERFACE_W
			ZeroMemory(&dbi, sizeof(dbi));
			dbi.dbcc_size = sizeof(dbi);
			dbi.dbcc_devicetype = DBT_DEVTYP_DEVICEINTERFACE;
			dbi.dbcc_classguid = GUID_DEVINTERFACE_HID;

			_glfw.win32.deviceNotificationHandle =
				RegisterDeviceNotificationW(_glfw.win32.helperWindowHandle,
					(DEV_BROADCAST_HDR*) &dbi,
					DEVICE_NOTIFY_WINDOW_HANDLE);
		}
	*/
	var msg Msg
	for PeekMessage(&msg, _glfw.win32.helperWindowHandle, 0, 0, PM_REMOVE) {
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
	return nil
}

const (
	DPI_AWARENESS_CONTEXT_UNAWARE              = 0xFFFFFFFFFFFFFFFF
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = 0xFFFFFFFFFFFFFFFE
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = 0xFFFFFFFFFFFFFFFD
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = 0xFFFFFFFFFFFFFFFC
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = 0xFFFFFFFFFFFFFFFB
	PROCESS_DPI_UNAWARE                        = 0
	PROCESS_SYSTEM_DPI_AWARE                   = 1
	PROCESS_PER_MONITOR_DPI_AWARE              = 2
	VER_MAJORVERSION                           = 0x0000002
	VER_MINORVERSION                           = 0x0000001
	VER_BUILDNUMBER                            = 0x0000004
	VER_SERVICEPACKMAJOR                       = 0x00000020
	WIN32_WINNT_WINBLUE                        = 0x0603
)

type _OSVERSIONINFOEXW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        uint8
	wReserved           uint8
}

func glfwIsWindows10Version1607OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 14393
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func isWindows10Version1703OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = 10
	osvi.dwMinorVersion = 0
	osvi.dwBuildNumber = 15063
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_BUILDNUMBER
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x80000000000000db))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func isWindows8Point1OrGreater() bool {
	var osvi _OSVERSIONINFOEXW
	osvi.dwOSVersionInfoSize = uint32(unsafe.Sizeof(osvi))
	osvi.dwMajorVersion = uint32(WIN32_WINNT_WINBLUE >> 8)
	osvi.dwMinorVersion = uint32(WIN32_WINNT_WINBLUE & 0xFF)
	osvi.wServicePackMajor = 0
	var mask uint32 = VER_MAJORVERSION | VER_MINORVERSION | VER_SERVICEPACKMAJOR
	// ULONGLONG cond = VerSetConditionMask(0, VER_MAJORVERSION, VER_GREATER_EQUAL);
	// cond = VerSetConditionMask(cond, VER_MINORVERSION, VER_GREATER_EQUAL);
	// cond = VerSetConditionMask(cond, VER_SERVICEPACKMAJOR, VER_GREATER_EQUAL);
	r, _, err := _RtlVerifyVersionInfo.Call(uintptr(unsafe.Pointer(&osvi)), uintptr(mask), uintptr(0x800000000001801b))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("SetProcessDpiAwarenessContext failed, " + err.Error())
	}
	return r == 0
}

func IsWindowsVistaOrGreater() bool {
	return true
}

func glfwGetWindowFrameSize(window *_GLFWwindow, left, top, right, bottom *int) {
	var rect RECT
	var width, height int
	glfwGetWindowSize(window, &width, &height)
	rect.Right = int32(width)
	rect.Bottom = int32(height)
	if glfwIsWindows10Version1607OrGreater() {
		// AdjustWindowRectExForDpi(&rect, getWindowStyle(window),	FALSE, getWindowExStyle(window),GetDpiForWindow(window->win32.hMonitor));
	} else {
		// AdjustWindowRectEx(&rect, getWindowStyle(window),FALSE, getWindowExStyle(window));
	}
	*left = int(-rect.Left)
	*top = int(-rect.Top)
	*right = int(rect.Right) - width
	*bottom = int(rect.Bottom) - height
}

func screenToClient(handle syscall.Handle, p *POINT) {
	_, _, err := _ScreenToClient.Call(uintptr(handle), uintptr(unsafe.Pointer(p)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic("GetCursorPos failed, " + err.Error())
	}
}

func glfwGetCursorPos(w *_GLFWwindow, x *int, y *int) {
	if w.cursorMode == GLFW_CURSOR_DISABLED {
		*x = int(w.virtualCursorPosX)
		*y = int(w.virtualCursorPosY)
	} else {
		var pos POINT
		_, _, err := _GetCursorPos.Call(uintptr(unsafe.Pointer(&pos)))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("GetCursorPos failed, " + err.Error())
		}
		screenToClient(w.Win32.handle, &pos)
		*x = int(pos.X)
		*y = int(pos.Y)
	}
}

func glfwGetWindowSize(window *_GLFWwindow, width *int, height *int) {
	var area RECT
	_, _, err := _GetClientRect.Call(uintptr(unsafe.Pointer(window.Win32.handle)), uintptr(unsafe.Pointer(&area)))
	if !errors.Is(err, syscall.Errno(0)) {
		panic(err)
	}
	// GetClientRect(window->win32.hMonitor, &area);
	*width = int(area.Right)
	*height = int(area.Bottom)
}

// GetClipboardString returns the contents of the system clipboard, if it
// contains or is convertible to a UTF-8 encoded string.
// This function may only be called from the main thread.
func glfwGetClipboardString() string {
	b := clipboard.Read(clipboard.FmtText)
	return string(b)
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
// This function may only be called from the main thread.
func glfwSetClipboardString(str string) {
	clipboard.Write(clipboard.FmtText, []byte(str))
}

func glfwCreateStandardCursorWin32(cursor *Cursor, shape int) {
	var id uint16
	switch shape {
	case ArrowCursor:
		id = IDC_ARROW
	case IbeamCursor:
		id = IDC_IBEAM
	case CrosshairCursor:
		id = IDC_CROSS
	case HResizeCursor:
		id = IDC_SIZEWE
	case VResizeCursor:
		id = IDC_SIZENS
	case HandCursor:
		id = IDC_HAND
	default:
		panic("Win32: Unknown or unsupported standard cursor")
	}
	cursor.handle = LoadCursor(id)
	if cursor.handle == 0 {
		panic("Win32: Failed to create standard cursor")
	}
}

func monitorFromWindow(handle syscall.Handle, flags uint32) syscall.Handle {
	r1, _, err := _MonitorFromWindow.Call(uintptr(handle), uintptr(flags))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("MonitorFromWindow failed, " + err.Error())
	}
	return syscall.Handle(r1)
}

func glfwGetContentScale(w *Window) (float32, float32) {
	var xscale, yscale float32
	var xdpi, ydpi int
	handle := monitorFromWindow(w.Win32.handle, MONITOR_DEFAULTTONEAREST)
	if isWindows8Point1OrGreater() {
		_, _, err := _GetDpiForMonitor.Call(uintptr(handle), uintptr(0),
			uintptr(unsafe.Pointer(&xdpi)), uintptr(unsafe.Pointer(&ydpi)))
		if !errors.Is(err, syscall.Errno(0)) {
			panic("GetDpiForMonitor failed, " + err.Error())
		}
	} else {
		dc := getDC(0)
		xdpi = GetDeviceCaps(dc, LOGPIXELSX)
		ydpi = GetDeviceCaps(dc, LOGPIXELSY)
		releaseDC(0, dc)
	}
	xscale = float32(xdpi) / USER_DEFAULT_SCREEN_DPI
	yscale = float32(ydpi) / USER_DEFAULT_SCREEN_DPI
	return xscale, yscale
}

func glfwSetWindowPos(window *_GLFWwindow, xpos, ypos int) {
	// SetWindowPos(window.Win32.handle, 0, xpos, ypos, 0, 0, SWP_NOACTIVATE|SWP_NOZORDER|SWP_NOSIZE)
	r1, _, err := _SetWindowPos.Call(uintptr(window.Win32.handle), uintptr(0), uintptr(xpos), uintptr(ypos), 0, 0, uintptr(SWP_NOACTIVATE|SWP_NOZORDER|SWP_NOSIZE))
	if err != nil && !errors.Is(err, syscall.Errno(0)) {
		panic("SetWindowPos failed, " + err.Error())
	}
	if r1 == 0 {
		panic("SetWindowPos failed")
	}
}
