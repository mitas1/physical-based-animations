package main

var state = HandledOptions{
	paused:  false,
	stopped: false,
}

func HandlePlayClick(state *HandledOptions) {
	state.paused = false
	state.stopped = false
}

func HandlePauseClick(state *HandledOptions) {
	state.paused = true
}

func HandleStopClick(state *HandledOptions) {
	state.stopped = true
}
