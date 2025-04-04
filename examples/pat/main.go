package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"runtime/debug"
	"strings"
)

var (
	tag  = "(developement build)"
	url  = "(developement build)"
	hash = "(developement build)"
)

func GetInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		slog.Error("Could not read build info")
		return
	}
	s := info.Main.Version
	if s != "" {
		words := strings.Split(s, "-")
		tag = words[0]
	}
	url = info.Main.Path
	for _, setting := range info.Settings {
		key := setting.Key
		if key == "vcs.revision" {
			hash = setting.Value[:8]
		}
	}
	slog.Info("Buildinfo", "hash", hash, "tag", tag, "url", url)
}

var CardName string
var CardList = []string{"RRADI16", "RRAIO16", "RRDIO15", "RRPT8", "RRLC2", "RREPS3"}
var Images []*wid.Img

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("IO-Card Production Acceptance Test", wid.H1C),
		wid.Row(nil,
			wid.Image(Images[0], wid.DefaultImgStyle.W(0.5), ""),
			wid.Col(wid.ContStyle.W(0.5),
				wid.Combo(&CardName, CardList, "Select card to test", nil),
				wid.Edit(&CardName, "Card", nil, nil),
			),
		),
	)
}

func main() {
	GetInfo()

	window := gpu.InitWindow(0, 0, "IO-Card PAT", 2)
	defer gpu.Shutdown()

	sys.Initialize(window, 14)
	im, _ := wid.New("rradi16.jpg")
	Images = append(Images, im)
	gpu.UserScale = 1.5

	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		ctx := wid.NewCtx()
		Form()(ctx)
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
