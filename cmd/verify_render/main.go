package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"os/exec"

	"github.com/arran4/gorillas"
)

func main() {
	// Create a sample game state
	g := gorillas.NewGame(800, 600, 10)
	g.Gorillas[0] = gorillas.Gorilla{X: 100, Y: 500}
	g.Gorillas[1] = gorillas.Gorilla{X: 700, Y: 500}
	g.Banana.Active = true
	g.Banana.X = 400
	g.Banana.Y = 300

	// Create a building with a specific color
	g.Buildings[0].Color = color.RGBA{255, 0, 0, 255}

	// Serialize
	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		panic(err)
	}
	stateFile := "verify_state.json"
	if err := os.WriteFile(stateFile, b, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Created %s\n", stateFile)

	// Run the headless renderer
	// We use "go run ./cmd/gorillia-ebiten" to run the main package
	cmd := exec.Command("go", "run", "./cmd/gorillia-ebiten", "-render-state", stateFile, "-output-image", "verify_output.png")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running renderer...")
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("renderer failed: %v", err))
	}

	// Check output
	info, err := os.Stat("verify_output.png")
	if err != nil {
		panic(fmt.Errorf("output file not found: %v", err))
	}
	if info.Size() == 0 {
		panic("output file is empty")
	}
	fmt.Printf("Success! Generated verify_output.png (%d bytes)\n", info.Size())
}
