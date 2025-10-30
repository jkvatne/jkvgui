package main

import (
	"log"
	"log/slog"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

var posV [16]wid.ResizerState
var posH [16]wid.ResizerState
var image [16]*wid.Img

func Form(n int32) wid.Wid {
	return wid.HorResizer(
		&posH[n], nil,
		wid.Image(image[n], nil, nil, "An image"),
		wid.VertResizer(&posV[n], nil,
			wid.Btn("Left", nil, func() {}, nil, ""),
			wid.Btn("Right", nil, func() {}, nil, ""),
		),
	)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	slog.Info("Resize demo")
	sys.Init()
	defer sys.Shutdown()
	sys.CreateWindow(100, 100, 500, 400, "Resizing1", 1, 1)
	sys.CreateWindow(200, 200, 500, 400, "Resizing2", 2, 1)
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for sys.Running() {
		for wno := range sys.WindowCount.Load() {
			sys.WindowList[wno].StartFrame()
			wid.Show(Form(wno))
			sys.WindowList[wno].EndFrame()
		}
		sys.PollEvents()
	}
}
