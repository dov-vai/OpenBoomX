package theme

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

// based on Catpuccin Mocha

var BaseColor = color.NRGBA{R: 0x1e, G: 0x1e, B: 0x2e, A: 0xff}
var TextColor = color.NRGBA{R: 0xcd, G: 0xd6, B: 0xf4, A: 0xff}
var MauveColor = color.NRGBA{R: 0xcb, G: 0xa6, B: 0xf7, A: 0xff}
var CrustColor = color.NRGBA{R: 0x11, G: 0x11, B: 0x1b, A: 0xff}
var MantleColor = color.NRGBA{R: 0x18, G: 0x18, B: 0x25, A: 0xff}
var Surface0Color = color.NRGBA{R: 0x31, G: 0x32, B: 0x44, A: 0xff}
var WarningColor = color.NRGBA{R: 0xe2, G: 0x6e, B: 0x8e, A: 0xff}

var Palette = material.Palette{
	// background (surface)
	Bg: BaseColor,
	// text
	Fg: TextColor,
	// components color
	ContrastBg: MauveColor,
	// elements text
	ContrastFg: TextColor,
}

var ButtonPalette = material.Palette{
	Bg:         BaseColor,
	Fg:         TextColor,
	ContrastBg: MantleColor,
	ContrastFg: TextColor,
}

var DeleteIcon, _ = widget.NewIcon(icons.ActionDelete)
var BatteryIcon, _ = widget.NewIcon(icons.DeviceBattery50)
var AddIcon, _ = widget.NewIcon(icons.ContentAdd)
var TuneIcon, _ = widget.NewIcon(icons.ImageTune)
var LightIcon, _ = widget.NewIcon(icons.ActionLightbulbOutline)
var ListIcon, _ = widget.NewIcon(icons.ActionList)
var StarIcon, _ = widget.NewIcon(icons.ToggleStar)
var SettingsIcon, _ = widget.NewIcon(icons.ActionSettings)
