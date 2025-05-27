package main

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	// Create a window with a title and size.
	sys.InitWindow(200, 100, "Hello world", 0, 2)
	defer sys.Shutdown()
	// Loop until the window is closed.
	for sys.Running() {
		sys.StartFrame(theme.Surface.Bg())
		// Show just a single widget and call it with a new Ctx.
		wid.Label("Hello world!", nil)(wid.NewCtx())
		// EndFrame will swap buffers and limit the maximum framerate.
		sys.EndFrame()
	}
}
