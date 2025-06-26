package gorillas

import (
	"math"
	"testing"
)

func newTestGame() *Game {
	g := &Game{Width: 100, Height: 100}
	g.Settings = DefaultSettings()
	bw := float64(g.Width) / BuildingCount
	for i := 0; i < BuildingCount; i++ {
		g.Buildings = append(g.Buildings, Building{X: float64(i) * bw, W: bw, H: 0})
	}
	g.Gorillas[0] = Gorilla{g.Buildings[1].X + bw/2, float64(g.Height) - g.Buildings[1].H}
	g.Gorillas[1] = Gorilla{g.Buildings[BuildingCount-2].X + bw/2, float64(g.Height) - g.Buildings[BuildingCount-2].H}
	return g
}

func almostEqual(a, b float64) bool {
	if a > b {
		return a-b < 1e-6
	}
	return b-a < 1e-6
}

func TestBananaTrajectoryAndOutOfBounds(t *testing.T) {
	g := newTestGame()
	g.Angle = 45
	g.Power = 100
	g.Current = 0
	startX := g.Gorillas[0].X
	startY := g.Gorillas[0].Y

	g.Throw()

	if !g.Banana.Active {
		t.Fatal("banana should be active after throw")
	}

	vx := math.Cos(g.Angle*math.Pi/180) * (g.Power / 2)
	vy := -math.Sin(g.Angle*math.Pi/180) * (g.Power / 2)

	if !almostEqual(g.Banana.VX, vx) || !almostEqual(g.Banana.VY, vy) {
		t.Fatalf("unexpected initial velocity got (%f,%f)", g.Banana.VX, g.Banana.VY)
	}

	g.Step()
	if !almostEqual(g.Banana.X, startX+vx) || !almostEqual(g.Banana.Y, startY+vy) {
		t.Fatalf("unexpected position after first step: (%f,%f)", g.Banana.X, g.Banana.Y)
	}
	if !almostEqual(g.Banana.VY, vy+0.5) {
		t.Fatalf("unexpected vy after first step: %f", g.Banana.VY)
	}
	if !g.Banana.Active {
		t.Fatal("banana should still be active after first step")
	}

	g.Step()
	if !g.Banana.Active {
		t.Fatal("banana should still be active after second step")
	}

	g.Step() // this should leave the screen
	if g.Banana.Active {
		t.Fatal("banana should be inactive after leaving screen")
	}
	if g.Current != 1 {
		t.Fatalf("turn should switch to player 2 after miss, got %d", g.Current)
	}
}

func TestBuildingCollisionEndsTurn(t *testing.T) {
	g := newTestGame()
	// make building 2 tall so banana will collide
	g.Buildings[2].H = 50

	g.Angle = 0
	g.Power = 20
	g.Current = 0

	g.Throw()
	g.Step()

	if g.Banana.Active {
		t.Fatal("banana should deactivate after hitting building")
	}
	if g.Current != 1 {
		t.Fatalf("turn should switch to player 2 after hitting building, got %d", g.Current)
	}
}

func TestGorillaHitIncrementsWin(t *testing.T) {
	g := newTestGame()
	g.Angle = 45
	g.Power = 100
	g.Current = 0

	startX := g.Gorillas[0].X
	startY := g.Gorillas[0].Y
	vx := math.Cos(g.Angle*math.Pi/180) * (g.Power / 2)
	vy := -math.Sin(g.Angle*math.Pi/180) * (g.Power / 2)
	// place second gorilla where the banana will be after one step
	g.Gorillas[1] = Gorilla{X: startX + vx, Y: startY + vy}

	g.Throw()
	g.Step()

	if g.Wins[0] != 1 {
		t.Fatalf("expected player 1 to score, wins: %v", g.Wins)
	}
	if g.Banana.Active {
		t.Fatal("banana should be inactive after hitting gorilla")
	}
	if !g.Explosion.Active {
		t.Fatal("explosion should be active after gorilla hit")
	}

	for g.Explosion.Active {
		g.Step()
	}

	if g.Current != 0 {
		t.Fatalf("current player should remain the same after win, got %d", g.Current)
	}
	if g.Wins[0] != 1 {
		t.Fatalf("wins should persist after reset, got %v", g.Wins)
	}
}
