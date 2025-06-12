package glfw

import "sync"

// Internal window list stuff
type windowList struct {
	l sync.Mutex
	m map[*_GLFWwindow]*Window
}

var windowMap = windowList{m: map[*_GLFWwindow]*Window{}}

func (w *windowList) put(wnd *Window) {
	w.l.Lock()
	defer w.l.Unlock()
	w.m[wnd.Data] = wnd
}

func (w *windowList) remove(wnd *_GLFWwindow) {
	w.l.Lock()
	defer w.l.Unlock()
	delete(w.m, wnd)
}

func (w *windowList) get(wnd *_GLFWwindow) *Window {
	w.l.Lock()
	defer w.l.Unlock()
	return w.m[wnd]
}
