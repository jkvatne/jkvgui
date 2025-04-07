package shader

/*
import (
	"fmt"
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gl"
	"github.com/jkvatne/jkvgui/gl/glutil"
	"strings"
)

var Programs []uint32

// CompileShader compiles the shader program and returns the program as integer.
func CompileShader(source string, shaderType Enum) uint32 {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, source)
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
func NewProgram(vertexShaderSource, fragmentShaderSource string) (gl.Program, error) {
	program, err := glutil.CreateProgram(vertexShaderSource, fragmentShaderSource)
	f32.ExitOn(err, "Error compiling shaders")
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program))
	}

	gl.UseProgram(program)

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
*/
