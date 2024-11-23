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
	Value         widget.Float
	Steps         int
	Title         string
	OnStepChanged func(step int) // returns step from 0 to Steps-1
}

func CreateBeepSlider(steps int, title string, onStepChanged func(step int)) *StepSlider {
	bs := &StepSlider{}
	bs.Steps = steps
	bs.Title = title
	bs.OnStepChanged = onStepChanged
	return bs
}

func (bs *StepSlider) Update(gtx layout.Context) {
	if !bs.Value.Dragging() {
		step := float32(1) / float32(bs.Steps-1)
		closest := findClosestStep(bs.Value.Value, step)
		if bs.Value.Value != closest {
			bs.Value.Value = closest
			bs.OnStepChanged(int(closest / step))
		}
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Millisecond)})
}

func (bs *StepSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, bs.Title)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, &bs.Value).Layout(gtx)
		}),
	)
}

func findClosestStep(value float32, step float32) float32 {
	if value <= 0 {
		return 0
	}
	return float32(math.Round(float64(value/step)) * float64(step))
}
