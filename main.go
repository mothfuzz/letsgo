package main

import (
	"embed"
	"math/rand"
	"time"

	"github.com/mothfuzz/dyndraw/framework/actors"
	"github.com/mothfuzz/dyndraw/framework/input"

	"github.com/mothfuzz/dyndraw/framework/app"
	//"github.com/mothfuzz/dyndraw/framework/events"
	"github.com/mothfuzz/dyndraw/framework/render"
	"github.com/mothfuzz/dyndraw/framework/transform"
)

//testing out rendering pipeline within an actor context
type Gopher struct {
	transform.Transform
	render.SpriteAnimation
	animationIndex int
}

func (g *Gopher) Init() {
	g.Transform.Translate(0, 0, -0.1)
	g.SpriteAnimation = render.SpriteAnimation{
		Frames: [][]float32{
			{0.0 / 3.0, 0, 1.0 / 3.0, 1},
			{1.0 / 3.0, 0, 1.0 / 3.0, 1},
			{2.0 / 3.0, 0, 1.0 / 3.0, 1},
		},
		Tags: map[string][]int{
			"idle": {0, 1, 2},
		},
	}
}
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
		actors.Spawn(&Gopher{Transform: t})
	}
	if input.IsKeyDown("g") {
		g.animationIndex += 1
		g.animationIndex %= 3
	}
}
func (g *Gopher) Draw() {
	render.DrawSpriteAnimated("gopog.png", g.Transform.Mat4(), g.SpriteAnimation.GetTexCoords("idle", g.animationIndex))
	//render.DrawSprite("gopog.png", g.Transform.Mat4())
}
func (g *Gopher) Destroy() {}

type Friendly struct {
	transform.Transform
}

func (f *Friendly) Init()    {}
func (f *Friendly) Update()  {}
func (f *Friendly) Destroy() {}
func (f *Friendly) Draw() {
	render.DrawSprite("friendly.png", f.Transform.Mat4())
}

type BnW struct {
	transform.Transform
}

func (b *BnW) Init()    {}
func (b *BnW) Update()  {}
func (b *BnW) Destroy() {}
func (b *BnW) Draw() {
	render.DrawSprite("circle.png", b.Transform.Mat4())
}

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
	app.Init()
	defer app.Quit()
	app.SetWindowSize(320, 240)

	for i := 0; i < 1; i++ {
		g := &Gopher{Transform: transform.Origin2D(32, 32)}
		g.Transform.Translate(float32(i)*64.0, 120, float32(i)*100.0)
		actors.Spawn(g)
	}
	for i := 0; i < 1; i++ {
		g := &Friendly{Transform: transform.Origin2D(32, 32)}
		g.Transform.Translate(float32(i)*100.0, 120, float32(i)*64.0)
		actors.Spawn(g)
	}
	for i := 1; i < 2; i++ {
		g := &BnW{Transform: transform.Origin2D(32, 32)}
		g.Transform.Translate(float32(i)*64.0, 120, float32(i)*-100.0)
		actors.Spawn(g)
	}

	actors.Spawn(&CameraController{})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}

}
