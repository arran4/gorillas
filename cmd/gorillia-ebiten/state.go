package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// State defines game state behaviour.
type State interface {
	Update(g *Game) error
	Draw(g *Game, screen *ebiten.Image)
}
