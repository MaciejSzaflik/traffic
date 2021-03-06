package gridMap

import (
	"time"

	"github.com/MaciejSzaflik/traffic/noise"
	"github.com/MaciejSzaflik/traffic/smallMath"
)

type Color struct {
	r, g, b float32
}

type GridMap struct {
	Vertices         []float32
	Indices          []uint32
	backgroundColors []Color
	groundValues     []float32
	Occupation       []int
	VerticesCount    int32
	Count            int
	vertDistanceSize float32
	VAO              uint32
	VAODirty         bool
}

func NewGridMap(vertDistanceSize float32, count int) *GridMap {
	gridMap := &GridMap{
		Count:            count,
		vertDistanceSize: vertDistanceSize,
		backgroundColors: make([]Color, count*count),
		Occupation:       make([]int, count*count),
	}
	gridMap.backgroundFromPerlin()
	gridMap.generateVertices()

	return gridMap
}

func (gm *GridMap) GetNeighbours() []int {
	return nil
}

func (gm *GridMap) backgroundFromPerlin() {
	gen := noise.NewNoiseGenerator(int(time.Now().Unix()))
	gm.groundValues = gen.PerlinArray(gm.Count, gm.Count, 6, 10)
	for i := range gm.groundValues {
		gm.groundValues[i] = smallMath.Map(gm.groundValues[i], 0, 1, 0.05, 0.4)
		gm.backgroundColors[i] = Color{gm.groundValues[i], gm.groundValues[i], gm.groundValues[i]}
	}
}

func (gm *GridMap) IncSpaceOccupaid(x, y int) {
	gm.Occupation[x*gm.Count+y]++
}

func (gm *GridMap) IncSpaceOccupaidIndex(index int) {
	gm.Occupation[index]++
}

func (gm *GridMap) DecSpaceOccupaidIndex(index int) {
	gm.Occupation[index]--
}

func (gm *GridMap) MoveToXYIfAllowed(x, y int, t *Traveler) {
	if gm.Occupation[x*gm.Count+y] > 0 {
		t.WaitAndCheckMyPatience()
		return
	}

	gm.MoveToXY(x, y, t.a)
	t.ClearWait()
}

func (gm *GridMap) MoveToXY(x, y int, agent *Agent) {
	index := agent.Pos.X*gm.Count + agent.Pos.Y
	gm.groundValues[index] += 0.01
	gm.backgroundColors[index] = Color{gm.groundValues[index], gm.groundValues[index], gm.groundValues[index]}

	c := gm.backgroundColors[index]
	gm.DecSpaceOccupaidIndex(index)
	gm.SetXYColor(agent.Pos.X, agent.Pos.Y, c)
	agent.Pos.X = x
	agent.Pos.Y = y

	gm.SetXYColor(agent.Pos.X, agent.Pos.Y, agent.color)
	gm.IncSpaceOccupaid(agent.Pos.X, agent.Pos.Y)
}

func (gm *GridMap) SetXYColor(x, y int, c Color) {
	i := x*gm.Count + y
	gm.Vertices[i*24+3] = c.r
	gm.Vertices[i*24+4] = c.g
	gm.Vertices[i*24+5] = c.b

	gm.Vertices[i*24+9] = c.r
	gm.Vertices[i*24+10] = c.g
	gm.Vertices[i*24+11] = c.b

	gm.Vertices[i*24+15] = c.r
	gm.Vertices[i*24+16] = c.g
	gm.Vertices[i*24+17] = c.b

	gm.Vertices[i*24+21] = c.r
	gm.Vertices[i*24+22] = c.g
	gm.Vertices[i*24+23] = c.b

	gm.VAODirty = true
}

func (gm *GridMap) SetXY(x, y int, r, g, b float32) {
	i := x*gm.Count + y
	gm.Vertices[i*24+3] = r
	gm.Vertices[i*24+4] = g
	gm.Vertices[i*24+5] = b

	gm.Vertices[i*24+9] = r
	gm.Vertices[i*24+10] = g
	gm.Vertices[i*24+11] = b

	gm.Vertices[i*24+15] = r
	gm.Vertices[i*24+16] = g
	gm.Vertices[i*24+17] = b

	gm.Vertices[i*24+21] = r
	gm.Vertices[i*24+22] = g
	gm.Vertices[i*24+23] = b

	gm.VAODirty = true
}

func (gm *GridMap) generateVertices() {
	gm.Vertices = []float32{}
	gm.Indices = []uint32{}
	gm.VerticesCount = int32(6 * gm.Count * gm.Count)
	start := -(gm.vertDistanceSize * float32(gm.Count)) / 2
	for i := 0; i < gm.Count; i++ {
		for j := 0; j < gm.Count; j++ {
			index := i*gm.Count + j
			fi := float32(i)*gm.vertDistanceSize + start
			fj := float32(j)*gm.vertDistanceSize + start
			c := gm.backgroundColors[index]
			i32 := uint32(index)*3 + uint32(index)
			gm.Vertices = append(gm.Vertices, []float32{
				fi, gm.vertDistanceSize + fj, 0.0,
				c.r, c.g, c.b,

				gm.vertDistanceSize + fi, fj, 0.0,
				c.r, c.g, c.b,

				fi, fj, 0.0,
				c.r, c.g, c.b,

				gm.vertDistanceSize + fi, gm.vertDistanceSize + fj, 0.0,
				c.r, c.g, c.b,
			}...)

			gm.Indices = append(gm.Indices, []uint32{
				i32, 1 + i32, 2 + i32,
				i32, 1 + i32, 3 + i32,
			}...)
		}
	}
}
