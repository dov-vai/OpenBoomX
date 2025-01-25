package pages

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LoadingPage struct {
	buttonTheme       *material.Theme
	retryConnection   widget.Clickable
	onRetryConnection func()
	appError          error
}

func NewLoadingPage(buttonTheme *material.Theme, onRetryConnection func()) *LoadingPage {
	return &LoadingPage{
		buttonTheme:       buttonTheme,
		onRetryConnection: onRetryConnection,
	}
}

func (l *LoadingPage) Layout(gtx layout.Context) layout.Dimensions {
	if l.retryConnection.Clicked(gtx) {
		l.appError = nil
		l.onRetryConnection()
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				text := "Loading..."
				if l.appError != nil {
					text = l.appError.Error()
				}
				label := material.H5(l.buttonTheme, text)
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 8}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if l.appError != nil {
					return layout.Dimensions{}
				}

				gtx.Constraints.Max.X = gtx.Dp(32)
				gtx.Constraints.Max.Y = gtx.Dp(32)
				return material.Loader(l.buttonTheme).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if l.appError == nil {
					return layout.Dimensions{}
				}
				return material.Button(l.buttonTheme, &l.retryConnection, "Retry").Layout(gtx)
			}),
		)
	})
}

func (l *LoadingPage) SetError(err error) {
	l.appError = err
}
