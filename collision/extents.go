package collision

import (
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
