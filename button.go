package vroom

// Button Entity
type Button struct {
	BaseComponent

	Width, Height int

	ClickSound string
	HoverSound string

	HoverSprite *Sprite
	IdleSprite  *Sprite
	ClickSprite *Sprite

	IsHover     bool
	IsMouseDown bool

	OnClick func()
}

func (b *Button) Init() {
	b.ClickSprite.SetEnabled(false)
	b.HoverSprite.SetEnabled(false)
}

func (b *Button) MouseEnter() {
	b.IsHover = true
	b.ClickSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(false)
	b.HoverSprite.SetEnabled(true)
	if b.HoverSound != "" {
		b.Parent.GetEngine().PlaySound(b.HoverSound)
	}
}

func (b *Button) MouseLeave() {
	b.IsHover = false
	b.IsMouseDown = false
	b.ClickSprite.SetEnabled(false)
	b.HoverSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(true)
}

func (b *Button) MouseDown(x, y, button int) {
	b.IsMouseDown = true
	b.HoverSprite.SetEnabled(false)
	b.IdleSprite.SetEnabled(false)
	b.ClickSprite.SetEnabled(true)
}

func (b *Button) MouseUp(x, y, button int) {
	if b.IsMouseDown {
		if b.OnClick != nil {
			if b.ClickSound != "" {
				b.Parent.GetEngine().PlaySound(b.ClickSound)
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

func (b *Button) MouseMove(x, y int) {

}

func (b *Button) Name() string {
	return "Button"
}
