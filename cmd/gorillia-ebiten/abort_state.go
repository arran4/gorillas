//go:build !test

package main

import ebiten "github.com/hajimehoshi/ebiten/v2"

// abortState ends the game when activated.
type abortState struct{}

func newAbortState() *abortState { return &abortState{} }

func (abortState) Update(g *Game) error {
	g.Game.Aborted = true
	return ebiten.Termination
}

func (abortState) Draw(_ *Game, _ *ebiten.Image) {}
