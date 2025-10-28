package sys

import (
	"reflect"

	"github.com/jkvatne/jkvgui/gpu"
)

func (win *Window) MoveByKey(forward bool) {
	if forward {
		win.MoveToNext = true
	} else {
		win.MoveToPrevious = true
	}
}

func (win *Window) At(tag interface{}) bool {
	if win.MoveToPrevious && gpu.TagsEqual(tag, win.CurrentTag) {
		win.CurrentTag = win.LastTag
		win.MoveToPrevious = false
		win.Invalidate()
	}
	if gpu.TagsEqual(tag, win.CurrentTag) {
		if win.MoveToNext {
			win.ToNext = true
			win.MoveToNext = false
			win.Invalidate()
		}
	} else if win.ToNext {
		win.ToNext = false
		win.CurrentTag = tag
		win.Invalidate()
	}
	win.LastTag = tag
	if !win.Focused {
		// return false
	}
	return gpu.TagsEqual(tag, win.CurrentTag) && !reflect.ValueOf(tag).IsNil()
}

func (win *Window) SetFocusedTag(action interface{}) {
	win.CurrentTag = action
	win.Invalidate()
}
