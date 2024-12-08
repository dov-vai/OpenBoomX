package components

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"strconv"
)

type EqSlider struct {
	SliderValues    []float32
	EditorValues    []string
	OnValuesChanged func(values []float32)
	Sliders         []widget.Float
	Editors         []widget.Editor
}

func CreateEqSlider(onValuesChanged func(values []float32)) *EqSlider {
	sliderValues := make([]float32, 10)
	sliders := make([]widget.Float, 10)
	editorValues := make([]string, 10)
	editors := make([]widget.Editor, 10)

	for i := 0; i < 10; i++ {
		editorValues[i] = "0.00"
		editors[i].SetText(editorValues[i])
	}

	return &EqSlider{SliderValues: sliderValues, EditorValues: editorValues, OnValuesChanged: onValuesChanged, Sliders: sliders, Editors: editors}
}

func (eq *EqSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	children := make([]layout.FlexChild, len(eq.Sliders))
	for i := range eq.Sliders {
		index := i

		children[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// FIXME: one of the sliders goes out of bounds to the left, so inset is quite large
			inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(40)}

			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						slider := material.Slider(th, &eq.Sliders[index])

						// rotate slider 90 degrees
						op.Affine(f32.Affine2D{}.Rotate(f32.Pt(0, 0), float32(3.14/2))).Add(gtx.Ops)

						if !eq.Sliders[index].Dragging() {
							if eq.SliderValues[index] != eq.Sliders[index].Value {
								eq.SliderValues[index] = eq.Sliders[index].Value
								sliderText := fmt.Sprintf("%.2f", eq.SliderValues[index])
								eq.Editors[index].SetText(sliderText)
								eq.EditorValues[index] = sliderText
								if eq.OnValuesChanged != nil {
									eq.OnValuesChanged(eq.SliderValues)
								}
							}
						}
						dims := slider.Layout(gtx)
						// swap width and height because they were rotated
						return layout.Dimensions{Size: image.Pt(dims.Size.Y, dims.Size.X)}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, &eq.Editors[index], "val")

						// update slider values
						if eq.Editors[index].Text() != eq.EditorValues[index] {
							if value, err := strconv.ParseFloat(eq.Editors[index].Text(), 32); err == nil {
								value = min(value, 1)
								value = max(value, 0)
								eq.EditorValues[index] = eq.Editors[index].Text()
								eq.SliderValues[index] = float32(value)
								eq.Sliders[index].Value = float32(value)
								if eq.OnValuesChanged != nil {
									eq.OnValuesChanged(eq.SliderValues)
								}
							}
						}

						return editor.Layout(gtx)
					}),
				)
			})
		})
	}

	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, children...)
}
