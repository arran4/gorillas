package gorillas

import (
	"math"
	"math/rand"
	"time"
)

type Building struct {
	X, W, H float64
}

type Gorilla struct {
	X, Y float64
}

type Banana struct {
	X, Y   float64
	VX, VY float64
	Active bool
}

type Game struct {
	Width, Height int
	Buildings     []Building
	Gorillas      [2]Gorilla
	Banana        Banana
	Angle         float64
	Power         float64
	Current       int
	Wins          [2]int
}

const BuildingCount = 10

func NewGame(width, height int) *Game {
	g := &Game{Width: width, Height: height, Angle: 45, Power: 50}
	rand.Seed(time.Now().UnixNano())
	bw := float64(width) / BuildingCount

	// create a sloping skyline similar to the original BASIC version
	slope := rand.Intn(6) + 1
	newHt := float64(height) * 0.3
	if slope == 2 || slope == 6 {
		newHt = float64(height) * 0.7
	}
	htInc := float64(height) / 20

	for i := 0; i < BuildingCount; i++ {
		x := float64(i) * bw
		switch slope {
		case 1:
			newHt += htInc
		case 2:
			newHt -= htInc
		case 3, 5:
			if x > float64(width)/2 {
				newHt -= 2 * htInc
			} else {
				newHt += 2 * htInc
			}
		case 4:
			if x > float64(width)/2 {
				newHt += 2 * htInc
			} else {
				newHt -= 2 * htInc
			}
		}

		h := newHt + rand.Float64()*float64(height)/8
		if h < float64(height)*0.1 {
			h = float64(height) * 0.1
		}
		if h > float64(height)*0.9 {
			h = float64(height) * 0.9
		}

		g.Buildings = append(g.Buildings, Building{
			X: x,
			W: bw,
			H: h,
		})
	}
	g.Gorillas[0] = Gorilla{g.Buildings[1].X + bw/2, float64(height) - g.Buildings[1].H}
	g.Gorillas[1] = Gorilla{g.Buildings[BuildingCount-2].X + bw/2, float64(height) - g.Buildings[BuildingCount-2].H}
	return g
}

func (g *Game) Reset() {
	wins := g.Wins
	*g = *NewGame(g.Width, g.Height)
	g.Wins = wins
}

func (g *Game) Throw() {
	start := g.Gorillas[g.Current]
	radians := g.Angle * math.Pi / 180
	speed := g.Power / 2
	g.Banana.X = start.X
	g.Banana.Y = start.Y
	if g.Current == 1 {
		g.Banana.VX = -math.Cos(radians) * speed
	} else {
		g.Banana.VX = math.Cos(radians) * speed
	}
	g.Banana.VY = -math.Sin(radians) * speed
	g.Banana.Active = true
}

func (g *Game) Step() {
	if !g.Banana.Active {
		return
	}
	g.Banana.X += g.Banana.VX
	g.Banana.Y += g.Banana.VY
	g.Banana.VY += 0.5
	idx := int(g.Banana.X / (float64(g.Width) / BuildingCount))
	if idx >= 0 && idx < BuildingCount {
		if g.Banana.Y > float64(g.Height)-g.Buildings[idx].H {
			g.Banana.Active = false
			g.Current = (g.Current + 1) % 2
		}
	}
	for i, gr := range g.Gorillas {
		if math.Abs(gr.X-g.Banana.X) < 5 && math.Abs(gr.Y-g.Banana.Y) < 10 {
			g.Banana.Active = false
			g.Wins[g.Current]++
			cur := g.Current
			g.Reset()
			g.Current = cur
			_ = i // to avoid unused variable if not used
		}
	}
	if g.Banana.Y > float64(g.Height) || g.Banana.X < 0 || g.Banana.X >= float64(g.Width) {
		g.Banana.Active = false
		g.Current = (g.Current + 1) % 2
	}
}
