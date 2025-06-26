package main

import (
	"image/color"
	"image/png"
	"os"

	ebdraw "github.com/arran4/gorillas/drawings/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	img := ebiten.NewImage(100, 100)
	ebdraw.DrawBASSun(img, 50, 50, 40, false, color.RGBA{255, 255, 0, 255})
	f, err := os.Create("sun.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
