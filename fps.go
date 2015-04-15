package vroom

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type FPSCounter struct {
	BaseEntity
	Font        string
	FontOutline string
	Label       *Label
}

func (fps *FPSCounter) Init() {
	transform := &Transform{}
	fps.AddComponent(transform)

	updater := &UpdateComp{
		Update: fps.Update,
	}
	fps.AddComponent(updater)

	label := &Label{
		Font:         fps.Font,
		FontOutline:  fps.FontOutline,
		IgnoreCamera: true,
		Color:        sdl.Color{100, 100, 100, 255},
		ColorOutline: sdl.Color{0, 0, 0, 255},
		Text:         "FPS: 0",
	}
	fps.AddChild(label, true)
	fps.Label = label
}

func (fps *FPSCounter) Update(dt float64) {
	numFps := 1000 / dt
	fps.Label.SetText(fmt.Sprintf("FPS: %.2f", numFps))
}
