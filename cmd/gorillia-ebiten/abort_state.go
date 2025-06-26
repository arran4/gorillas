//go:build !test

package main

import "github.com/hajimehoshi/ebiten/v2"

// abortState ends the game when activated.
type abortState struct{}

func newAbortState() *abortState { return &abortState{} }

func (abortState) Update(g *Game) error {
	g.Aborted = true
	return ebiten.Termination
}

func (abortState) Draw(g *Game, screen *ebiten.Image) {}
