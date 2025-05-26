//go:build noglfw

package sys

import (
	"encoding/binary"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jkvatne/jkvgui/f32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
	"log"
	"time"
)

const (
	KeyRight = iota
	KeyLeft
	KeyUp
	KeyDown
	KeySpace
	KeyEnter
	KeyEscape
	KeyBackspace
	KeyDelete
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyInsert
	KeyC
	KeyV
	KeyX
	ModShift
	ModControl
	ModAlt
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green          float32
	touchX         float32
	touchY         float32
	WindowHeightDp float32
	WindowWidthDp  float32
)

var (
	mousePos           f32.Pos
	leftBtnDown        bool
	leftBtnReleased    bool
	dragging           bool
	leftBtnDonwTime    time.Time
	LongPressTime      = time.Millisecond * 700
	DoubleClickTime    = time.Millisecond * 330
	leftBtnUpTime      = time.Now()
	leftBtnDoubleClick bool
)

var (
	LastRune rune
	LastKey  glfw.Key
	LastMods glfw.ModifierKey
)

func SetHresizeCursor() {

}
func SetVresizeCursor() {

}
func SetClipboardString(s string) {

}
func GetClipboardString() string {
	return ""
}

func SetCursorPos(x, y float32) {}

func RedrawsPrSec() int {
	return 1
}
func Running() bool {
	return false
}

func Reset() {

}
func LeftBtnDoubleClick(r f32.Rect) bool {
	return false
}

func Return() bool {
	return LastKey == glfw.KeyEnter || LastKey == glfw.KeyKPEnter
}
func Shutdonw() {

}

// Pos is the mouse pointer location in device-independent screen coordinates
func Pos() f32.Pos {
	return mousePos
}

// StartDrag is called when a widges wants to handle mouse events even
// outside its borders. Typically used when dragging a slider.
func StartDrag() f32.Pos {
	dragging = true
	return mousePos
}

// Hovered is true if the mouse pointer is inside the given rectangle
func Hovered(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && !dragging
}

// LeftBtnPressed is true if the mouse pointer is inside the
// given rectangle and the btn is pressed,
func LeftBtnPressed(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	return mousePos.Inside(r) && leftBtnDown && !dragging
}

// LeftBtnDown indicates that the user is holding the left btn down
// independent of the mouse pointer location
func LeftBtnDown() bool {
	if SuppressEvents {
		return false
	}
	return leftBtnDown
}

// LeftBtnClick returns true if the left btn has been clicked.
func LeftBtnClick(r f32.Rect) bool {
	if SuppressEvents {
		return false
	}
	if mousePos.Inside(r) && leftBtnReleased && time.Since(leftBtnDonwTime) < LongPressTime {
		leftBtnReleased = false
		return true
	}
	return false
}

func StartFrame(c f32.Color) {

}

func EndFrame(n int) {
}
func Shutdown() {

}
func InitWindow(wRequest, hRequest float32, name string, monitorNo int, userScale float32) {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	glctx.Uniform4f(color, 0, green, 0, 1)

	glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
