package mouse

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"log/slog"
	"time"
)

var (
	mousePos        f32.Pos
	leftBtnDown     bool
	leftBtnReleased bool
)

func Pos() f32.Pos {
	return mousePos
}

func Hovered(r f32.Rect) bool {
	return mousePos.Inside(r)
}

func PosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	mousePos.X = float32(xpos) / gpu.ScaleX
	mousePos.Y = float32(ypos) / gpu.ScaleY
	gpu.Invalidate(50 * time.Millisecond)
}

func LeftBtnPressed(r f32.Rect) bool {
	return mousePos.Inside(r) && leftBtnDown
}

func LeftBtnReleased(r f32.Rect) bool {
	if mousePos.Inside(r) && leftBtnReleased {
		leftBtnReleased = false
		return true
	}
	return false
}

func BtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	gpu.Invalidate(0)
	x, y := w.GetCursorPos()
	mousePos.X = float32(x) / gpu.ScaleX
	mousePos.Y = float32(y) / gpu.ScaleY
	slog.Debug("Mouse click:", "Button", button, "X", x, "Y", y, "Action", action)
	if action == glfw.Release {
		leftBtnDown = false
		leftBtnReleased = true
		for _, clickable := range gpu.Clickables {
			if mousePos.Inside(clickable.Rect) {
				if clickable.Action != nil {
					if f, ok := clickable.Action.(func()); ok {
						f()
					}
				}
			}
		}
	} else if action == glfw.Press {
		leftBtnDown = true
	}
}
