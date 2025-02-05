package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"obx/gui/services"
	"obx/gui/theme"
)

type PresetButtons struct {
	presetService *services.EqPresetService
	list          widget.List
	presetButtons []*PresetButton
	snackbar      *Snackbar
}

type PresetButton struct {
	clickable    widget.Clickable
	title        string
	removeButton widget.Clickable
}

func CreatePresetButtons(presetService *services.EqPresetService, snackbar *Snackbar) *PresetButtons {
	return &PresetButtons{
		presetService: presetService,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		presetButtons: createPresetButtons(presetService.ListPresets()),
		snackbar:      snackbar,
	}
}

func (pb *PresetButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	activePreset := pb.presetService.GetActivePreset()

	for _, btn := range pb.presetButtons {
		if btn.clickable.Clicked(gtx) {
			if err := pb.presetService.SetActivePreset(btn.title); err != nil {
				fmt.Printf("Error setting active preset: %v\n", err)
				pb.snackbar.ShowMessage(fmt.Sprintf("Failed setting preset to: %s", btn.title))
			} else {
				pb.snackbar.ShowMessage(fmt.Sprintf("Successfully set preset to: %s", btn.title))
			}
		}
		if btn.removeButton.Clicked(gtx) {
			if err := pb.presetService.DeletePreset(btn.title); err != nil {
				fmt.Printf("Error deleting preset: %v\n", err)
				pb.snackbar.ShowMessage(fmt.Sprintf("Failed removing preset: %s", btn.title))
			} else {
				pb.snackbar.ShowMessage(fmt.Sprintf("Successfully removed preset: %s", btn.title))
			}
		}
	}

	return material.List(th, &pb.list).Layout(gtx, len(pb.presetButtons), func(gtx layout.Context, index int) layout.Dimensions {
		btn := pb.presetButtons[index]
		return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btnStyle := material.Button(th, &btn.clickable, btn.title)
					if btn.title != activePreset {
						btnStyle.Background = theme.Surface0Color
					}

					return btnStyle.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btnStyle := material.IconButton(th, &btn.removeButton, theme.DeleteIcon, "Remove")
						btnStyle.Inset = layout.UniformInset(4)
						btnStyle.Background = theme.WarningColor
						return btnStyle.Layout(gtx)
					})
				}),
			)
		})
	})
}

func createPresetButtons(presetTitles []string) []*PresetButton {
	presetButtons := make([]*PresetButton, len(presetTitles))
	for i, title := range presetTitles {
		presetButtons[i] = &PresetButton{title: title}
	}
	return presetButtons
}

func (pb *PresetButtons) OnPresetChanged(newPreset string, values []float32) {
	pb.presetButtons = createPresetButtons(pb.presetService.ListPresets())
}
