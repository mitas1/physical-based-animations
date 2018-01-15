package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel/imdraw"
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

func updateParticles(
	particles []Particle,
	batch *pixel.Batch,
	dt float64,
	cam pixel.Matrix,
	positionIntegrator PositionIntegrationMethod,
	circle Circle) {
	for i := 0; i < len(particles); i++ {
		newPosition, err := getNewPosition(&particles[i], dt, positionIntegrator)
		if err != nil {
			fmt.Println(err.Error())
		}

		if circle.isPositionInside(newPosition) {
			unitNormalVector := particles[i].position.Sub(circle.position).Unit().Scaled(circle.radius)
			unitSpeed := particles[i].speed.Unit()

			newPosition = unitNormalVector.Scaled(1.1).Add(circle.position)
			particles[i].speed = particles[i].speed.Rotated(2 *
				(math.Atan2(unitSpeed.Y, unitSpeed.X) -
					math.Atan2(unitNormalVector.Y, unitNormalVector.X)))
		}

		particles[i].position = newPosition
		particles[i].alive += dt
		particles[i].sprite.Draw(batch, pixel.IM.Moved(cam.Unproject(particles[i].position)))
	}
}

func getNewPosition(
	particle *Particle,
	dt float64,
	mode PositionIntegrationMethod) (pixel.Vec, error) {
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

	guiCanvasWidth := 320.0

	particleSystem := ParticleSystem{
		position: pixel.V((win.Bounds().W()+win.Bounds().Min.X+guiCanvasWidth)/2, win.Bounds().H()/4.0),
		emitRate: &emitRate,
		angle:    &emitAngle,
	}

	circle := Circle{
		position: pixel.V(412, 400),
		radius:   50,
	}

	prevDt := 0.002

	last := time.Now()

	const timeControlButtonWidth = 60.0
	const timeControlButtonY = 420

	gui := GUI{
		win:    win,
		atlas:  text.NewAtlas(basicfont.Face7x13, text.ASCII),
		canvas: pixelgl.NewCanvas(pixel.R(0, 0, guiCanvasWidth, win.Bounds().Max.Y)),
	}

	gui.CreateBatch("assets/sprites/spritesheet.png")

	spaceBetweenButtons := (guiCanvasWidth - (3 * timeControlButtonWidth)) / 4

	playButton := &Button{
		position: pixel.V(spaceBetweenButtons, 10),
		croppingArea: pixel.R(0, timeControlButtonY, 60,
			timeControlButtonY+timeControlButtonWidth),
		bounds:  pixel.R(0, 0, timeControlButtonWidth, timeControlButtonWidth),
		onClick: handlePlayClick,
	}

	gui.NewButton(playButton)

	pauseButton := &Button{
		position: pixel.V(spaceBetweenButtons*2+timeControlButtonWidth, 10),
		croppingArea: pixel.R(60, timeControlButtonY, 120,
			timeControlButtonY+timeControlButtonWidth),
		bounds:  pixel.R(0, 0, timeControlButtonWidth, timeControlButtonWidth),
		onClick: handlePauseClick,
	}

	gui.NewButton(pauseButton)

	stopButton := &Button{
		position: pixel.V(spaceBetweenButtons*3+timeControlButtonWidth*2, 10),
		croppingArea: pixel.R(120, timeControlButtonY, 180,
			timeControlButtonY+timeControlButtonWidth),
		bounds:  pixel.R(0, 0, timeControlButtonWidth, timeControlButtonWidth),
		onClick: handleStopClick,
	}

	gui.NewButton(stopButton)

	emitRateSlider := SliderWannabe{
		y:           360,
		canvasWidth: guiCanvasWidth,
		parameter:   &emitRate,
		format:      "%.0f par/sec\n",
	}

	gui.NewSliderWannabe(emitRateSlider)

	emitAngleSlider := SliderWannabe{
		y:           450,
		canvasWidth: guiCanvasWidth,
		parameter:   &emitAngle,
		format:      "%.0f degrees\n",
	}

	gui.NewSliderWannabe(emitAngleSlider)

	particleLifeSlider := SliderWannabe{
		y:           540,
		canvasWidth: guiCanvasWidth,
		parameter:   &particleLife,
		format:      "lives %.1f s\n",
	}

	gui.NewSliderWannabe(particleLifeSlider)

	initialVelocitySlider := SliderWannabe{
		y:           630,
		canvasWidth: guiCanvasWidth,
		parameter:   &initialVelocity,
		format:      "%.1f m/s",
	}

	gui.NewSliderWannabe(initialVelocitySlider)

	positionIntegratorSwitch := SwitchWannabe{
		y:           100,
		canvasWidth: guiCanvasWidth,
	}

	gui.NewSwitchWannabe(&positionIntegratorSwitch)
	positionIntegratorSwitch.handleExplicitEuler(nil)

	cam := pixel.IM.Scaled(camPos, 1.0).Moved(win.Bounds().Center().Sub(camPos))

	win.SetMatrix(cam)

	gui.SetMatrix(cam)
	gui.BindState(&state)
	gui.MainLoop()

	fps := time.Tick(time.Second / 200)

	gui.canvas.Clear(colornames.White)

	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	imd.Push(pixel.V(0, 0).Sub(win.Bounds().Center()).Add(circle.position))
	imd.Circle(circle.radius, 0)

	for !win.Closed() {
		win.Update()
		gui.Draw()

		if !gui.GetState().paused && !gui.GetState().stopped {
			dt := time.Since(last).Seconds()
			last = time.Now()
			timeElapsed += dt

			batch.Clear()

			updateParticles(
				particleSystem.particles,
				batch,
				dt,
				cam,
				positionIntegratorSwitch.positionIntegrator,
				circle,
			)

			win.Clear(colornames.Whitesmoke)

			batch.Draw(win)

			imd.Draw(win)

			gui.canvas.Draw(
				win,
				pixel.IM.Moved(pixel.V((win.Bounds().W()/-2.0)+(gui.canvas.Bounds().W()/2.0), 0.0)),
			)
			gui.batch.Draw(win)
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
					collide:      false,
				}
				particleSystem.particles = append(particleSystem.particles, particle)
				timeElapsed = timeElapsed - timeForOneParticle
			}
		} else if gui.GetState().paused && !gui.GetState().stopped {
			last = time.Now()
			gui.batch.Draw(gui.win)
		} else {
			last = time.Now()
			timeElapsed = 0
			particleSystem.particles = particleSystem.particles[:0]
			batch.Clear()
			win.Clear(colornames.Whitesmoke)
			gui.canvas.Draw(
				win,
				pixel.IM.Moved(pixel.V((win.Bounds().W()/-2.0)+(gui.canvas.Bounds().W()/2.0), 0.0)),
			)
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
