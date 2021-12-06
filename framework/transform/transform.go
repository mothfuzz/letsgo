package transform

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Position Vec3
	Rotation Vec3
	Scale    Vec3
}

func (t *Transform) Translate(x, y, z float32) {
	t.Position[0] += x
	t.Position[1] += y
	t.Position[2] += z
}
func (t *Transform) Rotate(x, y, z float32) {
	t.Rotation[0] += x
	t.Rotation[1] += y
	t.Rotation[2] += z
}

func (t *Transform) Mat4() Mat4 {
	m := Ident4()
	m = m.Mul4(Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z()))
	m = m.Mul4(Rotate3DX(t.Rotation.X()).Mat4())
	m = m.Mul4(Rotate3DY(t.Rotation.Y()).Mat4())
	m = m.Mul4(Rotate3DZ(t.Rotation.Z()).Mat4())
	m = m.Mul4(Scale3D(t.Scale.X(), t.Scale.Y(), t.Scale.Z()))
	return m
}

var Origin Transform = Transform{Vec3{0, 0, 0}, Vec3{0, 0, 0}, Vec3{1, 1, 1}}

func MvpFromTransform(t Transform, v, p Mat4) [16]float32 {
	m := t.Mat4()
	mv := v.Mul4(m)
	mvp := p.Mul4(mv)
	return [16]float32(mvp)
}
