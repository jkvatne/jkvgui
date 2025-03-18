package focus

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
)

var (
	Current        interface{}
	MoveToNext     bool
	MoveToPrevious bool
	ToNext         bool
	Last           interface{}
)

func At(tag interface{}) bool {
	if !gpu.IsFocused {
		return false
	}
	return lib.TagsEqual(tag, Current)
}

func AddFocusable(rect f32.Rect, tag interface{}) {
	Last = tag
	gpu.Clickables = append(gpu.Clickables, gpu.Clickable{Rect: rect, Action: tag})
}

func Set(action interface{}) {
	Current = action
}

func Move(action interface{}) {
	if MoveToPrevious && lib.TagsEqual(action, Current) {
		Current = Last
		MoveToPrevious = false
		gpu.Invalidate(0)
	}
	if ToNext {
		ToNext = false
		Current = action
		gpu.Invalidate(0)
	}
	if lib.TagsEqual(action, Current) {
		if MoveToNext {
			ToNext = true
			MoveToNext = false
		}
		gpu.Invalidate(0)
	}
}
