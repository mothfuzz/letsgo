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

func (t *Transform) GetPositionV() Vec3 {
	return t.position
}
func (t *Transform) GetPosition() (float32, float32, float32) {
	return t.position.Elem()
}
func (t *Transform) X() float32 {
	return t.position.X()
}
func (t *Transform) Y() float32 {
	return t.position.Y()
}
func (t *Transform) Z() float32 {
	return t.position.Z()
}
func (t *Transform) GetRotationQ() Quat {
	return t.rotation
}
func (t *Transform) GetScaleV() Vec3 {
	return t.scale
}
func (t *Transform) GetScale() (float32, float32, float32) {
	return t.scale.Elem()
}

func (t *Transform) Translate(x, y, z float32) {
	t.position[0] += x
	t.position[1] += y
	t.position[2] += z
}
func (t *Transform) Translate2D(x, y float32) {
	t.position[0] += x
	t.position[1] += y
}
func (t *Transform) SetPosition(x, y, z float32) {
	t.position[0] = x
	t.position[1] = y
	t.position[2] = z
}
func (t *Transform) SetPosition2D(x, y float32) {
	t.position[0] = x
	t.position[1] = y
}
func (t *Transform) Rotate(x, y, z float32) {
	q := AnglesToQuat(x, y, z, XYZ)
	t.rotation = t.rotation.Mul(q)
}
func (t *Transform) SetRotation(x, y, z float32) {
	t.rotation = AnglesToQuat(x, y, z, XYZ)
}
func (t *Transform) Rotate2D(a float32) {
	q := AnglesToQuat(0, 0, a, XYZ)
	t.rotation = t.rotation.Mul(q)
}
func (t *Transform) SetRotation2D(a float32) {
	t.rotation = AnglesToQuat(0, 0, a, XYZ)
}
func (t *Transform) Scale(x, y, z float32) {
	t.scale[0] += x
	t.scale[1] += y
	t.scale[2] += z
}
func (t *Transform) Scale2D(x, y float32) {
	t.scale[0] += x
	t.scale[1] += y
}
func (t *Transform) SetScale(x, y, z float32) {
	t.scale[0] = x
	t.scale[1] = y
	t.scale[2] = z
}
func (t *Transform) SetScale2D(x, y float32) {
	t.scale[0] = x
	t.scale[1] = y
}

func (t *Transform) Mat4() Mat4 {
	m := Ident4()
	m = m.Mul4(Translate3D(t.position.X(), t.position.Y(), t.position.Z()))
	m = m.Mul4(t.rotation.Mat4())
	m = m.Mul4(Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z()))
	return m
}

func Origin() Transform {
	return Transform{Vec3{0, 0, 0}, QuatIdent(), Vec3{1, 1, 1}}
}

func Origin2D(w int, h int) Transform {
	return Transform{Vec3{0, 0, 0}, QuatIdent(), Vec3{float32(w), float32(h), 1}}
}

func MvpFromTransform(t Transform, v, p Mat4) [16]float32 {
	m := t.Mat4()
	mv := v.Mul4(m)
	mvp := p.Mul4(mv)
	return [16]float32(mvp)
}
