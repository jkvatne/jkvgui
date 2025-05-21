package main

import (
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"strconv"
	"time"
)

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
	go func() {
		time.Sleep(time.Second)
		for {
			if len(logText) < 13 {
				time.Sleep(1 * time.Second / 6)
			} else if len(logText) < 25 {
				time.Sleep(2 * time.Second)
			} else {
				time.Sleep(5 * time.Second)
			}
			gpu.Mutex.Lock()
			logText = append(logText, strconv.Itoa(len(logText))+
				" Some text with special characters æøåÆØÅ$€ÆØÅ and some more arbitary text to make a very long line that will be broken for wrap-around (or elipsis)")
			gpu.Mutex.Unlock()
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
				wid.Label("FPS="+strconv.Itoa(sys.RedrawsPrSec), nil),
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
	sys.Initialize()
	window := gpu.InitWindow(0, 0, "IO-Card PAT", 2, 1.5)
	defer sys.Shutdown()
	sys.InitializeWindow(window)
	img, _ := wid.NewImage("rradi16.jpg")
	Images = append(Images, img)
	slog.Info("Pat.exe is running4")
	DummyLogGenerator()
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface)
		Form()(wid.NewCtx())
		sys.EndFrame(2)
	}
}
