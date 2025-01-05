package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"log"
	"obx/gui/constants"
	"obx/gui/ui"
	"os"
)

func main() {
	ui := ui.NewUI()

	defer ui.Dispose()

	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(constants.AppName),
			app.Size(unit.Dp(580), unit.Dp(450)),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
