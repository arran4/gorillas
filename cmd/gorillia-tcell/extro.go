package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var extroPhrases = []string{
	"May the Schwarz be with you!",
	"Live long and prosper.",
	"Goodbye!",
	"So long!",
	"Adios!",
}

// showExtro displays a farewell message and waits briefly for user input.
func showExtro(s tcell.Screen) {
	w, h := s.Size()
	msg := "Thank you for playing Gorillas!"
	phrase := extroPhrases[rand.Intn(len(extroPhrases))]
	drawString(s, (w-len(msg))/2, h/2-1, msg)
	drawString(s, (w-len(phrase))/2, h/2+1, phrase)
	s.Show()
	end := time.Now().Add(4 * time.Second)
	for time.Now().Before(end) {
		if s.HasPendingEvent() {
			if _, ok := s.PollEvent().(*tcell.EventKey); ok {
				return
			}
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}
}
