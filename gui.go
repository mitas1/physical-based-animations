package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

// Button represents a gui button element
type Button struct {
	position           pixel.Vec
	bounds             pixel.Rect
	croppingArea       pixel.Rect
	croppingAreaActive pixel.Rect
	onClick            func(options *HandledOptions)
	sprite             *pixel.Sprite
	spriteActive       *pixel.Sprite
	isActive           bool
}

// Text represents parameters required to construct a struct object in gui
type Text struct {
	position pixel.Vec
	text     *Parameter
	widget   *text.Text
	format   string
}

// SliderWannabe represents abstract of slider that changes a given parameter
type SliderWannabe struct {
	y           float64
	canvasWidth float64 // width of the canvas SliderWannabe is rendered do so the internal objects can be properly spaced
	parameter   *Parameter
	format      string
}

// SwitchWannabe represents abstract of switch that can switch between various position integrator
// methods
type SwitchWannabe struct {
	y                  float64
	canvasWidth        float64 // width of the canvas SwitchWannabe is redered to so the internal objects can be properly spaced
	positionIntegrator PositionIntegrationMethod
	buttons            []*Button
}

// GUI represents an attributes of gui
type GUI struct {
	atlas       *text.Atlas
	win         *pixelgl.Window
	widgets     []*Button
	texts       []Text
	state       *HandledOptions
	matrix      pixel.Matrix
	batch       *pixel.Batch
	spritesheet pixel.Picture
	canvas      *pixelgl.Canvas
}

// HandledOptions define options which may be controlled by gui elements
type HandledOptions struct {
	paused  bool
	stopped bool
}

// CreateBatch creates a batch where gui be drawed to
func (gui *GUI) CreateBatch(imagePath string) {
	pic, err := loadPicture(imagePath)
	if err != nil {
		panic(err)
	}

	gui.spritesheet = pic
	gui.batch = pixel.NewBatch(&pixel.TrianglesData{}, pic)
}

// BindState binds a global state
func (gui *GUI) BindState(state *HandledOptions) {
	gui.state = state
}

// MainLoop starts a new goroutine where mouse events are handled
func (gui *GUI) MainLoop() {
	go func() {
		for !gui.win.Closed() {
			if gui.win.Pressed(pixelgl.MouseButtonLeft) {
				gui.handleClick(gui.win.MousePosition().X, gui.win.MousePosition().Y)
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
}

func (gui *GUI) handleClick(x, y float64) {
	y = winHeight - y
	for _, widget := range gui.widgets {
		if isInsideBoundingBox(x, y, widget.bounds, widget.position) {
			widget.onClick(gui.state)
		}
	}
}

// GetState returns a gui.state
func (gui *GUI) GetState() *HandledOptions {
	return gui.state
}

// SetMatrix sets a matrix of gui
func (gui *GUI) SetMatrix(matrix pixel.Matrix) {
	x, y := matrix.Project(pixel.ZV).XY()
	gui.matrix = pixel.IM.Moved(pixel.V(-x, y))
}

// NewButton creates a new button element
func (gui *GUI) NewButton(button *Button) {
	button.sprite = pixel.NewSprite(gui.spritesheet, button.croppingArea)
	button.spriteActive = pixel.NewSprite(gui.spritesheet, button.croppingAreaActive)

	gui.widgets = append(gui.widgets, button)
}

// NewText adds new text to the gui
func (gui *GUI) NewText(t Text) {
	gui.texts = append(gui.texts, t)
}

// NewSliderWannabe creates a slider which consists of two buttons and a text
func (gui *GUI) NewSliderWannabe(slider SliderWannabe) {
	// minusButton is placed 10 pixels from the left of the rendering canvas
	minusButton := Button{
		position:     pixel.V(10, slider.y),
		croppingArea: pixel.R(60, 360, 120, 420),
		bounds:       pixel.R(0, 0, 60, 60),
		onClick:      slider.parameter.handleMinus,
	}

	gui.NewButton(&minusButton)

	// plusButton is placed 10 pixels from the right of the rendering canvas
	// 60 pixels account for the width of the button itself
	plusButton := Button{
		position:     pixel.V(slider.canvasWidth-60-10, slider.y),
		croppingArea: pixel.R(0, 360, 60, 420),
		bounds:       pixel.R(0, 0, 60, 60),
		onClick:      slider.parameter.handlePlus,
	}

	gui.NewButton(&plusButton)

	txt := text.New(pixel.V(0, 0), gui.atlas)
	txt.Color = colornames.Black
	// textWidget is placed roughly to the midle of the two buttons
	// this is not the exact middle but looks fitting
	textWidget := Text{
		position: pixel.V((slider.canvasWidth-60)/2, slider.y+35),
		text:     slider.parameter,
		widget:   txt,
		format:   slider.format,
	}

	gui.NewText(textWidget)
}

// NewSwitchWannabe creates a switch that consists of three buttons
func (gui *GUI) NewSwitchWannabe(sw *SwitchWannabe) {
	sw.buttons = append(sw.buttons, &Button{
		position:           pixel.V((sw.canvasWidth-230)/2, sw.y),
		croppingArea:       pixel.R(0, 240, 230, 300),
		croppingAreaActive: pixel.R(0, 300, 230, 360),
		bounds:             pixel.R(0, 0, 230, 60),
		onClick:            sw.handleExplicitEuler,
	})

	sw.buttons = append(sw.buttons, &Button{
		position:           pixel.V((sw.canvasWidth-267)/2, sw.y+80),
		croppingArea:       pixel.R(0, 120, 267, 180),
		croppingAreaActive: pixel.R(0, 180, 267, 240),
		bounds:             pixel.R(0, 0, 267, 60),
		onClick:            sw.handleMidpoint,
	})

	sw.buttons = append(sw.buttons, &Button{
		position:           pixel.V((sw.canvasWidth-125)/2, sw.y+160),
		croppingArea:       pixel.R(0, 0, 125, 60),
		croppingAreaActive: pixel.R(0, 60, 125, 120),
		bounds:             pixel.R(0, 0, 125, 60),
		onClick:            sw.handleVerlet,
	})

	for _, button := range sw.buttons {
		gui.NewButton(button)
	}
}

// Draw draws a gui to gui batch
func (gui *GUI) Draw() {
	gui.batch.Clear()
	for _, widget := range gui.widgets {
		x0, y0 := widget.position.XY()
		x1, y1 := widget.bounds.Center().XY()
		if widget.isActive {
			widget.spriteActive.Draw(gui.batch, gui.matrix.Moved(pixel.V(x0, -y0).Sub(
				pixel.V(-x1, y1))))
		} else {
			widget.sprite.Draw(gui.batch, gui.matrix.Moved(pixel.V(x0, -y0).Sub(pixel.V(-x1, y1))))
		}
	}
}

// DrawText draws text widgets to gui text batch
func (gui *GUI) DrawText(target pixel.Target) {
	for _, t := range gui.texts {
		x0, y0 := t.position.XY()

		t.widget.Clear()
		t.widget.Dot = t.widget.Orig

		t.widget.WriteString(fmt.Sprintf(t.format, float64(t.text.value)))

		t.widget.Draw(gui.win, gui.matrix.Moved(pixel.V(x0, -y0)))
	}

}
