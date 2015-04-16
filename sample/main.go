package main

import (
	"fmt"
	"github.com/jonas747/vroom"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/vova616/chipmunk/vect"
)

var Engine *vroom.Engine

func main() {
	Engine = &vroom.Engine{}
	Engine.InitCoreSystems()

	fmt.Println("Initializing SDL")
	err := Engine.InitSDL(640, 400, "Sample game for vroom")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loading Assets")
	loadAssets()

	fmt.Println("Initializing scene")
	initScene()

	fmt.Println("Starting engine ")
	Engine.Start()
	Engine.Destroy()
}

type PathNamePair struct {
	name string
	path string
}

type FontLoadInfo struct {
	name    string
	path    string
	size    int
	outline int
}

func loadAssets() {
	textures := []PathNamePair{
		PathNamePair{"button_idle", "assets/button_idle.png"},
		PathNamePair{"button_hover", "assets/button_hover.png"},
		PathNamePair{"button_pressed", "assets/button_pressed.png"},
		PathNamePair{"box", "assets/box.png"},
	}

	for _, t := range textures {
		err := Engine.LoadTexture(t.path, t.name)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	fonts := []FontLoadInfo{
		FontLoadInfo{"mainfont", "assets/PatrickHand-Regular.ttf", 24, 0},
		FontLoadInfo{"mainfont_outline", "assets/PatrickHand-Regular.ttf", 24, 2},
	}

	for _, f := range fonts {
		err := Engine.LoadFont(f.path, f.name, f.size, f.outline)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	sounds := []PathNamePair{
		PathNamePair{"click", "assets/click.wav"},
		PathNamePair{"hover", "assets/hover.wav"},
	}

	for _, f := range sounds {
		err := Engine.LoadSound(f.path, f.name)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	mix.Volume(-1, sdl.MIX_MAXVOLUME/5)
}

func initScene() {
	button := &SimpleButton{
		Text: "BATTANN",
		X:    100,
		Y:    100,
		W:    150,
		H:    40,
	}
	Engine.AddEntity(button)
}

type SimpleButton struct {
	vroom.BaseEntity
	Text       string
	X, Y, W, H int
}

func (sb *SimpleButton) Init() {
	fmt.Println("Init Called")
	transform := vroom.NewTransform(vect.Float(sb.X), vect.Float(sb.Y), 0)
	sb.AddComponent(transform)

	mbox := &vroom.MouseBox{
		W: sb.W,
		H: sb.H,
	}
	sb.AddComponent(mbox)

	// Label has a different position so has to be in its own enity (but is a child of this)
	label := vroom.NewLabel(sb.Text, true, "mainfont", "mainfont_outline")
	lEntity := vroom.NewEntity(sb.W/2, sb.H/2)
	lEntity.AddComponent(label)

	sb.AddChild(lEntity, true)

	hs := Engine.NewSprite(sb.W, sb.H, true, "button_hover")
	is := Engine.NewSprite(sb.W, sb.H, true, "button_idle")
	cs := Engine.NewSprite(sb.W, sb.H, true, "button_pressed")

	sb.AddComponent(hs)
	sb.AddComponent(is)
	sb.AddComponent(cs)

	buttonComp := &vroom.Button{
		HoverSprite: hs,
		IdleSprite:  is,
		ClickSprite: cs,
		Width:       sb.W,
		Height:      sb.H,
		OnClick: func() {
			fmt.Println("Magic button pressed! Such press!")
		},
	}
	sb.AddComponent(buttonComp)
}
