package components

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type EqSaveButton struct {
	Clickable widget.Clickable
	OnSaved   func(title string)
	Editor    widget.Editor
}

func CreateEqSaveButton(onSaved func(title string)) *EqSaveButton {
	return &EqSaveButton{
		OnSaved: onSaved,
		Editor:  widget.Editor{SingleLine: true, Submit: true},
	}
}

func (btn *EqSaveButton) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if e, ok := btn.Editor.Update(gtx); ok {
		if _, ok := e.(widget.SubmitEvent); ok {
			btn.OnSaved(btn.Editor.Text())
		}
	}

	if btn.Clickable.Clicked(gtx) {
		btn.OnSaved(btn.Editor.Text())
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, &btn.Editor, "Preset title").Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Button(th, &btn.Clickable, "Save").Layout(gtx)
		}),
	)
}

func (btn *EqSaveButton) SetText(text string) {
	btn.Editor.SetText(text)
}
