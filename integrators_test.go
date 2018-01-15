package main

import (
	"testing"

	"github.com/faiface/pixel"
)

func createParticle(
	pos pixel.Vec,
	nextPos pixel.Vec,
	speed pixel.Vec,
	prevDT float64,
	lifespan float64) Particle {
	particleImage, err := loadPicture("assets/sprites/particle.png")
	if err != nil {
		panic(err)
	}

	particleSprite := pixel.NewSprite(
		particleImage,
		pixel.R(
			particleImage.Bounds().Min.X,
			particleImage.Bounds().Min.Y,
			particleImage.Bounds().Max.X,
			particleImage.Bounds().Max.Y,
		))

	p := Particle{
		position:     pos,
		nextPosition: nextPos,
		speed:        speed,
		prevDt:       prevDT,
		sprite:       *particleSprite,
		lifespan:     lifespan,
		alive:        0.0,
	}

	return p
}

// TestP01 tests particle initialization
func TestP01(t *testing.T) {
	var (
		pos           = pixel.V(0, 0)
		nextPos       = pixel.V(0, 0)
		speed         = pixel.V(0, 10)
		prevDt        = 0.0
		lifespan      = 10.0
		ePosition     = pixel.V(0, 0)
		eNextPosition = pixel.V(0, 0)
		eSpeed        = pixel.V(0, 10)
		eLifespan     = 10.0
	)

	p := createParticle(pos, nextPos, speed, prevDt, lifespan)

	if p.position != ePosition {
		t.Errorf("Particle initialization: Expected position of %f got %f", ePosition, p.position)
	}

	if p.nextPosition != eNextPosition {
		t.Errorf(
			"Particle initialization: Expected next position of %f got %f",
			eNextPosition, p.nextPosition,
		)
	}

	if p.speed != eSpeed {
		t.Errorf("Particle initialization: Expected speed of %f got %f", eSpeed, p.speed)
	}

	if p.lifespan != eLifespan {
		t.Errorf("Particle initialization: Expected lifespan of %f got %f", eLifespan, p.lifespan)
	}
}

// TestIee test Explicit euler position integration method
func TestIee(t *testing.T) {
	var (
		pos       = pixel.V(0, 0)
		nextPos   = pixel.V(0, 0)
		speed     = pixel.V(0, 10)
		prevDT    = 0.
		lifespan  = 10.0
		eSpeed    = pixel.V(0, 0.1899999999999995)
		ePosition = pixel.V(0, 18.99999999999995)
	)

	p := createParticle(pos, nextPos, speed, prevDT, lifespan)

	p.position = p.ExplicitEulerIntegrator(1)

	if p.speed.X != eSpeed.X || p.speed.Y != eSpeed.Y {
		t.Errorf("Explicit Euler Integrator DT=1: Expected speed of %f got %f", eSpeed, p.speed)
	}

	if p.position.X != ePosition.X || p.position.Y != ePosition.Y {
		t.Errorf(
			"Explicit Euler Integrator DT=1: Expected position of %f got %f", ePosition, p.position,
		)
	}

	p.position = p.ExplicitEulerIntegrator(0)

	if p.speed.X != eSpeed.X || p.speed.Y != eSpeed.Y {
		t.Errorf("Explicit Euler Integrator DT=0: Expected speed of %f got %f", eSpeed, p.speed)
	}

	if p.position.X != ePosition.X || p.position.Y != ePosition.Y {
		t.Errorf(
			"Explicit Euler Integrator DT=0: Expected position of %f got %f", ePosition, p.position,
		)
	}
}

// TestIv tests Verlet position integration method
func TestIv(t *testing.T) {
	var (
		pos            = pixel.V(0, 0)
		nextPos        = pixel.V(0, 1)
		speed          = pixel.V(0, 0)
		prevDT         = 1.0
		lifespan       = 10.0
		ePosition1     = pixel.V(0, 1)
		eNextPosition1 = pixel.V(0, -979)
		ePosition2     = pixel.V(0, -979)
		eNextPosition2 = pixel.V(0, -2940)
	)

	p := createParticle(pos, nextPos, speed, prevDT, lifespan)

	p.position = p.VerletIntegrator(1)

	if p.position.X != ePosition1.X || p.position.Y != ePosition1.Y {
		t.Errorf("Verlet Integrator DT=1: Expected position of %f got %f", ePosition1, p.position)
	}

	if p.nextPosition.X != eNextPosition1.X || p.nextPosition.Y != eNextPosition1.Y {
		t.Errorf(
			"Verlet Integrator DT=1: Expected next position of %f got %f",
			eNextPosition1, p.nextPosition,
		)
	}

	p.position = p.VerletIntegrator(1)

	if p.position.X != ePosition2.X || p.position.Y != ePosition2.Y {
		t.Errorf("Verlet Integrator DT=1: Expected position of %f got %f", ePosition1, p.position)
	}

	if p.nextPosition.X != eNextPosition2.X || p.nextPosition.Y != eNextPosition2.Y {
		t.Errorf(
			"Verlet Integrator DT=1: Expected next position of %f got %f",
			eNextPosition1, p.nextPosition,
		)
	}
}
