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

//testing out rendering pipeline within an actor context
type Gopher struct {
	render.Transform
	program	*render.Program
}

func (g *Gopher) Init() {
	g.program = render.CreateProgram("basic.vert.glsl", "basic.frag.glsl")
}
func (g *Gopher) Update() {
	g.Transform.Rotate(0, 0, 0.025)
}
func (g *Gopher) Draw() {
	//TODO: switching programs every time is slow.
	//cache draw calls & batch by program, then draw all at once during actors.Draw()
	//TODO: Draw() should not be part of actors
	g.program.BindBuffer("position", render.Quad.Position)
	g.program.BindBuffer("texcoord", render.Quad.TexCoord)
	g.program.Uniform("tex", render.Texture2D("gopog.png", false))
	g.program.Uniform("MVP", render.MvpFromTransform(g.Transform, render.View, render.Projection))
	g.program.Draw()
}
func (g *Gopher) Destroy() {}



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

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.DepthFunc(gl.LESS)
	gl.Viewport(0, 0, int32(width), int32(height))


	render.Projection = Perspective(DegToRad(60.0), float32(width)/float32(height), 0.1, 1000.0)
	render.View = Ident4()

	cameraPos := Vec3{0, 0, 2}

	fmt.Println("starting...")

	for i := 0; i < 3; i++ {
		g := &Gopher{Transform: render.Origin}
		g.Transform.Translate(float32(i), 0, float32(i) / 30.0)
		actors.Spawn(g)
	}

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

		render.View = LookAtV(cameraPos, cameraPos.Add(Vec3{0, 0, -1}), Vec3{0, 1, 0})

		actors.Draw()

		window.GLSwap()
	}

	actors.Quit()
}
