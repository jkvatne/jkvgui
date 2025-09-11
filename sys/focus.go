package sys

import (
	"reflect"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

func (w *Window) MoveByKey(forward bool) {
	if forward {
		w.MoveToNext = true
	} else {
		w.MoveToPrevious = true
	}
}

func (w *Window) At(rect f32.Rect, tag interface{}) bool {
	if w.MoveToPrevious && gpu.TagsEqual(tag, w.CurrentTag) {
		w.CurrentTag = w.LastTag
		w.MoveToPrevious = false
		w.Invalidate()
	}
	if gpu.TagsEqual(tag, w.CurrentTag) {
		if w.MoveToNext {
			w.ToNext = true
			w.MoveToNext = false
			w.Invalidate()
		}
	} else if w.ToNext {
		w.ToNext = false
		w.CurrentTag = tag
		w.Invalidate()
	}
	w.LastTag = tag
	if !w.Focused {
		return false
	}
	return gpu.TagsEqual(tag, w.CurrentTag) && !reflect.ValueOf(tag).IsNil()
}

func (w *Window) SetFocusedTag(action interface{}) {
	w.CurrentTag = action
	w.Invalidate()
}
