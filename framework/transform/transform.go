package transform

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	position Vec3
	rotation Quat
	scale    Vec3
}

func (t *Transform) Translate(x, y, z float32) {
	t.position[0] += x
	t.position[1] += y
	t.position[2] += z
}
func (t *Transform) Rotate(x, y, z float32) {
	q := AnglesToQuat(x, y, z, XYZ)
	t.rotation = t.rotation.Mul(q)
}

func (t *Transform) Mat4() Mat4 {
	m := Ident4()
	m = m.Mul4(Translate3D(t.position.X(), t.position.Y(), t.position.Z()))
	m = m.Mul4(t.rotation.Mat4())
	m = m.Mul4(Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z()))
	return m
}

var Origin Transform = Transform{Vec3{0, 0, 0}, QuatIdent(), Vec3{1, 1, 1}}

func MvpFromTransform(t Transform, v, p Mat4) [16]float32 {
	m := t.Mat4()
	mv := v.Mul4(m)
	mvp := p.Mul4(mv)
	return [16]float32(mvp)
}
