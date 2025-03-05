package gpu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/lib"
	"log"
	"unsafe"
)

func ptr(arg interface{}) unsafe.Pointer {
	return (*eface)(unsafe.Pointer(&arg)).val
}

var InFocus interface{}

type eface struct {
	typ, val unsafe.Pointer
}

func Focused(tag interface{}) bool {
	return ptr(tag) == ptr(InFocus)
}

func SetFocus(action interface{}) {
	InFocus = action
}

func MoveFocus(action interface{}) {
	if MoveFocusToPrevious && ptr(action) == ptr(InFocus) {
		InFocus = LastFocusable
		MoveFocusToPrevious = false
	}

	if FocusToNext {
		FocusToNext = false
		InFocus = action
	}
	if ptr(action) == ptr(InFocus) {
		if MoveFocusToNext {
			FocusToNext = true
			MoveFocusToNext = false
		}
	}
}

func Hovered(r lib.Rect) bool {
	return MousePos.Inside(r)
}

func MousePosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	MousePos.X = float32(xpos)
	MousePos.Y = float32(ypos)
}

func LeftMouseBtnPressed(r lib.Rect) bool {
	return MousePos.Inside(r) && MouseBtnDown
}

func LeftMouseBtnReleased(r lib.Rect) bool {
	return MousePos.Inside(r) && MouseBtnReleased
}

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	x, y := w.GetCursorPos()
	MousePos.X = float32(x)
	MousePos.Y = float32(y)
	var pos = lib.Pos{float32(x), float32(y)}
	log.Printf("Mouse btn %d clicked at %0.1f,%0.1f, Action %d\n", button, x, y, action)
	if action == glfw.Release {
		MouseBtnDown = false
		MouseBtnReleased = true
		for _, clickable := range Clickables {
			if pos.Inside(clickable.Rect) {
				clickable.Action()
			}
		}
	} else if action == glfw.Press {
		MouseBtnDown = true
	}
}

func AddFocusable(rect lib.Rect, action func()) {
	LastFocusable = action
	Clickables = append(Clickables, Clickable{Rect: rect, Action: action})
}
