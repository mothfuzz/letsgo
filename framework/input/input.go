package input

import (
	//"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func IsKeyDown(name string) bool {
	keystate := sdl.GetKeyboardState()
	return keystate[int(sdl.GetScancodeFromName(name))] != 0
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
