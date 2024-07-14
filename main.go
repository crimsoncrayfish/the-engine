package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"the_engine/cell"
)

const (
	width  = 500
	height = 500
	rows   = 10
	cols   = 10

	//starting position
	seed      = 420
	threshold = 0.15

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
	seedPtr := flag.Int("seed", seed, "an int")
	widthPtr := flag.Int("width", width, "an int")
	heightPtr := flag.Int("height", height, "an int")
	rowsPtr := flag.Int("rows", rows, "an int")
	colsPtr := flag.Int("cols", cols, "an int")
	flag.Parse()
	fmt.Printf("Starting Conways's game of life with seed %v\n", int64(*seedPtr))

	runtime.LockOSThread()

	window := initGlfw(*widthPtr, *heightPtr)
	defer glfw.Terminate()

	program := initOpenGl()

	previousTime := time.Now().UnixMilli()
	//direction := float32(1)

	cells := makeCells(*rowsPtr, *colsPtr, int64(*seedPtr))

	//run through blocks over time
	//	currentCol := 0
	//	currentRow := 0
	for !window.ShouldClose() {
		//rotating shapes need time to rotate
		//newTime := time.Now().UnixMilli()
		//milliseconds := newTime - previousTime
		//degrees := direction * float32(milliseconds) / 1000.0 * 6.0 * speed
		//if math.Abs(float64(degrees)) >= 360 {
		//previousTime = newTime
		//direction = direction * -1
		//}
		//currentTriangle := newVector(triangle, float32(degrees))
		//vao := makeVertexArrayObject(currentTriangle)

		//run through blocks over time
		//newTime := time.Now().UnixMilli()
		//milliseconds := newTime - previousTime
		//if milliseconds > 100 {
		//	currentCol, currentRow = nextBlock(cells, currentCol, currentRow)
		//	previousTime = newTime
		//}
		//drawACell(currentCol, currentRow, cells, window, program)

		//conway's game of life
		newTime := time.Now().UnixMilli()
		millisecondsPassed := newTime - previousTime
		if millisecondsPassed > 100 {
			for x := range cells {
				for _, c := range cells[x] {
					c.CheckState(cells)
				}
			}
			previousTime = newTime
		}
		draw(cells, window, program)
	}
}

func nextBlock(cells [][]*cell.Cell, currentCol, currentRow int) (int, int) {
	newCol := currentCol + 1
	newRow := currentRow

	if newCol > len(cells)-1 {
		newCol = 0
		newRow += 1
		if newRow > len(cells[0])-1 {
			newRow = 0
		}
	}

	return newCol, newRow

}

func initGlfw(width, height int) *glfw.Window {
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

func draw(cells [][]*cell.Cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// for 1 shap do the following
	//	gl.BindVertexArray(vao)
	//	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape)/3))

	//for a checkerboard
	//checkerBoard(cells)
	for row := range cells {
		for _, cell := range cells[row] {
			cell.Draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func drawACell(col, row int, cells [][]*cell.Cell, window *glfw.Window, prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	cells[col][row].Draw()
	glfw.PollEvents()
	window.SwapBuffers()
}

func checkerBoard(cells [][]*cell.Cell) {
	for row := range cells {
		for col, cell := range cells[row] {
			rowMod := math.Mod(float64(row), 2)
			colMod := math.Mod(float64(col), 2)
			if math.Mod(rowMod+colMod, 2) == 0 {
				cell.Draw()
			}
		}
	}
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

func makeCells(rows, cols int, seed int64) [][]*cell.Cell {
	cells := make([][]*cell.Cell, cols)
	r := rand.New(rand.NewSource(seed))
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			c := cell.NewCell(x, y, cols, rows)

			isAlive := r.Float64() < threshold
			c.SetAlive(isAlive)

			cells[x] = append(cells[x], c)
		}
	}
	return cells
}
