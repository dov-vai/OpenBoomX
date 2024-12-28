package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"log"
	"math"
	"time"
)

type StepSlider struct {
	value         widget.Float
	steps         int
	title         string
	OnStepChanged func(step int) // returns step from 0 to Steps-1
	stepLabels    []string
}

func CreateBeepSlider(steps int, title string, stepLabels []string, onStepChanged func(step int)) *StepSlider {
	bs := &StepSlider{}
	bs.steps = steps
	bs.title = title
	bs.OnStepChanged = onStepChanged

	if stepLabels != nil && len(stepLabels) != steps {
		log.Println("Warning: Number of step labels does not match the number of steps. Using step numbers instead.")
		bs.stepLabels = nil
	} else {
		bs.stepLabels = stepLabels
	}

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
	sliderWidth := gtx.Constraints.Max.X
	stepSizePx := float32(sliderWidth) / float32(bs.steps-1)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(th, bs.title)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, &bs.value).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := layout.Dimensions{}
			for i := 0; i < bs.steps; i++ {
				var labelText string
				if bs.stepLabels != nil && i < len(bs.stepLabels) {
					labelText = bs.stepLabels[i]
				} else {
					// Default to step number
					labelText = fmt.Sprintf("%d", i)
				}

				// subtract a little by an arbitrary value
				x := (float32(i) * stepSizePx) - float32(i*8)
				trans := op.Offset(image.Pt(int(x), -16)).Push(gtx.Ops)
				label := material.Body2(th, labelText)
				d := label.Layout(gtx)

				if i == bs.steps-1 {
					dims = d
				}

				trans.Pop()
			}
			return dims
		}),
		layout.Rigid(layout.Spacer{Height: 8}.Layout),
	)
}

func findClosestStep(value float32, step float32) float32 {
	if value <= 0 {
		return 0
	}
	return float32(math.Round(float64(value/step)) * float64(step))
}
