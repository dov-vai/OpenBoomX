// thanks: https://github.com/ankitkmrpatel/go-android-app-test/blob/master/internal/ui/components/snackbar.go

package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"obx/gui/theme"
	"time"
)

type Snackbar struct {
	message  string
	visible  bool
	timeout  time.Duration
	dismiss  widget.Clickable
	showTime time.Time
}

func CreateSnackbar() *Snackbar {
	return &Snackbar{}
}

func (s *Snackbar) show(message string, timeout time.Duration) {
	s.message = message
	s.timeout = timeout
	s.showTime = time.Now()
	s.visible = true
}

func (s *Snackbar) ShowMessageWithTimeout(message string, timeout time.Duration) {
	s.show(message, timeout)
}

func (s *Snackbar) ShowMessage(message string) {
	s.show(message, 3*time.Second)
}

func (s *Snackbar) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	if !s.visible {
		return layout.Dimensions{}
	}

	if time.Since(s.showTime) > s.timeout {
		s.visible = false
		return layout.Dimensions{}
	}

	// dismiss on click
	if s.dismiss.Clicked(gtx) {
		s.visible = false
		return layout.Dimensions{}
	}

	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							constraints := gtx.Constraints.Min
							constraints.X = gtx.Constraints.Max.X

							bounds := image.Rectangle{
								Max: constraints,
							}

							surfaceColor := theme.Surface0Color
							surfaceColor.A = 230

							paint.FillShape(gtx.Ops,
								surfaceColor,
								clip.UniformRRect(bounds, 4).Op(gtx.Ops),
							)

							return layout.Dimensions{Size: bounds.Max}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									button := s.dismiss.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										label := material.Body1(th, s.message)
										label.Color = th.Fg
										label.Alignment = text.Middle
										return label.Layout(gtx)
									})
									return button
								},
							)
						}),
					)
				},
			)
		}),
	)
}
