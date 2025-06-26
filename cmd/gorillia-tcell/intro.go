//go:build !test

package main

import (
	"fmt"
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
	SparklePause(s, 0)
}

func showInstructions(s tcell.Screen, sliding bool) {
	lines, err := gorillas.LoadInstructions()
	if err != nil {
		lines = []string{"Instructions unavailable"}
	}
	w, h := s.Size()
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	y := h/2 - len(lines)/2
	if sliding {
		for i := 1; i <= maxLen; i++ {
			for j, line := range lines {
				n := i
				if n > len(line) {
					n = len(line)
				}
				drawString(s, (w-len(line))/2, y+j, line[:n])
			}
			s.Show()
			time.Sleep(30 * time.Millisecond)
		}
	} else {
		for j, line := range lines {
			drawString(s, (w-len(line))/2, y+j, line)
		}
		s.Show()
	}
	SparklePause(s, 0)
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
		drawString(s, w/2-9, cy+3, "V/X - View Intro")
		drawString(s, w/2-9, cy+4, "I - Instructions")
		drawString(s, w/2-9, cy+5, "P/Start - Play Game")
		drawString(s, w/2-9, cy+6, "R - Replays")
		drawString(s, w/2-9, cy+7, "Q/B - Quit")
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
			case 'i', 'I':
				showInstructions(s, sliding)
			}
		}
	}
}

// SparklePause draws twinkling '*' borders for the given duration.  If
// duration is zero it waits until a key is pressed.
func SparklePause(s tcell.Screen, dur time.Duration) {
	for s.HasPendingEvent() { // clear pending keys
		s.PollEvent()
	}
	w, h := s.Size()
	start := time.Now()
	phase := 0
	pattern := []rune("*    ")
	for {
		for x := 0; x < w; x++ {
			ch1 := pattern[(phase+x)%5]
			ch2 := pattern[(4-phase+x)%5]
			s.SetContent(x, 0, ch1, nil, tcell.StyleDefault)
			s.SetContent(x, h-1, ch2, nil, tcell.StyleDefault)
		}
		for y := 1; y < h-1; y++ {
			ch := ' '
			if (phase+y)%5 == 0 {
				ch = '*'
			}
			s.SetContent(w-1, y, ch, nil, tcell.StyleDefault)
			s.SetContent(0, h-1-y, ch, nil, tcell.StyleDefault)
		}
		s.Show()
		time.Sleep(50 * time.Millisecond)
		phase = (phase + 1) % 5
		if dur > 0 && time.Since(start) > dur {
			return
		}
		if s.HasPendingEvent() {
			if _, ok := s.PollEvent().(*tcell.EventKey); ok {
				return
			}
		}
	}
}

// showStats prints the stats on the screen and waits for a key press.
func showStats(s tcell.Screen, stats string) {
	lines := strings.Split(stats, "\n")
	s.Clear()
	w, h := s.Size()
	y := h/2 - len(lines)/2
	for i, line := range lines {
		drawString(s, (w-len(line))/2, y+i, line)
	}
	msg := "Press any key to continue"
	drawString(s, (w-len(msg))/2, y+len(lines)+1, msg)
	s.Show()
	SparklePause(s, 0)
}

func showLeague(s tcell.Screen, l *gorillas.League) {
	if l == nil {
		return
	}
	rows := []string{"Player           Rounds Wins Accuracy"}
	for _, st := range l.Standings() {
		rows = append(rows, fmt.Sprintf("%-15s %6d %4d %8.1f", st.Name, st.Rounds, st.Wins, st.Accuracy))
	}
	s.Clear()
	w, h := s.Size()
	y := h/2 - len(rows)/2
	for i, line := range rows {
		drawString(s, (w-len(line))/2, y+i, line)
	}
	msg := "Press any key to continue"
	drawString(s, (w-len(msg))/2, y+len(rows)+1, msg)
	s.Show()
	SparklePause(s, 0)
}

// showGameAborted displays an aborted message and waits for a key press.
func showGameAborted(s tcell.Screen) {
	s.Clear()
	w, h := s.Size()
	msg := "Game aborted"
	drawString(s, (w-len(msg))/2, h/2, msg)
	s.Show()
	SparklePause(s, 0)
}
