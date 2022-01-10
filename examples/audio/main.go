package main

import (
	"embed"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/mothfuzz/letsgo/sound"
)

//go:embed resources
var Resources embed.FS

func main() {

	resources.Resources = Resources
	app.Init()
	defer app.Quit()
	app.SetWindowSize(320, 240)

	sound.PlayMusic("eh.mp3")
	sound.PlaySound("bonk.mp3")

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
