package cell

import (
	"github.com/go-gl/gl/v4.6-core/gl"
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

type cell struct {
	drawable uint32
	x        int
	y        int
}

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	//todo fix this dependency on square
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func newCell(x, y, totalX, totalY int) *cell {
	//todo is this a good idea -> len(square)? Hardcoded global variable
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var pos float32
		var size float32
		if isX(i) {
			size = 1.0 / float32(totalX)
			pos = float32(x) * size
		} else if isY(i) {
			size = 1.0 / float32(totalY)
			pos = float32(y) * size
		} else {
			//if is Z ignore
			continue
		}

		if points[i] < 0 {
			points[i] = (pos * 2) - 1
		} else {
			points[i] = ((pos + size) * 2) - 1
		}
	}

	return &cell{
		drawable: makeVertexArrayObject(points),
		x:        x,
		y:        y,
	}
}

func isY(index int) bool {
	return (index % 3) == 1
}
func isX(index int) bool {
	return (index % 3) == 0
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
