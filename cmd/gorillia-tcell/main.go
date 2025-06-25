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
	screen    tcell.Screen
	buildings []building
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
	return g
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

func (g *Game) run() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = s.Init(); err != nil {
		return err
	}
	g.screen = s
	defer s.Fini()

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		g.draw()
		select {
		case <-ticker.C:
			if g.Banana.Active {
				g.Step()
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
	g := newGame()
	if err := g.run(); err != nil {
		panic(err)
	}
}
