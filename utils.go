package main

import (
	"image"
	"os"

	"github.com/faiface/pixel"
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

func killOldParticles(slice []Particle) []Particle {
	i := 0
	for i < len(slice) {
		if slice[i].alive >= slice[i].lifespan {
			slice = append(slice[:i], slice[i+1:]...)
		} else {
			i++
		}
	}
	return slice
}
