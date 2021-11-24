package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	batch SpriteBatch
	id    uint64
	model Mat4
}

type SpriteBatch struct {
	*Program
	Sprites map[string]map[uint64]Sprite
}

func (s *SpriteBatch) Draw(view Mat4, proj Mat4) {
	if s.Program == nil {
		s.Program = CreateProgram("basic.vert.glsl", "basic.frag.glsl")
		s.Program.BufferData("position", []float32{
			-1, +1, 0,
			+1, +1, 0,
			+1, -1, 0,
			+1, -1, 0,
			-1, -1, 0,
			-1, +1, 0,
		})
		s.Program.BufferData("texcoord", []float32{
			0, 0,
			1, 0,
			1, 1,
			1, 1,
			0, 1,
			0, 0,
		})
	}
	s.Program.LoadAttributes()
	for texname, sprites := range s.Sprites {
		//sort by textures to reduce state changes
		s.Program.Uniform("tex", Texture2d(texname, false))
		for _, sprite := range sprites {
			mv := view.Mul4(sprite.model)
			mvp := proj.Mul4(mv)
			s.Program.Uniform("MVP", mvp)
			s.Program.DrawArrays()
		}
		s.Program.ClearTextureUnits()
	}
}
