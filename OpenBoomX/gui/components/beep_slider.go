package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"math"
)

type BeepSlider struct {
	Value widget.Float
}

func CreateBeepSlider() BeepSlider {
	bs := &BeepSlider{}
	return *bs
}

func (bs *BeepSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// FIXME: snapping doesn't always work
	if !bs.Value.Dragging() {
		lockValues := []float32{0, 0.25, 0.5, 0.75, 1}
		bs.Value.Value = findClosest(bs.Value.Value, lockValues)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, &bs.Value).Layout(gtx)
		}),
	)
}

func findClosest(value float32, targets []float32) float32 {
	var closest float32
	minDiff := float32(math.MaxFloat32)

	for _, target := range targets {
		diff := float32(math.Abs(float64(target - value)))
		if diff < minDiff {
			minDiff = diff
			closest = target
		}
	}

	return closest
}
