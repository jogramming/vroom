package vroom

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MOUSEBUTTONDOWN = int(sdl.MOUSEBUTTONDOWN)
	MOUSEBUTTONUP   = int(sdl.MOUSEBUTTONUP)
	MOUSEMOTION     = int(sdl.MOUSEMOTION)
	MOUSEWHEEL      = int(sdl.MOUSEWHEEL)

	KEYDOWN     = int(sdl.KEYDOWN)
	KEYUP       = int(sdl.KEYUP)
	TEXTEDITING = int(sdl.TEXTEDITING)
	TEXTINPUT   = int(sdl.TEXTINPUT)
)

// MOUSE BUTTONS

const (
	BUTTON_LEFT   = int(sdl.BUTTON_LEFT)
	BUTTON_MIDDLE = int(sdl.BUTTON_MIDDLE)
	BUTTON_RIGHT  = int(sdl.BUTTON_RIGHT)
	BUTTON_X1     = int(sdl.BUTTON_X1)
	BUTTON_X2     = int(sdl.BUTTON_X2)
)
