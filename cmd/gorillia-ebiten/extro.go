//go:build !test

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var extroPhrases = []string{
	"May the Schwarz be with you!",
	"Live long and prosper.",
	"Goodbye!",
	"So long!",
	"Adios!",
}

// showExtro prints a farewell message and waits briefly for user input.
func showExtro() {
	fmt.Println("Thank you for playing Gorillas!")
	fmt.Println(extroPhrases[rand.Intn(len(extroPhrases))])

	reader := bufio.NewReader(os.Stdin)
	done := make(chan struct{})
	go func() {
		reader.ReadString('\n')
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(4 * time.Second):
	}
}
