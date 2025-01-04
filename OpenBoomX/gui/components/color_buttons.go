package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

type ColorButtons struct {
	buttons        []*ColorButton
	list           *widget.List
	buttonDims     layout.Dimensions
	OnColorClicked func(color color.NRGBA)
	ButtonsPerRow  int
}

type ColorButton struct {
	clickable widget.Clickable
	color     color.NRGBA
}

func CreateColorButtons(colors []color.NRGBA, buttonsPerRow int, onColorClicked func(color color.NRGBA)) *ColorButtons {
	return &ColorButtons{
		buttons:        createColorButtons(colors),
		list:           &widget.List{List: layout.List{Axis: layout.Vertical}},
		OnColorClicked: onColorClicked,
		ButtonsPerRow:  buttonsPerRow,
	}
}

func (cb *ColorButtons) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	var buttons = cb.buildColorButtons(th, gtx)

	// limit to 2 rows
	if cb.buttonDims.Size.Y != 0 {
		gtx.Constraints.Max.Y = cb.buttonDims.Size.Y * 2
	}

	return material.List(th, cb.list).Layout(gtx, cb.getRowsNum(buttons),
		func(gtx layout.Context, index int) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx, cb.buildButtonColumns(index, buttons)...)
		})

}

func (cb *ColorButtons) buildColorButtons(th *material.Theme, gtx layout.Context) []layout.FlexChild {
	var buttons = make([]layout.FlexChild, len(cb.buttons))
	for i, btn := range cb.buttons {
		if btn.clickable.Clicked(gtx) {
			cb.OnColorClicked(btn.color)
		}

		btnStyle := material.IconButton(th, &btn.clickable, nil, "")
		btnStyle.Inset = layout.UniformInset(4)
		btnStyle.Background = btn.color
		buttons[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			cb.buttonDims = layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return btnStyle.Layout(gtx)
			})
			return cb.buttonDims
		})
	}

	return buttons
}

func (cb *ColorButtons) buildButtonColumns(rowIndex int, buttons []layout.FlexChild) []layout.FlexChild {
	var columns []layout.FlexChild
	emptyColumn := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return cb.buttonDims
	})

	for col := 0; col < cb.ButtonsPerRow; col++ {
		var columnButtons []layout.FlexChild
		index := rowIndex*cb.ButtonsPerRow + col
		if index < len(buttons) {
			columnButtons = append(columnButtons, buttons[index])
		} else {
			columnButtons = append(columnButtons, emptyColumn)
		}

		columns = append(columns, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, columnButtons...)
		}))
	}

	return columns
}

func (cb *ColorButtons) getRowsNum(buttons []layout.FlexChild) int {
	return (len(buttons) + cb.ButtonsPerRow - 1) / cb.ButtonsPerRow
}

func createColorButtons(colors []color.NRGBA) []*ColorButton {
	buttons := make([]*ColorButton, len(colors))
	for i, c := range colors {
		buttons[i] = &ColorButton{
			color: c,
		}
	}
	return buttons
}

func (cb *ColorButtons) OnColorListChanged(colors []color.NRGBA) {
	cb.buttons = createColorButtons(colors)
}
