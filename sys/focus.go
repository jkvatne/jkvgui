package sys

import (
	"reflect"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

func moveByKey(forward bool) {
	if forward {
		gpu.CurrentInfo.MoveToNext = true
	} else {
		gpu.CurrentInfo.MoveToPrevious = true
	}
}

func At(rect f32.Rect, tag interface{}) bool {
	if gpu.CurrentInfo.MoveToPrevious && gpu.TagsEqual(tag, gpu.CurrentInfo.CurrentTag) {
		gpu.CurrentInfo.CurrentTag = gpu.CurrentInfo.LastTag
		gpu.CurrentInfo.MoveToPrevious = false
		Invalidate()
	}
	if gpu.TagsEqual(tag, gpu.CurrentInfo.CurrentTag) {
		if gpu.CurrentInfo.MoveToNext {
			gpu.CurrentInfo.ToNext = true
			gpu.CurrentInfo.MoveToNext = false
			Invalidate()
		}
	} else if gpu.CurrentInfo.ToNext {
		gpu.CurrentInfo.ToNext = false
		gpu.CurrentInfo.CurrentTag = tag
		Invalidate()
	}
	gpu.CurrentInfo.LastTag = tag
	if !gpu.CurrentInfo.Focused {
		return false
	}
	return gpu.TagsEqual(tag, gpu.CurrentInfo.CurrentTag) && !reflect.ValueOf(tag).IsNil()
}

func SetFocusedTag(action interface{}) {
	gpu.CurrentInfo.CurrentTag = action
	Invalidate()
}
