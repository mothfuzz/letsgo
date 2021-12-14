package render

var (
	DefaultCamera = Camera{}
	ActiveCamera  = &DefaultCamera
	Quad          = struct {
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
			0, 1,
			1, 1,
			1, 0,
			1, 0,
			0, 0,
			0, 1,
		}},
	}
)
