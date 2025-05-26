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
	windowHasFocus = true
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

func MoveByKey(forward bool) {
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
	AddFocusable(rect, tag)
	if !windowHasFocus {
		return false
	}
	return windowHasFocus && gpu.TagsEqual(tag, currentTag) && !reflect.ValueOf(tag).IsNil()
}

func AddFocusable(rect f32.Rect, tag interface{}) {
	lastTag = tag
	clickables = append(clickables, clickable{Rect: rect, Action: tag})
}

func SetFocusedTag(action interface{}) {
	currentTag = action
	gpu.Invalidate(0)
}
