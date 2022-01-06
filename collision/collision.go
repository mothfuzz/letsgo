package collision

import (
	"math"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
)

func insideTriangleVertices(p Vec3, r float32, a, b, c Vec3) bool {
	r2 := r * r
	if p.Sub(a).LenSqr() <= r2 {
		return true
	}
	if p.Sub(b).LenSqr() <= r2 {
		return true
	}
	if p.Sub(c).LenSqr() <= r2 {
		return true
	}
	return false
}
func sphereEdge(p Vec3, r float32, a, b Vec3) bool {
	r2 := r * r
	//check a
	/*if p.Sub(a).LenSqr() <= r2 {
		return true
	}
	//check b
	if p.Sub(b).LenSqr() <= r2 {
		return true
	}*/
	//check parametric distance
	ab := b.Sub(a)
	t := p.Sub(a).Dot(ab.Normalize())
	if t > 0 && t < 1 {
		x := a.Add(ab.Mul(t))
		if p.Sub(x).LenSqr() <= r2 {
			return true
		}
	}
	return false
}
func insideTriangleEdges(p Vec3, r float32, a, b, c Vec3) bool {
	if sphereEdge(p, r, a, b) {
		return true
	}
	if sphereEdge(p, r, b, c) {
		return true
	}
	if sphereEdge(p, r, c, a) {
		return true
	}
	return false
}
func pointInTriangle(p Vec3, a, b, c Vec3) bool {
	ab := b.Sub(a)
	ac := c.Sub(a)
	ap := p.Sub(a)
	abac := ab.Cross(ac)
	//check if 3 points are coplanar first
	//the floating point errors are strong with this one...
	if math.Abs(float64(ap.Dot(abac))) <= 1e-2 {
		//compute barycentric coords
		//u = ||CAxCP|| / ||ABxAC||
		//v = ||ABxAP|| / ||ABxAC||
		//w = ||BCxBP|| / ||ABxAC||
		abacl := abac.LenSqr()
		u := c.Sub(a).Cross(c.Sub(p)).LenSqr() / abacl
		v := a.Sub(b).Cross(a.Sub(p)).LenSqr() / abacl
		w := b.Sub(c).Cross(b.Sub(p)).LenSqr() / abacl
		if u >= 0 && u <= 1 && v >= 0 && v <= 1 && w >= 0 && w <= 1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//moves a bounding sphere against a series of walls
func MoveAgainstPlanes(t *transform.Transform, planes []Plane, radius float32, xspeed, yspeed, zspeed float32) (float32, float32, float32) {
	velocity := Vec3{xspeed, yspeed, zspeed}
	for _, p := range planes {
		pos := t.GetPositionV().Add(velocity)
		//get vector from point to plane
		dist := pos.Sub(p.origin)
		//project it onto normal (assumed to be normalized already)
		//this gives us a vector from the point perpendicular to the plane
		//the length of which is the shortest possible distance
		v := p.normal.Mul(dist.Dot(p.normal))
		if v.LenSqr() <= radius*radius {
			a := p.points[0]
			b := p.points[1]
			c := p.points[2]
			//find the nearest point on the plane along that vector
			//then check if the point is actually within the bounds of the triangle
			if pointInTriangle(pos.Add(v), a, b, c) ||
				insideTriangleVertices(pos, radius, a, b, c) ||
				insideTriangleEdges(pos, radius, a, b, c) {
				//if colliding with a wall, subtract velocity going in the wall's direction
				//to prevent movement
				adj := p.normal.Mul(velocity.Dot(p.normal)) //.Mul(2) //bouncy :3
				preserve := velocity.LenSqr()
				velocity = velocity.Sub(adj)
				mag := velocity.LenSqr()
				//attempt to preserve momentum against slopes etc.
				//not physically accurate but it's more fun.
				if mag > 0 {
					velocity = velocity.Mul(preserve / mag / 1.0)
				}
			}
		}
	}
	return velocity.Elem()
}

type RayHit struct {
	actors.Actor
	Plane
	I Vec3
}

func RayCast(pos Vec3, ray Vec3) []RayHit {
	hits := []RayHit{}
	actors.All(func(ac actors.Actor) {
		if c, ok := ac.(HasCollider); ok {
			if c.GetCollider().IgnoreRaycast {
				return
			}
			for _, p := range c.GetCollider().Planes {
				if t, ok := ac.(transform.HasTransform); ok {
					p = TransformPlane(p, *t.GetTransform())
				}

				rdot1 := ray.Dot(p.normal)
				rdot2 := p.origin.Sub(pos).Dot(p.normal)
				t := rdot2 / rdot1
				i := pos.Add(ray.Mul(t))
				if t >= 0 && //t <= 1 && //t > 1 if plane exceeds distance
					pointInTriangle(i, p.points[0], p.points[1], p.points[2]) {
					hits = append(hits, RayHit{ac, p, i})
				}
			}
		}
	})
	return hits
}
func RayCastLen(pos Vec3, ray Vec3, l float32) (RayHit, bool) {
	ll := l * l
	shortest := ll
	ok, hit := false, RayHit{}
	for _, p := range RayCast(pos, ray) {
		dist := p.I.Sub(pos).LenSqr()
		if dist <= shortest {
			shortest = dist
			ok = true
			hit = p
		}
	}
	return hit, ok
}

func DistanceSqr(a actors.Actor, b actors.Actor) float32 {
	if at, ok := a.(transform.HasTransform); ok {
		if bt, ok := b.(transform.HasTransform); ok {
			at = at.GetTransform()
			bt = bt.GetTransform()
			ab := at.GetTransform().GetPositionV().Sub(bt.GetTransform().GetPositionV())
			return ab.LenSqr()
		}
	}
	return float32(math.MaxFloat32)
}

func Distance(a actors.Actor, b actors.Actor) float32 {
	if at, ok := a.(transform.HasTransform); ok {
		if bt, ok := b.(transform.HasTransform); ok {
			at = at.GetTransform()
			bt = bt.GetTransform()
			ab := at.GetTransform().GetPositionV().Sub(bt.GetTransform().GetPositionV())
			return ab.Len()
		}
	}
	return float32(math.MaxFloat32)
}

func SphereOverlap(ca Vec3, ra float32, cb Vec3, rb float32) bool {
	if cb.Sub(ca).Len() <= ra+rb {
		return true
	} else {
		return false
	}
}
func BoxOverlap(aMin Vec3, aMax Vec3, bMin Vec3, bMax Vec3) bool {
	return aMin.X() <= bMax.X() &&
		aMax.X() >= bMin.X() &&
		aMin.Y() <= bMax.Y() &&
		aMax.Y() >= bMin.Y() &&
		aMin.Z() <= bMax.Z() &&
		aMax.Z() >= bMin.Z()
}
func SphereBoxOverlap(ca Vec3, ra float32, bMin Vec3, bMax Vec3) bool {
	x := max2(bMin.X(), min2(ca.X(), bMax.X()))
	y := max2(bMin.Y(), min2(ca.Y(), bMax.Y()))
	z := max2(bMin.Z(), min2(ca.Z(), bMax.Z()))
	p := Vec3{x, y, z}
	return p.Sub(ca).LenSqr() <= ra*ra
}

func ActorOverlap(a actors.Actor, b actors.Actor) bool {
	if ca, ok := a.(HasCollider); ok {
		if cb, ok := b.(HasCollider); ok {
			aa := *ca.GetCollider()
			bb := *cb.GetCollider()
			//this could be less expensive, i.e. just transform extents
			//since we're not checking polygons anyway
			if at, ok := a.(transform.HasTransform); ok {
				aa = TransformCollider(aa, *at.GetTransform())
			}
			if bt, ok := b.(transform.HasTransform); ok {
				bb = TransformCollider(bb, *bt.GetTransform())
			}
			//just check boxes for now
			//we're not doing polygon-over-polygon collisions
			if aa.Shape != BoundingSphere && bb.Shape != BoundingSphere {
				return BoxOverlap(aa.Extents.Min, aa.Extents.Max, bb.Extents.Min, bb.Extents.Max)
			}
			//otherwise one of them is a sphere
			ra := (aa.Extents.Max.X() - aa.Extents.Min.X()) / 2.0
			rb := (bb.Extents.Max.X() - bb.Extents.Min.X()) / 2.0
			ca := aa.Extents.Min.Add(aa.Extents.Max).Mul(0.5)
			cb := aa.Extents.Min.Add(aa.Extents.Max).Mul(0.5)
			if aa.Shape == BoundingSphere && bb.Shape == BoundingSphere {
				return SphereOverlap(ca, ra, cb, rb)
			}
			if aa.Shape == BoundingSphere && bb.Shape != BoundingSphere {
				return SphereBoxOverlap(ca, ra, bb.Extents.Min, bb.Extents.Max)
			}
			if bb.Shape == BoundingSphere && aa.Shape != BoundingSphere {
				return SphereBoxOverlap(cb, rb, aa.Extents.Min, aa.Extents.Max)
			}
		}
	}
	return false
}
