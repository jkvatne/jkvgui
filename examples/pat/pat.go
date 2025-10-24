package main

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	Value2  = "Value2"
	Images  []*wid.Img
	logText []string
	win     *sys.Window
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
	sys.WinListMutex.RLock()
	defer sys.WinListMutex.RUnlock()
	return wid.Col(nil,
		wid.Label("IO-Card Production Acceptance Test", wid.H1C),
		wid.Row(wid.ContStyle.H(0.7),
			wid.Image(Images[0], wid.DefImg.W(0.7), ""),
			wid.Col(wid.ContStyle.W(0.3),
				wid.Edit(&Value2, "A long value here", nil, nil),
				wid.Label("FPS="+fmt.Sprintf("%0.2f", sys.WindowList[0].Fps()), nil),
				wid.Label("Log's last line="+getSize(), nil),
				wid.Btn("Add long line", nil, addLongLine, wid.Filled, ""),
				wid.Btn("Add short line", nil, addShortLine, wid.Filled, ""),
			),
		),
		wid.Memo(&logText, nil),
	)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	slog.Info("PAT example")
	sys.Init()
	defer sys.Shutdown()
	win = sys.CreateWindow(0, 0, 0, 0, "IO-Card PAT", 1, 1.5)
	img, _ := wid.NewImage("rradi16.jpg")
	Images = append(Images, img)
	dummyLogGenerator()
	for sys.Running() {
		win.StartFrame()
		wid.Show(Form())
		win.EndFrame()
		sys.PollEvents()
	}
}
