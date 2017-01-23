package gl

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
)

var goldenRatio = 1.61803398875
var goldenFraction = float32(goldenRatio / (goldenRatio + 1))

var (
	GlfwWindow *glfw.Window // deprecate eventually
	Texture    uint32
)

//gfx in cGfx.CurrAppWidth
//cGfx.InitFrustum

//only two gfx parameters should be eliminated
//gfx import should be eliminated
//settings in either app or gfx

func init() {
}

func WindowInit() {
	fmt.Printf("Gl: Init glfw \n")

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	//defer glfw.Terminate()
	/*
	   Go's defer statement schedules a function call (the deferred function)
	   to be run immediately before the function executing the defer returns.
	   It's an unusual but effective way to deal with situations such as resources
	   that must be released regardless of which path a function takes to return.
	   The canonical examples are unlocking a mutex or closing a file.
	*/

	fmt.Printf("Gl: set windowhint\n")
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var err error
	GlfwWindow, err = glfw.CreateWindow(gfx.InitAppWidth, gfx.InitAppHeight, app.Name, nil, nil)

	if err != nil {
		panic(err)
	}

	GlfwWindow.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

}

func LoadTextures() {
	fmt.Printf("GL: load texture \n")
	Texture = NewTexture("Bisasam_24x24_Shadowed.png")
}

func InitRenderer() {
	fmt.Println("gl.InitRenderer()")

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	//gl.Enable(gl.ALPHA_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.5, 0.5, 0.5, 1}
	diffuse := []float32{1, 1, 1, 1}
	lightPosition := []float32{-5, 5, 10, 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	SetFrustum(gfx.InitFrustum)
	//gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	//gl.Frustum(left, right, bottom, top, zNear, zFar)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func SetFrustum(r *app.Rectangle) {
	gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top), 1.0, 10.0)
}

func Update() {
	gfx.Curs.Update()
}

func DrawScene() {
	gl.Viewport(0, 0, gfx.CurrAppWidth, gfx.CurrAppHeight) // OPTIMIZEME?  could set flag upon frame buffer size change event
	if *gfx.PrevFrustum != *gfx.CurrFrustum {
		*gfx.PrevFrustum = *gfx.CurrFrustum
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		SetFrustum(gfx.CurrFrustum)
		fmt.Println("CHANGE OF FRUSTUM")
	}
	gl.MatrixMode(gl.MODELVIEW) //.PROJECTION)                   //.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Translatef(0, 0, -gfx.DistanceFromOrigin)

	gl.BindTexture(gl.TEXTURE_2D, Texture)

	gl.Begin(gl.QUADS)
	drawAll()
	gl.End()
}

func SwapDrawBuffer() {
	GlfwWindow.SwapBuffers()
}

func ScreenTeardown() {
	glfw.Terminate()
}