package vroom

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

// Label component
type Label struct {
	BaseComponent
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

// Helper function to create label struct
func NewLabel(text string, center bool, font, outline string) *Label {
	return &Label{
		Font:         font,
		FontOutline:  outline,
		CenterHor:    center,
		CenterVert:   center,
		Color:        sdl.Color{100, 100, 100, 255},
		ColorOutline: sdl.Color{0, 0, 0, 255},
		Text:         text,
	}
}

func NewSimpleLabelEntity(x, y float64, text string, center bool, font, outline string, ignoreCamera bool) (Entity, *Label) {
	label := NewLabel(text, center, font, outline)
	label.IgnoreCamera = ignoreCamera
	ent := NewEntity(x, y)
	ent.AddComponent(label)
	return ent, label
}

func (l *Label) Init() {
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
		position.Sub(l.Parent.GetEngine().Camera)
	}

	center := &sdl.Point{X: int32(l.Width / 2), Y: int32(l.Height / 2)}
	if l.CenterHor {
		position.X -= float64(l.Width / 2)
		center.X -= int32(l.Width / 2)
	}

	if l.CenterVert {
		position.Y -= float64(l.Height / 2)
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
		texture = l.Parent.GetEngine().CreateOutlinedTextTexture(l.Font, l.FontOutline, text, l.Color, l.ColorOutline)
	} else {
		texture = l.Parent.GetEngine().CreateTextTexture(l.Font, text, l.Color)
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

func (l *Label) Name() string {
	return "Label"
}

func (l *Label) GetLayer() int {
	return 1
}
