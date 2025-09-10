package sys

import (
	"reflect"

	"github.com/jkvatne/jkvgui/f32"
	"github.com/jkvatne/jkvgui/gpu"
)

func moveByKey(forward bool) {
	if forward {
		CurrentInfo.MoveToNext = true
	} else {
		CurrentInfo.MoveToPrevious = true
	}
}

func At(rect f32.Rect, tag interface{}) bool {
	if CurrentInfo.MoveToPrevious && gpu.TagsEqual(tag, CurrentInfo.CurrentTag) {
		CurrentInfo.CurrentTag = CurrentInfo.LastTag
		CurrentInfo.MoveToPrevious = false
		Invalidate()
	}
	if gpu.TagsEqual(tag, CurrentInfo.CurrentTag) {
		if CurrentInfo.MoveToNext {
			CurrentInfo.ToNext = true
			CurrentInfo.MoveToNext = false
			Invalidate()
		}
	} else if CurrentInfo.ToNext {
		CurrentInfo.ToNext = false
		CurrentInfo.CurrentTag = tag
		Invalidate()
	}
	CurrentInfo.LastTag = tag
	if !CurrentInfo.Focused {
		return false
	}
	return gpu.TagsEqual(tag, CurrentInfo.CurrentTag) && !reflect.ValueOf(tag).IsNil()
}

func SetFocusedTag(action interface{}) {
	CurrentInfo.CurrentTag = action
	Invalidate()
}
