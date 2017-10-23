package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func killOldParticles(slice []Particle, win *pixelgl.Window) []Particle {
	i := 0
	for i < len(slice) {
		if slice[i].alive >= slice[i].lifespan ||
			slice[i].pos.X < win.Bounds().Min.X ||
			slice[i].pos.X > win.Bounds().Max.X ||
			slice[i].pos.Y < win.Bounds().Min.Y {
			slice = append(slice[:i], slice[i+1:]...)
		} else {
			i++
		}
	}
	return slice
}
