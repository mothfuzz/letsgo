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
	//I don't like this either
	v1 := Vec3{-w / 2, -h / 2, -d / 2}
	v2 := Vec3{+w / 2, -h / 2, -d / 2}
	v3 := Vec3{-w / 2, -h / 2, +d / 2}
	v4 := Vec3{+w / 2, -h / 2, +d / 2}
	v5 := Vec3{-w / 2, +h / 2, -d / 2}
	v6 := Vec3{+w / 2, +h / 2, -d / 2}
	v7 := Vec3{-w / 2, +h / 2, +d / 2}
	v8 := Vec3{+w / 2, +h / 2, +d / 2}
	c := NewCollisionMesh([]Plane{
		//12
		//top
		NewPlane(v1, v2, v4),
		NewPlane(v1, v4, v3),
		//bottom
		NewPlane(v5, v6, v8),
		NewPlane(v5, v8, v7),
		//front
		NewPlane(v3, v4, v8),
		NewPlane(v3, v8, v7),
		//back
		NewPlane(v1, v2, v6),
		NewPlane(v1, v6, v5),
		//left
		NewPlane(v1, v3, v7),
		NewPlane(v1, v7, v5),
		//right
		NewPlane(v4, v2, v6),
		NewPlane(v4, v6, v8),
	})
	c.Shape = BoundingBox
	return c
}
func NewBoundingSphere(r float32) Collider {
	c := NewBoundingBox(r*2, r*2, r*2)
	c.Shape = BoundingSphere
	return c
}
