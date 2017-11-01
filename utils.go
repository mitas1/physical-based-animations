package main

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/faiface/pixel"
)

func loadPicture(data []byte) (pixel.Picture, error) {
	buf := new(bytes.Buffer)
	buf.Write(data)
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
