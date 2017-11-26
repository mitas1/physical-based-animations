package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Button represents a gui button element
type Button struct {
	position     pixel.Vec
	bounds       pixel.Rect
	croppingArea pixel.Rect
	onClick      func(options *HandledOptions)
	sprite       *pixel.Sprite
}

// GUI represents an attributes of gui
type GUI struct {
	win         *pixelgl.Window
	widgets     []Button
	state       *HandledOptions
	matrix      pixel.Matrix
	batch       *pixel.Batch
	spritesheet pixel.Picture
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
			if gui.win.JustPressed(pixelgl.MouseButtonLeft) {
				gui.handleClick(gui.win.MousePosition().X, gui.win.MousePosition().Y)
			}
		}
	}()
}

func (gui *GUI) handleClick(x, y float64) {
	y = winHeight - y
	for _, widget := range gui.widgets {
		x1, y1, x2, y2 := widget.bounds.Min.X, widget.bounds.Min.Y,
			widget.bounds.Max.X, widget.bounds.Max.Y

		x0, y0 := widget.position.XY()

		if x > x0+x1 && x < x0+x2 && y > y0+y1 && y < y0+y2 {
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
func (gui *GUI) NewButton(button Button) {
	button.sprite = pixel.NewSprite(gui.spritesheet, button.croppingArea)

	gui.widgets = append(gui.widgets, button)
}

// Draw draws a gui to gui batch
func (gui *GUI) Draw() {
	for _, widget := range gui.widgets {
		x0, y0 := widget.position.XY()
		x1, y1 := widget.bounds.Center().XY()
		widget.sprite.Draw(gui.batch, gui.matrix.Moved(pixel.V(x0, -y0).Sub(pixel.V(-x1, y1))))
	}
}
