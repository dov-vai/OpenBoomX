package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
	"image/color"
	"log"
	"obx/protocol"
)

type LightPicker struct {
	Picker            colorpicker.State
	CurrentColor      color.NRGBA
	Solid             bool
	RadioButtonsGroup widget.Enum
	DefaultButton     widget.Clickable
	OffButton         widget.Clickable
}

func CreateLightPicker() LightPicker {
	picker := &LightPicker{}
	picker.CurrentColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	picker.Picker.SetColor(picker.CurrentColor)
	picker.RadioButtonsGroup.Value = "dancing"
	return *picker
}

func (lp *LightPicker) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if lp.Picker.Update(gtx) {
		lp.CurrentColor = lp.Picker.Color()
	}
	if lp.OffButton.Clicked(gtx) {
		err := protocol.HandleLightAction("off", false, "")
		if err != nil {
			log.Printf("Failed to set light off mode: %v", err)
		}
	}
	if lp.DefaultButton.Clicked(gtx) {
		err := protocol.HandleLightAction("default", false, "")
		if err != nil {
			log.Printf("Failed to set light default mode: %v", err)
		}
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
