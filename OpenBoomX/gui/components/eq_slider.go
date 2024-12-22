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

// CreateEqSlider creates an EqSlider object, defaultValues can be nil
func CreateEqSlider(defaultValues []float32, onValuesChanged func(values []float32)) *EqSlider {
	sliderValues := make([]float32, 10)
	sliders := make([]widget.Float, 10)
	editorValues := make([]string, 10)
	editors := make([]widget.Editor, 10)

	if defaultValues == nil {
		defaultValues = []float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}
	}

	if len(defaultValues) != 10 {
		panic("defaultValues must have a length of 10")
	}

	for i := 0; i < 10; i++ {
		sliderValues[i] = defaultValues[i]
		sliders[i].Value = defaultValues[i]
		editorValues[i] = sliderToDb(defaultValues[i])
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
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								slider := material.Slider(th, &eq.Sliders[index])

								gtx.Constraints.Min, gtx.Constraints.Max = image.Pt(gtx.Constraints.Min.Y, gtx.Constraints.Min.X), image.Pt(gtx.Constraints.Max.Y, gtx.Constraints.Max.X)

								// rotate slider 90 degrees
								op.Affine(f32.Affine2D{}.Rotate(f32.Pt(0, 0), float32(3.14/2))).Add(gtx.Ops)

								if !eq.Sliders[index].Dragging() {
									if eq.SliderValues[index] != eq.Sliders[index].Value {
										eq.SliderValues[index] = eq.Sliders[index].Value
										eq.EditorValues[index] = sliderToDb(eq.SliderValues[index])
										eq.Editors[index].SetText(eq.EditorValues[index])

										if eq.OnValuesChanged != nil {
											eq.OnValuesChanged(eq.SliderValues)
										}
									}
								}
								return slider.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// use a horizontal Flex layout for editor and dB label
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								editor := material.Editor(th, &eq.Editors[index], "val")
								// update slider values
								if eq.Editors[index].Text() != eq.EditorValues[index] {
									if value, err := strconv.ParseFloat(eq.Editors[index].Text(), 64); err == nil {
										sliderValue := dbToSlider(value)

										eq.EditorValues[index] = eq.Editors[index].Text()
										eq.SliderValues[index] = float32(sliderValue)
										eq.Sliders[index].Value = float32(sliderValue)
										if eq.OnValuesChanged != nil {
											eq.OnValuesChanged(eq.SliderValues)
										}
									}
								}

								return editor.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(14), "dB")
								return label.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
	}
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, children...)
}

// dbToSlider maps dB value back to 0-1 range for the slider
func dbToSlider(value float64) float64 {
	sliderValue := (10 - value) / 20
	sliderValue = min(sliderValue, 1)
	sliderValue = max(sliderValue, 0)
	return sliderValue
}

// sliderToDb maps slider value to dB for the editor
func sliderToDb(value float32) string {
	dBValue := (1-value)*20 - 10
	return fmt.Sprintf("%.1f", dBValue)
}
