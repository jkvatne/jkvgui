package main

import (
	"log/slog"

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
)

func do() {
	slog.Info("Save clicked")
}

var MenuStyle = (&wid.ContainerStyle{}).W(0.3)

// Menu demonstrates how to show a list that is generated while drawing it.
func Menu() wid.Wid {
	return wid.Col(MenuStyle,
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

var ss = &wid.ScrollState{Width: 0.7}

// GetTotalCount returns the total number of articles
// This could be a database query for count(*)
func GetTotalCount() int {
	return 13
}

// GetItems could for example be a database query, reading article n
func GetItem(n int) wid.Wid {
	switch n {
	case 0:
		return wid.Label("Articles", &smallText)
	case 1:
		return wid.Col(&wid.Primary,
			wid.Label("Hiphop", nil),
			wid.Label("What Buttons are Artists Pushing When They Perform Live", &heading),
			wid.Label("12 hrs ago", &smallText),
			wid.Image(music, wid.DefImg.Bg(theme.PrimaryContainer), ""),
			wid.Row(nil,
				wid.Elastic(),
				wid.Btn("Save", gpu.ContentSave, do, nil, ""),
			),
		)
	case 2:
		return wid.Col(&wid.Primary,
			wid.Label("More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 3:
		return wid.Col(&wid.Primary,
			wid.Label("The new Beatles...", &heading),
		)
	case 4:
		return wid.Col(&wid.Primary,
			wid.Label("1 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 5:
		return wid.Col(&wid.Primary,
			wid.Label("2 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 6:
		return wid.Col(&wid.Primary,
			wid.Label("3 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 7:
		return wid.Col(&wid.Primary,
			wid.Label("4 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 8:
		return wid.Col(&wid.Primary,
			wid.Label("5 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 9:
		return wid.Col(&wid.Primary,
			wid.Label("6 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 10:
		return wid.Col(&wid.Primary,
			wid.Label("7 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 11:
		return wid.Col(&wid.Primary,
			wid.Label("8 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 12:
		return wid.Col(&wid.Primary,
			wid.Label("9 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	default:
		return nil
	}
}

func CachedItems() wid.Wid {
	return wid.CashedScroller(ss, GetItem, GetTotalCount)
}

func Items() wid.Wid {
	return wid.Scroller(ss,
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
		wid.Col(&wid.Primary,
			wid.Label("1 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("2 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("3 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("4 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("5 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("6 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("7 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("8 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
		wid.Col(&wid.Primary,
			wid.Label("9 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		),
	)
}

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(), CachedItems())
}

func main() {
	sys.Init()
	defer sys.Shutdown()
	w := sys.CreateWindow(-1, -1, 500, 500, "Material demo", 2, 1.0)
	sys.LoadOpenGl(w)

	music, _ = wid.NewImage("music.jpg")
	swift, _ = wid.NewImage("ts.jpg")
	smallText = wid.DefaultLabel
	smallText.FontNo = gpu.Normal10
	heading = *wid.H1L
	heading.Multiline = true
	heading.FontNo = gpu.Bold20
	theme.Colors[theme.OnPrimary] = f32.Yellow
	for sys.Running() {
		w.StartFrame(theme.Surface.Bg())
		wid.Show(Form())
		w.EndFrame()
		w.PollEvents()
	}
}
