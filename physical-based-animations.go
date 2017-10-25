package main

import (
	"errors"
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func updateParticles(particles []Particle, batch *pixel.Batch, dt float64, cam pixel.Matrix) {
	for i := 0; i < len(particles); i++ {
		newPos, err := newPosition(particles[i], dt, ExplicitEuler)
		if err != nil {
			fmt.Println(err.Error())
		}
		particles[i].position = newPos
		particles[i].alive += dt
		particles[i].sprite.Draw(batch, pixel.IM.Moved(cam.Unproject(particles[i].position)))
		particles[i].AddGravity(dt)
	}
}

func newPosition(particle Particle, dt float64, mode PositionIntegrationMethod) (pixel.Vec, error) {
	switch mode {
	case ExplicitEuler:
		return particle.ExplicitEulerIntegrator(dt), nil
	case MidPoint:
		return particle.position, errors.New("Unimplemented")
	case Verlet:
		return particle.position, errors.New("Unimplemented")
	default:
		return particle.position, errors.New("Unknown method of position integration")
	}

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Particle System",
		Bounds: pixel.R(0, 0, 1024, 768),
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	particleImage, err := loadPicture("assets/sprites/particle.png")
	if err != nil {
		panic(err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, particleImage)
	particleSprite := pixel.NewSprite(
		particleImage,
		pixel.R(
			particleImage.Bounds().Min.X,
			particleImage.Bounds().Min.Y,
			particleImage.Bounds().Max.X,
			particleImage.Bounds().Max.Y,
		))

	var (
		camPos      = pixel.ZV
		second      = time.Tick(time.Second)
		frames      = 0
		timeElapsed = 0.0
	)

	particleSystem := ParticleSystem{
		position: win.Bounds().Center().Sub(pixel.V(0.0, win.Bounds().H()/4.0)),
		pps:      1000,
		angle:    60.0,
	}

	timeForOneParticle := 1.0 / float64(particleSystem.pps)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		timeElapsed += dt

		cam := pixel.IM.Scaled(camPos, 1.0).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		win.Clear(colornames.Whitesmoke)

		batch.Clear()
		updateParticles(particleSystem.particles, batch, dt, cam)

		batch.Draw(win)

		for timeElapsed > timeForOneParticle {
			angle := (rand.Float64() - 0.5) * (particleSystem.angle * (math.Pi / 180))
			particle := Particle{
				position: particleSystem.position,
				speed:    pixel.V(0, 600.0).Rotated(angle),
				mass:     50,
				sprite:   *particleSprite,
				lifespan: 10.0,
				alive:    0.0,
			}
			particleSystem.particles = append(particleSystem.particles, particle)
			timeElapsed = timeElapsed - timeForOneParticle
		}

		win.Update()

		frames++
		select {
		case <-second:
			particleSystem.KillOldParticles(
				win.Bounds().Min.X,
				win.Bounds().Max.X,
				win.Bounds().Min.X,
			)
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | particles %d", cfg.Title, frames, len(particleSystem.particles)))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
