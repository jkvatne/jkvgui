package sys

import (
	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
	"reflect"
)

var (
	currentTag     interface{}
	moveToNext     bool
	moveToPrevious bool
	toNext         bool
	lastTag        interface{}
	SuppressEvents bool
)

type clickable struct {
	Rect   f32.Rect
	Action any
}

var clickables []clickable

func resetFocus() {
	clickables = clickables[0:0]
}

func moveByKey(forward bool) {
	if forward {
		moveToNext = true
	} else {
		moveToPrevious = true
	}
}

func At(rect f32.Rect, tag interface{}) bool {
	if moveToPrevious && gpu.TagsEqual(tag, currentTag) {
		currentTag = lastTag
		moveToPrevious = false
		gpu.Invalidate(0)
	}
	if gpu.TagsEqual(tag, currentTag) {
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
	lastTag = tag
	clickables = append(clickables, clickable{Rect: rect, Action: tag})
	if !gpu.CurrentInfo.Focused {
		return false
	}
	return gpu.CurrentInfo.Focused && gpu.TagsEqual(tag, currentTag) && !reflect.ValueOf(tag).IsNil()
}

func SetFocusedTag(action interface{}) {
	currentTag = action
	gpu.Invalidate(0)
}
