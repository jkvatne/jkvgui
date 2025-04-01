package main

import (
	"github.com/jkvatne/jkvgui/button"
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/icon"
	"github.com/jkvatne/jkvgui/img"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	lightMode     = true
	MenuContainer wid.ContainerStyle
	mainContainer wid.ContainerStyle
	smallText     wid.LabelStyle
	heading       wid.LabelStyle
	music         *img.Img
	entries       = []string{"Classic", "Jazz", "Rock", "Hiphop", ""}
)

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return func(ctx wid.Ctx) wid.Dim {
		widgets := make([]wid.Wid, len(entries))
		// Note that "no" has to be atomic because it is concurrently updated in "ticker"
		// entries[4] = strconv.Itoa(int(no.Load()))
		for i, s := range entries {
			widgets[i] = button.Text(s, icon.Home, nil, nil, "")
		}
		return wid.Col(nil, widgets...)(ctx)
	}
}

func Items() wid.Wid {
	return wid.Col(nil,
		wid.Col(&wid.Primary,
			wid.Label("Music", &smallText),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
			img.W(music, img.FIT, ""),
			wid.Row(nil,
				wid.Label("12 hrs ago", &smallText),
				button.Filled("Save", icon.ContentSave, nil, nil, ""),
			),
		),
		wid.Col(nil,
			wid.Label("Click Save button to test the confirmation dialog", nil),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(nil, Menu(), Items())
}

/*
func Form() wid.Wid {
	return wid.Col(nil,
		wid.Label("Containers", wid.H1C),
		wid.Row(nil,
			wid.Col(&wid.ContainerStyle,
				wid.Label("Containers", wid.H1C),
				wid.Label("Text", nil),
			),
			wid.Col(&wid.ContainerStyle,
				wid.Label("Containers", wid.H1C),
				wid.Label("Text", nil),
				wid.Label("Text", nil),
			),
			wid.Col(&wid.ContainerStyle,
				wid.Label("Containers", wid.H1C),
				wid.Label("Text", nil),
				wid.Label("Text", nil),
				wid.Label("Text", nil),
				wid.Label("Text", nil),
			),
		),
	)
}
*/
func main() {
	// Setting this true will draw a light blue frame around widgets.
	gpu.DebugWidgets = false
	theme.SetDefaultPallete(lightMode)
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 1)
	defer gpu.Shutdown()
	sys.Initialize(window, 14)
	music, _ = img.New("music.jpg")
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		Form()(wid.NewCtx())
		wid.ShowHint(nil)
		dialog.Show(nil)
		sys.EndFrame(50)
	}
}
