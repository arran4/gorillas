package main

import (
	"image/color"
	"time"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
}

func (g *introGame) Update() error {
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
			return ebiten.Termination
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

func showIntroMovie(useSound, sliding bool) {
	w, h := ebiten.WindowSize()
	if w == 0 || h == 0 {
		w, h = 800, 600
	}
	ig := &introGame{
		useSound: useSound,
		sliding:  sliding,
		lines:    []string{"QBasic GORILLAS", "", "Starring two gorillas"},
		width:    w,
		height:   h,
		next:     time.Now(),
	}
	_ = ebiten.RunGame(ig)
}
