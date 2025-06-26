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
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const sunRadius = 20
const sunMaxIntegrity = 4

type window struct {
	x, y, w, h float64
}

func drawFilledCircle(img *ebiten.Image, cx, cy, r float64, clr color.Color) {
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				ebitenutil.DrawRect(img, cx+dx, cy+dy, 1, 1, clr)
			}
		}
	}
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
	drawFilledCircle(img, g.sunX, g.sunY, r, clr)
	ebitenutil.DrawRect(img, g.sunX-6, g.sunY-4, 3, 3, color.Black)
	ebitenutil.DrawRect(img, g.sunX+3, g.sunY-4, 3, 3, color.Black)
	if g.sunHitTicks > 0 {
		drawFilledCircle(img, g.sunX, g.sunY+6, 5, color.Black)
		drawFilledCircle(img, g.sunX, g.sunY+6, 3, clr)
	} else {
		ebitenutil.DrawRect(img, g.sunX-4, g.sunY+4, 8, 2, color.Black)
	}
}

func createBananaSprite(mask []string) *ebiten.Image {
	h := len(mask)
	w := len(mask[0])
	img := ebiten.NewImage(w, h)
	clr := color.RGBA{255, 255, 0, 255}
	for y, row := range mask {
		for x, c := range row {
			if c != '.' {
				img.Set(x, y, clr)
			}
		}
	}
	return img
}

func createBananaSprites() (left, right, up, down *ebiten.Image) {
	left = createBananaSprite([]string{
		"..##.",
		".###.",
		"#####",
		".###.",
		"..##.",
	})
	right = createBananaSprite([]string{
		".##..",
		".###.",
		"#####",
		".###.",
		".##..",
	})
	up = createBananaSprite([]string{
		"..#..",
		".###.",
		"..#..",
		"..#..",
		"..#..",
	})
	down = createBananaSprite([]string{
		"..#..",
		"..#..",
		"..#..",
		".###.",
		"..#..",
	})
	return
}

func createGorillaSprite(mask []string, clr color.Color) *ebiten.Image {
	h := len(mask)
	w := len(mask[0])
	img := ebiten.NewImage(w, h)
	for y, row := range mask {
		for x, c := range row {
			if c != '.' {
				img.Set(x, y, clr)
			}
		}
	}
	return img
}

func defaultGorillaSprite() *ebiten.Image {
	mask := []string{
		"..##..",
		".####.",
		"######",
		"##..##",
		"######",
		"######",
		"##..##",
		"##..##",
		".#..#.",
		".####.",
	}
	return createGorillaSprite(mask, color.RGBA{150, 75, 0, 255})
}

type building struct {
	x, w, h float64
	color   color.Color
	windows []window
}

type Game struct {
	*gorillas.Game
	gamepads    []ebiten.GamepadID
	buildings    []building
	sunX, sunY   float64
	sunHitTicks  int
	sunIntegrity int
	angleInput   string
	powerInput   string
	enteringAng  bool
	enteringPow  bool
	abortPrompt  bool
	resumeAng    bool
	resumePow    bool
	bananaLeft   *ebiten.Image
	bananaRight  *ebiten.Image
	bananaUp     *ebiten.Image
	bananaDown   *ebiten.Image
	gorillaImg   *ebiten.Image
	gorillaArt   [][]string
	State        State
}

func newGame(settings gorillas.Settings, buildings int, wind float64) *Game {
	g := &Game{Game: gorillas.NewGame(800, 600, buildings)}
	if !math.IsNaN(wind) {
		g.Game.Wind = wind
	}
	g.Game.Settings = settings
	if art, err := gorillas.LoadGorillaArt("assets/gorilla.txt"); err == nil {
		g.gorillaArt = art
	} else {
		g.gorillaArt = [][]string{{" O ", "/|\\", "/ \\"}}
	}
	g.gorillaImg = defaultGorillaSprite()
	g.LoadScores()
	rand.Seed(time.Now().UnixNano())
	bw := float64(g.Width) / float64(g.Game.BuildingCount)
	for i := 0; i < g.Game.BuildingCount; i++ {
		h := g.Buildings[i].H
		b := building{
			x:     float64(i) * bw,
			w:     bw,
			h:     h,
			color: color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255},
		}
		for wx := b.x + 3; wx < b.x+b.w-3; wx += 6 {
			for wy := float64(g.Height) - 3; wy > float64(g.Height)-b.h+3; wy -= 6 {
				if rand.Intn(3) != 0 {
					b.windows = append(b.windows, window{wx, wy, 3, 3})
				}
			}
		}
		g.buildings = append(g.buildings, b)
	}
	g.sunX = float64(g.Width) - 40
	g.sunY = 40
	g.sunIntegrity = sunMaxIntegrity
	g.Game.ResetHook = func() {
		g.sunIntegrity = sunMaxIntegrity
	}
	g.bananaLeft, g.bananaRight, g.bananaUp, g.bananaDown = createBananaSprites()
	g.gamepads = ebiten.AppendGamepadIDs(nil)
	return g
}

func (g *Game) Update() error {
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
		op.GeoM.Translate(g.Gorillas[idx].X-float64(w)/2, g.Gorillas[idx].Y-float64(h))
		img.DrawImage(g.gorillaImg, op)
		return
	}
	if len(g.gorillaArt) == 0 {
		gr := g.Gorillas[idx]
		ebitenutil.DrawRect(img, gr.X-5, gr.Y-10, 10, 10, color.RGBA{255, 0, 0, 255})
		return
	}
	frame := g.gorillaArt[0]
	baseX := int(g.Gorillas[idx].X) - len(frame[0])/2
	baseY := int(g.Gorillas[idx].Y) - len(frame)
	for dy, line := range frame {
		for dx, ch := range line {
			if ch != ' ' {
				ebitenutil.DrawRect(img, float64(baseX+dx), float64(baseY+dy), 1, 1, color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func (g *Game) drawWindArrow(img *ebiten.Image) {
	if g.Wind == 0 {
		return
	}
	length := g.Wind * 3 * float64(g.Width) / 320
	y := float64(g.Height) - float64(g.Height)/40
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
	flag.BoolVar(&settings.UseSound, "sound", settings.UseSound, "enable sound")
	flag.BoolVar(&settings.WinnerFirst, "winnerfirst", settings.WinnerFirst, "winner starts next round")
	flag.Parse()
	settings.DefaultGravity = *gravity
	settings.DefaultRoundQty = *rounds
	game := newGame(settings, *buildings, *wind)
	game.Players = [2]string{*p1, *p2}
	if settings.ShowIntro {
		game.State = newIntroMovieState(settings.UseSound, settings.UseSlidingText)
	} else {
		game.State = newMenuState(settings.UseSound, settings.UseSlidingText)
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(fmt.Errorf("run game: %w", err))
	}
	game.SaveScores()
	fmt.Println(game.StatsString())
	showExtro()
}
