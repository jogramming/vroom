package vroom

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

func (e *Engine) LoadTexture(path string, name string) error {
	texture, err := img.LoadTexture(e.renderer, path)
	if err != nil {
		return err
	}
	if e.Textures == nil {
		e.Textures = make(map[string]*sdl.Texture)
	}
	e.Textures[name] = texture
	return nil
}

func (e *Engine) LoadSound(path, name string) error {
	if e.Sounds == nil {
		e.Sounds = make(map[string]*mix.Chunk)
	}
	chunk, err := mix.LoadWAV(path)
	if err != nil {
		return err
	}

	e.Sounds[name] = chunk
	return nil
}

func (e *Engine) PlaySound(name string) int {
	chunk := e.Sounds[name]
	if chunk == nil {
		return -1
	}

	chn, err := chunk.PlayChannel(-1, 0)
	if err != nil {
		fmt.Println("Error playing sound ", err)
	}
	return chn
}

func (e *Engine) LoadFont(path, name string, size int, outline int) error {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return err
	}
	if e.Fonts == nil {
		e.Fonts = make(map[string]*ttf.Font)
	}

	font.SetOutline(outline)

	e.Fonts[name] = font
	return nil
}

func (e *Engine) GetFont(name string) *ttf.Font {
	return e.Fonts[name]
}

func (e *Engine) CreateTextTexture(font string, text string, color sdl.Color) *sdl.Texture {
	f := e.GetFont(font)
	if f == nil {
		return nil
	}

	surface := f.RenderUTF8_Blended(text, color)
	texture, err := e.renderer.CreateTextureFromSurface(surface)
	surface.Free()
	if err != nil {
		return nil
	}
	return texture
}

func (e *Engine) CreateOutlinedTextTexture(font, outline string, text string, color sdl.Color, colorOutline sdl.Color) *sdl.Texture {
	f := e.GetFont(font)
	if f == nil {
		return nil
	}

	f2 := e.GetFont(outline)
	if f == nil {
		return nil
	}
	surface := f.RenderUTF8_Blended(text, color)
	surface2 := f2.RenderUTF8_Blended(text, colorOutline)

	surface.SetBlendMode(sdl.BLENDMODE_BLEND)
	surface.Blit(nil, surface2, &sdl.Rect{int32(f2.GetOutline()), int32(f2.GetOutline()), surface.W, surface.H})
	surface.Free()

	texture, err := e.renderer.CreateTextureFromSurface(surface2)
	surface2.Free()
	if err != nil {
		return nil
	}
	return texture
}

func (e *Engine) GetTexture(name string) *sdl.Texture {
	return e.Textures[name]
}

func (e *Engine) GetSpriteFromTexture(name string) *Sprite {
	texutre := e.GetTexture(name)
	if texutre == nil {
		return nil
	}
	_, _, w, h, _ := texutre.Query()

	sprite := &Sprite{
		Texture: texutre,
		Width:   w,
		Height:  h,
	}
	return sprite
}
