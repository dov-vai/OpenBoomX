package components

import (
	"fmt"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"log"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/services"
)

type PresetButtons struct {
	presetService *services.EqPresetService
	list          widget.List
	presetButtons []*PresetButton
}

type PresetButton struct {
	widget.Clickable
	Title        string
	RemoveButton widget.Clickable
}

var DeleteIcon *widget.Icon

func CreatePresetButtons(presetService *services.EqPresetService) *PresetButtons {
	// TODO: refactor, maybe should have some separate package with icons defined and loaded?
	var err error
	DeleteIcon, err = widget.NewIcon(icons.ActionDelete)
	if err != nil {
		log.Fatalf("Failed to create delete icon: %v", err)
	}

	return &PresetButtons{
		presetService: presetService,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
}

func (pb *PresetButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	presetTitles := pb.presetService.ListPresets()

	// Update the preset buttons if the list of presets has changed
	if len(pb.presetButtons) != len(presetTitles) {
		pb.presetButtons = make([]*PresetButton, len(presetTitles))
		for i, title := range presetTitles {
			pb.presetButtons[i] = &PresetButton{Title: title}
		}
	}

	activePreset := pb.presetService.GetActivePreset()

	for _, btn := range pb.presetButtons {
		if btn.Clicked(gtx) {
			if err := pb.presetService.SetActivePreset(btn.Title); err != nil {
				fmt.Printf("Error setting active preset: %v\n", err)
			}
		}
		if btn.RemoveButton.Clicked(gtx) {
			if err := pb.presetService.DeletePreset(btn.Title); err != nil {
				fmt.Printf("Error deleting preset: %v\n", err)
			}
		}
	}

	inactiveColor := color.NRGBA{R: 0x95, G: 0xb1, B: 0xb0, A: 0xff}
	removeColor := color.NRGBA{R: 0x8b, G: 0x1c, B: 0x00, A: 0xff}

	// FIXME: list order changes when removed or added, hashmap related
	return material.List(th, &pb.list).Layout(gtx, len(pb.presetButtons), func(gtx layout.Context, index int) layout.Dimensions {
		btn := pb.presetButtons[index]
		return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btnStyle := material.Button(th, &btn.Clickable, btn.Title)
					if btn.Title != activePreset {
						btnStyle.Background = inactiveColor
					}

					return btnStyle.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btnStyle := material.IconButton(th, &btn.RemoveButton, DeleteIcon, "Remove")
						btnStyle.Background = removeColor

						return btnStyle.Layout(gtx)
					})
				}),
			)
		})
	})
}
