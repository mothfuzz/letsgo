package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

var (
	View Mat4 = Ident4()
	Projection Mat4 = Ident4()
	Quad = struct {
		Position *Buffer
		TexCoord *Buffer
	}{
		&Buffer{Data: []float32{
					-1, +1, 0,
					+1, +1, 0,
					+1, -1, 0,
					+1, -1, 0,
					-1, -1, 0,
					-1, +1, 0,
		}},
		&Buffer{Data: []float32{
					0, 0,
					1, 0,
					1, 1,
					1, 1,
					0, 1,
					0, 0,
		}},
	}
	Origin Transform = Transform{Vec3{0, 0, 0}, Vec3{0, 0, 0}, Vec3{1, 1, 1}}
)
