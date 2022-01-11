package collision

import (
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

//actually a triangle
type Plane struct {
	origin Vec3
	normal Vec3
	points [3]Vec3
}

func (p *Plane) Origin() Vec3 {
	return p.origin
}
func (p *Plane) Normal() Vec3 {
	return p.normal
}
func (p *Plane) Points() [3]Vec3 {
	return p.points
}

func triNorm(a, b, c Vec3) Vec3 {
	//(B - A) x (C - A)
	return b.Sub(a).Cross(c.Sub(a)).Normalize()
}
func NewPlaneAt(x, y, z float32, a, b, c Vec3) Plane {
	t := Vec3{x, y, z}
	n := triNorm(a, b, c)
	return Plane{t, n, [3]Vec3{a, b, c}}
}
func NewPlane(a, b, c Vec3) Plane {
	t := a.Add(b).Add(c).Mul(1.0 / 3.0) //centroid
	n := triNorm(a, b, c)
	return Plane{t, n, [3]Vec3{a, b, c}}
}
func NewPlane2D(a Vec2, b Vec2) Plane {
	l := a.Sub(b).Len()
	t := a.Add(b).Mul(0.5)
	c := t.Vec3(l)
	n := triNorm(a.Vec3(0), b.Vec3(0), c)
	return Plane{t.Vec3(0), n, [3]Vec3{a.Vec3(0), b.Vec3(0), c}}
}

func TransformPlane(p Plane, t transform.Transform) Plane {
	model := t.Mat4()
	origin := model.Mul4x1(p.origin.Vec4(1.0)).Vec3()
	normal := model.Mul4x1(p.normal.Vec4(0.0)).Vec3()
	a := model.Mul4x1(p.points[0].Vec4(1.0)).Vec3()
	b := model.Mul4x1(p.points[1].Vec4(1.0)).Vec3()
	c := model.Mul4x1(p.points[2].Vec4(1.0)).Vec3()
	return Plane{origin, normal, [3]Vec3{a, b, c}}
}
