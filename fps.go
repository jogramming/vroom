package vroom

import (
	"fmt"
)

type FPSCounter struct {
	BaseComponent
	Label *Label
}

func (fps *FPSCounter) Update(dt float64) {
	if fps.Label == nil {
		return
	}

	numFps := 1000 / dt
	fps.Label.SetText(fmt.Sprintf("FPS: %.2f", numFps))
}

func (fps *FPSCounter) Name() string {
	return "FPSCounter"
}
