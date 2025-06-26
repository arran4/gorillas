//go:build !test

package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// menuState shows the intro menu.
type menuState struct {
	useSound bool
	sliding  bool
	stage    int
	frame    int
	next     time.Time
}

func newMenuState(useSound, sliding bool) *menuState {
	return &menuState{useSound: useSound, sliding: sliding, next: time.Now().Add(300 * time.Millisecond)}
}

func (m *menuState) Update(g *Game) error {
	now := time.Now()
	switch m.stage {
	case 0:
		if now.After(m.next) {
			m.frame++
			m.next = now.Add(300 * time.Millisecond)
			if m.frame >= 4 {
				m.stage = 1
				m.frame = 0
			}
		}
	case 1:
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			switch k {
			case ebiten.KeyQ:
				return ebiten.Termination
			case ebiten.KeyP:
				g.State = newSetupState(g)
				return nil
			case ebiten.KeyV:
				g.State = newIntroMovieState(m.useSound, m.sliding)
				return nil
			case ebiten.KeyI:
				g.State = newInstructionsState(m.sliding)
				return nil
			}
		}
	}
	return nil
}

func (m *menuState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	cx := g.Width/2 - 10*charW
	cy := g.Height/2 - 2*charH
	var f1, f2 []string
	if m.stage == 0 {
		f1 = gorillaFrames[m.frame%len(gorillaFrames)]
		f2 = gorillaFrames[(m.frame+1)%len(gorillaFrames)]
	} else {
		f1 = gorillaFrames[0]
		f2 = gorillaFrames[0]
	}
	for i, l := range f1 {
		ebitenutil.DebugPrintAt(screen, l, cx, cy+i*charH)
	}
	for i, l := range f2 {
		ebitenutil.DebugPrintAt(screen, l, cx+12*charW, cy+i*charH)
	}
	ebitenutil.DebugPrintAt(screen, "GORILLAS", (g.Width-8*charW)/2, cy-2*charH)
	if m.stage == 1 {
		line := "V/X - View Intro"
		ebitenutil.DebugPrintAt(screen, line, (g.Width-len(line)*charW)/2, cy+3*charH)
		line = "I - Instructions"
		ebitenutil.DebugPrintAt(screen, line, (g.Width-len(line)*charW)/2, cy+4*charH)
		line = "P/Start - Play Game"
		ebitenutil.DebugPrintAt(screen, line, (g.Width-len(line)*charW)/2, cy+5*charH)
		line = "Q/B - Quit"
		ebitenutil.DebugPrintAt(screen, line, (g.Width-len(line)*charW)/2, cy+6*charH)
	}
}
