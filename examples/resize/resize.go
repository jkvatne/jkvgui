package main

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var pos1 wid.ResizerState
var pos2 wid.ResizerState

func main() {
	sys.CreateWindow(f32.Rect{100, 100, 400, 400}, "Resizing1", 1, 2)
	// sys.CreateWindow(f32.Rect{400, 400, 400, 400}, "Resizing2", 2, 2)
	defer sys.Shutdown()
	image, _ := wid.NewImage("music.jpg")
	for {
		for wno, _ := range sys.WindowList {
			sys.StartFrame(wno, theme.Surface.Bg())
			ctx := wid.NewCtx(wno)
			wid.HorResizer(
				&pos1, nil,
				wid.Image(image, nil, ""),
				wid.VertResizer(&pos2, nil,
					wid.Btn("Left", nil, func() {}, nil, ""),
					wid.Btn("Right", nil, func() {}, nil, ""),
				),
			)(ctx)
			// EndFrame will swap buffers and limit the maximum framerate.
			sys.EndFrame(wno)
		}
		sys.PollEvents()
	}
}
