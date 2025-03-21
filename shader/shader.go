package shader

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"strings"
)

var Programs []uint32

// CompileShader compiles the shader program and returns the program as integer.
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
		txt := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(txt))
		panic(fmt.Sprintf("failed to compile %v: %v", source, txt))
	}
	return shader
}

// NewProgram links the frag and vertex shader programs
func NewProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	Programs = append(Programs, program)
	return program, nil
}
