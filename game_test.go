package gorillas

import (
	"math"
	"path/filepath"
	"testing"
)

func newTestGame() *Game {
	g := &Game{Width: 100, Height: 100}
	g.Settings = DefaultSettings()
	g.Gravity = g.Settings.DefaultGravity
	g.Wind = 0
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
	if !almostEqual(g.Banana.VY, vy+g.Gravity/34) {
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

func TestWindInfluencesVelocity(t *testing.T) {
	g := newTestGame()
	g.Angle = 0
	g.Power = 20
	g.Current = 0
	g.Wind = 4

	g.Throw()
	initialVX := g.Banana.VX
	g.Step()
	expectedVX := initialVX + g.Wind/20
	if !almostEqual(g.Banana.VX, expectedVX) {
		t.Fatalf("expected vx %f got %f", expectedVX, g.Banana.VX)
	}
}

func TestGravityInfluencesVelocity(t *testing.T) {
	g := newTestGame()
	g.Angle = 0
	g.Power = 20
	g.Current = 0
	g.Gravity = 34

	g.Throw()
	g.Step()
	if !almostEqual(g.Banana.VY, g.Gravity/34) {
		t.Fatalf("expected vy %f got %f", g.Gravity/34, g.Banana.VY)
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

func TestSecondPlayerThrowDirection(t *testing.T) {
	g := newTestGame()
	g.Current = 1
	g.Angle = 30
	g.Power = 40
	g.Throw()
	if g.Banana.VX >= 0 {
		t.Fatalf("player 2 banana should move left, got vx=%f", g.Banana.VX)
	}
}

func TestExplosionProgressAndReset(t *testing.T) {
	g := newTestGame()
	g.startGorillaExplosion(0)
	g.Explosion.radii = []float64{1, 2}
	if !g.Explosion.Active {
		t.Fatal("explosion should start active")
	}

	if g.Explosion.frame != 0 {
		t.Fatalf("initial explosion frame should be 0, got %d", g.Explosion.frame)
	}

	g.Step()
	if g.Explosion.frame != 1 {
		t.Fatalf("explosion frame should advance, got %d", g.Explosion.frame)
	}

	g.Step()
	if g.Explosion.Active {
		t.Fatal("explosion should finish and deactivate")
	}
}

func TestSaveAndLoadScores(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "scores.json")
	g1 := newTestGame()
	g1.ScoreFile = tmp
	g1.TotalWins = [2]int{2, 3}
	g1.SaveScores()

	g2 := newTestGame()
	g2.ScoreFile = tmp
	g2.LoadScores()

	if g2.TotalWins != g1.TotalWins {
		t.Fatalf("expected %v, got %v", g1.TotalWins, g2.TotalWins)
	}
}

func TestStatsString(t *testing.T) {
	g := newTestGame()
	g.Wins = [2]int{1, 2}
	g.TotalWins = [2]int{3, 4}
	expected := "Session - P1:1 P2:2\nOverall - P1:3 P2:4"
	if s := g.StatsString(); s != expected {
		t.Fatalf("unexpected stats string: %q", s)
	}
}
