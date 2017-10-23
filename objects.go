package main

import "github.com/faiface/pixel"

// Gravity is vector representing standard gravity force
var Gravity = pixel.Vec{
	X: 0.0,
	Y: -9.81,
}

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

func (p *Particle) addGravity(dt float64) {
	p.speed = p.speed.Add(Gravity.Scaled(dt * 100))
}

// ParticleSystem represents system of particles with and rate of particle generation per second
type ParticleSystem struct {
	pos       pixel.Vec
	pps       int
	angle     float64
	particles []Particle
}
