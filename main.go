package main

import (
	"embed"
	"fmt"
	"runtime"
	"strconv"

	//"math"
	"dyndraw/framework/actors"
	_ "embed"
	//"dyndraw/framework/events"
	"dyndraw/framework/render"
	"dyndraw/framework/transform"

	gl "github.com/go-gl/gl/v3.1/gles2"
	. "github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

//testing out rendering pipeline within an actor context
type Gopher struct {
	transform.Transform
	sprite *render.Sprite
}

func (g *Gopher) Init() {
	g.sprite = render.CreateSprite("gopog.png")
}
func (g *Gopher) Update() {
	g.Transform.Rotate(0, 0, 0.025)
}
func (g *Gopher) Draw() {
	g.sprite.Draw(g.Transform.Mat4())
}
func (g *Gopher) Destroy() {}

//go:embed resources
var Resources embed.FS

func main() {
	runtime.LockOSThread()

	render.Resources = Resources

	var width, height int32 = 1920, 1080
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
		sdl.WINDOW_FULLSCREEN|sdl.WINDOW_OPENGL)
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

	fmt.Println("Starting...")

	for i := 0; i < 3; i++ {
		g := &Gopher{Transform: transform.Origin}
		g.Transform.Translate(float32(i), 0, float32(i)/30.0)
		actors.Spawn(g)
	}

	timer := sdl.GetTicks()
	frame_ticks_passed := uint32(0)
	update_ticks_passed := uint32(0)
	frames_passed := 0

	running = true
	for running {
		//TODO: input manager for actors to poll
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

		now := sdl.GetTicks()
		frame_ticks_passed += now - timer
		update_ticks_passed += now - timer
		frames_passed += 1
		timer = now
		if frame_ticks_passed > 1000 {
			window.SetTitle(strconv.Itoa(int(frames_passed)))
			frame_ticks_passed = 0
			frames_passed = 0
		}

		//fixed timestep of 125fps for updates for the smooth
		if update_ticks_passed > 8 {
			actors.Update()
			update_ticks_passed = 0
		}

		//rendering at...  however fast it can go
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		render.View = LookAtV(cameraPos, cameraPos.Add(Vec3{0, 0, -1}), Vec3{0, 1, 0})

		actors.All(func(a actors.Actor) {
			if r, ok := a.(render.Render); ok {
				r.Draw()
			}
		})
		render.DrawSprites()

		window.GLSwap()
	}

	actors.Quit()
}
