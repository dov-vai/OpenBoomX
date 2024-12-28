package components

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"log"
	"math"
	"obx/gui/theme"
	"strconv"
)

type EqSlider struct {
	sliderValues    []float32
	editorValues    []string
	OnValuesChanged func(values []float32)
	sliders         []widget.Float
	editors         []widget.Editor
	defaultValues   []float32
}

func CreateEqSlider(onValuesChanged func(values []float32)) *EqSlider {
	sliderValues := make([]float32, 10)
	sliders := make([]widget.Float, 10)
	editorValues := make([]string, 10)
	editors := make([]widget.Editor, 10)

	eq := &EqSlider{
		sliderValues:    sliderValues,
		editorValues:    editorValues,
		OnValuesChanged: onValuesChanged,
		sliders:         sliders,
		editors:         editors,
		defaultValues:   []float32{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
	}

	err := eq.SetSliderValues(eq.defaultValues)
	if err != nil {
		return nil
	}

	return eq
}

func (eq *EqSlider) ResetValues() error {
	err := eq.SetSliderValues(eq.defaultValues)
	if err != nil {
		return err
	}
	eq.OnValuesChanged(eq.sliderValues)
	return nil
}

func (eq *EqSlider) GetSliderValues() []float32 {
	valuesCopy := make([]float32, len(eq.sliderValues))
	copy(valuesCopy, eq.sliderValues)
	return valuesCopy
}

// SetSliderValues sets the eq values, values must have a length of 10 (10 bands)
func (eq *EqSlider) SetSliderValues(values []float32) error {
	if len(values) != 10 {
		return fmt.Errorf("values must have a length of 10")
	}

	for i := 0; i < 10; i++ {
		eq.sliderValues[i] = values[i]
		eq.sliders[i].Value = values[i]
		eq.editorValues[i] = sliderToDb(values[i])
		eq.editors[i].SetText(eq.editorValues[i])
	}

	return nil
}

func (eq *EqSlider) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	children := make([]layout.FlexChild, len(eq.sliders))
	for i := range eq.sliders {
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					rec := op.Record(gtx.Ops)

					gtx.Constraints.Min, gtx.Constraints.Max = image.Pt(gtx.Constraints.Min.Y, gtx.Constraints.Min.X), image.Pt(gtx.Constraints.Max.Y, gtx.Constraints.Max.X)

					// yes.. magic number bad (approximate width of the text editor so it centers properly)
					estimateSize := 35

					op.Offset(image.Pt(estimateSize, estimateSize)).Add(gtx.Ops)

					// rotate slider 90 degrees
					op.Affine(f32.Affine2D{}.Rotate(f32.Pt(0, 0), float32(math.Pi/2))).Add(gtx.Ops)

					op.Offset(image.Pt(-estimateSize, -estimateSize)).Add(gtx.Ops)

					slider := material.Slider(th, &eq.sliders[i])

					dims := slider.Layout(gtx)

					dims.Size = image.Pt(dims.Size.Y, dims.Size.X)

					recorded := rec.Stop()

					recorded.Add(gtx.Ops)

					if !eq.sliders[i].Dragging() {
						if eq.sliderValues[i] != eq.sliders[i].Value {
							eq.sliderValues[i] = eq.sliders[i].Value
							eq.editorValues[i] = sliderToDb(eq.sliderValues[i])
							eq.editors[i].SetText(eq.editorValues[i])

							if eq.OnValuesChanged != nil {
								eq.OnValuesChanged(eq.sliderValues)
							}
						}
					}
					return dims
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Height: 8}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// update slider values
					if eq.editors[i].Text() != eq.editorValues[i] {
						if value, err := strconv.ParseFloat(eq.editors[i].Text(), 64); err == nil {
							sliderValue := dbToSlider(value)

							eq.editorValues[i] = eq.editors[i].Text()
							eq.sliderValues[i] = float32(sliderValue)
							eq.sliders[i].Value = float32(sliderValue)
							if eq.OnValuesChanged != nil {
								eq.OnValuesChanged(eq.sliderValues)
							}
						}
					}

					surfaceStyle := component.Surface(
						&material.Theme{
							Palette: material.Palette{
								Bg: theme.Surface0Color,
							},
						})

					surfaceStyle.CornerRadius = 4

					return surfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Editor(th, &eq.editors[i], "val").Layout(gtx)
						})
					})
				}),
			)
		})
	}
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, children...)
}

func (eq *EqSlider) OnPresetChanged(newPreset string, values []float32) {
	if newPreset == "" || values == nil {
		return
	}

	err := eq.SetSliderValues(values)
	if err != nil {
		log.Println(err)
	}
	eq.OnValuesChanged(eq.sliderValues)
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
