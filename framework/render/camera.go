package render

import (
	"math"

	. "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position      Vec3
	view          Mat4
	projection    Mat4
	z2D           float32
	width, height int32
}

func (c *Camera) SetViewSize(width, height int32) {
	c.projection = Perspective(DegToRad(60.0), float32(width)/float32(height), 1.0, 1000.0)
	c.z2D = float32(math.Sqrt(math.Pow(float64(height), 2) - math.Pow(float64(height)/2.0, 2)))
	c.width = width
	c.height = height
}
func (c *Camera) Look(position Vec3, orientation Quat) {
	center := Vec3{0, 0, 1}
	center = orientation.Rotate(center)
	center = position.Add(center)
	c.position = position
	c.view = LookAtV(position, center, Vec3{0, -1, 0})
}
func (c *Camera) Look2D(position Vec2) {
	eye := Vec3{position.X(), position.Y(), -float32(c.z2D)}
	c.position = position.Vec3(-float32(c.z2D))
	c.view = LookAtV(eye, eye.Add(Vec3{0, 0, 1}), Vec3{0, -1, 0})
}
func (c *Camera) Init2D(width, height int32) {
	c.SetViewSize(width, height)
	c.Look2D(Vec2{float32(width / 2), float32(height / 2)})
}
func (c *Camera) GetZ2D() float32 {
	return c.z2D
}
func (c *Camera) GetPosition() Vec3 {
	return c.position
}
func (c *Camera) GetView() Mat4 {
	return c.view
}
func (c *Camera) GetProjection() Mat4 {
	return c.projection
}

func RelativeToCamera(x, y int) Vec3 {
	//0 is near plane, 1 is far plane, so have to calculate projected Z first
	winZ := Project(Vec3{0, 0, ActiveCamera.position.Z() + ActiveCamera.z2D}, ActiveCamera.view, ActiveCamera.projection, 0, 0, int(ActiveCamera.width), int(ActiveCamera.height)).Z()
	//then unproject to get world space X & Y
	v, err := UnProject(Vec3{float32(x), float32(ActiveCamera.height - int32(y)), winZ}, ActiveCamera.view, ActiveCamera.projection, 0, 0, int(ActiveCamera.width), int(ActiveCamera.height))
	if err != nil {
		//panics if view*proj is not invertible.......
		panic(err)
	}
	return v
}
