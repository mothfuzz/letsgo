package render

import (
	"fmt"
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var textures = map[string]uint32{}

func Texture2D(name string, high_quality bool) uint32 {
	if tex, ok := textures[name]; ok {
		return tex
	} else {
		fmt.Println("loading new texture: " + name)
		gl.GenTextures(1, &tex)
		bytes, err := ReadResource("textures/" + name)
		if err != nil {
			panic(err)
		}
		rw, err := sdl.RWFromMem(bytes)
		if err != nil {
			panic(err)
		}
		image, err := img.LoadRW(rw, false)
		if err != nil {
			panic(err)
		}
		gl.BindTexture(gl.TEXTURE_2D, tex)
		var format int32 = gl.RGBA
		if image.BytesPerPixel() == 3 {
			format = gl.RGB
		}
		gl.TexImage2D(gl.TEXTURE_2D, 0, format, image.W, image.H, 0, uint32(format), gl.UNSIGNED_BYTE, image.Data())
		if high_quality {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_NEAREST)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		}
		gl.GenerateMipmap(gl.TEXTURE_2D)
		textures[name] = tex
		gl.BindTexture(gl.TEXTURE_2D, 0)
		return tex
	}
}
