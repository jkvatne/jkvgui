package main

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	_ "github.com/jkvatne/jkvgui/buildinfo" // Will print buildinfo at startup
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	MainRow   = wid.ContStyle.W(0.3)
	smallText wid.LabelStyle
	heading   wid.LabelStyle
	music     *wid.Img
	swift     *wid.Img
	entries   = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
	ss        = &wid.CachedScrollState{}
	MenuStyle = (&wid.ContainerStyle{}).W(0.3)
	// MyItemStyle is a container style for showing each record from the simulated database table
	MyItemStyle = wid.ContainerStyle{
		BorderRole:     theme.Outline,
		BorderWidth:    1,
		Role:           theme.PrimaryContainer,
		CornerRadius:   5.0,
		InsidePadding:  f32.Padding{L: 3, T: 3, R: 3, B: 3},
		OutsidePadding: f32.Padding{L: 3, T: 5, R: 3, B: 5},
	}
)

// dbCount returns the total number of articles
// This could be a database query for count(*)
func dbCount() int {
	// Simulation has items 0 to 20, for a total of 21 items.
	return 25
}

// dbRead simulates a database query, reading the contents of an article.
// Here we just simulate a database lookup by returning an item based on the given index n
func dbRead(n int) wid.Wid {
	switch n {
	case 0:
		return wid.Label("0 Articles "+strconv.Itoa(dbCount()), &smallText)
	case 1:
		return wid.Col(&MyItemStyle,
			wid.Label("1 Hiphop", nil),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
			wid.Label("12 hrs ago", &smallText),
			wid.Image(music, nil, wid.DefImg.Bg(theme.PrimaryContainer), ""),
			wid.Row(nil,
				wid.Flex(),
				wid.Btn("Save", gpu.ContentSave, do, nil, ""),
			),
		)
	case 2:
		return wid.Col(&MyItemStyle,
			wid.Label("2 More about Taylor Swift...", &heading),
			wid.Image(swift, nil, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 3:
		return wid.Col(&MyItemStyle,
			wid.Label("3 The new Beatles...", &heading),
		)
	case 4, 5, 6, 7:
		return wid.Col(&MyItemStyle,
			wid.Label(strconv.Itoa(n)+" More about Taylor Swift...", &heading),
			wid.Image(swift, nil, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return wid.Col(&MyItemStyle,
			wid.Label(strconv.Itoa(n)+" Some text here", &heading),
		)
	case 24:
		return wid.Col(&MyItemStyle,
			wid.Label(strconv.Itoa(n)+" Last item", &heading),
		)
	default:
		return nil
	}
}

func do() {
	slog.Info("Save clicked")
}

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return wid.Col(MenuStyle,
		wid.Label("Genre", &smallText),
		func(ctx wid.Ctx) wid.Dim {
			widgets := make([]wid.Wid, len(entries)+1)
			for i, s := range entries {
				widgets[i] = wid.Btn(s, gpu.Home, nil, wid.Text, "")
			}
			widgets[len(entries)] = wid.Label(fmt.Sprintf("MousePos = %5.0f, %5.0f ", sys.WindowList[0].MousePos().X, sys.WindowList[0].MousePos().Y), nil)
			return wid.Col(wid.Secondary.W(0.3), widgets...)(ctx)
		},
	)
}

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(),
		wid.CashedScroller(ss, &wid.DefaultScrollStyle, dbRead, dbCount))
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	slog.Info("Material example")
	sys.Init()
	defer sys.Shutdown()
	w := sys.CreateWindow(-1, -1, 500, 500, "Material demo", 1, 1.0)
	music, _ = wid.NewImage("music.jpg")
	swift, _ = wid.NewImage("ts.jpg")
	smallText = wid.DefaultLabel
	smallText.FontNo = gpu.Normal10
	heading = *wid.H1L
	heading.Multiline = true
	heading.FontNo = gpu.Bold20
	theme.Colors[theme.OnPrimary] = f32.Yellow
	ss.Width = 0.7
	for sys.Running() {
		w.StartFrame()
		wid.Show(Form())
		w.EndFrame()
		sys.PollEvents()
	}
}
