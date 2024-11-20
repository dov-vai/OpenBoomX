package components

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"math"
	"time"
)

type BeepSlider struct {
	Value widget.Float
}

func CreateBeepSlider() BeepSlider {
	bs := &BeepSlider{}
	return *bs
}

func (bs *BeepSlider) Update(gtx layout.Context) {
	if !bs.Value.Dragging() {
		snapPoint := findClosest(bs.Value.Value, 0.25)
		if bs.Value.Value != snapPoint {
			bs.Value.Value = snapPoint
		}
	}
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Millisecond)})
}

func (bs *BeepSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(th, &bs.Value).Layout(gtx)
		}),
	)
}

func findClosest(value float32, multiple float32) float32 {
	if value <= 0 {
		return 0
	}
	return float32(math.Round(float64(value/multiple)) * float64(multiple))
}
