package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	model Mat4
}

func (s *Sprite) Draw(m Mat4) {
	//this doesn't actually draw anything.
	s.model = m
}

type SpriteBatch struct {
	program *Program
	sprites map[string][]*Sprite
}

func (sb *SpriteBatch) CreateSprite(image string) *Sprite {
	if sb.program == nil {
		sb.program = CreateProgram("basic.vert.glsl", "basic.frag.glsl")
		sb.program.BindBuffer("position", Quad.Position)
		sb.program.BindBuffer("texcoord", Quad.TexCoord)
	}
	if sb.sprites == nil {
		sb.sprites = map[string][]*Sprite{}
	}
	sprite := new(Sprite)
	sprite.model = Ident4()
	sb.sprites[image] = append(sb.sprites[image], sprite)
	return sprite
}
func (sb *SpriteBatch) Draw() {
	sb.program.LoadAttributes()
	for image, spriteBatch := range defaultSpriteBatch.sprites {
		sb.program.Uniform("tex", Texture2D(image, false))
		for _, sprite := range spriteBatch {
			mv := View.Mul4(sprite.model)
			mvp := Projection.Mul4(mv)
			sb.program.Uniform("MVP", [16]float32(mvp))
			sb.program.DrawArrays()
		}
		sb.program.ClearTextureUnits()
	}
}

var defaultSpriteBatch = SpriteBatch{}

func CreateSprite(image string) *Sprite {
	return defaultSpriteBatch.CreateSprite(image)
}
func DrawSprites() {
	defaultSpriteBatch.Draw()
}
