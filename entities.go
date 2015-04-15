package vroom

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/vova616/chipmunk/vect"
	//"github.com/veandco/go-sdl2/sdl_ttf"
)

// Basic ui entities/components (labels, buttons and whatnot)

type Sprite struct {
	BaseEntity
	Texture       *sdl.Texture
	IgnoreCamera  bool
	Width, Height int
}

// creates a new sprite with x, y, and width height from texture name
// if w and h is 0 it will take that from the texture
func (e *Engine) NewSprite(x, y, w, h int, texture string) *Sprite {
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
		Texture: tex,
		Width:   w,
		Height:  h,
	}
	transform := &Transform{
		Position: vect.Vect{vect.Float(x), vect.Float(y)},
	}
	s.AddComponent(transform)
	return s
}

func (s *Sprite) CreatePhysBody(mass, moment float64, static bool) {
	physComp := &PhysBodyComp{}
	s.AddComponent(physComp)
	physComp.CreateBoxBody(s.Width, s.Height, vect.Float(mass), vect.Float(moment), static)
}

func (s *Sprite) Init() {
	if s.GetComponent("Transform") == nil {
		transform := &Transform{}
		s.AddComponent(transform)
	}

	physComp := s.GetComponent("PhysBodyComp")
	if physComp != nil {
		body := physComp.(*PhysBodyComp).Body
		if body != nil {
			s.GetEngine().Space.AddBody(body)
		}
	}

	draw := &DrawComp{
		Draw:  s.Draw,
		Layer: -5,
	}
	s.AddComponent(draw)
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
		position.Add(s.GetEngine().Camera)
	}

	angle := casted.CalcAngle()

	//center := &sdl.Point{X: int32(position.X + vect.Float(s.Width/2)), Y: int32(position.Y + vect.Float(s.Height/2))}
	//center := &sdl.Point{X: int32(s.Width / 2), Y: int32(s.Height / 2)}
	//center := &sdl.Point{X: int32(s.Width), Y: int32(s.Height)}
	center := &sdl.Point{}

	dstRect := &sdl.Rect{X: int32(position.X), Y: int32(position.Y), W: int32(s.Width), H: int32(s.Height)}
	renderer.CopyEx(s.Texture, nil, dstRect, float64(angle), center, sdl.FLIP_NONE)
}

// Label entity
type Label struct {
	BaseEntity
	Font         string
	FontOutline  string
	CenterHor    bool
	CenterVert   bool
	IgnoreCamera bool
	Color        sdl.Color
	ColorOutline sdl.Color

	Text          string
	Width, Height int
	Texture       *sdl.Texture
}

func (l *Label) Init() {
	transform := &Transform{}
	l.AddComponent(transform)

	draw := &DrawComp{
		Draw: l.Draw,
	}
	l.AddComponent(draw)

	if l.Text != "" {
		l.SetText(l.Text)
	}
}

func (l *Label) Draw(renderer *sdl.Renderer) {
	if l.Texture == nil {
		return
	}
	transform := l.GetComponent("Transform")
	if transform == nil {
		return
	}

	casted, ok := transform.(*Transform)
	if !ok {
		return
	}

	position := casted.CalcPos()
	if !l.IgnoreCamera {
		position.Add(l.GetEngine().Camera)
	}

	center := &sdl.Point{X: int32(position.X + vect.Float(l.Width/2)), Y: int32(position.Y + vect.Float(l.Height/2))}
	if l.CenterHor {
		position.X -= vect.Float(l.Width / 2)
		center.X -= int32(l.Width / 2)
	}

	if l.CenterVert {
		position.Y -= vect.Float(l.Height / 2)
		center.Y -= int32(l.Height / 2)
	}
	angle := casted.CalcAngle()

	dstRect := &sdl.Rect{X: int32(position.X), Y: int32(position.Y), W: int32(l.Width), H: int32(l.Height)}
	renderer.CopyEx(l.Texture, nil, dstRect, float64(angle), center, sdl.FLIP_NONE)
}

func (l *Label) SetText(text string) {
	if l.Texture != nil {
		l.Texture.Destroy()
	}

	var texture *sdl.Texture

	if l.FontOutline != "" {
		texture = l.GetEngine().CreateOutlinedTextTexture(l.Font, l.FontOutline, text, l.Color, l.ColorOutline)
	} else {
		texture = l.GetEngine().CreateTextTexture(l.Font, text, l.Color)
	}

	if texture == nil {
		fmt.Println("Texture not found!!")
		return
	}

	l.Texture = texture
	l.Text = text

	_, _, w, h, _ := texture.Query()
	l.Width = w
	l.Height = h
}

// Button Entity
type UIButton struct {
	BaseEntity

	Hover         string
	Idle          string
	Click         string
	Width, Height int

	Font        string
	FontOutline string
	Text        string

	ClickSound string
	HoverSound string

	HoverSprite *Sprite
	IdleSprite  *Sprite
	ClickSprite *Sprite

	Label *Label

	IsHover     bool
	IsMouseDown bool

	OnClick func()
}

func (b *UIButton) Init() {
	mbox := &MouseBox{
		W: b.Width,
		H: b.Height,
	}

	b.AddComponent(mbox)

	transform := &Transform{}
	b.AddComponent(transform)

	mcComp := &MouseClickComp{
		MouseUp:   b.MouseUp,
		MouseDown: b.MouseDown,
	}
	b.AddComponent(mcComp)

	mhComp := &MouseHoverComp{
		MouseEnter: b.MouseEnter,
		MouseLeave: b.MouseLeave,
	}

	b.AddComponent(mhComp)

	// Initialize the sprites
	hoverSprite := b.GetEngine().GetSpriteFromTexture(b.Hover)
	idleSprite := b.GetEngine().GetSpriteFromTexture(b.Idle)
	clickSprite := b.GetEngine().GetSpriteFromTexture(b.Click)

	b.AddChild(hoverSprite, true)
	b.AddChild(idleSprite, true)
	b.AddChild(clickSprite, true)

	b.HoverSprite = hoverSprite
	b.ClickSprite = clickSprite
	b.IdleSprite = idleSprite

	hoverSprite.Width = b.Width
	idleSprite.Width = b.Width
	clickSprite.Width = b.Width
	hoverSprite.Height = b.Height
	idleSprite.Height = b.Height
	clickSprite.Height = b.Height

	hoverSprite.SetEnabled(false)
	clickSprite.SetEnabled(false)

	label := &Label{
		Font:         b.Font,
		FontOutline:  b.FontOutline,
		CenterHor:    true,
		CenterVert:   true,
		IgnoreCamera: true,
		Color:        sdl.Color{100, 100, 100, 255},
		ColorOutline: sdl.Color{0, 0, 0, 255},
		Text:         b.Text,
	}

	b.AddChild(label, true)
	b.Label = label

	label.GetComponent("Transform").(*Transform).Position = vect.Vect{vect.Float(b.Width / 2), vect.Float(b.Height/2) - 10}
}

func (b *UIButton) MouseEnter() {
	b.IsHover = true
	b.ClickSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(false)
	b.HoverSprite.SetEnabled(true)
	if b.HoverSound != "" {
		b.GetEngine().PlaySound(b.HoverSound)
	}
}

func (b *UIButton) MouseLeave() {
	b.IsHover = false
	b.IsMouseDown = false
	b.ClickSprite.SetEnabled(false)
	b.HoverSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(true)
}

func (b *UIButton) MouseDown(x, y, button int) {
	b.IsMouseDown = true
	b.HoverSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(false)
	b.ClickSprite.SetEnabled(true)
}

func (b *UIButton) MouseUp(x, y, button int) {
	if b.IsMouseDown {
		if b.OnClick != nil {
			if b.ClickSound != "" {
				b.GetEngine().PlaySound(b.ClickSound)
			}
			b.OnClick()
		}
	}
	b.IsMouseDown = false
	b.ClickSprite.SetEnabled(false)
	if b.IsHover {
		b.HoverSprite.SetEnabled(true)
		b.IdleSprite.SetEnabled(false)
	} else {
		b.HoverSprite.SetEnabled(false)
		b.IdleSprite.SetEnabled(true)
	}
}
