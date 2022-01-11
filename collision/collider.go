package collision

import (
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

type CollisionShape uint8

//in order of complexity...
const (
	CollisionMesh CollisionShape = iota
	BoundingSphere
	BoundingBox
)

//can do narrow phase and broad phase collision in one convenient struct
type Collider struct {
	Planes        []Plane
	Shape         CollisionShape //defaults to mesh test but overridable for simpler collisions
	Extents                      //can be used for bounding box test & radius test
	IgnoreRaycast bool           //it's expensive.
}

func (c *Collider) GetCollider() *Collider {
	return c
}

type HasCollider interface {
	GetCollider() *Collider
}

func TransformCollider(c Collider, t transform.Transform) Collider {
	c2 := Collider{Shape: c.Shape, IgnoreRaycast: c.IgnoreRaycast}
	c2.Planes = make([]Plane, len(c.Planes))
	copy(c2.Planes, c.Planes)
	for i := range c.Planes {
		c2.Planes[i] = TransformPlane(c.Planes[i], t)
	}
	c2.Extents = CalculateExtents(c2.Planes)
	return c2
}

func NewCollisionMesh(planes []Plane) Collider {
	return Collider{planes, CollisionMesh, CalculateExtents(planes), false}
}
func NewBoundingBox(w, h, d float32) Collider {
	v := []Vec3{
		// front
		{-w / 2, -h / 2, +d / 2},
		{+w / 2, -h / 2, +d / 2},
		{+w / 2, +h / 2, +d / 2},
		{-w / 2, +h / 2, +d / 2},
		// back
		{-w / 2, -h / 2, -d / 2},
		{+w / 2, -h / 2, -d / 2},
		{+w / 2, +h / 2, -d / 2},
		{-w / 2, +h / 2, -d / 2},
	}
	c := NewCollisionMesh([]Plane{
		// front
		NewPlane(v[0], v[1], v[2]),
		NewPlane(v[2], v[3], v[0]),
		// right
		NewPlane(v[1], v[5], v[6]),
		NewPlane(v[6], v[2], v[1]),
		// back
		NewPlane(v[7], v[6], v[5]),
		NewPlane(v[5], v[4], v[7]),
		// left
		NewPlane(v[4], v[0], v[3]),
		NewPlane(v[3], v[7], v[4]),
		// bottom
		NewPlane(v[4], v[5], v[1]),
		NewPlane(v[1], v[0], v[4]),
		// top
		NewPlane(v[3], v[2], v[6]),
		NewPlane(v[6], v[7], v[3]),
	})
	c.Shape = BoundingBox
	return c
}
func NewBoundingSphere(r float32) Collider {
	c := NewBoundingBox(r*2, r*2, r*2)
	c.Shape = BoundingSphere
	return c
}
