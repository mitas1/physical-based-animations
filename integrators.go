package main

import "github.com/faiface/pixel"

// ExplicitEulerIntegrator calculates new position of a particle based on it's previous position
// and it's speed
func (p *Particle) ExplicitEulerIntegrator(dt float64) pixel.Vec {
	// v_{t+1} = v_{t} + h*(F/m) | v_{t+1} = v_{t} + h*g
	p.speed = p.speed.Add(Gravity.Scaled(dt))

	// p_{t+1} = p_{t} + h*v(t)
	return p.position.Add(p.speed.Scaled(dt).Scaled(PixelsPerMeter))
}

// ExplicitMidpointIntegrator calculates new position of a particle based on it's previous position
// and it's speed
func (p *Particle) ExplicitMidpointIntegrator(dt float64) pixel.Vec {
	// v_{t+(1/2)} = v_{t} + (h/2)*(F/m) | v_{t+(1/2)} = v_{t} + (h/2)*g
	p.speed = p.speed.Add(Gravity.Scaled(dt / 2))
	// p_{t+(1/2)} = p_{t} + (h/2)*v(t)
	position := p.position.Add(p.speed.Scaled(dt / 2).Scaled(PixelsPerMeter))

	p.speed = p.speed.Add(Gravity.Scaled(dt / 2))
	return position.Add(p.speed.Scaled(dt / 2).Scaled(PixelsPerMeter))
}

// VerletIntegrator calculates new position of a particle based on Verlet Integration Scheme
func (p *Particle) VerletIntegrator(dt float64) pixel.Vec {
	// While calculating next position using Verlet Integration scheme with changing time-step (Δt)
	// variable, Verlet scheme does not approximate the solution to the differencial equation.
	// This can be corrected using the following formula, where iteration rule becomes:
	// p_{t+1} = p_{t} + (p_{t} - p_{t-1}) * h_{i} / h_{i-1} + g * ((h_{i} + h_{i-1}) * h_{i}) / 2
	pNext := p.nextPosition.Add(
		p.nextPosition.Sub(p.position).Scaled(dt / p.prevDt)).Add(
		Gravity.Scaled(PixelsPerMeter).Scaled((dt + p.prevDt) * dt / 2))
	tmp := p.nextPosition
	p.prevDt = dt
	p.nextPosition = pNext
	p.speed = p.speed.Add(Gravity.Add(Gravity).Scaled(0.5).Scaled(dt))

	return tmp
}
