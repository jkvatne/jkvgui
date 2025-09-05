Origin
------

This subdirectory is a copy of https://github.com/neclepsio/gl/tree/master/all-core/gl
(Only the all-core version is copied)

Copyright (c) 2014 Eric Woroshow

neclepsio has released theese files under the MIT Licence, as shown at the end of this file.
The description is taken from his Readme file.

Description
------------

This repository holds Go bindings to various OpenGL versions. They are auto-generated using my fork of [Glow](https://github.com/neclepsio/glow).

The differences from [go-gl](https://github.com/go-gl/gl) are:
- WithOffset variants for some functions, so you don't have to pass pointers insteas of offsets (closes go-gl/gl issues [80](https://github.com/go-gl/gl/issues/80) and [124](https://github.com/go-gl/gl/issues/124)). Currently only functions `glDrawElements`, `glVertexAttribPointer`, `glGetVertexAttribPointerv` provide variants: let me know if you need more.
- No need to use cgo under Windows (much faster build times). It requires Go 1.12 for compatibilty profiles.

Features:
- Go functions that mirror the C specification using Go types.
- Support for multiple OpenGL APIs (GL/GLES/EGL/WGL/GLX/EGL), versions, and profiles.
- Support for extensions (including debug callbacks).

Requirements:
- A cgo compiler (typically gcc), only for non-Windows OS.
- On Ubuntu/Debian-based systems, the `libgl1-mesa-dev` package.

Licence
--------

The MIT License (MIT)

Copyright (c) 2014 Eric Woroshow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
