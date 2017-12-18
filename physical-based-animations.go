package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	winWidth  = 1024
	winHeight = 768
)

func updateParticles(particles []Particle, batch *pixel.Batch, dt float64, cam pixel.Matrix) {
	for i := 0; i < len(particles); i++ {
		newPos, err := newPosition(&particles[i], dt, MidPoint)
		if err != nil {
			fmt.Println(err.Error())
		}

		particles[i].position = newPos
		particles[i].alive += dt
		particles[i].sprite.Draw(batch, pixel.IM.Moved(cam.Unproject(particles[i].position)))
	}
}

func newPosition(particle *Particle, dt float64, mode PositionIntegrationMethod) (pixel.Vec, error) {
	switch mode {
	case ExplicitEuler:
		return particle.ExplicitEulerIntegrator(dt), nil
	case MidPoint:
		return particle.ExplicitMidpointIntegrator(dt), nil
	case Verlet:
		return particle.VerletIntegrator(dt), nil
	default:
		return particle.position, errors.New("Unknown method of position integration")
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Particle System",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
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

	emitRate := Parameter{
		value: 1000,
		step:  100,
		min:   0,
		max:   2200,
	}

	emitAngle := Parameter{
		value: 60,
		step:  5,
		min:   10,
		max:   360,
	}

	particleLife := Parameter{
		value: 2,
		step:  0.1,
		min:   0.1,
		max:   4,
	}

	initialVelocity := Parameter{
		value: 9.5,
		step:  0.5,
		min:   -2,
		max:   20,
	}

	particleSystem := ParticleSystem{
		position: win.Bounds().Center().Sub(pixel.V(0.0, win.Bounds().H()/4.0)),
		emitRate: &emitRate,
		angle:    &emitAngle,
	}

	prevDt := 0.002

	last := time.Now()

	gui := GUI{
		win:   win,
		atlas: text.NewAtlas(basicfont.Face7x13, text.ASCII),
	}

	gui.CreateBatch("assets/sprites/spritesheet.png")

	playButton := Button{
		position:     pixel.V(10, 10),
		croppingArea: pixel.R(0, 300, 60, 360),
		bounds:       pixel.R(0, 0, 60, 60),
		onClick:      handlePlayClick,
	}

	gui.NewButton(playButton)

	pauseButton := Button{
		position:     pixel.V(80, 10),
		croppingArea: pixel.R(60, 300, 120, 360),
		bounds:       pixel.R(0, 0, 60, 60),
		onClick:      handlePauseClick,
	}

	gui.NewButton(pauseButton)

	stopButton := Button{
		position:     pixel.V(160, 10),
		croppingArea: pixel.R(120, 300, 180, 360),
		bounds:       pixel.R(0, 0, 60, 60),
		onClick:      handleStopClick,
	}

	gui.NewButton(stopButton)

	emitRateSlider := SliderWannabe{
		y:         280,
		parameter: &emitRate,
		format:    "%.0f par/sec\n",
	}

	gui.NewSliderWannabe(emitRateSlider)

	emitAngleSlider := SliderWannabe{
		y:         370,
		parameter: &emitAngle,
		format:    "%.0f degrees\n",
	}

	gui.NewSliderWannabe(emitAngleSlider)

	particleLifeSlider := SliderWannabe{
		y:         460,
		parameter: &particleLife,
		format:    "lives %.1f s\n",
	}

	gui.NewSliderWannabe(particleLifeSlider)

	initialVelocitySlider := SliderWannabe{
		y:         550,
		parameter: &initialVelocity,
		format:    "%.1f m/s",
	}

	gui.NewSliderWannabe(initialVelocitySlider)

	cam := pixel.IM.Scaled(camPos, 1.0).Moved(win.Bounds().Center().Sub(camPos))

	win.SetMatrix(cam)

	gui.SetMatrix(cam)
	gui.BindState(&state)
	gui.MainLoop()
	gui.Draw()

	fps := time.Tick(time.Second / 200)

	for !win.Closed() {
		win.Update()

		if !gui.GetState().paused && !gui.GetState().stopped {
			dt := time.Since(last).Seconds()
			last = time.Now()
			timeElapsed += dt

			batch.Clear()
			updateParticles(particleSystem.particles, batch, dt, cam)

			win.Clear(colornames.Whitesmoke)

			gui.batch.Draw(win)
			batch.Draw(win)
			gui.DrawText(win)

			timeForOneParticle := 1.0 / float64(particleSystem.emitRate.value)

			for timeElapsed > timeForOneParticle {
				pos := particleSystem.position
				angle := (rand.Float64() - 0.5) * (particleSystem.angle.value * (math.Pi / 180))
				speed := pixel.V(0, initialVelocity.value).Rotated(angle)
				nextPost := pos.Add(speed.Scaled(PixelsPerMeter).Scaled(prevDt)).Add(
					Gravity.Scaled(PixelsPerMeter).Scaled(prevDt * prevDt * 0.5))

				particle := Particle{
					position:     pos,
					nextPosition: nextPost,
					speed:        speed,
					prevDt:       prevDt,
					sprite:       *particleSprite,
					lifespan:     particleLife.value,
					alive:        0.0,
				}
				particleSystem.particles = append(particleSystem.particles, particle)
				timeElapsed = timeElapsed - timeForOneParticle
			}
		} else if gui.GetState().paused && !gui.GetState().stopped {
			last = time.Now()
		} else {
			last = time.Now()
			timeElapsed = 0
			particleSystem.particles = particleSystem.particles[:0]
			batch.Clear()
			win.Clear(colornames.Whitesmoke)
			gui.batch.Draw(gui.win)
		}
		particleSystem.KillOldParticles(
			win.Bounds().Min.X,
			win.Bounds().Max.X,
			win.Bounds().Min.Y,
		)

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | particles %d", cfg.Title, frames,
				len(particleSystem.particles)))
			frames = 0
		default:
		}
		<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
