package components

import (
	"log"
	"obx/protocol"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type EqButtons struct {
	Buttons []EQButton
}

type EQButton struct {
	Mode      string
	Clickable widget.Clickable
}

func CreateEQButtons() []EQButton {
	var buttons []EQButton
	for mode := range protocol.EQModes {
		buttons = append(buttons, EQButton{Mode: mode})
	}
	return buttons
}

func (eq *EqButtons) LayoutEQButtons(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var buttons []layout.FlexChild

	for i := range eq.Buttons {
		btn := &eq.Buttons[i]
		for btn.Clickable.Clicked(gtx) {
			// TODO: implement connecting to speaker
			err := protocol.SetOluvMode(btn.Mode, "")
			if err != nil {
				log.Printf("Failed to set mode %s: %v", btn.Mode, err)
			} else {
				log.Printf("Mode set to: %s", btn.Mode)
			}
		}

		layout := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.Clickable, btn.Mode).Layout(gtx)
		})

		buttons = append(buttons, layout)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, buttons...)
}
