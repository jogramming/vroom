package vroom

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

// Simple sprite component for drawing sprites
type Sprite struct {
	BaseComponent
	Texture       *sdl.Texture
	IgnoreCamera  bool
	Width, Height int
}

// creates a new sprite with x, y, and width height from texture name
// if w and h is 0 it will take that from the texture
func (e *Engine) NewSprite(w, h int, ignoreCamera bool, texture string) *Sprite {
	tex := e.GetTexture(texture)
	if tex == nil {
		fmt.Println("Can't find texture: ", texture)
		return nil
	}

	_, _, rw, rh, _ := tex.Query()
	if w <= 0 {
		w = rw
	}
	if h <= 0 {
		h = rh
	}

	s := &Sprite{
		Texture:      tex,
		Width:        w,
		Height:       h,
		IgnoreCamera: ignoreCamera,
	}
	return s
}

// func (s *Sprite) CreatePhysBody(mass, moment float64, static bool) {
// 	physComp := &PhysBodyComp{}
// 	s.AddComponent(physComp)
// 	physComp.CreateBoxBody(s.Width, s.Height, vect.Float(mass), vect.Float(moment), static)
// }

func (s *Sprite) Init() {
	if s.GetComponent("Transform") == nil {
		transform := &Transform{}
		s.AddComponent(transform)
	}
}

func (s *Sprite) Draw(renderer *sdl.Renderer) {
	if s.Texture == nil {
		return
	}
	transform := s.GetComponent("Transform")
	if transform == nil {
		return
	}

	casted, ok := transform.(*Transform)
	if !ok {
		return
	}

	position := casted.CalcPos()
	if !s.IgnoreCamera {
		position.Sub(s.Parent.GetEngine().Camera)
	}

	angle := casted.CalcAngle()

	//center := &sdl.Point{X: int32(position.X + vect.Float(s.Width/2)), Y: int32(position.Y + vect.Float(s.Height/2))}
	center := &sdl.Point{X: int32(s.Width / 2), Y: int32(s.Height / 2)}
	//center := &sdl.Point{X: int32(s.Width), Y: int32(s.Height)}
	//center := &sdl.Point{}

	dstRect := &sdl.Rect{X: int32(position.X - float64(s.Width/2)), Y: int32(position.Y - float64(s.Height/2)), W: int32(s.Width), H: int32(s.Height)}
	renderer.CopyEx(s.Texture, nil, dstRect, float64(angle), center, sdl.FLIP_NONE)
}

func (s *Sprite) Name() string {
	return "Sprite"
}

func (s *Sprite) GetLayer() int {
	return 0
}

const (
	LOOPSTART   = iota + 1 // When cycle is finnished, starts from the beginning again
	LOOPREVERSE            // Reverses
	LOOPNONE               // Just stops..
	LOOPCB                 // Calls a callback
)

type AnimatedSprite struct {
	Sprite
	LoopMethod   int
	Frames       []string
	FrameTime    float64
	CurFrame     int
	CurFrameTime float64
	Dir          int
	Finnished    bool
	FinnishedCB  func()
}

func (a *AnimatedSprite) Init() {
	a.Sprite.Init()
	if a.Dir == 0 {
		a.Dir = 1
	}
	if a.LoopMethod == 0 {
		a.LoopMethod = LOOPREVERSE
	}
}

func (a *AnimatedSprite) Update(dt float64) {
	if !a.Finnished {
		a.CurFrameTime += dt
		if a.CurFrameTime >= a.FrameTime {
			a.CurFrameTime -= a.FrameTime
			a.CurFrame += a.Dir
			if a.CurFrame >= len(a.Frames) || a.CurFrame < 0 {
				switch a.LoopMethod {
				case LOOPSTART:
					a.CurFrame = 0
				case LOOPREVERSE:
					a.Dir = -a.Dir
					a.CurFrame += 2 * a.Dir
				case LOOPNONE:
					a.Finnished = true
					a.CurFrame = len(a.Frames) - 1
				case LOOPCB:
					a.Finnished = true
					a.CurFrame = len(a.Frames) - 1
					if a.FinnishedCB != nil {
						a.FinnishedCB()
					}
				}
			}
		}
		a.Texture = a.Parent.GetEngine().GetTexture(a.Frames[a.CurFrame])
	}
}

func (a *AnimatedSprite) Reset() {
	a.CurFrame = 0
	a.CurFrameTime = 0
}

func (e *Engine) NewAnimatedSprite(w, h int, ignoreCamera bool, textures []string, frameTime float64) *AnimatedSprite {
	baseSprite := e.NewSprite(w, h, ignoreCamera, textures[0])
	animated := AnimatedSprite{
		Sprite:    *baseSprite,
		Frames:    textures,
		FrameTime: frameTime,
	}

	return &animated
}
