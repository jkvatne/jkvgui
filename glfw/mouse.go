package glfw

// Mouse buttons.
const (
	MouseButton1      MouseButton = 0
	MouseButton2      MouseButton = 1
	MouseButton3      MouseButton = 2
	MouseButton4      MouseButton = 3
	MouseButton5      MouseButton = 4
	MouseButton6      MouseButton = 5
	MouseButton7      MouseButton = 6
	MouseButton8      MouseButton = 7
	MouseButtonLast   MouseButton = 7
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
)

type MouseButton int

// MouseButtonCallback is the mouse button callback.
type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

// SetMouseButtonCallback sets the mouse button callback which is called when a
// mouse button is pressed or released.
//
// When a window loses focus, it will generate synthetic mouse button release
// events for all pressed mouse buttons. You can tell these events from
// user-generated events by the fact that the synthetic ones are generated after
// the window has lost focus, i.e. Focused will be false and the focus
// callback will have already been called.
func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	/*previous = w.fMouseButtonHolder
	w.fMouseButtonHolder = cbfun
	if cbfun == nil {
		// C.glfwSetMouseButtonCallback(w.data, nil)
	} else {
		// C.glfwSetMouseButtonCallbackCB(w.data)
	}
	panicError()
	return previous */
	return nil
}

// ScrollCallback is the scroll callback.
type ScrollCallback func(w *Window, xoff float64, yoff float64)

// SetScrollCallback sets the scroll callback which is called when a scrolling
// device is used, such as a mouse wheel or scrolling area of a touchpad.
func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	/*previous = w.fScrollHolder
	w.fScrollHolder = cbfun
	if cbfun == nil {
		// C.glfwSetScrollCallback(w.data, nil)
	} else {
		// C.glfwSetScrollCallbackCB(w.data)
	}
	panicError()
	return previous */
	return nil
}
