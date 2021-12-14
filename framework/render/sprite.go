package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

type SpriteBatch struct {
	program   *Program
	drawCalls map[string][]Mat4
}

func (sb *SpriteBatch) DrawSprite(image string, model Mat4) {
	if sb.program == nil {
		sb.program = CreateProgram("basic.vert.glsl", "basic.frag.glsl")
		sb.program.BindBuffer("position", Quad.Position)
		sb.program.BindBuffer("texcoord", Quad.TexCoord)
	}
	if sb.drawCalls == nil {
		sb.drawCalls = map[string][]Mat4{}
	}
	sb.drawCalls[image] = append(sb.drawCalls[image], model)
}
func (sb *SpriteBatch) Draw() {
	sb.program.LoadAttributes()
	for image, models := range defaultSpriteBatch.drawCalls {
		sb.program.Uniform("tex", Texture2D(image, false))
		for _, model := range models {
			mv := ActiveCamera.GetView().Mul4(model)
			mvp := ActiveCamera.GetProjection().Mul4(mv)
			sb.program.Uniform("MVP", [16]float32(mvp))
			sb.program.DrawArrays()
		}
		sb.program.ClearTextureUnits()
	}
	for k := range sb.drawCalls {
		delete(sb.drawCalls, k)
	}
}

var defaultSpriteBatch = SpriteBatch{}

func DrawSprite(image string, model Mat4) {
	defaultSpriteBatch.DrawSprite(image, model)
}
func DrawSprites() {
	defaultSpriteBatch.Draw()
}
