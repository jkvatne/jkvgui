package wid

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/theme"
)

type ContainerStyle struct {
	Width          float32
	Height         float32
	BorderRole     theme.UIRole
	BorderWidth    float32
	Role           theme.UIRole
	CornerRadius   float32
	InsidePadding  f32.Padding
	OutsidePadding f32.Padding
	HasGrid        bool
}

var ContStyle = &ContainerStyle{
	BorderRole:     theme.Transparent,
	BorderWidth:    0.0,
	Role:           theme.Transparent,
	CornerRadius:   0.0,
	InsidePadding:  f32.Padding{},
	OutsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
}

var GridStyle = ContainerStyle{
	BorderRole:   theme.Outline,
	BorderWidth:  0.5,
	Role:         theme.Surface,
	CornerRadius: 0.0,
	HasGrid:      true,
}

var Primary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    1,
	Role:           theme.PrimaryContainer,
	CornerRadius:   0.0,
	InsidePadding:  f32.Padding{L: 2, T: 2, R: 2, B: 2},
	OutsidePadding: f32.Padding{L: 2, T: 2, R: 2, B: 2},
}

var Secondary = ContainerStyle{
	BorderRole:     theme.Outline,
	BorderWidth:    0,
	Role:           theme.SecondaryContainer,
	CornerRadius:   9.0,
	InsidePadding:  f32.Padding{L: 4, T: 4, R: 4, B: 4},
	OutsidePadding: f32.Padding{L: 4, T: 4, R: 4, B: 4},
}

func (style *ContainerStyle) Size(w, h, bw float32) *ContainerStyle {
	ss := *style
	ss.Width = w
	ss.Height = h
	ss.BorderWidth = bw
	return &ss
}

func (style *ContainerStyle) W(w float32) *ContainerStyle {
	rr := *style
	rr.Width = w
	return &rr
}

func (style *ContainerStyle) H(h float32) *ContainerStyle {
	rr := *style
	rr.Height = h
	return &rr
}

func (style *ContainerStyle) C(c theme.UIRole) *ContainerStyle {
	rr := *style
	rr.Role = c
	return &rr
}

func (style *ContainerStyle) R(c theme.UIRole) *ContainerStyle {
	rr := *style
	rr.Role = c
	return &rr
}

func (style *ContainerStyle) TotalVerticalPadding() float32 {
	return style.OutsidePadding.T + style.OutsidePadding.B + 2*style.BorderWidth + style.InsidePadding.T + style.InsidePadding.B
}
