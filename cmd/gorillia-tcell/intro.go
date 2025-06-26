package main

import (
	"strings"
	"time"

	"github.com/arran4/gorillas"
	"github.com/gdamore/tcell/v2"
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

func drawString(s tcell.Screen, x, y int, str string) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
}

func drawGorillaFrame(s tcell.Screen, x, y int, frame []string) {
	for i, line := range frame {
		drawString(s, x, y+i, line)
	}
}

// scrollText animates msg across the screen on row y when enabled.
// When disabled, it simply prints the message centred.
func scrollText(s tcell.Screen, y int, msg string, enabled bool) {
	w, _ := s.Size()
	if !enabled {
		drawString(s, (w-len(msg))/2, y, msg)
		s.Show()
		return
	}
	pad := strings.Repeat(" ", w)
	text := pad + msg + pad
	for i := 0; i <= len(msg)+w; i++ {
		drawString(s, 0, y, text[i:i+w])
		s.Show()
		time.Sleep(50 * time.Millisecond)
	}
}

func showIntroMovie(s tcell.Screen, useSound, sliding bool) {
	w, h := s.Size()
	lines := []string{
		"QBasic GORILLAS",
		"",
		"Starring two gorillas",
	}
	s.Clear()
	if sliding {
		for i, line := range lines {
			for j := 1; j <= len(line); j++ {
				drawString(s, (w-len(line))/2, h/2-1+i, line[:j])
				s.Show()
				time.Sleep(30 * time.Millisecond)
			}
		}
	} else {
		for i, line := range lines {
			drawString(s, (w-len(line))/2, h/2-1+i, line)
		}
		s.Show()
	}
	if useSound {
		gorillas.PlayIntroMusic()
	}
	time.Sleep(1500 * time.Millisecond)
	for i := 0; i < 4; i++ {
		drawGorillaFrame(s, w/2-10, h/2+2, gorillaFrames[i%len(gorillaFrames)])
		drawGorillaFrame(s, w/2+2, h/2+2, gorillaFrames[(i+1)%len(gorillaFrames)])
		s.Show()
		time.Sleep(300 * time.Millisecond)
	}
	time.Sleep(700 * time.Millisecond)
	scrollText(s, h-1, "Get ready to throw bananas!", sliding)
	time.Sleep(500 * time.Millisecond)
}

func introScreen(s tcell.Screen, useSound, sliding bool) bool {
	w, h := s.Size()
	cx := w/2 - 10
	cy := h/2 - 2
	for i := 0; i < 4; i++ {
		s.Clear()
		drawGorillaFrame(s, cx, cy, gorillaFrames[i%len(gorillaFrames)])
		drawGorillaFrame(s, cx+12, cy, gorillaFrames[(i+1)%len(gorillaFrames)])
		drawString(s, w/2-4, cy-2, "GORILLAS")
		s.Show()
		time.Sleep(300 * time.Millisecond)
	}
	for {
		s.Clear()
		drawGorillaFrame(s, cx, cy, gorillaFrames[0])
		drawGorillaFrame(s, cx+12, cy, gorillaFrames[0])
		drawString(s, w/2-4, cy-2, "GORILLAS")
		drawString(s, w/2-9, cy+3, "V - View Intro")
		drawString(s, w/2-9, cy+4, "P - Play Game")
		drawString(s, w/2-9, cy+5, "Q - Quit")
		s.Show()
		ev := s.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			switch key.Rune() {
			case 'q', 'Q':
				return false
			case 'p', 'P':
				return true
			case 'v', 'V':
				showIntroMovie(s, useSound, sliding)
			}
		}
	}
}
