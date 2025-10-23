package main

import (
	"log"
	"log/slog"
	"strconv"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
)

var (
	MainRow      = wid.ContStyle.W(0.3)
	smallText    wid.LabelStyle
	heading      wid.LabelStyle
	music        *wid.Img
	swift        *wid.Img
	entries      = []string{"Classic", "Jazz", "Rock", "Hiphop", "Opera", "Brass", "Soul"}
	ss           = &wid.ScrollState{Width: 0.7}
	CacheStart   int
	Cache        []wid.Wid
	BatchSize    = 8
	CacheMaxSize = 16
	DbTotalCount int
	MenuStyle    = (&wid.ContainerStyle{}).W(0.3)
)

func do() {
	slog.Info("Save clicked")
}

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

// GetItem implements a cache.
func GetItem(idx int) wid.Wid {
	DbTotalCount = GetTotalCount()
	if idx >= DbTotalCount {
		return nil
	}
	if idx-CacheStart > CacheMaxSize*2 {
		// We must fill again since the request is more that a cache size above end. Can not reuse anything
		Cache = Cache[0:0]
		// Fill up from idx and upwards
		CacheStart = idx
		w := GetRangeFromDb(0, BatchSize)
		Cache = append(Cache, w...)
	} else if idx >= CacheStart+len(Cache) {
		slog.Debug("Reading beyond end of Cache", "idx", idx, "CacheStart", CacheStart, "len(Cache)", len(Cache))
		start := CacheStart + len(Cache)
		w := GetRangeFromDb(start, BatchSize)
		Cache = append(Cache, w...)
		// IF adding data made the cache too large, throw out the beginning
		if len(Cache) > CacheMaxSize {
			slog.Debug("len(Cache)>CacheMaxSize, delete from starte", "n", idx, "start", start)
			start = len(Cache) - CacheMaxSize
			Cache = Cache[start:]
			CacheStart = CacheStart + start
			slog.Debug("New", "size", len(Cache), "start", start)
		}
	} else if idx < CacheStart {
		oldCacheStart := CacheStart
		// Read in either a full batch, or the number of items missing at the front.
		cnt := min(BatchSize, CacheStart)
		// Starting at either 0 or the numer
		CacheStart = max(0, CacheStart-cnt)
		temp := GetRangeFromDb(CacheStart, cnt)
		if len(temp) != cnt {
			slog.Error("GetRangeFromDb returned too few items")
		}
		slog.Debug("Fill Cache at front", "idx", idx, "CacheStart", CacheStart, "oldCacheStart", oldCacheStart, "cnt", cnt)
		Cache = append(temp, Cache...)
	}
	if idx-CacheStart < 0 {
		slog.Error("GetItem failed", "idx", idx, "CacheStart", CacheStart, "len(Cache)", len(Cache))
		return nil
	}
	if idx-CacheStart >= len(Cache) {
		return nil
	}
	return Cache[idx-CacheStart]
}

// GetTotalCount returns the total number of articles
// This could be a database query for count(*)
func GetTotalCount() int {
	// Simulation has items 0 to 20, for a total of 21 items.
	return 21
}

func GetRangeFromDb(start int, count int) []wid.Wid {
	var w []wid.Wid
	DbTotalCount = GetTotalCount()
	slog.Debug("GetRangeFromDb", "start", start, "DbTotalCount", DbTotalCount)
	if start >= DbTotalCount {
		return nil
	}
	for i := 0; i < count; i++ {
		if i+start >= DbTotalCount {
			break
		}
		w = append(w, getFromDb(start+i))
	}
	return w
}

// GetItems could for example be a database query, reading article n
func getFromDb(n int) wid.Wid {
	switch n {
	case 0:
		return wid.Label("0 Articles", &smallText)
	case 1:
		return wid.Col(&wid.Primary,
			wid.Label("1 Hiphop", nil),
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
			wid.Label("2 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 3:
		return wid.Col(&wid.Primary,
			wid.Label("3 The new Beatles...", &heading),
		)
	case 4:
		return wid.Col(&wid.Primary,
			wid.Label("4 More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 5, 6, 7:
		return wid.Col(&wid.Primary,
			wid.Label(strconv.Itoa(n)+" More about Taylor Swift...", &heading),
			wid.Image(swift, wid.DefImg.Bg(theme.PrimaryContainer), ""),
		)
	case 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19:
		return wid.Col(&wid.Primary,
			wid.Label(strconv.Itoa(n)+" Some text here", &heading),
		)
	default:
		return nil
	}
}

func CachedItems() wid.Wid {
	return wid.CashedScroller(ss, GetItem, GetTotalCount)
}

/*
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
			wid.Label("6 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("7 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("8 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("9 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("10 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("11 The new Beatles...", &heading),
		),
		wid.Col(&wid.Primary,
			wid.Label("23 The new Beatles...", &heading),
		),
	)
}
*/

func Form() wid.Wid {
	return wid.Row(MainRow, Menu(), CachedItems())
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
	for sys.Running() {
		w.StartFrame(theme.Surface.Bg())
		wid.Show(Form())
		w.EndFrame()
		sys.PollEvents()
	}
}
