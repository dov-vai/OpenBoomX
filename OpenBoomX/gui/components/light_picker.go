package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"obx/gui/obxcolorpicker"
)

type LightPicker struct {
	Picker            obxcolorpicker.State
	CurrentColor      color.NRGBA
	RadioButtonsGroup widget.Enum
	OnColorChanged    func(color color.NRGBA, solidColor bool)
}

func CreateLightPicker(onColorChanged func(color color.NRGBA, solidColor bool)) *LightPicker {
	picker := &LightPicker{}
	picker.CurrentColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	picker.Picker.SetColor(picker.CurrentColor)
	picker.RadioButtonsGroup.Value = "dancing"
	picker.OnColorChanged = onColorChanged
	return picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.Picker.Update(gtx) || lp.RadioButtonsGroup.Update(gtx) {
		lp.CurrentColor = lp.Picker.Color()
		solidColor := false
		if lp.RadioButtonsGroup.Value == "solid" {
			solidColor = true
		}
		lp.OnColorChanged(lp.CurrentColor, solidColor)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return obxcolorpicker.PickerStyle{
				Label:         "Lights Color",
				Theme:         th,
				State:         &lp.Picker,
				MonospaceFace: "Go Mono",
			}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.RadioButtonsGroup, "dancing", "Dancing").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.RadioButtonsGroup, "solid", "Solid").Layout(gtx)
				}),
			)
		}),
	)
}
