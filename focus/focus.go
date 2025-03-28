package focus

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"github.com/jkvatne/jkvgui/lib"
	"reflect"
)

var (
	currentTag     interface{}
	moveToNext     bool
	moveToPrevious bool
	toNext         bool
	lastTag        interface{}
)

func MoveByKey(forward bool) {
	if forward {
		moveToNext = true
	} else {
		moveToPrevious = true
	}
}

func At(rect f32.Rect, tag interface{}) bool {
	if moveToPrevious && lib.TagsEqual(tag, currentTag) {
		currentTag = lastTag
		moveToPrevious = false
		gpu.Invalidate(0)
	}
	if lib.TagsEqual(tag, currentTag) {
		if moveToNext {
			toNext = true
			moveToNext = false
			gpu.Invalidate(0)
		}
	} else if toNext {
		toNext = false
		currentTag = tag
		gpu.Invalidate(0)
	}
	AddFocusable(rect, tag)
	if !gpu.WindowHasFocus {
		return false
	}
	return gpu.WindowHasFocus && lib.TagsEqual(tag, currentTag) && !reflect.ValueOf(tag).IsNil()
}

func AddFocusable(rect f32.Rect, tag interface{}) {
	lastTag = tag
	gpu.Clickables = append(gpu.Clickables, gpu.Clickable{Rect: rect, Action: tag})
}

func Set(action interface{}) {
	currentTag = action
	gpu.Invalidate(0)
}
