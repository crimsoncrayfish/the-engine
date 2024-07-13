package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	acell "the-engine/cell"
)

const (
	width              = 500
	height             = 500
	speed              = 5
	rows               = 10
	cols               = 10
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

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGl()

	//previousTime := time.Now().UnixMilli()
	//direction := float32(1)

	cells := makeCells(rows, cols)

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

		draw(cells, window, program)
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

func draw(cells [][]*acell.cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	//	gl.BindVertexArray(vao)
	//	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape)/3))
	//TODO
	glfw.PollEvents()
	window.SwapBuffers()
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

func makeCells(rows, cols int) [][]*acell.cell {
	cells := make([][]*acell.cell, rows, cols)
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			c := acell.newCell(x, y, rows, cols)
			cells[x] = append(cells[x], c)
		}
	}
	return cells
}
