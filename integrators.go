package main

import "github.com/faiface/pixel"

// ExplicitEulerIntegrator calculates new position of a particle based on it's previous position
// and it's speed
func (p *Particle) ExplicitEulerIntegrator(dt float64) pixel.Vec {
	return p.position.Add(p.speed.Scaled(dt)) // p_{t+1} = p_{t} + h*v(t)
}
