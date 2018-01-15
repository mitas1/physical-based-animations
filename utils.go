package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func loadPicture(imagePath string) (pixel.Picture, error) {
	data, err := Asset(imagePath)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.Write(data)
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func isInsideBoundingBox(x, y float64, box pixel.Rect, boxPosition pixel.Vec) bool {
	x1, y1, x2, y2 := box.Min.X, box.Min.Y,
		box.Max.X, box.Max.Y

	x0, y0 := boxPosition.XY()

	return x > x0+x1 && x < x0+x2 && y > y0+y1 && y < y0+y2
}

func (circle *Circle) isPositionInside(position pixel.Vec) bool {
	return circle.position.To(position).Len() <= circle.radius
}

func (circle *Circle) draw(imd *imdraw.IMDraw, position pixel.Vec) {
	imd.Color = color.RGBA{0, 0, 0, 30}
	imd.Push(position)
	imd.Circle(circle.radius, 0)
	imd.Color = color.RGBA{0, 0, 0, 50}
	imd.Push(position)
	imd.Circle(circle.radius, 1)
}
