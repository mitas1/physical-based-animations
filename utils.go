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
