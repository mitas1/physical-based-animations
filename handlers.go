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

func (param *Parameter) handlePlus(state *HandledOptions) {
	if param.value+param.step <= param.max {
		param.value += param.step
	}
}

func (param *Parameter) handleMinus(state *HandledOptions) {
	if param.value-param.step >= param.min {
		param.value -= param.step
	}
}
