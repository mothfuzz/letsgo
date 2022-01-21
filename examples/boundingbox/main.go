package main

import (
	//"fmt"
	"embed"
	"math"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/collision"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/mothfuzz/letsgo/transform"
	"github.com/mothfuzz/letsgo/vecmath"
)

func DrawLine(x1, y1, x2, y2 float32) {
	distance := vecmath.Vec2{x2, y2}.Sub(vecmath.Vec2{x1, y1}).Len()
	origin := vecmath.Vec2{x1, y1}.Add(vecmath.Vec2{x2, y2}).Mul(0.5)
	angle := float32(math.Atan2(float64(y2-y1), float64(x2-x1)))
	t := transform.Origin2D()
	t.SetPosition2D(origin.Elem())
	t.SetRotation2D(angle)
	t.SetScale2D(distance, 1)
	render.DrawSprite("dot.png", t.Mat4())
}

type BoundingActor struct {
	transform.Transform
	collision.Collider
}

func (ba *BoundingActor) Init() {
	ba.Collider = collision.NewBoundingBox(32, 32, 1)
	ba.Transform.Scale2D(2, 2)
	ba.Transform.Translate(0, 0, -1)
}
func (ba *BoundingActor) Update() {
	ba.Transform.Rotate2D(0.02)
}
func (ba *BoundingActor) Destroy() {}
func (ba *BoundingActor) Draw() {
	render.DrawSprite("square.png", ba.Transform.Mat4())
	c := collision.TransformCollider(ba.Collider, ba.Transform)
	DrawLine(c.Min.X(), c.Min.Y(), c.Max.X(), c.Min.Y())
	DrawLine(c.Max.X(), c.Min.Y(), c.Max.X(), c.Max.Y())
	DrawLine(c.Max.X(), c.Max.Y(), c.Min.X(), c.Max.Y())
	DrawLine(c.Min.X(), c.Max.Y(), c.Min.X(), c.Min.Y())
}

//go:embed resources
var Resources embed.FS

func main() {
	resources.Resources = Resources
	app.Init()
	defer app.Quit()
	app.SetWindowSize(320, 240)

	actors.SpawnAt(&BoundingActor{}, transform.Location2D(320/2, 240/2))

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
