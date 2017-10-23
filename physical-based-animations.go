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
		particles[i].pos = newPos
		particles[i].alive += dt
		particles[i].sprite.Draw(batch, pixel.IM.Moved(cam.Unproject(particles[i].pos)))
		particles[i].addGravity(dt)
	}
}

func newPosition(particle Particle, dt float64, mode PositionIntegrationMethod) (pixel.Vec, error) {
	switch mode {
	case ExplicitEuler:
		return ExplicitEulerIntegrator(particle.pos, particle.speed, dt), nil
	case MidPoint:
		return particle.pos, errors.New("Unimplemented")
	case Verlet:
		return particle.pos, errors.New("Unimplemented")
	default:
		return particle.pos, errors.New("Unknown method of position integration")
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
		pos:   win.Bounds().Center().Sub(pixel.V(0.0, win.Bounds().H()/4.0)),
		pps:   1000,
		angle: 60.0,
	}

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

		if timeElapsed > 1.0/float64(particleSystem.pps) {
			angle := (rand.Float64() - 0.5) * (particleSystem.angle * (math.Pi / 180))
			particle := Particle{
				pos:      particleSystem.pos,
				speed:    pixel.V(0, 1000.0).Rotated(angle),
				sprite:   *particleSprite,
				lifespan: 10.0,
				alive:    0.0,
			}
			particleSystem.particles = append(particleSystem.particles, particle)
			timeElapsed = 1.0/float64(particleSystem.pps) - timeElapsed
		}

		win.Update()

		frames++
		select {
		case <-second:
			particleSystem.particles = killOldParticles(particleSystem.particles)
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | particles %d", cfg.Title, frames, len(particleSystem.particles)))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
