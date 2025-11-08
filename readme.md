# A minimal GUI for Windows and Linux

This GUI is inspired by the very good immediate mode GUI library called Gioui at
https://github.com/gioui/gio. But it does not share any code with it.

This library is much simpler and hopefully easier to understand.
It is fully immediate-mode, and persistant state is kept in a map, or in 
excplicite data strutctures.

Using a map simplifies generating forms. The key to the map is the value's address, 
and the state in the map is things like the edited text, cursor location etc.
(An exception is the scroller, which needs an explicit state).

It can be compiled on Windows without any CGO dependencies, by using the repoisitories
mentioned below.

# Examples

## Hello World

The complete code is as shown below. The function wid.Label() returns a function
that do the actual drawing.

```
package main

import (
	"github.com/jkvatne/jkvgui/sys"
	"github.com/jkvatne/jkvgui/wid"
)

func main() {
	sys.Init()
	defer sys.Shutdown()
	w := sys.CreateWindow(100, 100, 200, 100, "Hello world", 0, 2)
	for sys.Running() {
		w.StartFrame()
		wid.Show(wid.Label("Hello world!", nil))
		w.EndFrame()
		sys.PollEvents()
	}
}
```

To test, put this code into a file, f.ex. main.go.
Then type
```

    go mod init some.name.here
    go mod tidy
    go run main
```

For more examples, see the examples directory.

## Dependencies

### Freetype
A copy of the freetype code is included. It was cloned with
```sh
go get github.com/goki/freetype
```
the Freetype-Go source files are distributed under the BSD-style license

### Open-GL
To avoid dependency on GCC/CGO, I have used the bindings found in
https://github.com/neclepsio/gl/tree/master/all-core/gl

This version is identical to github.com/go-gl, but it does not use the c compiler.
This is much faster on Windows. Linux/Mac still needs GCC, but they have the c compiler installed by default.
A copy of the code is included.

## GLFW
To avoid dependency on GCC/CGO, I have made a translation of GLFW to pure go.
It is found in https://github.com/jkvatne/purego-glfw

GLFW is imported only once, in the sys/glfw.go file.
You can change to the standard version if you want.

## Installation on linux
You need to have open-gl installed with all developement libraries.
Make sure that git, gcc and pkg-config is installed and working ok.

The go-gl readme lists only the libgl1-mesa-dev package as dependency,
but when trying, I had to install all the packages below when I tested
it in WSL using a clean Ubuntu install.

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

## LICENSE

This software is released with the MIT license and it is found in the file LICENCE is in the root directory.
