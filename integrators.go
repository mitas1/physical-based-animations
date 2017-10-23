package main

import "github.com/faiface/pixel"

// ExplicitEulerIntegrator calculates new position of a particle based on it's previous position
// and it's speed
func ExplicitEulerIntegrator(position pixel.Vec, speed pixel.Vec, dt float64) pixel.Vec {
	return position.Add(speed.Scaled(dt))
}
