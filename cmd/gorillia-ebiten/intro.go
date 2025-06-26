//go:build !test

package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var gorillaFrames = [][]string{
	{
		" O ",
		"/|\\",
		"/ \\",
	},
	{
		" O ",
		"/| ",
		"/ \\",
	},
	{
		" O ",
		" |\\",
		"/ \\",
	},
}

const (
	charW = 6
	charH = 16
)

// introGame implements ebiten.Game to play the ASCII intro.
type introGame struct {
	useSound bool
	sliding  bool
	lines    []string
	width    int
	height   int
	stage    int
	lineIdx  int
	charIdx  int
	next     time.Time
	frame    int
	done     bool
}

func (g *introGame) Update() error {
	if g.done {
		return nil
	}
	now := time.Now()
	switch g.stage {
	case 0:
		if g.sliding {
			if now.After(g.next) {
				g.charIdx++
				g.next = now.Add(30 * time.Millisecond)
				if g.charIdx > len(g.lines[g.lineIdx]) {
					g.lineIdx++
					g.charIdx = 0
					if g.lineIdx >= len(g.lines) {
						g.stage = 1
						g.next = now.Add(1500 * time.Millisecond)
						if g.useSound {
							gorillas.PlayIntroMusic()
						}
					}
				}
			}
		} else {
			g.stage = 1
			g.next = now.Add(1500 * time.Millisecond)
			if g.useSound {
				gorillas.PlayIntroMusic()
			}
		}
	case 1:
		if now.After(g.next) {
			g.stage = 2
			g.next = now
		}
	case 2:
		if now.Sub(g.next) >= 300*time.Millisecond {
			g.frame++
			g.next = now
			if g.frame >= 4 {
				g.stage = 3
				g.next = now.Add(700 * time.Millisecond)
			}
		}
	case 3:
		if now.After(g.next) {
			g.done = true
		}
	}
	return nil
}

func (g *introGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	y0 := g.height/2 - charH
	for i, line := range g.lines {
		draw := line
		if g.stage == 0 && g.sliding {
			if i < g.lineIdx {
				draw = line
			} else if i == g.lineIdx {
				if g.charIdx <= len(line) {
					draw = line[:g.charIdx]
				} else {
					draw = line
				}
			} else {
				draw = ""
			}
		}
		ebitenutil.DebugPrintAt(screen, draw, (g.width-len(line)*charW)/2, y0+i*charH)
	}
	if g.stage >= 2 {
		f1 := gorillaFrames[g.frame%len(gorillaFrames)]
		f2 := gorillaFrames[(g.frame+1)%len(gorillaFrames)]
		x1 := g.width/2 - 10*charW
		x2 := g.width/2 + 2*charW
		y := g.height/2 + 2*charH
		for i, l := range f1 {
			ebitenutil.DebugPrintAt(screen, l, x1, y+i*charH)
		}
		for i, l := range f2 {
			ebitenutil.DebugPrintAt(screen, l, x2, y+i*charH)
		}
	}
}

func (g *introGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func newIntroGame(w, h int, useSound, sliding bool) *introGame {
	if w == 0 || h == 0 {
		w, h = 800, 600
	}
	return &introGame{
		useSound: useSound,
		sliding:  sliding,
		lines:    []string{"QBasic GORILLAS", "", "Starring two gorillas"},
		width:    w,
		height:   h,
		next:     time.Now(),
	}
}

// sparkleGame shows twinkling '*' borders and optional lines of text.
type sparkleGame struct {
	lines   []string
	width   int
	height  int
	timeout time.Duration
	start   time.Time
	phase   int
	next    time.Time
}

func (g *sparkleGame) Update() error {
	if g.start.IsZero() {
		g.start = time.Now()
		g.next = g.start.Add(50 * time.Millisecond)
		return nil
	}
	if g.timeout > 0 && time.Since(g.start) > g.timeout {
		return ebiten.Termination
	}
	if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
		return ebiten.Termination
	}
	if time.Now().After(g.next) {
		g.phase = (g.phase + 1) % 5
		g.next = time.Now().Add(50 * time.Millisecond)
	}
	return nil
}

func (g *sparkleGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	pattern := []rune("*    ")
	cols := g.width / charW
	rows := g.height / charH
	for x := 0; x < cols; x++ {
		ch1 := pattern[(g.phase+x)%5]
		ch2 := pattern[(4-g.phase+x)%5]
		ebitenutil.DebugPrintAt(screen, string(ch1), x*charW, 0)
		ebitenutil.DebugPrintAt(screen, string(ch2), x*charW, (rows-1)*charH)
	}
	for y := 1; y < rows-1; y++ {
		ch := ' '
		if (g.phase+y)%5 == 0 {
			ch = '*'
		}
		ebitenutil.DebugPrintAt(screen, string(ch), (cols-1)*charW, y*charH)
		ebitenutil.DebugPrintAt(screen, string(ch), 0, (rows-1-y)*charH)
	}
	if len(g.lines) > 0 {
		y0 := rows/2 - len(g.lines)/2
		for i, line := range g.lines {
			ebitenutil.DebugPrintAt(screen, line, (g.width-len(line)*charW)/2, (y0+i)*charH)
		}
	}
}

func (g *sparkleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

// SparklePause displays a star border for the specified duration. If lines are
// provided they are shown centred on the screen.
func SparklePause(lines []string, dur time.Duration) {
	w, h := ebiten.WindowSize()
	if w == 0 || h == 0 {
		w, h = 800, 600
	}
	sg := &sparkleGame{lines: lines, width: w, height: h, timeout: dur}
	_ = ebiten.RunGame(sg)
}

func showStats(stats string) {
	SparklePause(strings.Split(stats, "\n"), 0)
}

func showLeague(l *gorillas.League) {
	if l == nil {
		return
	}
	lines := []string{"Player           Rounds Wins Accuracy"}
	for _, s := range l.Standings() {
		lines = append(lines, fmt.Sprintf("%-15s %6d %4d %8.1f", s.Name, s.Rounds, s.Wins, s.Accuracy))
	}
	SparklePause(lines, 0)
}
