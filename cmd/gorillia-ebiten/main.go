package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type building struct {
	x, w, h float64
	color   color.Color
}

type Game struct {
	*gorillas.Game
	buildings []building
}

func newGame() *Game {
	g := &Game{Game: gorillas.NewGame(800, 600)}
	rand.Seed(time.Now().UnixNano())
	bw := float64(g.Width) / gorillas.BuildingCount
	for i := 0; i < gorillas.BuildingCount; i++ {
		h := g.Buildings[i].H
		g.buildings = append(g.buildings, building{
			x:     float64(i) * bw,
			w:     bw,
			h:     h,
			color: color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255},
		})
	}
	return g
}

func (g *Game) Update() error {
	if !g.Banana.Active {
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
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	for i, b := range g.buildings {
		ebitenutil.DrawRect(screen, b.x, float64(g.Height)-b.h, b.w-1, b.h, b.color)
		_ = i
	}
	for _, gr := range g.Gorillas {
		ebitenutil.DrawRect(screen, gr.X-5, gr.Y-10, 10, 10, color.RGBA{255, 0, 0, 255})
	}
	if g.Banana.Active {
		ebitenutil.DrawRect(screen, g.Banana.X-2, g.Banana.Y-2, 4, 4, color.RGBA{255, 255, 0, 255})
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("A:%2.0f P:%2.0f P%d %d-%d", g.Angle, g.Power, g.Current+1, g.Wins[0], g.Wins[1]))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Gorillas Ebiten")
	game := newGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
