package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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
var Value1 = "Value1"
var Value2 = "Value2"
var Value3 = "Value3"

var CardList = []string{"RRADI16", "RRAIO16", "RRDIO15", "RRPT8", "RRLC2", "RREPS3"}
var Images []*wid.Img
var logText []string

/*{"1 azdsfadf", "2 azdsfadf", "3 azdsfadf", "4 azdsfadf", "5 azdsfadf", "6 azdsfadf",
	"7 azdsfadf", "8 azdsfadf", "9 azdsfadf", "10 azdsfadf", "11 azdsfadf", "12 azdsfadf", "13 azdsfadf",
	"14 azdsfadf", "15 azdsfadf", "16 azdsfadf", "17 azdsfadf", "18 azdsfadf", "19 azdsfadf", "20 azdsfadf",
	"14 azdsfadf", "15 azdsfadf", "16 azdsfadf", "17 azdsfadf", "18 azdsfadf", "19 azdsfadf", "20 azdsfadf",
	"14 azdsfadf", "15 azdsfadf", "16 azdsfadf", "17 azdsfadf", "18 azdsfadf", "19 azdsfadf", "20 azdsfadf",
}

*/
var lenstr string

func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("IO-Card Production Acceptance Test", wid.H1C),
		wid.Row(nil,
			wid.Image(Images[0], wid.DefaultImgStyle.W(0.5), ""),
			wid.Col(wid.ContStyle.W(0.5),
				wid.Combo(&CardName, CardList, "Select card to test", nil),
				wid.Edit(&Value1, "Value1", nil, nil),
				wid.Edit(&Value2, "Value2", nil, nil),
				wid.Edit(&Value3, "Value3", nil, nil),
				wid.Label(lenstr, nil),
			),
		),
		wid.Memo(&logText, nil),
	)
}

func main() {
	GetInfo()
	// User scale is the zoom factor set by ctrl+Scroll wheel. Used to magnify fonts and widgets.
	gpu.UserScale = 1.5
	// Set DebugWidgets=true to draw rectangles showning widget sizes
	gpu.DebugWidgets = false
	window := gpu.InitWindow(0, 0, "IO-Card PAT", 2)
	defer gpu.Shutdown()

	sys.Initialize(window, 14)
	im, _ := wid.NewImage("rradi16.jpg")
	Images = append(Images, im)
	for _ = range 12 {
		logText = append(logText, "gggTTT qrtpåæØÆ asdfasdfasdfa asd adsf "+strconv.Itoa(len(logText)))
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			logText = append(logText, "gggTTT qrtpåæØÆ asdfasdfasdfa asd adsf  "+strconv.Itoa(len(logText)))
			lenstr = strconv.Itoa(len(logText))
			gpu.Invalidate(0)
		}
	}()
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		ctx := wid.NewCtx()
		Form()(ctx)
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
