package render

import (
	//"fmt"
	"github.com/go-gl/mathgl/mgl32"
	. "github.com/go-gl/mathgl/mgl32"
)

//internal sprite drawcall format
type drawCall struct {
	model     Mat4
	texcoords *Buffer
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
func (sb *SpriteBatch) DrawSpriteAnimated(image string, model Mat4, texcoords *Buffer) {
	sb.DrawSprite(image, model)
	sb.drawCalls[image][len(sb.drawCalls[image])-1].texcoords = texcoords
}
func (sb *SpriteBatch) Draw() {
	if sb.program == nil {
		sb.program = CreateProgram("basic.vert.glsl", "basic.frag.glsl")
		sb.program.BindBuffer("position", Quad.Position)
	}
	for image, sprites := range sb.drawCalls {
		sb.program.Uniform("tex", Texture2D(image, false))
		w := float32(TexWidth(image))
		h := float32(TexHeight(image))
		for _, sprite := range sprites {
			if sprite.texcoords == nil {
				//if not animated, use preexisting quad
				sb.program.BindBuffer("texcoord", Quad.TexCoord)
			} else {
				//if animated, use dynamically generated buffer
				sb.program.BindBuffer("texcoord", sprite.texcoords)
			}
			sb.program.LoadAttributes()
			m := sprite.model.Mul4(mgl32.Scale3D(w, h, 1))
			mv := ActiveCamera.GetView().Mul4(m)
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
func DrawSpriteAnimated(image string, model Mat4, texcoords *Buffer) {
	defaultSpriteBatch.DrawSpriteAnimated(image, model, texcoords)
}
func DrawSprites() {
	defaultSpriteBatch.Draw()
}

type Frame struct {
	X float32
	Y float32
	W float32
	H float32
}

type SpriteAnimation struct {
	Frames  []Frame
	Tags    map[string][]int
	buffers []Buffer
}

func (s *SpriteAnimation) GetTexCoords(tag string, frame int) *Buffer {
	if s.buffers == nil {
		//generate proper texcoords for the quad, store in buffer for reuse
		s.buffers = make([]Buffer, len(s.Frames))
		for f := range s.buffers {
			b := &s.buffers[f]
			b.Data = make([]float32, len(Quad.TexCoord.Data))
			copy(b.Data, Quad.TexCoord.Data)
			for i := 0; i < len(b.Data)-1; i += 2 {
				if b.Data[i] == 0 {
					b.Data[i] = s.Frames[f].X
				} else {
					b.Data[i] = s.Frames[f].X + s.Frames[f].W
				}
				if b.Data[i+1] == 0 {
					b.Data[i+1] = s.Frames[f].Y
				} else {
					b.Data[i+1] = s.Frames[f].Y + s.Frames[f].H
				}
			}
		}
	}
	return &s.buffers[s.Tags[tag][frame]]
}
