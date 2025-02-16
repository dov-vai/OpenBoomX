package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"obx/gui/theme"
)

type EqSaveButton struct {
	clickable widget.Clickable
	OnSaved   func(title string)
	editor    widget.Editor
}

func CreateEqSaveButton(onSaved func(title string)) *EqSaveButton {
	return &EqSaveButton{
		OnSaved: onSaved,
		editor:  widget.Editor{SingleLine: true, Submit: true},
	}
}

func (btn *EqSaveButton) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if e, ok := btn.editor.Update(gtx); ok {
		if _, ok := e.(widget.SubmitEvent); ok {
			btn.OnSaved(btn.editor.Text())
		}
	}

	if btn.clickable.Clicked(gtx) {
		btn.OnSaved(btn.editor.Text())
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			surfaceStyle := component.Surface(
				&material.Theme{
					Palette: material.Palette{
						Bg: theme.Surface0Color,
					},
				})

			surfaceStyle.CornerRadius = 4

			return surfaceStyle.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(8).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							editor := material.Editor(th, &btn.editor, "Preset title")
							selectionColor := theme.MauveColor
							selectionColor.A = selectionColor.A * 0x60
							editor.SelectionColor = selectionColor
							return editor.Layout(gtx)
						}))
				})
			})
		}),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.clickable, "Save").Layout(gtx)
		}),
	)
}

func (btn *EqSaveButton) SetText(text string) {
	btn.editor.SetText(text)
}

func (btn *EqSaveButton) OnPresetChanged(newPreset string, values []float32) {
	btn.SetText(newPreset)
}
