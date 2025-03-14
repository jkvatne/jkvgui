package focus

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"log/slog"
	"time"
)

type Clickable struct {
	Rect   f32.Rect
	Action func()
}

var (
	Current          interface{}
	MoveToNext       bool
	MoveToPrevious   bool
	ToNext           bool
	Last             interface{}
	Clickables       []Clickable
	MousePos         f32.Pos
	MouseBtnDown     bool
	MouseBtnReleased bool
)

func Update() {

}

func At(tag interface{}) bool {
	return lib.TagsEqual(tag, Current)
}

func Set(action interface{}) {
	Current = action
}

func Move(action interface{}) {
	if MoveToPrevious && lib.TagsEqual(action, Current) {
		Current = Last
		MoveToPrevious = false
		gpu.Invalidate(0)
	}
	if ToNext {
		ToNext = false
		Current = action
		gpu.Invalidate(0)
	}
	if lib.TagsEqual(action, Current) {
		if MoveToNext {
			ToNext = true
			MoveToNext = false
		}
		gpu.Invalidate(0)
	}
}

func AddFocusable(rect f32.Rect, action func()) {
	Last = action
	Clickables = append(Clickables, Clickable{Rect: rect, Action: action})
}

// Mouse

func Hovered(r f32.Rect) bool {
	return MousePos.Inside(r)
}

func MousePosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	MousePos.X = float32(xpos) / gpu.ScaleX
	MousePos.Y = float32(ypos) / gpu.ScaleY
	gpu.Invalidate(50 * time.Millisecond)
}

func LeftMouseBtnPressed(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnDown
}

func LeftMouseBtnReleased(r f32.Rect) bool {
	return MousePos.Inside(r) && MouseBtnReleased
}

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	gpu.Invalidate(0)
	x, y := w.GetCursorPos()
	MousePos.X = float32(x) / gpu.ScaleX
	MousePos.Y = float32(y) / gpu.ScaleY
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
