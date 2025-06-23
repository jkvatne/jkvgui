package glfw

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

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
	_ReleaseDC                     = user32.NewProc("ReleaseDC")
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
	var window *_GLFWwindow
	for PeekMessage(&msg, 0, 0, 0, PM_REMOVE) {
		if msg.Message == WM_QUIT {
			// NOTE: While GLFW does not itself post WM_QUIT, other processes
			//       may post it to this one, for example Task Manager
			// HACK: Treat WM_QUIT as a close on all windows
			window = _glfw.windowListHead
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
