package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Button struct {
	position 	pixel.Vec
	bounds 		pixel.Rect
	image		string
	onClick 	func(options *HandledOptions)
	sprite		*pixel.Sprite
}

type GUI struct {
	win			*pixelgl.Window
	widgets 	[]Button
	state 		*HandledOptions
	matrix		pixel.Matrix
}

type HandledOptions struct {
	running bool
}

const (
	menuWidth  = 200
	rowHeight = 60
)

func (gui *GUI) BindState(state *HandledOptions) {
	gui.state = state
}

func (gui *GUI) MainLoop() {
	go func(){
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

		if x > x0 + x1 && x < x0 + x2 && y > y0 + y1 && y < y0 + y2 {
			widget.onClick(gui.state)
		}
	}
}

func (gui *GUI) GetState() (*HandledOptions) {
	return gui.state
}

func (gui *GUI) SetMatrix(matrix pixel.Matrix) {
	x, y := matrix.Project(pixel.ZV).XY()
	gui.matrix = pixel.IM.Moved(pixel.V(-x,y))
}

func (gui *GUI) NewButton(button Button) {	
	pic, err := loadPicture(button.image)
	if err != nil {
		panic(err)
	}

	button.sprite = pixel.NewSprite(pic, button.bounds)

	gui.widgets = append(gui.widgets, button)
}

func (gui *GUI) Draw() {
	for _, widget := range gui.widgets {
		x0, y0 := widget.position.XY()
		x1, y1 := widget.bounds.Center().XY()
		widget.sprite.Draw(gui.win, gui.matrix.Moved(pixel.V(x0,-y0).Sub(pixel.V(-x1, y1))))
	}		
}