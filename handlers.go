package main

var state = HandledOptions{
	running: true,
}

func HandlePlayClick(state *HandledOptions) {
	state.running = true
}

func HandlePauseClick(state *HandledOptions) {
	state.running = false
}
