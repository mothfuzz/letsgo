package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

var keysCurrentFrame = map[sdl.Scancode]bool{}
var keysPreviousFrame = map[sdl.Scancode]bool{}

func UpdateKeys() {
	for sc, on := range keysCurrentFrame {
		keysPreviousFrame[sc] = on
	}
	sdl.PumpEvents()
	for sc, val := range sdl.GetKeyboardState() {
		keysCurrentFrame[sdl.Scancode(sc)] = (val != 0)
	}
}

func IsKeyDown(name string) bool {
	return keysCurrentFrame[sdl.GetScancodeFromName(name)]
}

func IsKeyPressed(name string) bool {
	return keysCurrentFrame[sdl.GetScancodeFromName(name)] &&
		!keysPreviousFrame[sdl.GetScancodeFromName(name)]
}
func IsKeyReleased(name string) bool {
	return !keysCurrentFrame[sdl.GetScancodeFromName(name)] &&
		keysPreviousFrame[sdl.GetScancodeFromName(name)]
}

func IsMouseButtonDown(button string) bool {
	_, _, mousestate := sdl.GetMouseState()
	switch button {
	case "left":
		return mousestate&sdl.ButtonLMask() != 0
	case "middle":
		return mousestate&sdl.ButtonMMask() != 0
	case "right":
		return mousestate&sdl.ButtonRMask() != 0
	case "x1":
		return mousestate&sdl.ButtonX1Mask() != 0
	case "x2":
		return mousestate&sdl.ButtonX2Mask() != 0
	default:
		return false
	}
}

func GetMousePosition() (int, int) {
	x, y, _ := sdl.GetMouseState()
	return int(x), int(y)
}

func GetMouseMovement() (int, int) {
	x, y, _ := sdl.GetRelativeMouseState()
	return int(x), int(y)
}
