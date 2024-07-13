package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width              = 500
	height             = 500
	speed              = 5
	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
)

var (
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
	triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGl()

	//previousTime := time.Now().UnixMilli()
	//direction := float32(1)

	for !window.ShouldClose() {
		//newTime := time.Now().UnixMilli()
		//milliseconds := newTime - previousTime
		//degrees := direction * float32(milliseconds) / 1000.0 * 6.0 * speed
		//if math.Abs(float64(degrees)) >= 360 {
		//previousTime = newTime
		//direction = direction * -1
		//}

		//currentTriangle := newVector(triangle, float32(degrees))
		//vao := makeVertexArrayObject(currentTriangle)

		vao := makeVertexArrayObject(square)
		draw(square, vao, window, program)
	}
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGl() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertxSh, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragSh, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertxSh)
	gl.AttachShader(prog, fragSh)

	gl.LinkProgram(prog)
	return prog
}

func draw(shape []float32, vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVertexArrayObject(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newVector(points []float32, degrees float32) []float32 {
	center := []float32{0, 0, 0}
	point1 := points[0:2]
	point2 := points[3:5]
	point3 := points[6:8]

	point1 = rotatePoint(point1, center, degrees)
	point2 = rotatePoint(point2, center, degrees)
	point3 = rotatePoint(point3, center, degrees)
	newSlice := append(point1, point2...)
	newSlice = append(newSlice, point3...)
	return newSlice
}

func rotatePoint(point, center []float32, angleInDegrees float32) []float32 {
	radians := angleInDegrees * (math.Pi / 180)
	cosTheta := float32(math.Cos(float64(radians)))
	sinTheta := float32(math.Sin(float64(radians)))
	newPoint := make([]float32, 3)
	newPoint[0] = (cosTheta*(point[0]-center[0]) -
		sinTheta*(point[1]-center[1]) + center[0])
	newPoint[1] = (sinTheta*(point[0]-center[0]) +
		cosTheta*(point[1]-center[1]) + center[1])
	return newPoint
}
