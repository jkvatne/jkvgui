package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
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
	win        *sys.Window
)

func getSize() string {
	win.Mutex.Lock()
	defer win.Mutex.Unlock()
	return strconv.Itoa(len(logText) - 1)
}

func addLine(s string) {
	win.Mutex.Lock()
	defer win.Mutex.Unlock()
	logText = append(logText, strconv.Itoa(len(logText))+" "+s)
	sys.Invalidate()
}

func dummyLogGenerator() {
	go func() {
		time.Sleep(time.Second)
		var n int
		for {
			if n < 13 {
				time.Sleep(time.Second / 5)
			} else if n < 25 {
				time.Sleep(time.Second / 5)
			} else {
				time.Sleep(99995 * time.Second)
			}
			addLine("Some text with special characters Ã¦Ã¸Ã¥Ã†Ã˜Ã…$â‚¬Ã†Ã˜Ã… and some more arbitary text to make a very long line that will be broken for wrap-around (or elipsis)")
		}
	}()
}

func addLongLine() {
	addLine(strconv.Itoa(len(logText)) + " Some text with special characters Ã¦Ã¸Ã¥Ã†Ã˜Ã…$â‚¬Ã†Ã˜Ã… and some more arbitary text to make a very long line that will be broken for wrap-around (or elipsis)")
}

func addShortLine() {
	addLine("A short line")
}

func Form() wid.Wid {
	// lenstr := fmt.Sprintf("%d", len(logText))
	// cardName := CardList[CardTypeNo]
	sys.WinListMutex.RLock()
	defer sys.WinListMutex.RUnlock()
	return wid.Col(nil,
		wid.Label("IO-Card Production Acceptance Test", wid.H1C),
		wid.Row(nil,
			wid.Image(Images[0], wid.DefImg.W(0.5), ""),
			wid.Col(wid.ContStyle.W(0.5),
				wid.Edit(&Value2, "A long value here", nil, nil),
				wid.Label("FPS="+fmt.Sprintf("%0.2f", sys.WindowList[0].Fps()), nil),
				wid.Label("Log's last line="+getSize(), nil),
				wid.Btn("Add long line", nil, addLongLine, wid.Filled, ""),
				wid.Btn("Add short line", nil, addShortLine, wid.Filled, ""),
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
	runtime.LockOSThread()
	sys.Init()
	defer sys.Shutdown()
	win = sys.CreateWindow(0, 0, 0, 0, "IO-Card PAT", 1, 1.5)
	img, _ := wid.NewImage("rradi16.jpg")
	Images = append(Images, img)
	dummyLogGenerator()
	for sys.Running() {
		win.StartFrame(theme.Surface.Bg())
		wid.Show(Form())
		win.EndFrame()
		sys.PollEvents()
	}
}
