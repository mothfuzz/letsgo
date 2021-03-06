package main

import (
	"embed"
	"math"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/app"
	"github.com/mothfuzz/letsgo/collision"
	"github.com/mothfuzz/letsgo/input"
	"github.com/mothfuzz/letsgo/render"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

//go:embed resources
var Resources embed.FS

const w = float32(800)
const h = float32(600)

type RayTest struct{}

func (r *RayTest) Update() {}
func (r *RayTest) Draw() {
	mx, my := input.GetMousePosition()
	startPoint := Vec3{w / 2, h / 2, 0}
	endPoint := render.RelativeToCamera(mx, my)

	t := transform.Origin2D()
	t.SetScale2D(4, 4)
	t.SetPosition(w/2, h/2, -0.1)
	render.DrawSprite("pointg.png", t.Mat4())

	ray := endPoint.Sub(startPoint).Normalize()
	for _, p := range collision.RayCast(startPoint, ray) {
		t.SetPosition(p.Point.X(), p.Point.Y(), -0.1)
		render.DrawSprite("point.png", t.Mat4())
	}
	if hit, ok := collision.RayCastLen(startPoint, ray, w/2); ok {
		t := transform.Origin2D()
		t.SetScale2D(4, 4)
		t.SetPosition(hit.Point.X(), hit.Point.Y(), -0.2)
		render.DrawSprite("pointg.png", t.Mat4())
	}

	render.ActiveCamera.Look2D(startPoint.Vec2().Add(endPoint.Vec2()).Mul(0.5))
}

type BnW struct {
	transform.Transform
	collision.Collider
}

func (b *BnW) Init() {
	b.Collider = collision.NewBoundingBox(16, 16, 1) //has extra long edges for some reason
}
func (b *BnW) Update() {}
func (b *BnW) Draw() {
	render.DrawSprite("bnw.png", b.Transform.Mat4())
	for _, p := range b.Collider.Planes {
		p = collision.TransformPlane(p, b.Transform)
		t := transform.Origin2D()
		t.SetScale2D(2, 2)
		t.SetPosition(p.Origin().X(), p.Origin().Y(), p.Origin().Z()-0.1)
		render.DrawSprite("point.png", t.Mat4())
		for _, v := range p.Points() {
			t := transform.Origin2D()
			t.SetScale2D(2, 2)
			t.SetPosition(v.X(), v.Y(), v.Z()-0.1)
			render.DrawSprite("pointg.png", t.Mat4())
		}
	}
}

func main() {
	resources.Resources = Resources
	app.Init()
	defer app.Quit()

	app.SetWindowSize(int32(w), int32(h))

	actors.Spawn(&RayTest{})
	total := 12
	for i := 0; i < total; i++ {
		div := math.Pi * 2.0 / float64(total)
		t := transform.Location2D(
			w/2+float32(math.Sin(float64(i)*div))*w/4,
			h/2+float32(math.Cos(float64(i)*div))*h/4,
		)
		t.SetScale2D(2, 2)
		actors.SpawnAt(&BnW{}, t)
		t = transform.Location2D(
			w/2+float32(math.Sin(float64(i)*div))*w/2,
			h/2+float32(math.Cos(float64(i)*div))*h/2,
		)
		t.SetScale2D(4, 4)
		actors.SpawnAt(&BnW{}, t)
	}

	for app.PollEvents() {
		app.Update()
		app.Draw()
	}
}
