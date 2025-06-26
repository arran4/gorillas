//go:build !test

package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const sunRadius = 20

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
	clr := color.RGBA{255, 255, 0, 255}
	if g.sunHitTicks > 0 {
		clr = color.RGBA{255, 100, 100, 255}
	}
	drawFilledCircle(img, g.sunX, g.sunY, sunRadius, clr)
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

type building struct {
	x, w, h float64
	color   color.Color
	windows []window
}

type Game struct {
	*gorillas.Game
	buildings   []building
	sunX, sunY  float64
	sunHitTicks int
	angleInput  string
	powerInput  string
	enteringAng bool
	enteringPow bool
	bananaLeft  *ebiten.Image
	bananaRight *ebiten.Image
	bananaUp    *ebiten.Image
	bananaDown  *ebiten.Image
	gorillaArt  [][]string
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
	g.bananaLeft, g.bananaRight, g.bananaUp, g.bananaDown = createBananaSprites()
	return g
}

func (g *Game) Update() error {
	if !g.Banana.Active && !g.Explosion.Active {
		if g.enteringAng || g.enteringPow {
			for _, k := range inpututil.AppendJustPressedKeys(nil) {
				switch k {
				case ebiten.KeyEnter:
					if g.enteringAng {
						if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Angle = float64(v)
						}
						g.enteringAng = false
						g.angleInput = ""
						g.enteringPow = true
					} else {
						if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Power = float64(v)
						}
						g.enteringPow = false
						g.powerInput = ""
						g.Throw()
					}
				case ebiten.KeyEscape:
					g.enteringAng = false
					g.enteringPow = false
					g.angleInput = ""
					g.powerInput = ""
				case ebiten.KeyBackspace:
					if g.enteringAng && len(g.angleInput) > 0 {
						g.angleInput = g.angleInput[:len(g.angleInput)-1]
					} else if g.enteringPow && len(g.powerInput) > 0 {
						g.powerInput = g.powerInput[:len(g.powerInput)-1]
					}
				default:
					if k >= ebiten.Key0 && k <= ebiten.Key9 {
						r := '0' + rune(k-ebiten.Key0)
						if g.enteringAng && len(g.angleInput) < 3 {
							g.angleInput += string(r)
						} else if g.enteringPow && len(g.powerInput) < 3 {
							g.powerInput += string(r)
						}
					}
				}
			}
			return nil
		}
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			if k >= ebiten.Key0 && k <= ebiten.Key9 {
				g.enteringAng = true
				g.angleInput = string('0' + rune(k-ebiten.Key0))
				return nil
			}
		}
		if g.Current == 0 {
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				g.Angle += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				g.Angle -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				g.Power += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				g.Power -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				g.Throw()
			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				g.Angle += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				g.Angle -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				g.Power += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				g.Power -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyF) {
				g.Throw()
			}
		}
	} else {
		g.Step()
		if g.Banana.Active {
			if g.Banana.X >= g.sunX-sunRadius && g.Banana.X <= g.sunX+sunRadius &&
				g.Banana.Y >= g.sunY-sunRadius && g.Banana.Y <= g.sunY+sunRadius {
				g.sunHitTicks = 10
			}
		}
	}
	if g.sunHitTicks > 0 {
		g.sunHitTicks--
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	for i, b := range g.buildings {
		ebitenutil.DrawRect(screen, b.x, float64(g.Height)-b.h, b.w-1, b.h, b.color)
		for _, w := range b.windows {
			ebitenutil.DrawRect(screen, w.x, w.y, w.w, w.h, color.RGBA{255, 255, 0, 255})
		}
		_ = i
	}
	for i := range g.Gorillas {
		g.drawGorilla(screen, i)
	}
	if g.Banana.Active {
		dir := 0
		if math.Abs(g.Banana.VX) > math.Abs(g.Banana.VY) {
			if g.Banana.VX < 0 {
				dir = 0
			} else {
				dir = 1
			}
		} else {
			if g.Banana.VY < 0 {
				dir = 2
			} else {
				dir = 3
			}
		}
		var img *ebiten.Image
		switch dir {
		case 0:
			img = g.bananaLeft
		case 1:
			img = g.bananaRight
		case 2:
			img = g.bananaUp
		case 3:
			img = g.bananaDown
		}
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			w, h := img.Size()
			op.GeoM.Translate(g.Banana.X-float64(w)/2, g.Banana.Y-float64(h)/2)
			screen.DrawImage(img, op)
		}
	}
	if g.Explosion.Active {
		clr := color.RGBA{255, 255, 0, 255}
		if len(g.Explosion.Colors) > g.Explosion.Frame {
			clr = color.RGBAModel.Convert(g.Explosion.Colors[g.Explosion.Frame]).(color.RGBA)
		}
		drawFilledCircle(screen, g.Explosion.X, g.Explosion.Y, g.Explosion.Radii[g.Explosion.Frame], clr)
	}
	g.drawSun(screen)
	g.drawWindArrow(screen)
	angleStr := fmt.Sprintf("%3.0f", g.Angle)
	if g.enteringAng {
		if g.angleInput == "" {
			angleStr = "_"
		} else {
			angleStr = g.angleInput
		}
	}
	powerStr := fmt.Sprintf("%3.0f", g.Power)
	if g.enteringPow {
		if g.powerInput == "" {
			powerStr = "_"
		} else {
			powerStr = g.powerInput
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("A:%3s P:%3s W:%+2.0f P%d %d-%d", angleStr, powerStr, g.Wind, g.Current+1, g.Wins[0], g.Wins[1]))
}

func (g *Game) drawGorilla(img *ebiten.Image, idx int) {
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
	if settings.ShowIntro {
		play, err := introScreen(settings.UseSound, settings.UseSlidingText)
		if err != nil {
			panic(fmt.Errorf("intro screen: %w", err))
		}
		if !play {
			return
		}
	}
	game := newGame(settings, *buildings, *wind)
	game.Players = [2]string{*p1, *p2}
	if err := ebiten.RunGame(game); err != nil {
		panic(fmt.Errorf("run game: %w", err))
	}
	game.SaveScores()
	if err := showStats(game.StatsString()); err != nil {
		fmt.Fprintf(os.Stderr, "show stats: %v\n", err)
	}
	if game.League != nil {
		if err := showLeague(game.League); err != nil {
			fmt.Fprintf(os.Stderr, "show league: %v\n", err)
		}
	}
	fmt.Println(game.StatsString())
	showExtro()
}
