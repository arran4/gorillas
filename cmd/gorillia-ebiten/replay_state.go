//go:build !test

package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type replayState struct {
	lines []string
	start time.Time
	phase int
	next  time.Time
}

func newReplayState(g *Game) *replayState {
	lines := []string{"Stored Shots:"}
	for i, sh := range g.ShotHistory {
		orig := g.Wind
		g.Wind = sh.Wind
		res := g.testShot(sh.Angle, sh.Power)
		g.Wind = orig
		lines = append(lines, fmt.Sprintf("%2d: %.0f %.0f wind %.1f -> %v", i+1, sh.Angle, sh.Power, sh.Wind, res))
	}
	if len(lines) == 1 {
		lines = append(lines, "No shots recorded")
	}
	lines = append(lines, "", "Press any key to continue")
	return &replayState{lines: lines}
}

func (s *replayState) Update(g *Game) error {
	if s.start.IsZero() {
		s.start = time.Now()
		s.next = s.start.Add(50 * time.Millisecond)
		return nil
	}
	if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
		g.State = newMenuState(g.Settings.UseSound, g.Settings.UseSlidingText)
		return nil
	}
	if time.Now().After(s.next) {
		s.phase = (s.phase + 1) % 5
		s.next = time.Now().Add(50 * time.Millisecond)
	}
	return nil
}

func (s *replayState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	pattern := []rune("*    ")
	cols := g.Width / charW
	rows := g.Height / charH
	for x := 0; x < cols; x++ {
		ch1 := pattern[(s.phase+x)%5]
		ch2 := pattern[(4-s.phase+x)%5]
		ebitenutil.DebugPrintAt(screen, string(ch1), x*charW, 0)
		ebitenutil.DebugPrintAt(screen, string(ch2), x*charW, (rows-1)*charH)
	}
	for y := 1; y < rows-1; y++ {
		ch := ' '
		if (s.phase+y)%5 == 0 {
			ch = '*'
		}
		ebitenutil.DebugPrintAt(screen, string(ch), (cols-1)*charW, y*charH)
		ebitenutil.DebugPrintAt(screen, string(ch), 0, (rows-1-y)*charH)
	}
	y0 := rows/2 - len(s.lines)/2
	for i, line := range s.lines {
		ebitenutil.DebugPrintAt(screen, line, (g.Width-len(line)*charW)/2, (y0+i)*charH)
	}
}
