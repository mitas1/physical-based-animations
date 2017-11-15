package main

import "github.com/faiface/pixel"

// Gravity is vector representing standard gravity force accelleration vector in m*s^{-2}
var Gravity = pixel.Vec{
	X: 0.0,
	Y: -9.81,
}

// PixelsPerMeter is the number of pixels on screen that represent one meter in real life
const PixelsPerMeter = 100.0

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
	position     pixel.Vec    // in pixels
	nextPosition pixel.Vec    // in pixels
	prevDt       float64      // in s
	speed        pixel.Vec    // in m*s^{-1}
	sprite       pixel.Sprite // particle display image
	lifespan     float64      // in s
	alive        float64      // in s
}

// KillOldParticles removes all particles that live up to their lifespan or are outside the
// boundaries of the view
func (particleSystem *ParticleSystem) KillOldParticles(minX float64, maxX float64, minY float64) {
	var aliveParticles []Particle
	for _, particle := range particleSystem.particles {
		if particle.alive < particle.lifespan &&
			particle.position.X >= minX &&
			particle.position.X <= maxX &&
			particle.position.Y >= minY {
			aliveParticles = append(aliveParticles, particle)
		}
	}
	particleSystem.particles = append([]Particle{}, aliveParticles...)
}

// ParticleSystem represents system of particles with and rate of particle generation per second
type ParticleSystem struct {
	position  pixel.Vec // in pixels
	pps       int
	angle     float64 // in degrees
	particles []Particle
}
