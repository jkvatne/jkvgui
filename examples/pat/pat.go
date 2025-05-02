package main

import (
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

var (
	CardTypeNo int
	CardName   string
	Value1     = "Value1"
	Value2     = "Value2"
	Value3     = "Value3"
	CardList   = []string{"Select card", "RRADI16", "RRAIO16", "RRDIO15", "RRPT8", "RRLC2", "RREPS3"}
	Images     []*wid.Img
	logText    []string
)

func DummyLogGenerator() {
	logText = append(logText, strconv.Itoa(len(logText))+
		" First line")

	go func() {
		for {
			if len(logText) < 13 {
				time.Sleep(1 * time.Second / 6)
			} else {
				time.Sleep(2 * time.Second)
			}
			logText = append(logText, strconv.Itoa(len(logText))+
				" Some text with special characters æøåÆØÅ$€ and some more arbitary text to make a very long line that will be broken for wrap-around (or elipsis)")
			gpu.Invalidate(0)
		}
	}()
}

func Form() wid.Wid {
	// lenstr := fmt.Sprintf("%d", len(logText))
	// cardName := CardList[CardTypeNo]
	return wid.Col(nil,
		wid.Label("IO-Card Production Acceptance Test", wid.H1C),
		wid.Row(nil,
			wid.Image(Images[0], wid.DefImg.W(0.5), ""),
			wid.Col(wid.ContStyle.W(0.5),
				wid.Edit(&Value2, "A long value here", nil, nil),
				/*
					wid.List(&CardTypeNo, CardList, "Select card to test", nil),
					wid.Edit(&CardTypeNo, "CardTypeNo", nil, nil),
					wid.Edit(&Value3, "Value3", nil, nil),
					wid.Label(lenstr, nil),
					wid.Label(cardName, nil),
				*/

			),
		),
		wid.Memo(&logText, nil),
	)
}

func main() {
	GetInfo()
	gpu.UserScale = 1.5 // User scale is the zoom factor set by ctrl+Scroll wheel
	window := gpu.InitWindow(0, 0, "IO-Card PAT", 2, 1.5)
	defer gpu.Shutdown()

	sys.Initialize(window)
	img, _ := wid.NewImage("rradi16.jpg")
	Images = append(Images, img)
	DummyLogGenerator()

	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		sys.EndFrame(20)
	}
}
