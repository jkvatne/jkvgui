package gpu

import (
	"fmt"
	"github.com/jkvatne/jkvgui/gl"
)

var Programs []uint32

// CompileShader compiles the shader program and returns the program as integer.
func CompileShader(source string, shaderType uint32) uint32 {
	shader := ctx.CreateShader(shaderType)
	ctx.ShaderSource(shader, source)
	ctx.CompileShader(shader)
	status := ctx.GetShaderi(shader, gl.COMPILE_STATUS)
	if status == gl.FALSE {
		log := ctx.GetShaderInfoLog(shader)
		panic(fmt.Sprintf("failed to compile %v: %v", source, log))
	}
	return shader
}

// NewProgram links the frag and vertex shader programs
func NewProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	program := ctx.CreateProgram()
	ctx.AttachShader(program, vertexShader)
	ctx.AttachShader(program, fragmentShader)
	ctx.LinkProgram(program)

	status := ctx.GetProgrami(program, gl.LINK_STATUS)
	if status == gl.FALSE {
		log := ctx.GetProgramInfoLog(program)
		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	ctx.DeleteShader(vertexShader)
	ctx.DeleteShader(fragmentShader)
	Programs = append(Programs, program)
	return program, nil
}
