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
	radioButtonsGroup widget.Enum
	OnColorChanged    func(color color.NRGBA, solidColor bool)
}

// keys for radio buttons
const (
	dancingLights = "dancing"
	solidLights   = "solid"
)

func CreateLightPicker(onColorChanged func(color color.NRGBA, solidColor bool)) *LightPicker {
	picker := &LightPicker{}
	picker.picker.SetColor(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	picker.radioButtonsGroup.Value = dancingLights
	picker.OnColorChanged = onColorChanged
	return picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.picker.Update(gtx) || lp.radioButtonsGroup.Update(gtx) {
		lp.sendUpdate()
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
					return material.RadioButton(th, &lp.radioButtonsGroup, dancingLights, "Dancing").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.RadioButton(th, &lp.radioButtonsGroup, solidLights, "Solid").Layout(gtx)
				}),
			)
		}),
	)
}

func (lp *LightPicker) SetColor(color color.NRGBA) {
	lp.picker.SetColor(color)
	lp.sendUpdate()
}

func (lp *LightPicker) GetColor() color.NRGBA {
	return lp.picker.Color()
}

func (lp *LightPicker) sendUpdate() {
	solidColor := false
	if lp.radioButtonsGroup.Value == solidLights {
		solidColor = true
	}
	lp.OnColorChanged(lp.picker.Color(), solidColor)
}
