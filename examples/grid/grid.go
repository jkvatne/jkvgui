// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls vertically and horizontally and implements highlighting of rows.

import (
	"github.com/jkvatne/jkvgui/dialog"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/theme"
	"github.com/jkvatne/jkvgui/wid"
	"log/slog"
	"sort"
)

var (
	Alternative = "Fractional"
	fontSize    = "Medium"
	// Column widths are given in units of approximately one average character width (en).
	// A width of zero means the widget's natural size should be used (f.ex. checkboxes)
	wideColWidth  = []float32{0, 60, 60, 10, 30}
	smallColWidth = []float32{0, 13, 13, 12, 12}
	fracColWidth  = []float32{0, 0.3, 0.3, .2, .2}
	selectAll     bool
	doOccupy      bool
	withoutHeader bool = false
	nameIcon      *gpu.Icon
	addressIcon   *gpu.Icon
	ageIcon       *gpu.Icon
	dir           bool
	line          string
	ro            *wid.EditStyle
)

type person struct {
	Selected bool
	Name     string
	Age      float32
	Address  string
	Status   int
}

var data = []person{
	{Name: "Ågåt Karlsen", Age: 21, Address: "Storgata 1", Status: 0},
	{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 2", Selected: true, Status: 0},
	{Name: "Nils Aure", Age: 23, Address: "Brogata 3"},
	{Name: "Kai Oppdal", Age: 24, Address: "Soleieveien 4"},
	{Name: "Gro Arneberg", Age: 25, Address: "Blomsterveien 5"},
	{Name: "Ole Kolås", Age: 26, Address: "Blåklokkevikua 6"},
	{Name: "Per Pedersen", Age: 27, Address: "Gamleveien 7"},
	{Name: "Nils Vukubråten", Age: 28, Address: "Nygata 8"},
	{Name: "Sindre Gratangen", Age: 29, Address: "Brosundet 9"},
	{Name: "Gro Nilsasveen", Age: 30, Address: "Blomsterveien 10"},
	{Name: "Petter Olsen", Age: 31, Address: "Katavågen 11"},
	{Name: "Per Pedersen", Age: 32, Address: "Nidelva 12"},
	{Name: "Oleg Karlsen", Age: 21, Address: "Storgata 1", Status: 0},
}

// makePersons will create a list of n persons for testing
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i <= m; i++ {
		data = append(data, data[i%len(data)])
	}
	for i := 0; i < len(data); i++ {
		data[i].Age = float32(i + 1)
	}
}

func doUpdate() {
	slog.Info("doUpdate()")
}

func onNameClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
		nameIcon = gpu.NavigationArrowDownward
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
		nameIcon = gpu.NavigationArrowUpward
	}
	addressIcon = gpu.NavigationUnfoldMore
	ageIcon = gpu.NavigationUnfoldMore
	dir = !dir
	wid.ClearBuffers()
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
		addressIcon = gpu.NavigationArrowDownward
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
		addressIcon = gpu.NavigationArrowUpward
	}
	nameIcon = gpu.NavigationUnfoldMore
	ageIcon = gpu.NavigationUnfoldMore
	dir = !dir
	wid.ClearBuffers()
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
		ageIcon = gpu.NavigationArrowDownward
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
		ageIcon = gpu.NavigationArrowUpward
	}
	nameIcon = gpu.NavigationUnfoldMore
	addressIcon = gpu.NavigationUnfoldMore
	dir = !dir
	wid.ClearBuffers()
}

// onCheck is called when the header checkbox is clicked. It will set or clear all rows.
func onCheck() {
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}

// gw is the grid line width
const gw = 1.0

var ss = &wid.ScrollState{Height: 0.5}
var GridStyle = wid.ContStyle

// Form is a widget that lays out the grid. This is all that is needed.
func Form() wid.Wid {
	nameIcon = gpu.NavigationUnfoldMore
	addressIcon = gpu.NavigationUnfoldMore
	ageIcon = gpu.NavigationUnfoldMore

	// Configure a grid with headings and several rows
	var gridLines []wid.Wid
	header := wid.Row(nil,
		wid.Btn("", nil, nil, wid.CbHeader, ""),
		wid.Btn("Name", nil, onNameClick, wid.Header, ""),
		wid.Btn("Address", nil, onAddressClick, wid.Header, ""),
		wid.Btn("Age", nil, onAgeClick, wid.Header, ""),
		wid.Btn("Gender", nil, nil, wid.Header, ""),
	)

	for i := 0; i < len(data); i++ {
		bgColor := theme.PrimaryContainer.Bg().MultAlpha(0.5)
		if i%2 == 0 {
			bgColor = theme.SecondaryContainer.Bg().MultAlpha(0.5)
		}
		gridLines = append(gridLines,
			wid.Row(GridStyle.C(bgColor),
				// One row of the grid is defined here
				wid.Checkbox("", &data[i].Selected, &wid.GridCb, ""),
				wid.Edit(&data[i].Name, "", nil, ro),
				wid.Edit(&data[i].Address, "", nil, &wid.GridEdit),
				wid.Edit(&data[i].Age, "", nil, &wid.GridEdit),
				wid.Combo(&data[i].Status, []string{"Male", "Female", "Other"}, "", &wid.GridCombo),
			))

	}
	return wid.Col(nil,
		wid.Label("Grid demo", wid.H1C),
		header,
		wid.Scroller(ss, gridLines...),
		wid.Line(0, 1.0, theme.Surface),
		wid.Row(nil,
			wid.Elastic(),
			wid.Btn("Update", nil, doUpdate, nil, "Click to update variables"),
			wid.Elastic(),
		),
	)
}

func main() {
	sys.Init()
	defer sys.Shutdown()

	makePersons(30)
	// Full monitor (maximize) on monitor 2 (if it is present), and with userScale=2
	sys.CreateWindow(0, 0, 880, 880, "Rounded rectangle demo", 2, 2.0)
	ro = wid.GridEdit.RO()
	for sys.Running() {
		sys.StartFrame(theme.Surface.Bg())
		// Paint a frame around the whole window
		gpu.Rect(gpu.Info[0].WindowRect.Reduce(1), 1, f32.Transparent, f32.Red)
		// Draw form
		Form()(wid.NewCtx())
		dialog.ShowDialogue()
		sys.EndFrame()
		sys.PollEvents()
	}
}
