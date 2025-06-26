package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/arran4/gorillas"
	"github.com/gdamore/tcell/v2"
)

type building struct {
	h       int
	windows []int
}

type Game struct {
	*gorillas.Game
	buildings   []building
	screen      tcell.Screen
	sunX, sunY  int
	sunHitTicks int
}

const buildingWidth = 8

func newGame(settings gorillas.Settings, buildings int, wind float64) *Game {
	g := &Game{Game: gorillas.NewGame(80, 24, buildings)}
	if !math.IsNaN(wind) {
		g.Game.Wind = wind
	}
	g.Game.Settings = settings
	g.LoadScores()
	rand.Seed(time.Now().UnixNano())
	for _, b := range g.Buildings {
		var wins []int
		top := g.Height - int(b.H) + 2
		for y := g.Height - 2; y > top; y -= 2 {
			if rand.Intn(3) != 0 {
				wins = append(wins, y)
			}
		}
		g.buildings = append(g.buildings, building{h: int(b.H), windows: wins})
	}
	g.sunX = g.Width - 4
	g.sunY = 1
	return g
}

var (
	sunHappy = []string{`\|/`, `-o-`, `/|\`}
	sunShock = []string{`\|/`, `-O-`, `/|\`}
)

func (g *Game) drawSun() {
	art := sunHappy
	if g.sunHitTicks > 0 {
		art = sunShock
		g.sunHitTicks--
	}
	for dy, line := range art {
		for dx, r := range line {
			if r != ' ' {
				g.screen.SetContent(g.sunX+dx, g.sunY+dy, r, nil, tcell.StyleDefault)
			}
		}
	}
}

func (g *Game) draw() {
	g.screen.Clear()
	for i, b := range g.buildings {
		x := i*buildingWidth + 4
		for y := g.Height - 1; y >= g.Height-b.h; y-- {
			g.screen.SetContent(x, y, '#', nil, tcell.StyleDefault)
		}
		for _, wy := range b.windows {
			g.screen.SetContent(x, wy, 'o', nil, tcell.StyleDefault)
		}
	}
	g.drawGorilla(0)
	g.drawGorilla(1)
	// draw a simple sun
	g.screen.SetContent(g.Width-2, 1, 'O', nil, tcell.StyleDefault)
	if g.Banana.Active {
		g.screen.SetContent(int(g.Banana.X), int(g.Banana.Y), 'o', nil, tcell.StyleDefault)
	}
	if g.Explosion.Active {
		r := int(g.Explosion.radii[g.Explosion.frame])
		ex := int(g.Explosion.X)
		ey := int(g.Explosion.Y)
		for dx := -r; dx <= r; dx++ {
			for dy := -r; dy <= r; dy++ {
				if dx*dx+dy*dy <= r*r {
					x := ex + dx
					y := ey + dy
					if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
						g.screen.SetContent(x, y, '*', nil, tcell.StyleDefault)
					}
				}
			}
		}
	}
	g.drawSun()
	s := fmt.Sprintf("A:%2.0f P:%2.0f W:%+2.0f P%d %d-%d", g.Angle, g.Power, g.Wind, g.Current+1, g.Wins[0], g.Wins[1])
	for i, r := range s {
		g.screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
	}
	g.screen.Show()
}

func (g *Game) drawGorilla(idx int) {
	x := int(g.Gorillas[idx].X)
	y := int(g.Gorillas[idx].Y) - 1
	style := tcell.StyleDefault
	g.screen.SetContent(x, y-2, 'O', nil, style)
	g.screen.SetContent(x-1, y-1, '/', nil, style)
	g.screen.SetContent(x, y-1, '|', nil, style)
	g.screen.SetContent(x+1, y-1, '\\', nil, style)
	g.screen.SetContent(x-1, y, '/', nil, style)
	g.screen.SetContent(x+1, y, '\\', nil, style)
}

func (g *Game) throw() {
	g.Throw()
}

func (g *Game) run(s tcell.Screen, ai bool) error {
	g.screen = s

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		g.draw()
		if g.Banana.Active || g.Explosion.Active {
			<-ticker.C
			g.Step()
			if int(g.Banana.X) >= g.sunX && int(g.Banana.X) < g.sunX+3 && int(g.Banana.Y) >= g.sunY && int(g.Banana.Y) < g.sunY+3 {
				g.sunHitTicks = 10
			}
			continue
		}

		if ai && g.Current == 1 {
			g.AutoShot()
			continue
		}

		ev := s.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			switch key.Key() {
			case tcell.KeyEscape:
				return nil
			case tcell.KeyLeft:
				g.Angle += 1
			case tcell.KeyRight:
				g.Angle -= 1
			case tcell.KeyUp:
				g.Power += 1
			case tcell.KeyDown:
				g.Power -= 1
			case tcell.KeyEnter:
				g.throw()
			}
		}
	}
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = s.Init(); err != nil {
		panic(err)
	}
	defer s.Fini()

	settings := gorillas.LoadSettings()
	wind := flag.Float64("wind", math.NaN(), "initial wind")
	gravity := flag.Float64("gravity", settings.DefaultGravity, "gravity")
	rounds := flag.Int("rounds", settings.DefaultRoundQty, "round count")
	buildings := flag.Int("buildings", gorillas.DefaultBuildingCount, "building count")
	flag.BoolVar(&settings.UseSound, "sound", settings.UseSound, "enable sound")
	ai := flag.Bool("ai", false, "enable computer opponent")
	flag.Parse()
	settings.DefaultGravity = *gravity
	settings.DefaultRoundQty = *rounds

	if !introScreen(s, settings.UseSound) {
		return
	}

	g := newGame(settings)
	if err := g.run(s, *ai); err != nil {
		panic(err)
	}
	g.SaveScores()
	fmt.Println(g.StatsString())
}
