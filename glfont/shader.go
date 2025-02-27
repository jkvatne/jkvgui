package glfont

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/jkvatne/jkvgui/shader"

	"fmt"
	"strings"
)

// newProgram links the frag and vertex shader programs
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader := shader.CompileShader(shader.VertexFontShader, gl.VERTEX_SHADER)
	fragmentShader := shader.CompileShader(shader.FragmentFontShader, gl.FRAGMENT_SHADER)

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

	return program, nil
}
