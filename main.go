package main

import (
	"embed"

	"dyndraw/framework/actors"
	"dyndraw/framework/input"
	_ "embed"
	"math/rand"
	"time"

	"dyndraw/framework/app"
	//"dyndraw/framework/events"
	"dyndraw/framework/render"
	"dyndraw/framework/transform"
)

//testing out rendering pipeline within an actor context
type Gopher struct {
	transform.Transform
}

func (g *Gopher) Init() {}
func (g *Gopher) Update() {
	if input.IsKeyDown("r") {
		//g.Transform.Rotate(0, 0.025, 0)
		g.Transform.Rotate2D(0.025)
	}
	if input.IsKeyReleased("h") {
		x := rand.Float32() * 320.0
		y := rand.Float32() * 240.0
		t := transform.Origin2D(16, 16)
		t.Translate2D(x, y)
		actors.Spawn(&Gopher{t})
	}
	//x, y := input.GetMousePosition()
	//fmt.Println(x, y)
	//g.Transform.SetPosition2D(float32(x), float32(y))
}
func (g *Gopher) Draw() {
	render.DrawSprite("gopog.png", g.Transform.Mat4())
}
func (g *Gopher) Destroy() {}

type CameraController struct {
	transform.Transform
	render.Camera
}

func (c *CameraController) Init() {
	c.Camera.SetViewSize(320, 240)
	c.Transform.SetRotation(0, 0, 0)
	c.Transform.SetPosition(160, 120, -c.Camera.GetZ2D())
	render.ActiveCamera = &c.Camera
}
func (c *CameraController) Update() {
	if input.IsKeyDown("w") {
		c.Transform.Translate(0, 0, 1)
	}
	if input.IsKeyDown("s") {
		c.Transform.Translate(0, 0, -1)
	}
	if input.IsKeyDown("a") {
		c.Transform.Translate(-1, 0, 0)
	}
	if input.IsKeyDown("d") {
		c.Transform.Translate(1, 0, 0)
	}
	if input.IsKeyDown("up") {
		c.Transform.Rotate(0.02, 0, 0)
	}
	if input.IsKeyDown("down") {
		c.Transform.Rotate(-0.02, 0, 0)
	}
	c.Camera.Look(c.Transform.GetPositionV(), c.Transform.GetRotationQ())
}
func (c *CameraController) Destroy() {}

//go:embed resources
var Resources embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	render.Resources = Resources
	var width, height int32 = 320, 240
	app.Init()
	defer app.Quit()
	app.SetWindowSize(width, height)

	for i := 0; i < 8100; i++ {
		g := &Gopher{Transform: transform.Origin2D(32, 32)}
		g.Transform.Translate(float32(i)*64.0, 120, float32(i)*100.0)
		actors.Spawn(g)
	}

	actors.Spawn(&CameraController{})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}

}
