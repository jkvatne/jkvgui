package wid

import (
	"github.com/jkvatne/jkvgui/theme"
)

type GridStyle struct {
	EvenRole         theme.UIRole
	OddRole          theme.UIRole
	HeaderRole       theme.UIRole
	GridLinesWidth   float32
	GridLinesSpacing float32
}

func Grid(style *GridStyle, header Wid, gridLines ...Wid) Wid {
	return Col(nil, gridLines...)
}
