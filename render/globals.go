package render

var (
	DefaultCamera = Camera{}
	ActiveCamera  = &DefaultCamera
	//it's useful to have a default Quad because lots of things are rectangles.
	Quad = struct {
		Position *Buffer
		TexCoord *Buffer
	}{
		&Buffer{Data: []float32{
			-0.5, +0.5, 0,
			+0.5, +0.5, 0,
			+0.5, -0.5, 0,
			+0.5, -0.5, 0,
			-0.5, -0.5, 0,
			-0.5, +0.5, 0,
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
