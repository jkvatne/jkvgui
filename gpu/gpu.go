package gpu

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"log"
	"strings"
)

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func CompileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
		s := fmt.Sprintf("Failed to compile %v: %v", source, infoLog)
		panic(s)
	}
	return shader
}

func CreateProgram(vert, frag string) uint32 {
	vertexShader := CompileShader(vert, gl.VERTEX_SHADER)
	fragmentShader := CompileShader(frag, gl.FRAGMENT_SHADER)
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic("Initialization error for OpenGL: " + err.Error())
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.SRC_ALPHA)
}
