package main

import "github.com/faiface/pixel"

// PositionIntegrationMethod is enum for choosing method for particle position integration
type PositionIntegrationMethod int

const (
	// ExplicitEuler is a position integration method based on Explicit Euler method
	ExplicitEuler PositionIntegrationMethod = iota
	// MidPoint is a position integration method based on Mid Point algorithm
	MidPoint PositionIntegrationMethod = iota
	// Verlet is a position integration method based on Verlet's algorithm
	Verlet PositionIntegrationMethod = iota
)

// Particle represents particle object
type Particle struct {
	pos      pixel.Vec
	speed    pixel.Vec
	sprite   pixel.Sprite
	lifespan float64
	alive    float64
}

// ParticleSystem represents system of particles with and rate of particle generation per second
type ParticleSystem struct {
	pos       pixel.Vec
	pps       int
	angle     float64
	particles []Particle
}
