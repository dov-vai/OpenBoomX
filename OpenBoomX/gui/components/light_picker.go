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
	OnColorChanged    func(color color.NRGBA, solidColor bool)
}

func CreateLightPicker(onColorChanged func(color color.NRGBA, solidColor bool)) *LightPicker {
	picker := &LightPicker{}
	picker.currentColor = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	picker.picker.SetColor(picker.currentColor)
	picker.radioButtonsGroup.Value = "dancing"
	picker.OnColorChanged = onColorChanged
	return picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.picker.Update(gtx) || lp.radioButtonsGroup.Update(gtx) {
		lp.currentColor = lp.picker.Color()
		solidColor := false
		if lp.radioButtonsGroup.Value == "solid" {
			solidColor = true
		}
		lp.OnColorChanged(lp.currentColor, solidColor)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
					return material.RadioButton(th, &lp.radioButtonsGroup, "dancing", "Dancing").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.radioButtonsGroup, "solid", "Solid").Layout(gtx)
				}),
			)
		}),
	)
}
