package gpu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/lib"
	"log/slog"
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

func mousePosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	MousePos.X = float32(xpos) / ScaleX
	MousePos.Y = float32(ypos) / ScaleY
}

func LeftMouseBtnPressed(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnDown
}

func LeftMouseBtnReleased(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnReleased
}

func mouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	Invalidate(0)
	x, y := w.GetCursorPos()
	MousePos.X = float32(x) / ScaleX
	MousePos.Y = float32(y) / ScaleY
	slog.Debug("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action)
	if action == glfw.Release {
		MouseBtnDown = false
		MouseBtnReleased = true
		for _, clickable := range Clickables {
			if MousePos.Inside(clickable.Rect) {
				if clickable.Action != nil {
					clickable.Action()
				}
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
