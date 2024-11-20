package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
	"image/color"
)

type LightPicker struct {
	Picker            colorpicker.State
	CurrentColor      color.NRGBA
	Solid             bool
	RadioButtonsGroup widget.Enum
	DefaultButton     widget.Clickable
	OffButton         widget.Clickable
	OnActionClicked   func(action string)
	OnColorChanged    func(color color.NRGBA, solidColor bool)
}

func CreateLightPicker(onActionClicked func(action string), onColorChanged func(color color.NRGBA, solidColor bool)) LightPicker {
	picker := &LightPicker{}
	picker.CurrentColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	picker.Picker.SetColor(picker.CurrentColor)
	picker.RadioButtonsGroup.Value = "dancing"
	picker.OnActionClicked = onActionClicked
	picker.OnColorChanged = onColorChanged
	return *picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.Picker.Update(gtx) {
		lp.CurrentColor = lp.Picker.Color()
		solidColor := false
		if lp.RadioButtonsGroup.Value == "solid" {
			solidColor = true
		}
		lp.OnColorChanged(lp.CurrentColor, solidColor)
	}
	if lp.OffButton.Clicked(gtx) {
		// TODO: fix hardcodes everywhere..
		lp.OnActionClicked("off")
	}
	if lp.DefaultButton.Clicked(gtx) {
		lp.OnActionClicked("default")
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &lp.DefaultButton, "Default").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Spacer{Width: unit.Dp(8)}.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &lp.OffButton, "Off").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return colorpicker.PickerStyle{
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
