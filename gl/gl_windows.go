// SPDX-License-Identifier: Unlicense OR MIT
// SPDX-License-Identifier: Unlicense OR MIT

package gl

import (
	"math"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	LibGLESv2                   = windows.NewLazyDLL("libGLESv2.dll")
	_glActiveTexture            = LibGLESv2.NewProc("glActiveTexture")
	_glAttachShader             = LibGLESv2.NewProc("glAttachShader")
	_glBeginQuery               = LibGLESv2.NewProc("glBeginQuery")
	_glBindAttribLocation       = LibGLESv2.NewProc("glBindAttribLocation")
	_glBindBuffer               = LibGLESv2.NewProc("glBindBuffer")
	_glBindBufferBase           = LibGLESv2.NewProc("glBindBufferBase")
	_glBindFramebuffer          = LibGLESv2.NewProc("glBindFramebuffer")
	_glBindRenderbuffer         = LibGLESv2.NewProc("glBindRenderbuffer")
	_glBindTexture              = LibGLESv2.NewProc("glBindTexture")
	_glBindVertexArray          = LibGLESv2.NewProc("glBindVertexArray")
	_glBlendEquation            = LibGLESv2.NewProc("glBlendEquation")
	_glBlendFuncSeparate        = LibGLESv2.NewProc("glBlendFuncSeparate")
	_glBufferData               = LibGLESv2.NewProc("glBufferData")
	_glBufferSubData            = LibGLESv2.NewProc("glBufferSubData")
	_glCheckFramebufferStatus   = LibGLESv2.NewProc("glCheckFramebufferStatus")
	_glClear                    = LibGLESv2.NewProc("glClear")
	_glClearColor               = LibGLESv2.NewProc("glClearColor")
	_glClearDepthf              = LibGLESv2.NewProc("glClearDepthf")
	_glDeleteQueries            = LibGLESv2.NewProc("glDeleteQueries")
	_glDeleteVertexArrays       = LibGLESv2.NewProc("glDeleteVertexArrays")
	_glCompileShader            = LibGLESv2.NewProc("glCompileShader")
	_glCopyTexSubImage2D        = LibGLESv2.NewProc("glCopyTexSubImage2D")
	_glGenerateMipmap           = LibGLESv2.NewProc("glGenerateMipmap")
	_glGenBuffers               = LibGLESv2.NewProc("glGenBuffers")
	_glGenFramebuffers          = LibGLESv2.NewProc("glGenFramebuffers")
	_glGenVertexArrays          = LibGLESv2.NewProc("glGenVertexArrays")
	_glGetUniformBlockIndex     = LibGLESv2.NewProc("glGetUniformBlockIndex")
	_glCreateProgram            = LibGLESv2.NewProc("glCreateProgram")
	_glGenRenderbuffers         = LibGLESv2.NewProc("glGenRenderbuffers")
	_glCreateShader             = LibGLESv2.NewProc("glCreateShader")
	_glGenTextures              = LibGLESv2.NewProc("glGenTextures")
	_glDeleteBuffers            = LibGLESv2.NewProc("glDeleteBuffers")
	_glDeleteFramebuffers       = LibGLESv2.NewProc("glDeleteFramebuffers")
	_glDeleteProgram            = LibGLESv2.NewProc("glDeleteProgram")
	_glDeleteShader             = LibGLESv2.NewProc("glDeleteShader")
	_glDeleteRenderbuffers      = LibGLESv2.NewProc("glDeleteRenderbuffers")
	_glDeleteTextures           = LibGLESv2.NewProc("glDeleteTextures")
	_glDepthFunc                = LibGLESv2.NewProc("glDepthFunc")
	_glDepthMask                = LibGLESv2.NewProc("glDepthMask")
	_glDisableVertexAttribArray = LibGLESv2.NewProc("glDisableVertexAttribArray")
	_glDisable                  = LibGLESv2.NewProc("glDisable")
	_glDrawArrays               = LibGLESv2.NewProc("glDrawArrays")
	_glDrawElements             = LibGLESv2.NewProc("glDrawElements")
	_glEnable                   = LibGLESv2.NewProc("glEnable")
	_glEnableVertexAttribArray  = LibGLESv2.NewProc("glEnableVertexAttribArray")
	_glEndQuery                 = LibGLESv2.NewProc("glEndQuery")
	_glFinish                   = LibGLESv2.NewProc("glFinish")
	_glFlush                    = LibGLESv2.NewProc("glFlush")
	_glFramebufferRenderbuffer  = LibGLESv2.NewProc("glFramebufferRenderbuffer")
	_glFramebufferTexture2D     = LibGLESv2.NewProc("glFramebufferTexture2D")
	_glGenQueries               = LibGLESv2.NewProc("glGenQueries")
	_glGetError                 = LibGLESv2.NewProc("glGetError")
	_glGetFloatv                = LibGLESv2.NewProc("glGetFloatv")
	_glGetIntegerv              = LibGLESv2.NewProc("glGetIntegerv")
	_glGetIntegeri_v            = LibGLESv2.NewProc("glGetIntegeri_v")
	_glGetProgramiv             = LibGLESv2.NewProc("glGetProgramiv")
	_glGetProgramInfoLog        = LibGLESv2.NewProc("glGetProgramInfoLog")
	_glGetQueryObjectuiv        = LibGLESv2.NewProc("glGetQueryObjectuiv")
	_glGetShaderiv              = LibGLESv2.NewProc("glGetShaderiv")
	_glGetShaderInfoLog         = LibGLESv2.NewProc("glGetShaderInfoLog")
	_glGetString                = LibGLESv2.NewProc("glGetString")
	_glGetUniformLocation       = LibGLESv2.NewProc("glGetUniformLocation")
	_glGetVertexAttribiv        = LibGLESv2.NewProc("glGetVertexAttribiv")
	_glGetVertexAttribPointerv  = LibGLESv2.NewProc("glGetVertexAttribPointerv")
	_glInvalidateFramebuffer    = LibGLESv2.NewProc("glInvalidateFramebuffer")
	_glIsEnabled                = LibGLESv2.NewProc("glIsEnabled")
	_glLinkProgram              = LibGLESv2.NewProc("glLinkProgram")
	_glPixelStorei              = LibGLESv2.NewProc("glPixelStorei")
	_glReadPixels               = LibGLESv2.NewProc("glReadPixels")
	_glRenderbufferStorage      = LibGLESv2.NewProc("glRenderbufferStorage")
	_glScissor                  = LibGLESv2.NewProc("glScissor")
	_glShaderSource             = LibGLESv2.NewProc("glShaderSource")
	_glTexImage2D               = LibGLESv2.NewProc("glTexImage2D")
	_glTexStorage2D             = LibGLESv2.NewProc("glTexStorage2D")
	_glTexSubImage2D            = LibGLESv2.NewProc("glTexSubImage2D")
	_glTexParameteri            = LibGLESv2.NewProc("glTexParameteri")
	_glUniformBlockBinding      = LibGLESv2.NewProc("glUniformBlockBinding")
	_glUniform1f                = LibGLESv2.NewProc("glUniform1f")
	_glUniform1i                = LibGLESv2.NewProc("glUniform1i")
	_glUniform2f                = LibGLESv2.NewProc("glUniform2f")
	_glUniform3f                = LibGLESv2.NewProc("glUniform3f")
	_glUniform4f                = LibGLESv2.NewProc("glUniform4f")
	_glUseProgram               = LibGLESv2.NewProc("glUseProgram")
	_glVertexAttribPointer      = LibGLESv2.NewProc("glVertexAttribPointer")
	_glViewport                 = LibGLESv2.NewProc("glViewport")
)
var (
	// Query caches.
	int32s   [100]int32
	float32s [100]float32
	uintptrs [100]uintptr
)

func Ptr(v []float32) uintptr {
	return uintptr(unsafe.Pointer(&v[0]))
}

func Str(s string) string {
	return s
}

type Context interface{}

func ActiveTexture(t Enum) {
	_, _, _ = syscall.SyscallN(_glActiveTexture.Addr(), 1, uintptr(t), 0, 0)
}
func AttachShader(p Program, s Shader) {
	_, _, _ = syscall.SyscallN(_glAttachShader.Addr(), 2, uintptr(p.V), uintptr(s.V), 0)
}
func BeginQuery(target Enum, query Query) {
	_, _, _ = syscall.SyscallN(_glBeginQuery.Addr(), 2, uintptr(target), uintptr(query.V), 0)
}
func BindAttribLocation(p Program, a Attrib, name string) {
	cname := cString(name)
	c0 := &cname[0]
	_, _, _ = syscall.SyscallN(_glBindAttribLocation.Addr(), 3, uintptr(p.V), uintptr(a), uintptr(unsafe.Pointer(c0)))
	// TODO issue34474KeepAlive(c)
}
func BindBuffer(target Enum, b uint32) {
	_, _, _ = syscall.SyscallN(_glBindBuffer.Addr(), 2, uintptr(target), uintptr(b), 0)
}
func BindBufferBase(target Enum, index int, b Buffer) {
	_, _, _ = syscall.SyscallN(_glBindBufferBase.Addr(), 3, uintptr(target), uintptr(index), uintptr(b.V))
}
func BindFramebuffer(target Enum, fb Framebuffer) {
	_, _, _ = syscall.SyscallN(_glBindFramebuffer.Addr(), 2, uintptr(target), uintptr(fb.V), 0)
}
func BindRenderbuffer(target Enum, rb Renderbuffer) {
	_, _, _ = syscall.SyscallN(_glBindRenderbuffer.Addr(), 2, uintptr(target), uintptr(rb.V), 0)
}
func BindImageTexture(unit int, t Texture, level int, layered bool, layer int, access, format Enum) {
	panic("not implemented")
}
func BindTexture(target Enum, t uint32) {
	_, _, _ = syscall.SyscallN(_glBindTexture.Addr(), 2, uintptr(target), uintptr(t), 0)
}
func BindVertexArray(a uint32) {
	_, _, _ = syscall.SyscallN(_glBindVertexArray.Addr(), 1, uintptr(a), 0, 0)
}
func BlendEquation(mode Enum) {
	_, _, _ = syscall.SyscallN(_glBlendEquation.Addr(), 1, uintptr(mode), 0, 0)
}
func BlendFuncSeparate(srcRGB, dstRGB, srcA, dstA Enum) {
	_, _, _ = syscall.SyscallN(_glBlendFuncSeparate.Addr(), 4, uintptr(srcRGB), uintptr(dstRGB), uintptr(srcA), uintptr(dstA), 0, 0)
}
func BlendFunc(srcRGB, dstRGB Enum) {
	BlendFuncSeparate(srcRGB, dstRGB, srcRGB, dstRGB)
}
func BufferData(target Enum, size int, data uintptr, usage Enum) {
	_, _, _ = syscall.SyscallN(_glBufferData.Addr(), 4, uintptr(target), uintptr(size), data, uintptr(usage))
}

// gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
func BufferSubData(target Enum, offset int, len int, src uintptr) {
	_, _, _ = syscall.SyscallN(_glBufferSubData.Addr(), 4, uintptr(target), uintptr(offset), uintptr(len), src, 0, 0)
	// TODO issue34474KeepAlive(s0)
}
func CheckFramebufferStatus(target Enum) Enum {
	s, _, _ := syscall.SyscallN(_glCheckFramebufferStatus.Addr(), 1, uintptr(target), 0, 0)
	return Enum(s)
}
func Clear(mask Enum) {
	_, _, _ = syscall.SyscallN(_glClear.Addr(), 1, uintptr(mask), 0, 0)
}
func ClearColor(red, green, blue, alpha float32) {
	_, _, _ = syscall.SyscallN(_glClearColor.Addr(), 4, uintptr(math.Float32bits(red)), uintptr(math.Float32bits(green)), uintptr(math.Float32bits(blue)), uintptr(math.Float32bits(alpha)), 0, 0)
}
func ClearDepthf(d float32) {
	_, _, _ = syscall.SyscallN(_glClearDepthf.Addr(), 1, uintptr(math.Float32bits(d)), 0, 0)
}
func CompileShader(s Shader) {
	_, _, _ = syscall.SyscallN(_glCompileShader.Addr(), 1, uintptr(s.V), 0, 0)
}
func CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	_, _, _ = syscall.SyscallN(_glCopyTexSubImage2D.Addr(), 8, uintptr(target), uintptr(level), uintptr(xoffset), uintptr(yoffset), uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0)
}
func GenerateMipmap(target Enum) {
	_, _, _ = syscall.SyscallN(_glGenerateMipmap.Addr(), 1, uintptr(target), 0, 0)
}
func GenBuffers(n uint32, b *uint32) {
	_, _, _ = syscall.SyscallN(_glGenBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(b)), 0)
}
func CreateBuffer() uint32 {
	var buf uintptr
	_, _, _ = syscall.SyscallN(_glGenBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&buf)), 0)
	return uint32(buf)
}
func CreateFramebuffer() uint32 {
	var fb uintptr
	_, _, _ = syscall.SyscallN(_glGenFramebuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&fb)), 0)
	return uint32(fb)
}
func CreateProgram() uint32 {
	p, _, _ := syscall.SyscallN(_glCreateProgram.Addr(), 0, 0, 0, 0)
	return uint32(p)
}
func CreateQuery() uint32 {
	var q uintptr
	_, _, _ = syscall.SyscallN(_glGenQueries.Addr(), 2, 1, uintptr(unsafe.Pointer(&q)), 0)
	return uint32(q)
}
func CreateRenderbuffer() Renderbuffer {
	var rb uintptr
	_, _, _ = syscall.SyscallN(_glGenRenderbuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&rb)), 0)
	return Renderbuffer{uint32(rb)}
}
func CreateShader(ty Enum) uint32 {
	s, _, _ := syscall.SyscallN(_glCreateShader.Addr(), 1, uintptr(ty))
	return uint32(s)
}
func GenTextures(n uint32, t *uint32) {
	_, _, _ = syscall.SyscallN(_glGenTextures.Addr(), 2, uintptr(n), uintptr(unsafe.Pointer(t)))
}
func GenVertexArrays(n uint32, t *uint32) {
	_, _, _ = syscall.SyscallN(_glGenVertexArrays.Addr(), 2, uintptr(n), uintptr(unsafe.Pointer(t)), 0)
}
func CreateVertexArray() uint32 {
	var t uintptr
	_, _, _ = syscall.SyscallN(_glGenVertexArrays.Addr(), 2, 1, uintptr(unsafe.Pointer(&t)), 0)
	return uint32(t)
}
func DeleteBuffer(v Buffer) {
	_, _, _ = syscall.SyscallN(_glDeleteBuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&v)), 0)
}
func DeleteFramebuffer(v Framebuffer) {
	_, _, _ = syscall.SyscallN(_glDeleteFramebuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&v.V)), 0)
}
func DeleteProgram(p Program) {
	_, _, _ = syscall.SyscallN(_glDeleteProgram.Addr(), 1, uintptr(p.V), 0, 0)
}
func DeleteQuery(query Query) {
	_, _, _ = syscall.SyscallN(_glDeleteQueries.Addr(), 2, 1, uintptr(unsafe.Pointer(&query.V)), 0)
}
func DeleteShader(s Shader) {
	_, _, _ = syscall.SyscallN(_glDeleteShader.Addr(), 1, uintptr(s.V), 0, 0)
}
func DeleteRenderbuffer(v Renderbuffer) {
	_, _, _ = syscall.SyscallN(_glDeleteRenderbuffers.Addr(), 2, 1, uintptr(unsafe.Pointer(&v.V)), 0)
}
func DeleteTexture(v Texture) {
	_, _, _ = syscall.SyscallN(_glDeleteTextures.Addr(), 2, 1, uintptr(unsafe.Pointer(&v.V)), 0)
}
func DeleteVertexArray(array VertexArray) {
	_, _, _ = syscall.SyscallN(_glDeleteVertexArrays.Addr(), 2, 1, uintptr(unsafe.Pointer(&array.V)), 0)
}
func DepthFunc(f Enum) {
	_, _, _ = syscall.SyscallN(_glDepthFunc.Addr(), 1, uintptr(f), 0, 0)
}
func DepthMask(mask bool) {
	var m uintptr
	if mask {
		m = 1
	}
	_, _, _ = syscall.SyscallN(_glDepthMask.Addr(), 1, m, 0, 0)
}
func DisableVertexAttribArray(a Attrib) {
	_, _, _ = syscall.SyscallN(_glDisableVertexAttribArray.Addr(), 1, uintptr(a), 0, 0)
}
func Disable(cap Enum) {
	_, _, _ = syscall.SyscallN(_glDisable.Addr(), 1, uintptr(cap), 0, 0)
}
func DrawArrays(mode Enum, first, count int) {
	_, _, _ = syscall.SyscallN(_glDrawArrays.Addr(), 3, uintptr(mode), uintptr(first), uintptr(count))
}
func DrawElements(mode Enum, count int, ty Enum, offset int) {
	_, _, _ = syscall.SyscallN(_glDrawElements.Addr(), 4, uintptr(mode), uintptr(count), uintptr(ty), uintptr(offset), 0, 0)
}

func Enable(cap Enum) {
	_, _, _ = syscall.SyscallN(_glEnable.Addr(), 1, uintptr(cap), 0, 0)
}
func EnableVertexAttribArray(a Attrib) {
	_, _, _ = syscall.SyscallN(_glEnableVertexAttribArray.Addr(), 1, uintptr(a), 0, 0)
}
func EndQuery(target Enum) {
	_, _, _ = syscall.SyscallN(_glEndQuery.Addr(), 1, uintptr(target), 0, 0)
}
func Finish() {
	_, _, _ = syscall.SyscallN(_glFinish.Addr(), 0, 0, 0, 0)
}
func Flush() {
	_, _, _ = syscall.SyscallN(_glFlush.Addr(), 0, 0, 0, 0)
}
func FramebufferRenderbuffer(target, attachment, renderbuffertarget Enum, renderbuffer Renderbuffer) {
	_, _, _ = syscall.SyscallN(_glFramebufferRenderbuffer.Addr(), 4, uintptr(target), uintptr(attachment), uintptr(renderbuffertarget), uintptr(renderbuffer.V), 0, 0)
}
func FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	_, _, _ = syscall.SyscallN(_glFramebufferTexture2D.Addr(), 5, uintptr(target), uintptr(attachment), uintptr(texTarget), uintptr(t.V), uintptr(level), 0)
}
func GetUniformBlockIndex(p Program, name string) uint {
	cname := cString(name)
	c0 := &cname[0]
	u, _, _ := syscall.SyscallN(_glGetUniformBlockIndex.Addr(), 2, uintptr(p.V), uintptr(unsafe.Pointer(c0)), 0)
	issue34474KeepAlive(c0)
	return uint(u)
}

func GetError() Enum {
	e, _, _ := syscall.SyscallN(_glGetError.Addr(), 0, 0, 0, 0)
	return Enum(e)
}
func GetInteger4(pname Enum) [4]int {
	_, _, _ = syscall.SyscallN(_glGetIntegerv.Addr(), 2, uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])), 0)
	var r [4]int
	for i := range r {
		r[i] = int(int32s[i])
	}
	return r
}
func GetInteger(pname Enum) int {
	_, _, _ = syscall.SyscallN(_glGetIntegerv.Addr(), 2, uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])), 0)
	return int(int32s[0])
}
func GetIntegeri(pname Enum, idx int) int {
	_, _, _ = syscall.SyscallN(_glGetIntegeri_v.Addr(), 3, uintptr(pname), uintptr(idx), uintptr(unsafe.Pointer(&int32s[0])))
	return int(int32s[0])
}
func GetFloat(pname Enum) float32 {
	_, _, _ = syscall.SyscallN(_glGetFloatv.Addr(), 2, uintptr(pname), uintptr(unsafe.Pointer(&float32s[0])), 0)
	return float32s[0]
}
func GetFloat4(pname Enum) [4]float32 {
	_, _, _ = syscall.SyscallN(_glGetFloatv.Addr(), 2, uintptr(pname), uintptr(unsafe.Pointer(&float32s[0])), 0)
	var r [4]float32
	copy(r[:], float32s[:])
	return r
}
func GetProgrami(p Program, pname Enum) int {
	_, _, _ = syscall.SyscallN(_glGetProgramiv.Addr(), 3, uintptr(p.V), uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])))
	return int(int32s[0])
}
func GetProgramInfoLog(p Program) string {
	n := GetProgrami(p, INFO_LOG_LENGTH)
	buf := make([]byte, n)
	_, _, _ = syscall.SyscallN(_glGetProgramInfoLog.Addr(), 4, uintptr(p.V), uintptr(len(buf)), 0, uintptr(unsafe.Pointer(&buf[0])), 0, 0)
	return string(buf)
}
func GetQueryObjectuiv(query Query, pname Enum) uint {
	_, _, _ = syscall.SyscallN(_glGetQueryObjectuiv.Addr(), 3, uintptr(query.V), uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])))
	return uint(int32s[0])
}
func GetShaderi(s Shader, pname Enum) int {
	_, _, _ = syscall.SyscallN(_glGetShaderiv.Addr(), 3, uintptr(s.V), uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])))
	return int(int32s[0])
}
func GetShaderInfoLog(s Shader) string {
	n := GetShaderi(s, INFO_LOG_LENGTH)
	buf := make([]byte, n)
	_, _, _ = syscall.SyscallN(_glGetShaderInfoLog.Addr(), 4, uintptr(s.V), uintptr(len(buf)), 0, uintptr(unsafe.Pointer(&buf[0])), 0, 0)
	return string(buf)
}
func GetString(pname Enum) string {
	s, _, _ := syscall.SyscallN(_glGetString.Addr(), 1, uintptr(pname), 0, 0)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(s)))
}
func GetUniformLocation(p uint32, name string) Uniform {
	cname := cString(name)
	c0 := &cname[0]
	u, _, _ := syscall.SyscallN(_glGetUniformLocation.Addr(), 2, uintptr(p), uintptr(unsafe.Pointer(c0)), 0)
	issue34474KeepAlive(c0)
	return Uniform{int(u)}
}
func GetVertexAttrib(index int, pname Enum) int {
	_, _, _ = syscall.SyscallN(_glGetVertexAttribiv.Addr(), 3, uintptr(index), uintptr(pname), uintptr(unsafe.Pointer(&int32s[0])))
	return int(int32s[0])
}

func GetVertexAttribPointer(index int, pname Enum) uintptr {
	_, _, _ = syscall.SyscallN(_glGetVertexAttribPointerv.Addr(), 3, uintptr(index), uintptr(pname), uintptr(unsafe.Pointer(&uintptrs[0])))
	return uintptrs[0]
}
func InvalidateFramebuffer(target, attachment Enum) {
	addr := _glInvalidateFramebuffer.Addr()
	if addr == 0 {
		// InvalidateFramebuffer is just a hint. Skip it if not supported.
		return
	}
	_, _, _ = syscall.SyscallN(addr, 3, uintptr(target), 1, uintptr(unsafe.Pointer(&attachment)))
}
func IsEnabled(cap Enum) bool {
	u, _, _ := syscall.SyscallN(_glIsEnabled.Addr(), 1, uintptr(cap), 0, 0)
	return u == TRUE
}
func LinkProgram(p Program) {
	_, _, _ = syscall.SyscallN(_glLinkProgram.Addr(), 1, uintptr(p.V), 0, 0)
}
func PixelStorei(pname Enum, param int) {
	_, _, _ = syscall.SyscallN(_glPixelStorei.Addr(), 2, uintptr(pname), uintptr(param), 0)
}
func MemoryBarrier(barriers Enum) {
	panic("not implemented")
}
func MapBufferRange(target Enum, offset, length int, access Enum) []byte {
	panic("not implemented")
}
func ReadPixels(x, y, width, height int32, format, ty Enum, d0 unsafe.Pointer) {
	_, _, _ = syscall.SyscallN(_glReadPixels.Addr(), 7, uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(format), uintptr(ty), uintptr(d0), 0, 0)
	issue34474KeepAlive(d0)
}
func RenderbufferStorage(target, internalformat Enum, width, height int) {
	_, _, _ = syscall.SyscallN(_glRenderbufferStorage.Addr(), 4, uintptr(target), uintptr(internalformat), uintptr(width), uintptr(height), 0, 0)
}
func Scissor(x, y, width, height int32) {
	_, _, _ = syscall.SyscallN(_glScissor.Addr(), 4, uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0, 0)
}
func ShaderSource(s Shader, src string) {
	var n = uintptr(len(src))
	psrc := &src
	_, _, _ = syscall.SyscallN(_glShaderSource.Addr(), 4, uintptr(s.V), 1, uintptr(unsafe.Pointer(psrc)), uintptr(unsafe.Pointer(&n)), 0, 0)
	issue34474KeepAlive(psrc)
}

// (GLenum  target, GLint  level, GLint  internalformat, GLsizei  width, GLsizei  height, GLint  border, GLenum  format, GLenum  type, const void * pixels);

func TexImage2D(target Enum, level int, internalFormat Enum, width int32, height int32, border uint, format Enum, ty Enum, pix uintptr) {
	_, _, _ = syscall.SyscallN(_glTexImage2D.Addr(), 9, uintptr(target), uintptr(level), uintptr(internalFormat), uintptr(width), uintptr(height), 0, uintptr(format), uintptr(ty), pix)
}
func TexStorage2D(target Enum, levels int, internalFormat Enum, width, height int) {
	_, _, _ = syscall.SyscallN(_glTexStorage2D.Addr(), 5, uintptr(target), uintptr(levels), uintptr(internalFormat), uintptr(width), uintptr(height), 0)
}
func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	d0 := &data[0]
	_, _, _ = syscall.SyscallN(_glTexSubImage2D.Addr(), 9, uintptr(target), uintptr(level), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(format), uintptr(ty), uintptr(unsafe.Pointer(d0)))
	issue34474KeepAlive(d0)
}
func TexParameteri(target, pname Enum, param int) {
	_, _, _ = syscall.SyscallN(_glTexParameteri.Addr(), 3, uintptr(target), uintptr(pname), uintptr(param))
}
func UniformBlockBinding(p Program, uniformBlockIndex uint, uniformBlockBinding uint) {
	_, _, _ = syscall.SyscallN(_glUniformBlockBinding.Addr(), 3, uintptr(p.V), uintptr(uniformBlockIndex), uintptr(uniformBlockBinding))
}
func Uniform1f(dst Uniform, v float32) {
	_, _, _ = syscall.SyscallN(_glUniform1f.Addr(), 2, uintptr(dst.V), uintptr(math.Float32bits(v)), 0)
}
func Uniform1i(dst Uniform, v int) {
	_, _, _ = syscall.SyscallN(_glUniform1i.Addr(), 2, uintptr(dst.V), uintptr(v), 0)
}
func Uniform2f(dst Uniform, v0, v1 float32) {
	_, _, _ = syscall.SyscallN(_glUniform2f.Addr(), 3, uintptr(dst.V), uintptr(math.Float32bits(v0)), uintptr(math.Float32bits(v1)))
}
func Uniform3f(dst Uniform, v0, v1, v2 float32) {
	_, _, _ = syscall.SyscallN(_glUniform3f.Addr(), 4, uintptr(dst.V), uintptr(math.Float32bits(v0)), uintptr(math.Float32bits(v1)), uintptr(math.Float32bits(v2)), 0, 0)
}
func Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	_, _, _ = syscall.SyscallN(_glUniform4f.Addr(), 5, uintptr(dst.V), uintptr(math.Float32bits(v0)), uintptr(math.Float32bits(v1)), uintptr(math.Float32bits(v2)), uintptr(math.Float32bits(v3)), 0)
}
func Uniform4fv(dst Uniform, n int, v *float32) {
	a := (*[4]float32)(unsafe.Pointer(v))
	v0 := (*a)[0]
	v1 := (*a)[1]
	v2 := (*a)[2]
	v3 := (*a)[3]
	_, _, _ = syscall.SyscallN(_glUniform4f.Addr(), 5, uintptr(dst.V), uintptr(math.Float32bits(v0)), uintptr(math.Float32bits(v1)), uintptr(math.Float32bits(v2)), uintptr(math.Float32bits(v3)), 0)
}
func UseProgram(p uint32) {
	_, _, _ = syscall.SyscallN(_glUseProgram.Addr(), 1, uintptr(p), 0, 0)
}
func UnmapBuffer(target Enum) bool {
	panic("not implemented")
}
func VertexAttribPointerWithOffset(dst Attrib, size int, ty Enum, normalized bool, stride uint32, offset unsafe.Pointer) {
	var norm uintptr
	if normalized {
		norm = 1
	}
	_, _, _ = syscall.SyscallN(_glVertexAttribPointer.Addr(), 6, uintptr(dst), uintptr(size), uintptr(ty), norm, uintptr(stride), uintptr(offset))
}
func Viewport(x, y, width, height int32) {
	_, _, _ = syscall.SyscallN(_glViewport.Addr(), 4, uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0, 0)
}

func cString(s string) []byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return b
}

// issue34474KeepAlive calls runtime.KeepAlive as a
// workaround for golang.org/issue/34474.
func issue34474KeepAlive(v interface{}) {
	runtime.KeepAlive(v)
}
