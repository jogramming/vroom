package main

import (
	"fmt"
	"github.com/jonas747/vroom"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/vova616/chipmunk/vect"
	"math/rand"
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
	// button := &vroom.UIButton{
	// 	Hover: "button_hover",
	// 	Idle:  "button_idle",
	// 	Click: "button_pressed",

	// 	Width:  200,
	// 	Height: 50,

	// 	Font:        "mainfont",
	// 	FontOutline: "mainfont_outline",
	// 	Text:        "This is a button",
	// 	ClickSound:  "click",
	// 	HoverSound:  "hover",

	// 	OnClick: func() {
	// 		fmt.Println("Magic button pressed! Such magic!")
	// 	},
	// }

	// Engine.AddEntity(button)
	// button.GetComponent("Transform").(*vroom.Transform).Position = vect.Vect{200, 100}

	ground := Engine.NewSprite(75, 350, 500, 0, "box")
	ground.CreatePhysBody(0, 0, true)
	Engine.AddEntity(ground)

	wl := Engine.NewSprite(75, 50, 0, 330, "box")
	wl.CreatePhysBody(0, 0, true)
	Engine.AddEntity(wl)

	wr := Engine.NewSprite(550, 50, 0, 330, "box")
	wr.CreatePhysBody(0, 0, true)
	Engine.AddEntity(wr)

	staticBox := Engine.NewSprite(250, 250, 40, 40, "box")
	//staticBox.GetComponent("Transform").(*vroom.Transform).Angle = 100
	staticBox.CreatePhysBody(0, 0, true)
	Engine.AddEntity(staticBox)

	fallingBox := Engine.NewSprite(207, 100, 50, 50, "box")
	fallingBox.CreatePhysBody(10, 10, false)
	Engine.AddEntity(fallingBox)

	b2 := Engine.NewSprite(250, 250, 20, 20, "box")
	Engine.AddEntity(b2)

	watcher := &SimpleEntity{
		box:  fallingBox,
		box2: b2,
	}
	Engine.AddEntity(watcher)

	// FPS Counter
	fpsCounter := &vroom.FPSCounter{
		Font:        "mainfont",
		FontOutline: "mainfont_outline",
	}

	Engine.AddEntity(fpsCounter)
}

type SimpleEntity struct {
	vroom.BaseEntity
	box  vroom.Entity
	box2 vroom.Entity
}

func (s *SimpleEntity) Init() {
	s.AddComponent(&vroom.UpdateComp{
		Update: func(dt float64) {
			if s.box != nil {
				transform := s.box.GetComponent("Transform")
				if transform != nil {
					pos := transform.(*vroom.Transform).CalcAngle()
					fmt.Println(pos)
				}
			}
			s.box2.GetComponent("Transform").(*vroom.Transform).Angle += 1
		},
	})
}

func addButtonssss(amount int) {
	for i := 0; i < amount; i++ {
		ic := i
		button := &vroom.UIButton{
			Hover: "button_hover",
			Idle:  "button_idle",
			Click: "button_pressed",

			Width:  200,
			Height: 50,

			Font:        "mainfont",
			FontOutline: "mainfont_outline",
			Text:        fmt.Sprintf("Button #%d", ic),

			OnClick: func() {
				fmt.Println("Magic button pressed! Such magic!", ic)
			},
		}

		Engine.AddEntity(button)

		button.GetComponent("Transform").(*vroom.Transform).Position = vect.Vect{vect.Float(rand.Intn(600)), vect.Float(rand.Intn(350))}
	}
}
