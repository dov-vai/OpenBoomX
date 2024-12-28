package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"math"
	"time"
)

type StepSlider struct {
	value         widget.Float
	steps         int
	title         string
	OnStepChanged func(step int) // returns step from 0 to Steps-1
}

func CreateBeepSlider(steps int, title string, onStepChanged func(step int)) *StepSlider {
	bs := &StepSlider{}
	bs.steps = steps
	bs.title = title
	bs.OnStepChanged = onStepChanged
	return bs
}

func (bs *StepSlider) Update(gtx layout.Context) {
	if !bs.value.Dragging() {
		step := float32(1) / float32(bs.steps-1)
		closest := findClosestStep(bs.value.Value, step)
		if bs.value.Value != closest {
			bs.value.Value = closest
			bs.OnStepChanged(int(closest / step))
		}
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Millisecond)})
}

func (bs *StepSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, bs.title)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, &bs.value).Layout(gtx)
		}),
	)
}

func findClosestStep(value float32, step float32) float32 {
	if value <= 0 {
		return 0
	}
	return float32(math.Round(float64(value/step)) * float64(step))
}
