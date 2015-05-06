package vroom

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
)

type Engine struct {

	// Core
	running      bool
	CurrentScene Scene
	Systems      []System
	Camera       vect.Vect

	// SDL
	window   *sdl.Window
	renderer *sdl.Renderer

	// Built in core systems
	DrawSystem       *DrawSystem
	UpdateSystem     *UpdateSystem
	MouseClickSystem *MouseClickSystem
	MouseHoverSystem *MouseHoverSystem
	Keyboardsystem   *KeyboardSystem

	// Assets
	Textures map[string]*sdl.Texture
	Fonts    map[string]*ttf.Font
	Sounds   map[string]*mix.Chunk

	//Physics
	Space *chipmunk.Space

	// Misc
	ClearColor sdl.Color
}

func (e *Engine) InitCoreSystems() {
	e.DrawSystem = &DrawSystem{}
	e.UpdateSystem = &UpdateSystem{}
	e.MouseClickSystem = &MouseClickSystem{}
	e.MouseHoverSystem = &MouseHoverSystem{}
	e.Keyboardsystem = &KeyboardSystem{}

	e.AddSystem(e.DrawSystem)
	e.AddSystem(e.UpdateSystem)
	e.AddSystem(e.MouseClickSystem)
	e.AddSystem(e.MouseHoverSystem)
	e.AddSystem(e.Keyboardsystem)

	e.Space = chipmunk.NewSpace()
	e.Space.Gravity = vect.Vect{0, 20}
}

func (e *Engine) InitSDL(w, h int, title string) error {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	e.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	e.renderer = renderer

	// Init sdl_image
	flags := img.Init(img.INIT_PNG)
	if flags&img.INIT_PNG != img.INIT_PNG {
		return img.GetError()
	}

	// init sdl ttf
	ttf.Init()

	// Init sound system
	err = mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 1024)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) AddSystem(sys System) {
	e.Systems = append(e.Systems, sys)
}

func (e *Engine) Start() {
	if e.running {
		return // Can't run the engine twice silly you
	}

	e.Loop()
}

func (e *Engine) Stop() {
	e.running = false
}

// Add this entity and all its children
func (e *Engine) AddEntity(entity Entity) {
	entity.SetEngine(e)
	if !entity.InitCalled() {
		entity.Init()
		entity.SetInitCalled()
	}

	children := entity.GetChildren(true)

	if len(children) > 0 {
		for _, child := range children {
			if !child.Added() {
				e.AddEntity(child)
			}
		}
	}

	entity.SetAdded(true)

	// Initialize all the components
	// And add them to system
	for _, compSlice := range entity.GetComponents() {
		for _, system := range e.Systems {
			for _, component := range compSlice {
				if !component.InitCalled() {
					component.SetInitCalled()
					component.Init()
				}
				system.AddComponent(component)
			}
		}
	}

	entity.Start()
}

// Removes the entity and all its children from the current scene
func (e *Engine) RemoveEntity(entity Entity) {
	// Call recusrively on children
	children := entity.GetChildren(false)
	if len(children) > 0 {
		for _, ent := range children {
			e.RemoveEntity(ent)
		}
	}
	entity.SetAdded(false)

	for _, compSlice := range entity.GetComponents() {
		for _, component := range compSlice {
			for _, system := range e.Systems {
				system.RemoveComponent(component)
			}
		}
	}
}

// Removes and destroys an entity and all its children
func (e *Engine) DestroyEntity(entity Entity) {
	e.RemoveEntity(entity)
	entity.Destroy()
}

// Remove all entities
func (e *Engine) Clear() {
	for _, v := range e.Systems {

	}
}

func (e *Engine) LoadScene(scene *Scene) {
	// Clear all the component buffer in the systems and load them up with components from this scene yo
}

func (e *Engine) ApplyCamera(x, y int) (int, int) {
	xo := x - int(e.Camera.X)
	yo := y - int(e.Camera.Y)
	return xo, yo
}

func (e *Engine) Destroy() {
	e.renderer.Destroy()
	e.window.Destroy()
	img.Quit()
	ttf.Quit()
	mix.CloseAudio()

	sdl.Quit()
}

func (e *Engine) GetKey(key sdl.Keycode) bool {
	return e.Keyboardsystem.Keys[key]
}

func RadiansToDeDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
