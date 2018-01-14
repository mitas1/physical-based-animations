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

func (sw *SwitchWannabe) handleExplicitEuler(state *HandledOptions) {
	sw.positionIntegrator = ExplicitEuler
	sw.explicitMidpointButton.isActive = false
	sw.verletButton.isActive = false
	sw.explicitEulerButton.isActive = true
}

func (sw *SwitchWannabe) handleMidpoint(state *HandledOptions) {
	sw.positionIntegrator = MidPoint
	sw.verletButton.isActive = false
	sw.explicitEulerButton.isActive = false
	sw.explicitMidpointButton.isActive = true
}

func (sw *SwitchWannabe) handleVerlet(state *HandledOptions) {
	sw.positionIntegrator = Verlet
	sw.explicitMidpointButton.isActive = false
	sw.explicitEulerButton.isActive = false
	sw.verletButton.isActive = true
}
