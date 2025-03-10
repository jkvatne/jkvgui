package gpu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/lib"
	"log"
)

var InFocus interface{}

func Focused(tag interface{}) bool {
	return lib.TagsEqual(tag, InFocus)
}

func SetFocus(action interface{}) {
	InFocus = action
}

func MoveFocus(action interface{}) {
	if MoveFocusToPrevious && lib.TagsEqual(action, InFocus) {
		InFocus = LastFocusable
		MoveFocusToPrevious = false
	}

	if FocusToNext {
		FocusToNext = false
		InFocus = action
	}
	if lib.TagsEqual(action, InFocus) {
		if MoveFocusToNext {
			FocusToNext = true
			MoveFocusToNext = false
		}
	}
}

func Hovered(r f32.Rect) bool {
	return MousePos.Inside(r)
}

func MousePosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	MousePos.X = float32(xpos) / Scale
	MousePos.Y = float32(ypos) / Scale
}

func LeftMouseBtnPressed(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnDown
}

func LeftMouseBtnReleased(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnReleased
}

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	x, y := w.GetCursorPos()
	MousePos.X = float32(x) / Scale
	MousePos.Y = float32(y) / Scale
	var pos = f32.Pos{float32(x), float32(y)}
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

func AddFocusable(rect f32.Rect, action func()) {
	LastFocusable = action
	Clickables = append(Clickables, Clickable{Rect: rect, Action: action})
}
