package collision

import (
	"github.com/mothfuzz/letsgo/transform"
	. "github.com/mothfuzz/letsgo/vecmath"
	"math"
)

type Extents struct {
	Min Vec3
	Max Vec3
}

func min3(a, b, c float32) float32 {
	min := float32(math.MaxFloat32)
	if a < min {
		min = a
	}
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}
func max3(a, b, c float32) float32 {
	max := float32(-math.MaxFloat32)
	if a > max {
		max = a
	}
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}
func min2(a float32, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}
func max2(a float32, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}
func CalculateExtents(planes []Plane) Extents {
	pInf := float32(+math.MaxFloat32)
	nInf := float32(-math.MaxFloat32)
	minx, miny, minz := pInf, pInf, pInf
	maxx, maxy, maxz := nInf, nInf, nInf
	for _, p := range planes {
		a := p.points[0]
		b := p.points[1]
		c := p.points[2]
		minx = min2(minx, min3(a.X(), b.X(), c.X()))
		miny = min2(miny, min3(a.Y(), b.Y(), c.Y()))
		minz = min2(minz, min3(a.Z(), b.Z(), c.Z()))
		maxx = max2(maxx, max3(a.X(), b.X(), c.X()))
		maxy = max2(maxy, max3(a.Y(), b.Y(), c.Y()))
		maxz = max2(maxz, max3(a.Z(), b.Z(), c.Z()))
	}
	return Extents{Vec3{minx, miny, minz}, Vec3{maxx, maxy, maxz}}
}
func TransformExtents(e Extents, t transform.Transform) Extents {
	m4 := t.Mat4()
	points := [8]Vec3{
		{e.Min.X(), e.Min.Y(), e.Min.Z()},
		{e.Max.X(), e.Min.Y(), e.Min.Z()},
		{e.Min.X(), e.Max.Y(), e.Min.Z()},
		{e.Max.X(), e.Max.Y(), e.Min.Z()},
		{e.Min.X(), e.Min.Y(), e.Max.Z()},
		{e.Max.X(), e.Min.Y(), e.Max.Z()},
		{e.Min.X(), e.Max.Y(), e.Max.Z()},
		{e.Max.X(), e.Max.Y(), e.Max.Z()},
	}
	min := Vec3{1, 1, 1}.Mul(+1 * math.MaxFloat32)
	max := Vec3{1, 1, 1}.Mul(-1 * math.MaxFloat32)
	for i := 0; i < 8; i++ {
		//multiply by homogenous coords
		p := m4.Mul4x1(points[i].Vec4(1)).Vec3()
		if p.X() < min.X() {
			min[0] = p.X()
		}
		if p.Y() < min.Y() {
			min[1] = p.Y()
		}
		if p.Z() < min.Z() {
			min[2] = p.Z()
		}
		if p.X() > max.X() {
			max[0] = p.X()
		}
		if p.Y() > max.Y() {
			max[1] = p.Y()
		}
		if p.Z() > max.Z() {
			max[2] = p.Z()
		}
	}
	return Extents{min, max}
}
