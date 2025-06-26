package gorillas

import (
	"math"
	"math/rand"
	"path/filepath"
	"testing"
)

func newTestGame() *Game {
	rand.Seed(1)
	g := NewGame(100, 100, DefaultBuildingCount)
	g.Settings = DefaultSettings()
	g.Gravity = g.Settings.DefaultGravity
	g.Wind = 0
	// tests expect no persistent league data
	g.League = nil
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
	expectedVY := vy + 0.5*(g.Settings.DefaultGravity/17)
	if !almostEqual(g.Banana.VY, expectedVY) {
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
	// make building 2 tall enough to block the banana
	g.Buildings[2].H = float64(g.Height) - g.Gorillas[0].Y + 5

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

func TestThrowAppliesWindFluctuation(t *testing.T) {
	g := newTestGame()
	g.Settings.WindFluctuations = true
	g.Wind = 5
	rand.Seed(2)
	g.Angle = 0
	g.Power = 20
	g.Current = 0
	g.Throw()
	if g.Wind != 4 {
		t.Fatalf("expected wind 4 got %f", g.Wind)
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
	g.Settings.WinnerFirst = true
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

func TestWinnerFirstDisabled(t *testing.T) {
	g := newTestGame()
	g.Angle = 45
	g.Power = 100
	g.Current = 0

	startX := g.Gorillas[0].X
	startY := g.Gorillas[0].Y
	vx := math.Cos(g.Angle*math.Pi/180) * (g.Power / 2)
	vy := -math.Sin(g.Angle*math.Pi/180) * (g.Power / 2)
	g.Gorillas[1] = Gorilla{X: startX + vx, Y: startY + vy}

	g.Throw()
	g.Step()
	for g.Explosion.Active {
		g.Step()
	}
	if g.Current != 1 {
		t.Fatalf("current should switch to player 2 when WinnerFirst is off, got %d", g.Current)
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
	g.Explosion.Radii = []float64{1, 2}
	if !g.Explosion.Active {
		t.Fatal("explosion should start active")
	}

	if g.Explosion.Frame != 0 {
		t.Fatalf("initial explosion Frame should be 0, got %d", g.Explosion.Frame)
	}

	g.Step()
	if g.Explosion.Frame != 1 {
		t.Fatalf("explosion Frame should advance, got %d", g.Explosion.Frame)
	}

	g.Step()
	if g.Explosion.Active {
		t.Fatal("explosion should finish and deactivate")
	}
}

func TestExplosionColorsMatchRadii(t *testing.T) {
	g := newTestGame()
	g.startGorillaExplosion(0)
	if g.Settings.UseOldExplosions {
		t.Skip("old explosions have no colours")
	}
	if len(g.Explosion.Colors) != len(g.Explosion.Radii) {
		t.Fatalf("colour frames %d do not match radii %d", len(g.Explosion.Colors), len(g.Explosion.Radii))
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

func TestForceCGAHalvesExplosionRadius(t *testing.T) {
	g := newTestGame()
	g.Settings.ForceCGA = true
	g.Settings.NewExplosionRadius = 20
	g.startGorillaExplosion(0)
	if len(g.Explosion.Radii) < 2 {
		t.Fatal("not enough explosion frames")
	}
	if g.Explosion.Radii[1] != 10 {
		t.Fatalf("expected radius 10 got %f", g.Explosion.Radii[1])
	}
}

func TestVictoryDanceStartsOnHit(t *testing.T) {
	g := newTestGame()
	g.Angle = 45
	g.Power = 100
	g.Current = 0

	startX := g.Gorillas[0].X
	startY := g.Gorillas[0].Y
	vx := math.Cos(g.Angle*math.Pi/180) * (g.Power / 2)
	vy := -math.Sin(g.Angle*math.Pi/180) * (g.Power / 2)
	g.Gorillas[1] = Gorilla{X: startX + vx, Y: startY + vy}

	g.Throw()
	g.Step()
	if !g.Dance.Active || g.Dance.idx != 0 {
		t.Fatal("victory dance should start for player 1")
	}
	baseY := g.Dance.baseY
	g.Step()
	if g.Gorillas[0].Y == baseY {
		t.Fatalf("expected gorilla Y to change during dance")
	}
	for g.Dance.Active {
		g.Step()
	}
	if g.Gorillas[0].Y != baseY {
		t.Fatalf("gorilla should return to base position")
	}
}

func TestNewGameWindUsesBasicAlgorithm(t *testing.T) {
	rand.Seed(1)
	g := NewGame(100, 100, DefaultBuildingCount)
	if g.Wind != -11 {
		t.Fatalf("expected wind -11 got %f", g.Wind)
	}
}

func TestVariableWindChangesEachRound(t *testing.T) {
	rand.Seed(1)
	g := NewGame(100, 100, DefaultBuildingCount)
	g.Settings = DefaultSettings()
	g.Settings.VariableWind = true
	initial := g.Wind
	// trigger round end immediately
	g.Explosion = Explosion{Active: true, Radii: []float64{1}}
	for g.Explosion.Active {
		g.Step()
	}
	if g.Wind == initial {
		t.Fatalf("wind should change each round")
	}
	if g.Wind != 10 {
		t.Fatalf("expected wind 10 got %f", g.Wind)
	}
}
