package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Position Vec3
	Rotation Vec3
	Scale Vec3
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

func MvpFromTransform(t Transform, v Mat4, p Mat4) [16]float32 {
	m := Ident4()
	m = m.Mul4(Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z()))
	m = m.Mul4(Rotate3DX(t.Rotation.X()).Mat4())
	m = m.Mul4(Rotate3DY(t.Rotation.Y()).Mat4())
	m = m.Mul4(Rotate3DZ(t.Rotation.Z()).Mat4())
	m = m.Mul4(Scale3D(t.Scale.X(), t.Scale.Y(), t.Scale.Z()))
	mv := v.Mul4(m)
	mvp := p.Mul4(mv)
	return [16]float32(mvp)
}
