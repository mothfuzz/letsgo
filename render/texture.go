package render

import (
	"fmt"
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/mothfuzz/letsgo/resources"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type texture2D struct {
	texId  uint32
	width  int32
	height int32
}

var textures = map[string]texture2D{}

func loadTexture2D(name string, high_quality bool) texture2D {
	if tex, ok := textures[name]; ok {
		return tex
	} else {
		fmt.Println("loading new texture: " + name)
		bytes, err := resources.ReadResource("textures/" + name)
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
		tex.width = image.W
		tex.height = image.H
		gl.GenTextures(1, &tex.texId)
		gl.BindTexture(gl.TEXTURE_2D, tex.texId)
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

func Texture2D(name string, high_quality bool) uint32 {
	return loadTexture2D(name, high_quality).texId
}
func TexWidth(name string) int {
	return int(loadTexture2D(name, true).width)
}
func TexHeight(name string) int {
	return int(loadTexture2D(name, true).height)
}
