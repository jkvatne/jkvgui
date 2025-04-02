package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/img"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"runtime/debug"
)

func GetInfo() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			key := setting.Key
			if key == "vcs.revision" {
				slog.Info("Revision", setting.Value[:8])
			}
		}
	}
}

var CardName string
var CardList = []string{"RRADI16", "RRAIO16", "RRDIO15", "RRPT8", "RRLC2", "RREPS3"}
var Images []*img.Img

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("IO-Card test", wid.H1C),
		wid.Combo(&CardName, CardList, "Select card to test", nil),
		wid.Row(nil,
			img.W(Images[0], img.SCALE, ""),
		),
	)
}

func Display(form wid.Wid) {
	ctx := wid.NewCtx()
	// First measure widgets
	form(ctx)
	// THen do drawing
	ctx.Draw = true
	form(ctx)
}

func main() {
	GetInfo()

	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()

	sys.Initialize(window, 14)
	im, _ := img.New("rradi16.jpg")
	Images = append(Images, im)
	gpu.UserScale = 1.5

	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
