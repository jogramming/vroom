package vroom

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

func (e *Engine) Loop() {
	e.running = true
	lastUpdate := time.Now()
	for e.running {
		now := time.Now()

		// Calculate deltatime
		deltatime := time.Since(lastUpdate)
		lastUpdate = time.Now()
		dt := float64(deltatime.Nanoseconds()) / float64(time.Second)

		// Clean up systems, maybe find a better way to do this later
		for _, v := range e.Systems {
			if time.Since(v.LastCleanUp()).Seconds() > 1 {
				v.CleanUp()
			}
		}

		e.ProcessEvents()
		e.StepPhysics(dt)
		e.Update(dt)
		e.Draw()

		elapsed := time.Since(now)
		milliseconds := elapsed.Seconds() * 100
		sleepAmount := (1000 / 60) - milliseconds
		if int(sleepAmount) > 0 {
			sdl.Delay(uint32(sleepAmount))
		}
	}
}

func (e *Engine) ProcessEvents() {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch evt := event.(type) {
		case *sdl.QuitEvent:
			e.running = false // byebye
		case *sdl.MouseMotionEvent:
			if e.window.GetID() != evt.WindowID {
				break
			}
			x := int(evt.X)
			y := int(evt.Y)

			e.MouseHoverSystem.MouseMove(x, y)
		case *sdl.MouseButtonEvent:
			if e.window.GetID() != evt.WindowID {
				break
			}
			x := int(evt.X)
			y := int(evt.Y)

			button := int(evt.Button)

			up := true
			if evt.Type == sdl.MOUSEBUTTONDOWN {
				up = false
			}
			e.MouseClickSystem.MouseButtonEvent(x, y, button, up)
		case *sdl.MouseWheelEvent:
			if e.window.GetID() != evt.WindowID {
				break
			}

		case *sdl.KeyUpEvent:
			if e.window.GetID() != evt.WindowID {
				break
			}
			e.Keyboardsystem.KeyboardEvent(evt.Keysym.Sym, true)
		case *sdl.KeyDownEvent:
			if e.window.GetID() != evt.WindowID {
				break
			}
			e.Keyboardsystem.KeyboardEvent(evt.Keysym.Sym, false)
		}
	}
}

func (e *Engine) StepPhysics(dt float64) {
	e.World.Step(dt)
}

func (e *Engine) Update(dt float64) {
	e.UpdateSystem.Update(dt)
}

func (e *Engine) Draw() {
	e.renderer.SetDrawColor(e.ClearColor.R, e.ClearColor.G, e.ClearColor.B, 255)
	e.renderer.Clear()
	e.DrawSystem.Draw(e.renderer)
	e.renderer.Present()
}
