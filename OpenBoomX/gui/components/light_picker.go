package components

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
	"image/color"
)

type LightPicker struct {
	Picker       colorpicker.State
	CurrentColor color.NRGBA
	Solid        bool
}

func CreateLightPicker() LightPicker {
	picker := &LightPicker{}
	picker.CurrentColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	picker.Picker.SetColor(picker.CurrentColor)
	return *picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.Picker.Update(gtx) {
		lp.CurrentColor = lp.Picker.Color()
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return colorpicker.PickerStyle{
				Label:         "Lights Color",
				Theme:         th,
				State:         &lp.Picker,
				MonospaceFace: "Go Mono",
			}.Layout(gtx)
		}),
	)
}
