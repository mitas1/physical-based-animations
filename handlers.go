package main

var state = HandledOptions{
	paused:  false,
	stopped: false,
}

func handlePlayClick(state *HandledOptions) {
	state.paused = false
	state.stopped = false
}

func handlePauseClick(state *HandledOptions) {
	state.paused = true
}

func handleStopClick(state *HandledOptions) {
	state.stopped = true
}
