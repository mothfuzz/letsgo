package main

import (
	"embed"
	"fmt"
	"runtime"
	"strconv"

	//"math"
	_ "embed"
	"dyndraw/framework/actors"
	//"dyndraw/framework/events"
	"dyndraw/framework/render"

	gl "github.com/go-gl/gl/v3.1/gles2"
	. "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed resources
var Resources embed.FS

func main() {
	runtime.LockOSThread()

	render.Resources = Resources

	var width, height int32 = 800, 600
	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_ES)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 0)

	window, err = sdl.CreateWindow("owo",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height,
		sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	if err = gl.Init(); err != nil {
		panic(err)
	}
	sdl.GLSetSwapInterval(1)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.DepthFunc(gl.LESS)
	gl.Viewport(0, 0, int32(width), int32(height))

	proj := Perspective(DegToRad(60.0), float32(width)/float32(height), 0.1, 1000.0)
	view := Ident4()

	cameraPos := Vec3{0, 0, 2}

	fmt.Println("starting...")

	p := render.CreateProgram("basic.vert.glsl", "basic.frag.glsl")
	pmod := Ident4()
	p.BufferData("position", []float32{
		-1, +1, 0,
		+1, +1, 0,
		+1, -1, 0,
		+1, -1, 0,
		-1, -1, 0,
		-1, +1, 0,
	})
	p.BufferData("texcoord", []float32{
		0, 0,
		1, 0,
		1, 1,
		1, 1,
		0, 1,
		0, 0,
	})

	delta := float32(0)
	timer := sdl.GetTicks()
	ticks_passed := uint32(0)
	frames_passed := 0

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
					running = false
				}
			}
		}

		actors.Update()
		delta++

		now := sdl.GetTicks()
		ticks_passed += now - timer
		frames_passed += 1
		timer = now
		if ticks_passed > 1000 {
			window.SetTitle(strconv.Itoa(int(frames_passed)))
			ticks_passed = 0
			frames_passed = 0
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		view = LookAtV(cameraPos, cameraPos.Add(Vec3{0, 0, -1}), Vec3{0, 1, 0})
		pmod = pmod.Mul4(Rotate3DZ(0.025).Mat4())
		mv := view.Mul4(pmod)
		mvp := proj.Mul4(mv)

		p.Uniform("tex", render.Texture2d("gopog.png", false))
		p.Uniform("MVP", [16]float32(mvp))
		p.Draw()

		window.GLSwap()
	}

	actors.Quit()
}
