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

func (sw *SwitchWannabe) setActiveButton(index int) {
	for _, button := range sw.buttons {
		button.isActive = false
	}
	sw.buttons[index].isActive = true
}

func (sw *SwitchWannabe) handleExplicitEuler(state *HandledOptions) {
	sw.positionIntegrator = ExplicitEuler
	sw.setActiveButton(0)
}

func (sw *SwitchWannabe) handleMidpoint(state *HandledOptions) {
	sw.positionIntegrator = MidPoint
	sw.setActiveButton(1)
}

func (sw *SwitchWannabe) handleVerlet(state *HandledOptions) {
	sw.positionIntegrator = Verlet
	sw.setActiveButton(2)
}
