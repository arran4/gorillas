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

// introMovieState replicates the original ASCII intro animation.
type introMovieState struct {
	useSound bool
	sliding  bool
	lines    []string
	stage    int
	lineIdx  int
	charIdx  int
	next     time.Time
	frame    int
	done     bool
}

func newIntroMovieState(useSound, sliding bool) *introMovieState {
	return &introMovieState{
		useSound: useSound,
		sliding:  sliding,
		lines:    []string{"QBasic GORILLAS", "", "Starring two gorillas"},
		next:     time.Now(),
	}
}

func (s *introMovieState) Update(g *Game) error {
	if s.done {
		g.State = newMenuState(s.useSound, s.sliding)
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.done = true
		g.State = newMenuState(s.useSound, s.sliding)
		return nil
	}
	now := time.Now()
	switch s.stage {
	case 0:
		if s.sliding {
			if now.After(s.next) {
				s.charIdx++
				s.next = now.Add(30 * time.Millisecond)
				if s.charIdx > len(s.lines[s.lineIdx]) {
					s.lineIdx++
					s.charIdx = 0
					if s.lineIdx >= len(s.lines) {
						s.stage = 1
						s.next = now.Add(1500 * time.Millisecond)
						if s.useSound {
							gorillas.PlayIntroMusic()
						}
					}
				}
			}
		} else {
			s.stage = 1
			s.next = now.Add(1500 * time.Millisecond)
			if s.useSound {
				gorillas.PlayIntroMusic()
			}
		}
	case 1:
		if now.After(s.next) {
			s.stage = 2
			s.next = now
		}
	case 2:
		if now.Sub(s.next) >= 300*time.Millisecond {
			s.frame++
			s.next = now
			if s.frame >= 4 {
				s.stage = 3
				s.next = now.Add(700 * time.Millisecond)
			}
		}
	case 3:
		if now.After(s.next) {
			s.done = true
		}
	}
	return nil
}

func (s *introMovieState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	y0 := g.Height/2 - charH
	for i, line := range s.lines {
		draw := line
		if s.stage == 0 && s.sliding {
			if i < s.lineIdx {
				draw = line
			} else if i == s.lineIdx {
				if s.charIdx <= len(line) {
					draw = line[:s.charIdx]
				} else {
					draw = line
				}
			} else {
				draw = ""
			}
		}
		ebitenutil.DebugPrintAt(screen, draw, (g.Width-len(line)*charW)/2, y0+i*charH)
	}
	if s.stage >= 2 {
		f1 := gorillaFrames[s.frame%len(gorillaFrames)]
		f2 := gorillaFrames[(s.frame+1)%len(gorillaFrames)]
		x1 := g.Width/2 - 10*charW
		x2 := g.Width/2 + 2*charW
		y := g.Height/2 + 2*charH
		for i, l := range f1 {
			ebitenutil.DebugPrintAt(screen, l, x1, y+i*charH)
		}
		for i, l := range f2 {
			ebitenutil.DebugPrintAt(screen, l, x2, y+i*charH)
		}
	}
}
