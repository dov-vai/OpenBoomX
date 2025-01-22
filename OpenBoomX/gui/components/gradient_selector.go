package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"obx/gui/theme"
	"time"
)

type GradientSelector struct {
	firstColor         color.NRGBA
	firstColorToggled  bool
	secondColor        color.NRGBA
	secondColorToggled bool
	gradient           []color.NRGBA
	firstColorButton   widget.Clickable
	secondColorButton  widget.Clickable
	startButton        widget.Clickable
	started            bool
	durationMil        int
	stepTimeMil        int
	OnColorChanged     func(color color.NRGBA)
}

func CreateGradientSelector(onColorChanged func(color color.NRGBA)) *GradientSelector {
	gs := &GradientSelector{firstColor: theme.PeachColor, secondColor: theme.GreenColor, OnColorChanged: onColorChanged, durationMil: 2000, stepTimeMil: 5}
	gs.gradient = createGradient(gs.firstColor, gs.secondColor, gs.getSteps())
	return gs
}

func (gs *GradientSelector) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if gs.startButton.Clicked(gtx) {
		gs.started = !gs.started

		if gs.started {
			go gs.startGradient()
		}
	}

	if gs.firstColorButton.Clicked(gtx) {
		gs.firstColorToggled = !gs.firstColorToggled
		gs.secondColorToggled = false
	}

	if gs.secondColorButton.Clicked(gtx) {
		gs.secondColorToggled = !gs.secondColorToggled
		gs.firstColorToggled = false
	}

	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEvenly, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H6(th, "Gradient Effect").Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.IconButton(th, &gs.firstColorButton, nil, "")
			btnStyle.Inset = layout.UniformInset(4)
			btnStyle.Background = gs.firstColor
			if gs.firstColorToggled {
				btnStyle.Icon = theme.HelpIcon
			}

			return btnStyle.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(12)
			gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

			dr := image.Rectangle{Max: gtx.Constraints.Max}
			paint.LinearGradientOp{
				Stop1:  layout.FPt(dr.Min),
				Stop2:  layout.FPt(dr.Max),
				Color1: gs.firstColor,
				Color2: gs.secondColor,
			}.Add(gtx.Ops)
			defer clip.Rect(dr).Push(gtx.Ops).Pop()
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{
				Size: gtx.Constraints.Max,
			}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnStyle := material.IconButton(th, &gs.secondColorButton, nil, "")
			btnStyle.Inset = layout.UniformInset(4)
			btnStyle.Background = gs.secondColor
			if gs.secondColorToggled {
				btnStyle.Icon = theme.HelpIcon
			}

			return btnStyle.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnText := "Start"
			if gs.started {
				btnText = "Stop"
			}
			return material.Button(th, &gs.startButton, btnText).Layout(gtx)
		}),
	)
}

func (gs *GradientSelector) OnColorSelected(color color.NRGBA) {
	if gs.firstColorToggled {
		gs.firstColor = color
	}

	if gs.secondColorToggled {
		gs.secondColor = color
	}

	if gs.secondColorToggled || gs.firstColorToggled {
		gs.started = false
		gs.gradient = createGradient(gs.firstColor, gs.secondColor, gs.getSteps())
	}

	gs.secondColorToggled = false
	gs.firstColorToggled = false
}

func (gs *GradientSelector) getSteps() int {
	return gs.durationMil / gs.stepTimeMil
}

func (gs *GradientSelector) startGradient() {
	ticker := time.NewTicker(time.Duration(gs.stepTimeMil) * time.Millisecond)
	defer ticker.Stop()

	i := 0
	ascending := true
	numSteps := gs.getSteps()
	for range ticker.C {
		if !gs.started {
			break
		}

		gs.OnColorChanged(gs.gradient[i])

		if i == numSteps-1 {
			ascending = false
		}

		if i == 0 {
			ascending = true
		}

		if ascending {
			i++
		} else {
			i--
		}
	}

}

// createGradient generates a list of colors that represent the gradient.
func createGradient(c1, c2 color.NRGBA, steps int) []color.NRGBA {
	colors := make([]color.NRGBA, steps)
	for i := 0; i < steps; i++ {
		progress := float64(i) / float64(steps-1)
		colors[i] = interpolateColor(c1, c2, progress)
	}
	return colors
}

// interpolateColor linearly interpolates between two colors.
func interpolateColor(c1, c2 color.NRGBA, progress float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(c1.R) + (float64(c2.R)-float64(c1.R))*progress),
		G: uint8(float64(c1.G) + (float64(c2.G)-float64(c1.G))*progress),
		B: uint8(float64(c1.B) + (float64(c2.B)-float64(c1.B))*progress),
		A: 255,
	}
}
