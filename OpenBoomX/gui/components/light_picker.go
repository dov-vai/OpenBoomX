package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"obx/gui/obxcolorpicker"
)

type LightPicker struct {
	picker            obxcolorpicker.State
	currentColor      color.NRGBA
	radioButtonsGroup widget.Enum
	presetButtons     []widget.Clickable
	OnColorChanged    func(color color.NRGBA, solidColor bool)
}

// keys for radio buttons
const (
	dancingLights = "dancing"
	solidLights   = "solid"
)

var colorPresets = []color.NRGBA{
	{R: 237, G: 34, B: 36, A: 255},   // Red
	{R: 243, G: 115, B: 33, A: 255},  // Orange
	{R: 249, G: 188, B: 40, A: 255},  // Yellow
	{R: 131, G: 197, B: 50, A: 255},  // Green
	{R: 73, G: 195, B: 176, A: 255},  // Teal
	{R: 75, G: 178, B: 251, A: 255},  // Light Blue
	{R: 0, G: 111, B: 249, A: 255},   // Blue
	{R: 81, G: 32, B: 223, A: 255},   // Indigo
	{R: 180, G: 44, B: 215, A: 255},  // Violet
	{R: 255, G: 255, B: 255, A: 255}, // White
}

func CreateLightPicker(onColorChanged func(color color.NRGBA, solidColor bool)) *LightPicker {
	picker := &LightPicker{}
	picker.currentColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	picker.picker.SetColor(picker.currentColor)
	picker.radioButtonsGroup.Value = dancingLights
	picker.OnColorChanged = onColorChanged
	picker.presetButtons = make([]widget.Clickable, len(colorPresets))
	return picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.picker.Update(gtx) || lp.radioButtonsGroup.Update(gtx) {
		lp.sendUpdate()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var buttons = make([]layout.FlexChild, len(lp.presetButtons))
			for i, c := range colorPresets {
				if lp.presetButtons[i].Clicked(gtx) {
					lp.picker.SetColor(c)
					lp.sendUpdate()
				}

				btnStyle := material.IconButton(th, &lp.presetButtons[i], nil, "")
				btnStyle.Inset = layout.UniformInset(4)
				btnStyle.Background = c
				buttons[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return btnStyle.Layout(gtx)
					})
				})
			}

			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, buttons...)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return obxcolorpicker.PickerStyle{
				Label:         "Hex color",
				Theme:         th,
				State:         &lp.picker,
				MonospaceFace: "Go Mono",
			}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.radioButtonsGroup, dancingLights, "Dancing").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.radioButtonsGroup, solidLights, "Solid").Layout(gtx)
				}),
			)
		}),
	)
}

func (lp *LightPicker) sendUpdate() {
	lp.currentColor = lp.picker.Color()
	solidColor := false
	if lp.radioButtonsGroup.Value == solidLights {
		solidColor = true
	}
	lp.OnColorChanged(lp.currentColor, solidColor)
}
