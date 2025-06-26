//go:build !test

package main

import (
	"image/color"
	"time"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// instructionsState displays credits and instructions.
type instructionsState struct {
	lines   []string
	sliding bool
	charIdx int
	maxLen  int
	next    time.Time
	done    bool
}

func newInstructionsState(sliding bool) *instructionsState {
	lines, err := gorillas.LoadInstructions()
	if err != nil {
		lines = []string{"Instructions unavailable"}
	}
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	return &instructionsState{lines: lines, sliding: sliding, maxLen: maxLen, next: time.Now()}
}

func (s *instructionsState) Update(g *Game) error {
	if s.done {
		if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
			g.State = newMenuState(g.Settings.UseSound, g.Settings.UseSlidingText)
		}
		return nil
	}
	if !s.sliding {
		s.charIdx = s.maxLen
		s.done = true
		return nil
	}
	if time.Now().After(s.next) {
		s.charIdx++
		s.next = time.Now().Add(30 * time.Millisecond)
		if s.charIdx > s.maxLen {
			s.charIdx = s.maxLen
			s.done = true
		}
	}
	return nil
}

func (s *instructionsState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	y0 := g.Height/2 - len(s.lines)*charH/2
	for i, line := range s.lines {
		draw := line
		if s.sliding && s.charIdx < len(line) {
			if s.charIdx <= 0 {
				draw = ""
			} else {
				draw = line[:s.charIdx]
			}
		}
		ebitenutil.DebugPrintAt(screen, draw, (g.Width-len(line)*charW)/2, y0+i*charH)
	}
}
