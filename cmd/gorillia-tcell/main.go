package main

import (
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
	buildings []building
	screen             tcell.Screen
	buildings          []building
	gorillas           [2]int
	bananaX, bananaY   float64
	bananaVX, bananaVY float64
	bananaActive       bool
	sunX, sunY         int
	sunHitTicks        int
}

const buildingWidth = 8

func newGame() *Game {
	g := &Game{Game: gorillas.NewGame(80, 24)}
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
	g.gorillas[0] = 1
	g.gorillas[1] = len(g.buildings) - 2
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
	g.drawSun()
	s := fmt.Sprintf("A:%2.0f P:%2.0f P%d %d-%d", g.Angle, g.Power, g.Current+1, g.Wins[0], g.Wins[1])
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

func (g *Game) run(s tcell.Screen) error {
	g.screen = s

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		g.draw()
		select {
		case <-ticker.C:
			if g.Banana.Active {
				g.Step()
				g.bananaX += g.bananaVX
				g.bananaY += g.bananaVY
				g.bananaVY += 0.2
				if int(g.bananaX) >= g.sunX && int(g.bananaX) < g.sunX+3 && int(g.bananaY) >= g.sunY && int(g.bananaY) < g.sunY+3 {
					g.sunHitTicks = 10
				}
				idx := int(g.bananaX) / buildingWidth
				if idx >= 0 && idx < len(g.buildings) {
					if int(g.bananaY) >= g.Height-g.buildings[idx].h {
						g.bananaActive = false
						g.Current = (g.Current + 1) % 2
					}
					for _, gidx := range g.gorillas {
						if idx == gidx && int(g.bananaY) >= g.Height-g.buildings[gidx].h-2 {
							g.bananaActive = false
							g.Wins[g.Current]++
							next := g.Current
							g.Reset()
							g.Current = next
						}
					}
				}
				if int(g.bananaY) >= g.Height || int(g.bananaX) < 0 || int(g.bananaX) >= g.Width {
					g.bananaActive = false
					g.Current = (g.Current + 1) % 2
				}
			}
		default:
			ev := s.PollEvent()
			switch e := ev.(type) {
			case *tcell.EventKey:
				switch e.Key() {
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

	if !introScreen(s) {
		return
	}

	g := newGame()
	if err := g.run(s); err != nil {
		panic(err)
	}
}
