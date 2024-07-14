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

type Cell struct {
	drawable uint32
	x        int
	y        int

	alive     bool
	aliveNext bool
}

func (c *Cell) Draw() {
	if !c.alive {
		return
	}
	gl.BindVertexArray(c.drawable)
	//todo fix this dependency on square
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func NewCell(x, y, totalX, totalY int) *Cell {
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

	return &Cell{
		drawable: makeVertexArrayObject(points),
		x:        x,
		y:        y,
	}
}

func (c *Cell) CheckState(cells [][]*Cell) {
	c.alive = c.aliveNext

	c.aliveNext = c.alive

	liveCount := c.liveSiblings(cells)
	if c.alive {
		if liveCount == 2 || liveCount == 3 {
			c.aliveNext = true
		} else {
			c.aliveNext = false
		}
	} else {
		if liveCount == 3 {
			c.aliveNext = true
		}
	}

}

func (c *Cell) liveSiblings(cells [][]*Cell) int {
	var liveCount int
	siblings := c.getSiblings(cells)
	for i := range siblings {
		if siblings[i].alive {
			liveCount += 1
		}
	}

	return liveCount
}

func (c *Cell) SetAlive(isAlive bool) {
	c.alive = isAlive
	c.aliveNext = isAlive
}

func (c *Cell) getSiblings(cells [][]*Cell) []*Cell {
	var siblings []*Cell
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if (x == 0) && (y == 0) {
				//ignore self
				continue
			}

			row := wrap(c.x+x, len(cells))
			col := wrap(c.y+y, len(cells[row]))
			siblings = append(siblings, cells[row][col])
		}
	}
	return siblings
}

func wrap(current, max int) int {
	if current >= max {
		return 0
	}
	if current < 0 {
		return max - 1
	}
	return current
}

func isY(index int) bool {
	return (index % 3) == 1
}
func isX(index int) bool {
	return (index % 3) == 0
}

func isZ(index int) bool {
	return (index % 3) == 2
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
