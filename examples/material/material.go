package main

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	lightMode = true
	MainRow   = wid.ContStyle.W(0.3)
	smallText wid.LabelStyle
	heading   wid.LabelStyle
	music     *wid.Img
	swift     *wid.Img
	entries   = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
)

func do() {
	// f, _ := os.Create("mem.pprof")
	// pprof.WriteHeapProfile(f)
	// f.Close()
	// go tool pprof mem.pprof
}

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return wid.Col((&wid.ContainerStyle{}).W(0.3),
		wid.Label("Genre", &smallText),
		func(ctx wid.Ctx) wid.Dim {
			widgets := make([]wid.Wid, len(entries))
			for i, s := range entries {
				widgets[i] = wid.Btn(s, gpu.Home, nil, wid.Text, "")
			}
			return wid.Col(wid.Secondary.W(0.3), widgets...)(ctx)
		},
	)
}

var ES wid.ContainerStyle
var ss = &wid.ScrollState{}

func Items() wid.Wid {
	return wid.Col((&wid.ContainerStyle{}).W(0.7),
		wid.Scroller(ss,
			wid.Label("Articles", &smallText),
			wid.Col(&wid.Primary,
				wid.Label("Hiphop", nil),
				wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
				wid.Label("12 hrs ago", &smallText),
				wid.Image(music, wid.DefImg.Bg(theme.PrimaryContainer), ""),
				wid.Row(nil,
					wid.Elastic(),
					wid.Btn("Save", gpu.ContentSave, do, nil, ""),
				),
			),
			wid.Col(&wid.Primary,
				wid.Label("More about Taylor Swift...", &heading),
				wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
			),
			wid.Col(&wid.Primary,
				wid.Label("The new Beatles...", &heading),
			),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(), Items())
}

func main() {
	// Start pprof server on port 6060
	// View at	http://localhost:6060/debug/pprof/heap
	/*
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				log.Printf("pprof server failed: %v", err)
			}
		}()
	*/
	// Setting DebugWidgets true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(500, 500, "Rounded rectangle demo", 2, 1.0)
	defer gpu.Shutdown()
	sys.Initialize(window)
	music, _ = wid.NewImage("music.jpg")
	swift, _ = wid.NewImage("ts.jpg")
	smallText = wid.DefaultLabel
	smallText.FontNo = gpu.Normal10
	heading = *wid.H1L
	heading.Multiline = true
	heading.FontNo = gpu.Bold20
	theme.Colors[theme.OnPrimary] = f32.Yellow
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.ShowDialogue(nil)
		sys.EndFrame(50)
	}
}
