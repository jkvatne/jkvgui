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
)

type person struct {
	Selected bool
	Name     string
	Age      float64
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
}

// makePersons will create a list of n persons.
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i < m; i++ {
		data[0].Age = data[0].Age + float64(i)
		data = append(data, data[0])
	}
	data = data[:n]
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
}

// onCheck is called when the header checkbox is clicked. It will set or clear all rows.
func onCheck() {
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}

// gw is the grid line width
const gw = 1.0

var GridStyle = wid.ContStyle

// GridDemo is a widget that lays out the grid. This is all that is needed.
func Form() wid.Wid {
	nameIcon = gpu.NavigationUnfoldMore
	addressIcon = gpu.NavigationUnfoldMore
	ageIcon = gpu.NavigationUnfoldMore

	// Configure a grid with headings and several rows
	var gridLines []wid.Wid
	header := wid.Row(nil,
		wid.Btn("", nil, onNameClick, wid.Text, ""),
		wid.Btn("Name", nil, onNameClick, wid.Text, ""),
		wid.Btn("Address", nil, onNameClick, wid.Text, ""),
		wid.Btn("Age", nil, onNameClick, wid.Text, ""),
		wid.Btn("Gender", nil, onNameClick, wid.Text, ""),
	)

	for i := 0; i < len(data); i++ {
		bgColor := f32.MultAlpha(theme.PrimaryContainer.Bg(), 50)
		if i%2 == 0 {
			bgColor = f32.MultAlpha(theme.SecondaryContainer.Bg(), 50)
		}
		gridLines = append(gridLines,
			wid.Row(GridStyle.C(bgColor),
				// One row of the grid is defined here
				wid.Checkbox("", &data[i].Selected, &wid.GridCb, ""),
				wid.Edit(&data[i].Name, "", nil, &wid.GridEdit),
				wid.Edit(&data[i].Address, "", nil, &wid.GridEdit),
				wid.Edit(&data[i].Age, "", nil, &wid.GridEdit),
				wid.Combo(&data[i].Status, []string{"Male", "Female", "Other"}, "", &wid.GridCombo),
			))

	}
	return wid.Col(nil,
		wid.Label("Grid demo", wid.H1C),
		wid.Grid(nil, header, gridLines...),
		wid.Separator(2, 0, theme.OnSurface),
		wid.Row(nil,
			wid.Elastic(),
			wid.Btn("Update", nil, nil, nil, "Click to update variables"),
			wid.Elastic(),
		),
	)
}

func main() {
	gpu.DebugWidgets = false
	makePersons(12)
	// Setting this true will draw a light blue frame around widgets.
	theme.SetDefaultPallete(true)
	// Full monitor (maximize) on monitor 2
	window := gpu.InitWindow(0, 0, "Rounded rectangle demo", 2, 2.0)
	defer gpu.Shutdown()
	sys.Initialize(window)
	wid.GridEdit.EditSize = 0.2
	wid.GridCombo.EditSize = 0.2
	for !window.ShouldClose() {
		sys.StartFrame(theme.Surface.Bg())
		// Paint a frame around the whole window
		gpu.Rect(gpu.WindowRect.Reduce(1), 1, f32.Transparent, f32.Red)
		// DrawIcon form
		Form()(wid.NewCtx())
		dialog.ShowDialogue()
		sys.EndFrame(50)
	}
}
