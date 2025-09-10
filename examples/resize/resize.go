package main

import (
	"log/slog"
	"os"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var posV [16]wid.ResizerState
var posH [16]wid.ResizerState
var image [16]*wid.Img

func main() {
	// Configure slog to print date and time in standard format
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
		}
		return a
	}})))
	slog.Info("Starting Resize demo")

	sys.Init()
	defer sys.Shutdown()
	sys.CreateWindow(100, 100, 500, 400, "Resizing1", 1, 1)
	sys.CreateWindow(200, 200, 500, 400, "Resizing2", 2, 1)
	image[0], _ = wid.NewImage("music.jpg")
	image[1], _ = wid.NewImage("ts.jpg")
	for sys.Running() {
		for sys.CurrentWno, _ = range sys.WindowList {
			sys.StartFrame(theme.Surface.Bg())
			ctx := wid.NewCtx()
			wid.HorResizer(
				&posH[sys.CurrentWno], nil,
				wid.Image(image[sys.CurrentWno], nil, ""),
				wid.VertResizer(&posV[sys.CurrentWno], nil,
					wid.Btn("Left", nil, func() {}, nil, ""),
					wid.Btn("Right", nil, func() {}, nil, ""),
				),
			)(ctx)
			// EndFrame will swap buffers and limit the maximum framerate.
			sys.EndFrame()
		}
		sys.PollEvents()
	}
}
