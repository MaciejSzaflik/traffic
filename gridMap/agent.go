package gridMap

import (
	"math"
	"math/rand"
	"time"

	"github.com/MaciejSzaflik/traffic/smallMath"
)

const (
	Up int = iota
	Down
	Left
	Right

	UpLeft
	UpRight
	DownLeft
	DownRight
)

const maxFloat32 = float32(math.MaxFloat32)

type PointXY struct {
	X, Y int
}

type Agent struct {
	Pos   *PointXY
	color Color
}

type Traveler struct {
	a           *Agent
	GoalPos     *PointXY
	Speed       int
	turnCounter int
	waitTime    int
	patience    int
}

type AgentDirector struct {
	Travelers []*Traveler
	Points    []*Agent
	GridMap   *GridMap
}

func NewAgentDirectorRandom(grid *GridMap, travelersCount int, pointsCount int) *AgentDirector {
	director := &AgentDirector{
		Travelers: make([]*Traveler, travelersCount),
		Points:    make([]*Agent, pointsCount),
		GridMap:   grid,
	}

	for i := 0; i < pointsCount; i++ {
		director.Points[i] = &Agent{
			Pos:   &PointXY{rand.Intn(grid.Count), rand.Intn(grid.Count)},
			color: Color{1, 1, 0},
		}
		director.Points[i].Init(grid)
	}

	for i := 0; i < travelersCount; i++ {
		startPos := director.Points[rand.Intn(len(director.Points))].Pos
		director.Travelers[i] = &Traveler{
			a: &Agent{
				Pos:   &PointXY{startPos.X, startPos.Y},
				color: Color{0.588, 0.449, 1},
			},
			Speed:    rand.Intn(10) + 1,
			patience: rand.Intn(20) + 1,
		}
		director.Travelers[i].Init(grid)
		director.SetRandomGoal(director.Travelers[i])
	}

	return director
}

func (aD *AgentDirector) SetRandomGoal(traveler *Traveler) {
	traveler.GoalPos = aD.Points[rand.Intn(len(aD.Points))].Pos
	traveler.ClearWait()
}

func (ad *AgentDirector) Update(d time.Duration) {
	for _, point := range ad.Points {
		ad.GridMap.SetXYColor(point.Pos.X, point.Pos.Y, point.color)
	}

	for _, traveler := range ad.Travelers {
		if traveler.turnCounter < traveler.Speed {
			traveler.turnCounter++
			continue
		}
		traveler.turnCounter = 0

		if traveler.GoalPos == nil {
			ad.SetRandomGoal(traveler)
		}

		if DistancePos(traveler.a.Pos.X, traveler.GoalPos.X, traveler.a.Pos.Y, traveler.GoalPos.Y) == 0 {
			ad.SetRandomGoal(traveler)
		}

		if DistancePos(traveler.a.Pos.X, traveler.GoalPos.X, traveler.a.Pos.Y, traveler.GoalPos.Y) == 0 {
			continue
		}

		switch _, min := smallMath.MinFloat(
			DistanceMul(traveler.a.Pos.X, traveler.GoalPos.X, traveler.a.Pos.Y+1, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X, traveler.GoalPos.X, traveler.a.Pos.Y-1, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X-1, traveler.GoalPos.X, traveler.a.Pos.Y, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X+1, traveler.GoalPos.X, traveler.a.Pos.Y, traveler.GoalPos.Y, ad.GridMap),

			DistanceMul(traveler.a.Pos.X-1, traveler.GoalPos.X, traveler.a.Pos.Y+1, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X+1, traveler.GoalPos.X, traveler.a.Pos.Y+1, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X-1, traveler.GoalPos.X, traveler.a.Pos.Y-1, traveler.GoalPos.Y, ad.GridMap),
			DistanceMul(traveler.a.Pos.X+1, traveler.GoalPos.X, traveler.a.Pos.Y-1, traveler.GoalPos.Y, ad.GridMap),
		); min {
		case Up:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X, traveler.a.Pos.Y+1, traveler)
		case Down:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X, traveler.a.Pos.Y-1, traveler)
		case Left:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X-1, traveler.a.Pos.Y, traveler)
		case Right:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X+1, traveler.a.Pos.Y, traveler)
		case UpLeft:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X-1, traveler.a.Pos.Y+1, traveler)
		case UpRight:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X+1, traveler.a.Pos.Y+1, traveler)
		case DownLeft:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X-1, traveler.a.Pos.Y-1, traveler)
		case DownRight:
			ad.GridMap.MoveToXYIfAllowed(traveler.a.Pos.X+1, traveler.a.Pos.Y-1, traveler)
		}
	}
}

func (t *Traveler) WaitAndCheckMyPatience() {
	t.waitTime++

	t.a.color.g = 0.449 - float32(t.waitTime)/float32(t.patience)*0.449
	t.a.color.b = 1 - float32(t.waitTime)/float32(t.patience)

	if t.waitTime > t.patience {
		t.GoalPos = nil
	}
}

func (t *Traveler) ClearWait() {
	t.waitTime = 0
	t.a.color.r = 0.588
	t.a.color.g = 0.449
	t.a.color.b = 1
}

func (a *Agent) Init(gd *GridMap) {
	gd.SetXYColor(a.Pos.X, a.Pos.Y, a.color)
}

func (t *Traveler) Init(gd *GridMap) {
	gd.IncSpaceOccupaid(t.a.Pos.X, t.a.Pos.Y)
	gd.SetXYColor(t.a.Pos.X, t.a.Pos.Y, t.a.color)
}

func DistanceMul(x1, x2, y1, y2 int, gridMap *GridMap) float32 {
	if x1 < 0 || x1 >= gridMap.Count {
		return maxFloat32
	}
	if y1 < 0 || y1 >= gridMap.Count {
		return maxFloat32
	}

	return float32(DistancePos(x1, x2, y1, y2)) * (1 - gridMap.groundValues[x1*gridMap.Count+y1]*0.2)
}

func DistancePos(x1, x2, y1, y2 int) int {
	return smallMath.Abs(x1-x2) + smallMath.Abs(y1-y2)
}
