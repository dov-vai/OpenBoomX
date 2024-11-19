package main

import (
	"log"
	"obx/gui/components"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	ui := newUI()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("OpenBoomX"),
			app.Size(unit.Dp(240), unit.Dp(300)),
		)
		if err := ui.run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

var defaultMargin = unit.Dp(10)

type UI struct {
	Theme       *material.Theme
	EqButtons   components.EqButtons
	LightPicker components.LightPicker
	BeepSlider  components.BeepSlider
}

func newUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	ui.EqButtons = components.CreateEQButtons()
	ui.LightPicker = components.CreateLightPicker()
	ui.BeepSlider = components.CreateBeepSlider()
	return ui
}

func (ui *UI) run(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.layout(gtx)
			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *UI) layout(gtx layout.Context) layout.Dimensions {
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return ui.EqButtons.Layout(ui.Theme, gtx)
			}),

			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return ui.LightPicker.Layout(ui.Theme, gtx)
			}),

			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return ui.BeepSlider.Layout(ui.Theme, gtx)
			}),
		)
	})
}
