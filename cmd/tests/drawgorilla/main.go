package main

import (
	"image/color"
	"image/png"
	"os"

	ebdraw "github.com/arran4/gorillas/drawings/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	scale := 2.0
	img := ebiten.NewImage(int(30*scale), int(30*scale))
	ebdraw.DrawBASGorilla(img, 15*scale, scale, scale, ebdraw.ArmsDown, color.RGBA{150, 75, 0, 255})
	f, err := os.Create("gorilla.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
