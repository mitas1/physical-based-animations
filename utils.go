package main

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/faiface/pixel"
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
