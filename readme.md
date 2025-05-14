# A minimal GUI for Windows and Linux

This GUI is inspired by the very good immediate mode GUI library called Gioui at
https://github.com/gioui/gio. But it does not share any code with it.

This library is much simpler and hopefully easier to understand.
It is fully immediate-mode, and persistant state is kept in a map, or in 
excplicite data strutctures.

Using a map simplifies generating forms. The key to the map is the value's address, 
and the state in the map is things like the edited text, cursor location etc.
(An exception is the scroller, which needs an explicit state).

# Examples

## Hello World

The complete code is as shown below. The function wid.Label() returns a function
that do the actual drawing. It has the signature `func(ctx wid.Ctx) wid.Dim` 
This drawing function needs a context, which is given in the second parenthesis.
If the context is empty, the widget returns the minimum dimension it needs.
If the context specifies an area, the widget will try to draw inside this area. 

```
func main() {
    window := gpu.InitWindow(150, 50, "Hello world", 0)
    callback.Initialize(window)
    for !window.ShouldClose() {
        gpu.StartFrame(f32.White)
        wid.Label("Hello world!", wid.H1C)(wid.NewCtx())
        gpu.EndFrame(50)
    }
}
```

## Dependencies

### Freetype
A copy of the freetype code is included. It was cloned with
```sh
go get github.com/goki/freetype
```
the Freetype-Go source files are distributed under the BSD-style license

### Open-gl
To avoid dependency on GCC/CGO, I have used the bindings found in https://github.com/neclepsio/gl/tree/master/all-core/gl
Linux/Mac still needs GCC, but they have it installed by default.

## GLFW

For window create etc. we need to use GLFW. it is found at https://github.com/go-gl/glfw/v3.3/glfw


## Installation on linux

You need to have open-gl installed with all developement libraries.
Make sure that git, gcc and pkg-config is installed and working ok.

The go-gl readme lists only the libgl1-mesa-dev package as dependency,
but when trying, I had to install all the packages below.

```
sudo apt install gcc
sudo apt install git
sudo apt install pkg-config

sudo apt install mesa-utils
sudo apt install freeglut3-dev
sudo apt install libgl1-mesa-dev
sudo apt install libxrandr-dev
sudo apt install libxcursor-dev
sudo apt install libxinerama-dev
sudo apt install libxinerama1git 
sudo apt install libxi-dev
sudo apt install libxxf86vm-dev
```

## Installation on Windows.

This package uses "C", so a gcc compiler must be available.<br>
A good version is found at https://www.mingw-w64.org/downloads/ <br>
By default it should install gcc at `C:\w64devkit\bin\gcc.exe`
and update the path.

You also need the open-gl drivers, which should be present by default.


## LICENSE

This software is released with the UNLICENSE.
See https://choosealicense.com/licenses/unlicense/
and the file UNLICENCE in the root directory.
