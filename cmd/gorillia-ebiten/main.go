//go:build !test

package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/arran4/gorillas"
	ebdraw "github.com/arran4/gorillas/drawings/ebiten"
	imgdraw "github.com/arran4/gorillas/drawings/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sunRadius          = 20 * sunScale
	sunMaxIntegrity    = 4
	digitBufferTimeout = 3 * time.Second
)

func drawVectorLines(img *ebiten.Image, pts []gorillas.VectorPoint, clr color.Color) {
	ebdraw.DrawVectorLines(img, pts, clr)
}

func (g *Game) drawSun(img *ebiten.Image) {
	if g.sunIntegrity <= 0 {
		return
	}
	clr := color.RGBA{255, 255, 0, 255}
	if g.sunHitTicks > 0 {
		clr = color.RGBA{255, 100, 100, 255}
	}
	r := float64(g.sunIntegrity) * sunRadius / sunMaxIntegrity
	ebdraw.DrawBASSun(img, g.sunX, g.sunY, r, g.sunHitTicks > 0, clr)
}

type Game struct {
	*gorillas.Game
	gamepads     []ebiten.GamepadID
	buildingBase []*ebiten.Image
	buildingImg  []*ebiten.Image
	sunX, sunY   float64
	sunHitTicks  int
	sunIntegrity int
	angleInput   string
	powerInput   string
	enteringAng  bool
	enteringPow  bool
	abortPrompt  bool
	selAngle     bool
	bananaLeft   *ebiten.Image
	bananaRight  *ebiten.Image
	bananaUp     *ebiten.Image
	bananaDown   *ebiten.Image
	gorillaImg   *ebiten.Image
	gorillaArt   [][]string
	AI           bool
	State        State
	lastDigit    time.Time
	// Closed indicates whether the window was closed by the user.
	Closed bool
}

func newGame(settings gorillas.Settings, buildings int, wind float64) *Game {
	g := &Game{Game: gorillas.NewGame(800, 600, buildings)}
	g.selAngle = true
	if !math.IsNaN(wind) {
		g.Game.Wind = wind
	}
	g.Game.Settings = settings
	if art, err := gorillas.LoadGorillaArt("assets/gorilla.txt"); err == nil {
		g.gorillaArt = art
	} else {
		g.gorillaArt = [][]string{{" O ", "/|\\", "/ \\"}}
	}
	gorillaBase := imgdraw.DefaultGorillaSprite(gorillaScale)
	g.gorillaImg = ebiten.NewImageFromImage(gorillaBase)
	if g.Game.HitMap != nil {
		for i, gr := range g.Game.Gorillas {
			g.Game.HitMap.ClearGorilla(int(gr.X), int(gr.Y), i, 4)
			g.Game.HitMap.DrawGorillaImage(int(gr.X), int(gr.Y), i, gorillaBase)
		}
	}
	g.LoadScores()
	rand.Seed(time.Now().UnixNano())
	bw := float64(g.Width) / float64(g.Game.BuildingCount)
	for i := 0; i < g.Game.BuildingCount; i++ {
		h := g.Buildings[i].H
		clr := color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255}
		base := ebdraw.CreateBuildingSprite(bw-1, h, clr)
		g.buildingBase = append(g.buildingBase, base)
		img := ebiten.NewImage(int(bw-1), int(h))
		g.buildingImg = append(g.buildingImg, img)
	}
	// centre the sun horizontally
	g.sunX = float64(g.Width) / 2
	g.sunY = 40
	g.sunIntegrity = sunMaxIntegrity
	g.Game.ResetHook = func() {
		g.sunIntegrity = sunMaxIntegrity
		g.buildingBase = g.buildingBase[:0]
		g.buildingImg = g.buildingImg[:0]
		bw := float64(g.Width) / float64(g.Game.BuildingCount)
		for i := 0; i < g.Game.BuildingCount; i++ {
			h := g.Buildings[i].H
			clr := color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255}
			base := ebdraw.CreateBuildingSprite(bw-1, h, clr)
			g.buildingBase = append(g.buildingBase, base)
			img := ebiten.NewImage(int(bw-1), int(h))
			g.buildingImg = append(g.buildingImg, img)
		}
	}
	g.bananaLeft, g.bananaRight, g.bananaUp, g.bananaDown = ebdraw.CreateBananaSprites()
	g.gamepads = ebiten.AppendGamepadIDs(nil)
	return g
}

func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Closed = true
		g.Aborted = true
		return ebiten.Termination
	}
	if g.State != nil {
		return g.State.Update(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.State != nil {
		g.State.Draw(g, screen)
	}
}

func (g *Game) drawGorilla(img *ebiten.Image, idx int) {
	if g.gorillaImg != nil {
		op := &ebiten.DrawImageOptions{}
		w, h := g.gorillaImg.Size()
		op.GeoM.Scale(gorillaScale, gorillaScale)
		op.GeoM.Translate(g.Gorillas[idx].X-float64(w)*gorillaScale/2, g.Gorillas[idx].Y-float64(h)*gorillaScale)
		img.DrawImage(g.gorillaImg, op)
		return
	}
	if len(g.gorillaArt) == 0 {
		gr := g.Gorillas[idx]
		ebitenutil.DrawRect(img, gr.X-5*gorillaScale, gr.Y-10*gorillaScale, 10*gorillaScale, 10*gorillaScale, color.RGBA{255, 0, 0, 255})
		return
	}
	frame := g.gorillaArt[0]
	width := gorillas.FrameWidth(frame)
	baseX := int(g.Gorillas[idx].X) - width*gorillaScale/2
	baseY := int(g.Gorillas[idx].Y) - len(frame)*gorillaScale
	for dy, line := range frame {
		for dx, ch := range line {
			if ch != ' ' {
				x := float64(baseX + dx*gorillaScale)
				y := float64(baseY + dy*gorillaScale)
				ebitenutil.DrawRect(img, x, y, gorillaScale, gorillaScale, color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func (g *Game) drawWindArrow(img *ebiten.Image) {
	if g.Wind == 0 {
		return
	}
	length := g.Wind * 3 * float64(g.Width) / 320
	// Position arrow near the top instead of the bottom
	y := float64(g.Height) / 40
	x := float64(g.Width) / 2
	end := x + length
	ebitenutil.DrawLine(img, x, y, end, y, color.RGBA{255, 255, 0, 255})
	head := 5.0
	if length > 0 {
		ebitenutil.DrawLine(img, end, y, end-head, y-3, color.RGBA{255, 255, 0, 255})
		ebitenutil.DrawLine(img, end, y, end-head, y+3, color.RGBA{255, 255, 0, 255})
	} else {
		ebitenutil.DrawLine(img, end, y, end+head, y-3, color.RGBA{255, 255, 0, 255})
		ebitenutil.DrawLine(img, end, y, end+head, y+3, color.RGBA{255, 255, 0, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

func main() {
	if err := increaseRLimit(); err != nil {
		fmt.Fprintf(os.Stderr, "increase rlimit: %v\n", err)
	}
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Gorillas Ebiten")
	settings := gorillas.LoadSettings()
	wind := flag.Float64("wind", math.NaN(), "initial wind")
	gravity := flag.Float64("gravity", settings.DefaultGravity, "gravity")
	rounds := flag.Int("rounds", settings.DefaultRoundQty, "round count")
	buildings := flag.Int("buildings", gorillas.DefaultBuildingCount, "building count")
	p1 := flag.String("player1", "Player 1", "name of player 1")
	p2 := flag.String("player2", "Player 2", "name of player 2")
	ai := flag.Bool("ai", false, "enable computer opponent")
	flag.BoolVar(&settings.UseSound, "sound", settings.UseSound, "enable sound")
	flag.BoolVar(&settings.WinnerFirst, "winnerfirst", settings.WinnerFirst, "winner starts next round")
	flag.Parse()
	settings.DefaultGravity = *gravity
	settings.DefaultRoundQty = *rounds
	game := newGame(settings, *buildings, *wind)
	game.AI = *ai
	game.Players = [2]string{*p1, *p2}
	if settings.ShowIntro {
		game.State = newIntroMovieState(settings.UseSound, settings.UseSlidingText)
	} else {
		game.State = newMenuState(settings.UseSound, settings.UseSlidingText)
	}
	winsBackup := game.TotalWins
	var playersBackup map[string]*gorillas.PlayerStats
	if game.League != nil {
		playersBackup = make(map[string]*gorillas.PlayerStats, len(game.League.Players))
		for n, ps := range game.League.Players {
			cp := *ps
			playersBackup[n] = &cp
		}
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(fmt.Errorf("run game: %w", err))
	}
	if game.Closed {
		return
	}
	if game.Aborted {
		game.TotalWins = winsBackup
		if game.League != nil {
			game.League.Players = playersBackup
			game.League.Save()
		}
		if err := SparklePause([]string{"Game aborted"}, 0); err != nil {
			panic(fmt.Errorf("sparkle pause: %w", err))
		}
		return
	}
	game.SaveScores()
	if err := showStats(game.StatsString()); err != nil {
		panic(fmt.Errorf("show stats: %w", err))
	}
	if game.League != nil {
		if err := showLeague(game.League); err != nil {
			panic(fmt.Errorf("show league: %w", err))
		}
	}
	fmt.Println(game.StatsString())
	showExtro()
}
