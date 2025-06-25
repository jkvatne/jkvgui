package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	// Create a window with a title and size.
	sys.CreateWindow(f32.Rect{X: 100, Y: 100, W: 200, H: 200},
		"Hello world", 0, 2)
	defer sys.Shutdown()
	// Loop until the window is closed.
	for sys.Running(0) {
		sys.StartFrame(0, theme.Surface.Bg())
		// Show just a single widget and call it with a new Ctx.
		wid.Label("Hello world!", nil)(wid.NewCtx(0))
		// EndFrame will swap buffers and limit the maximum framerate.
		sys.EndFrame(0)
		sys.PollEvents()
	}
}
