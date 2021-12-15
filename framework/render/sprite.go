package render

import (
	//"fmt"
	. "github.com/go-gl/mathgl/mgl32"
)

//internal sprite drawcall format
type drawCall struct {
	model     Mat4
	texcoords []float32
}

type SpriteBatch struct {
	program   *Program
	drawCalls map[string][]drawCall
}

func (sb *SpriteBatch) DrawSprite(image string, model Mat4) {
	if sb.drawCalls == nil {
		sb.drawCalls = map[string][]drawCall{}
	}
	sb.drawCalls[image] = append(sb.drawCalls[image], drawCall{model, nil})
}
func (sb *SpriteBatch) DrawSpriteAnimated(image string, model Mat4, texcoords []float32) {
	sb.DrawSprite(image, model)
	sb.drawCalls[image][len(sb.drawCalls[image])-1].texcoords = texcoords
}
func (sb *SpriteBatch) Draw() {
	if sb.program == nil {
		sb.program = CreateProgram("basic.vert.glsl", "basic.frag.glsl")
		sb.program.BindBuffer("position", Quad.Position)
	}
	sb.program.LoadAttributes()
	for image, sprites := range defaultSpriteBatch.drawCalls {
		sb.program.Uniform("tex", Texture2D(image, false))
		for _, sprite := range sprites {
			if sprite.texcoords == nil {
				//use preexisting quad (fastest case)
				sb.program.BindBuffer("texcoord", Quad.TexCoord)
			} else {
				//if animated, upload animated texcoords
				sb.program.BufferData("texcoord", sprite.texcoords)
			}
			mv := ActiveCamera.GetView().Mul4(sprite.model)
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
func DrawSpriteAnimated(image string, model Mat4, texcoords []float32) {
	defaultSpriteBatch.DrawSpriteAnimated(image, model, texcoords)
}
func DrawSprites() {
	defaultSpriteBatch.Draw()
}

func AnimateTexCoords(framesW, framesH int, index int) []float32 {
	newCoords := make([]float32, len(Quad.TexCoord.Data))
	copy(newCoords, Quad.TexCoord.Data)
	for i := 0; i < len(newCoords)-1; i += 2 {
		newCoords[i] /= float32(framesW)
		newCoords[i+1] /= float32(framesH)
		newCoords[i] += float32(index%framesW) / float32(framesW)
		newCoords[i+1] += float32(index/framesW) / float32(framesH)
	}
	return newCoords
}
