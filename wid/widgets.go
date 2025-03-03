package wid

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"log"
	"unsafe"
)

type Dim struct {
	w        float32
	h        float32
	baseline float32
}

type Ctx struct {
	Rect     lib.Rect
	Baseline float32
}

type Padding struct {
	L float32
	T float32
	R float32
	B float32
}

type Wid func(ctx Ctx) Dim

type Clickable struct {
	Rect   lib.Rect
	Action func()
}

type RowSetup struct {
	Height float32
}
type ColSetup struct {
	Widths []float32
}

var Clickables []Clickable
var MousePos lib.Pos
var MouseBtnDown bool
var MouseBtnReleased bool
var InFocus interface{}

type eface struct {
	typ, val unsafe.Pointer
}

func ptr(arg interface{}) unsafe.Pointer {
	return (*eface)(unsafe.Pointer(&arg)).val
}

func Row(setup RowSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		maxY := float32(0)
		maxB := float32(0)
		sumW := float32(0)
		ctx0 := Ctx{}
		ne := 0
		dims := make([]Dim, len(widgets))
		for i, w := range widgets {
			dims[i] = w(ctx0)
			maxY = max(maxY, dims[i].h)
			maxB = max(maxB, dims[i].baseline)
			sumW += dims[i].w
			if dims[i].w == 0 {
				ne++
			}
		}

		if ne > 0 {
			remaining := ctx.Rect.W - sumW
			for i, d := range dims {
				if d.w == 0 {
					dims[i].w = remaining / float32(ne)
				}
			}
		}
		ctx1 := ctx
		ctx1.Rect.H = maxY
		ctx1.Baseline = maxB
		for i, w := range widgets {
			_ = w(ctx1)
			ctx1.Rect.X += dims[i].w
		}
		gpu.RoundedRect(ctx.Rect.X, ctx.Rect.Y, ctx.Rect.W, maxY, 0, 1, gpu.Transparent, gpu.Color{0, 1, 0, 0.2})
		return Dim{w: sumW, h: maxY, baseline: maxB}
	}
}

func Col(setup ColSetup, widgets ...Wid) Wid {
	return func(ctx Ctx) Dim {
		TotHeight := float32(0.0)
		maxY := float32(0.0)
		if ctx.Rect.H == 0 {
			for _, w := range widgets {
				h := w(ctx).h
				maxY = max(maxY, h)
				TotHeight += h
			}
			return Dim{ctx.Rect.W, maxY * float32(len(widgets)), 0}
		} else {
			for _, w := range widgets {
				ctx.Rect.Y += w(ctx).h
			}
			return Dim{100, TotHeight, 0}
		}
	}
}

func Label(text string, size float32, p Padding, fontNo int) Wid {
	return func(ctx Ctx) Dim {
		if ctx.Rect.H == 0 {
			height := (gpu.Fonts[fontNo].Ascent+gpu.Fonts[fontNo].Descent)*size/gpu.InitialSize + p.T + p.B
			width := gpu.Fonts[fontNo].Width(size, text)/gpu.InitialSize + p.L + p.R
			return Dim{w: width, h: height, baseline: gpu.Fonts[fontNo].Ascent*size/gpu.InitialSize + p.T}
		} else {
			gpu.Fonts[fontNo].SetColor(gpu.Black)
			gpu.Fonts[fontNo].Printf(ctx.Rect.X+p.L, ctx.Rect.Y+p.T+ctx.Baseline, size, text)
			return Dim{}
		}
	}
}

func Elastic() Wid {
	return func(ctx Ctx) Dim {
		return Dim{}
	}
}

func MouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	x, y := w.GetCursorPos()
	MousePos.X = float32(x)
	MousePos.Y = float32(y)
	var pos = lib.Pos{float32(x), float32(y)}
	log.Printf("Mouse btn %d clicked at %0.1f,%0.1f, Action %d\n", button, x, y, action)
	if action == glfw.Release {
		MouseBtnDown = false
		MouseBtnReleased = true
		for _, clickable := range Clickables {
			if pos.Inside(clickable.Rect) {
				clickable.Action()
			}
		}
	} else if action == glfw.Press {
		MouseBtnDown = true
	}
}

func RR(r lib.Rect, t float32, fillColor gpu.Color, frameColor gpu.Color) {
	gpu.RoundedRect(r.X, r.Y, r.W, r.H, r.RR, t, fillColor, frameColor)
}

func Hovered(r lib.Rect) bool {
	return MousePos.Inside(r)
}

func MousePosCallback(xw *glfw.Window, xpos float64, ypos float64) {
	MousePos.X = float32(xpos)
	MousePos.Y = float32(ypos)
}

func Pressed(r lib.Rect) bool {
	return MousePos.Inside(r) && MouseBtnDown
}

func Released(r lib.Rect) bool {
	return MousePos.Inside(r) && MouseBtnReleased
}

func Focused(tag interface{}) bool {
	return ptr(tag) == ptr(InFocus)
}
