package components

import (
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op/clip"
	"image"
	"image/color"
	"image/draw"
	"math"

	"gioui.org/layout"
	"gioui.org/op/paint"
)

const eventTag = "ColorWheel"

type ColorWheel struct {
	prevMaxX        int
	image           *image.RGBA
	colorWheelImage *image.RGBA
	cursorImage     *image.RGBA
	OnColorChanged  func(color color.NRGBA)
}

func CreateColorWheel(onColorChanged func(color color.NRGBA)) *ColorWheel {
	return &ColorWheel{OnColorChanged: onColorChanged}
}

func (cw *ColorWheel) Layout(gtx layout.Context, radius float32) layout.Dimensions {
	// regenerate if window width changed
	if gtx.Constraints.Max.X != cw.prevMaxX {
		cw.colorWheelImage = drawColorWheel(radius)
		// FIXME: should maintain the cursor position when radius changes
		cw.image = cw.colorWheelImage
		cw.cursorImage = drawCursor(radius/12, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		cw.prevMaxX = gtx.Constraints.Max.X
	}

	gtx.Constraints.Min.Y = int(2 * radius)
	gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

	defer clip.Rect{Max: image.Pt(gtx.Constraints.Min.Y, gtx.Constraints.Min.Y)}.Push(gtx.Ops).Pop()

	// handle click and drag events
	event.Op(gtx.Ops, eventTag)
	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: eventTag,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag,
		})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			default:
				if rgb := cursorPosToRgb(x.Position, radius); rgb != nil {
					cw.OnColorChanged(*rgb)
					cw.image = drawCursorPosition(x.Position, cw.cursorImage, cw.colorWheelImage)
				}
			}
		}
	}

	// draw the color wheel
	imageOp := paint.NewImageOp(cw.image)
	imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: image.Point{X: gtx.Constraints.Min.Y, Y: gtx.Constraints.Min.Y},
	}
}

func drawCursorPosition(pos f32.Point, cursorImg *image.RGBA, imgData *image.RGBA) *image.RGBA {
	imgCopy := image.NewRGBA(imgData.Bounds())
	draw.Draw(imgCopy, imgCopy.Bounds(), imgData, imgData.Bounds().Min, draw.Src)

	cursorBounds := cursorImg.Bounds()
	halfWidth := float32(cursorBounds.Dx() / 2)
	halfHeight := float32(cursorBounds.Dy() / 2)

	destRect := image.Rect(
		int(pos.X-halfWidth),
		int(pos.Y-halfHeight),
		// size adjustment because it's getting cut off from the halving
		int(pos.X+halfWidth+2),
		int(pos.Y+halfHeight+2),
	)

	draw.Draw(imgCopy, destRect.Intersect(imgCopy.Bounds()), cursorImg, cursorBounds.Min, draw.Over)

	return imgCopy
}

func drawCursor(radius float32, col color.RGBA) *image.RGBA {
	size := int(2 * radius)
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	center := f32.Point{X: radius, Y: radius}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			dist := math.Sqrt(math.Pow(float64(x)-float64(center.X), 2) + math.Pow(float64(y)-float64(center.Y), 2))
			if dist <= float64(radius) {
				img.Set(x, y, col)
			}
		}
	}

	return img
}

func drawColorWheel(radius float32) *image.RGBA {
	diameter := 2 * radius

	imgData := image.NewRGBA(image.Rect(0, 0, int(diameter), int(diameter)))

	for x := -radius; x < radius; x++ {
		for y := -radius; y < radius; y++ {
			rgb := posToRgb(x, y, radius)
			if rgb == nil {
				continue
			}

			// offset to center the circle in the square image
			adjustedX := int(x + radius)
			adjustedY := int(y + radius)

			imgData.SetRGBA(adjustedX, adjustedY, color.RGBA{
				R: rgb.R,
				G: rgb.G,
				B: rgb.B,
				A: 255,
			})
		}
	}

	return imgData
}

func cursorPosToRgb(pos f32.Point, radius float32) *color.NRGBA {
	return posToRgb(pos.X-radius, pos.Y-radius, radius)
}

func posToRgb(x, y, radius float32) *color.NRGBA {
	r, phi := xyToPolar(x, y)

	// don't return color outside circle radius
	if r > radius {
		return nil
	}

	deg := radToDeg(phi)

	hue := float64(deg)
	saturation := float64(r / radius)
	value := 1.0
	red, green, blue := hsvToRgb(hue, saturation, value)

	return &color.NRGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 255,
	}
}

func xyToPolar(x, y float32) (float32, float32) {
	r := float32(math.Sqrt(float64(x*x + y*y)))
	phi := float32(math.Atan2(float64(y), float64(x)))
	return r, phi
}

func radToDeg(rad float32) float32 {
	return ((rad + math.Pi) / (2 * math.Pi)) * 360
}

func hsvToRgb(hue, saturation, value float64) (float64, float64, float64) {
	chroma := value * saturation
	hue1 := hue / 60.0
	x := chroma * (1 - math.Abs(math.Mod(hue1, 2)-1))
	m := value - chroma

	var r1, g1, b1 float64
	switch {
	case hue1 <= 1:
		r1, g1, b1 = chroma, x, 0
	case hue1 <= 2:
		r1, g1, b1 = x, chroma, 0
	case hue1 <= 3:
		r1, g1, b1 = 0, chroma, x
	case hue1 <= 4:
		r1, g1, b1 = 0, x, chroma
	case hue1 <= 5:
		r1, g1, b1 = x, 0, chroma
	case hue1 <= 6:
		r1, g1, b1 = chroma, 0, x
	}

	return 255 * (r1 + m), 255 * (g1 + m), 255 * (b1 + m)
}
