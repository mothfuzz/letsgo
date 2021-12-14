package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	view          Mat4
	projection    Mat4
	z2D           float32
	width, height int32
}

func (c *Camera) SetViewSize(width, height int32) {
	c.projection = Perspective(DegToRad(60.0), float32(width)/float32(height), 0.1, 1000.0)
	c.z2D = float32(math.Sqrt(math.Pow(float64(height), 2) - math.Pow(float64(height)/2.0, 2)))
	c.width = width
	c.height = height
}
func (c *Camera) Look(position Vec3, orientation Quat) {
	center := Vec3{0, 0, 1}
	center = orientation.Rotate(center)
	center = position.Add(center)
	c.view = LookAtV(position, center, Vec3{0, -1, 0})
}
func (c *Camera) Look2D(position Vec2) {
	eye := Vec3{float32(c.width) / 2, float32(c.height) / 2, -float32(c.z2D)}
	eye.Sub(position.Vec3(0))
	c.view = LookAtV(eye, eye.Add(Vec3{0, 0, 1}), Vec3{0, -1, 0})
}
func (c *Camera) Init2D(width, height int32) {
	c.SetViewSize(width, height)
	c.Look2D(Vec2{0, 0})
}
func (c *Camera) GetZ2D() float32 {
	return c.z2D
}
func (c *Camera) GetView() Mat4 {
	return c.view
}
func (c *Camera) GetProjection() Mat4 {
	return c.projection
}