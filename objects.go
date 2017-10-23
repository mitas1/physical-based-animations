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
	mass     float64
	sprite   pixel.Sprite
	lifespan float64
	alive    float64
}

func (p *Particle) addGravity(dt float64) {
	p.speed = p.speed.Add(Gravity.Scaled(dt * p.mass)) // v_{t+1} = v_{t} + h*(F/m)
}

func (particleSystem *ParticleSystem) killOldParticles(minX float64, maxX float64, minY float64) {
	var aliveParticles []Particle
	for _, particle := range particleSystem.particles {
		if particle.alive < particle.lifespan &&
			particle.pos.X >= minX &&
			particle.pos.X <= maxX &&
			particle.pos.Y >= minY {
			aliveParticles = append(aliveParticles, particle)
		}
	}
	particleSystem.particles = append([]Particle{}, aliveParticles...)
}

// ParticleSystem represents system of particles with and rate of particle generation per second
type ParticleSystem struct {
	pos       pixel.Vec
	pps       int
	angle     float64
	particles []Particle
}
