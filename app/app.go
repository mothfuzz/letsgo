package app

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/mothfuzz/letsgo/actors"
	"github.com/mothfuzz/letsgo/input"
	"github.com/mothfuzz/letsgo/render"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

//app lifecyle / sdl management

//there can only be one...
var window *sdl.Window
var context sdl.GLContext

func init() {
	//lol hacks
	runtime.LockOSThread()
}

func Init() {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := img.Init(img.INIT_JPG | img.INIT_PNG | img.INIT_TIF | img.INIT_WEBP); err != nil {
		panic(err)
	}

	if err := mix.Init(int(mix.INIT_MP3)); err != nil {
		panic(err)
	}
	if err = mix.OpenAudio(44100, uint16(mix.DEFAULT_FORMAT), 2, 512); err != nil {
		panic(err)
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)

	var w, h int32 = 640, 480
	window, err = sdl.CreateWindow("owo",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h,
		sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	if err = gl.Init(); err != nil {
		panic(err)
	}
	//sdl.GLSetSwapInterval(1)

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Disable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearColor(0.2, 0.2, 0.3, 1.0)
	gl.Viewport(0, 0, w, h)

	render.DefaultCamera.Init2D(w, h)

	fmt.Println("Starting...")
}

func SetWindowSize(width int32, height int32) {
	window.SetSize(width, height)
	gl.Viewport(0, 0, width, height)
	render.DefaultCamera.Init2D(width, height)
}
func GetWindowSize() (int32, int32) {
	return window.GetSize()
}

func SetFullScreen(fullscreen bool) {
	if fullscreen {
		window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	} else {
		window.SetFullscreen(0)
	}
}

func SetVSync(vsync bool) {
	if vsync {
		sdl.GLSetSwapInterval(1)
	} else {
		sdl.GLSetSwapInterval(0)
	}
}

func SetRelativeCursor(rc bool) {
	if rc {
		sdl.SetRelativeMouseMode(true)
	} else {
		sdl.SetRelativeMouseMode(false)
	}
}

func SetBackground(r, g, b float32) {
	gl.ClearColor(r, g, b, 1.0)
}

var timer uint32
var frameTicks uint32
var updateTicks uint32
var framesPassed uint32

func PollEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
				return false
			}
		}
	}
	return true
}
func Update() {
	now := sdl.GetTicks()
	frameTicks += now - timer
	updateTicks += now - timer
	framesPassed += 1
	timer = now
	if frameTicks > 1000 {
		window.SetTitle(strconv.Itoa(int(framesPassed)))
		frameTicks = 0
		framesPassed = 0
	}

	//fixed timestep of 125fps for updates for the smooth
	for updateTicks > 8 {
		//update input in here too, since it's per-tick
		input.UpdateKeys()
		updateTicks -= 8
		actors.Update()
	}
}

func Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	actors.All(func(a actors.Actor) {
		if r, ok := a.(render.Render); ok {
			r.Draw()
		}
	})
	render.DrawSprites()
	window.GLSwap()
}

func Quit() {
	window.Destroy()
	mix.CloseAudio()
	mix.Quit()
	img.Quit()
	sdl.GLDeleteContext(context)
	sdl.Quit()
	actors.DestroyAll()
}
