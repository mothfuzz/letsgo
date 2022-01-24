package main

import (
	"embed"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/input"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/mothfuzz/letsgo/transform"
)

//something we can see.
type Gopher struct {
	transform.Transform
	render.SpriteAnimation
	animationIndex int
}

func (g *Gopher) Init() {
	g.Transform.Translate(0, 0, -0.1)
	g.Transform.Scale2D(1.0/3.0, 1)
	g.SpriteAnimation = render.SpriteAnimation{
		Frames: []render.Frame{
			{X: 0.0 / 3.0, Y: 0, W: 1.0 / 3.0, H: 1},
			{X: 1.0 / 3.0, Y: 0, W: 1.0 / 3.0, H: 1},
			{X: 2.0 / 3.0, Y: 0, W: 1.0 / 3.0, H: 1},
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
		mx, my := input.GetMousePosition()
		x, y, z := render.RelativeToCamera(mx, my).Elem()
		t := transform.Origin2D()
		t.Translate(x, y, z)
		actors.Spawn(&Gopher{Transform: t})
	}
	if input.IsKeyDown("g") {
		g.animationIndex += 1
		g.animationIndex %= 3
	}
}
func (g *Gopher) Draw() {
	render.DrawSpriteAnimated("gopog.png", g.Transform.Mat4(), g.SpriteAnimation.GetTexCoords("idle", g.animationIndex))
}

type CameraController struct {
	transform.Transform
	render.Camera
}

func (c *CameraController) Init() {
	c.Camera.SetViewSize(800, 600)
	c.Transform.SetRotation(0, 0, 0)
	c.Transform.SetPosition(400, 300, -c.Camera.GetZ2D())
	render.ActiveCamera = &c.Camera
}
func (c *CameraController) Update() {
	motionVec := c.Transform.GetRotationQ().Rotate(mgl32.Vec3{0, 0, 1})
	if input.IsKeyDown("w") {
		c.Transform.Translate(motionVec.X(), 0, motionVec.Z())
	}
	if input.IsKeyDown("s") {
		c.Transform.Translate(-motionVec.X(), 0, -motionVec.Z())
	}
	if input.IsKeyDown("a") {
		c.Transform.Translate(-motionVec.Z(), 0, motionVec.X())
	}
	if input.IsKeyDown("d") {
		c.Transform.Translate(motionVec.Z(), 0, -motionVec.X())
	}
	dx, dy := input.GetMouseMovement()
	horiz := mgl32.QuatRotate(float32(dx)/800.0, mgl32.Vec3{0, 1, 0})
	vert := mgl32.QuatRotate(float32(dy)/600.0, mgl32.Vec3{-1, 0, 0})
	rotation := horiz.Mul(c.Transform.GetRotationQ()).Mul(vert)
	c.Transform.SetRotationQ(rotation)
	if input.IsKeyDown("up") {
		c.Transform.Rotate(0.02, 0, 0)
	}
	if input.IsKeyDown("down") {
		c.Transform.Rotate(-0.02, 0, 0)
	}
	c.Camera.Look(c.Transform.GetPositionV(), c.Transform.GetRotationQ())
}

//go:embed resources
var Resources embed.FS

func main() {
	resources.Resources = Resources
	app.Init()
	defer app.Quit()
	app.SetWindowSize(800, 600)
	//app.SetFullScreen(true)
	app.SetRelativeCursor(true)

	for i := 0; i < 5; i++ {
		g := &Gopher{Transform: transform.Origin2D()}
		g.Transform.Translate(float32(i)*64.0, 300, float32(i)*100.0)
		actors.Spawn(g)
	}

	actors.Spawn(&CameraController{})

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}

}
