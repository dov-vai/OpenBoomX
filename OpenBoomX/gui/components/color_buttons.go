package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

type ColorButtons struct {
	presetButtons  []widget.Clickable
	colors         []color.NRGBA
	OnColorClicked func(color color.NRGBA)
	ButtonsPerRow  int
}

func CreateColorButtons(colors []color.NRGBA, buttonsPerRow int, onColorClicked func(color color.NRGBA)) *ColorButtons {
	return &ColorButtons{
		presetButtons:  make([]widget.Clickable, len(colors)),
		colors:         colors,
		OnColorClicked: onColorClicked,
		ButtonsPerRow:  buttonsPerRow,
	}
}

func (cb *ColorButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var buttons = cb.buildColorButtons(th, gtx)

	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx, cb.buildButtonColumns(buttons)...)
}

func (cb *ColorButtons) buildColorButtons(th *material.Theme, gtx layout.Context) []layout.FlexChild {
	var buttons = make([]layout.FlexChild, len(cb.presetButtons))
	for i, c := range cb.colors {
		if cb.presetButtons[i].Clicked(gtx) {
			cb.OnColorClicked(c)
		}

		btnStyle := material.IconButton(th, &cb.presetButtons[i], nil, "")
		btnStyle.Inset = layout.UniformInset(4)
		btnStyle.Background = c
		buttons[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return btnStyle.Layout(gtx)
			})
		})
	}

	return buttons
}

func (cb *ColorButtons) buildButtonColumns(buttons []layout.FlexChild) []layout.FlexChild {
	numRows := (len(buttons) + cb.ButtonsPerRow - 1) / cb.ButtonsPerRow
	var columns []layout.FlexChild

	for col := 0; col < cb.ButtonsPerRow; col++ {
		var columnButtons []layout.FlexChild
		for row := 0; row < numRows; row++ {
			index := row*cb.ButtonsPerRow + col
			if index < len(buttons) {
				columnButtons = append(columnButtons, buttons[index])
			}
		}

		columns = append(columns, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, columnButtons...)
		}))
	}

	return columns
}
